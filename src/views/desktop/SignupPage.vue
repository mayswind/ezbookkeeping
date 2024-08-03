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
                    <steps-bar min-width="700" :steps="allSteps" :current-step="currentStep" @step:change="switchToTab" />

                    <v-window class="mt-5 disable-tab-transition" style="max-width: 700px" v-model="currentStep">
                        <v-form>
                            <v-window-item value="basicSetting">
                                <h4 class="text-h4 mb-1">{{ $t('Basic Information') }}</h4>
                                <p class="text-sm mt-2 mb-5">
                                    <span>{{ $t('Already have an account?') }}</span>
                                    <router-link class="ml-1" to="/login">{{ $t('Click here to log in') }}</router-link>
                                </p>
                                <v-row>
                                    <v-col cols="12" md="6">
                                        <v-text-field
                                            type="text"
                                            autocomplete="username"
                                            :disabled="submitting || navigateToHomePage"
                                            :label="$t('Username')"
                                            :placeholder="$t('Your username')"
                                            v-model="user.username"
                                        />
                                    </v-col>

                                    <v-col cols="12" md="6">
                                        <v-text-field
                                            type="text"
                                            autocomplete="nickname"
                                            :disabled="submitting || navigateToHomePage"
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
                                            :disabled="submitting || navigateToHomePage"
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
                                            type="password"
                                            :disabled="submitting || navigateToHomePage"
                                            :label="$t('Password')"
                                            :placeholder="$t('Your password, at least 6 characters')"
                                            v-model="user.password"
                                        />
                                    </v-col>
                                    <v-col cols="12" md="6">
                                        <v-text-field
                                            autocomplete="new-password"
                                            type="password"
                                            :disabled="submitting || navigateToHomePage"
                                            :label="$t('Confirm Password')"
                                            :placeholder="$t('Re-enter the password')"
                                            v-model="user.confirmPassword"
                                        />
                                    </v-col>
                                </v-row>

                                <v-row>
                                    <v-col cols="12" md="12">
                                        <v-select
                                            item-title="displayName"
                                            item-value="languageTag"
                                            :disabled="submitting || navigateToHomePage"
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
                                            :disabled="submitting || navigateToHomePage"
                                            :label="$t('Default Currency')"
                                            :placeholder="$t('Default Currency')"
                                            :items="allCurrencies"
                                            :no-data-text="$t('No results')"
                                            v-model="user.defaultCurrency"
                                        >
                                            <template #append-inner>
                                                <small class="text-field-append-text smaller">{{ user.defaultCurrency }}</small>
                                            </template>
                                        </v-autocomplete>
                                    </v-col>

                                    <v-col cols="12" md="6">
                                        <v-select
                                            item-title="displayName"
                                            item-value="type"
                                            :disabled="submitting || navigateToHomePage"
                                            :label="$t('First Day of Week')"
                                            :placeholder="$t('First Day of Week')"
                                            :items="allWeekDays"
                                            v-model="user.firstDayOfWeek"
                                        />
                                    </v-col>
                                </v-row>
                            </v-window-item>

                            <v-window-item value="presetCategories" class="signup-preset-categories">
                                <h4 class="text-h4 mb-1">{{ $t('Preset Categories') }}</h4>
                                <p class="text-sm mt-2 mb-5">{{ $t('Set whether to use preset transaction categories') }}</p>

                                <v-row>
                                    <v-col cols="12" sm="6">
                                        <v-switch :disabled="submitting || navigateToHomePage"
                                                  :label="$t('Use Preset Transaction Categories')"
                                                  v-model="usePresetCategories"/>
                                    </v-col>
                                    <v-col cols="12" sm="6" class="text-right-sm">
                                        <v-menu location="bottom">
                                            <template #activator="{ props }">
                                                <v-btn variant="text"
                                                       :disabled="submitting || navigateToHomePage"
                                                       v-bind="props">{{ currentLanguageName }}</v-btn>
                                            </template>
                                            <v-list>
                                                <v-list-item v-for="lang in allLanguages" :key="lang.languageTag">
                                                    <v-list-item-title
                                                        class="cursor-pointer"
                                                        @click="currentLocale = lang.languageTag">
                                                        {{ lang.displayName }}
                                                    </v-list-item-title>
                                                </v-list-item>
                                            </v-list>
                                        </v-menu>
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
                                <h4 class="text-h4 mb-1">{{ $t('Registration Completed') }}</h4>
                                <p class="my-5">{{ finalResultMessage }}</p>
                            </v-window-item>
                        </v-form>
                    </v-window>

                    <div class="d-flex justify-sm-space-between gap-4 flex-wrap justify-center mt-5">
                        <v-btn :color="(currentStep === 'basicSetting' || currentStep === 'finalResult') ? 'default' : 'primary'"
                               :disabled="currentStep === 'basicSetting' || currentStep === 'finalResult' || submitting || navigateToHomePage"
                               :prepend-icon="icons.previous"
                               @click="switchToPreviousTab">{{ $t('Previous') }}</v-btn>
                        <v-btn :color="(currentStep === 'presetCategories' || currentStep === 'finalResult') ? 'secondary' : 'primary'"
                               :disabled="currentStep === 'presetCategories' || currentStep === 'finalResult' || submitting || navigateToHomePage"
                               :append-icon="icons.next"
                               @click="switchToNextTab"
                               v-if="currentStep === 'basicSetting'">{{ $t('Next') }}</v-btn>
                        <v-btn color="teal"
                               :disabled="submitting || navigateToHomePage"
                               :append-icon="!submitting ? icons.submit : null"
                               @click="submit"
                               v-if="currentStep === 'presetCategories'">
                            {{ $t('Submit') }}
                            <v-progress-circular indeterminate size="22" class="ml-2" v-if="submitting"></v-progress-circular>
                        </v-btn>
                        <v-btn :append-icon="icons.next"
                               @click="navigateToLogin"
                               v-if="currentStep === 'finalResult'">{{ $t('Continue') }}</v-btn>
                    </div>
                </v-card>
            </v-col>
        </v-row>

        <snack-bar ref="snackbar" @update:show="onSnackbarShowStateChanged" />
    </div>
</template>

<script>
import { useTheme } from 'vuetify';

import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import assetConstants from '@/consts/asset.js';
import categoryConstants from '@/consts/category.js';
import { categorizedArrayToPlainArray } from '@/lib/common.js';
import { setExpenseAndIncomeAmountColor } from '@/lib/ui.js';

import {
    mdiArrowLeft,
    mdiArrowRight,
    mdiCheck
} from '@mdi/js';

export default {
    data() {
        const userStore = useUserStore();
        const newUser = userStore.generateNewUserModel(this.$locale.getCurrentLanguageTag());

        return {
            user: newUser,
            currentStep: 'basicSetting',
            submitting: false,
            usePresetCategories: false,
            finalResultMessage: null,
            navigateToHomePage: false,
            icons: {
                previous: mdiArrowLeft,
                next: mdiArrowRight,
                submit: mdiCheck
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
                return this.$locale.getCurrentLanguageTag();
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
        allSteps() {
            const allSteps = [
                {
                    name: 'basicSetting',
                    title: this.$t('User Information'),
                    subTitle: this.$t('Basic Information')
                },
                {
                    name: 'presetCategories',
                    title: this.$t('Transaction Categories'),
                    subTitle: this.$t('Preset Categories')
                }
            ];

            if (this.finalResultMessage) {
                allSteps.push({
                    name: 'finalResult',
                    title: this.$t('Complete'),
                    subTitle: this.$t('Registration Completed')
                });
            }

            return allSteps;
        },
        inputIsEmpty() {
            return !!this.inputEmptyProblemMessage;
        },
        inputIsInvalid() {
            return !!this.inputInvalidProblemMessage;
        },
        inputEmptyProblemMessage() {
            if (!this.user.username) {
                return 'Username cannot be blank';
            } else if (!this.user.nickname) {
                return 'Nickname cannot be blank';
            } else if (!this.user.email) {
                return 'Email address cannot be blank';
            } else if (!this.user.password) {
                return 'Password cannot be blank';
            } else if (!this.user.confirmPassword) {
                return 'Password confirmation cannot be blank';
            } else if (!this.user.defaultCurrency) {
                return 'Default currency cannot be blank';
            } else {
                return null;
            }
        },
        inputInvalidProblemMessage() {
            if (this.user.password && this.user.confirmPassword && this.user.password !== this.user.confirmPassword) {
                return 'Password and password confirmation do not match';
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
            if (this.submitting || this.currentStep === 'finalResult' || this.navigateToHomePage) {
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

            self.navigateToHomePage = false;
            self.submitting = true;

            let presetCategories = [];

            if (self.usePresetCategories) {
                presetCategories = categorizedArrayToPlainArray(self.allPresetCategories);
            }

            self.rootStore.register({
                user: self.user,
                presetCategories: presetCategories
            }).then(response => {
                if (!self.$user.isUserLogined()) {
                    self.submitting = false;

                    if (self.usePresetCategories && !response.presetCategoriesSaved) {
                        self.finalResultMessage = self.$t('You have been successfully registered, but there was an failure when adding preset categories. You can re-add preset categories in settings page anytime.');
                        self.currentStep = 'finalResult';
                    } else if (response.needVerifyEmail) {
                        self.finalResultMessage = self.$t('You have been successfully registered. An account activation link has been sent to your email address, please activate your account first.');
                        self.currentStep = 'finalResult';
                    } else {
                        self.$refs.snackbar.showMessage('You have been successfully registered');
                        self.navigateToHomePage = true;
                    }

                    return;
                }

                if (response.user) {
                    const localeDefaultSettings = self.$locale.setLanguage(response.user.language);
                    self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

                    setExpenseAndIncomeAmountColor(response.user.expenseAmountColor, response.user.incomeAmountColor);
                }

                if (self.settingsStore.appSettings.autoUpdateExchangeRatesData) {
                    self.exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
                }

                self.submitting = false;

                if (self.usePresetCategories && !response.presetCategoriesSaved) {
                    self.$refs.snackbar.showMessage('You have been successfully registered, but there was an failure when adding preset categories. You can re-add preset categories in settings page anytime.');
                } else {
                    self.$refs.snackbar.showMessage('You have been successfully registered');
                    self.$router.replace('/');
                }

                self.navigateToHomePage = true;
            }).catch(error => {
                self.submitting = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        navigateToLogin() {
            this.$router.push('/');
        },
        onSnackbarShowStateChanged(newValue) {
            if (!newValue && this.navigateToHomePage) {
                this.$router.replace('/');
            }
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
