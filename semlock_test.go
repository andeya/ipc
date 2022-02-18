package ipc

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSemLock1(t *testing.T) {
	// clean
	semLock, err := NewSemLock(1)
	assert.NoError(t, err)
	fmt.Println("semLock.Close()")
	semLock.Close()

	semLock, err = NewSemLock(1)
	assert.NoError(t, err)
	defer semLock.Close()

	t.Log(semLock, err)
	var wg sync.WaitGroup
	semLock.RLock()
	wg.Add(1)
	go func() {
		defer wg.Done()
		semLock.RLock()
		wg.Add(1)
		go func() {
			defer wg.Done()
			t.Log("semLock.Lock() 1")
			semLock.Lock()
			t.Log("semLock.Lock() 2")
			wg.Add(1)
			go func() {
				defer wg.Done()
				semLock.Lock()
				t.Log("4: Lock")
				semLock.Unlock()
				t.Log("4: Unlock")
			}()
			time.Sleep(time.Second * 1)
			t.Log("3: Lock")
			semLock.Unlock()
			t.Log("3: Unlock")
		}()
		time.Sleep(time.Second * 2)
		t.Log("1: RLock")
		semLock.RUnlock()
		t.Log("1: RUnlock")
	}()
	time.Sleep(time.Second * 3)
	t.Log("2: RLock")
	semLock.RUnlock()
	t.Log("2: RUnlock")
	wg.Wait()
}

func TestSemLock2(t *testing.T) {
	semLock, err := NewSemLock(1)
	assert.NoError(t, err)
	t.Log(semLock, err)
	semLock.Lock()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		semLock.Lock()
		t.Log("2: Lock")
		semLock.Unlock()
		t.Log("2: Unlock")
	}()
	time.Sleep(time.Second * 1)
	t.Log("1: Lock")
	semLock.Unlock()
	t.Log("1: Unlock")
	wg.Wait()
}

func TestSemLock3(t *testing.T) {
	semLock, err := NewSemLock(1)
	assert.NoError(t, err)
	defer semLock.Close()
	t.Log(semLock, err)
	semLock.Lock()
	go func() {
		time.Sleep(time.Second * 2)
		t.Log("1")
		semLock.Unlock()
	}()
	semLock.Lock()
	t.Log("2")
	semLock.Unlock()
}

func TestSemLock4(t *testing.T) {
	semLock, err := NewSemLock(1)
	assert.NoError(t, err)
	defer semLock.Close()
	t.Log(semLock, err)
	semLock.Lock()
	go func() {
		time.Sleep(time.Second * 2)
		t.Log("1")
		semLock.Unlock()
	}()
	semLock.RLock()
	t.Log("2")
	semLock.RUnlock()
}

func TestSemLock5(t *testing.T) {
	semLock, err := NewSemLock(1)
	assert.NoError(t, err)
	defer semLock.Close()
	t.Log(semLock, err)
	semLock.Lock()
	t.Log("1")
	semLock.Unlock()
	semLock.Lock()
	t.Log("2")
	semLock.Unlock()
}

func TestSemLock6(t *testing.T) {
	semLock, err := NewSemLock(1)
	assert.NoError(t, err)
	defer semLock.Close()
	t.Log(semLock, err)
	semLock.RLock()
	go func() {
		time.Sleep(time.Second * 2)
		t.Log("2")
		semLock.RUnlock()
		t.Log("3")
	}()
	t.Log("1")
	semLock.Lock()
	t.Log("4")
	semLock.Unlock()
}

func TestSemLock_Close(t *testing.T) {
	semLock, err := NewSemLock(1)
	assert.NoError(t, err)
	t.Log(semLock, err)
	semLock.Close()
}
