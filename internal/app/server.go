package app

import (
	"errors"
	"github.com/FredericoBento/HandGame/internal/middleware"
	"log/slog"
	"net/http"
	"strconv"
)

var (
	ErrAuthHandlerNotFound = errors.New("AuthHandler was not found")

	ErrServerCouldNotRan = errors.New("Server could not be ran")

	ErrServerCouldNotSetupRoutes = errors.New("Server could not setup http routes with handlers")
)

type ServerHandlers struct {
	AuthHandler  http.Handler
	HomeHandler  http.Handler
	AdminHandler http.Handler
}

type Server struct {
	Host        string
	Port        int
	HttpServer  *http.Server
	Router      *http.ServeMux
	AuthRouter  *http.ServeMux
	AdminRouter *http.ServeMux
	Handlers    *ServerHandlers
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
		s.Handlers = handlers
	}
}

func NewServer(opts ...ServerOption) *Server {
	server := &Server{
		HttpServer:  nil,
		Router:      http.NewServeMux(),
		AuthRouter:  http.NewServeMux(),
		AdminRouter: http.NewServeMux(),
		Handlers:    nil,
	}
	for _, option := range opts {
		option(server)
	}

	return server
}

func NewServerHandlers(authH http.Handler, homeH http.Handler, adminH http.Handler) *ServerHandlers {
	return &ServerHandlers{
		AuthHandler:  authH,
		HomeHandler:  homeH,
		AdminHandler: adminH,
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
	if s.Handlers.AuthHandler == nil {
		return ErrAuthHandlerNotFound
	}

	AuthHandlerMiddlewares := middleware.StackMiddleware(
		middleware.Logger,
		middleware.SecureHeadersMiddleware,
	)

	HomeHandlerMiddlewares := middleware.StackMiddleware(
		middleware.Logger,
		middleware.SecureHeadersMiddleware,
	)

	s.AuthRouter.Handle("/sign-in", s.Handlers.AuthHandler)
	s.AuthRouter.Handle("/sign-up", s.Handlers.AuthHandler)
	s.AuthRouter.Handle("/logout", s.Handlers.AuthHandler)

	s.Router.Handle("/", AuthHandlerMiddlewares(s.AuthRouter))

	fs := http.FileServer(http.Dir("./assets"))
	s.Router.Handle("/assets/", HomeHandlerMiddlewares(http.StripPrefix("/assets", fs)))

	return nil
}

func (s *Server) Run() error {
	addr := s.Host + ":" + strconv.Itoa(s.Port)

	slog.Info("Server is running on " + addr)

	s.HttpServer = &http.Server{
		Addr:    addr,
		Handler: middleware.BlockRoutes(s.Router),
	}

	go func() {
		err := s.HttpServer.ListenAndServe()
		// defer s.HttpServer.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	return nil
}

func (s *Server) Shutdown() error {
	slog.Info("Server is shuting down...")
	return s.HttpServer.Close()
}

func (s *Server) BlockAppRoutes(appPrefix string) {
	middleware.BlockRoute(appPrefix)
}

func (s *Server) UnblockAppRoutes(appPrefix string) {
	middleware.UnblockRoute(appPrefix)
}
