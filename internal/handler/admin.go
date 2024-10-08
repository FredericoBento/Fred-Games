package handler

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/FredericoBento/HandGame/internal/app"
	"github.com/FredericoBento/HandGame/internal/logger"
	"github.com/FredericoBento/HandGame/internal/models"
	"github.com/FredericoBento/HandGame/internal/services"
	"github.com/FredericoBento/HandGame/internal/views"
	"github.com/FredericoBento/HandGame/internal/views/admin_views"
	"github.com/a-h/templ"
)

var (
	ErrAppAlreadyStarted  = errors.New("app is already running")
	ErrAppCouldNotStart   = errors.New("a server error ocurred, could not start app")
	ErrAppAlreadyStopped  = errors.New("app is already inactive")
	ErrAppCouldNotStop    = errors.New("a server error ocurred, could not stop app")
	ErrAppAlreadyActive   = errors.New("app is already active, there is no need to resume")
	ErrAppCouldNotResume  = errors.New("a server error ocurred, could not resume app")
	ErrAppNotFound        = errors.New("app not found")
	ErrAppCouldNotGetMore = errors.New("could not get more info of app")
)

type AdminHandler struct {
	appManager  *app.AppsManager
	navbar      models.NavBarStructure
	userService *services.UserService
	log         *slog.Logger
}

func NewAdminHandler(am *app.AppsManager, userService *services.UserService) *AdminHandler {
	if am == nil {
		log.Fatal("AdminHandler: app manager not provided")
	}
	if userService == nil {
		log.Fatal("AdminHandler: user service not provided")
	}

	lo, err := logger.NewHandlerLogger("AdminHandler", "", false)
	if err != nil {
		lo = slog.Default()
		lo.Error("could not create admin handler logger, using slog.default")
	}

	h := &AdminHandler{
		appManager:  am,
		userService: userService,
		log:         lo,
	}

	h.setupNavbar()

	return h
}

func (h *AdminHandler) setupNavbar() {
	startBtns := []models.Button{
		{
			ButtonName:   "Games",
			Url:          "/home",
			NotHxRequest: true,
		},
		{
			ButtonName: "Dashboard",
			Url:        "/admin/dashboard",
		},
		{
			ButtonName: "Users",
			Url:        "/admin/users",
		},
	}

	endBtns := []models.Button{
		{
			ButtonName: "Account",
			Childs: []models.Button{
				{
					ButtonName: "Logout",
					Url:        "/logout",
				},
			},
		},
	}

	navbar := models.NavBarStructure{
		StartButtons: startBtns,
		EndButtons:   endBtns,
	}

	h.navbar = navbar
}

func (ah *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ah.Get(w, r)
	}
}

func (ah *AdminHandler) Get(w http.ResponseWriter, r *http.Request) {
	route := strings.Split(r.URL.Path, "/")
	switch route[len(route)-1] {

	case "":
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return

	case "dashboard":
		ah.GetDashboard(w, r)

	case "users":
		ah.GetUsers(w, r)

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 page not found")
		return

	}
}

func (ah *AdminHandler) Post(w http.ResponseWriter, r *http.Request) {
	route := strings.Split(r.URL.Path, "/")
	switch route[len(route)-1] {
	case "dashboard":
		ah.GetDashboard(w, r)
	}
}

func (ah *AdminHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	appID := r.URL.Query().Get("appid")
	if appID != "" {
		action := r.URL.Query().Get("action")
		switch action {
		case "start":
			ah.startApp(w, r, appID)
			return

		case "stop":
			ah.stopApp(w, r, appID)
			return

		case "resume":
			ah.resumeApp(w, r, appID)
			return

		case "more":
			ah.moreApp(w, r, appID)
			return

		case "goto":
			ah.gotoApp(w, r, appID)
			return

		default:
			http.Error(w, "Action not found", http.StatusBadRequest)
			return
		}
	}

	ah.View(w, r, AdminViewProps{
		title:   "Dashboard",
		content: admin_views.Dashboard(ah.appManager.Apps),
	})
}

func (ah *AdminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := ah.userService.GetAllUsers()
	if err != nil {
		ah.log.Error(err.Error())
		return
	}

	ah.View(w, r, AdminViewProps{
		title:   "Users",
		content: admin_views.UsersPage(&users),
	})

}

func (ah *AdminHandler) moreApp(w http.ResponseWriter, r *http.Request, appID string) {
	for name, app := range ah.appManager.Apps {
		if name == appID {
			ah.View(w, r, AdminViewProps{
				title:   "App Modal",
				content: admin_views.AppModal(app),
			})
			return
		}
	}

	http.Error(w, ErrAppNotFound.Error(), http.StatusNotFound)
	return
}

func (ah *AdminHandler) startApp(w http.ResponseWriter, r *http.Request, appID string) {
	for name, app := range ah.appManager.Apps {
		if name == appID {
			if app.GetStatus().IsActive() {
				http.Error(w, ErrAppAlreadyStarted.Error(), http.StatusBadRequest)
			} else {
				err := app.Start()
				if err != nil {
					ah.log.Error(err.Error())
					http.Error(w, ErrAppCouldNotStart.Error(), http.StatusInternalServerError)
				}
			}

			ah.View(w, r, AdminViewProps{
				title:   "Dashboard",
				content: admin_views.Dashboard(ah.appManager.Apps),
			})
			return
		}
	}

	http.Error(w, ErrAppCouldNotStart.Error(), http.StatusNotFound)
	return
}

func (ah *AdminHandler) stopApp(w http.ResponseWriter, r *http.Request, appID string) {
	for name, app := range ah.appManager.Apps {
		if name == appID {
			if app.GetStatus().IsInactive() {
				http.Error(w, ErrAppAlreadyStopped.Error(), http.StatusBadRequest)
			} else {
				err := app.Stop()
				if err != nil {
					ah.log.Error(err.Error())
					http.Error(w, ErrAppCouldNotStop.Error(), http.StatusInternalServerError)
				}
			}

			ah.View(w, r, AdminViewProps{
				title:   "Dashboard",
				content: admin_views.Dashboard(ah.appManager.Apps),
			})
			return
		}
	}

	http.Error(w, ErrAppCouldNotStop.Error(), http.StatusNotFound)
	return
}

func (ah *AdminHandler) resumeApp(w http.ResponseWriter, r *http.Request, appID string) {
	for name, app := range ah.appManager.Apps {
		if name == appID {
			if app.GetStatus().IsActive() {
				http.Error(w, ErrAppAlreadyActive.Error(), http.StatusBadRequest)
			} else {
				err := app.Resume()
				if err != nil {
					ah.log.Error(err.Error())
					http.Error(w, ErrAppCouldNotResume.Error(), http.StatusInternalServerError)
				}
			}

			ah.View(w, r, AdminViewProps{
				title:   "Dashboard",
				content: admin_views.Dashboard(ah.appManager.Apps),
			})
			return
		}
	}

	http.Error(w, ErrAppCouldNotResume.Error(), http.StatusNotFound)
	return
}

func (ah *AdminHandler) gotoApp(w http.ResponseWriter, r *http.Request, appID string) {
	for name, app := range ah.appManager.Apps {
		if name == appID {
			if app.GetStatus().IsInactive() {
				w.WriteHeader(http.StatusBadRequest)
				ah.ReturnError(w, r, "App is not active")
				return
			}

			w.Header().Add("Hx-Request", "false") // Because we want to replace the body with a full page
			Redirect(w, r, app.GetRoute()+"/home")
			// http.Redirect(w, r, app.GetRoute()+"/home", http.StatusSeeOther)
			return
		}
	}
}

type AdminViewProps struct {
	title   string
	content templ.Component
}

func (hh *AdminHandler) View(w http.ResponseWriter, r *http.Request, props AdminViewProps) {
	if IsHTMX(r) {
		props.content.Render(r.Context(), w)
	} else {
		views.Page(props.title, hh.navbar, props.content).Render(r.Context(), w)
	}
}

func (hh *AdminHandler) ReturnError(w http.ResponseWriter, r *http.Request, error string) {
	props := AdminViewProps{
		content: views.ErrorNotification(error),
	}
	hh.View(w, r, props)
}
