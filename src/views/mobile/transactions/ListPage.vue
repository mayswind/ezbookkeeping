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
                <f7-link icon-f7="plus" :class="{ 'disabled': !canAddTransaction }" :href="`/transaction/add?type=${query.type}&categoryId=${queryAllFilterCategoryIdsCount === 1 ? query.categoryIds : ''}&accountId=${queryAllFilterAccountIdsCount === 1 ? query.accountIds : ''}&tagIds=${query.tagIds || ''}`"></f7-link>
            </f7-nav-right>

            <f7-subnavbar :inner="false">
                <f7-searchbar
                    custom-searchs
                    :value="query.keyword"
                    :placeholder="$t('Search transaction description')"
                    :disable-button-text="$t('Cancel')"
                    @change="changeKeywordFilter($event.target.value)"
                ></f7-searchbar>
            </f7-subnavbar>
        </f7-navbar>

        <f7-toolbar tabbar bottom class="toolbar-item-auto-size transaction-list-toolbar">
            <f7-link :class="{ 'disabled': loading || query.dateType === allDateRanges.All.type }" @click="shiftDateRange(query.minTime, query.maxTime, -1)">
                <f7-icon f7="arrow_left_square"></f7-icon>
            </f7-link>
            <f7-link :class="{ 'tabbar-text-with-ellipsis': true, 'disabled': loading }" popover-open=".date-popover-menu">
                <span :class="{ 'tabbar-item-changed': query.dateType !== allDateRanges.All.type }">{{ queryDateRangeName }}</span>
            </f7-link>
            <f7-link :class="{ 'disabled': loading || query.dateType === allDateRanges.All.type }" @click="shiftDateRange(query.minTime, query.maxTime, 1)">
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
                                        <span>{{ getDisplayYearMonth(transactionMonthList) }}</span>
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
                                      :link="transaction.type !== allTransactionTypes.ModifyBalance ? `/transaction/detail?id=${transaction.id}&type=${transaction.type}` : null"
                                      :key="transaction.id"
                                      v-for="(transaction, idx) in transactionMonthList.items"
                        >
                            <template #media>
                                <div class="display-flex flex-direction-column transaction-date" :style="getTransactionDateStyle(transaction, idx > 0 ? transactionMonthList.items[idx - 1] : null)">
                                    <span class="transaction-day full-line flex-direction-column">
                                        {{ transaction.day }}
                                    </span>
                                    <span class="transaction-day-of-week full-line flex-direction-column">
                                        {{ getWeekdayShortName(transaction) }}
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
                                                    <span v-if="transaction.type === allTransactionTypes.ModifyBalance">
                                                        {{ $t('Modify Balance') }}
                                                    </span>
                                                        <span v-else-if="transaction.type !== allTransactionTypes.ModifyBalance && transaction.category">
                                                        {{ transaction.category.name }}
                                                    </span>
                                                        <span v-else-if="transaction.type !== allTransactionTypes.ModifyBalance && !transaction.category">
                                                        {{ getTransactionTypeName(transaction.type, 'Transaction') }}
                                                    </span>
                                                </div>
                                            </div>
                                            <div class="item-after">
                                                <div class="transaction-amount" v-if="transaction.sourceAccount"
                                                     :class="{ 'text-expense': transaction.type === allTransactionTypes.Expense, 'text-income': transaction.type === allTransactionTypes.Income }">
                                                    <span>{{ getTransactionDisplayAmount(transaction) }}</span>
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
                                                <f7-chip media-bg-color="black" class="transaction-tag"
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
                                                <f7-icon f7="arrow_right" class="transaction-account-arrow" v-if="transaction.sourceAccount && transaction.type === allTransactionTypes.Transfer && transaction.destinationAccount && transaction.sourceAccount.id !== transaction.destinationAccount.id"></f7-icon>
                                                <span v-if="transaction.sourceAccount && transaction.type === allTransactionTypes.Transfer && transaction.destinationAccount && transaction.sourceAccount.id !== transaction.destinationAccount.id">{{ transaction.destinationAccount.name }}</span>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </template>
                            <f7-swipeout-actions right>
                                <f7-swipeout-button color="primary" close
                                                    :text="$t('Duplicate')"
                                                    v-if="transaction.type !== allTransactionTypes.ModifyBalance"
                                                    @click="duplicate(transaction)"></f7-swipeout-button>
                                <f7-swipeout-button color="orange" close
                                                    :text="$t('Edit')"
                                                    v-if="transaction.editable && transaction.type !== allTransactionTypes.ModifyBalance"
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

        <f7-block class="text-align-center" :class="{ 'disabled': loadingMore }" v-show="!loading && hasMoreTransaction">
            <f7-link href="#" @click="loadMore(false)">{{ $t('Load More') }}</f7-link>
        </f7-block>

        <f7-popover class="date-popover-menu"
                    v-model:opened="showDatePopover"
                    @popover:open="scrollPopoverToSelectedItem">
            <f7-list dividers>
                <f7-list-item :title="dateRange.displayName"
                              :class="{ 'list-item-selected': query.dateType === dateRange.type }"
                              :key="dateRange.type"
                              v-for="dateRange in allDateRangesArray"
                              @click="changeDateFilter(dateRange.type)">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.dateType === dateRange.type"></f7-icon>
                    </template>
                    <template #footer>
                        <div v-if="dateRange.type === allDateRanges.Custom.type && query.dateType === allDateRanges.Custom.type && query.minTime && query.maxTime">
                            <span>{{ queryMinTime }}</span>
                            <span>&nbsp;-&nbsp;</span>
                            <br/>
                            <span>{{ queryMaxTime }}</span>
                        </div>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>

        <date-range-selection-sheet :title="$t('Custom Date Range')"
                                    :min-time="customMinDatetime"
                                    :max-time="customMaxDatetime"
                                    v-model:show="showCustomDateRangeSheet"
                                    @dateRange:change="changeCustomDateFilter">
        </date-range-selection-sheet>

        <f7-popover class="category-popover-menu"
                    v-model:opened="showCategoryPopover"
                    @popover:open="scrollPopoverToSelectedItem">
            <f7-list dividers accordion-list>
                <f7-list-item :class="{ 'list-item-selected': !query.categoryIds }" :title="$t('All')" @click="changeCategoryFilter('')">
                    <template #media>
                        <f7-icon f7="rectangle_grid_2x2"></f7-icon>
                    </template>
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="!query.categoryIds"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item :class="{ 'list-item-selected': query.categoryIds && queryAllFilterCategoryIdsCount > 1 }" :title="$t('Multiple Categories')" @click="filterMultipleCategories()">
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
            >
                <f7-list-item divider :title="getTransactionTypeName(getTransactionTypeFromCategoryType(categoryType), 'Type')"></f7-list-item>
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
                                          :title="$t('All')" @click="changeCategoryFilter(category.id)">
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
                    @popover:open="scrollPopoverToSelectedItem">
            <f7-list dividers>
                <f7-list-item :class="{ 'list-item-selected': !query.accountIds }" :title="$t('All')" @click="changeAccountFilter('')">
                    <template #media>
                        <f7-icon f7="rectangle_grid_2x2"></f7-icon>
                    </template>
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="!query.accountIds"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item :class="{ 'list-item-selected': query.accountIds && queryAllFilterAccountIdsCount > 1 }" :title="$t('Multiple Accounts')" @click="filterMultipleAccounts()">
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
                              v-show="(!account.hidden && (!allAccounts[account.parentId] || !allAccounts[account.parentId].hidden)) || query.accountIds === account.id"
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
                <f7-list-item group-title :title="$t('Type')" />
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

                <f7-list-item group-title :title="$t('Amount')" />
                <f7-list-item :class="{ 'list-item-selected': !query.amountFilter }" :title="$t('All')" @click="changeAmountFilter('')">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="!query.amountFilter"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item :key="filterType.type"
                              :class="{ 'list-item-selected': query.amountFilter && query.amountFilter.startsWith(`${filterType.type}:`) }"
                              :title="$t(filterType.name)"
                              v-for="filterType in allAmountFilterTypes"
                              @click="changeAmountFilter(filterType.type)">
                    <template #after>
                        <span class="margin-right-half" v-if="query.amountFilter && query.amountFilter.startsWith(`${filterType.type}:`)">{{ queryAmount }}</span>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.amountFilter && query.amountFilter.startsWith(`${filterType.type}:`)"></f7-icon>
                    </template>
                </f7-list-item>

                <f7-list-item group-title :title="$t('Tags')" />
                <f7-list-item :class="{ 'list-item-selected': !query.tagIds }" :title="$t('All')" @click="changeTagFilter('')">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="!query.tagIds"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item :class="{ 'list-item-selected': query.tagIds && queryAllFilterTagIdsCount > 1 }"
                              :title="$t('Multiple Tags')" @click="filterMultipleTags()" v-if="allAvailableTagsCount > 0">
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="query.tagIds && queryAllFilterTagIdsCount > 1"></f7-icon>
                    </template>
                </f7-list-item>
                <f7-list-item :title="transactionTag.name"
                              :class="{ 'list-item-selected': query.tagIds === transactionTag.id, 'item-in-multiple-selection': queryAllFilterTagIdsCount > 1 && queryAllFilterTagIds[transactionTag.id] }"
                              :key="transactionTag.id"
                              v-for="transactionTag in allTransactionTags"
                              v-show="!transactionTag.hidden || query.tagIds === transactionTag.id"
                              @click="changeTagFilter(transactionTag.id)"
                >
                    <template #before-title>
                        <f7-icon f7="number" class="transaction-tag-name"></f7-icon>
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
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useAccountsStore } from '@/stores/account.js';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';
import { useTransactionTagsStore } from '@/stores/transactionTag.js';
import { useTransactionsStore } from '@/stores/transaction.js';

import numeralConstants from '@/consts/numeral.js';
import datetimeConstants from '@/consts/datetime.js';
import accountConstants from '@/consts/account.js';
import transactionConstants from '@/consts/transaction.js';
import { getNameByKeyValue } from '@/lib/common.js';
import {
    getCurrentUnixTime,
    getSpecifiedDayFirstUnixTime,
    getUtcOffsetByUtcOffsetMinutes,
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffsetMinutes,
    getActualUnixTimeForStore,
    getShiftedDateRangeAndDateType,
    getDateTypeByDateRange,
    getDateRangeByDateType
} from '@/lib/datetime.js';
import { categoryTypeToTransactionType, transactionTypeToCategoryType } from '@/lib/category.js';
import { getUnifiedSelectedAccountsCurrencyOrDefaultCurrency } from '@/lib/account.js';
import { getTransactionDisplayAmount } from '@/lib/transaction.js';
import { onSwipeoutDeleted, scrollToSelectedItem } from '@/lib/ui.mobile.js';

export default {
    props: [
        'f7route',
        'f7router'
    ],
    data() {
        return {
            loading: true,
            loadingError: null,
            loadingMore: false,
            customMinDatetime: 0,
            customMaxDatetime: 0,
            transactionToDelete: null,
            showDatePopover: false,
            showCategoryPopover: false,
            showAccountPopover: false,
            showMorePopover: false,
            showCustomDateRangeSheet: false,
            showDeleteActionSheet: false
        };
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore, useAccountsStore, useTransactionCategoriesStore, useTransactionTagsStore, useTransactionsStore),
        defaultCurrency() {
            return getUnifiedSelectedAccountsCurrencyOrDefaultCurrency(this.allAccounts, this.queryAllFilterAccountIds, this.userStore.currentUserDefaultCurrency);
        },
        canAddTransaction() {
            if (this.query.accountIds && this.queryAllFilterAccountIdsCount === 1) {
                const account = this.allAccounts[this.query.accountIds];

                if (account && account.type === accountConstants.allAccountTypes.MultiSubAccounts) {
                    return false;
                }
            }

            return true;
        },
        currentTimezoneOffsetMinutes() {
            return getTimezoneOffsetMinutes(this.settingsStore.appSettings.timeZone);
        },
        firstDayOfWeek() {
            return this.userStore.currentUserFirstDayOfWeek;
        },
        query() {
            return this.transactionsStore.transactionsFilter;
        },
        queryDateRangeName() {
            if (this.query.dateType === this.allDateRanges.All.type) {
                return this.$t('Date');
            }

            return this.$locale.getDateRangeDisplayName(this.userStore, this.query.dateType, this.query.minTime, this.query.maxTime);
        },
        queryMinTime() {
            return this.$locale.formatUnixTimeToLongDateTime(this.userStore, this.query.minTime);
        },
        queryMaxTime() {
            return this.$locale.formatUnixTimeToLongDateTime(this.userStore, this.query.maxTime);
        },
        queryAllFilterCategoryIds() {
            return this.transactionsStore.allFilterCategoryIds;
        },
        queryAllFilterAccountIds() {
            return this.transactionsStore.allFilterAccountIds;
        },
        queryAllFilterTagIds() {
            return this.transactionsStore.allFilterTagIds;
        },
        queryAllFilterCategoryIdsCount() {
            return this.transactionsStore.allFilterCategoryIdsCount;
        },
        queryAllFilterAccountIdsCount() {
            return this.transactionsStore.allFilterAccountIdsCount;
        },
        queryAllFilterTagIdsCount() {
            return this.transactionsStore.allFilterTagIdsCount;
        },
        queryCategoryName() {
            if (this.queryAllFilterCategoryIdsCount > 1) {
                return this.$t('Multiple Categories');
            }

            return getNameByKeyValue(this.allCategories, this.query.categoryIds, null, 'name', this.$t('Category'));
        },
        queryAccountName() {
            if (this.queryAllFilterAccountIdsCount > 1) {
                return this.$t('Multiple Accounts');
            }

            return getNameByKeyValue(this.allAccounts, this.query.accountIds, null, 'name', this.$t('Account'));
        },
        queryAmount() {
            if (!this.query.amountFilter) {
                return '';
            }

            const amountFilterItems = this.query.amountFilter.split(':');

            if (amountFilterItems.length < 2) {
                return '';
            }

            const displayAmount = [];

            for (let i = 1; i < amountFilterItems.length; i++) {
                displayAmount.push(this.getDisplayCurrency(amountFilterItems[i], false));
            }

            return displayAmount.join(' ~ ');
        },
        transactions() {
            if (this.loading) {
                return [];
            }

            return this.transactionsStore.transactions;
        },
        noTransaction() {
            return this.transactionsStore.noTransaction;
        },
        hasMoreTransaction() {
            return this.transactionsStore.hasMoreTransaction;
        },
        allAmountFilterTypes() {
            return numeralConstants.allAmountFilterTypeArray;
        },
        allTransactionTypes() {
            return transactionConstants.allTransactionTypes;
        },
        allAccounts() {
            return this.accountsStore.allAccountsMap;
        },
        allCategories() {
            return this.transactionCategoriesStore.allTransactionCategoriesMap;
        },
        allPrimaryCategories() {
            const primaryCategories = {};

            for (const categoryType in this.transactionCategoriesStore.allTransactionCategories) {
                if (!Object.prototype.hasOwnProperty.call(this.transactionCategoriesStore.allTransactionCategories, categoryType)) {
                    continue;
                }

                if (this.query.type && this.getTransactionTypeFromCategoryType(categoryType) !== this.query.type) {
                    continue;
                }

                primaryCategories[categoryType] = this.transactionCategoriesStore.allTransactionCategories[categoryType];
            }

            return primaryCategories;
        },
        allTransactionTags() {
            return this.transactionTagsStore.allTransactionTagsMap;
        },
        allAvailableTagsCount() {
            return this.transactionTagsStore.allAvailableTagsCount;
        },
        allDateRanges() {
            return datetimeConstants.allDateRanges;
        },
        allDateRangesArray() {
            return this.$locale.getAllDateRanges(datetimeConstants.allDateRangeScenes.Normal, true);
        },
        showTotalAmountInTransactionListPage() {
            return this.settingsStore.appSettings.showTotalAmountInTransactionListPage;
        },
        showTagInTransactionListPage() {
            return this.settingsStore.appSettings.showTagInTransactionListPage;
        }
    },
    created() {
        const self = this;
        const query = self.f7route.query;

        let dateRange = getDateRangeByDateType(query.dateType ? parseInt(query.dateType) : undefined, self.firstDayOfWeek);

        if (!dateRange &&
            query.dateType === self.allDateRanges.Custom.type.toString() &&
            parseInt(query.maxTime) > 0 && parseInt(query.minTime) > 0) {
            dateRange = {
                dateType: parseInt(query.dateType),
                maxTime: parseInt(query.maxTime),
                minTime: parseInt(query.minTime)
            };
        }

        this.transactionsStore.initTransactionListFilter({
            dateType: dateRange ? dateRange.dateType : undefined,
            maxTime: dateRange ? dateRange.maxTime : undefined,
            minTime: dateRange ? dateRange.minTime : undefined,
            type: parseInt(query.type) > 0 ? parseInt(query.type) : undefined,
            categoryIds: query.categoryIds,
            accountIds: query.accountIds,
            tagIds: query.tagIds
        });

        this.reload(null);
    },
    methods: {
        onPageAfterIn() {
            if (this.transactionsStore.transactionListStateInvalid && !this.loading) {
                this.reload(null);
            }

            this.$routeBackOnError(this.f7router, 'loadingError');
        },
        reload(done) {
            const self = this;
            const force = !!done;

            if (!done) {
                self.loading = true;
            }

            Promise.all([
                self.accountsStore.loadAllAccounts({ force: false }),
                self.transactionCategoriesStore.loadAllCategories({ force: false }),
                self.transactionTagsStore.loadAllTags({ force: false })
            ]).then(() => {
                return self.transactionsStore.loadTransactions({
                    reload: true,
                    force: force,
                    autoExpand: true,
                    defaultCurrency: self.defaultCurrency
                });
            }).then(() => {
                if (done) {
                    done();
                }

                if (force) {
                    self.$toast('Data has been updated');
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

            self.transactionsStore.loadTransactions({
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
            this.transactionsStore.collapseMonthInTransactionList({
                month: month,
                collapse: collapse
            });
        },
        changeDateFilter(dateType) {
            if (dateType === this.allDateRanges.Custom.type) { // Custom
                if (!this.query.minTime || !this.query.maxTime) {
                    this.customMaxDatetime = getActualUnixTimeForStore(getCurrentUnixTime(), this.currentTimezoneOffsetMinutes, getBrowserTimezoneOffsetMinutes());
                    this.customMinDatetime = getSpecifiedDayFirstUnixTime(this.customMaxDatetime);
                } else {
                    this.customMaxDatetime = this.query.maxTime;
                    this.customMinDatetime = this.query.minTime;
                }

                this.showCustomDateRangeSheet = true;
                this.showDatePopover = false;
                return;
            } else if (this.query.dateType === dateType) {
                return;
            }

            const dateRange = getDateRangeByDateType(dateType, this.firstDayOfWeek);

            if (!dateRange) {
                return;
            }

            const changed = this.transactionsStore.updateTransactionListFilter({
                dateType: dateRange.dateType,
                maxTime: dateRange.maxTime,
                minTime: dateRange.minTime
            });

            this.showDatePopover = false;

            if (changed) {
                this.reload(null);
            }
        },
        changeCustomDateFilter(minTime, maxTime) {
            if (!minTime || !maxTime) {
                return;
            }

            const dateType = getDateTypeByDateRange(minTime, maxTime, this.firstDayOfWeek, datetimeConstants.allDateRangeScenes.Normal);

            const changed = this.transactionsStore.updateTransactionListFilter({
                dateType: dateType,
                maxTime: maxTime,
                minTime: minTime
            });

            this.showCustomDateRangeSheet = false;

            if (changed) {
                this.reload(null);
            }
        },
        changeTypeFilter(type) {
            if (this.query.type === type) {
                return;
            }

            let newCategoryFilter = undefined;

            if (type && this.query.categoryIds) {
                newCategoryFilter = '';

                for (let categoryId in this.queryAllFilterCategoryIds) {
                    if (!Object.prototype.hasOwnProperty.call(this.queryAllFilterCategoryIds, categoryId)) {
                        continue;
                    }

                    const category = this.allCategories[categoryId];

                    if (category && category.type === transactionTypeToCategoryType(type)) {
                        if (newCategoryFilter.length > 0) {
                            newCategoryFilter += ',';
                        }

                        newCategoryFilter += categoryId;
                    }
                }
            }

            const changed = this.transactionsStore.updateTransactionListFilter({
                type: type,
                categoryIds: newCategoryFilter
            });

            this.showMorePopover = false;

            if (changed) {
                this.reload(null);
            }
        },
        changeCategoryFilter(categoryIds) {
            if (this.query.categoryIds === categoryIds) {
                return;
            }

            const changed = this.transactionsStore.updateTransactionListFilter({
                categoryIds: categoryIds
            });

            this.showCategoryPopover = false;

            if (changed) {
                this.reload(null);
            }
        },
        changeAccountFilter(accountIds) {
            if (this.query.accountIds === accountIds) {
                return;
            }

            const changed = this.transactionsStore.updateTransactionListFilter({
                accountIds: accountIds
            });

            this.showAccountPopover = false;

            if (changed) {
                this.reload(null);
            }
        },
        changeTagFilter(tagIds) {
            if (this.query.tagIds === tagIds) {
                return;
            }

            const changed = this.transactionsStore.updateTransactionListFilter({
                tagIds: tagIds
            });

            this.showMorePopover = false;

            if (changed) {
                this.reload(null);
            }
        },
        changeAmountFilter(filterType) {
            if (this.query.amountFilter === filterType) {
                return;
            }

            if (filterType) {
                this.showMorePopover = false;
                this.f7router.navigate(`/transaction/filter/amount?type=${filterType}&value=${this.query.amountFilter}`);
                return;
            }

            const changed = this.transactionsStore.updateTransactionListFilter({
                amountFilter: filterType
            });

            this.showMorePopover = false;

            if (changed) {
                this.reload(null);
            }
        },
        changeKeywordFilter(keyword) {
            if (this.query.keyword === keyword) {
                return;
            }

            const changed = this.transactionsStore.updateTransactionListFilter({
                keyword: keyword
            });

            if (changed) {
                this.reload(null);
            }
        },
        filterMultipleCategories() {
            let navigateUrl = '/settings/filter/category?type=transactionListCurrent';

            if (this.allTransactionTypes.Income <= this.query.type && this.query.type <= this.allTransactionTypes.Transfer) {
               navigateUrl += '&allowCategoryTypes=' + transactionTypeToCategoryType(this.query.type);
            }

            this.f7router.navigate(navigateUrl);
        },
        filterMultipleAccounts() {
            this.f7router.navigate('/settings/filter/account?type=transactionListCurrent');
        },
        filterMultipleTags() {
            this.f7router.navigate('/settings/filter/tag?type=transactionListCurrent');
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
                self.$alert('An error occurred');
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

            self.transactionsStore.deleteTransaction({
                transaction: transaction,
                defaultCurrency: self.defaultCurrency,
                beforeResolve: (done) => {
                    onSwipeoutDeleted(self.getTransactionDomId(transaction), done);
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
        shiftDateRange(minTime, maxTime, scale) {
            if (this.query.dateType === this.allDateRanges.All.type) {
                return;
            }

            const newDateRange = getShiftedDateRangeAndDateType(minTime, maxTime, scale, this.firstDayOfWeek, datetimeConstants.allDateRangeScenes.Normal);

            const changed = this.transactionsStore.updateTransactionListFilter({
                dateType: newDateRange.dateType,
                maxTime: newDateRange.maxTime,
                minTime: newDateRange.minTime
            });

            if (changed) {
                this.reload(null);
            }
        },
        scrollPopoverToSelectedItem(event) {
            scrollToSelectedItem(event.$el, '.popover-inner', 'li.list-item-selected');
        },
        getDisplayYearMonth(transactionMonthList) {
            return this.$locale.formatTimeToLongYearMonth(this.userStore, transactionMonthList.yearMonth);
        },
        getDisplayTime(transaction) {
            return this.$locale.formatUnixTimeToShortTime(this.userStore, transaction.time, transaction.utcOffset, this.currentTimezoneOffsetMinutes);
        },
        getDisplayTimezone(transaction) {
            return `UTC${getUtcOffsetByUtcOffsetMinutes(transaction.utcOffset)}`;
        },
        getDisplayMonthTotalAmount(amount, currency, symbol, incomplete) {
            const displayAmount = this.getDisplayCurrency(amount, currency);
            return symbol + displayAmount + (incomplete ? '+' : '');
        },
        getDisplayCurrency(value, currencyCode) {
            return this.$locale.formatAmountWithCurrency(this.settingsStore, this.userStore, value, currencyCode);
        },
        getWeekdayShortName(transaction) {
            return this.$locale.getWeekdayShortName(transaction.dayOfWeek);
        },
        getTransactionTypeName(type, defaultName) {
            switch (type){
                case this.allTransactionTypes.ModifyBalance:
                    return this.$t('Modify Balance');
                case this.allTransactionTypes.Income:
                    return this.$t('Income');
                case this.allTransactionTypes.Expense:
                    return this.$t('Expense');
                case this.allTransactionTypes.Transfer:
                    return this.$t('Transfer');
                default:
                    return this.$t(defaultName);
            }
        },
        getTransactionTypeFromCategoryType(categoryType) {
            return categoryTypeToTransactionType(parseInt(categoryType));
        },
        getTransactionDisplayAmount(transaction) {
            return getTransactionDisplayAmount(transaction, this.queryAllFilterAccountIdsCount, this.queryAllFilterAccountIds, this.getDisplayCurrency);
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
        getCategoryListItemCheckedClass(category, queryCategoryIds) {
            if (queryCategoryIds && queryCategoryIds[category.id]) {
                return {
                    'list-item-checked': true
                };
            }

            for (let i = 0; i < category.subCategories.length; i++) {
                if (queryCategoryIds && queryCategoryIds[category.subCategories[i].id]) {
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

.list.transaction-info-list li.transaction-info .chip.transaction-tag > .chip-media {
    opacity: 0.6;
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
</style>
