# cs-server-manager

---

# [API documentation](api-documentation.md)

# Commands

> All commands should be run from the root of this repo

## Build

```
npm install --prefix frontend/
npm run build --prefix frontend/

cp frontend/dist backend/web

swag init --dir backend/ --output backend/docs

go mod tidy -C backend/
go mod verify -C backend/
go mod download -C backend/
go build -C backend/ -v -o cs-server-manager
```

## Docker

```
docker build -t cs-server-manager .
docker run -it --rm --name cs-server-manager --mount type=bind,source=/cs-server-manager,destination=/data -p 8080:8080 -p 27015:27015 cs-server-manager
```

## Generate swagger docs with [swaggo/swag](https://github.com/swaggo/swag)

```
go install github.com/swaggo/swag/cmd/swag@v1.16.3
swag init --dir backend/ --output backend/docs
```

## Generate api-documentation.md with [widdershins](https://github.com/Mermade/widdershins)

```
npm install widdershins@v4.0.0
npx widdershins --expandBody true --language_tabs 'http:HTTP' -o api-documentation.md backend/docs/swagger.json
```

# Environment variables

| KEY            | TYPE   | DEFAULT           | DESCRIPTION                                                                                                                    |
|----------------|--------|-------------------|--------------------------------------------------------------------------------------------------------------------------------|
| IP             | string | current public IP | The IP that gets reported with the status endpoint to generate the connection url.<br/>If no specified, die public ip is used. |
| HTTP_PORT      | string | 8080              | The port on witch the API server is bound to                                                                                   |
| CS_PORT        | string | 27015             | The Port that the CS server will use. This port will be reported with the status endpoint to generate the connection URL       |
| DATA_DIR       | string | /data             | The base data directory for all CS server files                                                                                |
| LOG_DIR        | string | DATA_DIR/logs     | The CS server logs will be stored in this folder                                                                               |
| SERVER_DIR     | string | DATA_DIR/server   | The directory in witch the CS server will be installed. The CS server will be around 30 GB                                     |
| STEAMCMD_DIR   | string | DATA_DIR/steamcmd | The directory in witch steamcmd will be installed                                                                              |
| ENABLE_WEB_UI  | bool   | true              | If set to true, the API will host the WEB UI                                                                                   |
| ENABLE_SWAGGER | bool   | true              | If set to true, the API will host the swagger UI                                                                               |

