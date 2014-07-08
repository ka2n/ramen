package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Queue struct {
	Name string
	ID   int
	Wait int32
}

func sendrepyo(ch chan *Queue) {
	wg := &sync.WaitGroup{}
	for queue := range ch {
		log.Println(queue)
		wg.Add(1)
		go func(q *Queue) {
			time.Sleep(time.Second * time.Duration(q.Wait))
			log.Println(fmt.Sprintf("reply to Yo %s, %d", q.Name, q.ID))
			wg.Done()
		}(queue)
	}
	wg.Wait()
}

func serveyo(ch chan *Queue) {
	// start server
	http.HandleFunc("/api/callback/5", genHandler(5, ch))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("listen http://localhost:%s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func genHandler(wait int32, ch chan *Queue) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		uname := r.FormValue("username")
		fmt.Fprint(w, "OK", uname)

		go (func() {
			cntn := atomic.AddInt64(&cnt, 1)
			ch <- &Queue{
				ID:   int(cntn),
				Name: uname,
				Wait: wait,
			}
		})()
	}
}

var q chan *Queue
var cnt int64

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	q = make(chan *Queue)
	defer close(q)

	go sendrepyo(q)
	serveyo(q)
}
