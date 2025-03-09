<template>
    <v-row>
        <v-col cols="12">
            <v-card :title="tt('Basic Settings')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="name"
                                    item-value="value"
                                    persistent-placeholder
                                    :label="tt('Theme')"
                                    :placeholder="tt('Theme')"
                                    :items="allThemes"
                                    v-model="currentTheme"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-autocomplete
                                    item-title="displayNameWithUtcOffset"
                                    item-value="name"
                                    auto-select-first
                                    persistent-placeholder
                                    :label="tt('Timezone')"
                                    :placeholder="tt('Timezone')"
                                    :items="allTimezones"
                                    :no-data-text="tt('No results')"
                                    v-model="timeZone"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="value"
                                    persistent-placeholder
                                    :label="tt('Auto-update Exchange Rates Data')"
                                    :placeholder="tt('Auto-update Exchange Rates Data')"
                                    :items="enableDisableOptions"
                                    v-model="isAutoUpdateExchangeRatesData"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="value"
                                    persistent-placeholder
                                    :label="tt('Show Account Balance')"
                                    :placeholder="tt('Show Account Balance')"
                                    :items="enableDisableOptions"
                                    v-model="showAccountBalance"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :title="tt('Navigation Bar')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="value"
                                    persistent-placeholder
                                    :label="tt('Show Add Transaction Button')"
                                    :placeholder="tt('Show Add Transaction Button')"
                                    :items="enableDisableOptions"
                                    v-model="showAddTransactionButtonInDesktopNavbar"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :title="tt('Overview Page')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="value"
                                    persistent-placeholder
                                    :label="tt('Show Amount')"
                                    :placeholder="tt('Show Amount')"
                                    :items="enableDisableOptions"
                                    v-model="showAmountInHomePage"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Timezone Used for Statistics')"
                                    :placeholder="tt('Timezone Used for Statistics')"
                                    :items="allTimezoneTypesUsedForStatistics"
                                    v-model="timezoneUsedForStatisticsInHomePage"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :title="tt('Transaction List Page')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    persistent-placeholder
                                    :label="tt('Transactions Per Page')"
                                    :placeholder="tt('Transactions Per Page')"
                                    :items="[ 5, 10, 15, 20, 25, 30, 50 ]"
                                    v-model="itemsCountInTransactionListPage"
                                />
                            </v-col>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="value"
                                    persistent-placeholder
                                    :label="tt('Show Monthly Total Amount')"
                                    :placeholder="tt('Show Monthly Total Amount')"
                                    :items="enableDisableOptions"
                                    v-model="showTotalAmountInTransactionListPage"
                                />
                            </v-col>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="value"
                                    persistent-placeholder
                                    :label="tt('Show Transaction Tag')"
                                    :placeholder="tt('Show Transaction Tag')"
                                    :items="enableDisableOptions"
                                    v-model="showTagInTransactionListPage"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :title="tt('Transaction Edit Page')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="name"
                                    item-value="value"
                                    persistent-placeholder
                                    :label="tt('Automatically Save Draft')"
                                    :placeholder="tt('Automatically Save Draft')"
                                    :items="allAutoSaveTransactionDraftTypes"
                                    v-model="autoSaveTransactionDraft"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="value"
                                    persistent-placeholder
                                    :label="tt('Automatically Add Geolocation')"
                                    :placeholder="tt('Automatically Add Geolocation')"
                                    :items="enableDisableOptions"
                                    v-model="isAutoGetCurrentGeoLocation"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :title="tt('Exchange Rates Data Page')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Sort by')"
                                    :placeholder="tt('Sort by')"
                                    :items="allCurrencySortingTypes"
                                    v-model="currencySortByInExchangeRatesPage"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>
    </v-row>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';
import { useAppSettingPageBase } from '@/views/base/settings/AppSettingsPageBase.ts';

import { useSettingsStore } from '@/stores/setting.ts';

import type { LocalizedSwitchOption } from '@/core/base.ts';
import { ThemeType } from '@/core/theme.ts';
import { getSystemTheme } from '@/lib/ui/common.ts';

const theme = useTheme();

const { tt, getAllEnableDisableOptions } = useI18n();
const {
    allThemes,
    allTimezones,
    allTimezoneTypesUsedForStatistics,
    allCurrencySortingTypes,
    allAutoSaveTransactionDraftTypes,
    timeZone,
    isAutoUpdateExchangeRatesData,
    showAccountBalance,
    showAmountInHomePage,
    itemsCountInTransactionListPage,
    timezoneUsedForStatisticsInHomePage,
    showTotalAmountInTransactionListPage,
    showTagInTransactionListPage,
    autoSaveTransactionDraft,
    isAutoGetCurrentGeoLocation,
    currencySortByInExchangeRatesPage
} = useAppSettingPageBase();

const settingsStore = useSettingsStore();

const enableDisableOptions = computed<LocalizedSwitchOption[]>(() => getAllEnableDisableOptions());

const currentTheme = computed<string>({
    get: () => settingsStore.appSettings.theme,
    set: (value: string) => {
        if (value !== settingsStore.appSettings.theme) {
            settingsStore.setTheme(value);

            if (value === ThemeType.Light || value === ThemeType.Dark) {
                theme.global.name.value = value;
            } else {
                theme.global.name.value = getSystemTheme();
            }
        }
    }
});

const showAddTransactionButtonInDesktopNavbar = computed<boolean>({
    get: () => settingsStore.appSettings.showAddTransactionButtonInDesktopNavbar,
    set: (value) => settingsStore.setShowAddTransactionButtonInDesktopNavbar(value)
});
</script>
