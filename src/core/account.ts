import type { TypeAndName } from './base.ts';

type AccountTypeName = 'SingleAccount' | 'MultiSubAccounts';

export class AccountType implements TypeAndName {
    private static readonly allInstances: AccountType[] = [];
    private static readonly allInstancesByTypeName: Record<string, AccountType> = {};

    public static readonly SingleAccount = new AccountType(1, 'SingleAccount', 'Single Account');
    public static readonly MultiSubAccounts = new AccountType(2, 'MultiSubAccounts', 'Multiple Sub-accounts');

    public readonly type: number;
    public readonly typeName: AccountTypeName;
    public readonly name: string;

    private constructor(type: number, typeName: AccountTypeName, name: string) {
        this.type = type;
        this.typeName = typeName;
        this.name = name;

        AccountType.allInstances.push(this);
        AccountType.allInstancesByTypeName[typeName] = this;
    }

    public static values(): AccountType[] {
        return AccountType.allInstances;
    }

    public static all(): Record<AccountTypeName, AccountType> {
        return AccountType.allInstancesByTypeName;
    }
}

type AccountCategoryTypeName = 'Cash' | 'CheckingAccount' | 'SavingsAccount' | 'CreditCard' | 'VirtualAccount' | 'DebtAccount' | 'Receivables' | 'CertificateOfDeposit' | 'InvestmentAccount';

export class AccountCategory implements TypeAndName {
    private static readonly allInstances: AccountCategory[] = [];
    private static readonly allInstancesByType: Record<number, AccountCategory> = {};
    private static readonly allInstancesByTypeName: Record<string, AccountCategory> = {};

    public static readonly Cash = new AccountCategory(1, 'Cash', 'Cash', '1');
    public static readonly CheckingAccount = new AccountCategory(2, 'CheckingAccount', 'Checking Account', '100');
    public static readonly SavingsAccount = new AccountCategory(8, 'SavingsAccount', 'Savings Account', '100');
    public static readonly CreditCard = new AccountCategory(3, 'CreditCard', 'Credit Card', '100');
    public static readonly VirtualAccount = new AccountCategory(4, 'VirtualAccount', 'Virtual Account', '500');
    public static readonly DebtAccount = new AccountCategory(5, 'DebtAccount', 'Debt Account', '600');
    public static readonly Receivables = new AccountCategory(6, 'Receivables', 'Receivables', '700');
    public static readonly CertificateOfDeposit = new AccountCategory(9, 'CertificateOfDeposit', 'Certificate of Deposit', '110');
    public static readonly InvestmentAccount = new AccountCategory(7, 'InvestmentAccount', 'Investment Account', '800');

    public static readonly Default = AccountCategory.Cash;

    public readonly type: number;
    public readonly typeName: AccountCategoryTypeName;
    public readonly name: string;
    public readonly defaultAccountIconId: string;

    private constructor(type: number, typeName: AccountCategoryTypeName, name: string, defaultAccountIconId: string) {
        this.type = type;
        this.typeName = typeName;
        this.name = name;
        this.defaultAccountIconId = defaultAccountIconId;

        AccountCategory.allInstances.push(this);
        AccountCategory.allInstancesByType[type] = this;
        AccountCategory.allInstancesByTypeName[typeName] = this;
    }

    public static values(): AccountCategory[] {
        return AccountCategory.allInstances;
    }

    public static valueOf(type: number): AccountCategory {
        return AccountCategory.allInstancesByType[type];
    }
}
