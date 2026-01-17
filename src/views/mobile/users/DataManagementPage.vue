<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar :class="{ 'disabled': loading }" :back-link="tt('Back')" :title="tt('Data Management')"></f7-navbar>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="loading">
            <f7-list-item title="Transactions" after="Count"></f7-list-item>
            <f7-list-item title="Transaction Pictures" after="Count"></f7-list-item>
            <f7-list-item title="Insights Explorers" after="Count"></f7-list-item>
            <f7-list-item title="Accounts" after="Count"></f7-list-item>
            <f7-list-item title="Transaction Categories" after="Count"></f7-list-item>
            <f7-list-item title="Transaction Tags" after="Count"></f7-list-item>
            <f7-list-item title="Transaction Templates" after="Count"></f7-list-item>
            <f7-list-item title="Scheduled Transactions" after="Count"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical" v-else-if="!loading">
            <f7-list-item :title="tt('Transactions')" :after="displayDataStatistics ? displayDataStatistics.totalTransactionCount : '-'"></f7-list-item>
            <f7-list-item :title="tt('Transaction Pictures')" :after="displayDataStatistics ? displayDataStatistics.totalTransactionPictureCount : '-'"></f7-list-item>
            <f7-list-item :title="tt('Insights Explorers')" :after="displayDataStatistics ? displayDataStatistics.totalInsightsExplorerCount : '-'"></f7-list-item>
            <f7-list-item :title="tt('Accounts')" :after="displayDataStatistics ? displayDataStatistics.totalAccountCount : '-'"></f7-list-item>
            <f7-list-item :title="tt('Transaction Categories')" :after="displayDataStatistics ? displayDataStatistics.totalTransactionCategoryCount : '-'"></f7-list-item>
            <f7-list-item :title="tt('Transaction Tags')" :after="displayDataStatistics ? displayDataStatistics.totalTransactionTagCount : '-'"></f7-list-item>
            <f7-list-item :title="tt('Transaction Templates')" :after="displayDataStatistics ? displayDataStatistics.totalTransactionTemplateCount : '-'"></f7-list-item>
            <f7-list-item :title="tt('Scheduled Transactions')" :after="displayDataStatistics ? displayDataStatistics.totalScheduledTransactionCount : '-'"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical" :class="{ 'disabled': loading }">
            <f7-list-button :class="{ 'disabled': !dataStatistics || !dataStatistics.totalTransactionCount || dataStatistics.totalTransactionCount === '0' }"
                            v-if="isDataExportingEnabled()"
                            @click="exportedData = null; showExportDataSheet = true">{{ tt('Export Data') }}</f7-list-button>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical" :class="{ 'disabled': loading }">
            <f7-list-button color="red" @click="clearAllTransactions(null)">{{ tt('Clear All Transactions') }}</f7-list-button>
            <f7-list-button color="red" @click="clearAllData(null)">{{ tt('Clear All Data') }}</f7-list-button>
        </f7-list>

        <f7-sheet swipe-handler=".swipe-handler" style="height:auto"
                  :swipe-to-close="!exportingData" :close-on-escape="!exportingData"
                  :close-by-backdrop-click="!exportingData" :close-by-outside-click="!exportingData"
                  :opened="showExportDataSheet" @sheet:closed="showExportDataSheet = false; exportedData = null;">
            <div class="swipe-handler"></div>
            <f7-page-content class="margin-top no-padding-top">
                <div class="display-flex padding justify-content-space-between align-items-center">
                    <div class="ebk-sheet-title"><b>{{ tt('Are you sure you want to export all transaction data to file?') }}</b></div>
                </div>
                <div class="padding-bottom padding-horizontal">
                    <f7-list class="export-file-type-list no-margin" dividers>
                        <f7-list-item radio radio-icon="start" :class="{ 'disabled': exportingData || exportedData }"
                                      :title="tt('CSV (Comma-separated values) File')"
                                      :checked="exportFileType === 'csv'" @change="exportFileType = 'csv'">
                        </f7-list-item>
                        <f7-list-item radio radio-icon="start" :class="{ 'disabled': exportingData || exportedData }"
                                      :title="tt('TSV (Tab-separated values) File')"
                                      :checked="exportFileType === 'tsv'" @change="exportFileType = 'tsv'">
                        </f7-list-item>
                    </f7-list>
                </div>
                <div class="padding-horizontal padding-bottom">
                    <p class="no-margin-top margin-bottom-half">{{ tt('It may take a long time, please wait for a few minutes.') }}</p>
                    <f7-button large fill :class="{ 'disabled': exportingData }" :text="tt('Continue')" @click="exportData" v-if="!exportedData"></f7-button>
                    <f7-button large fill external :text="tt('Save Data')" :download="exportFileName" :href="exportedData" target="_blank" v-if="exportedData"></f7-button>
                    <div class="margin-top text-align-center">
                        <f7-link :class="{ 'disabled': exportingData }" @click="showExportDataSheet = false" :text="tt('Cancel')"></f7-link>
                    </div>
                </div>
            </f7-page-content>
        </f7-sheet>

        <password-input-sheet :title="tt('Are you sure you want to clear all transactions?')"
                              :hint="tt('You CANNOT undo this action. This will clear your transactions data. Please enter your current password to confirm.')"
                              :confirm-disabled="clearingData"
                              :cancel-disabled="clearingData"
                              color="red"
                              v-model:show="showInputPasswordSheetForClearAllTransactions"
                              v-model="currentPasswordForClearData"
                              @password:confirm="clearAllTransactions">
        </password-input-sheet>

        <password-input-sheet :title="tt('Are you sure you want to clear all data?')"
                              :hint="tt('You CANNOT undo this action. This will clear your accounts, categories, tags and transactions data. Please enter your current password to confirm.')"
                              :confirm-disabled="clearingData"
                              :cancel-disabled="clearingData"
                              color="red"
                              v-model:show="showInputPasswordSheetForClearAllData"
                              v-model="currentPasswordForClearData"
                              @password:confirm="clearAllData">
        </password-input-sheet>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';
import { useDataManagementPageBase } from '@/views/base/users/DataManagementPageBase.ts';

import { useRootStore } from '@/stores/index.ts';
import { useUserStore } from '@/stores/user.ts';

import { isDataExportingEnabled } from '@/lib/server_settings.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt } = useI18n();
const { showToast, routeBackOnError } = useI18nUIComponents();
const { dataStatistics, displayDataStatistics, getExportFileName } = useDataManagementPageBase();

const rootStore = useRootStore();
const userStore = useUserStore();

const loading = ref<boolean>(true);
const loadingError = ref<unknown | null>(null);
const exportFileType = ref<string>('csv');
const exportingData = ref<boolean>(false);
const exportedData = ref<string | null>(null);
const currentPasswordForClearData = ref<string>('');
const clearingData = ref<boolean>(false);
const showExportDataSheet = ref<boolean>(false);
const showInputPasswordSheetForClearAllTransactions = ref<boolean>(false);
const showInputPasswordSheetForClearAllData = ref<boolean>(false);

const exportFileName = computed<string>(() => getExportFileName(exportFileType.value));

function reloadUserDataStatistics(): void {
    loading.value = true;

    userStore.getUserDataStatistics().then(dataStatisticsResponse => {
        dataStatistics.value = dataStatisticsResponse;
        loading.value = false;
    }).catch(error => {
        if (error.processed) {
            loading.value = false;
        } else {
            loadingError.value = error;
            showToast(error.message || error);
        }
    });
}

function exportData(): void {
    showLoading();
    exportingData.value = true;

    userStore.getExportedUserData(exportFileType.value).then(data => {
        exportedData.value = URL.createObjectURL(data);
        exportingData.value = false;
        hideLoading();
    }).catch(error => {
        exportedData.value = null;
        exportingData.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function clearAllTransactions(password: string | null): void {
    if (!password) {
        currentPasswordForClearData.value = '';
        showInputPasswordSheetForClearAllTransactions.value = true;
        return;
    }

    clearingData.value = true;
    showLoading(() => clearingData.value);

    rootStore.clearAllUserTransactions({
        password: password
    }).then(() => {
        clearingData.value = false;
        currentPasswordForClearData.value = '';
        hideLoading();

        showInputPasswordSheetForClearAllTransactions.value = false;
        showToast('All transactions has been cleared');

        reloadUserDataStatistics();
    }).catch(error => {
        clearingData.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function clearAllData(password: string | null): void {
    if (!password) {
        currentPasswordForClearData.value = '';
        showInputPasswordSheetForClearAllData.value = true;
        return;
    }

    clearingData.value = true;
    showLoading(() => clearingData.value);

    rootStore.clearAllUserData({
        password: password
    }).then(() => {
        clearingData.value = false;
        currentPasswordForClearData.value = '';
        hideLoading();

        showInputPasswordSheetForClearAllData.value = false;
        showToast('All user data has been cleared');

        reloadUserDataStatistics();
    }).catch(error => {
        clearingData.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function onPageAfterIn(): void {
    routeBackOnError(props.f7router, loadingError);
}

reloadUserDataStatistics();
</script>

<style>
.export-file-type-list.list > ul > li > .item-content {
    padding-inline-start: 0;
}
</style>
