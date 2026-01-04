import { type TypeAndName, type TypeAndDisplayName, itemAndIndex } from './base.ts';

export class AccountType implements TypeAndName {
    private static readonly allInstances: AccountType[] = [];

    public static readonly SingleAccount = new AccountType(1, 'Single Account');
    public static readonly MultiSubAccounts = new AccountType(2, 'Multiple Sub-accounts');

    public readonly type: number;
    public readonly name: string;

    private constructor(type: number, name: string) {
        this.type = type;
        this.name = name;

        AccountType.allInstances.push(this);
    }

    public static values(): AccountType[] {
        return AccountType.allInstances;
    }
}

export class AccountCategory implements TypeAndName {
    private static readonly allInstances: AccountCategory[] = [];
    private static readonly allInstancesByType: Record<number, AccountCategory> = {};

    public static readonly Cash = new AccountCategory(1, 1, 'Cash', true, false, '1');
    public static readonly CheckingAccount = new AccountCategory(2, 2, 'Checking Account', true, false, '100');
    public static readonly SavingsAccount = new AccountCategory(8, 3, 'Savings Account', true, false, '100');
    public static readonly CreditCard = new AccountCategory(3, 4, 'Credit Card', false, true, '100');
    public static readonly VirtualAccount = new AccountCategory(4, 5, 'Virtual Account', true, false, '500');
    public static readonly DebtAccount = new AccountCategory(5, 6, 'Debt Account', false, true, '600');
    public static readonly Receivables = new AccountCategory(6, 7, 'Receivables', true, false, '700');
    public static readonly CertificateOfDeposit = new AccountCategory(9, 8, 'Certificate of Deposit', true, false, '110');
    public static readonly InvestmentAccount = new AccountCategory(7, 9, 'Investment Account', true, false, '800');

    public static readonly Default = AccountCategory.Cash;

    public readonly type: number;
    public readonly defaultDisplayOrder: number;
    public readonly name: string;
    public readonly isAsset: boolean;
    public readonly isLiability: boolean
    public readonly defaultAccountIconId: string;

    private constructor(type: number, defaultDisplayOrder: number, name: string, isAsset: boolean, isLiability: boolean, defaultAccountIconId: string) {
        this.type = type;
        this.defaultDisplayOrder = defaultDisplayOrder;
        this.name = name;
        this.isAsset = isAsset;
        this.isLiability = isLiability;
        this.defaultAccountIconId = defaultAccountIconId;

        AccountCategory.allInstances.push(this);
        AccountCategory.allInstancesByType[type] = this;
    }

    public static values(customAccountCategoryOrder?: string): AccountCategory[] {
        if (!customAccountCategoryOrder) {
            return [...AccountCategory.allInstances];
        }

        const typeOrders: string[] = customAccountCategoryOrder.split(',');
        const orderedCategories: AccountCategory[] = [];
        const addedTypes: Record<string, boolean> = {};

        for (const type of typeOrders) {
            const category = AccountCategory.valueOf(parseInt(type.trim()));

            if (category) {
                orderedCategories.push(category);
                addedTypes[type] = true;
            }
        }

        for (const category of AccountCategory.allInstances) {
            if (!addedTypes[category.type]) {
                orderedCategories.push(category);
            }
        }

        return orderedCategories;
    }

    public static allDisplayOrders(customAccountCategoryOrder: string): Record<number, number> {
        const displayOrders: Record<number, number> = {};

        for (const [category, index] of itemAndIndex(AccountCategory.values(customAccountCategoryOrder))) {
            displayOrders[category.type] = index + 1;
        }

        return displayOrders;
    }

    public static valueOf(type: number): AccountCategory | undefined {
        return AccountCategory.allInstancesByType[type];
    }
}

export interface LocalizedAccountCategory extends TypeAndDisplayName {
    readonly type: number;
    readonly displayName: string;
    readonly defaultAccountIconId: string;
}
