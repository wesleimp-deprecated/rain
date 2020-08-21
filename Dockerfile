FROM golang:1.14-alpine AS base
WORKDIR /app

FROM base as builder
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o rain main.go

FROM scratch as final
WORKDIR /app
COPY --from=builder /app/rain /app
ENTRYPOINT ["./rain"]