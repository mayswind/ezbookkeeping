# NAS 部署指南 (Mac M2 建置 -> NAS 執行)

本指南說明如何在 Mac (M2/Apple Silicon) 上編譯 ezBookkeeping 的 Docker 映像檔，並部署到 x86 架構的 NAS 上。

## 1. Mac 端：建置與推送

**事前準備：**
*   請確保 Mac 已安裝 [Docker Desktop](https://www.docker.com/products/docker-desktop/)。
*   請確保擁有 [Docker Hub](https://hub.docker.com/) 帳號。

**步驟：**

1.  開啟終端機 (Terminal)。
2.  登入 Docker Hub：
    ```bash
    docker login
    ```
    (輸入你的帳號與密碼/Token)

3.  執行建置腳本：
    請將 `<你的Docker帳號>` 替換為你實際的 Docker Hub 使用者名稱。
    ```bash
    ./docker/build-push-nas.sh <你的Docker帳號>/ezbookkeeping
    ```
    例如：
    ```bash
    ./docker/build-push-nas.sh myuser/ezbookkeeping
    ```

    這段過程會稍微花一點時間，因為它會進行跨平台編譯 (`linux/amd64`) 並將映像檔上傳到 Docker Hub。

## 2. NAS 端：拉取與執行

**事前準備：**
*   請確保 NAS 已安裝 Docker (Container Manager)。
*   請確保 NAS 可以連線到網際網路。

**步驟：**

1.  SSH 進入 NAS，或者使用 NAS 的 Docker GUI 介面。
2.  拉取剛剛上傳的映像檔：
    ```bash
    docker pull <你的Docker帳號>/ezbookkeeping
    ```
3.  啟動容器：
    ```bash
    docker run -d \
      --name ezbookkeeping \
      -p 8080:8080 \
      -v /volume1/docker/ezbookkeeping/data:/ezbookkeeping/data \
      -v /volume1/docker/ezbookkeeping/log:/ezbookkeeping/log \
      -v /volume1/docker/ezbookkeeping/storage:/ezbookkeeping/storage \
      <你的Docker帳號>/ezbookkeeping
    ```
    *   請依據你的 NAS 實際路徑調整 `-v` 的掛載路徑 (例如 `/volume1/...`)。
    *   `-p 8080:8080` 表示將 NAS 的 8080 port 對應到容器的 8080 port。

**資料庫設定說明：**
本映像檔已內建以下資料庫連線資訊 (已修改 `conf/ezbookkeeping.ini`)：
*   **Host**: `10.0.4.20:3306`
*   **Database**: `ezbookkeeping`
*   **User**: `ezbk`
*   **Type**: `mysql` (MariaDB)

只要 NAS 能透過區域網路連線到 `10.0.4.20`，啟動後即可直接運作。
