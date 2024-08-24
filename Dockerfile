## building vuejs app
FROM node:lts AS vue-builder
WORKDIR /app

COPY frontend/. .
RUN npm install
RUN npm run build

## building go app
FROM golang:1.22 AS go-builder
WORKDIR /app
COPY backend/. ./

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

#RUN dpkg --add-architecture i386

COPY --from=go-builder /app/cs-server-manager /cs-server-manager
COPY --from=vue-builder /app/dist /dist

CMD ["/cs-server-manager"]
