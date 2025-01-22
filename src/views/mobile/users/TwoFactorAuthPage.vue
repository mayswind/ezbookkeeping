<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar :title="tt('Two-Factor Authentication')" :back-link="tt('Back')"></f7-navbar>

        <f7-list strong inset dividers class="margin-top skeleton-text" v-if="loading">
            <f7-list-item title="Status" after="Unknown"></f7-list-item>
            <f7-list-button class="disabled">Operate</f7-list-button>
        </f7-list>

        <f7-list strong inset dividers class="margin-top" v-else-if="!loading">
            <f7-list-item :title="tt('Status')" :after="tt(status ? 'Enabled' : 'Disabled')"></f7-list-item>
            <f7-list-button :class="{ 'disabled': regenerating }" v-if="status === true" @click="regenerateBackupCode(null)">{{ tt('Regenerate Backup Codes') }}</f7-list-button>
            <f7-list-button :class="{ 'disabled': disabling }" v-if="status === true" @click="disable(null)">{{ tt('Disable') }}</f7-list-button>
            <f7-list-button :class="{ 'disabled': enabling }" v-if="status === false" @click="enable">{{ tt('Enable') }}</f7-list-button>
        </f7-list>

        <passcode-input-sheet :title="tt('Enable Two-Factor Authentication')"
                              :hint="tt('Please use a two-factor authentication app to scan the qrcode below and enter the current passcode.')"
                              :confirm-disabled="enableConfirming"
                              :cancel-disabled="enableConfirming"
                              v-model:show="showInputPasscodeSheetForEnable"
                              v-model="currentPasscodeForEnable"
                              @passcode:confirm="enableConfirm">
            <div class="text-align-center">
                <img alt="qrcode" class="img-qrcode" :src="new2FAQRCode" />
            </div>
        </passcode-input-sheet>

        <password-input-sheet :title="tt('Disable Two-Factor Authentication')"
                              :hint="tt('Your current password is required to disable two-factor authentication.')"
                              :confirm-disabled="disabling"
                              :cancel-disabled="disabling"
                              v-model:show="showInputPasswordSheetForDisable"
                              v-model="currentPasswordForDisable"
                              @password:confirm="disable">
        </password-input-sheet>

        <password-input-sheet :title="tt('Regenerate Backup Codes')"
                              :hint="tt('Your current password is required to regenerate backup codes for two-factor authentication. If you regenerate backup codes, the previous ones will become invalid.')"
                              :confirm-disabled="regenerating"
                              :cancel-disabled="regenerating"
                              v-model:show="showInputPasswordSheetForRegenerate"
                              v-model="currentPasswordForRegenerate"
                              @password:confirm="regenerateBackupCode">
        </password-input-sheet>

        <information-sheet class="backup-code-sheet"
                           :title="tt('Backup Code')"
                           :hint="tt('Please copy these backup codes to safe place, the following backup codes will be displayed only once. If these codes were lost, you can regenerate them at any time.')"
                           :information="currentBackupCode"
                           :row-count="10"
                           :enable-copy="true"
                           v-model:show="showBackupCodeSheet"
                           @info:copied="onBackupCodeCopied">
        </information-sheet>
    </f7-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';

import { useTwoFactorAuthStore } from '@/stores/twoFactorAuth.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt } = useI18n();
const { showToast, routeBackOnError } = useI18nUIComponents();

const twoFactorAuthStore = useTwoFactorAuthStore();

const status = ref<boolean | null>(null);
const loading = ref<boolean>(true);
const loadingError = ref<unknown | null>(null);
const new2FASecret = ref<string>('');
const new2FAQRCode = ref<string>('');
const currentPasscodeForEnable = ref<string>('');
const currentPasswordForDisable = ref<string>('');
const currentPasswordForRegenerate = ref<string>('');
const currentBackupCode = ref<string>('');
const enabling = ref<boolean>(false);
const enableConfirming = ref<boolean>(false);
const disabling = ref<boolean>(false);
const regenerating = ref<boolean>(false);
const showInputPasscodeSheetForEnable = ref<boolean>(false);
const showInputPasswordSheetForDisable = ref<boolean>(false);
const showInputPasswordSheetForRegenerate = ref<boolean>(false);
const showBackupCodeSheet = ref<boolean>(false);

function init(): void {
    loading.value = true;

    twoFactorAuthStore.get2FAStatus().then(response => {
        status.value = response.enable;
        loading.value = false;
    }).catch(error => {
        if (error.processed) {
            loading.value = false;
        } else {
            loadingError.value = error;
            showToast(error.message || error);
        }
    });
}

function enable(): void {
    new2FAQRCode.value = '';
    new2FASecret.value = '';

    enabling.value = true;
    showLoading(() => enabling.value);

    twoFactorAuthStore.enable2FA().then(response => {
        enabling.value = false;
        hideLoading();

        new2FAQRCode.value = response.qrcode;
        new2FASecret.value = response.secret;

        showInputPasscodeSheetForEnable.value = true;
    }).catch(error => {
        enabling.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function enableConfirm(): void {
    enableConfirming.value = true;
    showLoading(() => enableConfirming.value);

    twoFactorAuthStore.confirmEnable2FA({
        secret: new2FASecret.value,
        passcode: currentPasscodeForEnable.value
    }).then(response => {
        enableConfirming.value = false;
        hideLoading();

        new2FAQRCode.value = '';
        new2FASecret.value = '';

        status.value = true;
        showInputPasscodeSheetForEnable.value = false;

        if (response.recoveryCodes && response.recoveryCodes.length) {
            currentBackupCode.value = response.recoveryCodes.join('\n');
            showBackupCodeSheet.value = true;
        }
    }).catch(error => {
        enableConfirming.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function disable(password: string | null): void {
    if (!password) {
        currentPasswordForDisable.value = '';
        showInputPasswordSheetForDisable.value = true;
        return;
    }

    disabling.value = true;
    showLoading(() => disabling.value);

    twoFactorAuthStore.disable2FA({
        password: password
    }).then(() => {
        disabling.value = false;
        hideLoading();

        status.value = false;
        showInputPasswordSheetForDisable.value = false;
        showToast('Two-factor authentication has been disabled');
    }).catch(error => {
        disabling.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function regenerateBackupCode(password: string | null): void {
    if (!password) {
        currentPasswordForRegenerate.value = '';
        showInputPasswordSheetForRegenerate.value = true;
        return;
    }

    regenerating.value = true;
    showLoading(() => regenerating.value);

    twoFactorAuthStore.regenerate2FARecoveryCode({
        password: password
    }).then(response => {
        regenerating.value = false;
        hideLoading();

        showInputPasswordSheetForRegenerate.value = false;

        currentBackupCode.value = response.recoveryCodes.join('\n');
        showBackupCodeSheet.value = true;
    }).catch(error => {
        regenerating.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function onPageAfterIn(): void {
    routeBackOnError(props.f7router, loadingError);
}

function onBackupCodeCopied(): void {
    showToast('Backup codes copied');
}

init();
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
