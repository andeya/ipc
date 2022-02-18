package ipc

import (
	"fmt"
	"syscall"
	"time"
)

// via:
//  https://blog.csdn.net/liufuchun111/article/details/87873130

var (
	hmsInit = []SemOp{{
		SemNum: 0,
		SemOp:  1,
		SemFlg: 0,
	}}
	hmsRl = []SemOp{{
		SemNum: 0,
		SemOp:  -1,
		SemFlg: SEM_UNDO,
	}, {
		SemNum: 0,
		SemOp:  2,
		SemFlg: SEM_UNDO,
	}}
	hmsRUl = []SemOp{{
		SemNum: 0,
		SemOp:  -1,
		SemFlg: SEM_UNDO,
	}}
	hmsWl = []SemOp{
		{
			SemNum: 0,
			SemOp:  -1,
			SemFlg: SEM_UNDO,
		},
		{
			SemNum: 0,
			SemOp:  0,
			SemFlg: SEM_UNDO,
		},
	}
	hmsWUl = []SemOp{{
		SemNum: 0,
		SemOp:  1,
		SemFlg: SEM_UNDO,
	}}
)

func NewSemLock(key uint64) (*SemLock, error) {
	semid, err := Semget(key, 1, IPC_CREAT|IPC_EXCL|1023)
	fmt.Println(semid, err)
	if err == nil {
		err = Semop(semid, hmsInit)
		if err != nil {
			Semctl(semid, IPC_RMID)
			return nil, err
		}
	} else if err == syscall.EEXIST {
		semid, err = Semget(key, 1, 0)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	return &SemLock{semid: semid}, nil
}

var _ IPLock = new(SemLock)

type SemLock struct {
	semid int
}

func (s *SemLock) RLock() {
	for {
		switch Semop(s.semid, hmsRl) {
		case nil:
			return
		case syscall.EINTR:
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func (s *SemLock) RUnlock() {
	_ = Semop(s.semid, hmsRUl)
}

func (s *SemLock) Lock() {
	for {
		switch Semop(s.semid, hmsWl) {
		case nil:
			return
		case syscall.EINTR:
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func (s *SemLock) Unlock() {
	_ = Semop(s.semid, hmsWUl)
}

func (s *SemLock) Close() {
	Semctl(s.semid, IPC_RMID)
}
