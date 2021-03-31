package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"oss.navercorp.com/metis/metis-server/internal/log"
	"oss.navercorp.com/metis/metis-server/server"
)

const gracefulTimeout = 10 * time.Second

func newServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Starts Metis Server and runs until an interrupt is received.",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := server.New()
			if err != nil {
				return err
			}

			if err := s.Start(); err != nil {
				return err
			}

			if code := handleSignals(s); code != 0 {
				return fmt.Errorf("exit code: %d", code)
			}

			return nil
		},
	}
}

func handleSignals(s *server.Server) int {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	var sig os.Signal
	select {
	case s := <-signalCh:
		sig = s
	case <-s.ShutdownCh():
		// Server is already shutdown
		return 0
	}

	// Check graceful shutdown
	graceful := false
	if sig == syscall.SIGINT || sig == syscall.SIGTERM {
		graceful = true
	}

	log.Logger.Infof("Caught signal: %v", sig)

	gracefulCh := make(chan struct{})
	go func() {
		if err := s.Shutdown(graceful); err != nil {
			return
		}
		close(gracefulCh)
	}()

	log.Logger.Info("Gracefully shutting down server...")

	// Wait for shutdown or another signal
	select {
	// This is a case that handles the Unix signal registered in `signal.Notify()`.
	// Registered signal: SIGINT, SIGTERM, SIGHUP
	case <-signalCh:
		return 1
	case <-time.After(gracefulTimeout):
		return 1
	case <-gracefulCh:
		return 0
	}
}

func init() {
	cmd := newServerCommand()

	rootCmd.AddCommand(cmd)
}
