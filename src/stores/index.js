import { defineStore } from 'pinia';

import { useUserStore } from './user.js';
import { useAccountsStore } from './account.js';
import { useTransactionCategoriesStore } from './transactionCategory.js';
import { useTransactionTagsStore } from './transactionTag.js';
import { useTransactionsStore } from './transaction.js';
import { useOverviewStore } from './overview.js';
import { useStatisticsStore } from './statistics.js';
import { useExchangeRatesStore } from './exchangeRates.js';

import userState from '@/lib/userstate.js';
import services from '@/lib/services.js';
import settings from '@/lib/settings.js';
import logger from '@/lib/logger.js';
import { isObject, isString } from '@/lib/common.js';

export const useRootStore = defineStore('root', {
    actions: {
        resetAllStates() {
            const exchangeRatesStore = useExchangeRatesStore();
            exchangeRatesStore.resetLatestExchangeRates();

            const statisticsStore = useStatisticsStore();
            statisticsStore.resetTransactionStatistics();

            const overviewStore = useOverviewStore();
            overviewStore.resetTransactionOverview();

            const transactionsStore = useTransactionsStore();
            transactionsStore.resetTransactions();

            const transactionTagsStore = useTransactionTagsStore();
            transactionTagsStore.resetTransactionTags();

            const transactionCategoriesStore = useTransactionCategoriesStore();
            transactionCategoriesStore.resetTransactionCategories();

            const accountsStore = useAccountsStore();
            accountsStore.resetAccounts();

            const userStore = useUserStore();
            userStore.resetUserInfo();
        },
        authorize({ loginName, password }) {
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

                    if (data.result.user && isObject(data.result.user)) {
                        const userStore = useUserStore();
                        userStore.storeUserInfo(data.result.user);
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
        },
        authorize2FA({ token, passcode, recoveryCode }) {
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

                    if (data.result.user && isObject(data.result.user)) {
                        const userStore = useUserStore();
                        userStore.storeUserInfo(data.result.user);
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
        },
        register({ user }) {
            return new Promise((resolve, reject) => {
                services.register({
                    username: user.username,
                    password: user.password,
                    email: user.email,
                    nickname: user.nickname,
                    language: user.language,
                    defaultCurrency: user.defaultCurrency,
                    firstDayOfWeek: user.firstDayOfWeek
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

                    if (data.result.token && isString(data.result.token)) {
                        userState.updateToken(data.result.token);
                    }

                    if (data.result.user && isObject(data.result.user)) {
                        const userStore = useUserStore();
                        userStore.storeUserInfo(data.result.user);
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
        },
        logout() {
            const self = this;

            return new Promise((resolve, reject) => {
                services.logout().then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to logout' });
                        return;
                    }

                    userState.clearTokenAndUserInfo(true);
                    userState.clearWebAuthnConfig();
                    self.resetAllStates();

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
        },
        forceLogout() {
            userState.clearTokenAndUserInfo(true);
            userState.clearWebAuthnConfig();
            this.resetAllStates();
        },
        updateUserProfile({ profile, currentPassword }) {
            return new Promise((resolve, reject) => {
                services.updateProfile({
                    password: profile.password,
                    oldPassword: currentPassword,
                    email: profile.email,
                    nickname: profile.nickname,
                    defaultAccountId: profile.defaultAccountId,
                    transactionEditScope: profile.transactionEditScope,
                    language: profile.language,
                    defaultCurrency: profile.defaultCurrency,
                    firstDayOfWeek: profile.firstDayOfWeek,
                    longDateFormat: profile.longDateFormat,
                    shortDateFormat: profile.shortDateFormat,
                    longTimeFormat: profile.longTimeFormat,
                    shortTimeFormat: profile.shortTimeFormat
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to update user profile' });
                        return;
                    }

                    if (data.result.newToken && isString(data.result.newToken)) {
                        userState.updateToken(data.result.newToken);
                    }

                    if (data.result.user && isObject(data.result.user)) {
                        const userStore = useUserStore();
                        userStore.storeUserInfo(data.result.user);
                    }

                    const accountsStore = useAccountsStore();
                    if (!accountsStore.accountListStateInvalid) {
                        accountsStore.updateAccountListInvalidState(true);
                    }

                    const overviewStore = useOverviewStore();
                    if (!overviewStore.transactionOverviewStateInvalid) {
                        overviewStore.updateTransactionOverviewInvalidState(true);
                    }

                    const statisticsStore = useStatisticsStore();
                    if (!statisticsStore.transactionStatisticsStateInvalid) {
                        statisticsStore.updateTransactionStatisticsInvalidState(true);
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
        },
        clearUserData({ password }) {
            return new Promise((resolve, reject) => {
                services.clearData({
                    password: password
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to clear user data' });
                        return;
                    }

                    const accountsStore = useAccountsStore();
                    if (!accountsStore.accountListStateInvalid) {
                        accountsStore.updateAccountListInvalidState(true);
                    }

                    const transactionCategoriesStore = useTransactionCategoriesStore();
                    if (!transactionCategoriesStore.transactionCategoryListStateInvalid) {
                        transactionCategoriesStore.updateTransactionCategoryListInvalidState(true);
                    }

                    const transactionTagsStore = useTransactionTagsStore();
                    if (!transactionTagsStore.transactionTagListStateInvalid) {
                        transactionTagsStore.updateTransactionTagListInvalidState(true);
                    }

                    const overviewStore = useOverviewStore();
                    if (!overviewStore.transactionOverviewStateInvalid) {
                        overviewStore.updateTransactionOverviewInvalidState(true);
                    }

                    const statisticsStore = useStatisticsStore();
                    if (!statisticsStore.transactionStatisticsStateInvalid) {
                        statisticsStore.updateTransactionStatisticsInvalidState(true);
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
    }
});
