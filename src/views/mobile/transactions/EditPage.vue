<template>
    <f7-page with-subnavbar @page:afterin="onPageAfterIn" @page:beforeout="onPageBeforeOut">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt(title)"></f7-nav-title>
            <f7-nav-right>
                <f7-link icon-f7="ellipsis" @click="showMoreActionSheet = true"></f7-link>
                <f7-link :class="{ 'disabled': inputIsEmpty || submitting }" :text="tt(saveButtonTitle)" @click="save" v-if="mode !== TransactionEditPageMode.View"></f7-link>
            </f7-nav-right>

            <f7-subnavbar>
                <f7-segmented strong :class="{ 'readonly': pageTypeAndMode?.type === TransactionEditPageType.Transaction && mode !== TransactionEditPageMode.Add }">
                    <f7-button :text="tt('Expense')" :active="transaction.type === TransactionType.Expense"
                               :disabled="pageTypeAndMode?.type === TransactionEditPageType.Transaction && mode !== TransactionEditPageMode.Add && transaction.type !== TransactionType.Expense"
                               v-if="transaction.type !== TransactionType.ModifyBalance"
                               @click="transaction.type = TransactionType.Expense"></f7-button>
                    <f7-button :text="tt('Income')" :active="transaction.type === TransactionType.Income"
                               :disabled="pageTypeAndMode?.type === TransactionEditPageType.Transaction && mode !== TransactionEditPageMode.Add && transaction.type !== TransactionType.Income"
                               v-if="transaction.type !== TransactionType.ModifyBalance"
                               @click="transaction.type = TransactionType.Income"></f7-button>
                    <f7-button :text="tt('Transfer')" :active="transaction.type === TransactionType.Transfer"
                               :disabled="pageTypeAndMode?.type === TransactionEditPageType.Transaction && mode !== TransactionEditPageMode.Add && transaction.type !== TransactionType.Transfer"
                               v-if="transaction.type !== TransactionType.ModifyBalance"
                               @click="transaction.type = TransactionType.Transfer"></f7-button>
                    <f7-button :text="tt('Modify Balance')" :active="transaction.type === TransactionType.ModifyBalance"
                               v-if="pageTypeAndMode?.type === TransactionEditPageType.Transaction && transaction.type === TransactionType.ModifyBalance"></f7-button>
                </f7-segmented>
            </f7-subnavbar>
        </f7-navbar>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="loading">
            <f7-list-input label="Template Name" placeholder="Template Name" v-if="pageTypeAndMode?.type === TransactionEditPageType.Template"></f7-list-input>
            <f7-list-item
                class="transaction-edit-amount ebk-large-amount"
                header="Expense Amount" title="0.00">
            </f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-title-hide-overflow" header="Category" title="Category Names"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title" header="Account" title="Account Name"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title" header="Transaction Time" title="YYYY/MM/DD HH:mm:ss" v-if="pageTypeAndMode?.type === TransactionEditPageType.Transaction"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title" header="Scheduled Transaction Frequency" title="Every XXXXX" v-if="pageTypeAndMode?.type === TransactionEditPageType.Template && transaction instanceof TransactionTemplate && transaction.templateType === TemplateType.Schedule.type"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-title-hide-overflow list-item-no-item-after" header="Transaction Timezone" title="(UTC XX:XX) System Default" link="#" :no-chevron="mode === TransactionEditPageMode.View" v-if="pageTypeAndMode?.type === TransactionEditPageType.Transaction || (pageTypeAndMode?.type === TransactionEditPageType.Template && transaction instanceof TransactionTemplate && transaction.templateType === TemplateType.Schedule.type)"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-title-hide-overflow" header="Geographic Location" title="No Location" v-if="pageTypeAndMode?.type === TransactionEditPageType.Transaction"></f7-list-item>
            <f7-list-item header="Tags">
                <template #footer>
                    <f7-block class="margin-top-half no-padding no-margin">
                        <f7-chip class="transaction-edit-tag" text="None"></f7-chip>
                    </f7-block>
                </template>
            </f7-list-item>
            <f7-list-input type="textarea" label="Description" placeholder="Your transaction description (optional)"></f7-list-input>
        </f7-list>

        <f7-list form strong inset dividers class="margin-vertical" v-else-if="!loading">
            <f7-list-input
                type="text"
                clear-button
                :label="tt('Template Name')"
                :placeholder="tt('Template Name')"
                v-model:value="transaction.name"
                v-if="pageTypeAndMode?.type === TransactionEditPageType.Template && transaction instanceof TransactionTemplate"
            ></f7-list-input>

            <f7-list-item
                class="transaction-edit-amount"
                link="#" no-chevron
                :class="sourceAmountClass"
                :header="sourceAmountTitle"
                :title="getDisplayAmount(transaction.sourceAmount, transaction.hideAmount, sourceAccountCurrency)"
                @click="showSourceAmountSheet = true"
            >
                <number-pad-sheet :min-value="TRANSACTION_MIN_AMOUNT"
                                  :max-value="TRANSACTION_MAX_AMOUNT"
                                  :currency="sourceAccountCurrency"
                                  v-model:show="showSourceAmountSheet"
                                  v-model="transaction.sourceAmount"
                ></number-pad-sheet>
            </f7-list-item>

            <f7-list-item
                class="transaction-edit-amount text-color-primary"
                link="#" no-chevron
                :class="destinationAmountClass"
                :header="transferInAmountTitle"
                :title="getDisplayAmount(transaction.destinationAmount, transaction.hideAmount, destinationAccountCurrency)"
                @click="showDestinationAmountSheet = true"
                v-if="transaction.type === TransactionType.Transfer"
            >
                <number-pad-sheet :min-value="TRANSACTION_MIN_AMOUNT"
                                  :max-value="TRANSACTION_MAX_AMOUNT"
                                  :currency="destinationAccountCurrency"
                                  v-model:show="showDestinationAmountSheet"
                                  v-model="transaction.destinationAmount"
                ></number-pad-sheet>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-title-hide-overflow"
                key="expenseCategorySelection"
                link="#" no-chevron
                :class="{ 'disabled': !hasAvailableExpenseCategories, 'readonly': mode === TransactionEditPageMode.View }"
                :header="tt('Category')"
                @click="showCategorySheet = true"
                v-if="transaction.type === TransactionType.Expense"
            >
                <template #title>
                    <div class="list-item-custom-title" v-if="hasAvailableExpenseCategories">
                        <span>{{ getTransactionPrimaryCategoryName(transaction.expenseCategoryId, allCategories[CategoryType.Expense]) }}</span>
                        <f7-icon class="category-separate-icon" f7="chevron_right"></f7-icon>
                        <span>{{ getTransactionSecondaryCategoryName(transaction.expenseCategoryId, allCategories[CategoryType.Expense]) }}</span>
                    </div>
                    <div class="list-item-custom-title" v-else-if="!hasAvailableExpenseCategories">
                        <span>{{ tt('None') }}</span>
                    </div>
                </template>
                <tree-view-selection-sheet primary-key-field="id" primary-title-field="name"
                                           primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                           primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                           secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                           secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                           secondary-hidden-field="hidden"
                                           :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                           :items="allCategories[CategoryType.Expense]"
                                           v-model:show="showCategorySheet"
                                           v-model="transaction.expenseCategoryId">
                </tree-view-selection-sheet>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-title-hide-overflow"
                key="incomeCategorySelection"
                link="#" no-chevron
                :class="{ 'disabled': !hasAvailableIncomeCategories, 'readonly': mode === TransactionEditPageMode.View }"
                :header="tt('Category')"
                @click="showCategorySheet = true"
                v-if="transaction.type === TransactionType.Income"
            >
                <template #title>
                    <div class="list-item-custom-title" v-if="hasAvailableIncomeCategories">
                        <span>{{ getTransactionPrimaryCategoryName(transaction.incomeCategoryId, allCategories[CategoryType.Income]) }}</span>
                        <f7-icon class="category-separate-icon" f7="chevron_right"></f7-icon>
                        <span>{{ getTransactionSecondaryCategoryName(transaction.incomeCategoryId, allCategories[CategoryType.Income]) }}</span>
                    </div>
                    <div class="list-item-custom-title" v-else-if="!hasAvailableIncomeCategories">
                        <span>{{ tt('None') }}</span>
                    </div>
                </template>
                <tree-view-selection-sheet primary-key-field="id" primary-title-field="name"
                                           primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                           primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                           secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                           secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                           secondary-hidden-field="hidden"
                                           :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                           :items="allCategories[CategoryType.Income]"
                                           v-model:show="showCategorySheet"
                                           v-model="transaction.incomeCategoryId">
                </tree-view-selection-sheet>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-title-hide-overflow"
                key="transferCategorySelection"
                link="#" no-chevron
                :class="{ 'disabled': !hasAvailableTransferCategories, 'readonly': mode === TransactionEditPageMode.View }"
                :header="tt('Category')"
                @click="showCategorySheet = true"
                v-if="transaction.type === TransactionType.Transfer"
            >
                <template #title>
                    <div class="list-item-custom-title" v-if="hasAvailableTransferCategories">
                        <span>{{ getTransactionPrimaryCategoryName(transaction.transferCategoryId, allCategories[CategoryType.Transfer]) }}</span>
                        <f7-icon class="category-separate-icon" f7="chevron_right"></f7-icon>
                        <span>{{ getTransactionSecondaryCategoryName(transaction.transferCategoryId, allCategories[CategoryType.Transfer]) }}</span>
                    </div>
                    <div class="list-item-custom-title" v-else-if="!hasAvailableTransferCategories">
                        <span>{{ tt('None') }}</span>
                    </div>
                </template>
                <tree-view-selection-sheet primary-key-field="id" primary-title-field="name"
                                           primary-icon-field="icon" primary-icon-type="category" primary-color-field="color"
                                           primary-hidden-field="hidden" primary-sub-items-field="subCategories"
                                           secondary-key-field="id" secondary-value-field="id" secondary-title-field="name"
                                           secondary-icon-field="icon" secondary-icon-type="category" secondary-color-field="color"
                                           secondary-hidden-field="hidden"
                                           :enable-filter="true" :filter-placeholder="tt('Find category')" :filter-no-items-text="tt('No available category')"
                                           :items="allCategories[CategoryType.Transfer]"
                                           v-model:show="showCategorySheet"
                                           v-model="transaction.transferCategoryId">
                </tree-view-selection-sheet>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title"
                link="#" no-chevron
                :class="{ 'disabled': !allVisibleAccounts.length, 'readonly': mode === TransactionEditPageMode.View }"
                :header="tt(sourceAccountTitle)"
                :title="sourceAccountName"
                @click="showSourceAccountSheet = true"
            >
                <two-column-list-item-selection-sheet primary-key-field="id" primary-value-field="category"
                                                      primary-title-field="name" primary-footer-field="displayBalance"
                                                      primary-icon-field="icon" primary-icon-type="account"
                                                      primary-sub-items-field="accounts"
                                                      :primary-title-i18n="true"
                                                      secondary-key-field="id" secondary-value-field="id"
                                                      secondary-title-field="name" secondary-footer-field="displayBalance"
                                                      secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                      :enable-filter="true" :filter-placeholder="tt('Find account')" :filter-no-items-text="tt('No available account')"
                                                      :items="allVisibleCategorizedAccounts"
                                                      v-model:show="showSourceAccountSheet"
                                                      v-model="transaction.sourceAccountId">
                </two-column-list-item-selection-sheet>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title"
                link="#" no-chevron
                :class="{ 'disabled': !allVisibleAccounts.length, 'readonly': mode === TransactionEditPageMode.View }"
                :header="tt('Destination Account')"
                :title="destinationAccountName"
                v-if="transaction.type === TransactionType.Transfer"
                @click="showDestinationAccountSheet = true"
            >
                <two-column-list-item-selection-sheet primary-key-field="id" primary-value-field="category"
                                                      primary-title-field="name" primary-footer-field="displayBalance"
                                                      primary-icon-field="icon" primary-icon-type="account"
                                                      primary-sub-items-field="accounts"
                                                      :primary-title-i18n="true"
                                                      secondary-key-field="id" secondary-value-field="id"
                                                      secondary-title-field="name" secondary-footer-field="displayBalance"
                                                      secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                      :enable-filter="true" :filter-placeholder="tt('Find account')" :filter-no-items-text="tt('No available account')"
                                                      :items="allVisibleCategorizedAccounts"
                                                      v-model:show="showDestinationAccountSheet"
                                                      v-model="transaction.destinationAccountId">
                </two-column-list-item-selection-sheet>
            </f7-list-item>

            <f7-list-item
                class="transaction-edit-datetime list-item-with-header-and-title"
                link="#" no-chevron
                :class="{ 'readonly': mode === TransactionEditPageMode.View && transaction.utcOffset === currentTimezoneOffsetMinutes }"
                v-if="pageTypeAndMode?.type === TransactionEditPageType.Transaction"
            >
                <template #header>
                    <div class="transaction-edit-datetime-header" @click="showDateTimeDialog('time')">{{ tt('Transaction Time') }}</div>
                </template>
                <template #title>
                    <div class="transaction-edit-datetime-title">
                        <div @click="showDateTimeDialog('date')">{{ transactionDisplayDate }}</div>&nbsp;<div class="transaction-edit-datetime-time" @click="showDateTimeDialog('time')">{{ transactionDisplayTime }}</div>
                    </div>
                </template>
                <date-time-selection-sheet :init-mode="transactionDateTimeSheetMode"
                                           v-model:show="showTransactionDateTimeSheet"
                                           v-model="transaction.time">
                </date-time-selection-sheet>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title"
                link="#" no-chevron
                :class="{ 'readonly': mode === TransactionEditPageMode.View }"
                :header="tt('Scheduled Transaction Frequency')"
                :title="transactionDisplayScheduledFrequency"
                @click="showTransactionScheduledFrequencySheet = true"
                v-if="pageTypeAndMode?.type === TransactionEditPageType.Template && transaction instanceof TransactionTemplate && transaction.templateType === TemplateType.Schedule.type"
            >
                <schedule-frequency-sheet v-model:show="showTransactionScheduledFrequencySheet"
                                          v-model:type="transaction.scheduledFrequencyType"
                                          v-model="transaction.scheduledFrequency">
                </schedule-frequency-sheet>
            </f7-list-item>

            <f7-list-item
                class="transaction-edit-datetime list-item-with-header-and-title"
                link="#" no-chevron
                :class="{ 'readonly': mode === TransactionEditPageMode.View }"
                :header="tt('Start Date')"
                :title="transactionDisplayScheduledStartDate"
                @click="showScheduledStartDateSheet = true"
                v-if="pageTypeAndMode?.type === TransactionEditPageType.Template && transaction instanceof TransactionTemplate && transaction.templateType === TemplateType.Schedule.type"
            >
                <date-selection-sheet v-model:show="showScheduledStartDateSheet"
                                      v-model="transaction.scheduledStartDate">
                </date-selection-sheet>
            </f7-list-item>

            <f7-list-item
                class="transaction-edit-datetime list-item-with-header-and-title"
                link="#" no-chevron
                :class="{ 'readonly': mode === TransactionEditPageMode.View }"
                :header="tt('End Date')"
                :title="transactionDisplayScheduledEndDate"
                @click="showScheduledEndDateSheet = true"
                v-if="pageTypeAndMode?.type === TransactionEditPageType.Template && transaction instanceof TransactionTemplate && transaction.templateType === TemplateType.Schedule.type"
            >
                <date-selection-sheet v-model:show="showScheduledEndDateSheet"
                                      v-model="transaction.scheduledEndDate">
                </date-selection-sheet>
            </f7-list-item>

            <f7-list-item
                :no-chevron="mode === TransactionEditPageMode.View"
                link="#"
                class="list-item-with-header-and-title list-item-title-hide-overflow list-item-no-item-after"
                :class="{ 'readonly': mode === TransactionEditPageMode.View }"
                :header="tt('Transaction Timezone')"
                v-if="pageTypeAndMode?.type === TransactionEditPageType.Transaction || (pageTypeAndMode?.type === TransactionEditPageType.Template && transaction instanceof TransactionTemplate && transaction.templateType === TemplateType.Schedule.type)"
                @click="showTimezonePopup = true"
            >
                <template #title>
                    <f7-block class="list-item-custom-title no-padding no-margin">
                        <span>{{ `(${transactionDisplayTimezone})` }}</span>
                        <span class="transaction-edit-timezone-name" v-if="transaction.timeZone || transaction.timeZone === ''">{{ transactionDisplayTimezoneName }}</span>
                        <span class="transaction-edit-timezone-name" v-else-if="!transaction.timeZone && transaction.timeZone !== ''">{{ transactionTimezoneTimeDifference }}</span>
                    </f7-block>
                </template>
                <list-item-selection-popup value-type="item"
                                           key-field="name" value-field="name"
                                           title-field="displayNameWithUtcOffset"
                                           :title="tt('Transaction Timezone')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Timezone')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allTimezones"
                                           v-model:show="showTimezonePopup"
                                           v-model="transaction.timeZone">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item
                link="#" no-chevron
                class="list-item-with-header-and-title list-item-title-hide-overflow"
                :class="{ 'readonly': mode === TransactionEditPageMode.View && !transaction.geoLocation }"
                :header="tt('Geographic Location')"
                @click="showGeoLocationActionSheet = true"
                v-if="pageTypeAndMode?.type === TransactionEditPageType.Transaction"
            >
                <template #title>
                    <f7-block class="list-item-custom-title no-padding no-margin">
                        <span v-if="transaction.geoLocation">{{ `(${transaction.geoLocation.longitude}, ${transaction.geoLocation.latitude})` }}</span>
                        <span v-else-if="!transaction.geoLocation">{{ geoLocationStatusInfo }}</span>
                    </f7-block>
                </template>

                <map-sheet v-model="transaction.geoLocation"
                           v-model:set-geo-location-by-click-map="setGeoLocationByClickMap"
                           v-model:show="showGeoLocationMapSheet">
                </map-sheet>
            </f7-list-item>

            <f7-list-item
                link="#" no-chevron
                :class="{ 'readonly': mode === TransactionEditPageMode.View }"
                :header="tt('Tags')"
                @click="showTransactionTagSheet = true"
            >
                <transaction-tag-selection-sheet :allow-add-new-tag="true" :enable-filter="true"
                                                 v-model:show="showTransactionTagSheet"
                                                 v-model="transaction.tagIds">
                </transaction-tag-selection-sheet>

                <template #footer>
                    <f7-block class="margin-top-half no-padding no-margin" v-if="transaction.tagIds && transaction.tagIds.length">
                        <f7-chip media-text-color="var(--f7-chip-text-color)" class="transaction-edit-tag"
                                 :text="getTagName(tagId)"
                                 :key="tagId"
                                 v-for="tagId in transaction.tagIds">
                            <template #media>
                                <f7-icon f7="number"></f7-icon>
                            </template>
                        </f7-chip>
                    </f7-block>
                    <f7-block class="margin-top-half no-padding no-margin" v-else-if="!transaction.tagIds || !transaction.tagIds.length">
                        <f7-chip class="transaction-edit-tag" :text="tt('None')">
                        </f7-chip>
                    </f7-block>
                </template>
            </f7-list-item>

            <f7-list-item
                link="#" no-chevron
                :header="tt('Pictures')"
                v-if="showTransactionPictures || (transaction.pictures && transaction.pictures.length > 0)"
            >
                <template #footer>
                    <f7-block class="margin-top-half no-padding no-margin" :class="{ 'readonly': submitting || uploadingPicture || removingPictureId }">
                        <swiper-container
                            :pagination="false"
                            :space-between="10"
                            :slides-per-view="'auto'"
                            class="transaction-pictures"
                        >
                            <swiper-slide class="transaction-picture-container" :key="picIdx"
                                          v-for="(pictureInfo, picIdx) in transaction.pictures"
                                          @click="viewOrRemovePicture(pictureInfo)">
                                <div class="transaction-picture">
                                    <div class="display-flex justify-content-center align-items-center transaction-picture-control-backdrop"
                                         v-if="mode === TransactionEditPageMode.Add || mode === TransactionEditPageMode.Edit">
                                        <f7-icon class="picture-control-icon picture-remove-icon" f7="trash" v-if="pictureInfo.pictureId !== removingPictureId"></f7-icon>
                                        <f7-preloader color="white" :size="28" v-if="pictureInfo.pictureId === removingPictureId" />
                                    </div>
                                    <img alt="picture" :src="getTransactionPictureUrl(pictureInfo)"/>
                                </div>
                            </swiper-slide>
                            <swiper-slide @click="showOpenPictureDialog" v-if="canAddTransactionPicture">
                                <div class="display-flex justify-content-center align-items-center transaction-picture transaction-picture-add">
                                    <f7-icon class="picture-control-icon" f7="plus" v-if="!uploadingPicture"></f7-icon>
                                    <f7-preloader :size="28" v-if="uploadingPicture" />
                                </div>
                            </swiper-slide>
                        </swiper-container>
                    </f7-block>
                </template>
            </f7-list-item>

            <f7-list-input
                type="textarea"
                class="transaction-edit-comment"
                style="height: auto"
                :class="{ 'readonly': mode === TransactionEditPageMode.View }"
                :label="tt('Description')"
                :placeholder="mode !== TransactionEditPageMode.View ? tt('Your transaction description (optional)') : ''"
                v-textarea-auto-size
                v-model:value="transaction.comment"
            ></f7-list-input>
        </f7-list>

        <f7-actions close-by-outside-click close-on-escape :opened="showGeoLocationActionSheet" @actions:closed="showGeoLocationActionSheet = false">
            <f7-actions-group>
                <f7-actions-button v-if="mode !== TransactionEditPageMode.View" @click="updateGeoLocation(true)">{{ tt('Update Geographic Location') }}</f7-actions-button>
                <f7-actions-button v-if="mode !== TransactionEditPageMode.View" @click="clearGeoLocation">{{ tt('Clear Geographic Location') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group v-if="!!getMapProvider()">
                <f7-actions-button :class="{ 'disabled': !transaction.geoLocation }" @click="showGeoLocationMapSheet = true">{{ tt('Show on the map') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group v-if="mode !== TransactionEditPageMode.View && transaction.type === TransactionType.Transfer">
                <f7-actions-button @click="swapTransactionData(true, false)">{{ tt('Swap Account') }}</f7-actions-button>
                <f7-actions-button @click="swapTransactionData(false, true)">{{ tt('Swap Amount') }}</f7-actions-button>
                <f7-actions-button @click="swapTransactionData(true, true)">{{ tt('Swap Account and Amount') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group v-if="mode !== TransactionEditPageMode.View">
                <f7-actions-button v-if="transaction.hideAmount" @click="transaction.hideAmount = false">{{ tt('Show Amount') }}</f7-actions-button>
                <f7-actions-button v-if="!transaction.hideAmount" @click="transaction.hideAmount = true">{{ tt('Hide Amount') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group v-if="pageTypeAndMode?.type === TransactionEditPageType.Transaction && (mode === TransactionEditPageMode.Add || mode === TransactionEditPageMode.Edit) && isTransactionPicturesEnabled() && !showTransactionPictures">
                <f7-actions-button @click="showTransactionPictures = true">{{ tt('Add Picture') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group v-if="pageTypeAndMode?.type === TransactionEditPageType.Transaction && mode === TransactionEditPageMode.View">
                <f7-actions-button @click="duplicate(false, false)">{{ tt('Duplicate') }}</f7-actions-button>
                <f7-actions-button @click="duplicate(true, false)">{{ tt('Duplicate (With Time)') }}</f7-actions-button>
                <f7-actions-button @click="duplicate(false, true)" v-if="transaction.geoLocation">{{ tt('Duplicate (With Geographic Location)') }}</f7-actions-button>
                <f7-actions-button @click="duplicate(true, true)" v-if="transaction.geoLocation">{{ tt('Duplicate (With Time and Geographic Location)') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-toolbar tabbar bottom v-if="mode !== TransactionEditPageMode.View">
            <f7-link :class="{ 'disabled': inputIsEmpty || submitting }" @click="save">
                <span class="tabbar-primary-link">{{ tt(saveButtonTitle) }}</span>
            </f7-link>
        </f7-toolbar>

        <f7-photo-browser ref="pictureBrowser" type="popup" navbar-of-text="/"
                          :theme="isDarkMode ? 'dark' : 'light'" :navbar-show-count="true" :exposition="false"
                          :photos="transactionPictures" :thumbs="transactionThumbs" />
        <input ref="pictureInput" type="file" style="display: none" :accept="SUPPORTED_IMAGE_EXTENSIONS" @change="uploadPicture($event)" />
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed, useTemplateRef } from 'vue';
import type { PhotoBrowser, Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';
import {
    TransactionEditPageMode,
    TransactionEditPageType,
    GeoLocationStatus,
    useTransactionEditPageBase
} from '@/views/base/transactions/TransactionEditPageBase.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useEnvironmentsStore } from '@/stores/environment.ts';
import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';
import { useTransactionsStore } from '@/stores/transaction.ts';
import { useTransactionTemplatesStore } from '@/stores/transactionTemplate.ts';

import { CategoryType } from '@/core/category.ts';
import { TransactionEditScopeType, TransactionType } from '@/core/transaction.ts';
import { ScheduledTemplateFrequencyType, TemplateType } from '@/core/template.ts';
import { TRANSACTION_MAX_AMOUNT, TRANSACTION_MIN_AMOUNT } from '@/consts/transaction.ts';
import { KnownErrorCode } from '@/consts/api.ts';
import { SUPPORTED_IMAGE_EXTENSIONS } from '@/consts/file.ts';

import { TransactionTemplate } from '@/models/transaction_template.ts';
import type { TransactionPictureInfoBasicResponse } from '@/models/transaction_picture_info.ts';
import { Transaction } from '@/models/transaction.ts';

import {
    getActualUnixTimeForStore,
    getBrowserTimezoneOffsetMinutes,
    getTimezoneOffset,
    getTimezoneOffsetMinutes
} from '@/lib/datetime.ts';
import { generateRandomUUID } from '@/lib/misc.ts';
import { getTransactionPrimaryCategoryName, getTransactionSecondaryCategoryName } from '@/lib/category.ts';
import { setTransactionModelByTransaction } from '@/lib/transaction.ts';
import { getMapProvider, isTransactionPicturesEnabled } from '@/lib/server_settings.ts';
import logger from '@/lib/logger.ts';

const props = defineProps<{
    f7route: Router.Route;
    f7router: Router.Router;
}>();

const query = props.f7route.query;
const pageTypeAndMode = getPageTypeNameMode();

const {
    tt,
    getMultiMonthdayShortNames,
    getMultiWeekdayLongNames,
    formatUnixTimeToLongDate,
    formatUnixTimeToLongTime,
    formatDateToLongDate
} = useI18n();
const { showAlert, showConfirm, showToast, routeBackOnError } = useI18nUIComponents();

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
    currentTimezoneOffsetMinutes,
    defaultCurrency,
    firstDayOfWeek,
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
    swapTransactionData,
    getDisplayAmount,
    getTransactionPictureUrl
} = useTransactionEditPageBase(pageTypeAndMode?.type || TransactionEditPageType.Transaction, pageTypeAndMode?.mode, query['type'] ? parseInt(query['type']) : undefined);

const settingsStore = useSettingsStore();
const environmentsStore = useEnvironmentsStore();
const userStore = useUserStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();
const transactionsStore = useTransactionsStore();
const transactionTemplatesStore = useTransactionTemplatesStore();

const pictureBrowser = useTemplateRef<PhotoBrowser.PhotoBrowser>('pictureBrowser');
const pictureInput = useTemplateRef<HTMLInputElement>('pictureInput');

const loadingError = ref<unknown | null>(null);
const submitted = ref<boolean>(false);
const removingPictureId = ref<string | null>(null);
const transactionDateTimeSheetMode = ref<string>('time');
const showTimeInDefaultTimezone = ref<boolean>(false);
const showTimezonePopup = ref<boolean>(false);
const showGeoLocationActionSheet = ref<boolean>(false);
const showMoreActionSheet = ref<boolean>(false);
const showSourceAmountSheet = ref<boolean>(false);
const showDestinationAmountSheet = ref<boolean>(false);
const showCategorySheet = ref<boolean>(false);
const showSourceAccountSheet = ref<boolean>(false);
const showDestinationAccountSheet = ref<boolean>(false);
const showTransactionDateTimeSheet = ref<boolean>(false);
const showTransactionScheduledFrequencySheet = ref<boolean>(false);
const showScheduledStartDateSheet = ref<boolean>(false);
const showScheduledEndDateSheet = ref<boolean>(false);
const showGeoLocationMapSheet = ref<boolean>(false);
const showTransactionTagSheet = ref<boolean>(false);
const showTransactionPictures = ref<boolean>(false);

const isDarkMode = computed<boolean>(() => environmentsStore.framework7DarkMode || false);

const sourceAmountClass = computed<Record<string, boolean>>(() => {
    const classes: Record<string, boolean> = {
        'readonly': mode.value === TransactionEditPageMode.View,
        'text-expense': transaction.value.type === TransactionType.Expense,
        'text-income': transaction.value.type === TransactionType.Income,
        'text-color-primary': transaction.value.type === TransactionType.Transfer
    };

    classes[getFontClassByAmount(transaction.value.sourceAmount)] = true;

    return classes;
});

const destinationAmountClass = computed<Record<string, boolean>>(() => {
    const classes: Record<string, boolean> = {
        'readonly': mode.value === TransactionEditPageMode.View
    };

    classes[getFontClassByAmount(transaction.value.destinationAmount)] = true;

    return classes;
});

const transactionDisplayDate = computed<string>(() => {
    if (mode.value !== TransactionEditPageMode.View || !showTimeInDefaultTimezone.value) {
        return formatUnixTimeToLongDate(getActualUnixTimeForStore(transaction.value.time, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()));
    }

    return formatUnixTimeToLongDate(getActualUnixTimeForStore(transaction.value.time, transaction.value.utcOffset, getBrowserTimezoneOffsetMinutes()));
});

const transactionDisplayTime = computed<string>(() => {
    if (mode.value !== TransactionEditPageMode.View || !showTimeInDefaultTimezone.value) {
        return formatUnixTimeToLongTime(getActualUnixTimeForStore(transaction.value.time, getTimezoneOffsetMinutes(), getBrowserTimezoneOffsetMinutes()));
    }

    return `${formatUnixTimeToLongTime(getActualUnixTimeForStore(transaction.value.time, transaction.value.utcOffset, getBrowserTimezoneOffsetMinutes()))} (UTC${getTimezoneOffset(settingsStore.appSettings.timeZone)})`;
});

const transactionDisplayTimezoneName = computed<string>(() => {
    for (const timezone of allTimezones.value) {
        if (timezone.name === transaction.value.timeZone) {
            return timezone.displayName;
        }
    }

    return '';
});

const transactionPictures = computed<Record<string, string | undefined>[]>(() => {
    const thumbs: Record<string, string | undefined>[] = [];

    if (!transaction.value.pictures || !transaction.value.pictures.length) {
        return thumbs;
    }

    for (let i = 0; i < transaction.value.pictures.length; i++) {
        thumbs.push({
            url: getTransactionPictureUrl(transaction.value.pictures[i])
        });
    }

    return thumbs;
});

const transactionThumbs = computed<(string | undefined)[]>(() => {
    const thumbs: (string | undefined)[] = [];

    if (!transaction.value.pictures || !transaction.value.pictures.length) {
        return thumbs;
    }

    for (let i = 0; i < transaction.value.pictures.length; i++) {
        thumbs.push(getTransactionPictureUrl(transaction.value.pictures[i]));
    }

    return thumbs;
});

const transactionDisplayScheduledFrequency = computed<string>(() => {
    if (pageTypeAndMode?.type !== TransactionEditPageType.Template) {
        return '';
    }

    const template = transaction.value as TransactionTemplate;

    if (template.scheduledFrequencyType === ScheduledTemplateFrequencyType.Disabled.type) {
        return tt('Disabled');
    }

    const items = (template.scheduledFrequency || '').split(',');
    const scheduledFrequencyValues: number[] = [];

    for (let i = 0; i < items.length; i++) {
        if (items[i]) {
            scheduledFrequencyValues.push(parseInt(items[i]));
        }
    }

    if (template.scheduledFrequencyType === ScheduledTemplateFrequencyType.Weekly.type) {
        if (scheduledFrequencyValues.length) {
            return tt('format.misc.everyMultiDaysOfWeek', {
                days: getMultiWeekdayLongNames(scheduledFrequencyValues, firstDayOfWeek.value)
            });
        } else {
            return tt('Weekly');
        }
    } else if (template.scheduledFrequencyType === ScheduledTemplateFrequencyType.Monthly.type) {
        if (scheduledFrequencyValues.length) {
            return tt('format.misc.everyMultiDaysOfMonth', {
                days: getMultiMonthdayShortNames(scheduledFrequencyValues)
            });
        } else {
            return tt('Monthly');
        }
    } else {
        return '';
    }
});

const transactionDisplayScheduledStartDate = computed<string>(() => {
    if (pageTypeAndMode?.type !== TransactionEditPageType.Template) {
        return '';
    }

    const template = transaction.value as TransactionTemplate;

    if (template.scheduledStartDate) {
        return formatDateToLongDate(template.scheduledStartDate);
    } else {
        return tt('Unspecified');
    }
});

const transactionDisplayScheduledEndDate = computed<string>(() => {
    if (pageTypeAndMode?.type !== TransactionEditPageType.Template) {
        return '';
    }

    const template = transaction.value as TransactionTemplate;

    if (template.scheduledEndDate) {
        return formatDateToLongDate(template.scheduledEndDate);
    } else {
        return tt('Unspecified');
    }
});

function getPageTypeNameMode(): { type: TransactionEditPageType, mode: TransactionEditPageMode } | null {
    if (props.f7route.path === '/transaction/add') {
        return {
            type: TransactionEditPageType.Transaction,
            mode: TransactionEditPageMode.Add
        };
    } else if (props.f7route.path === '/transaction/edit') {
        return {
            type: TransactionEditPageType.Transaction,
            mode: TransactionEditPageMode.Edit
        };
    } else if (props.f7route.path === '/transaction/detail') {
        return {
            type: TransactionEditPageType.Transaction,
            mode: TransactionEditPageMode.View
        };
    } else if (props.f7route.path === '/template/add') {
        return {
            type: TransactionEditPageType.Template,
            mode: TransactionEditPageMode.Add
        };
    } else if (props.f7route.path === '/template/edit') {
        return {
            type: TransactionEditPageType.Template,
            mode: TransactionEditPageMode.Edit
        };
    } else {
        return null;
    }
}

function getFontClassByAmount(amount: number): string {
    if (amount >= 100000000 || amount <= -100000000) {
        return 'ebk-small-amount';
    } else if (amount >= 1000000 || amount <= -1000000) {
        return 'ebk-normal-amount';
    } else {
        return 'ebk-large-amount';
    }
}

function getTagName(tagId: string): string {
    for (const tag of allTags.value) {
        if (tag.id === tagId) {
            return tag.name;
        }
    }

    return '';
}

function init(): void {
    if (!pageTypeAndMode) {
        showToast('Parameter Invalid');
        loadingError.value = 'Parameter Invalid';
        return;
    }

    loading.value = true;

    const promises: Promise<unknown>[] = [
        accountsStore.loadAllAccounts({ force: false }),
        transactionCategoriesStore.loadAllCategories({ force: false }),
        transactionTagsStore.loadAllTags({ force: false }),
        transactionTemplatesStore.loadAllTemplates({ force: false, templateType: TemplateType.Normal.type })
    ];

    if (pageTypeAndMode.type === TransactionEditPageType.Transaction) {
        if (query['id']) {
            if (mode.value === TransactionEditPageMode.Edit) {
                editId.value = query['id'];
            } else if (mode.value === TransactionEditPageMode.Add) {
                duplicateFromId.value = query['id'];
            }

            promises.push(transactionsStore.getTransaction({ transactionId: query['id'], withPictures: mode.value !== TransactionEditPageMode.Add }));
        }
    } else if (pageTypeAndMode.type === TransactionEditPageType.Template) {
        const template = TransactionTemplate.createNewTransactionTemplate(transaction.value);
        template.name = '';

        if (query['templateType']) {
            template.templateType = parseInt(query['templateType']);
        }

        if (template.templateType === TemplateType.Schedule.type) {
            template.scheduledFrequencyType = ScheduledTemplateFrequencyType.Disabled.type;
            template.scheduledFrequency = '';
        }

        transaction.value = template;

        if (query['id']) {
            if (mode.value === TransactionEditPageMode.Edit) {
                editId.value = query['id'];
            }

            promises.push(transactionTemplatesStore.getTemplate({ templateId: query['id'] }));
        }
    }

    const queryType = query['type'] ? parseInt(query['type']) : 0;

    if (queryType &&
        queryType >= TransactionType.Income &&
        queryType <= TransactionType.Transfer) {
        transaction.value.type = queryType;
    }

    if (mode.value === TransactionEditPageMode.Add) {
        clientSessionId.value = generateRandomUUID();
    }

    Promise.all(promises).then(function (responses) {
        if (query['id'] && !responses[4]) {
            if (pageTypeAndMode.type === TransactionEditPageType.Transaction) {
                showToast('Unable to retrieve transaction');
                loadingError.value = 'Unable to retrieve transaction';
            } else if (pageTypeAndMode.type === TransactionEditPageType.Template) {
                showToast('Unable to retrieve template');
                loadingError.value = 'Unable to retrieve template';
            }

            return;
        }

        let fromTransaction: Transaction | TransactionTemplate | null = null;

        if (pageTypeAndMode.type === TransactionEditPageType.Transaction) {
            if (query['id'] && responses[4] instanceof Transaction) {
                fromTransaction = responses[4];
            } else if (query['templateId'] && transactionTemplatesStore.allTransactionTemplatesMap && transactionTemplatesStore.allTransactionTemplatesMap[TemplateType.Normal.type]) {
                fromTransaction = transactionTemplatesStore.allTransactionTemplatesMap[TemplateType.Normal.type][query['templateId']];

                if (fromTransaction) {
                    addByTemplateId.value = fromTransaction.id;
                }
            } else if ((settingsStore.appSettings.autoSaveTransactionDraft === 'enabled' || settingsStore.appSettings.autoSaveTransactionDraft === 'confirmation') && transactionsStore.transactionDraft) {
                fromTransaction = Transaction.ofDraft(transactionsStore.transactionDraft);
            }
        } else if (pageTypeAndMode.type === TransactionEditPageType.Template && responses[4] instanceof TransactionTemplate) {
            if (query['id']) {
                fromTransaction = responses[4];
            }
        }

        setTransactionModelByTransaction(
            transaction.value,
            fromTransaction,
            allCategories.value,
            allCategoriesMap.value,
            allVisibleAccounts.value,
            allAccountsMap.value,
            allTagsMap.value,
            defaultAccountId.value,
            {
                type: queryType,
                categoryId: query['categoryId'],
                accountId: query['accountId'],
                destinationAccountId: query['destinationAccountId'],
                amount: query['amount'] ? parseInt(query['amount']) : undefined,
                destinationAmount: query['destinationAmount'] ? parseInt(query['destinationAmount']) : undefined,
                tagIds: query['tagIds'],
                comment: query['comment']
            },
            pageTypeAndMode.type === TransactionEditPageType.Transaction && (mode.value === TransactionEditPageMode.Edit || mode.value === TransactionEditPageMode.View),
            pageTypeAndMode.type === TransactionEditPageType.Transaction && (mode.value === TransactionEditPageMode.Edit || mode.value === TransactionEditPageMode.View)
        );

        if (pageTypeAndMode.type === TransactionEditPageType.Transaction && query['id'] && responses[4] instanceof Transaction) {
            if (fromTransaction && query['withTime'] && query['withTime'] === 'true') {
                transaction.value.time = fromTransaction.time;
                transaction.value.timeZone = fromTransaction.timeZone;
                transaction.value.utcOffset = fromTransaction.utcOffset;
            }

            if (fromTransaction && query['withGeoLocation'] && query['withGeoLocation'] === 'true') {
                transaction.value.setGeoLocation(fromTransaction.geoLocation);
            }
        } else if (pageTypeAndMode.type === TransactionEditPageType.Template && query['id'] && responses[4] instanceof TransactionTemplate) {
            const template = responses[4];
            transaction.value.id = template.id;

            if (!(transaction.value instanceof TransactionTemplate)) {
                transaction.value = TransactionTemplate.createNewTransactionTemplate(transaction.value);
            }

            (transaction.value as TransactionTemplate).fillFrom(template);
        }

        loading.value = false;
    }).catch(error => {
        logger.error('failed to load essential data for editing transaction', error);

        if (error.processed) {
            loading.value = false;
        } else {
            loadingError.value = error;
            showToast(error.message || error);
        }
    });
}

function save(): void {
    const router = props.f7router;

    if (mode.value === TransactionEditPageMode.View) {
        return;
    }

    const problemMessage = inputEmptyProblemMessage.value;

    if (problemMessage) {
        showAlert(problemMessage);
        return;
    }

    if (pageTypeAndMode?.type === TransactionEditPageType.Transaction && (mode.value === TransactionEditPageMode.Add || mode.value === TransactionEditPageMode.Edit)) {
        const doSubmit = function () {
            submitting.value = true;
            showLoading(() => submitting.value);

            transactionsStore.saveTransaction({
                transaction: transaction.value as Transaction,
                defaultCurrency: defaultCurrency.value,
                isEdit: mode.value === TransactionEditPageMode.Edit,
                clientSessionId: clientSessionId.value
            }).then(() => {
                submitting.value = false;
                hideLoading();

                if (mode.value === TransactionEditPageMode.Add) {
                    showToast('You have added a new transaction');
                } else if (mode.value === TransactionEditPageMode.Edit) {
                    showToast('You have saved this transaction');
                }

                if (mode.value === TransactionEditPageMode.Add && !addByTemplateId.value && !duplicateFromId.value) {
                    transactionsStore.clearTransactionDraft();
                }

                submitted.value = true;
                router.back();
            }).catch(error => {
                submitting.value = false;
                hideLoading();

                if (error.error && (error.error.errorCode === KnownErrorCode.TransactionCannotCreateInThisTime || error.error.errorCode === KnownErrorCode.TransactionCannotModifyInThisTime)) {
                    showConfirm('You have set this time range to prevent editing transactions. Would you like to change the editable transaction range to All?', () => {
                        submitting.value = true;
                        showLoading(() => submitting.value);

                        userStore.updateUserTransactionEditScope({
                            transactionEditScope: TransactionEditScopeType.All.type
                        }).then(() => {
                            submitting.value = false;
                            hideLoading();

                            showToast('Your editable transaction range has been set to All');
                        }).catch(error => {
                            submitting.value = false;
                            hideLoading();

                            if (!error.processed) {
                                showToast(error.message || error);
                            }
                        });
                    });
                } else if (!error.processed) {
                    showToast(error.message || error);
                }
            });
        };

        if (transaction.value.sourceAmount === 0) {
            showConfirm('Are you sure you want to save this transaction with a zero amount?', () => {
                doSubmit();
            });
        } else {
            doSubmit();
        }
    } else if (pageTypeAndMode?.type === TransactionEditPageType.Template && (mode.value === TransactionEditPageMode.Add || mode.value === TransactionEditPageMode.Edit)) {
        submitting.value = true;
        showLoading(() => submitting.value);

        transactionTemplatesStore.saveTemplateContent({
            template: transaction.value as TransactionTemplate,
            isEdit: mode.value === TransactionEditPageMode.Edit,
            clientSessionId: clientSessionId.value
        }).then(() => {
            submitting.value = false;
            hideLoading();

            if (mode.value === TransactionEditPageMode.Add) {
                showToast('You have added a new template');
            } else if (mode.value === TransactionEditPageMode.Edit) {
                showToast('You have saved this template');
            }

            submitted.value = true;
            router.back();
        }).catch(error => {
            submitting.value = false;
            hideLoading();

            if (!error.processed) {
                showToast(error.message || error);
            }
        });
    }
}

function updateGeoLocation(forceUpdate: boolean): void {
    if (!isSupportGeoLocation) {
        logger.warn('this browser does not support geo location');

        if (forceUpdate) {
            showToast('Unable to retrieve current position');
        }
        return;
    }

    navigator.geolocation.getCurrentPosition(function (position) {
        if (!position || !position.coords) {
            logger.error('current position is null');
            geoLocationStatus.value = GeoLocationStatus.Error;

            if (forceUpdate) {
                showToast('Unable to retrieve current position');
            }

            return;
        }

        geoLocationStatus.value = GeoLocationStatus.Success;

        transaction.value.setLatitudeAndLongitude(position.coords.latitude, position.coords.longitude);
    }, function (err) {
        logger.error('cannot retrieve current position', err);
        geoLocationStatus.value = GeoLocationStatus.Error;

        if (forceUpdate) {
            showToast('Unable to retrieve current position');
        }
    });

    geoLocationStatus.value = GeoLocationStatus.Getting;
}

function clearGeoLocation(): void {
    geoLocationStatus.value = null;
    transaction.value.removeGeoLocation();
}

function showDateTimeDialog(sheetMode: string): void {
    if (mode.value === TransactionEditPageMode.View) {
        showTimeInDefaultTimezone.value = !showTimeInDefaultTimezone.value;
    } else {
        transactionDateTimeSheetMode.value = sheetMode;
        showTransactionDateTimeSheet.value = true;
    }
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
            showToast(error.message || error);
        }
    });
}

function viewOrRemovePicture(pictureInfo: TransactionPictureInfoBasicResponse): void {
    if (mode.value !== TransactionEditPageMode.Add && mode.value !== TransactionEditPageMode.Edit && transaction.value.pictures && transaction.value.pictures.length) {
        pictureBrowser.value?.open();
        return;
    }

    showConfirm('Are you sure you want to remove this transaction picture?', () => {
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
                showToast(error.message || error);
            }

            removingPictureId.value = '';
            submitting.value = false;
        });
    });
}

function duplicate(withTime?: boolean, withGeoLocation?: boolean): void {
    props.f7router.navigate(`/transaction/add?id=${transaction.value.id}&type=${transaction.value.type}&withTime=${withTime ?? false}&withGeoLocation=${withGeoLocation ?? false}`);
}

function onPageAfterIn(): void {
    routeBackOnError(props.f7router, loadingError);

    if (settingsStore.appSettings.autoGetCurrentGeoLocation && mode.value === TransactionEditPageMode.Add
        && !geoLocationStatus.value && !transaction.value.geoLocation) {
        updateGeoLocation(false);
    }
}

function onPageBeforeOut(): void {
    if (submitted.value || pageTypeAndMode?.type !== TransactionEditPageType.Transaction || mode.value !== TransactionEditPageMode.Add || addByTemplateId.value || duplicateFromId.value) {
        return;
    }

    if (settingsStore.appSettings.autoSaveTransactionDraft === 'confirmation') {
        if (transactionsStore.isTransactionDraftModified(transaction.value, query['categoryId'], query['accountId'], query['tagIds'])) {
            showConfirm('Do you want to save this transaction draft?', () => {
                transactionsStore.saveTransactionDraft(transaction.value, query['categoryId'], query['accountId'], query['tagIds']);
            }, () => {
                transactionsStore.clearTransactionDraft();
            });
        } else {
            transactionsStore.clearTransactionDraft();
        }
    } else if (settingsStore.appSettings.autoSaveTransactionDraft === 'enabled') {
        transactionsStore.saveTransactionDraft(transaction.value, query['categoryId'], query['accountId'], query['tagIds']);
    }
}

init();
</script>

<style>
.category-separate-icon.icon {
    margin-left: 5px;
    margin-right: 5px;
    font-size: var(--ebk-category-separate-icon-font-size);
    line-height: 16px;
    color: var(--f7-color-gray-tint);
}

.transaction-edit-amount {
    line-height: 53px;
}

.transaction-edit-amount .item-title {
    font-weight: bolder;
}

.transaction-edit-amount .item-header {
    padding-top: calc(var(--f7-typography-padding) / 2);
}

.transaction-edit-datetime .item-title {
    width: 100%;
}

.transaction-edit-datetime .item-title > .item-header > .transaction-edit-datetime-header {
    display: block;
    width: 100%;
}

.transaction-edit-datetime .item-title > .transaction-edit-datetime-title {
    display: flex;
    width: 100%;
}

.transaction-edit-datetime .item-title > .transaction-edit-datetime-title > .transaction-edit-datetime-time {
    flex-grow: 1;
    overflow: hidden;
    text-overflow: ellipsis;
}

.transaction-edit-timezone-name {
    padding-left: 4px;
}

.transaction-edit-tag {
    --f7-chip-bg-color: var(--ebk-transaction-tag-chip-bg-color);
    margin-right: 4px;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
}

.chip.transaction-edit-tag .chip-media+.chip-label {
    margin-left: 0;
}

.chip.transaction-edit-tag .chip-media i.icon {
    font-size: calc(var(--f7-chip-media-size) - 12px);
    height: calc(var(--f7-chip-media-size) - 12px);
}

.transaction-pictures {
    height: var(--ebk-transaction-picture-size);
}

.transaction-picture-container,
.transaction-picture {
    width: var(--ebk-transaction-picture-size);
    height: var(--ebk-transaction-picture-size);
}

.transaction-picture .transaction-picture-control-backdrop {
    width: 100%;
    height: 100%;
    position: absolute;
    z-index: 10;
    background-color: rgba(0, 0, 0, 0.4);
    border-radius: 8px;
}

.transaction-picture .picture-control-icon {
    z-index: 15;
    font-size: var(--ebk-transaction-picture-add-icon-size);
}

.transaction-picture .picture-remove-icon {
    background-color: transparent;
    color: rgba(255, 255, 255, 0.8);
    font-size: var(--ebk-transaction-picture-remove-icon-size);
}

.transaction-picture > img {
    object-fit: cover;
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    border-radius: 8px;
}

.transaction-picture-add {
    width: calc(var(--ebk-transaction-picture-size) - 2px);
    height: calc(var(--ebk-transaction-picture-size) - 4px);
    border: 2px dashed #ccc;
    border-radius: 8px;
}
</style>
