import type { TypeAndName } from './base.ts';

export enum TransactionType {
    ModifyBalance = 1,
    Income = 2,
    Expense = 3,
    Transfer = 4
}

export enum TransactionRelatedAccountType {
    TransferFrom = 1,
    TransferTo = 2
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
    private static readonly allInstancesByType: Record<number, TransactionTagFilterType> = {};

    public static readonly HasAny = new TransactionTagFilterType(0, 'Include Any Selected Tags');
    public static readonly HasAll = new TransactionTagFilterType(1, 'Include All Selected Tags');
    public static readonly NotHasAny = new TransactionTagFilterType(2, 'Exclude Any Selected Tags');
    public static readonly NotHasAll = new TransactionTagFilterType(3, 'Exclude All Selected Tags');

    public readonly type: number;
    public readonly name: string;

    private constructor(type: number, name: string) {
        this.type = type;
        this.name = name;

        TransactionTagFilterType.allInstances.push(this);
        TransactionTagFilterType.allInstancesByType[type] = this;
    }

    public static values(): TransactionTagFilterType[] {
        return TransactionTagFilterType.allInstances;
    }

    public static parse(type: number): TransactionTagFilterType | undefined {
        return TransactionTagFilterType.allInstancesByType[type];
    }
}

export class TransactionQuickAddButtonActionType implements TypeAndName {
    private static readonly allInstances: TransactionQuickAddButtonActionType[] = [];
    private static readonly allInstancesByType: Record<number, TransactionQuickAddButtonActionType> = {};

    public static readonly SaveAndGoBack = new TransactionQuickAddButtonActionType(0, 'Save');
    public static readonly OpenMenu = new TransactionQuickAddButtonActionType(1, 'Open Menu');
    public static readonly SaveAndAddNewTransaction = new TransactionQuickAddButtonActionType(2, 'Save & New');
    public static readonly SaveAndKeepCurrentData = new TransactionQuickAddButtonActionType(3, 'Save & Duplicate');

    public static readonly Default = TransactionQuickAddButtonActionType.SaveAndGoBack;

    public readonly type: number;
    public readonly name: string;

    private constructor(type: number, name: string) {
        this.type = type;
        this.name = name;

        TransactionQuickAddButtonActionType.allInstances.push(this);
        TransactionQuickAddButtonActionType.allInstancesByType[type] = this;
    }

    public static values(): TransactionQuickAddButtonActionType[] {
        return TransactionQuickAddButtonActionType.allInstances;
    }

    public static valueOf(type: number): TransactionQuickAddButtonActionType | undefined {
        return TransactionQuickAddButtonActionType.allInstancesByType[type];
    }
}
