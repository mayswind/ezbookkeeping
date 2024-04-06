import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.js';

import userState from '@/lib/userstate.js';
import services from '@/lib/services.js';
import logger from '@/lib/logger.js';
import { isNumber } from '@/lib/common.js';

export const useUserStore = defineStore('user', {
    state: () => ({
        currentUserInfo: userState.getUserInfo()
    }),
    getters: {
        currentUserNickname(state) {
            const userInfo = state.currentUserInfo || {};
            return userInfo.nickname || userInfo.username || null;
        },
        currentUserAvatar(state) {
            const userInfo = state.currentUserInfo || {};
            return userInfo.avatar || null;
        },
        currentUserDefaultAccountId(state) {
            const userInfo = state.currentUserInfo || {};
            return userInfo.defaultAccountId || '';
        },
        currentUserLanguage(state) {
            const settingsStore = useSettingsStore();
            const userInfo = state.currentUserInfo || {};
            return userInfo.language || settingsStore.language;
        },
        currentUserDefaultCurrency(state) {
            const settingsStore = useSettingsStore();
            const userInfo = state.currentUserInfo || {};
            return userInfo.defaultCurrency || settingsStore.currency;
        },
        currentUserFirstDayOfWeek(state) {
            const settingsStore = useSettingsStore();
            const userInfo = state.currentUserInfo || {};
            return isNumber(userInfo.firstDayOfWeek) ? userInfo.firstDayOfWeek : settingsStore.firstDayOfWeek;
        },
        currentUserLongDateFormat(state) {
            const settingsStore = useSettingsStore();
            const userInfo = state.currentUserInfo || {};
            return isNumber(userInfo.longDateFormat) ? userInfo.longDateFormat : settingsStore.longDateFormat;
        },
        currentUserShortDateFormat(state) {
            const settingsStore = useSettingsStore();
            const userInfo = state.currentUserInfo || {};
            return isNumber(userInfo.shortDateFormat) ? userInfo.shortDateFormat : settingsStore.shortDateFormat;
        },
        currentUserLongTimeFormat(state) {
            const settingsStore = useSettingsStore();
            const userInfo = state.currentUserInfo || {};
            return isNumber(userInfo.longTimeFormat) ? userInfo.longTimeFormat : settingsStore.longTimeFormat;
        },
        currentUserShortTimeFormat(state) {
            const settingsStore = useSettingsStore();
            const userInfo = state.currentUserInfo || {};
            return isNumber(userInfo.shortTimeFormat) ? userInfo.shortTimeFormat : settingsStore.shortTimeFormat;
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
        storeUserInfo(userInfo) {
            this.currentUserInfo = userInfo;
            userState.updateUserInfo(userInfo);
        },
        resetUserInfo() {
            this.currentUserInfo = null;
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
    }
});
