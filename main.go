package main

import (
	"fmt"
	"github.com/ka2n/ramen3/yo"
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
	ID   int64
	Wait int32
}

func reply(ch chan *Queue) {
	client := yo.DefaultClient
	wg := &sync.WaitGroup{}

	for queue := range ch {
		wg.Add(1)
		go func(q *Queue, c *yo.Client) {
			time.Sleep(time.Second * time.Duration(q.Wait))

			if err := c.Yo(q.Name); err != nil {
				log.Println("Yo failed with error:", err, q.Name)
			} else {
				log.Println("Yo", q.Name)
			}
			wg.Done()
		}(queue, client)
	}

	wg.Wait()
}

func serve(ch chan *Queue) {
	// start server
	http.HandleFunc("/hook/3", genHandler(3*60, ch))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("listen http://localhost:%s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func genHandler(wait int32, ch chan *Queue) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.FormValue("username")
		if u == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("Queue <- Yo to %s after %dsec", u, wait)
		fmt.Fprint(w, "OK,", u)

		go (func() {
			cntn := atomic.AddInt64(&cnt, 1)
			ch <- &Queue{
				ID:   cntn,
				Name: u,
				Wait: wait,
			}
		})()
	}
}

var q chan *Queue
var cnt int64

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cnt = 0
	q = make(chan *Queue)
	defer close(q)

	go reply(q)
	serve(q)
}
