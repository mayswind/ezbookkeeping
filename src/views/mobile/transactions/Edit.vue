<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t(title)"></f7-nav-title>
            <f7-nav-right>
                <f7-link :class="{ 'disabled': inputIsEmpty || submitting }" :text="$t(saveButtonTitle)" @click="save"></f7-link>
            </f7-nav-right>

            <f7-subnavbar>
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
                    <f7-list-item class="transaction-edit-category" header="Category" link="#"></f7-list-item>
                    <f7-list-item class="transaction-edit-account" header="Account" title="Account Name" link="#"></f7-list-item>
                    <f7-list-input label="Transaction Time" placeholder="YYYY/MM/DD HH:mm"></f7-list-input>
                    <f7-list-item header="Tags" link="#"></f7-list-item>
                    <f7-list-input type="textarea" label="Description" placeholder="Your transaction description (optional)"></f7-list-input>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list form>
                    <f7-list-item
                        class="transaction-edit-amount padding-top-half padding-bottom-half"
                        :class="{ 'color-theme-teal': transaction.type === $constants.transaction.allTransactionTypes.Expense, 'color-theme-red': transaction.type === $constants.transaction.allTransactionTypes.Income }"
                        :style="{ fontSize: sourceAmountFontSize + 'px' }"
                        :header="$t(sourceAmountName)"
                        :title="transaction.sourceAmount | currency"
                        @click="transaction.showSourceAmountSheet = true"
                    >
                        <number-pad-sheet :min-value="$constants.transaction.minAmount"
                                          :max-value="$constants.transaction.maxAmount"
                                          :show.sync="transaction.showSourceAmountSheet"
                                          v-model="transaction.sourceAmount"
                        ></number-pad-sheet>
                    </f7-list-item>

                    <f7-list-item
                        class="transaction-edit-amount padding-top-half padding-bottom-half"
                        :style="{ fontSize: destinationAmountFontSize + 'px' }"
                        :header="$t('Transfer In Amount')"
                        :title="transaction.destinationAmount | currency"
                        @click="transaction.showDestinationAmountSheet = true"
                        v-if="transaction.type === $constants.transaction.allTransactionTypes.Transfer"
                    >
                        <number-pad-sheet :min-value="$constants.transaction.minAmount"
                                          :max-value="$constants.transaction.maxAmount"
                                          :show.sync="transaction.showDestinationAmountSheet"
                                          v-model="transaction.destinationAmount"
                        ></number-pad-sheet>
                    </f7-list-item>

                    <f7-list-item
                        class="transaction-edit-category"
                        link="#"
                        :header="$t('Category')"
                    >
                    </f7-list-item>

                    <f7-list-item
                        class="transaction-edit-account"
                        link="#"
                        :class="{ 'disabled': !plainAllAccounts.length }"
                        :header="$t(sourceAccountName)"
                        :title="transaction.sourceAccountId | accountName(plainAllAccounts)"
                        @click="transaction.showSourceAccountSheet = true"
                    >
                        <list-item-selection-sheet value-type="item"
                                                   key-field="id" value-field="id" title-field="name"
                                                   icon-field="icon" icon-type="account" color-field="color"
                                                   :items="plainAllAccounts"
                                                   :show.sync="transaction.showSourceAccountSheet"
                                                   v-model="transaction.sourceAccountId">
                        </list-item-selection-sheet>
                    </f7-list-item>

                    <f7-list-item
                        class="transaction-edit-account"
                        link="#"
                        :class="{ 'disabled': !plainAllAccounts.length }"
                        :header="$t('Destination Account')"
                        :title="transaction.destinationAccountId | accountName(plainAllAccounts)"
                        v-if="transaction.type === $constants.transaction.allTransactionTypes.Transfer"
                        @click="transaction.showDestinationAccountSheet = true"
                    >
                        <list-item-selection-sheet value-type="item"
                                                   key-field="id" value-field="id" title-field="name"
                                                   icon-field="icon" icon-type="account" color-field="color"
                                                   :items="plainAllAccounts"
                                                   :show.sync="transaction.showDestinationAccountSheet"
                                                   v-model="transaction.destinationAccountId">
                        </list-item-selection-sheet>
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
                        <f7-block class="margin-top-half no-padding" slot="footer" v-if="transaction.tagIds.length">
                            <f7-chip class="transaction-edit-tag" media-bg-color="black"
                                     v-for="tagId in transaction.tagIds"
                                     :key="tagId"
                                     :text="tagId | tagName(allTags)">
                                <f7-icon slot="media" f7="number"></f7-icon>
                            </f7-chip>
                        </f7-block>
                    </f7-list-item>

                    <f7-list-input
                        type="textarea"
                        :label="$t('Description')"
                        :placeholder="$t('Your transaction description (optional)')"
                        :value="transaction.comment"
                        @input="transaction.comment = $event.target.value"
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
        const now = new Date();

        return {
            editTransactionId: null,
            transaction: {
                type: 3,
                unixTime: self.$utilities.getUnixTime(now),
                time: self.$utilities.formatDate(now, 'YYYY-MM-DDTHH:mm'),
                sourceAccountId: '',
                destinationAccountId: '',
                sourceAmount: 0,
                destinationAmount: 0,
                tagIds: [],
                comment: '',
                showSourceAmountSheet: false,
                showDestinationAmountSheet: false,
                showSourceAccountSheet: false,
                showDestinationAccountSheet: false
            },
            allAccounts: [],
            allCategories: {},
            allTags: [],
            loading: true,
            submitting: false
        };
    },
    computed: {
        title() {
            if (!this.editTransactionId) {
                return 'Add Transaction';
            } else {
                return 'Edit Transaction';
            }
        },
        saveButtonTitle() {
            if (!this.editTransactionId) {
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
        sourceAmountFontSize() {
            return this.getFontSizeByAmount(this.transaction.sourceAmount);
        },
        destinationAmountFontSize() {
            return this.getFontSizeByAmount(this.transaction.destinationAmount);
        },
        plainAllAccounts() {
            const ret = [];

            for (let i = 0; i < this.allAccounts.length; i++) {
                const account = this.allAccounts[i];

                if (account.type === this.$constants.account.allAccountTypes.SingleAccount) {
                    ret.push(account);
                    continue;
                }

                for (let j = 0; j < account.subAccounts.length; j++) {
                    const subAccount = account.subAccounts[j];
                    ret.push(subAccount);
                }
            }

            return ret;
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

                for (let i = 0; i < this.plainAllAccounts.length; i++) {
                    if (this.plainAllAccounts[i].id === this.transaction.sourceAccountId) {
                        sourceAccount = this.plainAllAccounts[i];
                    }

                    if (this.plainAllAccounts[i].id === this.transaction.destinationAccountId) {
                        destinationAccount = this.plainAllAccounts[i];
                    }

                    if (sourceAccount && destinationAccount) {
                        break;
                    }
                }

                if (sourceAccount && destinationAccount && sourceAccount.currency !== destinationAccount.currency) {
                    const exchangedOldValue = this.$exchangeRates.getOtherCurrencyAmount(oldValue, sourceAccount.currency, destinationAccount.currency);
                    const exchangedNewValue = this.$exchangeRates.getOtherCurrencyAmount(newValue, sourceAccount.currency, destinationAccount.currency);

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
            this.transaction.unixTime = this.$utilities.getUnixTime(newValue);
        }
    },
    created() {
        const self = this;
        const query = self.$f7route.query;
        const router = self.$f7router;

        self.loading = true;

        const promises = [
            self.$services.getAllAccounts({ visibleOnly: true }),
            self.$services.getAllTransactionCategories({}),
            self.$services.getAllTransactionTags()
        ];

        if (query.id) {
            self.editTransactionId = query.id;
            promises.push(self.$services.getTransaction({ id: self.editTransactionId }));
        }

        Promise.all(promises).then(function (responses) {
            const accountDta = responses[0].data;
            const categoryData = responses[1].data;
            const tagData = responses[2].data;

            if (!accountDta || !accountDta.success || !accountDta.result) {
                self.$toast('Unable to get account list');
                router.back();
                return;
            }

            if (!categoryData || !categoryData.success || !categoryData.result) {
                self.$toast('Unable to get category list');
                router.back();
                return;
            }

            if (!tagData || !tagData.success || !tagData.result) {
                self.$toast('Unable to get tag list');
                router.back();
                return;
            }

            if (self.editTransactionId && (!responses[3] || !responses[3].data || !responses[3].data.success || !responses[3].data.result)) {
                self.$toast('Unable to get transaction');
                router.back();
                return;
            }

            self.allAccounts = accountDta.result;
            self.allCategories = categoryData.result;
            self.allTags = tagData.result;

            if (self.editTransactionId) {
                const transaction = responses[3].data.result;

                self.transaction.type = transaction.type;
                self.transaction.unixTime = transaction.time;
                self.transaction.time = self.$utilities.formatDate(transaction.time, 'YYYY-MM-DDTHH:mm');
                self.transaction.sourceAccountId = transaction.sourceAccountId;
                self.transaction.destinationAccountId = transaction.destinationAccountId;
                self.transaction.sourceAmount = transaction.sourceAmount;
                self.transaction.destinationAmount = transaction.destinationAmount;
                self.transaction.tagIds = transaction.tagIds;
                self.transaction.comment = transaction.comment;
            } else if (!self.editTransactionId && self.plainAllAccounts.length) {
                self.transaction.sourceAccountId = self.plainAllAccounts[0].id;
                self.transaction.destinationAccountId = self.plainAllAccounts[0].id;
            }

            self.loading = false;
        }).catch(errors => {
            self.$logger.error('failed to load essential data for editing transaction', errors);

            for (let i = 0; i < errors.length; i++) {
                const error = errors[i];

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                    router.back();
                    return;
                } else if (!error.processed) {
                    self.$toast('An error has occurred');
                    router.back();
                    return;
                }
            }
        });
    },
    methods: {
        save() {

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
.transaction-edit-amount {
    line-height: 53px;
    color: var(--f7-theme-color);
}

.transaction-edit-amount .item-title {
    font-weight: bolder;
}

.transaction-edit-account .item-header {
    margin-bottom: 11px;
}

.transaction-edit-time input[type="datetime-local"] {
    max-width: inherit;
}

.transaction-edit-tag {
    margin-right: 4px;
}
</style>
