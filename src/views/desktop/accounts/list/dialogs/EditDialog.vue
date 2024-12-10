<template>
    <v-dialog :width="account.type === allAccountTypes.MultiSubAccounts ? 1000 : 800" :persistent="!!persistent" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-8">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <div class="d-flex w-100 align-center justify-center">
                        <h4 class="text-h4">{{ $t(title) }}</h4>
                        <v-progress-circular indeterminate size="22" class="ml-2" v-if="loading"></v-progress-circular>
                    </div>
                    <v-btn density="comfortable" color="default" variant="text" class="ml-2" :icon="true"
                           :disabled="loading || submitting || !!editAccountId || account.type !== allAccountTypes.MultiSubAccounts">
                        <v-icon :icon="icons.more" />
                        <v-menu activator="parent">
                            <v-list>
                                <v-list-item :prepend-icon="icons.add"
                                             :title="$t('Add Sub-account')"
                                             @click="addSubAccount"></v-list-item>
                            </v-list>
                        </v-menu>
                    </v-btn>
                </div>
            </template>
            <v-card-text class="d-flex flex-column flex-md-row mt-md-4 pt-0">
                <div class="mb-4" v-if="account.type === allAccountTypes.MultiSubAccounts">
                    <v-tabs direction="vertical" :disabled="loading || submitting" v-model="currentAccountIndex">
                        <v-tab :value="-1">
                            <span>{{ $t('Main Account') }}</span>
                        </v-tab>
                        <template v-if="account.type === allAccountTypes.MultiSubAccounts">
                            <v-tab :key="idx" :value="idx" v-for="(subAccount, idx) in subAccounts">
                                <span>{{ $t('Sub Account') + ' #' + (idx + 1) }}</span>
                                <v-btn class="ml-2" color="error" size="24" variant="text"
                                       :icon="icons.delete" v-if="!editAccountId"
                                       @click="removeSubAccount(subAccount)"></v-btn>
                            </v-tab>
                        </template>
                    </v-tabs>
                </div>

                <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container"
                          :class="{ 'ml-md-5': account.type === allAccountTypes.MultiSubAccounts }"
                          v-model="activeTab">
                    <v-window-item value="account">
                        <v-form class="mt-2">
                            <v-row>
                                <v-col cols="12" md="12" v-if="account.type === allAccountTypes.SingleAccount || currentAccountIndex < 0">
                                    <v-select
                                        item-title="displayName"
                                        item-value="id"
                                        persistent-placeholder
                                        :disabled="loading || submitting"
                                        :label="$t('Account Category')"
                                        :placeholder="$t('Account Category')"
                                        :items="allAccountCategories"
                                        :no-data-text="$t('No results')"
                                        v-model="selectedAccount.category"
                                    >
                                        <template #item="{ props, item }">
                                            <v-list-item :value="item.value" v-bind="props">
                                                <template #title>
                                                    <v-list-item-title>
                                                        <div class="d-flex align-center">
                                                            <ItemIcon icon-type="account"
                                                                      :icon-id="item.raw.defaultAccountIconId"
                                                                      v-if="item.raw" />
                                                            <span class="ml-2">{{ item.title }}</span>
                                                        </div>
                                                    </v-list-item-title>
                                                </template>
                                            </v-list-item>
                                        </template>
                                    </v-select>
                                </v-col>
                                <v-col cols="12" md="12" v-if="account.type === allAccountTypes.SingleAccount || currentAccountIndex < 0">
                                    <v-select
                                        item-title="displayName"
                                        item-value="id"
                                        persistent-placeholder
                                        :disabled="loading || submitting || !!editAccountId"
                                        :label="$t('Account Type')"
                                        :placeholder="$t('Account Type')"
                                        :items="allAccountTypesArray"
                                        :no-data-text="$t('No results')"
                                        v-model="selectedAccount.type"
                                    />
                                </v-col>
                                <v-col cols="12" md="12">
                                    <v-text-field
                                        type="text"
                                        persistent-placeholder
                                        :disabled="loading || submitting"
                                        :label="currentAccountIndex < 0 ? $t('Account Name') : $t('Sub-account Name')"
                                        :placeholder="currentAccountIndex < 0 ? $t('Your account name') : $t('Your sub-account name')"
                                        v-model="selectedAccount.name"
                                    />
                                </v-col>
                                <v-col cols="12" md="6">
                                    <icon-select icon-type="account"
                                                 :all-icon-infos="allAccountIcons"
                                                 :label="currentAccountIndex < 0 ? $t('Account Icon') : $t('Sub-account Icon')"
                                                 :color="selectedAccount.color"
                                                 :disabled="loading || submitting"
                                                 v-model="selectedAccount.icon" />
                                </v-col>
                                <v-col cols="12" md="6">
                                    <color-select :all-color-infos="allAccountColors"
                                                  :label="currentAccountIndex < 0 ? $t('Account Color') : $t('Sub-account Color')"
                                                  :disabled="loading || submitting"
                                                  v-model="selectedAccount.color" />
                                </v-col>
                                <v-col cols="12" :md="isAccountSupportCreditCardStatementDate() ? 6 : 12" v-if="account.type === allAccountTypes.SingleAccount || currentAccountIndex >= 0">
                                    <v-autocomplete
                                        item-title="displayName"
                                        item-value="currencyCode"
                                        auto-select-first
                                        persistent-placeholder
                                        :disabled="loading || submitting || !!editAccountId"
                                        :label="$t('Currency')"
                                        :placeholder="$t('Currency')"
                                        :items="allCurrencies"
                                        :no-data-text="$t('No results')"
                                        v-model="selectedAccount.currency"
                                    >
                                        <template #append-inner>
                                            <small class="text-field-append-text smaller">{{ selectedAccount.currency }}</small>
                                        </template>
                                    </v-autocomplete>
                                </v-col>
                                <v-col cols="12" :md="account.type === allAccountTypes.SingleAccount || currentAccountIndex >= 0 ? 6 : 12" v-if="isAccountSupportCreditCardStatementDate()">
                                    <v-autocomplete
                                        item-title="displayName"
                                        item-value="day"
                                        auto-select-first
                                        persistent-placeholder
                                        :disabled="loading || submitting"
                                        :label="$t('Statement Date')"
                                        :placeholder="$t('Statement Date')"
                                        :items="allAvailableMonthDays"
                                        :no-data-text="$t('No results')"
                                        v-model="selectedAccount.creditCardStatementDate"
                                    ></v-autocomplete>
                                </v-col>
                                <v-col cols="12" :md="!editAccountId && selectedAccount.balance ? 6 : 12"
                                       v-if="account.type === allAccountTypes.SingleAccount || currentAccountIndex >= 0">
                                    <amount-input :disabled="loading || submitting || !!editAccountId"
                                                  :persistent-placeholder="true"
                                                  :currency="selectedAccount.currency"
                                                  :show-currency="true"
                                                  :label="currentAccountIndex < 0 ? $t('Account Balance') : $t('Sub-account Balance')"
                                                  :placeholder="currentAccountIndex < 0 ? $t('Account Balance') : $t('Sub-account Balance')"
                                                  v-model="selectedAccount.balance"/>
                                </v-col>
                                <v-col cols="12" md="6" v-show="selectedAccount.balance"
                                       v-if="!editAccountId && (account.type === allAccountTypes.SingleAccount || currentAccountIndex >= 0)">
                                    <date-time-select
                                        :disabled="loading || submitting"
                                        :label="$t('Balance Time')"
                                        v-model="selectedAccount.balanceTime"
                                        @error="showDateTimeError" />
                                </v-col>
                                <v-col cols="12" md="12">
                                    <v-textarea
                                        type="text"
                                        persistent-placeholder
                                        rows="3"
                                        :disabled="loading || submitting"
                                        :label="$t('Description')"
                                        :placeholder="currentAccountIndex < 0 ? $t('Your account description (optional)') : $t('Your sub-account description (optional)')"
                                        v-model="selectedAccount.comment"
                                    />
                                </v-col>
                                <v-col class="py-0" cols="12" md="12" v-if="editAccountId">
                                    <v-switch :disabled="loading || submitting"
                                              :label="$t('Visible')" v-model="selectedAccount.visible"/>
                                </v-col>
                            </v-row>
                        </v-form>
                    </v-window-item>
                </v-window>
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center mt-2 mt-sm-4 mt-md-6 gap-4">
                    <v-btn :disabled="isInputEmpty() || loading || submitting" @click="save">
                        {{ $t(saveButtonTitle) }}
                        <v-progress-circular indeterminate size="22" class="ml-2" v-if="submitting"></v-progress-circular>
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
import { useAccountsStore } from '@/stores/account.js';

import accountConstants from '@/consts/account.js';
import iconConstants from '@/consts/icon.js';
import colorConstants from '@/consts/color.js';
import { isNumber } from '@/lib/common.js';
import { generateRandomUUID } from '@/lib/misc.js';
import {
    setAccountModelByAnotherAccount,
    setAccountSuitableIcon
} from '@/lib/account.js';

import {
    mdiDotsVertical,
    mdiCreditCardPlusOutline,
    mdiDeleteOutline
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
        const accountsStore = useAccountsStore();
        const newAccount = accountsStore.generateNewAccountModel();

        return {
            showState: false,
            activeTab: 'account',
            editAccountId: null,
            clientSessionId: '',
            loading: false,
            account: newAccount,
            subAccounts: [],
            currentAccountIndex: -1,
            submitting: false,
            resolve: null,
            reject: null,
            icons: {
                more: mdiDotsVertical,
                add: mdiCreditCardPlusOutline,
                delete: mdiDeleteOutline
            }
        };
    },
    computed: {
        ...mapStores(useSettingsStore, useAccountsStore),
        title() {
            if (!this.editAccountId) {
                return 'Add Account';
            } else {
                return 'Edit Account';
            }
        },
        saveButtonTitle() {
            if (!this.editAccountId) {
                return 'Add';
            } else {
                return 'Save';
            }
        },
        allAccountTypes() {
            return accountConstants.allAccountTypes;
        },
        allAccountCategories() {
            return this.$locale.getAllAccountCategories();
        },
        allAccountTypesArray() {
            return this.$locale.getAllAccountTypes();
        },
        allAccountIcons() {
            return iconConstants.allAccountIcons;
        },
        allAccountColors() {
            return colorConstants.allAccountColors;
        },
        allCurrencies() {
            return this.$locale.getAllCurrencies();
        },
        allAvailableMonthDays() {
            const allAvailableDays = [];

            allAvailableDays.push({
                day: 0,
                displayName: this.$t('Not set'),
            });

            for (let i = 1; i <= 28; i++) {
                allAvailableDays.push({
                    day: i,
                    displayName: this.$locale.getMonthdayShortName(i),
                });
            }

            return allAvailableDays;
        },
        selectedAccount() {
            if (this.currentAccountIndex < 0) {
                return this.account;
            }

            return this.subAccounts[this.currentAccountIndex];
        }
    },
    watch: {
        'account.category': function (newValue, oldValue) {
            this.chooseSuitableIcon(oldValue, newValue);
        },
        'account.type': function () {
            if (this.subAccounts.length < 1) {
                this.addSubAccount();
            }
        }
    },
    methods: {
        open(options) {
            const self = this;
            self.showState = true;
            self.loading = true;
            self.submitting = false;

            const newAccount = self.accountsStore.generateNewAccountModel();
            setAccountModelByAnotherAccount(self.account, newAccount);
            self.subAccounts = [];
            self.currentAccountIndex = -1;

            if (options && options.id) {
                if (options.currentAccount) {
                    self.setAccount(options.currentAccount);
                }

                self.editAccountId = options.id;
                self.accountsStore.getAccount({
                    accountId: self.editAccountId
                }).then(account => {
                    self.setAccount(account);
                    self.loading = false;
                }).catch(error => {
                    self.loading = false;
                    self.showState = false;

                    if (!error.processed) {
                        if (self.reject) {
                            self.reject(error);
                        }
                    }
                });
            } else {
                if (isNumber(options.category)) {
                    self.account.category = options.category;
                    self.chooseSuitableIcon(1, options.category);
                }

                self.editAccountId = null;
                self.clientSessionId = generateRandomUUID();
                self.loading = false;
            }

            return new Promise((resolve, reject) => {
                self.resolve = resolve;
                self.reject = reject;
            });
        },
        addSubAccount() {
            if (this.account.type !== this.allAccountTypes.MultiSubAccounts) {
                return;
            }

            const subAccount = this.accountsStore.generateNewSubAccountModel(this.account);
            this.subAccounts.push(subAccount);
        },
        removeSubAccount(subAccount) {
            const self = this;

            self.$refs.confirmDialog.open('Are you sure you want to remove this sub-account?').then(() => {
                for (let i = 0; i < self.subAccounts.length; i++) {
                    if (self.subAccounts[i] === subAccount) {
                        self.subAccounts.splice(i, 1);

                        if (self.currentAccountIndex >= self.subAccounts.length) {
                            self.currentAccountIndex = self.subAccounts.length - 1;
                        }
                    }
                }
            });
        },
        save() {
            const self = this;

            let problemMessage = self.getInputEmptyProblemMessage(self.account, false);

            if (!problemMessage && self.account.type === self.allAccountTypes.MultiSubAccounts) {
                for (let i = 0; i < self.subAccounts.length; i++) {
                    problemMessage = self.getInputEmptyProblemMessage(self.subAccounts[i], true);

                    if (problemMessage) {
                        break;
                    }
                }
            }

            if (problemMessage) {
                self.$refs.snackbar.showMessage(problemMessage);
                return;
            }

            self.submitting = true;

            self.accountsStore.saveAccount({
                account: self.account,
                subAccounts: self.subAccounts,
                isEdit: !!self.editAccountId,
                clientSessionId: self.clientSessionId
            }).then(() => {
                self.submitting = false;

                let message = 'You have saved this account';

                if (!self.editAccountId) {
                    message = 'You have added a new account';
                }

                if (self.resolve) {
                    self.resolve({
                        message: message
                    });
                }

                self.showState = false;
            }).catch(error => {
                self.submitting = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        cancel() {
            if (this.reject) {
                this.reject();
            }

            this.showState = false;
        },
        isAccountSupportCreditCardStatementDate() {
            return this.account && this.account.category === accountConstants.creditCardCategoryType;
        },
        chooseSuitableIcon(oldCategory, newCategory) {
            setAccountSuitableIcon(this.account, oldCategory, newCategory);
        },
        isInputEmpty() {
            const isAccountEmpty = !!this.getInputEmptyProblemMessage(this.account, false);

            if (isAccountEmpty) {
                return true;
            }

            if (this.account.type === this.allAccountTypes.MultiSubAccounts) {
                for (let i = 0; i < this.subAccounts.length; i++) {
                    const isSubAccountEmpty = !!this.getInputEmptyProblemMessage(this.subAccounts[i], true);

                    if (isSubAccountEmpty) {
                        return true;
                    }
                }
            }

            return false;
        },
        getInputEmptyProblemMessage(account, isSubAccount) {
            if (!isSubAccount && !account.category) {
                return 'Account category cannot be blank';
            } else if (!isSubAccount && !account.type) {
                return 'Account type cannot be blank';
            } else if (!account.name) {
                return 'Account name cannot be blank';
            } else if (account.type === this.allAccountTypes.SingleAccount && !account.currency) {
                return 'Account currency cannot be blank';
            } else {
                return null;
            }
        },
        setAccount(account) {
            setAccountModelByAnotherAccount(this.account, account);
            this.subAccounts = [];

            if (account.subAccounts && account.subAccounts.length > 0) {
                for (let i = 0; i < account.subAccounts.length; i++) {
                    const subAccount = this.accountsStore.generateNewSubAccountModel(this.account);
                    setAccountModelByAnotherAccount(subAccount, account.subAccounts[i]);

                    this.subAccounts.push(subAccount);
                }
            }
        },
        showDateTimeError(error) {
            this.$refs.snackbar.showError(error);
        }
    }
}
</script>
