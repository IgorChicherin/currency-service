package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/IgorChicherin/currency-service/config"
	"github.com/IgorChicherin/currency-service/db"
	"github.com/IgorChicherin/currency-service/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	environment := flag.String("e", "develop", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	config.Init(*environment)
	db.Init()

	srv, shed, err := server.Run()

	if err != nil {
		log.Panicf("Couldn't start server")
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	shed.Stop()

	log.Println("Server exiting")
}
