# Build backend binary file
FROM golang:1.20.2-alpine3.17 AS be-builder
ARG RELEASE_BUILD
ENV RELEASE_BUILD=$RELEASE_BUILD
WORKDIR /go/src/github.com/mayswind/ezbookkeeping
COPY . .
RUN docker/backend-build-pre-setup.sh
RUN apk add git gcc g++ libc-dev
RUN ./build.sh backend

# Build frontend files
FROM node:18.15.0-alpine3.17 AS fe-builder
ARG RELEASE_BUILD
ENV RELEASE_BUILD=$RELEASE_BUILD
WORKDIR /go/src/github.com/mayswind/ezbookkeeping
COPY . .
RUN docker/frontend-build-pre-setup.sh
RUN apk add git
RUN ./build.sh frontend

# Package docker image
FROM alpine:3.17.2
LABEL maintainer="MaysWind <i@mayswind.net>"
RUN addgroup -S -g 1000 ezbookkeeping && adduser -S -G ezbookkeeping -u 1000 ezbookkeeping
RUN apk --no-cache add tzdata
COPY docker/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
RUN mkdir -p /ezbookkeeping && chown 1000:1000 /ezbookkeeping \
  && mkdir -p /ezbookkeeping/data && chown 1000:1000 /ezbookkeeping/data \
  && mkdir -p /ezbookkeeping/log && chown 1000:1000 /ezbookkeeping/log
WORKDIR /ezbookkeeping
COPY --from=be-builder --chown=1000:1000 /go/src/github.com/mayswind/ezbookkeeping/ezbookkeeping /ezbookkeeping/ezbookkeeping
COPY --from=fe-builder --chown=1000:1000 /go/src/github.com/mayswind/ezbookkeeping/dist /ezbookkeeping/public
COPY --chown=1000:1000 conf /ezbookkeeping/conf
COPY --chown=1000:1000 LICENSE /ezbookkeeping/LICENSE
USER 1000:1000
EXPOSE 8080
ENTRYPOINT ["/docker-entrypoint.sh"]
