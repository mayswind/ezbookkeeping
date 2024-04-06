import { defineStore } from 'pinia';

import { useUserStore } from './user.js';

import userState from '@/lib/userstate.js';
import services from '@/lib/services.js';
import logger from '@/lib/logger.js';
import { isObject } from '@/lib/common.js';

export const useTokensStore = defineStore('tokens', {
    actions: {
        getAllTokens() {
            return new Promise((resolve, reject) => {
                services.getTokens().then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to retrieve session list' });
                        return;
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to load token list', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to retrieve session list' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
        refreshTokenAndRevokeOldToken() {
            const self = this;

            return new Promise((resolve) => {
                services.refreshToken().then(response => {
                    const data = response.data;

                    if (data && data.success && data.result && data.result.newToken) {
                        userState.updateToken(data.result.newToken);

                        if (data.result.user && isObject(data.result.user)) {
                            const userStore = useUserStore();
                            userStore.storeUserInfo(data.result.user);
                        }

                        if (data.result.oldTokenId) {
                            self.revokeToken({
                                tokenId: data.result.oldTokenId,
                                ignoreError: true
                            });
                        }
                    }

                    resolve(data.result);
                });
            });
        },
        revokeToken({ tokenId, ignoreError }) {
            return new Promise((resolve, reject) => {
                services.revokeToken({
                    tokenId: tokenId,
                    ignoreError: !!ignoreError
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to logout from this session' });
                        return;
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to revoke token', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to logout from this session' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
        revokeAllTokens() {
            return new Promise((resolve, reject) => {
                services.revokeAllTokens().then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to logout all other sessions' });
                        return;
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to revoke all tokens', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to logout all other sessions' });
                    } else {
                        reject(error);
                    }
                });
            });
        }
    }
});
