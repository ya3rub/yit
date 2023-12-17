package internals

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Tree struct {
	objID []byte
	data  []Entry
}

func (t *Tree) New(entries []Entry) error {
	t.data = entries
	return nil
}

func (t *Tree) ToString() string {
	dc := make([]Entry, len(t.data))
	copy(dc, t.data)
	sort.SliceStable(dc, func(i, j int) bool {
		return dc[i].GetName() < dc[j].GetName()
	})
	// turn into byte slice

	// <mode> <name>\x00<objID><mode> <name>\x00<objID><mode> <name>\x00<objID>....
	var entries []byte
	// entries = append(entries, []byte(fmt.Sprintf("tree %d\x00", len(dc)))...)
	for _, e := range dc {
		// ex: 100644 go.mod\x00321bcf6efe93685a67b4127d8c00c3bc26145843
		entries = append(
			entries,
			[]byte(
				fmt.Sprintf(
					// 10 for the file type
					"10%04o %s\x00",
					YitDefaultPermissions,
					e.GetName(),
				),
			)...,
		)
		entries = append(entries, e.GetObjID()...)
	}
	return string(entries)
}

func (t *Tree) GetObjID() []byte {
	return t.objID
}

func (t *Tree) SetObjID(objID []byte) error {
	t.objID = objID
	return nil
}

func ParseTree(s Storer, treeObjID []byte) (*Tree, error) {
	treeData, err := s.ReadObj(treeObjID)
	if err != nil {
		return nil, err
	}
	i := bytes.Index(treeData, []byte(" "))
	objType := treeData[0:i]
	if string(objType) != "*internals.Tree" {
		fmt.Println(objType)
		return nil, errors.New("not a Tree type")
	}
	dataStartIdx := bytes.Index(treeData, []byte("\x00"))
	data := treeData[dataStartIdx+1:]
	var entries []Entry
	for p := 0; p < len(data); {
		d := bytes.Index(data[p:], []byte("\x00"))
		entryMetaData := data[p : p+d]
		objID := data[p+d+1 : p+d+21]
		entries = append(entries, Entry{
			name:  strings.Split(string(entryMetaData), " ")[1],
			objID: objID,
		})
		// objHex := fmt.Sprintf("%x", objID)
		// fmt.Println(
		// 	"ent is :",
		// 	strings.Split(string(entryMetaData), " ")[1],
		// 	"obj",
		// 	objHex,
		// )
		p = p + d + 21
	}
	// fmt.Printf("entries is :%+v", entries)
	return &Tree{
		data:  entries,
		objID: treeObjID,
	}, nil
}

func RestoreTree(s Storer, dst string, treeObjID []byte) error {
	tree, err := ParseTree(s, treeObjID)
	// fmt.Printf("tree: %v\n", tree)
	if err != nil {
		return err
	}
	for _, ent := range tree.data {
		data, err := s.ReadObj(ent.GetObjID())
		if err != nil {
			return err
		}
		fp := strings.Join(
			[]string{dst, ent.name},
			string(os.PathSeparator),
		)
		objType := strings.Split(string(data), " ")
		switch objType[0] {
		case "*internals.Tree":
			fmt.Println("Dir: ", fp)
			if err := os.MkdirAll(fp, YitDefaultDirPermissions); err != nil {
				return err
			}
			if err := RestoreTree(s, fp, ent.GetObjID()); err != nil {
				return err
			}
			// fmt.Printf("dir: %s\n", ent.name)
		case "*internals.Blob":
			fmt.Printf("file: %s\n", fp)
			err = os.WriteFile(
				fp,
				bytes.Split(data, []byte("\x00"))[1],
				YitDefaultPermissions,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
