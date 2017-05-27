package main

import (
	"github.com/liyue201/etcdsync"
	"log"
	"os"
	"sync"
	"time"
)

func testLock() {
	m := etcdsync.New("/mylock", 10, []string{"http://127.0.0.1:2379"})
	m.SetDebugLogger(os.Stdout)
	if m == nil {
		log.Printf("etcdsync.New failed")
	}
	err := m.Lock()
	if err != nil {
		log.Printf("etcdsync.Lock failed")
	} else {
		log.Printf("etcdsync.Lock OK")
	}

	log.Printf("Get the lock. Do something here.")

	err = m.Unlock()
	if err != nil {
		log.Printf("etcdsync.Unlock failed")
	} else {
		log.Printf("etcdsync.Unlock OK")
	}
}

func testTrylock() {
	wait := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		wait.Add(1)
		go func() {
			m := etcdsync.New("/mylock", 2, []string{"http://127.0.0.1:2379"})
			defer wait.Done()
			err := m.TryLock()
			if err != nil {
				log.Println("etcdsync.TryLock Failed:", err)
				return
			}
			log.Println("etcdsync.TryLock OK")
			time.Sleep(time.Second * 2)
			m.Unlock()
		}()
		time.Sleep(time.Second)
	}
	wait.Wait()
}

func main() {
	testLock()
	testTrylock()
}
