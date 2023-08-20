<template>
    <v-row>
        <v-col cols="12">
            <v-card :class="{ 'disabled': loading || saving }">
                <template #title>
                    <span>{{ $t('Basic Settings') }}</span>
                    <v-progress-circular indeterminate size="24" class="ml-2" v-if="loading"></v-progress-circular>
                </template>

                <v-card-text class="d-flex">
                    <v-avatar rounded="lg" color="primary" variant="tonal" size="100" class="me-4">
                        <v-img :src="oldProfile.avatar" v-if="oldProfile.avatar">
                            <template #placeholder>
                                <div class="d-flex align-center justify-center fill-height">
                                    <v-icon size="48" :icon="icons.user"/>
                                </div>
                            </template>
                        </v-img>
                        <v-icon size="48" :icon="icons.user" v-else-if="!oldProfile.avatar"/>
                    </v-avatar>
                    <div class="d-flex flex-column justify-center gap-5">
                        <p class="text-body-1 mb-0">
                            <span class="me-1">{{ $t('Username:') }}</span>
                            <span>{{ oldProfile.username }}</span>
                        </p>
                        <p class="text-body-1 mb-0">
                            <span class="me-1">{{ $t('Avatar Provider:') }}</span>
                            <span>{{ currentUserAvatarProvider }}</span>
                        </p>
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
                                    clearable
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
                                    clearable
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
                                                   :show-secondary-icon="true"
                                                   :label="$t('Default Account')"
                                                   :placeholder="$t('Default Account')"
                                                   :items="allCategorizedAccounts"
                                                   :no-item-text="$t('Not Specified')"
                                                   v-model="newProfile.defaultAccountId">
                                </two-column-select>
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="$t('Editable Transaction Scope')"
                                    :placeholder="$t('Editable Transaction Scope')"
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
                                    item-value="code"
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
                                    item-value="code"
                                    auto-select-first
                                    persistent-placeholder
                                    :disabled="loading || saving"
                                    :label="$t('Default Currency')"
                                    :placeholder="$t('Default Currency')"
                                    :items="allCurrencies"
                                    :no-data-text="$t('No results')"
                                    v-model="newProfile.defaultCurrency"
                                />
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

                    <v-card-text class="d-flex flex-wrap gap-4">
                        <v-btn :disabled="inputIsNotChanged || inputIsInvalid || saving" @click="save">
                            {{ $t('Save changes') }}
                            <v-progress-circular indeterminate size="24" class="ml-2" v-if="saving"></v-progress-circular>
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
</template>

<script>
import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useAccountsStore } from '@/stores/account.js';

import datetimeConstants from '@/consts/datetime.js';
import { getNameByKeyValue } from '@/lib/common.js';
import { getCategorizedAccounts } from '@/lib/account.js';

import {
    mdiAccount
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
                shortTimeFormat: 0
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
                shortTimeFormat: 0
            },
            loading: true,
            saving: false,
            icons: {
                user: mdiAccount
            }
        };
    },
    computed: {
        ...mapStores(useRootStore, useSettingsStore, useUserStore, useAccountsStore),
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
        allCategorizedAccounts() {
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
        allTransactionEditScopeTypes() {
            return this.$locale.getAllTransactionEditScopeTypes();
        },
        currentUserAvatarProvider() {
            if (this.oldProfile.avatarProvider === 'gravatar') {
                return 'Gravatar';
            } else {
                return this.$t('None');
            }
        },
        inputIsNotChanged() {
            return !!this.inputIsNotChangedProblemMessage;
        },
        inputIsInvalid() {
            return !!this.inputInvalidProblemMessage;
        },
        extendInputIsInvalid() {
            return !!this.extendInputInvalidProblemMessage;
        },
        langAndRegionInputIsInvalid() {
            return !!this.langAndRegionInputInvalidProblemMessage;
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
                this.newProfile.shortTimeFormat === this.oldProfile.shortTimeFormat) {
                return 'Nothing has been modified';
            } else {
                return null;
            }
        },
        inputInvalidProblemMessage() {
            if (!this.newProfile.email) {
                return 'Email address cannot be empty';
            } else if (!this.newProfile.nickname) {
                return 'Nickname cannot be empty';
            } else if (!this.newProfile.defaultCurrency) {
                return 'Default currency cannot be empty';
            } else {
                return null;
            }
        },
        extendInputInvalidProblemMessage() {
            return null;
        },
        langAndRegionInputInvalidProblemMessage() {
            if (!this.newProfile.defaultCurrency) {
                return 'Default currency cannot be empty';
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
                    self.setCurrentUserProfile(response.user);

                    const localeDefaultSettings = self.$locale.setLanguage(response.user.language);
                    self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
                }

                self.$refs.snackbar.showMessage('Your profile has been successfully updated');
            }).catch(error => {
                self.saving = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        reset() {
            this.setCurrentUserProfile(this.oldProfile);
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
        }
    }
};
</script>
