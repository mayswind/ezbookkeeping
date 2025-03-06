<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Sign Up')"></f7-nav-title>
            <f7-nav-right>
                <f7-link :class="{ 'disabled': inputIsEmpty || submitting }" :text="tt('Submit')" @click="submit"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list form strong inset dividers class="margin-vertical">
            <f7-list-input
                type="text"
                autocomplete="username"
                clear-button
                :label="tt('Username')"
                :placeholder="tt('Your username')"
                v-model:value="user.username"
            ></f7-list-input>

            <f7-list-input
                type="password"
                autocomplete="new-password"
                clear-button
                :label="tt('Password')"
                :placeholder="tt('Your password, at least 6 characters')"
                v-model:value="user.password"
            ></f7-list-input>

            <f7-list-input
                type="password"
                autocomplete="new-password"
                clear-button
                :label="tt('Confirm Password')"
                :placeholder="tt('Re-enter the password')"
                v-model:value="user.confirmPassword"
            ></f7-list-input>

            <f7-list-input
                type="email"
                autocomplete="email"
                clear-button
                :label="tt('E-mail')"
                :placeholder="tt('Your email address')"
                v-model:value="user.email"
            ></f7-list-input>

            <f7-list-input
                type="text"
                autocomplete="nickname"
                clear-button
                :label="tt('Nickname')"
                :placeholder="tt('Your nickname')"
                v-model:value="user.nickname"
            ></f7-list-input>

            <f7-list-item class="ebk-list-item-error-info" v-if="inputIsInvalid" :footer="tt(inputInvalidProblemMessage)"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical">
            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                :key="currentLocale + '_lang'"
                :header="languageTitle"
                :title="currentLanguageName"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: tt('Language'), searchbarDisableText: tt('Cancel'), appendSearchbarNotFound: tt('No results'), pageTitle: languageTitle, popupCloseLinkText: tt('Done') }"
            >
                <select v-model="currentLocale">
                    <option :value="lang.languageTag"
                            :key="lang.languageTag"
                            v-for="lang in allLanguages">{{ lang.nativeDisplayName }}</option>
                </select>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                :key="currentLocale + '_currency'"
                :header="tt('Default Currency')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: tt('Currency Name'), searchbarDisableText: tt('Cancel'), appendSearchbarNotFound: tt('No results'), pageTitle: tt('Default Currency'), popupCloseLinkText: tt('Done') }"
            >
                <template #title>
                    <f7-block class="no-padding no-margin">
                        <span>{{ getCurrencyName(user.defaultCurrency) }}&nbsp;</span>
                        <small class="smaller">{{ user.defaultCurrency }}</small>
                    </f7-block>
                </template>
                <select autocomplete="transaction-currency" v-model="user.defaultCurrency">
                    <option :value="currency.currencyCode"
                            :key="currency.currencyCode"
                            v-for="currency in allCurrencies">{{ currency.displayName }}</option>
                </select>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                :key="currentLocale + '_firstDayOfWeek'"
                :header="tt('First Day of Week')"
                :title="currentDayOfWeekName"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: tt('Date'), searchbarDisableText: tt('Cancel'), appendSearchbarNotFound: tt('No results'), pageTitle: tt('First Day of Week'), popupCloseLinkText: tt('Done') }"
            >
                <select v-model="user.firstDayOfWeek">
                    <option :value="weekDay.type"
                            :key="weekDay.type"
                            v-for="weekDay in allWeekDays">{{ weekDay.displayName }}</option>
                </select>
            </f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical">
            <f7-list-item :title="tt('Use preset transaction categories')" link="#" @click="showPresetCategories = true">
                <f7-toggle :checked="usePresetCategories" @toggle:change="usePresetCategories = $event"></f7-toggle>
            </f7-list-item>
        </f7-list>

        <f7-popup push :close-on-escape="false" :opened="showPresetCategories"
                  @popup:closed="showPresetCategories = false">
            <f7-page>
                <f7-navbar>
                    <f7-nav-left>
                        <f7-link popup-close :text="tt('Back')"></f7-link>
                    </f7-nav-left>
                    <f7-nav-title :title="tt('Preset Categories')"></f7-nav-title>
                    <f7-nav-right>
                        <f7-link icon-f7="ellipsis" @click="showPresetCategoriesMoreActionSheet = true"></f7-link>
                        <f7-link close @click="usePresetCategories = true; showPresetCategories = false" v-if="!usePresetCategories">{{ tt('Enable') }}</f7-link>
                        <f7-link close @click="usePresetCategories = false; showPresetCategories = false" v-if="usePresetCategories">{{ tt('Disable') }}</f7-link>
                    </f7-nav-right>
                </f7-navbar>
                <f7-block class="no-padding no-margin"
                          :key="categoryType" v-for="(categories, categoryType) in allPresetCategories">
                    <f7-block-title class="margin-top margin-horizontal">{{ getCategoryTypeName(categoryType) }}</f7-block-title>
                    <f7-list strong inset dividers v-if="showPresetCategories">
                        <f7-list-item :title="category.name"
                                      :accordion-item="!!category.subCategories.length"
                                      :key="idx"
                                      v-for="(category, idx) in categories">
                            <template #media>
                                <ItemIcon icon-type="category" :icon-id="category.icon" :color="category.color"></ItemIcon>
                            </template>

                            <f7-accordion-content v-if="category.subCategories.length" class="padding-left">
                                <f7-list>
                                    <f7-list-item :title="subCategory.name"
                                                  :key="subIdx"
                                                  v-for="(subCategory, subIdx) in category.subCategories">
                                        <template #media>
                                            <ItemIcon icon-type="category" :icon-id="subCategory.icon" :color="subCategory.color"></ItemIcon>
                                        </template>
                                    </f7-list-item>
                                </f7-list>
                            </f7-accordion-content>
                        </f7-list-item>
                    </f7-list>
                </f7-block>
            </f7-page>

            <f7-actions close-by-outside-click close-on-escape :opened="showPresetCategoriesMoreActionSheet" @actions:closed="showPresetCategoriesMoreActionSheet = false">
                <f7-actions-group>
                    <f7-actions-button @click="showPresetCategoriesChangeLocaleSheet = true">{{ tt('Change Language') }}</f7-actions-button>
                </f7-actions-group>
                <f7-actions-group>
                    <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
                </f7-actions-group>
            </f7-actions>

            <list-item-selection-sheet value-type="item"
                                       value-field="languageTag"
                                       title-field="nativeDisplayName"
                                       after-field="displayName"
                                       :items="allLanguages"
                                       v-model:show="showPresetCategoriesChangeLocaleSheet"
                                       v-model="currentLocale">
            </list-item-selection-sheet>
        </f7-popup>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Router } from 'framework7/types';

import type { LanguageOption } from '@/locales/index.ts';
import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';
import { useSignupPageBase } from '@/views/base/SignupPageBase.ts';

import { useRootStore } from '@/stores/index.ts';

import type { PartialRecord, TypeAndDisplayName } from '@/core/base.ts';
import type { LocalizedCurrencyInfo } from '@/core/currency.ts';
import { type LocalizedPresetCategory, CategoryType } from '@/core/category.ts';

import { findDisplayNameByType, categorizedArrayToPlainArray } from '@/lib/common.ts';
import { isUserLogined } from '@/lib/userstate.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt, getAllLanguageOptions, getAllCurrencies, getAllWeekDays, getAllTransactionDefaultCategories, getCurrencyName } = useI18n();
const { showAlert, showToast } = useI18nUIComponents();

const {
    user,
    submitting,
    languageTitle,
    currentLocale,
    currentLanguageName,
    inputEmptyProblemMessage,
    inputInvalidProblemMessage,
    inputIsEmpty,
    inputIsInvalid,
    getCategoryTypeName,
    doAfterSignupSuccess
} = useSignupPageBase();

const rootStore = useRootStore();

const usePresetCategories = ref<boolean>(false);
const showPresetCategories = ref<boolean>(false);
const showPresetCategoriesMoreActionSheet = ref<boolean>(false);
const showPresetCategoriesChangeLocaleSheet = ref<boolean>(false);

const allLanguages = computed<LanguageOption[]>(() => getAllLanguageOptions(false));
const allCurrencies = computed<LocalizedCurrencyInfo[]>(() => getAllCurrencies());
const allWeekDays = computed<TypeAndDisplayName[]>(() => getAllWeekDays());
const allPresetCategories = computed<PartialRecord<CategoryType, LocalizedPresetCategory[]>>(() => getAllTransactionDefaultCategories(0, currentLocale.value));
const currentDayOfWeekName = computed<string | null>(() => findDisplayNameByType(allWeekDays.value, user.value.firstDayOfWeek));

function submit(): void {
    const router = props.f7router;

    const problemMessage = inputEmptyProblemMessage.value || inputInvalidProblemMessage.value;

    if (problemMessage) {
        showAlert(problemMessage);
        return;
    }

    submitting.value = true;
    showLoading(() => submitting.value);

    let presetCategories: LocalizedPresetCategory[] = [];

    if (usePresetCategories.value) {
        presetCategories = categorizedArrayToPlainArray(allPresetCategories.value);
    }

    rootStore.register({
        user: user.value,
        presetCategories: presetCategories
    }).then(response => {
        if (!isUserLogined()) {
            submitting.value = false;
            hideLoading();

            if (usePresetCategories.value && !response.presetCategoriesSaved) {
                showToast('You have been successfully registered, but there was an failure when adding preset categories. You can re-add preset categories in settings page anytime.', 5000);
            } else if (response.needVerifyEmail) {
                showToast('You have been successfully registered. An account activation link has been sent to your email address, please activate your account first.', 5000);
            } else {
                showToast('You have been successfully registered');
            }

            router.navigate('/');
            return;
        }

        doAfterSignupSuccess(response);
        submitting.value = false;
        hideLoading();

        if (usePresetCategories.value && !response.presetCategoriesSaved) {
            showToast('You have been successfully registered, but there was an failure when adding preset categories. You can re-add preset categories in settings page anytime.');
        } else {
            showToast('You have been successfully registered');
        }

        router.navigate('/');
    }).catch(error => {
        submitting.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}
</script>
