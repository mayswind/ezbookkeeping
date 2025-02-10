# ezBookkeeping
[![License](https://img.shields.io/badge/license-MIT-green.svg)](https://github.com/mayswind/ezbookkeeping/blob/master/LICENSE)
[![Latest Build](https://img.shields.io/github/actions/workflow/status/mayswind/ezbookkeeping/docker-snapshot.yml?branch=main)](https://github.com/mayswind/ezbookkeeping/actions)
[![Go Report](https://goreportcard.com/badge/github.com/mayswind/ezbookkeeping)](https://goreportcard.com/report/github.com/mayswind/ezbookkeeping)
[![Latest Docker Image Size](https://img.shields.io/docker/image-size/mayswind/ezbookkeeping.svg?style=flat)](https://hub.docker.com/r/mayswind/ezbookkeeping)
[![Latest Release](https://img.shields.io/github/release/mayswind/ezbookkeeping.svg?style=flat)](https://github.com/mayswind/ezbookkeeping/releases)

## Introduction

ezBookkeeping is a lightweight, self-hosted personal bookkeeping app for Windows, macOS, and Linux. It supports x86, amd64, and ARM architectures and can also run on a Raspberry device. The application works with multiple databases, such as SQLite and MySQL. With Docker, it can be deployed easily using a single command.

Online Demo: [https://ezbookkeeping-demo.mayswind.net](https://ezbookkeeping-demo.mayswind.net)

## Features
- **Open source & Self-hosted**
- **Lightweight & Fast**
- **Simple Installation**
  - Supports Docker
  - Works with SQLite, MySQL, PostgreSQL, and more
  - Runs on Windows, macOS, and Linux (x86, amd64, ARM)
- **User-Friendly Interface**
  - Optimized for desktop and mobile
  - Two-level account and category management
  - Built-in categories
  - Location-based tracking with maps
  - Search and filter history records
  - Data statistics
  - Dark mode
- **Financial Management**
  - Multi-currency support with automatic exchange rate updates
  - Timezone support
- **Security & Data Management**
  - Two-factor authentication
  - App lock (PIN/WebAuthn)
  - Import & export data

## Screenshots
### Desktop Version
[![ezBookkeeping](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/desktop/en.png)](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/desktop/en.png)

### Mobile Version
[![ezBookkeeping](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/mobile/en.png)](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/mobile/en.png)

## Installation

### Using Docker
To deploy with Docker, run:

```sh
$ docker run -p8080:8080 mayswind/ezbookkeeping
```

For the latest development build:

```sh
$ docker run -p8080:8080 mayswind/ezbookkeeping:latest-snapshot
```

More images and tags are available on [Docker Hub](https://hub.docker.com/r/mayswind/ezbookkeeping).

### Installing from Binary
Download the latest release from [GitHub Releases](https://github.com/mayswind/ezbookkeeping/releases) and run:

**Linux / macOS:**
```sh
$ ./ezbookkeeping server run
```

**Windows:**
```sh
> .\ezbookkeeping.exe server run
```

The app runs on port 8080 by default. Access it via `http://{YOUR_HOST_ADDRESS}:8080/`.

### Building from Source
Ensure you have [Golang](https://golang.org/), [GCC](http://gcc.gnu.org/), [Node.js](https://nodejs.org/), and [NPM](https://www.npmjs.com/) installed. Then, clone the repository and build:

**Linux / macOS:**
```sh
$ ./build.sh package -o ezbookkeeping.tar.gz
```

**Windows:**
```sh
> .\build.bat package -o ezbookkeeping.zip
```

To build a Docker image, ensure [Docker](https://www.docker.com/) is installed, then run:

```sh
$ ./build.sh docker
```

## Documentation
- [English](http://ezbookkeeping.mayswind.net)
- [简体中文 (Simplified Chinese)](http://ezbookkeeping.mayswind.net/zh_Hans)

## License
[MIT](https://github.com/mayswind/ezbookkeeping/blob/master/LICENSE)
