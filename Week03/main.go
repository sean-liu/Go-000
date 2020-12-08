package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"

	"golang.org/x/sync/errgroup"
)

func main() {
	eg := errgroup.Group{}

	eg.Go(func() error {
		return http.ListenAndServe(":8090", nil)
	})

	eg.Go(func() error {
		sigs := make(chan os.Signal)
		signal.Notify(sigs, syscall.SIGINT|syscall.SIGTERM|syscall.SIGKILL)
		<-sigs
		return errors.New("received signal to exit")
	})

	eg.Wait()
}
