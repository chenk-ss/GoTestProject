package main

import (
	"fmt"
	"goTestProject/api"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Hello goTestProject")

	api.New().Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	fmt.Println("Server exiting")
}
