import { NormalizedText } from '@/core/text.ts';
import { AccountType } from '@/core/account.ts';
import { CategoryType } from '@/core/category.ts';
import { TransactionType } from '@/core/transaction.ts';
import {
    type ImportTransactionReplaceRuleApplyScope,
    type ImportTransactionReplaceRuleConditionValue,
    ImportTransactionReplaceRuleConditionFieldType,
    ImportTransactionReplaceRuleConditionOperatorType,
    ImportTransactionReplaceRuleActionType
} from '@/core/rule.ts';

import { ALL_TIMEZONES } from '@/consts/timezone.ts';
import { TRANSACTION_MAX_TAGS_COUNT } from '@/consts/transaction.ts';
import type { ImportTransaction } from '@/models/imported_transaction.ts';
import {
    type ImportTransactionReplaceRuleMatchData,
    type ImportTransactionReplaceRuleContext,
    type ImportTransactionReplaceRuleApplyResult,
    ImportTransactionReplaceRule
} from '@/models/rule.ts';

import { getTimezoneOffsetMinutes } from '@/lib/datetime.ts';
import { transactionTypeToCategoryType } from '@/lib/category.ts';

export function applyImportTransactionReplaceRules(transactions: ImportTransaction[], rules: ImportTransactionReplaceRule[], scope: ImportTransactionReplaceRuleApplyScope, context: ImportTransactionReplaceRuleContext): ImportTransactionReplaceRuleApplyResult {
    let matchedTransactionCount = 0;
    let updatedTransactionCount = 0;

    for (const transaction of transactions) {
        if (scope === 'selected' && !transaction.selected) {
            continue;
        }

        let matchData: ImportTransactionReplaceRuleMatchData = getTransactionMatchData(transaction, context);
        let shouldUpdateMatchData = false;
        let matched = false;
        let updated = false;

        for (const rule of rules) {
            if (shouldUpdateMatchData) {
                matchData = getTransactionMatchData(transaction, context);
                shouldUpdateMatchData = false;
            }

            if (!rule.isValid() || !isRuleTargetValid(rule, context) || !isRuleMatched(rule, matchData)) {
                continue;
            }

            matched = true;
            const ruleUpdated = applyRuleAction(transaction, rule, context);

            if (ruleUpdated) {
                updated = true;
                shouldUpdateMatchData = true;
            }
        }

        if (matched) {
            matchedTransactionCount++;
        }

        if (updated) {
            updatedTransactionCount++;
            context.updateTransaction(transaction);

            if (transaction.type === TransactionType.Transfer && transaction.sourceAccountId === transaction.destinationAccountId) {
                transaction.valid = false;
            }
        }
    }

    return {
        matchedTransactionCount: matchedTransactionCount,
        updatedTransactionCount: updatedTransactionCount
    };
}

function getTransactionMatchData(transaction: ImportTransaction, context: ImportTransactionReplaceRuleContext): ImportTransactionReplaceRuleMatchData {
    const category = transaction.categoryId ? context.allCategoriesMap[transaction.categoryId] : undefined;
    const sourceAccount = transaction.sourceAccountId ? context.allAccountsMap[transaction.sourceAccountId] : undefined;
    const destinationAccount = transaction.type === TransactionType.Transfer && transaction.destinationAccountId ? context.allAccountsMap[transaction.destinationAccountId] : undefined;

    return {
        type: transaction.type,
        originalCategoryName: category?.name ?? transaction.originalCategoryName,
        originalSourceAccountName: sourceAccount?.name ?? transaction.originalSourceAccountName,
        originalDestinationAccountName: (destinationAccount?.name ?? transaction.originalDestinationAccountName) || '',
        sourceAmount: transaction.sourceAmount,
        destinationAmount: transaction.destinationAmount,
        originalTagNames: transaction.originalTagNames ? transaction.originalTagNames.slice() : [],
        description: NormalizedText.of(transaction.comment || '')
    };
}

function isRuleTargetValid(rule: ImportTransactionReplaceRule, context: ImportTransactionReplaceRuleContext): boolean {
    if (rule.actionType === ImportTransactionReplaceRuleActionType.SetTransactionType) {
        return rule.targetValue === TransactionType.Income || rule.targetValue === TransactionType.Expense || rule.targetValue === TransactionType.Transfer;
    } else if (rule.actionType === ImportTransactionReplaceRuleActionType.SetSourceAmount || rule.actionType === ImportTransactionReplaceRuleActionType.SetDestinationAmount) {
        return rule.targetValue === 1 || rule.targetValue === -1;
    }

    if (typeof rule.targetValue !== 'string') {
        return false;
    }

    if (rule.actionType === ImportTransactionReplaceRuleActionType.SetExpenseCategory) {
        const category = context.allCategoriesMap[rule.targetValue];
        return !!category && !category.hidden && category.parentId !== '0' && category.type === CategoryType.Expense;
    } else if (rule.actionType === ImportTransactionReplaceRuleActionType.SetIncomeCategory) {
        const category = context.allCategoriesMap[rule.targetValue];
        return !!category && !category.hidden && category.parentId !== '0' && category.type === CategoryType.Income;
    } else if (rule.actionType === ImportTransactionReplaceRuleActionType.SetTransferCategory) {
        const category = context.allCategoriesMap[rule.targetValue];
        return !!category && !category.hidden && category.parentId !== '0' && category.type === CategoryType.Transfer;
    } else if (rule.actionType === ImportTransactionReplaceRuleActionType.SetSourceAccount
        || rule.actionType === ImportTransactionReplaceRuleActionType.SetDestinationAccount) {
        const account = context.allAccountsMap[rule.targetValue];
        return !!account && !account.hidden && account.type === AccountType.SingleAccount.type;
    } else if (rule.actionType === ImportTransactionReplaceRuleActionType.AddTag || rule.actionType === ImportTransactionReplaceRuleActionType.ReplaceTag) {
        const tag = context.allTagsMap[rule.targetValue];
        return !!tag && !tag.hidden;
    } else if (rule.actionType === ImportTransactionReplaceRuleActionType.DeleteTag) {
        return true;
    } else if (rule.actionType === ImportTransactionReplaceRuleActionType.SetTimezone) {
        return ALL_TIMEZONES.some(timezone => timezone.timezoneName === rule.targetValue);
    }

    return false;
}

function isRuleMatched(rule: ImportTransactionReplaceRule, transaction: ImportTransactionReplaceRuleMatchData): boolean {
    const value = rule.conditionValue;

    if (rule.conditionField === ImportTransactionReplaceRuleConditionFieldType.ExpenseCategory) {
        return transaction.type === TransactionType.Expense && transaction.originalCategoryName === value;
    } else if (rule.conditionField === ImportTransactionReplaceRuleConditionFieldType.IncomeCategory) {
        return transaction.type === TransactionType.Income && transaction.originalCategoryName === value;
    } else if (rule.conditionField === ImportTransactionReplaceRuleConditionFieldType.TransferCategory) {
        return transaction.type === TransactionType.Transfer && transaction.originalCategoryName === value;
    } else if (rule.conditionField === ImportTransactionReplaceRuleConditionFieldType.SourceAccount) {
        return transaction.originalSourceAccountName === value;
    } else if (rule.conditionField === ImportTransactionReplaceRuleConditionFieldType.DestinationAccount) {
        return transaction.type === TransactionType.Transfer && transaction.originalDestinationAccountName === value;
    } else if (rule.conditionField === ImportTransactionReplaceRuleConditionFieldType.Tag) {
        return typeof value === 'string' && transaction.originalTagNames.includes(value);
    } else if (rule.conditionField === ImportTransactionReplaceRuleConditionFieldType.Amount) {
        return matchAmount(transaction.sourceAmount, rule.conditionOperator, value);
    } else if (rule.conditionField === ImportTransactionReplaceRuleConditionFieldType.TransferInAmount) {
        return transaction.type === TransactionType.Transfer && matchAmount(transaction.destinationAmount, rule.conditionOperator, value);
    } else if (rule.conditionField === ImportTransactionReplaceRuleConditionFieldType.Description || rule.conditionField === ImportTransactionReplaceRuleConditionFieldType.DescriptionCaseInsensitive || rule.conditionField === ImportTransactionReplaceRuleConditionFieldType.DescriptionNormalized) {
        return typeof value === 'string' && matchDescription(transaction.description, rule.conditionField, rule.conditionOperator, value);
    }

    return false;
}

function matchAmount(amount: number, operator: ImportTransactionReplaceRuleConditionOperatorType, value: ImportTransactionReplaceRuleConditionValue): boolean {
    if (!Array.isArray(value)) {
        return false;
    }

    if (operator === ImportTransactionReplaceRuleConditionOperatorType.Equals) {
        return amount === value[0];
    } else if (operator === ImportTransactionReplaceRuleConditionOperatorType.NotEquals) {
        return amount !== value[0];
    } else if (operator === ImportTransactionReplaceRuleConditionOperatorType.GreaterThan) {
        return amount > value[0];
    } else if (operator === ImportTransactionReplaceRuleConditionOperatorType.LessThan) {
        return amount < value[0];
    } else if (operator === ImportTransactionReplaceRuleConditionOperatorType.Between) {
        return amount >= value[0] && amount <= value[1];
    } else if (operator === ImportTransactionReplaceRuleConditionOperatorType.NotBetween) {
        return amount < value[0] || amount > value[1];
    }

    return false;
}

function matchDescription(description: NormalizedText, field: ImportTransactionReplaceRuleConditionFieldType, operator: ImportTransactionReplaceRuleConditionOperatorType, value: string): boolean {
    const normalizedValue = NormalizedText.of(value);
    let descriptionText: string;
    let conditionText: string;

    if (field === ImportTransactionReplaceRuleConditionFieldType.DescriptionCaseInsensitive) {
        descriptionText = description.lowerCaseText;
        conditionText = normalizedValue.lowerCaseText;
    } else if (field === ImportTransactionReplaceRuleConditionFieldType.DescriptionNormalized) {
        descriptionText = description.normalizedText;
        conditionText = normalizedValue.normalizedText;
    } else {
        descriptionText = description.originalText;
        conditionText = normalizedValue.originalText;
    }

    if (operator === ImportTransactionReplaceRuleConditionOperatorType.IsEmpty) {
        return descriptionText.length === 0;
    } else if (operator === ImportTransactionReplaceRuleConditionOperatorType.IsNotEmpty) {
        return descriptionText.length > 0;
    } else if (operator === ImportTransactionReplaceRuleConditionOperatorType.Equals) {
        return descriptionText === conditionText;
    } else if (operator === ImportTransactionReplaceRuleConditionOperatorType.NotEquals) {
        return descriptionText !== conditionText;
    } else if (operator === ImportTransactionReplaceRuleConditionOperatorType.Contains) {
        return descriptionText.includes(conditionText);
    } else if (operator === ImportTransactionReplaceRuleConditionOperatorType.NotContains) {
        return !descriptionText.includes(conditionText);
    } else if (operator === ImportTransactionReplaceRuleConditionOperatorType.StartsWith) {
        return descriptionText.startsWith(conditionText);
    } else if (operator === ImportTransactionReplaceRuleConditionOperatorType.NotStartsWith) {
        return !descriptionText.startsWith(conditionText);
    } else if (operator === ImportTransactionReplaceRuleConditionOperatorType.EndsWith) {
        return descriptionText.endsWith(conditionText);
    } else if (operator === ImportTransactionReplaceRuleConditionOperatorType.NotEndsWith) {
        return !descriptionText.endsWith(conditionText);
    } else if (operator === ImportTransactionReplaceRuleConditionOperatorType.RegexMatch || operator === ImportTransactionReplaceRuleConditionOperatorType.NotRegexMatch) {
        let matched = false;

        try {
            const flags = (field === ImportTransactionReplaceRuleConditionFieldType.DescriptionCaseInsensitive || field === ImportTransactionReplaceRuleConditionFieldType.DescriptionNormalized) ? 'i' : undefined;
            matched = new RegExp(conditionText, flags).test(descriptionText);
        } catch {
            matched = false;
        }

        return operator === ImportTransactionReplaceRuleConditionOperatorType.RegexMatch ? matched : !matched;
    }

    return false;
}

function applyRuleAction(transaction: ImportTransaction, rule: ImportTransactionReplaceRule, context: ImportTransactionReplaceRuleContext): boolean {
    if (rule.actionType === ImportTransactionReplaceRuleActionType.SetTransactionType) {
        return setTransactionType(transaction, rule.targetValue as TransactionType, context);
    } else if (rule.actionType === ImportTransactionReplaceRuleActionType.SetExpenseCategory) {
        return setTransactionCategory(transaction, TransactionType.Expense, rule.targetValue as string);
    } else if (rule.actionType === ImportTransactionReplaceRuleActionType.SetIncomeCategory) {
        return setTransactionCategory(transaction, TransactionType.Income, rule.targetValue as string);
    } else if (rule.actionType === ImportTransactionReplaceRuleActionType.SetTransferCategory) {
        return setTransactionCategory(transaction, TransactionType.Transfer, rule.targetValue as string);
    } else if (rule.actionType === ImportTransactionReplaceRuleActionType.SetSourceAccount) {
        return setTransactionSourceAccount(transaction, rule.targetValue as string);
    } else if (rule.actionType === ImportTransactionReplaceRuleActionType.SetDestinationAccount) {
        return setTransactionDestinationAccount(transaction, rule.targetValue as string);
    } else if (rule.actionType === ImportTransactionReplaceRuleActionType.AddTag) {
        const tagId = rule.targetValue as string;
        return appendTransactionTag(transaction, tagId, context.allTagsMap[tagId]?.name || '');
    } else if (rule.actionType === ImportTransactionReplaceRuleActionType.ReplaceTag) {
        const tagId = rule.targetValue as string;
        return replaceMatchedTransactionTag(transaction, rule.conditionValue as string, tagId, context.allTagsMap[tagId]?.name || '');
    } else if (rule.actionType === ImportTransactionReplaceRuleActionType.DeleteTag) {
        return deleteMatchedTransactionTag(transaction, rule.conditionValue as string);
    } else if (rule.actionType === ImportTransactionReplaceRuleActionType.SetSourceAmount) {
        return setTransactionSourceAmountSign(transaction, rule.targetValue as number);
    } else if (rule.actionType === ImportTransactionReplaceRuleActionType.SetDestinationAmount) {
        return setTransactionDestinationAmountSign(transaction, rule.targetValue as number);
    } else if (rule.actionType === ImportTransactionReplaceRuleActionType.SetTimezone) {
        return setTransactionTimezone(transaction, rule.targetValue as string);
    }

    return false;
}

function setTransactionTimezone(transaction: ImportTransaction, timezone: string): boolean {
    const utcOffset = getTimezoneOffsetMinutes(transaction.time, timezone);

    if (transaction.utcOffset === utcOffset) {
        return false;
    }

    transaction.utcOffset = utcOffset;
    return true;
}

function setTransactionType(transaction: ImportTransaction, type: TransactionType, context: ImportTransactionReplaceRuleContext): boolean {
    const fromType = transaction.type;

    if (fromType === type) {
        return false;
    }

    const categoryType = transactionTypeToCategoryType(type);

    if (!categoryType) {
        return false;
    }

    transaction.type = type;
    transaction.categoryId = context.allSecondaryCategoriesMapByName[categoryType]?.[transaction.originalCategoryName]?.id || '0';

    if (type === TransactionType.Transfer) {
        transaction.destinationAccountId = context.allAccountsMapByName[transaction.originalDestinationAccountName || '']?.id || '0';
        transaction.destinationAmount = transaction.sourceAmount;
    } else {
        if (fromType === TransactionType.Transfer && type === TransactionType.Income) {
            transaction.sourceAccountId = transaction.destinationAccountId;
            transaction.sourceAmount = transaction.destinationAmount;
        }

        transaction.destinationAccountId = '0';
        transaction.destinationAmount = 0;
    }

    return true;
}

function setTransactionCategory(transaction: ImportTransaction, type: TransactionType, categoryId: string): boolean {
    if (transaction.type !== type || transaction.categoryId === categoryId) {
        return false;
    }

    transaction.categoryId = categoryId;
    return true;
}

function setTransactionSourceAccount(transaction: ImportTransaction, accountId: string): boolean {
    if (transaction.sourceAccountId === accountId) {
        return false;
    }

    transaction.sourceAccountId = accountId;
    return true;
}

function setTransactionDestinationAccount(transaction: ImportTransaction, accountId: string): boolean {
    if (transaction.type !== TransactionType.Transfer || transaction.destinationAccountId === accountId) {
        return false;
    }

    transaction.destinationAccountId = accountId;
    return true;
}

function appendTransactionTag(transaction: ImportTransaction, tagId: string, tagName: string): boolean {
    if (transaction.tagIds.includes(tagId) || transaction.tagIds.length >= TRANSACTION_MAX_TAGS_COUNT) {
        return false;
    }

    transaction.tagIds.push(tagId);
    transaction.originalTagNames.push(tagName);
    return true;
}

function replaceMatchedTransactionTag(transaction: ImportTransaction, sourceValue: string, tagId: string, tagName: string): boolean {
    let updated = false;

    for (let i = 0; i < transaction.originalTagNames.length; i++) {
        if (transaction.originalTagNames[i] === sourceValue && transaction.tagIds[i] !== tagId) {
            transaction.tagIds[i] = tagId;
            transaction.originalTagNames[i] = tagName;
            updated = true;
        }
    }

    return updated;
}

function deleteMatchedTransactionTag(transaction: ImportTransaction, sourceValue: string): boolean {
    let updated = false;

    for (let i = transaction.originalTagNames.length - 1; i >= 0; i--) {
        if (transaction.originalTagNames[i] === sourceValue) {
            transaction.tagIds.splice(i, 1);
            transaction.originalTagNames.splice(i, 1);
            updated = true;
        }
    }

    return updated;
}

function setTransactionSourceAmountSign(transaction: ImportTransaction, sign: number): boolean {
    const amount = Math.abs(transaction.sourceAmount) * sign;

    if (transaction.sourceAmount === amount) {
        return false;
    }

    transaction.sourceAmount = amount;
    return true;
}

function setTransactionDestinationAmountSign(transaction: ImportTransaction, sign: number): boolean {
    if (transaction.type !== TransactionType.Transfer) {
        return false;
    }

    const amount = Math.abs(transaction.destinationAmount) * sign;

    if (transaction.destinationAmount === amount) {
        return false;
    }

    transaction.destinationAmount = amount;
    return true;
}
