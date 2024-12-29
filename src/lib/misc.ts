import CryptoJS from 'crypto-js';

import { base64encode } from './common.ts';

export interface AsyncLoadAssetsResult {
    readonly type: string;
    readonly assetUrl: string;
}

export function asyncLoadAssets(type: string, assetUrl: string): Promise<AsyncLoadAssetsResult> {
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

        if (!el) {
            reject({
                type: type,
                assetUrl: assetUrl,
                error: 'unexpected'
            });
            return;
        }

        el.addEventListener('load', () => {
            el.setAttribute('data-loaded', 'true');
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

export function generateRandomString(): string {
    let baseString = 'ebk_' + new Date().getTime();

    if (crypto && crypto.getRandomValues) {
        const randoms = new Uint8Array(256);
        crypto.getRandomValues(randoms);
        baseString += '_' + base64encode(randoms.buffer);
    } else {
        baseString += '_' + Math.random();
    }

    return CryptoJS.SHA256(baseString).toString();
}

export function generateRandomUUID(): string {
    const randomString = generateRandomString();

    // convert hash string to UUID Version 8
    const uuid = randomString.substring(0, 8) + '-'
        + randomString.substring(8, 12) + '-'
        + '8' + randomString.substring(13, 16) + '-'
        + (0x8 | (parseInt(randomString.charAt(16), 16) & 0x3)).toString(16) + randomString.substring(17, 20) + '-'
        + randomString.substring(20, 32);

    return uuid;
}
