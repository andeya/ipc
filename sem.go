package ipc

import (
	"syscall"
	"unsafe"
)

// via
//  https://www.cnblogs.com/cangqinglang/p/13754849.html
//  https://code.woboq.org/userspace/glibc/sysdeps/unix/sysv/linux/bits/sem.h.html
//  https://blog.csdn.net/guoping16/article/details/6584043
//
//  https://blog.csdn.net/liufuchun111/article/details/87873130

const (
	/* ipcs ctl cmds */
	SEM_STAT     = 18
	SEM_INFO     = 19
	SEM_STAT_ANY = 20
	SEM_UNDO     = 0x1000 // undo the operation on exit

	/* Commands for `semctl'.  */
	GETPID  = 11 /* get sempid */
	GETVAL  = 12 /* get semval */
	GETALL  = 13 /* get all semval's */
	GETNCNT = 14 /* get semncnt */
	GETZCNT = 15 /* get semzcnt */
	SETVAL  = 16 /* set semval */
	SETALL  = 17 /* set all semval's */
)

// Semget create or open a semaphore set.
func Semget(key uint64, nsems int, semflg int) (semid int, err error) {
	_semid, _, errno := syscall.Syscall(syscall.SYS_SEMGET, uintptr(key), uintptr(nsems), uintptr(semflg))
	if errno != 0 {
		return 0, errno
	}
	return int(_semid), nil
}

type SemOp struct {
	SemNum uint16
	SemOp  int16
	SemFlg int16 // IPC_NOWAIT, SEM_UNDO
}

// Semop P-operation or V-operation on one or more semaphores.
func Semop(semid int, sops []SemOp) (err error) {
	_, _, errno := syscall.Syscall(syscall.SYS_SEMOP, uintptr(semid), uintptr(unsafe.Pointer(&sops[0])), uintptr(len(sops)))
	if errno != 0 {
		return errno
	}
	return nil
}

// Semctl get and set the properties of the message queue.
func Semctl(semid int, cmd int) error {
	var semnum int
	_, _, errno := syscall.Syscall(syscall.SYS_SEMCTL, uintptr(semid), uintptr(unsafe.Pointer(&semnum)), uintptr(cmd))
	if errno != 0 {
		return errno
	}
	return nil
}
