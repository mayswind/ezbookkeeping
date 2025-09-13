import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useStatisticsStore } from '@/stores/statistics.ts';

import { type TypeAndDisplayName, keys, keysIfValueEquals, values } from '@/core/base.ts';
import { TransactionTagFilterType } from '@/core/transaction.ts';
import type { TransactionTag } from '@/models/transaction_tag.ts';

export function useTransactionTagFilterSettingPageBase(type?: string) {
    const { getAllTransactionTagFilterTypes } = useI18n();

    const transactionTagsStore = useTransactionTagsStore();
    const transactionsStore = useTransactionsStore();
    const statisticsStore = useStatisticsStore();

    const loading = ref<boolean>(true);
    const showHidden = ref<boolean>(false);
    const filterTagIds = ref<Record<string, boolean>>({});
    const tagFilterType = ref<number>(TransactionTagFilterType.Default.type);

    const title = computed<string>(() => {
        return 'Filter Transaction Tags';
    });

    const applyText = computed<string>(() => {
        return 'Apply';
    });

    const allTags = computed<TransactionTag[]>(() => transactionTagsStore.allTransactionTags);
    const allTagFilterTypes = computed<TypeAndDisplayName[]>(() => getAllTransactionTagFilterTypes());
    const hasAnyAvailableTag = computed<boolean>(() => transactionTagsStore.allAvailableTagsCount > 0);
    const hasAnyVisibleTag = computed<boolean>(() => {
        if (showHidden.value) {
            return transactionTagsStore.allAvailableTagsCount > 0;
        } else {
            return transactionTagsStore.allVisibleTagsCount > 0;
        }
    });

    function loadFilterTagIds(): boolean {
        const allTransactionTagIds: Record<string, boolean> = {};

        for (const transactionTag of values(transactionTagsStore.allTransactionTagsMap)) {
            allTransactionTagIds[transactionTag.id] = true;
        }

        if (type === 'statisticsCurrent') {
            const transactionTagIds = statisticsStore.transactionStatisticsFilter.tagIds ? statisticsStore.transactionStatisticsFilter.tagIds.split(',') : [];

            for (const transactionTagId of transactionTagIds) {
                const transactionTag = transactionTagsStore.allTransactionTagsMap[transactionTagId];

                if (transactionTag) {
                    allTransactionTagIds[transactionTag.id] = false;
                }
            }
            filterTagIds.value = allTransactionTagIds;
            tagFilterType.value = statisticsStore.transactionStatisticsFilter.tagFilterType;
            return true;
        } else if (type === 'transactionListCurrent') {
            for (const transactionTagId of keysIfValueEquals(transactionsStore.allFilterTagIds, true)) {
                const transactionTag = transactionTagsStore.allTransactionTagsMap[transactionTagId];

                if (transactionTag) {
                    allTransactionTagIds[transactionTag.id] = false;
                }
            }
            filterTagIds.value = allTransactionTagIds;
            return true;
        } else {
            return false;
        }
    }

    function saveFilterTagIds(): boolean {
        const filteredTagIds: Record<string, boolean> = {};
        let finalTagIds = '';
        let changed = true;

        for (const transactionTagId of keys(filterTagIds.value)) {
            const transactionTag = transactionTagsStore.allTransactionTagsMap[transactionTagId];

            if (!transactionTag) {
                continue;
            }

            if (filterTagIds.value[transactionTag.id]) {
                filteredTagIds[transactionTag.id] = true;
            } else {
                if (finalTagIds.length > 0) {
                    finalTagIds += ',';
                }

                finalTagIds += transactionTag.id;
            }
        }

        if (type === 'statisticsCurrent') {
            changed = statisticsStore.updateTransactionStatisticsFilter({
                tagIds: finalTagIds,
                tagFilterType: tagFilterType.value
            });

            if (changed) {
                statisticsStore.updateTransactionStatisticsInvalidState(true);
            }
        } else if (type === 'transactionListCurrent') {
            changed = transactionsStore.updateTransactionListFilter({
                tagIds: finalTagIds
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
        tagFilterType,
        // computed states
        title,
        applyText,
        allTags,
        allTagFilterTypes,
        hasAnyAvailableTag,
        hasAnyVisibleTag,
        // functions
        loadFilterTagIds,
        saveFilterTagIds
    };
}
