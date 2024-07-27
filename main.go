package main

import (
	"context"
	"cs-server-controller/context_values"
	"cs-server-controller/handlers"
	"cs-server-controller/server"
	"cs-server-controller/steamcmd"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	configureLogger()

	var dataFolder = "/home/desk/programming/code/go/data/cs-server-controller"
	var steamcmdDir = filepath.Join(dataFolder, "steamcmd")
	var serverDir = filepath.Join(dataFolder, "server")

	// this lock is used to prevent collision between the server and steamcmd instance
	// Fox example the lock is used to prevent the server from being started while being updated.
	// This can occur if two http request are coming in at the same time
	ServerSteamcmdLock := sync.Mutex{}

	steamcmdInstance, err := steamcmd.NewInstance(steamcmdDir, serverDir, true)
	if err != nil {
		panic("failed to create new steamcmd instance. " + err.Error())
	}
	defer steamcmdInstance.Cancel()

	serverInstance, err := server.NewInstance(serverDir, "27015", true)
	if err != nil {
		panic("failed to create new server instance. " + err.Error())
	}
	defer serverInstance.Stop()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/v1", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), context_values.ServerSteamcmdLockKey, &ServerSteamcmdLock)
				ctx = context.WithValue(ctx, context_values.SteamCmdInstanceKey, steamcmdInstance)
				ctx = context.WithValue(ctx, context_values.ServerInstanceKey, serverInstance)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		})

		r.Post("/status", handlers.StatusHandler)

		r.Post("/start", handlers.StartHandler)
		r.Post("/stop", handlers.StopHandler)
		r.Post("/send-command", handlers.SendCommandHandler)

		r.Post("/update", handlers.UpdateHandler)
		r.Post("/cancel-update", handlers.CancelUpdateHandler)
	})

	address := ":8080"
	slog.Info("Starting http server at \"" + address + "\"")
	slog.Error("Failed to start http server: " + http.ListenAndServe(address, router).Error())
}

func configureLogger() {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(time.Now().UTC().Format(time.RFC3339Nano))
			}
			return a
		},
	})

	logger := slog.New(logHandler)
	slog.SetDefault(logger)
}
