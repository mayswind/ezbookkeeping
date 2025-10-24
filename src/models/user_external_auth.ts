export interface UserExternalAuthUnlinkRequest {
    readonly externalAuthType: string;
    readonly password: string;
}

export interface UserExternalAuthInfoResponse {
    readonly externalAuthCategory: string;
    readonly externalAuthType: string;
    readonly linked: boolean;
    readonly externalUsername?: string;
    readonly createdAt?: number;
}
