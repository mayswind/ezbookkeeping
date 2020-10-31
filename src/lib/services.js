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
    if (error.response && error.response.data && error.response.data.errorCode) {
        const errorCode = error.response.data.errorCode;

        if (202001 <= errorCode && errorCode <= 202007) { // unauthorized access or token is invalid
            userState.clearToken();
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
    register: ({ username, email, nickname, password }) => {
        return axios.post('register.json', {
            username,
            email,
            nickname,
            password
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

                if (data.result.oldTokenId) {
                    axios.post('v1/tokens/revoke.json', {
                        tokenId: data.result.oldTokenId
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
    getProfile: () => {
        return axios.get('v1/users/profile/get.json');
    },
    updateProfile: ({ email, nickname, password }) => {
        return axios.post('v1/users/profile/update.json', {
            email,
            nickname,
            password
        });
    },
};
