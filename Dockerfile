FROM golang:1.17 AS build-env

COPY . ${GOPATH}/src/github.com/Tomansru/keeneteus
WORKDIR ${GOPATH}/src/github.com/Tomansru/keeneteus

RUN go build -ldflags "-s -w" -trimpath -o keeneteus

FROM ubuntu:rolling
LABEL maintainer="stas@tomans.ru"

EXPOSE 2112

CMD ["/app/keeneteus/keeneteus"]

COPY --from=build-env /go/src/github.com/Tomansru/keeneteus/keeneteus /app/keeneteus/keeneteus