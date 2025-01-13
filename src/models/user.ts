import { LongDateFormat, ShortDateFormat, LongTimeFormat, ShortTimeFormat } from '@/core/datetime.ts';
import { DecimalSeparator, DigitGroupingSymbol, DigitGroupingType } from '@/core/numeral.ts';
import { CurrencyDisplayType } from '@/core/currency.ts';
import { PresetAmountColor } from '@/core/color.ts';
import { TransactionEditScopeType } from '@/core/transaction.ts';

import type { TransactionCategoryCreateBatchRequest } from './transaction_category.ts';

export class User {
    public username: string = '';
    public password: string = '';
    public confirmPassword: string = '';
    public email: string = '';
    public nickname: string = '';
    public language: string;
    public defaultCurrency: string;
    public firstDayOfWeek: number;

    private constructor(language: string, defaultCurrency: string, firstDayOfWeek: number) {
        this.language = language;
        this.defaultCurrency = defaultCurrency;
        this.firstDayOfWeek = firstDayOfWeek;
    }

    public static createNewUser(language: string, defaultCurrency: string, firstDayOfWeek: number): User {
        return new User(language, defaultCurrency, firstDayOfWeek);
    }
}

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

export const EMPTY_USER_BASIC_INFO: UserBasicInfo = {
    username: '',
    email: '',
    nickname: '',
    avatar: '',
    avatarProvider: undefined,
    defaultAccountId: '',
    transactionEditScope: TransactionEditScopeType.All.type,
    language: '',
    defaultCurrency: '',
    firstDayOfWeek: -1,
    longDateFormat: LongDateFormat.Default.type,
    shortDateFormat: ShortDateFormat.Default.type,
    longTimeFormat: LongTimeFormat.Default.type,
    shortTimeFormat: ShortTimeFormat.Default.type,
    decimalSeparator: DecimalSeparator.LanguageDefaultType,
    digitGroupingSymbol: DigitGroupingSymbol.LanguageDefaultType,
    digitGrouping: DigitGroupingType.LanguageDefaultType,
    currencyDisplayType: CurrencyDisplayType.Default.type,
    expenseAmountColor: PresetAmountColor.DefaultExpenseColor.type,
    incomeAmountColor: PresetAmountColor.DefaultIncomeColor.type,
    emailVerified: false
}
