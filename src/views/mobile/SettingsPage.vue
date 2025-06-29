<template>
    <f7-page>
        <f7-navbar :title="tt('Settings')" :back-link="tt('Back')"></f7-navbar>

        <f7-block-title class="margin-top">{{ currentNickName }}</f7-block-title>
        <f7-list strong inset dividers>
            <f7-list-item :title="tt('User Profile')" link="/user/profile"></f7-list-item>
            <f7-list-item :title="tt('Transaction Categories')" link="/category/all"></f7-list-item>
            <f7-list-item :title="tt('Transaction Tags')" link="/tag/list"></f7-list-item>
            <f7-list-item :title="tt('Transaction Templates')" link="/template/list"></f7-list-item>
            <f7-list-item :title="tt('Scheduled Transactions')" link="/schedule/list" v-if="isUserScheduledTransactionEnabled()"></f7-list-item>
            <f7-list-item :title="tt('Data Management')" link="/user/data/management"></f7-list-item>
            <f7-list-item :title="tt('Two-Factor Authentication')" link="/user/2fa"></f7-list-item>
            <f7-list-item :title="tt('Device & Sessions')" link="/user/sessions"></f7-list-item>
            <f7-list-button :class="{ 'disabled': logouting }" @click="logout">{{ tt('Log Out') }}</f7-list-button>
        </f7-list>

        <f7-block-title>{{ tt('Application') }}</f7-block-title>
        <f7-list strong inset dividers class="settings-list">
            <f7-list-item
                link="#"
                :title="tt('Theme')"
                :after="findNameByValue(allThemes, currentTheme)"
                @click="showThemePopup = true"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="value" value-field="value"
                                           title-field="name"
                                           :title="tt('Theme')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Theme')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allThemes"
                                           v-model:show="showThemePopup"
                                           v-model="currentTheme">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item :title="tt('Text Size')" link="/settings/textsize"></f7-list-item>

            <f7-list-item
                link="#"
                :title="tt('Timezone')"
                :after="currentTimezoneName"
                @click="showTimezonePopup = true"
            >
                <list-item-selection-popup value-type="item"
                                           key-field="name" value-field="name"
                                           title-field="displayNameWithUtcOffset"
                                           :title="tt('Timezone')"
                                           :enable-filter="true"
                                           :filter-placeholder="tt('Timezone')"
                                           :filter-no-items-text="tt('No results')"
                                           :items="allTimezones"
                                           v-model:show="showTimezonePopup"
                                           v-model="timeZone">
                </list-item-selection-popup>
            </f7-list-item>

            <f7-list-item :title="tt('Application Lock')" :after="isEnableApplicationLock ? tt('Enabled') : tt('Disabled')" link="/app_lock"></f7-list-item>

            <f7-list-item :title="tt('Exchange Rates Data')" :after="exchangeRatesLastUpdateDate" link="/exchange_rates"></f7-list-item>

            <f7-list-item>
                <span>{{ tt('Auto-update Exchange Rates Data') }}</span>
                <f7-toggle :checked="isAutoUpdateExchangeRatesData" @toggle:change="isAutoUpdateExchangeRatesData = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-item>
                <span>{{ tt('Show Account Balance') }}</span>
                <f7-toggle :checked="showAccountBalance" @toggle:change="showAccountBalance = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-item :title="tt('Page Settings')" link="/settings/page"></f7-list-item>
            <f7-list-item :title="tt('Statistics Settings')" link="/statistic/settings"></f7-list-item>
            <f7-list-item :title="tt('Settings Sync')" link="/settings/sync"></f7-list-item>

            <f7-list-item>
                <span>{{ tt('Enable Animation') }}</span>
                <f7-toggle :checked="isEnableAnimate" @toggle:change="isEnableAnimate = $event"></f7-toggle>
            </f7-list-item>

            <f7-list-item link="#" no-chevron :title="tt('Switch to Desktop Version')" @click="switchToDesktopVersion"></f7-list-item>

            <f7-list-item :title="tt('About')" link="/about" :after="version"></f7-list-item>
        </f7-list>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';
import { useAppSettingPageBase } from '@/views/base/settings/AppSettingsPageBase.ts';

import { useRootStore } from '@/stores/index.ts';
import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import { findNameByValue } from '@/lib/common.ts';
import { getVersion, getDesktopVersionPath } from '@/lib/version.ts';
import { isUserScheduledTransactionEnabled } from '@/lib/server_settings.ts';
import { setExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt, formatUnixTimeToLongDate, initLocale } = useI18n();
const { showToast, showConfirm } = useI18nUIComponents();
const { allThemes, allTimezones, timeZone, isAutoUpdateExchangeRatesData, showAccountBalance } = useAppSettingPageBase();

const rootStore = useRootStore();
const settingsStore = useSettingsStore();
const userStore = useUserStore();
const exchangeRatesStore = useExchangeRatesStore();

const version = `v${getVersion()}`;

const logouting = ref<boolean>(false);
const showThemePopup = ref<boolean>(false);
const showTimezonePopup = ref<boolean>(false);

const currentNickName = computed<string>(() => userStore.currentUserNickname || tt('User'));

const currentTheme = computed<string>({
    get: () => settingsStore.appSettings.theme,
    set: value => {
        if (value !== settingsStore.appSettings.theme) {
            settingsStore.setTheme(value);
            location.reload();
        }
    }
});

const currentTimezoneName = computed<string>(() => {
    for (const item of allTimezones.value) {
        if (item.name === timeZone.value) {
            return item.displayNameWithUtcOffset;
        }
    }

    return '';
});

const isEnableAnimate = computed<boolean>({
    get: () => settingsStore.appSettings.animate,
    set: value => {
        if (value !== settingsStore.appSettings.animate) {
            settingsStore.setEnableAnimate(value);
            location.reload();
        }
    }
});

const isEnableApplicationLock = computed<boolean>(() => settingsStore.appSettings.applicationLock);

const exchangeRatesLastUpdateDate = computed<string>(() => {
    const exchangeRatesLastUpdateTime = exchangeRatesStore.exchangeRatesLastUpdateTime;
    return exchangeRatesLastUpdateTime ? formatUnixTimeToLongDate(exchangeRatesLastUpdateTime) : '';
});

function switchToDesktopVersion(): void {
    showConfirm('Are you sure you want to switch to desktop version?', () => {
        window.location.replace(getDesktopVersionPath());
    });
}

function logout(): void {
    showConfirm('Are you sure you want to log out?', () => {
        logouting.value = true;
        showLoading(() => logouting.value);

        rootStore.logout().then(() => {
            logouting.value = false;
            hideLoading();

            settingsStore.clearAppSettings();

            const localeDefaultSettings = initLocale(userStore.currentUserLanguage, settingsStore.appSettings.timeZone);
            settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

            setExpenseAndIncomeAmountColor(userStore.currentUserExpenseAmountColor, userStore.currentUserIncomeAmountColor);

            props.f7router.navigate('/');
        }).catch(error => {
            logouting.value = false;
            hideLoading();

            if (!error.processed) {
                showToast(error.message || error);
            }
        });
    });
}
</script>
