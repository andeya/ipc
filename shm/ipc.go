package ipcshm

import (
	"fmt"
	"syscall"
	"unsafe"
)

// IPC through shared memory.
type IPC struct {
	key     uint64
	bufSize uint64
	shmaddr unsafe.Pointer
}

func (i IPC) Pointer() unsafe.Pointer {
	return i.shmaddr
}

func (i IPC) BufSize() uint64 {
	return i.bufSize
}

func (i IPC) Detach() {
	syscall.Syscall(syscall.SYS_SHMDT, uintptr(i.shmaddr), 0, 0)
}

func ShareByFtok(path string, id, bufSize uint64) (*IPC, error) {
	key, err := Ftok(path, id)
	if err != nil {
		return nil, err
	}
	return Share(key, bufSize)
}

// Ftok returns a probably-unique key that can be used by the Share.
// See ftok(3) and https://code.woboq.org/userspace/glibc/sysvipc/ftok.c.html
func Ftok(path string, id uint64) (uint64, error) {
	st := &syscall.Stat_t{}
	if err := syscall.Stat(path, st); err != nil {
		return 0, err
	}
	return uint64((st.Ino & 0xffff) | uint64((st.Dev&0xff)<<16) |
		((id & 0xff) << 24)), nil
}

func Share(key, bufSize uint64) (*IPC, error) {
	const IPC_CREAT = 01000 // Create key if key does not exist.
	const shmflg = IPC_CREAT | 0600
	shmid, _, errno := syscall.Syscall(syscall.SYS_SHMGET, uintptr(key), uintptr(bufSize), shmflg)
	if errno != 0 {
		return nil, errno
	}
	shmaddr, _, errno := syscall.Syscall(syscall.SYS_SHMAT, shmid, 0, 0)
	if errno != 0 {
		return nil, errno
	}
	return &IPC{
		key:     key,
		bufSize: bufSize,
		shmaddr: unsafe.Pointer(shmaddr),
	}, nil
}
