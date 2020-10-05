package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	syscall "golang.org/x/sys/unix"

	"github.com/voyagegroup/pidfile"
)

func main() {
	p, err := pidfile.Create("var/run/test.pid")
	if err != nil {
		log.Fatalln("failed creating pidfile:", err)
	}
	defer p.Remove()

	log.Println("server started with PID", syscall.Getpid())

	run()

	log.Println("server stopped")
}

func run() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	stopSigs := make(chan os.Signal, 1)
	contSig := make(chan os.Signal, 1)

	signal.Notify(stopSigs, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(contSig, syscall.SIGCONT)

	for {
		select {
		case <-ticker.C:
			log.Println("ticked!")
		case <-contSig:
			log.Println("got a signal to continue!")
		case <-stopSigs:
			log.Println("got a signal to stop!")
			return
		default:
		}
	}
}
