<template>
    <v-row>
        <v-col cols="12">
            <v-card :class="{ 'disabled': loading || saving }">
                <template #title>
                    <span>{{ $t('Basic Settings') }}</span>
                    <v-progress-circular indeterminate size="24" class="ml-2" v-if="loading"></v-progress-circular>
                </template>

                <v-form>
                    <v-card-text>
                        <v-row>
                            <v-col cols="12" md="6">
                                <v-text-field
                                    type="text"
                                    autocomplete="nickname"
                                    clearable
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
                                    :disabled="loading || saving"
                                    :label="$t('E-mail')"
                                    :placeholder="$t('Your email address')"
                                    v-model="newProfile.email"
                                />
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="name"
                                    item-value="id"
                                    :disabled="loading || saving"
                                    :label="$t('Default Account')"
                                    :placeholder="$t('Default Account')"
                                    :items="allVisibleAccounts"
                                    v-model="newProfile.defaultAccountId"
                                >
                                    <template v-slot:selection="{ item }">
                                        <v-label>{{ !item || item.value === 0 || item.value === '0' ? $t('Not Specified') : item.title }}</v-label>
                                    </template>

                                    <template v-slot:no-data>
                                        <div class="px-4">{{ $t('No results') }}</div>
                                    </template>
                                </v-select>
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="name"
                                    item-value="value"
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
                                    :disabled="loading || saving"
                                    :label="$t('Default Currency')"
                                    :placeholder="$t('Default Currency')"
                                    :items="allCurrencies"
                                    v-model="newProfile.defaultCurrency"
                                >
                                    <template v-slot:no-data>
                                        <div class="px-4">{{ $t('No results') }}</div>
                                    </template>
                                </v-autocomplete>
                            </v-col>

                            <v-col cols="12" md="6">
                                <v-select
                                    item-title="displayName"
                                    item-value="type"
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
                        <v-btn :disabled="saving" @click="save">
                            {{ $t('Save changes') }}
                            <v-progress-circular indeterminate size="24" class="ml-2" v-if="saving"></v-progress-circular>
                        </v-btn>
                    </v-card-text>
                </v-form>
            </v-card>
        </v-col>
    </v-row>

    <confirm-dialog ref="confirmDialog"/>

    <v-snackbar v-model="showSnackbar">
        {{ snackbarMessage }}

        <template #actions>
            <v-btn color="primary" variant="text" @click="showSnackbar = false">{{ $t('Close') }}</v-btn>
        </template>
    </v-snackbar>
</template>

<script>
import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useAccountsStore } from '@/stores/account.js';

import datetimeConstants from '@/consts/datetime.js';
import { getNameByKeyValue } from '@/lib/common.js';

export default {
    data() {
        const self = this;
        const defaultFirstDayOfWeekName = self.$locale.getDefaultFirstDayOfWeek();
        const defaultFirstDayOfWeek = datetimeConstants.allWeekDays[defaultFirstDayOfWeekName] ? datetimeConstants.allWeekDays[defaultFirstDayOfWeekName].type : datetimeConstants.defaultFirstDayOfWeek;

        return {
            newProfile: {
                email: ' ',
                nickname: self.$t('Your nickname'),
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
                nickname: self.$t('Your nickname'),
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
            showSnackbar: false,
            snackbarMessage: ''
        };
    },
    computed: {
        ...mapStores(useRootStore, useSettingsStore, useUserStore, useAccountsStore),
        allLanguages() {
            const ret = [];
            const allLanguageInfo = this.$locale.getAllLanguageInfos();

            ret.push({
                code: '',
                displayName: this.$t('System Default')
            });

            for (let code in allLanguageInfo) {
                if (!Object.prototype.hasOwnProperty.call(allLanguageInfo, code)) {
                    continue;
                }

                const languageInfo = allLanguageInfo[code];

                ret.push({
                    code: code,
                    displayName: languageInfo.displayName
                });
            }

            return ret;
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
            const self = this;

            return [{
                value: 0,
                name: self.$t('None')
            }, {
                value: 1,
                name: self.$t('All')
            }, {
                value: 2,
                name: self.$t('Today or later')
            }, {
                value: 3,
                name: self.$t('Recent 24 hours or later')
            }, {
                value: 4,
                name: self.$t('This week or later')
            }, {
                value: 5,
                name: self.$t('This month or later')
            }, {
                value: 6,
                name: self.$t('This year or later')
            }];
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
                self.showSnackbarMessage(self.$tError(error.message || error));
            }
        });
    },
    methods: {
        save() {
            const self = this;

            const problemMessage = self.inputIsNotChangedProblemMessage || self.inputInvalidProblemMessage || self.extendInputInvalidProblemMessage || self.langAndRegionInputInvalidProblemMessage;

            if (problemMessage) {
                self.showSnackbarMessage(self.$t(problemMessage));
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

                self.showSnackbarMessage(self.$t('Your profile has been successfully updated'));
            }).catch(error => {
                self.saving = false;

                if (!error.processed) {
                    self.showSnackbarMessage(self.$tError(error.message || error));
                }
            });
        },
        getNameByKeyValue(src, value, keyField, nameField, defaultName) {
            return getNameByKeyValue(src, value, keyField, nameField, defaultName);
        },
        setCurrentUserProfile(profile) {
            this.oldProfile.email = profile.email;
            this.oldProfile.nickname = profile.nickname;
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
        },
        showSnackbarMessage(message) {
            this.showSnackbar = true;
            this.snackbarMessage = message;
        }
    }
};
</script>
