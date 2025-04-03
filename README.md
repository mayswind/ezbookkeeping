# ezBookkeeping
[![License](https://img.shields.io/badge/license-MIT-green.svg)](https://github.com/mayswind/ezbookkeeping/blob/master/LICENSE)
[![Latest Build](https://img.shields.io/github/actions/workflow/status/mayswind/ezbookkeeping/docker-snapshot.yml?branch=main)](https://github.com/mayswind/ezbookkeeping/actions)
[![Go Report](https://goreportcard.com/badge/github.com/mayswind/ezbookkeeping)](https://goreportcard.com/report/github.com/mayswind/ezbookkeeping)
[![Latest Docker Image Size](https://img.shields.io/docker/image-size/mayswind/ezbookkeeping.svg?style=flat)](https://hub.docker.com/r/mayswind/ezbookkeeping)
[![Latest Release](https://img.shields.io/github/release/mayswind/ezbookkeeping.svg?style=flat)](https://github.com/mayswind/ezbookkeeping/releases)

## Introduction
ezBookkeeping is a lightweight self-hosted personal bookkeeping app with user-friendly interface for both desktop and mobile devices. It supports PWA, you can [add the app homepage to the home screen](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/mobile/add_to_home_screen.gif) of your mobile device and use it just like a native app. It's easily to be deployed and configured, you can just deploy it by a single command via Docker. It supports almost all platforms, including Windows, macOS, and Linux, and is compatible with x86, amd64 and ARM hardware architectures. It only requires very few system resources, and you can even run it on a Raspberry Pi device.

Online Demo: [https://ezbookkeeping-demo.mayswind.net](https://ezbookkeeping-demo.mayswind.net)

## Features
1. Open Source & Self-Hosted
2. Lightweight & Fast
3. Easy Installation
    * Support Docker
    * Support multiple databases (SQLite, MySQL, PostgreSQL, etc.)
    * Support multiple operation system & hardware architectures (Windows, macOS, Linux & x86, amd64, ARM)
4. User-Friendly Interface
    * Native UI for both desktop and mobile devices
    * Support PWA, providing near-native experience for mobile devices
    * Dark theme
5. Powerful Bookkeeping Features
    * Support two-level account
    * Support two-level transaction categories and predefined categories
    * Support transaction pictures
    * Support geographic location tracking and map
    * Support recurring transactions
    * Search and filtering transaction records
    * Data visualization and statistical analysis
6. Localization Support
    * Multi-language support
    * Multi-currency support with automatic exchange rate updates from various financial institutions
    * Multi-timezone support
    * Customizable date, time, number and currency display formats
7. Security & Reliability
    * Two-factor authentication (2FA)
    * Login rate limiting
    * Application lock (PIN code / WebAuthn)
8. Data Export & Import (CSV, OFX, QFX, QIF, IIF, GnuCash, FireFly III, Beancount, etc.)

## Screenshots
### Desktop Version
[![ezBookkeeping](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/desktop/en.png)](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/desktop/en.png)

### Mobile Version
[![ezBookkeeping](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/mobile/en.png)](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/mobile/en.png)

## Installation
### Ship with docker
Visit [Docker Hub](https://hub.docker.com/r/mayswind/ezbookkeeping) to see all images and tags.

Latest Release:

    $ docker run -p8080:8080 mayswind/ezbookkeeping

Latest Daily Build:

    $ docker run -p8080:8080 mayswind/ezbookkeeping:latest-snapshot

### Install from binary
Latest release: [https://github.com/mayswind/ezbookkeeping/releases](https://github.com/mayswind/ezbookkeeping/releases)

**Linux / macOS**

    $ ./ezbookkeeping server run

**Windows**

    > .\ezbookkeeping.exe server run

ezBookkeeping will listen at port 8080 as default. Then you can visit `http://{YOUR_HOST_ADDRESS}:8080/` .

### Build from source
Make sure you have [Golang](https://golang.org/), [GCC](http://gcc.gnu.org/), [Node.js](https://nodejs.org/) and [NPM](https://www.npmjs.com/) installed. Then download the source code, and follow these steps:

**Linux / macOS**

    $ ./build.sh package -o ezbookkeeping.tar.gz

All the files will be packaged in `ezbookkeeping.tar.gz`.

**Windows**

    > .\build.bat package -o ezbookkeeping.zip

All the files will be packaged in `ezbookkeeping.zip`.

You can also build docker image, make sure you have [docker](https://www.docker.com/) installed, then follow these steps:

**Linux**

    $ ./build.sh docker

## Documents
1. [English](http://ezbookkeeping.mayswind.net)
1. [中文 (简体)](http://ezbookkeeping.mayswind.net/zh_Hans)

## License
[MIT](https://github.com/mayswind/ezbookkeeping/blob/master/LICENSE)
