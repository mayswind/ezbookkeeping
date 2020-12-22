# lab
[![License](https://img.shields.io/github/license/mayswind/lab.svg?style=flat)](https://github.com/mayswind/lab/blob/master/LICENSE)
[![Latest Build](https://img.shields.io/docker/cloud/build/mayswind/lab.svg?style=flat)](https://hub.docker.com/r/mayswind/lab/builds)
[![Go Report](https://goreportcard.com/badge/github.com/mayswind/lab)](https://goreportcard.com/report/github.com/mayswind/lab)
[![Latest Docker Image Size](https://img.shields.io/docker/image-size/mayswind/lab.svg?style=flat)](https://hub.docker.com/r/mayswind/lab)
[![Latest Release](https://img.shields.io/github/release/mayswind/lab.svg?style=flat)](https://github.com/mayswind/lab/releases)

## Introduction
lab is a lightweight account book app hosted by yourself. This project is now **under construction**.

## Features
1. Open source & Self-hosted
2. Lightweight & Fast
3. Easy to install
4. User-friendly interface
    * Almost native app experience
    * Two-level account support
    * Two-level category support
    * Preset various categories
5. Two-factor authentication
6. Application lock (WebAuthn support)
7. Multi-language support
8. Dark theme

## Screenshots
### Mobile Device
(Coming soon...)

## Installation
### Ship with docker
Visit [Docker Hub](https://hub.docker.com/r/mayswind/lab) to see all images and tags.

Latest Release:

    $ docker run -p8080:8080 mayswind/lab

Latest Daily Build:

    $ docker run -p8080:8080 mayswind/lab:latest-snapshot

### Install from binary

Latest release: [https://github.com/mayswind/lab/releases](https://github.com/mayswind/lab/releases)

    $ ./lab server run

lab will listen at port 8080 as default. You can visit http://<YOUR_HOST_ADDRESS>:8080/ .

### Build from source

Make sure you have [Golang](https://golang.org/), [GCC](http://gcc.gnu.org/), [Node.js](https://nodejs.org/) and [NPM](https://www.npmjs.com/) installed. Then download the source code, and follow these steps.

    # Build backend binary file
    $ GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -a -v -i -trimpath -o lab lab.go

    # Build frontend static files
    $ npm install
    $ npm run build

    # Copy files to target path
    $ cp lab <target>/lab
    $ cp -R dist <target>/public
    $ cp -R conf <target>/conf

All the files will be placed in `<target>` directory.

For more information about how to install lab, please visit our documentation.

## License
[MIT](https://github.com/mayswind/lab/blob/master/LICENSE)
