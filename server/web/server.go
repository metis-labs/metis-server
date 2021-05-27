package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/yorkie-team/yorkie/pkg/document/key"
	yorkieTypes "github.com/yorkie-team/yorkie/pkg/types"

	"oss.navercorp.com/metis/metis-server/internal/log"
	"oss.navercorp.com/metis/metis-server/server/database"
	"oss.navercorp.com/metis/metis-server/server/types"
)

const (
	writeTimeout = 15 * time.Second
	readTimeout  = 15 * time.Second
)

// Config is the configuration for creating a Server instance.
type Config struct {
	Port int
}

// Server is a server that processes the web requested such as authentication webhook.
type Server struct {
	conf       *Config
	db         database.Database
	httpServer *http.Server
}

// NewServer creates a new instance of Server.
func NewServer(conf *Config, db database.Database) (*Server, error) {
	server := &Server{
		conf: conf,
		db:   db,
	}

	r := mux.NewRouter()
	r.HandleFunc("/auth", server.HandleAuth)
	r.Use(elapsedTimeMiddleware)

	server.httpServer = &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", conf.Port),
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
	}

	return server, nil
}

// Start starts to handle requests on incoming connections.
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.httpServer.Addr)
	if err != nil {
		return err
	}

	log.Logger.Infof("WebServer is running on %d", s.conf.Port)

	go func() {
		if err := s.httpServer.Serve(listener); err != nil {
			log.Logger.Errorf("fail to serve: %s", err.Error())
		}
	}()

	return nil
}

// HandleAuth handles the given authorization webhook request.
func (s *Server) HandleAuth(w http.ResponseWriter, r *http.Request) {
	req, err := yorkieTypes.NewAuthWebhookRequest(r.Body)
	if err != nil {
		log.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := s.handleAuth(req)
	if err != nil {
		log.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resBody, err := json.Marshal(resp)
	if err != nil {
		log.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(resBody); err != nil {
		log.Logger.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleAuth(req *yorkieTypes.AuthWebhookRequest) (*yorkieTypes.AuthWebhookResponse, error) {
	switch req.Method {
	case yorkieTypes.AttachDocument, yorkieTypes.DetachDocument, yorkieTypes.PushPull:
		docKey, err := key.FromBSONKey(req.Attributes[0].Key)
		if err != nil {
			return nil, err
		}

		project, err := s.db.FindProject(
			types.CtxWithUserID(context.Background(), req.Token),
			types.ID(docKey.Document),
		)
		if err != nil {
			return nil, err
		}

		if project.Owner != req.Token {
			return &yorkieTypes.AuthWebhookResponse{
				Allowed: false,
				Reason:  "user does not have permission to the document",
			}, nil
		}
	}

	return &yorkieTypes.AuthWebhookResponse{Allowed: true}, nil
}

// GracefulStop stops the server gracefully.
func (s *Server) GracefulStop() {
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		log.Logger.Error(err)
	}
}

// Stop stops the server. It immediately closes all open connections and listeners.
func (s *Server) Stop() {
	if err := s.httpServer.Close(); err != nil {
		log.Logger.Error(err)
	}
}
