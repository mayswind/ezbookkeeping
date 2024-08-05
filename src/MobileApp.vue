<template>
    <img style="display: none;" :src="devCookiePath" v-if="!isProduction" />
    <f7-app v-bind="f7params">
        <f7-view id="main-view" class="safe-areas" main url="/"></f7-view>
    </f7-app>
</template>

<script>
import { f7ready } from 'framework7-vue';
import routes from './router/mobile.js';

import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useTokensStore } from '@/stores/token.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import assetConstants from '@/consts/asset.js';
import { isProduction } from '@/lib/version.js';
import { getTheme, isEnableAnimate } from '@/lib/settings.js';
import { loadMapAssets } from '@/lib/map/index.js';
import { setExpenseAndIncomeAmountColor } from '@/lib/ui.js';
import { isModalShowing, setAppFontSize } from '@/lib/ui.mobile.js';

export default {
    data() {
        const self = this;
        let darkMode = 'auto';

        if (getTheme() === 'light') {
            darkMode = false;
        } else if (getTheme() === 'dark') {
            darkMode = true;
        }

        return {
            isProduction: isProduction(),
            devCookiePath: isProduction() ? '' : '/dev/cookies',
            notification: null,
            f7params: {
                name: 'ezBookkeeping',
                theme: 'ios',
                colors: {
                    primary: '#c67e48'
                },
                routes: routes,
                darkMode: darkMode,
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
                    browserHistory: !self.isiOSHomeScreenMode(),
                    browserHistoryInitialMatch: true,
                    browserHistoryAnimate: false,
                    iosSwipeBackAnimateShadow: false,
                    mdSwipeBackAnimateShadow: false
                }
            },
            isDarkMode: undefined,
            hasPushPopupBackdrop: undefined,
            hasBackdrop: undefined
        }
    },
    computed: {
        ...mapStores(useRootStore, useSettingsStore, useUserStore, useTokensStore, useExchangeRatesStore),
        currentNotificationContent() {
            return this.rootStore.currentNotification;
        }
    },
    watch: {
        currentNotificationContent: function (newValue) {
            const self = this;

            if (self.notification) {
                self.notification.close();
                self.notification.destroy();
                self.notification = null;
            }

            if (newValue) {
                f7ready((f7) => {
                    self.notification = f7.notification.create({
                        icon: `<img alt="logo" src="${assetConstants.ezBookkeepingLogoPath}" />`,
                        title: self.$t('global.app.title'),
                        text: newValue,
                        closeOnClick: true,
                        on: {
                            close() {
                                self.rootStore.setNotificationContent(null);
                            }
                        }
                    }).open();
                });
            }
        }
    },
    created() {
        const self = this;

        let localeDefaultSettings = self.$locale.initLocale(self.userStore.currentUserLanguage, self.settingsStore.appSettings.timeZone);
        self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

        setExpenseAndIncomeAmountColor(self.userStore.currentUserExpenseAmountColor, self.userStore.currentUserIncomeAmountColor);

        if (self.$user.isUserLogined()) {
            if (!self.settingsStore.appSettings.applicationLock || self.$user.isUserUnlocked()) {
                // refresh token if user is logined
                self.tokensStore.refreshTokenAndRevokeOldToken().then(response => {
                    if (response.user) {
                        localeDefaultSettings = self.$locale.setLanguage(response.user.language);
                        self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

                        setExpenseAndIncomeAmountColor(response.user.expenseAmountColor, response.user.incomeAmountColor);

                        if (response.notificationContent) {
                            self.rootStore.setNotificationContent(response.notificationContent);
                        }
                    }
                });

                // auto refresh exchange rates data
                if (self.settingsStore.appSettings.autoUpdateExchangeRatesData) {
                    self.exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
                }
            }
        }
    },
    mounted() {
        setAppFontSize(this.settingsStore.appSettings.fontSize);

        f7ready((f7) => {
            this.isDarkMode = f7.darkMode;
            this.setThemeColorMeta(f7.darkMode);

            f7.on('actionsOpen', (event) => this.onBackdropChanged(event));
            f7.on('actionsClose', (event) => this.onBackdropChanged(event));
            f7.on('dialogOpen', (event) => this.onBackdropChanged(event));
            f7.on('dialogClose', (event) => this.onBackdropChanged(event));
            f7.on('popoverOpen', (event) => this.onBackdropChanged(event));
            f7.on('popoverClose', (event) => this.onBackdropChanged(event));
            f7.on('popupOpen', (event) => this.onBackdropChanged(event));
            f7.on('popupClose', (event) => this.onBackdropChanged(event));
            f7.on('sheetOpen', (event) => this.onBackdropChanged(event));
            f7.on('sheetClose', (event) => this.onBackdropChanged(event));

            f7.on('pageBeforeOut',  () => {
                if (isModalShowing()) {
                    f7.actions.close('.actions-modal.modal-in', false);
                    f7.dialog.close('.dialog.modal-in', false);
                    f7.popover.close('.popover.modal-in', false);
                    f7.popup.close('.popup.modal-in', false);
                    f7.sheet.close('.sheet-modal.modal-in', false);
                }
            });

            f7.on('darkModeChange', (isDarkMode) => {
                this.isDarkMode = isDarkMode;
                this.setThemeColorMeta(isDarkMode);
            });
        });

        document.addEventListener('DOMContentLoaded', () => {
            const languageInfo = this.$locale.getCurrentLanguageInfo();
            loadMapAssets(languageInfo ? languageInfo.alternativeLanguageTag : null);
        });
    },
    methods: {
        isiOSHomeScreenMode() {
            if ((/iphone|ipod|ipad/gi).test(navigator.platform) && (/Safari/i).test(navigator.appVersion) &&
                window.matchMedia && window.matchMedia('(display-mode: standalone)').matches
            ) {
                return true;
            }

            return false;
        },
        onBackdropChanged(event) {
            if (event.push) {
                this.hasPushPopupBackdrop = event.opened;
            } else {
                this.hasBackdrop = event.opened;
            }

            this.setThemeColorMeta(this.isDarkMode);
        },
        setThemeColorMeta(isDarkMode) {
            if (this.hasPushPopupBackdrop) {
                document.querySelector('meta[name=theme-color]').setAttribute('content', '#000');
                return;
            }

            if (isDarkMode) {
                if (this.hasBackdrop) {
                    document.querySelector('meta[name=theme-color]').setAttribute('content', '#0b0b0b');
                } else {
                    document.querySelector('meta[name=theme-color]').setAttribute('content', '#121212');
                }
            } else {
                if (this.hasBackdrop) {
                    document.querySelector('meta[name=theme-color]').setAttribute('content', '#949495');
                } else {
                    document.querySelector('meta[name=theme-color]').setAttribute('content', '#f6f6f8');
                }
            }
        }
    }
}
</script>
