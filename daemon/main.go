package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"flag"
	"github.com/k3rn3l-p4n1c/goaway/agent"
	"github.com/k3rn3l-p4n1c/goaway/configuration"
	"github.com/k3rn3l-p4n1c/goaway/daemon/core"
	"github.com/takama/daemon"
	"net/rpc"
)

const (
	// name of the service
	name        = "goaway"
	description = "GoAway cloud platform"
)

//	dependencies that are NOT required by the service, but might be used
var dependencies = []string{"dummy.service"}

var stdlog, errlog *log.Logger

var (
	isBootstrap = flag.Bool("bootstrap", false, "Run server as bootstrap")
	joinAddress = flag.String("join", "", "Join address")
)

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {

	usage := "Usage: goawayd install | remove | start | stop | status"

	// if received any kind of grpc, do it
	if len(flag.Args()) > 1 {
		command := flag.Args()[0]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		case "run":
			break
		default:
			return usage, nil
		}
	}

	println(*isBootstrap, *joinAddress)
	configuration.GetInstance().Set("master.bootstrap", isBootstrap)
	configuration.GetInstance().Set("slave.join", joinAddress)

	// Do something, call your goroutines, etc
	rpc.Register(&core.Handler{})

	// Wait for incoming connections
	listener, err := net.Listen("unix", "/tmp/goaway.sock")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go rpc.Accept(listener)

	go agent.Start()

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// loop work cycle with accept connections or interrupt
	// by system signal
	for {
		select {
		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)
			stdlog.Println("Stoping listening on ", listener.Addr())
			listener.Close()
			if killSignal == os.Interrupt {
				return "Daemon was interruped by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}

	// never happen, but need to complete code
	return usage, nil
}

func init() {
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

func main() {
	flag.Parse()

	srv, err := daemon.New(name, description, dependencies...)

	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
