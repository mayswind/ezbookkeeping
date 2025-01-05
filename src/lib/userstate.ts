import CryptoJS from 'crypto-js';

import type { ApplicationLockState, WebAuthnConfig } from '@/core/setting.ts';
import type { UserBasicInfo } from '@/models/user.ts';

import { isString, isObject } from './common.ts';
import { isEnableApplicationLock } from './settings.ts';
import logger from './logger.ts';

const appLockSecretBaseStringPrefix: string = 'EBK_LOCK_SECRET_';

const tokenLocalStorageKey: string = 'ebk_user_token';
const webauthnConfigLocalStorageKey: string = 'ebk_user_webauthn_config';
const userInfoLocalStorageKey: string = 'ebk_user_info';
const transactionDraftLocalStorageKey: string = 'ebk_user_draft_transaction';

const tokenSessionStorageKey: string = 'ebk_user_session_token';
const encryptedTokenSessionStorageKey: string = 'ebk_user_session_encrypted_token';
const appLockStateSessionStorageKey: string = 'ebk_user_app_lock_state'; // { 'username': '', secret: '' }

function getAppLockSecret(pinCode: string): string {
    const hashedPinCode = CryptoJS.SHA256(appLockSecretBaseStringPrefix + pinCode).toString();
    return hashedPinCode.substring(0, 24); // put secret into user id of webauthn (user id total length must less 64 bytes)
}

function getEncryptedToken(token: string, appLockState: ApplicationLockState): string {
    const key = CryptoJS.SHA256(`${appLockSecretBaseStringPrefix}|${appLockState.username}|${appLockState.secret}`).toString();
    return CryptoJS.AES.encrypt(token, key).toString();
}

function getDecryptedToken(encryptedToken: string, appLockState: ApplicationLockState): string {
    const key = CryptoJS.SHA256(`${appLockSecretBaseStringPrefix}|${appLockState.username}|${appLockState.secret}`).toString();
    const bytes = CryptoJS.AES.decrypt(encryptedToken, key);
    return bytes.toString(CryptoJS.enc.Utf8);
}

function getToken(): string | null {
    if (isEnableApplicationLock()) {
        const usedEncryptedToken = sessionStorage.getItem(encryptedTokenSessionStorageKey);
        const currentEncryptedToken = localStorage.getItem(tokenLocalStorageKey);

        if (!usedEncryptedToken || !currentEncryptedToken) {
            return null;
        }

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

function getUserInfo(): UserBasicInfo | null {
    const data = localStorage.getItem(userInfoLocalStorageKey);

    if (!data) {
        return null;
    }

    return JSON.parse(data) as UserBasicInfo;
}

function getUserTransactionDraft(): unknown | null {
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

function getUserAppLockState(): ApplicationLockState {
    const data = sessionStorage.getItem(appLockStateSessionStorageKey);

    if (!data) {
        throw new Error('No app lock state in session storage');
    }

    const appLockState = JSON.parse(data);

    if (!appLockState || !appLockState.username || !appLockState.secret) {
        throw new Error('App lock state is invalid');
    }

    return appLockState as ApplicationLockState;
}

function isUserLogined(): boolean {
    return !!localStorage.getItem(tokenLocalStorageKey);
}

function isUserUnlocked(): boolean {
    if (!isUserLogined()) {
        return false;
    }

    if (!isEnableApplicationLock()) {
        return true;
    }

    return !!sessionStorage.getItem(appLockStateSessionStorageKey) && !!sessionStorage.getItem(tokenSessionStorageKey);
}

function getWebAuthnCredentialId(): string | undefined {
    const webauthnConfigData = localStorage.getItem(webauthnConfigLocalStorageKey);

    if (!webauthnConfigData) {
        return undefined;
    }

    const webauthnConfig = JSON.parse(webauthnConfigData) as WebAuthnConfig;

    return webauthnConfig.credentialId;
}

function saveWebAuthnConfig(credentialId: string): void {
    const webAuthnConfig: WebAuthnConfig = {
        credentialId: credentialId
    };

    localStorage.setItem(webauthnConfigLocalStorageKey, JSON.stringify(webAuthnConfig));
}

function clearWebAuthnConfig(): void {
    localStorage.removeItem(webauthnConfigLocalStorageKey);
}

function unlockTokenByWebAuthn(credentialId: string, userName: string, userSecret: string): void {
    const webauthnConfigData = localStorage.getItem(webauthnConfigLocalStorageKey);

    if (!webauthnConfigData) {
        throw new Error('WebAuthn credential is not set');
    }

    const webauthnConfig = JSON.parse(webauthnConfigData) as WebAuthnConfig;

    if (webauthnConfig.credentialId !== credentialId) {
        throw new Error('WebAuthn credential is invalid');
    }

    const encryptedToken = localStorage.getItem(tokenLocalStorageKey);

    if (!encryptedToken) {
        throw new Error('No token in local storage');
    }

    const appLockState: ApplicationLockState = {
        username: userName,
        secret: userSecret
    };
    const token = getDecryptedToken(encryptedToken, appLockState);

    sessionStorage.setItem(appLockStateSessionStorageKey, JSON.stringify(appLockState));
    sessionStorage.setItem(encryptedTokenSessionStorageKey, encryptedToken);
    sessionStorage.setItem(tokenSessionStorageKey, token);
}

function unlockTokenByPinCode(userName: string, pinCode: string): void {
    const encryptedToken = localStorage.getItem(tokenLocalStorageKey);

    if (!encryptedToken) {
        throw new Error('No token in local storage');
    }

    const appLockState: ApplicationLockState = {
        username: userName,
        secret: getAppLockSecret(pinCode)
    };
    const token = getDecryptedToken(encryptedToken, appLockState);

    sessionStorage.setItem(appLockStateSessionStorageKey, JSON.stringify(appLockState));
    sessionStorage.setItem(encryptedTokenSessionStorageKey, encryptedToken);
    sessionStorage.setItem(tokenSessionStorageKey, token);
}

function encryptToken(userName: string, pinCode: string): void {
    const token = localStorage.getItem(tokenLocalStorageKey);

    if (!token) {
        throw new Error('No token in local storage');
    }

    const appLockState: ApplicationLockState = {
        username: userName,
        secret: getAppLockSecret(pinCode)
    };
    const encryptedToken = getEncryptedToken(token, appLockState);

    sessionStorage.setItem(appLockStateSessionStorageKey, JSON.stringify(appLockState));
    sessionStorage.setItem(encryptedTokenSessionStorageKey, encryptedToken);
    sessionStorage.setItem(tokenSessionStorageKey, token);
    localStorage.setItem(tokenLocalStorageKey, encryptedToken);
}

function decryptToken(): void {
    const token = sessionStorage.getItem(tokenSessionStorageKey);

    if (!token) {
        throw new Error('No token in session storage');
    }

    localStorage.setItem(tokenLocalStorageKey, token);
    sessionStorage.removeItem(tokenSessionStorageKey);
    sessionStorage.removeItem(encryptedTokenSessionStorageKey);
    sessionStorage.removeItem(appLockStateSessionStorageKey);
}

function isCorrectPinCode(pinCode: string): boolean {
    const secret = getAppLockSecret(pinCode);
    const appLockState = getUserAppLockState();

    return appLockState && secret === appLockState.secret;
}

function updateToken(token: string): void {
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

function updateUserInfo(user: UserBasicInfo): void {
    if (isObject(user)) {
        localStorage.setItem(userInfoLocalStorageKey, JSON.stringify(user));
    }
}

function updateUserTransactionDraft(transaction: unknown): void {
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

function clearUserInfo(): void {
    localStorage.removeItem(userInfoLocalStorageKey);
}

function clearUserTransactionDraft(): void {
    localStorage.removeItem(transactionDraftLocalStorageKey);
}

function clearSessionToken(): void {
    sessionStorage.removeItem(tokenSessionStorageKey);
    sessionStorage.removeItem(encryptedTokenSessionStorageKey);
    sessionStorage.removeItem(appLockStateSessionStorageKey);
}

function clearTokenAndUserInfo(clearAppLockState: boolean): void {
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
