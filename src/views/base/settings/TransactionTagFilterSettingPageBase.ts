import { ref, computed } from 'vue';

import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useStatisticsStore } from '@/stores/statistics.ts';

import { entries, values } from '@/core/base.ts';
import { TransactionTagFilterType } from '@/core/transaction.ts';
import type { TransactionTag } from '@/models/transaction_tag.ts';
import { TransactionTagFilter } from '@/models/transaction.ts';

import { objectFieldWithValueToArrayItem } from '@/lib/common.ts';

export enum TransactionTagFilterState {
    Default = 0,
    Include = 1,
    Exclude = 2
}

export function useTransactionTagFilterSettingPageBase(type?: string) {
    const transactionTagsStore = useTransactionTagsStore();
    const transactionsStore = useTransactionsStore();
    const statisticsStore = useStatisticsStore();

    const loading = ref<boolean>(true);
    const showHidden = ref<boolean>(false);
    const filterTagIds = ref<Record<string, TransactionTagFilterState>>({});
    const includeTagFilterType = ref<number>(TransactionTagFilterType.HasAny.type);
    const excludeTagFilterType = ref<number>(TransactionTagFilterType.NotHasAny.type);

    const includeTagsCount = computed<number>(() => objectFieldWithValueToArrayItem(filterTagIds.value, TransactionTagFilterState.Include).length);
    const excludeTagsCount = computed<number>(() => objectFieldWithValueToArrayItem(filterTagIds.value, TransactionTagFilterState.Exclude).length);

    const title = computed<string>(() => {
        return 'Filter Transaction Tags';
    });

    const applyText = computed<string>(() => {
        return 'Apply';
    });

    const allTags = computed<TransactionTag[]>(() => transactionTagsStore.allTransactionTags);
    const hasAnyAvailableTag = computed<boolean>(() => transactionTagsStore.allAvailableTagsCount > 0);
    const hasAnyVisibleTag = computed<boolean>(() => {
        if (showHidden.value) {
            return transactionTagsStore.allAvailableTagsCount > 0;
        } else {
            return transactionTagsStore.allVisibleTagsCount > 0;
        }
    });

    function loadFilterTagIds(): boolean {
        let tagFilters: TransactionTagFilter[] = [];

        if (type === 'statisticsCurrent') {
            tagFilters = TransactionTagFilter.parse(statisticsStore.transactionStatisticsFilter.tagFilter);
        } else if (type === 'transactionListCurrent') {
            tagFilters = TransactionTagFilter.parse(transactionsStore.transactionsFilter.tagFilter);
        } else {
            return false;
        }

        const allTagIdsMap: Record<string, TransactionTagFilterState> = {};

        for (const transactionTag of values(transactionTagsStore.allTransactionTagsMap)) {
            allTagIdsMap[transactionTag.id] = TransactionTagFilterState.Default;
        }

        for (const tagFilter of tagFilters) {
            let state: TransactionTagFilterState = TransactionTagFilterState.Default;

            if (tagFilter.type === TransactionTagFilterType.HasAny || tagFilter.type === TransactionTagFilterType.HasAll) {
                state = TransactionTagFilterState.Include;
                includeTagFilterType.value = tagFilter.type.type;
            } else if (tagFilter.type === TransactionTagFilterType.NotHasAny || tagFilter.type === TransactionTagFilterType.NotHasAll) {
                state = TransactionTagFilterState.Exclude;
                excludeTagFilterType.value = tagFilter.type.type;
            } else {
                continue;
            }

            for (const tagId of tagFilter.tagIds) {
                allTagIdsMap[tagId] = state;
            }
        }

        filterTagIds.value = allTagIdsMap;
        return true;
    }

    function saveFilterTagIds(): boolean {
        const includeTagFilter: TransactionTagFilter = TransactionTagFilter.create(TransactionTagFilterType.parse(includeTagFilterType.value) ?? TransactionTagFilterType.HasAny);
        const excludeTagFilter: TransactionTagFilter = TransactionTagFilter.create(TransactionTagFilterType.parse(excludeTagFilterType.value) ?? TransactionTagFilterType.NotHasAny);
        let changed = true;

        for (const [transactionTagId, state] of entries(filterTagIds.value)) {
            const transactionTag = transactionTagsStore.allTransactionTagsMap[transactionTagId];

            if (!transactionTag) {
                continue;
            }

            if (state === TransactionTagFilterState.Include) {
                includeTagFilter.tagIds.push(transactionTag.id);
            } else if (state === TransactionTagFilterState.Exclude) {
                excludeTagFilter.tagIds.push(transactionTag.id);
            }
        }

        const tagFilters: TransactionTagFilter[] = [];

        if (includeTagFilter.tagIds.length > 0) {
            tagFilters.push(includeTagFilter);
        }

        if (excludeTagFilter.tagIds.length > 0) {
            tagFilters.push(excludeTagFilter);
        }

        if (type === 'statisticsCurrent') {
            changed = statisticsStore.updateTransactionStatisticsFilter({
                tagFilter: TransactionTagFilter.toTextualTagFilters(tagFilters)
            });

            if (changed) {
                statisticsStore.updateTransactionStatisticsInvalidState(true);
            }
        } else if (type === 'transactionListCurrent') {
            changed = transactionsStore.updateTransactionListFilter({
                tagFilter: TransactionTagFilter.toTextualTagFilters(tagFilters)
            });

            if (changed) {
                transactionsStore.updateTransactionListInvalidState(true);
            }
        }

        return changed;
    }

    return {
        // states
        loading,
        showHidden,
        filterTagIds,
        includeTagFilterType,
        excludeTagFilterType,
        // computed states
        includeTagsCount,
        excludeTagsCount,
        title,
        applyText,
        allTags,
        hasAnyAvailableTag,
        hasAnyVisibleTag,
        // functions
        loadFilterTagIds,
        saveFilterTagIds
    };
}
