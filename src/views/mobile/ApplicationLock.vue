<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Application Lock')"></f7-nav-title>
        </f7-navbar>

        <f7-card v-if="isEnableApplicationLock">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item :title="$t('Status')" :after="$t('Enabled')"></f7-list-item>
                    <f7-list-item v-if="isSupportedWebAuthn">
                        <span>{{ $t('Face ID / Touch ID') }}</span>
                        <f7-toggle :checked="isEnableApplicationLockWebAuthn" @toggle:change="isEnableApplicationLockWebAuthn = $event"></f7-toggle>
                    </f7-list-item>
                    <f7-list-button @click="disable(null)">{{ $t('Disable') }}</f7-list-button>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!isEnableApplicationLock">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item :title="$t('Status')" :after="$t('Disabled')"></f7-list-item>
                    <f7-list-button @click="enable(null)">{{ $t('Enable') }}</f7-list-button>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-sheet
            style="height:auto;"
            :opened="showInputPinCodeSheetForEnable" @sheet:closed="showInputPinCodeSheetForEnable = false; currentPinCodeForEnable = ''"
        >
            <f7-page-content>
                <div class="display-flex padding justify-content-space-between align-items-center">
                    <div style="font-size: 18px"><b>{{ $t('PIN Code') }}</b></div>
                </div>
                <div class="padding-horizontal padding-bottom">
                    <p class="no-margin-top margin-bottom-half">{{ $t('Please input a new PIN code. PIN code would encrypt your local data, so you need input this PIN code when you launch this app. If this PIN code is lost, you should re-login.') }}</p>
                    <f7-list no-hairlines class="no-margin-top margin-bottom">
                        <f7-list-item class="list-item-pincode-input">
                            <PincodeInput secure :length="6" v-model="currentPinCodeForEnable" />
                        </f7-list-item>
                    </f7-list>
                    <f7-button large fill :class="{ 'disabled': !currentPinCodeForEnableValid }" :text="$t('Continue')" @click="enable(currentPinCodeForEnable)"></f7-button>
                    <div class="margin-top text-align-center">
                        <f7-link @click="showInputPinCodeSheetForEnable = false" :text="$t('Cancel')"></f7-link>
                    </div>
                </div>
            </f7-page-content>
        </f7-sheet>

        <f7-sheet
            style="height:auto;"
            :opened="showInputPinCodeSheetForDisable" @sheet:closed="showInputPinCodeSheetForDisable = false; currentPinCodeForDisable = ''"
        >
            <f7-page-content>
                <div class="display-flex padding justify-content-space-between align-items-center">
                    <div style="font-size: 18px"><b>{{ $t('PIN Code') }}</b></div>
                </div>
                <div class="padding-horizontal padding-bottom">
                    <p class="no-margin-top margin-bottom-half">{{ $t('Please enter your current PIN code when disable application lock') }}</p>
                    <f7-list no-hairlines class="no-margin-top margin-bottom">
                        <f7-list-item class="list-item-pincode-input">
                            <PincodeInput secure :length="6" v-model="currentPinCodeForDisable" />
                        </f7-list-item>
                    </f7-list>
                    <f7-button large fill :class="{ 'disabled': !currentPinCodeForDisableValid }" :text="$t('Continue')" @click="disable(currentPinCodeForDisable)"></f7-button>
                    <div class="margin-top text-align-center">
                        <f7-link @click="showInputPinCodeSheetForDisable = false" :text="$t('Cancel')"></f7-link>
                    </div>
                </div>
            </f7-page-content>
        </f7-sheet>
    </f7-page>
</template>

<script>
export default {
    data() {
        return {
            isSupportedWebAuthn: false,
            isEnableApplicationLock: this.$settings.isEnableApplicationLock(),
            isEnableApplicationLockWebAuthn: this.$settings.isEnableApplicationLockWebAuthn(),
            currentPinCodeForEnable: '',
            currentPinCodeForDisable: '',
            showInputPinCodeSheetForEnable: false,
            showInputPinCodeSheetForDisable: false
        };
    },
    computed: {
        currentPinCodeForEnableValid() {
            return this.currentPinCodeForEnable && this.currentPinCodeForEnable.length === 6;
        },
        currentPinCodeForDisableValid() {
            return this.currentPinCodeForDisable && this.currentPinCodeForDisable.length === 6;
        },
    },
    watch: {
        isEnableApplicationLockWebAuthn: function (newValue) {
            const self = this;

            if (newValue) {
                self.$showLoading();

                self.$webauthn.registerCredential(
                    self.$user.getUserAppLockState(),
                    self.$user.getUserInfo(),
                ).then(({ id }) => {
                    self.$hideLoading();

                    self.$user.saveWebAuthnConfig(id);
                    self.$settings.setEnableApplicationLockWebAuthn(true);
                    self.$toast('You have enabled Face ID/Touch ID successfully');
                }).catch(error => {
                    self.$logger.error('failed to enable FaceID/Touch ID', error);

                    self.$hideLoading();

                    if (error.notSupported) {
                        self.$toast('This device does not support Face ID/Touch ID');
                    } else if (error.name === 'NotAllowedError') {
                        self.$toast('User has canceled authentication');
                    } else if (error.invalid) {
                        self.$toast('Failed to enable Face ID/Touch ID');
                    } else {
                        self.$toast('User has canceled or this device does not support Face ID/Touch ID');
                    }

                    self.isEnableApplicationLockWebAuthn = false;
                    self.$settings.setEnableApplicationLockWebAuthn(false);
                    self.$user.clearWebAuthnConfig();
                });
            } else {
                self.$settings.setEnableApplicationLockWebAuthn(false);
                self.$user.clearWebAuthnConfig();
            }
        }
    },
    created() {
        const self = this;
        self.$webauthn.isCompletelySupported().then(result => {
            self.isSupportedWebAuthn = result;
        });
    },
    methods: {
        enable(pinCode) {
            if (this.$settings.isEnableApplicationLock()) {
                this.$alert('Application lock has been enabled');
                return;
            }

            if (!pinCode) {
                this.showInputPinCodeSheetForEnable = true;
                return;
            }

            if (!this.currentPinCodeForEnableValid) {
                this.$alert('PIN code is invalid');
                return;
            }

            const user = this.$user.getUserInfo();

            if (!user || !user.username) {
                this.$alert('An error has occurred');
                return;
            }

            this.$user.encryptToken(user.username, pinCode);
            this.$settings.setEnableApplicationLock(true);
            this.isEnableApplicationLock = true;

            this.$settings.setEnableApplicationLockWebAuthn(false);
            this.$user.clearWebAuthnConfig();
            this.isEnableApplicationLockWebAuthn = false;

            this.showInputPinCodeSheetForEnable = false;
        },
        disable(pinCode) {
            if (!this.$settings.isEnableApplicationLock()) {
                this.$alert('Application lock is not enabled');
                return;
            }

            if (!pinCode) {
                this.showInputPinCodeSheetForDisable = true;
                return;
            }

            if (!this.$user.isCorrectPinCode(pinCode)) {
                this.$alert('PIN code is wrong');
                return;
            }

            this.$user.decryptToken();
            this.$settings.setEnableApplicationLock(false);
            this.isEnableApplicationLock = false;

            this.$settings.setEnableApplicationLockWebAuthn(false);
            this.$user.clearWebAuthnConfig();
            this.isEnableApplicationLockWebAuthn = false;

            this.showInputPinCodeSheetForDisable = false;
        }
    }
}
</script>
