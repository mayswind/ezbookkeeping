import type { TransactionCategoryCreateBatchRequest } from './transaction_category.ts';

export interface UserBasicInfo {
    readonly username: string;
    readonly email: string;
    readonly nickname: string;
    readonly avatar: string;
    readonly avatarProvider?: string;
    readonly defaultAccountId: string;
    readonly transactionEditScope: number;
    readonly language: string;
    readonly defaultCurrency: string;
    readonly firstDayOfWeek: number;
    readonly longDateFormat: number;
    readonly shortDateFormat: number;
    readonly longTimeFormat: number;
    readonly shortTimeFormat: number;
    readonly decimalSeparator: number;
    readonly digitGroupingSymbol: number;
    readonly digitGrouping: number;
    readonly currencyDisplayType: number;
    readonly expenseAmountColor: number;
    readonly incomeAmountColor: number;
    readonly emailVerified: boolean;
}

export interface UserLoginRequest {
    readonly loginName: string;
    readonly password: string;
}

export interface UserRegisterRequest {
    readonly username: string;
    readonly email: string;
    readonly nickname: string;
    readonly password: string;
    readonly language: string;
    readonly defaultCurrency: string;
    readonly firstDayOfWeek: number;
    readonly categories?: TransactionCategoryCreateBatchRequest;
}

export interface UserVerifyEmailResponse {
    readonly newToken?: string;
    readonly user: UserBasicInfo;
    readonly notificationContent?: string;
}

export interface UserResendVerifyEmailRequest {
    readonly email: string;
    readonly password: string;
}

export interface UserProfileUpdateRequest {
    readonly email?: string;
    readonly nickname?: string;
    readonly password?: string;
    readonly oldPassword?: string;
    readonly defaultAccountId?: string;
    readonly transactionEditScope?: number;
    readonly language?: string;
    readonly defaultCurrency?: string;
    readonly firstDayOfWeek?: number;
    readonly longDateFormat?: number;
    readonly shortDateFormat?: number;
    readonly longTimeFormat?: number;
    readonly shortTimeFormat?: number;
    readonly decimalSeparator?: number;
    readonly digitGroupingSymbol?: number;
    readonly digitGrouping?: number;
    readonly currencyDisplayType?: number;
    readonly expenseAmountColor?: number;
    readonly incomeAmountColor?: number;
}

export interface UserProfileUpdateResponse {
    readonly user: UserBasicInfo;
    readonly newToken?: string;
}

export interface UserProfileResponse extends UserBasicInfo {
    readonly lastLoginAt: number;
}
