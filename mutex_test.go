package etcdsync

import (
	"testing"
	"golang.org/x/net/context"
	"github.com/coreos/etcd/client"
	"log"
	"time"
)

func init() {
	debug = true
}

func TestMutex(t *testing.T) {
	log.SetFlags(log.Ltime|log.Ldate|log.Lshortfile)
	key := "/etcdsync"
	m := New(key, "100", []string{"http://127.0.0.1:2379"})
	if m == nil {
		t.Errorf("New Mutex ERROR")
	}
	err := m.Lock()
	if err != nil {
		t.Errorf("failed")
	}

	//do something here

	err = m.Unlock()
	if err != nil {
		t.Errorf("failed")
	}

	_, err = m.kapi.Get(context.Background(), key, nil)
	if e, ok := err.(client.Error); !ok {
		t.Errorf("Get key %v failed from etcd", key)
	} else if e.Code != client.ErrorCodeKeyNotFound {
		t.Errorf("ERROR %v", err)
	} else {
		println("test OK")
	}
}


func TestLockConcurrently(t *testing.T) {
	log.Println("\n\n")
	slice := make([]int, 0, 3)
	lockKey := "/etcd_sync"
	m1 := New(lockKey, "1", []string{"http://127.0.0.1:2379"})
	m2 := New(lockKey, "2", []string{"http://127.0.0.1:2379"})
	m3 := New(lockKey, "3", []string{"http://127.0.0.1:2379"})
	if m1 == nil || m2 == nil || m3 == nil {
		t.Errorf("New Mutex ERROR")
	}
	m1.Lock()
	ch1 := make(chan bool)
	go func() {
		ch2 := make(chan bool)
		m2.Lock()
		go func() {
			m3.Lock()
			slice = append(slice, 2)
			m3.Unlock()
			log.Println("\n\n")
			ch2 <- true
		}()
		slice = append(slice, 1)
		time.Sleep(1 * time.Second)
		m2.Unlock()
		log.Println("\n\n")
		<-ch2
		ch1 <- true
	}()
	slice = append(slice, 0)
	time.Sleep(1 * time.Second)
	m1.Unlock()
	log.Println("\n\n")
	<-ch1
	if len(slice) != 3 {
		t.Fail()
	}
	for n, i := range slice {
		println("n,i:", n, i)
		if n != i {
			t.Fail()
		}
	}
}