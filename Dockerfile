# Build backend binary file
FROM golang:1.16.4-alpine3.13 AS be-builder
WORKDIR /go/src/github.com/mayswind/ezbookkeeping
COPY . .
RUN docker/backend-build-pre-setup.sh
RUN apk add git gcc g++ libc-dev
RUN VERSION=`grep '"version": ' package.json | awk -F ':' '{print $2}' | tr -d ' ' | tr -d ',' | tr -d '"'` \
  && COMMIT_HASH=$(git rev-parse --short HEAD) \
  && BUILD_UNIXTIME="$(date '+%s')" \
  && VERSION_FLAGS="-X github.com/mayswind/ezbookkeeping/pkg/version.Version=${VERSION} -X github.com/mayswind/ezbookkeeping/pkg/version.CommitHash=${COMMIT_HASH} -X github.com/mayswind/ezbookkeeping/pkg/version.BuildUnixTime=${BUILD_UNIXTIME}" \
  && GOOS=linux \
  && GOARCH=amd64 \
  && CGO_ENABLED=1 \
  && go build -a -v -trimpath -ldflags "-w -linkmode external -extldflags '-static' ${VERSION_FLAGS}" -o ezbookkeeping ezbookkeeping.go
RUN chmod +x ezbookkeeping

# Build frontend files
FROM node:14.17.0-alpine3.13 AS fe-builder
WORKDIR /go/src/github.com/mayswind/ezbookkeeping
COPY . .
RUN docker/frontend-build-pre-setup.sh
RUN apk add git
RUN npm install && npm run build

# Package docker image
FROM alpine:3.13.5
LABEL maintainer="MaysWind <i@mayswind.net>"
RUN addgroup -S -g 1000 ezbookkeeping && adduser -S -G ezbookkeeping -u 1000 ezbookkeeping
RUN apk --no-cache add tzdata
COPY docker/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
RUN mkdir -p /usr/local/bin/ezbookkeeping && chown 1000:1000 /usr/local/bin/ezbookkeeping \
  && mkdir -p /usr/local/bin/ezbookkeeping/data && chown 1000:1000 /usr/local/bin/ezbookkeeping/data \
  && mkdir -p /var/log/ezbookkeeping && chown 1000:1000 /var/log/ezbookkeeping
WORKDIR /usr/local/bin/ezbookkeeping
COPY --from=be-builder --chown=1000:1000 /go/src/github.com/mayswind/ezbookkeeping/ezbookkeeping /usr/local/bin/ezbookkeeping/ezbookkeeping
COPY --from=fe-builder --chown=1000:1000 /go/src/github.com/mayswind/ezbookkeeping/dist /usr/local/bin/ezbookkeeping/public
COPY --chown=1000:1000 conf /usr/local/bin/ezbookkeeping/conf
USER 1000:1000
EXPOSE 8080
ENTRYPOINT ["/docker-entrypoint.sh"]
