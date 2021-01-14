# lab
[![License](https://img.shields.io/badge/license-MIT-green.svg)](https://github.com/mayswind/lab/blob/master/LICENSE)
[![Latest Build](https://img.shields.io/docker/cloud/build/mayswind/lab.svg?style=flat)](https://hub.docker.com/r/mayswind/lab/builds)
[![Go Report](https://goreportcard.com/badge/github.com/mayswind/lab)](https://goreportcard.com/report/github.com/mayswind/lab)
[![Latest Docker Image Size](https://img.shields.io/docker/image-size/mayswind/lab.svg?style=flat)](https://hub.docker.com/r/mayswind/lab)
[![Latest Release](https://img.shields.io/github/release/mayswind/lab.svg?style=flat)](https://github.com/mayswind/lab/releases)

## Introduction
The lab is a lightweight account book app hosted by yourself. It can be deployed on almost all platforms, including Windows, macOS and Linux on x86, amd64 and ARM architectures. You can even deploy it on an raspberry device. It also supports many different databases, including sqlite and mysql. With docker, you can just deploy it via one command without complicated configuration.

This project is still **under construction**.

## Features
1. Open source & Self-hosted
2. Lightweight & Fast
3. Easy to install
    * Docker support
    * Multiple database support (sqlite, mysql, etc.)
    * Multiple os & architecture support (Windows, macOS, Linux & x86, amd64, ARM)
4. User-friendly interface
    * Desktop (planning) and mobile support
    * Almost native app experience (for mobile device)
    * Two-level account & two-level category support
    * Plentiful preset categories
    * Searching & filtering history records
    * Data statistics
    * Dark theme
5. Multiple currency support & automatically updating exchange rates
6. Two-factor authentication
7. Application lock (WebAuthn support)
8. Data export
9. Multi-language support

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

The lab will listen at port 8080 as default. Then you can visit http://<YOUR_HOST_ADDRESS>:8080/ .

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
