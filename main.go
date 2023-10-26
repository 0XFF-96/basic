package main

import (
	_ "go.uber.org/automaxprocs"
	"go.uber.org/automaxprocs/maxprocs"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var build = "develop"

func main() {

	if _, err := maxprocs.Set(); err != nil {
		log.Printf("set procs err:%v", err)
	}
	
	g := runtime.GOMAXPROCS(0)
	log.Printf("service started env:%s update! CPU[%d]", build, g)
	defer log.Println("service ended!")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	log.Println("stopping service!", build)

}
