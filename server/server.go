package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Service struct {
	Ser  *http.Server
	Stop chan struct{}
}

func (s *Service) Server(ctx context.Context) error {
	http.HandleFunc("./", func(writer http.ResponseWriter, request *http.Request) {
		s.home(ctx, writer, request)
	})

	return http.ListenAndServe(":8080", nil)
}

func (s *Service) home(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	select {
	case <-ctx.Done():
		s.Stop <- struct{}{}
		return
	default:
		fmt.Fprintf(w, "hello world")
	}
}

func (s *Service) Shutdown(ctx context.Context) {
	defer func() {
		s.Ser.Shutdown(ctx)
		close(s.Stop)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	select {
	case <-signalChan:
		fmt.Println("stop")
		return
	case <-s.Stop:
		fmt.Println("error stop")
		return
	}
}
