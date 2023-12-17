package internals

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type Storer interface {
	New(string) error
	Store(Object) error
	ReadObj(objID []byte) ([]byte, error)
	RestoreFile(name string, data []byte) error
	LookupRef(ref string) (*Ref, error)
	UpdateRef(ref Refer) error
	GetTag(name string) ([]byte, error)
	GetBranch(name string) (*Branch, error)
	IsBranch(name string) bool
}

type Store struct {
	path string
}

func (db *Store) New(path string) error {
	db.path = path
	fmt.Println("New Store at ", db.path)
	err := db.createRefPath()
	if err != nil {
		return err
	}
	headPath := OsPathJoin(path, "HEAD")
	_, err = os.Stat(headPath)
	if os.IsNotExist(err) {
		err := os.WriteFile(
			headPath,
			[]byte(""),
			YitDefaultPermissions,
		)
		if err != nil {
			return err
		}

	}
	return os.MkdirAll(
		OsPathJoin(db.path, "refs", "heads"),
		YitDefaultDirPermissions,
	)
}

func (db *Store) RestoreFile(name string, data []byte) error {
	fp := strings.Join(
		[]string{db.path, name},
		string(os.PathSeparator),
	)
	return os.WriteFile(fp, data, YitDefaultPermissions)
}

func (db *Store) GetPath() string {
	return db.path
}

func (db *Store) Store(obj Object) error {
	var ed bytes.Buffer
	data := []byte(fmt.Sprintf("%T %d\x00", obj, len(obj.ToString())))
	data = append(data, []byte(obj.ToString())...)
	// Encode the Data
	if err := binary.Write(&ed, binary.BigEndian, data); err != nil {
		return err
	}
	// Create SHA-1 hash of the Encoded data
	objID := sha1.Sum(ed.Bytes())
	obj.SetObjID(objID[:])

	return db.writeObj(obj.GetObjID(), ed.Bytes())
}

func (s *Store) writeObj(objID []byte, data []byte) error {
	ObjIDHex := fmt.Sprintf("%x", objID)
	fmt.Println(ObjIDHex)
	objDir := OsPathJoin(s.path, "objects", ObjIDHex[:2])
	objPath := OsPathJoin(objDir, ObjIDHex[2:])

	if err := os.MkdirAll(objDir, YitDefaultDirPermissions); err != nil {
		return err
	}
	// create temp dir
	tmpDir := os.TempDir()
	tf, err := os.CreateTemp(tmpDir, genTmpObjName(ObjIDHex))
	if err != nil {
		return err
	}
	// compress the data
	var c bytes.Buffer
	zw := zlib.NewWriter(&c)
	if _, err := zw.Write(data); err != nil {
		return err
	}
	zw.Close()
	if _, err = tf.Write(c.Bytes()); err != nil {
		return err
	}
	tf.Close()
	if err := os.Rename(tf.Name(), objPath); err != nil {
		if err != nil &&
			strings.Contains(err.Error(), "invalid cross-device link") {
			return moveCrossDevice(tf.Name(), objPath)
		}
		return err
	}
	return nil
}

// https://gist.github.com/var23rav/23ae5d0d4d830aff886c3c970b8f6c6b
func moveCrossDevice(source, destination string) error {
	src, err := os.Open(source)
	if err != nil {
		return errors.Wrap(err, "Open(source)")
	}
	dst, err := os.Create(destination)
	if err != nil {
		src.Close()
		return errors.Wrap(err, "Create(destination)")
	}
	_, err = io.Copy(dst, src)
	src.Close()
	dst.Close()
	if err != nil {
		return errors.Wrap(err, "Copy")
	}
	fi, err := os.Stat(source)
	if err != nil {
		os.Remove(destination)
		return errors.Wrap(err, "Stat")
	}
	err = os.Chmod(destination, fi.Mode())
	if err != nil {
		os.Remove(destination)
		return errors.Wrap(err, "Stat")
	}
	os.Remove(source)
	return nil
}

func genTmpObjName(ObjIDHex string) string {
	return fmt.Sprintf("tmp_obj_%x", ObjIDHex[0:3])
}

func (s *Store) GetTag(name string) ([]byte, error) {
	tagPth := OsPathJoin("refs", "tags", name)
	ref, err := s.LookupRef(tagPth)
	if err != nil {
		return nil, err
	}
	return ref.GetValue(), nil
}

func OsPathJoin(entries ...string) string {
	return strings.Join(entries, string(os.PathSeparator))
}

type Branch struct {
	OID  []byte
	Name string
}

func (s *Store) GetBranch(name string) (*Branch, error) {
	brnchPth := OsPathJoin("refs", "heads", name)
	ref, err := s.LookupRef(brnchPth)
	if err != nil {
		return nil, err
	}
	return &Branch{
		OID:  ref.GetValue(),
		Name: brnchPth,
	}, nil
}

func (s *Store) CreateTag(name string, OID []byte) error {
	rootDir, err := os.Getwd()
	if err != nil {
		return err
	}
	refPth := strings.Join(
		[]string{rootDir, YitMetadataDir, "refs", "tags"},
		string(os.PathSeparator),
	)
	if err := os.MkdirAll(refPth, YitDefaultDirPermissions); err != nil {
		return err
	}
	tagPth := strings.Join([]string{"tags", name}, string(os.PathSeparator))

	ref := NewRef(tagPth, OID)
	return s.UpdateRef(ref)
}

func (s *Store) IsBranch(name string) bool {
	brnchPth := strings.Join(
		[]string{"heads", name},
		string(os.PathSeparator),
	)
	_, err := s.LookupRef(brnchPth)
	if err != nil {
		return false
	}
	return true
}

func (s *Store) CreateBranch(name string, OID []byte) error {
	headsPth := OsPathJoin(s.path, "refs", "heads")
	if err := os.MkdirAll(headsPth, YitDefaultDirPermissions); err != nil {
		return err
	}
	brnchPth := OsPathJoin("refs", "heads", name)
	ref := NewRef(brnchPth, OID)
	fmt.Printf(
		"creating branch for ref: %s, %x\n",
		ref.GetName(),
		ref.GetValue(),
	)
	return s.UpdateRef(ref)
}

func (s *Store) createRefPath() error {
	rootDir, err := os.Getwd()
	if err != nil {
		return err
	}
	refDir := strings.Join(
		[]string{rootDir, YitMetadataDir, "refs"},
		string(os.PathSeparator),
	)
	if err := os.MkdirAll(refDir, YitDefaultDirPermissions); err != nil {
		return err
	}
	return nil
}

func (s *Store) UpdateRef(ref Refer) error {
	fmt.Println("CALLED===========================")
	fmt.Printf(
		"updating ref name %s to match ref val: %x ",
		ref.GetName(),
		ref.GetValue(),
	)
	nref, err := s.LookupRef(ref.GetName())
	fmt.Println("found some nref", nref, err)

	fp := ""
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("ERROR:", nref, err)
			return err
		}
	}
	fp = OsPathJoin(s.path, ref.GetName())
	if nref != nil {
		fp = OsPathJoin(s.path, nref.GetName())
	}
	f, err := os.OpenFile(
		fp,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0644,
	)
	fmt.Println("f error", err)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write the OID to the HEAD file.
	val := ref.GetValue()
	if _, err := f.Write(val); err != nil {
		return err
	}

	return nil
}

var stck = 0

// return the direct ref only and resolve symRefs
func (s *Store) LookupRef(refName string) (*Ref, error) {
	fmt.Println("looking up...")
	if refName == OsPathJoin("refs", "heads", "main") {
		fmt.Println("is main !")
		if _, err := os.Stat(OsPathJoin(s.path, refName)); os.IsNotExist(err) {
			return NewRef(refName, []byte("")), nil
		}
	}
	fp := OsPathJoin(s.path, refName)

	if _, err := os.Stat(fp); err != nil {
		fmt.Println("refName:", refName)
		return nil, err
	}

	data, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}
	fmt.Printf(
		"looking for %x , str: %s ,with len %d %v,\n %s @%s\n",
		data,
		string(data),
		len(data),
		err,
		refName,
		fp,
	)
	strData := string(data)
	if strings.HasPrefix(strData, "ref: ") {
		str := strings.Split(strData, ":")[1]
		rfn := strings.TrimSpace(str)
		fmt.Println("rfn:", rfn)
		stck++
		if stck > 10 {
			log.Fatal("infinite loop...")
		}
		return s.LookupRef(rfn)
	}
	ref := NewRef(refName, data)
	return ref, nil
}

func (s *Store) ReadObj(objID []byte) ([]byte, error) {
	ObjIDHex := fmt.Sprintf("%x", objID)
	objDir := OsPathJoin(s.path, "objects", ObjIDHex[:2])
	objPath := OsPathJoin(objDir, ObjIDHex[2:])
	f, err := os.Open(objPath)
	if err != nil {
		return nil, err
	}
	zr, err := zlib.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dcmpData, err := io.ReadAll(zr)
	return dcmpData, nil
}
