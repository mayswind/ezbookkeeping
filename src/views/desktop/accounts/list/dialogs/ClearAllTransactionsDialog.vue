<template>
    <v-dialog width="640" :persistent="true" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-4">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <h4 class="text-h4 text-error text-wrap">{{ tt('Are you sure you want to clear all transactions?') }}</h4>
                </div>
            </template>
            <v-card-text class="pb-2 text-error">{{ tt('format.misc.clearTransactionsInAccountTip', { account: currentAccount?.name }) }}</v-card-text>
            <v-card-text class="mb-md-4 w-100 d-flex justify-center">
                <div class="w-100">
                    <v-text-field
                        autocomplete="current-password"
                        type="password"
                        variant="underlined"
                        color="error"
                        :disabled="clearingData"
                        :placeholder="tt('Current Password')"
                        v-model="currentPassword"
                    />
                </div>
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center gap-4">
                    <v-btn color="error" :disabled="!currentPassword || clearingData" @click="confirm">
                        {{ tt('Confirm') }}
                        <v-progress-circular indeterminate size="22" class="ms-2" v-if="clearingData"></v-progress-circular>
                    </v-btn>
                    <v-btn color="secondary" variant="tonal" :disabled="clearingData" @click="cancel">
                        {{ tt('Cancel') }}
                    </v-btn>
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

import { useRootStore } from '@/stores/index.ts';

import { Account } from '@/models/account.ts';

type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();

const rootStore = useRootStore();

let resolveFunc: (() => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const clearingData = ref<boolean>(false);
const currentPassword = ref<string>('');
const currentAccount = ref<Account | undefined>(undefined);

function open(account: Account): Promise<void> {
    showState.value = true;
    clearingData.value = false;
    currentPassword.value = '';
    currentAccount.value = account;

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function confirm(): void {
    if (!currentAccount.value || !currentPassword.value) {
        return;
    }

    clearingData.value = true;

    rootStore.clearAllUserTransactionsOfAccount({
        accountId: currentAccount.value.id,
        password: currentPassword.value
    }).then(() => {
        clearingData.value = false;

        resolveFunc?.();
        showState.value = false;
    }).catch(error => {
        clearingData.value = false;

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
