import Cookies from 'js-cookie'

const settingsLocalStorageKey = 'lab_user_settings';
const serverSettingsCookieKey = 'ACP_SETTINGS';

const defaultSettings = {
    lang: 'en',
    animate: true
};

function getOriginalSettings() {
    try {
        const storageData = localStorage.getItem(settingsLocalStorageKey) || '{}';
        return JSON.parse(storageData);
    } catch (ex) {
        console.warn('settings in local storage is invalid', ex);
        return {};
    }
}

function getFinalSettings() {
    return Object.assign({}, defaultSettings, getOriginalSettings());
}

function setSettings(settings) {
    const storageData = JSON.stringify(settings);
    return localStorage.setItem(settingsLocalStorageKey, storageData);
}

function getOriginalOption(key) {
    return getOriginalSettings()[key];
}

function setOption(key, value) {
    if (!Object.prototype.hasOwnProperty.call(defaultSettings, key)) {
        return;
    }

    const settings = getFinalSettings();
    settings[key] = value;

    return setSettings(settings);
}

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

export default {
    getLanguage: () => getOriginalOption('lang'),
    setLanguage: value => setOption('lang', value),
    isEnableAnimate: () => getOriginalOption('animate'),
    setEnableAnimate: value => setOption('animate', value),
    isUserRegistrationEnabled: () => getServerSetting('r') === '1'
};
