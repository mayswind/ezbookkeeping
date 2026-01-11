<template>
    <v-dialog width="800" :persistent="submitting || (!!defaultCurrencyAmount && defaultCurrencyAmount !== 1) || currency !== defaultCurrency || (!!targetCurrencyAmount && targetCurrencyAmount !== 1)" v-model="showState">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <div class="d-flex align-center">
                    <h4 class="text-h4">{{ tt('Update User Custom Exchange Rate') }}</h4>
                </div>
            </template>
            <v-card-text class="d-flex flex-column flex-md-row mt-2">
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
            </v-card-text>
            <v-card-text>
                <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                    <v-btn :disabled="submitting || !defaultCurrencyAmount || !currency || !targetCurrencyAmount" @click="confirm">
                        {{ tt('OK') }}
                        <v-progress-circular indeterminate size="22" class="ms-2" v-if="submitting"></v-progress-circular>
                    </v-btn>
                    <v-btn color="secondary" variant="tonal" :disabled="submitting" @click="cancel">{{ tt('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
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
