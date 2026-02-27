import { UAParser } from 'ua-parser-js';

import {
    type TokenInfoResponse,
    TOKEN_TYPE_API,
    TOKEN_TYPE_MCP,
    SessionDeviceType,
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
    const uaParseRet = new UAParser(ua).getResult();

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
    const isCreateForAPI = token.tokenType === TOKEN_TYPE_API;
    const isCreateForMCP = token.tokenType === TOKEN_TYPE_MCP;
    const uaInfo = parseUserAgent(token.userAgent);
    let deviceType: SessionDeviceType = SessionDeviceType.Default;
    let deviceName: string = 'Other Device';

    if (isCreateForAPI) {
        deviceType = SessionDeviceType.Api;
        deviceName = 'API Token';
    } else if (isCreateForMCP) {
        deviceType = SessionDeviceType.MCP;
        deviceName = 'MCP Token';
    } else {
        if (uaInfo && uaInfo.device) {
            if (uaInfo.device.type === 'mobile') {
                deviceType = SessionDeviceType.Phone;
            } else if (uaInfo.device.type === 'wearable') {
                deviceType = SessionDeviceType.Wearable;
            } else if (uaInfo.device.type === 'tablet') {
                deviceType = SessionDeviceType.Tablet;
            } else if (uaInfo.device.type === 'smarttv') {
                deviceType = SessionDeviceType.TV;
            } else {
                deviceType = SessionDeviceType.Default;
            }
        } else {
            deviceType = SessionDeviceType.Default;
        }

        if (token.isCurrent) {
            deviceName = 'Current';
        } else {
            deviceName = 'Other Device';
        }
    }

    return SessionInfo.of(
        token.tokenId,
        token.isCurrent,
        deviceType,
        isCreateForAPI || isCreateForMCP ? token.userAgent : parseDeviceInfo(uaInfo),
        deviceName,
        token.lastSeen
    );
}
