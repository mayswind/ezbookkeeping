const settingsLocalStorageKey = 'lab_user_settings';
const defaultSettings = {
    lang: 'en'
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

export default {
    getLanguage: () => getOriginalOption('lang'),
    setLanguage: value => setOption('lang', value)
};
