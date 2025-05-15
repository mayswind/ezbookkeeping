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
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title>
                <f7-link popover-open=".chart-data-type-popover-menu">
                    <span style="color: var(--f7-text-color)">{{ displayPageTypeName }}</span>
                    <f7-icon class="page-title-bar-icon" color="gray" style="opacity: 0.5" f7="chevron_down_circle_fill"></f7-icon>
                </f7-link>
            </f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link icon-f7="plus" :class="{ 'disabled': !canAddTransaction }" :href="`/transaction/add?type=${query.type}&categoryId=${queryAllFilterCategoryIdsCount === 1 ? query.categoryIds : ''}&accountId=${queryAllFilterAccountIdsCount === 1 ? query.accountIds : ''}&tagIds=${query.tagIds || ''}`"></f7-link>
            </f7-nav-right>

            <f7-subnavbar :inner="false">
                <f7-searchbar
                    custom-searchs
                    :value="query.keyword"
                    :placeholder="tt('Search transaction description')"
                    :disable-button-text="tt('Cancel')"
                    @change="changeKeywordFilter($event.target.value)"
                ></f7-searchbar>
            </f7-subnavbar>
        </f7-navbar>

        <f7-popover class="chart-data-type-popover-menu"
                    v-model:opened="showTransactionListPageTypePopover">
            <f7-list dividers>
                <f7-list-item :title="tt(type.name)"
                              :class="{ 'list-item-selected': pageType === type.type }"
                              :key="type.type"
                              v-for="type in TransactionListPageType.values()"
                              @click="changePageType(type.type); showTransactionListPageTypePopover = false">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="pageType === type.type"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-toolbar tabbar bottom class="toolbar-item-auto-size transaction-list-toolbar">
            <f7-link :class="{ 'disabled': loading || query.dateType === DateRange.All.type }" @click="shiftDateRange(query.minTime, query.maxTime, -1)">
                <f7-icon f7="arrow_left_square"></f7-icon>
            </f7-link>
            <f7-link :class="{ 'tabbar-text-with-ellipsis': true, 'disabled': loading }" popover-open=".date-popover-menu">
                <span :class="{ 'tabbar-item-changed': query.dateType !== DateRange.All.type }">{{ queryDateRangeName }}</span>
            </f7-link>
            <f7-link :class="{ 'disabled': loading || query.dateType === DateRange.All.type }" @click="shiftDateRange(query.minTime, query.maxTime, 1)">
                <f7-icon f7="arrow_right_square"></f7-icon>
            </f7-link>
            <f7-link class="tabbar-text-with-ellipsis" popover-open=".category-popover-menu" :class="{ 'disabled': query.type === 1 }">
                <span :class="{ 'tabbar-item-changed': query.categoryIds }">{{ queryCategoryName }}</span>
            </f7-link>
            <f7-link class="tabbar-text-with-ellipsis" popover-open=".account-popover-menu">
                <span :class="{ 'tabbar-item-changed': query.accountIds }">{{ queryAccountName }}</span>
            </f7-link>
            <f7-link popover-open=".more-popover-menu">
                <f7-icon f7="ellipsis_vertical" :class="{ 'tabbar-item-changed': query.type > 0 || query.amountFilter || query.tagIds }"></f7-icon>
            </f7-link>
        </f7-toolbar>

        <f7-block class="transaction-calendar-container margin-vertical" v-if="pageType === TransactionListPageType.Calendar.type">
            <vue-date-picker inline auto-apply model-type="yyyy-M-d"
                             month-name-format="long"
                             class="justify-content-center"
                             :config="{ noSwipe: true }"
                             :readonly="loading"
                             :disable-month-year-select="true"
                             :month-change-on-scroll="false"
                             :month-change-on-arrows="false"
                             :enable-time-picker="false"
                             :hide-offset-dates="true"
                             :min-date="transactionCalendarMinDate"
                             :max-date="transactionCalendarMaxDate"
                             :disabled-dates="noTransactionInMonthDay"
                             :prevent-min-max-navigation="true"
                             :clearable="false"
                             :dark="isDarkMode"
                             :week-start="firstDayOfWeek"
                             :day-names="dayNames"
                             v-model="currentCalendarDate">
                <template #day="{ day }">
                    <div class="transaction-calendar-daily-amounts display-flex flex-direction-column align-items-center justify-content-center">
                        <span>{{ day }}</span>
                        <small class="text-income" v-if="currentMonthTransactionData && currentMonthTransactionData.dailyTotalAmounts[day] && currentMonthTransactionData.dailyTotalAmounts[day].income">{{ getDisplayMonthTotalAmount(currentMonthTransactionData.dailyTotalAmounts[day].income, false, '', currentMonthTransactionData.dailyTotalAmounts[day].incompleteIncome) }}</small>
                        <small class="text-expense" v-if="currentMonthTransactionData && currentMonthTransactionData.dailyTotalAmounts[day] && currentMonthTransactionData.dailyTotalAmounts[day].expense">{{ getDisplayMonthTotalAmount(currentMonthTransactionData.dailyTotalAmounts[day].expense, false, '', currentMonthTransactionData.dailyTotalAmounts[day].incompleteExpense) }}</small>
                    </div>
                </template>
            </vue-date-picker>
        </f7-block>

        <div class="skeleton-text" v-if="loading">
            <f7-block class="combination-list-wrapper margin-vertical" :class="{ 'no-accordion-toggle': pageType !== TransactionListPageType.List.type }"
                      :key="blockIdx" v-for="blockIdx in (pageType === TransactionListPageType.List.type ? [ 1, 2 ] : [ 1 ])">
                <f7-accordion-item>
                    <f7-block-title v-if="pageType === TransactionListPageType.List.type">
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
                                          :key="itemIdx" v-for="itemIdx in (blockIdx === 1 ? [ 1, 2, 3, 4, 5, 6, 7 ] : [ 1, 2, 3 ])">
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
                                                    <div class="transaction-amount">
                                                        <span>0.00 USD</span>
                                                    </div>
                                                </div>
                                            </div>
                                            <div class="item-text">
                                                <div class="transaction-description">
                                                    <span>Transaction Description</span>
                                                </div>
                                            </div>
                                            <div class="item-footer">
                                                <div class="transaction-footer">
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
        </div>

        <f7-list strong inset dividers class="margin-vertical" v-if="!loading && noTransaction">
            <f7-list-item :title="tt('No transaction data')"></f7-list-item>
        </f7-list>

        <f7-block class="combination-list-wrapper margin-vertical" :class="{ 'no-accordion-toggle': pageType !== TransactionListPageType.List.type }"
                  :key="transactionMonthList.yearMonth" v-for="(transactionMonthList) in transactions">
            <f7-accordion-item :opened="transactionMonthList.opened"
                               @accordion:open="collapseTransactionMonthList(transactionMonthList, false)"
                               @accordion:close="collapseTransactionMonthList(transactionMonthList, true)">
                <f7-block-title v-if="pageType === TransactionListPageType.List.type">
                    <f7-accordion-toggle>
                        <f7-list strong inset dividers media-list
                                 class="transaction-amount-list combination-list-header"
                                 :class="transactionMonthList.opened ? 'combination-list-opened' : 'combination-list-closed'">
                            <f7-list-item>
                                <template #title>
                                    <small>
                                        <span>{{ getDisplayLongYearMonth(transactionMonthList) }}</span>
                                    </small>
                                    <small class="transaction-amount-statistics" v-if="showTotalAmountInTransactionListPage && transactionMonthList.totalAmount">
                                        <span class="text-income">
                                            {{ getDisplayMonthTotalAmount(transactionMonthList.totalAmount.income, defaultCurrency, '+', transactionMonthList.totalAmount.incompleteIncome) }}
                                        </span>
                                        <span class="text-expense">
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
                                      :link="transaction.type !== TransactionType.ModifyBalance ? `/transaction/detail?id=${transaction.id}&type=${transaction.type}` : null"
                                      :key="transaction.id"
                                      v-for="(transaction, idx) in transactionMonthList.items"
                        >
                            <template #media>
                                <div class="display-flex flex-direction-column transaction-date" :style="getTransactionDateStyle(transaction, idx > 0 ? transactionMonthList.items[idx - 1] : null)">
                                    <span class="transaction-day full-line flex-direction-column">
                                        {{ transaction.day }}
                                    </span>
                                    <span class="transaction-day-of-week full-line flex-direction-column" v-if="transaction.dayOfWeek">
                                        {{ getWeekdayShortName(transaction.dayOfWeek) }}
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
                                                    <span v-if="transaction.type === TransactionType.ModifyBalance">
                                                        {{ tt('Modify Balance') }}
                                                    </span>
                                                        <span v-else-if="transaction.type !== TransactionType.ModifyBalance && transaction.category">
                                                        {{ transaction.category.name }}
                                                    </span>
                                                        <span v-else-if="transaction.type !== TransactionType.ModifyBalance && !transaction.category">
                                                        {{ getTransactionTypeName(transaction.type, 'Transaction') }}
                                                    </span>
                                                </div>
                                            </div>
                                            <div class="item-after">
                                                <div class="transaction-amount" v-if="transaction.sourceAccount"
                                                     :class="{ 'text-expense': transaction.type === TransactionType.Expense, 'text-income': transaction.type === TransactionType.Income }">
                                                    <span>{{ getDisplayAmount(transaction) }}</span>
                                                </div>
                                            </div>
                                        </div>
                                        <div class="item-text">
                                            <div class="transaction-description" v-if="transaction.comment">
                                                <span>{{ transaction.comment }}</span>
                                            </div>
                                        </div>
                                        <div class="item-footer">
                                            <div class="transaction-tags" v-if="showTagInTransactionListPage && transaction.tagIds && transaction.tagIds.length">
                                                <f7-chip media-text-color="var(--f7-chip-text-color)" class="transaction-tag"
                                                         :text="allTransactionTags[tagId].name"
                                                         :key="tagId"
                                                         v-for="tagId in transaction.tagIds">
                                                    <template #media>
                                                        <f7-icon f7="number"></f7-icon>
                                                    </template>
                                                </f7-chip>
                                            </div>
                                            <div class="transaction-footer">
                                                <span>{{ getDisplayTime(transaction) }}</span>
                                                <span v-if="transaction.utcOffset !== currentTimezoneOffsetMinutes">{{ `(${getDisplayTimezone(transaction)})` }}</span>
                                                <span v-if="transaction.sourceAccount">·</span>
                                                <span v-if="transaction.sourceAccount">{{ transaction.sourceAccount.name }}</span>
                                                <f7-icon f7="arrow_right" class="transaction-account-arrow" v-if="transaction.sourceAccount && transaction.type === TransactionType.Transfer && transaction.destinationAccount && transaction.sourceAccount.id !== transaction.destinationAccount.id"></f7-icon>
                                                <span v-if="transaction.sourceAccount && transaction.type === TransactionType.Transfer && transaction.destinationAccount && transaction.sourceAccount.id !== transaction.destinationAccount.id">{{ transaction.destinationAccount.name }}</span>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </template>
                            <f7-swipeout-actions right>
                                <f7-swipeout-button color="primary" close
                                                    :text="tt('Duplicate')"
                                                    v-if="transaction.type !== TransactionType.ModifyBalance"
                                                    @click="duplicate(transaction)"></f7-swipeout-button>
                                <f7-swipeout-button color="orange" close
                                                    :text="tt('Edit')"
                                                    v-if="transaction.editable && transaction.type !== TransactionType.ModifyBalance"
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

        <f7-block class="text-align-center" :class="{ 'disabled': loadingMore }" v-show="!loading && hasMoreTransaction"
                  v-if="pageType === TransactionListPageType.List.type">
            <f7-link href="#" @click="loadMore(false)">{{ tt('Load More') }}</f7-link>
        </f7-block>

        <f7-popover class="date-popover-menu"
                    v-model:opened="showDatePopover"
                    @popover:open="onPopoverOpen">
            <f7-list dividers>
                <f7-list-item :title="dateRange.displayName"
                              :class="{ 'list-item-selected': query.dateType === dateRange.type }"
                              :key="dateRange.type"
                              v-for="dateRange in allDateRanges"
                              v-show="pageType === TransactionListPageType.List.type || dateRange.type === DateRange.ThisMonth.type || dateRange.type === DateRange.LastMonth.type || dateRange.type === DateRange.Custom.type"
                              @click="changeDateFilter(dateRange.type)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.dateType === dateRange.type"></f7-icon>
                    </template>
                    <template #footer>
                        <div v-if="((dateRange.isBillingCycle || dateRange.type === DateRange.Custom.type) && query.dateType === dateRange.type) && query.minTime && query.maxTime">
                            <span>{{ queryMinTime }}</span>
                            <span>&nbsp;-&nbsp;</span>
                            <br/>
                            <span>{{ queryMaxTime }}</span>
                        </div>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <date-range-selection-sheet :title="tt('Custom Date Range')"
                                    :min-time="customMinDatetime"
                                    :max-time="customMaxDatetime"
                                    v-model:show="showCustomDateRangeSheet"
                                    @dateRange:change="changeCustomDateFilter">
        </date-range-selection-sheet>

        <month-selection-sheet :title="tt('Custom Date Range')"
                               :model-value="queryMonth"
                               v-model:show="showCustomMonthSheet"
                               @update:modelValue="changeCustomMonthDateFilter">
        </month-selection-sheet>

        <f7-popover class="category-popover-menu"
                    v-model:opened="showCategoryPopover"
                    @popover:open="onPopoverOpen">
            <f7-list dividers accordion-list>
                <f7-list-item :class="{ 'list-item-selected': !query.categoryIds }" :title="tt('All')" @click="changeCategoryFilter('')">
                    <template #media>
                        <f7-icon f7="rectangle_grid_2x2"></f7-icon>
                    </template>
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="!query.categoryIds"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item :class="{ 'list-item-selected': query.categoryIds && queryAllFilterCategoryIdsCount > 1 }"
                              :title="tt('Multiple Categories')"
                              @click="filterMultipleCategories()"
                              v-if="allAvailableCategoriesCount > 0">
                    <template #media>
                        <f7-icon f7="rectangle_on_rectangle"></f7-icon>
                    </template>
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.categoryIds && queryAllFilterCategoryIdsCount > 1"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
            <f7-list dividers accordion-list
                     class="no-margin-vertical"
                     :key="categoryType"
                     v-for="(categories, categoryType) in allPrimaryCategories"
                     v-show="categories && categories.length"
            >
                <f7-list-item divider :title="getTransactionTypeName(categoryTypeToTransactionType(parseInt(categoryType)), 'Type')"></f7-list-item>
                <f7-list-item accordion-item
                              :title="category.name"
                              :class="getCategoryListItemCheckedClass(category, queryAllFilterCategoryIds)"
                              :key="category.id"
                              v-for="category in categories"
                              v-show="!category.hidden || query.categoryIds === category.id || (allCategories[query.categoryIds] && allCategories[query.categoryIds].parentId === category.id)"
                >
                    <template #media>
                        <ItemIcon icon-type="category" :icon-id="category.icon" :color="category.color"></ItemIcon>
                    </template>
                    <f7-accordion-content>
                        <f7-list dividers class="padding-left">
                            <f7-list-item :class="{ 'list-item-selected': query.categoryIds === category.id, 'item-in-multiple-selection': queryAllFilterCategoryIdsCount > 1 && queryAllFilterCategoryIds[category.id] }"
                                          :title="tt('All')" @click="changeCategoryFilter(category.id)">
                                <template #media>
                                    <f7-icon f7="rectangle_grid_2x2"></f7-icon>
                                </template>
                                <template #after>
                                    <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.categoryIds === category.id"></f7-icon>
                                </template>
                            </f7-list-item>
                            <f7-list-item :title="subCategory.name"
                                          :class="{ 'list-item-selected': query.categoryIds === subCategory.id, 'item-in-multiple-selection': queryAllFilterCategoryIdsCount > 1 && queryAllFilterCategoryIds[subCategory.id] }"
                                          :key="subCategory.id"
                                          v-for="subCategory in category.subCategories"
                                          v-show="!subCategory.hidden || query.categoryIds === subCategory.id"
                                          @click="changeCategoryFilter(subCategory.id)"
                            >
                                <template #media>
                                    <ItemIcon icon-type="category" :icon-id="subCategory.icon" :color="subCategory.color"></ItemIcon>
                                </template>
                                <template #after>
                                    <f7-icon class="list-item-checked-icon"
                                             f7="checkmark_alt"
                                             v-if="query.categoryIds === subCategory.id">
                                    </f7-icon>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-content>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-popover class="account-popover-menu"
                    v-model:opened="showAccountPopover"
                    @popover:open="onPopoverOpen">
            <f7-list dividers>
                <f7-list-item :class="{ 'list-item-selected': !query.accountIds }" :title="tt('All')" @click="changeAccountFilter('')">
                    <template #media>
                        <f7-icon f7="rectangle_grid_2x2"></f7-icon>
                    </template>
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="!query.accountIds"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item :class="{ 'list-item-selected': query.accountIds && queryAllFilterAccountIdsCount > 1 }"
                              :title="tt('Multiple Accounts')"
                              @click="filterMultipleAccounts()"
                              v-if="allAvailableAccountsCount > 0">
                    <template #media>
                        <f7-icon f7="rectangle_on_rectangle"></f7-icon>
                    </template>
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.accountIds && queryAllFilterAccountIdsCount > 1"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item :title="account.name"
                              :class="{ 'list-item-selected': query.accountIds === account.id, 'item-in-multiple-selection': queryAllFilterAccountIdsCount > 1 && queryAllFilterAccountIds[account.id] }"
                              :key="account.id"
                              v-for="account in allAccounts"
                              v-show="(!account.hidden && (!allAccountsMap[account.parentId] || !allAccountsMap[account.parentId].hidden)) || query.accountIds === account.id"
                              @click="changeAccountFilter(account.id)"
                >
                    <template #media>
                        <ItemIcon icon-type="account" :icon-id="account.icon" :color="account.color"></ItemIcon>
                    </template>
                    <template #after>
                        <f7-icon class="list-item-checked-icon"
                                 f7="checkmark_alt"
                                 v-if="query.accountIds === account.id">
                        </f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-popover class="more-popover-menu"
                    v-model:opened="showMorePopover">
            <f7-list dividers>
                <f7-list-item group-title :title="tt('Type')" />
                <f7-list-item :class="{ 'list-item-selected': query.type === 0 }" :title="tt('All')" @click="changeTypeFilter(0)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.type === 0"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item :class="{ 'list-item-selected': query.type === 1 }" :title="tt('Modify Balance')" @click="changeTypeFilter(1)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.type === 1"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item :class="{ 'list-item-selected': query.type === 2 }" :title="tt('Income')" @click="changeTypeFilter(2)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.type === 2"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item :class="{ 'list-item-selected': query.type === 3 }" :title="tt('Expense')" @click="changeTypeFilter(3)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.type === 3"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item :class="{ 'list-item-selected': query.type === 4 }" :title="tt('Transfer')" @click="changeTypeFilter(4)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.type === 4"></f7-icon>
                    </template>
                </f7-list-item>

                <f7-list-item group-title :title="tt('Amount')" />
                <f7-list-item :class="{ 'list-item-selected': !query.amountFilter }" :title="tt('All')" @click="changeAmountFilter('')">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="!query.amountFilter"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item :key="filterType.type"
                              :class="{ 'list-item-selected': query.amountFilter && query.amountFilter.startsWith(`${filterType.type}:`) }"
                              :title="tt(filterType.name)"
                              v-for="filterType in AmountFilterType.values()"
                              @click="changeAmountFilter(filterType.type)">
                    <template #after>
                        <span class="margin-right-half" v-if="query.amountFilter && query.amountFilter.startsWith(`${filterType.type}:`)">{{ queryAmount }}</span>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.amountFilter && query.amountFilter.startsWith(`${filterType.type}:`)"></f7-icon>
                    </template>
                </f7-list-item>

                <f7-list-item group-title :title="tt('Tags')" />
                <f7-list-item :class="{ 'list-item-selected': !query.tagIds }" :title="tt('All')" @click="changeTagFilter('')">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="!query.tagIds"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item :class="{ 'list-item-selected': query.tagIds === 'none' }" :title="tt('Without Tags')" @click="changeTagFilter('none')">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.tagIds === 'none'"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item :class="{ 'list-item-selected': query.tagIds && queryAllFilterTagIdsCount > 1 }"
                              :title="tt('Multiple Tags')" @click="filterMultipleTags()" v-if="allAvailableTagsCount > 0">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.tagIds && queryAllFilterTagIdsCount > 1"></f7-icon>
                    </template>
                </f7-list-item>

                <template v-if="query.tagIds && query.tagIds !== 'none'">
                    <f7-list-item :title="filterType.displayName"
                                  :key="filterType.type"
                                  v-for="filterType in allTransactionTagFilterTypes"
                                  @click="changeTagFilterType(filterType.type)"
                    >
                        <template #after>
                            <f7-icon class="list-item-checked-icon"
                                     f7="checkmark_alt"
                                     v-if="query.tagFilterType === filterType.type">
                            </f7-icon>
                        </template>
                    </f7-list-item>
                </template>

                <f7-list-item :title="transactionTag.name"
                              :class="{ 'list-item-selected': query.tagIds === transactionTag.id, 'item-in-multiple-selection': queryAllFilterTagIdsCount > 1 && queryAllFilterTagIds[transactionTag.id] }"
                              :key="transactionTag.id"
                              v-for="transactionTag in allTransactionTags"
                              v-show="!transactionTag.hidden || query.tagIds === transactionTag.id"
                              @click="changeTagFilter(transactionTag.id)"
                >
                    <template #before-title>
                        <f7-icon class="transaction-tag-name transaction-tag-icon" f7="number"></f7-icon>
                    </template>
                    <template #after>
                        <f7-icon class="list-item-checked-icon"
                                 f7="checkmark_alt"
                                 v-if="query.tagIds === transactionTag.id">
                        </f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <f7-actions close-by-outside-click close-on-escape :opened="showDeleteActionSheet" @actions:closed="showDeleteActionSheet = false">
            <f7-actions-group>
                <f7-actions-label>{{ tt('Are you sure you want to delete this transaction?') }}</f7-actions-label>
                <f7-actions-button color="red" @click="remove(transactionToDelete, true)">{{ tt('Delete') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { type Framework7Dom, useI18nUIComponents, showLoading, hideLoading, onSwipeoutDeleted, scrollToSelectedItem } from '@/lib/ui/mobile.ts';
import { TransactionListPageType, useTransactionListPageBase } from '@/views/base/transactions/TransactionListPageBase.ts';

import { useEnvironmentsStore } from '@/stores/environment.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { type TransactionMonthList, useTransactionsStore } from '@/stores/transaction.ts';

import type { TypeAndDisplayName } from '@/core/base.ts';
import {
    type TimeRangeAndDateType,
    DateRangeScene,
    DateRange,
} from '@/core/datetime.ts';
import { AmountFilterType } from '@/core/numeral.ts';
import { TransactionType } from '@/core/transaction.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';
import type { Transaction } from '@/models/transaction.ts';

import {
    arrangeArrayWithNewStartIndex
} from '@/lib/common.ts';
import {
    getCurrentUnixTime,
    parseDateFromUnixTime,
    getBrowserTimezoneOffsetMinutes,
    getActualUnixTimeForStore,
    getYear,
    getMonth,
    getSpecifiedDayFirstUnixTime,
    getYearMonthFirstUnixTime,
    getYearMonthLastUnixTime,
    getShiftedDateRangeAndDateType,
    getShiftedDateRangeAndDateTypeForBillingCycle,
    getDateTypeByDateRange,
    getDateTypeByBillingCycleDateRange,
    getDateRangeByDateType,
    getDateRangeByBillingCycleDateType,
    getFullMonthDateRange,
    getMonthFirstDayOrCurrentDayShortDate
} from '@/lib/datetime.ts';
import {
    categoryTypeToTransactionType,
    transactionTypeToCategoryType
} from '@/lib/category.ts';

const props = defineProps<{
    f7route: Router.Route;
    f7router: Router.Router;
}>();

const {
    tt,
    getAllShortWeekdayNames,
    getAllTransactionTagFilterTypes,
    getWeekdayShortName
} = useI18n();

const { showAlert, showToast, routeBackOnError } = useI18nUIComponents();

const {
    pageType,
    loading,
    customMinDatetime,
    customMaxDatetime,
    currentCalendarDate,
    currentTimezoneOffsetMinutes,
    firstDayOfWeek,
    fiscalYearStart,
    defaultCurrency,
    showTotalAmountInTransactionListPage,
    showTagInTransactionListPage,
    allDateRanges,
    allAccounts,
    allAccountsMap,
    allAvailableAccountsCount,
    allCategories,
    allPrimaryCategories,
    allAvailableCategoriesCount,
    allTransactionTags,
    allAvailableTagsCount,
    displayPageTypeName,
    query,
    queryDateRangeName,
    queryMinTime,
    queryMaxTime,
    queryMonthlyData,
    queryMonth,
    queryAllFilterCategoryIds,
    queryAllFilterAccountIds,
    queryAllFilterTagIds,
    queryAllFilterCategoryIdsCount,
    queryAllFilterAccountIdsCount,
    queryAllFilterTagIdsCount,
    queryAccountName,
    queryCategoryName,
    queryAmount,
    transactionCalendarMinDate,
    transactionCalendarMaxDate,
    currentMonthTransactionData,
    noTransactionInMonthDay,
    canAddTransaction,
    getDisplayTime,
    getDisplayLongYearMonth,
    getDisplayTimezone,
    getDisplayAmount,
    getDisplayMonthTotalAmount,
    getTransactionTypeName,
} = useTransactionListPageBase();

const environmentsStore = useEnvironmentsStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();
const transactionsStore = useTransactionsStore();

const isDarkMode = computed<boolean>(() => environmentsStore.framework7DarkMode || false);
const dayNames = computed<string[]>(() => arrangeArrayWithNewStartIndex(getAllShortWeekdayNames(), firstDayOfWeek.value));

const loadingError = ref<unknown | null>(null);
const loadingMore = ref<boolean>(false);
const transactionToDelete = ref<Transaction | null>(null);
const showTransactionListPageTypePopover = ref<boolean>(false);
const showDatePopover = ref<boolean>(false);
const showCategoryPopover = ref<boolean>(false);
const showAccountPopover = ref<boolean>(false);
const showMorePopover = ref<boolean>(false);
const showCustomDateRangeSheet = ref<boolean>(false);
const showCustomMonthSheet = ref<boolean>(false);
const showDeleteActionSheet = ref<boolean>(false);

const allTransactionTagFilterTypes = computed<TypeAndDisplayName[]>(() => getAllTransactionTagFilterTypes());

const transactions = computed<TransactionMonthList[]>(() => {
    if (loading.value) {
        return [];
    }

    if (pageType.value === TransactionListPageType.List.type) {
        return transactionsStore.transactions;
    } else if (pageType.value === TransactionListPageType.Calendar.type) {
        if (queryMonthlyData.value) {
            const transactionData = currentMonthTransactionData.value;

            if (!transactionData || !transactionData.items) {
                return [];
            }

            const transactions :Transaction[] = [];

            for (let i = 0; i < transactionData.items.length; i++) {
                const transaction = transactionData.items[i];

                if (transaction.date === currentCalendarDate.value) {
                    transactions.push(transaction);
                }
            }

            const dailyTransactionList: TransactionMonthList = {
                year: currentMonthTransactionData.value.year,
                month: currentMonthTransactionData.value.month,
                yearMonth: currentMonthTransactionData.value.yearMonth,
                opened: true,
                items: transactions,
                totalAmount: {
                    income: 0,
                    expense: 0,
                    incompleteIncome: false,
                    incompleteExpense: false
                },
                dailyTotalAmounts: {}
            };

            return [dailyTransactionList];
        } else {
            return [];
        }
    } else {
        return [];
    }
});

const noTransaction = computed<boolean>(() => {
    if (pageType.value === TransactionListPageType.List.type) {
        return transactionsStore.noTransaction;
    } else if (pageType.value === TransactionListPageType.Calendar.type) {
        return !transactions.value || !transactions.value.length || !transactions.value[0].items || !transactions.value[0].items.length;
    } else {
        return true;
    }
});

const hasMoreTransaction = computed<boolean>(() => transactionsStore.hasMoreTransaction);

function getTransactionDomId(transaction: Transaction): string {
    return 'transaction_' + transaction.id;
}

function getTransactionDateStyle(transaction: Transaction, previousTransaction: Transaction | null): Record<string, string> {
    if (!previousTransaction || transaction.day !== previousTransaction.day) {
        return {};
    }

    return {
        color: 'transparent'
    };
}

function getCategoryListItemCheckedClass(category: TransactionCategory, queryCategoryIds: Record<string, boolean>): Record<string, boolean> {
    if (queryCategoryIds && queryCategoryIds[category.id]) {
        return {
            'list-item-checked': true
        };
    }

    if (category.subCategories) {
        for (let i = 0; i < category.subCategories.length; i++) {
            if (queryCategoryIds && queryCategoryIds[category.subCategories[i].id]) {
                return {
                    'list-item-checked': true
                };
            }
        }
    }

    return {};
}

function init(): void {
    const initQuery = props.f7route.query;

    let dateRange: TimeRangeAndDateType | null = getDateRangeByDateType(initQuery['dateType'] ? parseInt(initQuery['dateType']) : undefined, firstDayOfWeek.value, fiscalYearStart.value);

    if (!dateRange && initQuery['dateType'] && initQuery['maxTime'] && initQuery['minTime'] &&
        (DateRange.isBillingCycle(parseInt(initQuery['dateType'])) || initQuery['dateType'] === DateRange.Custom.type.toString()) &&
        parseInt(initQuery['maxTime']) > 0 && parseInt(initQuery['minTime']) > 0) {
        dateRange = {
            dateType: parseInt(initQuery['dateType']),
            maxTime: parseInt(initQuery['maxTime']),
            minTime: parseInt(initQuery['minTime'])
        };
    }

    transactionsStore.initTransactionListFilter({
        dateType: dateRange ? dateRange.dateType : undefined,
        maxTime: dateRange ? dateRange.maxTime : undefined,
        minTime: dateRange ? dateRange.minTime : undefined,
        type: initQuery['type'] && parseInt(initQuery['type']) > 0 ? parseInt(initQuery['type']) : undefined,
        categoryIds: initQuery['categoryIds'],
        accountIds: initQuery['accountIds'],
        tagIds: initQuery['tagIds'],
        tagFilterType: initQuery['tagFilterType'] && parseInt(initQuery['tagFilterType']) >= 0 ? parseInt(initQuery['tagFilterType']) : undefined
    });

    reload();
}

function reload(done?: () => void): void {
    const force = !!done;

    if (!done) {
        loading.value = true;
    }

    Promise.all([
        accountsStore.loadAllAccounts({ force: false }),
        transactionCategoriesStore.loadAllCategories({ force: false }),
        transactionTagsStore.loadAllTags({ force: false })
    ]).then(() => {
        if (queryMonthlyData.value) {
            const currentMonthMinDate = parseDateFromUnixTime(query.value.minTime);
            const currentYear = getYear(currentMonthMinDate);
            const currentMonth = getMonth(currentMonthMinDate);

            return transactionsStore.loadMonthlyAllTransactions({
                year: currentYear,
                month: currentMonth,
                autoExpand: true,
                defaultCurrency: defaultCurrency.value
            });
        } else {
            return transactionsStore.loadTransactions({
                reload: true,
                autoExpand: true,
                defaultCurrency: defaultCurrency.value
            });
        }
    }).then(() => {
        done?.();

        if (force) {
            showToast('Data has been updated');
        }

        loading.value = false;
    }).catch(error => {
        if (error.processed || done) {
            loading.value = false;
        }

        done?.();

        if (!error.processed) {
            if (!done) {
                loadingError.value = error;
            }

            showToast(error.message || error);
        }
    });
}

function loadMore(autoExpand: boolean): void {
    if (!hasMoreTransaction.value) {
        return;
    }

    if (loadingMore.value || loading.value) {
        return;
    }

    loadingMore.value = true;

    transactionsStore.loadTransactions({
        reload: false,
        autoExpand: autoExpand,
        defaultCurrency: defaultCurrency.value
    }).then(() => {
        loadingMore.value = false;
    }).catch(error => {
        loadingMore.value = false;

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function changePageType(type: number): void {
    pageType.value = type;
    currentCalendarDate.value = getMonthFirstDayOrCurrentDayShortDate(query.value.minTime);

    if (pageType.value === TransactionListPageType.Calendar.type) {
        const dateRange = getFullMonthDateRange(query.value.minTime, query.value.maxTime, firstDayOfWeek.value);

        if (dateRange) {
            const changed = transactionsStore.updateTransactionListFilter({
                dateType: dateRange.dateType,
                maxTime: dateRange.maxTime,
                minTime: dateRange.minTime
            });

            if (changed) {
                currentCalendarDate.value = getMonthFirstDayOrCurrentDayShortDate(query.value.minTime);
                reload();
            }
        }
    }
}

function changeDateFilter(dateType: number): void {
    if (dateType === DateRange.Custom.type) { // Custom
        if (!query.value.minTime || !query.value.maxTime) {
            customMaxDatetime.value = getActualUnixTimeForStore(getCurrentUnixTime(), currentTimezoneOffsetMinutes.value, getBrowserTimezoneOffsetMinutes());
            customMinDatetime.value = getSpecifiedDayFirstUnixTime(customMaxDatetime.value);
        } else {
            customMaxDatetime.value = query.value.maxTime;
            customMinDatetime.value = query.value.minTime;
        }

        if (pageType.value === TransactionListPageType.Calendar.type) {
            showCustomMonthSheet.value = true;
        } else {
            showCustomDateRangeSheet.value = true;
        }

        showDatePopover.value = false;
        return;
    } else if (query.value.dateType === dateType) {
        return;
    }

    let dateRange: TimeRangeAndDateType | null = null;

    if (DateRange.isBillingCycle(dateType)) {
        dateRange = getDateRangeByBillingCycleDateType(dateType, firstDayOfWeek.value, fiscalYearStart.value, accountsStore.getAccountStatementDate(query.value.accountIds));
    } else {
        dateRange = getDateRangeByDateType(dateType, firstDayOfWeek.value, fiscalYearStart.value);
    }

    if (!dateRange) {
        return;
    }

    if (pageType.value === TransactionListPageType.Calendar.type) {
        currentCalendarDate.value = getMonthFirstDayOrCurrentDayShortDate(dateRange.minTime);
        const fullMonthDateRange = getFullMonthDateRange(dateRange.minTime, dateRange.maxTime, firstDayOfWeek.value);

        if (fullMonthDateRange) {
            dateRange = fullMonthDateRange;
            currentCalendarDate.value = getMonthFirstDayOrCurrentDayShortDate(dateRange.minTime);
        }
    }

    const changed = transactionsStore.updateTransactionListFilter({
        dateType: dateRange.dateType,
        maxTime: dateRange.maxTime,
        minTime: dateRange.minTime
    });

    showDatePopover.value = false;

    if (changed) {
        reload();
    }
}

function changeCustomDateFilter(minTime: number, maxTime: number): void {
    if (!minTime || !maxTime) {
        return;
    }

    let dateType: number | null = getDateTypeByBillingCycleDateRange(minTime, maxTime, firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.Normal, accountsStore.getAccountStatementDate(query.value.accountIds));

    if (!dateType) {
        dateType = getDateTypeByDateRange(minTime, maxTime, firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.Normal);
    }

    if (pageType.value === TransactionListPageType.Calendar.type) {
        currentCalendarDate.value = getMonthFirstDayOrCurrentDayShortDate(minTime);
        const dateRange = getFullMonthDateRange(minTime, maxTime, firstDayOfWeek.value);

        if (dateRange) {
            minTime = dateRange.minTime;
            maxTime = dateRange.maxTime;
            dateType = dateRange.dateType;
            currentCalendarDate.value = getMonthFirstDayOrCurrentDayShortDate(minTime);
        }
    }

    const changed = transactionsStore.updateTransactionListFilter({
        dateType: dateType,
        maxTime: maxTime,
        minTime: minTime
    });

    showCustomDateRangeSheet.value = false;

    if (changed) {
        reload();
    }
}

function changeCustomMonthDateFilter(yearMonth: string): void {
    if (!yearMonth) {
        return;
    }

    const minTime = getYearMonthFirstUnixTime(yearMonth);
    const maxTime = getYearMonthLastUnixTime(yearMonth);
    const dateType = getDateTypeByDateRange(minTime, maxTime, firstDayOfWeek.value, DateRangeScene.Normal);

    if (pageType.value === TransactionListPageType.Calendar.type) {
        currentCalendarDate.value = getMonthFirstDayOrCurrentDayShortDate(minTime);
    }

    const changed = transactionsStore.updateTransactionListFilter({
        dateType: dateType,
        maxTime: maxTime,
        minTime: minTime
    });

    showCustomMonthSheet.value = false;

    if (changed) {
        reload();
    }
}

function shiftDateRange(minTime: number, maxTime: number, scale: number): void {
    if (query.value.dateType === DateRange.All.type) {
        return;
    }

    let newDateRange: TimeRangeAndDateType | null = null;

    if (DateRange.isBillingCycle(query.value.dateType) || query.value.dateType === DateRange.Custom.type) {
        newDateRange = getShiftedDateRangeAndDateTypeForBillingCycle(minTime, maxTime, scale, firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.Normal, accountsStore.getAccountStatementDate(query.value.accountIds));
    }

    if (!newDateRange) {
        newDateRange = getShiftedDateRangeAndDateType(minTime, maxTime, scale, firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.Normal);
    }

    if (pageType.value === TransactionListPageType.Calendar.type) {
        currentCalendarDate.value = getMonthFirstDayOrCurrentDayShortDate(newDateRange.minTime);
        const fullMonthDateRange = getFullMonthDateRange(newDateRange.minTime, newDateRange.maxTime, firstDayOfWeek.value);

        if (fullMonthDateRange) {
            newDateRange = fullMonthDateRange;
            currentCalendarDate.value = getMonthFirstDayOrCurrentDayShortDate(newDateRange.minTime);
        }
    }

    const changed = transactionsStore.updateTransactionListFilter({
        dateType: newDateRange.dateType,
        maxTime: newDateRange.maxTime,
        minTime: newDateRange.minTime
    });

    if (changed) {
        reload();
    }
}

function changeTypeFilter(type: number): void {
    if (query.value.type === type) {
        return;
    }

    let newCategoryFilter = undefined;

    if (type && query.value.categoryIds) {
        newCategoryFilter = '';

        for (const categoryId in queryAllFilterCategoryIds.value) {
            if (!Object.prototype.hasOwnProperty.call(queryAllFilterCategoryIds.value, categoryId)) {
                continue;
            }

            const category = allCategories.value[categoryId];

            if (category && category.type === transactionTypeToCategoryType(type)) {
                if (newCategoryFilter.length > 0) {
                    newCategoryFilter += ',';
                }

                newCategoryFilter += categoryId;
            }
        }
    }

    const changed = transactionsStore.updateTransactionListFilter({
        type: type,
        categoryIds: newCategoryFilter
    });

    showMorePopover.value = false;

    if (changed) {
        reload();
    }
}

function changeCategoryFilter(categoryIds: string): void {
    if (query.value.categoryIds === categoryIds) {
        return;
    }

    const changed = transactionsStore.updateTransactionListFilter({
        categoryIds: categoryIds
    });

    showCategoryPopover.value = false;

    if (changed) {
        reload();
    }
}

function filterMultipleCategories(): void {
    let navigateUrl = '/settings/filter/category?type=transactionListCurrent';

    if (TransactionType.Income <= query.value.type && query.value.type <= TransactionType.Transfer) {
        navigateUrl += '&allowCategoryTypes=' + transactionTypeToCategoryType(query.value.type);
    }

    props.f7router.navigate(navigateUrl);
}

function changeAccountFilter(accountIds: string): void {
    if (query.value.accountIds === accountIds) {
        return;
    }

    const changed = transactionsStore.updateTransactionListFilter({
        accountIds: accountIds
    });

    showAccountPopover.value = false;

    if (changed) {
        reload();
    }
}

function filterMultipleAccounts(): void {
    props.f7router.navigate('/settings/filter/account?type=transactionListCurrent');
}

function changeTagFilter(tagIds: string): void {
    if (query.value.tagIds === tagIds) {
        return;
    }

    const changed = transactionsStore.updateTransactionListFilter({
        tagIds: tagIds
    });

    showMorePopover.value = false;

    if (changed) {
        reload();
    }
}

function filterMultipleTags(): void {
    props.f7router.navigate('/settings/filter/tag?type=transactionListCurrent');
}

function changeTagFilterType(filterType: number): void {
    if (query.value.tagFilterType === filterType) {
        return;
    }

    const changed = transactionsStore.updateTransactionListFilter({
        tagFilterType: filterType
    });

    showMorePopover.value = false;

    if (changed) {
        reload();
    }
}

function changeKeywordFilter(keyword: string): void {
    if (query.value.keyword === keyword) {
        return;
    }

    const changed = transactionsStore.updateTransactionListFilter({
        keyword: keyword
    });

    if (changed) {
        reload();
    }
}

function changeAmountFilter(filterType: string): void {
    if (query.value.amountFilter === filterType) {
        return;
    }

    if (filterType) {
        showMorePopover.value = false;
        props.f7router.navigate(`/transaction/filter/amount?type=${filterType}&value=${query.value.amountFilter}`);
        return;
    }

    const changed = transactionsStore.updateTransactionListFilter({
        amountFilter: filterType
    });

    showMorePopover.value = false;

    if (changed) {
        reload();
    }
}

function duplicate(transaction: Transaction): void {
    props.f7router.navigate(`/transaction/add?id=${transaction.id}&type=${transaction.type}`);
}

function edit(transaction: Transaction): void {
    props.f7router.navigate(`/transaction/edit?id=${transaction.id}&type=${transaction.type}`);
}

function remove(transaction: Transaction | null, confirm: boolean): void {
    if (!transaction) {
        showAlert('An error occurred');
        return;
    }

    if (!confirm) {
        transactionToDelete.value = transaction;
        showDeleteActionSheet.value = true;
        return;
    }

    showDeleteActionSheet.value = false;
    transactionToDelete.value = null;
    showLoading();

    transactionsStore.deleteTransaction({
        transaction: transaction,
        defaultCurrency: defaultCurrency.value,
        beforeResolve: (done) => {
            onSwipeoutDeleted(getTransactionDomId(transaction), done);
        }
    }).then(() => {
        hideLoading();
    }).catch(error => {
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function collapseTransactionMonthList(month: TransactionMonthList, collapse: boolean): void {
    transactionsStore.collapseMonthInTransactionList({
        month: month,
        collapse: collapse
    });
}

function onPopoverOpen(event: { $el: Framework7Dom }): void {
    scrollToSelectedItem(event.$el, '.popover-inner', 'li.list-item-selected');
}

function onPageAfterIn(): void {
    if (transactionsStore.transactionListStateInvalid && !loading.value) {
        reload();
    }

    routeBackOnError(props.f7router, loadingError);
}

init();
</script>

<style>
.transaction-list-toolbar .toolbar-inner {
    padding-right: 8px;
}

.list.transaction-amount-list .transaction-amount-statistics {
    overflow: hidden;
    text-overflow: ellipsis;
}

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
    overflow: hidden;
}

.list.transaction-info-list li.transaction-info .actual-item-inner .item-after {
    max-width: 100%;
}

.list.transaction-info-list li.transaction-info .transaction-date {
    width: var(--ebk-transaction-date-width);
    margin-right: 6px;
}

.list.transaction-info-list li.transaction-info .transaction-day {
    opacity: 0.6;
    font-size: var(--ebk-transaction-day-font-size);
    font-weight: bold;
    text-align: left;
}

.list.transaction-info-list li.transaction-info .transaction-day-of-week {
    opacity: 0.6;
    font-size: var(--ebk-transaction-day-of-week-font-size);
}

.list.transaction-info-list li.transaction-info .transaction-description {
    font-size: var(--ebk-large-footer-font-size);
    line-height: 20px;
    padding-top: 2px;
    padding-bottom: 2px;
}

.list.transaction-info-list li.transaction-info .chip.transaction-tag {
    --f7-chip-media-size: var(--ebk-transaction-tag-chip-media-size);
    --f7-chip-media-font-size: var(--ebk-transaction-tag-chip-font-size);
    --f7-chip-font-size: var(--ebk-transaction-tag-chip-font-size);
    --f7-chip-height: var(--ebk-transaction-tag-chip-height);
    --f7-chip-text-color: var(--f7-list-item-footer-text-color);
    --f7-chip-bg-color: var(--ebk-transaction-tag-chip-bg-color);
    margin-right: 4px;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
}

.list.transaction-info-list li.transaction-info .chip.transaction-tag .chip-media+.chip-label {
    margin-left: 0;
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

.list.transaction-info-list li.transaction-info .transaction-footer .transaction-account-arrow {
    font-size: var(--ebk-transaction-account-arrow-font-size);
    margin-right: 4px;
    margin-top: var(--ebk-transaction-account-arrow-margin-top);
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

.more-popover-menu .transaction-tag-name {
    padding-right: 4px;
    font-size: var(--f7-list-item-title-font-size);
}

.date-popover-menu .popover-inner,
.category-popover-menu .popover-inner,
.account-popover-menu .popover-inner,
.more-popover-menu .popover-inner{
    max-height: 400px;
    overflow-y: auto;
}

.transaction-calendar-container .dp__main .dp__menu {
    --dp-border-radius: var(--f7-list-inset-border-radius);
    --dp-menu-padding: 4px 6px;
    --dp-menu-border-color: transparent;
}

.transaction-calendar-container .dp__main .dp__menu.dp__theme_dark {
    --dp-background-color: var(--f7-list-strong-bg-color);
}

.transaction-calendar-container .dp__main .dp__calendar .dp__calendar_row {
    --dp-cell-size: var(--ebk-transaction-calendar-daily-amounts-height);
    --dp-cell-padding: 1px;
    --dp-primary-text-color: var(--f7-theme-color);
}

.transaction-calendar-container .dp__main .dp__calendar .dp__calendar_row > .dp__calendar_item .transaction-calendar-daily-amounts {
    width: 100%;
    height: 100%;
    background-color: var(--f7-list-group-title-bg-color);
    border-radius: 6px;
}

.transaction-calendar-container .dp__main .dp__calendar .dp__calendar_row > .dp__calendar_item > .dp__active_date {
    background-color: transparent;
}

.transaction-calendar-container .dp__main .dp__calendar .dp__calendar_row > .dp__calendar_item > .dp__today {
    border: inherit;
}

.transaction-calendar-container .dp__main .dp__calendar .dp__calendar_row > .dp__calendar_item > .dp__date_hover_end:hover,
.transaction-calendar-container .dp__main .dp__calendar .dp__calendar_row > .dp__calendar_item > .dp__date_hover_start:hover,
.transaction-calendar-container .dp__main .dp__calendar .dp__calendar_row > .dp__calendar_item > .dp__date_hover:hover {
    background-color: transparent;
}

.transaction-calendar-container .dp__main .dp__calendar .dp__calendar_row > .dp__calendar_item > .dp__active_date .transaction-calendar-daily-amounts {
    background-color: rgba(var(--ebk-primary-color), 0.16);
}

.transaction-calendar-container .dp__main .dp__calendar .dp__calendar_row > .dp__calendar_item > .dp__today .transaction-calendar-daily-amounts {
    border: 1px solid var(--dp-primary-color);
}

.transaction-calendar-container .dp__main .dp__calendar .dp__calendar_row > .dp__calendar_item > .dp__date_hover_end:hover .transaction-calendar-daily-amounts,
.transaction-calendar-container .dp__main .dp__calendar .dp__calendar_row > .dp__calendar_item > .dp__date_hover_start:hover .transaction-calendar-daily-amounts,
.transaction-calendar-container .dp__main .dp__calendar .dp__calendar_row > .dp__calendar_item > .dp__date_hover:hover .transaction-calendar-daily-amounts {
    background: var(--dp-hover-color);
}

.transaction-calendar-container .dp__main .dp__calendar .dp__calendar_row > .dp__calendar_item {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.transaction-calendar-container .dp__main .dp__calendar .dp__calendar_row > .dp__calendar_item .transaction-calendar-daily-amounts > small {
    display: block;
    width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: var(--ebk-transaction-calendar-amount-font-size);
}
</style>
