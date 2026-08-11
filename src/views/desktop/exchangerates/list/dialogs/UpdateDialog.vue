<template>
    <v-dialog width="800" :persistent="submitting || (!!defaultCurrencyAmount && defaultCurrencyAmount !== 1) || currency !== defaultCurrency || (!!targetCurrencyAmount && targetCurrencyAmount !== 1)" v-model="showState">
        <one-column-dialog-layout :disabled="submitting"
                                  :title="tt('Update User Custom Exchange Rate')" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #toolbar>
                <v-btn class="mx-2" density="comfortable" variant="outlined"
                       :disabled="submitting || !defaultCurrencyAmount || !currency || !targetCurrencyAmount"
                       @click="confirm">
                    {{ tt('OK') }}
                    <v-progress-circular indeterminate size="22" class="ms-2" v-if="submitting"></v-progress-circular>
                </v-btn>
            </template>

            <template #content>
                <v-form class="mt-2">
                    <v-row>
                        <v-col cols="12" md="6">
                            <number-input :autofocus="true"
                                          :disabled="submitting"
                                          :label="tt('Amount')"
                                          :placeholder="tt('Amount')"
                                          :persistent-placeholder="true"
                                          :min-value="USER_CUSTOM_EXCHANGE_RATE_MIN_VALUE"
                                          :max-value="USER_CUSTOM_EXCHANGE_RATE_MAX_VALUE"
                                          :max-decimal-count="4"
                                          v-model="defaultCurrencyAmount"
                                          @keyup.enter="targetAmountInput?.focus()" />
                        </v-col>
                        <v-col cols="12" md="6">
                            <currency-select :disabled="true"
                                             :label="tt('Currency')"
                                             :placeholder="tt('Currency')"
                                             v-model="defaultCurrency" />
                        </v-col>
                        <v-col cols="12" class="text-center">
                            <v-icon :icon="mdiSwapVertical" size="24" />
                        </v-col>
                        <v-col cols="12" md="6">
                            <number-input ref="targetAmountInput" :disabled="submitting"
                                          :label="tt('Amount')"
                                          :placeholder="tt('Amount')"
                                          :persistent-placeholder="true"
                                          :min-value="USER_CUSTOM_EXCHANGE_RATE_MIN_VALUE"
                                          :max-value="USER_CUSTOM_EXCHANGE_RATE_MAX_VALUE"
                                          :max-decimal-count="4"
                                          v-model="targetCurrencyAmount"
                                          @keyup.enter="confirm" />
                        </v-col>
                        <v-col cols="12" md="6">
                            <currency-select :disabled="submitting"
                                             :label="tt('Currency')"
                                             :placeholder="tt('Currency')"
                                             v-model="currency" />
                        </v-col>
                    </v-row>
                </v-form>
            </template>
        </one-column-dialog-layout>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import NumberInput from '@/components/desktop/NumberInput.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import {
    USER_CUSTOM_EXCHANGE_RATE_MAX_VALUE,
    USER_CUSTOM_EXCHANGE_RATE_MIN_VALUE
} from '@/consts/exchange_rate.ts';

import {
    mdiSwapVertical
} from '@mdi/js';

interface UserCustomExchangeRateUpdateResponse {
    message: string;
}

type NumberInputType = InstanceType<typeof NumberInput>;
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

const targetAmountInput = useTemplateRef<NumberInputType>('targetAmountInput');
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
    if (submitting.value || !defaultCurrencyAmount.value || !currency.value || !targetCurrencyAmount.value) {
        return;
    }

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
