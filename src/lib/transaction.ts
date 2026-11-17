import { CategoryType } from '@/core/category.ts';
import { TransactionType } from '@/core/transaction.ts';
import { Account } from '@/models/account.ts';
import { TransactionCategory } from '@/models/transaction_category.ts';
import { TransactionTag } from '@/models/transaction_tag.ts';
import { TransactionPicture, type TransactionPictureInfoBasicResponse } from '@/models/transaction_picture_info.ts';
import { type TransactionInfoResponse, Transaction } from '@/models/transaction.ts';

import {
    isDefined,
    isNumber
} from './common.ts';
import {
    getTimezoneOffsetMinutes
} from './datetime.ts';
import {
    categoryTypeToTransactionType,
    isSubCategoryIdAvailable,
    getFirstVisibleCategoryId,
    getFirstAvailableSubCategoryId
} from './category.ts';

export interface SetTransactionOptions {
    time?: number;
    type?: number;
    categoryId?: string;
    accountId?: string;
    destinationAccountId?: string;
    amount?: number;
    destinationAmount?: number;
    tagIds?: string;
    comment?: string;
}

export function* allTransactionPictures(transactions: Transaction[]): Iterable<[Transaction, TransactionPictureInfoBasicResponse]> {
    for (const transaction of transactions) {
        if (transaction.pictures) {
            for (const pictureInfo of transaction.pictures) {
                yield [transaction, pictureInfo];
            }
        }
    }
}

export function setTransactionModelByTransaction(transaction: Transaction, transaction2: Transaction | null | undefined, allCategories: Record<number, TransactionCategory[]>, allCategoriesMap: Record<string, TransactionCategory>, allVisibleAccounts: Account[], allAccountsMap: Record<string, Account>, allTagsMap: Record<string, TransactionTag>, defaultAccountId: string, options: SetTransactionOptions, setContextData: boolean): void {
    if (isDefined(options.time)) {
        transaction.time = options.time;
        transaction.utcOffset = getTimezoneOffsetMinutes(transaction.time, transaction.timeZone);
    }

    if (options.type === TransactionType.Income || options.type === TransactionType.Expense || options.type === TransactionType.Transfer) {
        transaction.type = options.type;
    } else if (!options.type && options.categoryId && options.categoryId !== '0' && allCategoriesMap[options.categoryId]) {
        const category = allCategoriesMap[options.categoryId] as TransactionCategory;
        const type = categoryTypeToTransactionType(category.type);

        if (isNumber(type)) {
            transaction.type = type;
        }
    }

    if (isDefined(options.amount)) {
        transaction.sourceAmount = options.amount;
    }

    if (isDefined(options.destinationAmount)) {
        transaction.destinationAmount = options.destinationAmount;
    }

    if (allCategories[CategoryType.Expense] &&
        allCategories[CategoryType.Expense].length) {
        if (options.categoryId && options.categoryId !== '0') {
            if (isSubCategoryIdAvailable(allCategories[CategoryType.Expense], options.categoryId)) {
                transaction.expenseCategoryId = options.categoryId;
            } else {
                transaction.expenseCategoryId = getFirstAvailableSubCategoryId(allCategories[CategoryType.Expense], options.categoryId);
            }
        }

        if (!transaction.expenseCategoryId) {
            transaction.expenseCategoryId = getFirstVisibleCategoryId(allCategories[CategoryType.Expense]);
        }
    }

    if (allCategories[CategoryType.Income] &&
        allCategories[CategoryType.Income].length) {
        if (options.categoryId && options.categoryId !== '0') {
            if (isSubCategoryIdAvailable(allCategories[CategoryType.Income], options.categoryId)) {
                transaction.incomeCategoryId = options.categoryId;
            } else {
                transaction.incomeCategoryId = getFirstAvailableSubCategoryId(allCategories[CategoryType.Income], options.categoryId);
            }
        }

        if (!transaction.incomeCategoryId) {
            transaction.incomeCategoryId = getFirstVisibleCategoryId(allCategories[CategoryType.Income]);
        }
    }

    if (allCategories[CategoryType.Transfer] &&
        allCategories[CategoryType.Transfer].length) {
        if (options.categoryId && options.categoryId !== '0') {
            if (isSubCategoryIdAvailable(allCategories[CategoryType.Transfer], options.categoryId)) {
                transaction.transferCategoryId = options.categoryId;
            } else {
                transaction.transferCategoryId = getFirstAvailableSubCategoryId(allCategories[CategoryType.Transfer], options.categoryId);
            }
        }

        if (!transaction.transferCategoryId) {
            transaction.transferCategoryId = getFirstVisibleCategoryId(allCategories[CategoryType.Transfer]);
        }
    }

    if (allVisibleAccounts.length) {
        if (options.accountId && options.accountId !== '0') {
            for (const account of allVisibleAccounts) {
                if (account.id === options.accountId) {
                    transaction.sourceAccountId = options.accountId;
                    transaction.destinationAccountId = options.accountId;
                    break;
                }
            }
        }

        if (options.destinationAccountId && options.destinationAccountId !== '0') {
            for (const account of allVisibleAccounts) {
                if (account.id === options.destinationAccountId) {
                    transaction.destinationAccountId = options.destinationAccountId;
                    break;
                }
            }
        }

        if (!transaction.sourceAccountId) {
            if (defaultAccountId && allAccountsMap[defaultAccountId] && !allAccountsMap[defaultAccountId].hidden) {
                transaction.sourceAccountId = defaultAccountId;
            } else {
                transaction.sourceAccountId = allVisibleAccounts[0]!.id;
            }
        }

        if (!transaction.destinationAccountId) {
            if (defaultAccountId && allAccountsMap[defaultAccountId] && !allAccountsMap[defaultAccountId].hidden) {
                transaction.destinationAccountId = defaultAccountId;
            } else {
                transaction.destinationAccountId = allVisibleAccounts[0]!.id;
            }
        }
    }

    if (allTagsMap && options.tagIds) {
        const tagIds = options.tagIds.split(',');
        const finalTagIds = [];

        for (const tagId of tagIds) {
            const tag = allTagsMap[tagId];

            if (tag && !tag.hidden) {
                finalTagIds.push(tag.id);
            }
        }

        transaction.tagIds = finalTagIds;
    }

    if (options.comment) {
        transaction.comment = options.comment;
    }

    if (transaction2) {
        if (setContextData) {
            transaction.id = transaction2.id;
        }

        transaction.type = transaction2.type;

        if (transaction.type === TransactionType.Expense) {
            transaction.expenseCategoryId = transaction2.categoryId || '';
        } else if (transaction.type === TransactionType.Income) {
            transaction.incomeCategoryId = transaction2.categoryId || '';
        } else if (transaction.type === TransactionType.Transfer) {
            transaction.transferCategoryId = transaction2.categoryId || '';
        }

        if (setContextData) {
            transaction.time = transaction2.time;
            transaction.timeZone = transaction2.timeZone;
            transaction.utcOffset = transaction2.utcOffset;
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
        transaction.tagIds = transaction2.tagIds ? [...transaction2.tagIds] : [];
        transaction.setPictures(TransactionPicture.ofMulti(transaction2.pictures || []));

        transaction.comment = transaction2.comment;

        if (setContextData) {
            transaction.setGeoLocation(transaction2.geoLocation);
        }
    }
}

export function isTransactionEquals(transaction1: TransactionInfoResponse, transaction2: TransactionInfoResponse): boolean {
    if (transaction1.id !== transaction2.id ||
        transaction1.timeSequenceId !== transaction2.timeSequenceId ||
        transaction1.type !== transaction2.type ||
        transaction1.categoryId !== transaction2.categoryId ||
        transaction1.time !== transaction2.time ||
        transaction1.utcOffset !== transaction2.utcOffset ||
        transaction1.sourceAccountId !== transaction2.sourceAccountId ||
        transaction1.destinationAccountId !== transaction2.destinationAccountId ||
        transaction1.sourceAmount !== transaction2.sourceAmount ||
        transaction1.destinationAmount !== transaction2.destinationAmount ||
        transaction1.hideAmount !== transaction2.hideAmount ||
        transaction1.comment !== transaction2.comment ||
        transaction1.editable !== transaction2.editable) {
        return false;
    }

    if (transaction1.tagIds.length !== transaction2.tagIds.length) {
        return false;
    }

    for (let i = 0; i < transaction1.tagIds.length; i++) {
        if (transaction1.tagIds[i] !== transaction2.tagIds[i]) {
            return false;
        }
    }

    if (transaction1.category && transaction2.category) {
        if (transaction1.category.id !== transaction2.category.id) {
            return false;
        }
    } else if (transaction1.category || transaction2.category) {
        return false;
    }

    const accounts1 = [transaction1.sourceAccount, transaction1.destinationAccount];
    const accounts2 = [transaction2.sourceAccount, transaction2.destinationAccount];

    for (let i = 0; i < accounts1.length; i++) {
        if (accounts1[i] && accounts2[i]) {
            if (accounts1[i]!.id !== accounts2[i]!.id) {
                return false;
            }
        } else if (accounts1[i] || accounts2[i]) {
            return false;
        }
    }

    const tags1 = transaction1.tags || [];
    const tags2 = transaction2.tags || [];

    if (tags1.length !== tags2.length) {
        return false;
    }

    for (let i = 0; i < tags1.length; i++) {
        const tag1 = tags1[i]!;
        const tag2 = tags2[i]!;

        if (tag1.id !== tag2.id) {
            return false;
        }
    }

    const pictures1 = transaction1.pictures || [];
    const pictures2 = transaction2.pictures || [];

    if (pictures1.length !== pictures2.length) {
        return false;
    }

    for (let i = 0; i < pictures1.length; i++) {
        const picture1 = pictures1[i]!;
        const picture2 = pictures2[i]!;

        if (picture1.pictureId !== picture2.pictureId) {
            return false;
        }
    }

    if (transaction1.geoLocation && transaction2.geoLocation) {
        if (transaction1.geoLocation.latitude !== transaction2.geoLocation.latitude || transaction1.geoLocation.longitude !== transaction2.geoLocation.longitude) {
            return false;
        }
    } else if (transaction1.geoLocation || transaction2.geoLocation) {
        return false;
    }

    return true;
}

export function isTransactionsEquals(transactions1: TransactionInfoResponse[], transactions2: TransactionInfoResponse[]): boolean {
    if (transactions1.length !== transactions2.length) {
        return false;
    }

    for (let i = 0; i < transactions1.length; i++) {
        if (!isTransactionEquals(transactions1[i]!, transactions2[i]!)) {
            return false;
        }
    }

    return true;
}
