<template>
    <v-dialog width="600" :persistent="true" v-model="showState">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <h4 class="text-h4 text-error text-wrap">{{ tt('Delete Transactions') }}</h4>
            </template>
            <v-card-text class="pb-2 text-error">{{ tt('format.misc.deleteTransactionsTip', { count: formatNumberToLocalizedNumerals(deleteIds?.length ?? 0) }) }}</v-card-text>
            <v-card-text class="w-100 d-flex justify-center">
                <div class="w-100">
                    <v-text-field
                        autocomplete="current-password"
                        type="password"
                        variant="underlined"
                        color="error"
                        :disabled="deleting"
                        :placeholder="tt('Current Password')"
                        v-model="currentPassword"
                    />
                </div>
            </v-card-text>
            <v-card-text>
                <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                    <v-btn color="error" :disabled="!currentPassword || deleting || deleteIds.length < 1" @click="confirm">
                        {{ tt('Confirm') }}
                        <v-progress-circular indeterminate size="22" class="ms-2" v-if="deleting"></v-progress-circular>
                    </v-btn>
                    <v-btn color="secondary" variant="tonal" :disabled="deleting" @click="cancel">{{ tt('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionsStore } from '@/stores/transaction.ts';

type SnackBarType = InstanceType<typeof SnackBar>;

const {
    tt,
    formatNumberToLocalizedNumerals
} = useI18n();

const transactionsStore = useTransactionsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const deleting = ref<boolean>(false);
const deleteIds = ref<string[]>([]);
const currentPassword = ref<string>('');

let resolveFunc: ((response: number) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

function open(options: { updateIds: string[] }): Promise<number> {
    deleteIds.value = options.updateIds;
    currentPassword.value = '';
    deleting.value = false;
    showState.value = true;

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function confirm(): void {
    deleting.value = true;

    transactionsStore.batchDeleteTransactions({
        transactionIds: deleteIds.value,
        password: currentPassword.value
    }).then(() => {
        deleting.value = false;
        showState.value = false;
        resolveFunc?.(deleteIds.value.length);
    }).catch(error => {
        deleting.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function cancel(): void {
    rejectFunc?.();
    showState.value = false;
}

defineExpose({
    open
});
</script>
