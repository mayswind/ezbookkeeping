import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.js';
import { useUserStore } from './user.js';
import { useAccountsStore } from './account.js';
import { useTransactionCategoriesStore } from './transactionCategory.js';
import { useTransactionTagsStore } from './transactionTag.js';
import { useTransactionTemplatesStore } from './transactionTemplate.js';
import { useTransactionsStore } from './transaction.js';
import { useOverviewStore } from './overview.js';
import { useStatisticsStore } from './statistics.js';
import { useExchangeRatesStore } from './exchangeRates.js';

import userState from '@/lib/userstate.ts';
import services from '@/lib/services.js';
import logger from '@/lib/logger.ts';
import { isObject, isString } from '@/lib/common.ts';

export const useRootStore = defineStore('root', {
    state: () => ({
        currentNotification: null
    }),
    actions: {
        resetAllStates(resetUserInfoAndSettings) {
            if (resetUserInfoAndSettings) {
                const exchangeRatesStore = useExchangeRatesStore();
                exchangeRatesStore.resetLatestExchangeRates();
            }

            this.setNotificationContent(null);

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

            const transactionTemplatesStore = useTransactionTemplatesStore();
            transactionTemplatesStore.resetTransactionTemplates();

            const accountsStore = useAccountsStore();
            accountsStore.resetAccounts();

            if (resetUserInfoAndSettings) {
                const userStore = useUserStore();
                userStore.resetUserBasicInfo();
            }
        },
        authorize({ loginName, password }) {
            const settingsStore = useSettingsStore();

            return new Promise((resolve, reject) => {
                services.authorize({
                    loginName: loginName,
                    password: password
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result || !data.result.token) {
                        reject({ message: 'Unable to log in' });
                        return;
                    }

                    if (data.result.need2FA) {
                        resolve(data.result);
                        return;
                    }

                    if (settingsStore.appSettings.applicationLock || userState.getUserAppLockState()) {
                        const appLockState = userState.getUserAppLockState();

                        if (!appLockState || appLockState.username !== data.result.user.username) {
                            userState.clearTokenAndUserInfo(true);
                            settingsStore.setEnableApplicationLock(false);
                            settingsStore.setEnableApplicationLockWebAuthn(false);
                            userState.clearWebAuthnConfig();
                        }
                    }

                    userState.updateToken(data.result.token);

                    if (data.result.user && isObject(data.result.user)) {
                        const userStore = useUserStore();
                        userStore.storeUserBasicInfo(data.result.user);
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to login', error);

                    if (error && error.processed) {
                        reject(error);
                    } else if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else {
                        reject({ message: 'Unable to log in' });
                    }
                });
            });
        },
        authorize2FA({ token, passcode, recoveryCode }) {
            const settingsStore = useSettingsStore();

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
                    reject({ message: 'An error occurred' });
                    return;
                }

                promise.then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result || !data.result.token) {
                        reject({ message: 'Unable to verify' });
                        return;
                    }

                    if (settingsStore.appSettings.applicationLock || userState.getUserAppLockState()) {
                        const appLockState = userState.getUserAppLockState();

                        if (!appLockState || appLockState.username !== data.result.user.username) {
                            userState.clearTokenAndUserInfo(true);
                            settingsStore.setEnableApplicationLock(false);
                            settingsStore.setEnableApplicationLockWebAuthn(false);
                            userState.clearWebAuthnConfig();
                        }
                    }

                    userState.updateToken(data.result.token);

                    if (data.result.user && isObject(data.result.user)) {
                        const userStore = useUserStore();
                        userStore.storeUserBasicInfo(data.result.user);
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
        register({ user, presetCategories }) {
            const settingsStore = useSettingsStore();

            return new Promise((resolve, reject) => {
                services.register({
                    username: user.username,
                    password: user.password,
                    email: user.email,
                    nickname: user.nickname,
                    language: user.language,
                    defaultCurrency: user.defaultCurrency,
                    firstDayOfWeek: user.firstDayOfWeek,
                    categories: presetCategories
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to sign up' });
                        return;
                    }

                    if (settingsStore.appSettings.applicationLock) {
                        settingsStore.setEnableApplicationLock(false);
                        settingsStore.setEnableApplicationLockWebAuthn(false);
                        userState.clearWebAuthnConfig();
                    }

                    if (data.result.token && isString(data.result.token)) {
                        userState.updateToken(data.result.token);
                    }

                    if (data.result.user && isObject(data.result.user)) {
                        const userStore = useUserStore();
                        userStore.storeUserBasicInfo(data.result.user);
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
        lock() {
            userState.clearSessionToken();
            this.resetAllStates(false);
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
                    self.resetAllStates(true);

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
            this.resetAllStates(true);
        },
        verifyEmail({ token, requestNewToken }) {
            return new Promise((resolve, reject) => {
                services.verifyEmail({
                    token,
                    requestNewToken
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to verify email' });
                        return;
                    }

                    if (data.result.newToken && isString(data.result.newToken)) {
                        userState.updateToken(data.result.newToken);
                    }

                    if (data.result.user && isObject(data.result.user)) {
                        const userStore = useUserStore();
                        userStore.storeUserBasicInfo(data.result.user);
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to verify email', error);

                    if (error && error.processed) {
                        reject(error);
                    } else if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else {
                        reject({ message: 'Unable to verify email' });
                    }
                });
            });
        },
        resendVerifyEmailByUnloginUser({ email, password }) {
            return new Promise((resolve, reject) => {
                services.resendVerifyEmailByUnloginUser({
                    email,
                    password
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to resend validation email' });
                        return;
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to resend verify email', error);

                    if (error && error.processed) {
                        reject(error);
                    } else if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else {
                        reject({ message: 'Unable to resend validation email' });
                    }
                });
            });
        },
        requestResetPassword({ email }) {
            return new Promise((resolve, reject) => {
                services.requestResetPassword({
                    email
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to send password reset email' });
                        return;
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to send password reset email', error);

                    if (error && error.processed) {
                        reject(error);
                    } else if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else {
                        reject({ message: 'Unable to send password reset email' });
                    }
                });
            });
        },
        resetPassword({ email, token, password }) {
            return new Promise((resolve, reject) => {
                services.resetPassword({
                    token,
                    email,
                    password
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to reset password' });
                        return;
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to reset password', error);

                    if (error && error.processed) {
                        reject(error);
                    } else if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else {
                        reject({ message: 'Unable to reset password' });
                    }
                });
            });
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
                    shortTimeFormat: profile.shortTimeFormat,
                    decimalSeparator: profile.decimalSeparator,
                    digitGroupingSymbol: profile.digitGroupingSymbol,
                    digitGrouping: profile.digitGrouping,
                    currencyDisplayType: profile.currencyDisplayType,
                    expenseAmountColor: profile.expenseAmountColor,
                    incomeAmountColor: profile.incomeAmountColor
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
                        userStore.storeUserBasicInfo(data.result.user);
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
        resendVerifyEmailByLoginedUser() {
            return new Promise((resolve, reject) => {
                services.resendVerifyEmailByLoginedUser().then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to resend validation email' });
                        return;
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to resend verify email', error);

                    if (error && error.processed) {
                        reject(error);
                    } else if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else {
                        reject({ message: 'Unable to resend validation email' });
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
        },
        setNotificationContent(content) {
            this.currentNotification = content;
        }
    }
});
