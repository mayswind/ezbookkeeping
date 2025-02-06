import type { TypeAndName, TypeAndDisplayName } from './base.ts';

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

    public static readonly Cash = new AccountCategory(1, 'Cash', '1');
    public static readonly CheckingAccount = new AccountCategory(2, 'Checking Account', '100');
    public static readonly SavingsAccount = new AccountCategory(8, 'Savings Account', '100');
    public static readonly CreditCard = new AccountCategory(3, 'Credit Card', '100');
    public static readonly VirtualAccount = new AccountCategory(4, 'Virtual Account', '500');
    public static readonly DebtAccount = new AccountCategory(5, 'Debt Account', '600');
    public static readonly Receivables = new AccountCategory(6, 'Receivables', '700');
    public static readonly CertificateOfDeposit = new AccountCategory(9, 'Certificate of Deposit', '110');
    public static readonly InvestmentAccount = new AccountCategory(7, 'Investment Account', '800');

    public static readonly Default = AccountCategory.Cash;

    public readonly type: number;
    public readonly name: string;
    public readonly defaultAccountIconId: string;

    private constructor(type: number, name: string, defaultAccountIconId: string) {
        this.type = type;
        this.name = name;
        this.defaultAccountIconId = defaultAccountIconId;

        AccountCategory.allInstances.push(this);
        AccountCategory.allInstancesByType[type] = this;
    }

    public static values(): AccountCategory[] {
        return AccountCategory.allInstances;
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
