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

        <f7-sheet
            style="height:auto;"
            :opened="showInputPasscodeSheetForEnable" @sheet:closed="showInputPasscodeSheetForEnable = false; currentPasscodeForEnable = ''"
        >
            <f7-page-content>
                <div class="display-flex padding justify-content-space-between align-items-center">
                    <div style="font-size: 18px"><b>{{ $t('Passcode') }}</b></div>
                </div>
                <div class="padding-horizontal padding-bottom">
                    <p class="no-margin-top margin-bottom-half">{{ $t('Please use two factor authentication app scan the below qrcode and input current passcode') }}</p>
                    <div class="row">
                        <div class="col-100 text-align-center">
                            <img alt="qrcode" width="240px" height="240px" :src="new2FAQRCode" />
                        </div>
                    </div>
                    <f7-list no-hairlines class="no-margin-top margin-bottom">
                        <f7-list-input
                            type="number"
                            outline
                            clear-button
                            :placeholder="$t('Passcode')"
                            :value="currentPasscodeForEnable"
                            @input="currentPasscodeForEnable = $event.target.value"
                        ></f7-list-input>
                    </f7-list>
                    <f7-button large fill :class="{ 'disabled': !currentPasscodeForEnable }" :text="$t('Continue')" @click="enableConfirm"></f7-button>
                    <div class="margin-top text-align-center">
                        <f7-link :class="{ 'disabled': enableConfirming }" @click="showInputPasscodeSheetForEnable = false" :text="$t('Cancel')"></f7-link>
                    </div>
                </div>
            </f7-page-content>
        </f7-sheet>

        <f7-sheet
            style="height:auto"
            :opened="showInputPasswordSheetForDisable" @sheet:closed="showInputPasswordSheetForDisable = false; currentPasswordForDisable = ''"
        >
            <f7-page-content>
                <div class="display-flex padding justify-content-space-between align-items-center">
                    <div style="font-size: 18px"><b>{{ $t('Current Password') }}</b></div>
                </div>
                <div class="padding-horizontal padding-bottom">
                    <p class="no-margin-top margin-bottom-half">{{ $t('Please enter your current password when disable two factor authentication') }}</p>
                    <f7-list no-hairlines class="no-margin-top margin-bottom">
                        <f7-list-input
                            type="password"
                            outline
                            clear-button
                            :placeholder="$t('Password')"
                            :value="currentPasswordForDisable"
                            @input="currentPasswordForDisable = $event.target.value"
                        ></f7-list-input>
                    </f7-list>
                    <f7-button large fill :class="{ 'disabled': !currentPasswordForDisable }" :text="$t('Continue')" @click="disable(currentPasswordForDisable)"></f7-button>
                    <div class="margin-top text-align-center">
                        <f7-link :class="{ 'disabled': disabling }" @click="showInputPasswordSheetForDisable = false" :text="$t('Cancel')"></f7-link>
                    </div>
                </div>
            </f7-page-content>
        </f7-sheet>

        <f7-sheet
            style="height:auto"
            :opened="showInputPasswordSheetForRegenerate" @sheet:closed="showInputPasswordSheetForRegenerate = false; currentPasswordForRegenerate= ''"
        >
            <f7-page-content>
                <div class="display-flex padding justify-content-space-between align-items-center">
                    <div style="font-size: 18px"><b>{{ $t('Current Password') }}</b></div>
                </div>
                <div class="padding-horizontal padding-bottom">
                    <p class="no-margin-top margin-bottom-half">{{ $t('Please enter your current password when regenerate two factor authentication backup codes. If you regenerate backup codes, the old codes will be invalidated.') }}</p>
                    <f7-list no-hairlines class="no-margin-top margin-bottom">
                        <f7-list-input
                            type="password"
                            outline
                            clear-button
                            :placeholder="$t('Password')"
                            :value="currentPasswordForRegenerate"
                            @input="currentPasswordForRegenerate = $event.target.value"
                        ></f7-list-input>
                    </f7-list>
                    <f7-button large fill :class="{ 'disabled': !currentPasswordForRegenerate }" :text="$t('Continue')" @click="regenerateBackupCode(currentPasswordForRegenerate)"></f7-button>
                    <div class="margin-top text-align-center">
                        <f7-link :class="{ 'disabled': regenerating }" @click="showInputPasswordSheetForRegenerate = false" :text="$t('Cancel')"></f7-link>
                    </div>
                </div>
            </f7-page-content>
        </f7-sheet>

        <f7-sheet
            style="height:auto"
            :opened="showBackupCodeSheet" @sheet:closed="showBackupCodeSheet = false; currentBackupCode = ''"
        >
            <f7-page-content>
                <div class="display-flex padding justify-content-space-between align-items-center">
                    <div style="font-size: 18px"><b>{{ $t('Backup Code') }}</b></div>
                </div>
                <div class="padding-horizontal padding-bottom">
                    <p class="no-margin-top margin-bottom-half">
                        <span>{{ $t('Please copy these backup codes to safe place, the below codes can only be shown once. If these codes were lost, you can regenerate backup codes at any time.') }}</span>
                        <f7-link icon-only icon-f7="doc_on_doc" icon-size="16px" class="icon-after-text"
                                 v-clipboard:copy="currentBackupCode" v-clipboard:success="onBackupCodeCopied"></f7-link>
                    </p>
                    <textarea class="full-line" rows="10" readonly="readonly" v-model="currentBackupCode"></textarea>
                    <div class="margin-top text-align-center">
                        <f7-link @click="showBackupCodeSheet = false" :text="$t('Close')"></f7-link>
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
