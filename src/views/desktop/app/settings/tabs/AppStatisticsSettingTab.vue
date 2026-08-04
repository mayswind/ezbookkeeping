<template>
    <v-row>
        <v-col cols="12">
            <v-card :title="tt('Statistics Settings')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Default Chart Data Type')"
                                    :placeholder="tt('Default Chart Data Type')"
                                    :items="allChartDataTypes"
                                    v-model="defaultChartDataType"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Timezone Used for Date Range')"
                                    :placeholder="tt('Timezone Used for Date Range')"
                                    :items="allTimezoneTypesUsedForStatistics"
                                    v-model="defaultTimezoneType"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Default Keyword Search Matching Mode')"
                                    :placeholder="tt('Default Keyword Search Matching Mode')"
                                    :items="allKeywordMatchModes"
                                    v-model="defaultKeywordMatchMode"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-text-field
                                    class="always-cursor-pointer"
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :loading="loadingAccounts"
                                    :readonly="true"
                                    :disabled="!hasAnyAccount"
                                    :label="tt('Default Account Filter')"
                                    :placeholder="tt('Default Account Filter')"
                                    :model-value="defaultAccountFilterDisplayContent"
                                    @click="showFilterAccountDialog = true"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-text-field
                                    class="always-cursor-pointer"
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :loading="loadingTransactionCategories"
                                    :readonly="true"
                                    :disabled="!hasAnyTransactionCategory"
                                    :label="tt('Default Transaction Category Filter')"
                                    :placeholder="tt('Default Transaction Category Filter')"
                                    :model-value="defaultTransactionCategoryFilterDisplayContent"
                                    @click="showFilterCategoryDialog = true"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Default Sort Order')"
                                    :placeholder="tt('Default Sort Order')"
                                    :items="allSortingTypes"
                                    v-model="defaultSortingType"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :title="tt('Categorical Analysis Settings')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Default Chart Type')"
                                    :placeholder="tt('Default Chart Type')"
                                    :items="allCategoricalChartTypes"
                                    v-model="defaultCategoricalChartType"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Default Date Range')"
                                    :placeholder="tt('Default Date Range')"
                                    :items="allCategoricalChartDateRanges"
                                    v-model="defaultCategoricalChartDateRange"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :title="tt('Trend Analysis Settings')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Default Chart Type')"
                                    :placeholder="tt('Default Chart Type')"
                                    :items="allTrendChartTypes"
                                    v-model="defaultTrendChartType"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Default Date Range')"
                                    :placeholder="tt('Default Date Range')"
                                    :items="allTrendChartDateRanges"
                                    v-model="defaultTrendChartDateRange"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :title="tt('Asset Trends Settings')">
                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Default Chart Type')"
                                    :placeholder="tt('Default Chart Type')"
                                    :items="allTrendChartTypes"
                                    v-model="defaultAssetTrendsChartType"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :label="tt('Default Date Range')"
                                    :placeholder="tt('Default Date Range')"
                                    :items="allAssetTrendsChartDateRanges"
                                    v-model="defaultAssetTrendsChartDateRange"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>
    </v-row>

    <account-filter-settings-dialog type="statisticsDefault"
                                    v-model:show="showFilterAccountDialog"
                                    @settings:change="showFilterAccountDialog = false" />


    <category-filter-settings-dialog type="statisticsDefault"
                                     v-model:show="showFilterCategoryDialog"
                                     @settings:change="showFilterCategoryDialog = false" />

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';
import AccountFilterSettingsDialog from '@/views/desktop/common/dialogs/AccountFilterSettingsDialog.vue';
import CategoryFilterSettingsDialog from '@/views/desktop/common/dialogs/CategoryFilterSettingsDialog.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useStatisticsSettingPageBase } from '@/views/base/statistics/StatisticsSettingPageBase.ts';

import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';

import { isObjectEmpty } from '@/lib/common.ts';

type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();
const {
    loadingAccounts,
    loadingTransactionCategories,
    allChartDataTypes,
    allTimezoneTypesUsedForStatistics,
    allKeywordMatchModes,
    allSortingTypes,
    allCategoricalChartTypes,
    allCategoricalChartDateRanges,
    allTrendChartTypes,
    allTrendChartDateRanges,
    allAssetTrendsChartDateRanges,
    defaultChartDataType,
    defaultTimezoneType,
    defaultKeywordMatchMode,
    defaultAccountFilterDisplayContent,
    defaultTransactionCategoryFilterDisplayContent,
    defaultSortingType,
    defaultCategoricalChartType,
    defaultCategoricalChartDateRange,
    defaultTrendChartType,
    defaultTrendChartDateRange,
    defaultAssetTrendsChartType,
    defaultAssetTrendsChartDateRange
} = useStatisticsSettingPageBase();

const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showFilterAccountDialog = ref<boolean>(false);
const showFilterCategoryDialog = ref<boolean>(false);

const hasAnyAccount = computed<boolean>(() => accountsStore.allPlainAccounts.length > 0);
const hasAnyTransactionCategory = computed<boolean>(() => !isObjectEmpty(transactionCategoriesStore.allTransactionCategoriesMap));

function init(): void {
    loadingAccounts.value = true;

    accountsStore.loadAllAccounts({
        force: false
    }).then(() => {
        loadingAccounts.value = false;
    }).catch(error => {
        loadingAccounts.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });

    transactionCategoriesStore.loadAllCategories({
        force: false
    }).then(() => {
        loadingTransactionCategories.value = false;
    }).catch(error => {
        loadingTransactionCategories.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

init();
</script>

