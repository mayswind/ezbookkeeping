<template>
    <v-row>
        <v-col cols="12">
            <v-card :class="{ 'disabled': loading }">
                <template #title>
                    <span>{{ $t('Two-Factor Authentication') }}</span>
                    <v-progress-circular indeterminate size="20" class="ml-3" v-if="loading"></v-progress-circular>
                </template>

                <v-card-text class="pb-0">
                    <v-skeleton-loader class="skeleton-no-margin pt-2 pb-5" type="text" style="width: 150px" :loading="true" v-if="loading"></v-skeleton-loader>
                    <p class="text-body-1 font-weight-semibold" v-if="!loading && !new2FAQRCode">
                        {{ status === true ? $t('Two-factor authentication is already enabled.') : $t('Two-factor authentication is not enabled yet.') }}
                    </p>
                    <p class="text-body-1" v-if="!loading && new2FAQRCode">
                        {{ $t('Please use a two-factor authentication app to scan the qrcode below and enter the current passcode.') }}
                    </p>
                    <p class="text-body-1" v-if="!loading && status === true">
                        {{ $t('Your current password is required to disable two-factor authentication or regenerate backup codes for two-factor authentication. If you regenerate backup codes, the previous ones will become invalid.') }}
                    </p>
                </v-card-text>

                <v-card-text v-if="status === false && new2FAQRCode">
                    <v-img alt="qrcode" class="img-qrcode" :src="new2FAQRCode" />
                    <v-row class="mb-3">
                        <v-col cols="12" md="3">
                            <v-text-field
                                type="number"
                                autocomplete="one-time-code"
                                clearable variant="underlined"
                                :disabled="loading || enabling || enableConfirming || disabling"
                                :placeholder="$t('Passcode')"
                                v-model="currentPasscode"
                                @keyup.enter="enableConfirm"
                            />
                        </v-col>
                    </v-row>
                </v-card-text>

                <v-card-text class="pb-0" v-if="status === true">
                    <v-row class="mb-3">
                        <v-col cols="12" md="6">
                            <v-text-field
                                autocomplete="current-password"
                                clearable variant="underlined"
                                :disabled="loading || enabling || enableConfirming || disabling"
                                :placeholder="$t('Current Password')"
                                :type="isCurrentPasswordVisible ? 'text' : 'password'"
                                :append-inner-icon="isCurrentPasswordVisible ? icons.eyeSlash : icons.eye"
                                v-model="currentPassword"
                                @click:append-inner="isCurrentPasswordVisible = !isCurrentPasswordVisible"
                            />
                        </v-col>
                    </v-row>
                </v-card-text>

                <v-card-text>
                    <v-row>
                        <v-col cols="12" class="d-flex flex-wrap gap-4">
                            <v-btn :disabled="!currentPassword || loading || disabling " v-if="status === true" @click="disable">
                                {{ $t('Disable Two-Factor Authentication') }}
                                <v-progress-circular indeterminate size="22" class="ml-2" v-if="disabling"></v-progress-circular>
                            </v-btn>
                            <v-btn :disabled="!currentPassword || loading || regenerating" v-if="status === true" @click="regenerateBackupCode()">
                                {{ $t('Regenerate Backup Codes') }}
                                <v-progress-circular indeterminate size="22" class="ml-2" v-if="regenerating"></v-progress-circular>
                            </v-btn>
                            <v-btn :disabled="loading || enabling" v-if="status === false && !new2FAQRCode" @click="enable">
                                {{ $t('Enable Two-Factor Authentication') }}
                                <v-progress-circular indeterminate size="22" class="ml-2" v-if="enabling"></v-progress-circular>
                            </v-btn>
                            <v-btn :disabled="!currentPasscode || loading || enableConfirming" v-if="status === false && new2FAQRCode" @click="enableConfirm">
                                {{ $t('Continue') }}
                                <v-progress-circular indeterminate size="22" class="ml-2" v-if="enableConfirming"></v-progress-circular>
                            </v-btn>
                        </v-col>
                    </v-row>
                </v-card-text>
            </v-card>
        </v-col>

        <v-col cols="12">
            <v-card v-if="currentBackupCode">
                <template #title>
                    <span>{{ $t('Backup Code') }}</span>
                    <v-btn id="copy-to-clipboard-icon" ref="copyToClipboardIcon"
                           density="compact" color="default" variant="text" size="24"
                           class="ml-2" :icon="true">
                        <v-icon :icon="icons.copy" size="20" />
                        <v-tooltip activator="parent">{{ $t('Copy') }}</v-tooltip>
                    </v-btn>
                </template>

                <v-card-text>
                    <p class="text-body-1" v-if="status === true">
                        {{ $t('Please copy these backup codes to safe place, the following backup codes will be displayed only once. If these codes were lost, you can regenerate them at any time.') }}
                    </p>
                    <v-textarea class="backup-code" readonly="readonly" :rows="10" :value="currentBackupCode"/>
                </v-card-text>
            </v-card>
        </v-col>
    </v-row>

    <snack-bar ref="snackbar" />
</template>

<script>
import { mapStores } from 'pinia';
import { useTwoFactorAuthStore } from '@/stores/twoFactorAuth.js';

import { makeButtonCopyToClipboard, changeClipboardObjectText } from '@/lib/misc.js';

import {
    mdiEyeOutline,
    mdiEyeOffOutline,
    mdiContentCopy
} from '@mdi/js';

export default {
    expose: [
        'reset'
    ],
    data() {
        return {
            status: null,
            loading: true,
            new2FASecret: '',
            new2FAQRCode: '',
            currentPassword: '',
            isCurrentPasswordVisible: false,
            currentPasscode: '',
            currentBackupCode: '',
            enabling: false,
            enableConfirming: false,
            disabling: false,
            regenerating: false,
            clipboardHolder: null,
            icons: {
                eye: mdiEyeOutline,
                eyeSlash: mdiEyeOffOutline,
                copy: mdiContentCopy
            }
        };
    },
    computed: {
        ...mapStores(useTwoFactorAuthStore),
    },
    created() {
        const self = this;

        self.loading = true;

        self.twoFactorAuthStore.get2FAStatus().then(response => {
            self.status = response.enable;
            self.loading = false;
        }).catch(error => {
            self.loading = false;

            if (!error.processed) {
                self.$refs.snackbar.showError(error);
            }
        });
    },
    watch: {
        'currentBackupCode': function (newValue) {
            if (this.clipboardHolder) {
                changeClipboardObjectText(this.clipboardHolder, newValue);
            }
        }
    },
    methods: {
        reset() {
            this.new2FASecret = '';
            this.new2FAQRCode = '';
            this.currentPassword = '';
            this.isCurrentPasswordVisible = false;
            this.currentPasscode = '';
            this.currentBackupCode = '';
            this.enabling = false;
            this.enableConfirming = false;
            this.disabling = false;
            this.regenerating = false;
        },
        enable() {
            const self = this;

            self.new2FAQRCode = '';
            self.new2FASecret = '';
            self.currentBackupCode = '';

            self.enabling = true;

            self.twoFactorAuthStore.enable2FA().then(response => {
                self.enabling = false;

                self.new2FAQRCode = response.qrcode;
                self.new2FASecret = response.secret;
            }).catch(error => {
                self.enabling = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        enableConfirm() {
            const self = this;

            if (!self.currentPasscode) {
                self.$refs.snackbar.showMessage('Passcode cannot be blank');
                return;
            }

            if (self.enableConfirming) {
                return;
            }

            const password = self.currentPasscode;

            self.currentBackupCode = '';
            self.currentPasscode = '';

            self.enableConfirming = true;

            self.twoFactorAuthStore.confirmEnable2FA({
                secret: self.new2FASecret,
                passcode: password
            }).then(response => {
                self.enableConfirming = false;

                self.new2FAQRCode = '';
                self.new2FASecret = '';

                self.status = true;

                if (response.recoveryCodes && response.recoveryCodes.length) {
                    self.currentBackupCode = response.recoveryCodes.join('\n');
                }

                self.$nextTick(() => {
                    self.makeCopyToClipboardClickable();
                });
            }).catch(error => {
                self.enableConfirming = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        disable() {
            const self = this;

            if (!self.currentPassword) {
                self.$refs.snackbar.showMessage('Current password cannot be blank');
                return;
            }

            if (self.disabling) {
                return;
            }

            const password = self.currentPassword;

            self.currentBackupCode = '';
            self.currentPassword = '';

            self.disabling = true;

            self.twoFactorAuthStore.disable2FA({
                password: password
            }).then(() => {
                self.disabling = false;

                self.status = false;
                self.$refs.snackbar.showMessage('Two-factor authentication has been disabled');
            }).catch(error => {
                self.disabling = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        regenerateBackupCode() {
            const self = this;

            if (!self.currentPassword) {
                self.$refs.snackbar.showMessage('Current password cannot be blank');
                return;
            }

            if (self.regenerating) {
                return;
            }

            const password = self.currentPassword;

            self.currentBackupCode = '';
            self.currentPassword = '';

            self.regenerating = true;

            self.twoFactorAuthStore.regenerate2FARecoveryCode({
                password: password
            }).then(response => {
                self.regenerating = false;

                self.currentBackupCode = response.recoveryCodes.join('\n');

                self.$nextTick(() => {
                    self.makeCopyToClipboardClickable();
                });
            }).catch(error => {
                self.regenerating = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        makeCopyToClipboardClickable() {
            const self = this;

            if (self.clipboardHolder) {
                return;
            }

            if (self.$refs.copyToClipboardIcon) {
                self.clipboardHolder = makeButtonCopyToClipboard({
                    el: '#copy-to-clipboard-icon',
                    text: self.currentBackupCode,
                    successCallback: function () {
                        self.$refs.snackbar.showMessage('Backup codes copied');
                    }
                });
            }
        }
    }
};
</script>

<style>
.img-qrcode {
    width: 240px;
    height: 240px
}

.backup-code {
    font-family: monospace;
}

.backup-code textarea {
    resize: none;
}
</style>
