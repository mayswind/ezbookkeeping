<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar :title="$t('Two-Factor Authentication')" :back-link="$t('Back')"></f7-navbar>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item title="Status" after="Unknown"></f7-list-item>
                    <f7-list-button class="disabled">Operate</f7-list-button>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading && status === true">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item :title="$t('Status')" :after="$t('Enabled')"></f7-list-item>
                    <f7-list-button :class="{ 'disabled': regenerating }" @click="regenerateBackupCode(null)">{{ $t('Regenerate Backup Codes') }}</f7-list-button>
                    <f7-list-button :class="{ 'disabled': disabling }" @click="disable(null)">{{ $t('Disable') }}</f7-list-button>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading && status === false">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item :title="$t('Status')" :after="$t('Disabled')"></f7-list-item>
                    <f7-list-button :class="{ 'disabled': enabling }" @click="enable">{{ $t('Enable') }}</f7-list-button>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <passcode-input-sheet :title="$t('Passcode')"
                              :hint="$t('Please use two factor authentication app scan the below qrcode and input current passcode')"
                              :show.sync="showInputPasscodeSheetForEnable"
                              :confirm-disabled="enableConfirming"
                              :cancel-disabled="enableConfirming"
                              v-model="currentPasscodeForEnable"
                              @passcode:confirm="enableConfirm">
            <div class="row">
                <div class="col-100 text-align-center">
                    <img alt="qrcode" width="240px" height="240px" :src="new2FAQRCode" />
                </div>
            </div>
        </passcode-input-sheet>

        <password-input-sheet :title="$t('Current Password')"
                              :hint="$t('Please enter your current password when disable two factor authentication')"
                              :show.sync="showInputPasswordSheetForDisable"
                              :confirm-disabled="disabling"
                              :cancel-disabled="disabling"
                              v-model="currentPasswordForDisable"
                              @password:confirm="disable">
        </password-input-sheet>

        <password-input-sheet :title="$t('Current Password')"
                              :hint="$t('Please enter your current password when regenerate two factor authentication backup codes. If you regenerate backup codes, the old codes will be invalidated.')"
                              :show.sync="showInputPasswordSheetForRegenerate"
                              :confirm-disabled="regenerating"
                              :cancel-disabled="regenerating"
                              v-model="currentPasswordForRegenerate"
                              @password:confirm="regenerateBackupCode">
        </password-input-sheet>

        <information-sheet class="backup-code-sheet"
                           :title="$t('Backup Code')"
                           :hint="$t('Please copy these backup codes to safe place, the below codes can only be shown once. If these codes were lost, you can regenerate backup codes at any time.')"
                           :information="currentBackupCode"
                           :row-count="10"
                           :enable-copy="true"
                           :show.sync="showBackupCodeSheet"
                           @info:copied="onBackupCodeCopied">
        </information-sheet>
    </f7-page>
</template>

<script>
export default {
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
    created() {
        const self = this;

        self.loading = true;

        self.$store.dispatch('get2FAStatus').then(response => {
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
            this.$routeBackOnError('loadingError');
        },
        enable() {
            const self = this;

            self.new2FAQRCode = '';
            self.new2FASecret = '';

            self.enabling = true;
            self.$showLoading(() => self.enabling);

            self.$store.dispatch('enable2FA').then(response => {
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

            self.$store.dispatch('confirmEnable2FA', {
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

            self.$store.dispatch('disable2FA', {
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

            self.$store.dispatch('regenerate2FARecoveryCode', {
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
.backup-code-sheet .information-content {
    font-family: monospace;
}
</style>
