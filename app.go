package main

import (
	"context"
	"strings"

	"github.com/NotKatsu/GoSearch/backend/json"

	"github.com/NotKatsu/GoSearch/backend/machine"

	"github.com/NotKatsu/GoSearch/backend"

	"github.com/NotKatsu/GoSearch/backend/dialog"

	"github.com/pterm/pterm"

	"github.com/NotKatsu/GoSearch/backend/search"
	"github.com/NotKatsu/GoSearch/database"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/NotKatsu/GoSearch/backend/keystroke"
)

type App struct {
	ctx context.Context
}

func GoSearch() *App {
	return &App{}
}

var (
	currentPage = "Search"
)

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	if database.SetupDatabase() == true {
		go keystroke.Listener(a.ctx)

		if json.SystemCached() == true {
			a.SetPage("Search")
		} else {
			a.SetPage("Welcome")
		}
	} else {
		runtime.Quit(a.ctx)
	}
}

func (a *App) HandleButtonClickEvent(application any) {
	runtime.Hide(a.ctx)
	applicationMap, successfulAssertion := application.(map[string]interface{})

	if successfulAssertion == true {
		applicationName := applicationMap["Name"].(string)
		applicationLocation := applicationMap["Location"].(string)

		if machine.OpenExecutable(applicationLocation) == false {
			errorMessage := "Failed to open " + applicationName
			dialog.ErrorDialog(errorMessage)
		}

	} else {
		pterm.Fatal.WithFatal(true).Println("Something went wrong while trying to complete a assertion.")
	}
}

func (a *App) ToggleFavorite(name string, location string, favorite bool) []backend.FileReturnStruct {
	database.UpdateFavorite(name, location, favorite)

	return search.GetRecommended()
}

func (a *App) CacheSystem() {
	if machine.CacheSystem() == true {
		json.UpdateCachedSetting(true)
		a.SetPage("Search")
	}
}

func (a *App) ClearCache() bool {
	if database.ClearDatabaseCache() == true {
		json.UpdateCachedSetting(false)
		return true
	} else {
		return false
	}
}

func (a *App) GetCurrentPage() string {
	return currentPage
}

func (a *App) SetPage(page string) {
	currentPage = page

	runtime.WindowReloadApp(a.ctx)
	runtime.WindowShow(a.ctx)
	keystroke.OverWriteState(true)
}

func (a *App) CloseApp() {
	runtime.Quit(a.ctx)
}

func (a *App) Search(query string) []backend.FileReturnStruct {
	var arrayWithEmptyStruct []backend.FileReturnStruct
	emptyStruct := backend.FileReturnStruct{}

	arrayWithEmptyStruct = append(arrayWithEmptyStruct, emptyStruct)

	if query == "" {
		return search.GetRecommended()

	} else if strings.HasPrefix(strings.ToLower(query), "/") {
		if strings.ToLower(query) == "/settings" {
			currentPage = "Settings"
			runtime.WindowReload(a.ctx)
		}
	} else {
		return database.RetrieveCachedResultsByQuery(query)
	}

	return arrayWithEmptyStruct
}
