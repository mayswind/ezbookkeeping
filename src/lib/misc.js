import Clipboard from 'clipboard';
import CryptoJS from 'crypto-js';
import uaParser from 'ua-parser-js';

import { base64encode } from './common.js';

export function asyncLoadAssets(type, assetUrl) {
    return new Promise(function (resolve, reject) {
        let addElement = false;
        let el = null;

        if (type === 'js') {
            el = document.querySelector('script[src="' + assetUrl + '"]');
        } else if (type === 'css') {
            el = document.querySelector('link[href="' + assetUrl + '"]');
        } else {
            reject({
                type: type,
                assetUrl: assetUrl,
                error: 'notsupport'
            });
            return;
        }

        if (!el) {
            if (type === 'js') {
                el = document.createElement('script');
                el.setAttribute('type', 'text/javascript');
                el.setAttribute('async', 'true');
                el.setAttribute('src', assetUrl);
            } else if (type === 'css') {
                el = document.createElement('link');
                el.setAttribute('rel', 'stylesheet');
                el.setAttribute('type', 'text/css');
                el.setAttribute('href', assetUrl);
            }

            addElement = true;
        } else if (el.hasAttribute('data-loaded')) {
            resolve({
                type: type,
                assetUrl: assetUrl
            });
            return;
        }

        el.addEventListener('load', () => {
            el.setAttribute('data-loaded', true);
            resolve({
                type: type,
                assetUrl: assetUrl
            });
        });
        el.addEventListener('error', () => {
            reject({
                type: type,
                assetUrl: assetUrl,
                error: 'error'
            });
        });
        el.addEventListener('abort', () => {
            reject({
                type: type,
                assetUrl: assetUrl,
                error: 'abort'
            });
        });

        if (addElement) {
            document.head.appendChild(el);
        }
    });
}

export function generateRandomString() {
    let baseString = 'ebk_' + new Date().getTime();

    if (crypto && crypto.getRandomValues) {
        const randoms = new Uint8Array(256);
        crypto.getRandomValues(randoms);
        baseString += '_' + base64encode(randoms);
    } else {
        baseString += '_' + Math.random();
    }

    return CryptoJS.SHA256(baseString).toString();
}

export function generateRandomUUID() {
    const randomString = generateRandomString();

    // convert hash string to UUID Version 8
    const uuid = randomString.substring(0, 8) + '-'
        + randomString.substring(8, 12) + '-'
        + '8' + randomString.substring(13, 16) + '-'
        + (0x8 | (parseInt(randomString.charAt(16), 16) & 0x3)).toString(16) + randomString.substring(17, 20) + '-'
        + randomString.substring(20, 32);

    return uuid;
}

export function isSessionUserAgentCreatedByCli(ua) {
    return ua === 'ezbookkeeping Cli';
}

export function parseUserAgent(ua) {
    const uaParseRet = uaParser(ua);

    return {
        device: {
            vendor: uaParseRet.device.vendor,
            model: uaParseRet.device.model,
            type: uaParseRet.device.type
        },
        os: {
            name: uaParseRet.os.name,
            version: uaParseRet.os.version
        },
        browser: {
            name: uaParseRet.browser.name,
            version: uaParseRet.browser.version
        }
    };
}

export function parseDeviceInfo(uaInfo) {
    if (!uaInfo) {
        return '';
    }

    let result = '';

    if (uaInfo.device && uaInfo.device.model) {
        result = uaInfo.device.model;
    } else if (uaInfo.os && uaInfo.os.name) {
        result = uaInfo.os.name;

        if (uaInfo.os.version) {
            result += ' ' + uaInfo.os.version;
        }
    }

    if (uaInfo.browser && uaInfo.browser.name) {
        let browserInfo = uaInfo.browser.name;

        if (uaInfo.browser.version) {
            browserInfo += ' ' + uaInfo.browser.version;
        }

        if (result) {
            result += ' (' + browserInfo + ')';
        } else {
            result = browserInfo;
        }
    }

    if (!result) {
        return 'Unknown Device';
    }

    return result;
}

export function parseSessionInfo(token) {
    const isCreatedByCli = isSessionUserAgentCreatedByCli(token.userAgent);
    const uaInfo = parseUserAgent(token.userAgent);
    let deviceType = '';

    if (isCreatedByCli) {
        deviceType = 'cli';
    } else {
        if (uaInfo && uaInfo.device) {
            if (uaInfo.device.type === 'mobile') {
                deviceType = 'phone';
            } else if (uaInfo.device.type === 'wearable') {
                deviceType = 'wearable';
            } else if (uaInfo.device.type === 'tablet') {
                deviceType = 'tablet';
            } else if (uaInfo.device.type === 'smarttv') {
                deviceType = 'tv';
            } else {
                deviceType = 'default';
            }
        } else {
            deviceType = 'default';
        }
    }

    return {
        tokenId: token.tokenId,
        isCurrent: token.isCurrent,
        deviceType: deviceType,
        deviceInfo: isCreatedByCli ? token.userAgent : parseDeviceInfo(uaInfo),
        createdByCli: isCreatedByCli,
        lastSeen: token.lastSeen
    }
}

export function makeButtonCopyToClipboard({ text, el, successCallback, errorCallback }) {
    const clipboard = new Clipboard(el, {
        text: function () {
            return text;
        }
    });

    clipboard.on('success', (e) => {
        if (successCallback) {
            successCallback(e);
        }
    });

    clipboard.on('error', (e) => {
        if (errorCallback) {
            errorCallback(e);
        }
    });

    return clipboard;
}

export function changeClipboardObjectText(clipboard, text) {
    if (!clipboard) {
        return;
    }

    clipboard.text = function () {
        return text;
    };
}
