package ctxex

import (
    "context"
    "cs-server-controller/config"
    "cs-server-controller/logwrt"
    "cs-server-controller/server"
    "cs-server-controller/steamcmd"
    "errors"
    "fmt"
    "sync"
)

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

type startParametersJsonFileKeyType uint

const StartParametersJsonFileKey startParametersJsonFileKeyType = 0

func Get[T any](ctx context.Context, key any) (T, error) {
    value, ok := ctx.Value(key).(T)
    if !ok {
        return *new(T), fmt.Errorf("failed to get %q value from context", key)
    }

    return value, nil
}

func GetConfig(ctx context.Context) (config.Config, error) {
    cfg, ok := ctx.Value(ConfigKey).(config.Config)
    if !ok {
        return config.Config{}, errors.New("failed to get config from context")
    }

    return cfg, nil
}

func GetSteamcmdAndServerInstance(ctx context.Context) (*sync.Mutex, *server.Instance, *steamcmd.Instance, error) {
    lock, ok := ctx.Value(ServerSteamcmdLockKey).(*sync.Mutex)
    if lock == nil || !ok {
        return nil, nil, nil, errors.New("failed to get server/steamcmd lock from context")
    }
    server, ok := ctx.Value(ServerInstanceKey).(*server.Instance)
    if server == nil || !ok {
        return nil, nil, nil, errors.New("failed to get server instance from context")
    }
    steamcmd, ok := ctx.Value(SteamCmdInstanceKey).(*steamcmd.Instance)
    if steamcmd == nil || !ok {
        return nil, nil, nil, errors.New("failed to get steamcmd instance from context")
    }

    return lock, server, steamcmd, nil
}

func GetUserLogWriter(ctx context.Context) (*logwrt.LogWriter, error) {
    logWriter, ok := ctx.Value(UserLogWriterKey).(*logwrt.LogWriter)
    if logWriter == nil || !ok {
        return nil, errors.New("failed to get user log writer from context")
    }

    return logWriter, nil
}
