// TODO: add documentation
package framework

import (
	"context"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"github.com/google/uuid"
)

type ctxKey int

const KeyRequestState ctxKey = 1

type RequestState struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

// A Handler is a type that handles a http request within our own little mini
// framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
	mw       []Middleware
}

func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	app := App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		mw:         mw,
	}

	return &app
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (app *App) Handle(method string, path string, handler Handler, mw ...Middleware) {
	// first wrap route specific middleware
	handler = wrapMiddleware(mw, handler)

	// then wrap app specific middleware
	handler = wrapMiddleware(app.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {
		requestState := RequestState{
			TraceID: uuid.New().String(),
			Now:     time.Now(),
		}

		ctx := context.WithValue(r.Context(), KeyRequestState, &requestState)

		if err := handler(ctx, w, r); err != nil {
			app.SignalShutdown()
			return
		}
	}

	app.ContextMux.Handle(method, path, h)
}

// SignalShutdown is used to gracefully shut down the app when an integrity
// issue is identified.
func (app *App) SignalShutdown() {
	app.shutdown <- syscall.SIGTERM
}
