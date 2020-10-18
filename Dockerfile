# Build backend binary file
FROM golang:1.14.10-alpine3.12 AS be-builder
RUN apk add gcc g++ libc-dev
WORKDIR /go/src/github.com/mayswind/lab
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -a -v -i -trimpath -o lab lab.go

# Build frontend files
FROM node:12.19.0-alpine3.12 AS fe-builder
WORKDIR /go/src/github.com/mayswind/lab
COPY . .
RUN npm install && npm run build

# Package docker image
FROM alpine:3.12.0
RUN addgroup -S -g 1000 labapp && adduser -S -G labapp -u 1000 labapp
RUN apk --no-cache add su-exec tzdata
COPY --from=be-builder /go/src/github.com/mayswind/lab/lab /usr/local/bin/labapp/lab
RUN chmod +x /usr/local/bin/labapp/lab
COPY --from=fe-builder /go/src/github.com/mayswind/lab/dist /usr/local/bin/labapp/public
COPY conf /usr/local/bin/labapp/conf
WORKDIR /usr/local/bin/labapp
COPY docker/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
EXPOSE 8080
ENTRYPOINT ["/docker-entrypoint.sh"]
