package ipc

import (
	"sync"
	"unsafe"
)

type SyncShm struct {
	key     uint64
	shmid   int
	flock   *Flock
	shmaddr unsafe.Pointer
	lock    sync.RWMutex
}

func AttachSyncShm(path string, sizeIfCreate uint64) (*SyncShm, error) {
	key, shmid, err := FtokAndShmget(path, 0, sizeIfCreate, IPC_CREAT|IPC_RW)
	if err != nil {
		return nil, err
	}
	flock, err := NewFlock(path)
	if err != nil {
		return nil, err
	}
	shmaddr, err := Shmat(shmid, SHM_REMAP)
	if err != nil {
		return nil, err
	}
	return &SyncShm{
		key:     key,
		shmid:   shmid,
		flock:   flock,
		shmaddr: shmaddr,
	}, nil
}

func (s *SyncShm) Read() []byte {
	s.lock.RLock()
	_ = s.flock.RLock()
	defer func() {
		_ = s.flock.Unlock()
		s.lock.RUnlock()
	}()
	return Shmread(s.shmaddr)
}

func (s *SyncShm) Write(data []byte) error {
	s.lock.Lock()
	_ = s.flock.Lock()
	defer func() {
		_ = s.flock.Unlock()
		s.lock.Unlock()
	}()
	return Shmwrite(s.shmaddr, data)
}

func (s *SyncShm) Detach() error {
	s.lock.Lock()
	_ = s.flock.Lock()
	defer func() {
		_ = s.flock.Unlock()
		s.lock.Unlock()
	}()
	err := Shmdt(s.shmaddr)
	_ = s.flock.Close()
	return err
}

func (s *SyncShm) Remove() error {
	s.lock.Lock()
	_ = s.flock.Lock()
	defer func() {
		_ = s.flock.Unlock()
		s.lock.Unlock()
	}()
	err := Shmctl(s.shmid, IPC_RMID)
	_ = s.flock.Close()
	return err
}
