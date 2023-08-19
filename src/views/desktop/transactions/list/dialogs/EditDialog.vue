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
                                <v-list-item :prepend-icon="icons.show"
                                             :title="$t('Show Amount')"
                                             v-if="transaction.hideAmount" @click="transaction.hideAmount = false"></v-list-item>
                                <v-list-item :prepend-icon="icons.hide"
                                             :title="$t('Hide Amount')"
                                             v-if="!transaction.hideAmount" @click="transaction.hideAmount = true"></v-list-item>
                            </v-list>
                        </v-menu>
                    </v-btn>
                </div>
            </template>
            <v-card-text class="d-flex flex-column flex-md-row mt-2 mt-md-4">
                <div class="mb-4">
                    <v-tabs class="v-tabs-pill" direction="vertical" :class="{ 'readonly': mode !== 'add' }"
                            :disabled="loading || submitting" v-model="transaction.type">
                        <v-tab :value="allTransactionTypes.Expense">
                            <span>{{ $t('Expense') }}</span>
                        </v-tab>
                        <v-tab :value="allTransactionTypes.Income">
                            <span>{{ $t('Income') }}</span>
                        </v-tab>
                        <v-tab :value="allTransactionTypes.Transfer">
                            <span>{{ $t('Transfer') }}</span>
                        </v-tab>
                    </v-tabs>
                    <v-divider class="my-2"/>
                    <v-tabs direction="vertical" :disabled="loading || submitting" v-model="activeTab">
                        <v-tab value="basicInfo">
                            <span>{{ $t('Basic Information') }}</span>
                        </v-tab>
                        <v-tab value="map" :disabled="!transaction.geoLocation" v-if="mapProvider">
                            <span>{{ $t('Location on Map') }}</span>
                        </v-tab>
                    </v-tabs>
                </div>

                <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container ml-md-5"
                          v-model="activeTab">
                    <v-window-item value="basicInfo">
                        <v-form class="mt-2">
                            <v-row>
                                <v-col cols="12" :md="transaction.type === allTransactionTypes.Transfer ? 6 : 12">
                                    <amount-input class="transaction-edit-amount font-weight-bold"
                                                  :color="sourceAmountColor"
                                                  :readonly="mode === 'view'"
                                                  :disabled="loading || submitting"
                                                  :persistent-placeholder="true"
                                                  :hide="transaction.hideAmount"
                                                  :label="$t(sourceAmountName)"
                                                  :placeholder="$t(sourceAmountName)"
                                                  v-model="transaction.sourceAmount"/>
                                </v-col>
                                <v-col cols="12" :md="6" v-if="transaction.type === allTransactionTypes.Transfer">
                                    <amount-input class="transaction-edit-amount font-weight-bold" color="primary"
                                                  :readonly="mode === 'view'"
                                                  :disabled="loading || submitting"
                                                  :persistent-placeholder="true"
                                                  :hide="transaction.hideAmount"
                                                  :label="$t('Transfer In Amount')"
                                                  :placeholder="$t('Transfer In Amount')"
                                                  v-model="transaction.destinationAmount"/>
                                </v-col>
                                <v-col cols="12" md="12">
                                    <v-text-field
                                        disabled
                                        persistent-placeholder
                                        :label="$t('Category')"
                                        :placeholder="$t('Category')" />
                                </v-col>
                                <v-col cols="12" :md="transaction.type === allTransactionTypes.Transfer ? 6 : 12">
                                    <v-select
                                        item-title="name"
                                        item-value="id"
                                        persistent-placeholder
                                        :readonly="mode === 'view'"
                                        :disabled="loading || submitting || !allVisibleAccounts.length"
                                        :label="$t(sourceAccountTitle)"
                                        :placeholder="$t(sourceAccountTitle)"
                                        :items="allVisibleAccounts"
                                        :no-data-text="$t('No results')"
                                        v-model="transaction.sourceAccountId"
                                    >
                                        <template #selection="{ item }">
                                            <v-label class="cursor-pointer" v-if="item && item.value !== 0 && item.value !== '0'">
                                                <ItemIcon class="mr-2" icon-type="account" size="23px"
                                                          :icon-id="getAccountNameByKeyValue(transaction.sourceAccountId, 'id', 'icon')"
                                                          :color="getAccountNameByKeyValue(transaction.sourceAccountId, 'id', 'color')"
                                                          v-if="getAccountNameByKeyValue(transaction.sourceAccountId, 'id', 'icon')" />
                                                <span>{{ getAccountNameByKeyValue(transaction.sourceAccountId, 'id', 'name', $t('Not Specified')) }}</span>
                                            </v-label>
                                            <v-label v-if="!item || item.value === 0 || item.value === '0'">{{ $t('Not Specified') }}</v-label>
                                        </template>

                                        <template #item="{ props, item }">
                                            <v-list-item :value="item.value" v-bind="props">
                                                <template #title>
                                                    <v-list-item-title>
                                                        <div class="d-flex align-center">
                                                            <ItemIcon icon-type="account"
                                                                      :icon-id="item.raw.icon" :color="item.raw.color"
                                                                      v-if="item.raw" />
                                                            <span class="ml-2">{{ item.title }}</span>
                                                        </div>
                                                    </v-list-item-title>
                                                </template>
                                            </v-list-item>
                                        </template>
                                    </v-select>
                                </v-col>
                                <v-col cols="12" md="6" v-if="transaction.type === allTransactionTypes.Transfer">
                                    <v-select
                                        item-title="name"
                                        item-value="id"
                                        persistent-placeholder
                                        :readonly="mode === 'view'"
                                        :disabled="loading || submitting || !allVisibleAccounts.length"
                                        :label="$t('Destination Account')"
                                        :placeholder="$t('Destination Account')"
                                        :items="allVisibleAccounts"
                                        :no-data-text="$t('No results')"
                                        v-model="transaction.destinationAccountId"
                                    >
                                        <template #selection="{ item }">
                                            <v-label class="cursor-pointer" v-if="item && item.value !== 0 && item.value !== '0'">
                                                <ItemIcon class="mr-2" icon-type="account" size="23px"
                                                          :icon-id="getAccountNameByKeyValue(transaction.destinationAccountId, 'id', 'icon')"
                                                          :color="getAccountNameByKeyValue(transaction.destinationAccountId, 'id', 'color')"
                                                          v-if="getAccountNameByKeyValue(transaction.destinationAccountId, 'id', 'icon')" />
                                                <span>{{ getAccountNameByKeyValue(transaction.destinationAccountId, 'id', 'name', $t('Not Specified')) }}</span>
                                            </v-label>
                                            <v-label v-if="!item || item.value === 0 || item.value === '0'">{{ $t('Not Specified') }}</v-label>
                                        </template>

                                        <template #item="{ props, item }">
                                            <v-list-item :value="item.value" v-bind="props">
                                                <template #title>
                                                    <v-list-item-title>
                                                        <div class="d-flex align-center">
                                                            <ItemIcon icon-type="account"
                                                                      :icon-id="item.raw.icon" :color="item.raw.color"
                                                                      v-if="item.raw" />
                                                            <span class="ml-2">{{ item.title }}</span>
                                                        </div>
                                                    </v-list-item-title>
                                                </template>
                                            </v-list-item>
                                        </template>
                                    </v-select>
                                </v-col>
                                <v-col cols="12" md="6">
                                    <date-time-select
                                        :readonly="mode === 'view'"
                                        :disabled="loading || submitting"
                                        :label="$t('Transaction Time')"
                                        v-model="transaction.time"
                                        @error="showDateTimeError" />
                                </v-col>
                                <v-col cols="12" md="6">
                                    <v-autocomplete
                                        class="transaction-edit-timezone"
                                        item-title="displayNameWithUtcOffset"
                                        item-value="name"
                                        auto-select-first
                                        persistent-placeholder
                                        :readonly="mode === 'view'"
                                        :disabled="loading || submitting"
                                        :label="$t('Timezone')"
                                        :placeholder="!transaction.timeZone && transaction.timeZone !== '' ? `(${transactionDisplayTimezone}) ${transactionTimezoneTimeDifference}` : $t('Timezone')"
                                        :items="allTimezones"
                                        :no-data-text="$t('No results')"
                                        v-model="transaction.timeZone"
                                    >
                                        <template #selection="{ item }">
                                            <span v-if="transaction.timeZone || transaction.timeZone === ''">
                                                {{ item.title }}
                                            </span>
                                        </template>
                                    </v-autocomplete>
                                </v-col>
                                <v-col cols="12" md="12">
                                    <v-select
                                        persistent-placeholder
                                        :readonly="mode === 'view'"
                                        :disabled="loading || submitting"
                                        :label="$t('Geographic Location')"
                                        v-model="transaction"
                                    >
                                        <template #selection>
                                            <span class="cursor-pointer" v-if="transaction.geoLocation">{{ `(${transaction.geoLocation.longitude}, ${transaction.geoLocation.latitude})` }}</span>
                                            <span class="cursor-pointer" v-else-if="!transaction.geoLocation">{{ geoLocationStatusInfo }}</span>
                                        </template>

                                        <template #no-data>
                                            <v-list>
                                                <v-list-item v-if="mode !== 'view'" @click="updateGeoLocation(true)">{{ $t('Update Geographic Location') }}</v-list-item>
                                                <v-list-item v-if="mode !== 'view'" @click="clearGeoLocation">{{ $t('Clear Geographic Location') }}</v-list-item>
                                            </v-list>
                                        </template>
                                    </v-select>
                                </v-col>
                                <v-col cols="12" md="12">
                                    <v-autocomplete
                                        item-title="name"
                                        item-value="id"
                                        auto-select-first
                                        persistent-placeholder
                                        multiple
                                        chips
                                        :closable-chips="mode !== 'view'"
                                        :readonly="mode === 'view'"
                                        :disabled="loading || submitting"
                                        :label="$t('Tags')"
                                        :placeholder="$t('None')"
                                        :items="allTags"
                                        :no-data-text="$t('No available tag')"
                                        v-model="transaction.tagIds"
                                    >
                                        <template #chip="{ props, item }">
                                            <v-chip :prepend-icon="icons.tag" :text="item.title" v-bind="props"/>
                                        </template>

                                        <template #item="{ props, item }">
                                            <v-list-item :value="item.value" v-bind="props">
                                                <template #title>
                                                    <v-list-item-title>
                                                        <div class="d-flex align-center">
                                                            <v-icon size="20" start :icon="icons.tag"/>
                                                            <span>{{ item.title }}</span>
                                                        </div>
                                                    </v-list-item-title>
                                                </template>
                                            </v-list-item>
                                        </template>
                                    </v-autocomplete>
                                </v-col>
                                <v-col cols="12" md="12">
                                    <v-textarea
                                        type="text"
                                        persistent-placeholder
                                        rows="3"
                                        :readonly="mode === 'view'"
                                        :disabled="loading || submitting"
                                        :label="$t('Description')"
                                        :placeholder="$t('Your transaction description (optional)')"
                                        v-model="transaction.comment"
                                    />
                                </v-col>
                            </v-row>
                        </v-form>
                    </v-window-item>
                    <v-window-item value="map">
                        <v-row>
                            <v-col cols="12" md="12">
                                <map-view ref="map" map-class="transaction-edit-map-view" :geo-location="transaction.geoLocation">
                                    <template #error-title="{ mapSupported, mapDependencyLoaded }">
                                        <span class="text-subtitle-1" v-if="!mapSupported"><b>{{ $t('Unsupported Map Provider') }}</b></span>
                                        <span class="text-subtitle-1" v-else-if="!mapDependencyLoaded"><b>{{ $t('Cannot Initialize Map') }}</b></span>
                                    </template>
                                    <template #error-content>
                                        <p class="text-body-1">
                                            {{ $t('Please refresh the page and try again. If the error is still displayed, make sure that server map settings are set correctly.') }}
                                        </p>
                                    </template>
                                </map-view>
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
                           :disabled="loading || submitting" @click="cancel">{{ $t(cancelButtonTitle) }}</v-btn>
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
import { setTransactionModelByTransaction } from '@/lib/transaction.js';
import { getMapProvider } from '@/lib/server_settings.js';

import {
    mdiDotsVertical,
    mdiEyeOffOutline,
    mdiEyeOutline,
    mdiPound
} from '@mdi/js';

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
                more: mdiDotsVertical,
                show: mdiEyeOutline,
                hide: mdiEyeOffOutline,
                tag: mdiPound
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
        cancelButtonTitle() {
            if (this.mode === 'view') {
                return 'Close';
            } else {
                return 'Cancel';
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
            return this.transactionTagsStore.allVisibleTags;
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
        transactionTimezoneTimeDifference() {
            return this.$locale.getTimezoneDifferenceDisplayText(this.transaction.utcOffset);
        },
        sourceAmountColor() {
            if (this.transaction.type === this.allTransactionTypes.Expense) {
                return 'expense';
            } else if (this.transaction.type === this.allTransactionTypes.Income) {
                return 'income';
            } else if (this.transaction.type === this.allTransactionTypes.Transfer) {
                return 'primary';
            }
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
        'activeTab': function (newValue) {
            if (newValue === 'map') {
                this.$nextTick(() => {
                    if (this.$refs.map) {
                        this.$refs.map.init();
                    }
                });
            }
        },
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
            self.setTransaction(newTransaction, options, true);

            const promises = [
                self.accountsStore.loadAllAccounts({ force: false }),
                self.transactionCategoriesStore.loadAllCategories({ force: false }),
                self.transactionTagsStore.loadAllTags({ force: false })
            ];

            if (options && options.id) {
                if (options.currentTransaction) {
                    self.setTransaction(options.currentTransaction, options, true);
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

                if (options.id && responses[3]) {
                    const transaction = responses[3];
                    self.setTransaction(transaction, options, true);
                } else {
                    self.setTransaction(null, options, true);
                }

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
        showDateTimeError(error) {
            this.$refs.snackbar.showError(error);
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
        clearGeoLocation() {
            this.geoLocationStatus = null;
            this.transaction.geoLocation = null;
        },
        getAccountNameByKeyValue(src, value, keyField, nameField, defaultName) {
            return getNameByKeyValue(this.allAccounts, src, value, keyField, nameField, defaultName);
        },
        setTransaction(transaction, options, setContextData) {
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
                setContextData
            );
        }
    }
}
</script>

<style>
.transaction-edit-amount .v-field__input > input {
    font-size: 1.25rem;
}

.transaction-edit-timezone.v-input input::placeholder {
    color: rgba(var(--v-theme-on-background), var(--v-high-emphasis-opacity)) !important;
}

.transaction-edit-map-view {
    height: 220px;
}

@media (min-height: 630px) {
    .transaction-edit-map-view {
        height: 300px;
    }
}

@media (min-height: 700px) {
    .transaction-edit-map-view {
        height: 350px;
    }
}

@media (min-height: 800px) {
    .transaction-edit-map-view {
        height: 450px;
    }
}

@media (min-height: 900px) {
    .transaction-edit-map-view {
        height: 550px;
    }
}
</style>
