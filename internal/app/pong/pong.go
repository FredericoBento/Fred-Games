package pong

import (
	"errors"
	"log/slog"

	"github.com/FredericoBento/HandGame/internal/app"
	"github.com/FredericoBento/HandGame/internal/logger"
	"github.com/FredericoBento/HandGame/internal/middleware"
	"github.com/FredericoBento/HandGame/internal/models"
	"github.com/FredericoBento/HandGame/internal/utils"
	"golang.org/x/net/websocket"
)

var (
	ErrServerRouterNotFound = errors.New("Could not find server router")
	ErrHasNotStartedYet     = errors.New("app hasnt started yet")
	ErrCouldNotCreateLogger = errors.New("could not create app logger")
	ErrInvalidCode          = errors.New("game with such code does not exist")
)

type PongApp struct {
	name        string
	routePrefix string
	status      *app.AppStatus
	server      *app.Server
	log         *slog.Logger
	games       map[string]*PongGame
}

type PongPlayer struct {
	username string
	socket   *websocket.Conn
}

type PongGame struct {
	p1   *PongPlayer
	p2   *PongPlayer
	code string
}

func NewPongApp(name, routePrefix string, server *app.Server) *PongApp {
	lo, err := logger.NewAppLogger(name, "", false)
	if err != nil {
		slog.Error(ErrCouldNotCreateLogger.Error() + " " + err.Error())
		lo = slog.Default()
	}

	return &PongApp{
		name:        name,
		routePrefix: routePrefix,
		status:      app.NewAppStatus(),
		server:      server,
		log:         lo,
		games:       make(map[string]*PongGame),
	}
}

func (a *PongApp) Start() error {
	a.log.Info("Starting Pong App...")
	err := a.setupRoutes()
	if err != nil {
		a.log.Error(" - Failed")
		a.log.Error(err.Error())
		return err
	}
	a.log.Info(" - Ok")
	a.status.SetActive()
	return nil
}

func (a *PongApp) setupRoutes() error {
	if a.server.Router == nil {
		return ErrServerRouterNotFound
	}

	appMiddlewares := middleware.StackMiddleware(
		middleware.Logger,
		middleware.SecureHeadersMiddleware,
		middleware.RequiredLogged,
	)

	a.server.Router.Handle(a.routePrefix+"/home", appMiddlewares(a.server.Handlers.PongHandler))

	// a.server.Router.Handle(a.routePrefix+"/ws", websocket.Handler(wsHandler))

	return nil
}

func (a *PongApp) Stop() error {
	a.log.Warn("Stopping Pong App...")
	a.server.BlockAppRoutes(a.routePrefix)

	a.log.Info(" - Ok")
	a.status.SetInactive()

	return nil
}

func (a *PongApp) Resume() error {
	if a.status.HasStartedOnce() == false {
		a.log.Error(ErrHasNotStartedYet.Error())
		return ErrHasNotStartedYet
	}
	a.log.Warn("Resuming Pong App...")
	a.server.UnblockAppRoutes(a.routePrefix)

	a.log.Info(" - Ok")
	a.status.SetActive()

	return nil
}

func (a *PongApp) GetName() string {
	return a.name
}

func (a *PongApp) GetRoute() string {
	return a.routePrefix
}

func (aa *PongApp) GetStatus() app.AppStatusChecker {
	return aa.status
}

func (da *PongApp) GetLogs() ([]logger.PrettyLogs, error) {
	logs, err := logger.GetAppLogs(da.name)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (a *PongApp) CreateGame(user *models.User, socket *websocket.Conn) (string, error) {
	p1 := &PongPlayer{
		username: user.Username,
		socket:   socket,
	}

	game := &PongGame{
		p1: p1,
		p2: nil,
	}

	randomCode := utils.RandomString(5)
	for a.games[randomCode] != nil {
		randomCode = utils.RandomString(5)
	}

	a.games[randomCode] = game

	return randomCode, nil
}

func (a *PongApp) JoinGame(code string, user *models.User, socket *websocket.Conn) error {
	game := a.games[code]
	if game == nil {
		return ErrInvalidCode
	}
	p2 := &PongPlayer{
		username: user.Username,
		socket:   socket,
	}

	game.p2 = p2

	return nil
}
