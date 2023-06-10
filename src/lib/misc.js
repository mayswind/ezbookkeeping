import Clipboard from 'clipboard';
import CryptoJS from 'crypto-js';
import uaParser from 'ua-parser-js';

export function generateRandomString() {
    const baseString = 'ebk_' + Math.round(new Date().getTime() / 1000) + '_' + Math.random();
    return CryptoJS.SHA256(baseString).toString();
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

export function parseDeviceInfo(ua) {
    const uaInfo = parseUserAgent(ua);
    let result = '';

    if (uaInfo.device.model) {
        result = uaInfo.device.model;
    } else if (uaInfo.os.name) {
        result = uaInfo.os.name;

        if (uaInfo.os.version) {
            result += ' ' + uaInfo.os.version;
        }
    }

    if (uaInfo.browser.name) {
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
