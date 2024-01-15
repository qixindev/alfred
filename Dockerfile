##
## bulid web
##
FROM node:14 AS build-web
WORKDIR /app
COPY ./web /app

RUN yarn install
RUN npm run generate


##
## bulid backend
##
FROM golang:1.20-alpine as build-back
WORKDIR /app
COPY . .

RUN go env -w GO111MODULE=on \
        && go env -w GOPROXY=https://goproxy.cn,direct \
        && go env -w CGO_ENABLED=0 \

RUN go build -o accounts .


##
## deploy
##
FROM alpine:latest
RUN mkdir -p /app/web
RUN mkdir -p /app/docs
WORKDIR /app
COPY --from=build-web       /app/.output ./web/.output
#COPY --from=build-web       /app/.nuxt   ./web/.nuxt
#COPY --from=build-web       /app/assets  ./web/assets

COPY --from=build-back   /app/accounts    ./
COPY --from=build-back   /app/config*.yaml ./
COPY --from=build-back   /app/backend/docs/*.md    ./backend/docs/

EXPOSE 8086
ENTRYPOINT ["/app/accounts"]
