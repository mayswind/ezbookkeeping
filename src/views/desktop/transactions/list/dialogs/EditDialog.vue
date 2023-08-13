<template>
    <v-dialog width="1000" :persistent="!!persistent && mode !== 'view'" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-8">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <div class="d-flex w-100 align-center justify-center">
                        <h5 class="text-h5">{{ $t(title) }}</h5>
                        <v-progress-circular indeterminate size="22" class="ml-2" v-if="loading"></v-progress-circular>
                    </div>
                    <v-btn density="comfortable" color="default" variant="text" class="ml-2" :icon="true"
                           :disabled="loading || submitting" v-if="mode !== 'view'">
                        <v-icon :icon="icons.more" />
                        <v-menu activator="parent">
                            <v-list>

                            </v-list>
                        </v-menu>
                    </v-btn>
                </div>
            </template>
            <v-card-text class="d-flex flex-column flex-md-row mt-2 mt-md-4">
                <div class="mb-4">
                    <v-tabs direction="vertical" :disabled="loading || submitting" v-model="activeTab">
                        <v-tab value="basicInfo">
                            <span>{{ $t('Basic Information') }}</span>
                        </v-tab>
                        <v-tab value="map" :disabled="!transaction.geoLocation">
                            <span>{{ $t('Location on Map') }}</span>
                        </v-tab>
                    </v-tabs>
                </div>

                <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container ml-md-5"
                          v-model="activeTab">
                    <v-window-item value="basicInfo">
                        <v-form class="mt-2">
                            <v-row>
                                <v-col cols="12" md="12">
                                </v-col>
                            </v-row>
                        </v-form>
                    </v-window-item>
                    <v-window-item value="map">
                        <v-row>
                            <v-col cols="12" md="12">

                            </v-col>
                        </v-row>
                    </v-window-item>
                </v-window>
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center mt-2 mt-sm-4 mt-md-6 gap-4">
                    <v-btn :disabled="inputIsEmpty || loading || submitting" v-if="mode !== 'view'" @click="save">
                        {{ $t(saveButtonTitle) }}
                        <v-progress-circular indeterminate size="24" class="ml-2" v-if="submitting"></v-progress-circular>
                    </v-btn>
                    <v-btn color="secondary" variant="tonal"
                           :disabled="loading || submitting" @click="cancel">{{ $t('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useAccountsStore } from '@/stores/account.js';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';
import { useTransactionTagsStore } from '@/stores/transactionTag.js';
import { useTransactionsStore } from '@/stores/transaction.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import categoryConstants from '@/consts/category.js';
import transactionConstants from '@/consts/transaction.js';
import logger from '@/lib/logger.js';
import {
    getNameByKeyValue
} from '@/lib/common.js';
import {
    getUtcOffsetByUtcOffsetMinutes
} from '@/lib/datetime.js';
import { getMapProvider } from '@/lib/server_settings.js';

import {
    mdiDotsVertical
} from '@mdi/js';
import { setTransactionModelByTransaction } from '@/lib/transaction.js';

export default {
    props: [
        'persistent',
        'show'
    ],
    expose: [
        'open'
    ],
    data() {
        const transactionsStore = useTransactionsStore();
        const newTransaction = transactionsStore.generateNewTransactionModel();

        return {
            showState: false,
            mode: 'add',
            activeTab: 'basicInfo',
            editTransactionId: null,
            loading: true,
            transaction: newTransaction,
            geoLocationStatus: null,
            submitting: false,
            isSupportGeoLocation: !!navigator.geolocation,
            resolve: null,
            reject: null,
            icons: {
                more: mdiDotsVertical
            }
        };
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore, useAccountsStore, useTransactionCategoriesStore, useTransactionTagsStore, useTransactionsStore, useExchangeRatesStore),
        title() {
            if (this.mode === 'add') {
                return 'Add Transaction';
            } else if (this.mode === 'edit') {
                return 'Edit Transaction';
            } else {
                return 'Transaction Detail';
            }
        },
        saveButtonTitle() {
            if (this.mode === 'add') {
                return 'Add';
            } else {
                return 'Save';
            }
        },
        sourceAmountName() {
            if (this.transaction.type === this.allTransactionTypes.Expense) {
                return 'Expense Amount';
            } else if (this.transaction.type === this.allTransactionTypes.Income) {
                return 'Income Amount';
            } else if (this.transaction.type === this.allTransactionTypes.Transfer) {
                return 'Transfer Out Amount';
            } else {
                return '';
            }
        },
        sourceAccountTitle() {
            if (this.transaction.type === this.allTransactionTypes.Expense || this.transaction.type === this.allTransactionTypes.Income) {
                return 'Account';
            } else if (this.transaction.type === this.allTransactionTypes.Transfer) {
                return 'Source Account';
            } else {
                return '';
            }
        },
        defaultCurrency() {
            return this.userStore.currentUserDefaultCurrency;
        },
        defaultAccountId() {
            return this.userStore.currentUserDefaultAccountId;
        },
        allTransactionTypes() {
            return transactionConstants.allTransactionTypes;
        },
        allCategoryTypes() {
            return categoryConstants.allCategoryTypes;
        },
        allTimezones() {
            return this.$locale.getAllTimezones(true);
        },
        allAccounts() {
            return this.accountsStore.allPlainAccounts;
        },
        allVisibleAccounts() {
            return this.accountsStore.allVisiblePlainAccounts;
        },
        allAccountsMap() {
            return this.accountsStore.allAccountsMap;
        },
        allCategories() {
            return this.transactionCategoriesStore.allTransactionCategories;
        },
        allCategoriesMap() {
            return this.transactionCategoriesStore.allTransactionCategoriesMap;
        },
        allTags() {
            return this.transactionTagsStore.allTransactionTags;
        },
        hasAvailableExpenseCategories() {
            if (!this.allCategories || !this.allCategories[this.allCategoryTypes.Expense] || !this.allCategories[this.allCategoryTypes.Expense].length) {
                return false;
            }

            const firstAvailableCategoryId = this.getFirstAvailableCategoryId(this.allCategories[this.allCategoryTypes.Expense]);
            return firstAvailableCategoryId !== '';
        },
        hasAvailableIncomeCategories() {
            if (!this.allCategories || !this.allCategories[this.allCategoryTypes.Income] || !this.allCategories[this.allCategoryTypes.Income].length) {
                return false;
            }

            const firstAvailableCategoryId = this.getFirstAvailableCategoryId(this.allCategories[this.allCategoryTypes.Income]);
            return firstAvailableCategoryId !== '';
        },
        hasAvailableTransferCategories() {
            if (!this.allCategories || !this.allCategories[this.allCategoryTypes.Transfer] || !this.allCategories[this.allCategoryTypes.Transfer].length) {
                return false;
            }

            const firstAvailableCategoryId = this.getFirstAvailableCategoryId(this.allCategories[this.allCategoryTypes.Transfer]);
            return firstAvailableCategoryId !== '';
        },
        transactionDisplayTimezone() {
            return `UTC${getUtcOffsetByUtcOffsetMinutes(this.transaction.utcOffset)}`;
        },
        transactionDisplayTimezoneName() {
            return getNameByKeyValue(this.allTimezones, this.transaction.timeZone, 'name', 'displayName');
        },
        transactionTimezoneTimeDifference() {
            return this.$locale.getTimezoneDifferenceDisplayText(this.transaction.utcOffset);
        },
        geoLocationStatusInfo() {
            if (this.geoLocationStatus === 'success') {
                return '';
            } else if (this.geoLocationStatus === 'getting') {
                return this.$t('Getting Location...');
            } else {
                return this.$t('No Location');
            }
        },
        mapProvider() {
            return getMapProvider();
        },
        inputIsEmpty() {
            return !!this.inputEmptyProblemMessage;
        },
        inputEmptyProblemMessage() {
            return null;
        }
    },
    watch: {
        'transaction.sourceAmount': function (newValue, oldValue) {
            if (this.mode === 'view') {
                return;
            }

            this.transactionsStore.setTransactionSuitableDestinationAmount(this.transaction, oldValue, newValue);
        },
        'transaction.destinationAmount': function (newValue) {
            if (this.mode === 'view') {
                return;
            }

            if (this.transaction.type === this.allTransactionTypes.Expense || this.transaction.type === this.allTransactionTypes.Income) {
                this.transaction.sourceAmount = newValue;
            }
        },
        'transaction.timeZone': function (newValue) {
            for (let i = 0; i < this.allTimezones.length; i++) {
                if (this.allTimezones[i].name === newValue) {
                    this.transaction.utcOffset = this.allTimezones[i].utcOffsetMinutes;
                    break;
                }
            }
        }
    },
    methods: {
        open(options) {
            const self = this;
            self.showState = true;
            self.activeTab = 'basicInfo';
            self.loading = true;
            self.submitting = false;

            const newTransaction = self.transactionsStore.generateNewTransactionModel(options.type);
            self.setTransaction(newTransaction, options);

            const promises = [
                self.accountsStore.loadAllAccounts({ force: false }),
                self.transactionCategoriesStore.loadAllCategories({ force: false }),
                self.transactionTagsStore.loadAllTags({ force: false })
            ];

            if (options && options.id) {
                if (options.currentTransaction) {
                    self.setTransaction(options.currentTransaction, options);
                }

                self.mode = 'view';
                self.editTransactionId = options.id;

                promises.push(self.transactionsStore.getTransaction({ transactionId: self.editTransactionId }));
            } else {
                self.mode = 'add';
                self.editTransactionId = null;

                if (self.settingsStore.appSettings.autoGetCurrentGeoLocation
                    && !self.geoLocationStatus && !self.transaction.geoLocation) {
                    self.updateGeoLocation(false);
                }
            }

            if (options.type && options.type !== '0' &&
                options.type >= self.allTransactionTypes.Income &&
                options.type <= self.allTransactionTypes.Transfer) {
                self.transaction.type = parseInt(options.type);
            }

            Promise.all(promises).then(function (responses) {
                if (self.editTransactionId && !responses[3]) {
                    if (self.reject) {
                        self.reject('Unable to get transaction');
                    }

                    return;
                }

                self.setTransaction(options.id ? responses[3] : null, options);

                self.loading = false;
            }).catch(error => {
                logger.error('failed to load essential data for editing transaction', error);

                self.loading = false;
                self.showState = false;

                if (!error.processed) {
                    if (self.reject) {
                        self.reject(error);
                    }
                }
            });

            return new Promise((resolve, reject) => {
                self.resolve = resolve;
                self.reject = reject;
            });
        },
        save() {
            const self = this;

            if (self.mode === 'view') {
                return;
            }


        },
        cancel() {
            if (this.reject) {
                this.reject();
            }

            this.showState = false;
        },
        updateGeoLocation(forceUpdate) {
            const self = this;

            if (!self.isSupportGeoLocation) {
                logger.warn('this browser does not support geo location');

                if (forceUpdate) {
                    self.$refs.snackbar.showMessage('Unable to get current position');
                }
                return;
            }

            navigator.geolocation.getCurrentPosition(function (position) {
                if (!position || !position.coords) {
                    logger.error('current position is null');
                    self.geoLocationStatus = 'error';

                    if (forceUpdate) {
                        self.$refs.snackbar.showMessage('Unable to get current position');
                    }

                    return;
                }

                self.geoLocationStatus = 'success';

                self.transaction.geoLocation = {
                    latitude: position.coords.latitude,
                    longitude: position.coords.longitude
                };
            }, function (err) {
                logger.error('cannot get current position', err);
                self.geoLocationStatus = 'error';

                if (forceUpdate) {
                    self.$refs.snackbar.showMessage('Unable to get current position');
                }
            });

            self.geoLocationStatus = 'getting';
        },
        getDisplayAmount(amount, hideAmount) {
            if (hideAmount) {
                return this.getDisplayCurrency('***');
            }

            return this.getDisplayCurrency(amount);
        },
        getDisplayCurrency(value, currencyCode) {
            return this.$locale.getDisplayCurrency(value, currencyCode, {
                currencyDisplayMode: this.settingsStore.appSettings.currencyDisplayMode,
                enableThousandsSeparator: this.settingsStore.appSettings.thousandsSeparator
            });
        },
        setTransaction(transaction, options) {
            setTransactionModelByTransaction(
                this.transaction,
                transaction,
                this.allCategories,
                this.allCategoriesMap,
                this.allVisibleAccounts,
                this.allAccountsMap,
                this.defaultAccountId,
                {
                    type: options.type,
                    categoryId: options.categoryId,
                    accountId: options.accountId
                },
                (this.mode !== 'edit' && this.mode !== 'view')
            );
        }
    }
}
</script>
