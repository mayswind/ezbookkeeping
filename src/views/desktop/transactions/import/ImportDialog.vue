<template>
    <v-dialog :persistent="!!persistent" v-model="showState">
        <v-card class="pa-6 pa-sm-10 pa-md-12">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <div class="d-flex w-100 align-center justify-center">
                        <h4 class="text-h4">{{ tt('Import Transactions') }}</h4>
                        <v-progress-circular indeterminate size="22" class="ml-2" v-if="loading"></v-progress-circular>
                    </div>
                    <v-btn density="comfortable" color="default" variant="text" class="ml-2"
                           :icon="true" :disabled="loading || submitting"
                           v-if="currentStep === 'checkData'">
                        <v-icon :icon="mdiFilterOutline" />
                        <v-menu activator="parent" max-height="500">
                            <v-list>
                                <v-list-subheader :title="tt('Date Range')"/>
                                <v-list-item :title="tt('All')"
                                             :append-icon="filters.minDatetime === null || filters.maxDatetime === null ? mdiCheck : undefined"
                                             @click="filters.minDatetime = filters.maxDatetime = null"></v-list-item>
                                <v-list-item :title="tt('Custom')"
                                             :subtitle="displayFilterCustomDateRange"
                                             :append-icon="filters.minDatetime !== null && filters.maxDatetime !== null ? mdiCheck : undefined"
                                             @click="showCustomDateRangeDialog = true"></v-list-item>
                                <v-divider class="my-2"/>
                                <v-list-subheader :title="tt('Type')"/>
                                <v-list-item :title="tt('All')"
                                             :append-icon="filters.transactionType === null ? mdiCheck : undefined"
                                             @click="filters.transactionType = null"></v-list-item>
                                <v-list-item :title="tt('Income')"
                                             :append-icon="filters.transactionType === TransactionType.Income ? mdiCheck : undefined"
                                             @click="filters.transactionType = TransactionType.Income"></v-list-item>
                                <v-list-item :title="tt('Expense')"
                                             :append-icon="filters.transactionType === TransactionType.Expense ? mdiCheck : undefined"
                                             @click="filters.transactionType = TransactionType.Expense"></v-list-item>
                                <v-list-item :title="tt('Transfer')"
                                             :append-icon="filters.transactionType === TransactionType.Transfer ? mdiCheck : undefined"
                                             @click="filters.transactionType = TransactionType.Transfer"></v-list-item>
                                <v-divider class="my-2"/>
                                <v-list-subheader :title="tt('Category')"/>
                                <v-list-item :title="tt('All')"
                                             :append-icon="filters.category === null ? mdiCheck : undefined"
                                             @click="filters.category = null"></v-list-item>
                                <v-list-item :title="tt('Invalid Category')"
                                             :append-icon="filters.category === undefined ? mdiCheck : undefined"
                                             @click="filters.category = undefined"></v-list-item>
                                <v-list-item :title="tt('None')"
                                             :append-icon="filters.category === '' ? mdiCheck : undefined"
                                             @click="filters.category = ''"></v-list-item>
                                <v-list-item :title="name" :key="name"
                                             :append-icon="filters.category === name ? mdiCheck : undefined"
                                             v-for="name in allUsedCategoryNames"
                                             @click="filters.category = name"></v-list-item>
                                <v-divider class="my-2"/>
                                <v-list-subheader :title="tt('Account')"/>
                                <v-list-item :title="tt('All')"
                                             :append-icon="filters.account === null ? mdiCheck : undefined"
                                             @click="filters.account = null"></v-list-item>
                                <v-list-item :title="tt('Invalid Account')"
                                             :append-icon="filters.account === undefined ? mdiCheck : undefined"
                                             @click="filters.account = undefined"></v-list-item>
                                <v-list-item :title="tt('None')"
                                             :append-icon="filters.account === '' ? mdiCheck : undefined"
                                             @click="filters.account = ''"></v-list-item>
                                <v-list-item :title="name" :key="name"
                                             :append-icon="filters.account === name ? mdiCheck : undefined"
                                             v-for="name in allUsedAccountNames"
                                             @click="filters.account = name"></v-list-item>
                                <v-divider class="my-2"/>
                                <v-list-subheader :title="tt('Tags')"/>
                                <v-list-item :title="tt('All')"
                                             :append-icon="filters.tag === null ? mdiCheck : undefined"
                                             @click="filters.tag = null"></v-list-item>
                                <v-list-item :title="tt('Invalid Tag')"
                                             :append-icon="filters.tag === undefined ? mdiCheck : undefined"
                                             @click="filters.tag = undefined"></v-list-item>
                                <v-list-item :title="tt('None')"
                                             :append-icon="filters.tag === '' ? mdiCheck : undefined"
                                             @click="filters.tag = ''"></v-list-item>
                                <v-list-item :title="name" :key="name"
                                             :append-icon="filters.tag === name ? mdiCheck : undefined"
                                             v-for="name in allUsedTagNames"
                                             @click="filters.tag = name"></v-list-item>
                                <v-divider class="my-2"/>
                                <v-list-subheader :title="tt('Description')"/>
                                <v-list-item :title="tt('All')"
                                             :append-icon="filters.description === null ? mdiCheck : undefined"
                                             @click="filters.description = null"></v-list-item>
                                <v-list-item :title="tt('None')"
                                             :append-icon="filters.description === '' ? mdiCheck : undefined"
                                             @click="filters.description = ''"></v-list-item>
                                <v-list-item :title="tt('Custom')"
                                             :subtitle="filters.description !== null ? filters.description : undefined"
                                             :append-icon="filters.description !== null && filters.description !== '' ? mdiCheck : undefined"
                                             @click="currentDescriptionFilterValue = filters.description || ''; showCustomDescriptionDialog = true"></v-list-item>
                            </v-list>
                        </v-menu>
                    </v-btn>
                    <v-btn density="comfortable" color="default" variant="text" class="ml-2"
                           :icon="true" :disabled="loading || submitting"
                           v-if="currentStep === 'checkData'">
                        <v-icon :icon="mdiDotsVertical" />
                        <v-menu activator="parent" max-height="500">
                            <v-list>
                                <v-list-item :prepend-icon="mdiFindReplace"
                                             :disabled="!!editingTransaction || selectedExpenseTransactionCount < 1"
                                             :title="tt('Batch Replace Selected Expense Categories')"
                                             @click="showBatchReplaceDialog('expenseCategory')"></v-list-item>
                                <v-list-item :prepend-icon="mdiFindReplace"
                                             :disabled="!!editingTransaction || selectedIncomeTransactionCount < 1"
                                             :title="tt('Batch Replace Selected Income Categories')"
                                             @click="showBatchReplaceDialog('incomeCategory')"></v-list-item>
                                <v-list-item :prepend-icon="mdiFindReplace"
                                             :disabled="!!editingTransaction || selectedTransferTransactionCount < 1"
                                             :title="tt('Batch Replace Selected Transfer Categories')"
                                             @click="showBatchReplaceDialog('transferCategory')"></v-list-item>
                                <v-list-item :prepend-icon="mdiFindReplace"
                                             :disabled="!!editingTransaction || selectedImportTransactionCount < 1"
                                             :title="tt('Batch Replace Selected Accounts')"
                                             @click="showBatchReplaceDialog('account')"></v-list-item>
                                <v-list-item :prepend-icon="mdiFindReplace"
                                             :disabled="!!editingTransaction || selectedTransferTransactionCount < 1"
                                             :title="tt('Batch Replace Selected Destination Accounts')"
                                             @click="showBatchReplaceDialog('destinationAccount')"></v-list-item>
                                <v-divider class="my-2"/>
                                <v-list-item :prepend-icon="mdiFindReplace"
                                             :disabled="!!editingTransaction || !allInvalidExpenseCategoryNames || allInvalidExpenseCategoryNames.length < 1"
                                             :title="tt('Replace Invalid Expense Categories')"
                                             @click="showReplaceInvalidItemDialog('expenseCategory', allInvalidExpenseCategoryNames)"></v-list-item>
                                <v-list-item :prepend-icon="mdiFindReplace"
                                             :disabled="!!editingTransaction || !allInvalidIncomeCategoryNames || allInvalidIncomeCategoryNames.length < 1"
                                             :title="tt('Replace Invalid Income Categories')"
                                             @click="showReplaceInvalidItemDialog('incomeCategory', allInvalidIncomeCategoryNames)"></v-list-item>
                                <v-list-item :prepend-icon="mdiFindReplace"
                                             :disabled="!!editingTransaction || !allInvalidTransferCategoryNames || allInvalidTransferCategoryNames.length < 1"
                                             :title="tt('Replace Invalid Transfer Categories')"
                                             @click="showReplaceInvalidItemDialog('transferCategory', allInvalidTransferCategoryNames)"></v-list-item>
                                <v-list-item :prepend-icon="mdiFindReplace"
                                             :disabled="!!editingTransaction || !allInvalidAccountNames || allInvalidAccountNames.length < 1"
                                             :title="tt('Replace Invalid Accounts')"
                                             @click="showReplaceInvalidItemDialog('account', allInvalidAccountNames)"></v-list-item>
                                <v-list-item :prepend-icon="mdiFindReplace"
                                             :disabled="!!editingTransaction || !allInvalidTransactionTagNames || allInvalidTransactionTagNames.length < 1"
                                             :title="tt('Replace Invalid Transaction Tags')"
                                             @click="showReplaceInvalidItemDialog('tag', allInvalidTransactionTagNames)"></v-list-item>
                                <v-divider class="my-2"/>
                                <v-list-item :prepend-icon="mdiShapePlusOutline"
                                             :disabled="!!editingTransaction || !allInvalidExpenseCategoryNames || allInvalidExpenseCategoryNames.length < 1"
                                             :title="tt('Create Nonexistent Expense Categories')"
                                             @click="showBatchCreateInvalidItemDialog('expenseCategory', allInvalidExpenseCategoryNames)"></v-list-item>
                                <v-list-item :prepend-icon="mdiShapePlusOutline"
                                             :disabled="!!editingTransaction || !allInvalidIncomeCategoryNames || allInvalidIncomeCategoryNames.length < 1"
                                             :title="tt('Create Nonexistent Income Categories')"
                                             @click="showBatchCreateInvalidItemDialog('incomeCategory', allInvalidIncomeCategoryNames)"></v-list-item>
                                <v-list-item :prepend-icon="mdiShapePlusOutline"
                                             :disabled="!!editingTransaction || !allInvalidTransferCategoryNames || allInvalidTransferCategoryNames.length < 1"
                                             :title="tt('Create Nonexistent Transfer Categories')"
                                             @click="showBatchCreateInvalidItemDialog('transferCategory', allInvalidTransferCategoryNames)"></v-list-item>
                                <v-list-item :prepend-icon="mdiShapePlusOutline"
                                             :disabled="!!editingTransaction || !allInvalidTransactionTagNames || allInvalidTransactionTagNames.length < 1"
                                             :title="tt('Create Nonexistent Transaction Tags')"
                                             @click="showBatchCreateInvalidItemDialog('tag', allInvalidTransactionTagNames)"></v-list-item>
                                <v-divider class="my-2"/>
                                <v-list-item :prepend-icon="mdiTransfer"
                                             :disabled="!!editingTransaction || selectedExpenseTransactionCount < 1"
                                             :title="tt('Batch Convert Expense Transaction to Income Transaction')"
                                             @click="convertTransactionType(TransactionType.Expense, TransactionType.Income)"></v-list-item>
                                <v-list-item :prepend-icon="mdiTransfer"
                                             :disabled="!!editingTransaction || selectedExpenseTransactionCount < 1"
                                             :title="tt('Batch Convert Expense Transaction to Transfer Transaction')"
                                             @click="convertTransactionType(TransactionType.Expense, TransactionType.Transfer)"></v-list-item>
                                <v-list-item :prepend-icon="mdiTransfer"
                                             :disabled="!!editingTransaction || selectedIncomeTransactionCount < 1"
                                             :title="tt('Batch Convert Income Transaction to Expense Transaction')"
                                             @click="convertTransactionType(TransactionType.Income, TransactionType.Expense)"></v-list-item>
                                <v-list-item :prepend-icon="mdiTransfer"
                                             :disabled="!!editingTransaction || selectedIncomeTransactionCount < 1"
                                             :title="tt('Batch Convert Income Transaction to Transfer Transaction')"
                                             @click="convertTransactionType(TransactionType.Income, TransactionType.Transfer)"></v-list-item>
                                <v-list-item :prepend-icon="mdiTransfer"
                                             :disabled="!!editingTransaction || selectedTransferTransactionCount < 1"
                                             :title="tt('Batch Convert Transfer Transaction to Expense Transaction')"
                                             @click="convertTransactionType(TransactionType.Transfer, TransactionType.Expense)"></v-list-item>
                                <v-list-item :prepend-icon="mdiTransfer"
                                             :disabled="!!editingTransaction || selectedTransferTransactionCount < 1"
                                             :title="tt('Batch Convert Transfer Transaction to Income Transaction')"
                                             @click="convertTransactionType(TransactionType.Transfer, TransactionType.Income)"></v-list-item>
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
                                :label="tt('File Type')"
                                :placeholder="tt('File Type')"
                                :items="allSupportedImportFileTypes"
                                v-model="fileType"
                            />
                        </v-col>

                        <v-col cols="12" md="12" v-if="allFileSubTypes">
                            <v-select
                                item-title="displayName"
                                item-value="type"
                                :disabled="submitting"
                                :label="tt('Format')"
                                :placeholder="tt('Format')"
                                :items="allFileSubTypes"
                                v-model="fileSubType"
                            />
                        </v-col>

                        <v-col cols="12" md="12" v-if="!isImportDataFromTextbox && allSupportedEncodings">
                            <v-select
                                item-title="displayName"
                                item-value="encoding"
                                :disabled="submitting"
                                :label="tt('File Encoding')"
                                :placeholder="tt('File Encoding')"
                                :items="allSupportedEncodings"
                                v-model="fileEncoding"
                            />
                        </v-col>

                        <v-col cols="12" md="12" v-if="!isImportDataFromTextbox">
                            <v-text-field
                                readonly
                                persistent-placeholder
                                type="text"
                                class="always-cursor-pointer"
                                :disabled="submitting"
                                :label="tt('Data File')"
                                :placeholder="tt('format.misc.clickToSelectedFile', { extensions: supportedImportFileExtensions })"
                                v-model="fileName"
                                @click="showOpenFileDialog"
                            />
                        </v-col>

                        <v-col cols="12" md="12" v-if="isImportDataFromTextbox">
                            <v-textarea
                                type="text"
                                persistent-placeholder
                                rows="5"
                                :disabled="submitting"
                                :placeholder="tt('Data to import')"
                                v-model="importData"
                            />
                        </v-col>

                        <v-col cols="12" md="12" class="mb-0 pb-0" v-if="exportFileGuideDocumentUrl">
                            <a :href="exportFileGuideDocumentUrl" :class="{ 'disabled': submitting }" target="_blank">
                                <v-icon :icon="mdiHelpCircleOutline" size="16" />
                                <span class="ml-1">{{ tt('How to export this file?') }}</span>
                                <span class="ml-1" v-if="exportFileGuideDocumentLanguageName">[{{ exportFileGuideDocumentLanguageName }}]</span>
                            </a>
                        </v-col>
                    </v-row>
                </v-window-item>
                <v-window-item value="defineColumn">
                    <v-data-table
                        fixed-header
                        fixed-footer
                        density="compact"
                        item-value="index"
                        :class="{ 'import-transaction-table': true, 'disabled': loading || submitting }"
                        :height="parsedFileLinesTableHeight"
                        :disable-sort="true"
                        :headers="parsedFileLinesHeaders"
                        :items="parsedFileLines"
                        :no-data-text="tt('No data to import')"
                        v-model:items-per-page="countPerPage"
                        v-model:page="currentPage"
                    >
                        <template #headers="{ columns }">
                            <tr>
                                <th class="text-no-wrap" :key="column.key ?? undefined" v-for="column in columns">
                                    <span v-if="!column.key || column.key === 'index'">{{ column.title }}</span>
                                    <div class="py-1" v-if="column.key && column.key !== 'index'">
                                        <span>{{ getParseDataMappedColumnDisplayName(parseInt(column.key)) }}</span>
                                        <br/>
                                        <span>({{ column.title }})</span>
                                        <v-menu activator="parent" location="bottom" max-height="500">
                                            <v-list>
                                                <v-list-item :key="columnType.type"
                                                             :append-icon="parsedFileDataColumnMapping[columnType.type] === parseInt(column.key) ? mdiCheck : undefined"
                                                             v-for="columnType in allImportTransactionColumnTypes"
                                                             @click="updateParseDataMappedColumn(parseInt(column.key), columnType.type)">
                                                    <v-list-item-title class="cursor-pointer">
                                                        {{ columnType.displayName }}
                                                    </v-list-item-title>
                                                </v-list-item>
                                            </v-list>
                                        </v-menu>
                                    </div>
                                </th>
                            </tr>
                        </template>
                        <template #bottom>
                            <div class="title-and-toolbar d-flex align-center text-no-wrap mt-2" v-if="parsedFileData">
                                <v-btn color="secondary" density="compact" variant="outlined"
                                       :append-icon="parsedFileIncludeHeader ? mdiCheck : mdiClose"
                                       @click="parsedFileIncludeHeader = !parsedFileIncludeHeader">{{ tt('Include Header Line') }}</v-btn>
                                <v-btn class="ml-2" color="secondary" density="compact" variant="outlined"
                                       :disabled="!parsedFileDataColumnMapping || !isNumber(parsedFileDataColumnMapping[ImportTransactionColumnType.TransactionType.type]) || !parsedFileAllTransactionTypes">
                                    <span>{{ tt('Transaction Type Mapping') }}</span>
                                    <span class="ml-1" v-if="parsedFileDataColumnMapping && isNumber(parsedFileDataColumnMapping[ImportTransactionColumnType.TransactionType.type]) && parsedFileAllTransactionTypes">({{ getObjectOwnFieldCount(parsedFileValidMappedTransactionTypes) || tt('None') }})</span>
                                    <v-menu eager activator="parent" location="bottom" max-height="500"
                                            :close-on-content-click="false">
                                        <v-list class="pa-0">
                                            <v-list-item class="pa-0">
                                                <v-table class="transaction-types-popup-menu">
                                                    <tbody>
                                                    <tr :key="typeName"
                                                        v-for="typeName in parsedFileAllTransactionTypes">
                                                        <td>{{ typeName }}</td>
                                                        <td>
                                                            <v-btn-toggle class="transaction-types-toggle" density="compact" variant="outlined"
                                                                          mandatory="force" divided
                                                                          v-model="parsedFileTransactionTypeMapping[typeName]">
                                                                <v-btn :value="undefined">{{ tt('None') }}</v-btn>
                                                                <v-btn :value="TransactionType.ModifyBalance">{{ tt('Modify Balance') }}</v-btn>
                                                                <v-btn :value="TransactionType.Income">{{ tt('Income') }}</v-btn>
                                                                <v-btn :value="TransactionType.Expense">{{ tt('Expense') }}</v-btn>
                                                                <v-btn :value="TransactionType.Transfer">{{ tt('Transfer') }}</v-btn>
                                                            </v-btn-toggle>
                                                        </td>
                                                    </tr>
                                                    </tbody>
                                                </v-table>
                                            </v-list-item>
                                        </v-list>
                                    </v-menu>
                                </v-btn>
                                <v-btn class="ml-2" color="secondary" density="compact" variant="outlined"
                                       :disabled="!parsedFileDataColumnMapping || !isNumber(parsedFileDataColumnMapping[ImportTransactionColumnType.TransactionTime.type])">
                                    <span>{{ tt('Time Format') }}</span>
                                    <span class="ml-1" v-if="parsedFileDataColumnMapping && isNumber(parsedFileDataColumnMapping[ImportTransactionColumnType.TransactionTime.type])">({{ parsedFileTimeFormat || parsedFileAutoDetectedTimeFormat || tt('Unknown') }})</span>
                                    <v-menu eager activator="parent" location="bottom" max-height="500">
                                        <v-list>
                                            <v-list-item key="auto"
                                                         :append-icon="parsedFileTimeFormat === '' ? mdiCheck : undefined"
                                                         @click="parsedFileTimeFormat = ''">
                                                <v-list-item-title class="cursor-pointer">
                                                    <span>{{ tt('Auto detect') }}</span>
                                                    <span class="ml-1" v-if="parsedFileAutoDetectedTimeFormat">({{ parsedFileAutoDetectedTimeFormat }})</span>
                                                    <span class="ml-1" v-if="!parsedFileAutoDetectedTimeFormat">({{ tt('Unknown') }})</span>
                                                </v-list-item-title>
                                            </v-list-item>
                                            <v-list-item :key="dateTimeFormat.format"
                                                         :append-icon="parsedFileTimeFormat === dateTimeFormat.format ? mdiCheck : undefined"
                                                         v-for="dateTimeFormat in KnownDateTimeFormat.values()"
                                                         @click="parsedFileTimeFormat = dateTimeFormat.format">
                                                <v-list-item-title class="cursor-pointer">
                                                    {{ dateTimeFormat.format }}
                                                </v-list-item-title>
                                            </v-list-item>
                                        </v-list>
                                    </v-menu>
                                </v-btn>
                                <v-btn class="ml-2" color="secondary" density="compact" variant="outlined"
                                       v-if="parsedFileDataColumnMapping && isNumber(parsedFileDataColumnMapping[ImportTransactionColumnType.TransactionTimezone.type])">
                                    <span>{{ tt('Timezone Format') }}</span>
                                    <span class="ml-1" v-if="parsedFileDataColumnMapping && isNumber(parsedFileDataColumnMapping[ImportTransactionColumnType.TransactionTimezone.type])">({{ KnownDateTimezoneFormat.valueOf(parsedFileTimezoneFormat || parsedFileAutoDetectedTimezoneFormat || '')?.name || tt('Unknown') }})</span>
                                    <v-menu eager activator="parent" location="bottom" max-height="500">
                                        <v-list>
                                            <v-list-item key="auto"
                                                         :append-icon="parsedFileTimezoneFormat === '' ? mdiCheck : undefined"
                                                         @click="parsedFileTimezoneFormat = ''">
                                                <v-list-item-title class="cursor-pointer">
                                                    <span>{{ tt('Auto detect') }}</span>
                                                    <span class="ml-1" v-if="parsedFileAutoDetectedTimezoneFormat && KnownDateTimezoneFormat.valueOf(parsedFileAutoDetectedTimezoneFormat || '')">({{ KnownDateTimezoneFormat.valueOf(parsedFileAutoDetectedTimezoneFormat || '')?.name }})</span>
                                                    <span class="ml-1" v-if="!parsedFileAutoDetectedTimezoneFormat || !KnownDateTimezoneFormat.valueOf(parsedFileAutoDetectedTimezoneFormat || '')">({{ tt('Unknown') }})</span>
                                                </v-list-item-title>
                                            </v-list-item>
                                            <v-list-item :key="timezoneFormat.value"
                                                         :append-icon="parsedFileTimezoneFormat === timezoneFormat.value ? mdiCheck : undefined"
                                                         v-for="timezoneFormat in KnownDateTimezoneFormat.values()"
                                                         @click="parsedFileTimezoneFormat = timezoneFormat.value">
                                                <v-list-item-title class="cursor-pointer">
                                                    {{ timezoneFormat.name }}
                                                </v-list-item-title>
                                            </v-list-item>
                                        </v-list>
                                    </v-menu>
                                </v-btn>
                                <v-btn class="ml-2" color="secondary" density="compact" variant="outlined"
                                       :disabled="!parsedFileDataColumnMapping || !isNumber(parsedFileDataColumnMapping[ImportTransactionColumnType.Amount.type])">
                                    <span>{{ tt('Amount Format') }}</span>
                                    <span class="ml-1" v-if="parsedFileDataColumnMapping && isNumber(parsedFileDataColumnMapping[ImportTransactionColumnType.Amount.type])">({{ KnownAmountFormat.valueOf(parsedFileAmountFormat || parsedFileAutoDetectedAmountFormat || '')?.format || tt('Unknown') }})</span>
                                    <v-menu eager activator="parent" location="bottom" max-height="500">
                                        <v-list>
                                            <v-list-item key="auto"
                                                         :append-icon="parsedFileAmountFormat === '' ? mdiCheck : undefined"
                                                         @click="parsedFileAmountFormat = ''">
                                                <v-list-item-title class="cursor-pointer">
                                                    <span>{{ tt('Auto detect') }}</span>
                                                    <span class="ml-1" v-if="parsedFileAutoDetectedAmountFormat && KnownAmountFormat.valueOf(parsedFileAutoDetectedAmountFormat || '')">({{ KnownAmountFormat.valueOf(parsedFileAutoDetectedAmountFormat || '')?.format }})</span>
                                                    <span class="ml-1" v-if="!parsedFileAutoDetectedAmountFormat || !KnownAmountFormat.valueOf(parsedFileAutoDetectedAmountFormat || '')">({{ tt('Unknown') }})</span>
                                                </v-list-item-title>
                                            </v-list-item>
                                            <v-list-item :key="amountFormat.type"
                                                         :append-icon="parsedFileAmountFormat === amountFormat.type ? mdiCheck : undefined"
                                                         v-for="amountFormat in KnownAmountFormat.values()"
                                                         @click="parsedFileAmountFormat = amountFormat.type">
                                                <v-list-item-title class="cursor-pointer">
                                                    {{ amountFormat.format }}
                                                </v-list-item-title>
                                            </v-list-item>
                                        </v-list>
                                    </v-menu>
                                </v-btn>
                                <v-btn class="ml-2" color="secondary" density="compact" variant="outlined"
                                       v-if="parsedFileDataColumnMapping && isNumber(parsedFileDataColumnMapping[ImportTransactionColumnType.GeographicLocation.type])">
                                    <span>{{ tt('Geographic Location Separator') }}</span>
                                    <span class="ml-1" v-if="parsedFileGeoLocationSeparator">({{ parsedFileGeoLocationSeparator }})</span>
                                    <v-menu eager activator="parent" location="bottom" max-height="500">
                                        <v-list>
                                            <v-list-item :key="separator.value"
                                                         :append-icon="parsedFileGeoLocationSeparator === separator.value ? mdiCheck : undefined"
                                                         v-for="separator in allSeparators"
                                                         @click="parsedFileGeoLocationSeparator = separator.value">
                                                <v-list-item-title class="cursor-pointer">
                                                    {{ separator.name }} ({{separator.value}})
                                                </v-list-item-title>
                                            </v-list-item>
                                        </v-list>
                                    </v-menu>
                                </v-btn>
                                <v-btn class="ml-2" color="secondary" density="compact" variant="outlined"
                                       v-if="parsedFileDataColumnMapping && isNumber(parsedFileDataColumnMapping[ImportTransactionColumnType.Tags.type])">
                                    <span>{{ tt('Transaction Tags Separator') }}</span>
                                    <span class="ml-1" v-if="parsedFileTagSeparator">({{ parsedFileTagSeparator }})</span>
                                    <v-menu eager activator="parent" location="bottom" max-height="500">
                                        <v-list>
                                            <v-list-item :key="separator.value"
                                                         :append-icon="parsedFileTagSeparator === separator.value ? mdiCheck : undefined"
                                                         v-for="separator in allSeparators"
                                                         @click="parsedFileTagSeparator = separator.value">
                                                <v-list-item-title class="cursor-pointer">
                                                    {{ separator.name }} ({{separator.value}})
                                                </v-list-item-title>
                                            </v-list-item>
                                        </v-list>
                                    </v-menu>
                                </v-btn>
                                <v-spacer/>
                                <span>{{ tt('Lines Per Page') }}</span>
                                <v-select class="ml-2" density="compact" max-width="100"
                                          item-title="title"
                                          item-value="value"
                                          :disabled="loading || submitting"
                                          :items="parsedFileLinesTablePageOptions"
                                          v-model="countPerPage"
                                />
                                <pagination-buttons density="compact"
                                                    :disabled="loading || submitting"
                                                    :totalPageCount="Math.ceil((parsedFileLines ? parsedFileLines.length : 0) / countPerPage)"
                                                    v-model="currentPage"></pagination-buttons>
                            </div>
                        </template>
                    </v-data-table>
                </v-window-item>
                <v-window-item value="checkData">
                    <v-data-table
                        fixed-header
                        fixed-footer
                        show-select
                        multi-sort
                        density="compact"
                        item-value="index"
                        :class="{ 'import-transaction-table': true, 'disabled': loading || submitting }"
                        :height="importTransactionsTableHeight"
                        :headers="importTransactionHeaders"
                        :items="importTransactions"
                        :search="JSON.stringify(filters)"
                        :custom-filter="importTransactionsFilter"
                        :no-data-text="tt('No data to import')"
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
                                        <v-list-item :prepend-icon="mdiSelectAll"
                                                     :title="tt('Select All Valid Items')"
                                                     :disabled="loading || submitting"
                                                     @click="selectAllValid"></v-list-item>
                                        <v-list-item :prepend-icon="mdiSelectAll"
                                                     :title="tt('Select All Invalid Items')"
                                                     :disabled="loading || submitting"
                                                     @click="selectAllInvalid"></v-list-item>
                                        <v-divider class="my-2"/>
                                        <v-list-item :prepend-icon="mdiSelectAll"
                                                     :title="tt('Select All')"
                                                     :disabled="loading || submitting"
                                                     @click="selectAll"></v-list-item>
                                        <v-list-item :prepend-icon="mdiSelect"
                                                     :title="tt('Select None')"
                                                     :disabled="loading || submitting"
                                                     @click="selectNone"></v-list-item>
                                        <v-list-item :prepend-icon="mdiSelectInverse"
                                                     :title="tt('Invert Selection')"
                                                     :disabled="loading || submitting"
                                                     @click="selectInvert"></v-list-item>
                                        <v-divider class="my-2"/>
                                        <v-list-item :prepend-icon="mdiSelectAll"
                                                     :title="tt('Select All in This Page')"
                                                     :disabled="loading || submitting"
                                                     @click="selectAllInThisPage"></v-list-item>
                                        <v-list-item :prepend-icon="mdiSelect"
                                                     :title="tt('Select None in This Page')"
                                                     :disabled="loading || submitting"
                                                     @click="selectNoneInThisPage"></v-list-item>
                                        <v-list-item :prepend-icon="mdiSelectInverse"
                                                     :title="tt('Invert Selection in This Page')"
                                                     :disabled="loading || submitting"
                                                     @click="selectInvertInThisPage"></v-list-item>
                                    </v-list>
                                </v-menu>
                            </v-checkbox>
                        </template>
                        <template #item.data-table-select="{ item }">
                            <v-checkbox density="compact"
                                        :color="!item.valid ? 'error' : 'primary'"
                                        :disabled="loading || submitting"
                                        v-model="item.selected"></v-checkbox>
                        </template>
                        <template #item.valid="{ item }">
                            <v-icon size="small" :class="{ 'text-error': !item.valid }"
                                    :disabled="loading || submitting"
                                    :icon="editingTransaction === item ? mdiCheck : mdiPencilOutline"
                                    @click="editTransaction(item)">
                            </v-icon>
                            <v-tooltip activator="parent" v-if="!loading && !submitting">{{ tt('Edit') }}</v-tooltip>
                        </template>
                        <template #item.time="{ item }">
                            <span>{{ getDisplayDateTime(item) }}</span>
                            <v-chip class="ml-1" variant="flat" color="secondary" size="x-small"
                                    v-if="item.utcOffset !== currentTimezoneOffsetMinutes">{{ getDisplayTimezone(item) }}</v-chip>
                        </template>
                        <template #item.type="{ value }">
                            <v-chip label color="secondary" variant="outlined" size="x-small" v-if="value === TransactionType.ModifyBalance">{{ tt('Modify Balance') }}</v-chip>
                            <v-chip label class="text-income" variant="outlined" size="x-small" v-else-if="value === TransactionType.Income">{{ tt('Income') }}</v-chip>
                            <v-chip label class="text-expense" variant="outlined" size="x-small" v-else-if="value === TransactionType.Expense">{{ tt('Expense') }}</v-chip>
                            <v-chip label color="primary" variant="outlined" size="x-small" v-else-if="value === TransactionType.Transfer">{{ tt('Transfer') }}</v-chip>
                            <v-chip label color="default" variant="outlined" size="x-small" v-else>{{ tt('Unknown') }}</v-chip>
                        </template>
                        <template #item.actualCategoryName="{ item }">
                            <div class="d-flex align-center" v-if="editingTransaction !== item || item.type === TransactionType.ModifyBalance">
                                <span v-if="item.type === TransactionType.ModifyBalance">-</span>
                                <ItemIcon size="24px" icon-type="category"
                                          :icon-id="allCategoriesMap[item.categoryId].icon"
                                          :color="allCategoriesMap[item.categoryId].color"
                                          v-if="item.type !== TransactionType.ModifyBalance && item.categoryId && item.categoryId !== '0' && allCategoriesMap[item.categoryId]"></ItemIcon>
                                <span class="ml-2" v-if="item.type !== TransactionType.ModifyBalance && item.categoryId && item.categoryId !== '0' && allCategoriesMap[item.categoryId]">
                                    {{ allCategoriesMap[item.categoryId].name }}
                                </span>
                                <div class="text-error font-italic" v-else-if="item.type !== TransactionType.ModifyBalance && (!item.categoryId || item.categoryId === '0' || !allCategoriesMap[item.categoryId])">
                                    <v-icon class="mr-1" :icon="mdiAlertOutline"/>
                                    <span>{{ item.originalCategoryName }}</span>
                                </div>
                            </div>
                            <div style="width: 260px" v-if="editingTransaction === item && item.type === TransactionType.Expense">
                                <two-column-select density="compact" variant="plain"
                                                   primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                   primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                   primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                   secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                   secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                   secondary-hidden-field="hidden"
                                                   :disabled="loading || submitting || !hasAvailableExpenseCategories"
                                                   :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                                   :show-selection-primary-text="true"
                                                   :custom-selection-primary-text="getTransactionPrimaryCategoryName(item.categoryId, allCategories[CategoryType.Expense])"
                                                   :custom-selection-secondary-text="getTransactionSecondaryCategoryName(item.categoryId, allCategories[CategoryType.Expense])"
                                                   :placeholder="tt('Category')"
                                                   :items="allCategories[CategoryType.Expense]"
                                                   v-model="item.categoryId">
                                </two-column-select>
                            </div>
                            <div style="width: 260px" v-if="editingTransaction === item && item.type === TransactionType.Income">
                                <two-column-select density="compact" variant="plain"
                                                   primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                   primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                   primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                   secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                   secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                   secondary-hidden-field="hidden"
                                                   :disabled="loading || submitting || !hasAvailableIncomeCategories"
                                                   :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                                   :show-selection-primary-text="true"
                                                   :custom-selection-primary-text="getTransactionPrimaryCategoryName(item.categoryId, allCategories[CategoryType.Income])"
                                                   :custom-selection-secondary-text="getTransactionSecondaryCategoryName(item.categoryId, allCategories[CategoryType.Income])"
                                                   :placeholder="tt('Category')"
                                                   :items="allCategories[CategoryType.Income]"
                                                   v-model="item.categoryId">
                                </two-column-select>
                            </div>
                            <div style="width: 260px" v-if="editingTransaction === item && item.type === TransactionType.Transfer">
                                <two-column-select density="compact" variant="plain"
                                                   primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                   primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                   primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                   secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                   secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                   secondary-hidden-field="hidden"
                                                   :disabled="loading || submitting || !hasAvailableTransferCategories"
                                                   :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                                   :show-selection-primary-text="true"
                                                   :custom-selection-primary-text="getTransactionPrimaryCategoryName(item.categoryId, allCategories[CategoryType.Transfer])"
                                                   :custom-selection-secondary-text="getTransactionSecondaryCategoryName(item.categoryId, allCategories[CategoryType.Transfer])"
                                                   :placeholder="tt('Category')"
                                                   :items="allCategories[CategoryType.Transfer]"
                                                   v-model="item.categoryId">
                                </two-column-select>
                            </div>
                        </template>
                        <template #item.sourceAmount="{ item }">
                            <span>{{ getTransactionDisplayAmount(item) }}</span>
                            <v-icon class="mx-1" size="13" :icon="mdiArrowRight" v-if="item.type === TransactionType.Transfer && item.sourceAccountId !== item.destinationAccountId"></v-icon>
                            <span v-if="item.type === TransactionType.Transfer && item.sourceAccountId !== item.destinationAccountId">{{ getTransactionDisplayDestinationAmount(item) }}</span>
                        </template>
                        <template #item.actualSourceAccountName="{ item }">
                            <div class="d-flex align-center" v-if="editingTransaction !== item">
                                <span v-if="item.sourceAccountId && item.sourceAccountId !== '0' && allAccountsMap[item.sourceAccountId]">{{ allAccountsMap[item.sourceAccountId].name }}</span>
                                <div class="text-error font-italic" v-else>
                                    <v-icon class="mr-1" :icon="mdiAlertOutline"/>
                                    <span>{{ item.originalSourceAccountName }}</span>
                                </div>
                                <v-icon class="mx-1" size="13" :icon="mdiArrowRight" v-if="item.type === TransactionType.Transfer"></v-icon>
                                <span v-if="item.type === TransactionType.Transfer && item.destinationAccountId && item.destinationAccountId !== '0' && allAccountsMap[item.destinationAccountId]">{{allAccountsMap[item.destinationAccountId].name }}</span>
                                <div class="text-error font-italic" v-else-if="item.type === TransactionType.Transfer && (!item.destinationAccountId || item.destinationAccountId === '0' || !allAccountsMap[item.destinationAccountId])">
                                    <v-icon class="mr-1" :icon="mdiAlertOutline"/>
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
                                                   :enable-filter="true" :filter-placeholder="tt('Find account')" :filter-no-items-text="tt('No available account')"
                                                   :custom-selection-primary-text="getSourceAccountDisplayName(item)"
                                                   :placeholder="getSourceAccountTitle(item)"
                                                   :items="allVisibleCategorizedAccounts"
                                                   v-model="item.sourceAccountId">
                                </two-column-select>
                                <v-icon class="mx-1" size="13" :icon="mdiArrowRight" v-if="item.type === TransactionType.Transfer"></v-icon>
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
                                                   :enable-filter="true" :filter-placeholder="tt('Find account')" :filter-no-items-text="tt('No available account')"
                                                   :custom-selection-primary-text="getDestinationAccountDisplayName(item)"
                                                   :placeholder="tt('Destination Account')"
                                                   :items="allVisibleCategorizedAccounts"
                                                   v-model="item.destinationAccountId"
                                                   v-if="item.type === TransactionType.Transfer">
                                </two-column-select>
                            </div>
                        </template>
                        <template #item.geoLocation="{ item }">
                            <span v-if="item.geoLocation">{{ `(${item.geoLocation.longitude}, ${item.geoLocation.latitude})` }}</span>
                            <span v-else-if="!item.geoLocation">{{ tt('None') }}</span>
                        </template>
                        <template #item.tagIds="{ item }">
                            <div v-if="editingTransaction !== item">
                                <v-chip class="transaction-tag" size="small"
                                        :class="{ 'font-italic': !tagId || tagId === '0' || !allTagsMap[tagId] }"
                                        :prepend-icon="tagId && tagId !== '0' && allTagsMap[tagId] ? mdiPound : mdiAlertOutline"
                                        :color="tagId && tagId !== '0' && allTagsMap[tagId] ? 'default' : 'error'"
                                        :text="tagId && tagId !== '0' && allTagsMap[tagId] ? allTagsMap[tagId].name : item.originalTagNames[index]"
                                        :key="tagId"
                                        v-for="(tagId, index) in item.tagIds"/>
                                <v-chip class="transaction-tag" size="small"
                                        :text="tt('None')"
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
                                    :placeholder="tt('None')"
                                    :items="allTags"
                                    :no-data-text="tt('No available tag')"
                                    v-model="editingTags"
                                >
                                    <template #chip="{ props, index }">
                                        <v-chip :class="{ 'font-italic': !isTagValid(editingTags, index) }"
                                                :prepend-icon="isTagValid(editingTags, index) ? mdiPound : mdiAlertOutline"
                                                :color="isTagValid(editingTags, index) ? 'default' : 'error'"
                                                :text="isTagValid(editingTags, index) ? allTagsMap[editingTags[index]].name : item.originalTagNames[index]"
                                                v-bind="props"/>
                                    </template>

                                    <template #item="{ props, item }">
                                        <v-list-item :value="item.value" v-bind="props" v-if="!item.raw.hidden">
                                            <template #title>
                                                <v-list-item-title>
                                                    <div class="d-flex align-center">
                                                        <v-icon size="20" start :icon="mdiPound"/>
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
                                    {{ tt('format.misc.selectedCount', { count: selectedImportTransactionCount, totalCount: importTransactions.length }) }}
                                </span>
                                <v-spacer/>
                                <span>{{ tt('Transactions Per Page') }}</span>
                                <v-select class="ml-2" density="compact" max-width="100"
                                          item-title="title"
                                          item-value="value"
                                          :disabled="loading || submitting"
                                          :items="importTransactionsTablePageOptions"
                                          v-model="countPerPage"
                                />
                                <pagination-buttons density="compact"
                                                    :disabled="loading || submitting"
                                                    :totalPageCount="totalPageCount"
                                                    v-model="currentPage"></pagination-buttons>
                            </div>
                        </template>
                    </v-data-table>
                </v-window-item>
                <v-window-item value="finalResult">
                    <h4 class="text-h4 mb-1">{{ tt('Data Import Completed') }}</h4>
                    <p class="my-5">{{ tt('format.misc.importTransactionResult', { count: importedCount }) }}</p>
                </v-window-item>
            </v-window>

            <div class="d-flex justify-sm-space-between gap-4 flex-wrap justify-center mt-10">
                <v-btn color="secondary" variant="tonal" :disabled="loading || submitting"
                       :prepend-icon="mdiClose" @click="close(false)"
                       v-if="currentStep !== 'finalResult'">{{ tt('Cancel') }}</v-btn>
                <v-btn color="primary" :disabled="loading || submitting || (!isImportDataFromTextbox && !importFile) || (isImportDataFromTextbox && !importData)"
                       :append-icon="!submitting ? mdiArrowRight : undefined" @click="parseData"
                       v-if="currentStep === 'defineColumn' || currentStep === 'uploadFile'">
                    {{ tt('Next') }}
                    <v-progress-circular indeterminate size="22" class="ml-2" v-if="submitting"></v-progress-circular>
                </v-btn>
                <v-btn color="teal" :disabled="submitting || !!editingTransaction || selectedImportTransactionCount < 1 || selectedInvalidTransactionCount > 0"
                       :append-icon="!submitting ? mdiArrowRight : undefined" @click="submit"
                       v-if="currentStep === 'checkData'">
                    {{ (submitting && importProcess > 0 ? tt('format.misc.importingTransactions', { process: formatNumber(importProcess, 2) }) : tt('Import')) }}
                    <v-progress-circular indeterminate size="22" class="ml-2" v-if="submitting"></v-progress-circular>
                </v-btn>
                <v-btn color="secondary" variant="tonal"
                       :append-icon="mdiCheck"
                       @click="close(true)"
                       v-if="currentStep === 'finalResult'">{{ tt('Close') }}</v-btn>
            </div>
        </v-card>
    </v-dialog>

    <v-dialog width="640" v-model="showCustomDescriptionDialog">
        <v-card class="pa-2 pa-sm-4 pa-md-4">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <h4 class="text-h4">{{ tt('Filter Description') }}</h4>
                </div>
            </template>
            <v-card-text class="mb-md-4 w-100 d-flex justify-center">
                <v-text-field
                    type="text"
                    persistent-placeholder
                    :label="tt('Description')"
                    :placeholder="tt('Description')"
                    v-model="currentDescriptionFilterValue"
                />
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center gap-4">
                    <v-btn :disabled="!currentDescriptionFilterValue" @click="showCustomDescriptionDialog = false; filters.description = currentDescriptionFilterValue">{{ tt('OK') }}</v-btn>
                    <v-btn color="secondary" variant="tonal" @click="showCustomDescriptionDialog = false; currentDescriptionFilterValue = ''">{{ tt('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <date-range-selection-dialog :title="tt('Custom Date Range')"
                                 :min-time="filters.minDatetime"
                                 :max-time="filters.maxDatetime"
                                 v-model:show="showCustomDateRangeDialog"
                                 @dateRange:change="changeCustomDateFilter"
                                 @error="onShowDateRangeError" />
    <batch-replace-dialog ref="batchReplaceDialog" />
    <BatchCreateDialog ref="batchCreateDialog" />
    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
    <input ref="fileInput" type="file" style="display: none" :accept="supportedImportFileExtensions" @change="setImportFile($event)" />
</template>

<script setup lang="ts">
import type { StepBarItem } from '@/components/desktop/StepsBar.vue';
import PaginationButtons from '@/components/desktop/PaginationButtons.vue';
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';
import BatchReplaceDialog, { type BatchReplaceDialogDataType } from './dialogs/BatchReplaceDialog.vue';
import BatchCreateDialog, { type BatchCreateDialogDataType } from './dialogs/BatchCreateDialog.vue';

import { ref, computed, useTemplateRef, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useOverviewStore } from '@/stores/overview.ts';
import { useStatisticsStore } from '@/stores/statistics.ts';

import type { NameValue, TypeAndDisplayName } from '@/core/base.ts';
import { KnownAmountFormat } from '@/core/numeral.ts';
import { KnownDateTimeFormat } from '@/core/datetime.ts';
import { KnownDateTimezoneFormat } from '@/core/timezone.ts';
import { CategoryType } from '@/core/category.ts';
import { TransactionType, ImportTransactionColumnType } from '@/core/transaction.ts';
import type { LocalizedImportFileType, LocalizedImportFileTypeSubType, LocalizedImportFileTypeSupportedEncodings } from '@/core/file.ts';
import { Account, type CategorizedAccountWithDisplayBalance } from '@/models/account.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';
import type { TransactionTag } from '@/models/transaction_tag.ts';
import { ImportTransaction } from '@/models/imported_transaction.ts';

import {
    isString,
    isNumber,
    isObjectEmpty,
    getObjectOwnFieldCount,
    findDisplayNameByType,
    objectFieldToArrayItem
} from '@/lib/common.ts';
import {
    findExtensionByType,
    isFileExtensionSupported
} from '@/lib/file.ts';
import { generateRandomUUID } from '@/lib/misc.ts';
import logger from '@/lib/logger.ts';
import {
    parseDateFromUnixTime,
    getUnixTime,
    getUtcOffsetByUtcOffsetMinutes,
    getTimezoneOffsetMinutes
} from '@/lib/datetime.ts';
import {
    getAccountMapByName
} from '@/lib/account.ts';
import {
    transactionTypeToCategoryType,
    getSecondaryTransactionMapByName,
    getTransactionPrimaryCategoryName,
    getTransactionSecondaryCategoryName
} from '@/lib/category.ts';

import {
    mdiFilterOutline,
    mdiCheck,
    mdiDotsVertical,
    mdiHelpCircleOutline,
    mdiFindReplace,
    mdiShapePlusOutline,
    mdiTransfer,
    mdiClose,
    mdiArrowRight,
    mdiSelectAll,
    mdiSelect,
    mdiSelectInverse,
    mdiPencilOutline,
    mdiAlertOutline,
    mdiPound
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;
type BatchReplaceDialogType = InstanceType<typeof BatchReplaceDialog>;
type BatchCreateDialogType = InstanceType<typeof BatchCreateDialog>;

type ImportTransactionDialogStep = 'uploadFile' | 'defineColumn' | 'checkData' | 'finalResult';

interface ImportTransactionDialogFilter {
    minDatetime: number | null; // minDatetime or maxDatetime is null for 'All Date Range', all are not null for 'Custom Date Range'
    maxDatetime: number | null;
    transactionType: TransactionType | null; // null for 'All Transaction Type'
    category: string | null | undefined; // null for 'All Category', undefined for 'Invalid Category'
    account: string | null | undefined; // null for 'All Account', undefined for 'Invalid Account'
    tag: string | null | undefined; // null for 'All Tag', undefined for 'Invalid Tag'
    description: string | null; // null for 'All Description'
}

interface ImportTransactionsDialogTablePageOption {
    value: number;
    title: string;
}

defineProps<{
    persistent?: boolean;
}>();

const {
    tt,
    getAllImportTransactionColumnTypes,
    getAllSupportedImportFileTypes,
    formatUnixTimeToLongDateTime,
    formatAmountWithCurrency,
    formatNumber,
    getCategorizedAccountsWithDisplayBalance
} = useI18n();

const settingsStore = useSettingsStore();
const userStore = useUserStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();
const transactionsStore = useTransactionsStore();
const overviewStore = useOverviewStore();
const statisticsStore = useStatisticsStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');
const batchReplaceDialog = useTemplateRef<BatchReplaceDialogType>('batchReplaceDialog');
const batchCreateDialog = useTemplateRef<BatchCreateDialogType>('batchCreateDialog');
const fileInput = useTemplateRef<HTMLInputElement>('fileInput');

const showState = ref<boolean>(false);
const clientSessionId = ref<string>('');
const currentStep = ref<ImportTransactionDialogStep>('uploadFile');
const importProcess = ref<number>(0);
const fileType = ref<string>('ezbookkeeping');
const fileSubType = ref<string>('ezbookkeeping_csv');
const fileEncoding = ref<string>('utf-8');
const importFile = ref<File | null>(null);
const importData = ref<string>('');
const parsedFileData = ref<string[][] | undefined>(undefined);
const parsedFileIncludeHeader = ref<boolean>(true);
const parsedFileDataColumnMapping = ref<Record<number, number>>({});
const parsedFileTransactionTypeMapping = ref<Record<string, TransactionType>>({});
const parsedFileTimeFormat = ref<string>('');
const parsedFileTimezoneFormat = ref<string>('');
const parsedFileAmountFormat = ref<string>('');
const parsedFileGeoLocationSeparator = ref<string>(' ');
const parsedFileTagSeparator = ref<string>(';');
const importTransactions = ref<ImportTransaction[] | undefined>(undefined);
const editingTransaction = ref<ImportTransaction | null>(null);
const editingTags = ref<string[]>([]);
const filters = ref<ImportTransactionDialogFilter>({
    minDatetime: null,
    maxDatetime: null,
    transactionType: null,
    category: null,
    account: null,
    tag: null,
    description: null
});

const currentPage = ref<number>(1);
const countPerPage = ref<number>(10);
const importedCount = ref<number | null>(null);
const showCustomDateRangeDialog = ref<boolean>(false);
const showCustomDescriptionDialog = ref<boolean>(false);
const currentDescriptionFilterValue = ref<string | null>(null);
const loading = ref<boolean>(true);
const submitting = ref<boolean>(false);

let resolveFunc: (() => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const showAccountBalance = computed<boolean>(() => settingsStore.appSettings.showAccountBalance);
const currentTimezoneOffsetMinutes = computed<number>(() => getTimezoneOffsetMinutes(settingsStore.appSettings.timeZone));

const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);

const allSteps = computed<StepBarItem[]>(() => {
    const steps: StepBarItem[] = [
        {
            name: 'uploadFile',
            title: tt('Upload File'),
            subTitle: tt('Upload Transaction Data File')
        }
    ];

    if (fileType.value === 'dsv' || fileType.value === 'dsv_data') {
        steps.push({
            name: 'defineColumn',
            title: tt('Define Column'),
            subTitle: tt('Define and Check Column Mapping')
        });
    }

    steps.push(...[
        {
            name: 'checkData',
            title: tt('Check & Modify'),
            subTitle: tt('Check and Modify Your Data')
        },
        {
            name: 'finalResult',
            title: tt('Complete'),
            subTitle: tt('Data Import Completed')
        }
    ]);

    return steps;
});

const allImportTransactionColumnTypes = computed<TypeAndDisplayName[]>(() => getAllImportTransactionColumnTypes());
const allSupportedImportFileTypes = computed<LocalizedImportFileType[]>(() => getAllSupportedImportFileTypes());

const allSeparators = computed<NameValue[]>(() => {
    const separators: NameValue[] = [
        {
            name: tt('Space'),
            value: ' '
        },
        {
            name: tt('Comma'),
            value: ','
        },
        {
            name: tt('Semicolon'),
            value: ';'
        },
        {
            name: tt('Tab'),
            value: '\t'
        },
        {
            name: tt('Vertical Bar'),
            value: '|'
        }
    ];

    return separators;
});

const isImportDataFromTextbox = computed<boolean>(() => {
    for (const importFileType of allSupportedImportFileTypes.value) {
        if (importFileType.type === fileType.value) {
            return !!importFileType.dataFromTextbox;
        }
    }

    return false;
});

const allFileSubTypes = computed<LocalizedImportFileTypeSubType[] | undefined>(() => {
    for (const importFileType of allSupportedImportFileTypes.value) {
        if (importFileType.type === fileType.value) {
            return importFileType.subTypes;
        }
    }

    return undefined;
});

const allSupportedEncodings = computed<LocalizedImportFileTypeSupportedEncodings[] | undefined>(() => {
    for (const importFileType of allSupportedImportFileTypes.value) {
        if (importFileType.type === fileType.value) {
            return importFileType.supportedEncodings;
        }
    }

    return undefined;
});

const allAccounts = computed<Account[]>(() => accountsStore.allPlainAccounts);
const allVisibleAccounts = computed<Account[]>(() => accountsStore.allVisiblePlainAccounts);
const allVisibleCategorizedAccounts = computed<CategorizedAccountWithDisplayBalance[]>(() => getCategorizedAccountsWithDisplayBalance(allVisibleAccounts.value, showAccountBalance.value));
const allAccountsMap = computed<Record<string, Account>>(() => accountsStore.allAccountsMap);
const allAccountsMapByName = computed<Record<string, Account>>(() => getAccountMapByName(accountsStore.allAccounts));
const allCategories = computed<Record<number, TransactionCategory[]>>(() => transactionCategoriesStore.allTransactionCategories);
const allCategoriesMap = computed<Record<string, TransactionCategory>>(() => transactionCategoriesStore.allTransactionCategoriesMap);
const allTags = computed<TransactionTag[]>(() => transactionTagsStore.allTransactionTags);
const allTagsMap = computed<Record<string, TransactionTag>>(() => transactionTagsStore.allTransactionTagsMap);

const hasAvailableExpenseCategories = computed<boolean>(() => transactionCategoriesStore.hasAvailableExpenseCategories);
const hasAvailableIncomeCategories = computed<boolean>(() => transactionCategoriesStore.hasAvailableIncomeCategories);
const hasAvailableTransferCategories = computed<boolean>(() => transactionCategoriesStore.hasAvailableTransferCategories);

const supportedImportFileExtensions = computed<string | undefined>(() => {
    if (allFileSubTypes.value && allFileSubTypes.value.length) {
        const subTypeExtensions = findExtensionByType(allFileSubTypes.value, fileSubType.value);

        if (subTypeExtensions) {
            return subTypeExtensions;
        }
    }

    return findExtensionByType(allSupportedImportFileTypes.value, fileType.value);
});

const exportFileGuideDocumentUrl = computed<string | undefined>(() => {
    for (const importFileType of allSupportedImportFileTypes.value) {
        if (importFileType.type === fileType.value) {
            const document = importFileType.document;

            if (!document) {
                return undefined;
            }

            const language = document.language ? document.language + '/' : '';
            const anchor = document.anchor ? '#' + document.anchor : '';
            return `https://ezbookkeeping.mayswind.net/${language}export_and_import${anchor}`;
        }
    }

    return undefined;
});

const exportFileGuideDocumentLanguageName = computed<string | undefined>(() => {
    for (const importFileType of allSupportedImportFileTypes.value) {
        if (importFileType.type === fileType.value) {
            const document = importFileType.document;
            return document?.displayLanguageName;
        }
    }

    return undefined;
});

const fileName = computed<string>(() => importFile.value?.name || '');

const parsedFileLines = computed<Record<string, string>[] | undefined>(() => {
    if (!parsedFileData.value) {
        return undefined;
    }

    const allLines: Record<string, string>[] = [];
    const startIndex = parsedFileIncludeHeader.value ? 1 : 0;

    for (let i = startIndex, index = 1; i < parsedFileData.value.length; i++, index++) {
        const line: Record<string, string> = {};
        const columns = parsedFileData.value[i];

        for (let j = 0; j < columns.length; j++) {
            line['index'] = index.toString();
            line[`column${j + 1}`] = columns[j];
        }

        allLines.push(line);
    }

    return allLines;
});

const parsedFileLinesTableHeight = computed<number | undefined>(() => {
    if (countPerPage.value <= 10 || !parsedFileLines.value || parsedFileLines.value.length <= 10) {
        return undefined;
    } else {
        return 400;
    }
});

const parsedFileLinesHeaders = computed<object[]>(() => {
    let maxColumnCount = 0;

    if (parsedFileData.value) {
        for (let i = 0; i < parsedFileData.value.length; i++) {
            if (parsedFileData.value[i].length > maxColumnCount) {
                maxColumnCount = parsedFileData.value[i].length;
            }
        }
    }

    const headers: object[] = [];

    headers.push({ key: 'index', value: 'index', title: '#', sortable: true, nowrap: true });

    for (let i = 0; i < maxColumnCount; i++) {
        let title = `#${i + 1}`;

        if (parsedFileIncludeHeader.value && parsedFileData.value && parsedFileData.value[0][i]) {
            title = parsedFileData.value[0][i];
        }

        headers.push({ key: i.toString(), value: `column${i + 1}`, title: title, sortable: true, nowrap: true });
    }

    return headers;
});

const parsedFileLinesTablePageOptions = computed<ImportTransactionsDialogTablePageOption[]>(() => getTablePageOptions(parsedFileLines.value?.length));

const parsedFileAllTransactionTypes = computed<string[]>(() => {
    if (!parsedFileData.value || !parsedFileData.value.length || !isNumber(parsedFileDataColumnMapping.value[ImportTransactionColumnType.TransactionType.type])) {
        return [];
    }

    const allTypeMap: Record<string, boolean> = {};
    const allTypes: string[] = [];
    const typeColumnIndex = parsedFileDataColumnMapping.value[ImportTransactionColumnType.TransactionType.type];

    const startIndex = parsedFileIncludeHeader.value ? 1 : 0;

    for (let i = startIndex; i < parsedFileData.value.length; i++) {
        if (parsedFileData.value[i].length <= typeColumnIndex) {
            continue;
        }

        const type = parsedFileData.value[i][typeColumnIndex];

        if (type && !allTypeMap[type]) {
            allTypes.push(type);
            allTypeMap[type] = true;
        }
    }

    return allTypes;
});

const parsedFileValidMappedTransactionTypes = computed<Record<string, TransactionType>>(() => {
    if (!parsedFileData.value || !parsedFileData.value.length || !isNumber(parsedFileDataColumnMapping.value[ImportTransactionColumnType.TransactionType.type])) {
        return {};
    }

    const result: Record<string, TransactionType> = {};

    if (!parsedFileTransactionTypeMapping.value) {
        return result;
    }

    for (const name in parsedFileTransactionTypeMapping.value) {
        if (!Object.prototype.hasOwnProperty.call(parsedFileTransactionTypeMapping.value, name)) {
            continue;
        }

        const type = parsedFileTransactionTypeMapping.value[name];

        if (TransactionType.ModifyBalance <= type && type <= TransactionType.Transfer) {
            result[name] = type;
        }
    }

    return result;
});

const parsedFileAutoDetectedTimeFormat = computed<string | undefined>(() => {
    if (!parsedFileData.value || !parsedFileData.value.length || !isNumber(parsedFileDataColumnMapping.value[ImportTransactionColumnType.TransactionTime.type])) {
        return undefined;
    }

    const allDateTimes: string[] = [];
    const dateTimeColumnIndex = parsedFileDataColumnMapping.value[ImportTransactionColumnType.TransactionTime.type];

    const startIndex = parsedFileIncludeHeader.value ? 1 : 0;

    for (let i = startIndex; i < parsedFileData.value.length; i++) {
        if (parsedFileData.value[i].length <= dateTimeColumnIndex) {
            continue;
        }

        const dateTime = parsedFileData.value[i][dateTimeColumnIndex];

        if (dateTime) {
            allDateTimes.push(dateTime);
        }
    }

    const detectedFormats = KnownDateTimeFormat.detectMulti(allDateTimes);

    if (!detectedFormats || !detectedFormats.length || detectedFormats.length > 1) {
        return undefined;
    }

    return detectedFormats[0].format;
});

const parsedFileAutoDetectedTimezoneFormat = computed<string | undefined>(() => {
    if (!parsedFileData.value || !parsedFileData.value.length || !isNumber(parsedFileDataColumnMapping.value[ImportTransactionColumnType.TransactionTimezone.type])) {
        return undefined;
    }

    const allTimezones: string[] = [];
    const timezoneColumnIndex = parsedFileDataColumnMapping.value[ImportTransactionColumnType.TransactionTimezone.type];

    const startIndex = parsedFileIncludeHeader.value ? 1 : 0;

    for (let i = startIndex; i < parsedFileData.value.length; i++) {
        if (parsedFileData.value[i].length <= timezoneColumnIndex) {
            continue;
        }

        const timezone = parsedFileData.value[i][timezoneColumnIndex];

        if (timezone) {
            allTimezones.push(timezone);
        }
    }

    const detectedFormats = KnownDateTimezoneFormat.detectMulti(allTimezones);

    if (!detectedFormats || !detectedFormats.length || detectedFormats.length > 1) {
        return undefined;
    }

    return detectedFormats[0].value;
});

const parsedFileAutoDetectedAmountFormat = computed<string | undefined>(() => {
    if (!parsedFileData.value || !parsedFileData.value.length || !isNumber(parsedFileDataColumnMapping.value[ImportTransactionColumnType.Amount.type])) {
        return undefined;
    }

    const allAmounts: string[] = [];
    const amountColumnIndex = parsedFileDataColumnMapping.value[ImportTransactionColumnType.Amount.type];

    const startIndex = parsedFileIncludeHeader.value ? 1 : 0;

    for (let i = startIndex; i < parsedFileData.value.length; i++) {
        if (parsedFileData.value[i].length <= amountColumnIndex) {
            continue;
        }

        const amount = parsedFileData.value[i][amountColumnIndex];

        if (amount) {
            allAmounts.push(amount);
        }
    }

    const detectedFormats = KnownAmountFormat.detectMulti(allAmounts);

    if (!detectedFormats || !detectedFormats.length) {
        return undefined;
    }

    return detectedFormats[0].type;
});

const importTransactionsTableHeight = computed<number | undefined>(() => {
    if (countPerPage.value <= 10 || !importTransactions.value || importTransactions.value.length <= 10) {
        return undefined;
    } else {
        return 400;
    }
});

const importTransactionHeaders = computed<object[]>(() => {
    return [
        { value: 'valid', sortable: true, nowrap: true, width: 35 },
        { value: 'time', title: tt('Transaction Time'), sortable: true, nowrap: true, maxWidth: 280 },
        { value: 'type', title: tt('Type'), sortable: true, nowrap: true, maxWidth: 140 },
        { value: 'actualCategoryName', title: tt('Category'), sortable: true, nowrap: true },
        { value: 'sourceAmount', title: tt('Amount'), sortable: true, nowrap: true },
        { value: 'actualSourceAccountName', title: tt('Account'), sortable: true, nowrap: true },
        { value: 'geoLocation', title: tt('Geographic Location'), sortable: true, nowrap: true },
        { value: 'tagIds', title: tt('Tags'), sortable: true, nowrap: true },
        { value: 'comment', title: tt('Description'), sortable: true, nowrap: true },
    ];
});

const importTransactionsTablePageOptions = computed<ImportTransactionsDialogTablePageOption[]>(() => getTablePageOptions(importTransactions.value?.length));

const totalPageCount = computed<number>(() => {
    if (!importTransactions.value || importTransactions.value.length < 1) {
        return 1;
    }

    let count = 0;

    for (let i = 0; i < importTransactions.value.length; i++) {
        if (isTransactionDisplayed(importTransactions.value[i])) {
            count++;
        }
    }

    return Math.ceil(count / countPerPage.value);
});

const currentPageTransactions = computed<ImportTransaction[]>(() => {
    const ret: ImportTransaction[] = [];

    if (!importTransactions.value || importTransactions.value.length < 1) {
        return ret;
    }

    const previousCount = Math.max(0, (currentPage.value - 1) * countPerPage.value);
    let count = 0;

    for (let i = 0; i < importTransactions.value.length; i++) {
        if (ret.length >= countPerPage.value) {
            break;
        }

        if (isTransactionDisplayed(importTransactions.value[i])) {
            if (count >= previousCount) {
                ret.push(importTransactions.value[i]);
            }

            count++;
        }
    }

    return ret;
});

const selectedImportTransactionCount = computed<number>(() => {
    let count = 0;

    if (!importTransactions.value || importTransactions.value.length < 1) {
        return count;
    }

    for (let i = 0; i < importTransactions.value.length; i++) {
        if (importTransactions.value[i].selected) {
            count++;
        }
    }

    return count;
});

const selectedExpenseTransactionCount = computed<number>(() => {
    let count = 0;

    if (!importTransactions.value || importTransactions.value.length < 1) {
        return count;
    }

    for (let i = 0; i < importTransactions.value.length; i++) {
        if (importTransactions.value[i].selected && importTransactions.value[i].type === TransactionType.Expense) {
            count++;
        }
    }

    return count;
});

const selectedIncomeTransactionCount = computed<number>(() => {
    let count = 0;

    if (!importTransactions.value || importTransactions.value.length < 1) {
        return count;
    }

    for (let i = 0; i < importTransactions.value.length; i++) {
        if (importTransactions.value[i].selected && importTransactions.value[i].type === TransactionType.Income) {
            count++;
        }
    }

    return count;
});

const selectedTransferTransactionCount = computed<number>(() => {
    let count = 0;

    if (!importTransactions.value || importTransactions.value.length < 1) {
        return count;
    }

    for (let i = 0; i < importTransactions.value.length; i++) {
        if (importTransactions.value[i].selected && importTransactions.value[i].type === TransactionType.Transfer) {
            count++;
        }
    }

    return count;
});

const selectedInvalidTransactionCount = computed<number>(() => {
    let count = 0;

    if (!importTransactions.value || importTransactions.value.length < 1) {
        return count;
    }

    for (let i = 0; i < importTransactions.value.length; i++) {
        if (!importTransactions.value[i].valid && importTransactions.value[i].selected) {
            count++;
        }
    }

    return count;
});

const anyButNotAllTransactionSelected = computed<boolean>(() => !!importTransactions.value && selectedImportTransactionCount.value > 0 && selectedImportTransactionCount.value !== importTransactions.value.length);
const allTransactionSelected = computed<boolean>(() => !!importTransactions.value && selectedImportTransactionCount.value === importTransactions.value.length);

const allUsedCategoryNames = computed<string[]>(() => {
    if (!importTransactions.value || importTransactions.value.length < 1) {
        return [];
    }

    const categoryNames: Record<string, boolean> = {};

    for (const transaction of importTransactions.value) {
        if (transaction.actualCategoryName && transaction.actualCategoryName !== '') {
            categoryNames[transaction.actualCategoryName] = true;
        }
    }

    return objectFieldToArrayItem(categoryNames);
});

const allUsedAccountNames = computed<string[]>(() => {
    if (!importTransactions.value || importTransactions.value.length < 1) {
        return [];
    }

    const accountNames: Record<string, boolean> = {};

    for (const transaction of importTransactions.value) {
        if (transaction.actualSourceAccountName && transaction.actualSourceAccountName !== '') {
            accountNames[transaction.actualSourceAccountName] = true;
        }

        if (transaction.actualDestinationAccountName && transaction.actualDestinationAccountName !== '') {
            accountNames[transaction.actualDestinationAccountName] = true;
        }
    }

    return objectFieldToArrayItem(accountNames);
});

const allUsedTagNames = computed<string[]>(() => {
    if (!importTransactions.value || importTransactions.value.length < 1) {
        return [];
    }

    const tagNames: Record<string, boolean> = {};

    for (const transaction of importTransactions.value) {
        if (!transaction.tagIds || !transaction.originalTagNames) {
            continue;
        }

        for (let j = 0; j < transaction.tagIds.length; j++) {
            const tagId = transaction.tagIds[j];
            const originalTagName = transaction.originalTagNames[j];

            if (tagId && tagId !== '0' && allTagsMap.value[tagId] && allTagsMap.value[tagId].name) {
                tagNames[allTagsMap.value[tagId].name] = true;
            } else if (originalTagName) {
                tagNames[originalTagName] = true;
            }
        }
    }

    return objectFieldToArrayItem(tagNames);
});

const allInvalidExpenseCategoryNames = computed<NameValue[]>(() => getCurrentInvalidCategoryNames(TransactionType.Expense));
const allInvalidIncomeCategoryNames = computed<NameValue[]>(() => getCurrentInvalidCategoryNames(TransactionType.Income));
const allInvalidTransferCategoryNames = computed<NameValue[]>(() => getCurrentInvalidCategoryNames(TransactionType.Transfer));
const allInvalidAccountNames = computed<NameValue[]>(() => getCurrentInvalidAccountNames());
const allInvalidTransactionTagNames = computed<NameValue[]>(() => getCurrentInvalidTagNames());

const displayFilterCustomDateRange = computed<string>(() => {
    if (filters.value.minDatetime === null || filters.value.maxDatetime === null) {
        return '';
    }

    const minDisplayTime = formatUnixTimeToLongDateTime(filters.value.minDatetime);
    const maxDisplayTime = formatUnixTimeToLongDateTime(filters.value.maxDatetime);

    return `${minDisplayTime} - ${maxDisplayTime}`
});

function getTablePageOptions(linesCount?: number): ImportTransactionsDialogTablePageOption[] {
    const pageOptions: ImportTransactionsDialogTablePageOption[] = [];

    if (!linesCount || linesCount < 1) {
        pageOptions.push({ value: -1, title: tt('All') });
        return pageOptions;
    }

    const availableCountPerPage = [ 5, 10, 15, 20, 25, 30, 50 ];

    for (let i = 0; i < availableCountPerPage.length; i++) {
        const count = availableCountPerPage[i];

        if (linesCount < count) {
            break;
        }

        pageOptions.push({ value: count, title: count.toString() });
    }

    pageOptions.push({ value: -1, title: tt('All') });

    return pageOptions;
}

function getParseDataMappedColumnDisplayName(columnIndex: number): string {
    for (const columnType in parsedFileDataColumnMapping.value) {
        if (parsedFileDataColumnMapping.value[columnType] === columnIndex) {
            return findDisplayNameByType(allImportTransactionColumnTypes.value, parseInt(columnType)) || tt('Unspecified');
        }
    }

    return tt('Unspecified');
}

function updateParseDataMappedColumn(columnIndex: number, columnType: number): void {
    if (parsedFileDataColumnMapping.value[columnType] === columnIndex) {
        delete parsedFileDataColumnMapping.value[columnType];
    } else {
        parsedFileDataColumnMapping.value[columnType] = columnIndex;
    }

    for (const otherColumnType in parsedFileDataColumnMapping.value) {
        if (otherColumnType !== columnType.toString() && parsedFileDataColumnMapping.value[otherColumnType] === columnIndex) {
            delete parsedFileDataColumnMapping.value[otherColumnType];
        }
    }
}

function isTransactionDisplayed(transaction: ImportTransaction): boolean {
    if (isNumber(filters.value.minDatetime) && isNumber(filters.value.maxDatetime) && (transaction.time < filters.value.minDatetime || transaction.time > filters.value.maxDatetime)) {
        return false;
    }

    if (isNumber(filters.value.transactionType) && transaction.type !== filters.value.transactionType) {
        return false;
    }

    if (isString(filters.value.category)) {
        if (filters.value.category === '' && transaction.actualCategoryName !== '') {
            return false;
        } else if (filters.value.category !== '' && transaction.actualCategoryName !== filters.value.category) {
            return false;
        }
    } else if (filters.value.category === undefined) {
        if (transaction.type !== TransactionType.ModifyBalance && transaction.categoryId && transaction.categoryId !== '0') {
            return false;
        }
    }

    if (isString(filters.value.account)) {
        if (filters.value.account === '' && (transaction.actualSourceAccountName !== '' || transaction.actualDestinationAccountName !== '')) {
            return false;
        } else if (filters.value.account !== '' && transaction.actualSourceAccountName !== filters.value.account && transaction.actualDestinationAccountName !== filters.value.account) {
            return false;
        }
    } else if (filters.value.account === undefined) {
        if (transaction.type !== TransactionType.Transfer && transaction.sourceAccountId && transaction.sourceAccountId !== '0') {
            return false;
        } else if (transaction.type === TransactionType.Transfer && transaction.sourceAccountId && transaction.sourceAccountId !== '0' && transaction.destinationAccountId && transaction.destinationAccountId !== '0') {
            return false;
        }
    }

    if (isString(filters.value.tag)) {
        if (filters.value.tag === '' && transaction.tagIds && transaction.tagIds.length) {
            return false;
        } else if (filters.value.tag !== '') {
            let hasTagName = false;

            if (transaction.tagIds && transaction.tagIds.length) {
                for (let i = 0; i < transaction.tagIds.length; i++) {
                    const tagId = transaction.tagIds[i];
                    let tagName = transaction.originalTagNames ? transaction.originalTagNames[i] : "";

                    if (tagId && tagId !== '0' && allTagsMap.value[tagId] && allTagsMap.value[tagId].name) {
                        tagName = allTagsMap.value[tagId].name;
                    }

                    if (tagName === filters.value.tag) {
                        hasTagName = true;
                        break;
                    }
                }
            }

            if (!hasTagName) {
                return false;
            }
        }
    } else if (filters.value.tag === undefined) {
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

    if (isString(filters.value.description)) {
        if (filters.value.description === '' && transaction.comment !== '') {
            return false;
        } else if (filters.value.description !== '' && transaction.comment.indexOf(filters.value.description) < 0) {
            return false;
        }
    }

    return true;
}

function isTagValid(tagIds: string[], tagIndex: number): boolean {
    if (!tagIds || !tagIds[tagIndex]) {
        return false;
    }

    if (tagIds[tagIndex] === '0') {
        return false;
    }

    const tagId = tagIds[tagIndex];
    return !!allTagsMap.value[tagId];
}

function getDisplayDateTime(transaction: ImportTransaction): string {
    const transactionTime = getUnixTime(parseDateFromUnixTime(transaction.time, transaction.utcOffset, currentTimezoneOffsetMinutes.value));
    return formatUnixTimeToLongDateTime(transactionTime);
}

function getDisplayTimezone(transaction: ImportTransaction): string {
    return `UTC${getUtcOffsetByUtcOffsetMinutes(transaction.utcOffset)}`;
}

function getDisplayCurrency(value: number, currencyCode: string): string {
    return formatAmountWithCurrency(value, currencyCode);
}

function getTransactionDisplayAmount(transaction: ImportTransaction): string {
    let currency = transaction.originalSourceAccountCurrency || defaultCurrency.value;

    if (transaction.sourceAccountId && transaction.sourceAccountId !== '0' && allAccountsMap.value[transaction.sourceAccountId]) {
        currency = allAccountsMap.value[transaction.sourceAccountId].currency;
    }

    return getDisplayCurrency(transaction.sourceAmount, currency);
}

function getTransactionDisplayDestinationAmount(transaction: ImportTransaction): string {
    if (transaction.type !== TransactionType.Transfer) {
        return '-';
    }

    let currency = transaction.originalDestinationAccountCurrency || defaultCurrency.value;

    if (transaction.destinationAccountId && transaction.destinationAccountId !== '0' && allAccountsMap.value[transaction.destinationAccountId]) {
        currency = allAccountsMap.value[transaction.destinationAccountId].currency;
    }

    return getDisplayCurrency(transaction.destinationAmount, currency);
}

function getSourceAccountTitle(transaction: ImportTransaction): string {
    if (transaction.type === TransactionType.Expense || transaction.type === TransactionType.Income) {
        return tt('Account');
    } else if (transaction.type === TransactionType.Transfer) {
        return tt('Source Account');
    } else {
        return tt('Account');
    }
}

function getSourceAccountDisplayName(transaction: ImportTransaction): string {
    if (transaction.sourceAccountId) {
        return Account.findAccountNameById(allAccounts.value, transaction.sourceAccountId) || '';
    } else {
        return tt('None');
    }
}

function getDestinationAccountDisplayName(transaction: ImportTransaction): string {
    if (transaction.destinationAccountId) {
        return Account.findAccountNameById(allAccounts.value, transaction.destinationAccountId) || '';
    } else {
        return tt('None');
    }
}

function getCurrentInvalidCategoryNames(transactionType: TransactionType): NameValue[] {
    const invalidCategoryNames: Record<string, boolean> = {};
    const invalidCategories: NameValue[] = [];

    if (!importTransactions.value || importTransactions.value.length < 1) {
        return invalidCategories;
    }

    for (let i = 0; i < importTransactions.value.length; i++) {
        const transaction: ImportTransaction = importTransactions.value[i];
        const categoryId = transaction.categoryId;

        if (transaction.type === transactionType && (!categoryId || categoryId === '0' || !allCategoriesMap.value[categoryId])) {
            invalidCategoryNames[transaction.originalCategoryName] = true;
        }
    }

    for (const name in invalidCategoryNames) {
        if (!Object.prototype.hasOwnProperty.call(invalidCategoryNames, name)) {
            continue;
        }

        invalidCategories.push({
            name: name || tt('(Empty)'),
            value: name
        });
    }

    return invalidCategories;
}

function getCurrentInvalidAccountNames(): NameValue[] {
    const invalidAccountNames: Record<string, boolean> = {};
    const invalidAccounts: NameValue[] = [];

    if (!importTransactions.value || importTransactions.value.length < 1) {
        return invalidAccounts;
    }

    for (let i = 0; i < importTransactions.value.length; i++) {
        const transaction: ImportTransaction = importTransactions.value[i];
        const sourceAccountId = transaction.sourceAccountId;
        const destinationAccountId = transaction.destinationAccountId;

        if (!sourceAccountId || sourceAccountId === '0' || !allAccountsMap.value[sourceAccountId]) {
            invalidAccountNames[transaction.originalSourceAccountName] = true;
        }

        if (transaction.type === TransactionType.Transfer && isString(transaction.originalDestinationAccountName) && (!destinationAccountId || destinationAccountId === '0' || !allAccountsMap.value[destinationAccountId])) {
            invalidAccountNames[transaction.originalDestinationAccountName] = true;
        }
    }

    for (const name in invalidAccountNames) {
        if (!Object.prototype.hasOwnProperty.call(invalidAccountNames, name)) {
            continue;
        }

        invalidAccounts.push({
            name: name || tt('(Empty)'),
            value: name
        });
    }

    return invalidAccounts;
}

function getCurrentInvalidTagNames(): NameValue[] {
    const invalidTagNames: Record<string, boolean> = {};
    const invalidTags: NameValue[] = [];

    if (!importTransactions.value || importTransactions.value.length < 1) {
        return invalidTags;
    }

    for (let i = 0; i < importTransactions.value.length; i++) {
        const transaction: ImportTransaction = importTransactions.value[i];

        if (!transaction.tagIds || !transaction.originalTagNames) {
            continue;
        }

        for (let j = 0; j < transaction.tagIds.length; j++) {
            const tagId = transaction.tagIds[j];
            const originalTagName = transaction.originalTagNames[j];

            if (!tagId || tagId === '0' || !allTagsMap.value[tagId]) {
                invalidTagNames[originalTagName] = true;
            }
        }
    }

    for (const name in invalidTagNames) {
        if (!Object.prototype.hasOwnProperty.call(invalidTagNames, name)) {
            continue;
        }

        invalidTags.push({
            name: name || tt('(Empty)'),
            value: name
        });
    }

    return invalidTags;
}

function open(): Promise<void> {
    fileType.value = 'ezbookkeeping';
    fileSubType.value = 'ezbookkeeping_csv';
    fileEncoding.value = 'utf-8';
    currentStep.value = 'uploadFile';
    importProcess.value = 0;
    importFile.value = null;
    importData.value = '';
    parsedFileData.value = undefined;
    parsedFileIncludeHeader.value = true;
    parsedFileDataColumnMapping.value = {};
    parsedFileTransactionTypeMapping.value = {};
    parsedFileTimeFormat.value = '';
    parsedFileTimezoneFormat.value = '';
    parsedFileAmountFormat.value = '';
    parsedFileGeoLocationSeparator.value = ' ';
    parsedFileTagSeparator.value = ';';
    importTransactions.value = undefined;
    editingTransaction.value = null;
    editingTags.value = [];
    filters.value.minDatetime = null;
    filters.value.maxDatetime = null;
    filters.value.transactionType = null;
    filters.value.category = null;
    filters.value.account = null;
    filters.value.tag = null;
    filters.value.description = null;
    currentPage.value = 1;
    countPerPage.value = 10;
    showState.value = true;
    clientSessionId.value = generateRandomUUID();

    const promises = [
        accountsStore.loadAllAccounts({ force: false }),
        transactionCategoriesStore.loadAllCategories({ force: false }),
        transactionTagsStore.loadAllTags({ force: false })
    ];

    Promise.all(promises).then(() => {
        loading.value = false;
    }).catch(error => {
        logger.error('failed to load essential data for importing transaction', error);

        loading.value = false;
        showState.value = false;

        if (!error.processed) {
            if (rejectFunc) {
                rejectFunc(error);
            }
        }
    });

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function showOpenFileDialog(): void {
    if (submitting.value) {
        return;
    }

    fileInput.value?.click();
}

function setImportFile(event: Event): void {
    if (!event || !event.target) {
        return;
    }

    const el = event.target as HTMLInputElement;

    if (!el.files || !el.files.length) {
        return;
    }

    importFile.value = el.files[0];
    el.value = '';
}

function parseData(): void {
    let uploadFile: File;
    let type: string = fileType.value;
    let encoding: string | undefined = undefined;

    if (allFileSubTypes.value) {
        type = fileSubType.value;
    }

    if (allSupportedEncodings.value) {
        encoding = fileEncoding.value;
    }

    if (!isImportDataFromTextbox.value) {
        if (!importFile.value) {
            snackbar.value?.showError('Please select a file to import');
            return;
        }

        uploadFile = importFile.value;
    } else if (isImportDataFromTextbox.value) {
        if (!importData.value) {
            snackbar.value?.showError('No data to import');
            return;
        }

        if (type === 'custom_csv') {
            uploadFile = new File([importData.value], 'import.csv', { type: 'text/csv' });
        } else if (type === 'custom_tsv') {
            uploadFile = new File([importData.value], 'import.tsv', { type: 'text/tab-separated-values' });
        } else {
            snackbar.value?.showError('Parameter Invalid');
            return;
        }

        encoding = 'utf-8';
    } else { // should not happen, but ts would check whether uploadFile has been assigned a value
        snackbar.value?.showMessage('An error occurred');
        return;
    }

    const isDsvFileType: boolean = fileType.value === 'dsv' || fileType.value === 'dsv_data';

    if (isDsvFileType && currentStep.value === 'uploadFile') {
        submitting.value = true;

        transactionsStore.parseImportDsvFile({
            fileType: type,
            fileEncoding: encoding,
            importFile: uploadFile
        }).then(response => {
            if (response && response.length) {
                parsedFileData.value = response;
                currentPage.value = 1;
                countPerPage.value = 10;
                currentStep.value = 'defineColumn';
            } else {
                parsedFileData.value = undefined;
                snackbar.value?.showError('No data to import');
            }

            submitting.value = false;
        }).catch(error => {
            submitting.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    } else {
        let columnMapping: Record<number, number> | undefined = undefined;
        let transactionTypeMapping: Record<string, TransactionType> | undefined = undefined;
        let hasHeaderLine: boolean | undefined = undefined;
        let timeFormat: string | undefined = undefined;
        let timezoneFormat: string | undefined = undefined;
        let amountFormat: string | undefined = undefined;
        let amountDecimalSeparator: string | undefined = undefined;
        let amountDigitGroupingSymbol: string | undefined = undefined;
        let geoLocationSeparator: string | undefined = undefined;
        let tagSeparator: string | undefined = undefined;

        if (isDsvFileType) {
            columnMapping = parsedFileDataColumnMapping.value;
            transactionTypeMapping = parsedFileValidMappedTransactionTypes.value;
            hasHeaderLine = parsedFileIncludeHeader.value;
            timeFormat = parsedFileTimeFormat.value;
            timezoneFormat = parsedFileTimezoneFormat.value;
            amountFormat = parsedFileAmountFormat.value;
            geoLocationSeparator = parsedFileGeoLocationSeparator.value;
            tagSeparator = parsedFileTagSeparator.value;

            if (!columnMapping
                || !isNumber(columnMapping[ImportTransactionColumnType.TransactionTime.type])
                || !isNumber(columnMapping[ImportTransactionColumnType.TransactionType.type])
                || !isNumber(columnMapping[ImportTransactionColumnType.Amount.type])) {
                snackbar.value?.showError('Missing transaction time, transaction type, or amount column mapping');
                return;
            }

            if (!transactionTypeMapping || isObjectEmpty(transactionTypeMapping)) {
                snackbar.value?.showError('Transaction type mapping is not set');
                return;
            }

            if (!parsedFileTimeFormat.value) {
                timeFormat = parsedFileAutoDetectedTimeFormat.value;
            }

            if (!parsedFileTimezoneFormat.value) {
                timezoneFormat = parsedFileAutoDetectedTimezoneFormat.value;
            }

            if (!parsedFileAmountFormat.value) {
                amountFormat = parsedFileAutoDetectedAmountFormat.value;
            }

            if (amountFormat) {
                const knownAmountFormat = KnownAmountFormat.valueOf(amountFormat);

                if (knownAmountFormat) {
                    amountDecimalSeparator = knownAmountFormat.decimalSeparator.symbol;
                    amountDigitGroupingSymbol = knownAmountFormat.digitGroupingSymbol?.symbol;
                }
            }

            if (!timeFormat) {
                snackbar.value?.showError('Transaction time format is not set');
                return;
            }

            if (!amountDecimalSeparator) {
                snackbar.value?.showError('Transaction amount format is not set');
                return;
            }
        }

        submitting.value = true;

        transactionsStore.parseImportTransaction({
            fileType: type,
            fileEncoding: encoding,
            importFile: uploadFile,
            columnMapping: columnMapping,
            transactionTypeMapping: transactionTypeMapping,
            hasHeaderLine: hasHeaderLine,
            timeFormat: timeFormat,
            timezoneFormat: timezoneFormat,
            amountDecimalSeparator: amountDecimalSeparator,
            amountDigitGroupingSymbol: amountDigitGroupingSymbol,
            geoSeparator: geoLocationSeparator,
            tagSeparator: tagSeparator
        }).then(response => {
            const parsedTransactions: ImportTransaction[] = [];

            if (response.items) {
                for (let i = 0; i < response.items.length; i++) {
                    const parsedTransaction = ImportTransaction.of(response.items[i], i);
                    parsedTransactions.push(parsedTransaction);
                }
            }

            importTransactions.value = parsedTransactions;
            editingTransaction.value = null;
            editingTags.value = [];
            currentPage.value = 1;

            if (importTransactions.value && importTransactions.value.length >= 0 && importTransactions.value.length < 10) {
                countPerPage.value = -1;
            } else {
                countPerPage.value = 10;
            }

            currentPage.value = 1;
            countPerPage.value = 10;
            currentStep.value = 'checkData';
            submitting.value = false;
        }).catch(error => {
            submitting.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    }
}

function submit(): void {
    if (editingTransaction.value) {
        return;
    }

    const transactions: ImportTransaction[] = [];

    if (importTransactions.value) {
        for (let i = 0; i < importTransactions.value.length; i++) {
            const transaction: ImportTransaction = importTransactions.value[i];

            if (transaction.valid && transaction.selected) {
                transactions.push(transaction);
            } else if (!transaction.valid && transaction.selected) {
                snackbar.value?.showError('Cannot import invalid transactions');
                return;
            }
        }
    }

    if (transactions.length < 1) {
        snackbar.value?.showError('No data to import');
        return;
    }

    confirmDialog.value?.open('format.misc.confirmImportTransactions', {
        count: transactions.length
    }).then(() => {
        editingTransaction.value = null;
        editingTags.value = [];
        submitting.value = true;

        let showProcessTimer : number | undefined = undefined;

        if (transactions.length > 100) {
            setTimeout(() => {
                if (!submitting.value) {
                    logger.warn('transaction import is not submitting');
                    return;
                }

                // @ts-expect-error the return value of setInterval is number, but lint shows it as NodeJS.Timer
                showProcessTimer = setInterval(() => {
                    if (submitting.value) {
                        transactionsStore.getImportTransactionsProcess({
                            clientSessionId: clientSessionId.value
                        }).then(response => {
                            if (isNumber(response) && 0 <= response && response < 100) {
                                importProcess.value = response;
                            } else {
                                importProcess.value = 0;
                                clearInterval(showProcessTimer);
                                showProcessTimer = undefined;
                            }
                        }).catch(() => {
                            importProcess.value = 0;
                            clearInterval(showProcessTimer);
                            showProcessTimer = undefined;
                        });
                    }
                }, 2000);
            }, 2000);
        }

        transactionsStore.importTransactions({
            transactions: transactions,
            clientSessionId: clientSessionId.value
        }).then(response => {
            if (showProcessTimer) {
                importProcess.value = 0;
                clearInterval(showProcessTimer);
                showProcessTimer = undefined;
            }

            importedCount.value = response;
            currentStep.value = 'finalResult';

            accountsStore.updateAccountListInvalidState(true);
            transactionsStore.updateTransactionListInvalidState(true);
            overviewStore.updateTransactionOverviewInvalidState(true);
            statisticsStore.updateTransactionStatisticsInvalidState(true);

            submitting.value = false;
        }).catch(error => {
            if (showProcessTimer) {
                importProcess.value = 0;
                clearInterval(showProcessTimer);
                showProcessTimer = undefined;
            }

            submitting.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    });
}

function close(completed: boolean): void {
    if (completed) {
        if (resolveFunc) {
            resolveFunc();
        }
    } else {
        if (rejectFunc) {
            rejectFunc();
        }
    }

    showState.value = false;
}

function changeCustomDateFilter(minTime: number, maxTime: number): void {
    filters.value.minDatetime = minTime;
    filters.value.maxDatetime = maxTime;
    showCustomDateRangeDialog.value = false;
}

function importTransactionsFilter(value: string, query: string, item?: { value: unknown, raw: ImportTransaction }): boolean {
    if (!item || !item.raw) {
        return false;
    }

    return isTransactionDisplayed(item.raw);
}

function selectAllValid(): void {
    if (!importTransactions.value || importTransactions.value.length < 1) {
        return;
    }

    for (let i = 0; i < importTransactions.value.length; i++) {
        if (importTransactions.value[i].valid && isTransactionDisplayed(importTransactions.value[i])) {
            importTransactions.value[i].selected = true;
        }
    }
}

function selectAllInvalid(): void {
    if (!importTransactions.value || importTransactions.value.length < 1) {
        return;
    }

    for (let i = 0; i < importTransactions.value.length; i++) {
        if (!importTransactions.value[i].valid && isTransactionDisplayed(importTransactions.value[i])) {
            importTransactions.value[i].selected = true;
        }
    }
}

function selectAll(): void {
    if (!importTransactions.value || importTransactions.value.length < 1) {
        return;
    }

    for (let i = 0; i < importTransactions.value.length; i++) {
        if (isTransactionDisplayed(importTransactions.value[i])) {
            importTransactions.value[i].selected = true;
        }
    }
}

function selectNone(): void {
    if (!importTransactions.value || importTransactions.value.length < 1) {
        return;
    }

    for (let i = 0; i < importTransactions.value.length; i++) {
        if (isTransactionDisplayed(importTransactions.value[i])) {
            importTransactions.value[i].selected = false;
        }
    }
}

function selectInvert(): void {
    if (!importTransactions.value || importTransactions.value.length < 1) {
        return;
    }

    for (let i = 0; i < importTransactions.value.length; i++) {
        if (isTransactionDisplayed(importTransactions.value[i])) {
            importTransactions.value[i].selected = !importTransactions.value[i].selected;
        }
    }
}

function selectAllInThisPage(): void {
    for (let i = 0; i < currentPageTransactions.value.length; i++) {
        currentPageTransactions.value[i].selected = true;
    }
}

function selectNoneInThisPage(): void {
    for (let i = 0; i < currentPageTransactions.value.length; i++) {
        currentPageTransactions.value[i].selected = false;
    }
}

function selectInvertInThisPage(): void {
    for (let i = 0; i < currentPageTransactions.value.length; i++) {
        currentPageTransactions.value[i].selected = !currentPageTransactions.value[i].selected;
    }
}

function editTransaction(transaction: ImportTransaction): void {
    if (editingTransaction.value) {
        editingTransaction.value.tagIds = editingTags.value;
        updateTransactionData(editingTransaction.value);
    }

    if (editingTransaction.value === transaction) {
        editingTags.value = [];
        editingTransaction.value = null;
    } else {
        editingTransaction.value = transaction;
        editingTags.value = editingTransaction.value.tagIds;
    }
}

function updateTransactionData(transaction: ImportTransaction): void {
    transaction.valid = transaction.isTransactionValid();

    if (transaction.categoryId && allCategoriesMap.value[transaction.categoryId]) {
        transaction.actualCategoryName = allCategoriesMap.value[transaction.categoryId].name;
    }

    if (transaction.sourceAccountId && allAccountsMap.value[transaction.sourceAccountId]) {
        transaction.actualSourceAccountName = allAccountsMap.value[transaction.sourceAccountId].name;
    }

    if (transaction.destinationAccountId && allAccountsMap.value[transaction.destinationAccountId]) {
        transaction.actualDestinationAccountName = allAccountsMap.value[transaction.destinationAccountId].name;
    }
}

function showBatchReplaceDialog(type: BatchReplaceDialogDataType): void {
    if (editingTransaction.value) {
        return;
    }

    batchReplaceDialog.value?.open({
        mode: 'batchReplace',
        type: type
    }).then(result => {
        if (!result || !result.targetItem) {
            return;
        }

        let updatedCount = 0;

        if (importTransactions.value) {
            for (let i = 0; i < importTransactions.value.length; i++) {
                const transaction: ImportTransaction = importTransactions.value[i];

                if (!transaction.selected) {
                    continue;
                }

                let updated = false;

                if (type === 'expenseCategory') {
                    if (transaction.type === TransactionType.Expense) {
                        transaction.categoryId = result.targetItem;
                        updated = true;
                    }
                } else if (type === 'incomeCategory') {
                    if (transaction.type === TransactionType.Income) {
                        transaction.categoryId = result.targetItem;
                        updated = true;
                    }
                } else if (type === 'transferCategory') {
                    if (transaction.type === TransactionType.Transfer) {
                        transaction.categoryId = result.targetItem;
                        updated = true;
                    }
                } else if (type === 'account') {
                    transaction.sourceAccountId = result.targetItem;
                    updated = true;
                } else if (type === 'destinationAccount') {
                    if (transaction.type === TransactionType.Transfer) {
                        transaction.destinationAccountId = result.targetItem;
                        updated = true;
                    }
                }

                if (updated) {
                    updatedCount++;
                    updateTransactionData(transaction);
                }
            }
        }

        if (updatedCount > 0) {
            snackbar.value?.showMessage('format.misc.youHaveUpdatedTransactions', {
                count: updatedCount
            });
        }
    });
}

function showReplaceInvalidItemDialog(type: BatchReplaceDialogDataType, invalidItems: NameValue[]): void {
    if (editingTransaction.value) {
        return;
    }

    batchReplaceDialog.value?.open({
        mode: 'replaceInvalidItems',
        type: type,
        invalidItems: invalidItems
    }).then(result => {
        if (!result || (!result.sourceItem && result.sourceItem !== '') || !result.targetItem) {
            return;
        }

        let updatedCount = 0;

        if (importTransactions.value) {
            for (let i = 0; i < importTransactions.value.length; i++) {
                const transaction: ImportTransaction = importTransactions.value[i];

                if (transaction.valid) {
                    continue;
                }

                let updated = false;

                if (type === 'expenseCategory' || type === 'incomeCategory' || type === 'transferCategory') {
                    const categoryId = transaction.categoryId;
                    const originalCategoryName = transaction.originalCategoryName;

                    if (transaction.type !== TransactionType.ModifyBalance && originalCategoryName === result.sourceItem && (!categoryId || categoryId === '0' || !allCategoriesMap.value[categoryId])) {
                        if (type === 'expenseCategory' && transaction.type === TransactionType.Expense) {
                            transaction.categoryId = result.targetItem;
                            updated = true;
                        } else if (type === 'incomeCategory' && transaction.type === TransactionType.Income) {
                            transaction.categoryId = result.targetItem;
                            updated = true;
                        } else if (type === 'transferCategory' && transaction.type === TransactionType.Transfer) {
                            transaction.categoryId = result.targetItem;
                            updated = true;
                        }
                    }
                } else if (type === 'account') {
                    const sourceAccountId = transaction.sourceAccountId;
                    const originalSourceAccountName = transaction.originalSourceAccountName;
                    const destinationAccountId = transaction.destinationAccountId;
                    const originalDestinationAccountName = transaction.originalDestinationAccountName;

                    if (originalSourceAccountName === result.sourceItem && (!sourceAccountId || sourceAccountId === '0' || !allAccountsMap.value[sourceAccountId])) {
                        transaction.sourceAccountId = result.targetItem;
                        updated = true;
                    }

                    if (transaction.type === TransactionType.Transfer && originalDestinationAccountName === result.sourceItem && (!destinationAccountId || destinationAccountId === '0' || !allAccountsMap.value[destinationAccountId])) {
                        transaction.destinationAccountId = result.targetItem;
                        updated = true;
                    }
                } else if (type === 'tag' && transaction.tagIds) {
                    for (let j = 0; j < transaction.tagIds.length; j++) {
                        const tagId = transaction.tagIds[j];
                        const originalTagName = transaction.originalTagNames ? transaction.originalTagNames[j] : "";

                        if (originalTagName === result.sourceItem && (!tagId || tagId === '0' || !allTagsMap.value[tagId])) {
                            transaction.tagIds[j] = result.targetItem;
                            updated = true;
                        }
                    }
                }

                if (updated) {
                    updatedCount++;
                    updateTransactionData(transaction);
                }
            }
        }

        if (updatedCount > 0) {
            snackbar.value?.showMessage('format.misc.youHaveUpdatedTransactions', {
                count: updatedCount
            });
        }
    });
}

function showBatchCreateInvalidItemDialog(type: BatchCreateDialogDataType, invalidItems: NameValue[]): void {
    if (editingTransaction.value) {
        return;
    }

    batchCreateDialog.value?.open({
        type: type,
        invalidItems: invalidItems
    }).then(result => {
        if (!result || !result.sourceTargetMap) {
            return;
        }

        let updatedCount = 0;

        if (importTransactions.value) {
            const sourceTargetMap: Record<string, string> = result.sourceTargetMap;

            for (let i = 0; i < importTransactions.value.length; i++) {
                const transaction: ImportTransaction = importTransactions.value[i];

                if (transaction.valid) {
                    continue;
                }

                let updated = false;

                if (type === 'expenseCategory' || type === 'incomeCategory' || type === 'transferCategory') {
                    const categoryId = transaction.categoryId;
                    const originalCategoryName = transaction.originalCategoryName;
                    const targetItem = sourceTargetMap[originalCategoryName];

                    if (transaction.type !== TransactionType.ModifyBalance && targetItem && (!categoryId || categoryId === '0' || !allCategoriesMap.value[categoryId])) {
                        if (type === 'expenseCategory' && transaction.type === TransactionType.Expense) {
                            transaction.categoryId = targetItem;
                            updated = true;
                        } else if (type === 'incomeCategory' && transaction.type === TransactionType.Income) {
                            transaction.categoryId = targetItem;
                            updated = true;
                        } else if (type === 'transferCategory' && transaction.type === TransactionType.Transfer) {
                            transaction.categoryId = targetItem;
                            updated = true;
                        }
                    }
                } else if (type === 'tag' && transaction.tagIds) {
                    for (let j = 0; j < transaction.tagIds.length; j++) {
                        const tagId = transaction.tagIds[j];
                        const originalTagName = transaction.originalTagNames ? transaction.originalTagNames[j] : "";
                        const targetItem = sourceTargetMap[originalTagName];

                        if (targetItem && (!tagId || tagId === '0' || !allTagsMap.value[tagId])) {
                            transaction.tagIds[j] = targetItem;
                            updated = true;
                        }
                    }
                }

                if (updated) {
                    updatedCount++;
                    updateTransactionData(transaction);
                }
            }
        }

        if (updatedCount > 0) {
            snackbar.value?.showMessage('format.misc.youHaveUpdatedTransactions', {
                count: updatedCount
            });
        }
    });
}

function convertTransactionType(fromType: TransactionType, toType: TransactionType): void {
    if (!importTransactions.value || importTransactions.value.length < 1) {
        return;
    }

    const categoryType = transactionTypeToCategoryType(toType);

    if (!categoryType) {
        return;
    }

    const categoryMapByName: Record<string, TransactionCategory> = getSecondaryTransactionMapByName(allCategories.value[categoryType]);

    for (let i = 0; i < importTransactions.value.length; i++) {
        const transaction: ImportTransaction = importTransactions.value[i];

        if (!transaction.selected || transaction.type !== fromType) {
            continue;
        }

        transaction.type = toType;
        transaction.categoryId = categoryMapByName[transaction.originalCategoryName]?.id || '0';

        if (transaction.type === TransactionType.Transfer) {
            transaction.destinationAccountId = allAccountsMapByName.value[transaction.originalDestinationAccountName || '']?.id || '0';
            transaction.destinationAmount = transaction.sourceAmount;
        } else {
            if (fromType === TransactionType.Transfer && toType === TransactionType.Income) {
                transaction.sourceAccountId = transaction.destinationAccountId;
                transaction.sourceAmount = transaction.destinationAmount;
            }

            transaction.destinationAccountId = '0';
            transaction.destinationAmount = 0;
        }

        updateTransactionData(transaction);
    }
}

function onShowDateRangeError(message: string): void {
    snackbar.value?.showError(message);
}

watch(fileType, () => {
    if (allFileSubTypes.value && allFileSubTypes.value.length) {
        fileSubType.value = allFileSubTypes.value[0].type;
    }

    importFile.value = null;
    importTransactions.value = undefined;
    editingTransaction.value = null;
    editingTags.value = [];
    currentPage.value = 1;
    countPerPage.value = 10;
});

watch(fileSubType, (newValue) => {
    let supportedExtensions: string | undefined = findExtensionByType(allFileSubTypes.value, newValue);

    if (!supportedExtensions) {
        supportedExtensions = findExtensionByType(allSupportedImportFileTypes.value, fileType.value);
    }

    if (importFile.value && importFile.value.name && !isFileExtensionSupported(importFile.value.name, supportedExtensions || '')) {
        importFile.value = null;
    }
});

defineExpose({
    open
});
</script>

<style>
.transaction-types-popup-menu .transaction-types-toggle {
    overflow-x: auto;
    white-space: nowrap;
}

.transaction-types-popup-menu .transaction-types-toggle.v-btn-toggle {
    height: auto !important;
    padding: 0;
    border: none;
}

.transaction-types-popup-menu .transaction-types-toggle.v-btn-toggle > .v-btn {
    border-color: rgba(var(--v-border-color), var(--v-border-opacity));
}

.transaction-types-popup-menu .transaction-types-toggle.v-btn-toggle > .v-btn:not(:first-child) {
    border-top-left-radius: 0;
    border-bottom-left-radius: 0;
    border-left: none;
}

.transaction-types-popup-menu .transaction-types-toggle.v-btn-toggle > .v-btn:not(:last-child) {
    border-top-right-radius: 0;
    border-bottom-right-radius: 0;
}

.transaction-types-popup-menu .transaction-types-toggle.v-btn-toggle > .v-btn {
    border: thin solid rgba(var(--v-border-color), var(--v-border-opacity));
}

.transaction-types-popup-menu .transaction-types-toggle.v-btn-toggle button.v-btn {
    width: auto !important;
}

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
