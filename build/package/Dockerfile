FROM golang AS build

WORKDIR /example/src

COPY . .

RUN go mod download

ENTRYPOINT ["go","run", "cmd/app/main.go", "-d", "-m"]