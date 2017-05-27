# etcdsync

etcdsync is a distributed lock library in Go using etcd. It easy to use like sync.Mutex.


In fact, there are many similar implementation which are all obsolete 
depending on library `github.com/coreos/go-etcd/etcd` which is official marked `deprecated`,
and the usage is a little bit complicated. 
Otherwise this library is very very simple. The usage is simple, the code is simple.

## Import
    
    go get github.com/zieckey/etcdsync

## Simplest usage

Steps:

1. m := etcdsync.New()
2. m.Lock()
3. Do your business here
4. m.Unlock()

```go
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

func testTryLock() {
	wait := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			m := etcdsync.New("/mylock", 2, []string{"http://127.0.0.1:2379"})
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
	testTryLock()
}
```

## Test

You need a etcd instance running on http://localhost:2379, then:

    go test
    
