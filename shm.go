package ipc

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"syscall"
	"unsafe"
)

// via:
//  https://code.woboq.org/userspace/glibc/sysdeps/unix/sysv/linux/bits/shm.h.html
//  https://blog.csdn.net/guoping16/article/details/6584058

const (
	/* Flags for `shmat'.  */
	SHM_RDONLY = 010000  // attach read-only else read-write
	SHM_RND    = 020000  // round attach address to SHMLBA
	SHM_REMAP  = 040000  // take-over region on attach
	SHM_EXEC   = 0100000 // execution access

	/* Commands for `shmctl'.  */
	SHM_LOCK   = 11 // lock segment (root only)
	SHM_UNLOCK = 12 // unlock segment (root only)

	/* shm ctl commands */
	SHM_STAT     = 13
	SHM_INFO     = 14
	SHM_STAT_ANY = 15

	/* shm_mode upper byte flags */
	SHM_DEST      = 01000  // segment will be destroyed on last detach
	SHM_LOCKED    = 02000  // segment will not be swapped
	SHM_HUGETLB   = 04000  // segment is mapped via hugetlb
	SHM_NORESERVE = 010000 // don't check for reservations
)

var shmInfo = struct {
	shmidSizes    map[int]uint64            // <shmid, size>
	shmaddrSizes  map[unsafe.Pointer]uint64 // <shmaddr, size>
	shmaddrShmids map[unsafe.Pointer]int    // <shmaddr, shmid>
	sync.RWMutex
}{
	shmidSizes:    make(map[int]uint64, 64),
	shmaddrSizes:  make(map[unsafe.Pointer]uint64, 64),
	shmaddrShmids: make(map[unsafe.Pointer]int, 64),
}

// FtokAndShmget returns a probably-unique key, and get a shared memory identifier,
// or create a shared memory object and return a shared memory identifier.
func FtokAndShmget(path string, id, size uint64, shmflg int) (key uint64, shmid int, err error) {
	key, err = Ftok(path, id)
	if err == nil {
		shmid, err = Shmget(key, size, shmflg)
	}
	return
}

// Shmget get a shared memory identifier,
// or create a shared memory object and return a shared memory identifier.
func Shmget(key, size uint64, shmflg int) (shmid int, err error) {
	_shmid, _, errno := syscall.Syscall(syscall.SYS_SHMGET, uintptr(key), uintptr(size), uintptr(shmflg))
	if errno != 0 {
		return 0, errno
	}
	shmid = int(_shmid)
	shmInfo.Lock()
	shmInfo.shmidSizes[shmid] = size
	shmInfo.Unlock()
	return shmid, nil
}

// Shmat connect the shared memory identifier to the shared memory of shmid.
// After the connection is successful, map the shared memory area object to the
// address space of the calling process, and then access it like the local space.
func Shmat(shmid int, shmflg int) (shmaddr unsafe.Pointer, err error) {
	_shmaddr, _, errno := syscall.Syscall(syscall.SYS_SHMAT, uintptr(shmid), 0, uintptr(shmflg))
	if errno != 0 {
		return nil, errno
	}
	shmaddr = unsafe.Pointer(_shmaddr)

	shmInfo.Lock()
	if size, ok := shmInfo.shmidSizes[shmid]; ok {
		shmInfo.shmaddrSizes[shmaddr] = size
		shmInfo.shmaddrShmids[shmaddr] = shmid
	}
	shmInfo.Unlock()

	return shmaddr, nil
}

// Shmdt contrary to the shmat function, it is used to disconnect the address with the
// shared memory attachment point, prohibiting the process from accessing the slice shared memory.
func Shmdt(shmaddr unsafe.Pointer) error {
	_, _, errno := syscall.Syscall(syscall.SYS_SHMDT, uintptr(shmaddr), 0, 0)
	if errno != 0 {
		return errno
	}
	shmInfo.Lock()
	delete(shmInfo.shmaddrSizes, shmaddr)
	delete(shmInfo.shmaddrShmids, shmaddr)
	shmInfo.Unlock()
	return nil
}

// Shmctl control shared memory
// NOTE:
//  cmd: IPC_STAT, IPC_SET, IPC_RMID
//  Currently only IPC_RMID is implemented!
func Shmctl(shmid, cmd int) error {
	var buf uintptr = 0
	_, _, errno := syscall.Syscall(syscall.SYS_SHMCTL, uintptr(shmid), uintptr(cmd), buf)
	if errno != 0 || cmd != IPC_RMID {
		return errno
	}
	shmInfo.Lock()
	delete(shmInfo.shmidSizes, shmid)
	for addr, id := range shmInfo.shmaddrShmids {
		if id == shmid {
			delete(shmInfo.shmaddrShmids, addr)
			delete(shmInfo.shmaddrSizes, addr)
		}
	}
	shmInfo.Unlock()
	return nil
}

// Shmwrite write data to the shared memory.
func Shmwrite(shmaddr unsafe.Pointer, data []byte) error {
	size := 4 + len(data)
	shmInfo.RLock()
	maxSize := shmInfo.shmaddrSizes[shmaddr]
	shmInfo.RUnlock()
	if uint64(size) > maxSize {
		return fmt.Errorf("data is too large, (4 + %d) > %d", len(data), maxSize)
	}
	buf := make([]byte, size)
	binary.BigEndian.PutUint32(buf, uint32(size))
	copy(buf[4:], data)
	ptr := unsafe.Pointer(*(*uintptr)(unsafe.Pointer(&buf)))
	t := reflect.ArrayOf(size, byteType)
	reflect.NewAt(t, shmaddr).Elem().Set(reflect.NewAt(t, ptr).Elem())
	return nil
}

// Shmread read data from the shared memory.
func Shmread(shmaddr unsafe.Pointer) []byte {
	bytesPtr := (*[4]byte)(shmaddr)
	if bytesPtr == nil {
		return nil
	}
	sizeBytes := *bytesPtr
	size := int(binary.BigEndian.Uint32(sizeBytes[:])) - 4
	if size <= 0 {
		return []byte{}
	}
	buf := make([]byte, size)
	copy(buf, *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(shmaddr) + 4,
		Len:  size,
		Cap:  size,
	})))
	runtime.KeepAlive(shmaddr)
	return buf
}
