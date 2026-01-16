import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useStatisticsStore } from '@/stores/statistics.ts';

import { entries, keys, values } from '@/core/base.ts';
import { TransactionTagFilterType } from '@/core/transaction.ts';
import { DEFAULT_TAG_GROUP_ID } from '@/consts/tag.ts';

import { TransactionTagGroup } from '@/models/transaction_tag_group.ts';
import type { TransactionTag } from '@/models/transaction_tag.ts';
import { TransactionTagFilter } from '@/models/transaction.ts';

import { objectFieldToArrayItem } from '@/lib/common.ts';

export enum TransactionTagFilterState {
    Default = 0,
    Include = 1,
    Exclude = 2
}

interface TransactionGroupTagFilterTypes {
    includeType: number;
    excludeType: number;
}

function getEmptyGroupTagFilterTypesMap(allTransactionTagsByGroupMap: Record<string, TransactionTag[]>): Record<string, TransactionGroupTagFilterTypes> {
    const ret: Record<string, TransactionGroupTagFilterTypes> = {};

    for (const groupId of keys(allTransactionTagsByGroupMap)) {
        ret[groupId] = {
            includeType: TransactionTagFilterType.HasAny.type,
            excludeType: TransactionTagFilterType.NotHasAny.type
        };
    }

    return ret;
}

export function useTransactionTagFilterSettingPageBase(type?: string) {
    const { tt } = useI18n();

    const transactionTagsStore = useTransactionTagsStore();
    const transactionsStore = useTransactionsStore();
    const statisticsStore = useStatisticsStore();

    const loading = ref<boolean>(true);
    const showHidden = ref<boolean>(false);
    const filterContent = ref<string>('');

    const tagFilterStateMap = ref<Record<string, TransactionTagFilterState>>({});
    const groupTagFilterTypesMap = ref<Record<string, TransactionGroupTagFilterTypes>>(getEmptyGroupTagFilterTypesMap(transactionTagsStore.allTransactionTagsByGroupMap));

    const lowerCaseFilterContent = computed<string>(() => filterContent.value.toLowerCase());

    const title = computed<string>(() => {
        return 'Filter Transaction Tags';
    });

    const applyText = computed<string>(() => {
        return 'Apply';
    });

    const groupTagFilterStateCountMap = computed<Record<string, Record<TransactionTagFilterState, number>>>(() => {
        const ret: Record<string, Record<TransactionTagFilterState, number>> = {};

        for (const [groupId, tags] of entries(transactionTagsStore.allTransactionTagsByGroupMap)) {
            const stateCountMap: Record<TransactionTagFilterState, number> = {
                [TransactionTagFilterState.Default]: 0,
                [TransactionTagFilterState.Include]: 0,
                [TransactionTagFilterState.Exclude]: 0
            };

            for (const tag of tags) {
                const state = tagFilterStateMap.value[tag.id] ?? TransactionTagFilterState.Default;
                stateCountMap[state] = (stateCountMap[state] || 0) + 1;
            }

            ret[groupId] = stateCountMap;
        }

        return ret;
    });

    const allTagGroupsWithDefault = computed<TransactionTagGroup[]>(() => {
        const allGroups: TransactionTagGroup[] = [];
        const tagsInDefaultGroup = transactionTagsStore.allTransactionTagsByGroupMap[DEFAULT_TAG_GROUP_ID];

        if (tagsInDefaultGroup && tagsInDefaultGroup.length) {
            const defaultGroup = TransactionTagGroup.createNewTagGroup(tt('Default Group'));
            defaultGroup.id = DEFAULT_TAG_GROUP_ID;
            allGroups.push(defaultGroup);
        }

        for (const tagGroup of transactionTagsStore.allTransactionTagGroups) {
            const tagsInGroup = transactionTagsStore.allTransactionTagsByGroupMap[tagGroup.id];

            if (tagsInGroup && tagsInGroup.length) {
                allGroups.push(tagGroup);
            }
        }

        return allGroups;
    });

    const allVisibleTags = computed<Record<string, TransactionTag[]>>(() => {
        const ret: Record<string, TransactionTag[]> = {};
        const allTagGroups = transactionTagsStore.allTransactionTagsByGroupMap;

        for (const [groupId, tags] of entries(allTagGroups)) {
            const visibleTags: TransactionTag[] = [];

            for (const tag of tags) {
                if (!showHidden.value && tag.hidden) {
                    continue;
                }

                if (lowerCaseFilterContent.value && !tag.name.toLowerCase().includes(lowerCaseFilterContent.value)) {
                    continue;
                }

                visibleTags.push(tag);
            }

            if (visibleTags.length > 0) {
                ret[groupId] = visibleTags;
            }
        }

        return ret;
    });

    const allVisibleTagGroupIds = computed<string[]>(() => objectFieldToArrayItem(allVisibleTags.value));
    const hasAnyAvailableTag = computed<boolean>(() => transactionTagsStore.allAvailableTagsCount > 0);
    const hasAnyVisibleTag = computed<boolean>(() => {
        for (const tags of values(allVisibleTags.value)) {
            if (tags.length > 0) {
                return true;
            }
        }
        return false;
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
        const allGroupTagFilterTypesMap: Record<string, TransactionGroupTagFilterTypes> = getEmptyGroupTagFilterTypesMap(transactionTagsStore.allTransactionTagsByGroupMap);

        for (const transactionTag of values(transactionTagsStore.allTransactionTagsMap)) {
            allTagIdsMap[transactionTag.id] = TransactionTagFilterState.Default;
        }

        for (const tagFilter of tagFilters) {
            let state: TransactionTagFilterState = TransactionTagFilterState.Default;

            if (tagFilter.type === TransactionTagFilterType.HasAny || tagFilter.type === TransactionTagFilterType.HasAll) {
                state = TransactionTagFilterState.Include;
            } else if (tagFilter.type === TransactionTagFilterType.NotHasAny || tagFilter.type === TransactionTagFilterType.NotHasAll) {
                state = TransactionTagFilterState.Exclude;
            } else {
                continue;
            }

            for (const tagId of tagFilter.tagIds) {
                const tag = transactionTagsStore.allTransactionTagsMap[tagId];

                if (!tag) {
                    continue;
                }

                const groupFilterTypes = allGroupTagFilterTypesMap[tag.groupId];

                if (groupFilterTypes) {
                    if (state === TransactionTagFilterState.Include) {
                        groupFilterTypes.includeType = tagFilter.type.type;
                    } else if (state === TransactionTagFilterState.Exclude) {
                        groupFilterTypes.excludeType = tagFilter.type.type;
                    }

                    allTagIdsMap[tagId] = state;
                }
            }
        }

        tagFilterStateMap.value = allTagIdsMap;
        groupTagFilterTypesMap.value = allGroupTagFilterTypesMap;
        return true;
    }

    function saveFilterTagIds(): boolean {
        const tagFilters: TransactionTagFilter[] = [];
        let changed = true;

        for (const [groupId, tags] of entries(transactionTagsStore.allTransactionTagsByGroupMap)) {
            const groupFilterTypes = groupTagFilterTypesMap.value[groupId];

            if (groupFilterTypes && tags && tags.length > 0) {
                const includeTagIds: string[] = [];
                const excludeTagIds: string[] = [];

                for (const tag of tags) {
                    const state = tagFilterStateMap.value[tag.id] ?? TransactionTagFilterState.Default;

                    if (state === TransactionTagFilterState.Include) {
                        includeTagIds.push(tag.id);
                    } else if (state === TransactionTagFilterState.Exclude) {
                        excludeTagIds.push(tag.id);
                    }
                }

                if (includeTagIds.length > 0) {
                    const includeTagFilter = TransactionTagFilter.create(includeTagIds, TransactionTagFilterType.parse(groupFilterTypes.includeType) ?? TransactionTagFilterType.HasAny);
                    tagFilters.push(includeTagFilter);
                }

                if (excludeTagIds.length > 0) {
                    const excludeTagFilter = TransactionTagFilter.create(excludeTagIds, TransactionTagFilterType.parse(groupFilterTypes.excludeType) ?? TransactionTagFilterType.NotHasAny);
                    tagFilters.push(excludeTagFilter);
                }
            }
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
        filterContent,
        tagFilterStateMap,
        groupTagFilterTypesMap,
        // computed states
        title,
        applyText,
        groupTagFilterStateCountMap,
        allTagGroupsWithDefault,
        allVisibleTags,
        allVisibleTagGroupIds,
        hasAnyAvailableTag,
        hasAnyVisibleTag,
        // functions
        loadFilterTagIds,
        saveFilterTagIds
    };
}
