# ezBookkeeping
[![License](https://img.shields.io/badge/license-MIT-green.svg)](https://github.com/mayswind/ezbookkeeping/blob/master/LICENSE)
[![Latest Build](https://img.shields.io/github/actions/workflow/status/mayswind/ezbookkeeping/docker-snapshot.yml?branch=main)](https://github.com/mayswind/ezbookkeeping/actions)
[![Go Report](https://goreportcard.com/badge/github.com/mayswind/ezbookkeeping)](https://goreportcard.com/report/github.com/mayswind/ezbookkeeping)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/mayswind/ezbookkeeping)
[![Latest Docker Image Size](https://img.shields.io/docker/image-size/mayswind/ezbookkeeping.svg?style=flat)](https://hub.docker.com/r/mayswind/ezbookkeeping)
[![Latest Release](https://img.shields.io/github/release/mayswind/ezbookkeeping.svg?style=flat)](https://github.com/mayswind/ezbookkeeping/releases)

## Introduction
ezBookkeeping is a lightweight, self-hosted personal finance app with a sleek, user-friendly interface and powerful bookkeeping features. Built with simplicity and portability in mind, it's easy to deploy, easy to use, and requires minimal system resources — perfect for microservers, NAS devices, and even Raspberry Pi.

The app is fully cross-platform and device-friendly — you can use it seamlessly on **mobile, tablet, and desktop devices**. With support for PWA (Progressive Web Apps), you can even [add it to your mobile home screen](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/mobile/add_to_home_screen.gif) and use it like a native app.

Live Demo: [https://ezbookkeeping-demo.mayswind.net](https://ezbookkeeping-demo.mayswind.net)

## Features
- **Open Source & Self-Hosted**
    - Built for privacy and control
- **Lightweight & Fast**
    - Optimized for performance, runs smoothly even on low-resource environments
- **Easy Installation**
    - Docker-ready
    - Supports SQLite, MySQL, PostgreSQL
    - Cross-platform (Windows, macOS, Linux)
    - Works on x86, amd64, ARM architectures
- **User-Friendly Interface**
    - UI optimized for both mobile and desktop
    - PWA support for native-like mobile experience
    - Dark mode
- **AI-Powered Features**
    - Supports MCP (Model Context Protocol) for AI integration
- **Powerful Bookkeeping**
    - Two-level accounts and categories
    - Attach images to transactions
    - Location tracking with maps
    - Recurring transactions
    - Advanced filtering, search, visualization, and analysis
- **Localization & Globalization**
    - Multi-language and multi-currency support
    - Automatic exchange rates
    - Multi-timezone awareness
    - Custom formats for dates, numbers, and currencies
- **Security**
    - Two-factor authentication (2FA)
    - Login rate limiting
    - Application lock (PIN code / WebAuthn)
- **Data Import/Export**
    - Supports CSV, OFX, QFX, QIF, IIF, Camt.053, MT940, GnuCash, Firefly III, Beancount, and more

## Screenshots
### Desktop Version
[![ezBookkeeping](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/desktop/en.png)](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/desktop/en.png)

### Mobile Version
[![ezBookkeeping](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/mobile/en.png)](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/mobile/en.png)

## Installation
### Run with Docker
Visit [Docker Hub](https://hub.docker.com/r/mayswind/ezbookkeeping) to see all images and tags.

**Latest Release:**

    $ docker run -p 8080:8080 mayswind/ezbookkeeping

**Latest Daily Build:**

    $ docker run -p 8080:8080 mayswind/ezbookkeeping:latest-snapshot

### Install from Binary
Download the latest release: [https://github.com/mayswind/ezbookkeeping/releases](https://github.com/mayswind/ezbookkeeping/releases)

**Linux / macOS**

    $ ./ezbookkeeping server run

**Windows**

    > .\ezbookkeeping.exe server run

By default, ezBookkeeping listens on port 8080. You can then visit `http://{YOUR_HOST_ADDRESS}:8080/` .

### Build from Source
Make sure you have [Golang](https://golang.org/), [GCC](http://gcc.gnu.org/), [Node.js](https://nodejs.org/) and [NPM](https://www.npmjs.com/) installed. Then download the source code, and follow these steps:

**Linux / macOS**

    $ ./build.sh package -o ezbookkeeping.tar.gz

All the files will be packaged in `ezbookkeeping.tar.gz`.

**Windows**

    > .\build.bat package -o ezbookkeeping.zip

All the files will be packaged in `ezbookkeeping.zip`.

You can also build a Docker image. Make sure you have [Docker](https://www.docker.com/) installed, then follow these steps:

**Linux**

    $ ./build.sh docker

## Contributing
We welcome contributions of all kinds!

Found a bug? [Submit an issue](https://github.com/mayswind/ezbookkeeping/issues)

Want to contribute code? Feel free to fork and send a pull request.

Contributions of all kinds — bug reports, feature suggestions, documentation improvements, or code — are highly appreciated.

Check out our [Contributor Graph](https://github.com/mayswind/ezbookkeeping/graphs/contributors) to see the amazing people who’ve already helped.

## Translating
Help make ezBookkeeping accessible to users around the world! If you want to contribute a translation, please refer to our [translation guide](https://ezbookkeeping.mayswind.net/translating).

Currently available translations:

| Tag | Language | Contributors |
| --- | --- | --- |
| de | Deutsch | [@chrgm](https://github.com/chrgm) |
| en | English | / |
| es | Español | [@Miguelonlonlon](https://github.com/Miguelonlonlon) |
| it | Italiano | [@waron97](https://github.com/waron97) |
| ja | 日本語 | [@tkymmm](https://github.com/tkymmm) |
| pt-BR | Português (Brasil) | [@thecodergus](https://github.com/thecodergus) |
| ru | Русский | [@artegoser](https://github.com/artegoser) |
| uk | Українська | [@nktlitvinenko](https://github.com/nktlitvinenko) |
| vi | Tiếng Việt | [@f97](https://github.com/f97) |
| zh-Hans | 中文 (简体) | / |
| zh-Hant | 中文 (繁體) | / |

Don't see your language? Help us add it!

## Documentation
1. [English](http://ezbookkeeping.mayswind.net)
1. [中文 (简体)](http://ezbookkeeping.mayswind.net/zh_Hans)

## License
[MIT](https://github.com/mayswind/ezbookkeeping/blob/master/LICENSE)
