<template>
    <f7-page>
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

        <PasscodeInputSheet :title="$t('Passcode')"
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
        </PasscodeInputSheet>

        <PasswordInputSheet :title="$t('Current Password')"
                            :hint="$t('Please enter your current password when disable two factor authentication')"
                            :show.sync="showInputPasswordSheetForDisable"
                            :confirm-disabled="disabling"
                            :cancel-disabled="disabling"
                            v-model="currentPasswordForDisable"
                            @password:confirm="disable">
        </PasswordInputSheet>

        <PasswordInputSheet :title="$t('Current Password')"
                            :hint="$t('Please enter your current password when regenerate two factor authentication backup codes. If you regenerate backup codes, the old codes will be invalidated.')"
                            :show.sync="showInputPasswordSheetForRegenerate"
                            :confirm-disabled="regenerating"
                            :cancel-disabled="regenerating"
                            v-model="currentPasswordForRegenerate"
                            @password:confirm="regenerateBackupCode">
        </PasswordInputSheet>

        <InformationSheet class="backup-code-sheet"
                          :title="$t('Backup Code')"
                          :hint="$t('Please copy these backup codes to safe place, the below codes can only be shown once. If these codes were lost, you can regenerate backup codes at any time.')"
                          :information="currentBackupCode"
                          :row-count="10"
                          :enable-copy="true"
                          :show.sync="showBackupCodeSheet"
                          @info:copied="onBackupCodeCopied">
        </InformationSheet>
    </f7-page>
</template>

<script>
export default {
    data() {
        return {
            status: null,
            loading: true,
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
        const router = self.$f7router;

        self.loading = true;

        self.$services.get2FAStatus().then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result || !self.$utilities.isBoolean(data.result.enable)) {
                self.$toast('Unable to get current two factor authentication status');
                router.back();
                return;
            }

            self.status = data.result.enable;
            self.loading = false;
        }).catch(error => {
            self.$logger.error('failed to get 2fa status', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                self.$toast({ error: error.response.data });
                router.back();
            } else if (!error.processed) {
                self.$toast('Unable to get current two factor authentication status');
                router.back();
            }
        });
    },
    methods: {
        enable() {
            const self = this;

            self.new2FAQRCode = '';
            self.new2FASecret = '';

            self.enabling = true;
            self.$showLoading(() => self.enabling);

            self.$services.enable2FA().then(response => {
                self.enabling = false;
                self.$hideLoading();
                const data = response.data;

                if (!data || !data.success || !data.result || !data.result.qrcode || !data.result.secret) {
                    self.$toast('Unable to enable two factor authentication');
                    return;
                }

                self.new2FAQRCode = data.result.qrcode;
                self.new2FASecret = data.result.secret;

                self.showInputPasscodeSheetForEnable = true;
            }).catch(error => {
                self.$logger.error('failed to request to enable 2fa', error);

                self.enabling = false;
                self.$hideLoading();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({error: error.response.data});
                } else if (!error.processed) {
                    self.$toast('Unable to enable two factor authentication');
                }
            });
        },
        enableConfirm() {
            const self = this;

            self.enableConfirming = true;
            self.$showLoading(() => self.enableConfirming);

            self.$services.confirmEnable2FA({
                secret: self.new2FASecret,
                passcode: self.currentPasscodeForEnable
            }).then(response => {
                self.enableConfirming = false;
                self.$hideLoading();
                const data = response.data;

                if (!data || !data.success || !data.result || !data.result.token) {
                    self.$toast('Unable to enable two factor authentication');
                    return;
                }

                self.new2FAQRCode = '';
                self.new2FASecret = '';

                self.status = true;
                self.showInputPasscodeSheetForEnable = false;

                self.$user.updateToken(data.result.token);

                if (data.result.recoveryCodes && data.result.recoveryCodes.length) {
                    self.currentBackupCode = data.result.recoveryCodes.join('\n');
                    self.showBackupCodeSheet = true;
                }
            }).catch(error => {
                self.$logger.error('failed to confirm to enable 2fa', error);

                self.enableConfirming = false;
                self.$hideLoading();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({error: error.response.data});
                } else if (!error.processed) {
                    self.$toast('Unable to enable two factor authentication');
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

            self.$services.disable2FA({
                password: password
            }).then(response => {
                self.disabling = false;
                self.$hideLoading();
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$toast('Unable to disable two factor authentication');
                    return;
                }

                self.status = false;
                self.showInputPasswordSheetForDisable = false;
                self.$toast('Two factor authentication has been disabled');
            }).catch(error => {
                self.$logger.error('failed to disable 2fa', error);

                self.disabling = false;
                self.$hideLoading();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({error: error.response.data});
                } else if (!error.processed) {
                    self.$toast('Unable to disable two factor authentication');
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

            self.$services.regenerate2FARecoveryCode({
                password: password
            }).then(response => {
                self.regenerating = false;
                self.$hideLoading();
                const data = response.data;

                if (!data || !data.success || !data.result || !data.result.recoveryCodes || !data.result.recoveryCodes.length) {
                    self.$toast('Unable to regenerate two factor authentication backup codes');
                    return;
                }

                self.showInputPasswordSheetForRegenerate = false;

                self.currentBackupCode = data.result.recoveryCodes.join('\n');
                self.showBackupCodeSheet = true;
            }).catch(error => {
                self.$logger.error('failed to regenerate 2fa recovery code', error);

                self.regenerating = false;
                self.$hideLoading();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({error: error.response.data});
                } else if (!error.processed) {
                    self.$toast('Unable to regenerate two factor authentication backup codes');
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
