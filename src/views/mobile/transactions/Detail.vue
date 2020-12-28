<template>
    <f7-page ptr infinite :infinite-preloader="loadingMore" @ptr:refresh="reload" @infinite="loadMore">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Transaction Details')"></f7-nav-title>
            <f7-nav-right v-if="false">
                <f7-link icon-f7="ellipsis" @click="showMoreActionSheet = true"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-header>
                <div class="full-line">
                    <small :style="{ opacity: 0.6 }">YYYY-MM</small>
                    <f7-icon class="transaction-month-card-chevron-icon float-right" f7="chevron_up"></f7-icon>
                </div>
            </f7-card-header>
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list media-list>
                    <f7-list-item class="transaction-info">
                        <div slot="media" class="display-flex no-padding-horizontal">
                            <div class="display-flex flex-direction-column transaction-date">
                                <span class="transaction-day full-line flex-direction-column">DD</span>
                                <span class="transaction-day-of-week full-line flex-direction-column">Sun</span>
                            </div>
                            <div class="transaction-icon display-flex align-items-center">
                                <f7-icon slot="media" f7="app_fill"></f7-icon>
                            </div>
                        </div>
                        <div slot="title" class="no-padding">
                            <span>Category</span>
                        </div>
                        <div slot="footer" class="no-padding-horizontal transaction-footer">
                            <span>HH:mm</span>
                            <span>·</span>
                            <span>Source Account</span>
                        </div>
                        <div slot="content-end" class="padding-right transaction-amount">
                            <span>0.00 USD</span>
                        </div>
                    </f7-list-item>
                    <f7-list-item class="transaction-info">
                        <div slot="media" class="display-flex no-padding-horizontal">
                            <div class="display-flex flex-direction-column transaction-date">
                                <span class="transaction-day full-line flex-direction-column">DD</span>
                                <span class="transaction-day-of-week full-line flex-direction-column">Sun</span>
                            </div>
                            <div class="transaction-icon display-flex align-items-center">
                                <f7-icon slot="media" f7="app_fill"></f7-icon>
                            </div>
                        </div>
                        <div slot="title" class="no-padding">
                            <span>Category 2</span>
                        </div>
                        <div slot="footer" class="no-padding-horizontal transaction-footer">
                            <span>HH:mm</span>
                            <span>·</span>
                            <span>Source Account</span>
                        </div>
                        <div slot="content-end" class="padding-right transaction-amount">
                            <span>0.00 USD</span>
                        </div>
                    </f7-list-item>
                    <f7-list-item class="transaction-info">
                        <div slot="media" class="display-flex no-padding-horizontal">
                            <div class="display-flex flex-direction-column transaction-date">
                                <span class="transaction-day full-line flex-direction-column">DD</span>
                                <span class="transaction-day-of-week full-line flex-direction-column">Sun</span>
                            </div>
                            <div class="transaction-icon display-flex align-items-center">
                                <f7-icon slot="media" f7="app_fill"></f7-icon>
                            </div>
                        </div>
                        <div slot="title" class="no-padding">
                            <span>Category 3</span>
                        </div>
                        <div slot="footer" class="no-padding-horizontal transaction-footer">
                            <span>HH:mm</span>
                            <span>·</span>
                            <span>Source Account</span>
                        </div>
                        <div slot="content-end" class="padding-right transaction-amount">
                            <span>0.00 USD</span>
                        </div>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-if="!loading && noTransaction">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item :title="$t('No transaction data')"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-for="transactionMonthList in transactions" :key="transactionMonthList.yearMonth">
            <f7-accordion-item :opened="transactionMonthList.opened"
                               @accordion:open="transactionMonthList.opened = true"
                               @accordion:close="transactionMonthList.opened = false">
                <f7-card-header>
                    <f7-accordion-toggle class="full-line">
                        <small :style="{ opacity: 0.6 }">
                            <span>{{ transactionMonthList.yearMonth | moment($t('format.date.yearMonth')) }}</span>
                        </small>
                        <small class="transaction-amount-statistics" v-if="transactionMonthList.totalAmount">
                            <span class="text-color-red">
                                {{ transactionMonthList.totalAmount.income | currency(defaultCurrency) | income(transactionMonthList.totalAmount.incompleteIncome) }}
                            </span>
                            <span class="text-color-teal">
                                {{ transactionMonthList.totalAmount.expense | currency(defaultCurrency) | expense(transactionMonthList.totalAmount.incompleteExpense) }}
                            </span>
                        </small>
                        <f7-icon class="transaction-month-card-chevron-icon float-right" :f7="transactionMonthList.opened ? 'chevron_up' : 'chevron_down'"></f7-icon>
                    </f7-accordion-toggle>
                </f7-card-header>
                <f7-card-content class="no-safe-areas" :padding="false" accordion-list>
                    <f7-accordion-content style="height: auto">
                        <f7-list media-list>
                            <f7-list-item class="transaction-info"
                                          v-for="(transaction, idx) in transactionMonthList.items"
                                          :key="transaction.id" :id="transaction | transactionDomId"
                                          swipeout
                            >
                                <div slot="media" class="display-flex no-padding-horizontal">
                                    <div class="display-flex flex-direction-column transaction-date" :style="transaction | transactionDateStyle(idx > 0 ? transactionMonthList.items[idx - 1] : null)">
                                        <span class="transaction-day full-line flex-direction-column">
                                            {{ transaction.day }}
                                        </span>
                                        <span class="transaction-day-of-week full-line flex-direction-column">
                                            {{ transaction.dayOfWeek }}
                                        </span>
                                    </div>
                                    <div class="transaction-icon display-flex align-items-center">
                                        <f7-icon v-if="transaction.category && transaction.category.color"
                                                 :icon="transaction.category.icon | categoryIcon"
                                                 :style="transaction.category.color | categoryIconStyle('var(--category-icon-color)')">
                                        </f7-icon>
                                    </div>
                                </div>
                                <div slot="title" class="no-padding">
                                    <span v-if="transaction.type === $constants.transaction.allTransactionTypes.ModifyBalance">
                                        {{ $t('Modify Balance') }}
                                    </span>
                                    <span v-else-if="transaction.type !== $constants.transaction.allTransactionTypes.ModifyBalance && transaction.category">
                                        {{ transaction.category.name }}
                                    </span>
                                    <span v-else-if="transaction.type !== $constants.transaction.allTransactionTypes.ModifyBalance && !transaction.category">
                                        {{ transaction.type | transactionTypeName($constants.transaction.allTransactionTypes) }}
                                    </span>
                                </div>
                                <div slot="footer" class="no-padding-horizontal transaction-footer">
                                    <span>{{ transaction.time | moment($t('format.time.hourMinute')) }}</span>
                                    <span v-if="transaction.sourceAccount">·</span>
                                    <span v-if="transaction.sourceAccount">{{ transaction.sourceAccount.name }}</span>
                                </div>
                                <div slot="content-end" class="padding-right transaction-amount" v-if="transaction.sourceAccount">
                                    <span :class="{ 'text-color-teal': transaction.type === $constants.transaction.allTransactionTypes.Expense, 'text-color-red': transaction.type === $constants.transaction.allTransactionTypes.Income }">
                                        {{ transaction.sourceAmount | currency(transaction.sourceAccount.currency) }}
                                    </span>
                                </div>
                                <f7-swipeout-actions right>
                                    <f7-swipeout-button color="orange" close :text="$t('Edit')" @click="edit(transaction)"></f7-swipeout-button>
                                    <f7-swipeout-button color="red" class="padding-left padding-right" @click="remove(transaction, false)">
                                        <f7-icon f7="trash"></f7-icon>
                                    </f7-swipeout-button>
                                </f7-swipeout-actions>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-content>
                </f7-card-content>
            </f7-accordion-item>
        </f7-card>

        <f7-block class="text-align-center" v-if="!loading && maxTime > 0">
            <f7-link :class="{ 'disabled': loadingMore }" href="#" @click="loadMore">{{ $t('Load More') }}</f7-link>
        </f7-block>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-actions close-by-outside-click close-on-escape :opened="showDeleteActionSheet" @actions:closed="showDeleteActionSheet = false">
            <f7-actions-group>
                <f7-actions-label>{{ $t('Are you sure you want to delete this transaction?') }}</f7-actions-label>
                <f7-actions-button color="red" @click="remove(transactionToDelete, true)">{{ $t('Delete') }}</f7-actions-button>
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
        return {
            transactions: [],
            allAccounts: {},
            allCategories: {},
            allTags: {},
            maxTime: 0,
            loading: true,
            loadingMore: false,
            transactionToDelete: null,
            showMoreActionSheet: false,
            showDeleteActionSheet: false
        };
    },
    computed: {
        defaultCurrency() {
            return this.$user.getUserInfo() ? this.$user.getUserInfo().defaultCurrency : this.$t('default.currency');
        },
        noTransaction() {
            for (let i = 0; i < this.transactions.length; i++) {
                const transactionMonthList = this.transactions[i];

                for (let j = 0; j < transactionMonthList.items.length; j++) {
                    if (transactionMonthList.items[j]) {
                        return false;
                    }
                }
            }

            return true;
        }
    },
    created() {
        this.reload(null);
    },
    methods: {
        reload(done) {
            const self = this;
            const router = self.$f7router;

            self.loading = true;
            self.maxTime = 0;

            const promises = [
                self.$services.getAllAccounts({ visibleOnly: false }),
                self.$services.getAllTransactionCategories({}),
                self.$services.getAllTransactionTags(),
                self.$services.getTransactions({
                    maxTime: self.maxTime
                })
            ];

            Promise.all(promises).then(responses => {
                if (done) {
                    done();
                }

                const accountData = responses[0].data;
                const categoryData = responses[1].data;
                const tagData = responses[2].data;
                const transactionListData = responses[3].data;

                if (!accountData || !accountData.success || !accountData.result) {
                    self.$toast('Unable to get account list');
                    if (!done) {
                        router.back();
                    }
                    return;
                }

                if (!categoryData || !categoryData.success || !categoryData.result) {
                    self.$toast('Unable to get category list');
                    if (!done) {
                        router.back();
                    }
                    return;
                }

                if (!tagData || !tagData.success || !tagData.result) {
                    self.$toast('Unable to get tag list');
                    if (!done) {
                        router.back();
                    }
                    return;
                }

                if (!transactionListData || !transactionListData.success || !transactionListData.result) {
                    self.$toast('Unable to get transaction list');
                    if (!done) {
                        router.back();
                    }
                    return;
                }

                const allAccounts = self.$utilities.getPlainAccounts(accountData.result);
                self.allAccounts = {};

                for (let i = 0; i < allAccounts.length; i++) {
                    const account = allAccounts[i];
                    self.allAccounts[account.id] = account;
                }

                const allCategories = categoryData.result;
                self.allCategories = {};

                for (let categoryType in allCategories) {
                    if (!Object.prototype.hasOwnProperty.call(allCategories, categoryType)) {
                        continue;
                    }

                    const categoryList = allCategories[categoryType];

                    for (let i = 0; i < categoryList.length; i++) {
                        const category = categoryList[i];
                        self.allCategories[category.id] = category;

                        for (let j = 0; j < category.subCategories.length; j++) {
                            const subCategory = category.subCategories[j];
                            self.allCategories[subCategory.id] = subCategory;
                        }
                    }
                }

                const allTags = tagData.result;
                self.allTags = {};

                for (let i = 0; i < allTags.length; i++) {
                    const tag = allTags[i];
                    self.allTags[tag.id] = tag;
                }

                self.transactions = [];
                self.setResult(transactionListData.result);

                self.loading = false;
            }).catch(error => {
                self.$logger.error('failed to load transaction list', error);

                if (done) {
                    done();
                }

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                    if (!done) {
                        router.back();
                    }
                } else if (!error.processed) {
                    self.$toast('Unable to get transaction list');
                    if (!done) {
                        router.back();
                    }
                }
            });
        },
        loadMore() {
            const self = this;

            if (self.maxTime <= 0) {
                return;
            }

            if (self.loadingMore) {
                return;
            }

            self.loadingMore = true;

            self.$services.getTransactions({
                maxTime: self.maxTime
            }).then(response => {
                self.loadingMore = false;

                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$toast('Unable to get transaction list');
                    return;
                }

                self.setResult(data.result);
            }).catch(error => {
                self.loadingMore = false;

                self.$logger.error('failed to reload transaction list', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                } else if (!error.processed) {
                    self.$toast('Unable to get account list');
                }
            });
        },
        edit(transaction) {
            this.$f7router.navigate('/transaction/edit?id=' + transaction.id);
        },
        remove(transaction, confirm) {
            const self = this;
            const app = self.$f7;
            const $$ = app.$;

            if (!transaction) {
                self.$alert('An error has occurred');
                return;
            }

            if (!confirm) {
                self.transactionToDelete = transaction;
                self.showDeleteActionSheet = true;
                return;
            }

            self.showDeleteActionSheet = false;
            self.transactionToDelete = null;
            self.$showLoading();

            self.$services.deleteTransaction({
                id: transaction.id
            }).then(response => {
                self.$hideLoading();
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$toast('Unable to delete this transaction');
                    return;
                }

                app.swipeout.delete($$(`#${self.$options.filters.transactionDomId(transaction)}`), () => {
                    for (let i = 0; i < self.transactions.length; i++) {
                        const transactionMonthList = self.transactions[i];

                        if (!transactionMonthList.items ||
                            transactionMonthList.items[0].time < transaction.time ||
                            transactionMonthList.items[transactionMonthList.items.length - 1].time > transaction.time) {
                            continue;
                        }

                        for (let j = 0; j < transactionMonthList.items.length; j++) {
                            if (transactionMonthList.items[j].id === transaction.id) {
                                transactionMonthList.items.splice(j, 1);
                            }
                        }

                        if (transactionMonthList.items.length < 1) {
                            self.transactions.splice(i, 1);
                        }
                    }
                });
            }).catch(error => {
                self.$logger.error('failed to delete transaction', error);

                self.$hideLoading();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                } else if (!error.processed) {
                    self.$toast('Unable to delete this transaction');
                }
            });
        },
        setResult(result) {
            if (result.items && result.items.length) {
                let currentMonthListIndex = -1;
                let currentMonthList = null;

                for (let i = 0; i < result.items.length; i++) {
                    const transaction = result.items[i];
                    const transactionTime = this.$utilities.parseDateFromUnixtime(transaction.time);
                    transaction.day = this.$utilities.getDay(transactionTime);
                    transaction.dayOfWeek = this.$t(`datetime.${this.$utilities.getDayOfWeek(transactionTime)}.short`);
                    transaction.sourceAccount = this.allAccounts[transaction.sourceAccountId];
                    transaction.destinationAccount = this.allAccounts[transaction.destinationAccountId];
                    transaction.category = this.allCategories[transaction.categoryId];
                    transaction.tags = [];

                    if (transaction.tagIds && transaction.tagIds.length) {
                        for (let j = 0; j < transaction.tagIds.length; j++) {
                            const tag = this.allTags[transaction.tagIds[j]];

                            if (tag) {
                                transaction.tags.push(tag);
                            }
                        }
                    }

                    const transactionYear = this.$utilities.getYear(transactionTime);
                    const transactionMonth = this.$utilities.getMonth(transactionTime);

                    if (currentMonthList && currentMonthList.year === transactionYear && currentMonthList.month === transactionMonth) {
                        currentMonthList.items.push(transaction);
                        this.calculateMonthTotalAmount(currentMonthList, true);
                        continue;
                    }

                    for (let j = currentMonthListIndex + 1; j < this.transactions.length; j++) {
                        if (this.transactions[j].year === transactionYear && this.transactions[j].month === transactionMonth) {
                            currentMonthListIndex = j;
                            currentMonthList = this.transactions[j];

                            if (j > 0) {
                                this.calculateMonthTotalAmount(this.transactions[j - 1], false);
                            }

                            break;
                        }
                    }

                    if (!currentMonthList && this.transactions.length > 0) {
                        this.calculateMonthTotalAmount(this.transactions[this.transactions.length - 1], false);
                    }

                    if (!currentMonthList || currentMonthList.year !== transactionYear || currentMonthList.month !== transactionMonth) {
                        this.calculateMonthTotalAmount(currentMonthList, false);

                        this.transactions.push({
                            year: transactionYear,
                            month: transactionMonth,
                            yearMonth: `${transactionYear}-${transactionMonth}`,
                            opened: true,
                            items: []
                        });

                        currentMonthListIndex = this.transactions.length - 1;
                        currentMonthList = this.transactions[this.transactions.length - 1];
                    }

                    currentMonthList.items.push(transaction);
                    this.calculateMonthTotalAmount(currentMonthList, true);
                }
            }

            if (result.nextTimeSequenceId) {
                this.maxTime = result.nextTimeSequenceId;
            } else {
                this.calculateMonthTotalAmount(this.transactions[this.transactions.length - 1], false);
                this.maxTime = -1;
            }
        },
        calculateMonthTotalAmount(transactionMonthList, incomplete) {
            if (!transactionMonthList) {
                return;
            }

            let totalExpense = 0;
            let totalIncome = 0;
            let hasUnCalculatedTotalExpense = false;
            let hasUnCalculatedTotalIncome = false;

            for (let i = 0; i < transactionMonthList.items.length; i++) {
                const transaction = transactionMonthList.items[i];

                if (!transaction.destinationAccount) {
                    continue;
                }

                let amount = transaction.destinationAmount;

                if (transaction.destinationAccount.currency !== this.defaultCurrency) {
                    const balance = this.$exchangeRates.getOtherCurrencyAmount(amount, transaction.destinationAccount.currency, this.defaultCurrency);

                    if (!this.$utilities.isNumber(balance)) {
                        if (transaction.type === this.$constants.transaction.allTransactionTypes.Expense) {
                            hasUnCalculatedTotalExpense = true;
                        } else if (transaction.type === this.$constants.transaction.allTransactionTypes.Income) {
                            hasUnCalculatedTotalIncome = true;
                        }

                        continue;
                    }

                    amount = Math.floor(balance);
                }

                if (transaction.type === this.$constants.transaction.allTransactionTypes.Expense) {
                    totalExpense += amount;
                } else if (transaction.type === this.$constants.transaction.allTransactionTypes.Income) {
                    totalIncome += amount;
                }
            }

            transactionMonthList.totalAmount = {
                expense: totalExpense,
                incompleteExpense: incomplete || hasUnCalculatedTotalExpense,
                income: totalIncome,
                incompleteIncome: incomplete || hasUnCalculatedTotalIncome
            };
        }
    },
    filters: {
        transactionTypeName(type, allTransactionTypes) {
            if (type === allTransactionTypes.Income) {
                return 'Income';
            } else if (type === allTransactionTypes.Expense) {
                return 'Expense';
            } else if (type === allTransactionTypes.Income) {
                return 'Transfer';
            } else {
                return 'Transaction';
            }
        },
        transactionDomId(transaction) {
            return 'transaction_' + transaction.id;
        },
        transactionDateStyle(transaction, previousTransaction) {
            if (!previousTransaction || transaction.day !== previousTransaction.day) {
                return {};
            }

            return {
                color: 'transparent'
            }
        },
        income(value, incomplete) {
            return '+' + value + (incomplete ? '+' : '');
        },
        expense(value, incomplete) {
            return '-' + value + (incomplete ? '+' : '');
        }
    }
};
</script>

<style>
.transaction-month-card-chevron-icon {
    color: var(--f7-list-chevron-icon-color);
    font-size: var(--f7-list-chevron-icon-font-size);
    font-weight: bolder;
}

.transaction-amount-statistics > span {
    margin-right: 4px;
}

.transaction-info .item-media {
    height: 64px;
}

.transaction-info .item-media + .item-inner {
    margin-left: 10px;
}

.transaction-date {
    margin-right: 6px;
}

.transaction-day {
    opacity: 0.6;
    font-size: 16px;
    font-weight: bold;
    text-align: left;
}

.transaction-day-of-week {
    opacity: 0.6;
    font-size: 12px;
    text-align: center;
}

.transaction-footer {
    padding-top: 4px;
}

.transaction-footer > span {
    margin-right: 4px;
}

.transaction-amount {
    color: var(--f7-list-item-after-text-color);
    white-space: nowrap;
}

.transaction-info .item-inner:after {
    background-color: transparent;
}

.transaction-info .transaction-icon:after {
    content: '';
    position: absolute;
    background-color: var(--f7-list-item-border-color);
    display: block;
    z-index: 15;
    top: auto;
    right: auto;
    bottom: 0;
    height: 1px;
    width: 100%;
    transform-origin: 50% 100%;
    transform: scaleY(calc(1 / var(--f7-device-pixel-ratio)));
}
</style>
