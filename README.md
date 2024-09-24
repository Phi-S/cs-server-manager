# cs-server-manager

### :warning: This project is under active development expect breaking changes

---

<br/>

# Installation

> To configure your installation see [Environment Variables](#environment-variables)

## Docker (recommended):

> [Docker Hub](https://hub.docker.com/r/phiis/cs-server-manager)

### Run:

```
docker run -p 8080:8080 -p 27015:27015/udp -v cs-server-manager_volume:/data phiis/cs-server-manager:latest
```

### Compose:

```
services:
  cs-server-manager:
    image: phiis/cs-server-manager:latest
    volumes:
      - cs-server-manager_volume:/data
    ports:
      - 27015:27015/udp
      - 8080:8080
volumes:
  cs-server-manager_volume:
```

## Binary:

Download the binary from [releases](https://github.com/Phi-S/cs-server-manager/releases)

Run with `./cs-server-manager`

<br/>

# API

### [API documentation](api-documentation.md)

Default API path: `/api/v1`

Default Swagger UI path: `/api/swagger/index.html`

<br/>

# Web UI

![Web UI](web-ui-start-server.gif)

<br/>

# Development

> This project is designed to only run on Linux (for now).<br/>
> For development with Windows, [WSL2](https://learn.microsoft.com/windows/wsl/install) can be used

> Before every commit please run `make ready`

## Prerequisites:

- Go version [1.22 or higher](https://go.dev/doc/install)
- NodeJs [v20.17.0](https://nodejs.org/en)
- Steamcmd requires `lib32gcc-s1` to run.
  `sudo apt install lib32gcc-s1`

## Setup local development environment

### Git:

```
git clone https://github.com/Phi-S/cs-server-manager.git
cd cs-server-manager

go get -C backend/
npm install --prefix frontend/
```

### Run

> Status endpoint: [http://localhost:8080/api/v1/status](http://localhost:8080/api/v1/status) <br/>

> WebUI: [http://localhost:8090](http://localhost:8090)

#### backend only:

```
make backend
```

#### frontend only:

```
make frontend
```

<br/>

# Build

### Binary:

```
make build
```

### Docker:

```
docker build -t cs-server-manager .
```

### Generate swagger docs with [swaggo/swag](https://github.com/swaggo/swag):

> If the command `swag` was not found, add the following lines to your `.bashrc` and restart your terminal
>
> ```
> export GOPATH="$HOME/go"
> export PATH=$PATH:$GOPATH/bin
> ```
>
> [source](https://stackoverflow.com/a/72166253/12487257)

```
go install github.com/swaggo/swag/cmd/swag@v1.16.3
swag init --dir backend -o . -ot json
```

### Generate api-documentation.md with [widdershins](https://github.com/Mermade/widdershins):

> npm install widdershins@v4.0.0

```
make doc
```

<br/>

# Environment variables

> Those environment variables can also be set via an `.env` file.
> <br/>
> It should be located in the same folder as the `cs-server-manager` binary or in the `backend` folder for development

| KEY            | TYPE   | DEFAULT                  | DESCRIPTION                                                                                                                           |
| -------------- | ------ | ------------------------ | ------------------------------------------------------------------------------------------------------------------------------------- |
| IP             | string | current public IP        | This IP is returned with the status endpoint to generate the connection url.<br/>If no specified, the current public ip will be used. |
| HTTP_PORT      | string | 8080                     | The API / WebSocket port                                                                                                              |
| CS_PORT        | string | 27015                    | CS 2 server port. This port will be reported with the status endpoint to generate the connection URL                                  |
| DATA_DIR       | string | {working directory}/data | The base data directory for all CS server files                                                                                       |
| LOG_DIR        | string | {DATA_DIR}/logs          | Location of the CS server logs                                                                                                        |
| SERVER_DIR     | string | {DATA_DIR}/server        | The CS 2 server directory.<br/>After installation this folder will be around 30 GB is size                                            |
| STEAMCMD_DIR   | string | {DATA_DIR}/steamcmd      | The steamcmd directory                                                                                                                |
| ENABLE_WEB_UI  | bool   | true                     | If set to true, the backend will host the WEB UI                                                                                      |
| ENABLE_SWAGGER | bool   | true                     | If set to true, the backend will host the swagger UI                                                                                  |

<br/>

# Plugins

## Default plugins list

By default, following plugins are available via one click install.

- [Cs2PracticeMode](https://github.com/Phi-S/cs2-practice-mode).

> If you want your plugin to be in this default list, please add it to the [default_plugins.go](/backend/plugins/default_plugins.go) file.

## Custom install actions

Some plugins require additional steps after installation for them to work correctly.
<br/>
For example: Metamod requires a line to be added in the `gameinfo.gi` for it to work.
<br/>
For Metamod, this action will be automatically executed if you try to install `metamod_source`(name must match) as dependency for your custom plugin.

If your plugin requires such an action add it to the [custom_install_actions.go](/backend/plugins/custom_install_actions.go) file.

## Custom `plugins.json`

It is also possible to create your own plugins list.

To create your own list create a file called `plugins.json` in the `{DATA_DIR}`(By Default `/data`) directory.

This will overwrite the [default plugin list](#default-plugins-list), so only the plugins defined in the new `plugin.json` will be available for installation.

> At the moment only `.tar.gz` and `.zip` files are supported

The `install_dir` field is the directory in which the downloaded content gets extracted to.
<br/>
`/` means `/{SERVER_DIR}/game/csgo`.

### Example `plugins.json`

```
[
  {
    "name": "CounterStrikeSharp",
    "description": "CounterStrikeSharp allows you to write server plugins in C# for Counter-Strike 2/Source2/CS2",
    "url": "https://github.com/roflmuffin/CounterStrikeSharp",
    "install_dir": "/",
    "versions": [
      {
        "name": "v264",
        "download_url": "https://github.com/roflmuffin/CounterStrikeSharp/releases/download/v264/counterstrikesharp-with-runtime-build-264-linux-8f59fd5.zip",
        "dependencies": [
          {
            "name": "metamod_source",
            "install_dir": "/",
            "version": "2.0.0-git1313",
            "download_url": "https://mms.alliedmods.net/mmsdrop/2.0/mmsource-2.0.0-git1313-linux.tar.gz",
            "dependencies": null
          }
        ]
      }
    ]
  },
  {
    "name": "Cs2PracticeMode",
    "description": "Practice mode for cs2 server based on CounterStrikeSharp",
    "url": "https://github.com/Phi-S/cs2-practice-mode",
    "install_dir": "/addons/counterstrikesharp/plugins/",
    "versions": [
      {
        "name": "0.0.16",
        "download_url": "https://github.com/Phi-S/cs2-practice-mode/releases/download/0.0.16/cs2-practice-mode-0.0.16.tar.gz",
        "dependencies": [
          {
            "name": "CounterStrikeSharp",
            "install_dir": "/",
            "version": "v264",
            "download_url": "https://github.com/roflmuffin/CounterStrikeSharp/releases/download/v264/counterstrikesharp-with-runtime-build-264-linux-8f59fd5.zip",
            "dependencies": [
              {
                "name": "metamod_source",
                "install_dir": "/",
                "version": "2.0.0-git1313",
                "download_url": "https://mms.alliedmods.net/mmsdrop/2.0/mmsource-2.0.0-git1313-linux.tar.gz",
                "dependencies": null
              }
            ]
          }
        ]
      }
    ]
  }
]

```

# Editor

The endpoints `/files` or `/configs` in the WebUI can be used to edit local files (for example the server.cfg file).

By default following files can be edited:

- Files witch end with `.cfg` in `/game/csgo/cfg`
- Files witch end with `.json` in `/game/csgo/addons/counterstrikesharp/configs`
- Files witch end with `.json` or `.cfg` in `/game/csgo/addons/counterstrikesharp/plugins`

Custom entries can be added by creating the file `editor-files.json` in the [DATA_DIR](#environment-variables)

> If a `editor-files.json` file exist, the defaults will be overwritten

> Files will only show if they exist

Example content of `editor-files.json`:

```
[
    {
        "path": "/path/to/file.txt",
        "extensions": null
    },
    {
        "path": "/game/csgo/cfg",
        "extensions": [
            ".cfg"
        ]
    },
    {
        "path": "/game/csgo/addons/counterstrikesharp/configs",
        "extensions": [
            ".json"
        ]
    },
    {
        "path": "/game/csgo/addons/counterstrikesharp/plugins",
        "extensions": [
            ".json",
            ".cfg"
        ]
    }
]
```
