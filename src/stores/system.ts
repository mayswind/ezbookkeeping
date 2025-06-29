import { defineStore } from 'pinia';

import type { VersionInfo } from '@/core/version.ts';

import logger from '@/lib/logger.ts';
import services from '@/lib/services.ts';
import { getClientVersionInfo } from '@/lib/version.ts';

export const useSystemsStore = defineStore('systems', () => {
    function checkIfClientVersionMatchServerVersion(): Promise<{ match: boolean, version: VersionInfo }> {
        return new Promise((resolve, reject) => {
            services.getServerVersion().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve server version' });
                    return;
                }

                const clientVersionInfo = getClientVersionInfo();

                if (data.result.version && clientVersionInfo.version !== data.result.version) {
                    logger.warn(`client version \"${clientVersionInfo.version}\" does not match server version \"${data.result.version}\"`);
                    resolve({
                        match: false,
                        version: data.result
                    });
                }

                if (data.result.commitHash && clientVersionInfo.commitHash !== data.result.commitHash) {
                    logger.warn(`client commit hash \"${clientVersionInfo.commitHash}\" does not match server commit hash \"${data.result.commitHash}\"`);
                    resolve({
                        match: false,
                        version: data.result
                    });
                }

                resolve({
                    match: true,
                    version: data.result
                });
            }).catch(error => {
                logger.error('failed to retrieve server version', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve server version' });
                } else {
                    reject(error);
                }
            });
        });
    }

    return {
        // functions
        checkIfClientVersionMatchServerVersion,
    };
});
