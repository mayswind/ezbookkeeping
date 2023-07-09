<template>
    <f7-page no-toolbar no-navbar no-swipeback login-screen>
        <f7-login-screen-title>
            <img alt="logo" class="login-page-logo" :src="ezBookkeepingLogoPath" />
            <f7-block class="login-page-tile margin-vertical-half">{{ $t('global.app.title') }}</f7-block>
        </f7-login-screen-title>

        <f7-list form>
            <f7-list-item class="no-padding no-margin">
                <template #inner>
                    <div class="display-flex justify-content-center full-line">{{ $t('Unlock Application') }}</div>
                </template>
            </f7-list-item>
            <f7-list-item class="list-item-pincode-input padding-horizontal margin-horizontal">
                <pin-code-input :secure="true" :length="6" v-model="pinCode" @pincode:confirm="unlockByPin" />
            </f7-list-item>
        </f7-list>

        <f7-list>
            <f7-list-button :class="{ 'disabled': !isPinCodeValid(pinCode) }" :text="$t('Unlock By PIN Code')" @click="unlockByPin"></f7-list-button>
            <f7-list-button v-if="isWebAuthnAvailable" :text="$t('Unlock By WebAuthn')" @click="unlockByWebAuthn"></f7-list-button>
            <f7-block-footer>
                <f7-link :text="$t('Re-login')" @click="relogin"></f7-link>
            </f7-block-footer>
            <f7-block-footer class="padding-bottom">
            </f7-block-footer>
        </f7-list>

        <f7-button small popover-open=".lang-popover-menu" :text="currentLanguageName"></f7-button>

        <f7-list>
            <f7-block-footer>
                <span>Powered by </span>
                <f7-link external href="https://github.com/mayswind/ezbookkeeping" target="_blank">ezBookkeeping</f7-link>&nbsp;
                <span>{{ version }}</span>
            </f7-block-footer>
            <f7-block-footer>
            </f7-block-footer>
        </f7-list>

        <f7-popover class="lang-popover-menu">
            <f7-list dividers>
                <f7-list-item
                    link="#" no-chevron popover-close
                    :title="lang.displayName"
                    :key="locale"
                    v-for="(lang, locale) in allLanguages"
                    @click="changeLanguage(locale)"
                >
                    <template #after>
                        <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="$i18n.locale === locale"></f7-icon>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-popover>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useTokensStore } from '@/stores/token.js';
import { useExchangeRatesStore } from '@/stores/exchangeRates.js';

import assetConstants from '@/consts/asset.js';
import logger from '@/lib/logger.js';
import webauthn from '@/lib/webauthn.js';
import { isModalShowing } from '@/lib/ui.mobile.js';

export default {
    props: [
        'f7router'
    ],
    data() {
        return {
            pinCode: ''
        }
    },
    computed: {
        ...mapStores(useRootStore, useSettingsStore, useUserStore, useTokensStore, useExchangeRatesStore),
        ezBookkeepingLogoPath() {
            return assetConstants.ezBookkeepingLogoPath;
        },
        version() {
            return 'v' + this.$version;
        },
        allLanguages() {
            return this.$locale.getAllLanguageInfos();
        },
        isWebAuthnAvailable() {
            return this.settingsStore.appSettings.applicationLockWebAuthn
                && this.$user.getWebAuthnCredentialId()
                && webauthn.isSupported();
        },
        currentLanguageName() {
            const currentLocale = this.$i18n.locale;
            let lang = this.$locale.getLanguageInfo(currentLocale);

            if (!lang) {
                lang = this.$locale.getLanguageInfo(this.$locale.getDefaultLanguage());
            }

            return lang.displayName;
        }
    },
    methods: {
        unlockByWebAuthn() {
            const self = this;
            const router = self.f7router;

            if (!self.settingsStore.appSettings.applicationLockWebAuthn || !self.$user.getWebAuthnCredentialId()) {
                self.$toast('WebAuthn is not enabled');
                return;
            }

            if (!webauthn.isSupported()) {
                self.$toast('This device does not support WebAuthn');
                return;
            }

            self.$showLoading();

            webauthn.verifyCredential(
                self.userStore.currentUserInfo,
                self.$user.getWebAuthnCredentialId()
            ).then(({ id, userName, userSecret }) => {
                self.$hideLoading();

                self.$user.unlockTokenByWebAuthn(id, userName, userSecret);
                self.tokensStore.refreshTokenAndRevokeOldToken().then(response => {
                    if (response.user) {
                        const localeDefaultSettings = self.$locale.setLanguage(response.user.language);
                        self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
                    }
                });

                if (self.settingsStore.appSettings.autoUpdateExchangeRatesData) {
                    self.exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
                }

                router.refreshPage();
            }).catch(error => {
                self.$hideLoading();
                logger.error('failed to use webauthn to verify', error);

                if (error.notSupported) {
                    self.$toast('This device does not support WebAuthn');
                } else if (error.name === 'NotAllowedError') {
                    self.$toast('User has canceled authentication');
                } else if (error.invalid) {
                    self.$toast('Failed to authenticate by WebAuthn');
                } else {
                    self.$toast('User has canceled or this device does not support WebAuthn');
                }
            });
        },
        unlockByPin(pinCode) {
            const self = this;

            if (!self.isPinCodeValid(pinCode)) {
                return;
            }

            if (isModalShowing()) {
                return;
            }

            const router = self.f7router;
            const user = self.userStore.currentUserInfo;

            if (!user || !user.username) {
                self.$alert('An error has occurred');
                return;
            }

            try {
                self.$user.unlockTokenByPinCode(user.username, pinCode);
                self.tokensStore.refreshTokenAndRevokeOldToken().then(response => {
                    if (response.user) {
                        const localeDefaultSettings = self.$locale.setLanguage(response.user.language);
                        self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
                    }
                });

                if (self.settingsStore.appSettings.autoUpdateExchangeRatesData) {
                    self.exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
                }

                router.refreshPage();
            } catch (ex) {
                logger.error('failed to unlock by pin code', ex);
                self.$toast('PIN code is wrong');
            }
        },
        relogin() {
            const self = this;
            const router = self.f7router;

            self.$confirm('Are you sure you want to re-login?', () => {
                self.rootStore.forceLogout();
                self.settingsStore.clearAppSettings();

                const localeDefaultSettings = self.$locale.initLocale(self.userStore.currentUserLanguage, self.settingsStore.appSettings.timeZone);
                self.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

                router.navigate('/login', {
                    clearPreviousHistory: true
                });
            });
        },
        isPinCodeValid(pinCode) {
            return pinCode && pinCode.length === 6;
        },
        changeLanguage(locale) {
            const localeDefaultSettings = this.$locale.setLanguage(locale);
            this.settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
        }
    }
}
</script>
