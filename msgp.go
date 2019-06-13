package ipc

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"unsafe"
)

// Msgp message packet
type Msgp struct {
	Mtype uint
	Mtext []byte
}

var (
	mtypeSize = int(unsafe.Sizeof(uint(0)))
	byteType  = reflect.TypeOf(byte(0))
)

func (m *Msgp) marshal() (ptr unsafe.Pointer, textSize int) {
	var w bytes.Buffer
	switch mtypeSize {
	case 4:
		binary.Write(&w, binary.BigEndian, uint32(m.Mtype))
	case 8:
		binary.Write(&w, binary.BigEndian, uint64(m.Mtype))
	}
	binary.Write(&w, binary.BigEndian, m.Mtext)
	data := w.Bytes()
	count := len(data)
	t := reflect.ArrayOf(count, byteType)
	v := reflect.New(t).Elem()
	for i := 0; i < count; i++ {
		v.Index(i).Set(reflect.ValueOf(data[i]))
	}
	return unsafe.Pointer(v.Addr().Pointer()), count - mtypeSize
}

func (m *Msgp) unmarshal(textSize int, ptr unsafe.Pointer) error {
	m.Mtext = make([]byte, 0, textSize)
	size := textSize + mtypeSize
	t := reflect.ArrayOf(size, byteType)
	v := reflect.NewAt(t, ptr).Elem()
	v.Interface() // implementation data
	mtypeBytes := make([]byte, mtypeSize)
	for i := 0; i < size; i++ {
		e := byte(v.Index(i).Uint())
		if i < mtypeSize {
			mtypeBytes[i] = e
		} else {
			m.Mtext = append(m.Mtext, e)
		}
	}
	switch mtypeSize {
	case 4:
		var u uint32
		err := binary.Read(bytes.NewReader(mtypeBytes), binary.BigEndian, &u)
		if err != nil {
			return err
		}
		m.Mtype = uint(u)
	case 8:
		var u uint64
		err := binary.Read(bytes.NewReader(mtypeBytes), binary.BigEndian, &u)
		if err != nil {
			return err
		}
		m.Mtype = uint(u)
	}
	return nil
}
