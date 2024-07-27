package context_values

import (
	"context"
	"cs-server-controller/server"
	"cs-server-controller/steamcmd"
	"errors"
	"sync"
)

const ServerSteamcmdLockKey = "server-steamcmd-lock"
const ServerInstanceKey = "server-instance"
const SteamCmdInstanceKey = "steamcmd-instance"

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
