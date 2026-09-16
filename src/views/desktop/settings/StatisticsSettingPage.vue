<template>
    <v-row>
        <v-col cols="12">
            <v-card class="setting-items" :title="tt('Common Settings')">
                <v-form>
                    <v-card-text class="pa-0 text-body-medium">
                        <div class="setting-item">
                            <span>{{ tt('Default Chart Data Type') }}</span>
                            <v-spacer/>
                            <v-select
                                class="ms-4"
                                density="compact"
                                item-title="displayName"
                                item-value="type"
                                persistent-placeholder
                                max-width="400px"
                                :placeholder="tt('Default Chart Data Type')"
                                :items="allChartDataTypes"
                                v-model="defaultChartDataType"
                            />
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Timezone Used for Date Range') }}</span>
                            <v-spacer/>
                            <v-select
                                class="ms-4"
                                density="compact"
                                item-title="displayName"
                                item-value="type"
                                persistent-placeholder
                                max-width="400px"
                                :placeholder="tt('Timezone Used for Date Range')"
                                :items="allTimezoneTypesUsedForStatistics"
                                v-model="defaultTimezoneType"
                            />
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Default Keyword Search Matching Mode') }}</span>
                            <v-spacer/>
                            <v-select
                                class="ms-4"
                                density="compact"
                                item-title="displayName"
                                item-value="type"
                                persistent-placeholder
                                max-width="400px"
                                :placeholder="tt('Default Keyword Search Matching Mode')"
                                :items="allKeywordMatchModes"
                                v-model="defaultKeywordMatchMode"
                            />
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Default Account Filter') }}</span>
                            <v-spacer/>
                            <v-btn class="ms-4" variant="outlined" color="default"
                                   :disabled="!hasAnyAccount" :loading="loadingAccounts"
                                   @click="showFilterAccountDialog = true">
                                {{ defaultAccountFilterDisplayContent || tt('All') }}
                                <template #loader>
                                    <v-progress-circular indeterminate size="20"/>
                                </template>
                            </v-btn>
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Default Transaction Category Filter') }}</span>
                            <v-spacer/>
                            <v-btn class="ms-4" variant="outlined" color="default"
                                   :disabled="!hasAnyTransactionCategory" :loading="loadingTransactionCategories"
                                   @click="showFilterCategoryDialog = true">
                                {{ defaultTransactionCategoryFilterDisplayContent || tt('All') }}
                                <template #loader>
                                    <v-progress-circular indeterminate size="20"/>
                                </template>
                            </v-btn>
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Default Sort Order') }}</span>
                            <v-spacer/>
                            <v-select
                                class="ms-4"
                                density="compact"
                                item-title="displayName"
                                item-value="type"
                                persistent-placeholder
                                max-width="400px"
                                :placeholder="tt('Default Sort Order')"
                                :items="allSortingTypes"
                                v-model="defaultSortingType"
                            />
                        </div>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card class="setting-items" :title="tt('Categorical Analysis Settings')">
                <v-form>
                    <v-card-text class="pa-0 text-body-medium">
                        <div class="setting-item">
                            <span>{{ tt('Default Chart Type') }}</span>
                            <v-spacer/>
                            <v-select
                                class="ms-4"
                                density="compact"
                                item-title="displayName"
                                item-value="type"
                                persistent-placeholder
                                max-width="400px"
                                :placeholder="tt('Default Chart Type')"
                                :items="allCategoricalChartTypes"
                                v-model="defaultCategoricalChartType"
                            />
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Default Date Range') }}</span>
                            <v-spacer/>
                            <v-select
                                class="ms-4"
                                density="compact"
                                item-title="displayName"
                                item-value="type"
                                persistent-placeholder
                                max-width="400px"
                                :placeholder="tt('Default Date Range')"
                                :items="allCategoricalChartDateRanges"
                                v-model="defaultCategoricalChartDateRange"
                            />
                        </div>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card class="setting-items" :title="tt('Trend Analysis Settings')">
                <v-form>
                    <v-card-text class="pa-0 text-body-medium">
                        <div class="setting-item">
                            <span>{{ tt('Default Chart Type') }}</span>
                            <v-spacer/>
                            <v-select
                                class="ms-4"
                                density="compact"
                                item-title="displayName"
                                item-value="type"
                                persistent-placeholder
                                max-width="400px"
                                :placeholder="tt('Default Chart Type')"
                                :items="allTrendChartTypes"
                                v-model="defaultTrendChartType"
                            />
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Default Date Range') }}</span>
                            <v-spacer/>
                            <v-select
                                class="ms-4"
                                density="compact"
                                item-title="displayName"
                                item-value="type"
                                persistent-placeholder
                                max-width="400px"
                                :placeholder="tt('Default Date Range')"
                                :items="allTrendChartDateRanges"
                                v-model="defaultTrendChartDateRange"
                            />
                        </div>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card class="setting-items" :title="tt('Asset Trends Settings')">
                <v-form>
                    <v-card-text class="pa-0 text-body-medium">
                        <div class="setting-item">
                            <span>{{ tt('Default Chart Type') }}</span>
                            <v-spacer/>
                            <v-select
                                class="ms-4"
                                density="compact"
                                item-title="displayName"
                                item-value="type"
                                persistent-placeholder
                                max-width="400px"
                                :placeholder="tt('Default Chart Type')"
                                :items="allTrendChartTypes"
                                v-model="defaultAssetTrendsChartType"
                            />
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Default Date Range') }}</span>
                            <v-spacer/>
                            <v-select
                                class="ms-4"
                                density="compact"
                                item-title="displayName"
                                item-value="type"
                                persistent-placeholder
                                max-width="400px"
                                :placeholder="tt('Default Date Range')"
                                :items="allAssetTrendsChartDateRanges"
                                v-model="defaultAssetTrendsChartDateRange"
                            />
                        </div>
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
    loadingTransactionCategories.value = true;

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

