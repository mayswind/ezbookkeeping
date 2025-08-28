<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <v-layout>
                    <v-navigation-drawer :permanent="alwaysShowNav" v-model="showNav">
                        <div class="mx-6 my-4">
                            <span class="text-subtitle-2">{{ tt('Data source') }}</span>
                            <p class="text-body-1 mt-1 mb-3">
                                <a tabindex="-1" target="_blank" :href="exchangeRatesData.referenceUrl" v-if="!loading && exchangeRatesData && !isUserCustomExchangeRates && exchangeRatesData.referenceUrl">{{ exchangeRatesData.dataSource }}</a>
                                <span v-else-if="!loading && exchangeRatesData && !isUserCustomExchangeRates && !exchangeRatesData.referenceUrl">{{ exchangeRatesData.dataSource }}</span>
                                <span v-else-if="!loading && exchangeRatesData && isUserCustomExchangeRates">{{ tt('User Custom') }}</span>
                                <span v-else-if="!loading && !exchangeRatesData">{{ tt('None') }}</span>
                                <span v-else-if="loading">
                                    <v-skeleton-loader class="skeleton-no-margin mt-3 mb-4" type="text" :loading="true"></v-skeleton-loader>
                                </span>
                            </p>
                            <span class="text-subtitle-2" v-if="exchangeRatesDataUpdateTime || loading">{{ tt('Last Updated') }}</span>
                            <p class="text-body-1 mt-1" v-if="exchangeRatesDataUpdateTime || loading">
                                <span v-if="!loading">{{ exchangeRatesDataUpdateTime }}</span>
                                <span v-if="loading">
                                    <v-skeleton-loader class="skeleton-no-margin mt-3 mb-4" type="text" :loading="true"></v-skeleton-loader>
                                </span>
                            </p>
                        </div>
                        <v-divider />
                        <div class="mx-6 mt-4">
                            <span class="text-subtitle-2">{{ tt('Base Amount') }}</span>
                            <amount-input class="mt-2" density="compact"
                                          :currency="baseCurrency"
                                          :disabled="loading || !exchangeRatesData || !exchangeRatesData.exchangeRates || !exchangeRatesData.exchangeRates.length"
                                          v-model="baseAmount"/>
                        </div>
                        <div class="mx-6 mt-4">
                            <span class="text-subtitle-2">{{ tt('Base Currency') }}</span>
                        </div>
                        <v-tabs show-arrows class="mb-4" direction="vertical"
                                :disabled="loading" v-model="baseCurrency"
                                v-if="exchangeRatesData && exchangeRatesData.exchangeRates && exchangeRatesData.exchangeRates.length">
                            <v-tab class="tab-text-truncate" :key="exchangeRate.currencyCode" :value="exchangeRate.currencyCode"
                                   v-for="exchangeRate in availableExchangeRates">
                                <div class="d-flex w-100">
                                    <span class="d-block text-truncate">{{ exchangeRate.currencyDisplayName }}</span>
                                    <small class="smaller ms-1">{{ exchangeRate.currencyCode }}</small>
                                </div>
                            </v-tab>
                        </v-tabs>
                        <div class="mx-6 mt-2 mb-4"
                             v-else-if="!exchangeRatesData || !exchangeRatesData.exchangeRates || !exchangeRatesData.exchangeRates.length">
                            <span v-if="!loading">{{ tt('None') }}</span>
                            <span v-else-if="loading">
                                <v-skeleton-loader class="skeleton-no-margin pt-2 pb-5" type="text"
                                                   :key="itemIdx" :loading="loading"
                                                   v-for="itemIdx in [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]"></v-skeleton-loader>
                            </span>
                        </div>
                    </v-navigation-drawer>
                    <v-main>
                        <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container" v-model="activeTab">
                            <v-window-item value="exchangeRatesPage">
                                <v-card variant="flat" min-height="680">
                                    <template #title>
                                        <div class="title-and-toolbar d-flex align-center">
                                            <v-btn class="me-3 d-md-none" density="compact" color="default" variant="plain"
                                                   :ripple="false" :icon="true" @click="showNav = !showNav">
                                                <v-icon :icon="mdiMenu" size="24" />
                                            </v-btn>
                                            <span>{{ tt('Exchange Rates Data') }}</span>
                                            <v-btn class="ms-3" color="default" variant="outlined"
                                                   :disabled="loading" @click="update"
                                                   v-if="isUserCustomExchangeRates">{{ tt('Update') }}</v-btn>
                                            <v-btn density="compact" color="default" variant="text" size="24"
                                                   class="ms-2" :icon="true" :loading="loading" @click="reload(true)">
                                                <template #loader>
                                                    <v-progress-circular indeterminate size="20"/>
                                                </template>
                                                <v-icon :icon="mdiRefresh" size="24" />
                                                <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                                            </v-btn>
                                        </div>
                                    </template>

                                    <v-table class="exchange-rates-table table-striped" :hover="!loading">
                                        <thead>
                                        <tr>
                                            <th>
                                                <div class="d-flex align-center">
                                                    <span>{{ tt('Currency') }}</span>
                                                    <v-spacer/>
                                                    <span>{{ tt('Amount') }}</span>
                                                </div>
                                            </th>
                                        </tr>
                                        </thead>

                                        <tbody>
                                        <tr :key="itemIdx"
                                            v-for="itemIdx in (loading && (!exchangeRatesData || !exchangeRatesData.exchangeRates || exchangeRatesData.exchangeRates.length < 1) ? [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12 ] : [])">
                                            <td class="px-0">
                                                <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                                            </td>
                                        </tr>

                                        <tr v-if="!loading && (!exchangeRatesData || !exchangeRatesData.exchangeRates || !exchangeRatesData.exchangeRates.length)">
                                            <td>{{ tt('No exchange rates data') }}</td>
                                        </tr>

                                        <tr class="exchange-rates-table-row-data" :key="exchangeRate.currencyCode"
                                            v-for="exchangeRate in availableExchangeRates">
                                            <td>
                                                <div class="d-flex align-center">
                                                    <span class="text-sm">{{ exchangeRate.currencyDisplayName }}</span>
                                                    <span class="text-caption ms-1">{{ exchangeRate.currencyCode }}</span>

                                                    <v-spacer/>

                                                    <v-btn class="px-2 ms-2" color="default"
                                                           density="comfortable" variant="text"
                                                           :class="{ 'd-none': loading, 'hover-display': !loading }"
                                                           v-if="exchangeRate.currencyCode !== baseCurrency"
                                                           @click="setAsBaseline(exchangeRate.currencyCode, getFinalConvertedAmount(exchangeRate, false))">
                                                        {{ tt('Set as Base') }}
                                                    </v-btn>
                                                    <v-btn class="px-2" color="default"
                                                           density="comfortable" variant="text"
                                                           :class="{ 'd-none': loading, 'hover-display': !loading }"
                                                           :prepend-icon="mdiDeleteOutline"
                                                           :loading="customExchangeRateRemoving[exchangeRate.currencyCode]"
                                                           :disabled="loading || updating"
                                                           v-if="exchangeRate.currencyCode !== defaultCurrency && isUserCustomExchangeRates"
                                                           @click="remove(exchangeRate.currencyCode)">
                                                        <template #loader>
                                                            <v-progress-circular indeterminate size="20" width="2"/>
                                                        </template>
                                                        {{ tt('Delete') }}
                                                    </v-btn>
                                                    <span class="ms-3">{{ getFinalConvertedAmount(exchangeRate, true) }}</span>
                                                </div>
                                            </td>
                                        </tr>
                                        </tbody>
                                    </v-table>
                                </v-card>
                            </v-window-item>
                        </v-window>
                    </v-main>
                </v-layout>
            </v-card>
        </v-col>
    </v-row>

    <update-dialog ref="updateDialog" />

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';
import UpdateDialog from './list/dialogs/UpdateDialog.vue';

import { ref, computed, useTemplateRef, watch } from 'vue';
import { useDisplay } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';
import { useExchangeRatesPageBase } from '@/views/base/ExchangeRatesPageBase.ts';

import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import { NumeralSystem } from '@/core/numeral.ts';

import type { LocalizedLatestExchangeRate } from '@/models/exchange_rate.ts';

import logger from '@/lib/logger.ts';

import {
    mdiRefresh,
    mdiMenu,
    mdiDeleteOutline
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;
type UpdateDialogType = InstanceType<typeof UpdateDialog>;

const { mdAndUp } = useDisplay();

const { tt, getCurrentNumeralSystemType, formatExchangeRateAmountToWesternArabicNumerals } = useI18n();
const {
    baseCurrency,
    baseAmount,
    defaultCurrency,
    exchangeRatesData,
    isUserCustomExchangeRates,
    exchangeRatesDataUpdateTime,
    availableExchangeRates,
    getConvertedAmount,
    setAsBaseline
} = useExchangeRatesPageBase();

const exchangeRatesStore = useExchangeRatesStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');
const updateDialog = useTemplateRef<UpdateDialogType>('updateDialog');

const activeTab = ref<string>('exchangeRatesPage');
const loading = ref<boolean>(true);
const updating = ref<boolean>(false);
const customExchangeRateRemoving = ref<Record<string, boolean>>({});
const alwaysShowNav = ref<boolean>(mdAndUp.value);
const showNav = ref<boolean>(mdAndUp.value);

const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());

function reload(force: boolean): void {
    loading.value = true;

    exchangeRatesStore.getLatestExchangeRates({
        silent: false,
        force: force
    }).then(() => {
        loading.value = false;

        if (exchangeRatesData.value && exchangeRatesData.value.exchangeRates) {
            const exchangeRates = exchangeRatesData.value.exchangeRates;
            let foundDefaultCurrency = false;

            for (let i = 0; i < exchangeRates.length; i++) {
                const exchangeRate = exchangeRates[i];
                if (exchangeRate.currency === baseCurrency.value) {
                    foundDefaultCurrency = true;
                    break;
                }
            }

            if (force) {
                snackbar.value?.showMessage('Exchange rates data has been updated');
            } else if (!foundDefaultCurrency) {
                snackbar.value?.showMessage('There is no exchange rates data for your default currency');
            }
        }
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function update(): void {
    updateDialog.value?.open().then(result => {
        if (result && result.message) {
            snackbar.value?.showMessage(result.message);
        }
    }).catch(error => {
        if (error) {
            snackbar.value?.showError(error);
        }
    });
}

function remove(currency: string): void {
    confirmDialog.value?.open('Are you sure you want to delete this user custom exchange rate?').then(() => {
        updating.value = true;
        customExchangeRateRemoving.value[currency] = true;

        exchangeRatesStore.deleteUserCustomExchangeRate({
            currency: currency
        }).then(() => {
            if (currency === baseCurrency.value) {
                baseCurrency.value = defaultCurrency.value;
            }

            updating.value = false;
            customExchangeRateRemoving.value[currency] = false;
        }).catch(error => {
            updating.value = false;
            customExchangeRateRemoving.value[currency] = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    });
}

function getFinalConvertedAmount(toExchangeRate: LocalizedLatestExchangeRate, displayLocalizedDigits: boolean): string {
    if (!baseCurrency.value) {
        if (displayLocalizedDigits) {
            return numeralSystem.value.digitZero;
        } else {
            return NumeralSystem.WesternArabicNumerals.digitZero;
        }
    }

    const fromExchangeRate = exchangeRatesStore.latestExchangeRateMap[baseCurrency.value];
    let exchangeRateAmount: number | '' | null = 0;

    try {
        exchangeRateAmount = getConvertedAmount(baseAmount.value / 100, fromExchangeRate, toExchangeRate);
    } catch (ex) {
        exchangeRateAmount = 0;
        logger.warn('failed to convert amount by exchange rates, original base amount is ' + baseAmount.value, ex)
    }

    if (!exchangeRateAmount) {
        if (displayLocalizedDigits) {
            return numeralSystem.value.digitZero;
        } else {
            return NumeralSystem.WesternArabicNumerals.digitZero;
        }
    }

    let ret = formatExchangeRateAmountToWesternArabicNumerals(exchangeRateAmount);

    if (displayLocalizedDigits) {
        ret = numeralSystem.value.replaceWesternArabicDigitsToLocalizedDigits(ret);
    }

    return ret;
}

watch(mdAndUp, (newValue) => {
    alwaysShowNav.value = newValue;

    if (!showNav.value) {
        showNav.value = newValue;
    }
});

reload(false);
</script>

<style>
.exchange-rates-table tr.exchange-rates-table-row-data .hover-display {
    display: none;
}

.exchange-rates-table tr.exchange-rates-table-row-data:hover .hover-display {
    display: grid;
}

</style>
