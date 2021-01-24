# Build backend binary file
FROM golang:1.14.10-alpine3.12 AS be-builder
WORKDIR /go/src/github.com/mayswind/lab
COPY . .
RUN docker/backend-build-pre-setup.sh
RUN apk add git gcc g++ libc-dev
RUN VERSION=`grep '"version": ' package.json | awk -F ':' '{print $2}' | tr -d ' ' | tr -d ',' | tr -d '"'` \
  && COMMIT_HASH=$(git rev-parse --short HEAD) \
  && GOOS=linux \
  && GOARCH=amd64 \
  && CGO_ENABLED=1 \
  && go build -a -v -i -trimpath -ldflags "-X main.version=${VERSION} -X main.commitHash=${COMMIT_HASH}" -o lab lab.go
RUN chmod +x lab

# Build frontend files
FROM node:12.19.0-alpine3.12 AS fe-builder
WORKDIR /go/src/github.com/mayswind/lab
COPY . .
RUN docker/frontend-build-pre-setup.sh
RUN apk add git
RUN npm install && npm run build

# Package docker image
FROM alpine:3.12.0
LABEL maintainer="MaysWind <i@mayswind.net>"
RUN addgroup -S -g 1000 labapp && adduser -S -G labapp -u 1000 labapp
COPY docker/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
RUN mkdir -p /usr/local/bin/labapp && chown 1000:1000 /usr/local/bin/labapp \
  && mkdir -p /usr/local/bin/labapp/data && chown 1000:1000 /usr/local/bin/labapp/data \
  && mkdir -p /var/log/labapp && chown 1000:1000 /var/log/labapp
WORKDIR /usr/local/bin/labapp
COPY --from=be-builder --chown=1000:1000 /go/src/github.com/mayswind/lab/lab /usr/local/bin/labapp/lab
COPY --from=fe-builder --chown=1000:1000 /go/src/github.com/mayswind/lab/dist /usr/local/bin/labapp/public
COPY --chown=1000:1000 conf /usr/local/bin/labapp/conf
USER 1000:1000
EXPOSE 8080
ENTRYPOINT ["/docker-entrypoint.sh"]
