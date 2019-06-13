package ipc

import (
	"syscall"
	"unsafe"
)

// via:
//  https://code.woboq.org/userspace/glibc/sysdeps/unix/sysv/linux/bits/shm.h.html
//  https://blog.csdn.net/guoping16/article/details/6584058

const (
	/* Permission flag for shmget.  */
	SHM_R = 0400 // read
	SHM_W = 0200 // write

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

// Shmget get a shared memory identifier,
// or create a shared memory object and return a shared memory identifier.
func Shmget(key, size uint64, shmflg int) (shmid int, err error) {
	_shmid, _, errno := syscall.Syscall(syscall.SYS_SHMGET, uintptr(key), uintptr(size), uintptr(shmflg))
	if errno != 0 {
		return 0, errno
	}
	return int(_shmid), nil
}

// Shmat connect the shared memory identifier to the shared memory of shmid.
// After the connection is successful, map the shared memory area object to the
// address space of the calling process, and then access it like the local space.
func Shmat(shmid int, shmflg int) (shmaddr unsafe.Pointer, err error) {
	_shmaddr, _, errno := syscall.Syscall(syscall.SYS_SHMAT, uintptr(shmid), 0, uintptr(shmflg))
	if errno != 0 {
		return nil, errno
	}
	return unsafe.Pointer(_shmaddr), nil
}

// Shmdt contrary to the shmat function, it is used to disconnect the address with the
// shared memory attachment point, prohibiting the process from accessing the slice shared memory.
func Shmdt(shmaddr unsafe.Pointer) error {
	_, _, errno := syscall.Syscall(syscall.SYS_SHMDT, uintptr(shmaddr), 0, 0)
	if errno != 0 {
		return errno
	}
	return nil
}

// Shmctl control shared memory
// NOTE:
//  cmd: IPC_STAT, IPC_SET, IPC_RMID
func Shmctl(shmid, cmd int) error {
	var buf uintptr = 0
	_, _, errno := syscall.Syscall(syscall.SYS_SHMCTL, uintptr(shmid), uintptr(cmd), buf)
	if errno != 0 {
		return errno
	}
	return nil
}
