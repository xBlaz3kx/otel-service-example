FROM golang:1.21 as build

WORKDIR /app/src

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Enable caching for go modules and build cache
ENV GOCACHE=/root/.cache/go-build
ENV GOMODCACHE=/root/.cache/go-build

# Leverage build cache by copying go.mod and go.sum first
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build go mod download
RUN go mod verify

# Copy the rest of the source code
COPY . .

RUN --mount=type=cache,target="/root/.cache/go-build" go build -o ../service ./cmd/app/main.go

# Debug stage using Delve
FROM build as debug

COPY --from=build /app/src /service/src

WORKDIR /service/src

# Build Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

RUN go build -gcflags="all=-N -l" -o ./service ./cmd/app/main.go

CMD ["/go/bin/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./service"]

FROM alpine as app

# Add curl for health checks
RUN apk --update --no-cache add curl

SHELL ["/bin/sh", "-l", "-c"]

ARG user=app
ARG group=docker

RUN addgroup $group
RUN adduser \
    --disabled-password \
    --gecos "" \
    --ingroup $group \
    $user

COPY --chown=$user:$group --from=build /app/service /usr/local/bin/service

RUN chmod +x /usr/local/bin/service

USER $user
WORKDIR /usr/local/service

# Todo add health check

# Todo expose port

CMD ["/usr/local/bin/service"]