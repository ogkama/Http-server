FROM golang:1.25.0-alpine AS build

WORKDIR /build

COPY go.mod go.sum /cmd/config/config.yml ./

RUN go mod download

COPY . .

RUN go build -o app ./cmd/consumer

FROM alpine AS runner

WORKDIR /app

COPY --from=build /build/app ./app

COPY --from=build /build/cmd/config/config.yml ./config.yml

COPY --from=build /build/cmd/consumer/moby.go ./moby.go

CMD ["/app/app", "--config=/app/config.yml"]