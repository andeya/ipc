package ipc

import (
	"sync"
	"sync/atomic"
)

// IPLock inter-process lock
type IPLock struct {
	file       *Flock
	local      sync.RWMutex
	rlockCount int32
}

func NewIPLock(path string) (*IPLock, error) {
	flock, err := NewFlock(path)
	if err != nil {
		return nil, err
	}
	return &IPLock{file: flock, local: sync.RWMutex{}}, nil
}

func (l *IPLock) RLock() {
	l.local.RLock()
	atomic.AddInt32(&l.rlockCount, 1)
	_ = l.file.ShareLock()
}

func (l *IPLock) RUnlock() {
	rlockCount := atomic.AddInt32(&l.rlockCount, -1)
	if rlockCount == 0 {
		_ = l.file.UnlockAll()
	}
	l.local.RUnlock()
}

func (l *IPLock) Lock() {
	l.local.Lock()
	_ = l.file.ExclusiveLock()
}

func (l *IPLock) Unlock() {
	_ = l.file.UnlockAll()
	l.local.Unlock()
}

func (l *IPLock) Close() {
	l.local.Lock()
	defer func() {
		l.local.Unlock()
		recover()
	}()
	_ = l.file.Close()
}
