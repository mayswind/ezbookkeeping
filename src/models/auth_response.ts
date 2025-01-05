import type { UserBasicInfo } from './user.ts';

export interface AuthResponse {
    token: string;
    need2FA: boolean;
    user?: UserBasicInfo;
    notificationContent?: string;
}

export interface RegisterResponse extends AuthResponse {
    needVerifyEmail: boolean;
    presetCategoriesSaved: boolean;
}
