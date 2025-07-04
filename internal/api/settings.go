package api

import (
	_ "embed"
	"net/http"
	"strconv"
	"text/template"

	"github.com/alexdglover/sage/internal/models"
	"github.com/alexdglover/sage/internal/utils"
)

type SettingsController struct {
	SettingsRepository *models.SettingsRepository
}

//go:embed settings.html
var settingsPageTmpl string

type SettingsPageDTO struct {
	ActivePage string
	// List of settings to be displayed on the page
	LaunchBrowserOnStartup bool

	SettingsUpdated        bool
	SettingsUpdatedMessage string
}

func (sc *SettingsController) generateSettingsView(w http.ResponseWriter, req *http.Request) {
	settings, err := sc.SettingsRepository.GetSettings()
	if err != nil {
		http.Error(w, "Unable to retrieve settings", http.StatusInternalServerError)
		return
	}

	dto := SettingsPageDTO{
		ActivePage:             "settings",
		LaunchBrowserOnStartup: settings.LaunchBrowserOnStartup,
	}
	tmpl := template.Must(template.New("settingsPage").Parse(pageComponents))
	tmpl = template.Must(tmpl.Parse(settingsPageTmpl))
	_ = utils.RenderTemplateAsHTML(w, tmpl, dto)
}

func (sc *SettingsController) upsertSettings(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, "Unable to Parse Form ", http.StatusBadRequest)
		return
	}

	settings, err := sc.SettingsRepository.GetSettings()
	if err != nil {
		http.Error(w, "Unable to retrieve settings", http.StatusInternalServerError)
		return
	}

	// Extract the settings from the form
	launchBrowserOnStartupInput, err := strconv.ParseBool(req.FormValue("launchBrowserOnStartup"))
	if err != nil {
		http.Error(w, "Invalid value for launchBrowserOnStartup", http.StatusBadRequest)
		return
	}
	settings.LaunchBrowserOnStartup = launchBrowserOnStartupInput
	err = sc.SettingsRepository.Save(settings)
	if err != nil {
		http.Error(w, "Error occurred while savings settings", http.StatusInternalServerError)
		return
	}

	dto := SettingsPageDTO{
		ActivePage:             "settings",
		LaunchBrowserOnStartup: settings.LaunchBrowserOnStartup,
		SettingsUpdated:        true,
		SettingsUpdatedMessage: "Settings saved successfully!",
	}
	tmpl := template.Must(template.New("settingsPage").Parse(pageComponents))
	tmpl = template.Must(tmpl.Parse(settingsPageTmpl))
	_ = utils.RenderTemplateAsHTML(w, tmpl, dto)
}
