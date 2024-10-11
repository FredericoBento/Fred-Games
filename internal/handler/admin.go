package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/FredericoBento/HandGame/internal/logger"
	"github.com/FredericoBento/HandGame/internal/services"
	"github.com/FredericoBento/HandGame/internal/services/admin_service"
	"github.com/FredericoBento/HandGame/internal/views"
	"github.com/FredericoBento/HandGame/internal/views/admin_views"
	"github.com/a-h/templ"
)

var (
	ErrGameAlreadyStarted  = errors.New("game is already running")
	ErrGameCouldNotStart   = errors.New("a server error ocurred, could not start game")
	ErrGameAlreadyStopped  = errors.New("game is already inactive")
	ErrGameCouldNotStop    = errors.New("a server error ocurred, could not stop game")
	ErrGameAlreadyActive   = errors.New("game is already active, there is no need to resume")
	ErrGameCouldNotResume  = errors.New("a server error ocurred, could not resume game")
	ErrGameNotFound        = errors.New("game not found")
	ErrGameCouldNotGetMore = errors.New("could not get more info of game")
	ErrGameIsInactive      = errors.New("game is inactive")
)

type AdminHandler struct {
	adminService *admin_service.AdminService
	userService  *services.UserService
	log          *slog.Logger
}

func NewAdminHandler(adminService *admin_service.AdminService, userService *services.UserService) *AdminHandler {
	lo, err := logger.NewHandlerLogger("AdminHandler", "", true)
	if err != nil {
		lo = slog.Default()
		lo.Error("could not create admin handler logger, using slog.default")
	}

	return &AdminHandler{
		adminService: adminService,
		userService:  userService,
		log:          lo,
	}
}

func (h *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	}
}

func (h *AdminHandler) Get(w http.ResponseWriter, r *http.Request) {
	route := strings.Split(r.URL.Path, "/")
	switch route[len(route)-1] {

	case "":
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return

	case "dashboard":
		h.GetDashboard(w, r)

	case "users":
		h.GetUsers(w, r)

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 page not found")
		return

	}
}

func (h *AdminHandler) Post(w http.ResponseWriter, r *http.Request) {
	route := strings.Split(r.URL.Path, "/")
	switch route[len(route)-1] {
	case "dashboard":
		h.GetDashboard(w, r)
	}
}

func (h *AdminHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	gameID := r.URL.Query().Get("gameid")
	if gameID != "" {
		action := r.URL.Query().Get("action")
		switch action {
		case "start":
			h.startGame(w, r, gameID)
			return

		case "stop":
			h.stopGame(w, r, gameID)
			return

		case "resume":
			h.resumeGame(w, r, gameID)
			return

		case "more":
			h.moreGame(w, r, gameID)
			return

		case "goto":
			h.gotoGame(w, r, gameID)
			return

		default:
			http.Error(w, "Action not found", http.StatusBadRequest)
			return
		}
	}

	h.View(w, r, AdminViewProps{
		title:   "Dashboard",
		content: admin_views.Dashboard(h.adminService.GameServices),
	})
}

func (h *AdminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		h.log.Error(err.Error())
		return
	}

	h.View(w, r, AdminViewProps{
		title:   "Users",
		content: admin_views.UsersPage(&users),
	})

}

func (h *AdminHandler) moreGame(w http.ResponseWriter, r *http.Request, gameID string) {
	game, ok := h.adminService.GetGame(gameID)
	if !ok {
		http.Error(w, ErrGameNotFound.Error(), http.StatusNotFound)
		return
	}

	h.View(w, r, AdminViewProps{
		content: admin_views.GameModal(game),
	})
	return
}

func (h *AdminHandler) startGame(w http.ResponseWriter, r *http.Request, gameID string) {
	game, ok := h.adminService.GetGame(gameID)
	if !ok {
		http.Error(w, ErrGameNotFound.Error(), http.StatusNotFound)
	}
	if game.GetStatus().IsActive() {
		http.Error(w, ErrGameAlreadyStarted.Error(), http.StatusBadRequest)
	} else {
		err := h.adminService.StartGame(gameID)
		if err != nil {
			h.log.Error(err.Error())
			http.Error(w, ErrGameCouldNotStart.Error(), http.StatusInternalServerError)
		}
	}

	h.View(w, r, AdminViewProps{
		title:   "Dashboard",
		content: admin_views.Dashboard(h.adminService.GameServices),
	})
	return
}

func (h *AdminHandler) stopGame(w http.ResponseWriter, r *http.Request, gameID string) {
	game, ok := h.adminService.GetGame(gameID)

	if !ok {
		h.log.Error(ErrGameNotFound.Error())
		http.Error(w, ErrGameNotFound.Error(), http.StatusNotFound)
		return
	}

	if game.GetStatus().IsInactive() {
		h.log.Error(ErrGameAlreadyStopped.Error())
		http.Error(w, ErrGameAlreadyStopped.Error(), http.StatusBadRequest)
		return
	}

	err := h.adminService.StopGame(gameID)
	if err != nil {
		h.log.Error(err.Error())
		http.Error(w, ErrGameCouldNotStop.Error(), http.StatusInternalServerError)
		return
	}

	h.View(w, r, AdminViewProps{
		title:   "Dashboard",
		content: admin_views.Dashboard(h.adminService.GameServices),
	})
	return
}

func (h *AdminHandler) resumeGame(w http.ResponseWriter, r *http.Request, gameID string) {
	game, ok := h.adminService.GetGame(gameID)

	if !ok {
		h.log.Error(ErrGameNotFound.Error())
		http.Error(w, ErrGameNotFound.Error(), http.StatusNotFound)
		return
	}

	if game.GetStatus().IsActive() {
		h.log.Error(ErrGameAlreadyActive.Error())
		http.Error(w, ErrGameAlreadyActive.Error(), http.StatusBadRequest)
		return
	}

	err := h.adminService.ResumeGame(gameID)
	if err != nil {
		h.log.Error(err.Error())
		http.Error(w, ErrGameCouldNotResume.Error(), http.StatusInternalServerError)
		return
	}

	h.View(w, r, AdminViewProps{
		title:   "Dashboard",
		content: admin_views.Dashboard(h.adminService.GameServices),
	})
	return
}

func (h *AdminHandler) gotoGame(w http.ResponseWriter, r *http.Request, gameID string) {
	game, ok := h.adminService.GetGame(gameID)
	if !ok {
		h.log.Error(ErrGameNotFound.Error())
		http.Error(w, ErrGameNotFound.Error(), http.StatusNotFound)
		return
	}

	if game.GetStatus().IsInactive() {
		h.log.Error(ErrGameIsInactive.Error())
		http.Error(w, ErrGameIsInactive.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Hx-Request", "false") // Because we want to replace the body with a full page
	Redirect(w, r, game.GetRoute()+"/home")
	return

}

type AdminViewProps struct {
	title   string
	content templ.Component
}

func (hh *AdminHandler) View(w http.ResponseWriter, r *http.Request, props AdminViewProps) {
	if IsHTMX(r) {
		props.content.Render(r.Context(), w)
	} else {
		views.Page(props.title, admin_views.Navbar(), props.content).Render(r.Context(), w)
	}
}

func (hh *AdminHandler) ReturnError(w http.ResponseWriter, r *http.Request, error string) {
	props := AdminViewProps{
		content: views.ErrorNotification(error),
	}
	hh.View(w, r, props)
}
