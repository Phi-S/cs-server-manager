package middleware

import (
	"context"
	"cs-server-controller/config"
	"cs-server-controller/logwrt"
	"cs-server-controller/server"
	"cs-server-controller/steamcmd"
	"errors"
	"net/http"
	"sync"
)

func ContextValues(
	h http.Handler,
	cfg config.Config,
	serverSteamcmdLock *sync.Mutex,
	serverInstance *server.ServerInstance,
	steamcmdInstance *steamcmd.SteamcmdInstance,
	userLogWriter *logwrt.LogWriter,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, ConfigKey, cfg)
		ctx = context.WithValue(ctx, ServerSteamcmdLockKey, serverSteamcmdLock)
		ctx = context.WithValue(ctx, ServerInstanceKey, serverInstance)
		ctx = context.WithValue(ctx, SteamCmdInstanceKey, steamcmdInstance)
		ctx = context.WithValue(ctx, UserLogWriterKey, userLogWriter)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

type traceIdKeyType uint

const TraceIdKey traceIdKeyType = 0

type configKeyType uint

const ConfigKey configKeyType = 0

type serverSteamcmdLockKeyType uint

const ServerSteamcmdLockKey serverSteamcmdLockKeyType = 0

type serverInstanceKeyType uint

const ServerInstanceKey serverInstanceKeyType = 0

type steamCmdInstanceKeyType uint

const SteamCmdInstanceKey steamCmdInstanceKeyType = 0

type userLogWriterKeyType uint

const UserLogWriterKey userLogWriterKeyType = 0

func GetSteamcmdAndServerInstance(ctx context.Context) (*sync.Mutex, *server.ServerInstance, *steamcmd.SteamcmdInstance, error) {
	lock, ok := ctx.Value(ServerSteamcmdLockKey).(*sync.Mutex)
	if lock == nil || !ok {
		return nil, nil, nil, errors.New("failed to get server/steamcmd lock from context")
	}
	server := ctx.Value(ServerInstanceKey).(*server.ServerInstance)
	if server == nil || !ok {
		return nil, nil, nil, errors.New("failed to get server instance from context")
	}
	steamcmd := ctx.Value(SteamCmdInstanceKey).(*steamcmd.SteamcmdInstance)
	if steamcmd == nil || !ok {
		return nil, nil, nil, errors.New("failed to get steamcmd instance from context")
	}

	return lock, server, steamcmd, nil
}
