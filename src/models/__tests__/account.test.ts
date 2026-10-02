import { describe, expect, it } from 'vitest';

import { AccountCategory } from '@/core/account.ts';

import { Account } from '@/models/account.ts';

function createAccount(id: string, category: number, displayOrder: number, parentId: string = '0'): Account {
    return Account.of({
        id: id,
        name: id,
        parentId: parentId,
        category: category,
        type: 1,
        icon: '1',
        iconType: 0,
        color: '000000',
        currency: 'USD',
        balance: '0',
        comment: '',
        displayOrder: displayOrder,
        hidden: false
    });
}

describe('Account.compareTo', () => {
    it('should return zero when comparing an account with itself', () => {
        const account = createAccount('1', AccountCategory.Cash.type, 10);

        expect(account.compareTo(account, { [AccountCategory.Cash.type]: 1 })).toBe(0);
    });

    it('should compare category display orders before account display orders', () => {
        const first = createAccount('1', AccountCategory.Cash.type, 1);
        const second = createAccount('2', AccountCategory.CheckingAccount.type, 10);
        const categoryDisplayOrders = { [AccountCategory.Cash.type]: 2, [AccountCategory.CheckingAccount.type]: 1 };

        expect(first.compareTo(second, categoryDisplayOrders)).toBe(1);
        expect(second.compareTo(first, categoryDisplayOrders)).toBe(-1);
    });

    it('should place an account with an unspecified category order after an account with a specified category order', () => {
        const first = createAccount('1', AccountCategory.Cash.type, 1);
        const second = createAccount('2', AccountCategory.CheckingAccount.type, 10);
        const categoryDisplayOrders = { [AccountCategory.CheckingAccount.type]: 1 };

        expect(first.compareTo(second, categoryDisplayOrders)).toBe(1);
        expect(second.compareTo(first, categoryDisplayOrders)).toBe(-1);
    });

    it.each([
        [1, 2, -1],
        [2, 1, 1],
        [2, 2, 0],
    ])('should compare top-level account display orders %j and %j', (firstOrder, secondOrder, expected) => {
        const first = createAccount('1', AccountCategory.Cash.type, firstOrder);
        const second = createAccount('2', AccountCategory.Cash.type, secondOrder);

        expect(first.compareTo(second, { [AccountCategory.Cash.type]: 1 })).toBe(expected);
    });

    it.each([
        [1, 2, -1],
        [2, 1, 1],
        [2, 2, 0],
    ])('should compare sibling account display orders %j and %j', (firstOrder, secondOrder, expected) => {
        const first = createAccount('11', AccountCategory.Cash.type, firstOrder, '1');
        const second = createAccount('12', AccountCategory.Cash.type, secondOrder, '1');

        expect(first.compareTo(second, { [AccountCategory.Cash.type]: 1 })).toBe(expected);
    });

    it('should place a parent account before its child regardless of their display orders', () => {
        const parent = createAccount('1', AccountCategory.Cash.type, 10);
        const child = createAccount('11', AccountCategory.Cash.type, 1, '1');

        expect(parent.compareTo(child, { [AccountCategory.Cash.type]: 1 })).toBe(-1);
        expect(child.compareTo(parent, { [AccountCategory.Cash.type]: 1 })).toBe(1);
    });

    it('should compare children of different parents using parent display orders', () => {
        const firstParent = createAccount('1', AccountCategory.Cash.type, 1);
        const secondParent = createAccount('2', AccountCategory.Cash.type, 2);
        const firstChild = createAccount('11', AccountCategory.Cash.type, 10, '1');
        const secondChild = createAccount('21', AccountCategory.Cash.type, 1, '2');
        const allAccountsMap = { '1': firstParent, '2': secondParent };

        expect(firstChild.compareTo(secondChild, { [AccountCategory.Cash.type]: 1 }, allAccountsMap)).toBe(-1);
        expect(secondChild.compareTo(firstChild, { [AccountCategory.Cash.type]: 1 }, allAccountsMap)).toBe(1);
    });

    it('should compare a child with an unrelated top-level account using its parent display order', () => {
        const parent = createAccount('1', AccountCategory.Cash.type, 1);
        const child = createAccount('11', AccountCategory.Cash.type, 10, '1');
        const other = createAccount('2', AccountCategory.Cash.type, 2);
        const allAccountsMap = { '1': parent };

        expect(child.compareTo(other, { [AccountCategory.Cash.type]: 1 }, allAccountsMap)).toBe(-1);
        expect(other.compareTo(child, { [AccountCategory.Cash.type]: 1 }, allAccountsMap)).toBe(1);
    });

    it('should return zero when different parents have equal display orders', () => {
        const firstParent = createAccount('1', AccountCategory.Cash.type, 2);
        const secondParent = createAccount('2', AccountCategory.Cash.type, 2);
        const firstChild = createAccount('11', AccountCategory.Cash.type, 10, '1');
        const secondChild = createAccount('21', AccountCategory.Cash.type, 1, '2');
        const allAccountsMap = { '1': firstParent, '2': secondParent };

        expect(firstChild.compareTo(secondChild, { [AccountCategory.Cash.type]: 1 }, allAccountsMap)).toBe(0);
    });

    it('should compare account IDs when the account map is not provided', () => {
        const first = createAccount('11', AccountCategory.Cash.type, 10, '1');
        const second = createAccount('21', AccountCategory.Cash.type, 1, '2');

        expect(first.compareTo(second, { [AccountCategory.Cash.type]: 1 })).toBeLessThan(0);
        expect(second.compareTo(first, { [AccountCategory.Cash.type]: 1 })).toBeGreaterThan(0);
    });

    it('should compare account IDs when either parent is missing from the account map', () => {
        const firstParent = createAccount('1', AccountCategory.Cash.type, 10);
        const secondParent = createAccount('2', AccountCategory.Cash.type, 1);
        const firstChild = createAccount('11', AccountCategory.Cash.type, 10, '1');
        const secondChild = createAccount('21', AccountCategory.Cash.type, 1, '2');

        expect(firstChild.compareTo(secondChild, { [AccountCategory.Cash.type]: 1 }, { '1': firstParent })).toBeLessThan(0);
        expect(secondChild.compareTo(firstChild, { [AccountCategory.Cash.type]: 1 }, { '1': firstParent })).toBeGreaterThan(0);
        expect(firstChild.compareTo(secondChild, { [AccountCategory.Cash.type]: 1 }, { '2': secondParent })).toBeLessThan(0);
        expect(secondChild.compareTo(firstChild, { [AccountCategory.Cash.type]: 1 }, { '2': secondParent })).toBeGreaterThan(0);
    });
});
