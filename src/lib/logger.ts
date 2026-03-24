import { isEnableDebug } from './settings.ts';

function logDebug(msg: string, obj?: unknown): void {
    if (isEnableDebug()) {
        if (obj) {
            console.debug('[oscar Debug] ' + msg, obj);
        } else {
            console.debug('[oscar Debug] ' + msg);
        }
    }
}

function logInfo(msg: string, obj?: unknown): void {
    if (obj) {
        console.info('[oscar Info] ' + msg, obj);
    } else {
        console.info('[oscar Info] ' + msg);
    }
}

function logWarn(msg: string, obj?: unknown): void {
    if (obj) {
        console.warn('[oscar Warn] ' + msg, obj);
    } else {
        console.warn('[oscar Warn] ' + msg);
    }
}

function logError(msg: string, obj?: unknown): void {
    if (obj) {
        console.error('[oscar Error] ' + msg, obj);
    } else {
        console.error('[oscar Error] ' + msg);
    }
}

export default {
    debug: logDebug,
    info: logInfo,
    warn: logWarn,
    error: logError
};
