import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import {
    type QueuedTransactionCreateRequest,
    QueuedTransaction
} from '@/models/queued_transaction.ts';

import logger from '@/lib/logger.ts';
import services from '@/lib/services.ts';

/** Highest count the queue badge shows as a number; anything above is rendered as "+99" */
export const MAX_DISPLAY_QUEUED_TRANSACTION_COUNT = 99;

export const useQueuedTransactionsStore = defineStore('queuedTransactions', () => {
    const allQueuedTransactions = ref<QueuedTransaction[]>([]);
    const pendingCount = ref<number>(0);
    const archivedCount = ref<number>(0);
    const queuedTransactionListStateInvalid = ref<boolean>(true);
    const loadedWithArchived = ref<boolean>(false);

    const pendingQueuedTransactions = computed<QueuedTransaction[]>(() => {
        return allQueuedTransactions.value.filter(item => !item.isArchived);
    });

    const archivedQueuedTransactions = computed<QueuedTransaction[]>(() => {
        return allQueuedTransactions.value.filter(item => item.isArchived);
    });

    /** The pending count rendered on the queue button, capped so the badge stays readable */
    const displayPendingCount = computed<string>(() => {
        if (pendingCount.value > MAX_DISPLAY_QUEUED_TRANSACTION_COUNT) {
            return '+' + MAX_DISPLAY_QUEUED_TRANSACTION_COUNT;
        }

        return pendingCount.value.toString();
    });

    function updateQueuedTransactionListInvalidState(invalidState: boolean): void {
        queuedTransactionListStateInvalid.value = invalidState;
    }

    function resetQueuedTransactions(): void {
        allQueuedTransactions.value = [];
        pendingCount.value = 0;
        archivedCount.value = 0;
        queuedTransactionListStateInvalid.value = true;
        loadedWithArchived.value = false;
    }

    function setCountsFromLoadedList(): void {
        pendingCount.value = pendingQueuedTransactions.value.length;

        // Only a listing that asked for archived items knows how many there are
        if (loadedWithArchived.value) {
            archivedCount.value = archivedQueuedTransactions.value.length;
        }
    }

    function removeQueuedTransactionFromList(id: string): void {
        allQueuedTransactions.value = allQueuedTransactions.value.filter(item => item.id !== id);
        setCountsFromLoadedList();
    }

    function loadQueuedTransactionCount(): Promise<number> {
        return new Promise((resolve, reject) => {
            services.getQueuedTransactionCount().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve the transaction queue count' });
                    return;
                }

                pendingCount.value = data.result.pendingCount;
                archivedCount.value = data.result.archivedCount;

                resolve(data.result.pendingCount);
            }).catch(error => {
                logger.error('failed to load transaction queue count', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve the transaction queue count' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function loadAllQueuedTransactions({ includeArchived, force }: { includeArchived: boolean, force?: boolean }): Promise<QueuedTransaction[]> {
        if (!force && !queuedTransactionListStateInvalid.value && loadedWithArchived.value === includeArchived) {
            return new Promise((resolve) => {
                resolve(allQueuedTransactions.value);
            });
        }

        return new Promise((resolve, reject) => {
            services.getAllQueuedTransactions({ includeArchived }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve the transaction queue' });
                    return;
                }

                allQueuedTransactions.value = QueuedTransaction.ofMulti(data.result);
                loadedWithArchived.value = includeArchived;
                updateQueuedTransactionListInvalidState(false);
                setCountsFromLoadedList();

                resolve(allQueuedTransactions.value);
            }).catch(error => {
                logger.error('failed to load the transaction queue', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve the transaction queue' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function addQueuedTransaction(req: QueuedTransactionCreateRequest): Promise<QueuedTransaction> {
        return new Promise((resolve, reject) => {
            services.addQueuedTransaction(req).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to add the draft to the transaction queue' });
                    return;
                }

                updateQueuedTransactionListInvalidState(true);

                resolve(QueuedTransaction.of(data.result));
            }).catch(error => {
                logger.error('failed to add a draft to the transaction queue', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to add the draft to the transaction queue' });
                } else {
                    reject(error);
                }
            });
        });
    }

    /** Rejects a draft: it is archived now and purged once its retention period elapses */
    function archiveQueuedTransaction({ id }: { id: string }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.archiveQueuedTransaction({ id }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to reject this draft' });
                    return;
                }

                updateQueuedTransactionListInvalidState(true);

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to reject a draft', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to reject this draft' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function archiveAllQueuedTransactions(): Promise<number> {
        return new Promise((resolve, reject) => {
            services.archiveAllQueuedTransactions().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to reject all drafts' });
                    return;
                }

                updateQueuedTransactionListInvalidState(true);

                resolve(data.result.affectedCount);
            }).catch(error => {
                logger.error('failed to reject all drafts', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to reject all drafts' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function restoreQueuedTransaction({ id }: { id: string }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.restoreQueuedTransaction({ id }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to restore this draft' });
                    return;
                }

                updateQueuedTransactionListInvalidState(true);

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to restore a draft', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to restore this draft' });
                } else {
                    reject(error);
                }
            });
        });
    }

    /**
     * Drops a draft from the queue. This is what accepting one does once the real transaction has
     * been created out of it, so the draft does not linger as a duplicate.
     */
    function deleteQueuedTransaction({ id }: { id: string }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.deleteQueuedTransaction({ id }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to remove this draft from the queue' });
                    return;
                }

                removeQueuedTransactionFromList(id);

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to remove a draft from the queue', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to remove this draft from the queue' });
                } else {
                    reject(error);
                }
            });
        });
    }

    return {
        // states
        allQueuedTransactions,
        pendingCount,
        archivedCount,
        queuedTransactionListStateInvalid,
        // computed states
        pendingQueuedTransactions,
        archivedQueuedTransactions,
        displayPendingCount,
        // functions
        updateQueuedTransactionListInvalidState,
        resetQueuedTransactions,
        loadQueuedTransactionCount,
        loadAllQueuedTransactions,
        addQueuedTransaction,
        archiveQueuedTransaction,
        archiveAllQueuedTransactions,
        restoreQueuedTransaction,
        deleteQueuedTransaction
    };
});
