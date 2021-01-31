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
                <span :class="{ 'tabbar-item-changed': query.maxTime > 0 || query.minTime > 0 }">{{ query.dateType | dateRangeName(allDateRanges, 'Date') | localized }}</span>
            </f7-link>
            <f7-link class="tabbar-text-with-ellipsis" popover-open=".type-popover-menu">
                <span :class="{ 'tabbar-item-changed': query.type > 0 }">{{ query.type | typeName('Type') | localized }}</span>
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
                            <span>{{ transactionMonthList.yearMonth | moment($t('format.yearMonth.long')) }}</span>
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
                                            {{ `datetime.${transaction.dayOfWeek}.short` | localized }}
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
                                        {{ transaction.type | transactionTypeName($constants.transaction.allTransactionTypes) | localized }}
                                    </span>
                                </div>
                                <div slot="text" class="transaction-comment" v-if="transaction.comment">
                                    <span>{{ transaction.comment }}</span>
                                </div>
                                <div slot="footer" class="transaction-footer">
                                    <span>{{ transaction.time | moment($t('format.hourMinute.long')) }}</span>
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
                <f7-list-item v-for="dateRange in allDateRanges"
                              :key="dateRange.type"
                              :title="dateRange.name | localized"
                              @click="changeDateFilter(dateRange.type)">
                    <f7-icon slot="after" class="list-item-checked-icon" f7="checkmark_alt" v-if="query.dateType === dateRange.type"></f7-icon>
                    <div slot="footer"
                         v-if="dateRange.type === $constants.datetime.allDateRanges.Custom.type && query.dateType === $constants.datetime.allDateRanges.Custom.type && query.minTime && query.maxTime">
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
                    <f7-icon slot="after" class="list-item-checked-icon" f7="checkmark_alt" v-if="query.type === 0"></f7-icon>
                </f7-list-item>
                <f7-list-item :title="$t('Modify Balance')" @click="changeTypeFilter(1)">
                    <f7-icon slot="after" class="list-item-checked-icon" f7="checkmark_alt" v-if="query.type === 1"></f7-icon>
                </f7-list-item>
                <f7-list-item :title="$t('Income')" @click="changeTypeFilter(2)">
                    <f7-icon slot="after" class="list-item-checked-icon" f7="checkmark_alt" v-if="query.type === 2"></f7-icon>
                </f7-list-item>
                <f7-list-item :title="$t('Expense')" @click="changeTypeFilter(3)">
                    <f7-icon slot="after" class="list-item-checked-icon" f7="checkmark_alt" v-if="query.type === 3"></f7-icon>
                </f7-list-item>
                <f7-list-item :title="$t('Transfer')" @click="changeTypeFilter(4)">
                    <f7-icon slot="after" class="list-item-checked-icon" f7="checkmark_alt" v-if="query.type === 4"></f7-icon>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-popover class="category-popover-menu" :opened="showCategoryPopover"
                    @popover:opened="showCategoryPopover = true" @popover:closed="showCategoryPopover = false">
            <f7-list accordion-list>
                <f7-list-item :title="$t('All')" @click="changeCategoryFilter('0')">
                    <f7-icon slot="media" f7="rectangle_badge_checkmark"></f7-icon>
                    <f7-icon slot="after" class="list-item-checked-icon" f7="checkmark_alt" v-if="query.categoryId === '0'"></f7-icon>
                </f7-list-item>
            </f7-list>
            <f7-list accordion-list
                     class="no-margin-vertical"
                     v-for="(categories, categoryType) in allPrimaryCategories"
                     v-show="!query.type || parseInt(categoryType) === query.type - 1"
                     :key="categoryType"
            >
                <f7-list-item divider :title="(parseInt(categoryType) + 1) | typeName('Type') | localized"></f7-list-item>
                <f7-list-item accordion-item
                              v-for="category in categories"
                              :key="category.id"
                              :class="category | categoryListItemCheckedClass(query.categoryId)"
                              :title="category.name"
                >
                    <f7-icon slot="media"
                             :icon="category.icon | categoryIcon"
                             :style="category.color | categoryIconStyle('var(--default-icon-color)')">
                    </f7-icon>
                    <f7-accordion-content>
                        <f7-list class="padding-left">
                            <f7-list-item :title="$t('All')" @click="changeCategoryFilter(category.id)">
                                <f7-icon slot="media" f7="rectangle_badge_checkmark"></f7-icon>
                                <f7-icon slot="after" class="list-item-checked-icon" f7="checkmark_alt" v-if="query.categoryId === category.id"></f7-icon>
                            </f7-list-item>
                            <f7-list-item v-for="subCategory in category.subCategories"
                                          :key="subCategory.id"
                                          :title="subCategory.name"
                                          @click="changeCategoryFilter(subCategory.id)"
                            >
                                <f7-icon slot="media"
                                         :icon="subCategory.icon | categoryIcon"
                                         :style="subCategory.color | categoryIconStyle('var(--default-icon-color)')">
                                </f7-icon>
                                <f7-icon slot="after"
                                         class="list-item-checked-icon"
                                         f7="checkmark_alt"
                                         v-if="query.categoryId === subCategory.id">
                                </f7-icon>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-content>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-popover class="account-popover-menu" :opened="showAccountPopover"
                    @popover:opened="showAccountPopover = true" @popover:closed="showAccountPopover = false">
            <f7-list>
                <f7-list-item :title="$t('All')" @click="changeAccountFilter('0')">
                    <f7-icon slot="media" f7="rectangle_badge_checkmark"></f7-icon>
                    <f7-icon slot="after" class="list-item-checked-icon" f7="checkmark_alt" v-if="query.accountId === '0'"></f7-icon>
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
                             class="list-item-checked-icon"
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
            if (this.query.accountId && this.query.accountId !== '0') {
                const account = this.allAccounts[this.query.accountId];

                if (account && account.currency && account.currency !== this.$constants.currency.parentAccountCurrencyPlaceholder) {
                    return account.currency;
                }
            }

            return this.$store.getters.currentUserDefaultCurrency || this.$t('default.currency');
        },
        firstDayOfWeek() {
            if (this.$utilities.isNumber(this.$store.getters.currentUserFirstDayOfWeek)) {
                return this.$store.getters.currentUserFirstDayOfWeek;
            }

            if (this.$constants.datetime.allWeekDays[this.$t('default.firstDayOfWeek')]) {
                return this.$constants.datetime.allWeekDays[this.$t('default.firstDayOfWeek')].type;
            }

            return 0;
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
        },
        allPrimaryCategories() {
            return this.$store.state.allTransactionCategories;
        },
        allDateRanges() {
            return this.$constants.datetime.allDateRanges;
        }
    },
    created() {
        const self = this;
        const query = self.$f7route.query;

        let dateRange = self.$utilities.getDateRangeByDateType(query.dateType ? parseInt(query.dateType) : undefined, self.firstDayOfWeek);

        if (!dateRange &&
            query.dateType === self.$constants.datetime.allDateRanges.Custom.type.toString() &&
            parseInt(query.maxTime) > 0 && parseInt(query.minTime) > 0) {
            dateRange = {
                dateType: parseInt(query.dateType),
                maxTime: parseInt(query.maxTime),
                minTime: parseInt(query.minTime)
            };
        }

        this.$store.dispatch('initTransactionListFilter', {
            dateType: dateRange ? dateRange.dateType : undefined,
            maxTime: dateRange ? dateRange.maxTime : undefined,
            minTime: dateRange ? dateRange.minTime : undefined,
            type: parseInt(query.type) > 0 ? parseInt(query.type) : undefined,
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
                return self.$store.dispatch('loadTransactions', {
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

            self.$store.dispatch('loadTransactions', {
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
            if (dateType === this.$constants.datetime.allDateRanges.Custom.type) { // Custom
                this.showCustomDateRangeSheet = true;
                this.showDatePopover = false;
                return;
            } else if (this.query.dateType === dateType) {
                return;
            }

            const dateRange = this.$utilities.getDateRangeByDateType(dateType, this.firstDayOfWeek);

            if (!dateRange) {
                return;
            }

            this.$store.dispatch('updateTransactionListFilter', {
                dateType: dateRange.dateType,
                maxTime: dateRange.maxTime,
                minTime: dateRange.minTime
            });

            this.showDatePopover = false;
            this.reload(null);
        },
        changeCustomDateFilter(minTime, maxTime) {
            if (!minTime || !maxTime) {
                return;
            }

            this.$store.dispatch('updateTransactionListFilter', {
                dateType: this.$constants.datetime.allDateRanges.Custom.type,
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
        dateRangeName(dateRangeType, allDateRanges, defaultName) {
            if (dateRangeType === allDateRanges.All.type) {
                return defaultName;
            }

            for (let dateRangeField in allDateRanges) {
                if (!Object.prototype.hasOwnProperty.call(allDateRanges, dateRangeField)) {
                    continue;
                }

                const dateRange = allDateRanges[dateRangeField];

                if (dateRange && dateRange.type === dateRangeType && dateRange.name) {
                    return dateRange.name;
                }
            }

            return defaultName;
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
        categoryListItemCheckedClass(category, queryCategoryId) {
            if (category.id === queryCategoryId) {
                return {
                    'list-item-checked': true
                };
            }

            for (let i = 0; i < category.subCategories.length; i++) {
                if (category.subCategories[i].id === queryCategoryId) {
                    return {
                        'list-item-checked': true
                    };
                }
            }

            return [];
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
