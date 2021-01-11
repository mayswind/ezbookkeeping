<template>
    <f7-page ptr
             infinite
             :infinite-preloader="loadingMore"
             :infinite-distance="400"
             @ptr:refresh="reload"
             @page:afterin="onPageAfterIn"
             @infinite="loadMore(true)">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Transaction List')"></f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link icon-f7="search" class="searchbar-enable" data-searchbar=".searchbar-keyword"></f7-link>
                <f7-link icon-f7="plus" :href="`/transaction/add?type=${query.type}&categoryId=${query.categoryId}&accountId=${query.accountId}`"></f7-link>
            </f7-nav-right>

            <f7-searchbar
                expandable custom-search
                class="searchbar-keyword"
                :value="query.keyword"
                :placeholder="$t('Search')"
                :disable-button-text="$t('Cancel')"
                @change="changeKeywordFilter($event.target.value)"
            ></f7-searchbar>
        </f7-navbar>

        <f7-toolbar tabbar bottom>
            <f7-link class="tabbar-text-with-ellipsis" popover-open=".date-popover-menu">
                <span :class="{ 'tabbar-item-changed': query.maxTime > 0 || query.minTime > 0 }">{{ query.dateType | dateName('Date') | t }}</span>
            </f7-link>
            <f7-link class="tabbar-text-with-ellipsis" popover-open=".type-popover-menu">
                <span :class="{ 'tabbar-item-changed': query.type > 0 }">{{ query.type | typeName('Type') | t }}</span>
            </f7-link>
            <f7-link class="tabbar-text-with-ellipsis" popover-open=".category-popover-menu" :class="{ 'disabled': query.type === 1 }">
                <span :class="{ 'tabbar-item-changed': query.categoryId > 0 }">{{ query.categoryId | categoryName(allCategories, $t('Category')) }}</span>
            </f7-link>
            <f7-link class="tabbar-text-with-ellipsis" popover-open=".account-popover-menu">
                <span :class="{ 'tabbar-item-changed': query.accountId > 0 }">{{ query.accountId | accountName(allAccounts, $t('Account')) }}</span>
            </f7-link>
        </f7-toolbar>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-header>
                <div class="full-line">
                    <small :style="{ opacity: 0.6 }">YYYY-MM</small>
                    <small class="transaction-amount-statistics">
                        <span>0.00 USD</span>
                        <span>0.00 USD</span>
                    </small>
                    <f7-icon class="transaction-month-card-chevron-icon float-right" f7="chevron_up"></f7-icon>
                </div>
            </f7-card-header>
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list media-list>
                    <f7-list-item class="transaction-info" link="#" chevron-center>
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
                        <div slot="after" class="no-padding transaction-amount">
                            <span>0.00 USD</span>
                        </div>
                    </f7-list-item>
                    <f7-list-item class="transaction-info" link="#" chevron-center>
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
                        <div slot="after" class="no-padding transaction-amount">
                            <span>0.00 USD</span>
                        </div>
                    </f7-list-item>
                    <f7-list-item class="transaction-info" link="#" chevron-center>
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
                        <div slot="after" class="no-padding transaction-amount">
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
                               @accordion:open="collapseTransactionMonthList(transactionMonthList, false)"
                               @accordion:close="collapseTransactionMonthList(transactionMonthList, true)">
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
                    <f7-accordion-content :style="{ height: transactionMonthList.opened ? 'auto' : '' }">
                        <f7-list media-list>
                            <f7-list-item class="transaction-info" chevron-center
                                          v-for="(transaction, idx) in transactionMonthList.items"
                                          :key="transaction.id" :id="transaction | transactionDomId"
                                          :link="transaction.type !== $constants.transaction.allTransactionTypes.ModifyBalance ? `/transaction/detail?id=${transaction.id}&type=${transaction.type}` : null"
                                          swipeout
                            >
                                <div slot="media" class="display-flex no-padding-horizontal">
                                    <div class="display-flex flex-direction-column transaction-date" :style="transaction | transactionDateStyle(idx > 0 ? transactionMonthList.items[idx - 1] : null)">
                                        <span class="transaction-day full-line flex-direction-column">
                                            {{ transaction.day }}
                                        </span>
                                        <span class="transaction-day-of-week full-line flex-direction-column">
                                            {{ `datetime.${transaction.dayOfWeek}.short` | t }}
                                        </span>
                                    </div>
                                    <div class="transaction-icon display-flex align-items-center">
                                        <f7-icon v-if="transaction.category && transaction.category.color"
                                                 :icon="transaction.category.icon | categoryIcon"
                                                 :style="transaction.category.color | categoryIconStyle('var(--category-icon-color)')">
                                        </f7-icon>
                                        <f7-icon v-else-if="!transaction.category || !transaction.category.color"
                                                 f7="pencil_ellipsis_rectangle">
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
                                        {{ transaction.type | transactionTypeName($constants.transaction.allTransactionTypes) | t }}
                                    </span>
                                </div>
                                <div slot="text" class="transaction-comment" v-if="transaction.comment">
                                    <span>{{ transaction.comment }}</span>
                                </div>
                                <div slot="footer" class="transaction-footer">
                                    <span>{{ transaction.time | moment($t('format.time.hourMinute')) }}</span>
                                    <span v-if="transaction.sourceAccount">·</span>
                                    <span v-if="transaction.sourceAccount">{{ transaction.sourceAccount.name }}</span>
                                    <span v-if="transaction.sourceAccount && transaction.type === $constants.transaction.allTransactionTypes.Transfer && transaction.destinationAccount && transaction.sourceAccount.id !== transaction.destinationAccount.id">→</span>
                                    <span v-if="transaction.sourceAccount && transaction.type === $constants.transaction.allTransactionTypes.Transfer && transaction.destinationAccount && transaction.sourceAccount.id !== transaction.destinationAccount.id">{{ transaction.destinationAccount.name }}</span>
                                </div>
                                <div slot="after" class="transaction-amount" v-if="transaction.sourceAccount"
                                     :class="{ 'text-color-teal': transaction.type === $constants.transaction.allTransactionTypes.Expense, 'text-color-red': transaction.type === $constants.transaction.allTransactionTypes.Income }">
                                    <span v-if="!query.accountId || query.accountId === '0' || (transaction.sourceAccount && transaction.sourceAccount.id === query.accountId)">{{ transaction.sourceAmount | currency(transaction.sourceAccount.currency) }}</span>
                                    <span v-else-if="query.accountId && query.accountId !== '0' && transaction.destinationAccount && transaction.destinationAccount.id === query.accountId">{{ transaction.destinationAmount | currency(transaction.destinationAccount.currency) }}</span>
                                </div>
                                <f7-swipeout-actions right>
                                    <f7-swipeout-button color="primary" close
                                                        :text="$t('Duplicate')"
                                                        v-if="transaction.type !== $constants.transaction.allTransactionTypes.ModifyBalance"
                                                        @click="duplicate(transaction)"></f7-swipeout-button>
                                    <f7-swipeout-button color="orange" close
                                                        :text="$t('Edit')"
                                                        v-if="transaction.type !== $constants.transaction.allTransactionTypes.ModifyBalance"
                                                        @click="edit(transaction)"></f7-swipeout-button>
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

        <f7-block class="text-align-center" v-if="!loading && hasMoreTransaction">
            <f7-link :class="{ 'disabled': loadingMore }" href="#" @click="loadMore(false)">{{ $t('Load More') }}</f7-link>
        </f7-block>

        <f7-popover class="date-popover-menu" :opened="showDatePopover"
                    @popover:opened="showDatePopover = true" @popover:closed="showDatePopover = false">
            <f7-list>
                <f7-list-item :title="$t('All')" @click="changeDateFilter(0)">
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.dateType === 0"></f7-icon>
                </f7-list-item>
                <f7-list-item :title="$t('Today')" @click="changeDateFilter(1)">
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.dateType === 1"></f7-icon>
                </f7-list-item>
                <f7-list-item :title="$t('Yesterday')" @click="changeDateFilter(2)">
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.dateType === 2"></f7-icon>
                </f7-list-item>
                <f7-list-item :title="$t('Recent 7 days')" @click="changeDateFilter(3)">
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.dateType === 3"></f7-icon>
                </f7-list-item>
                <f7-list-item :title="$t('Recent 30 days')" @click="changeDateFilter(4)">
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.dateType === 4"></f7-icon>
                </f7-list-item>
                <f7-list-item :title="$t('This week')" @click="changeDateFilter(5)">
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.dateType === 5"></f7-icon>
                </f7-list-item>
                <f7-list-item :title="$t('Last week')" @click="changeDateFilter(6)">
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.dateType === 6"></f7-icon>
                </f7-list-item>
                <f7-list-item :title="$t('This month')" @click="changeDateFilter(7)">
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.dateType === 7"></f7-icon>
                </f7-list-item>
                <f7-list-item :title="$t('Last month')" @click="changeDateFilter(8)">
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.dateType === 8"></f7-icon>
                </f7-list-item>
                <f7-list-item :title="$t('This year')" @click="changeDateFilter(9)">
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.dateType === 9"></f7-icon>
                </f7-list-item>
                <f7-list-item :title="$t('Last year')" @click="changeDateFilter(10)">
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.dateType === 10"></f7-icon>
                </f7-list-item>
                <f7-list-item :title="$t('Custom')" @click="changeDateFilter(11)">
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.dateType === 11"></f7-icon>
                    <div slot="footer" v-if="query.dateType === 11 && query.minTime && query.maxTime">
                        <span>{{ query.minTime | moment($t('format.datetime.long-without-second')) }}</span>
                        <span>&nbsp;-&nbsp;</span>
                        <br/>
                        <span>{{ query.maxTime | moment($t('format.datetime.long-without-second')) }}</span>
                    </div>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <date-range-selection-sheet :title="$t('Custom Date Range')"
                                    :show.sync="showCustomDateRangeSheet"
                                    :min-time="query.minTime"
                                    :max-time="query.maxTime"
                                    @dateRange:change="changeCustomDateFilter">
        </date-range-selection-sheet>

        <f7-popover class="type-popover-menu" :opened="showTypePopover"
                    @popover:opened="showTypePopover = true" @popover:closed="showTypePopover = false">
            <f7-list>
                <f7-list-item :title="$t('All')" @click="changeTypeFilter(0)">
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.type === 0"></f7-icon>
                </f7-list-item>
                <f7-list-item :title="$t('Modify Balance')" @click="changeTypeFilter(1)">
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.type === 1"></f7-icon>
                </f7-list-item>
                <f7-list-item :title="$t('Income')" @click="changeTypeFilter(2)">
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.type === 2"></f7-icon>
                </f7-list-item>
                <f7-list-item :title="$t('Expense')" @click="changeTypeFilter(3)">
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.type === 3"></f7-icon>
                </f7-list-item>
                <f7-list-item :title="$t('Transfer')" @click="changeTypeFilter(4)">
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.type === 4"></f7-icon>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-popover class="category-popover-menu" :opened="showCategoryPopover"
                    @popover:opened="showCategoryPopover = true" @popover:closed="showCategoryPopover = false">
            <f7-list>
                <f7-list-item :title="$t('All')" @click="changeCategoryFilter('0')">
                    <f7-icon slot="media" f7="rectangle_badge_checkmark"></f7-icon>
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.categoryId === '0'"></f7-icon>
                </f7-list-item>
                <f7-list-item v-for="category in allCategories"
                              v-show="category.parentId > 0 && (!query.type || category.type === query.type - 1)"
                              :key="category.id"
                              :title="category.name"
                              @click="changeCategoryFilter(category.id)"
                >
                    <f7-icon slot="media"
                             :icon="category.icon | categoryIcon"
                             :style="category.color | categoryIconStyle('var(--default-icon-color)')">
                    </f7-icon>
                    <f7-icon slot="after"
                             class="list-item-checked"
                             f7="checkmark_alt"
                             v-if="query.categoryId === category.id">
                    </f7-icon>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-popover class="account-popover-menu" :opened="showAccountPopover"
                    @popover:opened="showAccountPopover = true" @popover:closed="showAccountPopover = false">
            <f7-list>
                <f7-list-item :title="$t('All')" @click="changeAccountFilter('0')">
                    <f7-icon slot="media" f7="rectangle_badge_checkmark"></f7-icon>
                    <f7-icon slot="after" class="list-item-checked" f7="checkmark_alt" v-if="query.accountId === '0'"></f7-icon>
                </f7-list-item>
                <f7-list-item v-for="account in allAccounts"
                              v-show="!account.hidden"
                              :key="account.id"
                              :title="account.name"
                              @click="changeAccountFilter(account.id)"
                >
                    <f7-icon slot="media"
                             :icon="account.icon | accountIcon"
                             :style="account.color | accountIconStyle('var(--default-icon-color)')">
                    </f7-icon>
                    <f7-icon slot="after"
                             class="list-item-checked"
                             f7="checkmark_alt"
                             v-if="query.accountId === account.id">
                    </f7-icon>
                </f7-list-item>
            </f7-list>
        </f7-popover>

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
            loading: true,
            loadingMore: false,
            transactionToDelete: null,
            showDatePopover: false,
            showTypePopover: false,
            showCategoryPopover: false,
            showAccountPopover: false,
            showCustomDateRangeSheet: false,
            showDeleteActionSheet: false
        };
    },
    computed: {
        defaultCurrency() {
            return this.$store.getters.currentUserDefaultCurrency || this.$t('default.currency');
        },
        query() {
            return this.$store.state.transactionsFilter;
        },
        transactions() {
            if (this.loading) {
                return [];
            }

            return this.$store.state.transactions;
        },
        noTransaction() {
            return this.$store.getters.noTransaction;
        },
        hasMoreTransaction() {
            return this.$store.getters.hasMoreTransaction;
        },
        allAccounts() {
            return this.$store.state.allAccountsMap;
        },
        allCategories() {
            return this.$store.state.allTransactionCategoriesMap;
        }
    },
    created() {
        const self = this;
        const query = self.$f7route.query;

        const dateParam = self.getDateParamByDateType(query.dateType ? parseInt(query.dateType) : undefined);

        this.$store.dispatch('initTransactionListFilter', {
            dateType: dateParam ? dateParam.dateType : undefined,
            maxTime: dateParam ? dateParam.maxTime : undefined,
            minTime: dateParam ? dateParam.minTime : undefined,
            type: query.type,
            categoryId: query.categoryId,
            accountId: query.accountId
        });

        this.reload(null);
    },
    methods: {
        onPageAfterIn() {
            if (this.$store.state.transactionListStateInvalid && !this.loading) {
                this.reload(null);
            }
        },
        reload(done) {
            const self = this;
            const router = self.$f7router;

            if (!done) {
                self.loading = true;
            }

            Promise.all([
                self.$store.dispatch('loadAllAccounts', { force: false }),
                self.$store.dispatch('loadAllCategories', { force: false })
            ]).then(() => {
                return self.$store.dispatch('getTransactions', {
                    reload: true,
                    autoExpand: true,
                    defaultCurrency: self.defaultCurrency
                });
            }).then(() => {
                if (done) {
                    done();
                }

                self.loading = false;
            }).catch(error => {
                self.loading = false;

                if (done) {
                    done();
                }

                if (!error.processed) {
                    self.$toast(error.message || error);

                    if (!done) {
                        router.back();
                    }
                }
            });
        },
        loadMore(autoExpand) {
            const self = this;

            if (!self.hasMoreTransaction) {
                return;
            }

            if (self.loadingMore || self.loading) {
                return;
            }

            self.loadingMore = true;

            self.$store.dispatch('getTransactions', {
                reload: false,
                autoExpand: autoExpand,
                defaultCurrency: self.defaultCurrency
            }).then(() => {
                self.loadingMore = false;
            }).catch(error => {
                self.loadingMore = false;

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        collapseTransactionMonthList(month, collapse) {
            this.$store.dispatch('collapseMonthInTransactionList', {
                month: month,
                collapse: collapse
            });
        },
        changeDateFilter(dateType) {
            if (dateType === 11) { // Custom
                this.showCustomDateRangeSheet = true;
                this.showDatePopover = false;
                return;
            } else if (this.query.dateType === dateType) {
                return;
            }

            const dateParam = this.getDateParamByDateType(dateType);

            if (!dateParam) {
                return;
            }

            this.$store.dispatch('updateTransactionListFilter', {
                dateType: dateParam.dateType,
                maxTime: dateParam.maxTime,
                minTime: dateParam.minTime
            });

            this.showDatePopover = false;
            this.reload(null);
        },
        changeCustomDateFilter(minTime, maxTime) {
            if (!minTime || !maxTime) {
                return;
            }

            this.$store.dispatch('updateTransactionListFilter', {
                dateType: 11,
                maxTime: maxTime,
                minTime: minTime
            });

            this.showCustomDateRangeSheet = false;

            this.reload(null);
        },
        changeTypeFilter(type) {
            if (this.query.type === type) {
                return;
            }

            let removeCategoryFilter = false;

            if (type && this.query.categoryId) {
                const category = this.allCategories[this.query.categoryId];

                if (category && category.type !== type - 1) {
                    removeCategoryFilter = true;
                }
            }

            this.$store.dispatch('updateTransactionListFilter', {
                type: type,
                categoryId: removeCategoryFilter ? '0' : undefined
            });

            this.showTypePopover = false;
            this.reload(null);
        },
        changeCategoryFilter(categoryId) {
            if (this.query.categoryId === categoryId) {
                return;
            }

            this.$store.dispatch('updateTransactionListFilter', {
                categoryId: categoryId
            });

            this.showCategoryPopover = false;
            this.reload(null);
        },
        changeAccountFilter(accountId) {
            if (this.query.accountId === accountId) {
                return;
            }

            this.$store.dispatch('updateTransactionListFilter', {
                accountId: accountId
            });

            this.showAccountPopover = false;
            this.reload(null);
        },
        changeKeywordFilter(keyword) {
            if (this.query.keyword === keyword) {
                return;
            }

            this.$store.dispatch('updateTransactionListFilter', {
                keyword: keyword
            });

            this.reload(null);
        },
        duplicate(transaction) {
            this.$f7router.navigate(`/transaction/add?id=${transaction.id}&type=${transaction.type}`);
        },
        edit(transaction) {
            this.$f7router.navigate(`/transaction/edit?id=${transaction.id}&type=${transaction.type}`);
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

            self.$store.dispatch('deleteTransaction', {
                transaction: transaction,
                defaultCurrency: self.defaultCurrency,
                beforeResolve: (done) => {
                    app.swipeout.delete($$(`#${self.$options.filters.transactionDomId(transaction)}`), () => {
                        done();
                    });
                }
            }).then(() => {
                self.$hideLoading();
            }).catch(error => {
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        getDateParamByDateType(dateType) {
            let maxTime = 0;
            let minTime = 0;

            if (dateType === 0) { // All
                maxTime = 0;
                minTime = 0;
            } else if (dateType === 1) { // Today
                maxTime = this.$utilities.getTodayLastUnixTime();
                minTime = this.$utilities.getTodayFirstUnixTime();
            } else if (dateType === 2) { // Yesterday
                maxTime = this.$utilities.getUnixTimeBeforeUnixTime(this.$utilities.getTodayLastUnixTime(), 1, 'days');
                minTime = this.$utilities.getUnixTimeBeforeUnixTime(this.$utilities.getTodayFirstUnixTime(), 1, 'days');
            } else if (dateType === 3) { // Last 7 days
                maxTime = this.$utilities.getUnixTime(new Date());
                minTime = this.$utilities.getUnixTimeBeforeUnixTime(maxTime, 7, 'days');
            } else if (dateType === 4) { // Last 30 days
                maxTime = this.$utilities.getUnixTime(new Date());
                minTime = this.$utilities.getUnixTimeBeforeUnixTime(maxTime, 30, 'days');
            } else if (dateType === 5) { // This week
                maxTime = this.$utilities.getThisWeekLastUnixTime();
                minTime = this.$utilities.getThisWeekFirstUnixTime();
            } else if (dateType === 6) { // Last week
                maxTime = this.$utilities.getUnixTimeBeforeUnixTime(this.$utilities.getThisWeekLastUnixTime(), 7, 'days');
                minTime = this.$utilities.getUnixTimeBeforeUnixTime(this.$utilities.getThisWeekFirstUnixTime(), 7, 'days');
            } else if (dateType === 7) { // This month
                maxTime = this.$utilities.getThisMonthLastUnixTime();
                minTime = this.$utilities.getThisMonthFirstUnixTime();
            } else if (dateType === 8) { // Last month
                maxTime = this.$utilities.getUnixTimeBeforeUnixTime(this.$utilities.getThisMonthLastUnixTime(), 1, 'months');
                minTime = this.$utilities.getUnixTimeBeforeUnixTime(this.$utilities.getThisMonthFirstUnixTime(), 1, 'months');
            } else if (dateType === 9) { // This year
                maxTime = this.$utilities.getThisYearLastUnixTime();
                minTime = this.$utilities.getThisYearFirstUnixTime();
            } else if (dateType === 10) { // Last year
                maxTime = this.$utilities.getUnixTimeBeforeUnixTime(this.$utilities.getThisYearLastUnixTime(), 1, 'years');
                minTime = this.$utilities.getUnixTimeBeforeUnixTime(this.$utilities.getThisYearFirstUnixTime(), 1, 'years');
            } else {
                return null;
            }

            return {
                dateType: dateType,
                maxTime: maxTime,
                minTime: minTime
            }
        }
    },
    filters: {
        transactionTypeName(type, allTransactionTypes) {
            if (type === allTransactionTypes.Income) {
                return 'Income';
            } else if (type === allTransactionTypes.Expense) {
                return 'Expense';
            } else if (type === allTransactionTypes.Transfer) {
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
        dateName(dateType, defaultName) {
            switch (dateType){
                case 1:
                    return 'Today';
                case 2:
                    return 'Yesterday';
                case 3:
                    return 'Recent 7 days';
                case 4:
                    return 'Recent 30 days';
                case 5:
                    return 'This week';
                case 6:
                    return 'Last week';
                case 7:
                    return 'This month';
                case 8:
                    return 'Last month';
                case 9:
                    return 'This year';
                case 10:
                    return 'Last year';
                case 11:
                    return 'Custom Date';
                default:
                    return defaultName;
            }
        },
        typeName(type, defaultName) {
            switch (type){
                case 1:
                    return 'Modify Balance';
                case 2:
                    return 'Income';
                case 3:
                    return 'Expense';
                case 4:
                    return 'Transfer';
                default:
                    return defaultName;
            }
        },
        categoryName(categoryId, allCategories, defaultName) {
            if (allCategories[categoryId]) {
                return allCategories[categoryId].name;
            }

            return defaultName;
        },
        accountName(accountId, allAccounts, defaultName) {
            if (allAccounts[accountId]) {
                return allAccounts[accountId].name;
            }

            return defaultName;
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
    margin-left: 4px;
}

.transaction-info .item-media + .item-inner {
    margin-left: 10px;
}

.transaction-date {
    width: 25px;
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
}

.transaction-comment {
    font-size: 13px;
    line-height: 20px;
    padding-top: 2px;
    padding-bottom: 2px;
}

.transaction-footer {
    padding-top: 4px;
}

.transaction-info .item-text + .item-footer .transaction-footer {
    padding-top: 2px;
}

.transaction-footer > span {
    margin-right: 4px;
}

.transaction-amount {
    color: var(--f7-list-item-after-text-color);
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

.date-popover-menu .popover-inner, .category-popover-menu .popover-inner, .account-popover-menu .popover-inner {
    max-height: 400px;
    overflow-Y: auto;
}
</style>
