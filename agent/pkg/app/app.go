package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli"
)

var appInstance EdgeApp

// EdgeApp is app instance defination
type EdgeApp struct {
	appInstance *cli.App
	flags       []cli.Flag
}

// NewInstance create a new instance of the application
func NewInstance() *EdgeApp {

	return nil
}

// Run is the start function
func (app *EdgeApp) Run() error {
	fmt.Println("app running")
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			time.Sleep(time.Duration(1) * time.Second)
			fmt.Println("still running")
		}
	}()
	select {
	case <-quit:
		fmt.Printf("app quit")
	}
	return nil
}

func registerModules() error {
	return nil
}
