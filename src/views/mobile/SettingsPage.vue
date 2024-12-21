<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar :title="$t('Settings')" :back-link="$t('Back')"></f7-navbar>

        <f7-block-title class="margin-top">{{ currentNickName }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item :title="$t('User Profile')" link="/user/profile"></f7-list-item>
            <f7-list-item :title="$t('Transaction Categories')" link="/category/all"></f7-list-item>
            <f7-list-item :title="$t('Transaction Tags')" link="/tag/list"></f7-list-item>
            <f7-list-item :title="$t('Transaction Templates')" link="/template/list"></f7-list-item>
            <f7-list-item :title="$t('Scheduled Transactions')" link="/schedule/list" v-if="isUserScheduledTransactionEnabled"></f7-list-item>
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

            <f7-list-item :title="$t('Text Size')" link="/settings/textsize"></f7-list-item>

            <f7-list-item
                :key="currentLocale + '_timezone'"
                :title="$t('Timezone')"
                smart-select :smart-select-params="{ openIn: 'popup', popupPush: true, closeOnSelect: true, scrollToSelectedItem: true, searchbar: true, searchbarPlaceholder: $t('Timezone'), searchbarDisableText: $t('Cancel'), appendSearchbarNotFound: $t('No results'), popupCloseLinkText: $t('Done') }">
                <select v-model="timeZone">
                    <option :value="tz.name" :key="tz.name"
                            v-for="tz in allTimezones">{{ tz.displayNameWithUtcOffset }}</option>
                </select>
            </f7-list-item>

            <f7-list-item :title="$t('Application Lock')" :after="isEnableApplicationLock ? $t('Enabled') : $t('Disabled')" link="/app_lock"></f7-list-item>

            <f7-list-item :title="$t('Exchange Rates Data')" :after="exchangeRatesLastUpdateDate" link="/exchange_rates"></f7-list-item>

            <f7-list-item>
                <span>{{ $t('Auto-update Exchange Rates Data') }}</span>
                <f7-toggle :checked="isAutoUpdateExchangeRatesData" @toggle:change="isAutoUpdateExchangeRatesData = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-item>
                <span>{{ $t('Show Account Balance') }}</span>
                <f7-toggle :checked="showAccountBalance" @toggle:change="showAccountBalance = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-item :title="$t('Page Settings')" link="/settings/page"></f7-list-item>

            <f7-list-item :title="$t('Statistics Settings')" link="/statistic/settings"></f7-list-item>

            <f7-list-item>
                <span>{{ $t('Enable Animation') }}</span>
                <f7-toggle :checked="isEnableAnimate" @toggle:change="isEnableAnimate = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-item :title="$t('Switch to Desktop Version')" @click="switchToDesktopVersion"></f7-list-item>

            <f7-list-item :title="$t('About')" link="/about" :after="version"></f7-list-item>
        </f7-list>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useTransactionsStore } from '@/stores/transaction.js';
import { useOverviewStore } from '@/stores/overview.js';
import { useStatisticsStore } from '@/stores/statistics.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import { getDesktopVersionPath } from '@/lib/version.js';
import { isUserScheduledTransactionEnabled } from '@/lib/server_settings.js';
import { setExpenseAndIncomeAmountColor } from '@/lib/ui.js';

export default {
    props: [
        'f7router'
    ],
    data() {
        const self = this;

        return {
            currentLocale: self.$locale.getCurrentLanguageTag(),
            logouting: false
        };
    },
    computed: {
        ...mapStores(useRootStore, useSettingsStore, useUserStore, useTransactionsStore, useOverviewStore, useStatisticsStore, useExchangeRatesStore),
        version() {
            return 'v' + this.$version;
        },
        allTimezones() {
            return this.$locale.getAllTimezones(true);
        },
        currentNickName() {
            return this.userStore.currentUserNickname || this.$t('User');
        },
        isUserScheduledTransactionEnabled() {
            return isUserScheduledTransactionEnabled();
        },
        theme: {
            get: function () {
                return this.settingsStore.appSettings.theme;
            },
            set: function (value) {
                if (value !== this.settingsStore.appSettings.theme) {
                    this.settingsStore.setTheme(value);
                    location.reload();
                }
            }
        },
        timeZone: {
            get: function () {
                return this.settingsStore.appSettings.timeZone;
            },
            set: function (value) {
                this.settingsStore.setTimeZone(value);
                this.$locale.setTimeZone(value);
                this.transactionsStore.updateTransactionListInvalidState(true);
                this.overviewStore.updateTransactionOverviewInvalidState(true);
                this.statisticsStore.updateTransactionStatisticsInvalidState(true);
            }
        },
        exchangeRatesLastUpdateDate() {
            const exchangeRatesLastUpdateTime = this.exchangeRatesStore.exchangeRatesLastUpdateTime;
            return exchangeRatesLastUpdateTime ? this.$locale.formatUnixTimeToLongDate(this.userStore, exchangeRatesLastUpdateTime) : '';
        },
        isAutoUpdateExchangeRatesData: {
            get: function () {
                return this.settingsStore.appSettings.autoUpdateExchangeRatesData;
            },
            set: function (value) {
                this.settingsStore.setAutoUpdateExchangeRatesData(value);
            }
        },
        isEnableApplicationLock() {
            return this.settingsStore.appSettings.applicationLock;
        },
        showAccountBalance: {
            get: function () {
                return this.settingsStore.appSettings.showAccountBalance;
            },
            set: function (value) {
                this.settingsStore.setShowAccountBalance(value);
            }
        },
        isEnableAnimate: {
            get: function () {
                return this.settingsStore.appSettings.animate;
            },
            set: function (value) {
                if (value !== this.settingsStore.appSettings.animate) {
                    this.settingsStore.setEnableAnimate(value);
                    location.reload();
                }
            }
        }
    },
    methods: {
        onPageAfterIn() {
            this.currentLocale = this.$locale.getCurrentLanguageTag();
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

                    self.settingsStore.clearAppSettings();

                    const localeDefaultSettings = self.$locale.initLocale(self.userStore.currentUserLanguage, self.settingsStore.appSettings.timeZone);
                    self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

                    setExpenseAndIncomeAmountColor(self.userStore.currentUserExpenseAmountColor, self.userStore.currentUserIncomeAmountColor);

                    router.navigate('/');
                }).catch(error => {
                    self.logouting = false;
                    self.$hideLoading();

                    if (!error.processed) {
                        self.$toast(error.message || error);
                    }
                });
            });
        },
        switchToDesktopVersion() {
            this.$confirm('Are you sure you want to switch to desktop version?', () => {
                window.location.replace(getDesktopVersionPath());
            });
        }
    }
};
</script>
