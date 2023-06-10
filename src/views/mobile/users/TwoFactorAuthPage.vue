<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar :title="$t('Two-Factor Authentication')" :back-link="$t('Back')"></f7-navbar>

        <f7-list strong inset dividers class="margin-top skeleton-text" v-if="loading">
            <f7-list-item title="Status" after="Unknown"></f7-list-item>
            <f7-list-button class="disabled">Operate</f7-list-button>
        </f7-list>

        <f7-list strong inset dividers class="margin-top" v-else-if="!loading">
            <f7-list-item :title="$t('Status')" :after="$t(status ? 'Enabled' : 'Disabled')"></f7-list-item>
            <f7-list-button :class="{ 'disabled': regenerating }" v-if="status === true" @click="regenerateBackupCode(null)">{{ $t('Regenerate Backup Codes') }}</f7-list-button>
            <f7-list-button :class="{ 'disabled': disabling }" v-if="status === true" @click="disable(null)">{{ $t('Disable') }}</f7-list-button>
            <f7-list-button :class="{ 'disabled': enabling }" v-if="status === false" @click="enable">{{ $t('Enable') }}</f7-list-button>
        </f7-list>

        <passcode-input-sheet :title="$t('Enable Two-Factor Authentication')"
                              :hint="$t('Please use two factor authentication app scan the below qrcode and input current passcode')"
                              :confirm-disabled="enableConfirming"
                              :cancel-disabled="enableConfirming"
                              v-model:show="showInputPasscodeSheetForEnable"
                              v-model="currentPasscodeForEnable"
                              @passcode:confirm="enableConfirm">
            <div class="text-align-center">
                <img alt="qrcode" class="img-qrcode" :src="new2FAQRCode" />
            </div>
        </passcode-input-sheet>

        <password-input-sheet :title="$t('Disable Two-Factor Authentication')"
                              :hint="$t('Please enter your current password when disable two factor authentication')"
                              :confirm-disabled="disabling"
                              :cancel-disabled="disabling"
                              v-model:show="showInputPasswordSheetForDisable"
                              v-model="currentPasswordForDisable"
                              @password:confirm="disable">
        </password-input-sheet>

        <password-input-sheet :title="$t('Regenerate Backup Codes')"
                              :hint="$t('Please enter your current password when regenerate two factor authentication backup codes. If you regenerate backup codes, the old codes will be invalidated.')"
                              :confirm-disabled="regenerating"
                              :cancel-disabled="regenerating"
                              v-model:show="showInputPasswordSheetForRegenerate"
                              v-model="currentPasswordForRegenerate"
                              @password:confirm="regenerateBackupCode">
        </password-input-sheet>

        <information-sheet class="backup-code-sheet"
                           :title="$t('Backup Code')"
                           :hint="$t('Please copy these backup codes to safe place, the below codes can only be shown once. If these codes were lost, you can regenerate backup codes at any time.')"
                           :information="currentBackupCode"
                           :row-count="10"
                           :enable-copy="true"
                           v-model:show="showBackupCodeSheet"
                           @info:copied="onBackupCodeCopied">
        </information-sheet>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useTwoFactorAuthStore } from '@/stores/twoFactorAuth.js';

export default {
    props: [
        'f7router'
    ],
    data() {
        return {
            status: null,
            loading: true,
            loadingError: null,
            new2FASecret: '',
            new2FAQRCode: '',
            currentPasscodeForEnable: '',
            currentPasswordForDisable: '',
            currentPasswordForRegenerate: '',
            currentBackupCode: '',
            enabling: false,
            enableConfirming: false,
            disabling: false,
            regenerating: false,
            showInputPasscodeSheetForEnable: false,
            showInputPasswordSheetForDisable: false,
            showInputPasswordSheetForRegenerate: false,
            showBackupCodeSheet: false
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
            if (error.processed) {
                self.loading = false;
            } else {
                self.loadingError = error;
                self.$toast(error.message || error);
            }
        });
    },
    methods: {
        onPageAfterIn() {
            this.$routeBackOnError(this.f7router, 'loadingError');
        },
        enable() {
            const self = this;

            self.new2FAQRCode = '';
            self.new2FASecret = '';

            self.enabling = true;
            self.$showLoading(() => self.enabling);

            self.twoFactorAuthStore.enable2FA().then(response => {
                self.enabling = false;
                self.$hideLoading();

                self.new2FAQRCode = response.qrcode;
                self.new2FASecret = response.secret;

                self.showInputPasscodeSheetForEnable = true;
            }).catch(error => {
                self.enabling = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        enableConfirm() {
            const self = this;

            self.enableConfirming = true;
            self.$showLoading(() => self.enableConfirming);

            self.twoFactorAuthStore.confirmEnable2FA({
                secret: self.new2FASecret,
                passcode: self.currentPasscodeForEnable
            }).then(response => {
                self.enableConfirming = false;
                self.$hideLoading();

                self.new2FAQRCode = '';
                self.new2FASecret = '';

                self.status = true;
                self.showInputPasscodeSheetForEnable = false;

                if (response.recoveryCodes && response.recoveryCodes.length) {
                    self.currentBackupCode = response.recoveryCodes.join('\n');
                    self.showBackupCodeSheet = true;
                }
            }).catch(error => {
                self.enableConfirming = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        disable(password) {
            const self = this;

            if (!password) {
                self.currentPasswordForDisable = '';
                self.showInputPasswordSheetForDisable = true;
                return;
            }

            self.disabling = true;
            self.$showLoading(() => self.disabling);

            self.twoFactorAuthStore.disable2FA({
                password: password
            }).then(() => {
                self.disabling = false;
                self.$hideLoading();

                self.status = false;
                self.showInputPasswordSheetForDisable = false;
                self.$toast('Two factor authentication has been disabled');
            }).catch(error => {
                self.disabling = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        regenerateBackupCode(password) {
            const self = this;

            if (!password) {
                self.currentPasswordForRegenerate = '';
                self.showInputPasswordSheetForRegenerate = true;
                return;
            }

            self.regenerating = true;
            self.$showLoading(() => self.regenerating);

            self.twoFactorAuthStore.regenerate2FARecoveryCode({
                password: password
            }).then(response => {
                self.regenerating = false;
                self.$hideLoading();

                self.showInputPasswordSheetForRegenerate = false;

                self.currentBackupCode = response.recoveryCodes.join('\n');
                self.showBackupCodeSheet = true;
            }).catch(error => {
                self.regenerating = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        onBackupCodeCopied() {
            this.$toast('Backup codes copied');
        }
    }
};
</script>

<style>
.img-qrcode {
    width: 240px;
    height: 240px
}

.backup-code-sheet .information-content {
    font-family: monospace;
}
</style>
