import { keys } from '@/core/base.ts';
import { TextDirection } from '@/core/text.ts';

import type {
    ApplicationSettingKey,
    ApplicationSettingValue,
    ApplicationSettingSubValue,
    BaseApplicationSetting,
    ApplicationSettings,
    LocaleDefaultSettings
} from '@/core/setting.ts';

import {
    DEFAULT_APPLICATION_SETTINGS,
    DEFAULT_LOCALE_SETTINGS
} from '@/core/setting.ts';

import {
    isObject,
    isString
} from '@/lib/common.ts';

const settingsLocalStorageKey: string = 'ebk_app_settings';
const currentLanguageSessionStorageKey: string = 'ebk_current_language';
const languagePreviewSessionStorageKeyPrefix: string = 'ebk_language_preview_';

export interface SessionLanguagePreview {
    readonly data: string;
    readonly textDirection: TextDirection;
}

function getStoredApplicationSettings(): BaseApplicationSetting {
    try {
        const storageData = localStorage.getItem(settingsLocalStorageKey) || '{}';
        return JSON.parse(storageData);
    } catch (ex) {
        console.warn('settings in local storage is invalid', ex);
        return {};
    }
}

export function getApplicationSettings(): ApplicationSettings {
    const storedApplicationSettings = getStoredApplicationSettings();

    for (const key of keys(storedApplicationSettings)) {
        if (typeof(DEFAULT_APPLICATION_SETTINGS[key]) === 'object') {
            storedApplicationSettings[key] = Object.assign({}, DEFAULT_APPLICATION_SETTINGS[key], storedApplicationSettings[key]);
        }
    }

    return Object.assign({}, DEFAULT_APPLICATION_SETTINGS, storedApplicationSettings);
}

export function getLocaleDefaultSettings(): LocaleDefaultSettings {
    return Object.assign({}, DEFAULT_LOCALE_SETTINGS);
}

function updateApplicationSettings(settings: ApplicationSettings): void {
    const storageData = JSON.stringify(settings);
    return localStorage.setItem(settingsLocalStorageKey, storageData);
}

export function updateApplicationSettingsValue(key: ApplicationSettingKey, value: ApplicationSettingValue): void {
    if (!Object.prototype.hasOwnProperty.call(DEFAULT_APPLICATION_SETTINGS, key)) {
        return;
    }

    const settings = getApplicationSettings();
    settings[key] = value;

    return updateApplicationSettings(settings);
}

export function updateApplicationSettingsSubValue(key: ApplicationSettingKey, subKey: ApplicationSettingKey, value: ApplicationSettingSubValue): void {
    if (!Object.prototype.hasOwnProperty.call(DEFAULT_APPLICATION_SETTINGS, key)) {
        return;
    }

    if (!Object.prototype.hasOwnProperty.call(DEFAULT_APPLICATION_SETTINGS[key], subKey)) {
        return;
    }

    const settings = getApplicationSettings();
    let options = settings[key];

    if (!options) {
        options = {};
    }

    (options as Record<string, ApplicationSettingSubValue>)[subKey] = value;
    settings[key] = options;

    return updateApplicationSettings(settings);
}

export function isEnableDebug(): boolean {
    return getApplicationSettings().debug;
}

export function getTheme(): string {
    return getApplicationSettings().theme;
}

export function getTimeZone(): string {
    return getApplicationSettings().timeZone;
}

export function isEnableApplicationLock(): boolean {
    return getApplicationSettings().applicationLock;
}

export function isEnableSwipeBack(): boolean {
    return getApplicationSettings().swipeBack;
}

export function isEnableAnimate(): boolean {
    return getApplicationSettings().animate;
}

export function clearSettings(): void {
    localStorage.removeItem(settingsLocalStorageKey);
}

export function getSessionCurrentLanguageKey(): string {
    return sessionStorage.getItem(currentLanguageSessionStorageKey) || '';
}

export function setSessionCurrentLanguageKey(languageKey: string): void {
    if (!languageKey) {
        sessionStorage.removeItem(currentLanguageSessionStorageKey);
        return;
    }

    sessionStorage.setItem(currentLanguageSessionStorageKey, languageKey);
}

export function getSessionLanguagePreview(languageKey: string): SessionLanguagePreview | null {
    try {
        const storageData = sessionStorage.getItem(languagePreviewSessionStorageKeyPrefix + languageKey);

        if (!storageData) {
            return null;
        }

        const preview: unknown = JSON.parse(storageData) as unknown;

        if (!preview || !isObject(preview) || !('data' in preview) || !('textDirection' in preview)) {
            console.warn('language preview object is invalid', preview);
            return null;
        }

        const data = preview['data'];
        const textDirection = preview['textDirection'];

        if (!preview || !isString(data) || (textDirection !== TextDirection.LTR && textDirection !== TextDirection.RTL)) {
            console.warn('fields in language preview object are invalid', preview);
            return null;
        }

        return {
            data: data,
            textDirection: textDirection
        };
    } catch (ex) {
        console.warn('language preview in session storage is invalid', ex);
        return null;
    }
}

export function setSessionLanguagePreview(languageKey: string, preview: SessionLanguagePreview | null): void {
    const storageKey = languagePreviewSessionStorageKeyPrefix + languageKey;

    if (!preview) {
        sessionStorage.removeItem(storageKey);
    } else {
        sessionStorage.setItem(storageKey, JSON.stringify(preview));
    }
}
