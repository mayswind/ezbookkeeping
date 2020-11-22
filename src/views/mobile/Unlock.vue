<template>
    <f7-page no-toolbar no-navbar no-swipeback login-screen>
        <f7-login-screen-title>{{ $t('PIN Code') }}</f7-login-screen-title>
        <f7-list form>
            <f7-list-item class="list-item-pincode-input">
                <PincodeInput secure :length="6" v-model="pinCode" @keyup.native="unlock" />
            </f7-list-item>
        </f7-list>
        <f7-list>
            <f7-list-button :class="{ 'disabled': !pinCodeValid }" :text="$t('Unlock')" @click="unlock"></f7-list-button>
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
        pinCodeValid() {
            return this.pinCode && this.pinCode.length === 6;
        }
    },
    created() {
        const self = this;
        const router = self.$f7router;

        if (self.$settings.isEnableApplicationLockWebAuthn() && self.$user.getWebAuthnCredentialId()) {
            self.$webauthn.verifyCredential(
                self.$user.getUserInfo(),
                self.$user.getWebAuthnCredentialId()
            ).then(({ id, userSecret }) => {
                self.$user.unlockTokenByWebAuthn(id, userSecret);
                self.$services.refreshToken();

                if (self.$settings.isAutoUpdateExchangeRatesData()) {
                    self.$services.autoRefreshLatestExchangeRates();
                }

                router.refreshPage();
            }).catch(error => {
                self.$logger.error('failed to use webauthn to verify', error);

                if (error.notSupported) {
                    self.$toast('This device does not support Face ID/Touch ID');
                } else if (error.invalid) {
                    self.$toast('Failed to authenticate by Face ID/Touch ID');
                } else {
                    self.$toast('User has canceled or this device does not support Face ID/Touch ID');
                }
            });
        }
    },
    methods: {
        unlock() {
            const app = this.$f7;
            const $$ = app.$;

            if (!this.pinCodeValid) {
                return;
            }

            if ($$('.modal-in').length) {
                return;
            }

            const router = this.$f7router;

            try {
                this.$user.unlockTokenByPinCode(this.pinCode);
                this.$services.refreshToken();

                if (this.$settings.isAutoUpdateExchangeRatesData()) {
                    this.$services.autoRefreshLatestExchangeRates();
                }

                router.refreshPage();
            } catch (ex) {
                this.$logger.error('failed to unlock by pin code', ex);
                this.$alert('PIN code is wrong');
            }
        },
        relogin() {
            const router = this.$f7router;

            this.$user.clearTokenAndUserInfo();
            this.$settings.clearSettings();
            this.$exchangeRates.clearExchangeRates();

            router.navigate('/login', {
                clearPreviousHistory: true
            });
        }
    }
}
</script>
