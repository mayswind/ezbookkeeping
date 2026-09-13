import { describe, expect, test } from 'vitest';

import { TransactionType } from '@/core/transaction.ts';
import {
    ImportTransactionReplaceRuleConditionFieldType,
    ImportTransactionReplaceRuleConditionOperatorType,
    ImportTransactionReplaceRuleActionType
} from '@/core/rule.ts';

import {
    ImportTransactionReplaceRule,
    ImportTransactionReplaceRules
} from '@/models/rule.ts';

describe('rule', () => {
    test('should serialize and parse new rule format', () => {
        const rule = ImportTransactionReplaceRule.of(
            'id',
            'Coffee',
            ImportTransactionReplaceRuleConditionFieldType.DescriptionNormalized,
            ImportTransactionReplaceRuleConditionOperatorType.Contains,
            'coffee',
            ImportTransactionReplaceRuleActionType.SetExpenseCategory,
            'category-id'
        );
        const parsed = ImportTransactionReplaceRules.parseFromJson(ImportTransactionReplaceRules.of([rule]).toJson(), () => 'generated-id');

        expect(parsed).not.toBeNull();
        expect(parsed?.getRules()).toHaveLength(1);
        expect(parsed?.getRules()[0]).toEqual(rule);
    });

    test('should parse old rule format', () => {
        const json = JSON.stringify({
            ezBookkeepingImportTransactionReplaceRules: [
                {type: 'expenseCategory', sourceValue: 'Food', targetId: 'category-id'},
                {type: 'tag', sourceValue: 'Old', targetId: 'tag-id'}
            ]
        });
        const parsed = ImportTransactionReplaceRules.parseFromJson(json, () => 'generated-id');

        expect(parsed).not.toBeNull();
        expect(parsed?.getRules()).toHaveLength(2);
        expect(parsed?.getRules()[0]?.conditionField).toBe(ImportTransactionReplaceRuleConditionFieldType.ExpenseCategory);
        expect(parsed?.getRules()[0]?.conditionOperator).toBe(ImportTransactionReplaceRuleConditionOperatorType.Is);
        expect(parsed?.getRules()[0]?.actionType).toBe(ImportTransactionReplaceRuleActionType.SetExpenseCategory);
        expect(parsed?.getRules()[1]?.actionType).toBe(ImportTransactionReplaceRuleActionType.ReplaceTag);
    });

    test('should normalize the second amount value for non-range operators', () => {
        const rule = ImportTransactionReplaceRule.of(
            'id', '', ImportTransactionReplaceRuleConditionFieldType.Amount,
            ImportTransactionReplaceRuleConditionOperatorType.GreaterThan, [100, 200],
            ImportTransactionReplaceRuleActionType.SetTransactionType, TransactionType.Expense
        );
        const parsed = JSON.parse(ImportTransactionReplaceRules.of([rule]).toJson());

        expect(parsed.ezBookkeepingImportTransactionReplaceRules[0].conditionValue).toEqual([100, 100]);
    });

    test('should reject invalid condition and action values', () => {
        const invalidDescription = ImportTransactionReplaceRule.of(
            'id', '', ImportTransactionReplaceRuleConditionFieldType.Description,
            ImportTransactionReplaceRuleConditionOperatorType.Is, 'value',
            ImportTransactionReplaceRuleActionType.SetTransactionType, TransactionType.Expense
        );
        const invalidTransactionType = ImportTransactionReplaceRule.of(
            'id', '', ImportTransactionReplaceRuleConditionFieldType.Description,
            ImportTransactionReplaceRuleConditionOperatorType.Contains, 'value',
            ImportTransactionReplaceRuleActionType.SetTransactionType, TransactionType.ModifyBalance
        );
        const invalidReplaceTag = ImportTransactionReplaceRule.of(
            'id', '', ImportTransactionReplaceRuleConditionFieldType.Description,
            ImportTransactionReplaceRuleConditionOperatorType.Contains, 'value',
            ImportTransactionReplaceRuleActionType.ReplaceTag, 'tag-id'
        );
        const validDeleteTag = ImportTransactionReplaceRule.of(
            'id', '', ImportTransactionReplaceRuleConditionFieldType.Tag,
            ImportTransactionReplaceRuleConditionOperatorType.Contains, 'tag',
            ImportTransactionReplaceRuleActionType.DeleteTag, ''
        );
        const invalidTagOperator = ImportTransactionReplaceRule.of(
            'id', '', ImportTransactionReplaceRuleConditionFieldType.Tag,
            ImportTransactionReplaceRuleConditionOperatorType.Is, 'tag',
            ImportTransactionReplaceRuleActionType.DeleteTag, ''
        );
        const emptyTagCondition = ImportTransactionReplaceRule.of(
            'id', '', ImportTransactionReplaceRuleConditionFieldType.Tag,
            ImportTransactionReplaceRuleConditionOperatorType.Contains, '',
            ImportTransactionReplaceRuleActionType.DeleteTag, ''
        );
        const emptyItemCondition = ImportTransactionReplaceRule.of(
            'id', '', ImportTransactionReplaceRuleConditionFieldType.ExpenseCategory,
            ImportTransactionReplaceRuleConditionOperatorType.Is, '',
            ImportTransactionReplaceRuleActionType.SetExpenseCategory, 'category-id'
        );
        const emptyDescriptionCondition = ImportTransactionReplaceRule.of(
            'id', '', ImportTransactionReplaceRuleConditionFieldType.Description,
            ImportTransactionReplaceRuleConditionOperatorType.IsEmpty, '',
            ImportTransactionReplaceRuleActionType.SetTransactionType, TransactionType.Expense
        );
        const emptyDescriptionValue = ImportTransactionReplaceRule.of(
            'id', '', ImportTransactionReplaceRuleConditionFieldType.Description,
            ImportTransactionReplaceRuleConditionOperatorType.Equals, '',
            ImportTransactionReplaceRuleActionType.SetTransactionType, TransactionType.Expense
        );
        const regexDescriptionCondition = ImportTransactionReplaceRule.of(
            'id', '', ImportTransactionReplaceRuleConditionFieldType.DescriptionNormalized,
            ImportTransactionReplaceRuleConditionOperatorType.RegexMatch, '^coffee$',
            ImportTransactionReplaceRuleActionType.SetTransactionType, TransactionType.Expense
        );

        expect(invalidDescription.isValid()).toBe(false);
        expect(invalidTransactionType.isValid()).toBe(false);
        expect(invalidReplaceTag.isValid()).toBe(false);
        expect(validDeleteTag.isValid()).toBe(true);
        expect(invalidTagOperator.isValid()).toBe(false);
        expect(emptyTagCondition.isValid()).toBe(false);
        expect(emptyItemCondition.isValid()).toBe(false);
        expect(emptyDescriptionCondition.isValid()).toBe(true);
        expect(emptyDescriptionValue.isValid()).toBe(false);
        expect(regexDescriptionCondition.isValid()).toBe(true);
    });

    test('should reject malformed json root', () => {
        expect(ImportTransactionReplaceRules.parseFromJson('{}', () => 'id')).toBeNull();
        expect(ImportTransactionReplaceRules.parseFromJson('invalid', () => 'id')).toBeNull();
    });

    test('should reject a rule file containing an invalid rule', () => {
        const json = JSON.stringify({
            ezBookkeepingImportTransactionReplaceRules: [
                {
                    id: 'valid-id',
                    name: '',
                    conditionField: ImportTransactionReplaceRuleConditionFieldType.Description,
                    conditionOperator: ImportTransactionReplaceRuleConditionOperatorType.Contains,
                    conditionValue: 'value',
                    actionType: ImportTransactionReplaceRuleActionType.SetTransactionType,
                    targetValue: TransactionType.Expense
                },
                { conditionField: ImportTransactionReplaceRuleConditionFieldType.Tag }
            ]
        });

        expect(ImportTransactionReplaceRules.parseFromJson(json, () => 'id')).toBeNull();
    });
});
