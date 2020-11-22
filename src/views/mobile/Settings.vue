<template>
    <f7-page>
        <f7-navbar :title="$t('Settings')" :back-link="$t('Back')"></f7-navbar>

        <f7-block-title class="margin-top">{{ userNickName }}</f7-block-title>
        <f7-card>
            <f7-card-content :padding="false">
                <f7-list>
                    <f7-list-item :title="$t('User Profile')" link="/user/profile"></f7-list-item>
                    <f7-list-item :title="$t('Two-Factor Authentication')" link="/user/2fa"></f7-list-item>
                    <f7-list-item :title="$t('Device & Sessions')" link="/user/sessions"></f7-list-item>
                    <f7-list-button :class="{ 'disabled': logouting }" @click="logout">{{ $t('Log Out') }}</f7-list-button>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-block-title>{{ $t('Application') }}</f7-block-title>
        <f7-card>
            <f7-card-content :padding="false">
                <f7-list>
                    <f7-list-item
                        :key="currentLocale + '_lang'"
                        :title="$t('Language')"
                        smart-select :smart-select-params="{ openIn: 'sheet', closeOnSelect: true, sheetCloseLinkText: $t('Done') }">
                        <select v-model="currentLocale">
                            <option v-for="(lang, locale) in allLanguages"
                                    :key="locale"
                                    :value="locale">{{ lang.displayName }}</option>
                        </select>
                    </f7-list-item>

                    <f7-list-item :title="$t('Application Lock')" :after="isEnableApplicationLock ? $t('Enabled') : $t('Disabled')" link="/app_lock"></f7-list-item>

                    <f7-list-item :title="$t('Exchange Rates Data')" :after="exchangeRatesLastUpdateDate" link="/exchange_rates"></f7-list-item>

                    <f7-list-item>
                        <span>{{ $t('Auto Update Exchange Rates Data') }}</span>
                        <f7-toggle :checked="isAutoUpdateExchangeRatesData" @toggle:change="isAutoUpdateExchangeRatesData = $event"></f7-toggle>
                    </f7-list-item>

                    <f7-list-item>
                        <span>{{ $t('Enable Thousands Separator') }}</span>
                        <f7-toggle :checked="isEnableThousandsSeparator" @toggle:change="isEnableThousandsSeparator = $event"></f7-toggle>
                    </f7-list-item>

                    <f7-list-item
                        :key="currentLocale + '_currency_display'"
                        :title="$t('Currency Display Mode')"
                        smart-select :smart-select-params="{ openIn: 'sheet', closeOnSelect: true, sheetCloseLinkText: $t('Done') }">
                        <select v-model="currencyDisplayMode">
                            <option value="none">{{ $t('None') }}</option>
                            <option value="code">{{ $t('Currency Code') }}</option>
                            <option value="name">{{ $t('Currency Name') }}</option>
                        </select>
                    </f7-list-item>

                    <f7-list-item>
                        <span>{{ $t('Show Account Balance') }}</span>
                        <f7-toggle :checked="showAccountBalance" @toggle:change="showAccountBalance = $event"></f7-toggle>
                    </f7-list-item>

                    <f7-list-item>
                        <span>{{ $t('Enable Animate') }}</span>
                        <f7-toggle :checked="isEnableAnimate" @toggle:change="isEnableAnimate = $event"></f7-toggle>
                    </f7-list-item>

                    <f7-list-item>
                        <span>{{ $t('Enable Auto Dark Mode') }}</span>
                        <f7-toggle :checked="isEnableAutoDarkMode" @toggle:change="isEnableAutoDarkMode = $event"></f7-toggle>
                    </f7-list-item>

                    <f7-list-item :title="$t('About')" link="/about" :after="version"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>
    </f7-page>
</template>

<script>
export default {
    data() {
        return {
            logouting: false
        };
    },
    computed: {
        version() {
            return 'v' + this.$version();
        },
        userNickName() {
            const userInfo = this.$user.getUserInfo() || {};
            return userInfo.nickname || userInfo.username || this.$t('User');
        },
        allLanguages() {
            return this.$locale.getAllLanguages();
        },
        currentLocale: {
            get: function () {
                return this.$i18n.locale;
            },
            set: function (value) {
                this.$locale.setLanguage(value);
            }
        },
        isEnableApplicationLock() {
            return this.$settings.isEnableApplicationLock();
        },
        exchangeRatesLastUpdateDate() {
            const exchangeRates = this.$exchangeRates.getExchangeRates();
            return exchangeRates && exchangeRates.date ? this.$moment(exchangeRates.date).format(this.$t('format.date.long')) : '';
        },
        isAutoUpdateExchangeRatesData: {
            get: function () {
                return this.$settings.isAutoUpdateExchangeRatesData();
            },
            set: function (value) {
                this.$settings.setAutoUpdateExchangeRatesData(value);
            }
        },
        isEnableThousandsSeparator: {
            get: function () {
                return this.$settings.isEnableThousandsSeparator();
            },
            set: function (value) {
                this.$settings.setEnableThousandsSeparator(value);
            }
        },
        currencyDisplayMode: {
            get: function () {
                return this.$settings.getCurrencyDisplayMode();
            },
            set: function (value) {
                this.$settings.setCurrencyDisplayMode(value);
            }
        },
        showAccountBalance: {
            get: function () {
                return this.$settings.isShowAccountBalance();
            },
            set: function (value) {
                this.$settings.setShowAccountBalance(value);
            }
        },
        isEnableAnimate: {
            get: function () {
                return this.$settings.isEnableAnimate();
            },
            set: function (value) {
                if (value !== this.$settings.isEnableAnimate()) {
                    this.$settings.setEnableAnimate(value);
                    location.reload();
                }
            }
        },
        isEnableAutoDarkMode: {
            get: function () {
                return this.$settings.isEnableAutoDarkMode();
            },
            set: function (value) {
                if (value !== this.$settings.isEnableAutoDarkMode()) {
                    this.$settings.setEnableAutoDarkMode(value);
                    location.reload();
                }
            }
        }
    },
    methods: {
        logout() {
            const self = this;
            const router = self.$f7router;

            self.$confirm('Are you sure you want to log out?', () => {
                self.logouting = true;
                self.$showLoading(() => self.logouting);

                self.$services.logout().then(response => {
                    self.logouting = false;
                    self.$hideLoading();
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        self.$alert('Unable to logout');
                        return;
                    }

                    self.$user.clearTokenAndUserInfo();
                    self.$settings.clearSettings();
                    self.$exchangeRates.clearExchangeRates();
                    router.navigate('/');
                }).catch(error => {
                    self.$logger.error('failed to log out', error);

                    self.logouting = false;
                    self.$hideLoading();

                    if (error && error.processed) {
                        return;
                    }

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        self.$alert({ error: error.response.data });
                    } else if (!error.processed) {
                        self.$alert('Unable to logout');
                    }
                });
            });
        }
    }
};
</script>
