<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar :title="$t('Data Management')" :back-link="$t('Back')"></f7-navbar>

        <f7-list strong inset dividers class="margin-top skeleton-text" v-if="loading">
            <f7-list-item title="Accounts" after="Count"></f7-list-item>
            <f7-list-item title="Transaction Categories" after="Count"></f7-list-item>
            <f7-list-item title="Transaction Tags" after="Count"></f7-list-item>
            <f7-list-item title="Transactions" after="Count"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-top" v-else-if="!loading">
            <f7-list-item :title="$t('Accounts')" :after="dataStatistics.totalAccountCount"></f7-list-item>
            <f7-list-item :title="$t('Transaction Categories')" :after="dataStatistics.totalTransactionCategoryCount"></f7-list-item>
            <f7-list-item :title="$t('Transaction Tags')" :after="dataStatistics.totalTransactionTagCount"></f7-list-item>
            <f7-list-item :title="$t('Transactions')" :after="dataStatistics.totalTransactionCount"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers :class="{ 'disabled': loading }">
            <f7-list-button @click="exportedData = null; showExportDataSheet = true" v-if="isDataExportingEnabled">{{ $t('Export Data') }}</f7-list-button>
            <f7-list-button color="red" @click="clearData(null)">{{ $t('Clear User Data') }}</f7-list-button>
        </f7-list>

        <f7-sheet swipe-handler=".swipe-handler" style="height:auto"
                  :swipe-to-close="!exportingData" :close-on-escape="!exportingData"
                  :close-by-backdrop-click="!exportingData" :close-by-outside-click="!exportingData"
                  :opened="showExportDataSheet" @sheet:closed="showExportDataSheet = false; exportedData = null;">
            <div class="swipe-handler"></div>
            <f7-page-content class="margin-top no-padding-top">
                <div class="display-flex padding justify-content-space-between align-items-center">
                    <div style="font-size: 18px"><b>{{ $t('Are you sure you want to export all data to csv file?') }}</b></div>
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
                              :hint="$t('You CANNOT undo this action. This will clear your accounts, categories, tags and transactions data. Please input your current password to confirm.')"
                              :confirm-disabled="clearingData"
                              :cancel-disabled="clearingData"
                              v-model:show="showInputPasswordSheetForClearData"
                              v-model="currentPasswordForClearData"
                              @password:confirm="clearData">
        </password-input-sheet>
    </f7-page>
</template>

<script>
export default {
    props: [
        'f7router'
    ],
    data() {
        return {
            loading: true,
            loadingError: null,
            dataStatistics: null,
            exportingData: false,
            exportedData: null,
            currentPasswordForClearData: '',
            clearingData: false,
            showExportDataSheet: false,
            showInputPasswordSheetForClearData: false,
        };
    },
    computed: {
        isDataExportingEnabled() {
            return this.$settings.isDataExportingEnabled();
        },
        exportFileName() {
            const nickname = this.$store.getters.currentUserNickname;

            if (nickname) {
                return this.$t('dataExport.exportFilename', {
                    nickname: nickname
                }) + '.csv';
            }

            return this.$t('dataExport.defaultExportFilename') + '.csv';
        },
    },
    created() {
        const self = this;

        self.loading = true;

        self.$store.dispatch('getUserDataStatistics').then(dataStatistics => {
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

            self.$store.dispatch('getExportedUserData').then(data => {
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

            self.$store.dispatch('clearUserData', {
                password: password
            }).then(() => {
                self.clearingData = false;
                self.$hideLoading();

                self.showInputPasswordSheetForClearData = false;
                self.$toast('All user data has been cleared');

                self.loading = true;

                self.$store.dispatch('getUserDataStatistics').then(dataStatistics => {
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
