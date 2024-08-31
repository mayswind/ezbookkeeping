<template>
    <v-dialog width="1000" :persistent="!!persistent && mode !== 'view'" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-8">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <div class="d-flex w-100 align-center justify-center">
                        <h4 class="text-h4">{{ $t(title) }}</h4>
                        <v-progress-circular indeterminate size="22" class="ml-2" v-if="loading"></v-progress-circular>
                    </div>
                    <v-btn density="comfortable" color="default" variant="text" class="ml-2" :icon="true"
                           :disabled="loading || submitting" v-if="mode !== 'view'">
                        <v-icon :icon="icons.more" />
                        <v-menu activator="parent">
                            <v-list>
                                <v-list-item :prepend-icon="icons.swap"
                                             :title="$t('Swap Account')"
                                             v-if="transaction.type === allTransactionTypes.Transfer"
                                             @click="swapTransactionData(true, false)"></v-list-item>
                                <v-list-item :prepend-icon="icons.swap"
                                             :title="$t('Swap Amount')"
                                             v-if="transaction.type === allTransactionTypes.Transfer"
                                             @click="swapTransactionData(false, true)"></v-list-item>
                                <v-list-item :prepend-icon="icons.swap"
                                             :title="$t('Swap Account and Amount')"
                                             v-if="transaction.type === allTransactionTypes.Transfer"
                                             @click="swapTransactionData(true, true)"></v-list-item>
                                <v-divider v-if="transaction.type === allTransactionTypes.Transfer" />
                                <v-list-item :prepend-icon="icons.show"
                                             :title="$t('Show Amount')"
                                             v-if="transaction.hideAmount" @click="transaction.hideAmount = false"></v-list-item>
                                <v-list-item :prepend-icon="icons.hide"
                                             :title="$t('Hide Amount')"
                                             v-if="!transaction.hideAmount" @click="transaction.hideAmount = true"></v-list-item>
                            </v-list>
                        </v-menu>
                    </v-btn>
                </div>
            </template>
            <v-card-text class="d-flex flex-column flex-md-row mt-md-4 pt-0">
                <div class="mb-4">
                    <v-tabs class="v-tabs-pill" direction="vertical" :class="{ 'readonly': type === 'transaction' && mode !== 'add' }"
                            :disabled="loading || submitting" v-model="transaction.type">
                        <v-tab :value="allTransactionTypes.Expense" :disabled="type === 'transaction' && mode !== 'add' && transaction.type !== allTransactionTypes.Expense" v-if="transaction.type !== allTransactionTypes.ModifyBalance">
                            <span>{{ $t('Expense') }}</span>
                        </v-tab>
                        <v-tab :value="allTransactionTypes.Income" :disabled="type === 'transaction' && mode !== 'add' && transaction.type !== allTransactionTypes.Income" v-if="transaction.type !== allTransactionTypes.ModifyBalance">
                            <span>{{ $t('Income') }}</span>
                        </v-tab>
                        <v-tab :value="allTransactionTypes.Transfer" :disabled="type === 'transaction' && mode !== 'add' && transaction.type !== allTransactionTypes.Transfer" v-if="transaction.type !== allTransactionTypes.ModifyBalance">
                            <span>{{ $t('Transfer') }}</span>
                        </v-tab>
                        <v-tab :value="allTransactionTypes.ModifyBalance" v-if="type === 'transaction' && transaction.type === allTransactionTypes.ModifyBalance">
                            <span>{{ $t('Modify Balance') }}</span>
                        </v-tab>
                    </v-tabs>
                    <v-divider class="my-2"/>
                    <v-tabs direction="vertical" :disabled="loading || submitting" v-model="activeTab">
                        <v-tab value="basicInfo">
                            <span>{{ $t('Basic Information') }}</span>
                        </v-tab>
                        <v-tab value="map" :disabled="!transaction.geoLocation" v-if="type === 'transaction' && mapProvider">
                            <span>{{ $t('Location on Map') }}</span>
                        </v-tab>
                        <v-tab value="pictures" :disabled="mode !== 'add' && mode !== 'edit' && (!transaction.pictures || !transaction.pictures.length)" v-if="type === 'transaction' && isTransactionPicturesEnabled">
                            <span>{{ $t('Pictures') }}</span>
                        </v-tab>
                    </v-tabs>
                </div>

                <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container ml-md-5"
                          v-model="activeTab">
                    <v-window-item value="basicInfo">
                        <v-form class="mt-2">
                            <v-row>
                                <v-col cols="12" v-if="type === 'template'">
                                    <v-text-field
                                        type="text"
                                        persistent-placeholder
                                        :disabled="loading || submitting"
                                        :label="$t('Template Name')"
                                        :placeholder="$t('Template Name')"
                                        v-model="transaction.name"
                                    />
                                </v-col>
                                <v-col cols="12" :md="transaction.type === allTransactionTypes.Transfer ? 6 : 12">
                                    <amount-input class="transaction-edit-amount font-weight-bold"
                                                  :color="sourceAmountColor"
                                                  :readonly="mode === 'view'"
                                                  :disabled="loading || submitting"
                                                  :persistent-placeholder="true"
                                                  :hide="transaction.hideAmount"
                                                  :label="$t(sourceAmountName)"
                                                  :placeholder="$t(sourceAmountName)"
                                                  v-model="transaction.sourceAmount"/>
                                </v-col>
                                <v-col cols="12" :md="6" v-if="transaction.type === allTransactionTypes.Transfer">
                                    <amount-input class="transaction-edit-amount font-weight-bold" color="primary"
                                                  :readonly="mode === 'view'"
                                                  :disabled="loading || submitting"
                                                  :persistent-placeholder="true"
                                                  :hide="transaction.hideAmount"
                                                  :label="transferInAmountTitle"
                                                  :placeholder="$t('Transfer In Amount')"
                                                  v-model="transaction.destinationAmount"/>
                                </v-col>
                                <v-col cols="12" md="12" v-if="transaction.type === allTransactionTypes.Expense">
                                    <two-column-select primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                       primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                       primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                       secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                       secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                       secondary-hidden-field="hidden"
                                                       :readonly="mode === 'view'"
                                                       :disabled="loading || submitting || !hasAvailableExpenseCategories"
                                                       :show-selection-primary-text="true"
                                                       :custom-selection-primary-text="getPrimaryCategoryName(transaction.expenseCategory, allCategories[allCategoryTypes.Expense])"
                                                       :custom-selection-secondary-text="getSecondaryCategoryName(transaction.expenseCategory, allCategories[allCategoryTypes.Expense])"
                                                       :label="$t('Category')" :placeholder="$t('Category')"
                                                       :items="allCategories[allCategoryTypes.Expense]"
                                                       v-model="transaction.expenseCategory">
                                    </two-column-select>
                                </v-col>
                                <v-col cols="12" md="12" v-if="transaction.type === allTransactionTypes.Income">
                                    <two-column-select primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                       primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                       primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                       secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                       secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                       secondary-hidden-field="hidden"
                                                       :readonly="mode === 'view'"
                                                       :disabled="loading || submitting || !hasAvailableIncomeCategories"
                                                       :show-selection-primary-text="true"
                                                       :custom-selection-primary-text="getPrimaryCategoryName(transaction.incomeCategory, allCategories[allCategoryTypes.Income])"
                                                       :custom-selection-secondary-text="getSecondaryCategoryName(transaction.incomeCategory, allCategories[allCategoryTypes.Income])"
                                                       :label="$t('Category')" :placeholder="$t('Category')"
                                                       :items="allCategories[allCategoryTypes.Income]"
                                                       v-model="transaction.incomeCategory">
                                    </two-column-select>
                                </v-col>
                                <v-col cols="12" md="12" v-if="transaction.type === allTransactionTypes.Transfer">
                                    <two-column-select primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                       primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                       primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                       secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                       secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                       secondary-hidden-field="hidden"
                                                       :readonly="mode === 'view'"
                                                       :disabled="loading || submitting || !hasAvailableTransferCategories"
                                                       :show-selection-primary-text="true"
                                                       :custom-selection-primary-text="getPrimaryCategoryName(transaction.transferCategory, allCategories[allCategoryTypes.Transfer])"
                                                       :custom-selection-secondary-text="getSecondaryCategoryName(transaction.transferCategory, allCategories[allCategoryTypes.Transfer])"
                                                       :label="$t('Category')" :placeholder="$t('Category')"
                                                       :items="allCategories[allCategoryTypes.Transfer]"
                                                       v-model="transaction.transferCategory">
                                    </two-column-select>
                                </v-col>
                                <v-col cols="12" :md="transaction.type === allTransactionTypes.Transfer ? 6 : 12">
                                    <two-column-select primary-key-field="id" primary-value-field="category"
                                                       primary-title-field="name" primary-footer-field="displayBalance"
                                                       primary-icon-field="icon" primary-icon-type="account"
                                                       primary-sub-items-field="accounts"
                                                       :primary-title-i18n="true"
                                                       secondary-key-field="id" secondary-value-field="id"
                                                       secondary-title-field="name" secondary-footer-field="displayBalance"
                                                       secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                       :readonly="mode === 'view'"
                                                       :disabled="loading || submitting || !allVisibleAccounts.length"
                                                       :custom-selection-primary-text="sourceAccountName"
                                                       :label="$t(sourceAccountTitle)"
                                                       :placeholder="$t(sourceAccountTitle)"
                                                       :items="categorizedAccounts"
                                                       v-model="transaction.sourceAccountId">
                                    </two-column-select>
                                </v-col>
                                <v-col cols="12" md="6" v-if="transaction.type === allTransactionTypes.Transfer">
                                    <two-column-select primary-key-field="id" primary-value-field="category"
                                                       primary-title-field="name" primary-footer-field="displayBalance"
                                                       primary-icon-field="icon" primary-icon-type="account"
                                                       primary-sub-items-field="accounts"
                                                       :primary-title-i18n="true"
                                                       secondary-key-field="id" secondary-value-field="id"
                                                       secondary-title-field="name" secondary-footer-field="displayBalance"
                                                       secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                       :readonly="mode === 'view'"
                                                       :disabled="loading || submitting || !allVisibleAccounts.length"
                                                       :custom-selection-primary-text="destinationAccountName"
                                                       :label="$t('Destination Account')"
                                                       :placeholder="$t('Destination Account')"
                                                       :items="categorizedAccounts"
                                                       v-model="transaction.destinationAccountId">
                                    </two-column-select>
                                </v-col>
                                <v-col cols="12" md="6" v-if="type === 'transaction'">
                                    <date-time-select
                                        :readonly="mode === 'view'"
                                        :disabled="loading || submitting"
                                        :label="$t('Transaction Time')"
                                        v-model="transaction.time"
                                        @error="showDateTimeError" />
                                </v-col>
                                <v-col cols="12" md="6" v-if="type === 'template' && transaction.templateType === allTemplateTypes.Schedule">
                                    <schedule-frequency-select
                                        :readonly="mode === 'view'"
                                        :disabled="loading || submitting"
                                        :label="$t('Scheduled Transaction Frequency')"
                                        v-model:type="transaction.scheduledFrequencyType"
                                        v-model="transaction.scheduledFrequency" />
                                </v-col>
                                <v-col cols="12" md="6" v-if="type === 'transaction' || (type === 'template' && transaction.templateType === allTemplateTypes.Schedule)">
                                    <v-autocomplete
                                        class="transaction-edit-timezone"
                                        item-title="displayNameWithUtcOffset"
                                        item-value="name"
                                        auto-select-first
                                        persistent-placeholder
                                        :readonly="mode === 'view'"
                                        :disabled="loading || submitting"
                                        :label="$t('Transaction Timezone')"
                                        :placeholder="!transaction.timeZone && transaction.timeZone !== '' ? `(${transactionDisplayTimezone}) ${transactionTimezoneTimeDifference}` : $t('Timezone')"
                                        :items="allTimezones"
                                        :no-data-text="$t('No results')"
                                        v-model="transaction.timeZone"
                                    >
                                        <template #selection="{ item }">
                                            <span class="text-truncate" v-if="transaction.timeZone || transaction.timeZone === ''">
                                                {{ item.title }}
                                            </span>
                                        </template>
                                    </v-autocomplete>
                                </v-col>
                                <v-col cols="12" md="12" v-if="type === 'transaction'">
                                    <v-select
                                        persistent-placeholder
                                        :readonly="mode === 'view'"
                                        :disabled="loading || submitting"
                                        :label="$t('Geographic Location')"
                                        v-model="transaction"
                                        v-model:menu="geoMenuState"
                                    >
                                        <template #selection>
                                            <span class="cursor-pointer" v-if="transaction.geoLocation">{{ `(${transaction.geoLocation.longitude}, ${transaction.geoLocation.latitude})` }}</span>
                                            <span class="cursor-pointer" v-else-if="!transaction.geoLocation">{{ geoLocationStatusInfo }}</span>
                                        </template>

                                        <template #no-data>
                                            <v-list>
                                                <v-list-item v-if="mode !== 'view'" @click="updateGeoLocation(true)">{{ $t('Update Geographic Location') }}</v-list-item>
                                                <v-list-item v-if="mode !== 'view'" @click="clearGeoLocation">{{ $t('Clear Geographic Location') }}</v-list-item>
                                            </v-list>
                                        </template>
                                    </v-select>
                                </v-col>
                                <v-col cols="12" md="12">
                                    <v-autocomplete
                                        item-title="name"
                                        item-value="id"
                                        auto-select-first
                                        persistent-placeholder
                                        multiple
                                        chips
                                        :closable-chips="mode !== 'view'"
                                        :readonly="mode === 'view'"
                                        :disabled="loading || submitting"
                                        :label="$t('Tags')"
                                        :placeholder="$t('None')"
                                        :items="allTags"
                                        :no-data-text="$t('No available tag')"
                                        v-model="transaction.tagIds"
                                    >
                                        <template #chip="{ props, item }">
                                            <v-chip :prepend-icon="icons.tag" :text="item.title" v-bind="props"/>
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
                                <v-col cols="12" md="12">
                                    <v-textarea
                                        type="text"
                                        persistent-placeholder
                                        rows="3"
                                        :readonly="mode === 'view'"
                                        :disabled="loading || submitting"
                                        :label="$t('Description')"
                                        :placeholder="$t('Your transaction description (optional)')"
                                        v-model="transaction.comment"
                                    />
                                </v-col>
                            </v-row>
                        </v-form>
                    </v-window-item>
                    <v-window-item value="map">
                        <v-row>
                            <v-col cols="12" md="12">
                                <map-view ref="map" map-class="transaction-edit-map-view" :geo-location="transaction.geoLocation">
                                    <template #error-title="{ mapSupported, mapDependencyLoaded }">
                                        <span class="text-subtitle-1" v-if="!mapSupported"><b>{{ $t('Unsupported Map Provider') }}</b></span>
                                        <span class="text-subtitle-1" v-else-if="!mapDependencyLoaded"><b>{{ $t('Cannot Initialize Map') }}</b></span>
                                    </template>
                                    <template #error-content>
                                        <p class="text-body-1">
                                            {{ $t('Please refresh the page and try again. If the error persists, ensure that the server\'s map settings are correctly configured.') }}
                                        </p>
                                    </template>
                                </map-view>
                            </v-col>
                        </v-row>
                    </v-window-item>
                    <v-window-item value="pictures">
                        <v-row class="transaction-pictures align-content-start" :class="{ 'readonly': submitting || uploadingPicture || removingPictureId }">
                            <v-col :key="picIdx" cols="6" md="3" v-for="(pictureInfo, picIdx) in transaction.pictures">
                                <v-avatar rounded="lg" variant="tonal" size="160"
                                          class="cursor-pointer transaction-picture"
                                          color="rgba(0,0,0,0)" @click="viewOrRemovePicture(pictureInfo)">
                                    <v-img :src="getTransactionPictureUrl(pictureInfo)">
                                        <template #placeholder>
                                            <div class="d-flex align-center justify-center fill-height bg-light-primary">
                                                <v-progress-circular color="grey-500" indeterminate size="48"></v-progress-circular>
                                            </div>
                                        </template>
                                    </v-img>
                                    <div class="picture-control-icon" :class="{ 'show-control-icon': pictureInfo.pictureId === removingPictureId }">
                                        <v-icon size="64" :icon="icons.remove" v-if="(mode === 'add' || mode === 'edit') && pictureInfo.pictureId !== removingPictureId"/>
                                        <v-progress-circular color="grey-500" indeterminate size="48" v-if="(mode === 'add' || mode === 'edit') && pictureInfo.pictureId === removingPictureId"></v-progress-circular>
                                        <v-icon size="64" :icon="icons.fullscreen" v-if="mode !== 'add' && mode !== 'edit'"/>
                                    </div>
                                </v-avatar>
                            </v-col>
                            <v-col cols="6" md="3" v-if="canAddTransactionPicture">
                                <v-avatar rounded="lg" variant="tonal" size="160"
                                          class="transaction-picture transaction-picture-add"
                                          :class="{ 'enabled': !submitting, 'cursor-pointer': !submitting }"
                                          color="rgba(0,0,0,0)" @click="showOpenPictureDialog">
                                    <v-tooltip activator="parent" v-if="!submitting">{{ $t('Add Picture') }}</v-tooltip>
                                    <v-icon class="transaction-picture-add-icon" size="56" :icon="icons.add" v-if="!uploadingPicture"/>
                                    <v-progress-circular color="grey-500" indeterminate size="48" v-if="uploadingPicture"></v-progress-circular>
                                </v-avatar>
                            </v-col>
                        </v-row>
                    </v-window-item>
                </v-window>
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center flex-wrap mt-2 mt-sm-4 mt-md-6 gap-4">
                    <v-btn :disabled="inputIsEmpty || loading || submitting" v-if="mode !== 'view'" @click="save">
                        {{ $t(saveButtonTitle) }}
                        <v-progress-circular indeterminate size="22" class="ml-2" v-if="submitting"></v-progress-circular>
                    </v-btn>
                    <v-btn variant="tonal" :disabled="loading || submitting"
                           v-if="mode === 'view' && transaction.type !== allTransactionTypes.ModifyBalance"
                           @click="duplicate">{{ $t('Duplicate') }}</v-btn>
                    <v-btn color="warning" variant="tonal" :disabled="loading || submitting"
                           v-if="mode === 'view' && originalTransactionEditable && transaction.type !== allTransactionTypes.ModifyBalance"
                           @click="edit">{{ $t('Edit') }}</v-btn>
                    <v-btn color="error" variant="tonal" :disabled="loading || submitting"
                           v-if="mode === 'view' && originalTransactionEditable" @click="remove">
                        {{ $t('Delete') }}
                        <v-progress-circular indeterminate size="22" class="ml-2" v-if="submitting"></v-progress-circular>
                    </v-btn>
                    <v-btn color="secondary" variant="tonal" :disabled="loading || submitting"
                           @click="cancel">{{ $t(cancelButtonTitle) }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
    <input ref="pictureInput" type="file" style="display: none" :accept="supportedImageExtensions" @change="uploadPicture($event)" />
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useAccountsStore } from '@/stores/account.js';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';
import { useTransactionTagsStore } from '@/stores/transactionTag.js';
import { useTransactionsStore } from '@/stores/transaction.js';
import { useTransactionTemplatesStore } from '@/stores/transactionTemplate.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import fileConstants from '@/consts/file.js';
import categoryConstants from '@/consts/category.js';
import transactionConstants from '@/consts/transaction.js';
import templateConstants from '@/consts/template.js';
import logger from '@/lib/logger.js';
import {
    isArray,
    getNameByKeyValue
} from '@/lib/common.js';
import {
    getUtcOffsetByUtcOffsetMinutes,
    getTimezoneOffsetMinutes,
    getCurrentUnixTime
} from '@/lib/datetime.js';
import { generateRandomUUID } from '@/lib/misc.js';
import {
    getTransactionPrimaryCategoryName,
    getTransactionSecondaryCategoryName,
    getFirstAvailableCategoryId
} from '@/lib/category.js';
import { setTransactionModelByTransaction } from '@/lib/transaction.js';
import {
    isTransactionPicturesEnabled,
    getMapProvider
} from '@/lib/server_settings.js';

import {
    mdiDotsVertical,
    mdiEyeOffOutline,
    mdiEyeOutline,
    mdiSwapHorizontal,
    mdiPound,
    mdiImageOutline,
    mdiImagePlusOutline,
    mdiTrashCanOutline,
    mdiFullscreen
} from '@mdi/js';

export default {
    props: [
        'persistent',
        'type',
        'show'
    ],
    expose: [
        'open'
    ],
    data() {
        const transactionsStore = useTransactionsStore();
        const newTransaction = transactionsStore.generateNewTransactionModel();

        return {
            showState: false,
            mode: 'add',
            activeTab: 'basicInfo',
            editId: null,
            originalTransactionEditable: false,
            clientSessionId: '',
            loading: true,
            transaction: newTransaction,
            geoLocationStatus: null,
            geoMenuState: false,
            submitting: false,
            uploadingPicture: false,
            removingPictureId: '',
            isSupportGeoLocation: !!navigator.geolocation,
            resolve: null,
            reject: null,
            icons: {
                more: mdiDotsVertical,
                show: mdiEyeOutline,
                hide: mdiEyeOffOutline,
                swap: mdiSwapHorizontal,
                tag: mdiPound,
                picture: mdiImageOutline ,
                add: mdiImagePlusOutline,
                remove: mdiTrashCanOutline,
                fullscreen : mdiFullscreen
            }
        };
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore, useAccountsStore, useTransactionCategoriesStore, useTransactionTagsStore, useTransactionsStore, useTransactionTemplatesStore, useExchangeRatesStore),
        title() {
            if (this.type === 'transaction') {
                if (this.mode === 'add') {
                    return 'Add Transaction';
                } else if (this.mode === 'edit') {
                    return 'Edit Transaction';
                } else {
                    return 'Transaction Detail';
                }
            } else if (this.type === 'template' && this.transaction.templateType === templateConstants.allTemplateTypes.Normal) {
                if (this.mode === 'add') {
                    return 'Add Transaction Template';
                } else if (this.mode === 'edit') {
                    return 'Edit Transaction Template';
                }
            } else if (this.type === 'template' && this.transaction.templateType === templateConstants.allTemplateTypes.Schedule) {
                if (this.mode === 'add') {
                    return 'Add Scheduled Transaction';
                } else if (this.mode === 'edit') {
                    return 'Edit Scheduled Transaction';
                }
            }

            return '';
        },
        saveButtonTitle() {
            if (this.mode === 'add') {
                return 'Add';
            } else {
                return 'Save';
            }
        },
        cancelButtonTitle() {
            if (this.mode === 'view') {
                return 'Close';
            } else {
                return 'Cancel';
            }
        },
        sourceAmountName() {
            if (this.transaction.type === this.allTransactionTypes.Expense) {
                return 'Expense Amount';
            } else if (this.transaction.type === this.allTransactionTypes.Income) {
                return 'Income Amount';
            } else if (this.transaction.type === this.allTransactionTypes.Transfer) {
                return 'Transfer Out Amount';
            } else {
                return 'Amount';
            }
        },
        sourceAccountTitle() {
            if (this.transaction.type === this.allTransactionTypes.Expense || this.transaction.type === this.allTransactionTypes.Income) {
                return 'Account';
            } else if (this.transaction.type === this.allTransactionTypes.Transfer) {
                return 'Source Account';
            } else {
                return 'Account';
            }
        },
        transferInAmountTitle() {
            const sourceAccount = this.allAccountsMap[this.transaction.sourceAccountId];
            const destinationAccount = this.allAccountsMap[this.transaction.destinationAccountId];

            if (!sourceAccount || !destinationAccount || sourceAccount.currency === destinationAccount.currency) {
                return this.$t('Transfer In Amount');
            }

            const fromExchangeRate = this.exchangeRatesStore.latestExchangeRateMap[sourceAccount.currency];
            const toExchangeRate = this.exchangeRatesStore.latestExchangeRateMap[destinationAccount.currency];
            const amountRate = this.$locale.getAdaptiveAmountRate(this.userStore, this.transaction.sourceAmount, this.transaction.destinationAmount, fromExchangeRate, toExchangeRate);

            if (!amountRate) {
                return this.$t('Transfer In Amount');
            }

            return this.$t('Transfer In Amount') + ` (${amountRate})`;
        },
        defaultCurrency() {
            return this.userStore.currentUserDefaultCurrency;
        },
        defaultAccountId() {
            return this.userStore.currentUserDefaultAccountId;
        },
        allTransactionTypes() {
            return transactionConstants.allTransactionTypes;
        },
        allCategoryTypes() {
            return categoryConstants.allCategoryTypes;
        },
        allTemplateTypes() {
            return templateConstants.allTemplateTypes;
        },
        allTimezones() {
            return this.$locale.getAllTimezones(true);
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
        supportedImageExtensions() {
            return fileConstants.supportedImageExtensions;
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
        sourceAccountName() {
            if (this.transaction.sourceAccountId) {
                return getNameByKeyValue(this.allAccounts, this.transaction.sourceAccountId, 'id', 'name');
            } else {
                return this.$t('None');
            }
        },
        destinationAccountName() {
            if (this.transaction.destinationAccountId) {
                return getNameByKeyValue(this.allAccounts, this.transaction.destinationAccountId, 'id', 'name');
            } else {
                return this.$t('None');
            }
        },
        transactionDisplayTimezone() {
            return `UTC${getUtcOffsetByUtcOffsetMinutes(this.transaction.utcOffset)}`;
        },
        transactionTimezoneTimeDifference() {
            return this.$locale.getTimezoneDifferenceDisplayText(this.transaction.utcOffset);
        },
        sourceAmountColor() {
            if (this.transaction.type === this.allTransactionTypes.Expense) {
                return 'expense';
            } else if (this.transaction.type === this.allTransactionTypes.Income) {
                return 'income';
            } else if (this.transaction.type === this.allTransactionTypes.Transfer) {
                return 'primary';
            }

            return null;
        },
        geoLocationStatusInfo() {
            if (this.geoLocationStatus === 'success') {
                return '';
            } else if (this.geoLocationStatus === 'getting') {
                return this.$t('Getting Location...');
            } else {
                return this.$t('No Location');
            }
        },
        showAccountBalance() {
            return this.settingsStore.appSettings.showAccountBalance;
        },
        isTransactionPicturesEnabled() {
            return isTransactionPicturesEnabled();
        },
        canAddTransactionPicture() {
            if (this.type !== 'transaction' || (this.mode !== 'add' && this.mode !== 'edit')) {
                return false;
            }

            return !isArray(this.transaction.pictures) || this.transaction.pictures.length < 10;
        },
        mapProvider() {
            return getMapProvider();
        },
        inputIsEmpty() {
            return !!this.inputEmptyProblemMessage;
        },
        inputEmptyProblemMessage() {
            if (this.transaction.type === this.allTransactionTypes.Expense) {
                if (!this.transaction.expenseCategory || this.transaction.expenseCategory === '') {
                    return 'Transaction category cannot be blank';
                }

                if (!this.transaction.sourceAccountId || this.transaction.sourceAccountId === '') {
                    return 'Transaction account cannot be blank';
                }
            } else if (this.transaction.type === this.allTransactionTypes.Income) {
                if (!this.transaction.incomeCategory || this.transaction.incomeCategory === '') {
                    return 'Transaction category cannot be blank';
                }

                if (!this.transaction.sourceAccountId || this.transaction.sourceAccountId === '') {
                    return 'Transaction account cannot be blank';
                }
            } else if (this.transaction.type === this.allTransactionTypes.Transfer) {
                if (!this.transaction.transferCategory || this.transaction.transferCategory === '') {
                    return 'Transaction category cannot be blank';
                }

                if (!this.transaction.sourceAccountId || this.transaction.sourceAccountId === '') {
                    return 'Source account cannot be blank';
                }

                if (!this.transaction.destinationAccountId || this.transaction.destinationAccountId === '') {
                    return 'Destination account cannot be blank';
                }
            }

            if (this.type === 'template') {
                if (!this.transaction.name) {
                    return 'Template name cannot be blank';
                }
            }

            return null;
        }
    },
    watch: {
        'activeTab': function (newValue) {
            if (newValue === 'map') {
                this.$nextTick(() => {
                    if (this.$refs.map) {
                        this.$refs.map.init();
                    }
                });
            }
        },
        'transaction.sourceAmount': function (newValue, oldValue) {
            if (this.mode === 'view' || this.loading) {
                return;
            }

            this.transactionsStore.setTransactionSuitableDestinationAmount(this.transaction, oldValue, newValue);
        },
        'transaction.destinationAmount': function (newValue) {
            if (this.mode === 'view' || this.loading) {
                return;
            }

            if (this.transaction.type === this.allTransactionTypes.Expense || this.transaction.type === this.allTransactionTypes.Income) {
                this.transaction.sourceAmount = newValue;
            }
        },
        'transaction.timeZone': function (newValue) {
            for (let i = 0; i < this.allTimezones.length; i++) {
                if (this.allTimezones[i].name === newValue) {
                    this.transaction.utcOffset = this.allTimezones[i].utcOffsetMinutes;
                    break;
                }
            }
        }
    },
    methods: {
        open(options) {
            const self = this;
            self.showState = true;
            self.activeTab = 'basicInfo';
            self.loading = true;
            self.submitting = false;
            self.geoLocationStatus = null;
            self.originalTransactionEditable = false;

            const newTransaction = self.transactionsStore.generateNewTransactionModel(options.type);
            self.setTransaction(newTransaction, options, true, false);

            const promises = [
                self.accountsStore.loadAllAccounts({ force: false }),
                self.transactionCategoriesStore.loadAllCategories({ force: false }),
                self.transactionTagsStore.loadAllTags({ force: false })
            ];

            if (self.type === 'transaction') {
                if (options && options.id) {
                    if (options.currentTransaction) {
                        self.setTransaction(options.currentTransaction, options, true, true);
                    }

                    self.mode = 'view';
                    self.editId = options.id;

                    promises.push(self.transactionsStore.getTransaction({ transactionId: self.editId }));
                } else {
                    self.mode = 'add';
                    self.editId = null;

                    if (options.template) {
                        self.setTransaction(options.template, options, false, false);
                    }

                    if (self.settingsStore.appSettings.autoGetCurrentGeoLocation
                        && !self.geoLocationStatus && !self.transaction.geoLocation) {
                        self.updateGeoLocation(false);
                    }
                }
            } else if (self.type === 'template') {
                self.transaction.name = '';

                if (options && options.templateType) {
                    self.transaction.templateType = options.templateType;
                }

                if (self.transaction.templateType === templateConstants.allTemplateTypes.Schedule) {
                    self.transaction.scheduledFrequencyType = templateConstants.allTemplateScheduledFrequencyTypes.Disabled.type;
                    self.transaction.scheduledFrequency = '';
                }

                if (options && options.id) {
                    if (options.currentTemplate) {
                        self.setTransaction(options.currentTemplate, options, false, false);
                        self.transaction.templateType = options.currentTemplate.templateType;
                        self.transaction.name = options.currentTemplate.name;

                        if (self.transaction.templateType === templateConstants.allTemplateTypes.Schedule) {
                            self.transaction.scheduledFrequencyType = options.currentTemplate.scheduledFrequencyType;
                            self.transaction.scheduledFrequency = options.currentTemplate.scheduledFrequency;
                            self.transaction.utcOffset = options.currentTemplate.utcOffset;
                            self.transaction.timeZone = undefined;
                        }
                    }

                    self.mode = 'edit';
                    self.editId = options.id;
                    self.transaction.id = options.id;

                    promises.push(self.transactionTemplatesStore.getTemplate({ templateId: self.editId }));
                } else {
                    self.mode = 'add';
                    self.editId = null;
                    self.transaction.id = null;
                }
            }

            if (options.type && options.type !== '0' &&
                options.type >= self.allTransactionTypes.Income &&
                options.type <= self.allTransactionTypes.Transfer) {
                self.transaction.type = parseInt(options.type);
            }

            if (self.mode === 'add') {
                self.clientSessionId = generateRandomUUID();
            }

            Promise.all(promises).then(function (responses) {
                if (self.editId && !responses[3]) {
                    if (self.reject) {
                        if (self.type === 'transaction') {
                            self.reject('Unable to retrieve transaction');
                        } else if (self.type === 'template') {
                            self.reject('Unable to retrieve template');
                        }
                    }

                    return;
                }

                if (self.type === 'transaction' && options && options.id && responses[3]) {
                    const transaction = responses[3];
                    self.setTransaction(transaction, options, true, true);
                    self.originalTransactionEditable = transaction.editable;
                } else if (self.type === 'template' && options && options.id && responses[3]) {
                    const template = responses[3];
                    self.setTransaction(template, options, false, false);
                    self.transaction.templateType = template.templateType;
                    self.transaction.name = template.name;

                    if (self.transaction.templateType === templateConstants.allTemplateTypes.Schedule) {
                        self.transaction.scheduledFrequencyType = template.scheduledFrequencyType;
                        self.transaction.scheduledFrequency = template.scheduledFrequency;
                        self.transaction.utcOffset = template.utcOffset;
                        self.transaction.timeZone = undefined;
                    }
                } else {
                    self.setTransaction(null, options, true, true);
                }

                self.loading = false;
            }).catch(error => {
                logger.error('failed to load essential data for editing transaction', error);

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
        save() {
            const self = this;

            const problemMessage = self.inputEmptyProblemMessage;

            if (problemMessage) {
                self.$refs.snackbar.showMessage(problemMessage);
                return;
            }

            if (self.type === 'transaction' && (self.mode === 'add' || self.mode === 'edit')) {
                const doSubmit = function () {
                    self.submitting = true;

                    self.transactionsStore.saveTransaction({
                        transaction: self.transaction,
                        defaultCurrency: self.defaultCurrency,
                        isEdit: self.mode === 'edit',
                        clientSessionId: self.clientSessionId
                    }).then(() => {
                        self.submitting = false;

                        if (self.resolve) {
                            if (self.mode === 'add') {
                                self.resolve({
                                    message: 'You have added a new transaction'
                                });
                            } else if (self.mode === 'edit') {
                                self.resolve({
                                    message: 'You have saved this transaction'
                                });
                            }
                        }

                        self.showState = false;
                    }).catch(error => {
                        self.submitting = false;

                        if (!error.processed) {
                            self.$refs.snackbar.showError(error);
                        }
                    });
                };

                if (self.transaction.sourceAmount === 0) {
                    self.$refs.confirmDialog.open('Are you sure you want to save this transaction with a zero amount?').then(() => {
                        doSubmit();
                    });
                } else {
                    doSubmit();
                }
            } else if (self.type === 'template' && (self.mode === 'add' || self.mode === 'edit')) {
                self.submitting = true;

                self.transactionTemplatesStore.saveTemplateContent({
                    template: self.transaction,
                    isEdit: self.mode === 'edit',
                    clientSessionId: self.clientSessionId
                }).then(() => {
                    self.submitting = false;

                    if (self.resolve) {
                        if (self.mode === 'add') {
                            self.resolve({
                                message: 'You have added a new template'
                            });
                        } else if (self.mode === 'edit') {
                            self.resolve({
                                message: 'You have saved this template'
                            });
                        }
                    }

                    self.showState = false;
                }).catch(error => {
                    self.submitting = false;

                    if (!error.processed) {
                        self.$refs.snackbar.showError(error);
                    }
                });
            }
        },
        duplicate() {
            if (this.type !== 'transaction' || this.mode !== 'view') {
                return;
            }

            this.editId = null;
            this.transaction.id = null;
            this.transaction.time = getCurrentUnixTime();
            this.transaction.timeZone = this.settingsStore.appSettings.timeZone;
            this.transaction.utcOffset = getTimezoneOffsetMinutes(this.transaction.timeZone);
            this.transaction.geoLocation = null;
            this.transaction.pictures = [];
            this.mode = 'add';
        },
        edit() {
            if (this.type !== 'transaction' || this.mode !== 'view') {
                return;
            }

            this.mode = 'edit';
        },
        remove() {
            const self = this;

            if (this.type !== 'transaction' || self.mode !== 'view') {
                return;
            }

            self.$refs.confirmDialog.open('Are you sure you want to delete this transaction?').then(() => {
                self.submitting = true;

                self.transactionsStore.deleteTransaction({
                    transaction: self.transaction,
                    defaultCurrency: self.defaultCurrency
                }).then(() => {
                    if (self.resolve) {
                        self.resolve();
                    }

                    self.submitting = false;
                    self.showState = false;
                }).catch(error => {
                    self.submitting = false;

                    if (!error.processed) {
                        self.$refs.snackbar.showError(error);
                    }
                });
            });
        },
        cancel() {
            if (this.reject) {
                this.reject();
            }

            this.showState = false;
        },
        showDateTimeError(error) {
            this.$refs.snackbar.showError(error);
        },
        updateGeoLocation(forceUpdate) {
            const self = this;
            self.geoMenuState = false;

            if (!self.isSupportGeoLocation) {
                logger.warn('this browser does not support geo location');

                if (forceUpdate) {
                    self.$refs.snackbar.showMessage('Unable to retrieve current position');
                }
                return;
            }

            navigator.geolocation.getCurrentPosition(function (position) {
                if (!position || !position.coords) {
                    logger.error('current position is null');
                    self.geoLocationStatus = 'error';

                    if (forceUpdate) {
                        self.$refs.snackbar.showMessage('Unable to retrieve current position');
                    }

                    return;
                }

                self.geoLocationStatus = 'success';

                self.transaction.geoLocation = {
                    latitude: position.coords.latitude,
                    longitude: position.coords.longitude
                };
            }, function (err) {
                logger.error('cannot retrieve current position', err);
                self.geoLocationStatus = 'error';

                if (forceUpdate) {
                    self.$refs.snackbar.showMessage('Unable to retrieve current position');
                }
            });

            self.geoLocationStatus = 'getting';
        },
        clearGeoLocation() {
            this.geoMenuState = false;
            this.geoLocationStatus = null;
            this.transaction.geoLocation = null;
        },
        swapTransactionData(swapAccount, swapAmount) {
            if (swapAccount) {
                const oldSourceAccountId = this.transaction.sourceAccountId;
                this.transaction.sourceAccountId = this.transaction.destinationAccountId;
                this.transaction.destinationAccountId = oldSourceAccountId;
            }

            if (swapAmount) {
                const oldSourceAmount = this.transaction.sourceAmount;
                this.transaction.sourceAmount = this.transaction.destinationAmount;
                this.transaction.destinationAmount = oldSourceAmount;
            }
        },
        showOpenPictureDialog() {
            if (!this.canAddTransactionPicture || this.submitting) {
                return;
            }

            this.$refs.pictureInput.click();
        },
        uploadPicture(event) {
            if (!event || !event.target || !event.target.files || !event.target.files.length) {
                return;
            }

            const self = this;
            const pictureFile = event.target.files[0];

            event.target.value = null;

            self.uploadingPicture = true;
            self.submitting = true;

            self.transactionsStore.uploadTransactionPicture({ pictureFile }).then(response => {
                if (!isArray(self.transaction.pictures)) {
                    self.transaction.pictures = [];
                }

                self.transaction.pictures.push(response);

                self.uploadingPicture = false;
                self.submitting = false;
            }).catch(error => {
                self.uploadingPicture = false;
                self.submitting = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        viewOrRemovePicture(pictureInfo) {
            if (this.mode !== 'add' && this.mode !== 'edit') {
                window.open(this.getTransactionPictureUrl(pictureInfo), '_blank');
                return;
            }

            const self = this;

            self.$refs.confirmDialog.open('Are you sure you want to remove this transaction picture?').then(() => {
                self.removingPictureId = pictureInfo.pictureId;
                self.submitting = true;

                self.transactionsStore.removeUnusedTransactionPicture({ pictureInfo }).then(response => {
                    if (response && isArray(self.transaction.pictures)) {
                        for (let i = 0; i < self.transaction.pictures.length; i++) {
                            if (self.transaction.pictures[i].pictureId === pictureInfo.pictureId) {
                                self.transaction.pictures.splice(i, 1);
                            }
                        }
                    }

                    self.removingPictureId = '';
                    self.submitting = false;
                }).catch(error => {
                    if (error.error && error.error.errorCode === 211001) {
                        for (let i = 0; i < self.transaction.pictures.length; i++) {
                            if (self.transaction.pictures[i].pictureId === pictureInfo.pictureId) {
                                self.transaction.pictures.splice(i, 1);
                            }
                        }
                    } else if (!error.processed) {
                        self.$refs.snackbar.showError(error);
                    }

                    self.removingPictureId = '';
                    self.submitting = false;
                });
            });
        },
        getTransactionPictureUrl(pictureInfo) {
            return this.transactionsStore.getTransactionPictureUrl(pictureInfo);
        },
        getPrimaryCategoryName(categoryId, allCategories) {
            return getTransactionPrimaryCategoryName(categoryId, allCategories);
        },
        getSecondaryCategoryName(categoryId, allCategories) {
            return getTransactionSecondaryCategoryName(categoryId, allCategories);
        },
        setTransaction(transaction, options, setContextData, convertContextTime) {
            setTransactionModelByTransaction(
                this.transaction,
                transaction,
                this.allCategories,
                this.allCategoriesMap,
                this.allVisibleAccounts,
                this.allAccountsMap,
                this.allTagsMap,
                this.defaultAccountId,
                {
                    type: options.type,
                    categoryId: options.categoryId,
                    accountId: options.accountId,
                    destinationAccountId: options.destinationAccountId,
                    amount: options.amount,
                    destinationAmount: options.destinationAmount,
                    tagIds: options.tagIds,
                    comment: options.comment
                },
                setContextData,
                convertContextTime
            );
        }
    }
}
</script>

<style>
.transaction-edit-amount .v-field__field > input {
    font-size: 1.25rem;
}

.transaction-edit-timezone.v-input input::placeholder {
    color: rgba(var(--v-theme-on-background), var(--v-high-emphasis-opacity)) !important;
}

.transaction-edit-map-view {
    height: 220px;
}

@media (min-height: 630px) {
    .transaction-edit-map-view {
        height: 300px;
    }

    @media (min-width: 960px) {
        .transaction-pictures {
            min-height: 300px;
        }
    }
}

@media (min-height: 700px) {
    .transaction-edit-map-view {
        height: 350px;
    }

    @media (min-width: 960px) {
        .transaction-pictures {
            min-height: 350px;
        }
    }
}

@media (min-height: 800px) {
    .transaction-edit-map-view {
        height: 450px;
    }

    @media (min-width: 960px) {
        .transaction-pictures {
            min-height: 450px;
        }
    }
}

@media (min-height: 900px) {
    .transaction-edit-map-view {
        height: 550px;
    }

    @media (min-width: 960px) {
        .transaction-pictures {
            min-height: 550px;
        }
    }
}

.transaction-picture .picture-control-icon {
    display: none;
    position: absolute;
    width: 100% !important;
    height: 100% !important;
    background-color: rgba(0, 0, 0, 0.4);
}

.transaction-picture .picture-control-icon > i.v-icon {
    background-color: transparent;
    color: rgba(255, 255, 255, 0.8);
}

.transaction-picture:hover .picture-control-icon,
.transaction-picture .picture-control-icon.show-control-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    vertical-align: middle;
}

.transaction-picture:hover .transaction-picture-placeholder {
    display: none;
}

.transaction-picture-add {
    border: 2px dashed rgba(var(--v-theme-grey-500));

    .transaction-picture-add-icon {
        color: rgba(var(--v-theme-grey-500));
    }
}

.transaction-picture-add.enabled:hover {
    border: 2px dashed rgba(var(--v-theme-grey-700));

    .transaction-picture-add-icon {
        color: rgba(var(--v-theme-grey-700));
    }
}
</style>
