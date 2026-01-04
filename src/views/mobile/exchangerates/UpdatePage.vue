<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Update User Custom Exchange Rate')"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="checkmark_alt"
                         :class="{ 'disabled': submitting || !defaultCurrencyAmount || !currency || !targetCurrencyAmount }"
                         @click="confirm"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list form strong inset dividers class="margin-vertical">
            <template #list>
                <list-number-input
                    :disabled="submitting"
                    :label="tt('Amount')"
                    :placeholder="tt('Amount')"
                    :min-value="USER_CUSTOM_EXCHANGE_RATE_MIN_VALUE"
                    :max-value="USER_CUSTOM_EXCHANGE_RATE_MAX_VALUE"
                    :max-decimal-count="4"
                    v-model="defaultCurrencyAmount"
                ></list-number-input>
            </template>

            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                link="#"
                :class="{ 'disabled': true }"
                :header="tt('Currency')"
                :no-chevron="true"
            >
                <template #title>
                    <div class="no-padding no-margin">
                        <span>{{ getCurrencyName(defaultCurrency) }}&nbsp;</span>
                        <small class="smaller">{{ defaultCurrency }}</small>
                    </div>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-block class="display-flex justify-content-center full-line margin-vertical">
            <f7-icon class="separate-icon" f7="arrow_up_arrow_down"></f7-icon>
        </f7-block>

        <f7-list form strong inset dividers class="margin-vertical">
            <template #list>
                <list-number-input
                    :disabled="submitting"
                    :label="tt('Amount')"
                    :placeholder="tt('Amount')"
                    :min-value="USER_CUSTOM_EXCHANGE_RATE_MIN_VALUE"
                    :max-value="USER_CUSTOM_EXCHANGE_RATE_MAX_VALUE"
                    :max-decimal-count="4"
                    v-model="targetCurrencyAmount"
                ></list-number-input>
            </template>

            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                link="#"
                :class="{ 'disabled': submitting }"
                :header="tt('Currency')"
                @click="showCurrencyPopup = true"
            >
                <template #title>
                    <div class="no-padding no-margin">
                        <span>{{ getCurrencyName(currency) }}&nbsp;</span>
                        <small class="smaller">{{ currency }}</small>
                    </div>
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="currencyCode" value-field="currencyCode"
                                           title-field="displayName" after-field="currencyCode"
                                           :title="tt('Currency Name')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Currency')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allCurrencies"
                                           v-model:show="showCurrencyPopup"
                                           v-model="currency">
                </list-item-selection-popup>
            </f7-list-item>
        </f7-list>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';

import { useUserStore } from '@/stores/user.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import type { LocalizedCurrencyInfo } from '@/core/currency.ts';

import {
    USER_CUSTOM_EXCHANGE_RATE_MAX_VALUE,
    USER_CUSTOM_EXCHANGE_RATE_MIN_VALUE
} from '@/consts/exchange_rate.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt, getAllCurrencies, getCurrencyName } = useI18n();
const { showToast } = useI18nUIComponents();

const userStore = useUserStore();
const exchangeRatesStore = useExchangeRatesStore();

const submitting = ref<boolean>(false);
const defaultCurrency = ref<string>(userStore.currentUserDefaultCurrency);
const defaultCurrencyAmount = ref<number>(1);
const currency = ref<string>(userStore.currentUserDefaultCurrency);
const targetCurrencyAmount = ref<number>(1);
const showCurrencyPopup = ref<boolean>(false);

const allCurrencies = computed<LocalizedCurrencyInfo[]>(() => getAllCurrencies());

function init(): void {
    defaultCurrencyAmount.value = 1;
    currency.value = userStore.currentUserDefaultCurrency;
    targetCurrencyAmount.value = 1;
}

function confirm(): void {
    const router = props.f7router;

    submitting.value = true;
    showLoading(() => submitting.value);

    exchangeRatesStore.updateUserCustomExchangeRate({
        currency: currency.value,
        rate: targetCurrencyAmount.value / defaultCurrencyAmount.value
    }).then(() => {
        submitting.value = false;
        hideLoading();

        showToast('You have updated exchange rate');
        router.back();
    }).catch(error => {
        submitting.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

init();
</script>
