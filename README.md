# cs-server-manager

---

# [API documentation](api-documentation.md)

> All API endpoint should be prefixed with /api/v1/.
>
> Example:
>
> ```
> http://localhost:8080/api/v1/logs/100
> ```

# Commands

> All commands should be run from the root of this repo

## Build

```
npm install --prefix frontend/
npm run build --prefix frontend/

cp -R frontend/dist backend/web

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
npx widdershins -v --code --summary --expandBody --omitHeader -o api-documentation.md backend/docs/swagger.json
```

# Environment variables

When starting the application the following environment variables can be set.

> The environment variables can also be set via an `.env` file that should be located in the same folder as the `cs-server-manager` binary

| KEY            | TYPE   | DEFAULT           | DESCRIPTION                                                                                                                    |
| -------------- | ------ | ----------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| IP             | string | current public IP | The IP that gets reported with the status endpoint to generate the connection url.<br/>If no specified, die public ip is used. |
| HTTP_PORT      | string | 8080              | The port on witch the API server is bound to                                                                                   |
| CS_PORT        | string | 27015             | The Port that the CS server will use. This port will be reported with the status endpoint to generate the connection URL       |
| DATA_DIR       | string | /data             | The base data directory for all CS server files                                                                                |
| LOG_DIR        | string | DATA_DIR/logs     | The CS server logs will be stored in this folder                                                                               |
| SERVER_DIR     | string | DATA_DIR/server   | The directory in witch the CS server will be installed. The CS server will be around 30 GB                                     |
| STEAMCMD_DIR   | string | DATA_DIR/steamcmd | The directory in witch steamcmd will be installed                                                                              |
| ENABLE_WEB_UI  | bool   | true              | If set to true, the API will host the WEB UI                                                                                   |
| ENABLE_SWAGGER | bool   | true              | If set to true, the API will host the swagger UI                                                                               |
