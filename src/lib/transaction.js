import { CategoryType } from '@/core/category.ts';
import { TransactionType } from '@/core/transaction.ts';
import {
    isNumber
} from './common.ts';
import {
    getBrowserTimezoneOffsetMinutes,
    getDummyUnixTimeForLocalUsage
} from './datetime.js';
import {
    categoryTypeToTransactionType,
    isSubCategoryIdAvailable,
    getFirstAvailableCategoryId,
    getFirstAvailableSubCategoryId
} from './category.js';

function getDisplayAmount(amount, currency, hideAmount, formatAmountWithCurrencyFunc) {
    if (hideAmount) {
        return formatAmountWithCurrencyFunc('***', currency);
    }

    return formatAmountWithCurrencyFunc(amount, currency);
}

export function setTransactionModelByTransaction(transaction, transaction2, allCategories, allCategoriesMap, allVisibleAccounts, allAccountsMap, allTagsMap, defaultAccountId, options, setContextData, convertContextTime) {
    if ((!options.type || options.type === '0') && options.categoryId && options.categoryId !== '0' && allCategoriesMap[options.categoryId]) {
        const category = allCategoriesMap[options.categoryId];
        const type = categoryTypeToTransactionType(category.type);

        if (isNumber(type)) {
            transaction.type = type;
        }
    }

    if (allCategories[CategoryType.Expense] &&
        allCategories[CategoryType.Expense].length) {
        if (options.categoryId && options.categoryId !== '0') {
            if (isSubCategoryIdAvailable(allCategories[CategoryType.Expense], options.categoryId)) {
                transaction.expenseCategory = options.categoryId;
            } else {
                transaction.expenseCategory = getFirstAvailableSubCategoryId(allCategories[CategoryType.Expense], options.categoryId);
            }
        }

        if (!transaction.expenseCategory) {
            transaction.expenseCategory = getFirstAvailableCategoryId(allCategories[CategoryType.Expense]);
        }
    }

    if (allCategories[CategoryType.Income] &&
        allCategories[CategoryType.Income].length) {
        if (options.categoryId && options.categoryId !== '0') {
            if (isSubCategoryIdAvailable(allCategories[CategoryType.Income], options.categoryId)) {
                transaction.incomeCategory = options.categoryId;
            } else {
                transaction.incomeCategory = getFirstAvailableSubCategoryId(allCategories[CategoryType.Income], options.categoryId);
            }
        }

        if (!transaction.incomeCategory) {
            transaction.incomeCategory = getFirstAvailableCategoryId(allCategories[CategoryType.Income]);
        }
    }

    if (allCategories[CategoryType.Transfer] &&
        allCategories[CategoryType.Transfer].length) {
        if (options.categoryId && options.categoryId !== '0') {
            if (isSubCategoryIdAvailable(allCategories[CategoryType.Transfer], options.categoryId)) {
                transaction.transferCategory = options.categoryId;
            } else {
                transaction.transferCategory = getFirstAvailableSubCategoryId(allCategories[CategoryType.Transfer], options.categoryId);
            }
        }

        if (!transaction.transferCategory) {
            transaction.transferCategory = getFirstAvailableCategoryId(allCategories[CategoryType.Transfer]);
        }
    }

    if (allVisibleAccounts.length) {
        if (options.accountId && options.accountId !== '0') {
            for (let i = 0; i < allVisibleAccounts.length; i++) {
                if (allVisibleAccounts[i].id === options.accountId) {
                    transaction.sourceAccountId = options.accountId;
                    transaction.destinationAccountId = options.accountId;
                    break;
                }
            }
        }

        if (!transaction.sourceAccountId) {
            if (defaultAccountId && allAccountsMap[defaultAccountId]) {
                transaction.sourceAccountId = defaultAccountId;
            } else {
                transaction.sourceAccountId = allVisibleAccounts[0].id;
            }
        }

        if (!transaction.destinationAccountId) {
            if (defaultAccountId && allAccountsMap[defaultAccountId]) {
                transaction.destinationAccountId = defaultAccountId;
            } else {
                transaction.destinationAccountId = allVisibleAccounts[0].id;
            }
        }
    }

    if (allTagsMap && options.tagIds) {
        const tagIds = options.tagIds.split(',');
        const finalTagIds = [];

        for (let i = 0; i < tagIds.length; i++) {
            const tagId = tagIds[i];
            const tag = allTagsMap[tagId];

            if (tag && !tag.hidden) {
                finalTagIds.push(tag.id);
            }
        }

        transaction.tagIds = finalTagIds;
    }

    if (transaction2) {
        if (setContextData) {
            transaction.id = transaction2.id;
        }

        transaction.type = transaction2.type;

        if (transaction.type === TransactionType.Expense) {
            transaction.expenseCategory = transaction2.categoryId || '';
        } else if (transaction.type === TransactionType.Income) {
            transaction.incomeCategory = transaction2.categoryId || '';
        } else if (transaction.type === TransactionType.Transfer) {
            transaction.transferCategory = transaction2.categoryId || '';
        }

        if (setContextData) {
            transaction.utcOffset = transaction2.utcOffset;
            transaction.timeZone = transaction2.timeZone;

            if (convertContextTime) {
                transaction.time = getDummyUnixTimeForLocalUsage(transaction2.time, transaction.utcOffset, getBrowserTimezoneOffsetMinutes());
            } else {
                transaction.time = transaction2.time;
            }
        }

        transaction.sourceAccountId = transaction2.sourceAccountId;

        if (transaction2.destinationAccountId) {
            transaction.destinationAccountId = transaction2.destinationAccountId;
        } else {
            transaction.destinationAccountId = '';
        }

        transaction.sourceAmount = transaction2.sourceAmount;

        if (transaction2.destinationAmount) {
            transaction.destinationAmount = transaction2.destinationAmount;
        } else {
            transaction.destinationAmount = 0;
        }

        transaction.hideAmount = transaction2.hideAmount;
        transaction.tagIds = transaction2.tagIds || [];
        transaction.pictures = transaction2.pictures || [];

        transaction.comment = transaction2.comment;

        if (setContextData) {
            transaction.geoLocation = transaction2.geoLocation;
        }
    }
}

export function getTransactionDisplayAmount(transaction, allFilterAccountIdsCount, allFilterAccountIds, formatAmountWithCurrencyFunc) {
    if (allFilterAccountIdsCount < 1) {
        if (transaction.sourceAccount) {
            return getDisplayAmount(transaction.sourceAmount, transaction.sourceAccount.currency, transaction.hideAmount, formatAmountWithCurrencyFunc);
        }
    } else if (allFilterAccountIdsCount === 1) {
        if (transaction.sourceAccount && (allFilterAccountIds[transaction.sourceAccount.id] || allFilterAccountIds[transaction.sourceAccount.parentId])) {
            return getDisplayAmount(transaction.sourceAmount, transaction.sourceAccount.currency, transaction.hideAmount , formatAmountWithCurrencyFunc);
        } else if (transaction.destinationAccount && (allFilterAccountIds[transaction.destinationAccount.id] || allFilterAccountIds[transaction.destinationAccount.parentId])) {
            return getDisplayAmount(transaction.destinationAmount, transaction.destinationAccount.currency, transaction.hideAmount , formatAmountWithCurrencyFunc);
        }
    } else { // allFilterAccountIdsCount > 1
        if (transaction.sourceAccount && transaction.destinationAccount) {
            if ((allFilterAccountIds[transaction.sourceAccount.id] || allFilterAccountIds[transaction.sourceAccount.parentId])
                && !allFilterAccountIds[transaction.destinationAccount.id] && !allFilterAccountIds[transaction.destinationAccount.parentId]) {
                return getDisplayAmount(transaction.sourceAmount, transaction.sourceAccount.currency, transaction.hideAmount , formatAmountWithCurrencyFunc);
            } else if ((allFilterAccountIds[transaction.destinationAccount.id] || allFilterAccountIds[transaction.destinationAccount.parentId])
                && !allFilterAccountIds[transaction.sourceAccount.id] && !allFilterAccountIds[transaction.sourceAccount.parentId]) {
                return getDisplayAmount(transaction.destinationAmount, transaction.destinationAccount.currency, transaction.hideAmount , formatAmountWithCurrencyFunc);
            }
        }
    }

    if (transaction.sourceAccount) {
        return getDisplayAmount(transaction.sourceAmount, transaction.sourceAccount.currency, transaction.hideAmount, formatAmountWithCurrencyFunc);
    }

    return '';
}
