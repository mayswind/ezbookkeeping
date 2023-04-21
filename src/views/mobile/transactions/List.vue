<template>
    <f7-page with-subnavbar
             ptr
             infinite
             :infinite-preloader="loadingMore"
             :infinite-distance="600"
             @ptr:refresh="reload"
             @page:afterin="onPageAfterIn"
             @infinite="loadMore(true)">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Transaction List')"></f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link icon-f7="plus" v-if="canAddTransaction" :href="`/transaction/add?type=${query.type}&categoryId=${query.categoryId}&accountId=${query.accountId}`"></f7-link>
            </f7-nav-right>

            <f7-subnavbar :inner="false">
                <f7-searchbar
                    custom-searchs
                    :value="query.keyword"
                    :placeholder="$t('Search transaction comment')"
                    :disable-button-text="$t('Cancel')"
                    @change="changeKeywordFilter($event.target.value)"
                ></f7-searchbar>
            </f7-subnavbar>
        </f7-navbar>

        <f7-toolbar tabbar bottom>
            <f7-link class="tabbar-text-with-ellipsis" popover-open=".date-popover-menu">
                <span :class="{ 'tabbar-item-changed': query.maxTime > 0 || query.minTime > 0 }">{{ queryDateRangeName }}</span>
            </f7-link>
            <f7-link class="tabbar-text-with-ellipsis" popover-open=".type-popover-menu">
                <span :class="{ 'tabbar-item-changed': query.type > 0 }">{{ queryTransactionTypeName }}</span>
            </f7-link>
            <f7-link class="tabbar-text-with-ellipsis" popover-open=".category-popover-menu" :class="{ 'disabled': query.type === 1 }">
                <span :class="{ 'tabbar-item-changed': query.categoryId > 0 }">{{ queryCategoryName }}</span>
            </f7-link>
            <f7-link class="tabbar-text-with-ellipsis" popover-open=".account-popover-menu">
                <span :class="{ 'tabbar-item-changed': query.accountId > 0 }">{{ queryAccountName }}</span>
            </f7-link>
        </f7-toolbar>

        <f7-block class="combination-list-wrapper margin-vertical skeleton-text"
                  :key="blockIdx" v-for="blockIdx in [ 1, 2 ]" v-if="loading">
            <f7-accordion-item>
                <f7-block-title>
                    <f7-accordion-toggle>
                        <f7-list strong inset dividers media-list
                                 class="transaction-amount-list combination-list-header combination-list-opened">
                            <f7-list-item>
                                <template #title>
                                    <small>YYYY-MM</small>
                                    <small class="transaction-amount-statistics" v-if="showTotalAmountInTransactionListPage">
                                        <span>0.00 USD</span>
                                        <span>0.00 USD</span>
                                    </small>
                                    <f7-icon class="combination-list-chevron-icon" f7="chevron_up"></f7-icon>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-toggle>
                </f7-block-title>
                <f7-accordion-content style="height: auto">
                    <f7-list strong inset dividers media-list accordion-list class="transaction-info-list combination-list-content">
                        <f7-list-item link="#" chevron-center class="transaction-info"
                                      :key="itemIdx" v-for="itemIdx in (blockIdx === 1 ? [ 1, 2, 3, 4, 5 ] : [ 1, 2, 3 ])">
                            <template #media>
                                <div class="display-flex flex-direction-column transaction-date">
                                    <span class="transaction-day full-line flex-direction-column">DD</span>
                                    <span class="transaction-day-of-week full-line flex-direction-column">Sun</span>
                                </div>
                            </template>
                            <template #inner>
                                <div class="display-flex no-padding-horizontal">
                                    <div class="item-media">
                                        <div class="transaction-icon display-flex align-items-center">
                                            <f7-icon f7="app_fill"></f7-icon>
                                        </div>
                                    </div>
                                    <div class="actual-item-inner">
                                        <div class="item-title-row">
                                            <div class="item-title">
                                                <div class="transaction-category-name no-padding">
                                                    <span>Category</span>
                                                </div>
                                            </div>
                                            <div class="item-after">
                                                <div class="no-padding transaction-amount">
                                                    <span>0.00 USD</span>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="item-footer">
                                            <div class="no-padding-horizontal transaction-footer">
                                                <span>HH:mm</span>
                                                <span>·</span>
                                                <span>Source Account</span>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </template>
                        </f7-list-item>
                    </f7-list>
                </f7-accordion-content>
            </f7-accordion-item>
        </f7-block>

        <f7-list strong inset dividers class="margin-vertical" v-if="!loading && noTransaction">
            <f7-list-item :title="$t('No transaction data')"></f7-list-item>
        </f7-list>

        <f7-block class="combination-list-wrapper margin-vertical"
                  :key="transactionMonthList.yearMonth" v-for="(transactionMonthList) in transactions">
            <f7-accordion-item :opened="transactionMonthList.opened"
                               @accordion:open="collapseTransactionMonthList(transactionMonthList, false)"
                               @accordion:close="collapseTransactionMonthList(transactionMonthList, true)">
                <f7-block-title>
                    <f7-accordion-toggle>
                        <f7-list strong inset dividers media-list
                                 class="transaction-amount-list combination-list-header"
                                 :class="transactionMonthList.opened ? 'combination-list-opened' : 'combination-list-closed'">
                            <f7-list-item>
                                <template #title>
                                    <small>
                                        <span>{{ $utilities.formatTime(transactionMonthList.yearMonth, $t('format.yearMonth.long')) }}</span>
                                    </small>
                                    <small class="transaction-amount-statistics" v-if="showTotalAmountInTransactionListPage && transactionMonthList.totalAmount">
                                        <span class="text-color-red">
                                            {{ getDisplayMonthTotalAmount(transactionMonthList.totalAmount.income, defaultCurrency, '+', transactionMonthList.totalAmount.incompleteIncome) }}
                                        </span>
                                        <span class="text-color-teal">
                                            {{ getDisplayMonthTotalAmount(transactionMonthList.totalAmount.expense, defaultCurrency, '-', transactionMonthList.totalAmount.incompleteExpense) }}
                                        </span>
                                    </small>
                                    <f7-icon class="combination-list-chevron-icon" :f7="transactionMonthList.opened ? 'chevron_up' : 'chevron_down'"></f7-icon>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-toggle>
                </f7-block-title>
                <f7-accordion-content :style="{ height: transactionMonthList.opened ? 'auto' : '' }">
                    <f7-list strong inset dividers media-list accordion-list class="transaction-info-list combination-list-content">
                        <f7-list-item swipeout chevron-center
                                      class="transaction-info"
                                      :id="getTransactionDomId(transaction)"
                                      :link="transaction.type !== $constants.transaction.allTransactionTypes.ModifyBalance ? `/transaction/detail?id=${transaction.id}&type=${transaction.type}` : null"
                                      :key="transaction.id"
                                      v-for="(transaction, idx) in transactionMonthList.items"
                        >
                            <template #media>
                                <div class="display-flex flex-direction-column transaction-date" :style="getTransactionDateStyle(transaction, idx > 0 ? transactionMonthList.items[idx - 1] : null)">
                                    <span class="transaction-day full-line flex-direction-column">
                                        {{ transaction.day }}
                                    </span>
                                    <span class="transaction-day-of-week full-line flex-direction-column">
                                        {{ $t(`datetime.${transaction.dayOfWeek}.short`) }}
                                    </span>
                                </div>
                            </template>
                            <template #inner>
                                <div class="display-flex no-padding-horizontal">
                                    <div class="item-media">
                                        <div class="transaction-icon display-flex align-items-center">
                                            <ItemIcon icon-type="category"
                                                      :icon-id="transaction.category.icon"
                                                      :color="transaction.category.color"
                                                      v-if="transaction.category && transaction.category.color"></ItemIcon>
                                            <f7-icon v-else-if="!transaction.category || !transaction.category.color"
                                                     f7="pencil_ellipsis_rectangle">
                                            </f7-icon>
                                        </div>
                                    </div>
                                    <div class="actual-item-inner">
                                        <div class="item-title-row">
                                            <div class="item-title">
                                                <div class="transaction-category-name no-padding">
                                                    <span v-if="transaction.type === $constants.transaction.allTransactionTypes.ModifyBalance">
                                                        {{ $t('Modify Balance') }}
                                                    </span>
                                                        <span v-else-if="transaction.type !== $constants.transaction.allTransactionTypes.ModifyBalance && transaction.category">
                                                        {{ transaction.category.name }}
                                                    </span>
                                                        <span v-else-if="transaction.type !== $constants.transaction.allTransactionTypes.ModifyBalance && !transaction.category">
                                                        {{ getTransactionTypeName(transaction.type, 'Transaction') }}
                                                    </span>
                                                </div>
                                            </div>
                                            <div class="item-after">
                                                <div class="transaction-amount" v-if="transaction.sourceAccount"
                                                     :class="{ 'text-color-teal': transaction.type === $constants.transaction.allTransactionTypes.Expense, 'text-color-red': transaction.type === $constants.transaction.allTransactionTypes.Income }">
                                                    <span v-if="!query.accountId || query.accountId === '0' || (transaction.sourceAccount && (transaction.sourceAccount.id === query.accountId || transaction.sourceAccount.parentId === query.accountId))">{{ getDisplayAmount(transaction.sourceAmount, transaction.sourceAccount.currency, transaction.hideAmount) }}</span>
                                                    <span v-else-if="query.accountId && query.accountId !== '0' && transaction.destinationAccount && (transaction.destinationAccount.id === query.accountId || transaction.destinationAccount.parentId === query.accountId)">{{ getDisplayAmount(transaction.destinationAmount, transaction.destinationAccount.currency, transaction.hideAmount) }}</span>
                                                    <span v-else></span>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="item-text">
                                            <div class="transaction-comment" v-if="transaction.comment">
                                                <span>{{ transaction.comment }}</span>
                                            </div>
                                        </div>
                                        <div class="item-footer">
                                            <div class="transaction-footer">
                                                <span>{{ $utilities.formatUnixTime(transaction.time, $t('format.hourMinute.long'), transaction.utcOffset, currentTimezoneOffsetMinutes) }}</span>
                                                <span v-if="transaction.utcOffset !== currentTimezoneOffsetMinutes">{{ `(UTC${$utilities.getUtcOffsetByUtcOffsetMinutes(transaction.utcOffset)})` }}</span>
                                                <span v-if="transaction.sourceAccount">·</span>
                                                <span v-if="transaction.sourceAccount">{{ transaction.sourceAccount.name }}</span>
                                                <span v-if="transaction.sourceAccount && transaction.type === $constants.transaction.allTransactionTypes.Transfer && transaction.destinationAccount && transaction.sourceAccount.id !== transaction.destinationAccount.id">→</span>
                                                <span v-if="transaction.sourceAccount && transaction.type === $constants.transaction.allTransactionTypes.Transfer && transaction.destinationAccount && transaction.sourceAccount.id !== transaction.destinationAccount.id">{{ transaction.destinationAccount.name }}</span>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </template>
                            <f7-swipeout-actions right>
                                <f7-swipeout-button color="primary" close
                                                    :text="$t('Duplicate')"
                                                    v-if="transaction.type !== $constants.transaction.allTransactionTypes.ModifyBalance"
                                                    @click="duplicate(transaction)"></f7-swipeout-button>
                                <f7-swipeout-button color="orange" close
                                                    :text="$t('Edit')"
                                                    v-if="transaction.editable && transaction.type !== $constants.transaction.allTransactionTypes.ModifyBalance"
                                                    @click="edit(transaction)"></f7-swipeout-button>
                                <f7-swipeout-button color="red" class="padding-left padding-right"
                                                    v-if="transaction.editable"
                                                    @click="remove(transaction, false)">
                                    <f7-icon f7="trash"></f7-icon>
                                </f7-swipeout-button>
                            </f7-swipeout-actions>
                        </f7-list-item>
                    </f7-list>
                </f7-accordion-content>
            </f7-accordion-item>
        </f7-block>

        <f7-block class="text-align-center" v-if="!loading && hasMoreTransaction">
            <f7-link :class="{ 'disabled': loadingMore }" href="#" @click="loadMore(false)">{{ $t('Load More') }}</f7-link>
        </f7-block>

        <f7-popover class="date-popover-menu" :opened="showDatePopover"
                    @popover:open="scrollPopoverToSelectedItem"
                    @popover:opened="showDatePopover = true" @popover:closed="showDatePopover = false">
            <f7-list dividers>
                <f7-list-item :title="$t(dateRange.name)"
                              :class="{ 'list-item-selected': query.dateType === dateRange.type }"
                              :key="dateRange.type"
                              v-for="dateRange in allDateRanges"
                              @click="changeDateFilter(dateRange.type)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.dateType === dateRange.type"></f7-icon>
                    </template>
                    <template #footer>
                        <div v-if="dateRange.type === $constants.datetime.allDateRanges.Custom.type && query.dateType === $constants.datetime.allDateRanges.Custom.type && query.minTime && query.maxTime">
                            <span>{{ $utilities.formatUnixTime(query.minTime, $t('format.datetime.long-without-second')) }}</span>
                            <span>&nbsp;-&nbsp;</span>
                            <br/>
                            <span>{{ $utilities.formatUnixTime(query.maxTime, $t('format.datetime.long-without-second')) }}</span>
                        </div>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <date-range-selection-sheet :title="$t('Custom Date Range')"
                                    :min-time="query.minTime"
                                    :max-time="query.maxTime"
                                    v-model:show="showCustomDateRangeSheet"
                                    @dateRange:change="changeCustomDateFilter">
        </date-range-selection-sheet>

        <f7-popover class="type-popover-menu" :opened="showTypePopover"
                    @popover:open="scrollPopoverToSelectedItem"
                    @popover:opened="showTypePopover = true" @popover:closed="showTypePopover = false">
            <f7-list dividers>
                <f7-list-item :class="{ 'list-item-selected': query.type === 0 }" :title="$t('All')" @click="changeTypeFilter(0)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.type === 0"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item :class="{ 'list-item-selected': query.type === 1 }" :title="$t('Modify Balance')" @click="changeTypeFilter(1)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.type === 1"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item :class="{ 'list-item-selected': query.type === 2 }" :title="$t('Income')" @click="changeTypeFilter(2)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.type === 2"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item :class="{ 'list-item-selected': query.type === 3 }" :title="$t('Expense')" @click="changeTypeFilter(3)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.type === 3"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item :class="{ 'list-item-selected': query.type === 4 }" :title="$t('Transfer')" @click="changeTypeFilter(4)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.type === 4"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-popover class="category-popover-menu" :opened="showCategoryPopover"
                    @popover:open="scrollPopoverToSelectedItem"
                    @popover:opened="showCategoryPopover = true" @popover:closed="showCategoryPopover = false">
            <f7-list dividers accordion-list>
                <f7-list-item :class="{ 'list-item-selected': query.categoryId === '0' }" :title="$t('All')" @click="changeCategoryFilter('0')">
                    <template #media>
                        <f7-icon f7="rectangle_badge_checkmark"></f7-icon>
                    </template>
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.categoryId === '0'"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
            <f7-list dividers accordion-list
                     class="no-margin-vertical"
                     :key="categoryType"
                     v-for="(categories, categoryType) in allPrimaryCategories"
                     v-show="!query.type || $utilities.categroyTypeToTransactionType(parseInt(categoryType)) === query.type"
            >
                <f7-list-item divider :title="getTransactionTypeName($utilities.categroyTypeToTransactionType(parseInt(categoryType)), 'Type')"></f7-list-item>
                <f7-list-item accordion-item
                              :title="category.name"
                              :class="getCategoryListItemCheckedClass(category, query.categoryId)"
                              :key="category.id"
                              v-for="category in categories"
                >
                    <template #media>
                        <ItemIcon icon-type="category" :icon-id="category.icon" :color="category.color"></ItemIcon>
                    </template>
                    <f7-accordion-content>
                        <f7-list dividers class="padding-left">
                            <f7-list-item :class="{ 'list-item-selected': query.categoryId === category.id }" :title="$t('All')" @click="changeCategoryFilter(category.id)">
                                <template #media>
                                    <f7-icon f7="rectangle_badge_checkmark"></f7-icon>
                                </template>
                                <template #after>
                                    <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.categoryId === category.id"></f7-icon>
                                </template>
                            </f7-list-item>
                            <f7-list-item :title="subCategory.name"
                                          :class="{ 'list-item-selected': query.categoryId === subCategory.id }"
                                          :key="subCategory.id"
                                          v-for="subCategory in category.subCategories"
                                          @click="changeCategoryFilter(subCategory.id)"
                            >
                                <template #media>
                                    <ItemIcon icon-type="category" :icon-id="subCategory.icon" :color="subCategory.color"></ItemIcon>
                                </template>
                                <template #after>
                                    <f7-icon class="list-item-checked-icon"
                                             f7="checkmark_alt"
                                             v-if="query.categoryId === subCategory.id">
                                    </f7-icon>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-content>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-popover class="account-popover-menu" :opened="showAccountPopover"
                    @popover:open="scrollPopoverToSelectedItem"
                    @popover:opened="showAccountPopover = true" @popover:closed="showAccountPopover = false">
            <f7-list dividers>
                <f7-list-item :class="{ 'list-item-selected': query.accountId === '0' }" :title="$t('All')" @click="changeAccountFilter('0')">
                    <template #media>
                        <f7-icon f7="rectangle_badge_checkmark"></f7-icon>
                    </template>
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.accountId === '0'"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item :title="account.name"
                              :class="{ 'list-item-selected': query.accountId === account.id }"
                              :key="account.id"
                              v-for="account in allAccounts"
                              v-show="!account.hidden"
                              @click="changeAccountFilter(account.id)"
                >
                    <template #media>
                        <ItemIcon icon-type="account" :icon-id="account.icon" :color="account.color"></ItemIcon>
                    </template>
                    <template #after>
                        <f7-icon class="list-item-checked-icon"
                                 f7="checkmark_alt"
                                 v-if="query.accountId === account.id">
                        </f7-icon>
                    </template>
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
    props: [
        'f7route',
        'f7router'
    ],
    data() {
        const self = this;

        return {
            loading: true,
            loadingError: null,
            loadingMore: false,
            transactionToDelete: null,
            showTotalAmountInTransactionListPage: self.$settings.isShowTotalAmountInTransactionListPage(),
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

            return this.$store.getters.currentUserDefaultCurrency;
        },
        canAddTransaction() {
            if (this.query.accountId && this.query.accountId !== '0') {
                const account = this.allAccounts[this.query.accountId];

                if (account && account.type === this.$constants.account.allAccountTypes.MultiSubAccounts) {
                    return false;
                }
            }

            return true;
        },
        currentTimezoneOffsetMinutes() {
            return this.$utilities.getTimezoneOffsetMinutes();
        },
        firstDayOfWeek() {
            return this.$store.getters.currentUserFirstDayOfWeek;
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
        },
        queryDateRangeName() {
            if (this.query.dateType === this.allDateRanges.All.type) {
                return this.$t('Date');
            }

            for (let dateRangeField in this.allDateRanges) {
                if (!Object.prototype.hasOwnProperty.call(this.allDateRanges, dateRangeField)) {
                    continue;
                }

                const dateRange = this.allDateRanges[dateRangeField];

                if (dateRange && dateRange.type === this.query.dateType && dateRange.name) {
                    return this.$t(dateRange.name);
                }
            }

            return this.$t('Date');
        },
        queryTransactionTypeName() {
            return this.getTransactionTypeName(this.query.type, 'Type');
        },
        queryCategoryName() {
            return this.$utilities.getNameByKeyValue(this.allCategories, this.query.categoryId, null, 'name', this.$t('Category'));
        },
        queryAccountName() {
            return this.$utilities.getNameByKeyValue(this.allAccounts, this.query.accountId, null, 'name', this.$t('Account'));
        }
    },
    created() {
        const self = this;
        const query = self.f7route.query;

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

            this.$routeBackOnError(this.f7router, 'loadingError');
        },
        reload(done) {
            const self = this;

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
                if (error.processed || done) {
                    self.loading = false;
                }

                if (done) {
                    done();
                }

                if (!error.processed) {
                    if (!done) {
                        self.loadingError = error;
                    }

                    self.$toast(error.message || error);
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

                if (category && category.type !== this.$utilities.transactionTypeToCategroyType(type)) {
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
            this.f7router.navigate(`/transaction/add?id=${transaction.id}&type=${transaction.type}`);
        },
        edit(transaction) {
            this.f7router.navigate(`/transaction/edit?id=${transaction.id}&type=${transaction.type}`);
        },
        remove(transaction, confirm) {
            const self = this;

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
                    self.$ui.onSwipeoutDeleted(self.getTransactionDomId(transaction), done);
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
        scrollPopoverToSelectedItem(event) {
            if (!event || !event.$el || !event.$el.length) {
                return;
            }

            const container = event.$el.find('.popover-inner');
            const selectedItem = event.$el.find('li.list-item-selected');

            if (!container.length || !selectedItem.length) {
                return;
            }

            let targetPos = selectedItem.offset().top - container.offset().top - parseInt(container.css('padding-top'), 10)
                - (container.outerHeight() - selectedItem.outerHeight()) / 2;

            if (targetPos <= 0) {
                return;
            }

            container.scrollTop(targetPos);
        },
        getDisplayAmount(amount, currency, hideAmount) {
            if (hideAmount) {
                return this.$locale.getDisplayCurrency('***', currency);
            }

            return this.$locale.getDisplayCurrency(amount, currency);
        },
        getDisplayMonthTotalAmount(amount, currency, symbol, incomplete) {
            const displayAmount = this.$locale.getDisplayCurrency(amount, currency);
            return symbol + displayAmount + (incomplete ? '+' : '');
        },
        getTransactionTypeName(type, defaultName) {
            switch (type){
                case this.$constants.transaction.allTransactionTypes.ModifyBalance:
                    return this.$t('Modify Balance');
                case this.$constants.transaction.allTransactionTypes.Income:
                    return this.$t('Income');
                case this.$constants.transaction.allTransactionTypes.Expense:
                    return this.$t('Expense');
                case this.$constants.transaction.allTransactionTypes.Transfer:
                    return this.$t('Transfer');
                default:
                    return this.$t(defaultName);
            }
        },
        getTransactionDomId(transaction) {
            return 'transaction_' + transaction.id;
        },
        getTransactionDateStyle(transaction, previousTransaction) {
            if (!previousTransaction || transaction.day !== previousTransaction.day) {
                return {};
            }

            return {
                color: 'transparent'
            }
        },
        getCategoryListItemCheckedClass(category, queryCategoryId) {
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
        }
    }
};
</script>

<style>
.list.transaction-amount-list .transaction-amount-statistics > span {
    margin-left: 8px;
    font-weight: normal;
}

.list.transaction-info-list li.transaction-info .item-media + .item-inner {
    margin-left: 0;
}

.list.transaction-info-list li.transaction-info .actual-item-inner {
    width: 100%;
    margin-left: 10px;
}

.list.transaction-info-list li.transaction-info .transaction-date {
    width: 25px;
    margin-right: 6px;
}

.list.transaction-info-list li.transaction-info .transaction-day {
    opacity: 0.6;
    font-size: 16px;
    font-weight: bold;
    text-align: left;
}

.list.transaction-info-list li.transaction-info .transaction-day-of-week {
    opacity: 0.6;
    font-size: 12px;
}

.list.transaction-info-list li.transaction-info .transaction-comment {
    font-size: 13px;
    line-height: 20px;
    padding-top: 2px;
    padding-bottom: 2px;
}

.list.transaction-info-list li.transaction-info .transaction-footer {
    padding-top: 4px;
}

.list.transaction-info-list li.transaction-info .transaction-info .item-text + .item-footer .transaction-footer {
    padding-top: 2px;
}

.list.transaction-info-list li.transaction-info .transaction-footer > span {
    margin-right: 4px;
}

.list.transaction-info-list li.transaction-info .transaction-amount {
    color: var(--f7-list-item-after-text-color);
    overflow: hidden;
    text-overflow: ellipsis;
}

.list.transaction-info-list li.transaction-info .transaction-info .item-after {
    max-width: 70%;
}

.list.transaction-info-list li.transaction-info .transaction-category-name {
    overflow: hidden;
    text-overflow: ellipsis;
}

.date-popover-menu .popover-inner, .category-popover-menu .popover-inner, .account-popover-menu .popover-inner {
    max-height: 400px;
    overflow-Y: auto;
}
</style>
