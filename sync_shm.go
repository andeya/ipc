package ipc

import (
	"unsafe"
)

type SyncShm struct {
	key     uint64
	shmid   int
	shmaddr unsafe.Pointer
	lock    *IPLock
}

func AttachSyncShm(path string, sizeIfCreate uint64) (*SyncShm, error) {
	key, shmid, err := FtokAndShmget(path, 0, sizeIfCreate, IPC_CREAT|IPC_RW)
	if err != nil {
		return nil, err
	}
	lock, err := NewIPLock(path)
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
		lock:    lock,
		shmaddr: shmaddr,
	}, nil
}

func (s *SyncShm) Read() []byte {
	s.lock.RLock()
	defer func() {
		s.lock.RUnlock()
	}()
	return Shmread(s.shmaddr)
}

func (s *SyncShm) Write(data []byte) error {
	s.lock.Lock()
	defer func() {
		s.lock.Unlock()
	}()
	return Shmwrite(s.shmaddr, data)
}

func (s *SyncShm) Detach() error {
	s.lock.Lock()
	defer func() {
		s.lock.Unlock()
		s.lock.Close()
	}()
	err := Shmdt(s.shmaddr)
	return err
}

func (s *SyncShm) Remove() error {
	s.lock.Lock()
	defer func() {
		s.lock.Unlock()
		s.lock.Close()
	}()
	err := Shmctl(s.shmid, IPC_RMID)
	return err
}
