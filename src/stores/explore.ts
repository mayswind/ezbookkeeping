import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import { useSettingsStore } from './setting.ts';
import { useUserStore } from './user.ts';
import { useAccountsStore } from './account.ts';
import { useTransactionCategoriesStore } from './transactionCategory.ts';
import { useTransactionTagsStore } from './transactionTag.ts';

import { DateRangeScene, DateRange } from '@/core/datetime.ts';
import { DEFAULT_TRANSACTION_EXPLORE_DATE_RANGE } from '@/core/explore.ts';

import { type Account } from '@/models/account.ts';
import { type TransactionCategory } from '@/models/transaction_category.ts';
import { type TransactionTag } from '@/models/transaction_tag.ts';
import {
    type TransactionInfoResponse,
    type TransactionInsightDataItem
} from '@/models/transaction.ts';
import {
    TransactionExploreQuery
} from '@/models/explore.ts';

import { isInteger, isEquals } from '@/lib/common.ts';
import { getDateRangeByDateType } from '@/lib/datetime.ts';
import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

export interface TransactionExplorePartialFilter {
    dateRangeType?: number;
    startTime?: number;
    endTime?: number;
    queryId?: string;
}

export interface TransactionExploreFilter extends TransactionExplorePartialFilter {
    dateRangeType: number;
    startTime: number;
    endTime: number;
    query: TransactionExploreQuery[];
}

export const useExploresStore = defineStore('explores', () => {
    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const accountsStore = useAccountsStore();
    const transactionCategoriesStore = useTransactionCategoriesStore();
    const transactionTagsStore = useTransactionTagsStore();

    const transactionExploreFilter = ref<TransactionExploreFilter>({
        dateRangeType: DEFAULT_TRANSACTION_EXPLORE_DATE_RANGE.type,
        startTime: 0,
        endTime: 0,
        query: []
    });

    const transactionExploreAllData = ref<TransactionInfoResponse[]>([]);
    const transactionExploreStateInvalid = ref<boolean>(true);

    const allTransactions = computed<TransactionInsightDataItem[]>(() => {
        if (!transactionExploreAllData.value || transactionExploreAllData.value.length < 1) {
            return [];
        }

        const result: TransactionInsightDataItem[] = [];

        for (const transaction of transactionExploreAllData.value) {
            const sourceAccount: Account | undefined = accountsStore.allAccountsMap[transaction.sourceAccountId];

            if (!sourceAccount) {
                continue;
            }

            let destinationAccount: Account | undefined = undefined

            if (transaction.destinationAccountId && transaction.destinationAccountId !== '0') {
                destinationAccount = accountsStore.allAccountsMap[transaction.destinationAccountId];

                if (!destinationAccount) {
                    continue;
                }
            }

            const secondaryCategory: TransactionCategory | undefined = transactionCategoriesStore.allTransactionCategoriesMap[transaction.categoryId];

            if (!secondaryCategory) {
                continue;
            }

            const primaryCategory: TransactionCategory | undefined = transactionCategoriesStore.allTransactionCategoriesMap[secondaryCategory.parentId];

            if (!primaryCategory) {
                continue;
            }

            const tags: TransactionTag[] = [];

            for (const tagId of transaction.tagIds) {
                const tag: TransactionTag | undefined = transactionTagsStore.allTransactionTagsMap[tagId];

                if (tag) {
                    tags.push(tag);
                }
            }

            const item: TransactionInsightDataItem = {
                ...transaction,
                id: transaction.id,
                time: transaction.time,
                utcOffset: transaction.utcOffset,
                type: transaction.type,
                primaryCategory: primaryCategory,
                primaryCategoryName: primaryCategory.name,
                secondaryCategory: secondaryCategory,
                secondaryCategoryName: secondaryCategory.name,
                sourceAccount: sourceAccount,
                sourceAccountName: sourceAccount.name,
                destinationAccount: destinationAccount,
                destinationAccountName: destinationAccount?.name,
                sourceAmount: transaction.sourceAmount,
                destinationAmount: transaction.destinationAmount,
                hideAmount: transaction.hideAmount,
                tags: tags,
                comment: transaction.comment,
                geoLocation: transaction.geoLocation
            };

            result.push(item);
        }

        return result;
    });

    const filteredTransactions = computed<TransactionInsightDataItem[]>(() => {
        if (!allTransactions.value || allTransactions.value.length < 1) {
            return [];
        }

        if (!transactionExploreFilter.value.query || transactionExploreFilter.value.query.length < 1) {
            return allTransactions.value;
        }

        const result: TransactionInsightDataItem[] = [];

        for (const transaction of allTransactions.value) {
            for (const query of transactionExploreFilter.value.query) {
                if (query.match(transaction)) {
                    result.push(transaction);
                    break;
                }
            }
        }

        return result;
    });

    function updateTransactionExploreInvalidState(invalidState: boolean): void {
        transactionExploreStateInvalid.value = invalidState;
    }

    function resetTransactionExplores(): void {
        transactionExploreFilter.value.dateRangeType = DEFAULT_TRANSACTION_EXPLORE_DATE_RANGE.type;
        transactionExploreFilter.value.startTime = 0;
        transactionExploreFilter.value.endTime = 0;
        transactionExploreFilter.value.query = [];
        transactionExploreAllData.value = [];
        transactionExploreStateInvalid.value = true;
    }

    function initTransactionExploreFilter(filter?: TransactionExplorePartialFilter, resetQuery?: boolean): void {
        if (filter && isInteger(filter.dateRangeType)) {
            transactionExploreFilter.value.dateRangeType = filter.dateRangeType;
        } else {
            transactionExploreFilter.value.dateRangeType = settingsStore.appSettings.insightsExploreDefaultDateRangeType;
        }

        let dateRangeTypeValid = true;

        if (!DateRange.isAvailableForScene(transactionExploreFilter.value.dateRangeType, DateRangeScene.InsightsExplore)) {
            transactionExploreFilter.value.dateRangeType = DEFAULT_TRANSACTION_EXPLORE_DATE_RANGE.type;
            dateRangeTypeValid = false;
        }

        if (dateRangeTypeValid && transactionExploreFilter.value.dateRangeType === DateRange.Custom.type) {
            if (filter && isInteger(filter.startTime)) {
                transactionExploreFilter.value.startTime = filter.startTime;
            } else {
                transactionExploreFilter.value.startTime = 0;
            }

            if (filter && isInteger(filter.endTime)) {
                transactionExploreFilter.value.endTime = filter.endTime;
            } else {
                transactionExploreFilter.value.endTime = 0;
            }
        } else {
            const dateRange = getDateRangeByDateType(transactionExploreFilter.value.dateRangeType, userStore.currentUserFirstDayOfWeek, userStore.currentUserFiscalYearStart);

            if (dateRange) {
                transactionExploreFilter.value.dateRangeType = dateRange.dateType;
                transactionExploreFilter.value.startTime = dateRange.minTime;
                transactionExploreFilter.value.endTime = dateRange.maxTime;
            }
        }

        if (resetQuery) {
            transactionExploreFilter.value.query = [];
        }
    }

    function updateTransactionExploreFilter(filter: TransactionExplorePartialFilter): boolean {
        let changed = false;

        if (filter && isInteger(filter.dateRangeType) && transactionExploreFilter.value.dateRangeType !== filter.dateRangeType) {
            transactionExploreFilter.value.dateRangeType = filter.dateRangeType;
            changed = true;
        }

        if (filter && isInteger(filter.startTime) && transactionExploreFilter.value.startTime !== filter.startTime) {
            transactionExploreFilter.value.startTime = filter.startTime;
            changed = true;
        }

        if (filter && isInteger(filter.endTime) && transactionExploreFilter.value.endTime !== filter.endTime) {
            transactionExploreFilter.value.endTime = filter.endTime;
            changed = true;
        }

        return changed;
    }

    function getTransactionExplorePageParams(currentExploreId: string, activeTab: string): string {
        const querys: string[] = [];

        if (currentExploreId) {
            querys.push('id=' + currentExploreId);
        }

        if (activeTab) {
            querys.push('activeTab=' + activeTab);
        }

        querys.push('dateRangeType=' + transactionExploreFilter.value.dateRangeType);
        querys.push('startTime=' + transactionExploreFilter.value.startTime);
        querys.push('endTime=' + transactionExploreFilter.value.endTime);

        return querys.join('&');
    }

    function loadAllTransactions({ force }: { force: boolean }): Promise<TransactionInfoResponse[]> {
        return new Promise((resolve, reject) => {
            services.getAllTransactions({
                startTime: transactionExploreFilter.value.startTime,
                endTime: transactionExploreFilter.value.endTime
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve all transactions' });
                    return;
                }

                if (transactionExploreStateInvalid.value) {
                    updateTransactionExploreInvalidState(false);
                }

                if (force && data.result && isEquals(transactionExploreAllData.value, data.result)) {
                    reject({ message: 'Data is up to date', isUpToDate: true });
                    return;
                }

                transactionExploreAllData.value = data.result;

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to load all transactions', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve all transactions' });
                } else {
                    reject(error);
                }
            });
        });
    }

    return {
        // states
        transactionExploreFilter,
        transactionExploreStateInvalid,
        // computed
        filteredTransactions,
        // functions
        updateTransactionExploreInvalidState,
        resetTransactionExplores,
        initTransactionExploreFilter,
        updateTransactionExploreFilter,
        getTransactionExplorePageParams,
        loadAllTransactions
    };
});
