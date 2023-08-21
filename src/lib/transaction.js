import categoryConstants from '@/consts/category.js';
import transactionConstants from '@/consts/transaction.js';
import {
    isNumber
} from './common.js';
import {
    getBrowserTimezoneOffsetMinutes,
    getDummyUnixTimeForLocalUsage
} from './datetime.js';
import {
    categoryTypeToTransactionType,
    isSubCategoryIdAvailable,
    getFirstAvailableCategoryId,
    getFirstAvaiableSubCategoryId
} from './category.js';

export function setTransactionModelByTransaction(transaction, transaction2, allCategories, allCategoriesMap, allVisibleAccounts, allAccountsMap, defaultAccountId, options, setContextData) {
    if ((!options.type || options.type === '0') && options.categoryId && options.categoryId !== '0' && allCategoriesMap[options.categoryId]) {
        const category = allCategoriesMap[options.categoryId];
        const type = categoryTypeToTransactionType(category.type);

        if (isNumber(type)) {
            transaction.type = type;
        }
    }

    if (allCategories[categoryConstants.allCategoryTypes.Expense] &&
        allCategories[categoryConstants.allCategoryTypes.Expense].length) {
        if (options.categoryId && options.categoryId !== '0') {
            if (isSubCategoryIdAvailable(allCategories[categoryConstants.allCategoryTypes.Expense], options.categoryId)) {
                transaction.expenseCategory = options.categoryId;
            } else {
                transaction.expenseCategory = getFirstAvaiableSubCategoryId(allCategories[categoryConstants.allCategoryTypes.Expense], options.categoryId);
            }
        }

        if (!transaction.expenseCategory) {
            transaction.expenseCategory = getFirstAvailableCategoryId(allCategories[categoryConstants.allCategoryTypes.Expense]);
        }
    }

    if (allCategories[categoryConstants.allCategoryTypes.Income] &&
        allCategories[categoryConstants.allCategoryTypes.Income].length) {
        if (options.categoryId && options.categoryId !== '0') {
            if (isSubCategoryIdAvailable(allCategories[categoryConstants.allCategoryTypes.Income], options.categoryId)) {
                transaction.incomeCategory = options.categoryId;
            } else {
                transaction.incomeCategory = getFirstAvaiableSubCategoryId(allCategories[categoryConstants.allCategoryTypes.Income], options.categoryId);
            }
        }

        if (!transaction.incomeCategory) {
            transaction.incomeCategory = getFirstAvailableCategoryId(allCategories[categoryConstants.allCategoryTypes.Income]);
        }
    }

    if (allCategories[categoryConstants.allCategoryTypes.Transfer] &&
        allCategories[categoryConstants.allCategoryTypes.Transfer].length) {
        if (options.categoryId && options.categoryId !== '0') {
            if (isSubCategoryIdAvailable(allCategories[categoryConstants.allCategoryTypes.Transfer], options.categoryId)) {
                transaction.transferCategory = options.categoryId;
            } else {
                transaction.transferCategory = getFirstAvaiableSubCategoryId(allCategories[categoryConstants.allCategoryTypes.Transfer], options.categoryId);
            }
        }

        if (!transaction.transferCategory) {
            transaction.transferCategory = getFirstAvailableCategoryId(allCategories[categoryConstants.allCategoryTypes.Transfer]);
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

    if (transaction2) {
        if (setContextData) {
            transaction.id = transaction2.id;
        }

        transaction.type = transaction2.type;

        if (transaction.type === transactionConstants.allTransactionTypes.Expense) {
            transaction.expenseCategory = transaction2.categoryId || '';
        } else if (transaction.type === transactionConstants.allTransactionTypes.Income) {
            transaction.incomeCategory = transaction2.categoryId || '';
        } else if (transaction.type === transactionConstants.allTransactionTypes.Transfer) {
            transaction.transferCategory = transaction2.categoryId || '';
        }

        if (setContextData) {
            transaction.utcOffset = transaction2.utcOffset;
            transaction.timeZone = transaction2.timeZone;
            transaction.time = getDummyUnixTimeForLocalUsage(transaction2.time, transaction.utcOffset, getBrowserTimezoneOffsetMinutes());
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
        transaction.comment = transaction2.comment;

        if (setContextData) {
            transaction.geoLocation = transaction2.geoLocation;
        }
    }
}
