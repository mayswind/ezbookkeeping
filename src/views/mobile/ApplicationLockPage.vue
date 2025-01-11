<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Application Lock')"></f7-nav-title>
        </f7-navbar>

        <f7-list strong inset dividers class="margin-top">
            <f7-list-item :title="$t('Status')" :after="$t(isEnableApplicationLock ? 'Enabled' : 'Disabled')"></f7-list-item>
            <f7-list-item v-if="isEnableApplicationLock">
                <span>{{ $t('Unlock with PIN Code') }}</span>
                <f7-toggle checked disabled></f7-toggle>
            </f7-list-item>
            <f7-list-item v-if="isEnableApplicationLock && isSupportedWebAuthn">
                <span>{{ $t('Unlock with WebAuthn') }}</span>
                <f7-toggle :checked="isEnableApplicationLockWebAuthn" @toggle:change="isEnableApplicationLockWebAuthn = $event"></f7-toggle>
            </f7-list-item>
            <f7-list-button v-if="isEnableApplicationLock" @click="disable(null)">{{ $t('Disable') }}</f7-list-button>
            <f7-list-button v-if="!isEnableApplicationLock" @click="enable(null)">{{ $t('Enable') }}</f7-list-button>
        </f7-list>

        <pin-code-input-sheet :title="$t('PIN Code')"
                              :hint="$t('Please enter a new 6-digit PIN code. The PIN code would encrypt your local data, so you need to enter it every time you open this app. If this PIN code is lost, you will need to log in again.')"
                              v-model:show="showInputPinCodeSheetForEnable"
                              v-model="currentPinCodeForEnable"
                              @pincode:confirm="enable">
        </pin-code-input-sheet>

        <pin-code-input-sheet :title="$t('PIN Code')"
                              :hint="$t('Your current PIN code is required to disable application lock.')"
                              v-model:show="showInputPinCodeSheetForDisable"
                              v-model="currentPinCodeForDisable"
                              @pincode:confirm="disable">
        </pin-code-input-sheet>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useTransactionsStore } from '@/stores/transaction.js';

import {
    isWebAuthnCompletelySupported,
    registerWebAuthnCredential
} from '@/lib/webauthn.ts';
import {
    getUserAppLockState,
    encryptToken,
    decryptToken,
    isCorrectPinCode,
    saveWebAuthnConfig,
    clearWebAuthnConfig
} from '@/lib/userstate.ts';
import logger from '@/lib/logger.ts';

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
        ...mapStores(useSettingsStore, useUserStore, useTransactionsStore),
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

                registerWebAuthnCredential(
                    getUserAppLockState(),
                    self.userStore.currentUserBasicInfo,
                ).then(({ id }) => {
                    self.$hideLoading();

                    saveWebAuthnConfig(id);
                    self.settingsStore.setEnableApplicationLockWebAuthn(true);
                    self.$toast('You have enabled WebAuthn successfully');
                }).catch(error => {
                    logger.error('failed to enable WebAuthn', error);

                    self.$hideLoading();

                    if (error.notSupported) {
                        self.$toast('WebAuth is not supported on this device');
                    } else if (error.name === 'NotAllowedError') {
                        self.$toast('User has canceled authentication');
                    } else if (error.invalid) {
                        self.$toast('Failed to enable WebAuthn');
                    } else {
                        self.$toast('User has canceled or this device does not support WebAuthn');
                    }

                    self.isEnableApplicationLockWebAuthn = false;
                    self.settingsStore.setEnableApplicationLockWebAuthn(false);
                    clearWebAuthnConfig();
                });
            } else {
                self.settingsStore.setEnableApplicationLockWebAuthn(false);
                clearWebAuthnConfig();
            }
        }
    },
    created() {
        const self = this;
        isWebAuthnCompletelySupported().then(result => {
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
                this.$alert('Invalid PIN code');
                return;
            }

            const user = this.userStore.currentUserBasicInfo;

            if (!user || !user.username) {
                this.$alert('An error occurred');
                return;
            }

            encryptToken(user.username, pinCode);
            this.settingsStore.setEnableApplicationLock(true);
            this.transactionsStore.saveTransactionDraft();

            this.settingsStore.setEnableApplicationLockWebAuthn(false);
            clearWebAuthnConfig();

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

            if (!isCorrectPinCode(pinCode)) {
                this.$alert('Incorrect PIN code');
                return;
            }

            decryptToken();
            this.settingsStore.setEnableApplicationLock(false);
            this.transactionsStore.saveTransactionDraft();

            this.settingsStore.setEnableApplicationLockWebAuthn(false);
            clearWebAuthnConfig();

            this.showInputPinCodeSheetForDisable = false;
        }
    }
}
</script>
