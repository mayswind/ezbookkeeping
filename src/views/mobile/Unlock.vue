<template>
    <f7-page no-toolbar no-navbar no-swipeback login-screen>
        <f7-login-screen-title>{{ $t('Unlock') }}</f7-login-screen-title>
        <f7-list form>
            <f7-list-item class="list-item-pincode-input">
                <PincodeInput secure :length="6" v-model="pinCode" @keyup.native="unlockByPin" />
            </f7-list-item>
        </f7-list>
        <f7-list>
            <f7-list-button :class="{ 'disabled': !pinCodeValid }" :text="$t('Unlock By PIN Code')" @click="unlockByPin"></f7-list-button>
            <f7-list-button v-if="isWebAuthnAvailable" :text="$t('Unlock By Face ID/Touch ID')" @click="unlockByWebAuthn"></f7-list-button>
            <f7-block-footer>
                <f7-link :text="$t('Re-login')" @click="relogin"></f7-link>
            </f7-block-footer>
        </f7-list>
    </f7-page>
</template>

<script>
export default {
    data() {
        return {
            pinCode: ''
        }
    },
    computed: {
        isWebAuthnAvailable() {
            return this.$settings.isEnableApplicationLockWebAuthn()
                && this.$user.getWebAuthnCredentialId()
                && this.$webauthn.isSupported();
        },
        pinCodeValid() {
            return this.pinCode && this.pinCode.length === 6;
        }
    },
    methods: {
        unlockByWebAuthn() {
            const self = this;
            const router = self.$f7router;

            if (!self.$settings.isEnableApplicationLockWebAuthn() || !self.$user.getWebAuthnCredentialId()) {
                self.$toast('Face ID/Touch ID authentication is not enabled');
                return;
            }

            if (!self.$webauthn.isSupported()) {
                self.$toast('This device does not support Face ID/Touch ID');
                return;
            }

            self.$showLoading();

            self.$webauthn.verifyCredential(
                self.$user.getUserInfo(),
                self.$user.getWebAuthnCredentialId()
            ).then(({ id, userName, userSecret }) => {
                self.$hideLoading();

                self.$user.unlockTokenByWebAuthn(id, userName, userSecret);
                self.$services.refreshToken();

                if (self.$settings.isAutoUpdateExchangeRatesData()) {
                    self.$services.autoRefreshLatestExchangeRates();
                }

                router.refreshPage();
            }).catch(error => {
                self.$hideLoading();
                self.$logger.error('failed to use webauthn to verify', error);

                if (error.notSupported) {
                    self.$toast('This device does not support Face ID/Touch ID');
                } else if (error.name === 'NotAllowedError') {
                    self.$toast('User has canceled authentication');
                } else if (error.invalid) {
                    self.$toast('Failed to authenticate by Face ID/Touch ID');
                } else {
                    self.$toast('User has canceled or this device does not support Face ID/Touch ID');
                }
            });
        },
        unlockByPin() {
            const app = this.$f7;
            const $$ = app.$;

            if (!this.pinCodeValid) {
                return;
            }

            if ($$('.modal-in').length) {
                return;
            }

            const router = this.$f7router;
            const user = this.$user.getUserInfo();

            if (!user || !user.username) {
                this.$alert('An error has occurred');
                return;
            }

            try {
                this.$user.unlockTokenByPinCode(user.username, this.pinCode);
                this.$services.refreshToken();

                if (this.$settings.isAutoUpdateExchangeRatesData()) {
                    this.$services.autoRefreshLatestExchangeRates();
                }

                router.refreshPage();
            } catch (ex) {
                this.$logger.error('failed to unlock by pin code', ex);
                this.$toast('PIN code is wrong');
            }
        },
        relogin() {
            const self = this;
            const router = self.$f7router;

            self.$confirm('Are you sure you want to re-login?', () => {
                self.$user.clearTokenAndUserInfo(true);
                self.$user.clearWebAuthnConfig();
                self.$exchangeRates.clearExchangeRates();
                self.$settings.clearSettings();
                self.$locale.init();

                router.navigate('/login', {
                    clearPreviousHistory: true
                });
            });
        }
    }
}
</script>
