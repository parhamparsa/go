ARG ALPINE_VERSION="3.20"
ARG GOLANG_VERSION="1.23"

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} as base
RUN apk add --update --no-cache
WORKDIR /application
COPY go.mod go.sum ./
RUN go mod download

FROM base as builder
RUN echo "Building builder ..."
COPY . .
ARG GOARCH="amd64"
RUN GOOS=linux GOARCH=$GOARCH CGO_ENABLED=0
RUN go build -tags musl -o app cmd/api/main.go
RUN echo "application folder content..."
RUN ls -l /application

FROM base as dev
RUN apk add --no-cache \
    librdkafka-dev=~2.4.0-r0 \
    pkgconf=~2.2.0-r0
RUN --mount=type=cache,target=/root/.cache/go-build \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install github.com/boumenot/gocover-cobertura@latest && \
    go install github.com/jstemmer/go-junit-report@latest && \
    go install go.uber.org/mock/mockgen@latest

FROM alpine:${ALPINE_VERSION} as main
WORKDIR /application
COPY --from=builder /application/config/yaml config/yaml
COPY --from=builder /application/app .
ENTRYPOINT ["/application/app"]
