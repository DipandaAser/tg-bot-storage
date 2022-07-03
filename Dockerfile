FROM golang:alpine AS build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o bin/tg-bot-storage ./cmd/main.go


RUN apk add --no-cache ca-certificates

# build image with the binary
FROM scratch

# copy certificate to be able to make https request to telegram
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# copy the binary
COPY --from=build /build/bin/tg-bot-storage /

EXPOSE 7000

ENTRYPOINT ["/tg-bot-storage"]