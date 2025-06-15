<template>
    <v-row>
        <v-col cols="12">
            <v-card :class="{ 'disabled': loading }">
                <template #title>
                    <span>{{ tt('Two-Factor Authentication') }}</span>
                    <v-progress-circular indeterminate size="20" class="ml-3" v-if="loading"></v-progress-circular>
                </template>

                <v-card-text class="pb-0">
                    <v-skeleton-loader class="skeleton-no-margin pt-2 pb-5" type="text" style="width: 150px" :loading="true" v-if="loading"></v-skeleton-loader>
                    <p class="text-body-1 font-weight-semibold" v-if="!loading && !new2FAQRCode">
                        {{ status === true ? tt('Two-factor authentication is already enabled.') : tt('Two-factor authentication is not enabled yet.') }}
                    </p>
                    <p class="text-body-1" v-if="!loading && new2FAQRCode">
                        {{ tt('Please use a two-factor authentication app to scan the qrcode below and enter the current passcode.') }}
                    </p>
                    <p class="text-body-1" v-if="!loading && status === true">
                        {{ tt('Your current password is required to disable two-factor authentication or regenerate backup codes for two-factor authentication. If you regenerate backup codes, the previous ones will become invalid.') }}
                    </p>
                </v-card-text>

                <v-card-text v-if="status === false && new2FAQRCode">
                    <v-img alt="qrcode" class="img-qrcode" :src="new2FAQRCode" />
                    <v-row class="mb-3">
                        <v-col cols="12" md="3">
                            <v-text-field
                                type="number"
                                autocomplete="one-time-code"
                                variant="underlined"
                                :disabled="loading || enabling || enableConfirming || disabling"
                                :placeholder="tt('Passcode')"
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
                                type="password"
                                variant="underlined"
                                :disabled="loading || enabling || enableConfirming || disabling"
                                :placeholder="tt('Current Password')"
                                v-model="currentPassword"
                            />
                        </v-col>
                    </v-row>
                </v-card-text>

                <v-card-text>
                    <v-row>
                        <v-col cols="12" class="d-flex flex-wrap gap-4">
                            <v-btn :disabled="!currentPassword || loading || disabling " v-if="status === true" @click="disable">
                                {{ tt('Disable Two-Factor Authentication') }}
                                <v-progress-circular indeterminate size="22" class="ml-2" v-if="disabling"></v-progress-circular>
                            </v-btn>
                            <v-btn :disabled="!currentPassword || loading || regenerating" v-if="status === true" @click="regenerateBackupCode()">
                                {{ tt('Regenerate Backup Codes') }}
                                <v-progress-circular indeterminate size="22" class="ml-2" v-if="regenerating"></v-progress-circular>
                            </v-btn>
                            <v-btn :disabled="loading || enabling" v-if="status === false && !new2FAQRCode" @click="enable">
                                {{ tt('Enable Two-Factor Authentication') }}
                                <v-progress-circular indeterminate size="22" class="ml-2" v-if="enabling"></v-progress-circular>
                            </v-btn>
                            <v-btn :disabled="!currentPasscode || loading || enableConfirming" v-if="status === false && new2FAQRCode" @click="enableConfirm">
                                {{ tt('Continue') }}
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
                    <span>{{ tt('Backup Code') }}</span>
                    <v-btn density="compact" color="default" variant="text" size="24"
                           class="ml-2" :icon="true" @click="copyBackupCodes">
                        <v-icon :icon="mdiContentCopy" size="20" />
                        <v-tooltip activator="parent">{{ tt('Copy') }}</v-tooltip>
                    </v-btn>
                </template>

                <v-card-text>
                    <p class="text-body-1" v-if="status === true">
                        {{ tt('Please copy these backup codes to safe place, the following backup codes will be displayed only once. If these codes were lost, you can regenerate them at any time.') }}
                    </p>
                    <v-textarea class="backup-code" :readonly="true" :rows="10" :value="currentBackupCode"/>
                </v-card-text>
            </v-card>
        </v-col>
    </v-row>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTwoFactorAuthStore } from '@/stores/twoFactorAuth.ts';

import { copyTextToClipboard } from '@/lib/ui/common.ts';

import {
    mdiContentCopy
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();

const twoFactorAuthStore = useTwoFactorAuthStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const status = ref<boolean | null>(null);
const loading = ref<boolean>(true);
const new2FASecret = ref<string>('');
const new2FAQRCode = ref<string>('');
const currentPassword = ref<string>('');
const currentPasscode = ref<string>('');
const currentBackupCode = ref<string>('');
const enabling = ref<boolean>(false);
const enableConfirming = ref<boolean>(false);
const disabling = ref<boolean>(false);
const regenerating = ref<boolean>(false);

function init(): void {
    loading.value = true;

    twoFactorAuthStore.get2FAStatus().then(response => {
        status.value = response.enable;
        loading.value = false;
    }).catch(error => {
        loading.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function enable(): void {
    new2FAQRCode.value = '';
    new2FASecret.value = '';
    currentBackupCode.value = '';

    enabling.value = true;

    twoFactorAuthStore.enable2FA().then(response => {
        enabling.value = false;

        new2FAQRCode.value = response.qrcode;
        new2FASecret.value = response.secret;
    }).catch(error => {
        enabling.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function enableConfirm(): void {
    if (!currentPasscode.value) {
        snackbar.value?.showMessage('Passcode cannot be blank');
        return;
    }

    if (enableConfirming.value) {
        return;
    }

    const password = currentPasscode.value;

    currentBackupCode.value = '';
    currentPasscode.value = '';

    enableConfirming.value = true;

    twoFactorAuthStore.confirmEnable2FA({
        secret: new2FASecret.value,
        passcode: password
    }).then(response => {
        enableConfirming.value = false;

        new2FAQRCode.value = '';
        new2FASecret.value = '';

        status.value = true;

        if (response.recoveryCodes && response.recoveryCodes.length) {
            currentBackupCode.value = response.recoveryCodes.join('\n');
        }
    }).catch(error => {
        enableConfirming.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function disable(): void {
    if (!currentPassword.value) {
        snackbar.value?.showMessage('Current password cannot be blank');
        return;
    }

    if (disabling.value) {
        return;
    }

    const password = currentPassword.value;

    currentBackupCode.value = '';
    currentPassword.value = '';

    disabling.value = true;

    twoFactorAuthStore.disable2FA({
        password: password
    }).then(() => {
        disabling.value = false;

        status.value = false;
        snackbar.value?.showMessage('Two-factor authentication has been disabled');
    }).catch(error => {
        disabling.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function regenerateBackupCode(): void {
    if (!currentPassword.value) {
        snackbar.value?.showMessage('Current password cannot be blank');
        return;
    }

    if (regenerating.value) {
        return;
    }

    const password = currentPassword.value;

    currentBackupCode.value = '';
    currentPassword.value = '';

    regenerating.value = true;

    twoFactorAuthStore.regenerate2FARecoveryCode({
        password: password
    }).then(response => {
        regenerating.value = false;

        currentBackupCode.value = response.recoveryCodes.join('\n');
    }).catch(error => {
        regenerating.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function reset(): void {
    new2FASecret.value = '';
    new2FAQRCode.value = '';
    currentPassword.value = '';
    currentPasscode.value = '';
    currentBackupCode.value = '';
    enabling.value = false;
    enableConfirming.value = false;
    disabling.value = false;
    regenerating.value = false;
}

function copyBackupCodes(): void {
    copyTextToClipboard(currentBackupCode.value);
    snackbar.value?.showMessage('Backup codes copied');
}

defineExpose({
    reset
});

init();
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
