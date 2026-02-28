<template>
    <f7-app v-bind="f7params">
        <f7-view id="main-view" class="safe-areas" main url="/"></f7-view>
    </f7-app>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';

import type { Framework7Parameters, Notification, Actions, Dialog, Popover, Popup, Sheet } from 'framework7/types';
import { f7ready } from 'framework7-vue';
import routes from './router/mobile.ts';

import { useI18n } from '@/locales/helpers.ts';

import { useRootStore } from '@/stores/index.ts';
import { useSettingsStore } from '@/stores/setting.ts';
import { useEnvironmentsStore } from '@/stores/environment.ts';
import { useUserStore } from '@/stores/user.ts';
import { useTokensStore } from '@/stores/token.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import { APPLICATION_LOGO_PATH } from '@/consts/asset.ts';
import { ThemeType } from '@/core/theme.ts';
import { isProduction } from '@/lib/version.ts';
import { getTheme, isEnableSwipeBack, isEnableAnimate } from '@/lib/settings.ts';
import { initMapProvider } from '@/lib/map/index.ts';
import { isUserLogined, isUserUnlocked } from '@/lib/userstate.ts';
import { updateMapCacheExpiration } from '@/lib/cache.ts';
import { setExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';
import { isiOSHomeScreenMode, isModalShowing, setAppFontSize } from '@/lib/ui/mobile.ts';

const { tt, getCurrentLanguageInfo, setLanguage, initLocale } = useI18n();

const rootStore = useRootStore();
const settingsStore = useSettingsStore();
const environmentsStore = useEnvironmentsStore();
const userStore = useUserStore();
const tokensStore = useTokensStore();
const exchangeRatesStore = useExchangeRatesStore();

const f7params = ref<Framework7Parameters>({
    name: 'ezBookkeeping',
    theme: 'ios',
    colors: {
        primary: '#c67e48'
    },
    routes: routes,
    darkMode: (() => {
        let darkMode: boolean | string = 'auto';

        if (getTheme() === ThemeType.Light) {
            darkMode = false;
        } else if (getTheme() === ThemeType.Dark) {
            darkMode = true;
        }

        return darkMode;
    })(),
    touch: {
        disableContextMenu: true,
        tapHold: true
    },
    serviceWorker: {
        path: isProduction() ? './sw.js' : undefined,
        scope: './',
    },
    actions: {
        animate: isEnableAnimate(),
        backdrop: true,
        closeOnEscape: true
    },
    dialog: {
        // @ts-expect-error there is an "animate" field in dialog parameters, but it is not declared in the type definition file
        animate: isEnableAnimate(),
        backdrop: true
    },
    popover: {
        animate: isEnableAnimate(),
        backdrop: true,
        closeOnEscape: true
    },
    popup: {
        animate: isEnableAnimate(),
        backdrop: true,
        closeOnEscape: true,
        swipeToClose: true
    },
    sheet: {
        animate: isEnableAnimate(),
        backdrop: true,
        closeOnEscape: true
    },
    smartSelect: {
        routableModals: false
    },
    view: {
        animate: isEnableAnimate(),
        browserHistory: !isiOSHomeScreenMode(),
        browserHistoryInitialMatch: true,
        browserHistoryAnimate: false,
        iosSwipeBack: isEnableSwipeBack(),
        iosSwipeBackAnimateShadow: false,
        mdSwipeBack: isEnableSwipeBack(),
        mdSwipeBackAnimateShadow: false
    }
});

const notification = ref<Notification.Notification | null>(null);

const hasPushPopupBackdrop = ref<boolean | undefined>(undefined);
const hasBackdrop = ref<boolean | undefined>(undefined);
const currentNotificationContent = computed<string | null>(() => rootStore.currentNotification);

function setThemeColorMeta(darkMode: boolean | undefined): void {
    if (hasPushPopupBackdrop.value) {
        document.querySelector('meta[name=theme-color]')?.setAttribute('content', '#000');
        return;
    }

    if (darkMode) {
        if (hasBackdrop.value) {
            document.querySelector('meta[name=theme-color]')?.setAttribute('content', '#0b0b0b');
        } else {
            document.querySelector('meta[name=theme-color]')?.setAttribute('content', '#121212');
        }
    } else {
        if (hasBackdrop.value) {
            document.querySelector('meta[name=theme-color]')?.setAttribute('content', '#949495');
        } else {
            document.querySelector('meta[name=theme-color]')?.setAttribute('content', '#f6f6f8');
        }
    }
}

function onBackdropChanged(element: { push?: boolean, opened?: boolean }): void {
    if (element.push) {
        hasPushPopupBackdrop.value = element.opened;
    } else {
        hasBackdrop.value = element.opened;
    }

    setThemeColorMeta(environmentsStore.framework7DarkMode);
}

onMounted(() => {
    setAppFontSize(settingsStore.appSettings.fontSize);

    f7ready((f7) => {
        environmentsStore.framework7DarkMode = f7.darkMode;
        setThemeColorMeta(f7.darkMode);

        f7.on('actionsOpen', (actions: Actions.Actions) => onBackdropChanged(actions));
        f7.on('actionsClose', (actions: Actions.Actions) => onBackdropChanged(actions));
        f7.on('dialogOpen', (dialog: Dialog.Dialog) => onBackdropChanged(dialog));
        f7.on('dialogClose', (dialog: Dialog.Dialog) => onBackdropChanged(dialog));
        f7.on('popoverOpen', (popover: Popover.Popover) => onBackdropChanged(popover));
        f7.on('popoverClose', (popover: Popover.Popover) => onBackdropChanged(popover));
        f7.on('popupOpen', (popup: Popup.Popup) => onBackdropChanged(popup));
        f7.on('popupClose', (popup: Popup.Popup) => onBackdropChanged(popup));
        f7.on('sheetOpen', (sheet: Sheet.Sheet) => onBackdropChanged(sheet));
        f7.on('sheetClose', (sheet: Sheet.Sheet) => onBackdropChanged(sheet));

        f7.on('pageBeforeOut',  () => {
            if (isModalShowing()) {
                f7.actions.close('.actions-modal.modal-in', false);
                f7.dialog.close('.dialog.modal-in', false);
                f7.popover.close('.popover.modal-in', false);
                f7.popup.close('.popup.modal-in', false);
                f7.sheet.close('.sheet-modal.modal-in', false);
            }
        });

        f7.on('darkModeChange', (darkMode) => {
            environmentsStore.framework7DarkMode = darkMode;
            setThemeColorMeta(darkMode);
        });
    });

    document.addEventListener('DOMContentLoaded', () => {
        const languageInfo = getCurrentLanguageInfo();
        initMapProvider(languageInfo?.alternativeLanguageTag);
    });

    if ('serviceWorker' in navigator) {
        navigator.serviceWorker.addEventListener('controllerchange', () => {
            updateMapCacheExpiration(settingsStore.appSettings.mapCacheExpiration);
        });

        if (navigator.serviceWorker.controller) {
            updateMapCacheExpiration(settingsStore.appSettings.mapCacheExpiration);
        }
    }
});

watch(currentNotificationContent, (newValue) => {
    if (notification.value) {
        notification.value.close();
        // @ts-expect-error there is an "destroy" function in the Notification component, but it is not defined in the type definition file
        // see https://framework7.io/docs/notification
        notification.value.destroy();
        notification.value = null;
    }

    if (newValue) {
        f7ready((f7) => {
            notification.value = f7.notification.create({
                icon: `<img alt="logo" src="${APPLICATION_LOGO_PATH}" />`,
                title: tt('global.app.title'),
                text: newValue,
                closeOnClick: true,
                on: {
                    close() {
                        rootStore.setNotificationContent(null);
                    }
                }
            }).open();
        });
    }
});

let localeDefaultSettings = initLocale(userStore.currentUserLanguage, settingsStore.appSettings.timeZone);
settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

setExpenseAndIncomeAmountColor(userStore.currentUserExpenseAmountColor, userStore.currentUserIncomeAmountColor);

if (isUserLogined()) {
    if (!settingsStore.appSettings.applicationLock || isUserUnlocked()) {
        // refresh token if user is logined
        tokensStore.refreshTokenAndRevokeOldToken().then(response => {
            if (response.user) {
                localeDefaultSettings = setLanguage(response.user.language);
                settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

                setExpenseAndIncomeAmountColor(response.user.expenseAmountColor, response.user.incomeAmountColor);

                if (response.notificationContent) {
                    rootStore.setNotificationContent(response.notificationContent);
                }
            }
        });

        // auto refresh exchange rates data
        if (settingsStore.appSettings.autoUpdateExchangeRatesData) {
            exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
        }
    }
}
</script>
