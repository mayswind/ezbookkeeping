import CryptoJS from 'crypto-js';

import settings from './settings.js';
import utils from './utils.js';

const APP_LOCK_SECRET_BASE_STRING_PREFIX = 'LAB_LOCK_SECRET_';

const tokenLocalStorageKey = 'lab_user_token';
const userInfoLocalStorageKey = 'lab_user_info';
const tokenSessionStorageKey = 'lab_user_session_token';
const appLockSecretSessionStorageKey = 'lab_user_app_lock_secret';

function getAppLockSecret(pinCode) {
    return CryptoJS.SHA256(APP_LOCK_SECRET_BASE_STRING_PREFIX + pinCode).toString();
}

function getEncryptedToken(token, secret) {
    return CryptoJS.AES.encrypt(token, secret).toString();
}

function getDecryptedToken(encryptedToken, secret) {
    const bytes = CryptoJS.AES.decrypt(encryptedToken, secret);
    return bytes.toString(CryptoJS.enc.Utf8);
}

function getToken() {
    if (settings.isEnableApplicationLock()) {
        return sessionStorage.getItem(tokenSessionStorageKey);
    } else {
        return localStorage.getItem(tokenLocalStorageKey);
    }
}

function getUserInfo() {
    const data = localStorage.getItem(userInfoLocalStorageKey);
    return JSON.parse(data);
}

function isUserLogined() {
    return !!localStorage.getItem(tokenLocalStorageKey);
}

function isUserUnlocked() {
    if (!isUserLogined()) {
        return false;
    }

    if (!settings.isEnableApplicationLock()) {
        return true;
    }

    return !!sessionStorage.getItem(appLockSecretSessionStorageKey) && !!sessionStorage.getItem(tokenSessionStorageKey);
}

function unlockToken(pinCode) {
    const encryptedToken = localStorage.getItem(tokenLocalStorageKey);
    const secret = getAppLockSecret(pinCode);
    const token = getDecryptedToken(encryptedToken, secret);

    sessionStorage.setItem(appLockSecretSessionStorageKey, secret);
    sessionStorage.setItem(tokenSessionStorageKey, token);
}

function encryptToken(pinCode) {
    const token = localStorage.getItem(tokenLocalStorageKey);
    const secret = getAppLockSecret(pinCode);
    const encryptedToken = getEncryptedToken(token, secret);

    sessionStorage.setItem(appLockSecretSessionStorageKey, secret);
    sessionStorage.setItem(tokenSessionStorageKey, token);
    localStorage.setItem(tokenLocalStorageKey, encryptedToken);
}

function decryptToken() {
    const token = sessionStorage.getItem(tokenSessionStorageKey);

    localStorage.setItem(tokenLocalStorageKey, token);
    sessionStorage.removeItem(tokenSessionStorageKey);
    sessionStorage.removeItem(appLockSecretSessionStorageKey);
}

function updateToken(token) {
    if (utils.isString(token)) {
        if (settings.isEnableApplicationLock()) {
            sessionStorage.setItem(tokenSessionStorageKey, token);

            const secret = sessionStorage.getItem(appLockSecretSessionStorageKey);
            localStorage.setItem(tokenLocalStorageKey, getEncryptedToken(token, secret));
        } else {
            localStorage.setItem(tokenLocalStorageKey, token);
        }
    }
}

function updateUserInfo(user) {
    if (utils.isObject(user)) {
        localStorage.setItem(userInfoLocalStorageKey, JSON.stringify(user));
    }
}

function updateTokenAndUserInfo(item) {
    if (utils.isObject(item)) {
        if (item.token) {
            updateToken(item.token);
        }

        if (item.user) {
            updateUserInfo(item.user);
        }
    }
}

function clearTokenAndUserInfo() {
    sessionStorage.removeItem(tokenSessionStorageKey);
    sessionStorage.removeItem(appLockSecretSessionStorageKey);
    localStorage.removeItem(tokenLocalStorageKey);
    localStorage.removeItem(userInfoLocalStorageKey);
}

export default {
    getToken,
    getUserInfo,
    isUserLogined,
    isUserUnlocked,
    unlockToken,
    encryptToken,
    decryptToken,
    updateToken,
    updateUserInfo,
    updateTokenAndUserInfo,
    clearTokenAndUserInfo
};
