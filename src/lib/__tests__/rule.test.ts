import { describe, expect, test } from 'vitest';

import { AccountCategory } from '@/core/account.ts';
import { CategoryType } from '@/core/category.ts';
import { TransactionType } from '@/core/transaction.ts';
import {
    ImportTransactionReplaceRuleConditionFieldType,
    ImportTransactionReplaceRuleConditionOperatorType,
    ImportTransactionReplaceRuleActionType
} from '@/core/rule.ts';

import { Account } from '@/models/account.ts';
import { TransactionCategory } from '@/models/transaction_category.ts';
import { TransactionTag } from '@/models/transaction_tag.ts';
import { ImportTransaction } from '@/models/imported_transaction.ts';
import { type ImportTransactionReplaceRuleContext, ImportTransactionReplaceRule } from '@/models/rule.ts';

import { applyImportTransactionReplaceRules } from '@/lib/rule.ts';

function createTransaction(type: TransactionType): ImportTransaction {
    return ImportTransaction.of({
        type: type,
        categoryId: 'old-category',
        originalCategoryName: 'Original Category',
        time: 0,
        utcOffset: 0,
        sourceAccountId: 'source-account',
        originalSourceAccountName: 'Source Account',
        originalSourceAccountCurrency: 'USD',
        destinationAccountId: type === TransactionType.Transfer ? 'destination-account' : undefined,
        originalDestinationAccountName: type === TransactionType.Transfer ? 'Destination Account' : undefined,
        originalDestinationAccountCurrency: type === TransactionType.Transfer ? 'EUR' : undefined,
        sourceAmount: 100,
        destinationAmount: type === TransactionType.Transfer ? 80 : undefined,
        tagIds: ['old-tag'],
        originalTagNames: ['Original Tag'],
        comment: 'Café PAYMENT'
    }, 0);
}

function createCategory(id: string, name: string, type: CategoryType): TransactionCategory {
    return TransactionCategory.of({
        id: id,
        name: name,
        parentId: 'parent-id',
        type: type,
        icon: '1',
        iconType: 1,
        color: '#000000',
        comment: '',
        displayOrder: 0,
        hidden: false
    });
}

function createAccount(id: string, name: string): Account {
    const account = Account.createNewAccount(AccountCategory.Cash, 'USD', 0);
    account.id = id;
    account.name = name;
    return account;
}

function createTag(id: string, name: string): TransactionTag {
    const tag = TransactionTag.createNewTag(name);
    tag.id = id;
    return tag;
}

function createContext(updateTransaction: (transaction: ImportTransaction) => void = () => {}): ImportTransactionReplaceRuleContext {
    const expenseCategory = createCategory('expense-category', 'Expense Category', CategoryType.Expense);
    const incomeCategory = createCategory('income-category', 'Income Category', CategoryType.Income);
    const transferCategory = createCategory('transfer-category', 'Transfer Category', CategoryType.Transfer);
    const sourceAccount = createAccount('source-account', 'Source Account');
    const destinationAccount = createAccount('destination-account', 'Destination Account');
    const replacementAccount = createAccount('source-account-2', 'Replacement Account');
    const targetAccount = createAccount('target-account', 'Target Account');
    const oldTag = createTag('old-tag', 'Original Tag');
    const targetTag = createTag('target-tag', 'Target Tag');
    const replacementTag = createTag('target-tag-2', 'Replacement Tag');

    return {
        allCategoriesMap: {
            [expenseCategory.id]: expenseCategory,
            [incomeCategory.id]: incomeCategory,
            [transferCategory.id]: transferCategory
        },
        allSecondaryCategoriesMapByName: {
            [CategoryType.Expense]: {'Original Category': expenseCategory},
            [CategoryType.Income]: {'Original Category': incomeCategory},
            [CategoryType.Transfer]: {'Original Category': transferCategory}
        },
        allAccountsMap: {
            [sourceAccount.id]: sourceAccount,
            [destinationAccount.id]: destinationAccount,
            [targetAccount.id]: targetAccount,
            [replacementAccount.id]: replacementAccount
        },
        allAccountsMapByName: {
            [sourceAccount.name]: sourceAccount,
            [destinationAccount.name]: destinationAccount,
            [targetAccount.name]: targetAccount
        },
        allTagsMap: {
            [oldTag.id]: oldTag,
            [targetTag.id]: targetTag,
            [replacementTag.id]: replacementTag
        },
        updateTransaction: updateTransaction
    };
}

function createRule(field: ImportTransactionReplaceRuleConditionFieldType, operator: ImportTransactionReplaceRuleConditionOperatorType, value: string | [number, number], action: ImportTransactionReplaceRuleActionType, target: string | number): ImportTransactionReplaceRule {
    return ImportTransactionReplaceRule.of('id', '', field, operator, value, action, target);
}

describe('rule', () => {
    test('should match normalized description and execute rules in order', () => {
        const transaction = createTransaction(TransactionType.Expense);
        let updateCount = 0;
        const rules = [
            createRule(ImportTransactionReplaceRuleConditionFieldType.DescriptionNormalized, ImportTransactionReplaceRuleConditionOperatorType.StartsWith,
                'cafe', ImportTransactionReplaceRuleActionType.SetExpenseCategory, 'expense-category'),
            createRule(ImportTransactionReplaceRuleConditionFieldType.Amount, ImportTransactionReplaceRuleConditionOperatorType.Between,
                [50, 100], ImportTransactionReplaceRuleActionType.SetSourceAccount, 'source-account-2')
        ];

        const result = applyImportTransactionReplaceRules([transaction], rules, 'all', createContext(() => updateCount++));

        expect(result.updatedTransactionCount).toBe(1);
        expect(transaction.categoryId).toBe('expense-category');
        expect(transaction.sourceAccountId).toBe('source-account-2');
        expect(updateCount).toBe(1);
    });

    test('should match each rule against the result of preceding rules', () => {
        const transaction = createTransaction(TransactionType.Transfer);
        transaction.sourceAmount = 2500;
        transaction.destinationAmount = 2500;
        const rules = [
            createRule(ImportTransactionReplaceRuleConditionFieldType.Amount, ImportTransactionReplaceRuleConditionOperatorType.Between,
                [2000, 3000], ImportTransactionReplaceRuleActionType.SetSourceAmount, -1),
            createRule(ImportTransactionReplaceRuleConditionFieldType.Amount, ImportTransactionReplaceRuleConditionOperatorType.Between,
                [2000, 3000], ImportTransactionReplaceRuleActionType.SetDestinationAmount, -1)
        ];

        const result = applyImportTransactionReplaceRules([transaction], rules, 'all', createContext());

        expect(result.matchedTransactionCount).toBe(1);
        expect(result.updatedTransactionCount).toBe(1);
        expect(transaction.sourceAmount).toBe(-2500);
        expect(transaction.destinationAmount).toBe(2500);
    });

    test('should match updated account values in subsequent rules', () => {
        const transaction = createTransaction(TransactionType.Transfer);
        const rules = [
            createRule(ImportTransactionReplaceRuleConditionFieldType.SourceAccount, ImportTransactionReplaceRuleConditionOperatorType.Is,
                'Source Account', ImportTransactionReplaceRuleActionType.SetDestinationAccount, 'target-account'),
            createRule(ImportTransactionReplaceRuleConditionFieldType.DestinationAccount, ImportTransactionReplaceRuleConditionOperatorType.Is,
                'Target Account', ImportTransactionReplaceRuleActionType.SetSourceAccount, 'source-account-2')
        ];

        applyImportTransactionReplaceRules([transaction], rules, 'all', createContext());

        expect(transaction.sourceAccountId).toBe('source-account-2');
        expect(transaction.destinationAccountId).toBe('target-account');
    });

    test('should add tags only once', () => {
        const transaction = createTransaction(TransactionType.Expense);
        const appendRule = createRule(ImportTransactionReplaceRuleConditionFieldType.Description, ImportTransactionReplaceRuleConditionOperatorType.Contains,
            'PAYMENT', ImportTransactionReplaceRuleActionType.AddTag, 'target-tag');

        applyImportTransactionReplaceRules([transaction], [appendRule], 'all', createContext());
        applyImportTransactionReplaceRules([transaction], [appendRule], 'all', createContext());

        expect(transaction.tagIds).toEqual(['old-tag', 'target-tag']);
        expect(transaction.originalTagNames).toEqual(['Original Tag', 'Target Tag']);
    });

    test('should replace and delete matched tags', () => {
        const replacedTransaction = createTransaction(TransactionType.Expense);
        const deletedTransaction = createTransaction(TransactionType.Expense);
        const orderedTransaction = createTransaction(TransactionType.Expense);
        orderedTransaction.tagIds = ['old-tag', 'target-tag'];
        orderedTransaction.originalTagNames = ['Original Tag', 'Second Tag'];
        const replaceRule = createRule(ImportTransactionReplaceRuleConditionFieldType.Tag, ImportTransactionReplaceRuleConditionOperatorType.Contains,
            'Original Tag', ImportTransactionReplaceRuleActionType.ReplaceTag, 'target-tag-2');
        const deleteRule = createRule(ImportTransactionReplaceRuleConditionFieldType.Tag, ImportTransactionReplaceRuleConditionOperatorType.Contains,
            'Original Tag', ImportTransactionReplaceRuleActionType.DeleteTag, '');
        const replaceSecondRule = createRule(ImportTransactionReplaceRuleConditionFieldType.Tag, ImportTransactionReplaceRuleConditionOperatorType.Contains,
            'Second Tag', ImportTransactionReplaceRuleActionType.ReplaceTag, 'target-tag-2');

        applyImportTransactionReplaceRules([replacedTransaction], [replaceRule], 'all', createContext());
        applyImportTransactionReplaceRules([deletedTransaction], [deleteRule], 'all', createContext());
        applyImportTransactionReplaceRules([orderedTransaction], [deleteRule, replaceSecondRule], 'all', createContext());

        expect(replacedTransaction.tagIds).toEqual(['target-tag-2']);
        expect(replacedTransaction.originalTagNames).toEqual(['Replacement Tag']);
        expect(deletedTransaction.tagIds).toEqual([]);
        expect(deletedTransaction.originalTagNames).toEqual([]);
        expect(orderedTransaction.tagIds).toEqual(['target-tag-2']);
        expect(orderedTransaction.originalTagNames).toEqual(['Replacement Tag']);
    });

    test('should apply only to selected transactions and fully convert transaction type', () => {
        const selected = createTransaction(TransactionType.Transfer);
        const unselected = createTransaction(TransactionType.Transfer);
        selected.selected = true;
        const rule = createRule(ImportTransactionReplaceRuleConditionFieldType.TransferInAmount, ImportTransactionReplaceRuleConditionOperatorType.Equals,
            [80, 80], ImportTransactionReplaceRuleActionType.SetTransactionType, TransactionType.Income);

        const result = applyImportTransactionReplaceRules([selected, unselected], [rule], 'selected', createContext());

        expect(result.updatedTransactionCount).toBe(1);
        expect(selected.type).toBe(TransactionType.Income);
        expect(selected.categoryId).toBe('income-category');
        expect(selected.sourceAccountId).toBe('destination-account');
        expect(selected.sourceAmount).toBe(80);
        expect(selected.destinationAccountId).toBe('0');
        expect(unselected.type).toBe(TransactionType.Transfer);
    });
});
