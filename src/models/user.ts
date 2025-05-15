import { LongDateFormat, ShortDateFormat, LongTimeFormat, ShortTimeFormat } from '@/core/datetime.ts';
import { DecimalSeparator, DigitGroupingSymbol, DigitGroupingType } from '@/core/numeral.ts';
import { CurrencyDisplayType } from '@/core/currency.ts';
import { PresetAmountColor } from '@/core/color.ts';
import type { LocalizedPresetCategory } from '@/core/category.ts';
import { TransactionEditScopeType } from '@/core/transaction.ts';
import { FiscalYearFormat, FiscalYearStart } from '@/core/fiscalyear';

export class User {
    public username: string = '';
    public password: string = '';
    public confirmPassword: string = '';
    public email: string = '';
    public nickname: string = '';
    public language: string;
    public defaultCurrency: string;
    public firstDayOfWeek: number;

    public defaultAccountId: string = '';
    public transactionEditScope: number = 1;
    public fiscalYearStart: number = 0;
    public fiscalYearFormat: number = 0;
    public longDateFormat: number = 0;
    public shortDateFormat: number = 0;
    public longTimeFormat: number = 0;
    public shortTimeFormat: number = 0;
    public decimalSeparator: number = 0;
    public digitGroupingSymbol: number = 0;
    public digitGrouping: number = 0;
    public currencyDisplayType: number = 0;
    public expenseAmountColor: number = 0;
    public incomeAmountColor: number = 0;

    private constructor(language: string, defaultCurrency: string, firstDayOfWeek: number) {
        this.language = language;
        this.defaultCurrency = defaultCurrency;
        this.firstDayOfWeek = firstDayOfWeek;
    }

    public fillFrom(user: User | UserBasicInfo | UserProfileResponse): void {
        this.username = user.username;
        this.email = user.email;
        this.nickname = user.nickname;
        this.language = user.language;
        this.defaultCurrency = user.defaultCurrency;
        this.firstDayOfWeek = user.firstDayOfWeek;
        this.defaultAccountId = user.defaultAccountId;
        this.transactionEditScope = user.transactionEditScope;
        this.fiscalYearStart = user.fiscalYearStart;
        this.longDateFormat = user.longDateFormat;
        this.shortDateFormat = user.shortDateFormat;
        this.longTimeFormat = user.longTimeFormat;
        this.shortTimeFormat = user.shortTimeFormat;
        this.fiscalYearFormat = user.fiscalYearFormat;
        this.decimalSeparator = user.decimalSeparator;
        this.digitGroupingSymbol = user.digitGroupingSymbol;
        this.digitGrouping = user.digitGrouping;
        this.currencyDisplayType = user.currencyDisplayType;
        this.expenseAmountColor = user.expenseAmountColor;
        this.incomeAmountColor = user.incomeAmountColor;
    }

    public toRegisterRequest(categories?: LocalizedPresetCategory[]): UserRegisterRequest {
        return {
            username: this.username,
            email: this.email,
            nickname: this.nickname,
            password: this.password,
            language: this.language,
            defaultCurrency: this.defaultCurrency,
            firstDayOfWeek: this.firstDayOfWeek,
            categories: categories
        };
    }

    public toProfileUpdateRequest(currentPassword?: string): UserProfileUpdateRequest {
        return {
            email: this.email,
            nickname: this.nickname,
            password: this.password,
            oldPassword: currentPassword,
            defaultAccountId: this.defaultAccountId,
            transactionEditScope: this.transactionEditScope,
            language: this.language,
            defaultCurrency: this.defaultCurrency,
            firstDayOfWeek: this.firstDayOfWeek,
            fiscalYearStart: this.fiscalYearStart,
            longDateFormat: this.longDateFormat,
            shortDateFormat: this.shortDateFormat,
            longTimeFormat: this.longTimeFormat,
            shortTimeFormat: this.shortTimeFormat,
            fiscalYearFormat: this.fiscalYearFormat,
            decimalSeparator: this.decimalSeparator,
            digitGroupingSymbol: this.digitGroupingSymbol,
            digitGrouping: this.digitGrouping,
            currencyDisplayType: this.currencyDisplayType,
            expenseAmountColor: this.expenseAmountColor,
            incomeAmountColor: this.incomeAmountColor
        };
    }

    public static of(userInfo: UserBasicInfo): User {
        const user = new User(userInfo.language, userInfo.defaultCurrency, userInfo.firstDayOfWeek);
        user.defaultAccountId = userInfo.defaultAccountId;
        user.transactionEditScope = userInfo.transactionEditScope;
        user.fiscalYearStart = userInfo.fiscalYearStart;
        user.longDateFormat = userInfo.longDateFormat;
        user.shortDateFormat = userInfo.shortDateFormat;
        user.longTimeFormat = userInfo.longTimeFormat;
        user.shortTimeFormat = userInfo.shortTimeFormat;
        user.fiscalYearFormat = userInfo.fiscalYearFormat;
        user.decimalSeparator = userInfo.decimalSeparator;
        user.digitGroupingSymbol = userInfo.digitGroupingSymbol;
        user.digitGrouping = userInfo.digitGrouping;
        user.currencyDisplayType = userInfo.currencyDisplayType;
        user.expenseAmountColor = userInfo.expenseAmountColor;
        user.incomeAmountColor = userInfo.incomeAmountColor;

        return user;
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
    readonly fiscalYearStart: number;
    readonly firstDayOfWeek: number;
    readonly longDateFormat: number;
    readonly shortDateFormat: number;
    readonly longTimeFormat: number;
    readonly shortTimeFormat: number;
    readonly fiscalYearFormat: number;
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
    readonly categories?: LocalizedPresetCategory[];
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
    readonly fiscalYearStart?: number;
    readonly longDateFormat?: number;
    readonly shortDateFormat?: number;
    readonly longTimeFormat?: number;
    readonly shortTimeFormat?: number;
    readonly fiscalYearFormat?: number;
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
    fiscalYearStart: FiscalYearStart.Default.value,
    longDateFormat: LongDateFormat.Default.type,
    shortDateFormat: ShortDateFormat.Default.type,
    longTimeFormat: LongTimeFormat.Default.type,
    shortTimeFormat: ShortTimeFormat.Default.type,
    fiscalYearFormat: FiscalYearFormat.Default.type,
    decimalSeparator: DecimalSeparator.LanguageDefaultType,
    digitGroupingSymbol: DigitGroupingSymbol.LanguageDefaultType,
    digitGrouping: DigitGroupingType.LanguageDefaultType,
    currencyDisplayType: CurrencyDisplayType.Default.type,
    expenseAmountColor: PresetAmountColor.DefaultExpenseColor.type,
    incomeAmountColor: PresetAmountColor.DefaultIncomeColor.type,
    emailVerified: false
}
