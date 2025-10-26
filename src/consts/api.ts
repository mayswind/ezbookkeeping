export const BASE_API_URL_PATH: string = '/api';
export const BASE_QRCODE_PATH: string = '/qrcode';
export const BASE_PROXY_URL_PATH: string = '/proxy';
export const BASE_AMAP_API_PROXY_URL_PATH: string = '/_AMapService';

export const DEFAULT_API_TIMEOUT: number = 10000; // 10s
export const DEFAULT_UPLOAD_API_TIMEOUT: number = 30000; // 30s
export const DEFAULT_EXPORT_API_TIMEOUT: number = 180000; // 180s
export const DEFAULT_IMPORT_API_TIMEOUT: number = 1800000; // 1800s
export const DEFAULT_CLEAR_ALL_TRANSACTIONS_API_TIMEOUT: number = 1800000; // 1800s
export const DEFAULT_LLM_API_TIMEOUT: number = 600000; // 600s

export const GOOGLE_MAP_JAVASCRIPT_URL: string = 'https://maps.googleapis.com/maps/api/js';
export const BAIDU_MAP_JAVASCRIPT_URL: string = 'https://api.map.baidu.com/api?v=3.0';
export const AMAP_JAVASCRIPT_URL: string = 'https://webapi.amap.com/maps?v=2.0';

export enum KnownErrorCode {
    ApiNotFound = 100001,
    ValidatorError = 200000,
    UserEmailNotVerified = 201020,
    TwoFactorAuthorizationPasscodeEmpty = 203005,
    TransactionCannotCreateInThisTime = 205017,
    TransactionCannotModifyInThisTime = 205018,
    TransactionPictureNotFound = 211001
}

export interface SpecifiedApiError {
    readonly message: string;
}

export const SPECIFIED_API_NOT_FOUND_ERRORS: Record<string, SpecifiedApiError> = {
    '/api/register.json': {
        message: 'User registration is disabled'
    }
};

export interface ParameterizedError {
    readonly localeKey: string;
    readonly regex: RegExp;
    readonly parameters: {
        readonly field: string;
        readonly localized: boolean;
    }[];
}

export const PARAMETERIZED_ERRORS: ParameterizedError[] = [
    {
        localeKey: 'parameter invalid',
        regex: /^parameter "(\w+)" is invalid$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }]
    },
    {
        localeKey: 'parameter required',
        regex: /^parameter "(\w+)" is required$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }]
    },
    {
        localeKey: 'parameter too large',
        regex: /^parameter "(\w+)" must be less than (\d+)$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }, {
            field: 'number',
            localized: false
        }]
    },
    {
        localeKey: 'parameter too long',
        regex: /^parameter "(\w+)" must be less than (\d+) characters$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }, {
            field: 'length',
            localized: false
        }]
    },
    {
        localeKey: 'parameter too small',
        regex: /^parameter "(\w+)" must be more than (\d+)$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }, {
            field: 'number',
            localized: false
        }]
    },
    {
        localeKey: 'parameter too short',
        regex: /^parameter "(\w+)" must be more than (\d+) characters$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }, {
            field: 'length',
            localized: false
        }]
    },
    {
        localeKey: 'parameter length not equal',
        regex: /^parameter "(\w+)" length is not equal to (\d+)$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }, {
            field: 'length',
            localized: false
        }]
    },
    {
        localeKey: 'parameter cannot be blank',
        regex: /^parameter "(\w+)" cannot be blank$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }]
    },
    {
        localeKey: 'parameter invalid username format',
        regex: /^parameter "(\w+)" is invalid username format$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }]
    },
    {
        localeKey: 'parameter invalid email format',
        regex: /^parameter "(\w+)" is invalid email format$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }]
    },
    {
        localeKey: 'parameter invalid currency',
        regex: /^parameter "(\w+)" is invalid currency$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }]
    },
    {
        localeKey: 'parameter invalid color',
        regex: /^parameter "(\w+)" is invalid color$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }]
    },
    {
        localeKey: 'parameter invalid amount filter',
        regex: /^parameter "(\w+)" is invalid amount filter$/,
        parameters: [{
            field: 'parameter',
            localized: true
        }]
    }
];
