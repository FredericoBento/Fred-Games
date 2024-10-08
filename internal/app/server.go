package app

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/FredericoBento/HandGame/internal/logger"
	"github.com/FredericoBento/HandGame/internal/middleware"
	"golang.org/x/net/websocket"
)

var (
	ErrAuthHandlerNotFound       = errors.New("AuthHandler was not found")
	ErrServerCouldNotRan         = errors.New("Server could not be ran")
	ErrServerCouldNotSetupRoutes = errors.New("Server could not setup http routes with handlers")
)

type ServerHandlers struct {
	AuthHandler     http.Handler
	HomeHandler     http.Handler
	AdminHandler    http.Handler
	HandGameHandler http.Handler
	PongHandler     http.Handler
}

type Server struct {
	Host        string
	Port        int
	HttpServer  *http.Server
	Router      *http.ServeMux
	AuthRouter  *http.ServeMux
	AdminRouter *http.ServeMux
	Handlers    *ServerHandlers
	log         *slog.Logger
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
	lo, err := logger.NewServerLogger("server", "", true)
	if err != nil {
		lo = slog.Default()
	}
	server := &Server{
		HttpServer:  nil,
		Router:      http.NewServeMux(),
		AuthRouter:  http.NewServeMux(),
		AdminRouter: http.NewServeMux(),
		Handlers:    nil,
		log:         lo,
	}
	for _, option := range opts {
		option(server)
	}

	return server
}

func NewServerHandlers(authH http.Handler, adminH http.Handler, homeH http.Handler, handGameH http.Handler, pongH http.Handler) *ServerHandlers {
	return &ServerHandlers{
		AuthHandler:     authH,
		AdminHandler:    adminH,
		HomeHandler:     homeH,
		HandGameHandler: handGameH,
		PongHandler:     pongH,
	}
}

func (s *Server) Init() error {
	s.log.Info("Initiating Server...")

	err := s.setupRoutes()
	if err != nil {
		s.log.Error("Server could not be initiated")
		s.log.Error(err.Error())
		return ErrServerCouldNotSetupRoutes
	}

	return nil
}

func (s *Server) setupRoutes() error {
	if s.Handlers.AuthHandler == nil {
		return ErrAuthHandlerNotFound
	}

	standardMiddlewares := middleware.StackMiddleware(
		middleware.Logger,
		middleware.SecureHeadersMiddleware,
	)

	authHandlerMiddlewares := middleware.StackMiddleware(
		standardMiddlewares,
	)

	adminHandlerMiddlewares := middleware.StackMiddleware(
		standardMiddlewares,
		middleware.RequiredAdmin,
	)

	// Auth Routes
	s.AuthRouter.Handle("/sign-in", s.Handlers.AuthHandler)
	s.AuthRouter.Handle("/sign-up", s.Handlers.AuthHandler)
	s.AuthRouter.Handle("/logout", s.Handlers.AuthHandler)

	s.Router.Handle("/", authHandlerMiddlewares(s.AuthRouter))

	// Admin Routes
	s.AdminRouter.Handle("/dashboard", s.Handlers.AdminHandler)
	s.AdminRouter.Handle("/users", s.Handlers.AdminHandler)
	s.AdminRouter.Handle("/", s.Handlers.AdminHandler)

	s.Router.Handle("/admin/", adminHandlerMiddlewares(s.AdminRouter))

	// App Homepage
	s.Router.Handle("/home", authHandlerMiddlewares(s.Handlers.HomeHandler))

	// Fileserver
	fs := http.FileServer(http.Dir("./assets"))
	s.Router.Handle("/assets/", standardMiddlewares(http.StripPrefix("/assets", fs)))

	return nil
}

func (s *Server) Run() error {
	addr := s.Host + ":" + strconv.Itoa(s.Port)

	s.log.Info("Server is running on " + addr)

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
	s.log.Warn("Server is shuting down...")
	return s.HttpServer.Close()
}

func (s *Server) BlockAppRoutes(appPrefix string) {
	middleware.BlockRoute(appPrefix)
}

func (s *Server) UnblockAppRoutes(appPrefix string) {
	middleware.UnblockRoute(appPrefix)
}
