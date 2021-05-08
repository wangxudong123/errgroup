package main

import (
	"context"
	"errgroups/server"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

func main() {
	d := time.Now().Add(5 * time.Second)
	ctx, _ := context.WithDeadline(context.Background(), d)
	g, ctx := errgroup.WithContext(ctx)

	s := &server.Service{
		Ser:  &http.Server{Addr: ":8080"},
		Stop: make(chan struct{}),
	}

	g.Go(func() error {
		return s.Server(ctx)
	})

	g.Go(func() error {
		s.Shutdown(ctx)
		return nil
	})

	if err := g.Wait(); err != nil {
		fmt.Println(err)
		s.Stop <- struct{}{}
	}

}
