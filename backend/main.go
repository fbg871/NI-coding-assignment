package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	ff "github.com/peterbourgon/ff/v3"
)

func main() {
	// Read runtime configuration
	fs := flag.NewFlagSet("komp-registry", flag.ExitOnError)
	var (
		mysqlAddr     = fs.String("mysql-addr", ":3306", "")
		mysqlUser     = fs.String("mysql-user", "root", "")
		mysqlPassword = fs.String("mysql-password", "TopSecret", "")
		listenAddr    = fs.String("listen-port", ":3001", "")
	)

	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix())
	if err != nil {
		log.Printf("failed to parse environment configuration")
		return
	}

	// Connect to mysql database
	conn, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", *mysqlUser, *mysqlPassword, *mysqlAddr, "komp_registry"))
	if err != nil {
		log.Printf("failed to open connection to database: %s", err)
		return
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Printf("failed to open ping database: %s", err)
		return
	}

	// Setup http api server
	server := newHttpServer(*listenAddr, conn)

	// Asynchronously listen for stop signal
	go func() {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
		sig := <-signalChannel

		fmt.Printf("\n")
		log.Printf("received %s signal, backend closing...", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			log.Printf("failed to shutdown backend: %s", err)
		}
	}()

	// Start server (blocks until stop signal)
	log.Printf("backend listening on %s", server.Addr)
	err = server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Printf("failed to start backend: %s", err)
	} else {
		log.Print("backend closed successfully")
	}
}
