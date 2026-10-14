import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import { values } from '@/core/base.ts';
import { NormalizedText } from '@/core/text.ts';
import { DEFAULT_TAG_GROUP_ID } from '@/consts/tag.ts';

import { TransactionTag } from '@/models/transaction_tag.ts';

export type TransactionTagWithGroupHeader = TransactionTag | {
    type: 'subheader';
    title: string;
} | {
    type: 'addbutton';
}

export interface CommonTransactionTagSelectionProps {
    modelValue: string[];
    allowAddNewTag?: boolean;
}

export function useTransactionTagSelectionBase(props: CommonTransactionTagSelectionProps, supportAddButtonInList: boolean, useClonedModelValue?: boolean) {
    const { tt } = useI18n();

    const transactionTagsStore = useTransactionTagsStore();

    const clonedModelValue = ref<string[]>(useClonedModelValue ? Array.from(props.modelValue) : []);
    const tagSearchContent = ref<string>('');
    const pinnedTagIds = ref<Record<string, boolean>>({});

    const selectedTagIds = computed<Record<string, boolean>>(() => {
        const ret: Record<string, boolean> = {};

        if (useClonedModelValue) {
            for (const tagId of clonedModelValue.value) {
                ret[tagId] = true;
            }
        } else {
            for (const tagId of props.modelValue) {
                ret[tagId] = true;
            }
        }

        return ret;
    });

    const normalizedTagSearchContent = computed<string>(() => NormalizedText.normalizeForSearch(tagSearchContent.value));

    const allTagsWithGroupHeader = computed<TransactionTagWithGroupHeader[]>(() => getTagsWithGroupHeader(tag => {
        if (!tag.hidden) {
            return true;
        }

        if (selectedTagIds.value[tag.id]) {
            return true;
        }

        if (normalizedTagSearchContent.value && NormalizedText.normalizeForSearch(tag.name).indexOf(normalizedTagSearchContent.value) >= 0 && isAllFilteredTagHidden.value) {
            return true;
        }

        return false;
    }));

    const filteredTagsWithGroupHeader = computed<TransactionTagWithGroupHeader[]>(() => getTagsWithGroupHeader(tag => {
        if (normalizedTagSearchContent.value) {
            if (NormalizedText.normalizeForSearch(tag.name).indexOf(normalizedTagSearchContent.value) >= 0 && (!tag.hidden || isAllFilteredTagHidden.value)) {
                return true;
            } else {
                return false;
            }
        }

        return !tag.hidden || !!selectedTagIds.value[tag.id];
    }));

    const isAllFilteredTagHidden = computed<boolean>(() => {
        const lowerCaseTagSearchContent = NormalizedText.normalizeForSearch(tagSearchContent.value);
        let hiddenCount = 0;

        for (const tag of values(transactionTagsStore.allTransactionTagsMap)) {
            if (!lowerCaseTagSearchContent || NormalizedText.normalizeForSearch(tag.name).indexOf(lowerCaseTagSearchContent) >= 0) {
                if (!tag.hidden) {
                    return false;
                }

                hiddenCount++;
            }
        }

        return hiddenCount > 0;
    });

    // Remembers the tags which are selected when the selection list is opened, so they can be shown on top of the list.
    // The snapshot is not updated on every check or uncheck, so items don't jump around while the user is selecting them.
    function pinSelectedTags(): void {
        const ret: Record<string, boolean> = {};

        for (const tagId of (useClonedModelValue ? clonedModelValue.value : props.modelValue)) {
            ret[tagId] = true;
        }

        pinnedTagIds.value = ret;
    }

    function getTagsWithGroupHeader(tagFilterFn: (tag: TransactionTag) => boolean): TransactionTagWithGroupHeader[] {
        const result: TransactionTagWithGroupHeader[] = [];
        const tagsInDefaultGroup = transactionTagsStore.allTransactionTagsByGroupMap[DEFAULT_TAG_GROUP_ID];
        // keep the original order while searching, so the first matched tag is still the first one in its group
        const shouldPinTags = !normalizedTagSearchContent.value && Object.keys(pinnedTagIds.value).length > 0;
        const isTagPinned = (tag: TransactionTag): boolean => shouldPinTags && !!pinnedTagIds.value[tag.id];
        let addButtonAdded = false;

        if (shouldPinTags) {
            const pinnedTags: TransactionTag[] = [];

            if (tagsInDefaultGroup) {
                pinnedTags.push(...tagsInDefaultGroup.filter(tag => isTagPinned(tag) && tagFilterFn(tag)));
            }

            for (const tagGroup of transactionTagsStore.allTransactionTagGroups) {
                const tags = transactionTagsStore.allTransactionTagsByGroupMap[tagGroup.id];

                if (tags) {
                    pinnedTags.push(...tags.filter(tag => isTagPinned(tag) && tagFilterFn(tag)));
                }
            }

            if (pinnedTags.length > 0) {
                result.push({
                    type: 'subheader',
                    title: tt('Selected Tags')
                });

                result.push(...pinnedTags);
            }
        }

        if (tagsInDefaultGroup && tagsInDefaultGroup.length > 0) {
            const visibleTags = tagsInDefaultGroup.filter(tag => !isTagPinned(tag) && tagFilterFn(tag));

            if (visibleTags.length > 0) {
                result.push({
                    type: 'subheader',
                    title: tt('Default Group')
                });

                result.push(...visibleTags);

                if (supportAddButtonInList && props.allowAddNewTag && !addButtonAdded) {
                    result.push({
                        type: 'addbutton'
                    });
                    addButtonAdded = true;
                }
            }
        }

        for (const tagGroup of transactionTagsStore.allTransactionTagGroups) {
            const tags = transactionTagsStore.allTransactionTagsByGroupMap[tagGroup.id];

            if (!tags || tags.length < 1) {
                continue;
            }

            const visibleTags = tags.filter(tag => !isTagPinned(tag) && tagFilterFn(tag));

            if (visibleTags.length > 0) {
                result.push({
                    type: 'subheader',
                    title: tagGroup.name
                });

                result.push(...visibleTags);
            }
        }

        if (supportAddButtonInList && props.allowAddNewTag && !addButtonAdded) {
            result.push({
                type: 'addbutton'
            });
        }

        return result;
    }

    return {
        // states
        clonedModelValue,
        tagSearchContent,
        // computed states
        selectedTagIds,
        allTagsWithGroupHeader,
        filteredTagsWithGroupHeader,
        // functions
        pinSelectedTags
    };
}
