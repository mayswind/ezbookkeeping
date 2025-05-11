<template>
    <v-dialog width="1000" :persistent="isTransactionModified" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-8">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <div class="d-flex w-100 align-center justify-center">
                        <h4 class="text-h4">{{ tt(title) }}</h4>
                        <v-progress-circular indeterminate size="22" class="ml-2" v-if="loading"></v-progress-circular>
                    </div>
                    <v-btn density="comfortable" color="default" variant="text" class="ml-2" :icon="true"
                           :disabled="loading || submitting" v-if="mode !== TransactionEditPageMode.View && (activeTab === 'basicInfo' || (activeTab === 'map' && isSupportGetGeoLocationByClick()))">
                        <v-icon :icon="mdiDotsVertical" />
                        <v-menu activator="parent">
                            <v-list v-if="activeTab === 'basicInfo'">
                                <v-list-item :prepend-icon="mdiSwapHorizontal"
                                             :title="tt('Swap Account')"
                                             v-if="transaction.type === TransactionType.Transfer"
                                             @click="swapTransactionData(true, false)"></v-list-item>
                                <v-list-item :prepend-icon="mdiSwapHorizontal"
                                             :title="tt('Swap Amount')"
                                             v-if="transaction.type === TransactionType.Transfer"
                                             @click="swapTransactionData(false, true)"></v-list-item>
                                <v-list-item :prepend-icon="mdiSwapHorizontal"
                                             :title="tt('Swap Account and Amount')"
                                             v-if="transaction.type === TransactionType.Transfer"
                                             @click="swapTransactionData(true, true)"></v-list-item>
                                <v-divider v-if="transaction.type === TransactionType.Transfer" />
                                <v-list-item :prepend-icon="mdiEyeOutline"
                                             :title="tt('Show Amount')"
                                             v-if="transaction.hideAmount" @click="transaction.hideAmount = false"></v-list-item>
                                <v-list-item :prepend-icon="mdiEyeOffOutline"
                                             :title="tt('Hide Amount')"
                                             v-if="!transaction.hideAmount" @click="transaction.hideAmount = true"></v-list-item>
                            </v-list>
                            <v-list v-if="activeTab === 'map'">
                                <v-list-item key="setGeoLocationByClickMap" value="setGeoLocationByClickMap"
                                             :prepend-icon="mdiMapMarkerOutline"
                                             :disabled="!transaction.geoLocation" v-if="isSupportGetGeoLocationByClick()">
                                    <v-list-item-title class="cursor-pointer" @click="setGeoLocationByClickMap = !setGeoLocationByClickMap; geoMenuState = false">
                                        <div class="d-flex align-center">
                                            <span>{{ tt('Click on Map to Set Geographic Location') }}</span>
                                            <v-spacer/>
                                            <v-icon :icon="mdiCheck" v-if="setGeoLocationByClickMap" />
                                        </div>
                                    </v-list-item-title>
                                </v-list-item>
                            </v-list>
                        </v-menu>
                    </v-btn>
                </div>
            </template>
            <v-card-text class="d-flex flex-column flex-md-row mt-md-4 pt-0">
                <div class="mb-4">
                    <v-tabs class="v-tabs-pill" direction="vertical" :class="{ 'readonly': type === TransactionEditPageType.Transaction && mode !== TransactionEditPageMode.Add }"
                            :disabled="loading || submitting" v-model="transaction.type">
                        <v-tab :value="TransactionType.Expense" :disabled="type === TransactionEditPageType.Transaction && mode !== TransactionEditPageMode.Add && transaction.type !== TransactionType.Expense" v-if="transaction.type !== TransactionType.ModifyBalance">
                            <span>{{ tt('Expense') }}</span>
                        </v-tab>
                        <v-tab :value="TransactionType.Income" :disabled="type === TransactionEditPageType.Transaction && mode !== TransactionEditPageMode.Add && transaction.type !== TransactionType.Income" v-if="transaction.type !== TransactionType.ModifyBalance">
                            <span>{{ tt('Income') }}</span>
                        </v-tab>
                        <v-tab :value="TransactionType.Transfer" :disabled="type === TransactionEditPageType.Transaction && mode !== TransactionEditPageMode.Add && transaction.type !== TransactionType.Transfer" v-if="transaction.type !== TransactionType.ModifyBalance">
                            <span>{{ tt('Transfer') }}</span>
                        </v-tab>
                        <v-tab :value="TransactionType.ModifyBalance" v-if="type === TransactionEditPageType.Transaction && transaction.type === TransactionType.ModifyBalance">
                            <span>{{ tt('Modify Balance') }}</span>
                        </v-tab>
                    </v-tabs>
                    <v-divider class="my-2"/>
                    <v-tabs direction="vertical" :disabled="loading || submitting" v-model="activeTab">
                        <v-tab value="basicInfo">
                            <span>{{ tt('Basic Information') }}</span>
                        </v-tab>
                        <v-tab value="map" :disabled="!transaction.geoLocation" v-if="type === TransactionEditPageType.Transaction && !!getMapProvider()">
                            <span>{{ tt('Location on Map') }}</span>
                        </v-tab>
                        <v-tab value="pictures" :disabled="mode !== TransactionEditPageMode.Add && mode !== TransactionEditPageMode.Edit && (!transaction.pictures || !transaction.pictures.length)" v-if="type === TransactionEditPageType.Transaction && isTransactionPicturesEnabled()">
                            <span>{{ tt('Pictures') }}</span>
                        </v-tab>
                    </v-tabs>
                </div>

                <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container ml-md-5"
                          v-model="activeTab">
                    <v-window-item value="basicInfo">
                        <v-form class="mt-2">
                            <v-row>
                                <v-col cols="12" v-if="type === TransactionEditPageType.Template && transaction instanceof TransactionTemplate">
                                    <v-text-field
                                        type="text"
                                        persistent-placeholder
                                        :disabled="loading || submitting"
                                        :label="tt('Template Name')"
                                        :placeholder="tt('Template Name')"
                                        v-model="transaction.name"
                                    />
                                </v-col>
                                <v-col cols="12" :md="transaction.type === TransactionType.Transfer ? 6 : 12">
                                    <amount-input class="transaction-edit-amount font-weight-bold"
                                                  :color="sourceAmountColor"
                                                  :currency="sourceAccountCurrency"
                                                  :show-currency="true"
                                                  :readonly="mode === TransactionEditPageMode.View"
                                                  :disabled="loading || submitting"
                                                  :persistent-placeholder="true"
                                                  :hide="transaction.hideAmount"
                                                  :label="sourceAmountTitle"
                                                  :placeholder="tt(sourceAmountName)"
                                                  :enable-formula="mode !== TransactionEditPageMode.View"
                                                  v-model="transaction.sourceAmount"/>
                                </v-col>
                                <v-col cols="12" :md="6" v-if="transaction.type === TransactionType.Transfer">
                                    <amount-input class="transaction-edit-amount font-weight-bold" color="primary"
                                                  :currency="destinationAccountCurrency"
                                                  :show-currency="true"
                                                  :readonly="mode === TransactionEditPageMode.View"
                                                  :disabled="loading || submitting"
                                                  :persistent-placeholder="true"
                                                  :hide="transaction.hideAmount"
                                                  :label="transferInAmountTitle"
                                                  :placeholder="tt('Transfer In Amount')"
                                                  :enable-formula="mode !== TransactionEditPageMode.View"
                                                  v-model="transaction.destinationAmount"/>
                                </v-col>
                                <v-col cols="12" md="12" v-if="transaction.type === TransactionType.Expense">
                                    <two-column-select primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                       primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                       primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                       secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                       secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                       secondary-hidden-field="hidden"
                                                       :readonly="mode === TransactionEditPageMode.View"
                                                       :disabled="loading || submitting || !hasAvailableExpenseCategories"
                                                       :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                                       :show-selection-primary-text="true"
                                                       :custom-selection-primary-text="getTransactionPrimaryCategoryName(transaction.expenseCategoryId, allCategories[CategoryType.Expense])"
                                                       :custom-selection-secondary-text="getTransactionSecondaryCategoryName(transaction.expenseCategoryId, allCategories[CategoryType.Expense])"
                                                       :label="tt('Category')" :placeholder="tt('Category')"
                                                       :items="allCategories[CategoryType.Expense] || []"
                                                       v-model="transaction.expenseCategoryId">
                                    </two-column-select>
                                </v-col>
                                <v-col cols="12" md="12" v-if="transaction.type === TransactionType.Income">
                                    <two-column-select primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                       primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                       primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                       secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                       secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                       secondary-hidden-field="hidden"
                                                       :readonly="mode === TransactionEditPageMode.View"
                                                       :disabled="loading || submitting || !hasAvailableIncomeCategories"
                                                       :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                                       :show-selection-primary-text="true"
                                                       :custom-selection-primary-text="getTransactionPrimaryCategoryName(transaction.incomeCategoryId, allCategories[CategoryType.Income])"
                                                       :custom-selection-secondary-text="getTransactionSecondaryCategoryName(transaction.incomeCategoryId, allCategories[CategoryType.Income])"
                                                       :label="tt('Category')" :placeholder="tt('Category')"
                                                       :items="allCategories[CategoryType.Income] || []"
                                                       v-model="transaction.incomeCategoryId">
                                    </two-column-select>
                                </v-col>
                                <v-col cols="12" md="12" v-if="transaction.type === TransactionType.Transfer">
                                    <two-column-select primary-key-field="id" primary-value-field="id" primary-title-field="name"
                                                       primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                                       primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                                       secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                                       secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                                       secondary-hidden-field="hidden"
                                                       :readonly="mode === TransactionEditPageMode.View"
                                                       :disabled="loading || submitting || !hasAvailableTransferCategories"
                                                       :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                                       :show-selection-primary-text="true"
                                                       :custom-selection-primary-text="getTransactionPrimaryCategoryName(transaction.transferCategoryId, allCategories[CategoryType.Transfer])"
                                                       :custom-selection-secondary-text="getTransactionSecondaryCategoryName(transaction.transferCategoryId, allCategories[CategoryType.Transfer])"
                                                       :label="tt('Category')" :placeholder="tt('Category')"
                                                       :items="allCategories[CategoryType.Transfer] || []"
                                                       v-model="transaction.transferCategoryId">
                                    </two-column-select>
                                </v-col>
                                <v-col cols="12" :md="transaction.type === TransactionType.Transfer ? 6 : 12">
                                    <two-column-select primary-key-field="id" primary-value-field="category"
                                                       primary-title-field="name" primary-footer-field="displayBalance"
                                                       primary-icon-field="icon" primary-icon-type="account"
                                                       primary-sub-items-field="accounts"
                                                       :primary-title-i18n="true"
                                                       secondary-key-field="id" secondary-value-field="id"
                                                       secondary-title-field="name" secondary-footer-field="displayBalance"
                                                       secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                       :readonly="mode === TransactionEditPageMode.View"
                                                       :disabled="loading || submitting || !allVisibleAccounts.length"
                                                       :enable-filter="true" :filter-placeholder="tt('Find account')" :filter-no-items-text="tt('No available account')"
                                                       :custom-selection-primary-text="sourceAccountName"
                                                       :label="tt(sourceAccountTitle)"
                                                       :placeholder="tt(sourceAccountTitle)"
                                                       :items="allVisibleCategorizedAccounts"
                                                       v-model="transaction.sourceAccountId">
                                    </two-column-select>
                                </v-col>
                                <v-col cols="12" md="6" v-if="transaction.type === TransactionType.Transfer">
                                    <two-column-select primary-key-field="id" primary-value-field="category"
                                                       primary-title-field="name" primary-footer-field="displayBalance"
                                                       primary-icon-field="icon" primary-icon-type="account"
                                                       primary-sub-items-field="accounts"
                                                       :primary-title-i18n="true"
                                                       secondary-key-field="id" secondary-value-field="id"
                                                       secondary-title-field="name" secondary-footer-field="displayBalance"
                                                       secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                       :readonly="mode === TransactionEditPageMode.View"
                                                       :disabled="loading || submitting || !allVisibleAccounts.length"
                                                       :enable-filter="true" :filter-placeholder="tt('Find account')" :filter-no-items-text="tt('No available account')"
                                                       :custom-selection-primary-text="destinationAccountName"
                                                       :label="tt('Destination Account')"
                                                       :placeholder="tt('Destination Account')"
                                                       :items="allVisibleCategorizedAccounts"
                                                       v-model="transaction.destinationAccountId">
                                    </two-column-select>
                                </v-col>
                                <v-col cols="12" md="6" v-if="type === TransactionEditPageType.Transaction">
                                    <date-time-select
                                        :readonly="mode === TransactionEditPageMode.View"
                                        :disabled="loading || submitting"
                                        :label="tt('Transaction Time')"
                                        v-model="transaction.time"
                                        @error="onShowDateTimeError" />
                                </v-col>
                                <v-col cols="12" md="6" v-if="type === TransactionEditPageType.Template && transaction instanceof TransactionTemplate && transaction.templateType === TemplateType.Schedule.type">
                                    <schedule-frequency-select
                                        :readonly="mode === TransactionEditPageMode.View"
                                        :disabled="loading || submitting"
                                        :label="tt('Scheduled Transaction Frequency')"
                                        v-model:type="transaction.scheduledFrequencyType"
                                        v-model="transaction.scheduledFrequency" />
                                </v-col>
                                <v-col cols="12" md="6" v-if="type === TransactionEditPageType.Transaction || (type === TransactionEditPageType.Template && transaction instanceof TransactionTemplate && transaction.templateType === TemplateType.Schedule.type)">
                                    <v-autocomplete
                                        class="transaction-edit-timezone"
                                        item-title="displayNameWithUtcOffset"
                                        item-value="name"
                                        auto-select-first
                                        persistent-placeholder
                                        :readonly="mode === TransactionEditPageMode.View"
                                        :disabled="loading || submitting"
                                        :label="tt('Transaction Timezone')"
                                        :placeholder="!transaction.timeZone && transaction.timeZone !== '' ? `(${transactionDisplayTimezone}) ${transactionTimezoneTimeDifference}` : tt('Timezone')"
                                        :items="allTimezones"
                                        :no-data-text="tt('No results')"
                                        v-model="transaction.timeZone"
                                    >
                                        <template #selection="{ item }">
                                            <span class="text-truncate" v-if="transaction.timeZone || transaction.timeZone === ''">
                                                {{ item.title }}
                                            </span>
                                        </template>
                                    </v-autocomplete>
                                </v-col>
                                <v-col cols="12" md="6" v-if="type === TransactionEditPageType.Template && transaction instanceof TransactionTemplate && transaction.templateType === TemplateType.Schedule.type">
                                    <date-select
                                        :readonly="mode === TransactionEditPageMode.View"
                                        :disabled="loading || submitting"
                                        :clearable="true"
                                        :label="tt('Start Date')"
                                        v-model="transaction.scheduledStartDate" />
                                </v-col>
                                <v-col cols="12" md="6" v-if="type === TransactionEditPageType.Template && transaction instanceof TransactionTemplate && transaction.templateType === TemplateType.Schedule.type">
                                    <date-select
                                        :readonly="mode === TransactionEditPageMode.View"
                                        :disabled="loading || submitting"
                                        :clearable="true"
                                        :label="tt('End Date')"
                                        v-model="transaction.scheduledEndDate" />
                                </v-col>
                                <v-col cols="12" md="12" v-if="type === TransactionEditPageType.Transaction">
                                    <v-select
                                        persistent-placeholder
                                        :readonly="mode === TransactionEditPageMode.View"
                                        :disabled="loading || submitting"
                                        :label="tt('Geographic Location')"
                                        v-model="transaction"
                                        v-model:menu="geoMenuState"
                                    >
                                        <template #selection>
                                            <span class="cursor-pointer" v-if="transaction.geoLocation">{{ `(${transaction.geoLocation.longitude}, ${transaction.geoLocation.latitude})` }}</span>
                                            <span class="cursor-pointer" v-else-if="!transaction.geoLocation">{{ geoLocationStatusInfo }}</span>
                                        </template>

                                        <template #no-data>
                                            <v-list class="py-0">
                                                <v-list-item v-if="mode !== TransactionEditPageMode.View" @click="updateGeoLocation(true)">{{ tt('Update Geographic Location') }}</v-list-item>
                                                <v-list-item v-if="mode !== TransactionEditPageMode.View" @click="clearGeoLocation">{{ tt('Clear Geographic Location') }}</v-list-item>
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
                                        :closable-chips="mode !== TransactionEditPageMode.View"
                                        :readonly="mode === TransactionEditPageMode.View"
                                        :disabled="loading || submitting"
                                        :label="tt('Tags')"
                                        :placeholder="tt('None')"
                                        :items="allTags"
                                        v-model="transaction.tagIds"
                                        v-model:search="tagSearchContent"
                                    >
                                        <template #chip="{ props, item }">
                                            <v-chip :prepend-icon="mdiPound" :text="item.title" v-bind="props"/>
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
                                            <v-list-item :disabled="true" v-bind="props"
                                                         v-if="item.raw.hidden && item.raw.name.toLowerCase().indexOf(tagSearchContent.toLowerCase()) >= 0 && isAllFilteredTagHidden">
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

                                        <template #no-data>
                                            <v-list class="py-0">
                                                <v-list-item v-if="tagSearchContent" @click="saveNewTag(tagSearchContent)">{{ tt('format.misc.addNewTag', { tag: tagSearchContent }) }}</v-list-item>
                                                <v-list-item v-else-if="!tagSearchContent">{{ tt('No available tag') }}</v-list-item>
                                            </v-list>
                                        </template>
                                    </v-autocomplete>
                                </v-col>
                                <v-col cols="12" md="12">
                                    <v-textarea
                                        type="text"
                                        persistent-placeholder
                                        rows="3"
                                        :readonly="mode === TransactionEditPageMode.View"
                                        :disabled="loading || submitting"
                                        :label="tt('Description')"
                                        :placeholder="tt('Your transaction description (optional)')"
                                        v-model="transaction.comment"
                                    />
                                </v-col>
                            </v-row>
                        </v-form>
                    </v-window-item>
                    <v-window-item value="map">
                        <v-row>
                            <v-col cols="12" md="12">
                                <map-view ref="map" map-class="transaction-edit-map-view" :geo-location="transaction.geoLocation" @click="updateSpecifiedGeoLocation">
                                    <template #error-title="{ mapSupported, mapDependencyLoaded }">
                                        <span class="text-subtitle-1" v-if="!mapSupported"><b>{{ tt('Unsupported Map Provider') }}</b></span>
                                        <span class="text-subtitle-1" v-else-if="!mapDependencyLoaded"><b>{{ tt('Cannot Initialize Map') }}</b></span>
                                    </template>
                                    <template #error-content>
                                        <p class="text-body-1">
                                            {{ tt('Please refresh the page and try again. If the error persists, ensure that the server\'s map settings are correctly configured.') }}
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
                                        <template #error>
                                            <div class="d-flex align-center justify-center fill-height bg-light-primary">
                                                <span class="text-body-1">{{ tt('Failed to load image, please check whether the config "domain" and "root_url" are set correctly.') }}</span>
                                            </div>
                                        </template>
                                    </v-img>
                                    <div class="picture-control-icon" :class="{ 'show-control-icon': pictureInfo.pictureId === removingPictureId }">
                                        <v-icon size="64" :icon="mdiTrashCanOutline" v-if="(mode === TransactionEditPageMode.Add || mode === TransactionEditPageMode.Edit) && pictureInfo.pictureId !== removingPictureId"/>
                                        <v-progress-circular color="grey-500" indeterminate size="48" v-if="(mode === TransactionEditPageMode.Add || mode === TransactionEditPageMode.Edit) && pictureInfo.pictureId === removingPictureId"></v-progress-circular>
                                        <v-icon size="64" :icon="mdiFullscreen" v-if="mode !== TransactionEditPageMode.Add && mode !== TransactionEditPageMode.Edit"/>
                                    </div>
                                </v-avatar>
                            </v-col>
                            <v-col cols="6" md="3" v-if="canAddTransactionPicture">
                                <v-avatar rounded="lg" variant="tonal" size="160"
                                          class="transaction-picture transaction-picture-add"
                                          :class="{ 'enabled': !submitting, 'cursor-pointer': !submitting }"
                                          color="rgba(0,0,0,0)" @click="showOpenPictureDialog">
                                    <v-tooltip activator="parent" v-if="!submitting">{{ tt('Add Picture') }}</v-tooltip>
                                    <v-icon class="transaction-picture-add-icon" size="56" :icon="mdiImagePlusOutline" v-if="!uploadingPicture"/>
                                    <v-progress-circular color="grey-500" indeterminate size="48" v-if="uploadingPicture"></v-progress-circular>
                                </v-avatar>
                            </v-col>
                        </v-row>
                    </v-window-item>
                </v-window>
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center flex-wrap mt-2 mt-sm-4 mt-md-6 gap-4">
                    <v-btn :disabled="inputIsEmpty || loading || submitting" v-if="mode !== TransactionEditPageMode.View" @click="save">
                        {{ tt(saveButtonTitle) }}
                        <v-progress-circular indeterminate size="22" class="ml-2" v-if="submitting"></v-progress-circular>
                    </v-btn>
                    <v-btn-group variant="tonal" density="comfortable"
                                 v-if="mode === TransactionEditPageMode.View && transaction.type !== TransactionType.ModifyBalance">
                        <v-btn :disabled="loading || submitting"
                               @click="duplicate(false, false)">{{ tt('Duplicate') }}</v-btn>
                        <v-btn density="compact" :disabled="loading || submitting" :icon="true">
                            <v-icon :icon="mdiMenuDown" size="24" />
                            <v-menu activator="parent">
                                <v-list>
                                    <v-list-item :title="tt('Duplicate (With Time)')"
                                                 @click="duplicate(true, false)"></v-list-item>
                                    <v-list-item :title="tt('Duplicate (With Geographic Location)')"
                                                 @click="duplicate(false, true)"
                                                 v-if="transaction.geoLocation"></v-list-item>
                                    <v-list-item :title="tt('Duplicate (With Time and Geographic Location)')"
                                                 @click="duplicate(true, true)"
                                                 v-if="transaction.geoLocation"></v-list-item>
                                </v-list>
                            </v-menu>
                        </v-btn>
                    </v-btn-group>
                    <v-btn color="warning" variant="tonal" :disabled="loading || submitting"
                           v-if="mode === TransactionEditPageMode.View && originalTransactionEditable && transaction.type !== TransactionType.ModifyBalance"
                           @click="edit">{{ tt('Edit') }}</v-btn>
                    <v-btn color="error" variant="tonal" :disabled="loading || submitting"
                           v-if="mode === TransactionEditPageMode.View && originalTransactionEditable" @click="remove">
                        {{ tt('Delete') }}
                        <v-progress-circular indeterminate size="22" class="ml-2" v-if="submitting"></v-progress-circular>
                    </v-btn>
                    <v-btn color="secondary" variant="tonal" :disabled="loading || submitting"
                           @click="cancel">{{ tt(cancelButtonTitle) }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
    <input ref="pictureInput" type="file" style="display: none" :accept="SUPPORTED_IMAGE_EXTENSIONS" @change="uploadPicture($event)" />
</template>

<script setup lang="ts">
import MapView from '@/components/common/MapView.vue';
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef, watch, nextTick } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import {
    TransactionEditPageMode,
    TransactionEditPageType,
    GeoLocationStatus,
    useTransactionEditPageBase
} from '@/views/base/transactions/TransactionEditPageBase.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useTransactionTemplatesStore } from '@/stores/transactionTemplate.ts';

import type { MapPosition } from '@/core/map.ts';
import { CategoryType } from '@/core/category.ts';
import { TransactionType, TransactionEditScopeType } from '@/core/transaction.ts';
import { TemplateType, ScheduledTemplateFrequencyType } from '@/core/template.ts';
import { KnownErrorCode } from '@/consts/api.ts';
import { SUPPORTED_IMAGE_EXTENSIONS } from '@/consts/file.ts';

import { TransactionTag } from '@/models/transaction_tag.ts';
import { TransactionTemplate } from '@/models/transaction_template.ts';
import type { TransactionPictureInfoBasicResponse } from '@/models/transaction_picture_info.ts';
import { Transaction } from '@/models/transaction.ts';

import {
    getTimezoneOffsetMinutes,
    getCurrentUnixTime
} from '@/lib/datetime.ts';
import { generateRandomUUID } from '@/lib/misc.ts';
import {
    getTransactionPrimaryCategoryName,
    getTransactionSecondaryCategoryName
} from '@/lib/category.ts';
import { type SetTransactionOptions, setTransactionModelByTransaction } from '@/lib/transaction.ts';
import {
    isTransactionPicturesEnabled,
    getMapProvider
} from '@/lib/server_settings.ts';
import {
    isSupportGetGeoLocationByClick
} from '@/lib/map/index.ts';
import logger from '@/lib/logger.ts';

import {
    mdiDotsVertical,
    mdiEyeOffOutline,
    mdiEyeOutline,
    mdiSwapHorizontal,
    mdiMapMarkerOutline,
    mdiCheck,
    mdiPound,
    mdiMenuDown,
    mdiImagePlusOutline,
    mdiTrashCanOutline,
    mdiFullscreen
} from '@mdi/js';

export interface TransactionEditOptions extends SetTransactionOptions {
    id?: string;
    templateType?: number;
    template?: TransactionTemplate;
    currentTransaction?: Transaction;
    currentTemplate?: TransactionTemplate;
}

interface TransactionEditResponse {
    message: string;
}

type MapViewType = InstanceType<typeof MapView>;
type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;

const props = defineProps<{
    type: TransactionEditPageType;
    persistent?: boolean;
    show?: boolean;
}>();

const { tt } = useI18n();

const {
    mode,
    isSupportGeoLocation,
    editId,
    addByTemplateId,
    duplicateFromId,
    clientSessionId,
    loading,
    submitting,
    uploadingPicture,
    geoLocationStatus,
    setGeoLocationByClickMap,
    transaction,
    defaultCurrency,
    defaultAccountId,
    allTimezones,
    allVisibleAccounts,
    allAccountsMap,
    allVisibleCategorizedAccounts,
    allCategories,
    allCategoriesMap,
    allTags,
    allTagsMap,
    hasAvailableExpenseCategories,
    hasAvailableIncomeCategories,
    hasAvailableTransferCategories,
    canAddTransactionPicture,
    title,
    saveButtonTitle,
    cancelButtonTitle,
    sourceAmountName,
    sourceAmountTitle,
    sourceAccountTitle,
    transferInAmountTitle,
    sourceAccountName,
    destinationAccountName,
    sourceAccountCurrency,
    destinationAccountCurrency,
    transactionDisplayTimezone,
    transactionTimezoneTimeDifference,
    geoLocationStatusInfo,
    inputEmptyProblemMessage,
    inputIsEmpty,
    createNewTransactionModel,
    swapTransactionData,
    getTransactionPictureUrl
} = useTransactionEditPageBase(props.type);

const settingsStore = useSettingsStore();
const userStore = useUserStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();
const transactionsStore = useTransactionsStore();
const transactionTemplatesStore = useTransactionTemplatesStore();

const map = useTemplateRef<MapViewType>('map');
const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');
const pictureInput = useTemplateRef<HTMLInputElement>('pictureInput');

const showState = ref<boolean>(false);
const activeTab = ref<string>('basicInfo');
const originalTransactionEditable = ref<boolean>(false);
const geoMenuState = ref<boolean>(false);
const tagSearchContent = ref<string>('');
const removingPictureId = ref<string>('');

const initCategoryId = ref<string | undefined>(undefined);
const initAccountId = ref<string | undefined>(undefined);
const initTagIds = ref<string | undefined>(undefined);

let resolveFunc: ((response?: TransactionEditResponse) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const sourceAmountColor = computed<string | undefined>(() => {
    if (transaction.value.type === TransactionType.Expense) {
        return 'expense';
    } else if (transaction.value.type === TransactionType.Income) {
        return 'income';
    } else if (transaction.value.type === TransactionType.Transfer) {
        return 'primary';
    }

    return undefined;
});

const isAllFilteredTagHidden = computed<boolean>(() => {
    const lowerCaseTagSearchContent = tagSearchContent.value.toLowerCase();
    let hiddenCount = 0;

    for (const tag of allTags.value) {
        if (!lowerCaseTagSearchContent || tag.name.toLowerCase().indexOf(lowerCaseTagSearchContent) >= 0) {
            if (!tag.hidden) {
                return false;
            }

            hiddenCount++;
        }
    }

    return hiddenCount > 0;
});

const isTransactionModified = computed<boolean>(() => {
    if (mode.value === TransactionEditPageMode.Add) {
        return transactionsStore.isTransactionDraftModified(transaction.value, initCategoryId.value, initAccountId.value, initTagIds.value);
    } else if (mode.value === TransactionEditPageMode.Edit) {
        return true;
    } else {
        return false;
    }
});

function setTransaction(newTransaction: Transaction | null, options: SetTransactionOptions, setContextData: boolean, convertContextTime: boolean): void {
    setTransactionModelByTransaction(
        transaction.value,
        newTransaction,
        allCategories.value,
        allCategoriesMap.value,
        allVisibleAccounts.value,
        allAccountsMap.value,
        allTagsMap.value,
        defaultAccountId.value,
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

function open(options: TransactionEditOptions): Promise<TransactionEditResponse | undefined> {
    addByTemplateId.value = null;
    duplicateFromId.value = null;
    showState.value = true;
    activeTab.value = 'basicInfo';
    loading.value = true;
    submitting.value = false;
    geoLocationStatus.value = null;
    setGeoLocationByClickMap.value = false;
    originalTransactionEditable.value = false;

    initCategoryId.value = options.categoryId;
    initAccountId.value = options.accountId;
    initTagIds.value = options.tagIds;

    const newTransaction = createNewTransactionModel(options.type);
    setTransaction(newTransaction, options, true, false);

    const promises: Promise<unknown>[] = [
        accountsStore.loadAllAccounts({ force: false }),
        transactionCategoriesStore.loadAllCategories({ force: false }),
        transactionTagsStore.loadAllTags({ force: false })
    ];

    if (props.type === TransactionEditPageType.Transaction) {
        if (options && options.id) {
            if (options.currentTransaction) {
                setTransaction(options.currentTransaction, options, true, true);
            }

            mode.value = TransactionEditPageMode.View;
            editId.value = options.id;

            promises.push(transactionsStore.getTransaction({ transactionId: editId.value }));
        } else {
            mode.value = TransactionEditPageMode.Add;
            editId.value = null;

            if (options.template) {
                setTransaction(options.template, options, false, false);
                addByTemplateId.value = options.template.id;
            } else if ((settingsStore.appSettings.autoSaveTransactionDraft === 'enabled' || settingsStore.appSettings.autoSaveTransactionDraft === 'confirmation') && transactionsStore.transactionDraft) {
                setTransaction(Transaction.ofDraft(transactionsStore.transactionDraft), options, false, false);
            }

            if (settingsStore.appSettings.autoGetCurrentGeoLocation
                && !geoLocationStatus.value && !transaction.value.geoLocation) {
                updateGeoLocation(false);
            }
        }
    } else if (props.type === TransactionEditPageType.Template) {
        const template = TransactionTemplate.createNewTransactionTemplate(transaction.value);
        template.name = '';

        if (options && options.templateType) {
            template.templateType = options.templateType;
        }

        if (template.templateType === TemplateType.Schedule.type) {
            template.scheduledFrequencyType = ScheduledTemplateFrequencyType.Disabled.type;
            template.scheduledFrequency = '';
        }

        transaction.value = template;

        if (options && options.id) {
            if (options.currentTemplate) {
                setTransaction(options.currentTemplate, options, false, false);
                (transaction.value as TransactionTemplate).fillFrom(options.currentTemplate);
            }

            mode.value = TransactionEditPageMode.Edit;
            editId.value = options.id;
            transaction.value.id = options.id;

            promises.push(transactionTemplatesStore.getTemplate({ templateId: editId.value }));
        } else {
            mode.value = TransactionEditPageMode.Add;
            editId.value = null;
            transaction.value.id = '';
        }
    }

    if (options.type &&
        options.type >= TransactionType.Income &&
        options.type <= TransactionType.Transfer) {
        transaction.value.type = options.type;
    }

    if (mode.value === TransactionEditPageMode.Add) {
        clientSessionId.value = generateRandomUUID();
    }

    Promise.all(promises).then(function (responses) {
        if (editId.value && !responses[3]) {
            if (rejectFunc) {
                if (props.type === TransactionEditPageType.Transaction) {
                    rejectFunc('Unable to retrieve transaction');
                } else if (props.type === TransactionEditPageType.Template) {
                    rejectFunc('Unable to retrieve template');
                }
            }

            return;
        }

        if (props.type === TransactionEditPageType.Transaction && options && options.id && responses[3] && responses[3] instanceof Transaction) {
            const transaction: Transaction = responses[3];
            setTransaction(transaction, options, true, true);
            originalTransactionEditable.value = transaction.editable;
        } else if (props.type === TransactionEditPageType.Template && options && options.id && responses[3] && responses[3] instanceof TransactionTemplate) {
            const template: TransactionTemplate = responses[3];
            setTransaction(template, options, false, false);

            if (!(transaction.value instanceof TransactionTemplate)) {
                transaction.value = TransactionTemplate.createNewTransactionTemplate(transaction.value);
            }

            (transaction.value as TransactionTemplate).fillFrom(template);
        } else {
            setTransaction(null, options, true, true);
        }

        loading.value = false;
    }).catch(error => {
        logger.error('failed to load essential data for editing transaction', error);

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

function save(): void {
    const problemMessage = inputEmptyProblemMessage.value;

    if (problemMessage) {
        snackbar.value?.showMessage(problemMessage);
        return;
    }

    if (props.type === TransactionEditPageType.Transaction && (mode.value === TransactionEditPageMode.Add || mode.value === TransactionEditPageMode.Edit)) {
        const doSubmit = function () {
            submitting.value = true;

            transactionsStore.saveTransaction({
                transaction: transaction.value as Transaction,
                defaultCurrency: defaultCurrency.value,
                isEdit: mode.value === TransactionEditPageMode.Edit,
                clientSessionId: clientSessionId.value
            }).then(() => {
                submitting.value = false;

                if (resolveFunc) {
                    if (mode.value === TransactionEditPageMode.Add) {
                        resolveFunc({
                            message: 'You have added a new transaction'
                        });
                    } else if (mode.value === TransactionEditPageMode.Edit) {
                        resolveFunc({
                            message: 'You have saved this transaction'
                        });
                    }
                }

                if (mode.value === TransactionEditPageMode.Add && !addByTemplateId.value && !duplicateFromId.value) {
                    transactionsStore.clearTransactionDraft();
                }

                showState.value = false;
            }).catch(error => {
                submitting.value = false;

                if (error.error && (error.error.errorCode === KnownErrorCode.TransactionCannotCreateInThisTime || error.error.errorCode === KnownErrorCode.TransactionCannotModifyInThisTime)) {
                    confirmDialog.value?.open('You have set this time range to prevent editing transactions. Would you like to change the editable transaction range to All?').then(() => {
                        submitting.value = true;

                        userStore.updateUserTransactionEditScope({
                            transactionEditScope: TransactionEditScopeType.All.type
                        }).then(() => {
                            submitting.value = false;

                            snackbar.value?.showMessage('Your editable transaction range has been set to All');
                        }).catch(error => {
                            submitting.value = false;

                            if (!error.processed) {
                                snackbar.value?.showError(error);
                            }
                        });
                    });
                } else if (!error.processed) {
                    snackbar.value?.showError(error);
                }
            });
        };

        if (transaction.value.sourceAmount === 0) {
            confirmDialog.value?.open('Are you sure you want to save this transaction with a zero amount?').then(() => {
                doSubmit();
            });
        } else {
            doSubmit();
        }
    } else if (props.type === TransactionEditPageType.Template && (mode.value === TransactionEditPageMode.Add || mode.value === TransactionEditPageMode.Edit)) {
        submitting.value = true;

        transactionTemplatesStore.saveTemplateContent({
            template: transaction.value as TransactionTemplate,
            isEdit: mode.value === TransactionEditPageMode.Edit,
            clientSessionId: clientSessionId.value
        }).then(() => {
            submitting.value = false;

            if (resolveFunc) {
                if (mode.value === TransactionEditPageMode.Add) {
                    resolveFunc({
                        message: 'You have added a new template'
                    });
                } else if (mode.value === TransactionEditPageMode.Edit) {
                    resolveFunc({
                        message: 'You have saved this template'
                    });
                }
            }

            showState.value = false;
        }).catch(error => {
            submitting.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    }
}

function duplicate(withTime?: boolean, withGeoLocation?: boolean): void {
    if (props.type !== TransactionEditPageType.Transaction || mode.value !== TransactionEditPageMode.View) {
        return;
    }

    editId.value = null;
    duplicateFromId.value = transaction.value.id;
    clientSessionId.value = generateRandomUUID();
    activeTab.value = 'basicInfo';
    transaction.value.id = '';

    if (!withTime) {
        transaction.value.time = getCurrentUnixTime();
        transaction.value.timeZone = settingsStore.appSettings.timeZone;
        transaction.value.utcOffset = getTimezoneOffsetMinutes(transaction.value.timeZone);
    }

    if (!withGeoLocation) {
        transaction.value.removeGeoLocation();
    }

    transaction.value.clearPictures();
    mode.value = TransactionEditPageMode.Add;
}

function edit(): void {
    if (props.type !== TransactionEditPageType.Transaction || mode.value !== TransactionEditPageMode.View) {
        return;
    }

    mode.value = TransactionEditPageMode.Edit;
}

function remove(): void {
    if (props.type !== TransactionEditPageType.Transaction || mode.value !== TransactionEditPageMode.View) {
        return;
    }

    confirmDialog.value?.open('Are you sure you want to delete this transaction?').then(() => {
        submitting.value = true;

        transactionsStore.deleteTransaction({
            transaction: transaction.value as Transaction,
            defaultCurrency: defaultCurrency.value
        }).then(() => {
            if (resolveFunc) {
                resolveFunc();
            }

            submitting.value = false;
            showState.value = false;
        }).catch(error => {
            submitting.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    });
}

function cancel(): void {
    const doClose = function () {
        if (rejectFunc) {
            rejectFunc();
        }

        showState.value = false;
    };

    if (props.type !== TransactionEditPageType.Transaction || mode.value !== TransactionEditPageMode.Add || addByTemplateId.value || duplicateFromId.value) {
        doClose();
        return;
    }

    if (settingsStore.appSettings.autoSaveTransactionDraft === 'confirmation') {
        if (transactionsStore.isTransactionDraftModified(transaction.value, initCategoryId.value, initAccountId.value, initTagIds.value)) {
            confirmDialog.value?.open('Do you want to save this transaction draft?').then(() => {
                transactionsStore.saveTransactionDraft(transaction.value, initCategoryId.value, initAccountId.value, initTagIds.value);
                doClose();
            }).catch(() => {
                transactionsStore.clearTransactionDraft();
                doClose();
            });
        } else {
            transactionsStore.clearTransactionDraft();
            doClose();
        }
    } else if (settingsStore.appSettings.autoSaveTransactionDraft === 'enabled') {
        transactionsStore.saveTransactionDraft(transaction.value, initCategoryId.value, initAccountId.value, initTagIds.value);
        doClose();
    } else {
        doClose();
    }
}

function updateGeoLocation(forceUpdate: boolean): void {
    geoMenuState.value = false;

    if (!isSupportGeoLocation) {
        logger.warn('this browser does not support geo location');

        if (forceUpdate) {
            snackbar.value?.showMessage('Unable to retrieve current position');
        }
        return;
    }

    navigator.geolocation.getCurrentPosition(function (position) {
        if (!position || !position.coords) {
            logger.error('current position is null');
            geoLocationStatus.value = GeoLocationStatus.Error;

            if (forceUpdate) {
                snackbar.value?.showMessage('Unable to retrieve current position');
            }

            return;
        }

        geoLocationStatus.value = GeoLocationStatus.Success;

        transaction.value.setLatitudeAndLongitude(position.coords.latitude, position.coords.longitude);
    }, function (err) {
        logger.error('cannot retrieve current position', err);
        geoLocationStatus.value = GeoLocationStatus.Error;

        if (forceUpdate) {
            snackbar.value?.showMessage('Unable to retrieve current position');
        }
    });

    geoLocationStatus.value = GeoLocationStatus.Getting;
}

function updateSpecifiedGeoLocation(mapPosition: MapPosition): void {
    if (isSupportGetGeoLocationByClick() && setGeoLocationByClickMap.value) {
        transaction.value.setLatitudeAndLongitude(mapPosition.latitude, mapPosition.longitude);
        map.value?.setMarkerPosition(transaction.value.geoLocation);
    }
}

function clearGeoLocation(): void {
    geoMenuState.value = false;
    geoLocationStatus.value = null;
    transaction.value.removeGeoLocation();
}

function saveNewTag(tagName: string): void {
    submitting.value = true;

    transactionTagsStore.saveTag({
        tag: TransactionTag.createNewTag(tagName)
    }).then(tag => {
        submitting.value = false;

        if (tag && tag.id) {
            transaction.value.tagIds.push(tag.id);
        }
    }).catch(error => {
        submitting.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function showOpenPictureDialog(): void {
    if (!canAddTransactionPicture.value || submitting.value) {
        return;
    }

    pictureInput.value?.click();
}

function uploadPicture(event: Event): void {
    if (!event || !event.target) {
        return;
    }

    const el = event.target as HTMLInputElement;

    if (!el.files || !el.files.length) {
        return;
    }

    const pictureFile = el.files[0];

    el.value = '';

    uploadingPicture.value = true;
    submitting.value = true;

    transactionsStore.uploadTransactionPicture({ pictureFile }).then(response => {
        transaction.value.addPicture(response);
        uploadingPicture.value = false;
        submitting.value = false;
    }).catch(error => {
        uploadingPicture.value = false;
        submitting.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function viewOrRemovePicture(pictureInfo: TransactionPictureInfoBasicResponse): void {
    if (mode.value !== TransactionEditPageMode.Add && mode.value !== TransactionEditPageMode.Edit) {
        window.open(getTransactionPictureUrl(pictureInfo), '_blank');
        return;
    }

    confirmDialog.value?.open('Are you sure you want to remove this transaction picture?').then(() => {
        removingPictureId.value = pictureInfo.pictureId;
        submitting.value = true;

        transactionsStore.removeUnusedTransactionPicture({ pictureInfo }).then(response => {
            if (response) {
                transaction.value.removePicture(pictureInfo);
            }

            removingPictureId.value = '';
            submitting.value = false;
        }).catch(error => {
            if (error.error && error.error.errorCode === KnownErrorCode.TransactionPictureNotFound) {
                transaction.value.removePicture(pictureInfo);
            } else if (!error.processed) {
                snackbar.value?.showError(error);
            }

            removingPictureId.value = '';
            submitting.value = false;
        });
    });
}

function onShowDateTimeError(error: string): void {
    snackbar.value?.showError(error);
}

watch(activeTab, (newValue) => {
    if (newValue === 'map') {
        nextTick(() => {
            map.value?.initMapView();
        });
    }
});

defineExpose({
    open
});
</script>

<style>
.transaction-edit-amount .v-field__prepend-inner,
.transaction-edit-amount .v-field__append-inner,
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
