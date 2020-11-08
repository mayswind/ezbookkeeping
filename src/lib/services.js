import axios from 'axios';
import userState from "./userstate.js";

let needBlockRequest = false;
let blockedRequests = [];

axios.defaults.baseURL = '/api';
axios.interceptors.request.use(config => {
    const token = userState.getToken();

    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }

    if (needBlockRequest && !config.ignoreBlocked) {
        return new Promise(resolve => {
            blockedRequests.push(newToken => {
                if (newToken) {
                    config.headers.Authorization = `Bearer ${newToken}`;
                }

                resolve(config);
            });
        });
    }

    return config;
}, error => {
    return Promise.reject(error);
});

axios.interceptors.response.use(response => {
    return response;
}, error => {
    if (error.response && !error.response.config.ignoreError && error.response.data && error.response.data.errorCode) {
        const errorCode = error.response.data.errorCode;

        if (errorCode === 202001 // unauthorized access
            && errorCode <= 202002 // current token is invalid
            && errorCode <= 202003 // current token is expired
            && errorCode <= 202004 // current token type is invalid
            && errorCode <= 202005 // current token requires two factor authorization
            && errorCode <= 202006) { // current token does not require two factor authorization
            userState.clearTokenAndUserInfo();
            location.reload();
            return Promise.reject({ processed: true });
        }
    }

    return Promise.reject(error);
});

export default {
    setLocale: locale => {
        axios.defaults.headers.common['Accept-Language'] = locale;
    },
    authorize: ({ loginName, password }) => {
        return axios.post('authorize.json', {
            loginName,
            password
        });
    },
    authorize2FA: ({ passcode, token }) => {
        return axios.post('2fa/authorize.json', {
            passcode
        }, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
    },
    authorize2FAByBackupCode: ({ recoveryCode, token }) => {
        return axios.post('2fa/recovery.json', {
            recoveryCode
        }, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
    },
    register: ({ username, email, nickname, password, defaultCurrency }) => {
        return axios.post('register.json', {
            username,
            email,
            nickname,
            password,
            defaultCurrency
        });
    },
    logout: () => {
        return axios.get('logout.json');
    },
    refreshToken: () => {
        needBlockRequest = true;

        return axios.post('v1/tokens/refresh.json', {} , {
            ignoreBlocked: true
        }).then(response => {
            const data = response.data;

            if (data && data.success && data.result && data.result.newToken) {
                userState.updateToken(data.result.newToken);
                userState.updateUserInfo(data.result.user);

                if (data.result.oldTokenId) {
                    axios.post('v1/tokens/revoke.json', {
                        tokenId: data.result.oldTokenId
                    }, {
                        ignoreError: true
                    });
                }
            }

            needBlockRequest = false;
            return data.result.newToken;
        }).then(newToken => {
            blockedRequests.forEach(func => func(newToken));
            blockedRequests.length = 0;
        });
    },
    getTokens: () => {
        return axios.get('v1/tokens/list.json');
    },
    revokeToken: ({ tokenId }) => {
        return axios.post('v1/tokens/revoke.json', {
            tokenId
        });
    },
    getProfile: () => {
        return axios.get('v1/users/profile/get.json');
    },
    updateProfile: ({ email, nickname, password, oldPassword, defaultCurrency }) => {
        return axios.post('v1/users/profile/update.json', {
            email,
            nickname,
            password,
            oldPassword,
            defaultCurrency
        });
    },
    get2FAStatus: () => {
        return axios.get('v1/users/2fa/status.json');
    },
    enable2FA: () => {
        return axios.post('v1/users/2fa/enable/request.json');
    },
    confirmEnable2FA: ({ secret, passcode }) => {
        return axios.post('v1/users/2fa/enable/confirm.json', {
            secret,
            passcode
        });
    },
    disable2FA: ({ password }) => {
        return axios.post('v1/users/2fa/disable.json', {
            password
        });
    },
    regenerate2FARecoveryCode: ({ password }) => {
        return axios.post('v1/users/2fa/recovery/regenerate.json', {
            password
        });
    },
    getAllAccounts: () => {
        return axios.get('v1/accounts/list.json');
    },
    addAccount: ({ category, type, name, icon, currency, comment, subAccounts }) => {
        return axios.post('v1/accounts/add.json', {
            category,
            type,
            name,
            icon,
            currency,
            comment,
            subAccounts
        });
    },
};
