<template>
    <div class="layout-wrapper">
        <router-link to="/">
            <div class="auth-logo d-flex align-start gap-x-3">
                <img alt="logo" class="login-page-logo" :src="ezBookkeepingLogoPath" />
                <h1 class="font-weight-medium leading-normal text-2xl">{{ $t('global.app.title') }}</h1>
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
                    <StepsBar min-width="700" :steps="[
                        {
                            'name': 'basicSetting',
                            'title': $t('User Information'),
                            'subTitle': $t('Basic Information')
                        },
                        {
                            'name': 'presetCategories',
                            'title': $t('Transaction Categories'),
                            'subTitle': $t('Preset Categories')
                        }
                    ]" :current-step="currentStep" @step:change="switchToTab" />

                    <v-window class="mt-5 disable-tab-transition" style="max-width: 700px" v-model="currentStep">
                        <v-form>
                            <v-window-item value="basicSetting">
                                <h5 class="text-h5 mb-1">{{ $t('Basic Information') }}</h5>
                                <p class="text-sm mb-5">
                                    <span>{{ $t('Already have an account?') }}</span>
                                    <router-link class="ml-1" to="/login">{{ $t('Click here to log in') }}</router-link>
                                </p>
                                <v-row>
                                    <v-col cols="12" md="6">
                                        <v-text-field
                                            type="text"
                                            autocomplete="username"
                                            clearable
                                            :disabled="submitting"
                                            :label="$t('Username')"
                                            :placeholder="$t('Your username')"
                                            v-model="user.username"
                                        />
                                    </v-col>

                                    <v-col cols="12" md="6">
                                        <v-text-field
                                            type="text"
                                            autocomplete="nickname"
                                            clearable
                                            :disabled="submitting"
                                            :label="$t('Nickname')"
                                            :placeholder="$t('Your nickname')"
                                            v-model="user.nickname"
                                        />
                                    </v-col>
                                </v-row>
                                <v-row>
                                    <v-col cols="12" md="12">
                                        <v-text-field
                                            type="email"
                                            autocomplete="email"
                                            clearable
                                            :disabled="submitting"
                                            :label="$t('E-mail')"
                                            :placeholder="$t('Your email address')"
                                            v-model="user.email"
                                        />
                                    </v-col>
                                </v-row>
                                <v-row>
                                    <v-col cols="12" md="6">
                                        <v-text-field
                                            autocomplete="new-password"
                                            clearable
                                            :disabled="submitting"
                                            :label="$t('Password')"
                                            :placeholder="$t('Your password, at least 6 characters')"
                                            :type="isPasswordVisible ? 'text' : 'password'"
                                            :append-inner-icon="isPasswordVisible ? icons.eyeSlash : icons.eye"
                                            v-model="user.password"
                                            @click:append-inner="isPasswordVisible = !isPasswordVisible"
                                        />
                                    </v-col>
                                    <v-col cols="12" md="6">
                                        <v-text-field
                                            autocomplete="new-password"
                                            clearable
                                            :disabled="submitting"
                                            :label="$t('Confirmation Password')"
                                            :placeholder="$t('Re-enter the password')"
                                            :type="isConfirmPasswordVisible ? 'text' : 'password'"
                                            :append-inner-icon="isConfirmPasswordVisible ? icons.eyeSlash : icons.eye"
                                            v-model="user.confirmPassword"
                                            @click:append-inner="isConfirmPasswordVisible = !isConfirmPasswordVisible"
                                        />
                                    </v-col>
                                </v-row>

                                <v-row>
                                    <v-col cols="12" md="12">
                                        <v-select
                                            item-title="displayName"
                                            item-value="code"
                                            :disabled="submitting"
                                            :label="$t('Language')"
                                            :placeholder="$t('Language')"
                                            :items="allLanguages"
                                            v-model="currentLocale"
                                        />
                                    </v-col>
                                </v-row>

                                <v-row>
                                    <v-col cols="12" md="6">
                                        <v-autocomplete
                                            item-title="displayName"
                                            item-value="code"
                                            auto-select-first
                                            :disabled="submitting"
                                            :label="$t('Default Currency')"
                                            :placeholder="$t('Default Currency')"
                                            :items="allCurrencies"
                                            :no-data-text="$t('No results')"
                                            v-model="user.defaultCurrency"
                                        />
                                    </v-col>

                                    <v-col cols="12" md="6">
                                        <v-select
                                            item-title="displayName"
                                            item-value="type"
                                            :disabled="submitting"
                                            :label="$t('First Day of Week')"
                                            :placeholder="$t('First Day of Week')"
                                            :items="allWeekDays"
                                            v-model="user.firstDayOfWeek"
                                        />
                                    </v-col>
                                </v-row>
                            </v-window-item>

                            <v-window-item value="presetCategories" class="signup-preset-categories">
                                <h5 class="text-h5 mb-1">{{ $t('Preset Categories') }}</h5>
                                <p class="text-sm mb-5">{{ $t('Set Whether You Use The Preset Transaction Categories') }}</p>

                                <v-row class="mb-5">
                                    <v-col cols="12" sm="6">
                                        <v-switch inset
                                                  :disabled="submitting"
                                                  :label="$t('Use Preset Transaction Categories')"
                                                  v-model="usePresetCategories"/>
                                    </v-col>
                                    <v-col cols="12" sm="6" class="text-right-sm">
                                        <v-menu location="bottom">
                                            <template #activator="{ props }">
                                                <v-btn variant="text"
                                                       :disabled="submitting"
                                                       v-bind="props">{{ currentLanguageName }}</v-btn>
                                            </template>
                                            <v-list>
                                                <v-list-item v-for="lang in allLanguages" :key="lang.code">
                                                    <v-list-item-title
                                                        class="cursor-pointer"
                                                        @click="currentLocale = lang.code">
                                                        {{ lang.displayName }}
                                                    </v-list-item-title>
                                                </v-list-item>
                                            </v-list>
                                        </v-menu>
                                    </v-col>
                                </v-row>

                                <div class="overflow-y-auto px-3" :class="{ 'disabled': !usePresetCategories || submitting }" style="max-height: 323px">
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
                        </v-form>
                    </v-window>

                    <div class="d-flex justify-sm-space-between gap-4 flex-wrap justify-center mt-5">
                        <v-btn :color="currentStep === 'basicSetting' ? 'default' : 'primary'"
                               :disabled="currentStep === 'basicSetting' || submitting"
                               :prepend-icon="icons.previous"
                               @click="switchToPreviousTab">{{ $t('Previous') }}</v-btn>
                        <v-btn :color="currentStep === 'presetCategories' ? 'secondary' : 'primary'"
                               :disabled="currentStep === 'presetCategories' || submitting"
                               :append-icon="icons.next"
                               @click="switchToNextTab"
                               v-if="currentStep !== 'presetCategories'">{{ $t('Next') }}</v-btn>
                        <v-btn color="expense"
                               :disabled="submitting"
                               :append-icon="!submitting ? icons.submit : null"
                               @click="submit"
                               v-if="currentStep === 'presetCategories'">
                            {{ $t('Submit') }}
                            <v-progress-circular indeterminate size="24" class="ml-2" v-if="submitting"></v-progress-circular>
                        </v-btn>
                    </div>
                </v-card>
            </v-col>
        </v-row>

        <snack-bar ref="snackbar" />
    </div>
</template>

<script>
import { useTheme } from 'vuetify';

import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.js';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import assetConstants from '@/consts/asset.js';
import categoryConstants from '@/consts/category.js';
import { categoriedArrayToPlainArray } from '@/lib/common.js';

import {
    mdiArrowLeft,
    mdiArrowRight,
    mdiCheck,
    mdiEyeOffOutline,
    mdiEyeOutline
} from '@mdi/js';

export default {
    data() {
        const self = this;
        const settingsStore = useSettingsStore();

        return {
            user: {
                username: '',
                password: '',
                confirmPassword: '',
                email: '',
                nickname: '',
                language: self.$locale.getCurrentLanguageCode(),
                defaultCurrency: settingsStore.localeDefaultSettings.currency,
                firstDayOfWeek: settingsStore.localeDefaultSettings.firstDayOfWeek,
            },
            currentStep: 'basicSetting',
            isPasswordVisible: false,
            isConfirmPasswordVisible: false,
            submitting: false,
            usePresetCategories: false,
            icons: {
                previous: mdiArrowLeft,
                next: mdiArrowRight,
                submit: mdiCheck,
                eye: mdiEyeOutline,
                eyeSlash: mdiEyeOffOutline
            }
        };
    },
    computed: {
        ...mapStores(useRootStore, useSettingsStore, useTransactionCategoriesStore, useExchangeRatesStore),
        ezBookkeepingLogoPath() {
            return assetConstants.ezBookkeepingLogoPath;
        },
        allLanguages() {
            return this.$locale.getAllLanguageInfoArray(false);
        },
        allCurrencies() {
            return this.$locale.getAllCurrencies();
        },
        allWeekDays() {
            return this.$locale.getAllWeekDays();
        },
        allPresetCategories() {
            return this.$locale.getAllTransactionDefaultCategories(0, this.currentLocale);
        },
        currentLocale: {
            get: function () {
                return this.$locale.getCurrentLanguageCode();
            },
            set: function (value) {
                const isCurrencyDefault = this.user.defaultCurrency === this.settingsStore.localeDefaultSettings.currency;
                const isFirstWeekDayDefault = this.user.firstDayOfWeek === this.settingsStore.localeDefaultSettings.firstDayOfWeek;

                this.user.language = value;

                const localeDefaultSettings = this.$locale.setLanguage(value);
                this.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

                if (isCurrencyDefault) {
                    this.user.defaultCurrency = this.settingsStore.localeDefaultSettings.currency;
                }

                if (isFirstWeekDayDefault) {
                    this.user.firstDayOfWeek = this.settingsStore.localeDefaultSettings.firstDayOfWeek;
                }
            }
        },
        isDarkMode() {
            return this.globalTheme.global.name.value === 'dark';
        },
        currentLanguageName() {
            const languageInfo = this.$locale.getLanguageInfo(this.currentLocale);

            if (!languageInfo) {
                return '';
            }

            return languageInfo.displayName;
        },
        inputIsEmpty() {
            return !!this.inputEmptyProblemMessage;
        },
        inputIsInvalid() {
            return !!this.inputInvalidProblemMessage;
        },
        inputEmptyProblemMessage() {
            if (!this.user.username) {
                return 'Username cannot be empty';
            } else if (!this.user.nickname) {
                return 'Nickname cannot be empty';
            } else if (!this.user.email) {
                return 'Email address cannot be empty';
            } else if (!this.user.password) {
                return 'Password cannot be empty';
            } else if (!this.user.confirmPassword) {
                return 'Confirmation password cannot be empty';
            } else if (!this.user.defaultCurrency) {
                return 'Default currency cannot be empty';
            } else {
                return null;
            }
        },
        inputInvalidProblemMessage() {
            if (this.user.password && this.user.confirmPassword && this.user.password !== this.user.confirmPassword) {
                return 'Password and confirmation password do not match';
            } else {
                return null;
            }
        }
    },
    setup() {
        const theme = useTheme();

        return {
            globalTheme: theme
        };
    },
    methods: {
        switchToTab(tabName) {
            if (this.submitting) {
                return;
            }

            if (tabName === 'basicSetting') {
                this.currentStep = 'basicSetting';
            } else if (tabName === 'presetCategories') {
                const problemMessage = this.inputEmptyProblemMessage || this.inputInvalidProblemMessage;

                if (problemMessage) {
                    this.$refs.snackbar.showMessage(problemMessage);
                    return;
                }

                this.currentStep = 'presetCategories';
            }
        },
        switchToPreviousTab() {
            this.switchToTab('basicSetting');
        },
        switchToNextTab() {
            this.switchToTab('presetCategories');
        },
        submit() {
            const self = this;

            const problemMessage = self.inputEmptyProblemMessage || self.inputInvalidProblemMessage;

            if (problemMessage) {
                self.$refs.snackbar.showMessage(problemMessage);
                return;
            }

            self.submitting = true;

            let submitCategories = [];

            if (self.usePresetCategories) {
                submitCategories = categoriedArrayToPlainArray(self.allPresetCategories);
            }

            self.rootStore.register({
                user: self.user
            }).then(response => {
                if (!self.$user.isUserLogined()) {
                    self.submitting = false;

                    if (self.usePresetCategories) {
                        self.$refs.snackbar.showMessage('You have been successfully registered, but something wrong with adding preset categories. You can re-add preset categories in settings page anytime.');
                    } else {
                        self.$refs.snackbar.showMessage('You have been successfully registered');
                    }

                    self.$router.replace('/');
                    return;
                }

                if (response.user) {
                    const localeDefaultSettings = self.$locale.setLanguage(response.user.language);
                    self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
                }

                if (self.settingsStore.appSettings.autoUpdateExchangeRatesData) {
                    self.exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
                }

                if (!self.usePresetCategories) {
                    self.submitting = false;

                    self.$refs.snackbar.showMessage('You have been successfully registered');
                    self.$router.replace('/');
                    return;
                }

                self.transactionCategoriesStore.addCategories({
                    categories: submitCategories
                }).then(() => {
                    self.submitting = false;

                    self.$refs.snackbar.showMessage('You have been successfully registered');
                    self.$router.replace('/');
                }).catch(() => {
                    self.submitting = false;

                    self.$refs.snackbar.showMessage('You have been successfully registered, but something wrong with adding preset categories. You can re-add preset categories in settings page anytime.');
                    self.$router.replace('/');
                });
            }).catch(error => {
                self.submitting = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        getCategoryTypeName(categoryType) {
            switch (categoryType) {
                case categoryConstants.allCategoryTypes.Income.toString():
                    return this.$t('Income Categories');
                case categoryConstants.allCategoryTypes.Expense.toString():
                    return this.$t('Expense Categories');
                case categoryConstants.allCategoryTypes.Transfer.toString():
                    return this.$t('Transfer Categories');
                default:
                    return this.$t('Transaction Categories');
            }
        }
    }
}
</script>

<style>
.signup-preset-categories .v-expansion-panel-text__wrapper {
    padding: 0 0 0 20px;
}
</style>
