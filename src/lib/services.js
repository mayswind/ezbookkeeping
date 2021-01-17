import axios from 'axios';

import api from "../consts/api.js";
import userState from "./userstate.js";

let needBlockRequest = false;
let blockedRequests = [];

axios.defaults.baseURL = api.baseUrlPath;
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
            || errorCode === 202002 // current token is invalid
            || errorCode === 202003 // current token is expired
            || errorCode === 202004 // current token type is invalid
            || errorCode === 202005 // current token requires two factor authorization
            || errorCode === 202006 // current token does not require two factor authorization
        ) {
            userState.clearTokenAndUserInfo(false);
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
        return new Promise((resolve) => {
            needBlockRequest = true;

            axios.post('v1/tokens/refresh.json', {}, {
                ignoreBlocked: true
            }).then(response => {
                const data = response.data;

                resolve(response);
                needBlockRequest = false;

                return data.result.newToken;
            }).then(newToken => {
                blockedRequests.forEach(func => func(newToken));
                blockedRequests.length = 0;
            });
        });
    },
    getTokens: () => {
        return axios.get('v1/tokens/list.json');
    },
    revokeToken: ({ tokenId, ignoreError }) => {
        return axios.post('v1/tokens/revoke.json', {
            tokenId
        }, {
            ignoreError: !!ignoreError
        });
    },
    revokeAllTokens: () => {
        return axios.post('v1/tokens/revoke_all.json');
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
    clearData: ({ password }) => {
        return axios.post('v1/data/clear.json', {
            password
        });
    },
    getTransactionOverview: ( { today, thisWeek, thisMonth, thisYear } ) => {
        const queryParams = [];

        if (today) {
            queryParams.push(`today_${today.startTime}_${today.endTime}`);
        }

        if (thisWeek) {
            queryParams.push(`thisWeek_${thisWeek.startTime}_${thisWeek.endTime}`);
        }

        if (thisMonth) {
            queryParams.push(`thisMonth_${thisMonth.startTime}_${thisMonth.endTime}`);
        }

        if (thisYear) {
            queryParams.push(`thisYear_${thisYear.startTime}_${thisYear.endTime}`);
        }

        return axios.get('v1/overviews/transaction.json' + (queryParams.length ? '?query=' + queryParams.join('|') : ''));
    },
    getAllAccounts: ({ visibleOnly }) => {
        return axios.get('v1/accounts/list.json?visible_only=' + !!visibleOnly);
    },
    getAccount: ({ id }) => {
        return axios.get('v1/accounts/get.json?id=' + id);
    },
    addAccount: ({ category, type, name, icon, color, currency, balance, comment, subAccounts }) => {
        return axios.post('v1/accounts/add.json', {
            category,
            type,
            name,
            icon,
            color,
            currency,
            balance,
            comment,
            subAccounts
        });
    },
    modifyAccount: ({ id, category, name, icon, color, comment, hidden, subAccounts }) => {
        return axios.post('v1/accounts/modify.json', {
            id,
            category,
            name,
            icon,
            color,
            comment,
            hidden,
            subAccounts
        });
    },
    hideAccount: ({ id, hidden }) => {
        return axios.post('v1/accounts/hide.json', {
            id,
            hidden
        });
    },
    moveAccount: ({ newDisplayOrders }) => {
        return axios.post('v1/accounts/move.json', {
            newDisplayOrders,
        });
    },
    deleteAccount: ({ id }) => {
        return axios.post('v1/accounts/delete.json', {
            id
        });
    },
    getTransactions: ({ maxTime, minTime, type, categoryId, accountId, keyword }) => {
        return axios.get(`v1/transactions/list.json?max_time=${maxTime}&min_time=${minTime}&type=${type}&category_id=${categoryId}&account_id=${accountId}&keyword=${keyword}&count=50`);
    },
    getTransaction: ({ id }) => {
        return axios.get('v1/transactions/get.json?id=' + id);
    },
    addTransaction: ({ type, categoryId, time, sourceAccountId, destinationAccountId, sourceAmount, destinationAmount, tagIds, comment }) => {
        return axios.post('v1/transactions/add.json', {
            type,
            categoryId,
            time,
            sourceAccountId,
            destinationAccountId,
            sourceAmount,
            destinationAmount,
            tagIds,
            comment
        });
    },
    modifyTransaction: ({ id, type, categoryId, time, sourceAccountId, destinationAccountId, sourceAmount, destinationAmount, tagIds, comment }) => {
        return axios.post('v1/transactions/modify.json', {
            id,
            type,
            categoryId,
            time,
            sourceAccountId,
            destinationAccountId,
            sourceAmount,
            destinationAmount,
            tagIds,
            comment
        });
    },
    deleteTransaction: ({ id }) => {
        return axios.post('v1/transactions/delete.json', {
            id
        });
    },
    getAllTransactionCategories: () => {
        return axios.get('v1/transaction/categories/list.json');
    },
    getTransactionCategory: ({ id }) => {
        return axios.get('v1/transaction/categories/get.json?id=' + id);
    },
    addTransactionCategory: ({ name, type, parentId, icon, color, comment }) => {
        return axios.post('v1/transaction/categories/add.json', {
            name,
            type,
            parentId,
            icon,
            color,
            comment
        });
    },
    addTransactionCategoryBatch: ({ categories }) => {
        return axios.post('v1/transaction/categories/add_batch.json', {
            categories
        });
    },
    modifyTransactionCategory: ({ id, name, icon, color, comment, hidden }) => {
        return axios.post('v1/transaction/categories/modify.json', {
            id,
            name,
            icon,
            color,
            comment,
            hidden
        });
    },
    hideTransactionCategory: ({ id, hidden }) => {
        return axios.post('v1/transaction/categories/hide.json', {
            id,
            hidden
        });
    },
    moveTransactionCategory: ({ newDisplayOrders }) => {
        return axios.post('v1/transaction/categories/move.json', {
            newDisplayOrders,
        });
    },
    deleteTransactionCategory: ({ id }) => {
        return axios.post('v1/transaction/categories/delete.json', {
            id
        });
    },
    getAllTransactionTags: () => {
        return axios.get('v1/transaction/tags/list.json');
    },
    getTransactionTag: ({ id }) => {
        return axios.get('v1/transaction/tags/get.json?id=' + id);
    },
    addTransactionTag: ({ name }) => {
        return axios.post('v1/transaction/tags/add.json', {
            name
        });
    },
    modifyTransactionTag: ({ id, name }) => {
        return axios.post('v1/transaction/tags/modify.json', {
            id,
            name
        });
    },
    hideTransactionTag: ({ id, hidden }) => {
        return axios.post('v1/transaction/tags/hide.json', {
            id,
            hidden
        });
    },
    moveTransactionTag: ({ newDisplayOrders }) => {
        return axios.post('v1/transaction/tags/move.json', {
            newDisplayOrders,
        });
    },
    deleteTransactionTag: ({ id }) => {
        return axios.post('v1/transaction/tags/delete.json', {
            id
        });
    },
    getLatestExchangeRates: ({ ignoreError }) => {
        return axios.get('v1/exchange_rates/latest.json', {
            ignoreError: !!ignoreError
        });
    },
};
