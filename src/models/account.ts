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
    public visible: boolean;
    public subAccounts?: Account[];

    private readonly _isAsset?: boolean;
    private readonly _isLiability?: boolean;

    protected constructor(id: string, name: string, parentId: string, category: number, type: number, icon: string, color: string, currency: string, balance: number, comment: string, displayOrder: number, visible: boolean, balanceTime?: number, creditCardStatementDate?: number, isAsset?: boolean, isLiability?: boolean, subAccounts?: Account[]) {
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
        this._isAsset = isAsset;
        this._isLiability = isLiability;

        if (typeof(subAccounts) !== 'undefined') {
            this.subAccounts = subAccounts;
        } else {
            this.subAccounts = undefined;
        }
    }

    public get isAsset(): boolean {
        if (typeof(this._isAsset) !== 'undefined') {
            return this._isAsset;
        }

        const accountCategory = AccountCategory.valueOf(this.category);

        if (accountCategory) {
            return accountCategory.isAsset;
        }

        return false;
    }

    public get isLiability(): boolean {
        if (typeof(this._isLiability) !== 'undefined') {
            return this._isLiability;
        }

        const accountCategory = AccountCategory.valueOf(this.category);

        if (accountCategory) {
            return accountCategory.isLiability;
        }

        return false;
    }

    public get hidden(): boolean {
        return !this.visible;
    }

    public equals(other: Account): boolean {
        const isEqual = this.id === other.id &&
            this.name === other.name &&
            this.parentId === other.parentId &&
            this.category === other.category &&
            this.type === other.type &&
            this.icon === other.icon &&
            this.color === other.color &&
            this.currency === other.currency &&
            this.balance === other.balance &&
            this.balanceTime === other.balanceTime &&
            this.comment === other.comment &&
            this.displayOrder === other.displayOrder &&
            this.visible === other.visible &&
            this.creditCardStatementDate === other.creditCardStatementDate;

        if (!isEqual) {
            return false;
        }

        if (this.subAccounts && other.subAccounts) {
            if (this.subAccounts.length !== other.subAccounts.length) {
                return false;
            }

            for (let i = 0; i < this.subAccounts.length; i++) {
                if (!this.subAccounts[i].equals(other.subAccounts[i])) {
                    return false;
                }
            }
        } else if ((this.subAccounts && this.subAccounts.length) || (other.subAccounts && other.subAccounts.length)) {
            return false;
        }

        return true;
    }

    public fillFrom(other: Account): void {
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

    public setSuitableIcon(oldCategory: number, newCategory: number): void {
        const allCategories = AccountCategory.values();

        for (let i = 0; i < allCategories.length; i++) {
            if (allCategories[i].type === oldCategory) {
                if (this.icon !== allCategories[i].defaultAccountIconId) {
                    return;
                } else {
                    break;
                }
            }
        }

        for (let i = 0; i < allCategories.length; i++) {
            if (allCategories[i].type === newCategory) {
                this.icon = allCategories[i].defaultAccountIconId;
            }
        }
    }

    public toCreateRequest(clientSessionId: string, subAccounts?: Account[], parentAccount?: Account): AccountCreateRequest {
        let subAccountCreateRequests: AccountCreateRequest[] | undefined = undefined;

        if (this.type === AccountType.MultiSubAccounts.type) {
            subAccountCreateRequests = [];

            if (!subAccounts) {
                subAccounts = this.subAccounts;
            }

            if (subAccounts) {
                for (const subAccount of subAccounts) {
                    subAccountCreateRequests.push(subAccount.toCreateRequest(clientSessionId, undefined, this));
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
            subAccounts: !parentAccount ? subAccountCreateRequests : undefined,
            clientSessionId: !parentAccount ? clientSessionId : undefined
        };
    }

    public toModifyRequest(clientSessionId: string, subAccounts?: Account[], parentAccount?: Account): AccountModifyRequest {
        let subAccountModifyRequests: AccountModifyRequest[] | undefined = undefined;

        if (this.type === AccountType.MultiSubAccounts.type) {
            subAccountModifyRequests = [];

            if (!subAccounts) {
                subAccounts = this.subAccounts;
            }

            if (subAccounts) {
                for (const subAccount of subAccounts) {
                    subAccountModifyRequests.push(subAccount.toModifyRequest(clientSessionId, undefined, this));
                }
            }
        }

        return {
            id: this.id || '0',
            name: this.name,
            category: parentAccount ? parentAccount.category : this.category,
            icon: this.icon,
            color: this.color,
            currency: parentAccount && (!this.id || this.id === '0') ? this.currency : undefined,
            balance: parentAccount && (!this.id || this.id === '0') ? this.balance : undefined,
            balanceTime: parentAccount && (!this.id || this.id === '0') ? this.balanceTime : undefined,
            comment: this.comment,
            creditCardStatementDate: !parentAccount && this.category === AccountCategory.CreditCard.type ? this.creditCardStatementDate : undefined,
            hidden: !this.visible,
            subAccounts: !parentAccount ? subAccountModifyRequests : undefined,
            clientSessionId: !parentAccount ? clientSessionId : undefined
        };
    }

    public getAccountOrSubAccountId(subAccountId: string): string | null {
        if (this.type === AccountType.SingleAccount.type) {
            return this.id;
        } else if (this.type === AccountType.MultiSubAccounts.type && !subAccountId) {
            return this.id;
        } else if (this.type === AccountType.MultiSubAccounts.type && subAccountId) {
            if (!this.subAccounts || !this.subAccounts.length) {
                return null;
            }

            for (let i = 0; i < this.subAccounts.length; i++) {
                const subAccount = this.subAccounts[i];

                if (subAccountId && subAccountId === subAccount.id) {
                    return subAccount.id;
                }
            }

            return null;
        } else {
            return null;
        }
    }

    public isAccountOrSubAccountHidden(subAccountId: string): boolean {
        if (this.type === AccountType.SingleAccount.type) {
            return this.hidden;
        } else if (this.type === AccountType.MultiSubAccounts.type && !subAccountId) {
            return this.hidden;
        } else if (this.type === AccountType.MultiSubAccounts.type && subAccountId) {
            if (!this.subAccounts || !this.subAccounts.length) {
                return false;
            }

            for (let i = 0; i < this.subAccounts.length; i++) {
                const subAccount = this.subAccounts[i];

                if (subAccountId && subAccountId === subAccount.id) {
                    return subAccount.hidden;
                }
            }

            return false;
        } else {
            return false;
        }
    }

    public getAccountOrSubAccountComment(subAccountId: string): string | null {
        if (this.type === AccountType.SingleAccount.type) {
            return this.comment;
        } else if (this.type === AccountType.MultiSubAccounts.type && !subAccountId) {
            return this.comment;
        } else if (this.type === AccountType.MultiSubAccounts.type && subAccountId) {
            if (!this.subAccounts || !this.subAccounts.length) {
                return null;
            }

            for (let i = 0; i < this.subAccounts.length; i++) {
                const subAccount = this.subAccounts[i];

                if (subAccountId && subAccountId === subAccount.id) {
                    return subAccount.comment;
                }
            }

            return null;
        } else {
            return null;
        }
    }

    public getAccountOrSubAccount(subAccountId: string): Account | null {
        if (this.type === AccountType.SingleAccount.type) {
            return this;
        } else if (this.type === AccountType.MultiSubAccounts.type && !subAccountId) {
            return this;
        } else if (this.type === AccountType.MultiSubAccounts.type && subAccountId) {
            if (!this.subAccounts || !this.subAccounts.length) {
                return null;
            }

            for (let i = 0; i < this.subAccounts.length; i++) {
                const subAccount = this.subAccounts[i];

                if (subAccountId && subAccountId === subAccount.id) {
                    return subAccount;
                }
            }

            return null;
        } else {
            return null;
        }
    }

    public getSubAccount(subAccountId: string): Account | null {
        if (!this.subAccounts || !this.subAccounts.length) {
            return null;
        }

        for (let i = 0; i < this.subAccounts.length; i++) {
            const subAccount = this.subAccounts[i];

            if (subAccountId && subAccountId === subAccount.id) {
                return subAccount;
            }
        }

        return null;
    }

    public getSubAccountCurrencies(showHidden: boolean, subAccountId: string): string[] {
        if (!this.subAccounts || !this.subAccounts.length) {
            return [];
        }

        const subAccountCurrenciesMap: Record<string, boolean> = {};
        const subAccountCurrencies: string[] = [];

        for (let i = 0; i < this.subAccounts.length; i++) {
            const subAccount = this.subAccounts[i];

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

    public clone(): Account {
        return new Account(
            this.id,
            this.name,
            this.parentId,
            this.category,
            this.type,
            this.icon,
            this.color,
            this.currency,
            this.balance,
            this.comment,
            this.displayOrder,
            this.visible,
            this.balanceTime,
            this.creditCardStatementDate,
            this.isAsset,
            this.isLiability,
            typeof(this.subAccounts) !== 'undefined' ? Account.cloneAccounts(this.subAccounts) : undefined);
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
            accountResponse.subAccounts ? Account.ofMulti(accountResponse.subAccounts) : undefined
        );
    }

    public static ofMulti(accountResponses: AccountInfoResponse[]): Account[] {
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

    public static cloneAccounts(accounts: Account[]): Account[] {
        const clonedAccounts: Account[] = [];

        for (const account of accounts) {
            clonedAccounts.push(account.clone());
        }

        return clonedAccounts;
    }

    public static sortAccounts(accounts: Account[], allAccountsMap?: Record<string, Account>): Account[] {
        if (!accounts || !accounts.length) {
            return accounts;
        }

        return accounts.sort(function (account1, account2) {
            if (account1.category !== account2.category) {
                const account1Category = AccountCategory.valueOf(account1.category);
                const account2Category = AccountCategory.valueOf(account2.category);

                if (!account1Category) {
                    return 1;
                }

                if (!account2Category) {
                    return -1;
                }

                return account1Category.displayOrder - account2Category.displayOrder;
            }

            if (account1.parentId === account2.parentId) {
                return account1.displayOrder - account2.displayOrder;
            }

            if (account1.id === account2.parentId) {
                return -1;
            } else if (account2.id === account1.parentId) {
                return 1;
            }

            let account1DisplayOrder: number | null = account1.displayOrder;
            let account2DisplayOrder: number | null = account2.displayOrder;

            if (account1.parentId && account1.parentId !== '0') {
                if (allAccountsMap && allAccountsMap[account1.parentId]) {
                    account1DisplayOrder = allAccountsMap[account1.parentId].displayOrder;
                } else {
                    account1DisplayOrder = null;
                }
            }

            if (account2.parentId && account2.parentId !== '0') {
                if (allAccountsMap && allAccountsMap[account2.parentId]) {
                    account2DisplayOrder = allAccountsMap[account2.parentId].displayOrder;
                } else {
                    account2DisplayOrder = null;
                }
            }

            if (account1DisplayOrder !== null && account2DisplayOrder !== null) {
                return account1DisplayOrder - account2DisplayOrder;
            } else {
                return account1.id.localeCompare(account2.id);
            }
        });
    }
}

export class AccountWithDisplayBalance extends Account {
    public displayBalance: string;

    private constructor(account: Account, displayBalance: string) {
        super(
            account.id,
            account.name,
            account.parentId,
            account.category,
            account.type,
            account.icon,
            account.color,
            account.currency,
            account.balance,
            account.comment,
            account.displayOrder,
            account.visible,
            account.balanceTime,
            account.creditCardStatementDate,
            account.isAsset,
            account.isLiability,
            account.subAccounts
        );

        this.displayBalance = displayBalance;
    }

    public static fromAccount(account: Account, displayBalance: string): AccountWithDisplayBalance {
        return new AccountWithDisplayBalance(account, displayBalance);
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
    readonly currency?: string;
    readonly balance?: number;
    readonly balanceTime?: number;
    readonly comment: string;
    readonly creditCardStatementDate?: number;
    readonly hidden: boolean;
    readonly subAccounts?: AccountModifyRequest[];
    readonly clientSessionId?: string;
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

export class CategorizedAccountWithDisplayBalance {
    public category: number;
    public name: string;
    public icon: string;
    public accounts: AccountWithDisplayBalance[];
    public displayBalance: string;

    private constructor(category: number, name: string, icon: string, accounts: AccountWithDisplayBalance[], displayBalance: string) {
        this.category = category;
        this.name = name;
        this.icon = icon;
        this.accounts = accounts;
        this.displayBalance = displayBalance;
    }

    public static of(categorizedAccount: CategorizedAccount, accounts: AccountWithDisplayBalance[], displayBalance: string): CategorizedAccountWithDisplayBalance {
        return new CategorizedAccountWithDisplayBalance(categorizedAccount.category, categorizedAccount.name, categorizedAccount.icon, accounts, displayBalance);
    }
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
