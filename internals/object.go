package internals

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
)

type Object interface {
	ToString() string
	GetObjID() []byte
	SetObjID([]byte) error
}

func GenObjID(obj Object) ([]byte, error) {
	var ed bytes.Buffer
	data := []byte(fmt.Sprintf("%T %d\x00", obj, len(obj.ToString())))
	data = append(data, []byte(obj.ToString())...)
	// Encode the Data
	if err := binary.Write(&ed, binary.BigEndian, data); err != nil {
		return nil, err
	}
	// Create SHA-1 hash of the Encoded data
	objID := sha1.Sum(ed.Bytes())
	return objID[:], nil
}
