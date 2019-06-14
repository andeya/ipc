package ipc

import (
	"reflect"
	"unsafe"
)

// Msgp message packet
type Msgp struct {
	Mtype uint
	Mtext []byte
}

var (
	mtypeSize = unsafe.Sizeof(uint(0))
	byteType  = reflect.TypeOf(byte(0))
)

type msgpBuf struct {
	Mtype uint
	Mtext [0]byte
}

func (m *Msgp) marshal() (ptr unsafe.Pointer, textSize int) {
	count := len(m.Mtext)
	buf := make([]byte, count+int(mtypeSize))
	ptr = unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&buf)))
	mbuf := (*msgpBuf)(ptr)
	mbuf.Mtype = m.Mtype

	t := reflect.ArrayOf(count, byteType)
	mtext := reflect.NewAt(t, unsafe.Pointer(uintptr(ptr)+mtypeSize)).Elem()

	mtextPtr := uintptr(unsafe.Pointer(m)) + mtypeSize
	mtextData := reflect.NewAt(t, unsafe.Pointer(*(*uintptr)(unsafe.Pointer(mtextPtr)))).Elem()

	mtext.Set(mtextData)

	return ptr, count
}

func (m *Msgp) unmarshal(textSize int, ptr unsafe.Pointer) error {
	m.Mtext = make([]byte, textSize)
	copy(m.Mtext, *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(ptr) + mtypeSize,
		Len:  textSize,
		Cap:  textSize,
	})))
	m.Mtype = (*(*msgpBuf)(ptr)).Mtype
	return nil
}
