FROM postgres:13 AS develop

FROM golang:1.20 AS migrate

RUN go install github.com/rubenv/sql-migrate/...@latest
