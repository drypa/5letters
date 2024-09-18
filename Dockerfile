FROM golang:alpine3.20 AS build
COPY . /src
WORKDIR /src
RUN CGO_ENABLED=0 go build -o bot

FROM alpine:3.20.3
WORKDIR /app
COPY --from=build /src/bot /app/
ENTRYPOINT ./bot
