<template>
    <v-row>
        <v-col cols="12">
            <v-card :class="{ 'disabled': loadingDataStatistics }">
                <template #title>
                    <div class="d-flex align-center">
                        <span>{{ $t('Data Management') }}</span>
                        <v-btn density="compact" color="default" variant="text" size="24"
                               class="ml-2" :icon="true" :loading="loadingDataStatistics" @click="reloadUserDataStatistics(true)">
                            <template #loader>
                                <v-progress-circular indeterminate size="20"/>
                            </template>
                            <v-icon :icon="icons.refresh" size="24" />
                            <v-tooltip activator="parent">{{ $t('Refresh') }}</v-tooltip>
                        </v-btn>
                    </div>
                </template>

                <v-card-text>
                    <v-row>
                        <v-col cols="6" sm="3" :key="idx" v-for="(item, idx) in [
                            {
                                title: 'Transactions',
                                count: displayDataStatistics ? displayDataStatistics.totalTransactionCount : '-',
                                icon: icons.transactions,
                                color: 'info-darken-1'
                            },
                            {
                                title: 'Accounts',
                                count: displayDataStatistics ? displayDataStatistics.totalAccountCount : '-',
                                icon: icons.accounts,
                                color: 'primary'
                            },
                            {
                                title: 'Transaction Categories',
                                count: displayDataStatistics ? displayDataStatistics.totalTransactionCategoryCount : '-',
                                icon: icons.categories,
                                color: 'teal'
                            },
                            {
                                title: 'Transaction Tags',
                                count: displayDataStatistics ? displayDataStatistics.totalTransactionTagCount : '-',
                                icon: icons.tags,
                                color: 'secondary'
                            },
                            {
                                title: 'Transaction Pictures',
                                count: displayDataStatistics ? displayDataStatistics.totalTransactionPictureCount : '-',
                                icon: icons.pictures,
                                color: 'error-darken-1'
                            },
                            {
                                title: 'Transaction Templates',
                                count: displayDataStatistics ? displayDataStatistics.totalTransactionTemplateCount : '-',
                                icon: icons.templates,
                                color: 'secondary-darken-1'
                            },
                            {
                                title: 'Scheduled Transactions',
                                count: displayDataStatistics ? displayDataStatistics.totalScheduledTransactionCount : '-',
                                icon: icons.scheduledTransactions,
                                color: 'success-darken-1'
                            }
                        ]">
                            <div class="d-flex align-center">
                                <div class="me-3">
                                    <v-avatar rounded :color="item.color" size="42" class="elevation-1">
                                        <v-icon size="24" :icon="item.icon"/>
                                    </v-avatar>
                                </div>

                                <div class="d-flex flex-column">
                                    <span class="text-caption">{{ $t(item.title) }}</span>
                                    <v-skeleton-loader class="skeleton-no-margin pt-2 pb-2" type="text" style="width: 60px" :loading="true" v-if="loadingDataStatistics"></v-skeleton-loader>
                                    <span class="text-xl" v-if="!loadingDataStatistics">{{ item.count }}</span>
                                </div>
                            </div>
                        </v-col>
                    </v-row>
                </v-card-text>
            </v-card>
        </v-col>

        <v-col cols="12" v-if="isDataExportingEnabled">
            <v-card :class="{ 'disabled': exportingData }" :title="$t('Export Data')">
                <v-card-text>
                    <span class="text-body-1">{{ $t('Export all transaction data to file.') }}&nbsp;{{ $t('It may take a long time, please wait for a few minutes.') }}</span>
                </v-card-text>

                <v-card-text class="d-flex flex-wrap gap-4">
                    <v-btn-group variant="elevated" density="comfortable" color="primary"
                                  :disabled="loadingDataStatistics || exportingData || !dataStatistics || !dataStatistics.totalTransactionCount || dataStatistics.totalTransactionCount === '0'">
                        <v-btn>
                            {{ $t('Export Data') }}
                            <v-progress-circular indeterminate size="22" class="ml-2" v-if="exportingData"></v-progress-circular>
                            <v-menu activator="parent">
                                <v-list :disabled="loadingDataStatistics || exportingData || !dataStatistics || !dataStatistics.totalTransactionCount || dataStatistics.totalTransactionCount === '0'">
                                    <v-list-item @click="exportData('csv')">
                                        <v-list-item-title>{{ $t('CSV (Comma-separated values) File') }}</v-list-item-title>
                                    </v-list-item>
                                    <v-list-item @click="exportData('tsv')">
                                        <v-list-item-title>{{ $t('TSV (Tab-separated values) File') }}</v-list-item-title>
                                    </v-list-item>
                                </v-list>
                            </v-menu>
                        </v-btn>
                    </v-btn-group>
                </v-card-text>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card :class="{ 'disabled': clearingData }">
                <template #title>
                    <span class="text-error">{{ $t('Danger Zone') }}</span>
                </template>

                <v-form>
                    <v-card-text class="py-0">
                    <span class="text-body-1 text-error">
                        <v-icon :icon="icons.alert"/>
                        {{ $t('You CANNOT undo this action. This will clear your accounts, categories, tags and transactions data. Please enter your current password to confirm.') }}
                    </span>
                    </v-card-text>

                    <v-card-text class="pb-0">
                        <v-row class="mb-3">
                            <v-col cols="12" md="6">
                                <v-text-field
                                    autocomplete="current-password"
                                    ref="currentPasswordInput"
                                    type="password"
                                    variant="underlined"
                                    color="error"
                                    :disabled="loadingDataStatistics || clearingData"
                                    :placeholder="$t('Current Password')"
                                    v-model="currentPasswordForClearData"
                                    @keyup.enter="clearData"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>

                    <v-card-text class="d-flex flex-wrap gap-4">
                        <v-btn color="error" :disabled="loadingDataStatistics || !currentPasswordForClearData || clearingData" @click="clearData">
                            {{ $t('Clear User Data') }}
                            <v-progress-circular indeterminate size="22" class="ml-2" v-if="clearingData"></v-progress-circular>
                        </v-btn>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>
    </v-row>

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script>
import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.js';

import { isEquals } from '@/lib/common.ts';
import { isDataExportingEnabled } from '@/lib/server_settings.ts';
import { startDownloadFile } from '@/lib/ui/common.ts';

import {
    mdiRefresh,
    mdiListBoxOutline,
    mdiCreditCardOutline,
    mdiViewDashboardOutline,
    mdiTagOutline,
    mdiClipboardTextOutline,
    mdiImage,
    mdiClipboardTextClockOutline,
    mdiAlert
} from '@mdi/js';

export default {
    data() {
        return {
            loadingDataStatistics: true,
            dataStatistics: null,
            exportingData: false,
            currentPasswordForClearData: '',
            clearingData: false,
            icons: {
                refresh: mdiRefresh,
                transactions: mdiListBoxOutline,
                accounts: mdiCreditCardOutline,
                categories: mdiViewDashboardOutline,
                tags: mdiTagOutline,
                pictures: mdiImage,
                templates: mdiClipboardTextOutline,
                scheduledTransactions: mdiClipboardTextClockOutline,
                alert: mdiAlert
            }
        }
    },
    computed: {
        ...mapStores(useRootStore, useSettingsStore, useUserStore),
        displayDataStatistics() {
            const self = this;

            if (!self.dataStatistics) {
                return null;
            }

            return {
                totalTransactionCount: self.$locale.appendDigitGroupingSymbol(self.userStore, self.dataStatistics.totalTransactionCount),
                totalAccountCount: self.$locale.appendDigitGroupingSymbol(self.userStore, self.dataStatistics.totalAccountCount),
                totalTransactionCategoryCount: self.$locale.appendDigitGroupingSymbol(self.userStore, self.dataStatistics.totalTransactionCategoryCount),
                totalTransactionTagCount: self.$locale.appendDigitGroupingSymbol(self.userStore, self.dataStatistics.totalTransactionTagCount),
                totalTransactionPictureCount: self.$locale.appendDigitGroupingSymbol(self.userStore, self.dataStatistics.totalTransactionPictureCount),
                totalTransactionTemplateCount: self.$locale.appendDigitGroupingSymbol(self.userStore, self.dataStatistics.totalTransactionTemplateCount),
                totalScheduledTransactionCount: self.$locale.appendDigitGroupingSymbol(self.userStore, self.dataStatistics.totalScheduledTransactionCount)
            };
        },
        isDataExportingEnabled() {
            return isDataExportingEnabled();
        }
    },
    created() {
        this.reloadUserDataStatistics(false);
    },
    methods: {
        reloadUserDataStatistics(force) {
            const self = this;

            self.loadingDataStatistics = true;

            self.userStore.getUserDataStatistics().then(dataStatistics => {
                if (force) {
                    if (isEquals(self.dataStatistics, dataStatistics)) {
                        self.$refs.snackbar.showMessage('Data is up to date');
                    } else {
                        self.$refs.snackbar.showMessage('Data has been updated');
                    }
                }

                self.dataStatistics = dataStatistics;
                self.loadingDataStatistics = false;
            }).catch(error => {
                self.loadingDataStatistics = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        exportData(fileType) {
            const self = this;

            if (self.exportingData) {
                return;
            }

            self.exportingData = true;

            self.userStore.getExportedUserData(fileType).then(data => {
                startDownloadFile(self.getExportFileName(fileType), data);
                self.exportingData = false;
            }).catch(error => {
                self.exportingData = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        clearData() {
            const self = this;

            if (!self.currentPasswordForClearData) {
                self.$refs.snackbar.showMessage('Current password cannot be blank');
                return;
            }

            if (self.clearingData) {
                return;
            }

            self.$refs.confirmDialog.open('Are you sure you want to clear all data?', { color: 'error' }).then(() => {
                self.clearingData = true;

                self.rootStore.clearUserData({
                    password: self.currentPasswordForClearData
                }).then(() => {
                    self.clearingData = false;

                    self.$refs.snackbar.showMessage('All user data has been cleared');
                    self.reloadUserDataStatistics();
                }).catch(error => {
                    self.clearingData = false;

                    if (!error.processed) {
                        self.$refs.snackbar.showError(error);
                    }
                });
            });
        },
        getExportFileName(fileExtension) {
            const nickname = this.userStore.currentUserNickname;

            if (nickname) {
                return this.$t('dataExport.exportFilename', {
                    nickname: nickname
                }) + '.' + fileExtension;
            }

            return this.$t('dataExport.defaultExportFilename') + '.' + fileExtension;
        }
    }
}
</script>
