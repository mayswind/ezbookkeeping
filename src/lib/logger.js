import settings from './settings.js';

function logDebug(msg, obj) {
    if (settings.isEnableDebug()) {
        if (obj) {
            console.debug('[lab Debug] ' + msg, obj);
        } else {
            console.debug('[lab Debug] ' + msg);
        }
    }
}

function logInfo(msg, obj) {
    if (obj) {
        console.info('[lab Info] ' + msg, obj);
    } else {
        console.info('[lab Info] ' + msg);
    }
}

function logWarn(msg, obj) {
    if (obj) {
        console.warn('[lab Warn] ' + msg, obj);
    } else {
        console.warn('[lab Warn] ' + msg);
    }
}

function logError(msg, obj) {
    if (obj) {
        console.error('[lab Error] ' + msg, obj);
    } else {
        console.error('[lab Error] ' + msg);
    }
}

export default {
    debug: logDebug,
    info: logInfo,
    warn: logWarn,
    error: logError
};
