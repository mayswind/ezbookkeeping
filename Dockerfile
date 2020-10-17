# Build binary release
FROM golang:1.14.10-alpine3.12 AS builder
RUN apk add gcc g++ libc-dev
WORKDIR /go/src/github.com/mayswind/lab
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -a -v -i -trimpath -o lab lab.go

# Package docker image
FROM alpine:3.12.0
RUN addgroup -S -g 1000 labapp && adduser -S -G labapp -u 1000 labapp
RUN apk --no-cache add su-exec tzdata
COPY --from=builder /go/src/github.com/mayswind/lab/lab /usr/local/bin/labapp/lab
RUN chmod +x /usr/local/bin/labapp/lab
COPY conf /usr/local/bin/labapp/conf
COPY public /usr/local/bin/labapp/public
WORKDIR /usr/local/bin/labapp
COPY docker/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
EXPOSE 8080
ENTRYPOINT ["/docker-entrypoint.sh"]
