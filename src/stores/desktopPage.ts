import { ref } from 'vue';
import { defineStore } from 'pinia';

import { itemAndIndex } from '@/core/base.ts';
import type { ErrorResponse } from '@/core/api.ts';

import { isDefined, isObject, isString } from '@/lib/common.ts';
import { generateRandomUUID } from '@/lib/misc.ts';

interface DesktopMessageToast {
    readonly type: 'message'
    readonly id: string;
    readonly message: string;
    readonly options?: Record<string, unknown>;
}

interface DesktopErrorToast {
    readonly type: 'error'
    readonly id: string;
    readonly error: string | { message: string } | { error: ErrorResponse };
}

export type DesktopToast = DesktopMessageToast | DesktopErrorToast;

export const useDesktopPageStore = defineStore('desktopPages', () => {
    const pendingToastQueue = ref<DesktopToast[]>([]);
    const showingToastQueue = ref<DesktopToast[]>([]);

    const showAddTransactionDialogInTransactionList = ref<boolean>(false);

    function isToastEquals(toast1: DesktopToast, toast2: DesktopToast): boolean {
        if (toast1.type !== toast2.type) {
            return false;
        }

        if (toast1.type === 'message' && toast2.type === 'message') {
            return toast1.message === toast2.message && ((!isDefined(toast1.options) && !isDefined(toast2.options)) || JSON.stringify(toast1.options) === JSON.stringify(toast2.options));
        } else if (toast1.type === 'error' && toast2.type === 'error') {
            if (isString(toast1.error) && isString(toast2.error)) {
                return toast1.error === toast2.error;
            } else if (isObject(toast1.error) && isObject(toast2.error) && 'message' in toast1.error && 'message' in toast2.error) {
                return toast1.error.message === toast2.error.message;
            } else if (isObject(toast1.error) && isObject(toast2.error) && 'error' in toast1.error && 'error' in toast2.error) {
                return JSON.stringify(toast1.error.error) === JSON.stringify(toast2.error.error);
            }
        }

        return false;
    }

    function addToastToPendingQueue(toast: DesktopToast): void {
        if (showingToastQueue.value && showingToastQueue.value.length > 0) {
            for (const showingToast of showingToastQueue.value) {
                if (isToastEquals(showingToast, toast)) {
                    return;
                }
            }
        }

        if (pendingToastQueue.value.length > 1) {
            const lastPendingToast = pendingToastQueue.value[pendingToastQueue.value.length - 1];

            if (lastPendingToast && isToastEquals(lastPendingToast, toast)) {
                return;
            }
        }

        pendingToastQueue.value.push(toast);
    }

    function addToast(message: string, options?: Record<string, unknown>): void {
        if (!message || message.trim() === '') {
            return;
        }

        const newToast: DesktopMessageToast = {
            type: 'message',
            id: generateRandomUUID(),
            message: message,
            options: options
        };

        addToastToPendingQueue(newToast);
    }

    function addErrorToast(error: string | { message: string } | { error: ErrorResponse }): void {
        if (!error || (isString(error) && error.trim() === '')) {
            return;
        }

        const newToast: DesktopErrorToast = {
            type: 'error',
            id: generateRandomUUID(),
            error: error
        };

        addToastToPendingQueue(newToast);
    }

    function removePendingToast(count: number): void {
        const toasts: DesktopToast[] = pendingToastQueue.value.splice(0, count);
        showingToastQueue.value.push(...toasts);
    }

    function removeShowingToast(id: string): void {
        for (const [item, index] of itemAndIndex(showingToastQueue.value)) {
            if (item.id === id) {
                showingToastQueue.value.splice(index, 1);
                break;
            }
        }
    }

    function setShowAddTransactionDialogInTransactionList(): void {
        showAddTransactionDialogInTransactionList.value = true;
    }

    function resetShowAddTransactionDialogInTransactionList(): void {
        showAddTransactionDialogInTransactionList.value = false;
    }

    return {
        // states
        pendingToastQueue,
        showAddTransactionDialogInTransactionList,
        // functions
        addToast,
        addErrorToast,
        removePendingToast,
        removeShowingToast,
        setShowAddTransactionDialogInTransactionList,
        resetShowAddTransactionDialogInTransactionList
    }
});
