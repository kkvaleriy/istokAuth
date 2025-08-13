FROM golang:1.24.1-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN CGO_ENABLED=0 go build -o app ./cmd/main.go

FROM alpine:3.22.0 

WORKDIR /app

COPY --from=build /app/app .

CMD ["./app"]