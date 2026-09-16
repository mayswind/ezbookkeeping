<template>
    <v-row>
        <v-col cols="12">
            <v-card class="setting-items" :title="tt('General Settings')">
                <v-form>
                    <v-card-text class="pa-0 text-body-medium">
                        <div class="setting-item">
                            <span>{{ tt('Theme') }}</span>
                            <v-spacer/>
                            <v-select
                                class="ms-4"
                                density="compact"
                                item-title="name"
                                item-value="value"
                                persistent-placeholder
                                max-width="400px"
                                :placeholder="tt('Theme')"
                                :items="allThemes"
                                v-model="currentTheme"
                            />
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Timezone') }}</span>
                            <v-spacer/>
                            <v-autocomplete
                                class="ms-4"
                                density="compact"
                                item-title="displayNameWithUtcOffset"
                                item-value="name"
                                auto-select-first
                                persistent-placeholder
                                max-width="600px"
                                :placeholder="tt('Timezone')"
                                :items="allTimezones"
                                :no-data-text="tt('No results')"
                                v-model="timeZone"
                            />
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Show Account Balance') }}</span>
                            <v-spacer/>
                            <v-switch class="ms-4" v-model="showAccountBalance"/>
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Account Category Order') }}</span>
                            <v-spacer/>
                            <v-btn class="ms-4" variant="outlined" color="default"
                                   @click="accountCategorysDisplayOrderDialog?.open().catch(()=>{})">
                                {{ accountCategorysDisplayOrderContent }}
                            </v-btn>
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Chart Color Scheme') }}</span>
                            <v-spacer/>
                            <v-btn class="ms-4" variant="outlined" color="default"
                                   @click="chartColorSchemeDialog?.open().catch(()=>{})">
                                {{ chartColorSchemeContent }}
                            </v-btn>
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Auto-update Exchange Rates Data') }}</span>
                            <v-spacer/>
                            <v-switch class="ms-4" v-model="isAutoUpdateExchangeRatesData"/>
                        </div>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card class="setting-items" :title="tt('Navigation Bar')">
                <v-form>
                    <v-card-text class="pa-0 text-body-medium">
                        <div class="setting-item">
                            <span>{{ tt('Show Add Transaction Button') }}</span>
                            <v-spacer/>
                            <v-switch class="ms-4" v-model="showAddTransactionButtonInDesktopNavbar"/>
                        </div>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card class="setting-items" :title="tt('Overview Page')">
                <v-form>
                    <v-card-text class="pa-0 text-body-medium">
                        <div class="setting-item">
                            <span>{{ tt('Home Page Layout') }}</span>
                            <v-spacer/>
                            <v-btn class="ms-4" variant="outlined" color="default"
                                   @click="router.push('/overview/edit')">
                                {{ desktopOverviewPageLayoutDisplayContent }}
                            </v-btn>
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Show Amount') }}</span>
                            <v-spacer/>
                            <v-switch class="ms-4" v-model="showAmountInHomePage"/>
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Timezone Used for Statistics') }}</span>
                            <v-spacer/>
                            <v-select
                                class="ms-4"
                                density="compact"
                                item-title="displayName"
                                item-value="type"
                                persistent-placeholder
                                max-width="400px"
                                :placeholder="tt('Timezone Used for Statistics')"
                                :items="allTimezoneTypesUsedForStatistics"
                                v-model="timezoneUsedForStatisticsInHomePage"
                            />
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Accounts Included in Overview Statistics') }}</span>
                            <v-spacer/>
                            <v-btn class="ms-4" variant="outlined" color="default"
                                   :disabled="!hasAnyAccount" :loading="loadingAccounts"
                                   @click="showAccountsIncludedInHomePageOverviewDialog = true">
                                {{ accountsIncludedInHomePageOverviewDisplayContent || tt('All') }}
                                <template #loader>
                                    <v-progress-circular indeterminate size="20"/>
                                </template>
                            </v-btn>
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Transaction Categories Included in Overview Statistics') }}</span>
                            <v-spacer/>
                            <v-btn class="ms-4" variant="outlined" color="default"
                                   :disabled="!hasAnyTransactionCategory" :loading="loadingTransactionCategories"
                                   @click="showTransactionCategoriesIncludedInHomePageOverviewDialog = true">
                                {{ transactionCategoriesIncludedInHomePageOverviewDisplayContent || tt('All') }}
                                <template #loader>
                                    <v-progress-circular indeterminate size="20"/>
                                </template>
                            </v-btn>
                        </div>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card class="setting-items" :title="tt('Transaction List Page')">
                <v-form>
                    <v-card-text class="pa-0 text-body-medium">
                        <div class="setting-item">
                            <span>{{ tt('Transactions Per Page') }}</span>
                            <v-spacer/>
                            <v-select
                                class="ms-4"
                                density="compact"
                                item-title="name"
                                item-value="value"
                                persistent-placeholder
                                max-width="200px"
                                :placeholder="tt('Transactions Per Page')"
                                :items="allPageCounts"
                                v-model="itemsCountInTransactionListPage"
                            />
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Show Monthly Total Amount') }}</span>
                            <v-spacer/>
                            <v-switch class="ms-4" v-model="showTotalAmountInTransactionListPage"/>
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Show Transaction Tags') }}</span>
                            <v-spacer/>
                            <v-switch class="ms-4" v-model="showTagInTransactionListPage"/>
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
                                v-model="defaultKeywordMatchModeInTransactionListPage"
                            />
                        </div>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card class="setting-items" :title="tt('Transaction Edit Page')">
                <v-form>
                    <v-card-text class="pa-0 text-body-medium">
                        <div class="setting-item">
                            <span>{{ tt('Automatically Save Draft') }}</span>
                            <v-spacer/>
                            <v-select
                                class="ms-4"
                                density="compact"
                                item-title="name"
                                item-value="value"
                                persistent-placeholder
                                max-width="400px"
                                :placeholder="tt('Automatically Save Draft')"
                                :items="allAutoSaveTransactionDraftTypes"
                                v-model="autoSaveTransactionDraft"
                            />
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Automatically Add Geolocation') }}</span>
                            <v-spacer/>
                            <v-switch class="ms-4" v-model="isAutoGetCurrentGeoLocation"/>
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Transaction Picture Upload Quality') }}</span>
                            <v-spacer/>
                            <v-select
                                class="ms-4"
                                density="compact"
                                item-title="displayName"
                                item-value="type"
                                persistent-placeholder
                                max-width="400px"
                                :placeholder="tt('Transaction Picture Upload Quality')"
                                :items="allImageUploadQualityTypes"
                                v-model="transactionPictureQuality"
                            />
                        </div>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card class="setting-items" :title="tt('AI Clipboard Text Recognition')">
                <v-form>
                    <v-card-text class="pa-0 text-body-medium">
                        <div class="setting-item">
                            <span>{{ tt('Always Require Confirmation of Clipboard Content Before Submission') }}</span>
                            <v-spacer/>
                            <v-switch class="ms-4" v-model="isAlwaysRequireConfirmationOfClipboardContentBeforeSubmission"/>
                        </div>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card class="setting-items" :title="tt('AI Image Recognition')">
                <v-form>
                    <v-card-text class="pa-0 text-body-medium">
                        <div class="setting-item">
                            <span>{{ tt('Auto Upload AI Recognition Image as Transaction Picture') }}</span>
                            <v-spacer/>
                            <v-switch class="ms-4" v-model="isAutoUploadTransactionPictureForAIRecognition"/>
                        </div>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card class="setting-items" :title="tt('Import Transaction Dialog')">
                <v-form>
                    <v-card-text class="pa-0 text-body-medium">
                        <div class="setting-item">
                            <span>{{ tt('Remember Last Selected File Type') }}</span>
                            <v-spacer/>
                            <v-switch class="ms-4" v-model="rememberLastSelectedFileTypeInImportTransactionDialog"/>
                        </div>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card class="setting-items" :title="tt('Insights Explorer Page')">
                <v-form>
                    <v-card-text class="pa-0 text-body-medium">
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
                                :items="allInsightsExplorerDefaultDateRanges"
                                v-model="insightsExplorerDefaultDateRangeType"
                            />
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Show Transaction Tags') }}</span>
                            <v-spacer/>
                            <v-switch class="ms-4" v-model="showTagInInsightsExplorerPage"/>
                        </div>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card class="setting-items" :title="tt('Account List Page')">
                <v-form>
                    <v-card-text class="pa-0 text-body-medium">
                        <div class="setting-item">
                            <span>{{ tt('Accounts Included in Total') }}</span>
                            <v-spacer/>
                            <v-btn class="ms-4" variant="outlined" color="default"
                                   :disabled="!hasAnyVisibleAccount" :loading="loadingAccounts"
                                   @click="showAccountsIncludedInTotalDialog = true">
                                {{ accountsIncludedInTotalDisplayContent || tt('All') }}
                                <template #loader>
                                    <v-progress-circular indeterminate size="20"/>
                                </template>
                            </v-btn>
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Hide Categories Without Accounts') }}</span>
                            <v-spacer/>
                            <v-switch class="ms-4" v-model="hideCategoriesWithoutAccounts"/>
                        </div>

                        <div class="setting-item">
                            <span>{{ tt('Default Date Range for Reconciliation Statement Button') }}</span>
                            <v-spacer/>
                            <v-select
                                class="ms-4"
                                density="compact"
                                item-title="displayName"
                                item-value="type"
                                persistent-placeholder
                                max-width="400px"
                                :placeholder="tt('Default Date Range for Reconciliation Statement Button')"
                                :items="allReconciliationStatementDateRanges"
                                v-model="reconciliationStatementButtonDefaultDateRangeTypeInDesktop"
                            />
                        </div>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card class="setting-items" :title="tt('Exchange Rates Data Page')">
                <v-form>
                    <v-card-text class="pa-0 text-body-medium">
                        <div class="setting-item">
                            <span>{{ tt('Sort by') }}</span>
                            <v-spacer/>
                            <v-select
                                class="ms-4"
                                density="compact"
                                item-title="displayName"
                                item-value="type"
                                persistent-placeholder
                                max-width="400px"
                                :placeholder="tt('Sort by')"
                                :items="allCurrencySortingTypes"
                                v-model="currencySortByInExchangeRatesPage"
                            />
                        </div>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>
    </v-row>

    <account-filter-settings-dialog type="homePageOverview"
                                    v-model:show="showAccountsIncludedInHomePageOverviewDialog"
                                    @settings:change="showAccountsIncludedInHomePageOverviewDialog = false" />

    <category-filter-settings-dialog type="homePageOverview"
                                     :category-types="`${CategoryType.Income},${CategoryType.Expense}`"
                                     v-model:show="showTransactionCategoriesIncludedInHomePageOverviewDialog"
                                     @settings:change="showTransactionCategoriesIncludedInHomePageOverviewDialog = false" />

    <account-filter-settings-dialog type="accountListTotalAmount"
                                    v-model:show="showAccountsIncludedInTotalDialog"
                                    @settings:change="showAccountsIncludedInTotalDialog = false" />

    <chart-color-scheme-dialog ref="chartColorSchemeDialog" />
    <account-category-display-order-dialog ref="accountCategorysDisplayOrderDialog" />

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';
import AccountFilterSettingsDialog from '@/views/desktop/common/dialogs/AccountFilterSettingsDialog.vue';
import CategoryFilterSettingsDialog from '@/views/desktop/common/dialogs/CategoryFilterSettingsDialog.vue';
import ChartColorSchemeDialog from './dialogs/ChartColorSchemeDialog.vue';
import AccountCategoryDisplayOrderDialog from './dialogs/AccountCategoryDisplayOrderDialog.vue';

import { ref, computed, useTemplateRef } from 'vue';
import { useRouter } from 'vue-router';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';
import { useAppSettingPageBase } from '@/views/base/settings/AppSettingsPageBase.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';

import type { NameNumeralValue } from '@/core/base.ts';
import { ThemeType } from '@/core/theme.ts';
import { type LocalizedDateRange, DateRangeScene } from '@/core/datetime.ts';
import { CategoryType } from '@/core/category.ts';
import { DEFAULT_RECONCILIATION_STATEMENT_DATE_RANGE_IN_DESKTOP } from '@/core/statistics.ts';

import { DEFAULT_PAGE_COUNTS } from '@/consts/page.ts';

import { isDefaultDesktopOverviewLayout, parseDesktopOverviewLayout } from '@/lib/overview_layout.ts';
import { getSystemTheme } from '@/lib/ui/common.ts';

type SnackBarType = InstanceType<typeof SnackBar>;
type ChartColorSchemeDialogType = InstanceType<typeof ChartColorSchemeDialog>;
type AccountCategoryDisplayOrderDialogType = InstanceType<typeof AccountCategoryDisplayOrderDialog>;

const theme = useTheme();
const router = useRouter();

const { tt, getAllDateRanges, getTablePageOptions } = useI18n();
const {
    loadingAccounts,
    loadingTransactionCategories,
    allThemes,
    allTimezones,
    allTimezoneTypesUsedForStatistics,
    allCurrencySortingTypes,
    allKeywordMatchModes,
    allAutoSaveTransactionDraftTypes,
    allImageUploadQualityTypes,
    allReconciliationStatementDateRanges,
    hasAnyAccount,
    hasAnyVisibleAccount,
    hasAnyTransactionCategory,
    timeZone,
    isAutoUpdateExchangeRatesData,
    showAccountBalance,
    showAmountInHomePage,
    itemsCountInTransactionListPage,
    timezoneUsedForStatisticsInHomePage,
    showTotalAmountInTransactionListPage,
    showTagInTransactionListPage,
    defaultKeywordMatchModeInTransactionListPage,
    autoSaveTransactionDraft,
    isAutoGetCurrentGeoLocation,
    transactionPictureQuality,
    isAlwaysRequireConfirmationOfClipboardContentBeforeSubmission,
    isAutoUploadTransactionPictureForAIRecognition,
    currencySortByInExchangeRatesPage,
    chartColorSchemeContent,
    accountsIncludedInHomePageOverviewDisplayContent,
    accountsIncludedInTotalDisplayContent,
    accountCategorysDisplayOrderContent,
    transactionCategoriesIncludedInHomePageOverviewDisplayContent,
    getValidReconciliationStatementPageDefaultDateRangeType
} = useAppSettingPageBase();

const settingsStore = useSettingsStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');
const chartColorSchemeDialog = useTemplateRef<ChartColorSchemeDialogType>('chartColorSchemeDialog');
const accountCategorysDisplayOrderDialog = useTemplateRef<AccountCategoryDisplayOrderDialogType>('accountCategorysDisplayOrderDialog');

const showAccountsIncludedInHomePageOverviewDialog = ref<boolean>(false);
const showTransactionCategoriesIncludedInHomePageOverviewDialog = ref<boolean>(false);
const showAccountsIncludedInTotalDialog = ref<boolean>(false);

const allPageCounts = computed<NameNumeralValue[]>(() => getTablePageOptions(DEFAULT_PAGE_COUNTS, undefined, false, true));
const allInsightsExplorerDefaultDateRanges = computed<LocalizedDateRange[]>(() => getAllDateRanges(DateRangeScene.InsightsExplorer, {}));

const currentTheme = computed<string>({
    get: () => settingsStore.appSettings.theme,
    set: (value: string) => {
        if (value !== settingsStore.appSettings.theme) {
            settingsStore.setTheme(value);

            if (value === ThemeType.Light || value === ThemeType.Dark) {
                theme.change(value);
            } else {
                theme.change(getSystemTheme());
            }
        }
    }
});

const desktopOverviewPageLayoutDisplayContent = computed(() => {
    try {
        return tt(isDefaultDesktopOverviewLayout(parseDesktopOverviewLayout(settingsStore.appSettings.desktopOverviewPageLayout)) ? 'Default' : 'Custom');
    } catch {
        return tt('Custom');
    }
});

const showAddTransactionButtonInDesktopNavbar = computed<boolean>({
    get: () => settingsStore.appSettings.showAddTransactionButtonInDesktopNavbar,
    set: (value) => settingsStore.setShowAddTransactionButtonInDesktopNavbar(value)
});

const rememberLastSelectedFileTypeInImportTransactionDialog = computed<boolean>({
    get: () => settingsStore.appSettings.rememberLastSelectedFileTypeInImportTransactionDialog,
    set: (value) => settingsStore.setRememberLastSelectedFileTypeInImportTransactionDialog(value)
});

const insightsExplorerDefaultDateRangeType = computed<number>({
    get: () => settingsStore.appSettings.insightsExplorerDefaultDateRangeType,
    set: (value) => settingsStore.setInsightsExplorerDefaultDateRangeType(value)
});

const showTagInInsightsExplorerPage = computed<boolean>({
    get: () => settingsStore.appSettings.showTagInInsightsExplorerPage,
    set: (value) => settingsStore.setShowTagInInsightsExplorerPage(value)
});

const hideCategoriesWithoutAccounts = computed<boolean>({
    get: () => settingsStore.appSettings.hideCategoriesWithoutAccounts,
    set: (value) => settingsStore.setHideCategoriesWithoutAccounts(value)
});

const reconciliationStatementButtonDefaultDateRangeTypeInDesktop = computed<number>({
    get: () => getValidReconciliationStatementPageDefaultDateRangeType(settingsStore.appSettings.reconciliationStatementButtonDefaultDateRangeTypeInDesktop, DEFAULT_RECONCILIATION_STATEMENT_DATE_RANGE_IN_DESKTOP.type),
    set: (value: number) => settingsStore.setReconciliationStatementButtonDefaultDateRangeTypeInDesktop(value)
});

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
