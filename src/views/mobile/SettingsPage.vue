<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar :title="$t('Settings')" :back-link="$t('Back')"></f7-navbar>

        <f7-block-title class="margin-top">{{ currentNickName }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item :title="$t('User Profile')" link="/user/profile"></f7-list-item>
            <f7-list-item :title="$t('Transaction Categories')" link="/category/all"></f7-list-item>
            <f7-list-item :title="$t('Transaction Tags')" link="/tag/list"></f7-list-item>
            <f7-list-item :title="$t('Data Management')" link="/user/data/management"></f7-list-item>
            <f7-list-item :title="$t('Two-Factor Authentication')" link="/user/2fa"></f7-list-item>
            <f7-list-item :title="$t('Device & Sessions')" link="/user/sessions"></f7-list-item>
            <f7-list-button :class="{ 'disabled': logouting }" @click="logout">{{ $t('Log Out') }}</f7-list-button>
        </f7-list>

        <f7-block-title>{{ $t('Application') }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item
                :key="currentLocale + '_theme'"
                :title="$t('Theme')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Theme'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">
                <select v-model="theme">
                    <option value="auto">{{ $t('System Default') }}</option>
                    <option value="light">{{ $t('Light') }}</option>
                    <option value="dark">{{ $t('Dark') }}</option>
                </select>
            </f7-list-item>

            <f7-list-item
                :key="currentLocale + '_fontSize'"
                :title="$t('Font Size')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Font Size'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">
                <select v-model="fontSize">
                    <option value="default">{{ $t('Default') }}</option>
                    <option value="large">{{ $t('Large') }}</option>
                    <option value="extraLarge">{{ $t('Extra Large') }}</option>
                </select>
            </f7-list-item>

            <f7-list-item
                :key="currentLocale + '_timezone'"
                :title="$t('Timezone')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Timezone'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">
                <select v-model="currentTimezone">
                    <option :value="timezone.name"
                            :key="timezone.name"
                            v-for="timezone in allTimezones">{{ `(UTC${timezone.utcOffset}) ${timezone.displayName}` }}</option>
                </select>
            </f7-list-item>

            <f7-list-item :title="$t('Application Lock')" :after="isEnableApplicationLock ? $t('Enabled') : $t('Disabled')" link="/app_lock"></f7-list-item>

            <f7-list-item :title="$t('Exchange Rates Data')" :after="exchangeRatesLastUpdateDate" link="/exchange_rates"></f7-list-item>

            <f7-list-item>
                <span>{{ $t('Auto Update Exchange Rates Data') }}</span>
                <f7-toggle :checked="isAutoUpdateExchangeRatesData" @toggle:change="isAutoUpdateExchangeRatesData = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-item>
                <span>{{ $t('Auto Get Current Geographic Location') }}</span>
                <f7-toggle :checked="isAutoGetCurrentGeoLocation" @toggle:change="isAutoGetCurrentGeoLocation = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-item>
                <span>{{ $t('Enable Thousands Separator') }}</span>
                <f7-toggle :checked="isEnableThousandsSeparator" @toggle:change="isEnableThousandsSeparator = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-item
                :key="currentLocale + '_currency_display'"
                :title="$t('Currency Display Mode')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Currency Display Mode'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">
                <select v-model="currencyDisplayMode">
                    <option :value="allCurrencyDisplayModes.None">{{ $t('None') }}</option>
                    <option :value="allCurrencyDisplayModes.Symbol">{{ $t('Currency Symbol') }}</option>
                    <option :value="allCurrencyDisplayModes.Code">{{ $t('Currency Code') }}</option>
                    <option :value="allCurrencyDisplayModes.Name">{{ $t('Currency Name') }}</option>
                </select>
            </f7-list-item>

            <f7-list-item>
                <span>{{ $t('Show Amount In Home Page') }}</span>
                <f7-toggle :checked="showAmountInHomePage" @toggle:change="showAmountInHomePage = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-item>
                <span>{{ $t('Show Account Balance') }}</span>
                <f7-toggle :checked="showAccountBalance" @toggle:change="showAccountBalance = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-item>
                <span>{{ $t('Show Total Amount In Transaction List Page') }}</span>
                <f7-toggle :checked="showTotalAmountInTransactionListPage" @toggle:change="showTotalAmountInTransactionListPage = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-item :title="$t('Statistics Settings')" link="/statistic/settings"></f7-list-item>

            <f7-list-item>
                <span>{{ $t('Enable Animate') }}</span>
                <f7-toggle :checked="isEnableAnimate" @toggle:change="isEnableAnimate = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-item :title="$t('About')" link="/about" :after="version"></f7-list-item>
        </f7-list>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useUserStore } from '@/stores/user.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import currencyConstants from '@/consts/currency.js';
import { setAppFontSize } from '@/lib/ui.mobile.js';

export default {
    props: [
        'f7router'
    ],
    data() {
        const self = this;

        return {
            currentLocale: this.$i18n.locale,
            isEnableApplicationLock: self.$settings.isEnableApplicationLock(),
            logouting: false
        };
    },
    computed: {
        ...mapStores(useRootStore, useUserStore, useExchangeRatesStore),
        version() {
            return 'v' + this.$version;
        },
        allTimezones() {
            return this.$locale.getAllTimezones(true);
        },
        theme: {
            get: function () {
                return this.$settings.getTheme();
            },
            set: function (value) {
                if (value !== this.$settings.getTheme()) {
                    this.$settings.setTheme(value);
                    location.reload();
                }
            }
        },
        fontSize: {
            get: function () {
                return this.$settings.getFontSize();
            },
            set: function (value) {
                if (value !== this.$settings.getFontSize()) {
                    this.$settings.setFontSize(value);
                    setAppFontSize(value);
                }
            }
        },
        currentTimezone: {
            get: function () {
                return this.$locale.getTimezone();
            },
            set: function (value) {
                this.$locale.setTimezone(value);
            }
        },
        currentNickName() {
            return this.userStore.currentUserNickname || this.$t('User');
        },
        exchangeRatesLastUpdateDate() {
            const exchangeRatesLastUpdateTime = this.exchangeRatesStore.exchangeRatesLastUpdateTime;
            return exchangeRatesLastUpdateTime ? this.$locale.formatUnixTimeToLongDate(this.userStore, exchangeRatesLastUpdateTime) : '';
        },
        allCurrencyDisplayModes() {
            return currencyConstants.allCurrencyDisplayModes;
        },
        isAutoUpdateExchangeRatesData: {
            get: function () {
                return this.$settings.isAutoUpdateExchangeRatesData();
            },
            set: function (value) {
                this.$settings.setAutoUpdateExchangeRatesData(value);
            }
        },
        isAutoGetCurrentGeoLocation: {
            get: function () {
                return this.$settings.isAutoGetCurrentGeoLocation();
            },
            set: function (value) {
                this.$settings.setAutoGetCurrentGeoLocation(value);
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
        showAmountInHomePage: {
            get: function () {
                return this.$settings.isShowAmountInHomePage();
            },
            set: function (value) {
                this.$settings.setShowAmountInHomePage(value);
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
        showTotalAmountInTransactionListPage: {
            get: function () {
                return this.$settings.isShowTotalAmountInTransactionListPage();
            },
            set: function (value) {
                this.$settings.setShowTotalAmountInTransactionListPage(value);
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
        }
    },
    methods: {
        onPageAfterIn() {
            this.currentLocale = this.$i18n.locale;
            this.isEnableApplicationLock = this.$settings.isEnableApplicationLock();
        },
        logout() {
            const self = this;
            const router = self.f7router;

            self.$confirm('Are you sure you want to log out?', () => {
                self.logouting = true;
                self.$showLoading(() => self.logouting);

                self.rootStore.logout().then(() => {
                    self.logouting = false;
                    self.$hideLoading();

                    self.$settings.clearSettings();
                    self.$locale.initLocale();

                    router.navigate('/');
                }).catch(error => {
                    self.logouting = false;
                    self.$hideLoading();

                    if (!error.processed) {
                        self.$toast(error.message || error);
                    }
                });
            });
        }
    }
};
</script>
