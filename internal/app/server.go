package app

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/FredericoBento/HandGame/internal/handlers"
)

var (
	ErrAuthHandlerNotFound = errors.New("AuthHandler was not found")

	ErrServerCouldNotRan = errors.New("Server could not be ran")

	ErrServerCouldNotSetupRoutes = errors.New("Server could not setup http routes with handlers")
)

type ServerHandlers struct {
	authHandler *handlers.AuthHandler
}

type Server struct {
	Host        string
	Port        int
	httpServer  *http.Server
	router      *http.ServeMux
	authRouter  *http.ServeMux
	adminRouter *http.ServeMux
	handlers    *ServerHandlers
}

type ServerOption func(*Server)

func WithHost(host string) ServerOption {
	return func(s *Server) {
		s.Host = host
	}
}

func WithPort(port int) ServerOption {
	return func(s *Server) {
		s.Port = port
	}
}

func WithHandlers(handlers *ServerHandlers) ServerOption {
	return func(s *Server) {
		s.handlers = handlers
	}
}

func NewServer(opts ...ServerOption) *Server {
	server := &Server{
		httpServer:  nil,
		router:      http.NewServeMux(),
		authRouter:  http.NewServeMux(),
		adminRouter: http.NewServeMux(),
		handlers:    nil,
	}
	for _, option := range opts {
		option(server)
	}

	return server
}

func NewServerHandlers(authH *handlers.AuthHandler) *ServerHandlers {
	return &ServerHandlers{
		authHandler: authH,
	}
}

func (s *Server) Init() error {
	slog.Info("Initiating Server...")

	err := s.setupRoutes()
	if err != nil {
		slog.Error("Server could not be initiated")
		slog.Error(err.Error())
		return ErrServerCouldNotSetupRoutes
	}

	return nil
}

func (s *Server) setupRoutes() error {
	if s.handlers.authHandler == nil {
		return ErrAuthHandlerNotFound
	}

	s.authRouter.HandleFunc("GET /sign-in", s.handlers.authHandler.GetSignIn)
	s.authRouter.HandleFunc("GET /sign-up", s.handlers.authHandler.GetSignUp)

	// s.authRouter.HandleFunc("/dashboard", s..Dashboard)

	s.router.Handle("/", s.authRouter)
	// s.router.Handle("/admin/", http.StripPrefix("/admin", s.adminRouter))

	return nil
}

func (s *Server) Run() error {
	addr := s.Host + ":" + strconv.Itoa(s.Port)

	slog.Info("Server is running on " + addr)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	err := s.httpServer.ListenAndServe()
	defer s.httpServer.Close()
	if err != nil {
		slog.Error(err.Error())
		return ErrServerCouldNotRan
	}
	return nil
}

func (s *Server) Shutdown() error {
	slog.Info("Server is shuting down...")
	return nil
}
