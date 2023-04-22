<template>
    <f7-page no-toolbar no-navbar no-swipeback login-screen>
        <f7-login-screen-title>
            <img alt="logo" class="login-page-logo" src="/img/ezbookkeeping-192.png" />
            <f7-block class="margin-vertical-half">{{ $t('global.app.title') }}</f7-block>
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
        version() {
            return 'v' + this.$version;
        },
        allLanguages() {
            return this.$locale.getAllLanguageInfos();
        },
        isWebAuthnAvailable() {
            return this.$settings.isEnableApplicationLockWebAuthn()
                && this.$user.getWebAuthnCredentialId()
                && this.$webauthn.isSupported();
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

            if (!self.$settings.isEnableApplicationLockWebAuthn() || !self.$user.getWebAuthnCredentialId()) {
                self.$toast('WebAuthn is not enabled');
                return;
            }

            if (!self.$webauthn.isSupported()) {
                self.$toast('This device does not support WebAuthn');
                return;
            }

            self.$showLoading();

            self.$webauthn.verifyCredential(
                self.$store.state.currentUserInfo,
                self.$user.getWebAuthnCredentialId()
            ).then(({ id, userName, userSecret }) => {
                self.$hideLoading();

                self.$user.unlockTokenByWebAuthn(id, userName, userSecret);
                self.$store.dispatch('refreshTokenAndRevokeOldToken');

                if (self.$settings.isAutoUpdateExchangeRatesData()) {
                    self.$store.dispatch('getLatestExchangeRates', { silent: true, force: false });
                }

                router.refreshPage();
            }).catch(error => {
                self.$hideLoading();
                self.$logger.error('failed to use webauthn to verify', error);

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
            if (!this.isPinCodeValid(pinCode)) {
                return;
            }

            if (this.$ui.isModalShowing()) {
                return;
            }

            const router = this.f7router;
            const user = this.$store.state.currentUserInfo;

            if (!user || !user.username) {
                this.$alert('An error has occurred');
                return;
            }

            try {
                this.$user.unlockTokenByPinCode(user.username, pinCode);
                this.$store.dispatch('refreshTokenAndRevokeOldToken');

                if (this.$settings.isAutoUpdateExchangeRatesData()) {
                    this.$store.dispatch('getLatestExchangeRates', { silent: true, force: false });
                }

                router.refreshPage();
            } catch (ex) {
                this.$logger.error('failed to unlock by pin code', ex);
                this.$toast('PIN code is wrong');
            }
        },
        relogin() {
            const self = this;
            const router = self.f7router;

            self.$confirm('Are you sure you want to re-login?', () => {
                self.$user.clearTokenAndUserInfo(true);
                self.$user.clearWebAuthnConfig();
                self.$store.dispatch('clearUserInfoState');
                self.$store.dispatch('resetState');
                self.$settings.clearSettings();
                self.$locale.initLocale();

                router.navigate('/login', {
                    clearPreviousHistory: true
                });
            });
        },
        isPinCodeValid(pinCode) {
            return pinCode && pinCode.length === 6;
        },
        changeLanguage(locale) {
            this.$locale.setLanguage(locale);
        }
    }
}
</script>
