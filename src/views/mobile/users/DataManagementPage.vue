<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar :title="$t('Data Management')" :back-link="$t('Back')"></f7-navbar>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="loading">
            <f7-list-item title="Accounts" after="Count"></f7-list-item>
            <f7-list-item title="Transaction Categories" after="Count"></f7-list-item>
            <f7-list-item title="Transaction Tags" after="Count"></f7-list-item>
            <f7-list-item title="Transactions" after="Count"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical" v-else-if="!loading">
            <f7-list-item :title="$t('Accounts')" :after="displayDataStatistics.totalAccountCount"></f7-list-item>
            <f7-list-item :title="$t('Transaction Categories')" :after="displayDataStatistics.totalTransactionCategoryCount"></f7-list-item>
            <f7-list-item :title="$t('Transaction Tags')" :after="displayDataStatistics.totalTransactionTagCount"></f7-list-item>
            <f7-list-item :title="$t('Transactions')" :after="displayDataStatistics.totalTransactionCount"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical" :class="{ 'disabled': loading }">
            <f7-list-button :class="{ 'disabled': !dataStatistics || !dataStatistics.totalTransactionCount || dataStatistics.totalTransactionCount === '0' }"
                            v-if="isDataExportingEnabled"
                            @click="exportedData = null; showExportDataSheet = true">{{ $t('Export Data') }}</f7-list-button>
            <f7-list-button color="red" @click="clearData(null)">{{ $t('Clear User Data') }}</f7-list-button>
        </f7-list>

        <f7-sheet swipe-handler=".swipe-handler" style="height:auto"
                  :swipe-to-close="!exportingData" :close-on-escape="!exportingData"
                  :close-by-backdrop-click="!exportingData" :close-by-outside-click="!exportingData"
                  :opened="showExportDataSheet" @sheet:closed="showExportDataSheet = false; exportedData = null;">
            <div class="swipe-handler" style="z-index: 10"></div>
            <f7-page-content class="margin-top no-padding-top">
                <div class="display-flex padding justify-content-space-between align-items-center">
                    <div class="ebk-sheet-title"><b>{{ $t('Are you sure you want to export all transaction data to file?') }}</b></div>
                </div>
                <div class="padding-bottom padding-horizontal">
                    <f7-list class="export-file-type-list no-margin" dividers>
                        <f7-list-item radio radio-icon="start" :class="{ 'disabled': exportingData || exportedData }"
                                      :title="$t('CSV (Comma-separated values) File')"
                                      :checked="exportFileType === 'csv'" @change="exportFileType = 'csv'">
                        </f7-list-item>
                        <f7-list-item radio radio-icon="start" :class="{ 'disabled': exportingData || exportedData }"
                                      :title="$t('TSV (Tab-separated values) File')"
                                      :checked="exportFileType === 'tsv'" @change="exportFileType = 'tsv'">
                        </f7-list-item>
                    </f7-list>
                </div>
                <div class="padding-horizontal padding-bottom">
                    <p class="no-margin-top margin-bottom-half">{{ $t('It may take a long time, please wait for a few minutes.') }}</p>
                    <f7-button large fill :class="{ 'disabled': exportingData }" :text="$t('Continue')" @click="exportData" v-if="!exportedData"></f7-button>
                    <f7-button large fill external :text="$t('Save Data')" :download="exportFileName" :href="exportedData" target="_blank" v-if="exportedData"></f7-button>
                    <div class="margin-top text-align-center">
                        <f7-link :class="{ 'disabled': exportingData }" @click="showExportDataSheet = false" :text="$t('Cancel')"></f7-link>
                    </div>
                </div>
            </f7-page-content>
        </f7-sheet>

        <password-input-sheet :title="$t('Are you sure you want to clear all data?')"
                              :hint="$t('You CANNOT undo this action. This will clear your accounts, categories, tags and transactions data. Please enter your current password to confirm.')"
                              :confirm-disabled="clearingData"
                              :cancel-disabled="clearingData"
                              color="red"
                              v-model:show="showInputPasswordSheetForClearData"
                              v-model="currentPasswordForClearData"
                              @password:confirm="clearData">
        </password-input-sheet>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';

import { appendThousandsSeparator } from '@/lib/common.js';
import { isDataExportingEnabled } from '@/lib/server_settings.js';

export default {
    props: [
        'f7router'
    ],
    data() {
        return {
            loading: true,
            loadingError: null,
            dataStatistics: null,
            exportFileType: 'csv',
            exportingData: false,
            exportedData: null,
            currentPasswordForClearData: '',
            clearingData: false,
            showExportDataSheet: false,
            showInputPasswordSheetForClearData: false,
        };
    },
    computed: {
        ...mapStores(useRootStore, useSettingsStore, useUserStore),
        isEnableThousandsSeparator() {
            return this.settingsStore.appSettings.thousandsSeparator;
        },
        displayDataStatistics() {
            const self = this;

            if (!self.dataStatistics) {
                return null;
            }

            return {
                totalAccountCount: appendThousandsSeparator(self.dataStatistics.totalAccountCount, self.isEnableThousandsSeparator),
                totalTransactionCategoryCount: appendThousandsSeparator(self.dataStatistics.totalTransactionCategoryCount, self.isEnableThousandsSeparator),
                totalTransactionTagCount: appendThousandsSeparator(self.dataStatistics.totalTransactionTagCount, self.isEnableThousandsSeparator),
                totalTransactionCount: appendThousandsSeparator(self.dataStatistics.totalTransactionCount, self.isEnableThousandsSeparator)
            };
        },
        isDataExportingEnabled() {
            return isDataExportingEnabled();
        },
        exportFileName() {
            const nickname = this.userStore.currentUserNickname;

            if (nickname) {
                return this.$t('dataExport.exportFilename', {
                    nickname: nickname
                }) + '.' + this.exportFileType;
            }

            return this.$t('dataExport.defaultExportFilename') + '.' + this.exportFileType;
        },
    },
    created() {
        const self = this;

        self.loading = true;

        self.userStore.getUserDataStatistics().then(dataStatistics => {
            self.dataStatistics = dataStatistics;
            self.loading = false;
        }).catch(error => {
            if (error.processed) {
                self.loading = false;
            } else {
                self.loadingError = error;
                self.$toast(error.message || error);
            }
        });
    },
    methods: {
        onPageAfterIn() {
            this.$routeBackOnError(this.f7router, 'loadingError');
        },
        exportData() {
            const self = this;

            self.$showLoading();
            self.exportingData = true;

            self.userStore.getExportedUserData(self.exportFileType).then(data => {
                self.exportedData = URL.createObjectURL(data);
                self.exportingData = false;
                self.$hideLoading();
            }).catch(error => {
                self.exportedData = null;
                self.exportingData = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        clearData(password) {
            const self = this;

            if (!password) {
                self.currentPasswordForClearData = '';
                self.showInputPasswordSheetForClearData = true;
                return;
            }

            self.clearingData = true;
            self.$showLoading(() => self.clearingData);

            self.rootStore.clearUserData({
                password: password
            }).then(() => {
                self.clearingData = false;
                self.$hideLoading();

                self.showInputPasswordSheetForClearData = false;
                self.$toast('All user data has been cleared');

                self.loading = true;

                self.userStore.getUserDataStatistics().then(dataStatistics => {
                    self.dataStatistics = dataStatistics;
                    self.loading = false;
                }).catch(error => {
                    if (error.processed) {
                        self.loading = false;
                    } else {
                        self.loadingError = error;
                        self.$toast(error.message || error);
                    }
                });
            }).catch(error => {
                self.clearingData = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
    }
};
</script>

<style>
.export-file-type-list.list > ul > li > .item-content {
    padding-left: 0;
}
</style>
