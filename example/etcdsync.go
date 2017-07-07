package main

import (
	"github.com/liyue201/etcdsync"
	"golang.org/x/net/context"
	"log"
	"os"
	"sync"
	"time"
)

func testLock() {
	m := etcdsync.New("/mylock", 10, []string{"http://120.24.44.201:4001"})
	m.SetDebugLogger(os.Stdout)
	if m == nil {
		log.Printf("etcdsync.New failed")
	}
	err := m.Lock(context.Background())
	if err != nil {
		log.Printf("etcdsync.Lock failed")
	} else {
		log.Printf("etcdsync.Lock OK")
	}

	log.Printf("Get the lock. Do something here.")

	err = m.Unlock(context.Background())
	if err != nil {
		log.Printf("etcdsync.Unlock failed")
	} else {
		log.Printf("etcdsync.Unlock OK")
	}
}

func testTryLock() {
	wait := sync.WaitGroup{}

	factory := etcdsync.NewMutexFactory([]string{"http://120.24.44.201:4001"})
	for i := 0; i < 20; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			m := factory.GetMutex("/mylock", 2)
			defer factory.ReleaseMutex(m)
			err := m.TryLock(context.Background())
			if err != nil {
				log.Println("etcdsync.TryLock Failed:", err)
				return
			}
			log.Println("etcdsync.TryLock OK")
			time.Sleep(time.Second * 2)
			m.Unlock(context.Background())
		}()
		time.Sleep(time.Second)
	}
	wait.Wait()
}

func main() {
	testLock()
	testTryLock()
}
