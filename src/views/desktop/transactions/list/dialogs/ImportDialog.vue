<template>
    <v-dialog width="1000" :persistent="!!persistent" v-model="showState">
        <v-card class="pa-6 pa-sm-10 pa-md-12">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <div class="d-flex w-100 align-center justify-center">
                        <h4 class="text-h4">{{ $t('Import Transactions') }}</h4>
                        <v-progress-circular indeterminate size="22" class="ml-2" v-if="loading"></v-progress-circular>
                    </div>
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

                        <v-col cols="12" md="12">
                            <v-text-field
                                readonly
                                persistent-placeholder
                                type="text"
                                class="always-cursor-pointer"
                                :disabled="submitting"
                                :label="$t('Data File')"
                                :placeholder="$t('Click to select import file')"
                                v-model="fileName"
                                @click="showOpenFileDialog"
                            />
                        </v-col>
                    </v-row>
                </v-window-item>
                <v-window-item value="checkData">
                    <v-data-table
                        fixed-header
                        fixed-footer
                        show-select
                        show-expand
                        class="import-transaction-table"
                        density="compact"
                        item-value="index"
                        :height="importTransactionsTableHeight"
                        :headers="importTransactionHeaders"
                        :items="importTransactions"
                        :no-data-text="$t('No data to import')"
                        v-model:items-per-page="countPerPage"
                        v-model:page="currentPage"
                        v-model:expanded="expandedTransactions"
                    >
                        <template #header.data-table-select>
                            <v-checkbox readonly class="cursor-pointer"
                                        density="compact" width="28"
                                        :indeterminate="anyButNotAllTransactionSelected"
                                        v-model="allTransactionSelected"
                            >
                                <v-menu activator="parent" location="bottom">
                                    <v-list>
                                        <v-list-item :prepend-icon="icons.selectAll"
                                                     :title="$t('Select All')"
                                                     @click="selectAll"></v-list-item>
                                        <v-list-item :prepend-icon="icons.selectNone"
                                                     :title="$t('Select None')"
                                                     @click="selectNone"></v-list-item>
                                        <v-list-item :prepend-icon="icons.selectInverse"
                                                     :title="$t('Invert Selection')"
                                                     @click="selectInvert"></v-list-item>
                                        <v-divider class="my-2"/>
                                        <v-list-item :prepend-icon="icons.selectAll"
                                                     :title="$t('Select All in This Page')"
                                                     @click="selectAllInThisPage"></v-list-item>
                                        <v-list-item :prepend-icon="icons.selectNone"
                                                     :title="$t('Select None in This Page')"
                                                     @click="selectNoneInThisPage"></v-list-item>
                                        <v-list-item :prepend-icon="icons.selectInverse"
                                                     :title="$t('Invert Selection in This Page')"
                                                     @click="selectInvertInThisPage"></v-list-item>
                                    </v-list>
                                </v-menu>
                            </v-checkbox>
                        </template>
                        <template #item.data-table-select="{ item }">
                            <v-checkbox density="compact"
                                        :disabled="!item.valid"
                                        v-model="item.selected"></v-checkbox>
                        </template>
                        <template #item.data-table-expand="{ item, internalItem, toggleExpand }">
                            <v-icon size="small" :class="{ 'text-error': !item.valid }"
                                    :icon="icons.edit" @click="toggleExpand(internalItem)">
                            </v-icon>
                            <v-tooltip activator="parent">{{ $t('Edit') }}</v-tooltip>
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
                        <template #item.categoryId="{ item }">
                            <div class="d-flex align-center">
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
                        </template>
                        <template #item.sourceAmount="{ item }">
                            <span>{{ getTransactionDisplayAmount(item) }}</span>
                            <v-icon class="mx-1" size="13" :icon="icons.arrowRight" v-if="item.type === allTransactionTypes.Transfer && item.sourceAccountId !== item.destinationAccountId"></v-icon>
                            <span v-if="item.type === allTransactionTypes.Transfer && item.sourceAccountId !== item.destinationAccountId">{{ getTransactionDisplayDestinationAmount(item) }}</span>
                        </template>
                        <template #item.sourceAccountId="{ item }">
                            <div class="d-flex align-center">
                                <span v-if="item.sourceAccountId && item.sourceAccountId !== '0' && allAccountsMap[item.sourceAccountId]">{{ allAccountsMap[item.sourceAccountId].name }}</span>
                                <div class="text-error font-italic" v-else>
                                    <v-icon class="mr-1" :icon="icons.alert"/>
                                    <span>{{ item.originalSourceAccountName }}</span>
                                </div>
                                <v-icon class="mx-1" size="13" :icon="icons.arrowRight" v-if="item.type === allTransactionTypes.Transfer && item.sourceAccountId !== item.destinationAccountId"></v-icon>
                                <span v-if="item.type === allTransactionTypes.Transfer && item.destinationAccountId && item.destinationAccountId !== '0' && allAccountsMap[item.destinationAccountId]">{{allAccountsMap[item.destinationAccountId].name }}</span>
                                <div class="text-error font-italic" v-else-if="item.type === allTransactionTypes.Transfer && (!item.destinationAccountId || item.destinationAccountId === '0' || !allAccountsMap[item.destinationAccountId])">
                                    <v-icon class="mr-1" :icon="icons.alert"/>
                                    <span>{{ item.originalDestinationAccountName }}</span>
                                </div>
                            </div>
                        </template>
                        <template #item.geoLocation="{ item }">
                            <span class="cursor-pointer" v-if="item.geoLocation">{{ `(${item.geoLocation.longitude}, ${item.geoLocation.latitude})` }}</span>
                            <span class="cursor-pointer" v-else-if="!item.geoLocation">{{ $t('None') }}</span>
                        </template>
                        <template #item.tagIds="{ item }">
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
                        </template>
                        <template #expanded-row="{ columns, item }">
                            <tr>
                                <td :colspan="columns.length">
                                    <v-row class="py-4" style="width: 400px">
                                        <v-col cols="12" v-if="item.type === allTransactionTypes.Expense">
                                            <two-column-select primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                               primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                               primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                               secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                               secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                               secondary-hidden-field="hidden"
                                                               :disabled="loading || submitting || !hasAvailableExpenseCategories"
                                                               :show-selection-primary-text="true"
                                                               :custom-selection-primary-text="getPrimaryCategoryName(item.categoryId, allCategories[allCategoryTypes.Expense])"
                                                               :custom-selection-secondary-text="getSecondaryCategoryName(item.categoryId, allCategories[allCategoryTypes.Expense])"
                                                               :label="$t('Category')" :placeholder="$t('Category')"
                                                               :items="allCategories[allCategoryTypes.Expense]"
                                                               v-model="item.categoryId"
                                                               @update:model-value="updateTransactionData(item)">
                                            </two-column-select>
                                        </v-col>
                                        <v-col cols="12" v-if="item.type === allTransactionTypes.Income">
                                            <two-column-select primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                               primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                               primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                               secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                               secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                               secondary-hidden-field="hidden"
                                                               :disabled="loading || submitting || !hasAvailableIncomeCategories"
                                                               :show-selection-primary-text="true"
                                                               :custom-selection-primary-text="getPrimaryCategoryName(item.categoryId, allCategories[allCategoryTypes.Income])"
                                                               :custom-selection-secondary-text="getSecondaryCategoryName(item.categoryId, allCategories[allCategoryTypes.Income])"
                                                               :label="$t('Category')" :placeholder="$t('Category')"
                                                               :items="allCategories[allCategoryTypes.Income]"
                                                               v-model="item.categoryId"
                                                               @update:model-value="updateTransactionData(item)">
                                            </two-column-select>
                                        </v-col>
                                        <v-col cols="12" v-if="item.type === allTransactionTypes.Transfer">
                                            <two-column-select primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                               primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                               primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                               secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                               secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                               secondary-hidden-field="hidden"
                                                               :disabled="loading || submitting || !hasAvailableTransferCategories"
                                                               :show-selection-primary-text="true"
                                                               :custom-selection-primary-text="getPrimaryCategoryName(item.categoryId, allCategories[allCategoryTypes.Transfer])"
                                                               :custom-selection-secondary-text="getSecondaryCategoryName(item.categoryId, allCategories[allCategoryTypes.Transfer])"
                                                               :label="$t('Category')" :placeholder="$t('Category')"
                                                               :items="allCategories[allCategoryTypes.Transfer]"
                                                               v-model="item.categoryId"
                                                               @update:model-value="updateTransactionData(item)">
                                            </two-column-select>
                                        </v-col>
                                        <v-col cols="12">
                                            <two-column-select primary-key-field="id" primary-value-field="category"
                                                               primary-title-field="name" primary-footer-field="displayBalance"
                                                               primary-icon-field="icon" primary-icon-type="account"
                                                               primary-sub-items-field="accounts"
                                                               :primary-title-i18n="true"
                                                               secondary-key-field="id" secondary-value-field="id"
                                                               secondary-title-field="name" secondary-footer-field="displayBalance"
                                                               secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                               :disabled="loading || submitting || !allVisibleAccounts.length"
                                                               :custom-selection-primary-text="getSourceAccountDisplayName(item)"
                                                               :label="getSourceAccountTitle(item)"
                                                               :placeholder="getSourceAccountTitle(item)"
                                                               :items="categorizedAccounts"
                                                               v-model="item.sourceAccountId"
                                                               @update:model-value="updateTransactionData(item)">
                                            </two-column-select>
                                        </v-col>
                                        <v-col cols="12" v-if="item.type === allTransactionTypes.Transfer">
                                            <two-column-select primary-key-field="id" primary-value-field="category"
                                                               primary-title-field="name" primary-footer-field="displayBalance"
                                                               primary-icon-field="icon" primary-icon-type="account"
                                                               primary-sub-items-field="accounts"
                                                               :primary-title-i18n="true"
                                                               secondary-key-field="id" secondary-value-field="id"
                                                               secondary-title-field="name" secondary-footer-field="displayBalance"
                                                               secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                               :disabled="loading || submitting || !allVisibleAccounts.length"
                                                               :custom-selection-primary-text="getDestinationAccountDisplayName(item)"
                                                               :label="$t('Destination Account')"
                                                               :placeholder="$t('Destination Account')"
                                                               :items="categorizedAccounts"
                                                               v-model="item.destinationAccountId"
                                                               @update:model-value="updateTransactionData(item)">
                                            </two-column-select>
                                        </v-col>
                                        <v-col cols="12">
                                            <v-autocomplete
                                                item-title="name"
                                                item-value="id"
                                                auto-select-first
                                                persistent-placeholder
                                                multiple
                                                chips
                                                closable-chips
                                                :disabled="loading || submitting"
                                                :label="$t('Tags')"
                                                :placeholder="$t('None')"
                                                :items="allTags"
                                                :no-data-text="$t('No available tag')"
                                                v-model="item.tagIds"
                                                @update:model-value="updateTransactionData(item)"
                                            >
                                                <template #chip="{ props, index }">
                                                    <v-chip :class="{ 'font-italic': !isTagValid(item, index) }"
                                                            :prepend-icon="isTagValid(item, index) ? icons.tag : icons.alert"
                                                            :color="isTagValid(item, index) ? 'default' : 'error'"
                                                            :text="isTagValid(item, index) ? allTagsMap[item.tagIds[index]].name : item.originalTagNames[index]"
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
                                        </v-col>
                                    </v-row>
                                </td>
                            </tr>
                        </template>
                        <template #bottom>
                            <div class="title-and-toolbar d-flex align-center text-no-wrap mt-2"
                                 v-if="importTransactions && importTransactions.length > 10">
                                <span>{{ $t('format.misc.selectedCount', { count: selectedImportTransactionCount, totalCount: importTransactions.length }) }}</span>
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
                                              v-model="currentPage"></v-pagination>
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
                <v-btn color="teal" :disabled="submitting || selectedImportTransactionCount < 1"
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

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
    <input ref="fileInput" type="file" style="display: none" :accept="supportedImportFileExtensions" @change="setImportFile($event)" />
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useAccountsStore } from '@/stores/account.js';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';
import { useTransactionTagsStore } from '@/stores/transactionTag.js';
import { useTransactionsStore } from '@/stores/transaction.js';
import { useOverviewStore } from '@/stores/overview.js';
import { useStatisticsStore } from '@/stores/statistics.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import categoryConstants from '@/consts/category.js';
import transactionConstants from '@/consts/transaction.js';
import { getNameByKeyValue } from '@/lib/common.js';
import { generateRandomUUID } from '@/lib/misc.js';
import logger from '@/lib/logger.js';
import {
    parseDateFromUnixTime,
    getUnixTime,
    getUtcOffsetByUtcOffsetMinutes,
    getTimezoneOffsetMinutes
} from '@/lib/datetime.js';
import {
    getTransactionPrimaryCategoryName,
    getTransactionSecondaryCategoryName,
    getFirstAvailableCategoryId
} from '@/lib/category.js';

import {
    mdiClose,
    mdiArrowRight,
    mdiCheck,
    mdiSelectAll,
    mdiSelect,
    mdiSelectInverse,
    mdiPencilOutline,
    mdiAlertOutline,
    mdiPound
} from '@mdi/js';

export default {
    props: [
        'persistent',
        'show'
    ],
    expose: [
        'open'
    ],
    data() {
        return {
            showState: false,
            clientSessionId: '',
            currentStep: 'uploadFile',
            fileType: 'ezbookkeeping_csv',
            importFile: null,
            importTransactions: null,
            expandedTransactions: [],
            currentPage: 1,
            countPerPage: 10,
            importedCount: null,
            loading: true,
            submitting: false,
            resolve: null,
            reject: null,
            icons: {
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
        allTransactionTypes() {
            return transactionConstants.allTransactionTypes;
        },
        allCategoryTypes() {
            return categoryConstants.allCategoryTypes;
        },
        allAccounts() {
            return this.accountsStore.allPlainAccounts;
        },
        allVisibleAccounts() {
            return this.accountsStore.allVisiblePlainAccounts;
        },
        categorizedAccounts() {
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
            return getNameByKeyValue(this.allSupportedImportFileTypes, this.fileType, 'type', 'extensions');
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
                { value: 'valid', key: 'data-table-expand', sortable: true, nowrap: true },
                { value: 'time', title: this.$t('Transaction Time'), sortable: true, nowrap: true },
                { value: 'type', title: this.$t('Type'), sortable: true, nowrap: true },
                { value: 'categoryId', title: this.$t('Category'), sortable: true, nowrap: true },
                { value: 'sourceAmount', title: this.$t('Amount'), sortable: true, nowrap: true },
                { value: 'sourceAccountId', title: this.$t('Account'), sortable: true, nowrap: true },
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
        totalPageCount() {
            if (!this.importTransactions || this.importTransactions.length < 1) {
                return 1;
            }

            return Math.ceil(this.importTransactions.length / this.countPerPage);
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
        anyButNotAllTransactionSelected: {
            get: function () {
                return this.selectedImportTransactionCount > 0 && this.selectedImportTransactionCount !== this.importTransactions.length;
            }
        },
        allTransactionSelected: {
            get: function () {
                return this.selectedImportTransactionCount === this.importTransactions.length;
            }
        }
    },
    watch: {
        fileType: function () {
            this.importFile = null;
            this.importTransactions = null;
            this.expandedTransactions = [];
            this.currentPage = 1;
            this.countPerPage = 10;
        }
    },
    methods: {
        open() {
            const self = this;
            self.fileType = 'ezbookkeeping_csv';
            self.currentStep = 'uploadFile';
            self.importFile = null;
            self.importTransactions = null;
            self.expandedTransactions = [];
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

            self.transactionsStore.parseImportTransaction({
                fileType: self.fileType,
                importFile: self.importFile
            }).then(response => {
                const parsedTransactions = response.items;

                if (parsedTransactions) {
                    for (let i = 0; i < parsedTransactions.length; i++) {
                        const transaction = parsedTransactions[i];
                        transaction.index = i;
                        transaction.selected = false;
                        transaction.valid = self.isTransactionValid(transaction);
                    }
                }

                self.importTransactions = parsedTransactions;
                self.expandedTransactions = [];
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
            const transactions = [];

            for (let i = 0; i < self.importTransactions.length; i++) {
                const transaction = self.importTransactions[i];

                if (transaction.valid && transaction.selected) {
                    transactions.push(transaction);
                }
            }

            if (transactions.length < 1) {
                self.$refs.snackbar.showError('No data to import');
                return;
            }

            self.$refs.confirmDialog.open('format.misc.confirmImportTransactions', {
                count: transactions.length
            }).then(() => {
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
        selectAll() {
            for (let i = 0; i < this.importTransactions.length; i++) {
                if (this.importTransactions[i].valid) {
                    this.importTransactions[i].selected = true;
                }
            }
        },
        selectNone() {
            for (let i = 0; i < this.importTransactions.length; i++) {
                this.importTransactions[i].selected = false;
            }
        },
        selectInvert() {
            for (let i = 0; i < this.importTransactions.length; i++) {
                if (this.importTransactions[i].valid || this.importTransactions[i].selected) {
                    this.importTransactions[i].selected = !this.importTransactions[i].selected;
                }
            }
        },
        selectAllInThisPage() {
            for (let i = Math.max(0, (this.currentPage - 1) * this.countPerPage); i < Math.min(this.importTransactions.length, this.currentPage * this.countPerPage); i++) {
                if (this.importTransactions[i] && this.importTransactions[i].valid) {
                    this.importTransactions[i].selected = true;
                }
            }
        },
        selectNoneInThisPage() {
            for (let i = Math.max(0, (this.currentPage - 1) * this.countPerPage); i < Math.min(this.importTransactions.length, this.currentPage * this.countPerPage); i++) {
                if (this.importTransactions[i]) {
                    this.importTransactions[i].selected = false;
                }
            }
        },
        selectInvertInThisPage() {
            for (let i = Math.max(0, (this.currentPage - 1) * this.countPerPage); i < Math.min(this.importTransactions.length, this.currentPage * this.countPerPage); i++) {
                if (this.importTransactions[i] && (this.importTransactions[i].valid || this.importTransactions[i].selected)) {
                    this.importTransactions[i].selected = !this.importTransactions[i].selected;
                }
            }
        },
        updateTransactionData(transaction) {
            transaction.valid = this.isTransactionValid(transaction);
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
        isTagValid(transaction, tagIndex) {
            if (!transaction || !transaction.tagIds || !transaction.tagIds[tagIndex]) {
                return false;
            }

            if (transaction.tagIds[tagIndex] === '0') {
                return false;
            }

            const tagId = transaction.tagIds[tagIndex];
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
