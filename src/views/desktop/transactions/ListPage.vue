<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <v-layout>
                    <v-navigation-drawer :permanent="alwaysShowNav" v-model="showNav">
                        <div class="mx-6 my-4">
                            <btn-vertical-group :disabled="loading" :buttons="[
                                { name: $t('All Types'), value: 0 },
                                { name: $t('Modify Balance'), value: 1 },
                                { name: $t('Income'), value: 2 },
                                { name: $t('Expense'), value: 3 },
                                { name: $t('Transfer'), value: 4 }
                            ]" v-model="queryType" />
                        </div>
                        <v-divider />
                        <div class="mx-6 mt-4">
                            <span class="text-subtitle-2">{{ $t('Transactions Per Page') }}</span>
                            <v-select class="mt-2" density="compact" :disabled="loading"
                                      :items="[ 5, 10, 15, 20, 25, 30, 50 ]"
                                      v-model="countPerPage"
                            />
                        </div>
                        <v-tabs show-arrows class="my-4" direction="vertical"
                                :disabled="loading" v-model="recentDateRangeType">
                            <v-tab class="tab-text-truncate" :key="idx" :value="idx" v-for="(recentDateRange, idx) in recentMonthDateRanges"
                                   @click="changeDateFilter(recentDateRange)">
                                <span class="text-truncate">{{ recentDateRange.displayName }}</span>
                            </v-tab>
                        </v-tabs>
                    </v-navigation-drawer>
                    <v-main>
                        <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container" v-model="activeTab">
                            <v-window-item value="transactionPage">
                                <v-card variant="flat" min-height="920">
                                    <template #title>
                                        <div class="title-and-toolbar d-flex align-center text-no-wrap">
                                            <v-btn class="mr-3 d-md-none" density="compact" color="default" variant="plain"
                                                   :ripple="false" :icon="true" @click="showNav = !showNav">
                                                <v-icon :icon="icons.menu" size="24" />
                                            </v-btn>
                                            <span>{{ $t('Transaction List') }}</span>
                                            <v-btn class="ml-3" color="default" variant="outlined"
                                                   :disabled="loading || !canAddTransaction" @click="add">{{ $t('Add') }}</v-btn>
                                            <v-btn density="compact" color="default" variant="text" size="24"
                                                   class="ml-2" :icon="true" :loading="loading" @click="reload">
                                                <template #loader>
                                                    <v-progress-circular indeterminate size="20"/>
                                                </template>
                                                <v-icon :icon="icons.refresh" size="24" />
                                                <v-tooltip activator="parent">{{ $t('Refresh') }}</v-tooltip>
                                            </v-btn>
                                            <v-spacer/>
                                            <div class="transaction-keyword-filter ml-2">
                                                <v-text-field density="compact" :disabled="loading"
                                                              :prepend-inner-icon="icons.search"
                                                              :append-inner-icon="searchKeyword !== query.keyword ? icons.check : null"
                                                              :placeholder="$t('Search transaction description')"
                                                              v-model="searchKeyword"
                                                              @click:append-inner="changeKeywordFilter(searchKeyword)"
                                                              @keyup.enter="changeKeywordFilter(searchKeyword)"
                                                />
                                            </div>
                                        </div>
                                    </template>

                                    <v-card-text class="pt-0">
                                        <div class="transaction-list-datetime-range d-flex align-center">
                                            <span class="text-body-1">{{ $t('Date Range') }}</span>
                                            <span class="text-body-1 transaction-list-datetime-range-text ml-2"
                                                  v-if="!query.minTime && !query.maxTime">
                                                <span class="text-sm">{{ $t('All') }}</span>
                                            </span>
                                            <span class="text-body-1 transaction-list-datetime-range-text ml-2"
                                                  v-else-if="query.minTime || query.maxTime">
                                                <v-btn class="mr-1" size="x-small"
                                                       density="compact" color="default" variant="outlined"
                                                       :icon="icons.arrowLeft" :disabled="loading"
                                                       @click="shiftDateRange(query.minTime, query.maxTime, -1)"/>
                                                <span class="text-sm">{{ `${queryMinTime} - ${queryMaxTime}` }}</span>
                                                <v-btn class="ml-1" size="x-small"
                                                       density="compact" color="default" variant="outlined"
                                                       :icon="icons.arrowRight" :disabled="loading"
                                                       @click="shiftDateRange(query.minTime, query.maxTime, 1)"/>
                                            </span>
                                            <v-spacer/>
                                            <div class="skeleton-no-margin d-flex align-center" v-if="showTotalAmountInTransactionListPage && currentMonthTotalAmount">
                                                <span class="ml-2 text-subtitle-1">{{ $t('Total Income') }}</span>
                                                <span class="text-income ml-2" v-if="loading">
                                                    <v-skeleton-loader type="text" style="width: 60px" :loading="true"></v-skeleton-loader>
                                                </span>
                                                <span class="text-income ml-2" v-else-if="!loading">
                                                    {{ currentMonthTotalAmount.income }}
                                                </span>
                                                <span class="text-subtitle-1 ml-3">{{ $t('Total Expense') }}</span>
                                                <span class="text-expense ml-2" v-if="loading">
                                                    <v-skeleton-loader type="text" style="width: 60px" :loading="true"></v-skeleton-loader>
                                                </span>
                                                <span class="text-expense ml-2" v-else-if="!loading">
                                                    {{ currentMonthTotalAmount.expense }}
                                                </span>
                                            </div>
                                        </div>
                                    </v-card-text>

                                    <v-table class="transaction-table" :hover="!loading">
                                        <thead>
                                        <tr>
                                            <th class="transaction-table-column-time">{{ $t('Time') }}</th>
                                            <th class="transaction-table-column-category">
                                                <v-menu ref="categoryFilterMenu" class="transaction-category-menu"
                                                        eager location="bottom" max-height="500"
                                                        :disabled="query.type === 1"
                                                        :close-on-content-click="false"
                                                        v-model="categoryMenuState"
                                                        @update:model-value="scrollCategoryMenuToSelectedItem">
                                                    <template #activator="{ props }">
                                                        <div class="d-flex align-center"
                                                            :class="{ 'readonly': loading, 'cursor-pointer': query.type !== 1, 'text-primary': query.categoryIds }" v-bind="props">
                                                            <span>{{ queryCategoryName }}</span>
                                                            <v-icon :icon="icons.dropdownMenu" v-show="query.type !== 1" />
                                                        </div>
                                                    </template>
                                                    <v-list :selected="[queryAllSelectedFilterCategoryIds]">
                                                        <v-list-item key="" value="" class="text-sm" density="compact"
                                                                     :class="{ 'list-item-selected': !query.categoryIds }"
                                                                     :append-icon="(!query.categoryIds ? icons.check : null)">
                                                            <v-list-item-title class="cursor-pointer"
                                                                               @click="changeCategoryFilter('')">
                                                                <div class="d-flex align-center">
                                                                    <v-icon :icon="icons.all" />
                                                                    <span class="text-sm ml-3">{{ $t('All') }}</span>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>
                                                        <v-list-item key="multiple" value="multiple" class="text-sm" density="compact"
                                                                     :class="{ 'list-item-selected': query.categoryIds && queryAllFilterCategoryIdsCount > 1 }"
                                                                     :append-icon="(query.categoryIds && queryAllFilterCategoryIdsCount > 1 ? icons.check : null)">
                                                            <v-list-item-title class="cursor-pointer"
                                                                               @click="showFilterCategoryDialog = true">
                                                                <div class="d-flex align-center">
                                                                    <v-icon :icon="icons.multiple" />
                                                                    <span class="text-sm ml-3">{{ $t('Multiple Categories') }}</span>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>

                                                        <template :key="categoryType"
                                                                  v-for="(categories, categoryType) in allPrimaryCategories">
                                                            <v-list-item density="compact">
                                                                <v-list-item-title>
                                                                    <span class="text-sm">{{ getTransactionTypeName(getTransactionTypeFromCategoryType(categoryType), 'Type') }}</span>
                                                                </v-list-item-title>
                                                            </v-list-item>

                                                            <v-list-group :key="category.id" v-for="category in categories">
                                                                <template #activator="{ props }" v-if="!category.hidden || query.categoryIds === category.id || (allCategories[query.categoryIds] && allCategories[query.categoryIds].parentId === category.id)">
                                                                    <v-divider />
                                                                    <v-list-item class="text-sm" density="compact"
                                                                                 :class="getCategoryListItemCheckedClass(category, queryAllFilterCategoryIds)"
                                                                                 v-bind="props">
                                                                        <v-list-item-title>
                                                                            <div class="d-flex align-center">
                                                                                <ItemIcon icon-type="category" size="24px" :icon-id="category.icon" :color="category.color"></ItemIcon>
                                                                                <span class="text-sm ml-3">{{ category.name }}</span>
                                                                            </div>
                                                                        </v-list-item-title>
                                                                    </v-list-item>
                                                                </template>

                                                                <v-divider />
                                                                <v-list-item class="text-sm" density="compact"
                                                                             :class="{ 'item-in-multiple-selection': queryAllFilterCategoryIdsCount > 1 && queryAllFilterCategoryIds[category.id] }"
                                                                             :value="category.id"
                                                                             :append-icon="(query.categoryIds === category.id ? icons.check : null)">
                                                                    <v-list-item-title class="cursor-pointer"
                                                                                       @click="changeCategoryFilter(category.id)">
                                                                        <div class="d-flex align-center">
                                                                            <v-icon :icon="icons.all" />
                                                                            <span class="text-sm ml-3">{{ $t('All') }}</span>
                                                                        </div>
                                                                    </v-list-item-title>
                                                                </v-list-item>

                                                                <template :key="subCategory.id"
                                                                          v-for="subCategory in category.subCategories">
                                                                    <v-divider v-if="!subCategory.hidden || query.categoryIds === subCategory.id" />
                                                                    <v-list-item class="text-sm" density="compact"
                                                                                 :value="subCategory.id"
                                                                                 :class="{ 'list-item-selected': query.categoryIds === subCategory.id, 'item-in-multiple-selection': queryAllFilterCategoryIdsCount > 1 && queryAllFilterCategoryIds[subCategory.id] }"
                                                                                 :append-icon="(query.categoryIds === subCategory.id ? icons.check : null)"
                                                                                 v-if="!subCategory.hidden || query.categoryIds === subCategory.id">
                                                                        <v-list-item-title class="cursor-pointer"
                                                                                           @click="changeCategoryFilter(subCategory.id)">
                                                                            <div class="d-flex align-center">
                                                                                <ItemIcon icon-type="category" size="24px" :icon-id="subCategory.icon" :color="subCategory.color"></ItemIcon>
                                                                                <span class="text-sm ml-3">{{ subCategory.name }}</span>
                                                                            </div>
                                                                        </v-list-item-title>
                                                                    </v-list-item>
                                                                </template>
                                                            </v-list-group>
                                                        </template>
                                                    </v-list>
                                                </v-menu>
                                            </th>
                                            <th class="transaction-table-column-amount">
                                                <v-menu ref="amountFilterMenu" class="transaction-amount-menu"
                                                        eager location="bottom" max-height="500"
                                                        :close-on-content-click="false"
                                                        v-model="amountMenuState"
                                                        @update:model-value="scrollAmountMenuToSelectedItem">
                                                    <template #activator="{ props }">
                                                        <div class="d-flex align-center cursor-pointer"
                                                             :class="{ 'readonly': loading, 'text-primary': query.amountFilter }" v-bind="props">
                                                            <span>{{ $t('Amount') }}</span>
                                                            <v-icon :icon="icons.dropdownMenu" />
                                                        </div>
                                                    </template>
                                                    <v-list :selected="[query.amountFilter.split(':')[0]]">
                                                        <v-list-item key="" value="" class="text-sm" density="compact"
                                                                     :class="{ 'list-item-selected': !query.amountFilter }"
                                                                     :append-icon="(!query.amountFilter && !currentAmountFilterType ? icons.check : null)">
                                                            <v-list-item-title class="cursor-pointer"
                                                                               @click="changeAmountFilter('')">
                                                                <div class="d-flex align-center">
                                                                    <span class="text-sm ml-3">{{ $t('All') }}</span>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>
                                                        <template :key="filterType.type"
                                                                  v-for="filterType in allAmountFilterTypes">
                                                            <v-list-item class="text-sm" density="compact"
                                                                         :value="filterType.type"
                                                                         :class="{ 'list-item-selected': query.amountFilter && query.amountFilter.startsWith(`${filterType.type}:`) }"
                                                                         :append-icon="(query.amountFilter && query.amountFilter.startsWith(`${filterType.type}:`) && currentAmountFilterType !== filterType.type ? icons.check : null)">
                                                                <v-list-item-title class="cursor-pointer"
                                                                                   @click="currentAmountFilterType = filterType.type">
                                                                    <div class="d-flex align-center">
                                                                        <span class="text-sm ml-3">{{ $t(filterType.name) }}</span>
                                                                        <span class="text-sm ml-4" v-if="query.amountFilter && query.amountFilter.startsWith(`${filterType.type}:`) && currentAmountFilterType !== filterType.type">{{ queryAmount }}</span>
                                                                        <amount-input class="transaction-amount-filter-value ml-4" density="compact" v-model="currentAmountFilterValue1"
                                                                                      v-if="currentAmountFilterType === filterType.type"/>
                                                                        <span class="ml-2 mr-2" v-if="currentAmountFilterType === filterType.type && filterType.paramCount === 2">~</span>
                                                                        <amount-input class="transaction-amount-filter-value" density="compact" v-model="currentAmountFilterValue2"
                                                                                      v-if="currentAmountFilterType === filterType.type && filterType.paramCount === 2"/>
                                                                        <v-btn class="ml-2" density="compact" color="primary" variant="tonal"
                                                                               @click="changeAmountFilter(filterType.type)"
                                                                               v-if="currentAmountFilterType === filterType.type">{{ $t('Apply') }}</v-btn>
                                                                    </div>
                                                                </v-list-item-title>
                                                            </v-list-item>
                                                        </template>
                                                    </v-list>
                                                </v-menu>
                                            </th>
                                            <th class="transaction-table-column-account">
                                                <v-menu ref="accountFilterMenu" class="transaction-account-menu"
                                                        eager location="bottom" max-height="500"
                                                        @update:model-value="scrollAccountMenuToSelectedItem">
                                                    <template #activator="{ props }">
                                                        <div class="d-flex align-center cursor-pointer"
                                                             :class="{ 'readonly': loading, 'text-primary': query.accountIds }" v-bind="props">
                                                            <span>{{ queryAccountName }}</span>
                                                            <v-icon :icon="icons.dropdownMenu" />
                                                        </div>
                                                    </template>
                                                    <v-list :selected="[queryAllSelectedFilterAccountIds]">
                                                        <v-list-item key="" value="" class="text-sm" density="compact"
                                                                     :class="{ 'list-item-selected': !query.accountIds }"
                                                                     :append-icon="(!query.accountIds ? icons.check : null)">
                                                            <v-list-item-title class="cursor-pointer"
                                                                               @click="changeAccountFilter('')">
                                                                <div class="d-flex align-center">
                                                                    <v-icon :icon="icons.all" />
                                                                    <span class="text-sm ml-3">{{ $t('All') }}</span>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>
                                                        <v-list-item key="multiple" value="multiple" class="text-sm" density="compact"
                                                                     :class="{ 'list-item-selected': query.accountIds && queryAllFilterAccountIdsCount > 1 }"
                                                                     :append-icon="(query.accountIds && queryAllFilterAccountIdsCount > 1 ? icons.check : null)">
                                                            <v-list-item-title class="cursor-pointer"
                                                                               @click="showFilterAccountDialog = true">
                                                                <div class="d-flex align-center">
                                                                    <v-icon :icon="icons.multiple" />
                                                                    <span class="text-sm ml-3">{{ $t('Multiple Accounts') }}</span>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>
                                                        <template :key="account.id"
                                                                  v-for="account in allAccounts">
                                                            <v-divider v-if="!account.hidden || query.accountIds === account.id" />
                                                            <v-list-item class="text-sm" density="compact"
                                                                         :value="account.id"
                                                                         :class="{ 'list-item-selected': query.accountIds === account.id, 'item-in-multiple-selection': queryAllFilterAccountIdsCount > 1 && queryAllFilterAccountIds[account.id] }"
                                                                         :append-icon="(query.accountIds === account.id ? icons.check : null)"
                                                                         v-if="!account.hidden || query.accountIds === account.id">
                                                                <v-list-item-title class="cursor-pointer"
                                                                                   @click="changeAccountFilter(account.id)">
                                                                    <div class="d-flex align-center">
                                                                        <ItemIcon icon-type="account" size="24px" :icon-id="account.icon" :color="account.color"></ItemIcon>
                                                                        <span class="text-sm ml-3">{{ account.name }}</span>
                                                                    </div>
                                                                </v-list-item-title>
                                                            </v-list-item>
                                                        </template>
                                                    </v-list>
                                                </v-menu>
                                            </th>
                                            <th class="transaction-table-column-tags" v-if="showTagInTransactionListPage">
                                                <v-menu ref="tagFilterMenu" class="transaction-tag-menu"
                                                        eager location="bottom" max-height="500"
                                                        @update:model-value="scrollTagMenuToSelectedItem">
                                                    <template #activator="{ props }">
                                                        <div class="d-flex align-center cursor-pointer"
                                                             :class="{ 'readonly': loading, 'text-primary': query.tagIds }" v-bind="props">
                                                            <span>{{ queryTagName }}</span>
                                                            <v-icon :icon="icons.dropdownMenu" />
                                                        </div>
                                                    </template>
                                                    <v-list :selected="[queryAllSelectedFilterTagIds]">
                                                        <v-list-item key="" value="" class="text-sm" density="compact"
                                                                     :class="{ 'list-item-selected': !query.tagIds }"
                                                                     :append-icon="(!query.tagIds ? icons.check : null)">
                                                            <v-list-item-title class="cursor-pointer"
                                                                               @click="changeTagFilter('')">
                                                                <div class="d-flex align-center">
                                                                    <v-icon :icon="icons.all" />
                                                                    <span class="text-sm ml-3">{{ $t('All') }}</span>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>
                                                        <v-list-item key="multiple" value="multiple" class="text-sm" density="compact"
                                                                     :class="{ 'list-item-selected': query.tagIds && queryAllFilterTagIdsCount > 1 }"
                                                                     :append-icon="(query.tagIds && queryAllFilterTagIdsCount > 1 ? icons.check : null)"
                                                                     v-if="allAvailableTagsCount > 0">
                                                            <v-list-item-title class="cursor-pointer"
                                                                               @click="showFilterTagDialog = true">
                                                                <div class="d-flex align-center">
                                                                    <v-icon :icon="icons.multiple" />
                                                                    <span class="text-sm ml-3">{{ $t('Multiple Tags') }}</span>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>
                                                        <template :key="transactionTag.id"
                                                                  v-for="transactionTag in allTransactionTags">
                                                            <v-divider v-if="!transactionTag.hidden || query.tagIds === transactionTag.id" />
                                                            <v-list-item class="text-sm" density="compact"
                                                                         :value="transactionTag.id"
                                                                         :class="{ 'list-item-selected': query.tagIds === transactionTag.id, 'item-in-multiple-selection': queryAllFilterTagIdsCount > 1 && queryAllFilterTagIds[transactionTag.id] }"
                                                                         :append-icon="(query.tagIds === transactionTag.id ? icons.check : null)"
                                                                         v-if="!transactionTag.hidden || query.tagIds === transactionTag.id">
                                                                <v-list-item-title class="cursor-pointer"
                                                                                   @click="changeTagFilter(transactionTag.id)">
                                                                    <div class="d-flex align-center">
                                                                        <v-icon size="24" :icon="icons.tag"/>
                                                                        <span class="text-sm ml-3">{{ transactionTag.name }}</span>
                                                                    </div>
                                                                </v-list-item-title>
                                                            </v-list-item>
                                                        </template>
                                                    </v-list>
                                                </v-menu>
                                            </th>
                                            <th class="transaction-table-column-description">{{ $t('Description') }}</th>
                                        </tr>
                                        </thead>

                                        <tbody v-if="loading && (!transactions || !transactions.length || transactions.length < 1)">
                                        <tr :key="itemIdx" v-for="itemIdx in skeletonData">
                                            <td class="px-0" :colspan="showTagInTransactionListPage ? 6 : 5">
                                                <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                                            </td>
                                        </tr>
                                        </tbody>

                                        <tbody v-if="!loading && (!transactions || !transactions.length || transactions.length < 1)">
                                        <tr>
                                            <td :colspan="showTagInTransactionListPage ? 6 : 5">{{ $t('No transaction data') }}</td>
                                        </tr>
                                        </tbody>

                                        <tbody :key="transaction.id"
                                               :class="{ 'disabled': loading, 'has-bottom-border': idx < transactions.length - 1 }"
                                               v-for="(transaction, idx) in transactions">
                                            <tr class="transaction-list-row-date no-hover text-sm"
                                                v-if="idx === 0 || (idx > 0 && (transaction.date !== transactions[idx - 1].date))">
                                                <td :colspan="showTagInTransactionListPage ? 6 : 5" class="font-weight-bold">
                                                    <div class="d-flex align-center">
                                                        <span>{{ getLongDate(transaction) }}</span>
                                                        <v-chip class="ml-1" color="default" size="x-small">
                                                            {{ getWeekdayLongName(transaction) }}
                                                        </v-chip>
                                                    </div>
                                                </td>
                                            </tr>
                                            <tr class="transaction-table-row-data text-sm"
                                                :class="{ 'cursor-pointer': transaction.type !== allTransactionTypes.ModifyBalance }"
                                                @click="show(transaction)">
                                                <td class="transaction-table-column-time">
                                                    <div class="d-flex flex-column">
                                                        <span>{{ getDisplayTime(transaction) }}</span>
                                                        <span class="text-caption" v-if="transaction.utcOffset !== currentTimezoneOffsetMinutes">{{ getDisplayTimezone(transaction) }}</span>
                                                        <v-tooltip activator="parent" v-if="transaction.utcOffset !== currentTimezoneOffsetMinutes">{{ getDisplayTimeInDefaultTimezone(transaction) }}</v-tooltip>
                                                    </div>
                                                </td>
                                                <td class="transaction-table-column-category">
                                                    <div class="d-flex align-center">
                                                        <ItemIcon size="24px" icon-type="category"
                                                                  :icon-id="transaction.category.icon"
                                                                  :color="transaction.category.color"
                                                                  v-if="transaction.category && transaction.category.color"></ItemIcon>
                                                        <v-icon size="24" :icon="icons.modifyBalance" v-else-if="!transaction.category || !transaction.category.color" />
                                                        <span class="ml-2" v-if="transaction.type === allTransactionTypes.ModifyBalance">
                                                            {{ $t('Modify Balance') }}
                                                        </span>
                                                        <span class="ml-2" v-else-if="transaction.type !== allTransactionTypes.ModifyBalance && transaction.category">
                                                            {{ transaction.category.name }}
                                                        </span>
                                                        <span class="ml-2" v-else-if="transaction.type !== allTransactionTypes.ModifyBalance && !transaction.category">
                                                            {{ getTransactionTypeName(transaction.type, 'Transaction') }}
                                                        </span>
                                                    </div>
                                                </td>
                                                <td class="transaction-table-column-amount" :class="{ 'text-expense': transaction.type === allTransactionTypes.Expense, 'text-income': transaction.type === allTransactionTypes.Income }">
                                                    <div v-if="transaction.sourceAccount">
                                                        <span>{{ getTransactionDisplayAmount(transaction) }}</span>
                                                    </div>
                                                </td>
                                                <td class="transaction-table-column-account">
                                                    <div class="d-flex align-center">
                                                        <span v-if="transaction.sourceAccount">{{ transaction.sourceAccount.name }}</span>
                                                        <v-icon class="mx-1" size="13" :icon="icons.arrowRight" v-if="transaction.sourceAccount && transaction.type === allTransactionTypes.Transfer && transaction.destinationAccount && transaction.sourceAccount.id !== transaction.destinationAccount.id"></v-icon>
                                                        <span v-if="transaction.sourceAccount && transaction.type === allTransactionTypes.Transfer && transaction.destinationAccount && transaction.sourceAccount.id !== transaction.destinationAccount.id">{{ transaction.destinationAccount.name }}</span>
                                                    </div>
                                                </td>
                                                <td class="transaction-table-column-tags" v-if="showTagInTransactionListPage">
                                                    <v-chip class="transaction-tag" size="small" :prepend-icon="icons.tag"
                                                            :text="allTransactionTags[tagId].name"
                                                            :key="tagId"
                                                            v-for="tagId in transaction.tagIds"/>
                                                    <v-chip class="transaction-tag" size="small" :prepend-icon="icons.tag"
                                                            :text="$t('None')"
                                                            v-if="!transaction.tagIds || !transaction.tagIds.length"/>
                                                </td>
                                                <td class="transaction-table-column-description text-truncate">
                                                    {{ transaction.comment }}
                                                </td>
                                            </tr>
                                        </tbody>
                                    </v-table>

                                    <div class="mt-2 mb-4">
                                        <v-pagination :total-visible="6" :length="totalPageCount"
                                                      v-model="paginationCurrentPage"></v-pagination>
                                    </div>
                                </v-card>
                            </v-window-item>
                        </v-window>
                    </v-main>
                </v-layout>
            </v-card>
        </v-col>
    </v-row>

    <date-range-selection-dialog :title="$t('Custom Date Range')"
                                 :min-time="customMinDatetime"
                                 :max-time="customMaxDatetime"
                                 v-model:show="showCustomDateRangeDialog"
                                 @dateRange:change="changeCustomDateFilter" />
    <edit-dialog ref="editDialog" :persistent="true" />

    <v-dialog width="800" v-model="showFilterAccountDialog">
        <account-filter-settings-card type="transactionListCurrent" :dialog-mode="true"
                                      @settings:change="changeMultipleAccountsFilter" />
    </v-dialog>

    <v-dialog width="800" v-model="showFilterCategoryDialog">
        <category-filter-settings-card type="transactionListCurrent" :dialog-mode="true" :category-types="allowCategoryTypes"
                                       @settings:change="changeMultipleCategoriesFilter" />
    </v-dialog>

    <v-dialog width="800" v-model="showFilterTagDialog">
        <transaction-tag-filter-settings-card type="transactionListCurrent" :dialog-mode="true"
                                       @settings:change="changeMultipleTagsFilter" />
    </v-dialog>

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script>
import EditDialog from './list/dialogs/EditDialog.vue';
import AccountFilterSettingsCard from '@/views/desktop/common/cards/AccountFilterSettingsCard.vue';
import CategoryFilterSettingsCard from '@/views/desktop/common/cards/CategoryFilterSettingsCard.vue';
import TransactionTagFilterSettingsCard from '@/views/desktop/common/cards/TransactionTagFilterSettingsCard.vue';

import { useDisplay } from 'vuetify';

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
import { isString, getNameByKeyValue } from '@/lib/common.js';
import logger from '@/lib/logger.js';
import {
    getCurrentUnixTime,
    parseDateFromUnixTime,
    getYear,
    getMonth,
    getUnixTime,
    getSpecifiedDayFirstUnixTime,
    getUtcOffsetByUtcOffsetMinutes,
    getTimezoneOffset,
    getTimezoneOffsetMinutes,
    getBrowserTimezoneOffsetMinutes,
    getActualUnixTimeForStore,
    getShiftedDateRangeAndDateType,
    getDateTypeByDateRange,
    getDateRangeByDateType,
    getRecentDateRangeType,
    isDateRangeMatchOneMonth
} from '@/lib/datetime.js';
import {
    categoryTypeToTransactionType,
    transactionTypeToCategoryType
} from '@/lib/category.js';
import { getUnifiedSelectedAccountsCurrencyOrDefaultCurrency } from '@/lib/account.js';
import { getTransactionDisplayAmount } from '@/lib/transaction.js';
import { scrollToSelectedItem } from '@/lib/ui.desktop.js';

import {
    mdiMagnify,
    mdiCheck,
    mdiViewGridOutline,
    mdiVectorArrangeBelow,
    mdiRefresh,
    mdiMenu,
    mdiMenuDown,
    mdiPencilBoxOutline,
    mdiArrowLeft,
    mdiArrowRight,
    mdiPound,
    mdiDotsVertical
} from '@mdi/js';

export default {
    components: {
        TransactionTagFilterSettingsCard,
        EditDialog,
        AccountFilterSettingsCard,
        CategoryFilterSettingsCard
    },
    props: [
        'initDateType',
        'initMaxTime',
        'initMinTime',
        'initType',
        'initCategoryIds',
        'initAccountIds',
        'initTagIds',
        'initAmountFilter',
        'initKeyword'
    ],
    data() {
        const { mdAndUp } = useDisplay();

        return {
            loading: true,
            updating: false,
            activeTab: 'transactionPage',
            currentPage: 1,
            temporaryCountPerPage: null,
            totalCount: 1,
            searchKeyword: '',
            customMinDatetime: 0,
            customMaxDatetime: 0,
            currentAmountFilterType: '',
            currentAmountFilterValue1: '0',
            currentAmountFilterValue2: '0',
            currentPageTransactions: [],
            categoryMenuState: false,
            amountMenuState: false,
            alwaysShowNav: mdAndUp.value,
            showNav: mdAndUp.value,
            showCustomDateRangeDialog: false,
            showFilterAccountDialog: false,
            showFilterCategoryDialog: false,
            showFilterTagDialog: false,
            icons: {
                search: mdiMagnify,
                check: mdiCheck,
                all: mdiViewGridOutline,
                multiple: mdiVectorArrangeBelow,
                refresh: mdiRefresh,
                menu: mdiMenu,
                dropdownMenu: mdiMenuDown,
                modifyBalance: mdiPencilBoxOutline,
                arrowLeft: mdiArrowLeft,
                arrowRight: mdiArrowRight,
                tag: mdiPound,
                more: mdiDotsVertical
            }
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
        recentDateRangeType: {
            get: function () {
                return getRecentDateRangeType(this.recentMonthDateRanges, this.query.dateType, this.query.minTime, this.query.maxTime, this.firstDayOfWeek);
            },
            set: function (value) {
                if (value < 0 || value >= this.recentMonthDateRanges.length) {
                    value = 0;
                }

                this.changeDateFilter(this.recentMonthDateRanges[value]);
            }
        },
        query() {
            return this.transactionsStore.transactionsFilter;
        },
        queryType: {
            get: function () {
                return this.query.type;
            },
            set: function(value) {
                this.changeTypeFilter(value);
            }
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
        queryAllSelectedFilterCategoryIds() {
            if (this.queryAllFilterCategoryIdsCount === 0) {
                return '';
            } else if (this.queryAllFilterCategoryIdsCount === 1) {
                return this.query.categoryIds;
            } else { // this.queryAllFilterCategoryIdsCount > 1
                return 'multiple';
            }
        },
        queryAllSelectedFilterAccountIds() {
            if (this.queryAllFilterAccountIdsCount === 0) {
                return '';
            } else if (this.queryAllFilterAccountIdsCount === 1) {
                return this.query.accountIds;
            } else { // this.queryAllFilterAccountIdsCount > 1
                return 'multiple';
            }
        },
        queryAllSelectedFilterTagIds() {
            if (this.queryAllFilterTagIdsCount === 0) {
                return '';
            } else if (this.queryAllFilterTagIdsCount === 1) {
                return this.query.tagIds;
            } else { // this.queryAllFilterTagIdsCount > 1
                return 'multiple';
            }
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
        queryTagName() {
            if (this.queryAllFilterTagIdsCount > 1) {
                return this.$t('Multiple Tags');
            }

            return getNameByKeyValue(this.allTransactionTags, this.query.tagIds, null, 'name', this.$t('Tags'));
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
        queryMonthlyData() {
            return isDateRangeMatchOneMonth(this.query.minTime, this.query.maxTime);
        },
        allowCategoryTypes() {
            if (this.allTransactionTypes.Income <= this.query.type && this.query.type <= this.allTransactionTypes.Transfer) {
                return transactionTypeToCategoryType(this.query.type).toString();
            }

            return '';
        },
        countPerPage: {
            get: function () {
                if (this.temporaryCountPerPage) {
                    return this.temporaryCountPerPage;
                }

                return this.settingsStore.appSettings.itemsCountInTransactionListPage;
            },
            set: function(value) {
                const newTotalPageCount = Math.ceil(this.totalCount / value);

                if (this.currentPage > newTotalPageCount) {
                    this.currentPage = newTotalPageCount;
                }

                this.temporaryCountPerPage = value;

                if (!this.queryMonthlyData) {
                    this.reload(false);
                }
            }
        },
        totalPageCount() {
            return Math.ceil(this.totalCount / this.countPerPage);
        },
        paginationCurrentPage: {
            get: function () {
                return this.currentPage;
            },
            set: function (value) {
                this.currentPage = value;

                if (!this.queryMonthlyData) {
                    this.reload(false);
                }
            }
        },
        skeletonData() {
            const data = [];

            for (let i = 0; i < this.countPerPage; i++) {
                data.push(i);
            }

            return data;
        },
        currentMonthTransactionData() {
            const allTransactions = this.transactionsStore.transactions;

            if (!allTransactions || !allTransactions.length) {
                return null;
            }

            const currentMonthMinDate = parseDateFromUnixTime(this.query.minTime);
            const currentYear = getYear(currentMonthMinDate);
            const currentMonth = getMonth(currentMonthMinDate);

            for (let i = 0; i < allTransactions.length; i++) {
                if (allTransactions[i].year === currentYear && allTransactions[i].month === currentMonth) {
                    return allTransactions[i];
                }
            }

            return null;
        },
        transactions() {
            if (this.queryMonthlyData) {
                const transactionData = this.currentMonthTransactionData;

                if (!transactionData || !transactionData.items) {
                    return [];
                }

                const firstIndex = (this.currentPage - 1) * this.countPerPage;
                const lastIndex = this.currentPage * this.countPerPage;

                return transactionData.items.slice(firstIndex, lastIndex);
            } else {
                return this.currentPageTransactions;
            }
        },
        currentMonthTotalAmount() {
            if (this.queryMonthlyData) {
                const transactionData = this.currentMonthTransactionData;

                if (!transactionData) {
                    return null;
                }

                return {
                    income: this.getDisplayMonthTotalAmount(transactionData.totalAmount.income, this.defaultCurrency, '', transactionData.totalAmount.incompleteIncome),
                    expense: this.getDisplayMonthTotalAmount(transactionData.totalAmount.expense, this.defaultCurrency, '', transactionData.totalAmount.incompleteExpense)
                };
            } else {
                return null;
            }
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
        recentMonthDateRanges() {
            return this.$locale.getAllRecentMonthDateRanges(this.userStore, true, true);
        },
        showTotalAmountInTransactionListPage() {
            return this.settingsStore.appSettings.showTotalAmountInTransactionListPage;
        },
        showTagInTransactionListPage() {
            return this.settingsStore.appSettings.showTagInTransactionListPage;
        }
    },
    created() {
        this.init({
            dateType: this.initDateType,
            minTime: this.initMinTime,
            maxTime: this.initMaxTime,
            type: this.initType,
            categoryIds: this.initCategoryIds,
            accountIds: this.initAccountIds,
            tagIds: this.initTagIds,
            amountFilter: this.initAmountFilter,
            keyword: this.initKeyword
        });
    },
    setup() {
        const display = useDisplay();

        return {
            display: display
        };
    },
    watch: {
        'display.mdAndUp.value': function (newValue) {
            this.alwaysShowNav = newValue;

            if (!this.showNav) {
                this.showNav = newValue;
            }
        }
    },
    beforeRouteUpdate(to) {
        if (to.query) {
            this.init({
                dateType: to.query.dateType,
                minTime: to.query.minTime,
                maxTime: to.query.maxTime,
                type: to.query.type,
                categoryIds: to.query.categoryIds,
                accountIds: to.query.accountIds,
                tagIds: to.query.tagIds,
                amountFilter: to.query.amountFilter,
                keyword: to.query.keyword
            });
        }
    },
    methods: {
        init(query) {
            let dateRange = getDateRangeByDateType(query.dateType ? parseInt(query.dateType) : undefined, this.firstDayOfWeek);

            if (!dateRange &&
                query.dateType === datetimeConstants.allDateRanges.Custom.type.toString() &&
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
                tagIds: query.tagIds,
                amountFilter: query.amountFilter || '',
                keyword: query.keyword || ''
            });

            this.searchKeyword = query.keyword || '';
            this.currentAmountFilterType = '';

            this.currentPage = 1;
            this.reload(false);
        },
        reload(force) {
            const self = this;

            self.loading = true;

            const page = self.currentPage;

            Promise.all([
                self.accountsStore.loadAllAccounts({ force: false }),
                self.transactionCategoriesStore.loadAllCategories({ force: false }),
                self.transactionTagsStore.loadAllTags({ force: false })
            ]).then(() => {
                if (this.queryMonthlyData) {
                    const currentMonthMinDate = parseDateFromUnixTime(this.query.minTime);
                    const currentYear = getYear(currentMonthMinDate);
                    const currentMonth = getMonth(currentMonthMinDate);

                    return self.transactionsStore.loadMonthlyAllTransactions({
                        year: currentYear,
                        month: currentMonth,
                        force: force,
                        autoExpand: true,
                        defaultCurrency: self.defaultCurrency
                    });
                } else {
                    return self.transactionsStore.loadTransactions({
                        reload: true,
                        force: force,
                        count: self.countPerPage,
                        page: page,
                        withCount: page <= 1,
                        autoExpand: true,
                        defaultCurrency: self.defaultCurrency
                    });
                }
            }).then(data => {
                self.loading = false;
                self.currentPageTransactions = data && data.items && data.items.length ? data.items : [];

                if (page <= 1) {
                    self.totalCount = data && data.totalCount ? data.totalCount : 1;
                }

                if (force) {
                    self.$refs.snackbar.showMessage('Data has been updated');
                }
            }).catch(error => {
                self.loading = false;
                self.currentPageTransactions = [];
                self.totalCount = 1;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        shiftDateRange(startTime, endTime, scale) {
            if (this.recentDateRangeType === datetimeConstants.allDateRanges.All.type) {
                return;
            }

            const newDateRange = getShiftedDateRangeAndDateType(startTime, endTime, scale, this.firstDayOfWeek, datetimeConstants.allDateRangeScenes.Normal);

            const changed = this.transactionsStore.updateTransactionListFilter({
                dateType: newDateRange.dateType,
                maxTime: newDateRange.maxTime,
                minTime: newDateRange.minTime
            });

            if (changed) {
                this.loading = true;
                this.currentPageTransactions = [];
                this.transactionsStore.clearTransactions();
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        changeDateFilter(recentDateRange) {
            if (recentDateRange.dateType === datetimeConstants.allDateRanges.Custom.type &&
                !recentDateRange.minTime && !recentDateRange.maxTime) { // Custom
                if (!this.query.minTime || !this.query.maxTime) {
                    this.customMaxDatetime = getActualUnixTimeForStore(getCurrentUnixTime(), this.currentTimezoneOffsetMinutes, getBrowserTimezoneOffsetMinutes());
                    this.customMinDatetime = getSpecifiedDayFirstUnixTime(this.customMaxDatetime);
                } else {
                    this.customMaxDatetime = this.query.maxTime;
                    this.customMinDatetime = this.query.minTime;
                }

                this.showCustomDateRangeDialog = true;
                return;
            }

            if (this.query.dateType === recentDateRange.dateType && this.query.maxTime === recentDateRange.maxTime && this.query.minTime === recentDateRange.minTime) {
                return;
            }

            const changed = this.transactionsStore.updateTransactionListFilter({
                dateType: recentDateRange.dateType,
                maxTime: recentDateRange.maxTime,
                minTime: recentDateRange.minTime
            });

            if (changed) {
                this.loading = true;
                this.currentPageTransactions = [];
                this.transactionsStore.clearTransactions();
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        changeCustomDateFilter(minTime, maxTime) {
            if (!minTime || !maxTime) {
                return;
            }

            const dateType = getDateTypeByDateRange(minTime, maxTime, this.firstDayOfWeek, datetimeConstants.allDateRangeScenes.Normal);

            if (this.query.dateType === dateType && this.query.maxTime === maxTime && this.query.minTime === minTime) {
                this.showCustomDateRangeDialog = false;
                return;
            }

            const changed = this.transactionsStore.updateTransactionListFilter({
                dateType: dateType,
                maxTime: maxTime,
                minTime: minTime
            });

            this.showCustomDateRangeDialog = false;

            if (changed) {
                this.loading = true;
                this.currentPageTransactions = [];
                this.transactionsStore.clearTransactions();
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        changeTypeFilter(type) {
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

            if (changed) {
                this.loading = true;
                this.currentPageTransactions = [];
                this.transactionsStore.clearTransactions();
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        changeCategoryFilter(categoryIds) {
            this.categoryMenuState = false;

            if (this.query.categoryIds === categoryIds) {
                return;
            }

            const changed = this.transactionsStore.updateTransactionListFilter({
                categoryIds: categoryIds
            });

            if (changed) {
                this.loading = true;
                this.currentPageTransactions = [];
                this.transactionsStore.clearTransactions();
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        changeMultipleCategoriesFilter(changed) {
            this.categoryMenuState = false;
            this.showFilterCategoryDialog = false;

            if (changed) {
                this.loading = true;
                this.currentPageTransactions = [];
                this.transactionsStore.clearTransactions();
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        changeAmountFilter(filterType) {
            this.currentAmountFilterType = '';
            this.amountMenuState = false;

            if (this.query.amountFilter === filterType) {
                return;
            }

            let amountFilter = filterType;

            if (filterType) {
                const amountCount = this.getAmountFilterParameterCount(filterType);

                if (!amountCount) {
                    return;
                }

                if (amountCount === 1) {
                    amountFilter += ':' + this.currentAmountFilterValue1;
                } else if (amountCount === 2) {
                    if (this.currentAmountFilterValue2 < this.currentAmountFilterValue1) {
                        this.$refs.snackbar.showMessage('Incorrect amount range');
                        return;
                    }

                    amountFilter += ':' + this.currentAmountFilterValue1 + ':' + this.currentAmountFilterValue2;
                } else {
                    return;
                }
            }

            if (this.query.amountFilter === amountFilter) {
                return;
            }

            const changed = this.transactionsStore.updateTransactionListFilter({
                amountFilter: amountFilter
            });

            if (changed) {
                this.loading = true;
                this.currentPageTransactions = [];
                this.transactionsStore.clearTransactions();
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        changeAccountFilter(accountIds) {
            if (this.query.accountIds === accountIds) {
                return;
            }

            const changed = this.transactionsStore.updateTransactionListFilter({
                accountIds: accountIds
            });

            if (changed) {
                this.loading = true;
                this.currentPageTransactions = [];
                this.transactionsStore.clearTransactions();
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        changeMultipleAccountsFilter(changed) {
            this.showFilterAccountDialog = false;

            if (changed) {
                this.loading = true;
                this.currentPageTransactions = [];
                this.transactionsStore.clearTransactions();
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        changeTagFilter(tagIds) {
            if (this.query.tagIds === tagIds) {
                return;
            }

            const changed = this.transactionsStore.updateTransactionListFilter({
                tagIds: tagIds
            });

            if (changed) {
                this.loading = true;
                this.currentPageTransactions = [];
                this.transactionsStore.clearTransactions();
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        changeMultipleTagsFilter(changed) {
            this.showFilterTagDialog = false;

            if (changed) {
                this.loading = true;
                this.currentPageTransactions = [];
                this.transactionsStore.clearTransactions();
                this.$router.push(this.getFilterLinkUrl());
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
                this.loading = true;
                this.currentPage = 1;
                this.currentPageTransactions = [];
                this.transactionsStore.clearTransactions();
                this.$router.push(this.getFilterLinkUrl());
            }
        },
        add() {
            const self = this;

            self.$refs.editDialog.open({
                type: self.query.type,
                categoryId: self.queryAllFilterCategoryIdsCount === 1 ? self.query.categoryIds : '',
                accountId: self.queryAllFilterAccountIdsCount === 1 ? self.query.accountIds : '',
                tagIds: self.query.tagIds || ''
            }).then(result => {
                if (result && result.message) {
                    self.$refs.snackbar.showMessage(result.message);
                }

                self.reload(false);
            }).catch(error => {
                if (error) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        show(transaction) {
            const self = this;

            if (transaction.type === self.allTransactionTypes.ModifyBalance) {
                return;
            }

            self.$refs.editDialog.open({
                id: transaction.id,
                currentTransaction: transaction
            }).then(result => {
                if (result && result.message) {
                    self.$refs.snackbar.showMessage(result.message);
                }

                self.reload(false);
            }).catch(error => {
                if (error) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        scrollCategoryMenuToSelectedItem(opened) {
            if (opened) {
                this.scrollMenuToSelectedItem(this.$refs.categoryFilterMenu);
            }
        },
        scrollAmountMenuToSelectedItem(opened) {
            if (opened) {
                this.currentAmountFilterType = '';

                let amount1 = 0, amount2 = 0;

                if (isString(this.query.amountFilter)) {
                    try {
                        const filterItems = this.query.amountFilter.split(':');
                        const amountCount = this.getAmountFilterParameterCount(filterItems[0]);

                        if (filterItems.length === 2 && amountCount === 1) {
                            amount1 = parseInt(filterItems[1]);
                        } else if (filterItems.length === 3 && amountCount === 2) {
                            amount1 = parseInt(filterItems[1]);
                            amount2 = parseInt(filterItems[2]);
                        }
                    } catch (ex) {
                        logger.warn('cannot parse amount from filter value, original value is ' + this.query.amountFilter);
                    }
                }

                this.currentAmountFilterValue1 = amount1;
                this.currentAmountFilterValue2 = amount2;

                this.scrollMenuToSelectedItem(this.$refs.amountFilterMenu);
            }
        },
        scrollAccountMenuToSelectedItem(opened) {
            if (opened) {
                this.scrollMenuToSelectedItem(this.$refs.accountFilterMenu);
            }
        },
        scrollTagMenuToSelectedItem(opened) {
            if (opened) {
                this.scrollMenuToSelectedItem(this.$refs.tagFilterMenu);
            }
        },
        scrollMenuToSelectedItem(menu) {
            this.$nextTick(() => {
                scrollToSelectedItem(menu.contentEl, 'div.v-list', 'div.v-list-item.list-item-selected');
            });
        },
        getDisplayTime(transaction) {
            return this.$locale.formatUnixTimeToShortTime(this.userStore, transaction.time, transaction.utcOffset, this.currentTimezoneOffsetMinutes);
        },
        getDisplayTimeInDefaultTimezone(transaction) {
            return `${this.$locale.formatUnixTimeToLongDateTime(this.userStore, transaction.time)} (UTC${getTimezoneOffset(this.settingsStore.appSettings.timeZone)})`;
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
        getLongDate(transaction) {
            const transactionTime = getUnixTime(parseDateFromUnixTime(transaction.time, transaction.utcOffset, this.currentTimezoneOffsetMinutes));
            return this.$locale.formatUnixTimeToLongDate(this.userStore, transactionTime);
        },
        getWeekdayLongName(transaction) {
            return this.$locale.getWeekdayLongName(transaction.dayOfWeek);
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
        getCategoryListItemCheckedClass(category, queryCategoryIds) {
            if (queryCategoryIds && queryCategoryIds[category.id]) {
                return {
                    'list-item-selected': true,
                    'has-children-item-selected': true
                };
            }

            for (let i = 0; i < category.subCategories.length; i++) {
                if (queryCategoryIds && queryCategoryIds[category.subCategories[i].id]) {
                    return {
                        'list-item-selected': true,
                        'has-children-item-selected': true
                    };
                }
            }

            return [];
        },
        getAmountFilterParameterCount(filterType) {
            const amountFilterType = numeralConstants.allAmountFilterTypeMap[filterType];
            return amountFilterType ? amountFilterType.paramCount : 0;
        },
        getFilterLinkUrl() {
            return `/transaction/list?${this.transactionsStore.getTransactionListPageParams()}`;
        }
    }
};
</script>

<style>
.transaction-keyword-filter .v-input--density-compact {
    --v-input-control-height: 36px !important;
    --v-input-padding-top: 5px !important;
    --v-input-padding-bottom: 5px !important;
    --v-input-chips-margin-top: 0px !important;
    --v-input-chips-margin-bottom: 0px !important;
    inline-size: 20rem;
}

.transaction-list-datetime-range {
    min-height: 28px;
    flex-wrap: wrap;
    row-gap: 1rem;
}

.transaction-list-datetime-range .transaction-list-datetime-range-text {
    color: rgba(var(--v-theme-on-background), var(--v-medium-emphasis-opacity)) !important;
}

.v-table.transaction-table .transaction-list-row-date > td {
    height: 40px !important;
}

.transaction-table .transaction-table-column-time {
    width: 110px;
    white-space: nowrap;
}

.transaction-table .transaction-table-column-category {
    width: 140px;
    white-space: nowrap;
}

.transaction-table .transaction-table-column-amount {
    width: 120px;
    white-space: nowrap;
}

.transaction-table .transaction-table-column-account {
    width: 160px;
    white-space: nowrap;
}

.transaction-table .transaction-table-column-tags {
    width: 90px;
    max-width: 300px;
}

.transaction-table-column-description {
    max-width: 300px;
}

.transaction-table .transaction-table-column-category .v-btn,
.transaction-table .transaction-table-column-account .v-btn {
    font-size: 0.75rem;
}

.transaction-table .transaction-table-column-category .v-btn .v-btn__append,
.transaction-table .transaction-table-column-account .v-btn .v-btn__append {
    margin-left: 0in;
}

.transaction-table .transaction-table-column-tags .v-chip.transaction-tag {
    margin-right: 4px;
    margin-top: 2px;
    margin-bottom: 2px;
}

.transaction-table .transaction-table-column-tags .v-chip.transaction-tag > .v-chip__content {
    display: block;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
}

.transaction-category-menu .item-icon,
.transaction-amount-menu .item-icon,
.transaction-account-menu .item-icon,
.transaction-tag-menu .item-icon,
.transaction-table .item-icon {
    padding-bottom: 3px;
}

.transaction-amount-filter-value {
    width: 100px;
}

.transaction-amount-filter-value input.v-field__input {
    min-height: 32px !important;
    padding: 0 8px 0 8px;
}

.transaction-category-menu .has-children-item-selected span,
.transaction-category-menu .item-in-multiple-selection span,
.transaction-account-menu .item-in-multiple-selection span,
.transaction-tag-menu .item-in-multiple-selection span {
    font-weight: bold;
}
</style>
