FROM docker.io/golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

ADD . /app
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o /server /app/cmd/server

FROM scratch
COPY --from=builder /server /
ENTRYPOINT ["/server"]
