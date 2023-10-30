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

// 1. 如果没有 package main 作为程序的入口，就会出现 format error ， 尽管能够成功构建。
//

var build = "develop"

func main() {
	if _, err := maxprocs.Set(); err != nil {
		log.Printf("set procs err:%v", err)
	}

	g := runtime.GOMAXPROCS(0)
	log.Printf("services started env:%s update! CPU[%d]", build, g)
	defer log.Println("services ended!")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	log.Println("stopping services!", build)

}
