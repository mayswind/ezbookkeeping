import axios from 'axios';

import apiConstants from '@/consts/api.js';
import userState from './userstate.js';
import {
    getGoogleMapAPIKey,
    getBaiduMapAK,
    getAmapApplicationKey
} from './server_settings.js';
import { getTimezoneOffsetMinutes } from './datetime.js';

let needBlockRequest = false;
let blockedRequests = [];

axios.defaults.baseURL = apiConstants.baseApiUrlPath;
axios.defaults.timeout = apiConstants.defaultTimeout;
axios.interceptors.request.use(config => {
    const token = userState.getToken();

    if (token && !config.noAuth) {
        config.headers.Authorization = `Bearer ${token}`;
    }

    config.headers['X-Timezone-Offset'] = getTimezoneOffsetMinutes();

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
            || errorCode === 202005 // current token requires two-factor authorization
            || errorCode === 202006 // current token does not require two-factor authorization
            || errorCode === 202012 // token is empty
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
    register: ({ username, email, nickname, password, language, defaultCurrency, firstDayOfWeek, categories }) => {
        return axios.post('register.json', {
            username,
            email,
            nickname,
            password,
            language,
            defaultCurrency,
            firstDayOfWeek,
            categories
        });
    },
    verifyEmail: ({ token, requestNewToken }) => {
        return axios.post('verify_email/by_token.json?token=' + token, {
            requestNewToken
        }, {
            noAuth: true,
            ignoreError: true
        });
    },
    resendVerifyEmailByUnloginUser: ({ email, password }) => {
        return axios.post('verify_email/resend.json', {
            email,
            password
        });
    },
    requestResetPassword: ({ email }) => {
        return axios.post('forget_password/request.json', {
            email
        });
    },
    resetPassword: ({ email, token, password }) => {
        return axios.post('forget_password/reset/by_token.json?token=' + token, {
            email,
            password
        }, {
            noAuth: true,
            ignoreError: true
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
    updateProfile: ({ email, nickname, password, oldPassword, defaultAccountId, transactionEditScope, language, defaultCurrency, firstDayOfWeek, longDateFormat, shortDateFormat, longTimeFormat, shortTimeFormat, decimalSeparator, digitGroupingSymbol, digitGrouping, currencyDisplayType, expenseAmountColor, incomeAmountColor }) => {
        return axios.post('v1/users/profile/update.json', {
            email,
            nickname,
            password,
            oldPassword,
            defaultAccountId,
            transactionEditScope,
            language,
            defaultCurrency,
            firstDayOfWeek,
            longDateFormat,
            shortDateFormat,
            longTimeFormat,
            shortTimeFormat,
            decimalSeparator,
            digitGroupingSymbol,
            digitGrouping,
            currencyDisplayType,
            expenseAmountColor,
            incomeAmountColor
        });
    },
    resendVerifyEmailByLoginedUser: () => {
        return axios.post('v1/users/verify_email/resend.json');
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
    getUserDataStatistics: () => {
        return axios.get('v1/data/statistics.json');
    },
    getExportedUserData: (fileType) => {
        if (fileType === 'csv') {
            return axios.get('v1/data/export.csv');
        } else if (fileType === 'tsv') {
            return axios.get('v1/data/export.tsv');
        } else {
            return Promise.reject('Parameter Invalid');
        }
    },
    clearData: ({ password }) => {
        return axios.post('v1/data/clear.json', {
            password
        });
    },
    getAllAccounts: ({ visibleOnly }) => {
        return axios.get('v1/accounts/list.json?visible_only=' + !!visibleOnly);
    },
    getAccount: ({ id }) => {
        return axios.get('v1/accounts/get.json?id=' + id);
    },
    addAccount: ({ category, type, name, icon, color, currency, balance, comment, subAccounts, clientSessionId }) => {
        return axios.post('v1/accounts/add.json', {
            category,
            type,
            name,
            icon,
            color,
            currency,
            balance,
            comment,
            subAccounts,
            clientSessionId
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
    getTransactions: ({ maxTime, minTime, count, page, withCount, type, categoryIds, accountIds, tagIds, amountFilter, keyword }) => {
        amountFilter = encodeURIComponent(amountFilter);
        keyword = encodeURIComponent(keyword);
        return axios.get(`v1/transactions/list.json?max_time=${maxTime}&min_time=${minTime}&type=${type}&category_ids=${categoryIds}&account_ids=${accountIds}&tag_ids=${tagIds}&amount_filter=${amountFilter}&keyword=${keyword}&count=${count}&page=${page}&with_count=${withCount}&trim_account=true&trim_category=true&trim_tag=true`);
    },
    getAllTransactionsByMonth: ({ year, month, type, categoryIds, accountIds, tagIds, amountFilter, keyword }) => {
        amountFilter = encodeURIComponent(amountFilter);
        keyword = encodeURIComponent(keyword);
        return axios.get(`v1/transactions/list/by_month.json?year=${year}&month=${month}&type=${type}&category_ids=${categoryIds}&account_ids=${accountIds}&tag_ids=${tagIds}&amount_filter=${amountFilter}&keyword=${keyword}&trim_account=true&trim_category=true&trim_tag=true`);
    },
    getTransactionStatistics: ({ startTime, endTime, useTransactionTimezone }) => {
        const queryParams = [];

        if (startTime) {
            queryParams.push(`start_time=${startTime}`);
        }

        if (endTime) {
            queryParams.push(`end_time=${endTime}`);
        }

        return axios.get(`v1/transactions/statistics.json?use_transaction_timezone=${useTransactionTimezone}` + (queryParams.length ? '&' + queryParams.join('&') : ''));
    },
    getTransactionStatisticsTrends: ({ startYearMonth, endYearMonth, useTransactionTimezone }) => {
        const queryParams = [];

        if (startYearMonth) {
            queryParams.push(`start_year_month=${startYearMonth}`);
        }

        if (endYearMonth) {
            queryParams.push(`end_year_month=${endYearMonth}`);
        }

        return axios.get(`v1/transactions/statistics/trends.json?use_transaction_timezone=${useTransactionTimezone}` + (queryParams.length ? '&' + queryParams.join('&') : ''));
    },
    getTransactionAmounts: ({ useTransactionTimezone, today, thisWeek, thisMonth, thisYear, lastMonth, monthBeforeLastMonth, monthBeforeLast2Months, monthBeforeLast3Months, monthBeforeLast4Months, monthBeforeLast5Months, monthBeforeLast6Months, monthBeforeLast7Months, monthBeforeLast8Months, monthBeforeLast9Months, monthBeforeLast10Months }) => {
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

        if (lastMonth) {
            queryParams.push(`lastMonth_${lastMonth.startTime}_${lastMonth.endTime}`);
        }

        if (monthBeforeLastMonth) {
            queryParams.push(`monthBeforeLastMonth_${monthBeforeLastMonth.startTime}_${monthBeforeLastMonth.endTime}`);
        }

        if (monthBeforeLast2Months) {
            queryParams.push(`monthBeforeLast2Months_${monthBeforeLast2Months.startTime}_${monthBeforeLast2Months.endTime}`);
        }

        if (monthBeforeLast3Months) {
            queryParams.push(`monthBeforeLast3Months_${monthBeforeLast3Months.startTime}_${monthBeforeLast3Months.endTime}`);
        }

        if (monthBeforeLast4Months) {
            queryParams.push(`monthBeforeLast4Months_${monthBeforeLast4Months.startTime}_${monthBeforeLast4Months.endTime}`);
        }

        if (monthBeforeLast5Months) {
            queryParams.push(`monthBeforeLast5Months_${monthBeforeLast5Months.startTime}_${monthBeforeLast5Months.endTime}`);
        }

        if (monthBeforeLast6Months) {
            queryParams.push(`monthBeforeLast6Months_${monthBeforeLast6Months.startTime}_${monthBeforeLast6Months.endTime}`);
        }

        if (monthBeforeLast7Months) {
            queryParams.push(`monthBeforeLast7Months_${monthBeforeLast7Months.startTime}_${monthBeforeLast7Months.endTime}`);
        }

        if (monthBeforeLast8Months) {
            queryParams.push(`monthBeforeLast8Months_${monthBeforeLast8Months.startTime}_${monthBeforeLast8Months.endTime}`);
        }

        if (monthBeforeLast9Months) {
            queryParams.push(`monthBeforeLast9Months_${monthBeforeLast9Months.startTime}_${monthBeforeLast9Months.endTime}`);
        }

        if (monthBeforeLast10Months) {
            queryParams.push(`monthBeforeLast10Months_${monthBeforeLast10Months.startTime}_${monthBeforeLast10Months.endTime}`);
        }

        return axios.get(`v1/transactions/amounts.json?use_transaction_timezone=${useTransactionTimezone}` + (queryParams.length ? '&query=' + queryParams.join('|') : ''));
    },
    getTransaction: ({ id }) => {
        return axios.get(`v1/transactions/get.json?id=${id}&trim_account=true&trim_category=true&trim_tag=true`);
    },
    addTransaction: ({ type, categoryId, time, sourceAccountId, destinationAccountId, sourceAmount, destinationAmount, hideAmount, tagIds, comment, geoLocation, utcOffset, clientSessionId }) => {
        return axios.post('v1/transactions/add.json', {
            type,
            categoryId,
            time,
            sourceAccountId,
            destinationAccountId,
            sourceAmount,
            destinationAmount,
            hideAmount,
            tagIds,
            comment,
            geoLocation,
            utcOffset,
            clientSessionId
        });
    },
    modifyTransaction: ({ id, type, categoryId, time, sourceAccountId, destinationAccountId, sourceAmount, destinationAmount, hideAmount, tagIds, comment, geoLocation, utcOffset }) => {
        return axios.post('v1/transactions/modify.json', {
            id,
            type,
            categoryId,
            time,
            sourceAccountId,
            destinationAccountId,
            sourceAmount,
            destinationAmount,
            hideAmount,
            tagIds,
            comment,
            geoLocation,
            utcOffset
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
    addTransactionCategory: ({ name, type, parentId, icon, color, comment, clientSessionId }) => {
        return axios.post('v1/transaction/categories/add.json', {
            name,
            type,
            parentId,
            icon,
            color,
            comment,
            clientSessionId
        });
    },
    addTransactionCategoryBatch: ({ categories }) => {
        return axios.post('v1/transaction/categories/add_batch.json', {
            categories
        });
    },
    modifyTransactionCategory: ({ id, name, parentId, icon, color, comment, hidden }) => {
        return axios.post('v1/transaction/categories/modify.json', {
            id,
            name,
            parentId,
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
    generateQrCodeUrl: (qrCodeName) => {
        return `${apiConstants.baseQrcodePath}/${qrCodeName}.png`;
    },
    generateMapProxyTileImageUrl: (mapProvider, language) => {
        const token = userState.getToken();
        let url = `${apiConstants.baseProxyUrlPath}/map/tile/{z}/{x}/{y}.png?provider=${mapProvider}&token=${token}`;

        if (language) {
            url = url + `&language=${language}`;
        }

        return url;
    },
    generateMapProxyAnnotationImageUrl: (mapProvider, language) => {
        const token = userState.getToken();
        let url = `${apiConstants.baseProxyUrlPath}/map/annotation/{z}/{x}/{y}.png?provider=${mapProvider}&token=${token}`;

        if (language) {
            url = url + `&language=${language}`;
        }

        return url;
    },
    generateGoogleMapJavascriptUrl: (language, callbackFnName) => {
        let url = `${apiConstants.googleMapJavascriptUrl}?key=${getGoogleMapAPIKey()}&libraries=core,marker&callback=${callbackFnName}`;

        if (language) {
            url = url + `&language=${language}`;
        }

        return url;
    },
    generateBaiduMapJavascriptUrl: (callbackFnName) => {
        return `${apiConstants.baiduMapJavascriptUrl}&ak=${getBaiduMapAK()}&callback=${callbackFnName}`;
    },
    generateAmapJavascriptUrl: (callbackFnName) => {
        return `${apiConstants.amapJavascriptUrl}&key=${getAmapApplicationKey()}&plugin=AMap.ToolBar&callback=${callbackFnName}`;
    },
    generateAmapApiInternalProxyUrl: () => {
        return `${window.location.origin}${apiConstants.baseAmapApiProxyUrlPath}`;
    }
};
