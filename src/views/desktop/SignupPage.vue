<template>
    <div class="layout-wrapper">
        <router-link to="/">
            <div class="auth-logo d-flex align-start gap-x-3">
                <img alt="logo" class="login-page-logo" :src="APPLICATION_LOGO_PATH" />
                <h1 class="font-weight-medium leading-normal text-2xl">{{ tt('global.app.title') }}</h1>
            </div>
        </router-link>
        <v-row no-gutters class="auth-wrapper">
            <v-col cols="12" md="4" class="d-none d-md-flex align-center justify-center position-relative">
                <div class="d-flex auth-img-footer" v-if="!isDarkMode">
                    <v-img src="img/desktop/background.svg"/>
                </div>
                <div class="d-flex auth-img-footer" v-if="isDarkMode">
                    <v-img src="img/desktop/background-dark.svg"/>
                </div>
                <div class="d-flex align-center justify-center w-100 pt-10">
                    <v-img max-width="320px" src="img/desktop/people2.svg"/>
                </div>
            </v-col>
            <v-col cols="12" md="8" class="auth-card d-flex align-center justify-center pa-10">
                <v-card variant="flat" class="mt-12 mt-sm-0 pt-sm-12 pt-md-0">
                    <steps-bar min-width="700" :steps="allSteps" :current-step="currentStep" @step:change="switchToTab" />

                    <v-window class="mt-5 disable-tab-transition" style="max-width: 700px" v-model="currentStep">
                        <v-form>
                            <v-window-item value="basicSetting">
                                <h4 class="text-h4 mb-1">{{ tt('Basic Information') }}</h4>
                                <p class="text-sm mt-2 mb-5">
                                    <span>{{ tt('Already have an account?') }}</span>
                                    <router-link class="ml-1" to="/login">{{ tt('Click here to log in') }}</router-link>
                                </p>
                                <v-row>
                                    <v-col cols="12" md="6">
                                        <v-text-field
                                            type="text"
                                            autocomplete="username"
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
                                        <language-select :disabled="submitting || navigateToHomePage"
                                                         :label="languageTitle"
                                                         :placeholder="languageTitle"
                                                         :use-model-value="true" v-model="currentLocale" />
                                    </v-col>
                                </v-row>

                                <v-row>
                                    <v-col cols="12" md="6">
                                        <currency-select :disabled="submitting || navigateToHomePage"
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
                            </v-window-item>

                            <v-window-item value="presetCategories" class="signup-preset-categories">
                                <h4 class="text-h4 mb-1">{{ tt('Preset Categories') }}</h4>
                                <p class="text-sm mt-2 mb-2">{{ tt('Set whether to use preset transaction categories') }}</p>

                                <v-row>
                                    <v-col cols="12" sm="6">
                                        <v-switch class="mb-2" :disabled="submitting || navigateToHomePage"
                                                  :label="tt('Use Preset Transaction Categories')"
                                                  v-model="usePresetCategories"/>
                                    </v-col>
                                    <v-col cols="12" sm="6" class="text-right-sm">
                                        <language-select-button :disabled="submitting || navigateToHomePage"
                                                                :use-model-value="true" v-model="currentLocale" />
                                    </v-col>
                                </v-row>

                                <div class="overflow-y-auto px-3" :class="{ 'disabled': !usePresetCategories || submitting || navigateToHomePage }" style="max-height: 323px">
                                    <v-row :key="categoryType" v-for="(categories, categoryType) in allPresetCategories">
                                        <v-col cols="12" md="12">
                                            <h4 class="mb-3">{{ getCategoryTypeName(categoryType) }}</h4>

                                            <v-expansion-panels class="border rounded" variant="accordion" multiple>
                                                <v-expansion-panel :key="idx" v-for="(category, idx) in categories">
                                                    <v-expansion-panel-title class="py-0">
                                                        <ItemIcon icon-type="category" :icon-id="category.icon" :color="category.color"></ItemIcon>
                                                        <span class="ml-3">{{ category.name }}</span>
                                                    </v-expansion-panel-title>
                                                    <v-expansion-panel-text v-if="category.subCategories.length">
                                                        <v-list rounded density="comfortable" class="pa-0">
                                                            <template :key="subIdx"
                                                                      v-for="(subCategory, subIdx) in category.subCategories">
                                                                <v-list-item>
                                                                    <template #prepend>
                                                                        <ItemIcon icon-type="category" :icon-id="subCategory.icon" :color="subCategory.color"></ItemIcon>
                                                                    </template>
                                                                    <span class="ml-3">{{ subCategory.name }}</span>
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
                                <h4 class="text-h4 mb-1">{{ tt('Registration Completed') }}</h4>
                                <p class="my-5">{{ finalResultMessage }}</p>
                            </v-window-item>
                        </v-form>
                    </v-window>

                    <div class="d-flex justify-sm-space-between gap-4 flex-wrap justify-center mt-5">
                        <v-btn :color="(currentStep === 'basicSetting' || currentStep === 'finalResult') ? 'default' : 'primary'"
                               :disabled="currentStep === 'basicSetting' || currentStep === 'finalResult' || submitting || navigateToHomePage"
                               :prepend-icon="mdiArrowLeft"
                               @click="switchToPreviousTab">{{ tt('Previous') }}</v-btn>
                        <v-btn color="primary"
                               :disabled="submitting || navigateToHomePage"
                               :append-icon="mdiArrowRight"
                               @click="switchToNextTab"
                               v-if="currentStep === 'basicSetting'">{{ tt('Next') }}</v-btn>
                        <v-btn color="teal"
                               :disabled="submitting || navigateToHomePage"
                               :append-icon="!submitting ? mdiCheck : undefined"
                               @click="submit"
                               v-if="currentStep === 'presetCategories'">
                            {{ tt('Submit') }}
                            <v-progress-circular indeterminate size="22" class="ml-2" v-if="submitting"></v-progress-circular>
                        </v-btn>
                        <v-btn :append-icon="mdiArrowRight"
                               @click="navigateToLogin"
                               v-if="currentStep === 'finalResult'">{{ tt('Continue') }}</v-btn>
                    </div>
                </v-card>
            </v-col>
        </v-row>

        <snack-bar ref="snackbar" @update:show="onSnackbarShowStateChanged" />
    </div>
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';
import type { StepBarItem } from '@/components/desktop/StepsBar.vue';

import { ref, computed, useTemplateRef } from 'vue';
import { useRouter } from 'vue-router';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';
import { useSignupPageBase } from '@/views/base/SignupPageBase.ts';

import { useRootStore } from '@/stores/index.ts';

import type { PartialRecord, TypeAndDisplayName } from '@/core/base.ts';
import { type LocalizedPresetCategory, CategoryType } from '@/core/category.ts';
import { ThemeType } from '@/core/theme.ts';
import { APPLICATION_LOGO_PATH } from '@/consts/asset.ts';

import { categorizedArrayToPlainArray } from '@/lib/common.ts';
import { isUserLogined } from '@/lib/userstate.ts';

import {
    mdiArrowLeft,
    mdiArrowRight,
    mdiCheck
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const router = useRouter();
const theme = useTheme();

const { tt, getAllWeekDays, getAllTransactionDefaultCategories } = useI18n();

const {
    user,
    submitting,
    languageTitle,
    currentLocale,
    inputEmptyProblemMessage,
    inputInvalidProblemMessage,
    getCategoryTypeName,
    doAfterSignupSuccess
} = useSignupPageBase();

const rootStore = useRootStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const currentStep = ref<string>('basicSetting');
const usePresetCategories = ref<boolean>(false);
const finalResultMessage = ref<string | null>(null);
const navigateToHomePage = ref<boolean>(false);

const allWeekDays = computed<TypeAndDisplayName[]>(() => getAllWeekDays());
const allPresetCategories = computed<PartialRecord<CategoryType, LocalizedPresetCategory[]>>(() => getAllTransactionDefaultCategories(0, currentLocale.value));
const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

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
            snackbar.value?.showMessage(problemMessage);
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
        snackbar.value?.showMessage(problemMessage);
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
                snackbar.value?.showMessage('You have been successfully registered');
                navigateToHomePage.value = true;
            }

            return;
        }

        doAfterSignupSuccess(response);
        submitting.value = false;

        if (usePresetCategories.value && !response.presetCategoriesSaved) {
            snackbar.value?.showMessage('You have been successfully registered, but there was an failure when adding preset categories. You can re-add preset categories in settings page anytime.');
        } else {
            snackbar.value?.showMessage('You have been successfully registered');
            router.replace('/');
        }

        navigateToHomePage.value = true;
    }).catch(error => {
        submitting.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function navigateToLogin(): void {
    router.push('/');
}

function onSnackbarShowStateChanged(newValue: boolean): void {
    if (!newValue && navigateToHomePage.value) {
        router.replace('/');
    }
}
</script>

<style>
.signup-preset-categories .v-expansion-panel-text__wrapper {
    padding: 0 0 0 20px;
}
</style>
