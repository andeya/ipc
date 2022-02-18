package ipc

import (
	"reflect"
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
func Semctl(semid int, cmd int) (uintptr, error) {
	var semnum int = 0
	var arg uintptr = 0
	r1, _, errno := syscall.Syscall6(syscall.SYS_SEMCTL, uintptr(semid), *(*uintptr)(unsafe.Pointer(&semnum)), uintptr(cmd), arg, 0, 0)
	if errno != 0 {
		return 0, errno
	}
	return r1, nil
}

type SemidDs struct {
	SemPerm  uintptr
	SemInfos []SemInfo
	SemOtime int64
	SemCtime int64
}

type SemInfo struct {
	SemVal  uint16 // 信号量的值
	SemPID  int    // 最近一次执行操作的进程的进程号
	SemNcnt uint16 // 等待信号值增长，即等待可利用资源出现的进程数
	SemZcnt uint16 // 等待信号值减少，即等待全部资源可被独占的进程数
}

// SemAllInfo get all the info.
func SemAllInfo(semid int) (*SemidDs, error) {
	type _SemidDs struct {
		SemPerm  uintptr
		SemBase  uintptr
		SemNsems uint16 // []SemInfo ptr
		SemOtime int64
		SemCtime int64
	}
	r, err := Semctl(semid, GETALL)
	if err != nil {
		return nil, err
	}
	buf := (*_SemidDs)(unsafe.Pointer(r))
	return (*SemidDs)(unsafe.Pointer(&reflect.SliceHeader{
		Data: buf.SemBase,
		Len:  int(buf.SemNsems),
		Cap:  int(buf.SemNsems),
	})), nil
}
