# ezBookkeeping
[![License](https://img.shields.io/badge/license-MIT-green.svg)](https://github.com/mayswind/ezbookkeeping/blob/master/LICENSE)
[![Go Report](https://goreportcard.com/badge/github.com/mayswind/ezbookkeeping)](https://goreportcard.com/report/github.com/mayswind/ezbookkeeping)
[![Latest Release](https://img.shields.io/github/release/mayswind/ezbookkeeping.svg?style=flat)](https://github.com/mayswind/ezbookkeeping/releases)
[![Latest Build](https://img.shields.io/github/actions/workflow/status/mayswind/ezbookkeeping/build-snapshot.yml?branch=main)](https://github.com/mayswind/ezbookkeeping/actions)
[![Latest Docker Image Size](https://img.shields.io/docker/image-size/mayswind/ezbookkeeping.svg?style=flat)](https://hub.docker.com/r/mayswind/ezbookkeeping)
[![Docker Pulls](https://img.shields.io/docker/pulls/mayswind/ezbookkeeping)](https://hub.docker.com/r/mayswind/ezbookkeeping)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/mayswind/ezbookkeeping)

[![Recommend By HelloGitHub](https://api.hellogithub.com/v1/widgets/recommend.svg?rid=ded5af09da574ec1811ddb154f1b2093&claim_uid=LT7EZxeBukCnh0K)](https://hellogithub.com/en/repository/mayswind/ezbookkeeping)
[![Trending](https://trendshift.io/api/badge/repositories/12917)](https://trendshift.io/repositories/12917)

## Introduction
ezBookkeeping is a lightweight, self-hosted personal finance app with a user-friendly interface and powerful bookkeeping features. It's easy to deploy, and you can start it with just one single Docker command. Designed to be resource-efficient and highly scalable, it can run smoothly on devices as small as a Raspberry Pi, or scale up to NAS, MicroServers, and even large cluster environments.

ezBookkeeping offers tailored interfaces for both mobile and desktop devices. With support for PWA (Progressive Web Apps), you can even [add it to your mobile home screen](https://raw.githubusercontent.com/wiki/mayswind/ezbookkeeping/img/mobile/add_to_home_screen.gif) and use it like a native app.

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
    - Receipt image recognition
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

    $ docker run -p8080:8080 mayswind/ezbookkeeping

**Latest Daily Build:**

    $ docker run -p8080:8080 mayswind/ezbookkeeping:latest-snapshot

### Install from Binary
Download the latest release: [https://github.com/mayswind/ezbookkeeping/releases](https://github.com/mayswind/ezbookkeeping/releases)

**Linux / macOS**

    $ ./ezbookkeeping server run

**Windows**

    > .\ezbookkeeping.exe server run

By default, ezBookkeeping listens on port 8080. You can then visit `http://{YOUR_HOST_ADDRESS}:8080/` .

### Build from Source
Make sure you have [Golang](https://golang.org/), [GCC](https://gcc.gnu.org/), [Node.js](https://nodejs.org/) and [NPM](https://www.npmjs.com/) installed. Then download the source code, and follow these steps:

**Linux / macOS**

    $ ./build.sh package -o ezbookkeeping.tar.gz

All the files will be packaged in `ezbookkeeping.tar.gz`.

**Windows**

    > .\build.bat package -o ezbookkeeping.zip

or

    PS > .\build.ps1 package -Output ezbookkeeping.zip

All the files will be packaged in `ezbookkeeping.zip`.

You can also build a Docker image. Make sure you have [Docker](https://www.docker.com/) installed, then follow these steps:

**Linux**

    $ ./build.sh docker

## Contributing
We welcome contributions of all kinds.

Found a bug? [Submit an issue](https://github.com/mayswind/ezbookkeeping/issues)

Want to contribute code? Feel free to fork and send a pull request.

Contributions of all kinds — bug reports, feature suggestions, documentation improvements, or code — are highly appreciated.

Check out our [Contributor Graph](https://github.com/mayswind/ezbookkeeping/graphs/contributors) to see the amazing people who've already helped.

## Translating
Help make ezBookkeeping accessible to users around the world. If you want to contribute a translation, please refer to our [translation guide](https://ezbookkeeping.mayswind.net/translating).

Currently available translations:

| Tag | Language | Contributors |
| --- | --- | --- |
| de | Deutsch | [@chrgm](https://github.com/chrgm) |
| en | English | / |
| es | Español | [@Miguelonlonlon](https://github.com/Miguelonlonlon), [@abrugues](https://github.com/abrugues), [@AndresTeller](https://github.com/AndresTeller) |
| fr | Français | [@brieucdlf](https://github.com/brieucdlf) |
| it | Italiano | [@waron97](https://github.com/waron97) |
| ja | 日本語 | [@tkymmm](https://github.com/tkymmm) |
| kn | ಕನ್ನಡ | [@Darshanbm05](https://github.com/Darshanbm05) |
| ko | 한국어 | [@overworks](https://github.com/overworks) |
| nl | Nederlands | [@automagics](https://github.com/automagics) |
| pt-BR | Português (Brasil) | [@thecodergus](https://github.com/thecodergus) |
| ru | Русский | [@artegoser](https://github.com/artegoser) |
| sl | Slovenščina | [@thehijacker](https://github.com/thehijacker) |
| ta | தமிழ் | [@hhharsha36](https://github.com/hhharsha36) |
| th | ไทย | [@natthavat28](https://github.com/natthavat28) |
| tr | Türkçe | [@aydnykn](https://github.com/aydnykn) |
| uk | Українська | [@nktlitvinenko](https://github.com/nktlitvinenko) |
| vi | Tiếng Việt | [@f97](https://github.com/f97) |
| zh-Hans | 中文 (简体) | / |
| zh-Hant | 中文 (繁體) | / |

Don't see your language? Help us add it.

## Documentation
1. [English](https://ezbookkeeping.mayswind.net)
1. [中文 (简体)](https://ezbookkeeping.mayswind.net/zh_Hans)

## License
[MIT](https://github.com/mayswind/ezbookkeeping/blob/master/LICENSE)
