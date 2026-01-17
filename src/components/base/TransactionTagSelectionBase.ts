import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import { values } from '@/core/base.ts';
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

    const lowerCaseTagSearchContent = computed<string>(() => tagSearchContent.value.toLowerCase());

    const allTagsWithGroupHeader = computed<TransactionTagWithGroupHeader[]>(() => getTagsWithGroupHeader(tag => {
        if (!tag.hidden) {
            return true;
        }

        if (selectedTagIds.value[tag.id]) {
            return true;
        }

        if (lowerCaseTagSearchContent.value && tag.name.toLowerCase().indexOf(lowerCaseTagSearchContent.value) >= 0 && isAllFilteredTagHidden.value) {
            return true;
        }

        return false;
    }));

    const filteredTagsWithGroupHeader = computed<TransactionTagWithGroupHeader[]>(() => getTagsWithGroupHeader(tag => {
        if (lowerCaseTagSearchContent.value) {
            if (tag.name.toLowerCase().indexOf(lowerCaseTagSearchContent.value) >= 0 && (!tag.hidden || isAllFilteredTagHidden.value)) {
                return true;
            } else {
                return false;
            }
        }

        return !tag.hidden || !!selectedTagIds.value[tag.id];
    }));

    const isAllFilteredTagHidden = computed<boolean>(() => {
        const lowerCaseTagSearchContent = tagSearchContent.value.toLowerCase();
        let hiddenCount = 0;

        for (const tag of values(transactionTagsStore.allTransactionTagsMap)) {
            if (!lowerCaseTagSearchContent || tag.name.toLowerCase().indexOf(lowerCaseTagSearchContent) >= 0) {
                if (!tag.hidden) {
                    return false;
                }

                hiddenCount++;
            }
        }

        return hiddenCount > 0;
    });

    function getTagsWithGroupHeader(tagFilterFn: (tag: TransactionTag) => boolean): TransactionTagWithGroupHeader[] {
        const result: TransactionTagWithGroupHeader[] = [];
        const tagsInDefaultGroup = transactionTagsStore.allTransactionTagsByGroupMap[DEFAULT_TAG_GROUP_ID];
        let addButtonAdded = false;

        if (tagsInDefaultGroup && tagsInDefaultGroup.length > 0) {
            const visibleTags = tagsInDefaultGroup.filter(tag => tagFilterFn(tag));

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

            const visibleTags = tags.filter(tag => tagFilterFn(tag));

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
        filteredTagsWithGroupHeader
    };
}
