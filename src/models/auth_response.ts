import type { UserBasicInfo } from './user.ts';

export interface AuthResponse {
    readonly token: string;
    readonly need2FA: boolean;
    readonly user?: UserBasicInfo;
    readonly notificationContent?: string;
}

export interface RegisterResponse extends AuthResponse {
    readonly needVerifyEmail: boolean;
    readonly presetCategoriesSaved: boolean;
}
