package ipc

import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"
)

// via:
//  https://blog.csdn.net/liufuchun111/article/details/87873130

var (
	hmsRl = []SemOp{{
		SemNum: 0,
		SemOp:  0,
		SemFlg: SEM_UNDO,
	}, {
		SemNum: 1,
		SemOp:  1,
		SemFlg: SEM_UNDO,
	}}
	hmsRUl = []SemOp{{
		SemNum: 1,
		SemOp:  -1,
		SemFlg: SEM_UNDO,
	}}
	hmsWl = []SemOp{
		{
			SemNum: 1,
			SemOp:  0,
			SemFlg: SEM_UNDO,
		},
		{
			SemNum: 0,
			SemOp:  0,
			SemFlg: SEM_UNDO,
		},
		{
			SemNum: 0,
			SemOp:  1,
			SemFlg: SEM_UNDO,
		},
	}
	hmsWUl = []SemOp{{
		SemNum: 0,
		SemOp:  -1,
		SemFlg: SEM_UNDO,
	}}
)

func NewSemLock(key uint64) (*SemLock, error) {
	semid, err := Semget(key, 2, IPC_CREAT|IPC_EXCL|1023)
	// fmt.Println(semid, err)
	if err == syscall.EEXIST {
		semid, err = Semget(key, 2, 0)
	}
	if err != nil {
		return nil, err
	}
	return &SemLock{semid: semid}, nil
}

var _ IPLock = new(SemLock)

type SemLock struct {
	semid int
	local sync.RWMutex
}

func (s *SemLock) RLock() {
	for {
		switch err := Semop(s.semid, hmsRl); err {
		case nil:
			return
		case syscall.EINTR:
		case syscall.EINVAL:
			fmt.Fprintf(os.Stderr, "*SemLock.Lock: error(%v), so downgrade to local lock\n", err)
			s.local.RLock()
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

func (s *SemLock) RUnlock() {
	err := Semop(s.semid, hmsRUl)
	if err == syscall.EINVAL {
		fmt.Fprintf(os.Stderr, "*SemLock.Unlock: error(%v), so downgrade to local lock\n", err)
		s.local.RUnlock()
	}
}

func (s *SemLock) Lock() {
	for {
		switch err := Semop(s.semid, hmsWl); err {
		case nil:
			return
		case syscall.EINTR:
		case syscall.EINVAL:
			fmt.Fprintf(os.Stderr, "*SemLock.Lock: error(%v), so downgrade to local lock\n", err)
			s.local.Lock()
			return
		default:
			time.Sleep(time.Millisecond)
			fmt.Println("Lock", err)
		}
	}
}

func (s *SemLock) Unlock() {
	err := Semop(s.semid, hmsWUl)
	if err == syscall.EINVAL {
		fmt.Fprintf(os.Stderr, "*SemLock.Unlock: error(%v), so downgrade to local lock\n", err)
		s.local.Unlock()
	}
}

func (s *SemLock) Close() {
	err := Semctl(s.semid, IPC_RMID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "*SemLock.Close: error(%v), but ignore error\n", err)
	}
}
