package api

import (
	"context"
	"flag"
	"fmt"
	"goTestProject/api/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

type Api struct{}

func New() *Api {
	return &Api{}
}

func (a *Api) Run() {
	r := router.Register()
	port := 7003
	flag.Parse()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("start listen : %s\n", err)
		}
	}()
	// if have two quit signal , this signal will priority capture ,also can graceful shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	logrus.Infof("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("Server Shutdown:", err)
	}
	logrus.Infof("Server exiting")
	os.Exit(0)
}
