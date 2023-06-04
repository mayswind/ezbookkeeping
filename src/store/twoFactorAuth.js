import userState from '../lib/userstate.js';
import services from '../lib/services.js';
import logger from '../lib/logger.js';
import utilities from '../lib/utilities/index.js';

export function get2FAStatus() {
    return new Promise((resolve, reject) => {
        services.get2FAStatus().then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result || !utilities.isBoolean(data.result.enable)) {
                reject({ message: 'Unable to get current two factor authentication status' });
                return;
            }

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to get 2fa status', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                reject({ message: 'Unable to get current two factor authentication status' });
            } else {
                reject(error);
            }
        });
    });
}

export function enable2FA() {
    return new Promise((resolve, reject) => {
        services.enable2FA().then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result || !data.result.qrcode || !data.result.secret) {
                reject({ message: 'Unable to enable two factor authentication' });
                return;
            }

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to request to enable 2fa', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                reject({ message: 'Unable to enable two factor authentication' });
            } else {
                reject(error);
            }
        });
    });
}

export function confirmEnable2FA(context, { secret, passcode }) {
    return new Promise((resolve, reject) => {
        services.confirmEnable2FA({
            secret: secret,
            passcode: passcode
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result || !data.result.token) {
                reject({ message: 'Unable to enable two factor authentication' });
                return;
            }

            if (data.result.token) {
                userState.updateToken(data.result.token);
            }

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to confirm to enable 2fa', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                reject({ message: 'Unable to enable two factor authentication' });
            } else {
                reject(error);
            }
        });
    });
}

export function disable2FA(context, { password }) {
    return new Promise((resolve, reject) => {
        services.disable2FA({
            password: password
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                reject({ message: 'Unable to disable two factor authentication' });
                return;
            }

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to disable 2fa', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                reject({ message: 'Unable to disable two factor authentication' });
            } else {
                reject(error);
            }
        });
    });
}

export function regenerate2FARecoveryCode(context, { password }) {
    return new Promise((resolve, reject) => {
        services.regenerate2FARecoveryCode({
            password: password
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result || !data.result.recoveryCodes || !data.result.recoveryCodes.length) {
                reject({ message: 'Unable to regenerate two factor authentication backup codes' });
                return;
            }

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to regenerate 2fa recovery code', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                reject({ message: 'Unable to regenerate two factor authentication backup codes' });
            } else {
                reject(error);
            }
        });
    });
}
