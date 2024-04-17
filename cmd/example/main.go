package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/eliofery/closer"
	"log"
	"math/rand/v2"
	"net/http"
	"time"
)

// Example of usage.
func main() {
	// Create new instance of Closer.
	clr := closer.New()

	// Create http server.
	rt := http.NewServeMux()
	rt.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("test"))
	})

	srv := http.Server{
		Addr:    ":8080",
		Handler: rt,
	}

	go func() {
		fmt.Println("server started")
		if err := srv.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	// Close connection to http server.
	clr.Add(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		return srv.Shutdown(ctx)
	})

	// Close connection to postgres database.
	clr.Add(func() error {
		time.Sleep(1 * time.Second)

		if rand.IntN(2) == 1 {
			return nil
		}

		return errors.New("error close postgres")
	})

	// Close connection to redis database.
	clr.Add(func() error {
		time.Sleep(2 * time.Second)

		if rand.IntN(2) == 1 {
			return nil
		}

		return errors.New("error close redis")
	})

	// Close connection to rabbitmq.
	clr.Add(func() error {
		return nil
	})

	// Wait signals about closing.
	clr.Wait()

	// Close all functions.
	log.Printf("\n%v", clr.Close())

	// or
	// defer func() {
	//   log.Printf("\n%v", clr.Close())
	// }()
	// clr.Wait()
}
