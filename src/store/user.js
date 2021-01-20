import userState from '../lib/userstate.js';
import services from '../lib/services.js';
import settings from '../lib/settings.js';
import logger from '../lib/logger.js';
import utils from '../lib/utils.js';

import {
    RESET_STATE,

    STORE_USER_INFO,
    CLEAR_USER_INFO,

    UPDATE_ACCOUNT_LIST_INVALID_STATE,
    UPDATE_TRANSACTION_CATEGORY_LIST_INVALID_STATE,
    UPDATE_TRANSACTION_TAG_LIST_INVALID_STATE,
    UPDATE_TRANSACTION_OVERVIEW_INVALID_STATE,
    UPDATE_TRANSACTION_STATISTICS_INVALID_STATE
} from './mutations.js';

export function authorize(context, { loginName, password }) {
    return new Promise((resolve, reject) => {
        services.authorize({
            loginName: loginName,
            password: password
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result || !data.result.token) {
                reject({ message: 'Unable to login' });
                return;
            }

            if (data.result.need2FA) {
                resolve(data.result);
                return;
            }

            if (settings.isEnableApplicationLock() || userState.getUserAppLockState()) {
                const appLockState = userState.getUserAppLockState();

                if (!appLockState || appLockState.username !== data.result.user.username) {
                    userState.clearTokenAndUserInfo(true);
                    settings.setEnableApplicationLock(false);
                    settings.setEnableApplicationLockWebAuthn(false);
                    userState.clearWebAuthnConfig();
                }
            }

            userState.updateToken(data.result.token);

            if (data.result.user && utils.isObject(data.result.user)) {
                context.commit(STORE_USER_INFO, data.result.user);
            }

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to login', error);

            if (error && error.processed) {
                reject(error);
            } else if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else {
                reject({ message: 'Unable to login' });
            }
        });
    });
}

export function authorize2FA(context, { token, passcode, recoveryCode }) {
    return new Promise((resolve, reject) => {
        let promise = null;

        if (passcode) {
            promise = services.authorize2FA({
                passcode: passcode,
                token: token
            });
        } else if (recoveryCode) {
            promise = services.authorize2FAByBackupCode({
                recoveryCode: recoveryCode,
                token: token
            });
        } else {
            reject({ message: 'An error has occurred' });
            return;
        }

        promise.then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result || !data.result.token) {
                reject({ message: 'Unable to verify' });
                return;
            }

            if (settings.isEnableApplicationLock() || userState.getUserAppLockState()) {
                const appLockState = userState.getUserAppLockState();

                if (!appLockState || appLockState.username !== data.result.user.username) {
                    userState.clearTokenAndUserInfo(true);
                    settings.setEnableApplicationLock(false);
                    settings.setEnableApplicationLockWebAuthn(false);
                    userState.clearWebAuthnConfig();
                }
            }

            userState.updateToken(data.result.token);

            if (data.result.user && utils.isObject(data.result.user)) {
                context.commit(STORE_USER_INFO, data.result.user);
            }

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to verify 2fa', error);

            if (error && error.processed) {
                reject(error);
            } else if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else {
                reject({ message: 'Unable to verify' });
            }
        });
    });
}

export function register(context, { user }) {
    return new Promise((resolve, reject) => {
        services.register({
            username: user.username,
            password: user.password,
            email: user.email,
            nickname: user.nickname,
            defaultCurrency: user.defaultCurrency
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                reject({ message: 'Unable to sign up' });
                return;
            }

            if (settings.isEnableApplicationLock()) {
                settings.setEnableApplicationLock(false);
                settings.setEnableApplicationLockWebAuthn(false);
                userState.clearWebAuthnConfig();
            }

            if (data.result.token && utils.isString(data.result.token)) {
                userState.updateToken(data.result.token);
            }

            if (data.result.user && utils.isObject(data.result.user)) {
                context.commit(STORE_USER_INFO, data.result.user);
            }

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to sign up', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                reject({ message: 'Unable to sign up' });
            } else {
                reject(error);
            }
        });
    });
}

export function logout(context) {
    return new Promise((resolve, reject) => {
        services.logout().then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                reject({ message: 'Unable to logout' });
                return;
            }

            context.commit(CLEAR_USER_INFO);
            userState.clearTokenAndUserInfo(true);
            userState.clearWebAuthnConfig();

            context.commit(RESET_STATE);

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to log out', error);

            if (error && error.processed) {
                reject(error);
            } else if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else {
                reject({ message: 'Unable to logout' });
            }
        });
    });
}

export function getCurrentUserProfile() {
    return new Promise((resolve, reject) => {
        services.getProfile().then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                reject({ message: 'Unable to get user profile' });
                return;
            }

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to get user profile', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                reject({ message: 'Unable to get user profile' });
            } else {
                reject(error);
            }
        });
    });
}

export function updateUserProfile(context, { profile, currentPassword }) {
    return new Promise((resolve, reject) => {
        services.updateProfile({
            password: profile.password,
            oldPassword: currentPassword,
            email: profile.email,
            nickname: profile.nickname,
            defaultCurrency: profile.defaultCurrency
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                reject({ message: 'Unable to update user profile' });
                return;
            }

            if (data.result.newToken && utils.isString(data.result.newToken)) {
                userState.updateToken(data.result.newToken);
            }

            if (data.result.user && utils.isObject(data.result.user)) {
                context.commit(STORE_USER_INFO, data.result.user);
            }

            if (!context.state.accountListStateInvalid) {
                context.commit(UPDATE_ACCOUNT_LIST_INVALID_STATE, true);
            }

            if (!context.state.transactionOverviewStateInvalid) {
                context.commit(UPDATE_TRANSACTION_OVERVIEW_INVALID_STATE, true);
            }

            if (!context.state.transactionStatisticsStateInvalid) {
                context.commit(UPDATE_TRANSACTION_STATISTICS_INVALID_STATE, true);
            }

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to save user profile', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                reject({ message: 'Unable to update user profile' });
            } else {
                reject(error);
            }
        });
    });
}

export function clearUserData(context, { password }) {
    return new Promise((resolve, reject) => {
        services.clearData({
            password: password
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                reject({ message: 'Unable to clear user data' });
                return;
            }

            if (!context.state.accountListStateInvalid) {
                context.commit(UPDATE_ACCOUNT_LIST_INVALID_STATE, true);
            }

            if (!context.state.transactionCategoryListStateInvalid) {
                context.commit(UPDATE_TRANSACTION_CATEGORY_LIST_INVALID_STATE, true);
            }

            if (!context.state.transactionTagListStateInvalid) {
                context.commit(UPDATE_TRANSACTION_TAG_LIST_INVALID_STATE, true);
            }

            if (!context.state.transactionOverviewStateInvalid) {
                context.commit(UPDATE_TRANSACTION_OVERVIEW_INVALID_STATE, true);
            }

            if (!context.state.transactionStatisticsStateInvalid) {
                context.commit(UPDATE_TRANSACTION_STATISTICS_INVALID_STATE, true);
            }

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to clear user data', error);

            if (error && error.processed) {
                reject(error);
            } else if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else {
                reject({ message: 'Unable to clear user data' });
            }
        });
    });
}

export function clearUserInfoState(context) {
    context.commit(CLEAR_USER_INFO);
}

export function resetState(context) {
    context.commit(RESET_STATE);
}

export function currentUserNickname(state) {
    const userInfo = state.currentUserInfo || {};
    return userInfo.nickname || userInfo.username || null;
}

export function currentUserDefaultCurrency(state) {
    const userInfo = state.currentUserInfo || {};
    return userInfo.defaultCurrency || null;
}
