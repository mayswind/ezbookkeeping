import type { TypeAndName } from './base.ts';

export enum TransactionType {
    ModifyBalance = 1,
    Income = 2,
    Expense = 3,
    Transfer = 4
}

export class TransactionEditScopeType implements TypeAndName {
    private static readonly allInstances: TransactionEditScopeType[] = [];

    public static readonly None = new TransactionEditScopeType(0, 'None');
    public static readonly All = new TransactionEditScopeType(1, 'All');
    public static readonly TodayOrLater = new TransactionEditScopeType(2, 'Today or later');
    public static readonly Recent24HoursOrLater = new TransactionEditScopeType(3, 'Recent 24 hours or later');
    public static readonly ThisWeekOrLater = new TransactionEditScopeType(4, 'This week or later');
    public static readonly ThisMonthOrLater = new TransactionEditScopeType(5, 'This month or later');
    public static readonly ThisYearOrLater = new TransactionEditScopeType(6, 'This year or later');

    public readonly type: number;
    public readonly name: string;

    private constructor(type: number, name: string) {
        this.type = type;
        this.name = name;

        TransactionEditScopeType.allInstances.push(this);
    }

    public static values(): TransactionEditScopeType[] {
        return TransactionEditScopeType.allInstances;
    }
}

export class TransactionTagFilterType implements TypeAndName {
    private static readonly allInstances: TransactionTagFilterType[] = [];

    public static readonly HasAny = new TransactionTagFilterType(0, 'With Any Selected Tags');
    public static readonly HasAll = new TransactionTagFilterType(1, 'With All Selected Tags');
    public static readonly NotHasAny = new TransactionTagFilterType(2, 'Without Any Selected Tags');
    public static readonly NotHasAll = new TransactionTagFilterType(3, 'Without All Selected Tags');

    public static readonly Default = TransactionTagFilterType.HasAny;

    public readonly type: number;
    public readonly name: string;

    private constructor(type: number, name: string) {
        this.type = type;
        this.name = name;

        TransactionTagFilterType.allInstances.push(this);
    }

    public static values(): TransactionTagFilterType[] {
        return TransactionTagFilterType.allInstances;
    }
}

export class ImportTransactionColumnType implements TypeAndName {
    private static readonly allInstances: ImportTransactionColumnType[] = [];

    public static readonly TransactionTime = new ImportTransactionColumnType(1, 'Transaction Time');
    public static readonly TransactionTimezone = new ImportTransactionColumnType(2, 'Transaction Timezone');
    public static readonly TransactionType = new ImportTransactionColumnType(3, 'Transaction Type');
    public static readonly Category = new ImportTransactionColumnType(4, 'Category');
    public static readonly SubCategory = new ImportTransactionColumnType(5, 'Secondary Category');
    public static readonly AccountName = new ImportTransactionColumnType(6, 'Account Name');
    public static readonly AccountCurrency = new ImportTransactionColumnType(7, 'Currency');
    public static readonly Amount = new ImportTransactionColumnType(8, 'Amount');
    public static readonly RelatedAccountName = new ImportTransactionColumnType(9, 'Transfer In Account Name');
    public static readonly RelatedAccountCurrency = new ImportTransactionColumnType(10, 'Transfer In Currency');
    public static readonly RelatedAmount = new ImportTransactionColumnType(11, 'Transfer In Amount');
    public static readonly GeographicLocation = new ImportTransactionColumnType(12, 'Geographic Location');
    public static readonly Tags = new ImportTransactionColumnType(13, 'Tags');
    public static readonly Description = new ImportTransactionColumnType(14, 'Description');

    public readonly type: number;
    public readonly name: string;

    private constructor(type: number, name: string) {
        this.type = type;
        this.name = name;

        ImportTransactionColumnType.allInstances.push(this);
    }

    public static values(): ImportTransactionColumnType[] {
        return ImportTransactionColumnType.allInstances;
    }
}
