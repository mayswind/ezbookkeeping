const defaultTimeout = 10000; // 10s
const uploadTimeout = 30000; // 30s
const importTimeout = 1800000; // 1800s
const baseApiUrlPath = '/api';
const baseQrcodePath = '/qrcode';
const baseProxyUrlPath = '/proxy';
const baseAmapApiProxyUrlPath = '/_AMapService';
const apiNotFoundErrorCode = 100001;
const validatorErrorCode = 200000;
const userEmailNotVerifiedErrorCode = 201020;
const transactionCannotCreateInThisTimeErrorCode = 205017;
const transactionCannotModifyInThisTimeErrorCode = 205018;
const transactionPictureNotFoundErrorCode = 211001;
const googleMapJavascriptUrl = 'https://maps.googleapis.com/maps/api/js';
const baiduMapJavascriptUrl = 'https://api.map.baidu.com/api?v=3.0';
const amapJavascriptUrl = 'https://webapi.amap.com/maps?v=2.0';

const specifiedApiNotFoundErrors = {
    '/api/register.json': {
        message: 'User registration is disabled'
    }
};

const parameterizedErrors = [
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

export default {
    defaultTimeout: defaultTimeout,
    uploadTimeout: uploadTimeout,
    importTimeout: importTimeout,
    baseApiUrlPath: baseApiUrlPath,
    baseQrcodePath: baseQrcodePath,
    baseProxyUrlPath: baseProxyUrlPath,
    baseAmapApiProxyUrlPath: baseAmapApiProxyUrlPath,
    apiNotFoundErrorCode: apiNotFoundErrorCode,
    validatorErrorCode: validatorErrorCode,
    userEmailNotVerifiedErrorCode: userEmailNotVerifiedErrorCode,
    transactionCannotCreateInThisTimeErrorCode: transactionCannotCreateInThisTimeErrorCode,
    transactionCannotModifyInThisTimeErrorCode: transactionCannotModifyInThisTimeErrorCode,
    transactionPictureNotFoundErrorCode: transactionPictureNotFoundErrorCode,
    specifiedApiNotFoundErrors: specifiedApiNotFoundErrors,
    parameterizedErrors: parameterizedErrors,
    googleMapJavascriptUrl: googleMapJavascriptUrl,
    baiduMapJavascriptUrl: baiduMapJavascriptUrl,
    amapJavascriptUrl: amapJavascriptUrl
}
