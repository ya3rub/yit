package internals

import (
	"fmt"
)

type Entry struct {
	objID []byte
	// TODO:
	// mode  os.FileMode
	name string
}

func (e *Entry) New(objID []byte, name string) error {
	if len(name) == 0 || len(objID) == 0 {
		return fmt.Errorf("name , objID cannot be empty")
	}
	e.objID = objID
	e.name = name
	return nil
}

func (e *Entry) GetName() string {
	return e.name
}

func (e *Entry) SetName(name string) error {
	e.name = name
	return nil
}

func (e *Entry) ToString() string {
	return e.name
}

func (e *Entry) GetObjID() []byte {
	return e.objID
}

func (e *Entry) SetObjID(objID []byte) error {
	e.objID = objID
	return nil
}
