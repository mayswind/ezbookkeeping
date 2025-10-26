import { ref } from 'vue';
import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.ts';
import { useUserStore } from './user.ts';
import { useAccountsStore } from './account.ts';
import { useTransactionCategoriesStore } from './transactionCategory.ts';
import { useTransactionTagsStore } from './transactionTag.ts';
import { useTransactionTemplatesStore } from './transactionTemplate.ts';
import { useTransactionsStore } from './transaction.ts';
import { useOverviewStore } from './overview.ts';
import { useStatisticsStore } from './statistics.ts';
import { useExchangeRatesStore } from './exchangeRates.ts';

import type { AuthResponse, RegisterResponse } from '@/models/auth_response.ts';
import type {
    User,
    UserLoginRequest,
    UserResendVerifyEmailRequest,
    UserVerifyEmailResponse,
    UserProfileUpdateRequest,
    UserProfileUpdateResponse
} from '@/models/user.ts';
import type { ForgetPasswordRequest } from '@/models/forget_password.ts';
import type { LocalizedPresetCategory } from '@/core/category.ts';

import {
    isObject,
    isString
} from '@/lib/common.ts';
import {
    hasUserAppLockState,
    getUserAppLockState,
    updateCurrentToken,
    clearWebAuthnConfig,
    clearCurrentSessionToken,
    clearCurrentTokenAndUserInfo
} from '@/lib/userstate.ts';
import services, { type ApiResponsePromise } from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

export const useRootStore = defineStore('root', () => {
    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const accountsStore = useAccountsStore();
    const transactionCategoriesStore = useTransactionCategoriesStore();
    const transactionTagsStore = useTransactionTagsStore();
    const transactionTemplatesStore = useTransactionTemplatesStore();
    const transactionsStore = useTransactionsStore();
    const overviewStore = useOverviewStore();
    const statisticsStore = useStatisticsStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const currentNotification = ref<string | null>(null);

    function resetAllStates(resetUserInfoAndSettings: boolean): void {
        if (resetUserInfoAndSettings) {
            exchangeRatesStore.resetLatestExchangeRates();
        }

        setNotificationContent(null);

        statisticsStore.resetTransactionStatistics();
        overviewStore.resetTransactionOverview();
        transactionsStore.resetTransactions();
        transactionTagsStore.resetTransactionTags();
        transactionCategoriesStore.resetTransactionCategories();
        transactionTemplatesStore.resetTransactionTemplates();
        accountsStore.resetAccounts();

        if (resetUserInfoAndSettings) {
            userStore.resetUserBasicInfo();
        }
    }

    function setNotificationContent(content: string | null): void {
        currentNotification.value = content;
    }

    function generateOAuth2LoginUrl(platform: 'mobile' | 'desktop', clientSessionId: string): string {
        return services.generateOAuth2LoginUrl(platform, clientSessionId);
    }

    function authorize(req: UserLoginRequest): Promise<AuthResponse> {
        return new Promise((resolve, reject) => {
            services.authorize(req).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result || !data.result.token) {
                    reject({ message: 'Unable to log in' });
                    return;
                }

                if (data.result.need2FA) {
                    resolve(data.result);
                    return;
                }

                if (settingsStore.appSettings.applicationLock || hasUserAppLockState()) {
                    const appLockState = getUserAppLockState();

                    if (!appLockState || appLockState.username !== data.result.user?.username) {
                        clearCurrentTokenAndUserInfo(true);
                        settingsStore.setEnableApplicationLock(false);
                        settingsStore.setEnableApplicationLockWebAuthn(false);
                        clearWebAuthnConfig();
                    }
                }

                settingsStore.setApplicationSettingsFromCloudSettings(data.result.applicationCloudSettings);

                updateCurrentToken(data.result.token);

                if (data.result.user && isObject(data.result.user)) {
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
    }

    function authorize2FA({ token, passcode, recoveryCode }: { token: string, passcode: string | null, recoveryCode: string | null }): Promise<AuthResponse> {
        return new Promise((resolve, reject) => {
            let promise: ApiResponsePromise<AuthResponse>;

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

                if (settingsStore.appSettings.applicationLock || hasUserAppLockState()) {
                    const appLockState = getUserAppLockState();

                    if (!appLockState || appLockState.username !== data.result.user?.username) {
                        clearCurrentTokenAndUserInfo(true);
                        settingsStore.setEnableApplicationLock(false);
                        settingsStore.setEnableApplicationLockWebAuthn(false);
                        clearWebAuthnConfig();
                    }
                }

                settingsStore.setApplicationSettingsFromCloudSettings(data.result.applicationCloudSettings);

                updateCurrentToken(data.result.token);

                if (data.result.user && isObject(data.result.user)) {
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
    }

    function authorizeOAuth2({ password, passcode, token }: { password?: string, passcode?: string, token: string }): Promise<AuthResponse> {
        return new Promise((resolve, reject) => {
            services.authorizeOAuth2({
                req: {
                    password,
                    passcode
                },
                token
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result || !data.result.token) {
                    reject({ message: 'Unable to log in' });
                    return;
                }

                if (settingsStore.appSettings.applicationLock || hasUserAppLockState()) {
                    const appLockState = getUserAppLockState();

                    if (!appLockState || appLockState.username !== data.result.user?.username) {
                        clearCurrentTokenAndUserInfo(true);
                        settingsStore.setEnableApplicationLock(false);
                        settingsStore.setEnableApplicationLockWebAuthn(false);
                        clearWebAuthnConfig();
                    }
                }

                settingsStore.setApplicationSettingsFromCloudSettings(data.result.applicationCloudSettings);

                updateCurrentToken(data.result.token);

                if (data.result.user && isObject(data.result.user)) {
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
    }

    function register({ user, presetCategories }: { user: User, presetCategories?: LocalizedPresetCategory[] }): Promise<RegisterResponse> {
        return new Promise((resolve, reject) => {
            services.register(user.toRegisterRequest(presetCategories)).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to sign up' });
                    return;
                }

                if (settingsStore.appSettings.applicationLock) {
                    settingsStore.setEnableApplicationLock(false);
                    settingsStore.setEnableApplicationLockWebAuthn(false);
                    clearWebAuthnConfig();
                }

                if (data.result.token && isString(data.result.token)) {
                    updateCurrentToken(data.result.token);
                }

                if (data.result.user && isObject(data.result.user)) {
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
    }

    function lock(): void {
        clearCurrentSessionToken();
        resetAllStates(false);
    }

    function logout(): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.logout().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to logout' });
                    return;
                }

                clearCurrentTokenAndUserInfo(true);
                clearWebAuthnConfig();
                resetAllStates(true);

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

    function forceLogout(): void {
        clearCurrentTokenAndUserInfo(true);
        clearWebAuthnConfig();
        resetAllStates(true);
    }

    function verifyEmail({ token, requestNewToken }: { token: string, requestNewToken: boolean }): Promise<UserVerifyEmailResponse> {
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
                    updateCurrentToken(data.result.newToken);
                }

                if (data.result.user && isObject(data.result.user)) {
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
    }

    function resendVerifyEmailByUnloginUser(req: UserResendVerifyEmailRequest): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.resendVerifyEmailByUnloginUser(req).then(response => {
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
    }

    function requestResetPassword(req: ForgetPasswordRequest): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.requestResetPassword(req).then(response => {
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
    }

    function resetPassword({ email, token, password }: { email: string, token: string, password: string }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.resetPassword({
                email,
                token,
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
    }

    function updateUserProfile(req: UserProfileUpdateRequest): Promise<UserProfileUpdateResponse> {
        const userDefaultCurrency = userStore.currentUserDefaultCurrency;

        return new Promise((resolve, reject) => {
            services.updateProfile(req).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to update user profile' });
                    return;
                }

                if (data.result.newToken && isString(data.result.newToken)) {
                    updateCurrentToken(data.result.newToken);
                }

                if (data.result.user && isObject(data.result.user)) {
                    userStore.storeUserBasicInfo(data.result.user);
                }

                if (!accountsStore.accountListStateInvalid) {
                    accountsStore.updateAccountListInvalidState(true);
                }

                if (!overviewStore.transactionOverviewStateInvalid) {
                    overviewStore.updateTransactionOverviewInvalidState(true);
                }

                if (!statisticsStore.transactionStatisticsStateInvalid) {
                    statisticsStore.updateTransactionStatisticsInvalidState(true);
                }

                if (data.result.user && data.result.user.defaultCurrency !== userDefaultCurrency) {
                    exchangeRatesStore.resetLatestExchangeRates();
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

    function resendVerifyEmailByLoginedUser(): Promise<boolean> {
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
    }

    function clearAllUserTransactionsOfAccount({ accountId, password }: { accountId: string, password: string }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.clearAllTransactionsOfAccount({
                accountId: accountId,
                password: password
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to clear user data' });
                    return;
                }

                if (!accountsStore.accountListStateInvalid) {
                    accountsStore.updateAccountListInvalidState(true);
                }

                if (!overviewStore.transactionOverviewStateInvalid) {
                    overviewStore.updateTransactionOverviewInvalidState(true);
                }

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

    function clearAllUserTransactions({ password }: { password: string }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.clearAllTransactions({
                password: password
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to clear user data' });
                    return;
                }

                if (!accountsStore.accountListStateInvalid) {
                    accountsStore.updateAccountListInvalidState(true);
                }

                if (!overviewStore.transactionOverviewStateInvalid) {
                    overviewStore.updateTransactionOverviewInvalidState(true);
                }

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

    function clearAllUserData({ password }: { password: string }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.clearAllData({
                password: password
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to clear user data' });
                    return;
                }

                if (!accountsStore.accountListStateInvalid) {
                    accountsStore.updateAccountListInvalidState(true);
                }

                if (!transactionCategoriesStore.transactionCategoryListStateInvalid) {
                    transactionCategoriesStore.updateTransactionCategoryListInvalidState(true);
                }

                if (!transactionTagsStore.transactionTagListStateInvalid) {
                    transactionTagsStore.updateTransactionTagListInvalidState(true);
                }

                if (!overviewStore.transactionOverviewStateInvalid) {
                    overviewStore.updateTransactionOverviewInvalidState(true);
                }

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

    return {
        // states
        currentNotification,
        // functions
        setNotificationContent,
        generateOAuth2LoginUrl,
        authorize,
        authorize2FA,
        authorizeOAuth2,
        register,
        lock,
        logout,
        forceLogout,
        verifyEmail,
        resendVerifyEmailByUnloginUser,
        requestResetPassword,
        resetPassword,
        updateUserProfile,
        resendVerifyEmailByLoginedUser,
        clearAllUserTransactionsOfAccount,
        clearAllUserTransactions,
        clearAllUserData
    };
});
