FROM docker.io/golang:1.22-alpine AS builder


RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

ADD . /app
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o /server /app/cmd/server

FROM scratch
COPY --from=builder /server /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/server"]
