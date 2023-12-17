package internals

type Blob struct {
	objID []byte
	data  []byte
	name  string
}

func (b *Blob) New(data []byte, name string) error {
	b.data = data
	b.name = name
	return nil
}

func (b *Blob) ToString() string {
	return string(b.data)
}

func (b *Blob) GetName() string {
	return b.name
}

func (b *Blob) GetObjID() []byte {
	return b.objID
}

func (b *Blob) SetObjID(objID []byte) error {
	b.objID = objID
	return nil
}
