import { defineStore } from 'pinia';

import type {
    TwoFactorEnableConfirmRequest,
    TwoFactorEnableResponse,
    TwoFactorEnableConfirmResponse,
    TwoFactorDisableRequest,
    TwoFactorRegenerateRecoveryCodeRequest,
    TwoFactorStatusResponse
} from '@/models/two_factor.ts';

import { isBoolean } from '@/lib/common.ts';
import { updateCurrentToken } from '@/lib/userstate.ts';

import logger from '@/lib/logger.ts';
import services from '@/lib/services.ts';

export const useTwoFactorAuthStore = defineStore('twoFactorAuth', () => {
    function get2FAStatus(): Promise<TwoFactorStatusResponse> {
        return new Promise((resolve, reject) => {
            services.get2FAStatus().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result || !isBoolean(data.result.enable)) {
                    reject({ message: 'Unable to retrieve current two-factor authentication status' });
                    return;
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to retrieve 2fa status', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve current two-factor authentication status' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function enable2FA(): Promise<TwoFactorEnableResponse> {
        return new Promise((resolve, reject) => {
            services.enable2FA().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result || !data.result.qrcode || !data.result.secret) {
                    reject({ message: 'Unable to enable two-factor authentication' });
                    return;
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to request to enable 2fa', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to enable two-factor authentication' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function confirmEnable2FA(req: TwoFactorEnableConfirmRequest): Promise<TwoFactorEnableConfirmResponse> {
        return new Promise((resolve, reject) => {
            services.confirmEnable2FA(req).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result || !data.result.token) {
                    reject({ message: 'Unable to enable two-factor authentication' });
                    return;
                }

                if (data.result.token) {
                    updateCurrentToken(data.result.token);
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to confirm to enable 2fa', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to enable two-factor authentication' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function disable2FA(req: TwoFactorDisableRequest): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.disable2FA(req).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to disable two-factor authentication' });
                    return;
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to disable 2fa', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to disable two-factor authentication' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function regenerate2FARecoveryCode(req: TwoFactorRegenerateRecoveryCodeRequest): Promise<TwoFactorEnableConfirmResponse> {
        return new Promise((resolve, reject) => {
            services.regenerate2FARecoveryCode(req).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result || !data.result.recoveryCodes || !data.result.recoveryCodes.length) {
                    reject({ message: 'Unable to regenerate two-factor authentication backup codes' });
                    return;
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to regenerate 2fa recovery code', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to regenerate two-factor authentication backup codes' });
                } else {
                    reject(error);
                }
            });
        });
    }

    return {
        // functions
        get2FAStatus,
        enable2FA,
        confirmEnable2FA,
        disable2FA,
        regenerate2FARecoveryCode
    };
});
