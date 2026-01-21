import { defineStore } from 'pinia';

import type {UserExternalAuthInfoResponse, UserExternalAuthUnlinkRequest} from '@/models/user_external_auth.ts';

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

    return {
        // functions
        getExternalAuths,
        unlinkExternalAuth
    };
});
