<template>
    <f7-page no-toolbar no-navbar no-swipeback login-screen>
        <f7-login-screen-title>{{ $t('PIN Code') }}</f7-login-screen-title>
        <f7-list>
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
    methods: {
        unlock() {
            if (!this.pinCodeValid) {
                return;
            }

            const router = this.$f7router;

            try {
                this.$user.unlockToken(this.pinCode);
                this.$services.refreshToken();

                if (this.$settings.isAutoUpdateExchangeRatesData()) {
                    this.$services.autoRefreshLatestExchangeRates();
                }

                router.refreshPage();
            } catch (ex) {
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
