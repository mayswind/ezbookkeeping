import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import { DEFAULT_TAG_GROUP_ID } from '@/consts/tag.ts';

import { TransactionTagGroup } from '@/models/transaction_tag_group.ts';
import { TransactionTag } from '@/models/transaction_tag.ts';

import { isNoAvailableTag } from '@/lib/tag.ts';

export function useTagListPageBase() {
    const { tt } = useI18n();

    const transactionTagsStore = useTransactionTagsStore();

    const activeTagGroupId = ref<string>(DEFAULT_TAG_GROUP_ID);
    const newTag = ref<TransactionTag | null>(null);
    const editingTag = ref<TransactionTag>(TransactionTag.createNewTag());
    const loading = ref<boolean>(true);
    const showHidden = ref<boolean>(false);
    const displayOrderModified = ref<boolean>(false);

    const allTagGroupsWithDefault = computed<TransactionTagGroup[]>(() => {
        const allGroups: TransactionTagGroup[] = [];
        const defaultGroup = TransactionTagGroup.createNewTagGroup(tt('Default Group'));
        defaultGroup.id = DEFAULT_TAG_GROUP_ID;
        allGroups.push(defaultGroup);
        allGroups.push(...transactionTagsStore.allTransactionTagGroups);
        return allGroups;
    });
    const tags = computed<TransactionTag[]>(() => transactionTagsStore.allTransactionTagsByGroupMap[activeTagGroupId.value] || []);

    const noAvailableTag = computed<boolean>(() => isNoAvailableTag(tags.value, showHidden.value));
    const hasEditingTag = computed<boolean>(() => !!(newTag.value || (editingTag.value.id && editingTag.value.id !== '')));

    function isTagModified(tag: TransactionTag): boolean {
        if (tag.id) {
            return editingTag.value.name !== '' && editingTag.value.name !== tag.name;
        } else {
            return tag.name !== '';
        }
    }

    function switchTagGroup(tagGroupId: string): void {
        activeTagGroupId.value = tagGroupId;

        if (newTag.value) {
            newTag.value.groupId = tagGroupId;
        }
    }

    function add(): void {
        newTag.value = TransactionTag.createNewTag('', activeTagGroupId.value);
    }

    function edit(tag: TransactionTag): void {
        editingTag.value.id = tag.id;
        editingTag.value.groupId = tag.groupId;
        editingTag.value.name = tag.name;
    }

    return {
        // states
        activeTagGroupId,
        newTag,
        editingTag,
        loading,
        showHidden,
        displayOrderModified,
        // computed states
        allTagGroupsWithDefault,
        tags,
        noAvailableTag,
        hasEditingTag,
        // functions
        isTagModified,
        switchTagGroup,
        add,
        edit
    };
}
