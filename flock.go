package ipc

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
)

// Flock file lock
type Flock struct {
	path string
	f    *os.File
}

// NewFlock new a flock
func NewFlock(path string) (*Flock, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &Flock{
		path: path,
		f:    f,
	}, nil
}

// ExclusiveLock set an exclusive lock
func (l *Flock) ExclusiveLock(nonBlock ...bool) error {
	return l.lock(true, nonBlock...)
}

// ShareLock set a sharing lock
func (l *Flock) ShareLock(nonBlock ...bool) error {
	return l.lock(false, nonBlock...)
}

func (l *Flock) lock(exclusive bool, nonBlock ...bool) error {
	how := syscall.LOCK_SH
	if exclusive {
		how = syscall.LOCK_EX
	}
	if len(nonBlock) > 0 && nonBlock[0] {
		how |= syscall.LOCK_NB
	}
	err := syscall.Flock(int(l.f.Fd()), how)
	if err != nil {
		return fmt.Errorf("cannot flock pathectory %s - %s", l.path, err)
	}
	return nil
}

// UnlockAll release all locks in the current process
func (l *Flock) UnlockAll() error {
	return syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
}

// Close closes the file and free all locks
func (l *Flock) Close() error {
	return l.f.Close()
}

func (l *Flock) RWMutex() IPLock {
	return &FRWMutex{file: l, local: sync.RWMutex{}}
}

// FRWMutex inter-process lock by flock
type FRWMutex struct {
	file       *Flock
	local      sync.RWMutex
	rlockCount int32
}

func (l *FRWMutex) RLock() {
	l.local.RLock()
	atomic.AddInt32(&l.rlockCount, 1)
	_ = l.file.ShareLock()
}

func (l *FRWMutex) RUnlock() {
	rlockCount := atomic.AddInt32(&l.rlockCount, -1)
	if rlockCount == 0 {
		_ = l.file.UnlockAll()
	}
	l.local.RUnlock()
}

func (l *FRWMutex) Lock() {
	l.local.Lock()
	_ = l.file.ExclusiveLock()
}

func (l *FRWMutex) Unlock() {
	_ = l.file.UnlockAll()
	l.local.Unlock()
}

func (l *FRWMutex) Close() {
	l.local.Lock()
	defer func() {
		l.local.Unlock()
		recover()
	}()
	_ = l.file.Close()
}
