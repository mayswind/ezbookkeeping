import { type Ref, ref, computed } from 'vue';
import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.ts';

import { WeekDay } from '@/core/datetime.ts';

import {
    type UserBasicInfo,
    type UserProfileResponse,
    type UserProfileUpdateResponse,
    User,
    EMPTY_USER_BASIC_INFO
} from '@/models/user.ts';

import type {
    DataStatisticsResponse
} from '@/models/data_management.ts';

import {
    isObject,
    isNumber
} from '@/lib/common.ts';

import {
    getCurrentUserInfo,
    updateCurrentUserInfo,
    clearCurrentUserInfo
} from '@/lib/userstate.ts';

import logger from '@/lib/logger.ts';
import services from '@/lib/services.ts';

export const useUserStore = defineStore('user', () => {
    const settingsStore = useSettingsStore();
    const currentUserBasicInfo: Ref<UserBasicInfo | null> = ref(getCurrentUserInfo());

    const currentUserNickname = computed<string | null>(() => {
        const userInfo = currentUserBasicInfo.value || EMPTY_USER_BASIC_INFO;
        return userInfo.nickname || userInfo.username || null;
    });

    const currentUserAvatar = computed<string | null>(() => {
        const userInfo = currentUserBasicInfo.value || EMPTY_USER_BASIC_INFO;
        return getUserAvatarUrl(userInfo, false);
    });

    const currentUserDefaultAccountId = computed<string>(() => {
        const userInfo = currentUserBasicInfo.value || EMPTY_USER_BASIC_INFO;
        return userInfo.defaultAccountId;
    });

    const currentUserLanguage = computed<string>(() => {
        const userInfo = currentUserBasicInfo.value || EMPTY_USER_BASIC_INFO;
        return userInfo.language;
    });

    const currentUserDefaultCurrency = computed<string>(() => {
        const userInfo = currentUserBasicInfo.value || EMPTY_USER_BASIC_INFO;
        return userInfo.defaultCurrency || settingsStore.localeDefaultSettings.currency;
    });

    const currentUserFirstDayOfWeek = computed<number>(() => {
        const userInfo = currentUserBasicInfo.value || EMPTY_USER_BASIC_INFO;
        return isNumber(userInfo.firstDayOfWeek) && WeekDay.valueOf(userInfo.firstDayOfWeek) ? userInfo.firstDayOfWeek : settingsStore.localeDefaultSettings.firstDayOfWeek;
    });

    const currentUserLongDateFormat = computed<number>(() => {
        const userInfo = currentUserBasicInfo.value || EMPTY_USER_BASIC_INFO;
        return userInfo.longDateFormat;
    });

    const currentUserShortDateFormat = computed<number>(() => {
        const userInfo = currentUserBasicInfo.value || EMPTY_USER_BASIC_INFO;
        return userInfo.shortDateFormat;
    });

    const currentUserLongTimeFormat = computed<number>(() => {
        const userInfo = currentUserBasicInfo.value || EMPTY_USER_BASIC_INFO;
        return userInfo.longTimeFormat;
    });

    const currentUserShortTimeFormat = computed<number>(() => {
        const userInfo = currentUserBasicInfo.value || EMPTY_USER_BASIC_INFO;
        return userInfo.shortTimeFormat;
    });

    const currentUserDecimalSeparator = computed<number>(() => {
        const userInfo = currentUserBasicInfo.value || EMPTY_USER_BASIC_INFO;
        return userInfo.decimalSeparator;
    });

    const currentUserDigitGroupingSymbol = computed<number>(() => {
        const userInfo = currentUserBasicInfo.value || EMPTY_USER_BASIC_INFO;
        return userInfo.digitGroupingSymbol;
    });

    const currentUserDigitGrouping = computed<number>(() => {
        const userInfo = currentUserBasicInfo.value || EMPTY_USER_BASIC_INFO;
        return userInfo.digitGrouping;
    });

    const currentUserCurrencyDisplayType = computed<number>(() => {
        const userInfo = currentUserBasicInfo.value || EMPTY_USER_BASIC_INFO;
        return userInfo.currencyDisplayType;
    });

    const currentUserExpenseAmountColor = computed<number>(() => {
        const userInfo = currentUserBasicInfo.value || EMPTY_USER_BASIC_INFO;
        return userInfo.expenseAmountColor;
    });

    const currentUserIncomeAmountColor = computed<number>(() => {
        const userInfo = currentUserBasicInfo.value || EMPTY_USER_BASIC_INFO;
        return userInfo.incomeAmountColor;
    });

    function generateNewUserModel(language: string): User {
        return User.createNewUser(language, settingsStore.localeDefaultSettings.currency, settingsStore.localeDefaultSettings.firstDayOfWeek);
    }

    function storeUserBasicInfo(userInfo: UserBasicInfo): void {
        currentUserBasicInfo.value = userInfo;
        updateCurrentUserInfo(userInfo);
    }

    function resetUserBasicInfo(): void {
        currentUserBasicInfo.value = null;
        clearCurrentUserInfo();
    }

    function getCurrentUserProfile(): Promise<UserProfileResponse> {
        return new Promise((resolve, reject) => {
            services.getProfile().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve user profile' });
                    return;
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to retrieve user profile', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve user profile' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function updateUserTransactionEditScope(req: { transactionEditScope: number }): Promise<UserProfileUpdateResponse> {
        return new Promise((resolve, reject) => {
            services.updateProfile({
                transactionEditScope: req.transactionEditScope,
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result || !data.result.user || !isObject(data.result.user)) {
                    reject({ message: 'Unable to update editable transaction range' });
                    return;
                }

                storeUserBasicInfo(data.result.user);

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to save editable transaction range', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to update editable transaction range' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function updateUserAvatar(req: { avatarFile: unknown }): Promise<UserProfileResponse> {
        return new Promise((resolve, reject) => {
            services.updateAvatar({ avatarFile: req.avatarFile }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to update user avatar' });
                    return;
                }

                storeUserBasicInfo(data.result);

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to update user avatar', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to update user avatar' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function removeUserAvatar(): Promise<UserProfileResponse> {
        return new Promise((resolve, reject) => {
            services.removeAvatar().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to remove user avatar' });
                    return;
                }

                storeUserBasicInfo(data.result);

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to remove user avatar', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to remove user avatar' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function getUserDataStatistics(): Promise<DataStatisticsResponse> {
        return new Promise((resolve, reject) => {
            services.getUserDataStatistics().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve user statistics data' });
                    return;
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to retrieve user statistics data', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve user statistics data' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function getExportedUserData(fileType: string): Promise<Blob> {
        return new Promise((resolve, reject) => {
            services.getExportedUserData(fileType).then(response => {
                if (response && response.headers) {
                    if (fileType === 'csv' && response.headers['content-type'] !== 'text/csv') {
                        reject({ message: 'Unable to retrieve exported user data' });
                        return;
                    } else if (fileType === 'tsv' && response.headers['content-type'] !== 'text/tab-separated-values') {
                        reject({ message: 'Unable to retrieve exported user data' });
                        return;
                    }
                }

                const blob = new Blob([response.data], { type: response.headers['content-type'] });
                resolve(blob);
            }).catch(error => {
                logger.error('failed to retrieve user statistics data', error);

                if (error.response && error.response.headers['content-type'] === 'text/text' && error.response && error.response.data) {
                    reject({ message: 'error.' + error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve exported user data' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function getUserAvatarUrl(userInfo: UserBasicInfo, disableBrowserCache: boolean | string): string | null {
        if (!userInfo || !userInfo.avatar) {
            return null;
        }

        return services.getInternalAvatarUrlWithToken(userInfo.avatar, disableBrowserCache);
    }

    return {
        // states
        currentUserBasicInfo,
        // computed states
        currentUserNickname,
        currentUserAvatar,
        currentUserDefaultAccountId,
        currentUserLanguage,
        currentUserDefaultCurrency,
        currentUserFirstDayOfWeek,
        currentUserLongDateFormat,
        currentUserShortDateFormat,
        currentUserLongTimeFormat,
        currentUserShortTimeFormat,
        currentUserDecimalSeparator,
        currentUserDigitGroupingSymbol,
        currentUserDigitGrouping,
        currentUserCurrencyDisplayType,
        currentUserExpenseAmountColor,
        currentUserIncomeAmountColor,
        // functions
        generateNewUserModel,
        storeUserBasicInfo,
        resetUserBasicInfo,
        getCurrentUserProfile,
        updateUserTransactionEditScope,
        updateUserAvatar,
        removeUserAvatar,
        getUserDataStatistics,
        getExportedUserData,
        getUserAvatarUrl
    };
});