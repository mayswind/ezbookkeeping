import axios, { type AxiosRequestConfig, type AxiosRequestHeaders, type AxiosResponse } from 'axios';

import {
    BASE_API_URL_PATH,
    BASE_QRCODE_PATH,
    BASE_PROXY_URL_PATH,
    BASE_AMAP_API_PROXY_URL_PATH,
    DEFAULT_API_TIMEOUT,
    DEFAULT_UPLOAD_API_TIMEOUT,
    DEFAULT_IMPORT_API_TIMEOUT,
    GOOGLE_MAP_JAVASCRIPT_URL,
    BAIDU_MAP_JAVASCRIPT_URL,
    AMAP_JAVASCRIPT_URL
} from '@/consts/api.ts';

import type {
    AccountCreateRequest,
    AccountModifyRequest,
    AccountInfoResponse,
    AccountHideRequest,
    AccountMoveRequest,
    AccountDeleteRequest
} from '@/models/account.ts';
import type {
    AuthResponse,
    RegisterResponse
} from '@/models/auth_response.ts';
import type {
    ClearDataRequest,
    DataStatisticsResponse
} from '@/models/data_management.ts';
import type {
    LatestExchangeRateResponse
} from '@/models/exchange_rate.ts';
import type {
    ForgetPasswordRequest
} from '@/models/forget_password.ts';
import type {
    ImportTransactionResponsePageWrapper
} from '@/models/imported_transaction.ts';
import type {
    TransactionCreateRequest,
    TransactionModifyRequest,
    TransactionDeleteRequest,
    TransactionImportRequest,
    TransactionListByMaxTimeRequest,
    TransactionListInMonthByPageRequest,
    TransactionInfoResponse,
    TransactionInfoPageWrapperResponse,
    TransactionInfoPageWrapperResponse2,
    TransactionStatisticRequest,
    TransactionStatisticResponse,
    TransactionStatisticTrendsRequest,
    TransactionStatisticTrendsItem,
    TransactionAmountsRequestParams,
    TransactionAmountsResponseItem
} from '@/models/transaction.ts';
import {
    TransactionAmountsRequest
} from '@/models/transaction.ts';
import type {
    TransactionCategoryCreateRequest,
    TransactionCategoryCreateBatchRequest,
    TransactionCategoryModifyRequest,
    TransactionCategoryHideRequest,
    TransactionCategoryMoveRequest,
    TransactionCategoryDeleteRequest,
    TransactionCategoryInfoResponse
} from '@/models/transaction_category.ts';
import type {
    TransactionPictureUnusedDeleteRequest,
    TransactionPictureInfoBasicResponse
} from '@/models/transaction_picture_info.ts';
import type {
    TransactionTagCreateRequest,
    TransactionTagModifyRequest,
    TransactionTagHideRequest,
    TransactionTagMoveRequest,
    TransactionTagDeleteRequest,
    TransactionTagInfoResponse
} from '@/models/transaction_tag.ts';
import type {
    TransactionTemplateCreateRequest,
    TransactionTemplateModifyRequest,
    TransactionTemplateHideRequest,
    TransactionTemplateMoveRequest,
    TransactionTemplateDeleteRequest,
    TransactionTemplateInfoResponse
} from '@/models/transaction_template.ts';
import type {
    TokenRefreshResponse,
    TokenInfoResponse
} from '@/models/token.ts';
import type {
    TwoFactorEnableConfirmRequest,
    TwoFactorEnableResponse,
    TwoFactorEnableConfirmResponse,
    TwoFactorDisableRequest,
    TwoFactorRegenerateRecoveryCodeRequest,
    TwoFactorStatusResponse
} from '@/models/two_factor.ts';
import type {
    UserLoginRequest,
    UserRegisterRequest,
    UserVerifyEmailResponse,
    UserResendVerifyEmailRequest,
    UserProfileResponse,
    UserProfileUpdateRequest,
    UserProfileUpdateResponse
} from '@/models/user.ts';

import userState from './userstate.ts';

import {
    isDefined,
    isBoolean
} from './common.ts';
import {
    getGoogleMapAPIKey,
    getBaiduMapAK,
    getAmapApplicationKey,
    getExchangeRatesRequestTimeout
} from './server_settings.ts';
import { getTimezoneOffsetMinutes } from './datetime.ts';
import { generateRandomUUID } from './misc.ts';

interface ApiRequestConfig extends AxiosRequestConfig {
    headers: AxiosRequestHeaders;
    noAuth?: boolean;
    ignoreBlocked?: boolean;
    ignoreError?: boolean;
    timeout?: number;
}

interface ApiResponse<T> {
    readonly success: boolean;
    readonly result: T;
}

type ApiResponsePromise<T> = Promise<AxiosResponse<ApiResponse<T>>>;

let needBlockRequest = false;
const blockedRequests: ((token: string | undefined) => void)[] = [];

axios.defaults.baseURL = BASE_API_URL_PATH;
axios.defaults.timeout = DEFAULT_API_TIMEOUT;
axios.interceptors.request.use((config: ApiRequestConfig) => {
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
    setLocale: (locale: string) => {
        axios.defaults.headers.common['Accept-Language'] = locale;
    },
    authorize: (data: UserLoginRequest): ApiResponsePromise<AuthResponse> => {
        return axios.post<ApiResponse<AuthResponse>>('authorize.json', data);
    },
    authorize2FA: (req: { passcode: string, token: string }): ApiResponsePromise<AuthResponse> => {
        return axios.post<ApiResponse<AuthResponse>>('2fa/authorize.json', {
            passcode: req.passcode
        }, {
            headers: {
                Authorization: `Bearer ${req.token}`
            }
        });
    },
    authorize2FAByBackupCode: (req: { recoveryCode: string, token: string }): ApiResponsePromise<AuthResponse> => {
        return axios.post<ApiResponse<AuthResponse>>('2fa/recovery.json', {
            recoveryCode: req.recoveryCode
        }, {
            headers: {
                Authorization: `Bearer ${req.token}`
            }
        });
    },
    register: (req: UserRegisterRequest): ApiResponsePromise<RegisterResponse> => {
        return axios.post<ApiResponse<RegisterResponse>>('register.json', req);
    },
    verifyEmail: (req: { token: string, requestNewToken: boolean }): ApiResponsePromise<UserVerifyEmailResponse> => {
        return axios.post<ApiResponse<UserVerifyEmailResponse>>('verify_email/by_token.json?token=' + req.token, {
            requestNewToken: req.requestNewToken
        }, {
            noAuth: true,
            ignoreError: true
        } as ApiRequestConfig);
    },
    resendVerifyEmailByUnloginUser: (req: UserResendVerifyEmailRequest): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('verify_email/resend.json', req);
    },
    requestResetPassword: (req: ForgetPasswordRequest): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('forget_password/request.json', req);
    },
    resetPassword: (req: { email: string, token: string, password: string }): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('forget_password/reset/by_token.json?token=' + req.token, {
            email: req.email,
            password: req.password
        }, {
            noAuth: true,
            ignoreError: true
        } as ApiRequestConfig);
    },
    logout: (): ApiResponsePromise<boolean> => {
        return axios.get<ApiResponse<boolean>>('logout.json');
    },
    refreshToken: (): ApiResponsePromise<TokenRefreshResponse> => {
        return new Promise((resolve) => {
            needBlockRequest = true;

            axios.post<ApiResponse<TokenRefreshResponse>>('v1/tokens/refresh.json', {}, {
                ignoreBlocked: true
            } as ApiRequestConfig).then(response => {
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
    getTokens: (): ApiResponsePromise<TokenInfoResponse[]> => {
        return axios.get<ApiResponse<TokenInfoResponse[]>>('v1/tokens/list.json');
    },
    revokeToken: (req: { tokenId: string, ignoreError: boolean }): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('v1/tokens/revoke.json', {
            tokenId: req.tokenId
        }, {
            ignoreError: req.ignoreError
        } as ApiRequestConfig);
    },
    revokeAllTokens: (): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('v1/tokens/revoke_all.json');
    },
    getProfile: (): ApiResponsePromise<UserProfileResponse> => {
        return axios.get<ApiResponse<UserProfileResponse>>('v1/users/profile/get.json');
    },
    updateProfile: (req: UserProfileUpdateRequest): ApiResponsePromise<UserProfileUpdateResponse> => {
        return axios.post<ApiResponse<UserProfileUpdateResponse>>('v1/users/profile/update.json', req);
    },
    updateAvatar: (req: { avatarFile: unknown }): ApiResponsePromise<UserProfileResponse> => {
        return axios.postForm<ApiResponse<UserProfileResponse>>('v1/users/avatar/update.json', {
            avatar: req.avatarFile
        }, {
            timeout: DEFAULT_UPLOAD_API_TIMEOUT
        });
    },
    removeAvatar: (): ApiResponsePromise<UserProfileResponse> => {
        return axios.post<ApiResponse<UserProfileResponse>>('v1/users/avatar/remove.json');
    },
    resendVerifyEmailByLoginedUser: (): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('v1/users/verify_email/resend.json');
    },
    get2FAStatus: (): ApiResponsePromise<TwoFactorStatusResponse> => {
        return axios.get<ApiResponse<TwoFactorStatusResponse>>('v1/users/2fa/status.json');
    },
    enable2FA: (): ApiResponsePromise<TwoFactorEnableResponse> => {
        return axios.post<ApiResponse<TwoFactorEnableResponse>>('v1/users/2fa/enable/request.json');
    },
    confirmEnable2FA: (req: TwoFactorEnableConfirmRequest): ApiResponsePromise<TwoFactorEnableConfirmResponse> => {
        return axios.post<ApiResponse<TwoFactorEnableConfirmResponse>>('v1/users/2fa/enable/confirm.json', req);
    },
    disable2FA: (req: TwoFactorDisableRequest): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('v1/users/2fa/disable.json', req);
    },
    regenerate2FARecoveryCode: (req: TwoFactorRegenerateRecoveryCodeRequest): ApiResponsePromise<TwoFactorEnableConfirmResponse> => {
        return axios.post<ApiResponse<TwoFactorEnableConfirmResponse>>('v1/users/2fa/recovery/regenerate.json', req);
    },
    getUserDataStatistics: (): ApiResponsePromise<DataStatisticsResponse> => {
        return axios.get<ApiResponse<DataStatisticsResponse>>('v1/data/statistics.json');
    },
    getExportedUserData: (fileType: string): Promise<AxiosResponse<unknown>> => {
        if (fileType === 'csv') {
            return axios.get<unknown>('v1/data/export.csv');
        } else if (fileType === 'tsv') {
            return axios.get<unknown>('v1/data/export.tsv');
        } else {
            return Promise.reject('Parameter Invalid');
        }
    },
    clearData: (req: ClearDataRequest): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('v1/data/clear.json', req);
    },
    getAllAccounts: (req: { visibleOnly: boolean }): ApiResponsePromise<AccountInfoResponse[]> => {
        return axios.get<ApiResponse<AccountInfoResponse[]>>('v1/accounts/list.json?visible_only=' + req.visibleOnly);
    },
    getAccount: (req: { id: string }): ApiResponsePromise<AccountInfoResponse> => {
        return axios.get<ApiResponse<AccountInfoResponse>>('v1/accounts/get.json?id=' + req.id);
    },
    addAccount: (req: AccountCreateRequest): ApiResponsePromise<AccountInfoResponse> => {
        return axios.post<ApiResponse<AccountInfoResponse>>('v1/accounts/add.json', req);
    },
    modifyAccount: (req: AccountModifyRequest): ApiResponsePromise<AccountInfoResponse> => {
        return axios.post<ApiResponse<AccountInfoResponse>>('v1/accounts/modify.json', req);
    },
    hideAccount: (req: AccountHideRequest): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('v1/accounts/hide.json', req);
    },
    moveAccount: (req: AccountMoveRequest): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('v1/accounts/move.json', req);
    },
    deleteAccount: (req: AccountDeleteRequest): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('v1/accounts/delete.json', req);
    },
    getTransactions: (req: TransactionListByMaxTimeRequest): ApiResponsePromise<TransactionInfoPageWrapperResponse> => {
        const amountFilter = encodeURIComponent(req.amountFilter);
        const keyword = encodeURIComponent(req.keyword);
        return axios.get<ApiResponse<TransactionInfoPageWrapperResponse>>(`v1/transactions/list.json?max_time=${req.maxTime}&min_time=${req.minTime}&type=${req.type}&category_ids=${req.categoryIds}&account_ids=${req.accountIds}&tag_ids=${req.tagIds}&tag_filter_type=${req.tagFilterType}&amount_filter=${amountFilter}&keyword=${keyword}&count=${req.count}&page=${req.page}&with_count=${req.withCount}&trim_account=true&trim_category=true&trim_tag=true`);
    },
    getAllTransactionsByMonth: (req: TransactionListInMonthByPageRequest): ApiResponsePromise<TransactionInfoPageWrapperResponse2> => {
        const amountFilter = encodeURIComponent(req.amountFilter);
        const keyword = encodeURIComponent(req.keyword);
        return axios.get<ApiResponse<TransactionInfoPageWrapperResponse2>>(`v1/transactions/list/by_month.json?year=${req.year}&month=${req.month}&type=${req.type}&category_ids=${req.categoryIds}&account_ids=${req.accountIds}&tag_ids=${req.tagIds}&tag_filter_type=${req.tagFilterType}&amount_filter=${amountFilter}&keyword=${keyword}&trim_account=true&trim_category=true&trim_tag=true`);
    },
    getTransactionStatistics: (req: TransactionStatisticRequest): ApiResponsePromise<TransactionStatisticResponse> => {
        const queryParams = [];

        if (req.startTime) {
            queryParams.push(`start_time=${req.startTime}`);
        }

        if (req.endTime) {
            queryParams.push(`end_time=${req.endTime}`);
        }

        if (req.tagIds) {
            queryParams.push(`tag_ids=${req.tagIds}`);
        }

        if (req.tagFilterType) {
            queryParams.push(`tag_filter_type=${req.tagFilterType}`);
        }

        return axios.get<ApiResponse<TransactionStatisticResponse>>(`v1/transactions/statistics.json?use_transaction_timezone=${req.useTransactionTimezone}` + (queryParams.length ? '&' + queryParams.join('&') : ''));
    },
    getTransactionStatisticsTrends: (req: TransactionStatisticTrendsRequest): ApiResponsePromise<TransactionStatisticTrendsItem[]> => {
        const queryParams = [];

        if (req.startYearMonth) {
            queryParams.push(`start_year_month=${req.startYearMonth}`);
        }

        if (req.endYearMonth) {
            queryParams.push(`end_year_month=${req.endYearMonth}`);
        }

        if (req.tagIds) {
            queryParams.push(`tag_ids=${req.tagIds}`);
        }

        if (req.tagFilterType) {
            queryParams.push(`tag_filter_type=${req.tagFilterType}`);
        }

        return axios.get<ApiResponse<TransactionStatisticTrendsItem[]>>(`v1/transactions/statistics/trends.json?use_transaction_timezone=${req.useTransactionTimezone}` + (queryParams.length ? '&' + queryParams.join('&') : ''));
    },
    getTransactionAmounts: (params: TransactionAmountsRequestParams): ApiResponsePromise<TransactionAmountsResponseItem> => {
        const req = TransactionAmountsRequest.of(params);
        return axios.get<ApiResponse<TransactionAmountsResponseItem>>(`v1/transactions/amounts.json?${req.buildQuery()}`);
    },
    getTransaction: (req: { id: string, withPictures: boolean | undefined }): ApiResponsePromise<TransactionInfoResponse> => {
        if (!isDefined(req.withPictures)) {
            req.withPictures = true;
        }

        return axios.get<ApiResponse<TransactionInfoResponse>>(`v1/transactions/get.json?id=${req.id}&with_pictures=${req.withPictures}&trim_account=true&trim_category=true&trim_tag=true`);
    },
    addTransaction: (req: TransactionCreateRequest): ApiResponsePromise<TransactionInfoResponse> => {
        return axios.post<ApiResponse<TransactionInfoResponse>>('v1/transactions/add.json', req);
    },
    modifyTransaction: (req: TransactionModifyRequest): ApiResponsePromise<TransactionInfoResponse> => {
        return axios.post<ApiResponse<TransactionInfoResponse>>('v1/transactions/modify.json', req);
    },
    deleteTransaction: (req: TransactionDeleteRequest): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('v1/transactions/delete.json', req);
    },
    parseImportTransaction: (req: { fileType: string, importFile: unknown }): ApiResponsePromise<ImportTransactionResponsePageWrapper> => {
        return axios.postForm<ApiResponse<ImportTransactionResponsePageWrapper>>('v1/transactions/parse_import.json', {
            fileType: req.fileType,
            file: req.importFile
        }, {
            timeout: DEFAULT_UPLOAD_API_TIMEOUT
        } as ApiRequestConfig);
    },
    importTransactions: (req: TransactionImportRequest): ApiResponsePromise<number> => {
        return axios.post<ApiResponse<number>>('v1/transactions/import.json', req, {
            timeout: DEFAULT_IMPORT_API_TIMEOUT
        } as ApiRequestConfig);
    },
    uploadTransactionPicture: (req: { pictureFile: unknown, clientSessionId: string }): ApiResponsePromise<TransactionPictureInfoBasicResponse> => {
        return axios.postForm<ApiResponse<TransactionPictureInfoBasicResponse>>('v1/transaction/pictures/upload.json', {
            picture: req.pictureFile,
            clientSessionId: req.clientSessionId
        }, {
            timeout: DEFAULT_UPLOAD_API_TIMEOUT
        } as ApiRequestConfig);
    },
    removeUnusedTransactionPicture: (req: TransactionPictureUnusedDeleteRequest): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('v1/transaction/pictures/remove_unused.json', req);
    },
    getAllTransactionCategories: (): ApiResponsePromise<TransactionCategoryInfoResponse[]> => {
        return axios.get<ApiResponse<TransactionCategoryInfoResponse[]>>('v1/transaction/categories/list.json');
    },
    getTransactionCategory: (req: { id: string }): ApiResponsePromise<TransactionCategoryInfoResponse> => {
        return axios.get<ApiResponse<TransactionCategoryInfoResponse>>('v1/transaction/categories/get.json?id=' + req.id);
    },
    addTransactionCategory: (req: TransactionCategoryCreateRequest): ApiResponsePromise<TransactionCategoryInfoResponse> => {
        return axios.post<ApiResponse<TransactionCategoryInfoResponse>>('v1/transaction/categories/add.json', req);
    },
    addTransactionCategoryBatch: (req: TransactionCategoryCreateBatchRequest): ApiResponsePromise<Record<number, TransactionCategoryInfoResponse[]>> => {
        return axios.post<ApiResponse<Record<number, TransactionCategoryInfoResponse[]>>>('v1/transaction/categories/add_batch.json', req);
    },
    modifyTransactionCategory: (req: TransactionCategoryModifyRequest): ApiResponsePromise<TransactionCategoryInfoResponse> => {
        return axios.post<ApiResponse<TransactionCategoryInfoResponse>>('v1/transaction/categories/modify.json', req);
    },
    hideTransactionCategory: (req: TransactionCategoryHideRequest): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('v1/transaction/categories/hide.json', req);
    },
    moveTransactionCategory: (req: TransactionCategoryMoveRequest): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('v1/transaction/categories/move.json', req);
    },
    deleteTransactionCategory: (req: TransactionCategoryDeleteRequest): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('v1/transaction/categories/delete.json', req);
    },
    getAllTransactionTags: (): ApiResponsePromise<TransactionTagInfoResponse[]> => {
        return axios.get<ApiResponse<TransactionTagInfoResponse[]>>('v1/transaction/tags/list.json');
    },
    getTransactionTag: (req: { id: string }): ApiResponsePromise<TransactionTagInfoResponse> => {
        return axios.get<ApiResponse<TransactionTagInfoResponse>>('v1/transaction/tags/get.json?id=' + req.id);
    },
    addTransactionTag: (req: TransactionTagCreateRequest): ApiResponsePromise<TransactionTagInfoResponse> => {
        return axios.post<ApiResponse<TransactionTagInfoResponse>>('v1/transaction/tags/add.json', req);
    },
    modifyTransactionTag: (req: TransactionTagModifyRequest): ApiResponsePromise<TransactionTagInfoResponse> => {
        return axios.post<ApiResponse<TransactionTagInfoResponse>>('v1/transaction/tags/modify.json', req);
    },
    hideTransactionTag: (req: TransactionTagHideRequest): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('v1/transaction/tags/hide.json', req);
    },
    moveTransactionTag: (req: TransactionTagMoveRequest): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('v1/transaction/tags/move.json', req);
    },
    deleteTransactionTag: (req: TransactionTagDeleteRequest): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('v1/transaction/tags/delete.json', req);
    },
    getAllTransactionTemplates: (req: { templateType: number }): ApiResponsePromise<TransactionTemplateInfoResponse[]> => {
        return axios.get<ApiResponse<TransactionTemplateInfoResponse[]>>('v1/transaction/templates/list.json?templateType=' + req.templateType);
    },
    getTransactionTemplate: (req: { id: string }): ApiResponsePromise<TransactionTemplateInfoResponse> => {
        return axios.get<ApiResponse<TransactionTemplateInfoResponse>>('v1/transaction/templates/get.json?id=' + req.id);
    },
    addTransactionTemplate: (req: TransactionTemplateCreateRequest): ApiResponsePromise<TransactionTemplateInfoResponse> => {
        return axios.post<ApiResponse<TransactionTemplateInfoResponse>>('v1/transaction/templates/add.json', req);
    },
    modifyTransactionTemplate: (req: TransactionTemplateModifyRequest): ApiResponsePromise<TransactionTemplateInfoResponse> => {
        return axios.post<ApiResponse<TransactionTemplateInfoResponse>>('v1/transaction/templates/modify.json', req);
    },
    hideTransactionTemplate: (req: TransactionTemplateHideRequest): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('v1/transaction/templates/hide.json', req);
    },
    moveTransactionTemplate: (req: TransactionTemplateMoveRequest): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('v1/transaction/templates/move.json', req);
    },
    deleteTransactionTemplate: (req: TransactionTemplateDeleteRequest): ApiResponsePromise<boolean> => {
        return axios.post<ApiResponse<boolean>>('v1/transaction/templates/delete.json', req);
    },
    getLatestExchangeRates: (param: { ignoreError: boolean }): ApiResponsePromise<LatestExchangeRateResponse> => {
        return axios.get<ApiResponse<LatestExchangeRateResponse>>('v1/exchange_rates/latest.json', {
            ignoreError: param.ignoreError,
            timeout: getExchangeRatesRequestTimeout() || DEFAULT_API_TIMEOUT
        } as ApiRequestConfig);
    },
    generateQrCodeUrl: (qrCodeName: string): string => {
        return `${BASE_QRCODE_PATH}/${qrCodeName}.png`;
    },
    generateMapProxyTileImageUrl: (mapProvider: string, language: string): string => {
        const token = userState.getToken();
        let url = `${BASE_PROXY_URL_PATH}/map/tile/{z}/{x}/{y}.png?provider=${mapProvider}&token=${token}`;

        if (language) {
            url = url + `&language=${language}`;
        }

        return url;
    },
    generateMapProxyAnnotationImageUrl: (mapProvider: string, language: string): string => {
        const token = userState.getToken();
        let url = `${BASE_PROXY_URL_PATH}/map/annotation/{z}/{x}/{y}.png?provider=${mapProvider}&token=${token}`;

        if (language) {
            url = url + `&language=${language}`;
        }

        return url;
    },
    generateGoogleMapJavascriptUrl: (language: string, callbackFnName: string): string => {
        let url = `${GOOGLE_MAP_JAVASCRIPT_URL}?key=${getGoogleMapAPIKey()}&libraries=core,marker&callback=${callbackFnName}`;

        if (language) {
            url = url + `&language=${language}`;
        }

        return url;
    },
    generateBaiduMapJavascriptUrl: (callbackFnName: string): string => {
        return `${BAIDU_MAP_JAVASCRIPT_URL}&ak=${getBaiduMapAK()}&callback=${callbackFnName}`;
    },
    generateAmapJavascriptUrl: (callbackFnName: string): string => {
        return `${AMAP_JAVASCRIPT_URL}&key=${getAmapApplicationKey()}&plugin=AMap.ToolBar&callback=${callbackFnName}`;
    },
    generateAmapApiInternalProxyUrl: (): string => {
        return `${window.location.origin}${BASE_AMAP_API_PROXY_URL_PATH}`;
    },
    getInternalAvatarUrlWithToken(avatarUrl: string, disableBrowserCache: boolean | string): string {
        if (!avatarUrl) {
            return avatarUrl;
        }

        const params = [];
        params.push('token=' + userState.getToken());

        if (disableBrowserCache) {
            if (isBoolean(disableBrowserCache)) {
                params.push('_nocache=' + generateRandomUUID());
            } else {
                params.push('_nocache=' + disableBrowserCache);
            }
        }

        if (avatarUrl.indexOf('?') >= 0) {
            return avatarUrl + '&' + params.join('&');
        } else {
            return avatarUrl + '?' + params.join('&');
        }
    },
    getTransactionPictureUrlWithToken(pictureUrl: string, disableBrowserCache: boolean | string): string {
        if (!pictureUrl) {
            return pictureUrl;
        }

        const params = [];
        params.push('token=' + userState.getToken());

        if (disableBrowserCache) {
            if (isBoolean(disableBrowserCache)) {
                params.push('_nocache=' + generateRandomUUID());
            } else {
                params.push('_nocache=' + disableBrowserCache);
            }
        }

        if (pictureUrl.indexOf('?') >= 0) {
            return pictureUrl + '&' + params.join('&');
        } else {
            return pictureUrl + '?' + params.join('&');
        }
    }
};
