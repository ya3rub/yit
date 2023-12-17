package internals

// import (
//
//	"strings"
//
// )
type Refer interface {
	// resolve() ([]byte, error)
	GetValue() []byte
	GetName() string
}

type Ref struct {
	name  string
	value []byte
}

type symRef struct {
	name  string
	value string
}

func NewRef(name string, val []byte) *Ref {
	return &Ref{
		name:  name,
		value: val,
	}
}

func (r *Ref) GetName() string {
	return r.name
}

func (r *Ref) GetValue() []byte {
	return r.value
}

//	func (r *ref) resolve() ([]byte, error) {
//		if strings.HasPrefix(string(r.value), "ref: ") {
//			// s:= strings.Split(string(r.value), ":")[1]
//			// strings.TrimSpace(s)
//			// return
//		}
//		return r.value, nil
//	}
func NewSymRef(name string, value string) *symRef {
	return &symRef{
		name:  name,
		value: "ref: " + value,
	}
}

func (r *symRef) GetValue() []byte {
	return []byte(r.value)
}

func (r *symRef) GetName() string {
	return r.name
}

// func (r *symRef) resolve() ([]byte, error) {
// 	if strings.HasPrefix(string(r.value), "ref: ") {
// 		str := strings.Split(string(r.value), ":")[1]
// 		refName := strings.TrimSpace(str)
// 		ref, err := r.store.GetRef(refName)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return ref.resolve()
// 	}
// 	return r.resolve()
// }
//
// func (r *symRef) getName() string {
// 	return r.name
// }
//
// func isSymRef(data string) bool {
// }
