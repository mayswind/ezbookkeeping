import CryptoJS from 'crypto-js';

import { isString, isObject } from './common.js';
import { isEnableApplicationLock } from './settings.js';
import logger from './logger.js';

const appLockSecretBaseStringPrefix = 'EBK_LOCK_SECRET_';

const tokenLocalStorageKey = 'ebk_user_token';
const webauthnConfigLocalStorageKey = 'ebk_user_webauthn_config';
const userInfoLocalStorageKey = 'ebk_user_info';
const transactionDraftLocalStorageKey = 'ebk_user_draft_transaction';

const tokenSessionStorageKey = 'ebk_user_session_token';
const encryptedTokenSessionStorageKey = 'ebk_user_session_encrypted_token';
const appLockStateSessionStorageKey = 'ebk_user_app_lock_state'; // { 'username': '', secret: '' }

function getAppLockSecret(pinCode) {
    const hashedPinCode = CryptoJS.SHA256(appLockSecretBaseStringPrefix + pinCode).toString();
    return hashedPinCode.substring(0, 24); // put secret into user id of webauthn (user id total length must less 64 bytes)
}

function getEncryptedToken(token, appLockState) {
    const key = CryptoJS.SHA256(`${appLockSecretBaseStringPrefix}|${appLockState.username}|${appLockState.secret}`).toString();
    return CryptoJS.AES.encrypt(token, key).toString();
}

function getDecryptedToken(encryptedToken, appLockState) {
    const key = CryptoJS.SHA256(`${appLockSecretBaseStringPrefix}|${appLockState.username}|${appLockState.secret}`).toString();
    const bytes = CryptoJS.AES.decrypt(encryptedToken, key);
    return bytes.toString(CryptoJS.enc.Utf8);
}

function getToken() {
    if (isEnableApplicationLock()) {
        const usedEncryptedToken = sessionStorage.getItem(encryptedTokenSessionStorageKey);
        const currentEncryptedToken = localStorage.getItem(tokenLocalStorageKey);

        if (usedEncryptedToken === currentEncryptedToken) {
            return sessionStorage.getItem(tokenSessionStorageKey);
        }

        // re-decrypt token
        logger.warn(`encrypted token in local storage does not equal to the one in session storage, need to re-decrypt`);

        const appLockState = getUserAppLockState();
        const token = getDecryptedToken(currentEncryptedToken, appLockState);

        sessionStorage.setItem(encryptedTokenSessionStorageKey, currentEncryptedToken);
        sessionStorage.setItem(tokenSessionStorageKey, token);

        return token;
    } else {
        return localStorage.getItem(tokenLocalStorageKey);
    }
}

function getUserInfo() {
    const data = localStorage.getItem(userInfoLocalStorageKey);
    return JSON.parse(data);
}

function getUserTransactionDraft() {
    let data = localStorage.getItem(transactionDraftLocalStorageKey);

    if (!data) {
        return null;
    }

    if (isEnableApplicationLock()) {
        const appLockState = getUserAppLockState();
        data = getDecryptedToken(data, appLockState);
    }

    return JSON.parse(data);
}

function getUserAppLockState() {
    const data = sessionStorage.getItem(appLockStateSessionStorageKey);
    return JSON.parse(data);
}

function isUserLogined() {
    return !!localStorage.getItem(tokenLocalStorageKey);
}

function isUserUnlocked() {
    if (!isUserLogined()) {
        return false;
    }

    if (!isEnableApplicationLock()) {
        return true;
    }

    return !!sessionStorage.getItem(appLockStateSessionStorageKey) && !!sessionStorage.getItem(tokenSessionStorageKey);
}

function getWebAuthnCredentialId() {
    const webauthnConfigData = localStorage.getItem(webauthnConfigLocalStorageKey);
    const webauthnConfig = JSON.parse(webauthnConfigData);

    return webauthnConfig.credentialId;
}

function saveWebAuthnConfig(credentialId) {
    const webAuthnConfig = {
        credentialId: credentialId
    };

    localStorage.setItem(webauthnConfigLocalStorageKey, JSON.stringify(webAuthnConfig));
}

function clearWebAuthnConfig() {
    localStorage.removeItem(webauthnConfigLocalStorageKey);
}

function unlockTokenByWebAuthn(credentialId, userName, userSecret) {
    const webauthnConfigData = localStorage.getItem(webauthnConfigLocalStorageKey);
    const webauthnConfig = JSON.parse(webauthnConfigData);

    if (webauthnConfig.credentialId !== credentialId) {
        return false;
    }

    const encryptedToken = localStorage.getItem(tokenLocalStorageKey);
    const appLockState = {
        username: userName,
        secret: userSecret
    };
    const token = getDecryptedToken(encryptedToken, appLockState);

    sessionStorage.setItem(appLockStateSessionStorageKey, JSON.stringify(appLockState));
    sessionStorage.setItem(encryptedTokenSessionStorageKey, encryptedToken);
    sessionStorage.setItem(tokenSessionStorageKey, token);
}

function unlockTokenByPinCode(userName, pinCode) {
    const encryptedToken = localStorage.getItem(tokenLocalStorageKey);
    const appLockState = {
        username: userName,
        secret: getAppLockSecret(pinCode)
    };
    const token = getDecryptedToken(encryptedToken, appLockState);

    sessionStorage.setItem(appLockStateSessionStorageKey, JSON.stringify(appLockState));
    sessionStorage.setItem(encryptedTokenSessionStorageKey, encryptedToken);
    sessionStorage.setItem(tokenSessionStorageKey, token);
}

function encryptToken(userName, pinCode) {
    const token = localStorage.getItem(tokenLocalStorageKey);
    const appLockState = {
        username: userName,
        secret: getAppLockSecret(pinCode)
    };
    const encryptedToken = getEncryptedToken(token, appLockState);

    sessionStorage.setItem(appLockStateSessionStorageKey, JSON.stringify(appLockState));
    sessionStorage.setItem(encryptedTokenSessionStorageKey, encryptedToken);
    sessionStorage.setItem(tokenSessionStorageKey, token);
    localStorage.setItem(tokenLocalStorageKey, encryptedToken);
}

function decryptToken() {
    const token = sessionStorage.getItem(tokenSessionStorageKey);

    localStorage.setItem(tokenLocalStorageKey, token);
    sessionStorage.removeItem(tokenSessionStorageKey);
    sessionStorage.removeItem(encryptedTokenSessionStorageKey);
    sessionStorage.removeItem(appLockStateSessionStorageKey);
}

function isCorrectPinCode(pinCode) {
    const secret = getAppLockSecret(pinCode);
    const appLockState = getUserAppLockState();

    return appLockState && secret === appLockState.secret;
}

function updateToken(token) {
    if (isString(token)) {
        if (isEnableApplicationLock()) {
            const appLockState = getUserAppLockState();
            const encryptedToken = getEncryptedToken(token, appLockState);

            sessionStorage.setItem(encryptedTokenSessionStorageKey, encryptedToken);
            sessionStorage.setItem(tokenSessionStorageKey, token);
            localStorage.setItem(tokenLocalStorageKey, encryptedToken);
        } else {
            localStorage.setItem(tokenLocalStorageKey, token);
        }
    }
}

function updateUserInfo(user) {
    if (isObject(user)) {
        localStorage.setItem(userInfoLocalStorageKey, JSON.stringify(user));
    }
}

function updateUserTransactionDraft(transaction) {
    if (!isObject(transaction)) {
        return;
    }

    let data = JSON.stringify(transaction);

    if (isEnableApplicationLock()) {
        const appLockState = getUserAppLockState();
        data = getEncryptedToken(data, appLockState);
    }

    localStorage.setItem(transactionDraftLocalStorageKey, data);
}

function clearUserInfo() {
    localStorage.removeItem(userInfoLocalStorageKey);
}

function clearUserTransactionDraft() {
    localStorage.removeItem(transactionDraftLocalStorageKey);
}

function clearSessionToken() {
    sessionStorage.removeItem(tokenSessionStorageKey);
    sessionStorage.removeItem(encryptedTokenSessionStorageKey);
    sessionStorage.removeItem(appLockStateSessionStorageKey);
}

function clearTokenAndUserInfo(clearAppLockState) {
    if (clearAppLockState) {
        sessionStorage.removeItem(appLockStateSessionStorageKey);
    }

    sessionStorage.removeItem(tokenSessionStorageKey);
    sessionStorage.removeItem(encryptedTokenSessionStorageKey);
    localStorage.removeItem(tokenLocalStorageKey);
    clearUserTransactionDraft();
    clearUserInfo();
}

export default {
    getToken,
    getUserInfo,
    getUserTransactionDraft,
    getUserAppLockState,
    isUserLogined,
    isUserUnlocked,
    getWebAuthnCredentialId,
    saveWebAuthnConfig,
    clearWebAuthnConfig,
    unlockTokenByWebAuthn,
    unlockTokenByPinCode,
    encryptToken,
    decryptToken,
    isCorrectPinCode,
    updateToken,
    updateUserTransactionDraft,
    updateUserInfo,
    clearUserInfo,
    clearUserTransactionDraft,
    clearSessionToken,
    clearTokenAndUserInfo
};
