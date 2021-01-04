import userState from '../lib/userstate.js';
import services from '../lib/services.js';
import logger from '../lib/logger.js';
import utils from '../lib/utils.js';

import {
    STORE_USER_INFO
} from './mutations.js';

function getAllTokens() {
    return new Promise((resolve, reject) => {
        services.getTokens().then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                reject({ message: 'Unable to get session list' });
                return;
            }

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to load token list', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                reject({ message: 'Unable to get session list' });
            } else {
                reject(error);
            }
        });
    });
}

function refreshTokenAndRevokeOldToken(context) {
    return new Promise((resolve) => {
        services.refreshToken().then(response => {
            const data = response.data;

            if (data && data.success && data.result && data.result.newToken) {
                userState.updateToken(data.result.newToken);

                if (data.result.user && utils.isObject(data.result.user)) {
                    context.commit(STORE_USER_INFO, data.result.user);
                }

                if (data.result.oldTokenId) {
                    revokeToken(context, {
                        tokenId: data.result.oldTokenId,
                        ignoreError: true
                    });
                }
            }

            resolve(data.result);
        });
    });
}

function revokeToken(context, { tokenId, ignoreError }) {
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
}

function revokeAllTokens() {
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

export default {
    getAllTokens,
    refreshTokenAndRevokeOldToken,
    revokeToken,
    revokeAllTokens
}
