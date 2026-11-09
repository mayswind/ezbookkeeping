import { ref } from 'vue';
import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.ts';
import { useUserStore } from './user.ts';
import { useUserCustomIconsStore } from './userCustomIcon.ts';

import type {
    TokenGenerateAPIResponse,
    TokenGenerateMCPResponse,
    TokenRefreshResponse,
    TokenInfoResponse
} from '@/models/token.ts';

import { isObject } from '@/lib/common.ts';
import { updateCurrentToken } from '@/lib/userstate.ts';

import logger from '@/lib/logger.ts';
import services from '@/lib/services.ts';

export const useTokensStore = defineStore('tokens', () => {
    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const userCustomIconsStore = useUserCustomIconsStore();

    const mapProxyToken = ref<string | null>(null);

    function updateMapProxyToken(token?: string): void {
        if (token) {
            mapProxyToken.value = token;
        }
    }

    function resetMapProxyToken(): void {
        mapProxyToken.value = null;
    }

    function getMapProxyToken(): string | null {
        return mapProxyToken.value;
    }

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
        return new Promise((resolve, reject) => {
            services.refreshToken().then(response => {
                const data = response.data;

                if (data && data.success && data.result) {
                    settingsStore.setApplicationSettingsFromCloudSettings(data.result.applicationCloudSettings);
                    updateMapProxyToken(data.result.newMapProxyToken);
                }

                if (data && data.success && data.result && data.result.user && isObject(data.result.user)) {
                    userStore.storeUserBasicInfo(data.result.user);
                }

                if (data && data.success && data.result && data.result.newToken) {
                    updateCurrentToken(data.result.newToken);

                    if (!userCustomIconsStore.customIconListStateInvalid) {
                        userCustomIconsStore.updateCustomIconListInvalidState(true);
                    }

                    if (data.result.oldTokenId) {
                        revokeToken({
                            tokenId: data.result.oldTokenId,
                            ignoreError: true
                        }).catch(error => {
                            logger.warn('failed to revoke old token after refreshing', error);
                        });
                    }
                }

                resolve(data.result);
            }).catch(reject);
        });
    }

    function generateToken<T extends 'api' | 'mcp'>({ type, expiresInSeconds, password, enabledMCPTools }: { type: T, expiresInSeconds: number, password: string, enabledMCPTools?: string[] }): Promise<{ 'api': TokenGenerateAPIResponse, 'mcp': TokenGenerateMCPResponse }[T]> {
        return new Promise((resolve, reject) => {
            let promise = null;

            if (type === 'api') {
                promise = services.generateAPIToken({
                    expiresInSeconds: expiresInSeconds,
                    password: password
                });
            } else if (type === 'mcp') {
                promise = services.generateMCPToken({
                    expiresInSeconds: expiresInSeconds,
                    password: password,
                    enabledTools: enabledMCPTools
                });
            } else {
                reject({ message: 'An error occurred' });
                return;
            }

            promise.then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to generate token' });
                    return;
                }

                resolve(data.result as { 'api': TokenGenerateAPIResponse, 'mcp': TokenGenerateMCPResponse }[T]);
            }).catch(error => {
                logger.error('failed to generate token', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to generate token' });
                } else {
                    reject(error);
                }
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
        updateMapProxyToken,
        resetMapProxyToken,
        getMapProxyToken,
        getAllTokens,
        refreshTokenAndRevokeOldToken,
        generateToken,
        revokeToken,
        revokeAllTokens
    };
});
