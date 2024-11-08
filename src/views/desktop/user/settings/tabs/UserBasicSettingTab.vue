<template>
    <v-row>
        <v-col cols="12">
            <v-card :class="{ 'disabled': loading || saving }">
                <template #title>
                    <span>{{ $t('Basic Settings') }}</span>
                    <v-progress-circular indeterminate size="20" class="ml-3" v-if="loading"></v-progress-circular>
                </template>

                <v-card-text class="d-flex">
                    <v-avatar rounded="lg" variant="tonal" size="100" class="me-4 user-profile-avatar-icon"
                              :class="{ 'cursor-pointer': oldProfile.avatarProvider === 'internal', 'user-profile-avatar-icon-modifiable': oldProfile.avatarProvider === 'internal' }"
                              :color="currentUserAvatar ? 'rgba(0,0,0,0)' : 'primary'">
                        <v-img :src="currentUserAvatar" v-if="currentUserAvatar">
                            <template #placeholder>
                                <div class="d-flex align-center justify-center fill-height bg-light-primary">
                                    <v-icon color="primary" size="48" class="user-profile-avatar-placeholder" :icon="icons.user"/>
                                </div>
                            </template>
                        </v-img>
                        <v-icon size="48" class="user-profile-avatar-placeholder" :icon="icons.user" v-else-if="!currentUserAvatar"/>
                        <div class="avatar-edit-icon" v-if="oldProfile.avatarProvider === 'internal'">
                            <v-icon size="48" :icon="icons.pencil"/>
                        </div>
                        <v-menu activator="parent" width="200" location="bottom" offset="14px" v-if="oldProfile.avatarProvider === 'internal'">
                            <v-list>
                                <v-list-item :disabled="saving" :title="$t('Update Avatar')" @click="showOpenAvatarDialog"></v-list-item>
                                <v-list-item :disabled="!currentUserAvatar || saving" :title="$t('Remove Avatar')" @click="removeAvatar"></v-list-item>
                            </v-list>
                        </v-menu>
                    </v-avatar>
                    <div class="d-flex flex-column justify-center gap-3">
                        <div class="d-flex text-body-1">
                            <span class="me-1">{{ $t('Username:') }}</span>
                            <v-skeleton-loader class="skeleton-no-margin" type="text" style="width: 100px" :loading="true" v-if="loading"></v-skeleton-loader>
                            <span v-if="!loading">{{ oldProfile.username }}</span>
                        </div>
                        <div class="d-flex text-body-1 align-center" style="height: 40px;">
                            <span v-if="!loading && emailVerified">{{ $t('Email address is verified') }}</span>
                            <span v-if="!loading && !emailVerified">{{ $t('Email address is not verified') }}</span>
                            <v-btn class="ml-2 px-2" size="small" variant="text" :disabled="loading || resending"
                                   @click="resendVerifyEmail" v-if="isUserVerifyEmailEnabled && !loading && !emailVerified">
                                {{ $t('Resend Validation Email') }}
                                <v-progress-circular indeterminate size="18" class="ml-2" v-if="resending"></v-progress-circular>
                            </v-btn>
                            <v-skeleton-loader class="skeleton-no-margin mt-2 mb-1" type="text" style="width: 160px" :loading="true" v-if="loading"></v-skeleton-loader>
                        </div>
                    </div>
                </v-card-text>

                <v-divider />

                <v-form class="mt-6">
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-text-field
                                    type="text"
                                    autocomplete="nickname"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="$t('Nickname')"
                                    :placeholder="$t('Your nickname')"
                                    v-model="newProfile.nickname"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-text-field
                                    type="email"
                                    autocomplete="email"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="$t('E-mail')"
                                    :placeholder="$t('Your email address')"
                                    v-model="newProfile.email"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <two-column-select primary-key-field="id" primary-value-field="category"
                                                   primary-title-field="name"
                                                   primary-icon-field="icon" primary-icon-type="account"
                                                   primary-sub-items-field="accounts"
                                                   :primary-title-i18n="true"
                                                   secondary-key-field="id" secondary-value-field="id"
                                                   secondary-title-field="name"
                                                   secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                   :disabled="loading || saving || !allVisibleAccounts.length"
                                                   :label="$t('Default Account')"
                                                   :placeholder="$t('Default Account')"
                                                   :items="allVisibleCategorizedAccounts"
                                                   :no-item-text="$t('Unspecified')"
                                                   v-model="newProfile.defaultAccountId">
                                </two-column-select>
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="$t('Editable Transaction Range')"
                                    :placeholder="$t('Editable Transaction Range')"
                                    :items="allTransactionEditScopeTypes"
                                    v-model="newProfile.transactionEditScope"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>

                    <v-divider />

                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="languageTag"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="$t('Language')"
                                    :placeholder="$t('Language')"
                                    :items="allLanguages"
                                    v-model="newProfile.language"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-autocomplete
                                    item-title="displayName"
                                    item-value="currencyCode"
                                    auto-select-first
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="$t('Default Currency')"
                                    :placeholder="$t('Default Currency')"
                                    :items="allCurrencies"
                                    :no-data-text="$t('No results')"
                                    v-model="newProfile.defaultCurrency"
                                >
                                    <template #append-inner>
                                        <small class="text-field-append-text smaller">{{ newProfile.defaultCurrency }}</small>
                                    </template>
                                </v-autocomplete>
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="$t('First Day of Week')"
                                    :placeholder="$t('First Day of Week')"
                                    :items="allWeekDays"
                                    v-model="newProfile.firstDayOfWeek"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>

                    <v-divider />

                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="$t('Long Date Format')"
                                    :placeholder="$t('Long Date Format')"
                                    :items="allLongDateFormats"
                                    v-model="newProfile.longDateFormat"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="$t('Short Date Format')"
                                    :placeholder="$t('Short Date Format')"
                                    :items="allShortDateFormats"
                                    v-model="newProfile.shortDateFormat"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="$t('Long Time Format')"
                                    :placeholder="$t('Long Time Format')"
                                    :items="allLongTimeFormats"
                                    v-model="newProfile.longTimeFormat"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="$t('Short Time Format')"
                                    :placeholder="$t('Short Time Format')"
                                    :items="allShortTimeFormats"
                                    v-model="newProfile.shortTimeFormat"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>

                    <v-divider />

                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="$t('Currency Display Mode')"
                                    :placeholder="$t('Currency Display Mode')"
                                    :items="allCurrencyDisplayTypes"
                                    v-model="newProfile.currencyDisplayType"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="$t('Digit Grouping')"
                                    :placeholder="$t('Digit Grouping')"
                                    :items="allDigitGroupingTypes"
                                    v-model="newProfile.digitGrouping"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving || !getNameByKeyValue(allDigitGroupingTypes, newProfile.digitGrouping, 'type', 'enabled')"
                                    :label="$t('Digit Grouping Symbol')"
                                    :placeholder="$t('Digit Grouping Symbol')"
                                    :items="allDigitGroupingSymbols"
                                    v-model="newProfile.digitGroupingSymbol"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="$t('Decimal Separator')"
                                    :placeholder="$t('Decimal Separator')"
                                    :items="allDecimalSeparators"
                                    v-model="newProfile.decimalSeparator"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>

                    <v-divider />

                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="$t('Expense Amount Color')"
                                    :placeholder="$t('Expense Amount Color')"
                                    :items="allExpenseAmountColorTypes"
                                    v-model="newProfile.expenseAmountColor"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="$t('Income Amount Color')"
                                    :placeholder="$t('Income Amount Color')"
                                    :items="allIncomeAmountColorTypes"
                                    v-model="newProfile.incomeAmountColor"
                                />
                            </v-col>
                        </v-row>
                    </v-card-text>

                    <v-card-text class="d-flex flex-wrap gap-4">
                        <v-btn :disabled="inputIsNotChanged || inputIsInvalid || saving" @click="save">
                            {{ $t('Save Changes') }}
                            <v-progress-circular indeterminate size="22" class="ml-2" v-if="saving"></v-progress-circular>
                        </v-btn>

                        <v-btn color="default" variant="tonal" @click="reset">
                            {{ $t('Reset') }}
                        </v-btn>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>
    </v-row>

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
    <input ref="avatarInput" type="file" style="display: none" :accept="supportedImageExtensions" @change="updateAvatar($event)" />
</template>

<script>
import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useAccountsStore } from '@/stores/account.js';
import { useOverviewStore } from '@/stores/overview.js';

import datetimeConstants from '@/consts/datetime.js';
import fileConstants from '@/consts/file.js';
import { getNameByKeyValue } from '@/lib/common.js';
import { generateRandomUUID } from '@/lib/misc.js';
import { getCategorizedAccounts } from '@/lib/account.js';
import { isUserVerifyEmailEnabled } from '@/lib/server_settings.js';
import { setExpenseAndIncomeAmountColor } from '@/lib/ui.js';

import {
    mdiAccount,
    mdiAccountEditOutline
} from '@mdi/js';

export default {
    data() {
        const self = this;
        const defaultFirstDayOfWeekName = self.$locale.getDefaultFirstDayOfWeek();
        const defaultFirstDayOfWeek = datetimeConstants.allWeekDays[defaultFirstDayOfWeekName] ? datetimeConstants.allWeekDays[defaultFirstDayOfWeekName].type : datetimeConstants.defaultFirstDayOfWeek;

        return {
            newProfile: {
                email: '',
                nickname: '',
                defaultAccountId: 0,
                transactionEditScope: 1,
                language: '',
                defaultCurrency: self.$locale.getDefaultCurrency(),
                firstDayOfWeek: defaultFirstDayOfWeek,
                longDateFormat: 0,
                shortDateFormat: 0,
                longTimeFormat: 0,
                shortTimeFormat: 0,
                decimalSeparator: 0,
                digitGroupingSymbol: 0,
                digitGrouping: 0,
                currencyDisplayType: 0,
                expenseAmountColor: 0,
                incomeAmountColor: 0
            },
            oldProfile: {
                email: '',
                nickname: '',
                defaultAccountId: 0,
                transactionEditScope: 1,
                language: '',
                defaultCurrency: self.$locale.getDefaultCurrency(),
                firstDayOfWeek: defaultFirstDayOfWeek,
                longDateFormat: 0,
                shortDateFormat: 0,
                longTimeFormat: 0,
                shortTimeFormat: 0,
                decimalSeparator: 0,
                digitGroupingSymbol: 0,
                digitGrouping: 0,
                currencyDisplayType: 0,
                expenseAmountColor: 0,
                incomeAmountColor: 0
            },
            avatarNoCacheId: '',
            emailVerified: false,
            loading: true,
            resending: false,
            saving: false,
            icons: {
                user: mdiAccount,
                pencil: mdiAccountEditOutline,
            }
        };
    },
    computed: {
        ...mapStores(useRootStore, useSettingsStore, useUserStore, useAccountsStore, useOverviewStore),
        allLanguages() {
            return this.$locale.getAllLanguageInfoArray(true);
        },
        allCurrencies() {
            return this.$locale.getAllCurrencies();
        },
        allAccounts() {
            return this.accountsStore.allPlainAccounts;
        },
        allVisibleAccounts() {
            return this.accountsStore.allVisiblePlainAccounts;
        },
        allVisibleCategorizedAccounts() {
            return getCategorizedAccounts(this.allVisibleAccounts);
        },
        allWeekDays() {
            return this.$locale.getAllWeekDays();
        },
        allLongDateFormats() {
            return this.$locale.getAllLongDateFormats();
        },
        allShortDateFormats() {
            return this.$locale.getAllShortDateFormats();
        },
        allLongTimeFormats() {
            return this.$locale.getAllLongTimeFormats();
        },
        allShortTimeFormats() {
            return this.$locale.getAllShortTimeFormats();
        },
        allDecimalSeparators() {
            return this.$locale.getAllDecimalSeparators();
        },
        allDigitGroupingSymbols() {
            return this.$locale.getAllDigitGroupingSymbols();
        },
        allDigitGroupingTypes() {
            return this.$locale.getAllDigitGroupingTypes();
        },
        allCurrencyDisplayTypes() {
            return this.$locale.getAllCurrencyDisplayTypes(this.settingsStore, this.userStore);
        },
        allExpenseAmountColorTypes() {
            return this.$locale.getAllExpenseAmountColors();
        },
        allIncomeAmountColorTypes() {
            return this.$locale.getAllIncomeAmountColors();
        },
        allTransactionEditScopeTypes() {
            return this.$locale.getAllTransactionEditScopeTypes();
        },
        supportedImageExtensions() {
            return fileConstants.supportedImageExtensions;
        },
        currentUserAvatar() {
            return this.userStore.getUserAvatarUrl(this.oldProfile, this.avatarNoCacheId);
        },
        isUserVerifyEmailEnabled() {
            return isUserVerifyEmailEnabled();
        },
        inputIsNotChanged() {
            return !!this.inputIsNotChangedProblemMessage;
        },
        inputIsInvalid() {
            return !!this.inputInvalidProblemMessage;
        },
        inputIsNotChangedProblemMessage() {
            if (!this.newProfile.email && !this.newProfile.nickname) {
                return 'Nothing has been modified';
            } else if (this.newProfile.email === this.oldProfile.email &&
                this.newProfile.nickname === this.oldProfile.nickname &&
                this.newProfile.defaultAccountId === this.oldProfile.defaultAccountId &&
                this.newProfile.transactionEditScope === this.oldProfile.transactionEditScope &&
                this.newProfile.language === this.oldProfile.language &&
                this.newProfile.defaultCurrency === this.oldProfile.defaultCurrency &&
                this.newProfile.firstDayOfWeek === this.oldProfile.firstDayOfWeek &&
                this.newProfile.longDateFormat === this.oldProfile.longDateFormat &&
                this.newProfile.shortDateFormat === this.oldProfile.shortDateFormat &&
                this.newProfile.longTimeFormat === this.oldProfile.longTimeFormat &&
                this.newProfile.shortTimeFormat === this.oldProfile.shortTimeFormat &&
                this.newProfile.decimalSeparator === this.oldProfile.decimalSeparator &&
                this.newProfile.digitGroupingSymbol === this.oldProfile.digitGroupingSymbol &&
                this.newProfile.digitGrouping === this.oldProfile.digitGrouping &&
                this.newProfile.currencyDisplayType === this.oldProfile.currencyDisplayType &&
                this.newProfile.expenseAmountColor === this.oldProfile.expenseAmountColor &&
                this.newProfile.incomeAmountColor === this.oldProfile.incomeAmountColor) {
                return 'Nothing has been modified';
            } else {
                return null;
            }
        },
        inputInvalidProblemMessage() {
            if (!this.newProfile.email) {
                return 'Email address cannot be blank';
            } else if (!this.newProfile.nickname) {
                return 'Nickname cannot be blank';
            } else if (!this.newProfile.defaultCurrency) {
                return 'Default currency cannot be blank';
            } else {
                return null;
            }
        },
        extendInputInvalidProblemMessage() {
            return null;
        },
        langAndRegionInputInvalidProblemMessage() {
            if (!this.newProfile.defaultCurrency) {
                return 'Default currency cannot be blank';
            } else {
                return null;
            }
        }
    },
    created() {
        const self = this;

        self.loading = true;

        const promises = [
            self.accountsStore.loadAllAccounts({ force: false }),
            self.userStore.getCurrentUserProfile()
        ];

        Promise.all(promises).then(responses => {
            const profile = responses[1];
            self.setCurrentUserProfile(profile);
            self.emailVerified = profile.emailVerified;
            self.loading = false;
        }).catch(error => {
            self.oldProfile.nickname = '';
            self.oldProfile.email = '';
            self.newProfile.nickname = '';
            self.newProfile.email = '';
            self.loading = false;

            if (!error.processed) {
                self.$refs.snackbar.showError(error);
            }
        });
    },
    methods: {
        save() {
            const self = this;

            const problemMessage = self.inputIsNotChangedProblemMessage || self.inputInvalidProblemMessage || self.extendInputInvalidProblemMessage || self.langAndRegionInputInvalidProblemMessage;

            if (problemMessage) {
                self.$refs.snackbar.showMessage(problemMessage);
                return;
            }

            self.saving = true;

            self.rootStore.updateUserProfile({
                profile: self.newProfile
            }).then(response => {
                self.saving = false;

                if (response.user) {
                    if (response.user.firstDayOfWeek !== self.oldProfile.firstDayOfWeek) {
                        this.overviewStore.resetTransactionOverview();
                    }

                    self.setCurrentUserProfile(response.user);
                    self.emailVerified = response.user.emailVerified;

                    const localeDefaultSettings = self.$locale.setLanguage(response.user.language);
                    self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

                    setExpenseAndIncomeAmountColor(response.user.expenseAmountColor, response.user.incomeAmountColor);
                }

                self.$refs.snackbar.showMessage('Your profile has been successfully updated');
            }).catch(error => {
                self.saving = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        showOpenAvatarDialog() {
            this.$refs.avatarInput.click();
        },
        updateAvatar(event) {
            if (!event || !event.target || !event.target.files || !event.target.files.length) {
                return;
            }

            const self = this;
            const avatarFile = event.target.files[0];

            event.target.value = null;

            self.saving = true;

            self.userStore.updateUserAvatar({ avatarFile }).then(response => {
                self.saving = false;

                if (response) {
                    self.avatarNoCacheId = generateRandomUUID();
                    self.setCurrentUserProfile(response);
                }

                self.$refs.snackbar.showMessage('Your avatar has been successfully updated');
            }).catch(error => {
                self.saving = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        removeAvatar() {
            const self = this;

            self.$refs.confirmDialog.open('Are you sure you want to remove avatar?').then(() => {
                self.saving = true;

                self.userStore.removeUserAvatar().then(response => {
                    self.saving = false;

                    if (response) {
                        self.setCurrentUserProfile(response);
                    }

                    self.$refs.snackbar.showMessage('Your profile has been successfully updated');
                }).catch(error => {
                    self.saving = false;

                    if (!error.processed) {
                        self.$refs.snackbar.showError(error);
                    }
                });
            });
        },
        reset() {
            this.setCurrentUserProfile(this.oldProfile);
        },
        resendVerifyEmail() {
            const self = this;

            self.resending = true;

            self.rootStore.resendVerifyEmailByLoginedUser().then(() => {
                self.resending = false;
                self.$refs.snackbar.showMessage('Validation email has been sent');
            }).catch(error => {
                self.resending = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        getNameByKeyValue(src, value, keyField, nameField, defaultName) {
            return getNameByKeyValue(src, value, keyField, nameField, defaultName);
        },
        setCurrentUserProfile(profile) {
            this.oldProfile.username = profile.username;
            this.oldProfile.email = profile.email;
            this.oldProfile.nickname = profile.nickname;
            this.oldProfile.avatar = profile.avatar;
            this.oldProfile.avatarProvider = profile.avatarProvider;
            this.oldProfile.defaultAccountId = profile.defaultAccountId;
            this.oldProfile.transactionEditScope = profile.transactionEditScope;
            this.oldProfile.language = profile.language;
            this.oldProfile.defaultCurrency = profile.defaultCurrency;
            this.oldProfile.firstDayOfWeek = profile.firstDayOfWeek;
            this.oldProfile.longDateFormat = profile.longDateFormat;
            this.oldProfile.shortDateFormat = profile.shortDateFormat;
            this.oldProfile.longTimeFormat = profile.longTimeFormat;
            this.oldProfile.shortTimeFormat = profile.shortTimeFormat;
            this.oldProfile.decimalSeparator = profile.decimalSeparator;
            this.oldProfile.digitGroupingSymbol = profile.digitGroupingSymbol;
            this.oldProfile.digitGrouping = profile.digitGrouping;
            this.oldProfile.currencyDisplayType = profile.currencyDisplayType;
            this.oldProfile.expenseAmountColor = profile.expenseAmountColor;
            this.oldProfile.incomeAmountColor = profile.incomeAmountColor;

            this.newProfile.email = this.oldProfile.email
            this.newProfile.nickname = this.oldProfile.nickname;
            this.newProfile.defaultAccountId = this.oldProfile.defaultAccountId;
            this.newProfile.transactionEditScope = this.oldProfile.transactionEditScope;
            this.newProfile.language = this.oldProfile.language;
            this.newProfile.defaultCurrency = this.oldProfile.defaultCurrency;
            this.newProfile.firstDayOfWeek = this.oldProfile.firstDayOfWeek;
            this.newProfile.longDateFormat = this.oldProfile.longDateFormat;
            this.newProfile.shortDateFormat = this.oldProfile.shortDateFormat;
            this.newProfile.longTimeFormat = this.oldProfile.longTimeFormat;
            this.newProfile.shortTimeFormat = this.oldProfile.shortTimeFormat;
            this.newProfile.decimalSeparator = this.oldProfile.decimalSeparator;
            this.newProfile.digitGroupingSymbol = this.oldProfile.digitGroupingSymbol;
            this.newProfile.digitGrouping = this.oldProfile.digitGrouping;
            this.newProfile.currencyDisplayType = this.oldProfile.currencyDisplayType;
            this.newProfile.expenseAmountColor = this.oldProfile.expenseAmountColor;
            this.newProfile.incomeAmountColor = this.oldProfile.incomeAmountColor;
        }
    }
};
</script>

<style>
.user-profile-avatar-icon .avatar-edit-icon {
    display: none;
    position: absolute;
    width: 100% !important;
    height: 100% !important;
    background-color: rgba(0, 0, 0, 0.4);
}

.user-profile-avatar-icon .avatar-edit-icon > i.v-icon {
    background-color: transparent;
    color: rgba(255, 255, 255, 0.8);
}

.user-profile-avatar-icon:hover .avatar-edit-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    vertical-align: middle;
}

.user-profile-avatar-icon-modifiable:hover .user-profile-avatar-placeholder {
    display: none;
}
</style>
