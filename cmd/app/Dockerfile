FROM golang:1.25.0-alpine AS build

WORKDIR /build

COPY go.mod go.sum /cmd/config/config.yml ./

RUN go mod download

COPY . .

RUN go build -o app cmd/app/main.go

FROM alpine AS runner

WORKDIR /app

COPY --from=build /build/app ./app

COPY --from=build /build/cmd/config/config.yml ./config.yml

CMD ["/app/app", "--config=/app/config.yml"]
