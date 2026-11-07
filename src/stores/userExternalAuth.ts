import { defineStore } from 'pinia';

import type {
    UserExternalAuthUnlinkRequest,
    UserExternalAuthInfoResponse,
    OAuth2LinkRequest
} from '@/models/user_external_auth.ts';

import logger from '@/lib/logger.ts';
import services from '@/lib/services.ts';

export const useUserExternalAuthStore = defineStore('userExternalAUth', () => {
    function getExternalAuths(): Promise<UserExternalAuthInfoResponse[]> {
        return new Promise((resolve, reject) => {
            services.getExternalAuths().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve third-party login list' });
                    return;
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to load third-party login list', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve third-party login list' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function unlinkExternalAuth(req: UserExternalAuthUnlinkRequest): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.unlinkExternalAuth(req).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to unlink third-party login' });
                    return;
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to revoke token', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to unlink third-party login' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function linkExternalOAuth2(req: OAuth2LinkRequest): Promise<string> {
        return new Promise((resolve, reject) => {
            services.linkExternalOAuth2(req).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'An error occurred' });
                    return;
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to get external auth link url', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'An error occurred' });
                } else {
                    reject(error);
                }
            });
        });
    }

    return {
        // functions
        getExternalAuths,
        unlinkExternalAuth,
        linkExternalOAuth2
    };
});
