import Cookies from 'js-cookie';

import { base64decode } from '@/lib/common.js';

const serverSettingsCookieKey = 'ebk_server_settings';

function getServerSetting(key) {
    const settings = Cookies.get(serverSettingsCookieKey) || '';
    const settingsArr = settings.split('_');

    for (let i = 0; i < settingsArr.length; i++) {
        const pairs = settingsArr[i].split('.');

        if (pairs[0] === key) {
            return pairs[1];
        }
    }

    return undefined;
}

function getServerDecodedSetting(key) {
    const value = getServerSetting(key);

    if (!value) {
        return value;
    }

    return decodeURIComponent(base64decode(value));
}

export function isUserRegistrationEnabled() {
    return getServerSetting('r') === '1';
}

export function isDataExportingEnabled() {
    return getServerSetting('e') === '1';
}

export function getMapProvider() {
    return getServerSetting('m');
}

export function isMapDataFetchProxyEnabled() {
    return getServerSetting('mp') === '1';
}

export function getTomTomMapAPIKey() {
    return getServerDecodedSetting('tmak');
}

export function getGoogleMapAPIKey() {
    return getServerDecodedSetting('gmak');
}

export function getBaiduMapAK() {
    return getServerDecodedSetting('bmak');
}

export function getAmapApplicationKey() {
    return getServerDecodedSetting('amak');
}

export function getAmapSecurityVerificationMethod() {
    return getServerSetting('amsv');
}

export function getAmapApiExternalProxyUrl() {
    return getServerDecodedSetting('amep');
}

export function getAmapApplicationSecret() {
    return getServerDecodedSetting('amas');
}
