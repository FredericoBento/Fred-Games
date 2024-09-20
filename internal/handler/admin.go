package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/FredericoBento/HandGame/internal/app"
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
	menu        map[string]string
	userService *services.UserService
}

func NewAdminHandler(am *app.AppsManager, userService *services.UserService) *AdminHandler {
	navlinks := make(map[string]string, 0)
	navlinks["Dashboard"] = "/admin/dashboard"
	navlinks["Users"] = "/admin/users"

	return &AdminHandler{
		appManager:  am,
		menu:        navlinks,
		userService: userService,
	}
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

		default:
			http.Error(w, "Action not found", http.StatusBadRequest)
			return
		}
	}

	ah.View(w, r, adminViewProps{
		title:   "Dashboard",
		content: admin_views.Dashboard(ah.appManager.Apps),
	})
}

func (ah *AdminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := ah.userService.GetAllUsers()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	ah.View(w, r, adminViewProps{
		title:   "Users",
		content: admin_views.UsersPage(&users),
	})

}

func (ah *AdminHandler) moreApp(w http.ResponseWriter, r *http.Request, appID string) {
	for name, app := range ah.appManager.Apps {
		if name == appID {
			ah.View(w, r, adminViewProps{
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
					slog.Error(err.Error())
					http.Error(w, ErrAppCouldNotStart.Error(), http.StatusInternalServerError)
				}
			}

			ah.View(w, r, adminViewProps{
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
					slog.Error(err.Error())
					http.Error(w, ErrAppCouldNotStop.Error(), http.StatusInternalServerError)
				}
			}

			ah.View(w, r, adminViewProps{
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
					slog.Error(err.Error())
					http.Error(w, ErrAppCouldNotResume.Error(), http.StatusInternalServerError)
				}
			}

			ah.View(w, r, adminViewProps{
				title:   "Dashboard",
				content: admin_views.Dashboard(ah.appManager.Apps),
			})
			return
		}
	}

	http.Error(w, ErrAppCouldNotResume.Error(), http.StatusNotFound)
	return
}

type adminViewProps struct {
	title   string
	content templ.Component
}

func (hh *AdminHandler) View(w http.ResponseWriter, r *http.Request, props adminViewProps) {

	if IsHTMX(r) {
		props.content.Render(r.Context(), w)
		slog.Warn("here")
	} else {
		views.Page(props.title, "", hh.menu, props.content).Render(r.Context(), w)
	}
}
