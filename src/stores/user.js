import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.js';

import userState from '@/lib/userstate.js';
import services from '@/lib/services.js';
import logger from '@/lib/logger.js';
import { isNumber } from '@/lib/common.js';

export const useUserStore = defineStore('user', {
    state: () => ({
        currentUserBasicInfo: userState.getUserInfo()
    }),
    getters: {
        currentUserNickname(state) {
            const userInfo = state.currentUserBasicInfo || {};
            return userInfo.nickname || userInfo.username || null;
        },
        currentUserAvatar(state) {
            const userInfo = state.currentUserBasicInfo || {};
            return state.getUserAvatarUrl(userInfo, false);
        },
        currentUserDefaultAccountId(state) {
            const userInfo = state.currentUserBasicInfo || {};
            return userInfo.defaultAccountId || '';
        },
        currentUserLanguage(state) {
            const userInfo = state.currentUserBasicInfo || {};
            return userInfo.language;
        },
        currentUserDefaultCurrency(state) {
            const settingsStore = useSettingsStore();
            const userInfo = state.currentUserBasicInfo || {};
            return userInfo.defaultCurrency || settingsStore.localeDefaultSettings.currency;
        },
        currentUserFirstDayOfWeek(state) {
            const settingsStore = useSettingsStore();
            const userInfo = state.currentUserBasicInfo || {};
            return isNumber(userInfo.firstDayOfWeek) ? userInfo.firstDayOfWeek : settingsStore.localeDefaultSettings.firstDayOfWeek;
        },
        currentUserLongDateFormat(state) {
            const userInfo = state.currentUserBasicInfo || {};
            return userInfo.longDateFormat;
        },
        currentUserShortDateFormat(state) {
            const userInfo = state.currentUserBasicInfo || {};
            return userInfo.shortDateFormat;
        },
        currentUserLongTimeFormat(state) {
            const userInfo = state.currentUserBasicInfo || {};
            return userInfo.longTimeFormat;
        },
        currentUserShortTimeFormat(state) {
            const userInfo = state.currentUserBasicInfo || {};
            return userInfo.shortTimeFormat;
        },
        currentUserDecimalSeparator(state) {
            const userInfo = state.currentUserBasicInfo || {};
            return userInfo.decimalSeparator;
        },
        currentUserDigitGroupingSymbol(state) {
            const userInfo = state.currentUserBasicInfo || {};
            return userInfo.digitGroupingSymbol;
        },
        currentUserDigitGrouping(state) {
            const userInfo = state.currentUserBasicInfo || {};
            return userInfo.digitGrouping;
        },
        currentUserCurrencyDisplayType(state) {
            const userInfo = state.currentUserBasicInfo || {};
            return userInfo.currencyDisplayType;
        },
        currentUserExpenseAmountColor(state) {
            const userInfo = state.currentUserBasicInfo || {};
            return userInfo.expenseAmountColor;
        },
        currentUserIncomeAmountColor(state) {
            const userInfo = state.currentUserBasicInfo || {};
            return userInfo.incomeAmountColor;
        }
    },
    actions: {
        generateNewUserModel(language) {
            const settingsStore = useSettingsStore();

            return {
                username: '',
                password: '',
                confirmPassword: '',
                email: '',
                nickname: '',
                language: language,
                defaultCurrency: settingsStore.localeDefaultSettings.currency,
                firstDayOfWeek: settingsStore.localeDefaultSettings.firstDayOfWeek,
            };
        },
        storeUserBasicInfo(userInfo) {
            this.currentUserBasicInfo = userInfo;
            userState.updateUserInfo(userInfo);
        },
        resetUserBasicInfo() {
            this.currentUserBasicInfo = null;
            userState.clearUserInfo();
        },
        getCurrentUserProfile() {
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
        },
        updateUserAvatar({ avatarFile }) {
            return new Promise((resolve, reject) => {
                services.updateAvatar({ avatarFile }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to update user avatar' });
                        return;
                    }

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
        },
        removeUserAvatar() {
            return new Promise((resolve, reject) => {
                services.removeAvatar().then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to remove user avatar' });
                        return;
                    }

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
        },
        getUserDataStatistics() {
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
        },
        getExportedUserData(fileType) {
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
        },
        getUserAvatarUrl(userInfo, disableBrowserCache) {
            if (!userInfo || !userInfo.avatar) {
                return null;
            }

            return services.getInternalAvatarUrlWithToken(userInfo.avatar, disableBrowserCache);
        }
    }
});
