import type { ColorValue } from '@/core/color.ts';
import { AccountType, AccountCategory } from '@/core/account.ts';
import { PARENT_ACCOUNT_CURRENCY_PLACEHOLDER } from '@/consts/currency.ts';
import { DEFAULT_ACCOUNT_ICON_ID } from '@/consts/icon.ts';
import { DEFAULT_ACCOUNT_COLOR } from '@/consts/color.ts';

export class Account implements AccountInfoResponse {
    public id: string;
    public name: string;
    public parentId: string;
    public category: number;
    public type: number;
    public icon: string;
    public color: ColorValue;
    public currency: string;
    public balance: number;
    public balanceTime?: number;
    public comment: string;
    public creditCardStatementDate?: number;
    public displayOrder: number;
    public isAsset?: boolean;
    public isLiability?: boolean;
    public visible: boolean;
    public childrenAccounts?: Account[];

    private constructor(id: string, name: string, parentId: string, category: number, type: number, icon: string, color: string, currency: string, balance: number, comment: string, displayOrder: number, visible: boolean, balanceTime?: number, creditCardStatementDate?: number, isAsset?: boolean, isLiability?: boolean, childrenAccounts?: Account[]) {
        this.id = id;
        this.name = name;
        this.parentId = parentId;
        this.category = category;
        this.type = type;
        this.icon = icon;
        this.color = color;
        this.currency = currency;
        this.balance = balance;
        this.balanceTime = balanceTime;
        this.comment = comment;
        this.displayOrder = displayOrder;
        this.visible = visible;
        this.creditCardStatementDate = creditCardStatementDate;
        this.isAsset = isAsset;
        this.isLiability = isLiability;

        if (typeof(childrenAccounts) !== 'undefined') {
            this.childrenAccounts = childrenAccounts;
        } else {
            this.childrenAccounts = undefined;
        }
    }

    public get hidden(): boolean {
        return !this.visible;
    }

    public get subAccounts(): AccountInfoResponse[] | undefined {
        if (typeof(this.childrenAccounts) === 'undefined') {
            return undefined;
        }

        const ret: AccountInfoResponse[] = [];

        if (this.childrenAccounts) {
            for (const subCategory of this.childrenAccounts) {
                ret.push(subCategory);
            }
        }

        return ret;
    }

    public from(other: Account): void {
        this.id = other.id;
        this.category = other.category;
        this.type = other.type;
        this.name = other.name;
        this.icon = other.icon;
        this.color = other.color;
        this.currency = other.currency;
        this.balance = other.balance;
        this.balanceTime = other.balanceTime;
        this.comment = other.comment;
        this.creditCardStatementDate = other.creditCardStatementDate;
        this.visible = other.visible;
    }

    public toCreateRequest(clientSessionId: string, childrenAccounts?: Account[], parentAccount?: Account): AccountCreateRequest {
        let subAccounts: AccountCreateRequest[] | undefined = undefined;

        if (this.type === AccountType.MultiSubAccounts.type) {
            subAccounts = [];

            if (!childrenAccounts) {
                childrenAccounts = this.childrenAccounts;
            }

            if (childrenAccounts) {
                for (const subAccount of childrenAccounts) {
                    subAccounts.push(subAccount.toCreateRequest(clientSessionId, undefined, this));
                }
            }
        }

        return {
            name: this.name,
            category: parentAccount ? parentAccount.category : this.category,
            type: parentAccount ? AccountType.SingleAccount.type : this.type,
            icon: this.icon,
            color: this.color,
            currency: parentAccount || this.type === AccountType.SingleAccount.type ? this.currency : PARENT_ACCOUNT_CURRENCY_PLACEHOLDER,
            balance: parentAccount || this.type === AccountType.SingleAccount.type ? this.balance : 0,
            balanceTime: (parentAccount || this.type === AccountType.SingleAccount.type) && this.balanceTime ? this.balanceTime : 0,
            comment: this.comment,
            creditCardStatementDate: !parentAccount && this.category === AccountCategory.CreditCard.type ? this.creditCardStatementDate : undefined,
            subAccounts: !parentAccount ? subAccounts : undefined,
            clientSessionId: !parentAccount ? clientSessionId : undefined
        };
    }

    public toModifyRequest(childrenAccounts?: Account[], parentAccount?: Account): AccountModifyRequest {
        let subAccounts: AccountModifyRequest[] | undefined = undefined;

        if (this.type === AccountType.MultiSubAccounts.type) {
            subAccounts = [];

            if (!childrenAccounts) {
                childrenAccounts = this.childrenAccounts;
            }

            if (childrenAccounts) {
                for (const subAccount of childrenAccounts) {
                    subAccounts.push(subAccount.toModifyRequest(undefined, this));
                }
            }
        }

        return {
            id: this.id,
            name: this.name,
            category: parentAccount ? parentAccount.category : this.category,
            icon: this.icon,
            color: this.color,
            comment: this.comment,
            creditCardStatementDate: !parentAccount && this.category === AccountCategory.CreditCard.type ? this.creditCardStatementDate : undefined,
            hidden: !this.visible,
            subAccounts: !parentAccount ? subAccounts : undefined,
        };
    }

    public getAccountOrSubAccountId(subAccountId: string): string | null {
        if (this.type === AccountType.SingleAccount.type) {
            return this.id;
        } else if (this.type === AccountType.MultiSubAccounts.type && !subAccountId) {
            return this.id;
        } else if (this.type === AccountType.MultiSubAccounts.type && subAccountId) {
            if (!this.childrenAccounts || !this.childrenAccounts.length) {
                return null;
            }

            for (let i = 0; i < this.childrenAccounts.length; i++) {
                const subAccount = this.childrenAccounts[i];

                if (subAccountId && subAccountId === subAccount.id) {
                    return subAccount.id;
                }
            }

            return null;
        } else {
            return null;
        }
    }

    public getAccountOrSubAccountComment(subAccountId: string): string | null {
        if (this.type === AccountType.SingleAccount.type) {
            return this.comment;
        } else if (this.type === AccountType.MultiSubAccounts.type && !subAccountId) {
            return this.comment;
        } else if (this.type === AccountType.MultiSubAccounts.type && subAccountId) {
            if (!this.childrenAccounts || !this.childrenAccounts.length) {
                return null;
            }

            for (let i = 0; i < this.childrenAccounts.length; i++) {
                const subAccount = this.childrenAccounts[i];

                if (subAccountId && subAccountId === subAccount.id) {
                    return subAccount.comment;
                }
            }

            return null;
        } else {
            return null;
        }
    }

    public getSubAccountCurrencies(showHidden: boolean, subAccountId: string): string[] {
        if (!this.childrenAccounts || !this.childrenAccounts.length) {
            return [];
        }

        const subAccountCurrenciesMap: Record<string, boolean> = {};
        const subAccountCurrencies: string[] = [];

        for (let i = 0; i < this.childrenAccounts.length; i++) {
            const subAccount = this.childrenAccounts[i];

            if (!showHidden && subAccount.hidden) {
                continue;
            }

            if (subAccountId && subAccountId === subAccount.id) {
                return [subAccount.currency];
            } else {
                if (!subAccountCurrenciesMap[subAccount.currency]) {
                    subAccountCurrenciesMap[subAccount.currency] = true;
                    subAccountCurrencies.push(subAccount.currency);
                }
            }
        }

        return subAccountCurrencies;
    }

    public createNewSubAccount(currency: string, balanceTime: number): Account {
        return new Account(
            '', // id
            '', // name
            '', // parentId
            0, // category
            0, // type
            this.icon, // icon
            this.color, // color
            currency, // currency
            0, // balance
            '', // comment
            0, // displayOrder
            true, // visible
            balanceTime, // balanceTime
            0 // creditCardStatementDate
        );
    }

    public static createNewAccount(currency: string, balanceTime: number): Account {
        return new Account(
            '', // id
            '', // name
            '', // parentId
            AccountCategory.Cash.type, // category
            AccountType.SingleAccount.type, // type
            DEFAULT_ACCOUNT_ICON_ID, // icon
            DEFAULT_ACCOUNT_COLOR, // color
            currency, // currency
            0, // balance
            '', // comment
            0, // displayOrder
            true, // visible
            balanceTime, // balanceTime
            0 // creditCardStatementDate
        );
    }

    public static of(accountResponse: AccountInfoResponse): Account {
        return new Account(
            accountResponse.id,
            accountResponse.name,
            accountResponse.parentId,
            accountResponse.category,
            accountResponse.type,
            accountResponse.icon,
            accountResponse.color,
            accountResponse.currency,
            accountResponse.balance,
            accountResponse.comment,
            accountResponse.displayOrder,
            !accountResponse.hidden,
            undefined,
            accountResponse.creditCardStatementDate,
            accountResponse.isAsset,
            accountResponse.isLiability,
            accountResponse.subAccounts ? Account.ofMany(accountResponse.subAccounts) : undefined
        );
    }

    public static ofMany(accountResponses: AccountInfoResponse[]): Account[] {
        const accounts: Account[] = [];

        for (const accountResponse of accountResponses) {
            accounts.push(Account.of(accountResponse));
        }

        return accounts;
    }

    public static findAccountNameById(accounts: Account[], accountId: string, defaultName?: string): string | undefined {
        for (const account of accounts) {
            if (account.id === accountId) {
                return account.name;
            }
        }

        return defaultName;
    }
}

export interface AccountCreateRequest {
    readonly name: string;
    readonly category: number;
    readonly type: number;
    readonly icon: string;
    readonly color: string;
    readonly currency: string;
    readonly balance: number;
    readonly balanceTime: number;
    readonly comment: string;
    readonly creditCardStatementDate?: number;
    readonly subAccounts?: AccountCreateRequest[];
    readonly clientSessionId?: string;
}

export interface AccountModifyRequest {
    readonly id: string;
    readonly name: string;
    readonly category: number;
    readonly icon: string;
    readonly color: string;
    readonly comment: string;
    readonly creditCardStatementDate?: number;
    readonly hidden: boolean;
    readonly subAccounts?: AccountModifyRequest[];
}

export interface AccountInfoResponse {
    readonly id: string;
    readonly name: string;
    readonly parentId: string;
    readonly category: number;
    readonly type: number;
    readonly icon: string;
    readonly color: string;
    readonly currency: string;
    readonly balance: number;
    readonly comment: string;
    readonly creditCardStatementDate?: number;
    readonly displayOrder: number;
    readonly isAsset?: boolean;
    readonly isLiability?: boolean;
    readonly hidden: boolean;
    readonly subAccounts?: AccountInfoResponse[];
}

export interface AccountHideRequest {
    readonly id: string;
    readonly hidden: boolean;
}

export interface AccountMoveRequest {
    readonly newDisplayOrders: AccountNewDisplayOrderRequest[];
}

export interface AccountNewDisplayOrderRequest {
    readonly id: string;
    readonly displayOrder: number;
}

export interface AccountDeleteRequest {
    readonly id: string;
}

export interface AccountBalance {
    readonly balance: number;
    readonly isAsset: boolean;
    readonly isLiability: boolean;
    readonly currency: string;
}

export interface AccountDisplayBalance {
    readonly balance: string;
    readonly currency: string;
}

export interface CategorizedAccount {
    readonly category: number;
    readonly name: string;
    readonly icon: string;
    readonly accounts: Account[];
}

export interface AccountCategoriesWithVisibleCount {
    readonly category: number;
    readonly name: string;
    readonly icon: string;
    readonly allAccounts: Account[];
    readonly allVisibleAccountCount: number;
    readonly firstVisibleAccountIndex: number;
    readonly allSubAccounts: Record<string, Account[]>;
    readonly allVisibleSubAccountCounts: Record<string, number>;
    readonly allFirstVisibleSubAccountIndexes: Record<string, number>;
}

export interface AccountShowingIds {
    readonly accounts: Record<number, string>;
    readonly subAccounts: Record<string, string>;

}
