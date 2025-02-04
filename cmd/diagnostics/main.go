package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ledgerwatch/diagnostics/api"
	"github.com/ledgerwatch/diagnostics/internal/logging"
	"github.com/ledgerwatch/diagnostics/internal/sessions"
)

func main() {
	// Retrieving flags and configurations
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//set up logger to implement log rotation
	logging.SetupLogger(logDirPath, logFileName, logFileSizeMax, logFilesAgeMax, logFilesMax, logCompress)

	// Use of system calls SIGINT and SIGTERM signals that cause a gracefully  stop.
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	// Initialize Services
	cache, err := sessions.NewCache(5, 5)

	if err != nil {
		log.Fatalf("session cache creation  failed: %v", err)
	}

	// Passing in the services to REST layer
	handlers := api.NewHandler(
		api.APIServices{
			StoreSession: cache,
		})

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", listenAddr, listenPort),
		Handler:           handlers,
		MaxHeaderBytes:    1 << 20,
		ReadHeaderTimeout: 1 * time.Minute,
	}

	go func() {
		err := srv.ListenAndServe()

		if err != nil {
			log.Fatal(err)
		}
	}()

	open(fmt.Sprintf("http://%s:%d", listenAddr, listenPort))

	// Graceful and eager terminations
	switch s := <-signalCh; s {
	case syscall.SIGTERM:
		log.Println("Terminating gracefully.")
		if err := srv.Shutdown(context.Background()); err != http.ErrServerClosed {
			log.Println("Failed to shutdown server:", err)
		}
	case syscall.SIGINT:
		log.Println("Terminating eagerly.")
		os.Exit(-int(syscall.SIGINT))
	}
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
