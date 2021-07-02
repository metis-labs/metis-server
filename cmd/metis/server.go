/*
 * Copyright 2021-present NAVER Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/metis-labs/metis-server/internal/log"
	"github.com/metis-labs/metis-server/server"
)

const gracefulTimeout = 10 * time.Second

var (
	mongoConnectionTimeoutSec int
	mongoPingTimeoutSec       int
	conf                      = server.NewConfig()
)

func newServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Starts Metis Server and runs until an interrupt is received.",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf.Mongo.ConnectionTimeoutSec = time.Duration(mongoConnectionTimeoutSec)
			conf.Mongo.PingTimeoutSec = time.Duration(mongoPingTimeoutSec)
			s, err := server.New(conf)
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
	cmd.Flags().IntVar(
		&conf.RPC.Port,
		"rpc-port",
		server.DefaultRPCPort,
		"RPC port",
	)

	cmd.Flags().StringVar(
		&conf.Yorkie.RPCAddr,
		"yorkie-rpc-addr",
		server.DefaultYorkieRPCAddr,
		"Yorkie's RPC Address",
	)

	cmd.Flags().StringVar(
		&conf.Mongo.ConnectionURI,
		"mongo-connection-uri",
		server.DefaultMongoConnectionURI,
		"MongoDB's connection URI",
	)
	cmd.Flags().IntVar(
		&mongoConnectionTimeoutSec,
		"mongo-connection-timeout-sec",
		server.DefaultMongoConnectionTimeoutSec,
		"Mongo DB's connection timeout in seconds",
	)
	cmd.Flags().StringVar(
		&conf.Mongo.Database,
		"mongo-database",
		server.DefaultMongoDatabase,
		"Metis database name in MongoDB",
	)
	cmd.Flags().IntVar(
		&mongoPingTimeoutSec,
		"mongo-ping-timeout-sec",
		server.DefaultMongoPingTimeoutSec,
		"Mongo DB's ping timeout in seconds",
	)

	rootCmd.AddCommand(cmd)
}
