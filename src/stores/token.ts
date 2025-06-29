import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.ts';
import { useUserStore } from './user.ts';

import type { TokenRefreshResponse, TokenInfoResponse } from '@/models/token.ts';

import { isObject } from '@/lib/common.ts';
import { updateCurrentToken } from '@/lib/userstate.ts';

import logger from '@/lib/logger.ts';
import services from '@/lib/services.ts';

export const useTokensStore = defineStore('tokens', () => {
    const settingsStore = useSettingsStore();
    const userStore = useUserStore();

    function getAllTokens(): Promise<TokenInfoResponse[]> {
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
    }

    function refreshTokenAndRevokeOldToken(): Promise<TokenRefreshResponse> {
        return new Promise((resolve) => {
            services.refreshToken().then(response => {
                const data = response.data;

                if (data && data.success && data.result) {
                    settingsStore.setApplicationSettingsFromCloudSettings(data.result.applicationCloudSettings);
                }

                if (data && data.success && data.result && data.result.user && isObject(data.result.user)) {
                    userStore.storeUserBasicInfo(data.result.user);
                }

                if (data && data.success && data.result && data.result.newToken) {
                    updateCurrentToken(data.result.newToken);

                    if (data.result.oldTokenId) {
                        revokeToken({
                            tokenId: data.result.oldTokenId,
                            ignoreError: true
                        });
                    }
                }

                resolve(data.result);
            });
        });
    }

    function revokeToken({ tokenId, ignoreError }: { tokenId: string, ignoreError?: boolean }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.revokeToken({ tokenId, ignoreError }).then(response => {
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
    }

    function revokeAllTokens(): Promise<boolean> {
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

    return {
        // functions
        getAllTokens,
        refreshTokenAndRevokeOldToken,
        revokeToken,
        revokeAllTokens
    };
});
