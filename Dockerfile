##
## bulid web
##
FROM node:14 AS build-web
WORKDIR /app
COPY ./web /app


RUN yarn install
RUN npm run build:prod


##
## bulid backend
##
FROM golang:1.20-alpine as build-back
WORKDIR /app
COPY . .

RUN go env -w GO111MODULE=on \
        && go env -w GOPROXY=https://goproxy.cn,direct \
        && go env -w CGO_ENABLED=0

RUN go build -o accounts .


##
## deploy
##
FROM alpine:latest
RUN mkdir -p /app/web
WORKDIR /app
COPY --from=build-web       /app/.output ./web/.output
COPY --from=build-web       /app/.nuxt   ./web/.nuxt
COPY --from=build-web       /app/assets  ./web/assets

COPY --from=build-back   /app/accounts    ./

EXPOSE 8086
ENV dsn = "host=143.64.18.19 port=5432 dbname=accounts user=qixin password=at9z9?gntsLPv/_Jk/,pyIrX"
ENTRYPOINT ["/app/accounts"]

