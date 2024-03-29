<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Application Lock')"></f7-nav-title>
        </f7-navbar>

        <f7-list strong inset dividers class="margin-top">
            <f7-list-item :title="$t('Status')" :after="$t(isEnableApplicationLock ? 'Enabled' : 'Disabled')"></f7-list-item>
            <f7-list-item v-if="isEnableApplicationLock">
                <span>{{ $t('Unlock By PIN Code') }}</span>
                <f7-toggle checked disabled></f7-toggle>
            </f7-list-item>
            <f7-list-item v-if="isEnableApplicationLock && isSupportedWebAuthn">
                <span>{{ $t('Unlock By WebAuthn') }}</span>
                <f7-toggle :checked="isEnableApplicationLockWebAuthn" @toggle:change="isEnableApplicationLockWebAuthn = $event"></f7-toggle>
            </f7-list-item>
            <f7-list-button v-if="isEnableApplicationLock" @click="disable(null)">{{ $t('Disable') }}</f7-list-button>
            <f7-list-button v-if="!isEnableApplicationLock" @click="enable(null)">{{ $t('Enable') }}</f7-list-button>
        </f7-list>

        <pin-code-input-sheet :title="$t('PIN Code')"
                              :hint="$t('Please input a new 6-digit PIN code. PIN code would encrypt your local data, so you need input this PIN code when you launch this app. If this PIN code is lost, you should re-login.')"
                              v-model:show="showInputPinCodeSheetForEnable"
                              v-model="currentPinCodeForEnable"
                              @pincode:confirm="enable">
        </pin-code-input-sheet>

        <pin-code-input-sheet :title="$t('PIN Code')"
                              :hint="$t('Please enter your current PIN code when disable application lock.')"
                              v-model:show="showInputPinCodeSheetForDisable"
                              v-model="currentPinCodeForDisable"
                              @pincode:confirm="disable">
        </pin-code-input-sheet>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';

import logger from '@/lib/logger.js';
import webauthn from '@/lib/webauthn.js';

export default {
    data() {
        return {
            isSupportedWebAuthn: false,
            currentPinCodeForEnable: '',
            currentPinCodeForDisable: '',
            showInputPinCodeSheetForEnable: false,
            showInputPinCodeSheetForDisable: false
        };
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore),
        isEnableApplicationLock: {
            get: function () {
                return this.settingsStore.appSettings.applicationLock;
            },
            set: function (value) {
                this.settingsStore.setEnableApplicationLock(value);
            }
        },
        isEnableApplicationLockWebAuthn: {
            get: function () {
                return this.settingsStore.appSettings.applicationLockWebAuthn;
            },
            set: function (value) {
                this.settingsStore.setEnableApplicationLockWebAuthn(value);
            }
        }
    },
    watch: {
        isEnableApplicationLockWebAuthn: function (newValue) {
            const self = this;

            if (newValue) {
                self.$showLoading();

                webauthn.registerCredential(
                    self.$user.getUserAppLockState(),
                    self.userStore.currentUserInfo,
                ).then(({ id }) => {
                    self.$hideLoading();

                    self.$user.saveWebAuthnConfig(id);
                    self.settingsStore.setEnableApplicationLockWebAuthn(true);
                    self.$toast('You have enabled WebAuthn successfully');
                }).catch(error => {
                    logger.error('failed to enable WebAuthn', error);

                    self.$hideLoading();

                    if (error.notSupported) {
                        self.$toast('This device does not support WebAuthn');
                    } else if (error.name === 'NotAllowedError') {
                        self.$toast('User has canceled authentication');
                    } else if (error.invalid) {
                        self.$toast('Failed to enable WebAuthn');
                    } else {
                        self.$toast('User has canceled or this device does not support WebAuthn');
                    }

                    self.isEnableApplicationLockWebAuthn = false;
                    self.settingsStore.setEnableApplicationLockWebAuthn(false);
                    self.$user.clearWebAuthnConfig();
                });
            } else {
                self.settingsStore.setEnableApplicationLockWebAuthn(false);
                self.$user.clearWebAuthnConfig();
            }
        }
    },
    created() {
        const self = this;
        webauthn.isCompletelySupported().then(result => {
            self.isSupportedWebAuthn = result;
        });
    },
    methods: {
        enable(pinCode) {
            if (this.settingsStore.appSettings.applicationLock) {
                this.$alert('Application lock has been enabled');
                return;
            }

            if (!pinCode) {
                this.showInputPinCodeSheetForEnable = true;
                return;
            }

            if (!this.currentPinCodeForEnable || this.currentPinCodeForEnable.length !== 6) {
                this.$alert('PIN code is invalid');
                return;
            }

            const user = this.userStore.currentUserInfo;

            if (!user || !user.username) {
                this.$alert('An error has occurred');
                return;
            }

            this.$user.encryptToken(user.username, pinCode);
            this.settingsStore.setEnableApplicationLock(true);

            this.settingsStore.setEnableApplicationLockWebAuthn(false);
            this.$user.clearWebAuthnConfig();

            this.showInputPinCodeSheetForEnable = false;
        },
        disable(pinCode) {
            if (!this.settingsStore.appSettings.applicationLock) {
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
            this.settingsStore.setEnableApplicationLock(false);

            this.settingsStore.setEnableApplicationLockWebAuthn(false);
            this.$user.clearWebAuthnConfig();

            this.showInputPinCodeSheetForDisable = false;
        }
    }
}
</script>
