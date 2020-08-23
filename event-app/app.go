package event_app

import (
	"encoding/json"
	"io"
	"notify"
	"sync"
)

type Event struct {
	Name string
	Data string
}

type App struct {
	mu sync.RWMutex
	events map[string]func(client *notify.Client, d string)
	*notify.App
}

func (app *App) On(name string, h func(*notify.Client, string)) {
	if h == nil || name == ""{
		return
	}

	app.mu.Lock()
	defer app.mu.Unlock()
	app.events[name] = h
}

func (app *App) Off(name string) {
	app.mu.Lock()
	defer app.mu.Unlock()
	delete(app.events, name)
}

func (app *App) RawOn() *notify.EventsHandler {
	return app.App.On()
}

func (app *App) Run() {
	app.App.On().Data(func(client *notify.Client, reader io.Reader) error {
		var event Event
		err := json.NewDecoder(reader).Decode(&event)
		if err != nil {
			return err
		}
		app.mu.RLock()
		defer app.mu.RUnlock()
		if handler, ok := app.events[event.Name]; ok {
			handler(client, event.Data)
		}
		return nil
	})
	app.App.Run()
}

func New(config notify.AppConfig) *App {
	return &App{
		events: map[string]func(client *notify.Client, d string){},
		App:    notify.NewApp(config),
	}
}