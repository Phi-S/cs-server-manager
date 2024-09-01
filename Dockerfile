## building frontend
FROM node:22.4.0 AS webui-builder
WORKDIR /app

COPY frontend .

RUN npm install
RUN npm run build

## building backend
FROM golang:1.22 AS backend-builder
WORKDIR /app

RUN export GOPATH=$HOME/go
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3

COPY backend .

RUN $GOPATH/bin/swag init
RUN cp docs/swagger.json swagger-ui/swagger.json

COPY --from=webui-builder /app/dist web

RUN go mod tidy
RUN go mod verify
RUN go mod download
RUN go build -v -o cs-server-manager

## building debian with dependencies for cs server and steamcmd
FROM debian:12.6-slim
RUN set -x \
	&& apt-get update \
	&& apt-get install -y --no-install-recommends --no-install-suggests \
		ca-certificates \
		lib32z1 \
    && apt-get autoremove \
    && apt-get clean \
    && find /var/lib/apt/lists/ -type f -delete

COPY --from=backend-builder /app/cs-server-manager .
CMD ["/cs-server-manager"]
