<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t(title)"></f7-nav-title>
            <f7-nav-right>
                <f7-link :class="{ 'disabled': inputIsEmpty || submitting }" :text="$t(saveButtonTitle)" @click="save" v-if="mode !== 'view'"></f7-link>
            </f7-nav-right>

            <f7-subnavbar :class="{ 'disabled': mode !== 'add' }">
                <f7-segmented strong>
                    <f7-button :text="$t('Expense')" :active="transaction.type === $constants.transaction.allTransactionTypes.Expense" @click="transaction.type = $constants.transaction.allTransactionTypes.Expense"></f7-button>
                    <f7-button :text="$t('Income')" :active="transaction.type === $constants.transaction.allTransactionTypes.Income" @click="transaction.type = $constants.transaction.allTransactionTypes.Income"></f7-button>
                    <f7-button :text="$t('Transfer')" :active="transaction.type === $constants.transaction.allTransactionTypes.Transfer" @click="transaction.type = $constants.transaction.allTransactionTypes.Transfer"></f7-button>
                </f7-segmented>
            </f7-subnavbar>
        </f7-navbar>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item
                        class="transaction-edit-amount padding-top-half padding-bottom-half"
                        style="font-size: 40px;"
                        header="Expense Amount" title="0.00">
                    </f7-list-item>
                    <f7-list-item class="transaction-edit-category" header="Category" title="Category Names" link="#"></f7-list-item>
                    <f7-list-item class="transaction-edit-account" header="Account" title="Account Name" link="#"></f7-list-item>
                    <f7-list-input label="Transaction Time" placeholder="YYYY/MM/DD HH:mm"></f7-list-input>
                    <f7-list-item header="Tags" link="#">
                        <f7-block class="margin-top-half no-padding" slot="footer">
                            <f7-chip class="transaction-edit-tag" text="None"></f7-chip>
                        </f7-block>
                    </f7-list-item>
                    <f7-list-input type="textarea" label="Description" placeholder="Your transaction description (optional)"></f7-list-input>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list form :class="{ 'disabled': mode === 'view' }">
                    <f7-list-item
                        class="transaction-edit-amount padding-top-half padding-bottom-half"
                        :class="{ 'color-theme-teal': transaction.type === $constants.transaction.allTransactionTypes.Expense, 'color-theme-red': transaction.type === $constants.transaction.allTransactionTypes.Income }"
                        :style="{ fontSize: sourceAmountFontSize + 'px' }"
                        :header="$t(sourceAmountName)"
                        :title="transaction.sourceAmount | currency"
                        @click="showSourceAmountSheet = true"
                    >
                        <number-pad-sheet :min-value="$constants.transaction.minAmount"
                                          :max-value="$constants.transaction.maxAmount"
                                          :show.sync="showSourceAmountSheet"
                                          v-model="transaction.sourceAmount"
                        ></number-pad-sheet>
                    </f7-list-item>

                    <f7-list-item
                        class="transaction-edit-amount padding-top-half padding-bottom-half"
                        :style="{ fontSize: destinationAmountFontSize + 'px' }"
                        :header="$t('Transfer In Amount')"
                        :title="transaction.destinationAmount | currency"
                        @click="showDestinationAmountSheet = true"
                        v-if="transaction.type === $constants.transaction.allTransactionTypes.Transfer"
                    >
                        <number-pad-sheet :min-value="$constants.transaction.minAmount"
                                          :max-value="$constants.transaction.maxAmount"
                                          :show.sync="showDestinationAmountSheet"
                                          v-model="transaction.destinationAmount"
                        ></number-pad-sheet>
                    </f7-list-item>

                    <f7-list-item
                        class="transaction-edit-category"
                        key="expenseCategorySelection"
                        link="#"
                        :class="{ 'disabled': !hasAvailableExpenseCategories }"
                        :header="$t('Category')"
                        @click="showCategorySheet = true"
                        v-if="transaction.type === $constants.transaction.allTransactionTypes.Expense"
                    >
                        <div slot="title">
                            <span>{{ transaction.expenseCategory | primaryCategoryName(allCategories[$constants.category.allCategoryTypes.Expense]) }}</span>
                            <f7-icon class="category-separate-icon" f7="chevron_right"></f7-icon>
                            <span>{{ transaction.expenseCategory | secondaryCategoryName(allCategories[$constants.category.allCategoryTypes.Expense]) }}</span>
                        </div>
                        <tree-view-selection-sheet primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                   primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                   primary-sub-items-field="subCategories"
                                                   secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                   secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                   :items="allCategories[$constants.category.allCategoryTypes.Expense]"
                                                   :show.sync="showCategorySheet"
                                                   v-model="transaction.expenseCategory">
                        </tree-view-selection-sheet>
                    </f7-list-item>

                    <f7-list-item
                        class="transaction-edit-category"
                        key="incomeCategorySelection"
                        link="#"
                        :class="{ 'disabled': !hasAvailableIncomeCategories }"
                        :header="$t('Category')"
                        @click="showCategorySheet = true"
                        v-if="transaction.type === $constants.transaction.allTransactionTypes.Income"
                    >
                        <div slot="title">
                            <span>{{ transaction.incomeCategory | primaryCategoryName(allCategories[$constants.category.allCategoryTypes.Income]) }}</span>
                            <f7-icon class="category-separate-icon" f7="chevron_right"></f7-icon>
                            <span>{{ transaction.incomeCategory | secondaryCategoryName(allCategories[$constants.category.allCategoryTypes.Income]) }}</span>
                        </div>
                        <tree-view-selection-sheet primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                   primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                   primary-sub-items-field="subCategories"
                                                   secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                   secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                   :items="allCategories[$constants.category.allCategoryTypes.Income]"
                                                   :show.sync="showCategorySheet"
                                                   v-model="transaction.incomeCategory">
                        </tree-view-selection-sheet>
                    </f7-list-item>

                    <f7-list-item
                        class="transaction-edit-category"
                        key="transferCategorySelection"
                        link="#"
                        :class="{ 'disabled': !hasAvailableTransferCategories }"
                        :header="$t('Category')"
                        @click="showCategorySheet = true"
                        v-if="transaction.type === $constants.transaction.allTransactionTypes.Transfer"
                    >
                        <div slot="title">
                            <span>{{ transaction.transferCategory | primaryCategoryName(allCategories[$constants.category.allCategoryTypes.Transfer]) }}</span>
                            <f7-icon class="category-separate-icon" f7="chevron_right"></f7-icon>
                            <span>{{ transaction.transferCategory | secondaryCategoryName(allCategories[$constants.category.allCategoryTypes.Transfer]) }}</span>
                        </div>
                        <tree-view-selection-sheet primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                   primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                   primary-sub-items-field="subCategories"
                                                   secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                   secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                   :items="allCategories[$constants.category.allCategoryTypes.Transfer]"
                                                   :show.sync="showCategorySheet"
                                                   v-model="transaction.transferCategory">
                        </tree-view-selection-sheet>
                    </f7-list-item>

                    <f7-list-item
                        class="transaction-edit-account"
                        link="#"
                        :class="{ 'disabled': !allVisibleAccounts.length }"
                        :header="$t(sourceAccountName)"
                        :title="transaction.sourceAccountId | accountName(allAccounts)"
                        @click="showSourceAccountSheet = true"
                    >
                        <two-column-list-item-selection-sheet primary-key-field="id" primary-value-field="category"
                                                              primary-title-field="name" primary-footer-field="displayBalance"
                                                              primary-icon-field="icon" primary-icon-type="account"
                                                              primary-sub-items-field="accounts"
                                                              :primary-title-i18n="true"
                                                              secondary-key-field="id" secondary-value-field="id"
                                                              secondary-title-field="name" secondary-footer-field="displayBalance"
                                                              secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                              :items="categorizedAccounts"
                                                              :show.sync="showSourceAccountSheet"
                                                              v-model="transaction.sourceAccountId">
                        </two-column-list-item-selection-sheet>
                    </f7-list-item>

                    <f7-list-item
                        class="transaction-edit-account"
                        link="#"
                        :class="{ 'disabled': !allVisibleAccounts.length }"
                        :header="$t('Destination Account')"
                        :title="transaction.destinationAccountId | accountName(allAccounts)"
                        v-if="transaction.type === $constants.transaction.allTransactionTypes.Transfer"
                        @click="showDestinationAccountSheet = true"
                    >
                        <two-column-list-item-selection-sheet primary-key-field="id" primary-value-field="category"
                                                              primary-title-field="name" primary-footer-field="displayBalance"
                                                              primary-icon-field="icon" primary-icon-type="account"
                                                              primary-sub-items-field="accounts"
                                                              :primary-title-i18n="true"
                                                              secondary-key-field="id" secondary-value-field="id"
                                                              secondary-title-field="name" secondary-footer-field="displayBalance"
                                                              secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                              :items="categorizedAccounts"
                                                              :show.sync="showDestinationAccountSheet"
                                                              v-model="transaction.destinationAccountId">
                        </two-column-list-item-selection-sheet>
                    </f7-list-item>

                    <f7-list-input
                        :label="$t('Transaction Time')"
                        type="datetime-local"
                        class="transaction-edit-time"
                        :value="transaction.time"
                        @input="transaction.time = $event.target.value"
                    >
                    </f7-list-input>

                    <f7-list-item :header="$t('Tags')" link="#"
                                  smart-select :smart-select-params="{ openIn: 'sheet', setValueText: false, closeOnSelect: true, sheetCloseLinkText: $t('Close') }">
                        <select multiple v-model="transaction.tagIds">
                            <option v-for="tag in allTags"
                                    :key="tag.id"
                                    :value="tag.id">{{ tag.name }}</option>
                        </select>
                        <f7-block class="margin-top-half no-padding" slot="footer" v-if="transaction.tagIds && transaction.tagIds.length">
                            <f7-chip class="transaction-edit-tag" media-bg-color="black"
                                     v-for="tagId in transaction.tagIds"
                                     :key="tagId"
                                     :text="tagId | tagName(allTags)">
                                <f7-icon slot="media" f7="number"></f7-icon>
                            </f7-chip>
                        </f7-block>
                        <f7-block class="margin-top-half no-padding" slot="footer" v-else-if="!transaction.tagIds || !transaction.tagIds.length">
                            <f7-chip class="transaction-edit-tag" :text="$t('None')">
                            </f7-chip>
                        </f7-block>
                    </f7-list-item>

                    <f7-list-input
                        type="textarea"
                        style="height: auto"
                        :label="$t('Description')"
                        :placeholder="mode !== 'view' ? $t('Your transaction description (optional)') : ''"
                        :value="transaction.comment"
                        @input="transaction.comment = $event.target.value; changeSize($event)"
                    ></f7-list-input>
                </f7-list>
            </f7-card-content>
        </f7-card>
    </f7-page>
</template>

<script>
export default {
    data() {
        const self = this;
        const query = self.$f7route.query;
        const now = new Date();

        let defaultType = self.$constants.transaction.allTransactionTypes.Expense;

        if (query.type === self.$constants.transaction.allTransactionTypes.Income.toString()) {
            defaultType = self.$constants.transaction.allTransactionTypes.Income;
        } else if (query.type === self.$constants.transaction.allTransactionTypes.Transfer.toString()) {
            defaultType = self.$constants.transaction.allTransactionTypes.Transfer;
        }

        return {
            mode: 'add',
            editTransactionId: null,
            transaction: {
                type: defaultType,
                unixTime: self.$utilities.getUnixTime(now),
                time: self.$utilities.formatDate(now, 'YYYY-MM-DDTHH:mm'),
                expenseCategory: '',
                incomeCategory: '',
                transferCategory: '',
                sourceAccountId: '',
                destinationAccountId: '',
                sourceAmount: 0,
                destinationAmount: 0,
                tagIds: [],
                comment: ''
            },
            loading: true,
            submitting: false,
            showAccountBalance: self.$settings.isShowAccountBalance(),
            showSourceAmountSheet: false,
            showDestinationAmountSheet: false,
            showCategorySheet: false,
            showSourceAccountSheet: false,
            showDestinationAccountSheet: false
        };
    },
    computed: {
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
            if (this.transaction.type === this.$constants.transaction.allTransactionTypes.Expense) {
                return 'Expense Amount';
            } else if (this.transaction.type === this.$constants.transaction.allTransactionTypes.Income) {
                return 'Income Amount';
            } else if (this.transaction.type === this.$constants.transaction.allTransactionTypes.Transfer) {
                return 'Transfer Out Amount';
            } else {
                return '';
            }
        },
        sourceAccountName() {
            if (this.transaction.type === this.$constants.transaction.allTransactionTypes.Expense || this.transaction.type === this.$constants.transaction.allTransactionTypes.Income) {
                return 'Account';
            } else if (this.transaction.type === this.$constants.transaction.allTransactionTypes.Transfer) {
                return 'Source Account';
            } else {
                return '';
            }
        },
        allAccounts() {
            return this.$store.getters.allPlainAccounts;
        },
        allVisibleAccounts() {
            return this.$store.getters.allVisiblePlainAccounts;
        },
        categorizedAccounts() {
            const categorizedAccounts = this.$utilities.copyObjectTo(this.$utilities.getCategorizedAccounts(this.allVisibleAccounts), {});

            for (let category in categorizedAccounts) {
                if (!Object.prototype.hasOwnProperty.call(categorizedAccounts, category)) {
                    continue;
                }

                const accountCategory = categorizedAccounts[category];

                if (accountCategory.accounts) {
                    for (let i = 0; i < accountCategory.accounts.length; i++) {
                        const account = accountCategory.accounts[i];

                        if (this.showAccountBalance && account.isAsset) {
                            account.displayBalance = this.$options.filters.currency(account.balance, account.currency);
                        } else if (this.showAccountBalance && account.isLiability) {
                            account.displayBalance = this.$options.filters.currency(-account.balance, account.currency);
                        } else {
                            account.displayBalance = '***';
                        }
                    }
                }

                if (this.showAccountBalance) {
                    const accountsBalance = this.$utilities.getAllFilteredAccountsBalance(categorizedAccounts, account => account.category === accountCategory.category);
                    let totalBalance = 0;
                    let hasUnCalculatedAmount = false;

                    for (let i = 0; i < accountsBalance.length; i++) {
                        if (accountsBalance[i].currency === this.defaultCurrency) {
                            if (accountsBalance[i].isAsset) {
                                totalBalance += accountsBalance[i].balance;
                            } else if (accountsBalance[i].isLiability) {
                                totalBalance -= accountsBalance[i].balance;
                            }
                        } else {
                            const balance = this.$store.getters.getExchangedAmount(accountsBalance[i].balance, accountsBalance[i].currency, this.defaultCurrency);

                            if (!this.$utilities.isNumber(balance)) {
                                hasUnCalculatedAmount = true;
                                continue;
                            }

                            if (accountsBalance[i].isAsset) {
                                totalBalance += Math.floor(balance);
                            } else if (accountsBalance[i].isLiability) {
                                totalBalance -= Math.floor(balance);
                            }
                        }
                    }

                    if (hasUnCalculatedAmount) {
                        totalBalance = totalBalance + '+';
                    }

                    accountCategory.displayBalance = this.$options.filters.currency(totalBalance, this.defaultCurrency);
                } else {
                    accountCategory.displayBalance = '***';
                }
            }

            return categorizedAccounts;
        },
        allCategories() {
            return this.$store.state.allTransactionCategories;
        },
        allTags() {
            return this.$store.state.allTransactionTags;
        },
        defaultCurrency() {
            return this.$store.getters.currentUserDefaultCurrency || this.$t('default.currency');
        },
        hasAvailableExpenseCategories() {
            if (!this.allCategories || !this.allCategories[this.$constants.category.allCategoryTypes.Expense] || !this.allCategories[this.$constants.category.allCategoryTypes.Expense].length) {
                return false;
            }

            const firstAvailableCategoryId = this.getFirstAvailableCategoryId(this.allCategories[this.$constants.category.allCategoryTypes.Expense]);
            return firstAvailableCategoryId !== '';
        },
        hasAvailableIncomeCategories() {
            if (!this.allCategories || !this.allCategories[this.$constants.category.allCategoryTypes.Income] || !this.allCategories[this.$constants.category.allCategoryTypes.Income].length) {
                return false;
            }

            const firstAvailableCategoryId = this.getFirstAvailableCategoryId(this.allCategories[this.$constants.category.allCategoryTypes.Income]);
            return firstAvailableCategoryId !== '';
        },
        hasAvailableTransferCategories() {
            if (!this.allCategories || !this.allCategories[this.$constants.category.allCategoryTypes.Transfer] || !this.allCategories[this.$constants.category.allCategoryTypes.Transfer].length) {
                return false;
            }

            const firstAvailableCategoryId = this.getFirstAvailableCategoryId(this.allCategories[this.$constants.category.allCategoryTypes.Transfer]);
            return firstAvailableCategoryId !== '';
        },
        sourceAmountFontSize() {
            return this.getFontSizeByAmount(this.transaction.sourceAmount);
        },
        destinationAmountFontSize() {
            return this.getFontSizeByAmount(this.transaction.destinationAmount);
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
            if (this.transaction.type === this.$constants.transaction.allTransactionTypes.Expense || this.transaction.type === this.$constants.transaction.allTransactionTypes.Income) {
                this.transaction.destinationAmount = newValue;
            } else if (this.transaction.type === this.$constants.transaction.allTransactionTypes.Transfer) {
                let sourceAccount, destinationAccount = null;

                for (let i = 0; i < this.allVisibleAccounts.length; i++) {
                    if (this.allVisibleAccounts[i].id === this.transaction.sourceAccountId) {
                        sourceAccount = this.allVisibleAccounts[i];
                    }

                    if (this.allVisibleAccounts[i].id === this.transaction.destinationAccountId) {
                        destinationAccount = this.allVisibleAccounts[i];
                    }

                    if (sourceAccount && destinationAccount) {
                        break;
                    }
                }

                if (sourceAccount && destinationAccount && sourceAccount.currency !== destinationAccount.currency) {
                    const exchangedOldValue = this.$store.getters.getExchangedAmount(oldValue, sourceAccount.currency, destinationAccount.currency);
                    const exchangedNewValue = this.$store.getters.getExchangedAmount(newValue, sourceAccount.currency, destinationAccount.currency);

                    if (this.$utilities.isNumber(exchangedOldValue)) {
                        oldValue = Math.floor(exchangedOldValue);
                    }

                    if (this.$utilities.isNumber(exchangedNewValue)) {
                        newValue = Math.floor(exchangedNewValue);
                    }
                }

                if ((!sourceAccount || !destinationAccount || this.transaction.destinationAmount === oldValue) &&
                    (this.$utilities.stringCurrencyToNumeric(this.$constants.transaction.minAmount) <= newValue &&
                        newValue <= this.$utilities.stringCurrencyToNumeric(this.$constants.transaction.maxAmount))) {
                    this.transaction.destinationAmount = newValue;
                }
            }
        },
        'transaction.destinationAmount': function (newValue) {
            if (this.transaction.type === this.$constants.transaction.allTransactionTypes.Expense || this.transaction.type === this.$constants.transaction.allTransactionTypes.Income) {
                this.transaction.sourceAmount = newValue;
            }
        },
        'transaction.time': function (newValue) {
            if (!newValue) {
                newValue = this.$utilities.formatDate(new Date(), 'YYYY-MM-DDTHH:mm');
                this.transaction.time = newValue;
            }

            this.transaction.unixTime = this.$utilities.getUnixTime(newValue);
        }
    },
    created() {
        const self = this;
        const query = self.$f7route.query;
        const router = self.$f7router;

        if (self.$f7route.path === '/transaction/edit') {
            self.mode = 'edit';
        } else if (self.$f7route.path === '/transaction/detail') {
            self.mode = 'view';
        }

        self.loading = true;

        const promises = [
            self.$store.dispatch('loadAllAccounts', { force: false }),
            self.$store.dispatch('loadAllCategories', { force: false }),
            self.$store.dispatch('loadAllTags', { force: false })
        ];

        if (query.id) {
            if (self.mode === 'edit') {
                self.editTransactionId = query.id;
            }

            promises.push(self.$store.dispatch('getTransaction', { transactionId: query.id }));
        }

        if (query.type && query.type !== '0' &&
            query.type >= self.$constants.transaction.allTransactionTypes.Income &&
            query.type <= self.$constants.transaction.allTransactionTypes.Transfer) {
            self.transaction.type = parseInt(query.type);
        }

        Promise.all(promises).then(function (responses) {
            if (query.id && !responses[3]) {
                self.$toast('Unable to get transaction');
                router.back();
                return;
            }

            if (self.allCategories[self.$constants.category.allCategoryTypes.Expense] &&
                self.allCategories[self.$constants.category.allCategoryTypes.Expense].length) {
                if (query.categoryId && query.categoryId !== '0' && self.isCategoryIdAvailable(self.allCategories[self.$constants.category.allCategoryTypes.Expense], query.categoryId)) {
                    self.transaction.expenseCategory = query.categoryId;
                }

                if (!self.transaction.expenseCategory) {
                    self.transaction.expenseCategory = self.getFirstAvailableCategoryId(self.allCategories[self.$constants.category.allCategoryTypes.Expense]);
                }
            }

            if (self.allCategories[self.$constants.category.allCategoryTypes.Income] &&
                self.allCategories[self.$constants.category.allCategoryTypes.Income].length) {
                if (query.categoryId && query.categoryId !== '0' && self.isCategoryIdAvailable(self.allCategories[self.$constants.category.allCategoryTypes.Income], query.categoryId)) {
                    self.transaction.incomeCategory = query.categoryId;
                }

                if (!self.transaction.incomeCategory) {
                    self.transaction.incomeCategory = self.getFirstAvailableCategoryId(self.allCategories[self.$constants.category.allCategoryTypes.Income]);
                }
            }

            if (self.allCategories[self.$constants.category.allCategoryTypes.Transfer] &&
                self.allCategories[self.$constants.category.allCategoryTypes.Transfer].length) {
                if (query.categoryId && query.categoryId !== '0' && self.isCategoryIdAvailable(self.allCategories[self.$constants.category.allCategoryTypes.Transfer], query.categoryId)) {
                    self.transaction.transferCategory = query.categoryId;
                }

                if (!self.transaction.transferCategory) {
                    self.transaction.transferCategory = self.getFirstAvailableCategoryId(self.allCategories[self.$constants.category.allCategoryTypes.Transfer]);
                }
            }

            if (self.allVisibleAccounts.length) {
                if (query.accountId && query.accountId !== '0') {
                    for (let i = 0; i < self.allVisibleAccounts.length; i++) {
                        if (self.allVisibleAccounts[i].id === query.accountId) {
                            self.transaction.sourceAccountId = query.accountId;
                            self.transaction.destinationAccountId = query.accountId;
                            break;
                        }
                    }
                }

                if (!self.transaction.sourceAccountId) {
                    self.transaction.sourceAccountId = self.allVisibleAccounts[0].id;
                }

                if (!self.transaction.destinationAccountId) {
                    self.transaction.destinationAccountId = self.allVisibleAccounts[0].id;
                }
            }

            if (query.id) {
                const transaction = responses[3];

                if (self.mode === 'edit') {
                    self.transaction.id = transaction.id;
                }

                self.transaction.type = transaction.type;

                if (self.transaction.type === self.$constants.transaction.allTransactionTypes.Expense) {
                    self.transaction.expenseCategory = transaction.categoryId;
                } else if (self.transaction.type === self.$constants.transaction.allTransactionTypes.Income) {
                    self.transaction.incomeCategory = transaction.categoryId;
                } else if (self.transaction.type === self.$constants.transaction.allTransactionTypes.Transfer) {
                    self.transaction.transferCategory = transaction.categoryId;
                }

                if (self.mode === 'edit') {
                    self.transaction.unixTime = transaction.time;
                    self.transaction.time = self.$utilities.formatUnixTime(transaction.time, 'YYYY-MM-DDTHH:mm');
                }

                self.transaction.sourceAccountId = transaction.sourceAccountId;

                if (transaction.destinationAccountId) {
                    self.transaction.destinationAccountId = transaction.destinationAccountId;
                }

                self.transaction.sourceAmount = transaction.sourceAmount;

                if (transaction.destinationAmount) {
                    self.transaction.destinationAmount = transaction.destinationAmount;
                }

                self.transaction.tagIds = transaction.tagIds || [];
                self.transaction.comment = transaction.comment;
            }

            self.loading = false;
        }).catch(error => {
            self.$logger.error('failed to load essential data for editing transaction', error);

            if (!error.processed) {
                self.$toast(error.message || error);
                router.back();
            }
        });
    },
    methods: {
        save() {
            const self = this;
            const router = self.$f7router;

            if (self.mode === 'view') {
                return;
            }

            self.submitting = true;
            self.$showLoading(() => self.submitting);

            const submitTransaction = {
                type: self.transaction.type,
                time: self.transaction.unixTime,
                sourceAccountId: self.transaction.sourceAccountId,
                sourceAmount: self.transaction.sourceAmount,
                destinationAccountId: '0',
                destinationAmount: 0,
                tagIds: self.transaction.tagIds,
                comment: self.transaction.comment
            };

            if (self.transaction.type === self.$constants.transaction.allTransactionTypes.Expense) {
                submitTransaction.categoryId = self.transaction.expenseCategory;
            } else if (self.transaction.type === self.$constants.transaction.allTransactionTypes.Income) {
                submitTransaction.categoryId = self.transaction.incomeCategory;
            } else if (self.transaction.type === self.$constants.transaction.allTransactionTypes.Transfer) {
                submitTransaction.categoryId = self.transaction.transferCategory;
                submitTransaction.destinationAccountId = self.transaction.destinationAccountId;
                submitTransaction.destinationAmount = self.transaction.destinationAmount;
            } else {
                self.$toast('An error has occurred');
                return;
            }

            if (self.mode === 'edit') {
                submitTransaction.id = self.transaction.id;
            }

            self.$store.dispatch('saveTransaction', {
                transaction: submitTransaction
            }).then(() => {
                self.submitting = false;
                self.$hideLoading();

                if (self.mode === 'add') {
                    self.$toast('You have added a new transaction');
                } else if (self.mode === 'edit') {
                    self.$toast('You have saved this transaction');
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
        changeSize(event) {
            const textarea = event.target;

            if (!textarea) {
                return;
            }

            textarea.scrollTop = 0;
            textarea.style.height = textarea.scrollHeight + 'px';
        },
        isCategoryIdAvailable(categories, categoryId) {
            if (!categories || !categories.length) {
                return false;
            }

            for (let i = 0; i < categories.length; i++) {
                for (let j = 0; j < categories[i].subCategories.length; j++) {
                    if (categories[i].subCategories[j].id === categoryId) {
                        return true;
                    }
                }
            }

            return false;
        },
        getFirstAvailableCategoryId(categories) {
            if (!categories || !categories.length) {
                return '';
            }

            for (let i = 0; i < categories.length; i++) {
                for (let j = 0; j < categories[i].subCategories.length; j++) {
                    return categories[i].subCategories[j].id;
                }
            }
        },
        getFontSizeByAmount(amount) {
            if (amount >= 100000000 || amount <= -100000000) {
                return 32;
            } else if (amount >= 1000000 || amount <= -1000000) {
                return 36;
            } else {
                return 40;
            }
        }
    },
    filters: {
        primaryCategoryName(categoryId, allCategories) {
            for (let i = 0; i < allCategories.length; i++) {
                for (let j = 0; j < allCategories[i].subCategories.length; j++) {
                    const subCategory = allCategories[i].subCategories[j];
                    if (subCategory.id === categoryId) {
                        return allCategories[i].name;
                    }
                }
            }

            return '';
        },
        secondaryCategoryName(categoryId, allCategories) {
            for (let i = 0; i < allCategories.length; i++) {
                for (let j = 0; j < allCategories[i].subCategories.length; j++) {
                    const subCategory = allCategories[i].subCategories[j];
                    if (subCategory.id === categoryId) {
                        return subCategory.name;
                    }
                }
            }

            return '';
        },
        accountName(accountId, allAccounts) {
            for (let i = 0; i < allAccounts.length; i++) {
                if (allAccounts[i].id === accountId) {
                    return allAccounts[i].name;
                }
            }

            return '';
        },
        tagName(tagId, allTags) {
            for (let i = 0; i < allTags.length; i++) {
                if (allTags[i].id === tagId) {
                    return allTags[i].name;
                }
            }

            return '';
        }
    }
};
</script>

<style>
.category-separate-icon {
    margin-left: 5px;
    margin-right: 5px;
    font-size: 18px;
    line-height: 16px;
    color: var(--f7-color-gray-tint);
}

.transaction-edit-amount {
    line-height: 53px;
    color: var(--f7-theme-color);
}

.transaction-edit-amount .item-title {
    font-weight: bolder;
}

.transaction-edit-category .item-header, .transaction-edit-account .item-header {
    margin-bottom: 11px;
}

.transaction-edit-time input[type="datetime-local"] {
    max-width: inherit;
}

.transaction-edit-tag {
    margin-right: 4px;
}
</style>
