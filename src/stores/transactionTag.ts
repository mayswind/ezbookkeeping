import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import { type BeforeResolveFunction, itemAndIndex, values } from '@/core/base.ts';

import {
    type TransactionTagGroupInfoResponse,
    type TransactionTagGroupNewDisplayOrderRequest,
    TransactionTagGroup
} from '@/models/transaction_tag_group.ts';

import {
    type TransactionTagCreateBatchRequest,
    type TransactionTagInfoResponse,
    type TransactionTagNewDisplayOrderRequest,
    TransactionTag
} from '@/models/transaction_tag.ts';

import { isEquals } from '@/lib/common.ts';

import logger from '@/lib/logger.ts';
import services, { type ApiResponsePromise } from '@/lib/services.ts';

export const useTransactionTagsStore = defineStore('transactionTags', () => {
    const allTransactionTagGroups = ref<TransactionTagGroup[]>([]);
    const allTransactionTagGroupsMap = ref<Record<string, TransactionTagGroup>>({});
    const allTransactionTags = ref<TransactionTag[]>([]);
    const allTransactionTagsMap = ref<Record<string, TransactionTag>>({});
    const allTransactionTagsByGroupMap = ref<Record<string, TransactionTag[]>>({});
    const transactionTagGroupListStateInvalid = ref<boolean>(true);
    const transactionTagListStateInvalid = ref<boolean>(true);

    const allAvailableTagsCount = computed<number>(() => allTransactionTags.value.length);

    function loadTransactionTagGroupList(tagGroups: TransactionTagGroup[]): void {
        allTransactionTagGroups.value = tagGroups;
        allTransactionTagGroupsMap.value = {};

        for (const tagGroup of tagGroups) {
            allTransactionTagGroupsMap.value[tagGroup.id] = tagGroup;
        }
    }

    function loadTransactionTagList(tags: TransactionTag[]): void {
        allTransactionTags.value = tags;
        allTransactionTagsMap.value = {};
        allTransactionTagsByGroupMap.value = {};

        for (const tag of tags) {
            allTransactionTagsMap.value[tag.id] = tag;
        }

        for (const tag of tags) {
            let tagsInGroup = allTransactionTagsByGroupMap.value[tag.groupId];

            if (!tagsInGroup) {
                tagsInGroup = [];
                allTransactionTagsByGroupMap.value[tag.groupId] = tagsInGroup;
            }

            tagsInGroup.push(tag);
        }
    }

    function addTagGroupToTransactionTagGroupList(tagGroup: TransactionTagGroup): void {
        allTransactionTagGroups.value.push(tagGroup);
        allTransactionTagGroupsMap.value[tagGroup.id] = tagGroup;
    }

    function addTagToTransactionTagList(tag: TransactionTag): void {
        allTransactionTags.value.push(tag);
        allTransactionTagsMap.value[tag.id] = tag;

        let tagsInGroup = allTransactionTagsByGroupMap.value[tag.groupId];

        if (!tagsInGroup) {
            tagsInGroup = [];
            allTransactionTagsByGroupMap.value[tag.groupId] = tagsInGroup;
        }

        tagsInGroup.push(tag);
    }

    function updateTagGroupInTransactionTagGroupList(currentTagGroup: TransactionTagGroup): void {
        for (const [transactionTagGroup, index] of itemAndIndex(allTransactionTagGroups.value)) {
            if (transactionTagGroup.id === currentTagGroup.id) {
                allTransactionTagGroups.value.splice(index, 1, currentTagGroup);
                break;
            }
        }

        allTransactionTagGroupsMap.value[currentTagGroup.id] = currentTagGroup;
    }

    function updateTagInTransactionTagList(currentTag: TransactionTag, oldTagGroupId?: string): void {
        // update in the main list
        for (const [transactionTag, index] of itemAndIndex(allTransactionTags.value)) {
            if (transactionTag.id === currentTag.id) {
                if (oldTagGroupId && oldTagGroupId !== currentTag.groupId) {
                    allTransactionTags.value.splice(index, 1);
                } else {
                    allTransactionTags.value.splice(index, 1, currentTag);
                }
                break;
            }
        }

        if (oldTagGroupId && oldTagGroupId !== currentTag.groupId) {
            let insertIndex = allTransactionTags.value.length;

            for (const [tag, index] of itemAndIndex(allTransactionTags.value)) {
                if (tag.groupId === currentTag.groupId) {
                    insertIndex = index;
                    break;
                }
            }

            allTransactionTags.value.splice(insertIndex, 0, currentTag);
        }

        // update in the map
        allTransactionTagsMap.value[currentTag.id] = currentTag;

        // update in the group list
        for (const tags of values(allTransactionTagsByGroupMap.value)) {
            for (const [transactionTag, index] of itemAndIndex(tags)) {
                if (transactionTag.id === currentTag.id) {
                    if (oldTagGroupId && oldTagGroupId !== currentTag.groupId) {
                        tags.splice(index, 1);
                    } else {
                        tags.splice(index, 1, currentTag);
                    }
                    break;
                }
            }
        }

        if (oldTagGroupId && oldTagGroupId !== currentTag.groupId) {
            let newGroupTags = allTransactionTagsByGroupMap.value[currentTag.groupId];

            if (!newGroupTags) {
                newGroupTags = [];
                allTransactionTagsByGroupMap.value[currentTag.groupId] = newGroupTags;
            }

            newGroupTags.push(currentTag);
        }
    }

    function updateTagGroupDisplayOrderInTransactionTagList({ from, to }: { from: number, to: number }): void {
        allTransactionTagGroups.value.splice(to, 0, allTransactionTagGroups.value.splice(from, 1)[0] as TransactionTagGroup);
    }

    function updateTagDisplayOrderInTransactionTagList({ groupId, from, to }: { groupId: string, from: number, to: number }): void {
        // update in the group list
        const tagsInGroup = allTransactionTagsByGroupMap.value[groupId];

        if (!tagsInGroup) {
            return;
        }

        const fromTag = tagsInGroup[from];

        if (!fromTag) {
            return;
        }

        const toTag = tagsInGroup[to];

        if (!toTag) {
            return;
        }

        tagsInGroup.splice(to, 0, tagsInGroup.splice(from, 1)[0] as TransactionTag);

        // update in the main list
        let mainListFromIndex = -1;
        let mainListToIndex = -1;

        for (const [tag, index] of itemAndIndex(allTransactionTags.value)) {
            if (tag.id === fromTag.id) {
                mainListFromIndex = index;
            }

            if (tag.id === toTag.id) {
                mainListToIndex = index;
            }
        }

        if (mainListFromIndex === -1 || mainListToIndex === -1) {
            return;
        }

        allTransactionTags.value.splice(mainListToIndex, 0, allTransactionTags.value.splice(mainListFromIndex, 1)[0] as TransactionTag);
    }

    function updateTagVisibilityInTransactionTagList({ tag, hidden }: { tag: TransactionTag, hidden: boolean }): void {
        if (allTransactionTagsMap.value[tag.id]) {
            allTransactionTagsMap.value[tag.id]!.hidden = hidden;
        }
    }

    function removeTagGroupFromTransactionTagGroupList(currentTagGroup: TransactionTagGroup): void {
        for (const [transactionTagGroup, index] of itemAndIndex(allTransactionTagGroups.value)) {
            if (transactionTagGroup.id === currentTagGroup.id) {
                allTransactionTagGroups.value.splice(index, 1);
                break;
            }
        }

        if (allTransactionTagGroupsMap.value[currentTagGroup.id]) {
            delete allTransactionTagGroupsMap.value[currentTagGroup.id];
        }
    }

    function removeTagFromTransactionTagList(currentTag: TransactionTag): void {
        for (const [transactionTag, index] of itemAndIndex(allTransactionTags.value)) {
            if (transactionTag.id === currentTag.id) {
                allTransactionTags.value.splice(index, 1);
                break;
            }
        }

        if (allTransactionTagsMap.value[currentTag.id]) {
            delete allTransactionTagsMap.value[currentTag.id];
        }

        for (const tags of values(allTransactionTagsByGroupMap.value)) {
            for (const [transactionTag, index] of itemAndIndex(tags)) {
                if (transactionTag.id === currentTag.id) {
                    tags.splice(index, 1);
                    break;
                }
            }
        }
    }

    function updateTransactionTagGroupListInvalidState(invalidState: boolean): void {
        transactionTagGroupListStateInvalid.value = invalidState;
    }

    function updateTransactionTagListInvalidState(invalidState: boolean): void {
        transactionTagListStateInvalid.value = invalidState;
    }

    function resetTransactionTags(): void {
        allTransactionTagGroups.value = [];
        allTransactionTagGroupsMap.value = {};
        allTransactionTags.value = [];
        allTransactionTagsMap.value = {};
        allTransactionTagsByGroupMap.value = {};
        transactionTagGroupListStateInvalid.value = true;
        transactionTagListStateInvalid.value = true;
    }

    function loadAllTagGroups({ force }: { force?: boolean }): Promise<TransactionTagGroup[]> {
        if (!force && !transactionTagGroupListStateInvalid.value) {
            return new Promise((resolve) => {
                resolve(allTransactionTagGroups.value);
            });
        }

        return new Promise((resolve, reject) => {
            services.getAllTransactionTagGroups().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve tag group list' });
                    return;
                }

                if (transactionTagGroupListStateInvalid.value) {
                    updateTransactionTagGroupListInvalidState(false);
                }

                const transactionTagGroups = TransactionTagGroup.ofMulti(data.result);

                loadTransactionTagGroupList(transactionTagGroups);

                resolve(transactionTagGroups);
            }).catch(error => {
                logger.error('failed to load tag group list', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve tag group list' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function loadAllTags({ force }: { force?: boolean }): Promise<TransactionTag[]> {
        if (!force && !transactionTagGroupListStateInvalid.value && !transactionTagListStateInvalid.value) {
            return new Promise((resolve) => {
                resolve(allTransactionTags.value);
            });
        }

        return new Promise((resolve, reject) => {
            services.getAllTransactionTagGroups().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve tag list' });
                    return;
                }

                if (transactionTagGroupListStateInvalid.value) {
                    updateTransactionTagGroupListInvalidState(false);
                }

                const transactionTagGroups = TransactionTagGroup.ofMulti(data.result);

                loadTransactionTagGroupList(transactionTagGroups);

                return services.getAllTransactionTags();
            }).then(response => {
                if (!response) {
                    reject({ message: 'Unable to retrieve tag list' });
                    return;
                }

                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve tag list' });
                    return;
                }

                if (transactionTagListStateInvalid.value) {
                    updateTransactionTagListInvalidState(false);
                }

                const transactionTags = TransactionTag.ofMulti(data.result);

                if (force && data.result && isEquals(allTransactionTags.value, transactionTags)) {
                    reject({ message: 'Tag list is up to date', isUpToDate: true });
                    return;
                }

                loadTransactionTagList(transactionTags);

                resolve(transactionTags);
            }).catch(error => {
                if (force) {
                    logger.error('failed to force load tag list', error);
                } else {
                    logger.error('failed to load tag list', error);
                }

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve tag list' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function saveTagGroup({ tagGroup }: { tagGroup: TransactionTagGroup }): Promise<TransactionTagGroup> {
        return new Promise((resolve, reject) => {
            let promise: ApiResponsePromise<TransactionTagGroupInfoResponse>;

            if (!tagGroup.id) {
                promise = services.addTransactionTagGroup(tagGroup.toCreateRequest());
            } else {
                promise = services.modifyTransactionTagGroup(tagGroup.toModifyRequest());
            }

            promise.then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (!tagGroup.id) {
                        reject({ message: 'Unable to add tag group' });
                    } else {
                        reject({ message: 'Unable to save tag group' });
                    }
                    return;
                }

                const transactionTagGroup = TransactionTagGroup.of(data.result);

                if (!tagGroup.id) {
                    addTagGroupToTransactionTagGroupList(transactionTagGroup);
                } else {
                    updateTagGroupInTransactionTagGroupList(transactionTagGroup);
                }

                resolve(transactionTagGroup);
            }).catch(error => {
                logger.error('failed to save tag group', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    if (!tagGroup.id) {
                        reject({ message: 'Unable to add tag group' });
                    } else {
                        reject({ message: 'Unable to save tag group' });
                    }
                } else {
                    reject(error);
                }
            });
        });
    }

    function changeTagGroupDisplayOrder({ tagGroupId, from, to }: { tagGroupId: string, from: number, to: number }): Promise<void> {
        return new Promise((resolve, reject) => {
            let currentTagGroup: TransactionTagGroup | null = null;

            for (const transactionTagGroup of allTransactionTagGroups.value) {
                if (transactionTagGroup.id === tagGroupId) {
                    currentTagGroup = transactionTagGroup;
                    break;
                }
            }

            if (!currentTagGroup || !allTransactionTagGroups.value[to]) {
                reject({ message: 'Unable to move tag group' });
                return;
            }

            if (!transactionTagGroupListStateInvalid.value) {
                updateTransactionTagGroupListInvalidState(true);
            }

            updateTagGroupDisplayOrderInTransactionTagList({ from, to });

            resolve();
        });
    }

    function updateTagGroupDisplayOrders(): Promise<boolean> {
        const newDisplayOrders: TransactionTagGroupNewDisplayOrderRequest[] = [];

        for (const [transactionTagGroup, index] of itemAndIndex(allTransactionTagGroups.value)) {
            newDisplayOrders.push({
                id: transactionTagGroup.id,
                displayOrder: index + 1
            });
        }

        return new Promise((resolve, reject) => {
            services.moveTransactionTagGroup({
                newDisplayOrders: newDisplayOrders
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to move tag group' });
                    return;
                }

                if (transactionTagGroupListStateInvalid.value) {
                    updateTransactionTagGroupListInvalidState(false);
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to save tag groups display order', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to move tag group' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function deleteTagGroup({ tagGroup, beforeResolve }: { tagGroup: TransactionTagGroup, beforeResolve?: BeforeResolveFunction }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.deleteTransactionTagGroup({
                id: tagGroup.id
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to delete this tag group' });
                    return;
                }

                if (beforeResolve) {
                    beforeResolve(() => {
                        removeTagGroupFromTransactionTagGroupList(tagGroup);
                    });
                } else {
                    removeTagGroupFromTransactionTagGroupList(tagGroup);
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to delete tag group', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to delete this tag group' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function saveTag({ tag }: { tag: TransactionTag }): Promise<TransactionTag> {
        const oldTagGroupId = allTransactionTagsMap.value[tag.id]?.groupId;

        return new Promise((resolve, reject) => {
            let promise: ApiResponsePromise<TransactionTagInfoResponse>;

            if (!tag.id) {
                promise = services.addTransactionTag(tag.toCreateRequest());
            } else {
                promise = services.modifyTransactionTag(tag.toModifyRequest());
            }

            promise.then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (!tag.id) {
                        reject({ message: 'Unable to add tag' });
                    } else {
                        reject({ message: 'Unable to save tag' });
                    }
                    return;
                }

                const transactionTag = TransactionTag.of(data.result);

                if (!tag.id) {
                    addTagToTransactionTagList(transactionTag);
                } else {
                    updateTagInTransactionTagList(transactionTag, oldTagGroupId);
                }

                resolve(transactionTag);
            }).catch(error => {
                logger.error('failed to save tag', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    if (!tag.id) {
                        reject({ message: 'Unable to add tag' });
                    } else {
                        reject({ message: 'Unable to save tag' });
                    }
                } else {
                    reject(error);
                }
            });
        });
    }

    function addTags(req: TransactionTagCreateBatchRequest): Promise<TransactionTag[]> {
        return new Promise((resolve, reject) => {
            services.addTransactionTagBatch(req).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to add tag' });
                    return;
                }

                if (!transactionTagListStateInvalid.value) {
                    updateTransactionTagListInvalidState(true);
                }

                const transactionTags = TransactionTag.ofMulti(data.result);

                resolve(transactionTags);
            }).catch(error => {
                logger.error('failed to add tags', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to add tag' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function changeTagDisplayOrder({ tagId, from, to }: { tagId: string, from: number, to: number }): Promise<void> {
        return new Promise((resolve, reject) => {
            let currentTag: TransactionTag | null = null;

            for (const transactionTag of allTransactionTags.value) {
                if (transactionTag.id === tagId) {
                    currentTag = transactionTag;
                    break;
                }
            }

            if (!currentTag || !allTransactionTags.value[to]) {
                reject({ message: 'Unable to move tag' });
                return;
            }

            if (!transactionTagListStateInvalid.value) {
                updateTransactionTagListInvalidState(true);
            }

            updateTagDisplayOrderInTransactionTagList({
                groupId: currentTag.groupId,
                from: from,
                to: to
            });

            resolve();
        });
    }

    function updateTagDisplayOrders(groupId: string): Promise<boolean> {
        const tagsInGroup = allTransactionTagsByGroupMap.value[groupId];

        if (!tagsInGroup) {
            return Promise.reject('Unable to move tag');
        }

        const newDisplayOrders: TransactionTagNewDisplayOrderRequest[] = [];

        for (const [transactionTag, index] of itemAndIndex(tagsInGroup)) {
            newDisplayOrders.push({
                id: transactionTag.id,
                displayOrder: index + 1
            });
        }

        return new Promise((resolve, reject) => {
            services.moveTransactionTag({
                newDisplayOrders: newDisplayOrders
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to move tag' });
                    return;
                }

                if (transactionTagListStateInvalid.value) {
                    updateTransactionTagListInvalidState(false);
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to save tags display order', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to move tag' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function hideTag({ tag, hidden }: { tag: TransactionTag, hidden: boolean }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.hideTransactionTag({
                id: tag.id,
                hidden: hidden
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (hidden) {
                        reject({ message: 'Unable to hide this tag' });
                    } else {
                        reject({ message: 'Unable to unhide this tag' });
                    }
                    return;
                }

                updateTagVisibilityInTransactionTagList({ tag, hidden });

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to change tag visibility', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    if (hidden) {
                        reject({ message: 'Unable to hide this tag' });
                    } else {
                        reject({ message: 'Unable to unhide this tag' });
                    }
                } else {
                    reject(error);
                }
            });
        });
    }

    function deleteTag({ tag, beforeResolve }: { tag: TransactionTag, beforeResolve?: BeforeResolveFunction }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.deleteTransactionTag({
                id: tag.id
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to delete this tag' });
                    return;
                }

                if (beforeResolve) {
                    beforeResolve(() => {
                        removeTagFromTransactionTagList(tag);
                    });
                } else {
                    removeTagFromTransactionTagList(tag);
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to delete tag', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to delete this tag' });
                } else {
                    reject(error);
                }
            });
        });
    }

    return {
        // states
        allTransactionTagGroups,
        allTransactionTagGroupsMap,
        allTransactionTags,
        allTransactionTagsMap,
        allTransactionTagsByGroupMap,
        transactionTagGroupListStateInvalid,
        transactionTagListStateInvalid,
        // computed states
        allAvailableTagsCount,
        // functions
        updateTransactionTagListInvalidState,
        resetTransactionTags,
        loadAllTagGroups,
        loadAllTags,
        saveTagGroup,
        changeTagGroupDisplayOrder,
        updateTagGroupDisplayOrders,
        deleteTagGroup,
        saveTag,
        addTags,
        changeTagDisplayOrder,
        updateTagDisplayOrders,
        hideTag,
        deleteTag
    }
});
