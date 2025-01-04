import { isEnableDebug } from './settings.ts';

function logDebug(msg: string, obj?: unknown): void {
    if (isEnableDebug()) {
        if (obj) {
            console.debug('[ezBookkeeping Debug] ' + msg, obj);
        } else {
            console.debug('[ezBookkeeping Debug] ' + msg);
        }
    }
}

function logInfo(msg: string, obj?: unknown): void {
    if (obj) {
        console.info('[ezBookkeeping Info] ' + msg, obj);
    } else {
        console.info('[ezBookkeeping Info] ' + msg);
    }
}

function logWarn(msg: string, obj?: unknown): void {
    if (obj) {
        console.warn('[ezBookkeeping Warn] ' + msg, obj);
    } else {
        console.warn('[ezBookkeeping Warn] ' + msg);
    }
}

function logError(msg: string, obj?: unknown): void {
    if (obj) {
        console.error('[ezBookkeeping Error] ' + msg, obj);
    } else {
        console.error('[ezBookkeeping Error] ' + msg);
    }
}

export default {
    debug: logDebug,
    info: logInfo,
    warn: logWarn,
    error: logError
};
