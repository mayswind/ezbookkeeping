export interface TwoFactorEnableConfirmRequest {
    readonly secret: string;
    readonly passcode: string;
}

export interface TwoFactorEnableResponse {
    readonly secret: string;
    readonly qrcode: string;
}

export interface TwoFactorEnableConfirmResponse {
    readonly token?: string;
    readonly recoveryCodes: string[];
}

export interface TwoFactorDisableRequest {
    readonly password?: string;
}

export interface TwoFactorRegenerateRecoveryCodeRequest {
    readonly password?: string;
}

export interface TwoFactorStatusResponse {
    readonly enable: boolean;
    readonly createdAt?: number;
}
