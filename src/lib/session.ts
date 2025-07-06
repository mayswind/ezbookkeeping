import uaParser from 'ua-parser-js';

import {
    TOKEN_TYPE_MCP,
    TOKEN_CLI_USER_AGENT,
    type TokenInfoResponse,
    SessionInfo
} from '@/models/token.ts';

interface UserAgentInfo {
    readonly device: {
        readonly vendor?: string;
        readonly model?: string;
        readonly type?: string;
    };
    readonly os: {
        readonly name?: string;
        readonly version?: string;
    };
    readonly browser: {
        readonly name?: string;
        readonly version?: string;
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
    return ua === TOKEN_CLI_USER_AGENT;
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
    const isCreateForMCP = token.tokenType === TOKEN_TYPE_MCP;
    const isCreatedByCli = isSessionUserAgentCreatedByCli(token.userAgent);
    const uaInfo = parseUserAgent(token.userAgent);
    let deviceType = '';

    if (isCreateForMCP) {
        deviceType = 'mcp';
    } else if (isCreatedByCli) {
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
        isCreateForMCP || isCreatedByCli ? token.userAgent : parseDeviceInfo(uaInfo),
        isCreatedByCli,
        token.lastSeen
    );
}
