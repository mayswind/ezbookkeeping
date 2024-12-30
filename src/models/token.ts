export const TOKEN_CLI_USER_AGENT: string = 'ezbookkeeping Cli';

export interface TokenInfoResponse {
    readonly tokenId: string;
    readonly tokenType: number;
    readonly userAgent: string;
    readonly lastSeen: number;
    readonly isCurrent: boolean;
}

export class SessionInfo {
    public readonly tokenId: string;
    public readonly isCurrent: boolean;
    public readonly deviceType: string;
    public readonly deviceInfo: string;
    public readonly createdByCli: boolean;
    public readonly lastSeen: number;

    private constructor(tokenId: string, isCurrent: boolean, deviceType: string, deviceInfo: string, createdByCli: boolean, lastSeen: number) {
        this.tokenId = tokenId;
        this.isCurrent = isCurrent;
        this.deviceType = deviceType;
        this.deviceInfo = deviceInfo;
        this.createdByCli = createdByCli;
        this.lastSeen = lastSeen;
    }

    public static of(tokenId: string, isCurrent: boolean, deviceType: string, deviceInfo: string, createdByCli: boolean, lastSeen: number): SessionInfo {
        return new SessionInfo(tokenId, isCurrent, deviceType, deviceInfo, createdByCli, lastSeen);
    }
}
