<template>
    <v-dialog width="600" :persistent="true" v-model="showState">
        <one-column-dialog-layout :disabled="loading || submitting" :title="title" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #after-title>
                <v-btn density="compact" color="default" variant="text" size="22"
                       class="ms-2" :icon="true" :disabled="loading || submitting"
                       :loading="loading" @click="reload">
                    <template #loader>
                        <v-progress-circular indeterminate size="20"/>
                    </template>
                    <v-icon :icon="mdiRefresh" size="22" />
                    <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                </v-btn>
            </template>

            <template #content>
                <div class="mt-4" v-if="type === 'clear'">{{ tt('format.misc.clearTransactionsTagsTip', { count: formatNumberToLocalizedNumerals(updateIds?.length ?? 0) }) }}</div>
                <div class="mt-5" v-if="type !== 'clear'">
                    <transaction-tag-auto-complete
                        :disabled="loading || submitting"
                        :show-label="true"
                        :allow-add-new-tag="type === 'add'"
                        v-model="tagIds"
                        @tag:saving="onSavingTag"
                    />
                </div>
            </template>

            <template #footer>
                <v-btn color="secondary" variant="tonal" :disabled="loading || submitting" @click="cancel">{{ tt('Cancel') }}</v-btn>
                <v-spacer/>
                <v-btn :disabled="loading || submitting || updateIds.length < 1 || (type !== 'clear' && (!tagIds || tagIds.length < 1))" @click="confirm">
                    {{ tt('OK') }}
                    <v-progress-circular indeterminate size="22" class="ms-2" v-if="submitting"></v-progress-circular>
                </v-btn>
            </template>
        </one-column-dialog-layout>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import {computed, ref, useTemplateRef} from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';

import {
    mdiRefresh
} from '@mdi/js'

export type BatchUpdateTagsOperationType = 'add' | 'remove' | 'clear';

type SnackBarType = InstanceType<typeof SnackBar>;

const {
    tt,
    formatNumberToLocalizedNumerals
} = useI18n();

const transactionTagsStore = useTransactionTagsStore();
const transactionsStore = useTransactionsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const loading = ref<boolean>(false);
const submitting = ref<boolean>(false);
const type = ref<BatchUpdateTagsOperationType>('add');
const updateIds = ref<string[]>([]);
const tagIds = ref<string[]>([]);

let resolveFunc: ((response: number) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const title = computed<string>(() => {
    if (type.value === 'add') {
        return tt('Add Tags to Transactions');
    } else if (type.value === 'remove') {
        return tt('Remove Tags from Transactions');
    } else if (type.value === 'clear') {
        return tt('Clear All Tags from Transactions');
    }

    return '';
});

function open(options: { type: BatchUpdateTagsOperationType; updateIds: string[] }): Promise<number> {
    type.value = options.type;
    updateIds.value = options.updateIds;
    tagIds.value = [];
    submitting.value = false;
    showState.value = true;

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function reload(): void {
    loading.value = true;

    transactionTagsStore.loadAllTags({ force: true }).then(() => {
        loading.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function confirm(): void {
    if (type.value === 'add') {
        submitting.value = true;

        transactionsStore.batchAddTagsToTransaction({
            transactionIds: updateIds.value,
            tagIds: tagIds.value
        }).then(() => {
            submitting.value = false;
            showState.value = false;
            resolveFunc?.(updateIds.value.length);
        }).catch(error => {
            submitting.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    } else if (type.value === 'remove') {
        submitting.value = true;

        transactionsStore.batchRemoveTagsFromTransaction({
            transactionIds: updateIds.value,
            tagIds: tagIds.value
        }).then(() => {
            submitting.value = false;
            showState.value = false;
            resolveFunc?.(updateIds.value.length);
        }).catch(error => {
            submitting.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    } else if (type.value === 'clear') {
        submitting.value = true;

        transactionsStore.batchClearAllTagsFromTransaction({
            transactionIds: updateIds.value
        }).then(() => {
            submitting.value = false;
            showState.value = false;
            resolveFunc?.(updateIds.value.length);
        }).catch(error => {
            submitting.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    }
}

function cancel(): void {
    rejectFunc?.();
    showState.value = false;
}

function onSavingTag(state: boolean): void {
    submitting.value = state;
}

defineExpose({
    open
});
</script>
