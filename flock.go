package ipc

import (
	"fmt"
	"os"
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

// Lock set a writing lock
func (l *Flock) Lock(nonBlock ...bool) error {
	return l.lock(true, nonBlock...)
}

// RLock set a reading lock
func (l *Flock) RLock(nonBlock ...bool) error {
	return l.lock(false, nonBlock...)
}

func (l *Flock) lock(writable bool, nonBlock ...bool) error {
	how := syscall.LOCK_SH
	if writable {
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

// Unlock free a lock
func (l *Flock) Unlock() error {
	return syscall.Flock(int(l.f.Fd()), syscall.LOCK_UN)
}

// Close closes the file and free all locks
func (l *Flock) Close() error {
	return l.f.Close()
}
