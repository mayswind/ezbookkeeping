import uaParser from 'ua-parser-js';

import { CliUserAgent, type TokenInfoResponse, SessionInfo } from '@/models/token.ts';

interface UserAgentInfo {
    device: {
        vendor?: string;
        model?: string;
        type?: string;
    };
    os: {
        name?: string;
        version?: string;
    };
    browser: {
        name?: string;
        version?: string;
    };
}

function parseUserAgent(ua: string): UserAgentInfo {
    const uaParseRet = uaParser(ua);

    return {
        device: {
            vendor: uaParseRet.device.vendor,
            model: uaParseRet.device.model,
            type: uaParseRet.device.type
        },
        os: {
            name: uaParseRet.os.name,
            version: uaParseRet.os.version
        },
        browser: {
            name: uaParseRet.browser.name,
            version: uaParseRet.browser.version
        }
    };
}

function isSessionUserAgentCreatedByCli(ua: string): boolean {
    return ua === CliUserAgent;
}

function parseDeviceInfo(uaInfo: UserAgentInfo): string {
    if (!uaInfo) {
        return '';
    }

    let result = '';

    if (uaInfo.device && uaInfo.device.model) {
        result = uaInfo.device.model;
    } else if (uaInfo.os && uaInfo.os.name) {
        result = uaInfo.os.name;

        if (uaInfo.os.version) {
            result += ' ' + uaInfo.os.version;
        }
    }

    if (uaInfo.browser && uaInfo.browser.name) {
        let browserInfo = uaInfo.browser.name;

        if (uaInfo.browser.version) {
            browserInfo += ' ' + uaInfo.browser.version;
        }

        if (result) {
            result += ' (' + browserInfo + ')';
        } else {
            result = browserInfo;
        }
    }

    if (!result) {
        return 'Unknown Device';
    }

    return result;
}

export function parseSessionInfo(token: TokenInfoResponse): SessionInfo {
    const isCreatedByCli = isSessionUserAgentCreatedByCli(token.userAgent);
    const uaInfo = parseUserAgent(token.userAgent);
    let deviceType = '';

    if (isCreatedByCli) {
        deviceType = 'cli';
    } else {
        if (uaInfo && uaInfo.device) {
            if (uaInfo.device.type === 'mobile') {
                deviceType = 'phone';
            } else if (uaInfo.device.type === 'wearable') {
                deviceType = 'wearable';
            } else if (uaInfo.device.type === 'tablet') {
                deviceType = 'tablet';
            } else if (uaInfo.device.type === 'smarttv') {
                deviceType = 'tv';
            } else {
                deviceType = 'default';
            }
        } else {
            deviceType = 'default';
        }
    }

    return SessionInfo.of(
        token.tokenId,
        token.isCurrent,
        deviceType,
        isCreatedByCli ? token.userAgent : parseDeviceInfo(uaInfo),
        isCreatedByCli,
        token.lastSeen
    );
}
