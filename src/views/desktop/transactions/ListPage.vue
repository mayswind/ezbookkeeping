<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <v-layout>
                    <v-navigation-drawer :permanent="alwaysShowNav" v-model="showNav">
                        <div class="mx-6 my-4">
                            <btn-vertical-group :disabled="loading" :buttons="TransactionListPageType.values().map(item => {
                                return {
                                    name: tt(item.name),
                                    value: item.type
                                }
                            })" v-model="queryPageType" />
                        </div>
                        <v-divider />
                        <div class="mx-6 mt-4">
                            <span class="text-subtitle-2">{{ tt('Transaction Type') }}</span>
                            <v-select
                                item-title="displayName"
                                item-value="type"
                                class="mt-2"
                                density="compact"
                                :disabled="loading"
                                :items="[
                                    { displayName: tt('All Types'), type: 0 },
                                    { displayName: tt('Modify Balance'), type: 1 },
                                    { displayName: tt('Income'), type: 2 },
                                    { displayName: tt('Expense'), type: 3 },
                                    { displayName: tt('Transfer'), type: 4 }
                                ]"
                                v-model="queryType"
                            />
                        </div>
                        <div class="mx-6 mt-4" v-if="pageType === TransactionListPageType.List.type">
                            <span class="text-subtitle-2">{{ tt('Transactions Per Page') }}</span>
                            <v-select class="mt-2" density="compact" :disabled="loading"
                                      :items="[ 5, 10, 15, 20, 25, 30, 50 ]"
                                      v-model="countPerPage"
                            />
                        </div>
                        <v-tabs show-arrows class="my-4" direction="vertical"
                                :disabled="loading" v-model="recentDateRangeIndex">
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
                                            <v-btn class="me-3 d-md-none" density="compact" color="default" variant="plain"
                                                   :ripple="false" :icon="true" @click="showNav = !showNav">
                                                <v-icon :icon="mdiMenu" size="24" />
                                            </v-btn>
                                            <span>{{ tt('Transaction List') }}</span>
                                            <v-btn class="ms-3" color="default" variant="outlined"
                                                   :disabled="loading || !canAddTransaction" @click="add()">
                                                {{ tt('Add') }}
                                                <v-menu activator="parent" :open-on-hover="true" v-if="allTransactionTemplates && allTransactionTemplates.length">
                                                    <v-list>
                                                        <v-list-item :title="template.name"
                                                                     :prepend-icon="mdiTextBoxOutline"
                                                                     :key="template.id"
                                                                     v-for="template in allTransactionTemplates"
                                                                     @click="add(template)"></v-list-item>
                                                    </v-list>
                                                </v-menu>
                                            </v-btn>
                                            <v-btn class="ms-3" color="default" variant="outlined"
                                                   :disabled="loading" @click="importTransaction"
                                                   v-if="isDataImportingEnabled()">
                                                {{ tt('Import') }}
                                                <v-menu activator="parent" :open-on-hover="true" v-if="isDataExportingEnabled()">
                                                    <v-list>
                                                        <v-list-item :disabled="loading || exportingData || !transactions || !transactions.length || transactions.length < 1"
                                                                     @click="exportTransactions('csv')">
                                                            <v-list-item-title>{{ tt('Export to CSV (Comma-separated values) File') }}</v-list-item-title>
                                                        </v-list-item>
                                                        <v-list-item :disabled="loading || exportingData || !transactions || !transactions.length || transactions.length < 1"
                                                                     @click="exportTransactions('tsv')">
                                                            <v-list-item-title>{{ tt('Export to TSV (Tab-separated values) File') }}</v-list-item-title>
                                                        </v-list-item>
                                                    </v-list>
                                                </v-menu>
                                            </v-btn>
                                            <v-btn class="ms-3" color="default" variant="outlined"
                                                   :disabled="loading || exportingData || !transactions || !transactions.length || transactions.length < 1" v-if="!isDataImportingEnabled() && isDataExportingEnabled()">
                                                {{ tt('Export') }}
                                                <v-menu activator="parent">
                                                    <v-list>
                                                        <v-list-item :disabled="loading || exportingData || !transactions || !transactions.length || transactions.length < 1"
                                                                     @click="exportTransactions('csv')">
                                                            <v-list-item-title>{{ tt('Export to CSV (Comma-separated values) File') }}</v-list-item-title>
                                                        </v-list-item>
                                                        <v-list-item :disabled="loading || exportingData || !transactions || !transactions.length || transactions.length < 1"
                                                                     @click="exportTransactions('tsv')">
                                                            <v-list-item-title>{{ tt('Export to TSV (Tab-separated values) File') }}</v-list-item-title>
                                                        </v-list-item>
                                                    </v-list>
                                                </v-menu>
                                            </v-btn>
                                            <v-btn density="compact" color="default" variant="text" size="24"
                                                   class="ms-2" :icon="true" :loading="loading" @click="reload(true, false)">
                                                <template #loader>
                                                    <v-progress-circular indeterminate size="20"/>
                                                </template>
                                                <v-icon :icon="mdiRefresh" size="24" />
                                                <v-tooltip activator="parent">{{ tt('Refresh') }}</v-tooltip>
                                            </v-btn>
                                            <v-spacer/>
                                            <div class="transaction-keyword-filter ms-2">
                                                <v-text-field density="compact" :disabled="loading"
                                                              :prepend-inner-icon="mdiMagnify"
                                                              :append-inner-icon="searchKeyword !== query.keyword ? mdiCheck : undefined"
                                                              :placeholder="tt('Search transaction description')"
                                                              v-model="searchKeyword"
                                                              @click:append-inner="changeKeywordFilter(searchKeyword)"
                                                              @keyup.enter="changeKeywordFilter(searchKeyword)"
                                                />
                                            </div>
                                        </div>
                                    </template>

                                    <v-card-text class="pt-0">
                                        <div class="transaction-list-datetime-range d-flex align-center">
                                            <span class="text-body-1">{{ tt('Date Range') }}</span>
                                            <span class="text-body-1 transaction-list-datetime-range-text ms-2"
                                                  v-if="!query.minTime && !query.maxTime">
                                                <span class="text-sm">{{ tt('All') }}</span>
                                            </span>
                                            <span class="text-body-1 transaction-list-datetime-range-text ms-2"
                                                  v-else-if="query.minTime || query.maxTime">
                                                <v-btn class="button-icon-with-direction me-1" size="x-small"
                                                       density="compact" color="default" variant="outlined"
                                                       :icon="mdiArrowLeft" :disabled="loading"
                                                       @click="shiftDateRange(query.minTime, query.maxTime, -1)"/>
                                                <span class="text-sm">{{ `${queryMinTime} - ${queryMaxTime}` }}</span>
                                                <v-btn class="button-icon-with-direction ms-1" size="x-small"
                                                       density="compact" color="default" variant="outlined"
                                                       :icon="mdiArrowRight" :disabled="loading"
                                                       @click="shiftDateRange(query.minTime, query.maxTime, 1)"/>
                                            </span>
                                            <v-spacer/>
                                            <div class="skeleton-no-margin d-flex align-center" v-if="showTotalAmountInTransactionListPage && currentMonthTotalAmount">
                                                <span class="ms-2 text-subtitle-1">{{ queryAllFilterAccountIdsCount ? tt('Total Inflows') : tt('Total Income') }}</span>
                                                <span class="text-income ms-2" v-if="loading">
                                                    <v-skeleton-loader type="text" style="width: 60px" :loading="true"></v-skeleton-loader>
                                                </span>
                                                <span class="text-income ms-2" v-else-if="!loading">
                                                    {{ currentMonthTotalAmount.income }}
                                                </span>
                                                <span class="text-subtitle-1 ms-3">{{ queryAllFilterAccountIdsCount ? tt('Total Outflows') : tt('Total Expense') }}</span>
                                                <span class="text-expense ms-2" v-if="loading">
                                                    <v-skeleton-loader type="text" style="width: 60px" :loading="true"></v-skeleton-loader>
                                                </span>
                                                <span class="text-expense ms-2" v-else-if="!loading">
                                                    {{ currentMonthTotalAmount.expense }}
                                                </span>
                                            </div>
                                        </div>
                                    </v-card-text>

                                    <v-card-text class="transaction-calendar-container pt-0" v-if="pageType === TransactionListPageType.Calendar.type">
                                        <vue-date-picker inline auto-apply model-type="yyyy-MM-dd"
                                                         month-name-format="long"
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
                                                <div class="transaction-calendar-daily-amounts d-flex flex-column align-center justify-center w-100">
                                                    <span :class="{ 'font-weight-bold': currentMonthTransactionData && currentMonthTransactionData.dailyTotalAmounts[day] }">{{ day }}</span>
                                                    <span class="text-income" v-if="currentMonthTransactionData && currentMonthTransactionData.dailyTotalAmounts[day] && currentMonthTransactionData.dailyTotalAmounts[day].income">{{ getDisplayMonthTotalAmount(currentMonthTransactionData.dailyTotalAmounts[day].income, defaultCurrency, '', currentMonthTransactionData.dailyTotalAmounts[day].incompleteIncome) }}</span>
                                                    <span class="text-expense" v-if="currentMonthTransactionData && currentMonthTransactionData.dailyTotalAmounts[day] && currentMonthTransactionData.dailyTotalAmounts[day].expense">{{ getDisplayMonthTotalAmount(currentMonthTransactionData.dailyTotalAmounts[day].expense, defaultCurrency, '', currentMonthTransactionData.dailyTotalAmounts[day].incompleteExpense) }}</span>
                                                </div>
                                            </template>
                                        </vue-date-picker>
                                    </v-card-text>

                                    <v-table class="transaction-table" :hover="!loading">
                                        <thead>
                                        <tr>
                                            <th class="transaction-table-column-time text-no-wrap">
                                                <v-menu ref="timeFilterMenu" class="transaction-time-menu"
                                                        eager location="bottom" max-height="500"
                                                        @update:model-value="scrollTimeMenuToSelectedItem">
                                                    <template #activator="{ props }">
                                                        <div class="d-flex align-center cursor-pointer"
                                                             :class="{ 'readonly': loading, 'text-primary': query.dateType !== DateRange.ThisMonth.type }" v-bind="props">
                                                            <span>{{ tt('Time') }}</span>
                                                            <v-icon :icon="mdiMenuDown" />
                                                        </div>
                                                    </template>
                                                    <v-list :selected="[query.dateType]">
                                                        <v-list-item class="text-sm" density="compact"
                                                                     :key="dateRange.type" :value="dateRange.type"
                                                                     :class="{ 'list-item-selected': query.dateType === dateRange.type }"
                                                                     :append-icon="(query.dateType === dateRange.type ? mdiCheck : undefined)"
                                                                     v-for="dateRange in allDateRanges">
                                                            <v-list-item-title class="cursor-pointer"
                                                                               @click="changeDateFilter(dateRange.type)">
                                                                <div class="d-flex align-center">
                                                                    <span class="text-sm ms-3">{{ dateRange.displayName }}</span>
                                                                </div>
                                                                <div class="transaction-list-custom-datetime-range ms-3 smaller" v-if="dateRange.isUserCustomRange && query.dateType === dateRange.type && query.minTime && query.maxTime">
                                                                    <span>{{ queryMinTime }}</span>
                                                                    <span>&nbsp;-&nbsp;</span>
                                                                    <br/>
                                                                    <span>{{ queryMaxTime }}</span>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>
                                                    </v-list>
                                                </v-menu>
                                            </th>
                                            <th class="transaction-table-column-category text-no-wrap">
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
                                                            <v-icon :icon="mdiMenuDown" v-show="query.type !== 1" />
                                                        </div>
                                                    </template>
                                                    <v-list :selected="[queryAllSelectedFilterCategoryIds]">
                                                        <v-list-item key="" value="" class="text-sm" density="compact"
                                                                     :class="{ 'list-item-selected': !query.categoryIds }"
                                                                     :append-icon="(!query.categoryIds ? mdiCheck : undefined)">
                                                            <v-list-item-title class="cursor-pointer"
                                                                               @click="changeCategoryFilter('')">
                                                                <div class="d-flex align-center">
                                                                    <v-icon :icon="mdiViewGridOutline" />
                                                                    <span class="text-sm ms-3">{{ tt('All') }}</span>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>
                                                        <v-list-item key="multiple" value="multiple" class="text-sm" density="compact"
                                                                     :class="{ 'list-item-selected': query.categoryIds && queryAllFilterCategoryIdsCount > 1 }"
                                                                     :append-icon="(query.categoryIds && queryAllFilterCategoryIdsCount > 1 ? mdiCheck : undefined)"
                                                                     v-if="allAvailableCategoriesCount > 0">
                                                            <v-list-item-title class="cursor-pointer"
                                                                               @click="showFilterCategoryDialog = true">
                                                                <div class="d-flex align-center">
                                                                    <v-icon :icon="mdiVectorArrangeBelow" />
                                                                    <span class="text-sm ms-3">{{ tt('Multiple Categories') }}</span>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>

                                                        <template :key="categoryType"
                                                                  v-for="(categories, categoryType) in allPrimaryCategories">
                                                            <v-list-item density="compact" v-show="categories && categories.length">
                                                                <v-list-item-title>
                                                                    <span class="text-sm">{{ getTransactionTypeName(categoryTypeToTransactionType(parseInt(categoryType)), 'Type') }}</span>
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
                                                                                <span class="text-sm ms-3">{{ category.name }}</span>
                                                                            </div>
                                                                        </v-list-item-title>
                                                                    </v-list-item>
                                                                </template>

                                                                <v-divider />
                                                                <v-list-item class="text-sm" density="compact"
                                                                             :class="{ 'item-in-multiple-selection': queryAllFilterCategoryIdsCount > 1 && queryAllFilterCategoryIds[category.id] }"
                                                                             :value="category.id"
                                                                             :append-icon="(query.categoryIds === category.id ? mdiCheck : undefined)">
                                                                    <v-list-item-title class="cursor-pointer"
                                                                                       @click="changeCategoryFilter(category.id)">
                                                                        <div class="d-flex align-center">
                                                                            <v-icon :icon="mdiViewGridOutline" />
                                                                            <span class="text-sm ms-3">{{ tt('All') }}</span>
                                                                        </div>
                                                                    </v-list-item-title>
                                                                </v-list-item>

                                                                <template :key="subCategory.id"
                                                                          v-for="subCategory in category.subCategories">
                                                                    <v-divider v-if="!subCategory.hidden || query.categoryIds === subCategory.id" />
                                                                    <v-list-item class="text-sm" density="compact"
                                                                                 :value="subCategory.id"
                                                                                 :class="{ 'list-item-selected': query.categoryIds === subCategory.id, 'item-in-multiple-selection': queryAllFilterCategoryIdsCount > 1 && queryAllFilterCategoryIds[subCategory.id] }"
                                                                                 :append-icon="(query.categoryIds === subCategory.id ? mdiCheck : undefined)"
                                                                                 v-if="!subCategory.hidden || query.categoryIds === subCategory.id">
                                                                        <v-list-item-title class="cursor-pointer"
                                                                                           @click="changeCategoryFilter(subCategory.id)">
                                                                            <div class="d-flex align-center">
                                                                                <ItemIcon icon-type="category" size="24px" :icon-id="subCategory.icon" :color="subCategory.color"></ItemIcon>
                                                                                <span class="text-sm ms-3">{{ subCategory.name }}</span>
                                                                            </div>
                                                                        </v-list-item-title>
                                                                    </v-list-item>
                                                                </template>
                                                            </v-list-group>
                                                        </template>
                                                    </v-list>
                                                </v-menu>
                                            </th>
                                            <th class="transaction-table-column-amount text-no-wrap">
                                                <v-menu ref="amountFilterMenu" class="transaction-amount-menu"
                                                        eager location="bottom" max-height="500"
                                                        :close-on-content-click="false"
                                                        v-model="amountMenuState"
                                                        @update:model-value="scrollAmountMenuToSelectedItem">
                                                    <template #activator="{ props }">
                                                        <div class="d-flex align-center cursor-pointer"
                                                             :class="{ 'readonly': loading, 'text-primary': query.amountFilter }" v-bind="props">
                                                            <span>{{ tt('Amount') }}</span>
                                                            <v-icon :icon="mdiMenuDown" />
                                                        </div>
                                                    </template>
                                                    <v-list :selected="[query.amountFilter.split(':')[0]]">
                                                        <v-list-item key="" value="" class="text-sm" density="compact"
                                                                     :class="{ 'list-item-selected': !query.amountFilter }"
                                                                     :append-icon="(!query.amountFilter && !currentAmountFilterType ? mdiCheck : undefined)">
                                                            <v-list-item-title class="cursor-pointer"
                                                                               @click="changeAmountFilter('')">
                                                                <div class="d-flex align-center">
                                                                    <span class="text-sm ms-3">{{ tt('All') }}</span>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>
                                                        <template :key="filterType.type"
                                                                  v-for="filterType in AmountFilterType.values()">
                                                            <v-list-item class="text-sm" density="compact"
                                                                         :value="filterType.type"
                                                                         :class="{ 'list-item-selected': query.amountFilter && query.amountFilter.startsWith(`${filterType.type}:`) }"
                                                                         :append-icon="(query.amountFilter && query.amountFilter.startsWith(`${filterType.type}:`) && currentAmountFilterType !== filterType.type ? mdiCheck : undefined)">
                                                                <v-list-item-title class="cursor-pointer"
                                                                                   @click="currentAmountFilterType = filterType.type">
                                                                    <div class="d-flex align-center">
                                                                        <span class="text-sm ms-3">{{ tt(filterType.name) }}</span>
                                                                        <span class="text-sm ms-4" v-if="query.amountFilter && query.amountFilter.startsWith(`${filterType.type}:`) && currentAmountFilterType !== filterType.type">{{ queryAmount }}</span>
                                                                        <amount-input class="transaction-amount-filter-value ms-4" density="compact"
                                                                                      :currency="defaultCurrency"
                                                                                      v-model="currentAmountFilterValue1"
                                                                                      v-if="currentAmountFilterType === filterType.type"/>
                                                                        <span class="ms-2 me-2" v-if="currentAmountFilterType === filterType.type && filterType.paramCount === 2">~</span>
                                                                        <amount-input class="transaction-amount-filter-value" density="compact"
                                                                                      :currency="defaultCurrency"
                                                                                      v-model="currentAmountFilterValue2"
                                                                                      v-if="currentAmountFilterType === filterType.type && filterType.paramCount === 2"/>
                                                                        <v-btn class="ms-2" density="compact" color="primary" variant="tonal"
                                                                               @click="changeAmountFilter(filterType.type)"
                                                                               v-if="currentAmountFilterType === filterType.type">{{ tt('Apply') }}</v-btn>
                                                                    </div>
                                                                </v-list-item-title>
                                                            </v-list-item>
                                                        </template>
                                                    </v-list>
                                                </v-menu>
                                            </th>
                                            <th class="transaction-table-column-account text-no-wrap">
                                                <v-menu ref="accountFilterMenu" class="transaction-account-menu"
                                                        eager location="bottom" max-height="500"
                                                        @update:model-value="scrollAccountMenuToSelectedItem">
                                                    <template #activator="{ props }">
                                                        <div class="d-flex align-center cursor-pointer"
                                                             :class="{ 'readonly': loading, 'text-primary': query.accountIds }" v-bind="props">
                                                            <span>{{ queryAccountName }}</span>
                                                            <v-icon :icon="mdiMenuDown" />
                                                        </div>
                                                    </template>
                                                    <v-list :selected="[queryAllSelectedFilterAccountIds]">
                                                        <v-list-item key="" value="" class="text-sm" density="compact"
                                                                     :class="{ 'list-item-selected': !query.accountIds }"
                                                                     :append-icon="(!query.accountIds ? mdiCheck : undefined)">
                                                            <v-list-item-title class="cursor-pointer"
                                                                               @click="changeAccountFilter('')">
                                                                <div class="d-flex align-center">
                                                                    <v-icon :icon="mdiViewGridOutline" />
                                                                    <span class="text-sm ms-3">{{ tt('All') }}</span>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>
                                                        <v-list-item key="multiple" value="multiple" class="text-sm" density="compact"
                                                                     :class="{ 'list-item-selected': query.accountIds && queryAllFilterAccountIdsCount > 1 }"
                                                                     :append-icon="(query.accountIds && queryAllFilterAccountIdsCount > 1 ? mdiCheck : undefined)"
                                                                     v-if="allAvailableAccountsCount > 0">
                                                            <v-list-item-title class="cursor-pointer"
                                                                               @click="showFilterAccountDialog = true">
                                                                <div class="d-flex align-center">
                                                                    <v-icon :icon="mdiVectorArrangeBelow" />
                                                                    <span class="text-sm ms-3">{{ tt('Multiple Accounts') }}</span>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>
                                                        <template :key="account.id"
                                                                  v-for="account in allAccounts">
                                                            <v-divider v-if="(!account.hidden && (!allAccountsMap[account.parentId] || !allAccountsMap[account.parentId].hidden)) || query.accountIds === account.id" />
                                                            <v-list-item class="text-sm" density="compact"
                                                                         :value="account.id"
                                                                         :class="{ 'list-item-selected': query.accountIds === account.id, 'item-in-multiple-selection': queryAllFilterAccountIdsCount > 1 && queryAllFilterAccountIds[account.id] }"
                                                                         :append-icon="(query.accountIds === account.id ? mdiCheck : undefined)"
                                                                         v-if="(!account.hidden && (!allAccountsMap[account.parentId] || !allAccountsMap[account.parentId].hidden)) || query.accountIds === account.id">
                                                                <v-list-item-title class="cursor-pointer"
                                                                                   @click="changeAccountFilter(account.id)">
                                                                    <div class="d-flex align-center">
                                                                        <ItemIcon icon-type="account" size="24px" :icon-id="account.icon" :color="account.color"></ItemIcon>
                                                                        <span class="text-sm ms-3">{{ account.name }}</span>
                                                                    </div>
                                                                </v-list-item-title>
                                                            </v-list-item>
                                                        </template>
                                                    </v-list>
                                                </v-menu>
                                            </th>
                                            <th class="transaction-table-column-tags text-no-wrap" v-if="showTagInTransactionListPage">
                                                <v-menu ref="tagFilterMenu" class="transaction-tag-menu"
                                                        eager location="bottom" max-height="500"
                                                        @update:model-value="scrollTagMenuToSelectedItem">
                                                    <template #activator="{ props }">
                                                        <div class="d-flex align-center cursor-pointer"
                                                             :class="{ 'readonly': loading, 'text-primary': query.tagIds }" v-bind="props">
                                                            <span>{{ queryTagName }}</span>
                                                            <v-icon :icon="mdiMenuDown" />
                                                        </div>
                                                    </template>
                                                    <v-list :selected="[queryAllSelectedFilterTagIds]">
                                                        <v-list-item key="" value="" class="text-sm" density="compact"
                                                                     :class="{ 'list-item-selected': !query.tagIds }"
                                                                     :append-icon="(!query.tagIds ? mdiCheck : undefined)">
                                                            <v-list-item-title class="cursor-pointer"
                                                                               @click="changeTagFilter('')">
                                                                <div class="d-flex align-center">
                                                                    <v-icon :icon="mdiViewGridOutline" />
                                                                    <span class="text-sm ms-3">{{ tt('All') }}</span>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>
                                                        <v-list-item key="none" value="none" class="text-sm" density="compact"
                                                                     :class="{ 'list-item-selected': query.tagIds === 'none' }"
                                                                     :append-icon="(query.tagIds === 'none' ? mdiCheck : undefined)">
                                                            <v-list-item-title class="cursor-pointer"
                                                                               @click="changeTagFilter('none')">
                                                                <div class="d-flex align-center">
                                                                    <v-icon :icon="mdiBorderNoneVariant" />
                                                                    <span class="text-sm ms-3">{{ tt('Without Tags') }}</span>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>
                                                        <v-list-item key="multiple" value="multiple" class="text-sm" density="compact"
                                                                     :class="{ 'list-item-selected': query.tagIds && queryAllFilterTagIdsCount > 1 }"
                                                                     :append-icon="(query.tagIds && queryAllFilterTagIdsCount > 1 ? mdiCheck : undefined)"
                                                                     v-if="allAvailableTagsCount > 0">
                                                            <v-list-item-title class="cursor-pointer"
                                                                               @click="showFilterTagDialog = true">
                                                                <div class="d-flex align-center">
                                                                    <v-icon :icon="mdiVectorArrangeBelow" />
                                                                    <span class="text-sm ms-3">{{ tt('Multiple Tags') }}</span>
                                                                </div>
                                                            </v-list-item-title>
                                                        </v-list-item>

                                                        <v-divider v-if="query.tagIds && query.tagIds !== 'none'" />

                                                        <template v-if="query.tagIds && query.tagIds !== 'none'">
                                                            <v-list-item class="text-sm" density="compact"
                                                                         :key="filterType.type"
                                                                         :value="filterType.type"
                                                                         :append-icon="(query.tagFilterType === filterType.type ? mdiCheck : undefined)"
                                                                         v-for="filterType in allTransactionTagFilterTypes">
                                                                <v-list-item-title class="cursor-pointer"
                                                                                   @click="changeTagFilterType(filterType.type)">
                                                                    <div class="d-flex align-center">
                                                                        <v-icon size="24" :icon="filterType.icon"/>
                                                                        <span class="text-sm ms-3">{{ filterType.displayName }}</span>
                                                                    </div>
                                                                </v-list-item-title>
                                                            </v-list-item>
                                                        </template>

                                                        <template :key="transactionTag.id"
                                                                  v-for="transactionTag in allTransactionTags">
                                                            <v-divider v-if="!transactionTag.hidden || query.tagIds === transactionTag.id" />
                                                            <v-list-item class="text-sm" density="compact"
                                                                         :value="transactionTag.id"
                                                                         :class="{ 'list-item-selected': query.tagIds === transactionTag.id, 'item-in-multiple-selection': queryAllFilterTagIdsCount > 1 && queryAllFilterTagIds[transactionTag.id] }"
                                                                         :append-icon="(query.tagIds === transactionTag.id ? mdiCheck : undefined)"
                                                                         v-if="!transactionTag.hidden || query.tagIds === transactionTag.id">
                                                                <v-list-item-title class="cursor-pointer"
                                                                                   @click="changeTagFilter(transactionTag.id)">
                                                                    <div class="d-flex align-center">
                                                                        <v-icon size="24" :icon="mdiPound"/>
                                                                        <span class="text-sm ms-3">{{ transactionTag.name }}</span>
                                                                    </div>
                                                                </v-list-item-title>
                                                            </v-list-item>
                                                        </template>
                                                    </v-list>
                                                </v-menu>
                                            </th>
                                            <th class="transaction-table-column-description text-no-wrap">{{ tt('Description') }}</th>
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
                                            <td :colspan="showTagInTransactionListPage ? 6 : 5">{{ tt('No transaction data') }}</td>
                                        </tr>
                                        </tbody>

                                        <tbody :key="transaction.id"
                                               :class="{ 'disabled': loading, 'has-bottom-border': idx < transactions.length - 1 }"
                                               v-for="(transaction, idx) in transactions">
                                            <tr class="transaction-list-row-date no-hover text-sm"
                                                v-if="pageType === TransactionListPageType.List.type && (idx === 0 || (idx > 0 && (transaction.gregorianCalendarYearDashMonthDashDay !== transactions[idx - 1].gregorianCalendarYearDashMonthDashDay)))">
                                                <td :colspan="showTagInTransactionListPage ? 6 : 5" class="font-weight-bold">
                                                    <div class="d-flex align-center">
                                                        <span>{{ getDisplayLongDate(transaction) }}</span>
                                                        <v-chip class="ms-1" color="default" size="x-small"
                                                                v-if="transaction.displayDayOfWeek">
                                                            {{ getWeekdayLongName(transaction.displayDayOfWeek) }}
                                                        </v-chip>
                                                    </div>
                                                </td>
                                            </tr>
                                            <tr class="transaction-table-row-data text-sm cursor-pointer"
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
                                                        <v-icon size="24" :icon="mdiPencilBoxOutline" v-else-if="!transaction.category || !transaction.category.color" />
                                                        <span class="ms-2" v-if="transaction.type === TransactionType.ModifyBalance">
                                                            {{ tt('Modify Balance') }}
                                                        </span>
                                                        <span class="ms-2" v-else-if="transaction.type !== TransactionType.ModifyBalance && transaction.category">
                                                            {{ transaction.category.name }}
                                                        </span>
                                                        <span class="ms-2" v-else-if="transaction.type !== TransactionType.ModifyBalance && !transaction.category">
                                                            {{ getTransactionTypeName(transaction.type, 'Transaction') }}
                                                        </span>
                                                    </div>
                                                </td>
                                                <td class="transaction-table-column-amount" :class="{ 'text-expense': transaction.type === TransactionType.Expense, 'text-income': transaction.type === TransactionType.Income }">
                                                    <div v-if="transaction.sourceAccount">
                                                        <span>{{ getDisplayAmount(transaction) }}</span>
                                                    </div>
                                                </td>
                                                <td class="transaction-table-column-account">
                                                    <div class="d-flex align-center">
                                                        <span v-if="transaction.sourceAccount">{{ transaction.sourceAccount.name }}</span>
                                                        <v-icon class="icon-with-direction mx-1" size="13" :icon="mdiArrowRight" v-if="transaction.sourceAccount && transaction.type === TransactionType.Transfer && transaction.destinationAccount && transaction.sourceAccount.id !== transaction.destinationAccount.id"></v-icon>
                                                        <span v-if="transaction.sourceAccount && transaction.type === TransactionType.Transfer && transaction.destinationAccount && transaction.sourceAccount.id !== transaction.destinationAccount.id">{{ transaction.destinationAccount.name }}</span>
                                                    </div>
                                                </td>
                                                <td class="transaction-table-column-tags" v-if="showTagInTransactionListPage">
                                                    <v-chip class="transaction-tag" size="small" :prepend-icon="mdiPound"
                                                            :text="allTransactionTags[tagId].name"
                                                            :key="tagId"
                                                            v-for="tagId in transaction.tagIds"/>
                                                    <v-chip class="transaction-tag" size="small"
                                                            :text="tt('None')"
                                                            v-if="!transaction.tagIds || !transaction.tagIds.length"/>
                                                </td>
                                                <td class="transaction-table-column-description text-truncate">
                                                    {{ transaction.comment }}
                                                </td>
                                            </tr>
                                        </tbody>
                                    </v-table>

                                    <div class="mt-2 mb-4" v-if="pageType === TransactionListPageType.List.type">
                                        <pagination-buttons :totalPageCount="totalPageCount"
                                                            v-model="paginationCurrentPage"></pagination-buttons>
                                    </div>
                                </v-card>
                            </v-window-item>
                        </v-window>
                    </v-main>
                </v-layout>
            </v-card>
        </v-col>
    </v-row>

    <date-range-selection-dialog :title="tt('Custom Date Range')"
                                 :min-time="customMinDatetime"
                                 :max-time="customMaxDatetime"
                                 v-model:show="showCustomDateRangeDialog"
                                 @dateRange:change="changeCustomDateFilter"
                                 @error="onShowDateRangeError" />

    <month-selection-dialog :title="tt('Custom Date Range')"
                            :model-value="queryMonth"
                            v-model:show="showCustomMonthDialog"
                            @update:modelValue="changeCustomMonthDateFilter"
                            @error="onShowDateRangeError" />

    <edit-dialog ref="editDialog" :type="TransactionEditPageType.Transaction" />
    <import-dialog ref="importDialog" :persistent="true" />

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

<script setup lang="ts">
import { VMenu } from 'vuetify/components/VMenu';
import PaginationButtons from '@/components/desktop/PaginationButtons.vue';
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';
import EditDialog from './list/dialogs/EditDialog.vue';
import ImportDialog from './import/ImportDialog.vue';
import AccountFilterSettingsCard from '@/views/desktop/common/cards/AccountFilterSettingsCard.vue';
import CategoryFilterSettingsCard from '@/views/desktop/common/cards/CategoryFilterSettingsCard.vue';
import TransactionTagFilterSettingsCard from '@/views/desktop/common/cards/TransactionTagFilterSettingsCard.vue';
import { TransactionEditPageType } from '@/views/base/transactions/TransactionEditPageBase.ts';

import { ref, computed, useTemplateRef, watch, nextTick } from 'vue';
import { useRouter, onBeforeRouteUpdate } from 'vue-router';
import { useDisplay, useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';
import { TransactionListPageType, useTransactionListPageBase } from '@/views/base/transactions/TransactionListPageBase.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useTransactionTemplatesStore } from '@/stores/transactionTemplate.ts';
import { useDesktopPageStore } from '@/stores/desktopPage.ts';

import type { TypeAndDisplayName } from '@/core/base.ts';
import {
    type Year0BasedMonth,
    type LocalizedRecentMonthDateRange,
    type TimeRangeAndDateType,
    DateRangeScene,
    DateRange
} from '@/core/datetime.ts';
import { AmountFilterType } from '@/core/numeral.ts';
import { ThemeType } from '@/core/theme.ts';
import { TransactionType, TransactionTagFilterType } from '@/core/transaction.ts';
import { TemplateType }  from '@/core/template.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';
import type { Transaction } from '@/models/transaction.ts';
import type { TransactionTemplate } from '@/models/transaction_template.ts';

import {
    isObject,
    isString,
    isNumber,
    arrangeArrayWithNewStartIndex
} from '@/lib/common.ts';
import {
    getCurrentUnixTime,
    parseDateTimeFromUnixTime,
    getBrowserTimezoneOffsetMinutes,
    getActualUnixTimeForStore,
    getDayFirstUnixTimeBySpecifiedUnixTime,
    getYearMonthFirstUnixTime,
    getYearMonthLastUnixTime,
    getShiftedDateRangeAndDateType,
    getShiftedDateRangeAndDateTypeForBillingCycle,
    getDateTypeByDateRange,
    getDateTypeByBillingCycleDateRange,
    getDateRangeByDateType,
    getDateRangeByBillingCycleDateType,
    getRecentDateRangeIndex,
    getFullMonthDateRange,
    getValidMonthDayOrCurrentDayShortDate
} from '@/lib/datetime.ts';
import {
    categoryTypeToTransactionType,
    transactionTypeToCategoryType
} from '@/lib/category.ts';
import { isDataExportingEnabled, isDataImportingEnabled } from '@/lib/server_settings.ts';
import { startDownloadFile } from '@/lib/ui/common.ts';
import { scrollToSelectedItem } from '@/lib/ui/desktop.ts';
import logger from '@/lib/logger.ts';

import {
    mdiMagnify,
    mdiCheck,
    mdiViewGridOutline,
    mdiBorderNoneVariant,
    mdiVectorArrangeBelow,
    mdiRefresh,
    mdiMenu,
    mdiMenuDown,
    mdiPencilBoxOutline,
    mdiArrowLeft,
    mdiArrowRight,
    mdiPlusBoxMultipleOutline,
    mdiCheckboxMultipleOutline,
    mdiMinusBoxMultipleOutline,
    mdiCloseBoxMultipleOutline,
    mdiPound,
    mdiTextBoxOutline
} from '@mdi/js';

interface TransactionListProps {
    initPageType?: string;
    initDateType?: string,
    initMaxTime?: string,
    initMinTime?: string,
    initType?: string,
    initCategoryIds?: string,
    initAccountIds?: string,
    initTagIds?: string,
    initTagFilterType?: string,
    initAmountFilter?: string,
    initKeyword?: string
}

const props = defineProps<TransactionListProps>();

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;
type EditDialogType = InstanceType<typeof EditDialog>;
type ImportDialogType = InstanceType<typeof ImportDialog>;

interface TransactionTemplateWithIcon {
    type: number;
    displayName: string;
    icon: string;
}

interface TransactionListDisplayTotalAmount {
    income: string;
    expense: string;
}

const router = useRouter();
const display = useDisplay();
const theme = useTheme();

const {
    tt,
    getAllLongWeekdayNames,
    getAllRecentMonthDateRanges,
    getAllTransactionTagFilterTypes,
    getWeekdayLongName
} = useI18n();

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
    query,
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
    queryTagName,
    queryAmount,
    transactionCalendarMinDate,
    transactionCalendarMaxDate,
    currentMonthTransactionData,
    noTransactionInMonthDay,
    canAddTransaction,
    getDisplayTime,
    getDisplayLongDate,
    getDisplayTimezone,
    getDisplayTimeInDefaultTimezone,
    getDisplayAmount,
    getDisplayMonthTotalAmount,
    getTransactionTypeName,
} = useTransactionListPageBase();

const settingsStore = useSettingsStore();
const userStore = useUserStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();
const transactionsStore = useTransactionsStore();
const transactionTemplatesStore = useTransactionTemplatesStore();
const desktopPageStore = useDesktopPageStore();

const tagFilterIconMap: Record<number, string> = {
    [TransactionTagFilterType.HasAny.type]: mdiPlusBoxMultipleOutline,
    [TransactionTagFilterType.HasAll.type]: mdiCheckboxMultipleOutline,
    [TransactionTagFilterType.NotHasAny.type]: mdiMinusBoxMultipleOutline,
    [TransactionTagFilterType.NotHasAll.type]: mdiCloseBoxMultipleOutline
};

const timeFilterMenu = useTemplateRef<VMenu>('timeFilterMenu');
const categoryFilterMenu = useTemplateRef<VMenu>('categoryFilterMenu');
const amountFilterMenu = useTemplateRef<VMenu>('amountFilterMenu');
const accountFilterMenu = useTemplateRef<VMenu>('accountFilterMenu');
const tagFilterMenu = useTemplateRef<VMenu>('tagFilterMenu');

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');
const editDialog = useTemplateRef<EditDialogType>('editDialog');
const importDialog = useTemplateRef<ImportDialogType>('importDialog');

const activeTab = ref<string>('transactionPage');
const currentPage = ref<number>(1);
const temporaryCountPerPage = ref<number | null>(null);
const totalCount = ref<number>(1);
const searchKeyword = ref<string>('');
const currentAmountFilterType = ref<string>('');
const currentAmountFilterValue1 = ref<number>(0);
const currentAmountFilterValue2 = ref<number>(0);
const currentPageTransactions = ref<Transaction[]>([]);
const categoryMenuState = ref<boolean>(false);
const amountMenuState = ref<boolean>(false);
const exportingData = ref<boolean>(false);
const alwaysShowNav = ref<boolean>(display.mdAndUp.value);
const showNav = ref<boolean>(display.mdAndUp.value);
const showCustomDateRangeDialog = ref<boolean>(false);
const showCustomMonthDialog = ref<boolean>(false);
const showFilterAccountDialog = ref<boolean>(false);
const showFilterCategoryDialog = ref<boolean>(false);
const showFilterTagDialog = ref<boolean>(false);

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);
const dayNames = computed<string[]>(() => arrangeArrayWithNewStartIndex(getAllLongWeekdayNames(), firstDayOfWeek.value));

const recentMonthDateRanges = computed<LocalizedRecentMonthDateRange[]>(() => getAllRecentMonthDateRanges(pageType.value === TransactionListPageType.List.type, true));

const allTransactionTemplates = computed<TransactionTemplate[]>(() => {
    const allTemplates = transactionTemplatesStore.allVisibleTemplates;
    return allTemplates[TemplateType.Normal.type] || [];
});

const allTransactionTagFilterTypes = computed<TransactionTemplateWithIcon[]>(() => {
    const allTagFilterTypes: TypeAndDisplayName[] = getAllTransactionTagFilterTypes();
    const allTagFilterTypesWithIcon: TransactionTemplateWithIcon[] = [];

    for (let i = 0; i < allTagFilterTypes.length; i++) {
        allTagFilterTypesWithIcon.push({
            type: allTagFilterTypes[i].type,
            displayName: allTagFilterTypes[i].displayName,
            icon: tagFilterIconMap[allTagFilterTypes[i].type]
        });
    }

    return allTagFilterTypesWithIcon;
});

const allowCategoryTypes = computed<string>(() => {
    if (TransactionType.Income <= query.value.type && query.value.type <= TransactionType.Transfer) {
        return transactionTypeToCategoryType(query.value.type)?.toString() ?? '';
    }

    return '';
});

const transactions = computed<Transaction[]>(() => {
    if (pageType.value === TransactionListPageType.List.type) {
        if (queryMonthlyData.value) {
            const transactionData = currentMonthTransactionData.value;

            if (!transactionData || !transactionData.items) {
                return [];
            }

            const firstIndex = (currentPage.value - 1) * countPerPage.value;
            const lastIndex = currentPage.value * countPerPage.value;

            return transactionData.items.slice(firstIndex, lastIndex);
        } else {
            return currentPageTransactions.value;
        }
    } else if (pageType.value === TransactionListPageType.Calendar.type) {
        if (queryMonthlyData.value) {
            const transactionData = currentMonthTransactionData.value;

            if (!transactionData || !transactionData.items) {
                return [];
            }

            const transactions :Transaction[] = [];

            for (let i = 0; i < transactionData.items.length; i++) {
                const transaction = transactionData.items[i];

                if (transaction.gregorianCalendarYearDashMonthDashDay === currentCalendarDate.value) {
                    transactions.push(transaction);
                }
            }

            return transactions;
        } else {
            return [];
        }
    } else {
        return [];
    }
});

const recentDateRangeIndex = computed<number>({
    get: () => getRecentDateRangeIndex(recentMonthDateRanges.value, query.value.dateType, query.value.minTime, query.value.maxTime, firstDayOfWeek.value, fiscalYearStart.value),
    set: (value) => {
        if (value < 0 || value >= recentMonthDateRanges.value.length) {
            value = 0;
        }

        changeDateFilter(recentMonthDateRanges.value[value]);
    }
});

const queryPageType = computed<number>({
    get: () => pageType.value,
    set: (value) => changePageType(value)
});

const queryType = computed<number>({
    get: () => query.value.type,
    set: (value) => changeTypeFilter(value)
});

const queryAllSelectedFilterCategoryIds = computed<string>(() => {
    if (queryAllFilterCategoryIdsCount.value === 0) {
        return '';
    } else if (queryAllFilterCategoryIdsCount.value === 1) {
        return query.value.categoryIds;
    } else { // queryAllFilterCategoryIdsCount.value > 1
        return 'multiple';
    }
});

const queryAllSelectedFilterAccountIds = computed<string>(() => {
    if (queryAllFilterAccountIdsCount.value === 0) {
        return '';
    } else if (queryAllFilterAccountIdsCount.value === 1) {
        return query.value.accountIds;
    } else { // queryAllFilterAccountIdsCount.value > 1
        return 'multiple';
    }
});

const queryAllSelectedFilterTagIds = computed<string>(() => {
    if (queryAllFilterTagIdsCount.value === 0) {
        return '';
    } else if (queryAllFilterTagIdsCount.value === 1) {
        return query.value.tagIds;
    } else { // queryAllFilterTagIdsCount.value > 1
        return 'multiple';
    }
});

const countPerPage = computed<number>({
    get: () => {
        if (temporaryCountPerPage.value) {
            return temporaryCountPerPage.value;
        }

        return settingsStore.appSettings.itemsCountInTransactionListPage;
    },
    set: (value) => {
        const newTotalPageCount = Math.ceil(totalCount.value / value);

        if (currentPage.value > newTotalPageCount) {
            currentPage.value = newTotalPageCount;
        }

        temporaryCountPerPage.value = value;

        if (!queryMonthlyData.value) {
            reload(false, false);
        }
    }
});

const totalPageCount = computed<number>(() => Math.ceil(totalCount.value / countPerPage.value));

const paginationCurrentPage = computed<number>({
    get: () => currentPage.value,
    set: (value) => {
        currentPage.value = value;

        if (!queryMonthlyData.value) {
            reload(false, false);
        }
    }
});

const skeletonData = computed<number[]>(() => {
    const data: number[] = [];
    const totalCount = (pageType.value === TransactionListPageType.List.type ? countPerPage.value : 3);

    for (let i = 0; i < totalCount; i++) {
        data.push(i);
    }

    return data;
});

const currentMonthTotalAmount = computed<TransactionListDisplayTotalAmount | null>(() => {
    if (queryMonthlyData.value) {
        const transactionData = currentMonthTransactionData.value;

        if (!transactionData) {
            return null;
        }

        return {
            income: getDisplayMonthTotalAmount(transactionData.totalAmount.income, defaultCurrency.value, '', transactionData.totalAmount.incompleteIncome),
            expense: getDisplayMonthTotalAmount(transactionData.totalAmount.expense, defaultCurrency.value, '', transactionData.totalAmount.incompleteExpense)
        };
    } else {
        return null;
    }
});

function getCategoryListItemCheckedClass(category: TransactionCategory, queryCategoryIds: Record<string, boolean>): Record<string, boolean> {
    if (queryCategoryIds && queryCategoryIds[category.id]) {
        return {
            'list-item-selected': true,
            'has-children-item-selected': true
        };
    }

    if (category.subCategories) {
        for (let i = 0; i < category.subCategories.length; i++) {
            if (queryCategoryIds && queryCategoryIds[category.subCategories[i].id]) {
                return {
                    'list-item-selected': true,
                    'has-children-item-selected': true
                };
            }
        }
    }

    return {};
}

function getAmountFilterParameterCount(filterType: string): number {
    const amountFilterType = AmountFilterType.valueOf(filterType);
    return amountFilterType ? amountFilterType.paramCount : 0;
}

function updateUrlWhenChanged(changed: boolean): void {
    if (changed) {
        loading.value = true;
        currentPageTransactions.value = [];
        transactionsStore.clearTransactions();
        router.push(`/transaction/list?${transactionsStore.getTransactionListPageParams(pageType.value)}`);
    }
}

function init(initProps: TransactionListProps): void {
    let dateRange: TimeRangeAndDateType | null = getDateRangeByDateType(initProps.initDateType ? parseInt(initProps.initDateType) : undefined, firstDayOfWeek.value, fiscalYearStart.value);

    if (!dateRange && initProps.initDateType && initProps.initMaxTime && initProps.initMinTime &&
        (DateRange.isBillingCycle(parseInt(initProps.initDateType)) || initProps.initDateType === DateRange.Custom.type.toString()) &&
        parseInt(initProps.initMaxTime) > 0 && parseInt(initProps.initMinTime) > 0) {
        dateRange = {
            dateType: parseInt(initProps.initDateType),
            maxTime: parseInt(initProps.initMaxTime),
            minTime: parseInt(initProps.initMinTime)
        };
    }

    transactionsStore.initTransactionListFilter({
        dateType: dateRange ? dateRange.dateType : undefined,
        maxTime: dateRange ? dateRange.maxTime : undefined,
        minTime: dateRange ? dateRange.minTime : undefined,
        type: initProps.initType && parseInt(initProps.initType) > 0 ? parseInt(initProps.initType) : undefined,
        categoryIds: initProps.initCategoryIds,
        accountIds: initProps.initAccountIds,
        tagIds: initProps.initTagIds,
        tagFilterType: initProps.initTagFilterType && parseInt(initProps.initTagFilterType) >= 0 ? parseInt(initProps.initTagFilterType) : undefined,
        amountFilter: initProps.initAmountFilter || '',
        keyword: initProps.initKeyword || ''
    });

    if (initProps.initPageType) {
        const type = TransactionListPageType.valueOf(parseInt(initProps.initPageType));

        if (type) {
            pageType.value = type.type;
            currentCalendarDate.value = getValidMonthDayOrCurrentDayShortDate(query.value.minTime, currentCalendarDate.value);

            if (pageType.value === TransactionListPageType.Calendar.type) {
                const dateRange = getFullMonthDateRange(query.value.minTime, query.value.maxTime, firstDayOfWeek.value, fiscalYearStart.value);

                if (dateRange) {
                    const changed = transactionsStore.updateTransactionListFilter({
                        dateType: dateRange.dateType,
                        maxTime: dateRange.maxTime,
                        minTime: dateRange.minTime
                    });

                    if (changed) {
                        currentCalendarDate.value = getValidMonthDayOrCurrentDayShortDate(query.value.minTime, currentCalendarDate.value);
                        updateUrlWhenChanged(changed);
                        return;
                    }
                }
            }
        }
    }

    searchKeyword.value = initProps.initKeyword || '';
    currentAmountFilterType.value = '';

    currentPage.value = 1;
    reload(false, true);

    transactionTemplatesStore.loadAllTemplates({
        templateType: TemplateType.Normal.type,
        force: false
    });
}

function reload(force: boolean, init: boolean): void {
    loading.value = true;

    const page = currentPage.value;

    Promise.all([
        accountsStore.loadAllAccounts({ force: false }),
        transactionCategoriesStore.loadAllCategories({ force: false }),
        transactionTagsStore.loadAllTags({ force: false })
    ]).then(() => {
        if (init) {
            if (desktopPageStore.showAddTransactionDialogInTransactionList) {
                desktopPageStore.resetShowAddTransactionDialogInTransactionList();
                add();
            }
        }

        if (queryMonthlyData.value) {
            const currentMonthMinDate = parseDateTimeFromUnixTime(query.value.minTime);
            const currentYear = currentMonthMinDate.getGregorianCalendarYear();
            const currentMonth = currentMonthMinDate.getGregorianCalendarMonth();

            return transactionsStore.loadMonthlyAllTransactions({
                year: currentYear,
                month: currentMonth,
                autoExpand: true,
                defaultCurrency: defaultCurrency.value
            });
        } else {
            return transactionsStore.loadTransactions({
                reload: true,
                count: countPerPage.value,
                page: page,
                withCount: page <= 1,
                autoExpand: true,
                defaultCurrency: defaultCurrency.value
            });
        }
    }).then(data => {
        loading.value = false;
        currentPageTransactions.value = data && data.items && data.items.length ? data.items : [];

        if (page <= 1) {
            totalCount.value = data && data.totalCount ? data.totalCount : 1;
        }

        if (force) {
            snackbar.value?.showMessage('Data has been updated');
        }
    }).catch(error => {
        loading.value = false;
        currentPageTransactions.value = [];
        totalCount.value = 1;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function changePageType(type: number): void {
    pageType.value = type;
    currentCalendarDate.value = getValidMonthDayOrCurrentDayShortDate(query.value.minTime, currentCalendarDate.value);

    if (pageType.value === TransactionListPageType.Calendar.type) {
        const dateRange = getFullMonthDateRange(query.value.minTime, query.value.maxTime, firstDayOfWeek.value, fiscalYearStart.value);

        if (dateRange) {
            transactionsStore.updateTransactionListFilter({
                dateType: dateRange.dateType,
                maxTime: dateRange.maxTime,
                minTime: dateRange.minTime
            });
            currentCalendarDate.value = getValidMonthDayOrCurrentDayShortDate(query.value.minTime, currentCalendarDate.value);
        }
    }

    updateUrlWhenChanged(true);
}

function changeDateFilter(dateRange: TimeRangeAndDateType | number | null): void {
    if (dateRange === DateRange.Custom.type || (isObject(dateRange) && dateRange.dateType === DateRange.Custom.type && !dateRange.minTime && !dateRange.maxTime)) { // Custom
        if (!query.value.minTime || !query.value.maxTime) {
            customMaxDatetime.value = getActualUnixTimeForStore(getCurrentUnixTime(), currentTimezoneOffsetMinutes.value, getBrowserTimezoneOffsetMinutes());
            customMinDatetime.value = getDayFirstUnixTimeBySpecifiedUnixTime(customMaxDatetime.value);
        } else {
            customMaxDatetime.value = query.value.maxTime;
            customMinDatetime.value = query.value.minTime;
        }

        if (pageType.value === TransactionListPageType.Calendar.type) {
            showCustomMonthDialog.value = true;
        } else {
            showCustomDateRangeDialog.value = true;
        }

        return;
    }

    if (isNumber(dateRange)) {
        if (DateRange.isBillingCycle(dateRange)) {
            dateRange = getDateRangeByBillingCycleDateType(dateRange, firstDayOfWeek.value, fiscalYearStart.value, accountsStore.getAccountStatementDate(query.value.accountIds));
        } else {
            dateRange = getDateRangeByDateType(dateRange, firstDayOfWeek.value, fiscalYearStart.value);
        }
    }

    if (!dateRange) {
        return;
    }

    if (pageType.value === TransactionListPageType.Calendar.type) {
        currentCalendarDate.value = getValidMonthDayOrCurrentDayShortDate(dateRange.minTime, currentCalendarDate.value);
        const fullMonthDateRange = getFullMonthDateRange(dateRange.minTime, dateRange.maxTime, firstDayOfWeek.value, fiscalYearStart.value);

        if (fullMonthDateRange) {
            dateRange = fullMonthDateRange;
            currentCalendarDate.value = getValidMonthDayOrCurrentDayShortDate(dateRange.minTime, currentCalendarDate.value);
        }
    }

    if (query.value.dateType === dateRange.dateType && query.value.maxTime === dateRange.maxTime && query.value.minTime === dateRange.minTime) {
        return;
    }

    const changed = transactionsStore.updateTransactionListFilter({
        dateType: dateRange.dateType,
        maxTime: dateRange.maxTime,
        minTime: dateRange.minTime
    });

    updateUrlWhenChanged(changed);
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
        currentCalendarDate.value = getValidMonthDayOrCurrentDayShortDate(minTime, currentCalendarDate.value);
        const dateRange = getFullMonthDateRange(minTime, maxTime, firstDayOfWeek.value, fiscalYearStart.value);

        if (dateRange) {
            minTime = dateRange.minTime;
            maxTime = dateRange.maxTime;
            dateType = dateRange.dateType;
            currentCalendarDate.value = getValidMonthDayOrCurrentDayShortDate(minTime, currentCalendarDate.value);
        }
    }

    if (query.value.dateType === dateType && query.value.maxTime === maxTime && query.value.minTime === minTime) {
        showCustomDateRangeDialog.value = false;
        return;
    }

    const changed = transactionsStore.updateTransactionListFilter({
        dateType: dateType,
        maxTime: maxTime,
        minTime: minTime
    });

    showCustomDateRangeDialog.value = false;
    updateUrlWhenChanged(changed);
}

function changeCustomMonthDateFilter(yearMonth: Year0BasedMonth): void {
    if (!yearMonth) {
        return;
    }

    const minTime = getYearMonthFirstUnixTime(yearMonth);
    const maxTime = getYearMonthLastUnixTime(yearMonth);
    const dateType = getDateTypeByDateRange(minTime, maxTime, firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.Normal);

    if (pageType.value === TransactionListPageType.Calendar.type) {
        currentCalendarDate.value = getValidMonthDayOrCurrentDayShortDate(minTime, currentCalendarDate.value);
    }

    if (query.value.dateType === dateType && query.value.maxTime === maxTime && query.value.minTime === minTime) {
        showCustomMonthDialog.value = false;
        return;
    }

    const changed = transactionsStore.updateTransactionListFilter({
        dateType: dateType,
        maxTime: maxTime,
        minTime: minTime
    });

    showCustomMonthDialog.value = false;
    updateUrlWhenChanged(changed);
}

function shiftDateRange(startTime: number, endTime: number, scale: number): void {
    if (recentMonthDateRanges.value[recentDateRangeIndex.value].dateType === DateRange.All.type) {
        return;
    }

    let newDateRange: TimeRangeAndDateType | null = null;

    if (DateRange.isBillingCycle(query.value.dateType) || query.value.dateType === DateRange.Custom.type) {
        newDateRange = getShiftedDateRangeAndDateTypeForBillingCycle(startTime, endTime, scale, firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.Normal, accountsStore.getAccountStatementDate(query.value.accountIds));
    }

    if (!newDateRange) {
        newDateRange = getShiftedDateRangeAndDateType(startTime, endTime, scale, firstDayOfWeek.value, fiscalYearStart.value, DateRangeScene.Normal);
    }

    if (pageType.value === TransactionListPageType.Calendar.type) {
        currentCalendarDate.value = getValidMonthDayOrCurrentDayShortDate(newDateRange.minTime, currentCalendarDate.value);
        const fullMonthDateRange = getFullMonthDateRange(newDateRange.minTime, newDateRange.maxTime, firstDayOfWeek.value, fiscalYearStart.value);

        if (fullMonthDateRange) {
            newDateRange = fullMonthDateRange;
            currentCalendarDate.value = getValidMonthDayOrCurrentDayShortDate(newDateRange.minTime, currentCalendarDate.value);
        }
    }

    const changed = transactionsStore.updateTransactionListFilter({
        dateType: newDateRange.dateType,
        maxTime: newDateRange.maxTime,
        minTime: newDateRange.minTime
    });

    updateUrlWhenChanged(changed);
}

function changeTypeFilter(type: number): void {
    let newCategoryFilter: string | undefined = undefined;

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

    updateUrlWhenChanged(changed);
}

function changeCategoryFilter(categoryIds: string): void {
    categoryMenuState.value = false;

    if (query.value.categoryIds === categoryIds) {
        return;
    }

    const changed = transactionsStore.updateTransactionListFilter({
        categoryIds: categoryIds
    });

    updateUrlWhenChanged(changed);
}

function changeMultipleCategoriesFilter(changed: boolean): void {
    categoryMenuState.value = false;
    showFilterCategoryDialog.value = false;
    updateUrlWhenChanged(changed);
}

function changeAccountFilter(accountIds: string): void {
    if (query.value.accountIds === accountIds) {
        return;
    }

    const changed = transactionsStore.updateTransactionListFilter({
        accountIds: accountIds
    });

    updateUrlWhenChanged(changed);
}

function changeMultipleAccountsFilter(changed: boolean): void {
    showFilterAccountDialog.value = false;
    updateUrlWhenChanged(changed);
}

function changeTagFilter(tagIds: string): void {
    if (query.value.tagIds === tagIds) {
        return;
    }

    const changed = transactionsStore.updateTransactionListFilter({
        tagIds: tagIds
    });

    updateUrlWhenChanged(changed);
}

function changeMultipleTagsFilter(changed: boolean): void {
    showFilterTagDialog.value = false;

    updateUrlWhenChanged(changed);
}

function changeTagFilterType(filterType: number): void {
    if (query.value.tagFilterType === filterType) {
        return;
    }

    const changed = transactionsStore.updateTransactionListFilter({
        tagFilterType: filterType
    });

    updateUrlWhenChanged(changed);
}

function changeKeywordFilter(keyword: string): void {
    if (query.value.keyword === keyword) {
        return;
    }

    const changed = transactionsStore.updateTransactionListFilter({
        keyword: keyword
    });

    updateUrlWhenChanged(changed);
}

function changeAmountFilter(filterType: string): void {
    currentAmountFilterType.value = '';
    amountMenuState.value = false;

    if (query.value.amountFilter === filterType) {
        return;
    }

    let amountFilter = filterType;

    if (filterType) {
        const amountCount = getAmountFilterParameterCount(filterType);

        if (!amountCount) {
            return;
        }

        if (amountCount === 1) {
            amountFilter += ':' + currentAmountFilterValue1.value;
        } else if (amountCount === 2) {
            if (currentAmountFilterValue2.value < currentAmountFilterValue1.value) {
                snackbar.value?.showMessage('Incorrect amount range');
                return;
            }

            amountFilter += ':' + currentAmountFilterValue1.value + ':' + currentAmountFilterValue2.value;
        } else {
            return;
        }
    }

    if (query.value.amountFilter === amountFilter) {
        return;
    }

    const changed = transactionsStore.updateTransactionListFilter({
        amountFilter: amountFilter
    });

    updateUrlWhenChanged(changed);
}

function add(template?: TransactionTemplate): void {
    const currentUnixTime = getCurrentUnixTime();

    let newTransactionTime: number | undefined = undefined;

    if (query.value.maxTime && query.value.minTime) {
        if (query.value.maxTime < currentUnixTime) {
            newTransactionTime = query.value.maxTime;
        } else if (currentUnixTime < query.value.minTime) {
            newTransactionTime = query.value.minTime;
        }
    }

    editDialog.value?.open({
        time: newTransactionTime,
        type: query.value.type,
        categoryId: queryAllFilterCategoryIdsCount.value === 1 ? query.value.categoryIds : '',
        accountId: queryAllFilterAccountIdsCount.value === 1 ? query.value.accountIds : '',
        tagIds: query.value.tagIds || '',
        template: template
    }).then(result => {
        if (result && result.message) {
            snackbar.value?.showMessage(result.message);
        }

        reload(false, false);
    }).catch(error => {
        if (error) {
            snackbar.value?.showError(error);
        }
    });
}

function importTransaction(): void {
    importDialog.value?.open().then(() => {
        reload(false, false);
    }).catch(error => {
        if (error) {
            snackbar.value?.showError(error);
        }
    });
}

function exportTransactions(fileExtension: string): void {
    if (exportingData.value) {
        return;
    }

    const nickname = userStore.currentUserNickname;
    let exportFileName = '';

    if (nickname) {
        exportFileName = tt('dataExport.exportFilename', {
            nickname: nickname
        }) + '.' + fileExtension;
    } else {
        exportFileName = tt('dataExport.defaultExportFilename') + '.' + fileExtension;
    }

    const exportTransactionReq = transactionsStore.getExportTransactionDataRequestByTransactionFilter();

    exportingData.value = true;

    userStore.getExportedUserData(fileExtension, exportTransactionReq).then(data => {
        startDownloadFile(exportFileName, data);
        exportingData.value = false;
    }).catch(error => {
        exportingData.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function show(transaction: Transaction): void {
    editDialog.value?.open({
        id: transaction.id,
        currentTransaction: transaction
    }).then(result => {
        if (result && result.message) {
            snackbar.value?.showMessage(result.message);
        }

        reload(false, false);
    }).catch(error => {
        if (error) {
            snackbar.value?.showError(error);
        }
    });
}

function scrollTimeMenuToSelectedItem(opened: boolean): void {
    if (opened) {
        scrollMenuToSelectedItem(timeFilterMenu.value);
    }
}

function scrollCategoryMenuToSelectedItem(opened: boolean): void {
    if (opened) {
        scrollMenuToSelectedItem(categoryFilterMenu.value);
    }
}

function scrollAmountMenuToSelectedItem(opened: boolean): void {
    if (opened) {
        currentAmountFilterType.value = '';

        let amount1 = 0, amount2 = 0;

        if (isString(query.value.amountFilter)) {
            try {
                const filterItems = query.value.amountFilter.split(':');
                const amountCount = getAmountFilterParameterCount(filterItems[0]);

                if (filterItems.length === 2 && amountCount === 1) {
                    amount1 = parseInt(filterItems[1]);
                } else if (filterItems.length === 3 && amountCount === 2) {
                    amount1 = parseInt(filterItems[1]);
                    amount2 = parseInt(filterItems[2]);
                }
            } catch (ex) {
                logger.warn('cannot parse amount from filter value, original value is ' + query.value.amountFilter, ex);
            }
        }

        currentAmountFilterValue1.value = amount1;
        currentAmountFilterValue2.value = amount2;

        scrollMenuToSelectedItem(amountFilterMenu.value);
    }
}

function scrollAccountMenuToSelectedItem(opened: boolean): void {
    if (opened) {
        scrollMenuToSelectedItem(accountFilterMenu.value);
    }
}

function scrollTagMenuToSelectedItem(opened: boolean): void {
    if (opened) {
        scrollMenuToSelectedItem(tagFilterMenu.value);
    }
}

function scrollMenuToSelectedItem(menu: VMenu | null): void {
    nextTick(() => {
        scrollToSelectedItem(menu?.contentEl, 'div.v-list', 'div.v-list-item.list-item-selected');
    });
}

function onShowDateRangeError(message: string): void {
    snackbar.value?.showError(message);
}

onBeforeRouteUpdate((to) => {
    if (to.query) {
        init({
            initDateType: (to.query['dateType'] as string | null) || undefined,
            initMinTime: (to.query['minTime'] as string | null) || undefined,
            initMaxTime: (to.query['maxTime'] as string | null) || undefined,
            initType: (to.query['type'] as string | null) || undefined,
            initCategoryIds: (to.query['categoryIds'] as string | null) || undefined,
            initAccountIds: (to.query['accountIds'] as string | null) || undefined,
            initTagIds: (to.query['tagIds'] as string | null) || undefined,
            initTagFilterType: (to.query['tagFilterType'] as string | null) || undefined,
            initAmountFilter: (to.query['amountFilter'] as string | null) || undefined,
            initKeyword: (to.query['keyword'] as string | null) || undefined
        });
    } else {
        init({});
    }
});

watch(() => display.mdAndUp.value, (newValue) => {
    alwaysShowNav.value = newValue;

    if (!showNav.value) {
        showNav.value = newValue;
    }
});

watch(() => desktopPageStore.showAddTransactionDialogInTransactionList, (newValue) => {
    if (newValue) {
        desktopPageStore.resetShowAddTransactionDialogInTransactionList();
        add();
    }
});

init(props);
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

.transaction-list-custom-datetime-range {
    line-height: 1rem;
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
    margin-inline-start: 0in;
}

.transaction-table .transaction-table-column-tags .v-chip.transaction-tag {
    margin-inline-end: 4px;
    margin-top: 2px;
    margin-bottom: 2px;
}

.transaction-table .transaction-table-column-tags .v-chip.transaction-tag > .v-chip__content {
    display: block;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
}

.transaction-time-menu .item-icon,
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

.transaction-calendar-container .dp__main .dp__menu {
    --dp-border-radius: 6px;
    --dp-menu-border-color: rgba(var(--v-border-color), var(--v-border-opacity));
}

.transaction-calendar-container .dp__main .dp__calendar {
    --dp-border-color: rgba(var(--v-border-color), var(--v-border-opacity));
}

.transaction-calendar-container .dp__main .dp__calendar .dp__calendar_row {
    --dp-cell-size: 80px;
    --dp-primary-color: rgba(var(--v-theme-primary), var(--v-activated-opacity));
    --dp-primary-text-color: rgb(var(--v-theme-primary));
}

.transaction-calendar-container .dp__main .dp__calendar .dp__calendar_row > .dp__calendar_item {
    overflow: hidden;
}

.transaction-calendar-container .dp__main .dp__calendar .dp__calendar_row > .dp__calendar_item .transaction-calendar-daily-amounts > span {
    display: block;
    width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}
</style>
