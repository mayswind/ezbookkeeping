<template>
    <v-row>
        <v-col cols="12">
            <v-card :title="$t('Application Lock')">
                <v-card-text class="pb-0">
                    <p class="text-body-1 font-weight-semibold" v-if="!isEnableApplicationLock">
                        {{ $t('Application lock is not enabled') }}
                    </p>
                    <p class="text-body-1" v-if="isEnableApplicationLock">
                        {{ $t('Application lock has been enabled') }}
                    </p>
                </v-card-text>

                <v-card-text v-if="isEnableApplicationLock">
                    <v-switch :disabled="true"
                              :label="$t('Unlock with PIN Code')"
                              v-model="isEnableApplicationLock"/>
                    <v-switch :label="$t('Unlock with WebAuthn')"
                              :loading="enablingWebAuthn"
                              v-model="isEnableApplicationLockWebAuthn"/>
                </v-card-text>

                <v-card-text class="pb-0">
                    <p class="text-body-1 font-weight-semibold" v-if="!isEnableApplicationLock">
                        {{ $t('Please enter a new 6-digit PIN code. The PIN code would encrypt your local data, so you need to enter it every time you open this app. If this PIN code is lost, you will need to log in again.') }}
                    </p>
                    <p class="text-body-1 font-weight-semibold" v-if="isEnableApplicationLock">
                        {{ $t('Your current PIN code is required to disable application lock.') }}
                    </p>
                </v-card-text>

                <v-card-text class="pb-0">
                    <v-row class="mb-3">
                        <v-col cols="12" md="12">
                            <div style="max-width: 428px">
                                <pin-code-input :secure="true" :length="6" v-model="pinCode" />
                            </div>
                        </v-col>
                    </v-row>
                </v-card-text>

                <v-card-text>
                    <v-row>
                        <v-col cols="12" class="d-flex flex-wrap gap-4">
                            <v-btn :disabled="!pinCodeValid"
                                   v-if="!isEnableApplicationLock" @click="enable">
                                {{ $t('Enable Application Lock') }}
                            </v-btn>
                            <v-btn :disabled="!pinCodeValid"
                                   v-if="isEnableApplicationLock" @click="disable">
                                {{ $t('Disable Application Lock') }}
                            </v-btn>
                        </v-col>
                    </v-row>
                </v-card-text>
            </v-card>
        </v-col>
    </v-row>

    <snack-bar ref="snackbar" />
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';
import { useTransactionsStore } from '@/stores/transaction.js';

import logger from '@/lib/logger.js';
import webauthn from '@/lib/webauthn.js';

export default {
    data() {
        return {
            isSupportedWebAuthn: false,
            enablingWebAuthn: false,
            pinCode: ''
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
        },
        pinCodeValid() {
            return this.pinCode && this.pinCode.length === 6;
        }
    },
    watch: {
        isEnableApplicationLockWebAuthn: function (newValue) {
            const self = this;

            if (newValue) {
                self.enablingWebAuthn = true;

                webauthn.registerCredential(
                    self.$user.getUserAppLockState(),
                    self.userStore.currentUserBasicInfo,
                ).then(({ id }) => {
                    self.enablingWebAuthn = false;

                    self.$user.saveWebAuthnConfig(id);
                    self.settingsStore.setEnableApplicationLockWebAuthn(true);
                    self.$refs.snackbar.showMessage('You have enabled WebAuthn successfully');
                }).catch(error => {
                    logger.error('failed to enable WebAuthn', error);

                    self.enablingWebAuthn = false;

                    if (error.notSupported) {
                        self.$refs.snackbar.showMessage('WebAuth is not supported on this device');
                    } else if (error.name === 'NotAllowedError') {
                        self.$refs.snackbar.showMessage('User has canceled authentication');
                    } else if (error.invalid) {
                        self.$refs.snackbar.showMessage('Failed to enable WebAuthn');
                    } else {
                        self.$refs.snackbar.showMessage('User has canceled or this device does not support WebAuthn');
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
        enable() {
            if (this.settingsStore.appSettings.applicationLock) {
                this.$refs.snackbar.showMessage('Application lock has been enabled');
                return;
            }

            if (!this.pinCode || this.pinCode.length !== 6) {
                this.pinCode = '';
                this.$refs.snackbar.showMessage('Invalid PIN code');
                return;
            }

            const user = this.userStore.currentUserBasicInfo;

            if (!user || !user.username) {
                this.pinCode = '';
                this.$refs.snackbar.showMessage('An error occurred');
                return;
            }

            this.$user.encryptToken(user.username, this.pinCode);
            this.settingsStore.setEnableApplicationLock(true);
            this.transactionsStore.saveTransactionDraft();

            this.settingsStore.setEnableApplicationLockWebAuthn(false);
            this.$user.clearWebAuthnConfig();

            this.pinCode = '';
        },
        disable() {
            if (!this.settingsStore.appSettings.applicationLock) {
                this.$refs.snackbar.showMessage('Application lock is not enabled');
                return;
            }

            if (!this.$user.isCorrectPinCode(this.pinCode)) {
                this.pinCode = '';
                this.$refs.snackbar.showMessage('Incorrect PIN code');
                return;
            }

            this.pinCode = '';

            this.$user.decryptToken();
            this.settingsStore.setEnableApplicationLock(false);
            this.transactionsStore.saveTransactionDraft();

            this.settingsStore.setEnableApplicationLockWebAuthn(false);
            this.$user.clearWebAuthnConfig();
        }
    }
}
</script>
