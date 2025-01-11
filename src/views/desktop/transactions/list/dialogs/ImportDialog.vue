<template>
    <v-dialog :persistent="!!persistent" v-model="showState">
        <v-card class="pa-6 pa-sm-10 pa-md-12">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <div class="d-flex w-100 align-center justify-center">
                        <h4 class="text-h4">{{ $t('Import Transactions') }}</h4>
                        <v-progress-circular indeterminate size="22" class="ml-2" v-if="loading"></v-progress-circular>
                    </div>
                    <v-btn density="comfortable" color="default" variant="text" class="ml-2" :icon="true"
                           v-if="currentStep === 'checkData'">
                        <v-icon :icon="icons.filter" />
                        <v-menu activator="parent" max-height="500">
                            <v-list>
                                <v-list-subheader :title="$t('Date Range')"/>
                                <v-list-item :title="$t('All')"
                                             :append-icon="filters.minDatetime === null || filters.maxDatetime === null ? icons.checked : null"
                                             @click="filters.minDatetime = filters.maxDatetime = null"></v-list-item>
                                <v-list-item :title="$t('Custom')"
                                             :subtitle="displayFilterCustomDateRange"
                                             :append-icon="filters.minDatetime !== null && filters.maxDatetime !== null ? icons.checked : null"
                                             @click="showCustomDateRangeDialog = true"></v-list-item>
                                <v-divider class="my-2"/>
                                <v-list-subheader :title="$t('Type')"/>
                                <v-list-item :title="$t('All')"
                                             :append-icon="filters.transactionType === null ? icons.checked : null"
                                             @click="filters.transactionType = null"></v-list-item>
                                <v-list-item :title="$t('Income')"
                                             :append-icon="filters.transactionType === allTransactionTypes.Income ? icons.checked : null"
                                             @click="filters.transactionType = allTransactionTypes.Income"></v-list-item>
                                <v-list-item :title="$t('Expense')"
                                             :append-icon="filters.transactionType === allTransactionTypes.Expense ? icons.checked : null"
                                             @click="filters.transactionType = allTransactionTypes.Expense"></v-list-item>
                                <v-list-item :title="$t('Transfer')"
                                             :append-icon="filters.transactionType === allTransactionTypes.Transfer ? icons.checked : null"
                                             @click="filters.transactionType = allTransactionTypes.Transfer"></v-list-item>
                                <v-divider class="my-2"/>
                                <v-list-subheader :title="$t('Category')"/>
                                <v-list-item :title="$t('All')"
                                             :append-icon="filters.category === null ? icons.checked : null"
                                             @click="filters.category = null"></v-list-item>
                                <v-list-item :title="$t('Invalid Category')"
                                             :append-icon="filters.category === undefined ? icons.checked : null"
                                             @click="filters.category = undefined"></v-list-item>
                                <v-list-item :title="$t('None')"
                                             :append-icon="filters.category === '' ? icons.checked : null"
                                             @click="filters.category = ''"></v-list-item>
                                <v-list-item :title="name" :key="name"
                                             :append-icon="filters.category === name ? icons.checked : null"
                                             v-for="name in allUsedCategoryNames"
                                             @click="filters.category = name"></v-list-item>
                                <v-divider class="my-2"/>
                                <v-list-subheader :title="$t('Account')"/>
                                <v-list-item :title="$t('All')"
                                             :append-icon="filters.account === null ? icons.checked : null"
                                             @click="filters.account = null"></v-list-item>
                                <v-list-item :title="$t('Invalid Account')"
                                             :append-icon="filters.account === undefined ? icons.checked : null"
                                             @click="filters.account = undefined"></v-list-item>
                                <v-list-item :title="$t('None')"
                                             :append-icon="filters.account === '' ? icons.checked : null"
                                             @click="filters.account = ''"></v-list-item>
                                <v-list-item :title="name" :key="name"
                                             :append-icon="filters.account === name ? icons.checked : null"
                                             v-for="name in allUsedAccountNames"
                                             @click="filters.account = name"></v-list-item>
                                <v-divider class="my-2"/>
                                <v-list-subheader :title="$t('Tags')"/>
                                <v-list-item :title="$t('All')"
                                             :append-icon="filters.tag === null ? icons.checked : null"
                                             @click="filters.tag = null"></v-list-item>
                                <v-list-item :title="$t('Invalid Tag')"
                                             :append-icon="filters.tag === undefined ? icons.checked : null"
                                             @click="filters.tag = undefined"></v-list-item>
                                <v-list-item :title="$t('None')"
                                             :append-icon="filters.tag === '' ? icons.checked : null"
                                             @click="filters.tag = ''"></v-list-item>
                                <v-list-item :title="name" :key="name"
                                             :append-icon="filters.tag === name ? icons.checked : null"
                                             v-for="name in allUsedTagNames"
                                             @click="filters.tag = name"></v-list-item>
                                <v-divider class="my-2"/>
                                <v-list-subheader :title="$t('Description')"/>
                                <v-list-item :title="$t('All')"
                                             :append-icon="filters.description === null ? icons.checked : null"
                                             @click="filters.description = null"></v-list-item>
                                <v-list-item :title="$t('None')"
                                             :append-icon="filters.description === '' ? icons.checked : null"
                                             @click="filters.description = ''"></v-list-item>
                                <v-list-item :title="$t('Custom')"
                                             :subtitle="filters.description"
                                             :append-icon="filters.description !== null && filters.description !== '' ? icons.checked : null"
                                             @click="currentDescriptionFilterValue = filters.description || ''; showCustomDescriptionDialog = true"></v-list-item>
                            </v-list>
                        </v-menu>
                    </v-btn>
                    <v-btn density="comfortable" color="default" variant="text" class="ml-2" :icon="true"
                           v-if="currentStep === 'checkData'">
                        <v-icon :icon="icons.more" />
                        <v-menu activator="parent">
                            <v-list>
                                <v-list-item :prepend-icon="icons.replace"
                                             :disabled="editingTransaction || selectedExpenseTransactionCount < 1"
                                             :title="$t('Batch Replace Selected Expense Categories')"
                                             @click="showBatchReplaceDialog('expenseCategory')"></v-list-item>
                                <v-list-item :prepend-icon="icons.replace"
                                             :disabled="editingTransaction || selectedIncomeTransactionCount < 1"
                                             :title="$t('Batch Replace Selected Income Categories')"
                                             @click="showBatchReplaceDialog('incomeCategory')"></v-list-item>
                                <v-list-item :prepend-icon="icons.replace"
                                             :disabled="editingTransaction || selectedTransferTransactionCount < 1"
                                             :title="$t('Batch Replace Selected Transfer Categories')"
                                             @click="showBatchReplaceDialog('transferCategory')"></v-list-item>
                                <v-list-item :prepend-icon="icons.replace"
                                             :disabled="editingTransaction || selectedImportTransactionCount < 1"
                                             :title="$t('Batch Replace Selected Accounts')"
                                             @click="showBatchReplaceDialog('account')"></v-list-item>
                                <v-list-item :prepend-icon="icons.replace"
                                             :disabled="editingTransaction || selectedTransferTransactionCount < 1"
                                             :title="$t('Batch Replace Selected Destination Accounts')"
                                             @click="showBatchReplaceDialog('destinationAccount')"></v-list-item>
                                <v-divider class="my-2"/>
                                <v-list-item :prepend-icon="icons.replace"
                                             :disabled="editingTransaction || allInvalidExpenseCategoryNames < 1"
                                             :title="$t('Replace Invalid Expense Categories')"
                                             @click="showReplaceInvalidItemDialog('expenseCategory', allInvalidExpenseCategoryNames)"></v-list-item>
                                <v-list-item :prepend-icon="icons.replace"
                                             :disabled="editingTransaction || allInvalidIncomeCategoryNames < 1"
                                             :title="$t('Replace Invalid Income Categories')"
                                             @click="showReplaceInvalidItemDialog('incomeCategory', allInvalidIncomeCategoryNames)"></v-list-item>
                                <v-list-item :prepend-icon="icons.replace"
                                             :disabled="editingTransaction || allInvalidTransferCategoryNames < 1"
                                             :title="$t('Replace Invalid Transfer Categories')"
                                             @click="showReplaceInvalidItemDialog('transferCategory', allInvalidTransferCategoryNames)"></v-list-item>
                                <v-list-item :prepend-icon="icons.replace"
                                             :disabled="editingTransaction || allInvalidAccountNames < 1"
                                             :title="$t('Replace Invalid Accounts')"
                                             @click="showReplaceInvalidItemDialog('account', allInvalidAccountNames)"></v-list-item>
                                <v-list-item :prepend-icon="icons.replace"
                                             :disabled="editingTransaction || allInvalidTransactionTagNames < 1"
                                             :title="$t('Replace Invalid Transaction Tags')"
                                             @click="showReplaceInvalidItemDialog('tag', allInvalidTransactionTagNames)"></v-list-item>
                            </v-list>
                        </v-menu>
                    </v-btn>
                </div>
            </template>

            <div class="mt-4 cursor-default">
                <steps-bar min-width="700" :clickable="false" :steps="allSteps" :current-step="currentStep" />
            </div>

            <v-window class="disable-tab-transition" v-model="currentStep">
                <v-window-item value="uploadFile">
                    <v-row>
                        <v-col cols="12" md="12">
                            <v-select
                                item-title="displayName"
                                item-value="type"
                                :disabled="submitting"
                                :label="$t('File Type')"
                                :placeholder="$t('File Type')"
                                :items="allSupportedImportFileTypes"
                                v-model="fileType"
                            />
                        </v-col>

                        <v-col cols="12" md="12" v-if="allFileSubTypes">
                            <v-select
                                item-title="displayName"
                                item-value="type"
                                :disabled="submitting"
                                :label="$t('Format')"
                                :placeholder="$t('Format')"
                                :items="allFileSubTypes"
                                v-model="fileSubType"
                            />
                        </v-col>

                        <v-col cols="12" md="12">
                            <v-text-field
                                readonly
                                persistent-placeholder
                                type="text"
                                class="always-cursor-pointer"
                                :disabled="submitting"
                                :label="$t('Data File')"
                                :placeholder="$t('format.misc.clickToSelectedFile', { extensions: supportedImportFileExtensions })"
                                v-model="fileName"
                                @click="showOpenFileDialog"
                            />
                        </v-col>

                        <v-col cols="12" md="12" class="mb-0 pb-0" v-if="exportFileGuideDocumentUrl">
                            <a :href="exportFileGuideDocumentUrl" :class="{ 'disabled': submitting }" target="_blank">
                                <v-icon :icon="icons.document" size="16" />
                                <span class="ml-1">{{ $t('How to export this file?') }}</span>
                                <span class="ml-1" v-if="exportFileGuideDocumentLanguageName">({{ exportFileGuideDocumentLanguageName }})</span>
                            </a>
                        </v-col>
                    </v-row>
                </v-window-item>
                <v-window-item value="checkData">
                    <v-data-table
                        fixed-header
                        fixed-footer
                        show-select
                        multi-sort
                        class="import-transaction-table"
                        density="compact"
                        item-value="index"
                        :height="importTransactionsTableHeight"
                        :headers="importTransactionHeaders"
                        :items="importTransactions"
                        :search="JSON.stringify(filters)"
                        :custom-filter="importTransactionsFilter"
                        :no-data-text="$t('No data to import')"
                        v-model:items-per-page="countPerPage"
                        v-model:page="currentPage"
                    >
                        <template #header.data-table-select>
                            <v-checkbox readonly class="always-cursor-pointer"
                                        density="compact" width="28"
                                        :disabled="loading || submitting"
                                        :indeterminate="anyButNotAllTransactionSelected"
                                        v-model="allTransactionSelected"
                            >
                                <v-menu activator="parent" location="bottom">
                                    <v-list>
                                        <v-list-item :prepend-icon="icons.selectAll"
                                                     :title="$t('Select All Valid Items')"
                                                     :disabled="loading || submitting"
                                                     @click="selectAllValid"></v-list-item>
                                        <v-list-item :prepend-icon="icons.selectAll"
                                                     :title="$t('Select All Invalid Items')"
                                                     :disabled="loading || submitting"
                                                     @click="selectAllInvalid"></v-list-item>
                                        <v-divider class="my-2"/>
                                        <v-list-item :prepend-icon="icons.selectAll"
                                                     :title="$t('Select All')"
                                                     :disabled="loading || submitting"
                                                     @click="selectAll"></v-list-item>
                                        <v-list-item :prepend-icon="icons.selectNone"
                                                     :title="$t('Select None')"
                                                     :disabled="loading || submitting"
                                                     @click="selectNone"></v-list-item>
                                        <v-list-item :prepend-icon="icons.selectInverse"
                                                     :title="$t('Invert Selection')"
                                                     :disabled="loading || submitting"
                                                     @click="selectInvert"></v-list-item>
                                        <v-divider class="my-2"/>
                                        <v-list-item :prepend-icon="icons.selectAll"
                                                     :title="$t('Select All in This Page')"
                                                     :disabled="loading || submitting"
                                                     @click="selectAllInThisPage"></v-list-item>
                                        <v-list-item :prepend-icon="icons.selectNone"
                                                     :title="$t('Select None in This Page')"
                                                     :disabled="loading || submitting"
                                                     @click="selectNoneInThisPage"></v-list-item>
                                        <v-list-item :prepend-icon="icons.selectInverse"
                                                     :title="$t('Invert Selection in This Page')"
                                                     :disabled="loading || submitting"
                                                     @click="selectInvertInThisPage"></v-list-item>
                                    </v-list>
                                </v-menu>
                            </v-checkbox>
                        </template>
                        <template #item.data-table-select="{ item }">
                            <v-checkbox density="compact"
                                        :disabled="loading || submitting"
                                        v-model="item.selected"></v-checkbox>
                        </template>
                        <template #item.valid="{ item }">
                            <v-icon size="small" :class="{ 'text-error': !item.valid }"
                                    :disabled="loading || submitting"
                                    :icon="editingTransaction === item ? icons.complete : icons.edit"
                                    @click="editTransaction(item)">
                            </v-icon>
                            <v-tooltip activator="parent" v-if="!loading && !submitting">{{ $t('Edit') }}</v-tooltip>
                        </template>
                        <template #item.time="{ item }">
                            <span>{{ getDisplayDateTime(item) }}</span>
                            <v-chip class="ml-1" variant="flat" color="secondary" size="x-small"
                                    v-if="item.utcOffset !== currentTimezoneOffsetMinutes">{{ getDisplayTimezone(item) }}</v-chip>
                        </template>
                        <template #item.type="{ value }">
                            <v-chip label color="secondary" variant="outlined" size="x-small" v-if="value === allTransactionTypes.ModifyBalance">{{ $t('Modify Balance') }}</v-chip>
                            <v-chip label class="text-income" variant="outlined" size="x-small" v-else-if="value === allTransactionTypes.Income">{{ $t('Income') }}</v-chip>
                            <v-chip label class="text-expense" variant="outlined" size="x-small" v-else-if="value === allTransactionTypes.Expense">{{ $t('Expense') }}</v-chip>
                            <v-chip label color="primary" variant="outlined" size="x-small" v-else-if="value === allTransactionTypes.Transfer">{{ $t('Transfer') }}</v-chip>
                            <v-chip label color="default" variant="outlined" size="x-small" v-else>{{ $t('Unknown') }}</v-chip>
                        </template>
                        <template #item.actualCategoryName="{ item }">
                            <div class="d-flex align-center" v-if="editingTransaction !== item || item.type === allTransactionTypes.ModifyBalance">
                                <span v-if="item.type === allTransactionTypes.ModifyBalance">-</span>
                                <ItemIcon size="24px" icon-type="category"
                                          :icon-id="allCategoriesMap[item.categoryId].icon"
                                          :color="allCategoriesMap[item.categoryId].color"
                                          v-if="item.type !== allTransactionTypes.ModifyBalance && item.categoryId && item.categoryId !== '0' && allCategoriesMap[item.categoryId]"></ItemIcon>
                                <span class="ml-2" v-if="item.type !== allTransactionTypes.ModifyBalance && item.categoryId && item.categoryId !== '0' && allCategoriesMap[item.categoryId]">
                                    {{ allCategoriesMap[item.categoryId].name }}
                                </span>
                                <div class="text-error font-italic" v-else-if="item.type !== allTransactionTypes.ModifyBalance && (!item.categoryId || item.categoryId === '0' || !allCategoriesMap[item.categoryId])">
                                    <v-icon class="mr-1" :icon="icons.alert"/>
                                    <span>{{ item.originalCategoryName }}</span>
                                </div>
                            </div>
                            <div style="width: 260px" v-if="editingTransaction === item && item.type === allTransactionTypes.Expense">
                                <two-column-select density="compact" variant="plain"
                                                   primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                   primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                   primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                   secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                   secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                   secondary-hidden-field="hidden"
                                                   :disabled="loading || submitting || !hasAvailableExpenseCategories"
                                                   :show-selection-primary-text="true"
                                                   :custom-selection-primary-text="getPrimaryCategoryName(item.categoryId, allCategories[allCategoryTypes.Expense])"
                                                   :custom-selection-secondary-text="getSecondaryCategoryName(item.categoryId, allCategories[allCategoryTypes.Expense])"
                                                   :placeholder="$t('Category')"
                                                   :items="allCategories[allCategoryTypes.Expense]"
                                                   v-model="item.categoryId">
                                </two-column-select>
                            </div>
                            <div style="width: 260px" v-if="editingTransaction === item && item.type === allTransactionTypes.Income">
                                <two-column-select density="compact" variant="plain"
                                                   primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                   primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                   primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                   secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                   secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                   secondary-hidden-field="hidden"
                                                   :disabled="loading || submitting || !hasAvailableIncomeCategories"
                                                   :show-selection-primary-text="true"
                                                   :custom-selection-primary-text="getPrimaryCategoryName(item.categoryId, allCategories[allCategoryTypes.Income])"
                                                   :custom-selection-secondary-text="getSecondaryCategoryName(item.categoryId, allCategories[allCategoryTypes.Income])"
                                                   :placeholder="$t('Category')"
                                                   :items="allCategories[allCategoryTypes.Income]"
                                                   v-model="item.categoryId">
                                </two-column-select>
                            </div>
                            <div style="width: 260px" v-if="editingTransaction === item && item.type === allTransactionTypes.Transfer">
                                <two-column-select density="compact" variant="plain"
                                                   primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                   primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                   primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                   secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                   secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                   secondary-hidden-field="hidden"
                                                   :disabled="loading || submitting || !hasAvailableTransferCategories"
                                                   :show-selection-primary-text="true"
                                                   :custom-selection-primary-text="getPrimaryCategoryName(item.categoryId, allCategories[allCategoryTypes.Transfer])"
                                                   :custom-selection-secondary-text="getSecondaryCategoryName(item.categoryId, allCategories[allCategoryTypes.Transfer])"
                                                   :placeholder="$t('Category')"
                                                   :items="allCategories[allCategoryTypes.Transfer]"
                                                   v-model="item.categoryId">
                                </two-column-select>
                            </div>
                        </template>
                        <template #item.sourceAmount="{ item }">
                            <span>{{ getTransactionDisplayAmount(item) }}</span>
                            <v-icon class="mx-1" size="13" :icon="icons.arrowRight" v-if="item.type === allTransactionTypes.Transfer && item.sourceAccountId !== item.destinationAccountId"></v-icon>
                            <span v-if="item.type === allTransactionTypes.Transfer && item.sourceAccountId !== item.destinationAccountId">{{ getTransactionDisplayDestinationAmount(item) }}</span>
                        </template>
                        <template #item.actualSourceAccountName="{ item }">
                            <div class="d-flex align-center" v-if="editingTransaction !== item">
                                <span v-if="item.sourceAccountId && item.sourceAccountId !== '0' && allAccountsMap[item.sourceAccountId]">{{ allAccountsMap[item.sourceAccountId].name }}</span>
                                <div class="text-error font-italic" v-else>
                                    <v-icon class="mr-1" :icon="icons.alert"/>
                                    <span>{{ item.originalSourceAccountName }}</span>
                                </div>
                                <v-icon class="mx-1" size="13" :icon="icons.arrowRight" v-if="item.type === allTransactionTypes.Transfer"></v-icon>
                                <span v-if="item.type === allTransactionTypes.Transfer && item.destinationAccountId && item.destinationAccountId !== '0' && allAccountsMap[item.destinationAccountId]">{{allAccountsMap[item.destinationAccountId].name }}</span>
                                <div class="text-error font-italic" v-else-if="item.type === allTransactionTypes.Transfer && (!item.destinationAccountId || item.destinationAccountId === '0' || !allAccountsMap[item.destinationAccountId])">
                                    <v-icon class="mr-1" :icon="icons.alert"/>
                                    <span>{{ item.originalDestinationAccountName }}</span>
                                </div>
                            </div>
                            <div class="d-flex align-center" style="width: 200px" v-if="editingTransaction === item">
                                <two-column-select density="compact" variant="plain"
                                                   primary-key-field="id" primary-value-field="category"
                                                   primary-title-field="name" primary-footer-field="displayBalance"
                                                   primary-icon-field="icon" primary-icon-type="account"
                                                   primary-sub-items-field="accounts"
                                                   :primary-title-i18n="true"
                                                   secondary-key-field="id" secondary-value-field="id"
                                                   secondary-title-field="name" secondary-footer-field="displayBalance"
                                                   secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                   :disabled="loading || submitting || !allVisibleAccounts.length"
                                                   :custom-selection-primary-text="getSourceAccountDisplayName(item)"
                                                   :placeholder="getSourceAccountTitle(item)"
                                                   :items="allVisibleCategorizedAccounts"
                                                   v-model="item.sourceAccountId">
                                </two-column-select>
                                <v-icon class="mx-1" size="13" :icon="icons.arrowRight" v-if="item.type === allTransactionTypes.Transfer"></v-icon>
                                <two-column-select density="compact" variant="plain"
                                                   primary-key-field="id" primary-value-field="category"
                                                   primary-title-field="name" primary-footer-field="displayBalance"
                                                   primary-icon-field="icon" primary-icon-type="account"
                                                   primary-sub-items-field="accounts"
                                                   :primary-title-i18n="true"
                                                   secondary-key-field="id" secondary-value-field="id"
                                                   secondary-title-field="name" secondary-footer-field="displayBalance"
                                                   secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                   :disabled="loading || submitting || !allVisibleAccounts.length"
                                                   :custom-selection-primary-text="getDestinationAccountDisplayName(item)"
                                                   :placeholder="$t('Destination Account')"
                                                   :items="allVisibleCategorizedAccounts"
                                                   v-model="item.destinationAccountId"
                                                   v-if="item.type === allTransactionTypes.Transfer">
                                </two-column-select>
                            </div>
                        </template>
                        <template #item.geoLocation="{ item }">
                            <span v-if="item.geoLocation">{{ `(${item.geoLocation.longitude}, ${item.geoLocation.latitude})` }}</span>
                            <span v-else-if="!item.geoLocation">{{ $t('None') }}</span>
                        </template>
                        <template #item.tagIds="{ item }">
                            <div v-if="editingTransaction !== item">
                                <v-chip class="transaction-tag" size="small"
                                        :class="{ 'font-italic': !tagId || tagId === '0' || !allTagsMap[tagId] }"
                                        :prepend-icon="tagId && tagId !== '0' && allTagsMap[tagId] ? icons.tag : icons.alert"
                                        :color="tagId && tagId !== '0' && allTagsMap[tagId] ? 'default' : 'error'"
                                        :text="tagId && tagId !== '0' && allTagsMap[tagId] ? allTagsMap[tagId].name : item.originalTagNames[index]"
                                        :key="tagId"
                                        v-for="(tagId, index) in item.tagIds"/>
                                <v-chip class="transaction-tag" size="small"
                                        :text="$t('None')"
                                        v-if="!item.tagIds || !item.tagIds.length"/>
                            </div>
                            <div style="width: 200px" v-if="editingTransaction === item">
                                <v-autocomplete
                                    item-title="name"
                                    item-value="id"
                                    auto-select-first
                                    persistent-placeholder
                                    multiple
                                    chips
                                    closable-chips
                                    density="compact" variant="plain"
                                    :disabled="loading || submitting"
                                    :placeholder="$t('None')"
                                    :items="allTags"
                                    :no-data-text="$t('No available tag')"
                                    v-model="editingTags"
                                >
                                    <template #chip="{ props, index }">
                                        <v-chip :class="{ 'font-italic': !isTagValid(editingTags, index) }"
                                                :prepend-icon="isTagValid(editingTags, index) ? icons.tag : icons.alert"
                                                :color="isTagValid(editingTags, index) ? 'default' : 'error'"
                                                :text="isTagValid(editingTags, index) ? allTagsMap[editingTags[index]].name : item.originalTagNames[index]"
                                                v-bind="props"/>
                                    </template>

                                    <template #item="{ props, item }">
                                        <v-list-item :value="item.value" v-bind="props" v-if="!item.raw.hidden">
                                            <template #title>
                                                <v-list-item-title>
                                                    <div class="d-flex align-center">
                                                        <v-icon size="20" start :icon="icons.tag"/>
                                                        <span>{{ item.title }}</span>
                                                    </div>
                                                </v-list-item-title>
                                            </template>
                                        </v-list-item>
                                    </template>
                                </v-autocomplete>
                            </div>
                        </template>
                        <template #bottom>
                            <div class="title-and-toolbar d-flex align-center text-no-wrap mt-2"
                                 v-if="importTransactions && importTransactions.length > 10">
                                <span :class="{ 'text-error': selectedInvalidTransactionCount > 0 }">
                                    {{ $t('format.misc.selectedCount', { count: selectedImportTransactionCount, totalCount: importTransactions.length }) }}
                                </span>
                                <v-spacer/>
                                <span>{{ $t('Transactions Per Page') }}</span>
                                <v-select class="ml-2" density="compact" max-width="100"
                                          item-title="title"
                                          item-value="value"
                                          :disabled="loading"
                                          :items="importTransactionsTablePageOptions"
                                          v-model="countPerPage"
                                />
                                <v-pagination density="compact"
                                              :total-visible="6"
                                              :length="totalPageCount"
                                              v-model="currentPage">
                                    <template #item="{ page, isActive }">
                                        <v-btn density="compact"
                                               variant="text"
                                               :icon="true"
                                               :color="isActive ? 'primary' : 'default'"
                                               @click="currentPage = parseInt(page)"
                                               v-if="page !== '...'"
                                        >
                                            <span>{{ page }}</span>
                                        </v-btn>
                                        <v-btn density="compact"
                                               variant="text"
                                               color="default"
                                               :icon="true"
                                               v-if="page === '...'"
                                        >
                                            <span>{{ page }}</span>
                                            <v-menu :close-on-content-click="false" activator="parent">
                                                <v-list>
                                                    <v-list-item class="text-sm" density="compact">
                                                        <v-list-item-title class="cursor-pointer">
                                                            <v-autocomplete density="compact"
                                                                            width="100"
                                                                            item-title="page"
                                                                            item-value="page"
                                                                            :items="allPages"
                                                                            :no-data-text="$t('No results')"
                                                                            v-model="inputCurrentPage"
                                                            />
                                                        </v-list-item-title>
                                                    </v-list-item>
                                                </v-list>
                                            </v-menu>
                                        </v-btn>
                                    </template>
                                </v-pagination>
                            </div>
                        </template>
                    </v-data-table>
                </v-window-item>
                <v-window-item value="finalResult">
                    <h4 class="text-h4 mb-1">{{ $t('Data Import Completed') }}</h4>
                    <p class="my-5">{{ $t('format.misc.importTransactionResult', { count: importedCount }) }}</p>
                </v-window-item>
            </v-window>

            <div class="d-flex justify-sm-space-between gap-4 flex-wrap justify-center mt-10">
                <v-btn color="secondary" variant="tonal" :disabled="loading || submitting"
                       :prepend-icon="icons.previous" @click="cancel(false)"
                       v-if="currentStep !== 'finalResult'">{{ $t('Cancel') }}</v-btn>
                <v-btn color="primary" :disabled="loading || submitting || !importFile"
                       :append-icon="!submitting ? icons.next : null" @click="parseData"
                       v-if="currentStep === 'uploadFile'">
                    {{ $t('Next') }}
                    <v-progress-circular indeterminate size="22" class="ml-2" v-if="submitting"></v-progress-circular>
                </v-btn>
                <v-btn color="teal" :disabled="submitting || editingTransaction || selectedImportTransactionCount < 1 || selectedInvalidTransactionCount > 0"
                       :append-icon="!submitting ? icons.next : null" @click="submit"
                       v-if="currentStep === 'checkData'">
                    {{ $t('Import') }}
                    <v-progress-circular indeterminate size="22" class="ml-2" v-if="submitting"></v-progress-circular>
                </v-btn>
                <v-btn color="secondary" variant="tonal"
                       :append-icon="icons.complete"
                       @click="cancel(true)"
                       v-if="currentStep === 'finalResult'">{{ $t('Close') }}</v-btn>
            </div>
        </v-card>
    </v-dialog>

    <v-dialog width="640" v-model="showCustomDescriptionDialog">
        <v-card class="pa-2 pa-sm-4 pa-md-4">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <h4 class="text-h4">{{ $t('Filter Description') }}</h4>
                </div>
            </template>
            <v-card-text class="mb-md-4 w-100 d-flex justify-center">
                <v-text-field
                    type="text"
                    persistent-placeholder
                    :label="$t('Description')"
                    :placeholder="$t('Description')"
                    v-model="currentDescriptionFilterValue"
                />
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center gap-4">
                    <v-btn :disabled="!currentDescriptionFilterValue" @click="showCustomDescriptionDialog = false; filters.description = currentDescriptionFilterValue">{{ $t('OK') }}</v-btn>
                    <v-btn color="secondary" variant="tonal" @click="showCustomDescriptionDialog = false; currentDescriptionFilterValue = ''">{{ $t('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <date-range-selection-dialog :title="$t('Custom Date Range')"
                                 :min-time="filters.minDatetime"
                                 :max-time="filters.maxDatetime"
                                 v-model:show="showCustomDateRangeDialog"
                                 @dateRange:change="changeCustomDateFilter"
                                 @error="showError" />
    <batch-replace-dialog ref="batchReplaceDialog" />
    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
    <input ref="fileInput" type="file" style="display: none" :accept="supportedImportFileExtensions" @change="setImportFile($event)" />
</template>

<script>
import BatchReplaceDialog from './BatchReplaceDialog.vue';

import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.js';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';
import { useTransactionTagsStore } from '@/stores/transactionTag.js';
import { useTransactionsStore } from '@/stores/transaction.js';
import { useOverviewStore } from '@/stores/overview.ts';
import { useStatisticsStore } from '@/stores/statistics.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import { CategoryType } from '@/core/category.ts';
import { TransactionType } from '@/core/transaction.ts';
import {
    isString,
    isNumber,
    getNameByKeyValue,
    objectFieldToArrayItem
} from '@/lib/common.ts';
import { isFileExtensionSupported } from '@/lib/file.ts';
import { generateRandomUUID } from '@/lib/misc.ts';
import logger from '@/lib/logger.ts';
import {
    parseDateFromUnixTime,
    getUnixTime,
    getUtcOffsetByUtcOffsetMinutes,
    getTimezoneOffsetMinutes
} from '@/lib/datetime.ts';
import {
    getTransactionPrimaryCategoryName,
    getTransactionSecondaryCategoryName,
    getFirstAvailableCategoryId
} from '@/lib/category.js';

import {
    mdiFilterOutline,
    mdiCheck,
    mdiDotsVertical,
    mdiHelpCircleOutline,
    mdiFindReplace,
    mdiClose,
    mdiArrowRight,
    mdiSelectAll,
    mdiSelect,
    mdiSelectInverse,
    mdiPencilOutline,
    mdiAlertOutline,
    mdiPound
} from '@mdi/js';

export default {
    components: {
        BatchReplaceDialog
    },
    props: [
        'persistent'
    ],
    expose: [
        'open'
    ],
    data() {
        return {
            showState: false,
            clientSessionId: '',
            currentStep: 'uploadFile',
            fileType: 'ezbookkeeping',
            fileSubType: 'ezbookkeeping_csv',
            importFile: null,
            importTransactions: null,
            editingTransaction: null,
            editingTags: [],
            filters: {
                minDatetime: null,
                maxDatetime: null,
                transactionType: null,
                category: null,
                account: null,
                tag: null,
                description: null
            },
            currentPage: 1,
            countPerPage: 10,
            importedCount: null,
            showCustomDateRangeDialog: false,
            showCustomDescriptionDialog: false,
            currentDescriptionFilterValue: null,
            loading: true,
            submitting: false,
            resolve: null,
            reject: null,
            icons: {
                filter: mdiFilterOutline,
                checked: mdiCheck,
                more: mdiDotsVertical,
                document: mdiHelpCircleOutline,
                replace: mdiFindReplace,
                previous: mdiClose,
                next: mdiArrowRight,
                complete: mdiCheck,
                select: mdiSelect,
                selectAll: mdiSelectAll,
                selectNone: mdiSelect,
                selectInverse: mdiSelectInverse,
                edit: mdiPencilOutline,
                arrowRight: mdiArrowRight,
                alert: mdiAlertOutline,
                tag: mdiPound
            }
        }
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore, useAccountsStore, useTransactionCategoriesStore, useTransactionTagsStore, useTransactionsStore, useOverviewStore, useStatisticsStore, useExchangeRatesStore),
        defaultCurrency() {
            return this.userStore.currentUserDefaultCurrency;
        },
        allSteps() {
            return [
                {
                    name: 'uploadFile',
                    title: this.$t('Upload File'),
                    subTitle: this.$t('Upload Transaction Data File')
                },
                {
                    name: 'checkData',
                    title: this.$t('Check & Modify'),
                    subTitle: this.$t('Check and Modify Your Data')
                },
                {
                    name: 'finalResult',
                    title: this.$t('Complete'),
                    subTitle: this.$t('Data Import Completed')
                }
            ];
        },
        allSupportedImportFileTypes() {
            return this.$locale.getAllSupportedImportFileTypes();
        },
        allFileSubTypes() {
            return getNameByKeyValue(this.allSupportedImportFileTypes, this.fileType, 'type', 'subTypes');
        },
        allTransactionTypes() {
            return TransactionType;
        },
        allCategoryTypes() {
            return CategoryType;
        },
        allAccounts() {
            return this.accountsStore.allPlainAccounts;
        },
        allVisibleAccounts() {
            return this.accountsStore.allVisiblePlainAccounts;
        },
        allVisibleCategorizedAccounts() {
            return this.$locale.getCategorizedAccountsWithDisplayBalance(this.allVisibleAccounts, this.showAccountBalance, this.defaultCurrency, this.settingsStore, this.userStore, this.exchangeRatesStore);
        },
        allAccountsMap() {
            return this.accountsStore.allAccountsMap;
        },
        allCategories() {
            return this.transactionCategoriesStore.allTransactionCategories;
        },
        allCategoriesMap() {
            return this.transactionCategoriesStore.allTransactionCategoriesMap;
        },
        allTags() {
            return this.transactionTagsStore.allTransactionTags;
        },
        allTagsMap() {
            return this.transactionTagsStore.allTransactionTagsMap;
        },
        hasAvailableExpenseCategories() {
            if (!this.allCategories || !this.allCategories[this.allCategoryTypes.Expense] || !this.allCategories[this.allCategoryTypes.Expense].length) {
                return false;
            }

            const firstAvailableCategoryId = getFirstAvailableCategoryId(this.allCategories[this.allCategoryTypes.Expense]);
            return firstAvailableCategoryId !== '';
        },
        hasAvailableIncomeCategories() {
            if (!this.allCategories || !this.allCategories[this.allCategoryTypes.Income] || !this.allCategories[this.allCategoryTypes.Income].length) {
                return false;
            }

            const firstAvailableCategoryId = getFirstAvailableCategoryId(this.allCategories[this.allCategoryTypes.Income]);
            return firstAvailableCategoryId !== '';
        },
        hasAvailableTransferCategories() {
            if (!this.allCategories || !this.allCategories[this.allCategoryTypes.Transfer] || !this.allCategories[this.allCategoryTypes.Transfer].length) {
                return false;
            }

            const firstAvailableCategoryId = getFirstAvailableCategoryId(this.allCategories[this.allCategoryTypes.Transfer]);
            return firstAvailableCategoryId !== '';
        },
        showAccountBalance() {
            return this.settingsStore.appSettings.showAccountBalance;
        },
        currentTimezoneOffsetMinutes() {
            return getTimezoneOffsetMinutes(this.settingsStore.appSettings.timeZone);
        },
        supportedImportFileExtensions() {
            if (this.allFileSubTypes && this.allFileSubTypes.length) {
                const subTypeExtensions = getNameByKeyValue(this.allFileSubTypes, this.fileSubType, 'type', 'extensions');

                if (subTypeExtensions) {
                    return subTypeExtensions;
                }
            }

            return getNameByKeyValue(this.allSupportedImportFileTypes, this.fileType, 'type', 'extensions');
        },
        exportFileGuideDocumentUrl() {
            const document = getNameByKeyValue(this.allSupportedImportFileTypes, this.fileType, 'type', 'document');

            if (!document) {
                return null;
            }

            const language = document.language ? document.language + '/' : '';
            const anchor = document.anchor ? '#' + document.anchor : '';
            return `https://ezbookkeeping.mayswind.net/${language}export_and_import${anchor}`;
        },
        exportFileGuideDocumentLanguageName() {
            const document = getNameByKeyValue(this.allSupportedImportFileTypes, this.fileType, 'type', 'document');

            if (!document) {
                return null;
            }

            return document.displayLanguageName;
        },
        fileName: {
            get: function () {
                if (this.importFile == null) {
                    return '';
                } else {
                    return this.importFile.name;
                }
            }
        },
        importTransactionsTableHeight() {
            if (this.countPerPage <= 10 || !this.importTransactions || this.importTransactions.length <= 10) {
                return undefined;
            } else {
                return 400;
            }
        },
        importTransactionHeaders() {
            return [
                { value: 'valid', sortable: true, nowrap: true, width: 35 },
                { value: 'time', title: this.$t('Transaction Time'), sortable: true, nowrap: true, maxWidth: 280 },
                { value: 'type', title: this.$t('Type'), sortable: true, nowrap: true, maxWidth: 140 },
                { value: 'actualCategoryName', title: this.$t('Category'), sortable: true, nowrap: true },
                { value: 'sourceAmount', title: this.$t('Amount'), sortable: true, nowrap: true },
                { value: 'actualSourceAccountName', title: this.$t('Account'), sortable: true, nowrap: true },
                { value: 'geoLocation', title: this.$t('Geographic Location'), sortable: true, nowrap: true },
                { value: 'tagIds', title: this.$t('Tags'), sortable: true, nowrap: true },
                { value: 'comment', title: this.$t('Description'), sortable: true, nowrap: true },
            ];
        },
        importTransactionsTablePageOptions() {
            const pageOptions = [];

            if (!this.importTransactions || this.importTransactions.length < 1) {
                pageOptions.push({ value: -1, title: this.$t('All') });
                return pageOptions;
            }

            const availableCountPerPage = [ 5, 10, 15, 20, 25, 30, 50 ];

            for (let i = 0; i < availableCountPerPage.length; i++) {
                const count = availableCountPerPage[i];

                if (this.importTransactions.length < count) {
                    break;
                }

                pageOptions.push({ value: count, title: count.toString() });
            }

            pageOptions.push({ value: -1, title: this.$t('All') });

            return pageOptions;
        },
        allPages() {
            const pages = [];

            for (let i = 1; i < this.totalPageCount; i++) {
                pages.push({
                    page: i
                });
            }

            return pages;
        },
        inputCurrentPage: {
            get: function () {
                return this.currentPage;
            },
            set: function (value) {
                if (value && value >= 1 && value < this.totalPageCount) {
                    this.currentPage = value;
                }
            }
        },
        totalPageCount() {
            if (!this.importTransactions || this.importTransactions.length < 1) {
                return 1;
            }

            let count = 0;

            for (let i = 0; i < this.importTransactions.length; i++) {
                if (this.isTransactionDisplayed(this.importTransactions[i])) {
                    count++;
                }
            }

            return Math.ceil(count / this.countPerPage);
        },
        currentPageTransactions() {
            const ret = [];
            const previousCount = Math.max(0, (this.currentPage - 1) * this.countPerPage);
            let count = 0;

            for (let i = 0; i < this.importTransactions.length; i++) {
                if (ret.length >= this.countPerPage) {
                    break;
                }

                if (this.isTransactionDisplayed(this.importTransactions[i])) {
                    if (count >= previousCount) {
                        ret.push(this.importTransactions[i]);
                    }

                    count++;
                }
            }

            return ret;
        },
        selectedImportTransactionCount() {
            let count = 0;

            for (let i = 0; i < this.importTransactions.length; i++) {
                if (this.importTransactions[i].selected) {
                    count++;
                }
            }

            return count;
        },
        selectedExpenseTransactionCount() {
            let count = 0;

            for (let i = 0; i < this.importTransactions.length; i++) {
                if (this.importTransactions[i].selected && this.importTransactions[i].type === this.allTransactionTypes.Expense) {
                    count++;
                }
            }

            return count;
        },
        selectedIncomeTransactionCount() {
            let count = 0;

            for (let i = 0; i < this.importTransactions.length; i++) {
                if (this.importTransactions[i].selected && this.importTransactions[i].type === this.allTransactionTypes.Income) {
                    count++;
                }
            }

            return count;
        },
        selectedTransferTransactionCount() {
            let count = 0;

            for (let i = 0; i < this.importTransactions.length; i++) {
                if (this.importTransactions[i].selected && this.importTransactions[i].type === this.allTransactionTypes.Transfer) {
                    count++;
                }
            }

            return count;
        },
        selectedInvalidTransactionCount() {
            let count = 0;

            for (let i = 0; i < this.importTransactions.length; i++) {
                if (!this.importTransactions[i].valid && this.importTransactions[i].selected) {
                    count++;
                }
            }

            return count;
        },
        anyButNotAllTransactionSelected: {
            get: function () {
                return this.selectedImportTransactionCount > 0 && this.selectedImportTransactionCount !== this.importTransactions.length;
            }
        },
        allTransactionSelected: {
            get: function () {
                return this.selectedImportTransactionCount === this.importTransactions.length;
            }
        },
        allUsedCategoryNames() {
            return this.getAllUsedCategoryNames();
        },
        allUsedAccountNames() {
            return this.getAllUsedAccountNames();
        },
        allUsedTagNames() {
            return this.getAllUsedTagNames();
        },
        allInvalidExpenseCategoryNames() {
            return this.getCurrentInvalidCategoryNames(this.allTransactionTypes.Expense);
        },
        allInvalidIncomeCategoryNames() {
            return this.getCurrentInvalidCategoryNames(this.allTransactionTypes.Income);
        },
        allInvalidTransferCategoryNames() {
            return this.getCurrentInvalidCategoryNames(this.allTransactionTypes.Transfer);
        },
        allInvalidAccountNames() {
            return this.getCurrentInvalidAccountNames();
        },
        allInvalidTransactionTagNames() {
            return this.getCurrentInvalidTagNames();
        },
        displayFilterCustomDateRange() {
            if (this.filters.minDatetime === null || this.filters.maxDatetime === null) {
                return '';
            }

            const minDisplayTime = this.$locale.formatUnixTimeToLongDateTime(this.userStore, this.filters.minDatetime);
            const maxDisplayTime = this.$locale.formatUnixTimeToLongDateTime(this.userStore, this.filters.maxDatetime);

            return `${minDisplayTime} - ${maxDisplayTime}`
        }
    },
    watch: {
        fileType: function () {
            if (this.allFileSubTypes && this.allFileSubTypes.length) {
                this.fileSubType = this.allFileSubTypes[0].type;
            }

            this.importFile = null;
            this.importTransactions = null;
            this.editingTransaction = null;
            this.editingTags = [];
            this.currentPage = 1;
            this.countPerPage = 10;
        },
        fileSubType: function (newValue) {
            let supportedExtensions = getNameByKeyValue(this.allFileSubTypes, newValue, 'type', 'extensions');

            if (!supportedExtensions) {
                supportedExtensions = getNameByKeyValue(this.allSupportedImportFileTypes, this.fileType, 'type', 'extensions');
            }

            if (this.importFile && this.importFile.name && !isFileExtensionSupported(this.importFile.name, supportedExtensions)) {
                this.importFile = null;
            }
        }
    },
    methods: {
        open() {
            const self = this;
            self.fileType = 'ezbookkeeping';
            self.fileSubType = 'ezbookkeeping_csv';
            self.currentStep = 'uploadFile';
            self.importFile = null;
            self.importTransactions = null;
            self.editingTransaction = null;
            self.editingTags = [];
            self.filters.minDatetime = null;
            self.filters.maxDatetime = null;
            self.filters.transactionType = null;
            self.filters.category = null;
            self.filters.account = null;
            self.filters.tag = null;
            self.filters.description = null;
            self.currentPage = 1;
            self.countPerPage = 10;
            self.showState = true;
            self.clientSessionId = generateRandomUUID();

            const promises = [
                self.accountsStore.loadAllAccounts({ force: false }),
                self.transactionCategoriesStore.loadAllCategories({ force: false }),
                self.transactionTagsStore.loadAllTags({ force: false })
            ];

            Promise.all(promises).then(function () {
                self.loading = false;
            }).catch(error => {
                logger.error('failed to load essential data for importing transaction', error);

                self.loading = false;
                self.showState = false;

                if (!error.processed) {
                    if (self.reject) {
                        self.reject(error);
                    }
                }
            });

            return new Promise((resolve, reject) => {
                self.resolve = resolve;
                self.reject = reject;
            });
        },
        showOpenFileDialog() {
            if (this.submitting) {
                return;
            }

            this.$refs.fileInput.click();
        },
        setImportFile(event) {
            if (!event || !event.target || !event.target.files || !event.target.files.length) {
                return;
            }

            this.importFile = event.target.files[0];
            event.target.value = null;
        },
        parseData() {
            const self = this;
            self.submitting = true;

            let fileType = self.fileType;

            if (self.allFileSubTypes) {
                fileType = self.fileSubType;
            }

            self.transactionsStore.parseImportTransaction({
                fileType: fileType,
                importFile: self.importFile
            }).then(response => {
                const parsedTransactions = response.items;

                if (parsedTransactions) {
                    for (let i = 0; i < parsedTransactions.length; i++) {
                        const transaction = parsedTransactions[i];
                        transaction.index = i;
                        transaction.selected = false;
                        transaction.valid = self.isTransactionValid(transaction);
                        transaction.actualCategoryName = transaction.originalCategoryName;
                        transaction.actualSourceAccountName = transaction.originalSourceAccountName;
                        transaction.actualDestinationAccountName = transaction.originalDestinationAccountName;
                    }
                }

                self.importTransactions = parsedTransactions;
                self.editingTransaction = null;
                self.editingTags = [];
                self.currentPage = 1;

                if (self.importTransactions && self.importTransactions.length >= 0 && self.importTransactions.length < 10) {
                    self.countPerPage = -1;
                } else {
                    self.countPerPage = 10;
                }

                self.currentStep = 'checkData';
                self.submitting = false;
            }).catch(error => {
                self.submitting = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        submit() {
            const self = this;

            if (self.editingTransaction) {
                return;
            }

            const transactions = [];

            for (let i = 0; i < self.importTransactions.length; i++) {
                const transaction = self.importTransactions[i];

                if (transaction.valid && transaction.selected) {
                    transactions.push(transaction);
                } else if (!transaction.valid && transaction.selected) {
                    self.$refs.snackbar.showError('Cannot import invalid transactions');
                    return;
                }
            }

            if (transactions.length < 1) {
                self.$refs.snackbar.showError('No data to import');
                return;
            }

            self.$refs.confirmDialog.open('format.misc.confirmImportTransactions', {
                count: transactions.length
            }).then(() => {
                self.editingTransaction = null;
                self.editingTags = [];
                self.submitting = true;

                self.transactionsStore.importTransactions({
                    transactions: transactions,
                    clientSessionId: self.clientSessionId
                }).then(response => {
                    self.importedCount = response;
                    self.currentStep = 'finalResult';

                    self.accountsStore.updateAccountListInvalidState(true);
                    self.transactionsStore.updateTransactionListInvalidState(true);
                    self.overviewStore.updateTransactionOverviewInvalidState(true);
                    self.statisticsStore.updateTransactionStatisticsInvalidState(true);

                    self.submitting = false;
                }).catch(error => {
                    self.submitting = false;

                    if (!error.processed) {
                        self.$refs.snackbar.showError(error);
                    }
                });
            });
        },
        cancel(completed) {
            if (completed) {
                if (this.resolve) {
                    this.resolve();
                }
            } else {
                if (this.reject) {
                    this.reject();
                }
            }

            this.showState = false;
        },
        changeCustomDateFilter(minTime, maxTime) {
            this.filters.minDatetime = minTime;
            this.filters.maxDatetime = maxTime;
            this.showCustomDateRangeDialog = false;
        },
        importTransactionsFilter(value, query, item) {
            if (!item || !item.raw) {
                return false;
            }

            return this.isTransactionDisplayed(item.raw);
        },
        isTransactionDisplayed(transaction) {
            if (isNumber(this.filters.minDatetime) && isNumber(this.filters.maxDatetime) && (transaction.time < this.filters.minDatetime || transaction.time > this.filters.maxDatetime)) {
                return false;
            }

            if (isNumber(this.filters.transactionType) && transaction.type !== this.filters.transactionType) {
                return false;
            }

            if (isString(this.filters.category)) {
                if (this.filters.category === '' && transaction.actualCategoryName !== '') {
                    return false;
                } else if (this.filters.category !== '' && transaction.actualCategoryName !== this.filters.category) {
                    return false;
                }
            } else if (this.filters.category === undefined) {
                if (transaction.type !== this.allTransactionTypes.ModifyBalance && transaction.categoryId && transaction.categoryId !== '0') {
                    return false;
                }
            }

            if (isString(this.filters.account)) {
                if (this.filters.account === '' && (transaction.actualSourceAccountName !== '' || transaction.actualDestinationAccountName !== '')) {
                    return false;
                } else if (this.filters.account !== '' && transaction.actualSourceAccountName !== this.filters.account && transaction.actualDestinationAccountName !== this.filters.account) {
                    return false;
                }
            } else if (this.filters.account === undefined) {
                if (transaction.type !== this.allTransactionTypes.Transfer && transaction.sourceAccountId && transaction.sourceAccountId !== '0') {
                    return false;
                } else if (transaction.type === this.allTransactionTypes.Transfer && transaction.sourceAccountId && transaction.sourceAccountId !== '0' && transaction.destinationAccountId && transaction.destinationAccountId !== '0') {
                    return false;
                }
            }

            if (isString(this.filters.tag)) {
                if (this.filters.tag === '' && transaction.tagIds && transaction.tagIds.length) {
                    return false;
                } else if (this.filters.tag !== '') {
                    let hasTagName = false;

                    if (transaction.tagIds && transaction.tagIds.length) {
                        for (let i = 0; i < transaction.tagIds.length; i++) {
                            const tagId = transaction.tagIds[i];
                            let tagName = transaction.originalTagNames ? transaction.originalTagNames[i] : "";

                            if (tagId && tagId !== '0' && this.allTagsMap[tagId] && this.allTagsMap[tagId].name) {
                                tagName = this.allTagsMap[tagId].name;
                            }

                            if (tagName === this.filters.tag) {
                                hasTagName = true;
                                break;
                            }
                        }
                    }

                    if (!hasTagName) {
                        return false;
                    }
                }
            } else if (this.filters.tag === undefined) {
                if (transaction.tagIds && transaction.tagIds.length) {
                    let hasInvalidTag = false;

                    for (let i = 0; i < transaction.tagIds.length; i++) {
                        if (!transaction.tagIds[i] || transaction.tagIds[i] === '0') {
                            hasInvalidTag = true;
                            break;
                        }
                    }

                    if (!hasInvalidTag) {
                        return false;
                    }
                } else {
                    return false;
                }
            }

            if (isString(this.filters.description)) {
                if (this.filters.description === '' && transaction.comment !== '') {
                    return false;
                } else if (this.filters.description !== '' && transaction.comment.indexOf(this.filters.description) < 0) {
                    return false;
                }
            }

            return true;
        },
        selectAllValid() {
            for (let i = 0; i < this.importTransactions.length; i++) {
                if (this.importTransactions[i].valid && this.isTransactionDisplayed(this.importTransactions[i])) {
                    this.importTransactions[i].selected = true;
                }
            }
        },
        selectAllInvalid() {
            for (let i = 0; i < this.importTransactions.length; i++) {
                if (!this.importTransactions[i].valid && this.isTransactionDisplayed(this.importTransactions[i])) {
                    this.importTransactions[i].selected = true;
                }
            }
        },
        selectAll() {
            for (let i = 0; i < this.importTransactions.length; i++) {
                if (this.isTransactionDisplayed(this.importTransactions[i])) {
                    this.importTransactions[i].selected = true;
                }
            }
        },
        selectNone() {
            for (let i = 0; i < this.importTransactions.length; i++) {
                if (this.isTransactionDisplayed(this.importTransactions[i])) {
                    this.importTransactions[i].selected = false;
                }
            }
        },
        selectInvert() {
            for (let i = 0; i < this.importTransactions.length; i++) {
                if (this.isTransactionDisplayed(this.importTransactions[i])) {
                    this.importTransactions[i].selected = !this.importTransactions[i].selected;
                }
            }
        },
        selectAllInThisPage() {
            for (let i = 0; i < this.currentPageTransactions.length; i++) {
                this.currentPageTransactions[i].selected = true;
            }
        },
        selectNoneInThisPage() {
            for (let i = 0; i < this.currentPageTransactions.length; i++) {
                this.currentPageTransactions[i].selected = false;
            }
        },
        selectInvertInThisPage() {
            for (let i = 0; i < this.currentPageTransactions.length; i++) {
                this.currentPageTransactions[i].selected = !this.currentPageTransactions[i].selected;
            }
        },
        editTransaction(transaction) {
            if (this.editingTransaction) {
                this.editingTransaction.tagIds = this.editingTags;
                this.updateTransactionData(this.editingTransaction);
            }

            if (this.editingTransaction === transaction) {
                this.editingTags = [];
                this.editingTransaction = null;
            } else {
                this.editingTransaction = transaction;
                this.editingTags = this.editingTransaction.tagIds;
            }
        },
        updateTransactionData(transaction) {
            transaction.valid = this.isTransactionValid(transaction);

            if (transaction.categoryId && this.allCategoriesMap[transaction.categoryId]) {
                transaction.actualCategoryName = this.allCategoriesMap[transaction.categoryId].name;
            }

            if (transaction.sourceAccountId && this.allAccountsMap[transaction.sourceAccountId]) {
                transaction.actualSourceAccountName = this.allAccountsMap[transaction.sourceAccountId].name;
            }

            if (transaction.destinationAccountId && this.allAccountsMap[transaction.destinationAccountId]) {
                transaction.actualDestinationAccountName = this.allAccountsMap[transaction.destinationAccountId].name;
            }
        },
        showBatchReplaceDialog(type) {
            const self = this;

            if (self.editingTransaction) {
                return;
            }

            self.$refs.batchReplaceDialog.open({
                mode: 'batchReplace',
                type: type
            }).then(result => {
                if (!result || !result.targetItem) {
                    return;
                }

                let updatedCount = 0;

                for (let i = 0; i < self.importTransactions.length; i++) {
                    const transaction = self.importTransactions[i];

                    if (!transaction.selected) {
                        continue;
                    }

                    let updated = false;

                    if (type === 'expenseCategory') {
                        if (transaction.type === self.allTransactionTypes.Expense) {
                            transaction.categoryId = result.targetItem;
                            updated = true;
                        }
                    } else if (type === 'incomeCategory') {
                        if (transaction.type === self.allTransactionTypes.Income) {
                            transaction.categoryId = result.targetItem;
                            updated = true;
                        }
                    } else if (type === 'transferCategory') {
                        if (transaction.type === self.allTransactionTypes.Transfer) {
                            transaction.categoryId = result.targetItem;
                            updated = true;
                        }
                    } else if (type === 'account') {
                        transaction.sourceAccountId = result.targetItem;
                        updated = true;
                    } else if (type === 'destinationAccount') {
                        if (transaction.type === self.allTransactionTypes.Transfer) {
                            transaction.destinationAccountId = result.targetItem;
                            updated = true;
                        }
                    }

                    if (updated) {
                        updatedCount++;
                        self.updateTransactionData(transaction);
                    }
                }

                if (updatedCount > 0) {
                    self.$refs.snackbar.showMessage('format.misc.youHaveUpdatedTransactions', {
                        count: updatedCount
                    });
                }
            });
        },
        showReplaceInvalidItemDialog(type, invalidItems) {
            const self = this;

            if (self.editingTransaction) {
                return;
            }

            self.$refs.batchReplaceDialog.open({
                mode: 'replaceInvalidItems',
                type: type,
                invalidItems: invalidItems
            }).then(result => {
                if (!result || (!result.sourceItem && result.sourceItem !== '') || !result.targetItem) {
                    return;
                }

                let updatedCount = 0;

                for (let i = 0; i < self.importTransactions.length; i++) {
                    const transaction = self.importTransactions[i];

                    if (transaction.valid) {
                        continue;
                    }

                    let updated = false;

                    if (type === 'expenseCategory' || type === 'incomeCategory' || type === 'transferCategory') {
                        const categoryId = transaction.categoryId;
                        const originalCategoryName = transaction.originalCategoryName;

                        if (transaction.type !== self.allTransactionTypes.ModifyBalance && originalCategoryName === result.sourceItem && (!categoryId || categoryId === '0' || !self.allCategoriesMap[categoryId])) {
                            if (type === 'expenseCategory' && transaction.type === self.allTransactionTypes.Expense) {
                                transaction.categoryId = result.targetItem;
                                updated = true;
                            } else if (type === 'incomeCategory' && transaction.type === self.allTransactionTypes.Income) {
                                transaction.categoryId = result.targetItem;
                                updated = true;
                            } else if (type === 'transferCategory' && transaction.type === self.allTransactionTypes.Transfer) {
                                transaction.categoryId = result.targetItem;
                                updated = true;
                            }
                        }
                    } else if (type === 'account') {
                        const sourceAccountId = transaction.sourceAccountId;
                        const originalSourceAccountName = transaction.originalSourceAccountName;
                        const destinationAccountId = transaction.destinationAccountId;
                        const originalDestinationAccountName = transaction.originalDestinationAccountName;

                        if (originalSourceAccountName === result.sourceItem && (!sourceAccountId || sourceAccountId === '0' || !self.allAccountsMap[sourceAccountId])) {
                            transaction.sourceAccountId = result.targetItem;
                            updated = true;
                        }

                        if (transaction.type === self.allTransactionTypes.Transfer && originalDestinationAccountName === result.sourceItem && (!destinationAccountId || destinationAccountId === '0' || !self.allAccountsMap[destinationAccountId])) {
                            transaction.destinationAccountId = result.targetItem;
                            updated = true;
                        }
                    } else if (type === 'tag' && transaction.tagIds) {
                        for (let j = 0; j < transaction.tagIds.length; j++) {
                            const tagId = transaction.tagIds[j];
                            const originalTagName = transaction.originalTagNames ? transaction.originalTagNames[j] : "";

                            if (originalTagName === result.sourceItem && (!tagId || tagId === '0' || !self.allTagsMap[tagId])) {
                                transaction.tagIds[j] = result.targetItem;
                                updated = true;
                            }
                        }
                    }

                    if (updated) {
                        updatedCount++;
                        self.updateTransactionData(transaction);
                    }
                }

                if (updatedCount > 0) {
                    self.$refs.snackbar.showMessage('format.misc.youHaveUpdatedTransactions', {
                        count: updatedCount
                    });
                }
            });
        },
        showError(message) {
            this.$refs.snackbar.showError(message);
        },
        getAllUsedCategoryNames() {
            const categoryNames = {};

            for (let i = 0; i < this.importTransactions.length; i++) {
                const transaction = this.importTransactions[i];

                if (transaction.actualCategoryName && transaction.actualCategoryName !== '') {
                    categoryNames[transaction.actualCategoryName] = true;
                }
            }

            return objectFieldToArrayItem(categoryNames);
        },
        getAllUsedAccountNames() {
            const accountNames = {};

            for (let i = 0; i < this.importTransactions.length; i++) {
                const transaction = this.importTransactions[i];

                if (transaction.actualSourceAccountName && transaction.actualSourceAccountName !== '') {
                    accountNames[transaction.actualSourceAccountName] = true;
                }

                if (transaction.actualDestinationAccountName && transaction.actualDestinationAccountName !== '') {
                    accountNames[transaction.actualDestinationAccountName] = true;
                }
            }

            return objectFieldToArrayItem(accountNames);
        },
        getAllUsedTagNames(){
            const tagNames = {};

            for (let i = 0; i < this.importTransactions.length; i++) {
                const transaction = this.importTransactions[i];

                if (!transaction.tagIds || !transaction.originalTagNames) {
                    continue;
                }

                for (let j = 0; j < transaction.tagIds.length; j++) {
                    const tagId = transaction.tagIds[j];
                    const originalTagName = transaction.originalTagNames[j];

                    if (tagId && tagId !== '0' && this.allTagsMap[tagId] && this.allTagsMap[tagId].name) {
                        tagNames[this.allTagsMap[tagId].name] = true;
                    } else if (originalTagName) {
                        tagNames[originalTagName] = true;
                    }
                }
            }

            return objectFieldToArrayItem(tagNames);
        },
        getCurrentInvalidCategoryNames(transactionType) {
            const invalidCategoryNames = {};
            const invalidCategories = [];

            for (let i = 0; i < this.importTransactions.length; i++) {
                const transaction = this.importTransactions[i];
                const categoryId = transaction.categoryId;

                if (transaction.type === transactionType && (!categoryId || categoryId === '0' || !this.allCategoriesMap[categoryId])) {
                    invalidCategoryNames[transaction.originalCategoryName] = true;
                }
            }

            for (const name in invalidCategoryNames) {
                if (!Object.prototype.hasOwnProperty.call(invalidCategoryNames, name)) {
                    continue;
                }

                invalidCategories.push({
                    name: name || this.$t('(Empty)'),
                    value: name
                });
            }

            return invalidCategories;
        },
        getCurrentInvalidAccountNames() {
            const invalidAccountNames = {};
            const invalidAccounts = [];

            for (let i = 0; i < this.importTransactions.length; i++) {
                const transaction = this.importTransactions[i];
                const sourceAccountId = transaction.sourceAccountId;
                const destinationAccountId = transaction.destinationAccountId;

                if (!sourceAccountId || sourceAccountId === '0' || !this.allAccountsMap[sourceAccountId]) {
                    invalidAccountNames[transaction.originalSourceAccountName] = true;
                }

                if (transaction.type === this.allTransactionTypes.Transfer && (!destinationAccountId || destinationAccountId === '0' || !this.allAccountsMap[destinationAccountId])) {
                    invalidAccountNames[transaction.originalDestinationAccountName] = true;
                }
            }

            for (const name in invalidAccountNames) {
                if (!Object.prototype.hasOwnProperty.call(invalidAccountNames, name)) {
                    continue;
                }

                invalidAccounts.push({
                    name: name || this.$t('(Empty)'),
                    value: name
                });
            }

            return invalidAccounts;
        },
        getCurrentInvalidTagNames() {
            const invalidTagNames = {};
            const invalidTags = [];

            for (let i = 0; i < this.importTransactions.length; i++) {
                const transaction = this.importTransactions[i];

                if (!transaction.tagIds || !transaction.originalTagNames) {
                    continue;
                }

                for (let j = 0; j < transaction.tagIds.length; j++) {
                    const tagId = transaction.tagIds[j];
                    const originalTagName = transaction.originalTagNames[j];

                    if (!tagId || tagId === '0' || !this.allTagsMap[tagId]) {
                        invalidTagNames[originalTagName] = true;
                    }
                }
            }

            for (const name in invalidTagNames) {
                if (!Object.prototype.hasOwnProperty.call(invalidTagNames, name)) {
                    continue;
                }

                invalidTags.push({
                    name: name || this.$t('(Empty)'),
                    value: name
                });
            }

            return invalidTags;
        },
        isTransactionValid(transaction) {
            if (!transaction) {
                return false;
            }

            if (transaction.type !== this.allTransactionTypes.ModifyBalance && (!transaction.categoryId || transaction.categoryId === '0')) {
                return false;
            }

            if (!transaction.sourceAccountId || transaction.sourceAccountId === '0') {
                return false;
            }

            if (transaction.type === this.allTransactionTypes.Transfer && (!transaction.destinationAccountId || transaction.destinationAccountId === '0')) {
                return false;
            }

            if (transaction.tagIds && transaction.tagIds.length) {
                for (let j = 0; j < transaction.tagIds.length; j++) {
                    if (!transaction.tagIds[j] || transaction.tagIds[j] === '0') {
                        return false;
                    }
                }
            }

            return true;
        },
        isTagValid(tagIds, tagIndex) {
            if (!tagIds || !tagIds[tagIndex]) {
                return false;
            }

            if (tagIds[tagIndex] === '0') {
                return false;
            }

            const tagId = tagIds[tagIndex];
            return !!this.allTagsMap[tagId];
        },
        getDisplayDateTime(transaction) {
            const transactionTime = getUnixTime(parseDateFromUnixTime(transaction.time, transaction.utcOffset, this.currentTimezoneOffsetMinutes));
            return this.$locale.formatUnixTimeToLongDateTime(this.userStore, transactionTime);
        },
        getDisplayTimezone(transaction) {
            return `UTC${getUtcOffsetByUtcOffsetMinutes(transaction.utcOffset)}`;
        },
        getDisplayCurrency(value, currencyCode) {
            return this.$locale.formatAmountWithCurrency(this.settingsStore, this.userStore, value, currencyCode);
        },
        getTransactionDisplayAmount(transaction) {
            let currency = transaction.originalSourceAccountCurrency || this.userStore.currentUserDefaultCurrency;

            if (transaction.sourceAccountId && transaction.sourceAccountId !== '0' && this.allAccountsMap[transaction.sourceAccountId]) {
                currency = this.allAccountsMap[transaction.sourceAccountId].currency;
            }

            return this.getDisplayCurrency(transaction.sourceAmount, currency);
        },
        getTransactionDisplayDestinationAmount(transaction) {
            if (transaction.type !== this.allTransactionTypes.Transfer) {
                return '-';
            }

            let currency = transaction.originalDestinationAccountCurrency || this.userStore.currentUserDefaultCurrency;

            if (transaction.destinationAccountId && transaction.destinationAccountId !== '0' && this.allAccountsMap[transaction.destinationAccountId]) {
                currency = this.allAccountsMap[transaction.destinationAccountId].currency;
            }

            return this.getDisplayCurrency(transaction.destinationAmount, currency);
        },
        getPrimaryCategoryName(categoryId, allCategories) {
            return getTransactionPrimaryCategoryName(categoryId, allCategories);
        },
        getSecondaryCategoryName(categoryId, allCategories) {
            return getTransactionSecondaryCategoryName(categoryId, allCategories);
        },
        getSourceAccountTitle(transaction) {
            if (transaction.type === this.allTransactionTypes.Expense || transaction.type === this.allTransactionTypes.Income) {
                return this.$t('Account');
            } else if (transaction.type === this.allTransactionTypes.Transfer) {
                return this.$t('Source Account');
            } else {
                return this.$t('Account');
            }
        },
        getSourceAccountDisplayName(transaction) {
            if (transaction.sourceAccountId) {
                return getNameByKeyValue(this.allAccounts, transaction.sourceAccountId, 'id', 'name');
            } else {
                return this.$t('None');
            }
        },
        getDestinationAccountDisplayName(transaction) {
            if (transaction.destinationAccountId) {
                return getNameByKeyValue(this.allAccounts, transaction.destinationAccountId, 'id', 'name');
            } else {
                return this.$t('None');
            }
        }
    }
}
</script>

<style>
.import-transaction-table .v-autocomplete.v-input.v-input--density-compact:not(.v-textarea) .v-field__input,
.import-transaction-table .v-select.v-input.v-input--density-compact:not(.v-textarea) .v-field__input {
    min-height: inherit;
    padding-top: 4px;
}

.import-transaction-table .v-chip.transaction-tag {
    margin-right: 4px;
    margin-top: 2px;
    margin-bottom: 2px;
}

.import-transaction-table .v-chip.transaction-tag > .v-chip__content {
    display: block;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>
