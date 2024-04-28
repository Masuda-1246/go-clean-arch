# develop
FROM golang:1.22 AS develop

WORKDIR /go/src/github.com/Masuda-1246/go-clean-arch

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/cosmtrek/air@latest

CMD ["air", "-c", ".air.toml"]

## Build
FROM golang:1.22 AS build

WORKDIR /go/src/

COPY app .

RUN go mod download

COPY db db

COPY .env .env

RUN [ ! -f .env ] || export $(grep -v '^#' .env | xargs)

############################################################################### Migrate
# RUN go install github.com/rubenv/sql-migrate/...@latest

# RUN sql-migrate up -env="production" -config=./db/postgres/dbconfig.prod.yml
################################################################################

RUN CGO_ENABLED=0 go build -o /server ./cmd

## Deploy
FROM golang:1.22 AS deploy

WORKDIR /

COPY --from=build /server /server

EXPOSE 8080

CMD ["/server"]
