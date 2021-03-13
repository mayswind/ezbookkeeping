<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t(title)"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" v-if="!editAccountId && account.type === $constants.account.allAccountTypes.MultiSubAccounts" @click="showMoreActionSheet = true"></f7-link>
                <f7-link :class="{ 'disabled': isInputEmpty() || submitting }" :text="$t(saveButtonTitle)" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item link="#" title="Account Category" after="Category"></f7-list-item>
                    <f7-list-item link="#" title="Account Type" after="Account Type"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list form>
                    <f7-list-item
                        link="#"
                        :title="$t('Account Category')"
                        :after="account.category | accountCategoryName(allAccountCategories) | localized"
                        @click="showAccountCategorySheet = true"
                    >
                        <list-item-selection-sheet value-type="item"
                                                   key-field="id" value-field="id" title-field="name"
                                                   icon-field="defaultAccountIconId" icon-type="account"
                                                   :title-i18n="true"
                                                   :items="allAccountCategories"
                                                   :show.sync="showAccountCategorySheet"
                                                   v-model="account.category">
                        </list-item-selection-sheet>
                    </f7-list-item>

                    <f7-list-item
                        link="#"
                        :class="{ 'disabled': editAccountId }"
                        :title="$t('Account Type')"
                        :after="account.type | accountTypeName | localized"
                        @click="showAccountTypeSheet = true"
                    >
                        <list-item-selection-sheet value-type="item"
                                                   key-field="id" value-field="id" title-field="name"
                                                   :items="[{ id: 1, name: 'Single Account' }, { id: 2, name: 'Multi Sub Accounts' }]"
                                                   :title-i18n="true"
                                                   :show.sync="showAccountTypeSheet"
                                                   v-model="account.type">
                        </list-item-selection-sheet>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-input inline-label label="Account Name" placeholder="Your account name"></f7-list-input>
                    <f7-list-item title="Account Icon" link="#">
                        <f7-icon f7="app_fill"></f7-icon>
                    </f7-list-item>
                    <f7-list-item title="Account Color" link="#">
                        <f7-icon f7="app_fill"></f7-icon>
                    </f7-list-item>
                    <f7-list-item title="Currency" after="Currency" link="#"></f7-list-item>
                    <f7-list-item title="Account Balance" after="Balance" link="#"></f7-list-item>
                    <f7-list-item title="Visible" after="Visible"></f7-list-item>
                    <f7-list-input type="textarea" placeholder="Your account description (optional)"></f7-list-input>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading && account.type === $constants.account.allAccountTypes.SingleAccount">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list form>
                    <f7-list-input
                        type="text"
                        inline-label
                        clear-button
                        :label="$t('Account Name')"
                        :placeholder="$t('Your account name')"
                        :value="account.name"
                        @input="account.name = $event.target.value"
                    ></f7-list-input>

                    <f7-list-item :title="$t('Account Icon')" link="#"
                                  @click="account.showIconSelectionSheet = true">
                        <f7-icon slot="after"
                                 :icon="account.icon | accountIcon"
                                 :style="account.color | accountIconStyle('var(--default-icon-color)')"></f7-icon>
                        <icon-selection-sheet :all-icon-infos="allAccountIcons"
                                              :show.sync="account.showIconSelectionSheet"
                                              :color="account.color"
                                              v-model="account.icon"
                        ></icon-selection-sheet>
                    </f7-list-item>

                    <f7-list-item :title="$t('Account Color')" link="#"
                                  @click="account.showColorSelectionSheet = true">
                        <f7-icon slot="after"
                                 f7="app_fill"
                                 :style="account.color | accountIconStyle('var(--default-icon-color)')"></f7-icon>
                        <color-selection-sheet :all-color-infos="allAccountColors"
                                               :show.sync="account.showColorSelectionSheet"
                                               v-model="account.color"
                        ></color-selection-sheet>
                    </f7-list-item>

                    <f7-list-item
                        :class="{ 'disabled': editAccountId }"
                        :title="$t('Currency')"
                        smart-select :smart-select-params="{ openIn: 'popup', searchbar: true, searchbarPlaceholder: $t('Currency Name'), searchbarDisableText: $t('Cancel'), closeOnSelect: true, popupCloseLinkText: $t('Done'), scrollToSelectedItem: true }"
                    >
                        <select autocomplete="transaction-currency" v-model="account.currency">
                            <option v-for="currency in allCurrencies"
                                    :key="currency.code"
                                    :value="currency.code">{{ currency.displayName }}</option>
                        </select>
                    </f7-list-item>

                    <f7-list-item
                        link="#"
                        :class="{ 'disabled': editAccountId }"
                        :title="$t('Account Balance')"
                        :after="account.balance | currency(account.currency)"
                        @click="account.showBalanceSheet = true"
                    >
                        <number-pad-sheet :min-value="$constants.transaction.minAmount"
                                          :max-value="$constants.transaction.maxAmount"
                                          :show.sync="account.showBalanceSheet"
                                          v-model="account.balance"
                        ></number-pad-sheet>
                    </f7-list-item>

                    <f7-list-item :title="$t('Visible')" v-if="editAccountId">
                        <f7-toggle :checked="account.visible" @toggle:change="account.visible = $event"></f7-toggle>
                    </f7-list-item>

                    <f7-list-input
                        type="textarea"
                        class="textarea-auto-size"
                        style="height: auto"
                        :placeholder="$t('Your account description (optional)')"
                        :value="account.comment"
                        @input="account.comment = $event.target.value"
                    ></f7-list-input>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading && account.type === $constants.account.allAccountTypes.MultiSubAccounts">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list form>
                    <f7-list-input
                        type="text"
                        inline-label
                        clear-button
                        :label="$t('Account Name')"
                        :placeholder="$t('Your account name')"
                        :value="account.name"
                        @input="account.name = $event.target.value"
                    ></f7-list-input>

                    <f7-list-item :title="$t('Account Icon')" link="#"
                                  @click="account.showIconSelectionSheet = true">
                        <f7-icon slot="after"
                                 :icon="account.icon | accountIcon"
                                 :style="account.color | accountIconStyle('var(--default-icon-color)')"></f7-icon>
                        <icon-selection-sheet :all-icon-infos="allAccountIcons"
                                              :show.sync="account.showIconSelectionSheet"
                                              :color="account.color"
                                              v-model="account.icon"
                        ></icon-selection-sheet>
                    </f7-list-item>

                    <f7-list-item :title="$t('Account Color')" link="#"
                                  @click="account.showColorSelectionSheet = true">
                        <f7-icon slot="after"
                                 f7="app_fill"
                                 :style="account.color | accountIconStyle('var(--default-icon-color)')"></f7-icon>
                        <color-selection-sheet :all-color-infos="allAccountColors"
                                               :show.sync="account.showColorSelectionSheet"
                                               v-model="account.color"
                        ></color-selection-sheet>
                    </f7-list-item>

                    <f7-list-item :title="$t('Visible')" v-if="editAccountId">
                        <f7-toggle :checked="account.visible" @toggle:change="account.visible = $event"></f7-toggle>
                    </f7-list-item>

                    <f7-list-input
                        type="textarea"
                        class="textarea-auto-size"
                        style="height: auto"
                        :placeholder="$t('Your account description (optional)')"
                        :value="account.comment"
                        @input="account.comment = $event.target.value"
                    ></f7-list-input>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-block class="no-padding no-margin" v-if="!loading && account.type === $constants.account.allAccountTypes.MultiSubAccounts">
            <f7-card v-for="(subAccount, idx) in subAccounts" :key="idx">
                <f7-card-header v-if="!editAccountId">
                    <f7-button rasied fill color="red"
                               icon-f7="trash" icon-size="16px"
                               :tooltip="$t('Remove Sub Account')"
                               @click="removeSubAccount(subAccount, false)">
                    </f7-button>
                </f7-card-header>
                <f7-card-content class="no-safe-areas" :padding="false">
                    <f7-list>
                        <f7-list-input
                            type="text"
                            inline-label
                            clear-button
                            :label="$t('Sub Account Name')"
                            :placeholder="$t('Your sub account name')"
                            :value="subAccount.name"
                            @input="subAccount.name = $event.target.value"
                        ></f7-list-input>

                        <f7-list-item :title="$t('Sub Account Icon')" link="#"
                                      @click="subAccount.showIconSelectionSheet = true">
                            <f7-icon slot="after"
                                     :icon="subAccount.icon | accountIcon"
                                     :style="subAccount.color | accountIconStyle('var(--default-icon-color)')"></f7-icon>
                            <icon-selection-sheet :all-icon-infos="allAccountIcons"
                                                  :show.sync="subAccount.showIconSelectionSheet"
                                                  :color="subAccount.color"
                                                  v-model="subAccount.icon"
                            ></icon-selection-sheet>
                        </f7-list-item>

                        <f7-list-item :title="$t('Sub Account Color')" link="#"
                                      @click="subAccount.showColorSelectionSheet = true">
                            <f7-icon slot="after"
                                     f7="app_fill"
                                     :style="subAccount.color | accountIconStyle('var(--default-icon-color)')"></f7-icon>
                            <color-selection-sheet :all-color-infos="allAccountColors"
                                                   :show.sync="subAccount.showColorSelectionSheet"
                                                   v-model="subAccount.color"
                            ></color-selection-sheet>
                        </f7-list-item>

                        <f7-list-item
                            :class="{ 'disabled': editAccountId }"
                            :title="$t('Currency')"
                            smart-select :smart-select-params="{ openIn: 'popup', searchbar: true, searchbarPlaceholder: $t('Currency Name'), searchbarDisableText: $t('Cancel'), closeOnSelect: true, popupCloseLinkText: $t('Done'), scrollToSelectedItem: true }"
                        >
                            <select autocomplete="transaction-currency" v-model="subAccount.currency">
                                <option v-for="currency in allCurrencies"
                                        :key="currency.code"
                                        :value="currency.code">{{ currency.displayName }}</option>
                            </select>
                        </f7-list-item>

                        <f7-list-item
                            link="#"
                            :class="{ 'disabled': editAccountId }"
                            :title="$t('Sub Account Balance')"
                            :after="subAccount.balance | currency(subAccount.currency)"
                            @click="subAccount.showBalanceSheet = true"
                        >
                            <number-pad-sheet :min-value="$constants.transaction.minAmount"
                                              :max-value="$constants.transaction.maxAmount"
                                              :show.sync="subAccount.showBalanceSheet"
                                              v-model="subAccount.balance"
                            ></number-pad-sheet>
                        </f7-list-item>

                        <f7-list-item :title="$t('Visible')" v-if="editAccountId">
                            <f7-toggle :checked="subAccount.visible" @toggle:change="subAccount.visible = $event"></f7-toggle>
                        </f7-list-item>

                        <f7-list-input
                            type="textarea"
                            class="textarea-auto-size"
                            style="height: auto"
                            :placeholder="$t('Your sub account description (optional)')"
                            :value="subAccount.comment"
                            @input="subAccount.comment = $event.target.value"
                        ></f7-list-input>
                    </f7-list>
                </f7-card-content>
            </f7-card>
        </f7-block>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="addSubAccount">{{ $t('Add Sub Account') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-actions close-by-outside-click close-on-escape :opened="showDeleteActionSheet" @actions:closed="showDeleteActionSheet = false">
            <f7-actions-group>
                <f7-actions-label>{{ $t('Are you sure you want to remove this sub account?') }}</f7-actions-label>
                <f7-actions-button color="red" @click="removeSubAccount(subAccountToDelete, true)">{{ $t('Remove') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script>
export default {
    data() {
        const self = this;

        return {
            editAccountId: null,
            loading: false,
            loadingError: null,
            account: {
                category: 1,
                type: self.$constants.account.allAccountTypes.SingleAccount,
                name: '',
                icon: self.$constants.icons.defaultAccountIconId,
                color: self.$constants.colors.defaultAccountColor,
                currency: self.$store.getters.currentUserDefaultCurrency || self.$t('default.currency'),
                balance: 0,
                comment: '',
                visible: true,
                showIconSelectionSheet: false,
                showColorSelectionSheet: false,
                showBalanceSheet: false
            },
            subAccounts: [],
            subAccountToDelete: null,
            submitting: false,
            showAccountCategorySheet: false,
            showAccountTypeSheet: false,
            showMoreActionSheet: false,
            showDeleteActionSheet: false
        };
    },
    computed: {
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
        allAccountCategories() {
            return this.$constants.account.allCategories;
        },
        allAccountIcons() {
            return this.$constants.icons.allAccountIcons;
        },
        allAccountColors() {
            return this.$constants.colors.allAccountColors;
        },
        allCurrencies() {
            return this.$locale.getAllCurrencies();
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
    created() {
        const self = this;
        const query = self.$f7route.query;

        if (query.id) {
            self.loading = true;

            self.editAccountId = query.id;

            self.$store.dispatch('getAccount', {
                accountId: self.editAccountId
            }).then(account => {
                self.account.id = account.id;
                self.account.category = account.category;
                self.account.type = account.type;
                self.account.name = account.name;
                self.account.icon = account.icon;
                self.account.color = account.color;
                self.account.currency = account.currency;
                self.account.balance = account.balance;
                self.account.comment = account.comment;
                self.account.visible = !account.hidden;
                self.subAccounts = [];

                if (account.subAccounts && account.subAccounts.length > 0) {
                    for (let i = 0; i < account.subAccounts.length; i++) {
                        const subAccount = account.subAccounts[i];

                        self.subAccounts.push({
                            id: subAccount.id,
                            category: subAccount.category,
                            type: subAccount.type,
                            name: subAccount.name,
                            icon: subAccount.icon,
                            color: subAccount.color,
                            currency: subAccount.currency,
                            balance: subAccount.balance,
                            comment: subAccount.comment,
                            visible: !subAccount.hidden,
                            showIconSelectionSheet: false,
                            showColorSelectionSheet: false,
                            showBalanceSheet: false
                        });
                    }
                }

                self.loading = false;
            }).catch(error => {
                if (error.processed) {
                    self.loading = false;
                } else {
                    self.loadingError = error;
                    self.$toast(error.message || error);
                }
            });
        } else {
            self.loading = false;
        }
    },
    updated: function () {
        this.autoChangeCommentTextareaSize();
    },
    methods: {
        onPageAfterIn() {
            this.$routeBackOnError('loadingError');
        },
        addSubAccount() {
            const self = this;

            if (self.account.type !== self.$constants.account.allAccountTypes.MultiSubAccounts) {
                return;
            }

            this.subAccounts.push({
                category: null,
                type: null,
                name: '',
                icon: self.account.icon,
                color: self.account.color,
                currency: self.$store.getters.currentUserDefaultCurrency || self.$t('default.currency'),
                balance: 0,
                comment: '',
                visible: true,
                showIconSelectionSheet: false,
                showColorSelectionSheet: false,
                showBalanceSheet: false
            });
        },
        removeSubAccount(subAccount, confirm) {
            if (!subAccount) {
                this.$alert('An error has occurred');
                return;
            }

            if (!confirm) {
                this.subAccountToDelete = subAccount;
                this.showDeleteActionSheet = true;
                return;
            }

            this.showDeleteActionSheet = false;
            this.subAccountToDelete = null;

            for (let i = 0; i < this.subAccounts.length; i++) {
                if (this.subAccounts[i] === subAccount) {
                    this.subAccounts.splice(i, 1);
                }
            }
        },
        save() {
            const self = this;
            const router = self.$f7router;

            let problemMessage = self.getInputEmptyProblemMessage(self.account, false);

            if (!problemMessage && self.account.type === self.$constants.account.allAccountTypes.MultiSubAccounts) {
                for (let i = 0; i < self.subAccounts.length; i++) {
                    problemMessage = self.getInputEmptyProblemMessage(self.subAccounts[i], true);

                    if (problemMessage) {
                        break;
                    }
                }
            }

            if (problemMessage) {
                self.$alert(problemMessage);
                return;
            }

            self.submitting = true;
            self.$showLoading(() => self.submitting);

            const subAccounts = [];

            if (self.account.type === self.$constants.account.allAccountTypes.MultiSubAccounts) {
                for (let i = 0; i < self.subAccounts.length; i++) {
                    const subAccount = self.subAccounts[i];
                    const submitAccount = {
                        category: self.account.category,
                        type: self.$constants.account.allAccountTypes.SingleAccount,
                        name: subAccount.name,
                        icon: subAccount.icon,
                        color: subAccount.color,
                        currency: subAccount.currency,
                        balance: subAccount.balance,
                        comment: subAccount.comment
                    };

                    if (self.editAccountId) {
                        submitAccount.id = subAccount.id;
                        submitAccount.hidden = !subAccount.visible;
                    }

                    subAccounts.push(submitAccount);
                }
            }

            const submitAccount = {
                category: self.account.category,
                type: self.account.type,
                name: self.account.name,
                icon: self.account.icon,
                color: self.account.color,
                currency: self.account.type === self.$constants.account.allAccountTypes.SingleAccount ? self.account.currency : self.$constants.currency.parentAccountCurrencyPlaceholder,
                balance: self.account.type === self.$constants.account.allAccountTypes.SingleAccount ? self.account.balance : 0,
                comment: self.account.comment,
                subAccounts: self.account.type === self.$constants.account.allAccountTypes.SingleAccount ? null : subAccounts,
            };

            if (self.editAccountId) {
                submitAccount.id = self.account.id;
                submitAccount.hidden = !self.account.visible;
            }

            self.$store.dispatch('saveAccount', {
                account: submitAccount
            }).then(() => {
                self.submitting = false;
                self.$hideLoading();

                if (!self.editAccountId) {
                    self.$toast('You have added a new account');
                } else {
                    self.$toast('You have saved this account');
                }

                router.back();
            }).catch(error => {
                self.submitting = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        autoChangeCommentTextareaSize() {
            const app = this.$f7;
            const $$ = app.$;

            $$('.textarea-auto-size textarea').each((idx, el) => {
                el.scrollTop = 0;
                el.style.height = '';
                el.style.height = el.scrollHeight + 'px';
            });
        },
        chooseSuitableIcon(oldCategory, newCategory) {
            const allCategories = this.$constants.account.allCategories;

            for (let i = 0; i < allCategories.length; i++) {
                if (allCategories[i].id === oldCategory) {
                    if (this.account.icon !== allCategories[i].defaultAccountIconId) {
                        return;
                    } else {
                        break;
                    }
                }
            }

            for (let i = 0; i < allCategories.length; i++) {
                if (allCategories[i].id === newCategory) {
                    this.account.icon = allCategories[i].defaultAccountIconId;
                }
            }
        },
        isInputEmpty() {
            const isAccountEmpty = !!this.getInputEmptyProblemMessage(this.account, false);

            if (isAccountEmpty) {
                return true;
            }

            if (this.account.type === this.$constants.account.allAccountTypes.MultiSubAccounts) {
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
                return 'Account category cannot be empty';
            } else if (!isSubAccount && !account.type) {
                return 'Account type cannot be empty';
            } else if (!account.name) {
                return 'Account name cannot be empty';
            } else if (account.type === this.$constants.account.allAccountTypes.SingleAccount && !account.currency) {
                return 'Account currency cannot be empty';
            } else {
                return null;
            }
        }
    },
    filters: {
        accountCategoryName(categoryId, allCategories) {
            for (let i = 0; i < allCategories.length; i++) {
                if (allCategories[i].id === categoryId) {
                    return allCategories[i].name;
                }
            }

            return '';
        },
        accountTypeName(type) {
            if (type === 1) {
                return 'Single Account';
            } else if (type === 2) {
                return 'Multi Sub Accounts';
            }

            return '';
        }
    }
}
</script>
