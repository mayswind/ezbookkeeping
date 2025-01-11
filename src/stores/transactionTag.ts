import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import type { BeforeResolveFunction } from '@/core/base.ts';

import {
    type TransactionTagInfoResponse,
    type TransactionTagNewDisplayOrderRequest,
    TransactionTag
} from '@/models/transaction_tag.ts';

import { isEquals } from '@/lib/common.ts';

import logger from '@/lib/logger.ts';
import services, { type ApiResponsePromise } from '@/lib/services.ts';

export const useTransactionTagsStore = defineStore('transactionTags', () => {
    const allTransactionTags = ref<TransactionTag[]>([]);
    const allTransactionTagsMap = ref<Record<string, TransactionTag>>({});
    const transactionTagListStateInvalid = ref<boolean>(true);

    const allVisibleTags = computed<TransactionTag[]>(() => {
        const visibleTags: TransactionTag[] = [];

        for (const tag of allTransactionTags.value) {
            if (!tag.hidden) {
                visibleTags.push(tag);
            }
        }

        return visibleTags;
    });

    const allAvailableTagsCount = computed<number>(() => allTransactionTags.value.length);
    const allVisibleTagsCount = computed<number>(() => allVisibleTags.value.length);

    function loadTransactionTagList(tags: TransactionTag[]): void {
        allTransactionTags.value = tags;
        allTransactionTagsMap.value = {};

        for (const tag of tags) {
            allTransactionTagsMap.value[tag.id] = tag;
        }
    }

    function addTagToTransactionTagList(tag: TransactionTag): void {
        allTransactionTags.value.push(tag);
        allTransactionTagsMap.value[tag.id] = tag;
    }

    function updateTagInTransactionTagList(tag: TransactionTag): void {
        for (let i = 0; i < allTransactionTags.value.length; i++) {
            if (allTransactionTags.value[i].id === tag.id) {
                allTransactionTags.value.splice(i, 1, tag);
                break;
            }
        }

        allTransactionTagsMap.value[tag.id] = tag;
    }

    function updateTagDisplayOrderInTransactionTagList(params: { from: number, to: number }): void {
        allTransactionTags.value.splice(params.to, 0, allTransactionTags.value.splice(params.from, 1)[0]);
    }

    function updateTagVisibilityInTransactionTagList(params: { tag: TransactionTag, hidden: boolean }): void {
        if (allTransactionTagsMap.value[params.tag.id]) {
            allTransactionTagsMap.value[params.tag.id].hidden = params.hidden;
        }
    }

    function removeTagFromTransactionTagList(tag: TransactionTag): void {
        for (let i = 0; i < allTransactionTags.value.length; i++) {
            if (allTransactionTags.value[i].id === tag.id) {
                allTransactionTags.value.splice(i, 1);
                break;
            }
        }

        if (allTransactionTagsMap.value[tag.id]) {
            delete allTransactionTagsMap.value[tag.id];
        }
    }

    function updateTransactionTagListInvalidState(invalidState: boolean): void {
        transactionTagListStateInvalid.value = invalidState;
    }

    function resetTransactionTags(): void {
        allTransactionTags.value = [];
        allTransactionTagsMap.value = {};
        transactionTagListStateInvalid.value = true;
    }

    function loadAllTags(params: { force?: boolean }): Promise<TransactionTag[]> {
        if (!params.force && !transactionTagListStateInvalid.value) {
            return new Promise((resolve) => {
                resolve(allTransactionTags.value);
            });
        }

        return new Promise((resolve, reject) => {
            services.getAllTransactionTags().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve tag list' });
                    return;
                }

                if (transactionTagListStateInvalid.value) {
                    updateTransactionTagListInvalidState(false);
                }

                const transactionTags = TransactionTag.ofMany(data.result);

                if (params.force && data.result && isEquals(allTransactionTags.value, transactionTags)) {
                    reject({ message: 'Tag list is up to date' });
                    return;
                }

                loadTransactionTagList(transactionTags);

                resolve(transactionTags);
            }).catch(error => {
                if (params.force) {
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

    function saveTag(params: { tag: TransactionTag }): Promise<TransactionTag> {
        return new Promise((resolve, reject) => {
            let promise: ApiResponsePromise<TransactionTagInfoResponse>;

            if (!params.tag.id) {
                promise = services.addTransactionTag(params.tag.toCreateRequest());
            } else {
                promise = services.modifyTransactionTag(params.tag.toModifyRequest());
            }

            promise.then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (!params.tag.id) {
                        reject({ message: 'Unable to add tag' });
                    } else {
                        reject({ message: 'Unable to save tag' });
                    }
                    return;
                }

                const transactionTag = TransactionTag.of(data.result);

                if (!params.tag.id) {
                    addTagToTransactionTagList(transactionTag);
                } else {
                    updateTagInTransactionTagList(transactionTag);
                }

                resolve(transactionTag);
            }).catch(error => {
                logger.error('failed to save tag', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    if (!params.tag.id) {
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

    function changeTagDisplayOrder(params: { tagId: string, from: number, to: number }): Promise<void> {
        return new Promise((resolve, reject) => {
            let tag: TransactionTag | null = null;

            for (let i = 0; i < allTransactionTags.value.length; i++) {
                if (allTransactionTags.value[i].id === params.tagId) {
                    tag = allTransactionTags.value[i];
                    break;
                }
            }

            if (!tag || !allTransactionTags.value[params.to]) {
                reject({ message: 'Unable to move tag' });
                return;
            }

            if (!transactionTagListStateInvalid.value) {
                updateTransactionTagListInvalidState(true);
            }

            updateTagDisplayOrderInTransactionTagList({
                from: params.from,
                to: params.to
            });

            resolve();
        });
    }

    function updateTagDisplayOrders(): Promise<boolean> {
        const newDisplayOrders: TransactionTagNewDisplayOrderRequest[] = [];

        for (let i = 0; i < allTransactionTags.value.length; i++) {
            newDisplayOrders.push({
                id: allTransactionTags.value[i].id,
                displayOrder: i + 1
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

    function hideTag(params: { tag: TransactionTag, hidden: boolean }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.hideTransactionTag({
                id: params.tag.id,
                hidden: params.hidden
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (params.hidden) {
                        reject({ message: 'Unable to hide this tag' });
                    } else {
                        reject({ message: 'Unable to unhide this tag' });
                    }
                    return;
                }

                updateTagVisibilityInTransactionTagList({
                    tag: params.tag,
                    hidden: params.hidden
                });

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to change tag visibility', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    if (params.hidden) {
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

    function deleteTag(params: { tag: TransactionTag, beforeResolve?: BeforeResolveFunction }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.deleteTransactionTag({
                id: params.tag.id
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to delete this tag' });
                    return;
                }

                if (params.beforeResolve) {
                    params.beforeResolve(() => {
                        removeTagFromTransactionTagList(params.tag);
                    });
                } else {
                    removeTagFromTransactionTagList(params.tag);
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
        allTransactionTags,
        allTransactionTagsMap,
        transactionTagListStateInvalid,
        // computed states
        allVisibleTags,
        allAvailableTagsCount,
        allVisibleTagsCount,
        // functions
        updateTransactionTagListInvalidState,
        resetTransactionTags,
        loadAllTags,
        saveTag,
        changeTagDisplayOrder,
        updateTagDisplayOrders,
        hideTag,
        deleteTag
    }
});
