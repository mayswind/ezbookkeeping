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
