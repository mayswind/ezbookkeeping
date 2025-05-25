<template>
    <v-dialog width="800" :persistent="submitting || (!!defaultCurrencyAmount && defaultCurrencyAmount !== 1) || currency !== defaultCurrency || (!!targetCurrencyAmount && targetCurrencyAmount !== 1)" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-4">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <div class="d-flex w-100 align-center justify-center">
                        <h4 class="text-h4">{{ tt('Update User Custom Exchange Rate') }}</h4>
                    </div>
                </div>
            </template>
            <v-card-text class="my-md-4 w-100 d-flex justify-center">
                <v-row>
                    <v-col cols="12" md="6">
                        <v-text-field type="number"
                                      :disabled="submitting"
                                      :label="tt('Amount')"
                                      :placeholder="tt('Amount')"
                                      :persistent-placeholder="true"
                                      v-model="defaultCurrencyAmount"/>
                    </v-col>
                    <v-col cols="12" md="6">
                        <currency-select :disabled="true"
                                         :label="tt('Currency')"
                                         :placeholder="tt('Currency')"
                                         v-model="defaultCurrency" />
                    </v-col>
                    <v-col cols="12" class="text-center my-2">
                        <v-icon :icon="mdiSwapVertical" size="24" />
                    </v-col>
                    <v-col cols="12" md="6">
                        <v-text-field type="number"
                                      :disabled="submitting"
                                      :label="tt('Amount')"
                                      :placeholder="tt('Amount')"
                                      :persistent-placeholder="true"
                                      v-model="targetCurrencyAmount"/>
                    </v-col>
                    <v-col cols="12" md="6">
                        <currency-select :disabled="submitting"
                                         :label="tt('Currency')"
                                         :placeholder="tt('Currency')"
                                         v-model="currency" />
                    </v-col>
                </v-row>
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center gap-4">
                    <v-btn :disabled="submitting || !defaultCurrencyAmount || !currency || !targetCurrencyAmount" @click="confirm">
                        {{ tt('OK') }}
                        <v-progress-circular indeterminate size="22" class="ml-2" v-if="submitting"></v-progress-circular>
                    </v-btn>
                    <v-btn color="secondary" variant="tonal" :disabled="submitting" @click="cancel">{{ tt('Cancel') }}</v-btn>
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

import { useUserStore } from '@/stores/user.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import {
    mdiSwapVertical
} from '@mdi/js';

interface UserCustomExchangeRateUpdateResponse {
    message: string;
}

type SnackBarType = InstanceType<typeof SnackBar>;

defineProps<{
    show?: boolean;
}>();

const { tt } = useI18n();

const userStore = useUserStore();
const exchangeRatesStore = useExchangeRatesStore();

const showState = ref<boolean>(false);
const submitting = ref<boolean>(false);
const defaultCurrency = ref<string>(userStore.currentUserDefaultCurrency);
const defaultCurrencyAmount = ref<number>(1);
const currency = ref<string>(userStore.currentUserDefaultCurrency);
const targetCurrencyAmount = ref<number>(1);

const snackbar = useTemplateRef<SnackBarType>('snackbar');

let resolveFunc: ((response: UserCustomExchangeRateUpdateResponse) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

function open(): Promise<UserCustomExchangeRateUpdateResponse> {
    showState.value = true;
    defaultCurrencyAmount.value = 1;
    currency.value = userStore.currentUserDefaultCurrency;
    targetCurrencyAmount.value = 1;

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function confirm(): void {
    submitting.value = true;

    exchangeRatesStore.updateUserCustomExchangeRate({
        currency: currency.value,
        rate: targetCurrencyAmount.value / defaultCurrencyAmount.value
    }).then(() => {
        submitting.value = false;
        resolveFunc?.({ message: 'You have updated exchange rate' });
        showState.value = false;
    }).catch(error => {
        submitting.value = false;

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
