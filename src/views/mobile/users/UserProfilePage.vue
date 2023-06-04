<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('User Profile')"></f7-nav-title>
            <f7-nav-right>
                <f7-link :class="{ 'disabled': inputIsNotChanged || inputIsInvalid || saving }" :text="$t('Save')" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="loading">
            <f7-list-input label="Password" placeholder="Your password"></f7-list-input>
            <f7-list-input label="Confirmation Password" placeholder="Re-enter the password"></f7-list-input>
            <f7-list-input label="E-mail" placeholder="Your email address"></f7-list-input>
            <f7-list-input label="Nickname" placeholder="Your nickname"></f7-list-input>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="loading">
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Default Account" title="Not Specified"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Editable Transaction Scope" title="All" link="#"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-vertical skeleton-text" v-if="loading">
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Default Language" title="Language" link="#"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Default Currency" title="Currency" link="#"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="First Day of Week" title="Week Day" link="#"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Long Date Format" title="YYYY-MM-DD" link="#"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Short Date Format" title="YYYY-MM-DD" link="#"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Long Time Format" title="HH:mm:ss" link="#"></f7-list-item>
            <f7-list-item class="list-item-with-header-and-title list-item-no-item-after" header="Short Time Format" title="HH:mm" link="#"></f7-list-item>
        </f7-list>

        <f7-list form strong inset dividers class="margin-vertical" v-if="!loading">
            <f7-list-input
                type="password"
                autocomplete="new-password"
                clear-button
                :label="$t('Password')"
                :placeholder="$t('Your password')"
                v-model:value="newProfile.password"
            ></f7-list-input>

            <f7-list-input
                type="password"
                autocomplete="new-password"
                clear-button
                :label="$t('Confirmation Password')"
                :placeholder="$t('Re-enter the password')"
                v-model:value="newProfile.confirmPassword"
            ></f7-list-input>

            <f7-list-input
                type="email"
                autocomplete="email"
                clear-button
                :label="$t('E-mail')"
                :placeholder="$t('Your email address')"
                v-model:value="newProfile.email"
            ></f7-list-input>

            <f7-list-input
                type="text"
                autocomplete="nickname"
                clear-button
                :label="$t('Nickname')"
                :placeholder="$t('Your nickname')"
                v-model:value="newProfile.nickname"
            ></f7-list-input>

            <f7-list-item class="ebk-list-item-error-info" v-if="inputIsInvalid" :footer="$t(inputInvalidProblemMessage)"></f7-list-item>
        </f7-list>

        <f7-list form strong inset dividers class="margin-vertical" v-if="!loading">
            <f7-list-item
                class="list-item-with-header-and-title"
                link="#" no-chevron
                :class="{ 'disabled': !allVisibleAccounts.length }"
                :header="$t('Default Account')"
                :title="$utilities.getNameByKeyValue(allAccounts, newProfile.defaultAccountId, 'id', 'name', $t('Not Specified'))"
                @click="showAccountSheet = true"
            >
                <two-column-list-item-selection-sheet primary-key-field="id" primary-value-field="category"
                                                      primary-title-field="name"
                                                      primary-icon-field="icon" primary-icon-type="account"
                                                      primary-sub-items-field="accounts"
                                                      :primary-title-i18n="true"
                                                      secondary-key-field="id" secondary-value-field="id"
                                                      secondary-title-field="name"
                                                      secondary-icon-field="icon" secondary-icon-type="account" secondary-color-field="color"
                                                      :items="allCategorizedAccounts"
                                                      v-model:show="showAccountSheet"
                                                      v-model="newProfile.defaultAccountId">
                </two-column-list-item-selection-sheet>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="$t('Editable Transaction Scope')"
                :title="$t($utilities.getNameByKeyValue(allTransactionEditScopeTypes, newProfile.transactionEditScope, 'value', 'name'))"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Date Range'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), pageTitle: $t('Editable Transaction Scope'), popupCloseLinkText: $t('Done') }"
            >
                <select v-model="newProfile.transactionEditScope">
                    <option :value="option.value"
                            :key="option.value"
                            v-for="option in allTransactionEditScopeTypes">{{ $t(option.name) }}</option>
                </select>
            </f7-list-item>

            <f7-list-item class="ebk-list-item-error-info" v-if="extendInputIsInvalid" :footer="$t(extendInputInvalidProblemMessage)"></f7-list-item>
        </f7-list>

        <f7-list form strong inset dividers class="margin-vertical" v-if="!loading">
            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="$t('Language')"
                :title="currentLanguageName"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Language'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), pageTitle: $t('Language'), popupCloseLinkText: $t('Done') }">
                <select v-model="newProfile.language">
                    <option :value="language.code"
                            :key="language.code"
                            v-for="language in allLanguages">{{ language.displayName }}</option>
                </select>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="$t('Default Currency')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Currency Name'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), pageTitle: $t('Default Currency'), popupCloseLinkText: $t('Done') }"
            >
                <template #title>
                    <f7-block class="no-padding no-margin">
                        <span>{{ $t(`currency.${newProfile.defaultCurrency}`) }}&nbsp;</span>
                        <small class="smaller">{{ newProfile.defaultCurrency }}</small>
                    </f7-block>
                </template>
                <select autocomplete="transaction-currency" v-model="newProfile.defaultCurrency">
                    <option :value="currency.code"
                            :key="currency.code"
                            v-for="currency in allCurrencies">{{ currency.displayName }}</option>
                </select>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="$t('First Day of Week')"
                :title="currentDayOfWeekName"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Date'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), pageTitle: $t('First Day of Week'), popupCloseLinkText: $t('Done') }"
            >
                <select v-model="newProfile.firstDayOfWeek">
                    <option :value="weekDay.type"
                            :key="weekDay.type"
                            v-for="weekDay in allWeekDays">{{ $t(`datetime.${weekDay.name}.long`) }}</option>
                </select>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="$t('Long Date Format')"
                :title="$utilities.getNameByKeyValue(allLongDateFormats, newProfile.longDateFormat, 'type', 'displayName')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Long Date Format'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), pageTitle: $t('Long Date Format'), popupCloseLinkText: $t('Done') }"
            >
                <select v-model="newProfile.longDateFormat">
                    <option :value="format.type"
                            :key="format.type"
                            v-for="format in allLongDateFormats">{{ format.displayName }}</option>
                </select>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="$t('Short Date Format')"
                :title="$utilities.getNameByKeyValue(allShortDateFormats, newProfile.shortDateFormat, 'type', 'displayName')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Short Date Format'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), pageTitle: $t('Short Date Format'), popupCloseLinkText: $t('Done') }"
            >
                <select v-model="newProfile.shortDateFormat">
                    <option :value="format.type"
                            :key="format.type"
                            v-for="format in allShortDateFormats">{{ format.displayName }}</option>
                </select>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="$t('Long Time Format')"
                :title="$utilities.getNameByKeyValue(allLongTimeFormats, newProfile.longTimeFormat, 'type', 'displayName')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Long Time Format'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), pageTitle: $t('Long Time Format'), popupCloseLinkText: $t('Done') }"
            >
                <select v-model="newProfile.longTimeFormat">
                    <option :value="format.type"
                            :key="format.type"
                            v-for="format in allLongTimeFormats">{{ format.displayName }}</option>
                </select>
            </f7-list-item>

            <f7-list-item
                class="list-item-with-header-and-title list-item-no-item-after"
                :header="$t('Short Time Format')"
                :title="$utilities.getNameByKeyValue(allShortTimeFormats, newProfile.shortTimeFormat, 'type', 'displayName')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Long Time Format'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), pageTitle: $t('Short Time Format'), popupCloseLinkText: $t('Done') }"
            >
                <select v-model="newProfile.shortTimeFormat">
                    <option :value="format.type"
                            :key="format.type"
                            v-for="format in allShortTimeFormats">{{ format.displayName }}</option>
                </select>
            </f7-list-item>

            <f7-list-item class="ebk-list-item-error-info" v-if="langAndRegionInputIsInvalid" :footer="$t(langAndRegionInputInvalidProblemMessage)"></f7-list-item>
        </f7-list>

        <password-input-sheet :title="$t('Modify Password')"
                              :hint="$t('Please enter your current password when modifying your password')"
                              :confirm-disabled="saving"
                              :cancel-disabled="saving"
                              v-model:show="showInputPasswordSheet"
                              v-model="currentPassword"
                              @password:confirm="save()">
        </password-input-sheet>
    </f7-page>
</template>

<script>
export default {
    props: [
        'f7router'
    ],
    data() {
        return {
            newProfile: {
                password: '',
                confirmPassword: '',
                email: '',
                nickname: '',
                defaultAccountId: '',
                transactionEditScope: 1,
                language: '',
                defaultCurrency: '',
                firstDayOfWeek: 0,
                longDateFormat: 0,
                shortDateFormat: 0,
                longTimeFormat: 0,
                shortTimeFormat: 0
            },
            oldProfile: {
                email: '',
                nickname: '',
                defaultAccountId: '',
                transactionEditScope: 1,
                language: '',
                defaultCurrency: '',
                firstDayOfWeek: 0,
                longDateFormat: 0,
                shortDateFormat: 0,
                longTimeFormat: 0,
                shortTimeFormat: 0
            },
            currentPassword: '',
            loading: true,
            loadingError: null,
            saving: false,
            showInputPasswordSheet: false,
            showAccountSheet: false
        };
    },
    computed: {
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
            return this.$store.getters.allPlainAccounts;
        },
        allVisibleAccounts() {
            return this.$store.getters.allVisiblePlainAccounts;
        },
        allCategorizedAccounts() {
            return this.$utilities.getCategorizedAccounts(this.allVisibleAccounts);
        },
        allWeekDays() {
            return this.$constants.datetime.allWeekDays;
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
            return [{
                value: 0,
                name: 'None'
            }, {
                value: 1,
                name: 'All'
            }, {
                value: 2,
                name: 'Today or later'
            }, {
                value: 3,
                name: 'Recent 24 hours or later'
            }, {
                value: 4,
                name: 'This week or later'
            }, {
                value: 5,
                name: 'This month or later'
            }, {
                value: 6,
                name: 'This year or later'
            }];
        },
        currentLanguageName() {
            for (let i = 0; i < this.allLanguages.length; i++) {
                if (this.allLanguages[i].code === this.newProfile.language) {
                    return this.allLanguages[i].displayName;
                }
            }

            return this.$t('Unknown');
        },
        currentDayOfWeekName() {
            const weekName = this.$utilities.getNameByKeyValue(this.$constants.datetime.allWeekDays, this.newProfile.firstDayOfWeek, 'type', 'name');
            const i18nWeekNameKey = `datetime.${weekName}.long`;
            return this.$t(i18nWeekNameKey);
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
            if (!this.newProfile.password && !this.newProfile.confirmPassword && !this.newProfile.email && !this.newProfile.nickname) {
                return 'Nothing has been modified';
            } else if (!this.newProfile.password && !this.newProfile.confirmPassword &&
                this.newProfile.email === this.oldProfile.email &&
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
            } else if (!this.newProfile.password && this.newProfile.confirmPassword) {
                return 'Password cannot be empty';
            } else if (this.newProfile.password && !this.newProfile.confirmPassword) {
                return 'Confirmation password cannot be empty';
            } else {
                return null;
            }
        },
        inputInvalidProblemMessage() {
            if (this.newProfile.password && this.newProfile.confirmPassword && this.newProfile.password !== this.newProfile.confirmPassword) {
                return 'Password and confirmation password do not match';
            } else if (!this.newProfile.email) {
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
            self.$store.dispatch('loadAllAccounts', { force: false }),
            self.$store.dispatch('getCurrentUserProfile')
        ];

        Promise.all(promises).then(responses => {
            const profile = responses[1];
            self.setCurrentUserProfile(profile);
            self.loading = false;
        }).catch(error => {
            if (error.processed) {
                self.loading = false;
            } else {
                self.loadingError = error;
                self.$toast(error.message || error);
            }
        });
    },
    methods: {
        onPageAfterIn() {
            this.$routeBackOnError(this.f7router, 'loadingError');
        },
        save() {
            const self = this;
            const router = self.f7router;

            self.showInputPasswordSheet = false;

            let problemMessage = self.inputIsNotChangedProblemMessage || self.inputInvalidProblemMessage || self.extendInputInvalidProblemMessage || self.langAndRegionInputInvalidProblemMessage;

            if (problemMessage) {
                self.$alert(problemMessage);
                return;
            }

            if (self.newProfile.password && !self.currentPassword) {
                self.showInputPasswordSheet = true;
                return;
            }

            self.saving = true;
            self.$showLoading(() => self.saving);

            self.$store.dispatch('updateUserProfile', {
                profile: self.newProfile,
                currentPassword: self.currentPassword
            }).then(response => {
                self.saving = false;
                self.$hideLoading();
                self.currentPassword = '';

                if (response.user) {
                    self.setCurrentUserProfile(response.user);
                    self.$locale.setLanguage(response.user.language);
                }

                self.$toast('Your profile has been successfully updated');
                router.back();
            }).catch(error => {
                self.saving = false;
                self.$hideLoading();
                self.currentPassword = '';

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
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
        }
    }
};
</script>
