<template>
    <div class="layout-wrapper">
        <router-link to="/">
            <div class="auth-logo d-flex align-start gap-x-3">
                <img alt="logo" class="login-page-logo" :src="APPLICATION_LOGO_PATH" />
                <span class="auth-app-title">{{ tt('global.app.title') }}</span>
            </div>
        </router-link>
        <v-row no-gutters class="auth-wrapper">
            <v-col cols="12" md="4" class="auth-image-background d-none d-md-flex align-center justify-center position-relative">
                <auth-illustration variant="signup" />
            </v-col>
            <v-col cols="12" md="8" class="auth-card d-flex align-center justify-center pa-10">
                <v-card variant="flat" class="mt-12 mt-sm-0 pt-sm-12 pt-md-0">
                    <steps-bar min-width="700" :steps="allSteps" :current-step="currentStep" @step:change="switchToTab" />

                    <v-window class="mt-5 mt-md-6 disable-tab-transition" style="max-width: 700px" v-model="currentStep">
                        <v-form>
                            <v-window-item value="basicSetting">
                                <div class="text-headline-small mb-1">{{ tt('Basic Information') }}</div>
                                <div class="auth-message text-body-medium my-3">
                                    <span>{{ tt('Already have an account?') }}</span>
                                    <router-link class="ms-1" to="/login">{{ tt('Click here to log in') }}</router-link>
                                </div>
                                <v-row>
                                    <v-col cols="12" md="6">
                                        <v-text-field
                                            type="text"
                                            autocomplete="username"
                                            autocapitalize="none"
                                            autocorrect="off"
                                            spellcheck="false"
                                            inputmode="email"
                                            :disabled="submitting || navigateToHomePage"
                                            :label="tt('Username')"
                                            :placeholder="tt('Your username')"
                                            v-model="user.username"
                                        />
                                    </v-col>

                                    <v-col cols="12" md="6">
                                        <v-text-field
                                            type="text"
                                            autocomplete="nickname"
                                            :disabled="submitting || navigateToHomePage"
                                            :label="tt('Nickname')"
                                            :placeholder="tt('Your nickname')"
                                            v-model="user.nickname"
                                        />
                                    </v-col>
                                </v-row>
                                <v-row>
                                    <v-col cols="12" md="12">
                                        <v-text-field
                                            type="email"
                                            autocomplete="email"
                                            :disabled="submitting || navigateToHomePage"
                                            :label="tt('E-mail')"
                                            :placeholder="tt('Your email address')"
                                            v-model="user.email"
                                        />
                                    </v-col>
                                </v-row>
                                <v-row>
                                    <v-col cols="12" md="6">
                                        <v-text-field
                                            autocomplete="new-password"
                                            type="password"
                                            :disabled="submitting || navigateToHomePage"
                                            :label="tt('Password')"
                                            :placeholder="tt('Your password, at least 6 characters')"
                                            v-model="user.password"
                                        />
                                    </v-col>
                                    <v-col cols="12" md="6">
                                        <v-text-field
                                            autocomplete="new-password"
                                            type="password"
                                            :disabled="submitting || navigateToHomePage"
                                            :label="tt('Confirm Password')"
                                            :placeholder="tt('Re-enter the password')"
                                            v-model="user.confirmPassword"
                                        />
                                    </v-col>
                                </v-row>

                                <v-row>
                                    <v-col cols="12" md="12">
                                        <language-select :disabled="submitting || settingLanguage || navigateToHomePage"
                                                         :label="languageTitle"
                                                         :placeholder="languageTitle"
                                                         :loading="settingLanguage"
                                                         :use-model-value="true"
                                                         :model-value="user.language"
                                                         @update:model-value="changeLanguage" />
                                    </v-col>
                                </v-row>

                                <v-row>
                                    <v-col cols="12" md="6">
                                        <currency-select :menu-props="{ contentClass: 'select-menu-large' }"
                                                         :disabled="submitting || navigateToHomePage"
                                                         :label="tt('Default Currency')"
                                                         :placeholder="tt('Default Currency')"
                                                         v-model="user.defaultCurrency" />
                                    </v-col>

                                    <v-col cols="12" md="6">
                                        <v-select
                                            item-title="displayName"
                                            item-value="type"
                                            :disabled="submitting || navigateToHomePage"
                                            :label="tt('First Day of Week')"
                                            :placeholder="tt('First Day of Week')"
                                            :items="allWeekDays"
                                            v-model="user.firstDayOfWeek"
                                        />
                                    </v-col>
                                </v-row>

                                <v-row>
                                    <v-col cols="12" md="6">
                                        <v-select
                                            item-title="displayName"
                                            item-value="type"
                                            :disabled="submitting || navigateToHomePage"
                                            :label="tt('Expense Amount Color')"
                                            :placeholder="tt('Expense Amount Color')"
                                            :items="allExpenseAmountColorTypes"
                                            v-model="user.expenseAmountColor"
                                        />
                                    </v-col>

                                    <v-col cols="12" md="6">
                                        <v-select
                                            item-title="displayName"
                                            item-value="type"
                                            :disabled="submitting || navigateToHomePage"
                                            :label="tt('Income Amount Color')"
                                            :placeholder="tt('Income Amount Color')"
                                            :items="allIncomeAmountColorTypes"
                                            v-model="user.incomeAmountColor"
                                        />
                                    </v-col>
                                </v-row>
                            </v-window-item>

                            <v-window-item value="presetCategories">
                                <div class="text-headline-small mb-1">{{ tt('Preset Categories') }}</div>
                                <div class="auth-message text-body-medium my-3">{{ tt('Set whether to use preset transaction categories') }}</div>

                                <v-row>
                                    <v-col cols="12" sm="6">
                                        <v-switch class="mb-2" :disabled="submitting || navigateToHomePage"
                                                  :label="tt('Use Preset Transaction Categories')"
                                                  v-model="usePresetCategories"/>
                                    </v-col>
                                    <v-col cols="12" sm="6" class="text-end-sm">
                                        <language-select-button :disabled="submitting || settingLanguage || navigateToHomePage"
                                                                :loading="settingLanguage"
                                                                :use-model-value="true"
                                                                :model-value="user.language"
                                                                @update:model-value="changeLanguage" />
                                    </v-col>
                                </v-row>

                                <div class="overflow-y-auto pe-1" :class="{ 'disabled': !usePresetCategories || submitting || settingLanguage || navigateToHomePage }" style="max-height: 323px">
                                    <v-row :key="categoryType" v-for="(categories, categoryType) in allPresetCategories">
                                        <v-col cols="12" md="12">
                                            <span class="text-body-large">{{ getCategoryTypeName(parseInt(categoryType)) }}</span>

                                            <v-expansion-panels class="compacted-expansion-panels mt-2" variant="accordion" multiple>
                                                <v-expansion-panel :key="idx" v-for="(category, idx) in categories">
                                                    <v-expansion-panel-title class="py-0 px-4">
                                                        <ItemIcon :icon-type="getCategoryIconType(category.iconType)" :icon-id="category.icon" :color="category.color"></ItemIcon>
                                                        <span class="text-body-medium ms-2">{{ category.name }}</span>
                                                    </v-expansion-panel-title>
                                                    <v-expansion-panel-text v-if="category.subCategories.length">
                                                        <v-list rounded density="comfortable" class="pa-0">
                                                            <template :key="subIdx"
                                                                      v-for="(subCategory, subIdx) in category.subCategories">
                                                                <v-list-item>
                                                                    <template #prepend>
                                                                        <ItemIcon :icon-type="getCategoryIconType(subCategory.iconType)" :icon-id="subCategory.icon" :color="subCategory.color"></ItemIcon>
                                                                    </template>
                                                                    <span class="text-body-medium ms-2">{{ subCategory.name }}</span>
                                                                </v-list-item>
                                                                <v-divider v-if="subIdx !== category.subCategories.length - 1"/>
                                                            </template>
                                                        </v-list>
                                                    </v-expansion-panel-text>
                                                </v-expansion-panel>
                                            </v-expansion-panels>
                                        </v-col>
                                    </v-row>
                                </div>
                            </v-window-item>

                            <v-window-item value="finalResult" v-if="finalResultMessage">
                                <div class="text-headline-large mb-1">{{ tt('Registration Completed') }}</div>
                                <div class="my-3">{{ finalResultMessage }}</div>
                            </v-window-item>
                        </v-form>
                    </v-window>

                    <div class="d-flex justify-sm-space-between gap-4 flex-wrap justify-center mt-5">
                        <v-btn class="button-icon-with-direction"
                               :color="(currentStep === 'basicSetting' || currentStep === 'finalResult') ? 'default' : 'primary'"
                               :disabled="currentStep === 'basicSetting' || currentStep === 'finalResult' || submitting || settingLanguage || navigateToHomePage"
                               :prepend-icon="mdiArrowLeft"
                               @click="switchToPreviousTab">{{ tt('Previous') }}</v-btn>
                        <v-btn class="button-icon-with-direction" color="primary"
                               :disabled="submitting || settingLanguage || navigateToHomePage"
                               :append-icon="mdiArrowRight"
                               @click="switchToNextTab"
                               v-if="currentStep === 'basicSetting'">{{ tt('Next') }}</v-btn>
                        <v-btn color="teal"
                               :disabled="submitting || settingLanguage || navigateToHomePage"
                               :append-icon="!submitting ? mdiCheck : undefined"
                               @click="submit"
                               v-if="currentStep === 'presetCategories'">
                            {{ tt('Submit') }}
                            <v-progress-circular indeterminate size="22" class="ms-2" v-if="submitting"></v-progress-circular>
                        </v-btn>
                        <v-btn class="button-icon-with-direction"
                               :append-icon="mdiArrowRight"
                               @click="navigateToLogin"
                               v-if="currentStep === 'finalResult'">{{ tt('Continue') }}</v-btn>
                    </div>
                </v-card>
            </v-col>
        </v-row>
    </div>
</template>

<script setup lang="ts">
import type { StepBarItem } from '@/components/desktop/StepsBar.vue';

import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';

import { useI18n } from '@/locales/helpers.ts';
import { useSignupPageBase } from '@/views/base/SignupPageBase.ts';

import { useRootStore } from '@/stores/index.ts';
import { useDesktopPageStore } from '@/stores/desktopPage.ts';

import type { TypeAndDisplayName } from '@/core/base.ts';
import { type LocalizedPresetCategory } from '@/core/category.ts';
import { APPLICATION_LOGO_PATH } from '@/consts/asset.ts';

import { categorizedArrayToPlainArray } from '@/lib/common.ts';
import { getCategoryIconType } from '@/lib/icon.ts';
import { isUserLogined } from '@/lib/userstate.ts';
import logger from '@/lib/logger.ts';

import {
    mdiArrowLeft,
    mdiArrowRight,
    mdiCheck
} from '@mdi/js';

const router = useRouter();

const { tt, getAllWeekDays, getAllTransactionDefaultCategories } = useI18n();

const {
    user,
    settingLanguage,
    submitting,
    allExpenseAmountColorTypes,
    allIncomeAmountColorTypes,
    languageTitle,
    inputEmptyProblemMessage,
    inputInvalidProblemMessage,
    updateLanguage,
    getCategoryTypeName,
    doAfterSignupSuccess
} = useSignupPageBase();

const rootStore = useRootStore();
const desktopPageStore = useDesktopPageStore();

const currentStep = ref<string>('basicSetting');
const usePresetCategories = ref<boolean>(false);
const finalResultMessage = ref<string | null>(null);
const navigateToHomePage = ref<boolean>(false);

const allWeekDays = computed<TypeAndDisplayName[]>(() => getAllWeekDays());
const allPresetCategories = computed<Record<string, LocalizedPresetCategory[]>>(() => getAllTransactionDefaultCategories(0, user.value.language));

const allSteps = computed<StepBarItem[]>(() => {
    const allSteps = [
        {
            name: 'basicSetting',
            title: tt('User Information'),
            subTitle: tt('Basic Information')
        },
        {
            name: 'presetCategories',
            title: tt('Transaction Categories'),
            subTitle: tt('Preset Categories')
        }
    ];

    if (finalResultMessage.value) {
        allSteps.push({
            name: 'finalResult',
            title: tt('Complete'),
            subTitle: tt('Registration Completed')
        });
    }

    return allSteps;
});

function switchToTab(tabName: string): void {
    if (submitting.value || currentStep.value === 'finalResult' || navigateToHomePage.value) {
        return;
    }

    if (tabName === 'basicSetting') {
        currentStep.value = 'basicSetting';
    } else if (tabName === 'presetCategories') {
        const problemMessage = inputEmptyProblemMessage.value || inputInvalidProblemMessage.value;

        if (problemMessage) {
            desktopPageStore.addToast(problemMessage);
            return;
        }

        currentStep.value = 'presetCategories';
    }
}

function switchToPreviousTab(): void {
    switchToTab('basicSetting');
}

function switchToNextTab(): void {
    switchToTab('presetCategories');
}

function submit(): void {
    const problemMessage = inputEmptyProblemMessage.value || inputInvalidProblemMessage.value;

    if (problemMessage) {
        desktopPageStore.addToast(problemMessage);
        return;
    }

    navigateToHomePage.value = false;
    submitting.value = true;

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

            if (usePresetCategories.value && !response.presetCategoriesSaved) {
                finalResultMessage.value = tt('You have been successfully registered, but there was an failure when adding preset categories. You can re-add preset categories in settings page anytime.');
                currentStep.value = 'finalResult';
            } else if (response.needVerifyEmail) {
                finalResultMessage.value = tt('You have been successfully registered. An account activation link has been sent to your email address, please activate your account first.');
                currentStep.value = 'finalResult';
            } else {
                desktopPageStore.addToast('You have been successfully registered');
                navigateToHomePage.value = true;
                router.replace('/');
            }

            return;
        }

        doAfterSignupSuccess(response).then(() => {
            submitting.value = false;

            if (usePresetCategories.value && !response.presetCategoriesSaved) {
                desktopPageStore.addToast('You have been successfully registered, but there was an failure when adding preset categories. You can re-add preset categories in settings page anytime.');
            } else {
                desktopPageStore.addToast('You have been successfully registered');
            }

            navigateToHomePage.value = true;
            router.replace('/');
        });
    }).catch(error => {
        submitting.value = false;

        if (!error.processed) {
            desktopPageStore.addErrorToast(error);
        }
    });
}

function changeLanguage(languageTag: string): void {
    updateLanguage(languageTag).catch(error => {
        logger.error('failed to set language', error);
        desktopPageStore.addErrorToast(error);
    });
}

function navigateToLogin(): void {
    router.push('/');
}
</script>
