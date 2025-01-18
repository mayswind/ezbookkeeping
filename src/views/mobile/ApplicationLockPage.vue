<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Application Lock')"></f7-nav-title>
        </f7-navbar>

        <f7-list strong inset dividers class="margin-top">
            <f7-list-item :title="tt('Status')" :after="tt(isEnableApplicationLock ? 'Enabled' : 'Disabled')"></f7-list-item>
            <f7-list-item v-if="isEnableApplicationLock">
                <span>{{ tt('Unlock with PIN Code') }}</span>
                <f7-toggle checked disabled></f7-toggle>
            </f7-list-item>
            <f7-list-item v-if="isEnableApplicationLock && isSupportedWebAuthn">
                <span>{{ tt('Unlock with WebAuthn') }}</span>
                <f7-toggle :checked="isEnableApplicationLockWebAuthn" @toggle:change="isEnableApplicationLockWebAuthn = $event"></f7-toggle>
            </f7-list-item>
            <f7-list-button v-if="isEnableApplicationLock" @click="disable(null)">{{ tt('Disable') }}</f7-list-button>
            <f7-list-button v-if="!isEnableApplicationLock" @click="enable(null)">{{ tt('Enable') }}</f7-list-button>
        </f7-list>

        <pin-code-input-sheet :title="tt('PIN Code')"
                              :hint="tt('Please enter a new 6-digit PIN code. The PIN code would encrypt your local data, so you need to enter it every time you open this app. If this PIN code is lost, you will need to log in again.')"
                              v-model:show="showInputPinCodeSheetForEnable"
                              v-model="currentPinCodeForEnable"
                              @pincode:confirm="enable">
        </pin-code-input-sheet>

        <pin-code-input-sheet :title="tt('PIN Code')"
                              :hint="tt('Your current PIN code is required to disable application lock.')"
                              v-model:show="showInputPinCodeSheetForDisable"
                              v-model="currentPinCodeForDisable"
                              @pincode:confirm="disable">
        </pin-code-input-sheet>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';
import { useAppLockPageBase } from '@/views/base/settings/AppLockPageBase.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useTransactionsStore } from '@/stores/transaction.js';

import { registerWebAuthnCredential } from '@/lib/webauthn.ts';
import {
    getUserAppLockState,
    encryptToken,
    decryptToken,
    isCorrectPinCode,
    saveWebAuthnConfig,
    clearWebAuthnConfig
} from '@/lib/userstate.ts';
import logger from '@/lib/logger.ts';

const { tt } = useI18n();
const { showToast } = useI18nUIComponents();
const { isSupportedWebAuthn, isEnableApplicationLock, isEnableApplicationLockWebAuthn } = useAppLockPageBase();

const settingsStore = useSettingsStore();
const userStore = useUserStore();
const transactionsStore = useTransactionsStore();

const currentPinCodeForEnable = ref<string>('');
const currentPinCodeForDisable = ref<string>('');
const showInputPinCodeSheetForEnable = ref<boolean>(false);
const showInputPinCodeSheetForDisable = ref<boolean>(false);

function enable(pinCode: string | null): void {
    if (settingsStore.appSettings.applicationLock) {
        showToast('Application lock has been enabled');
        return;
    }

    if (!pinCode) {
        currentPinCodeForEnable.value = '';
        showInputPinCodeSheetForEnable.value = true;
        return;
    }

    if (!currentPinCodeForEnable.value || currentPinCodeForEnable.value.length !== 6) {
        nextTick(() => {
            currentPinCodeForEnable.value = '';
        });
        showToast('Invalid PIN code');
        return;
    }

    const user = userStore.currentUserBasicInfo;

    if (!user || !user.username) {
        nextTick(() => {
            currentPinCodeForEnable.value = '';
        });
        showToast('An error occurred');
        return;
    }

    encryptToken(user.username, pinCode);
    settingsStore.setEnableApplicationLock(true);
    transactionsStore.saveTransactionDraft();

    settingsStore.setEnableApplicationLockWebAuthn(false);
    clearWebAuthnConfig();

    showInputPinCodeSheetForEnable.value = false;
}

function disable(pinCode: string | null): void {
    if (!settingsStore.appSettings.applicationLock) {
        showToast('Application lock is not enabled');
        return;
    }

    if (!pinCode) {
        currentPinCodeForDisable.value = '';
        showInputPinCodeSheetForDisable.value = true;
        return;
    }

    if (!isCorrectPinCode(pinCode)) {
        nextTick(() => {
            currentPinCodeForDisable.value = '';
        });
        showToast('Incorrect PIN code');
        return;
    }

    decryptToken();
    settingsStore.setEnableApplicationLock(false);
    transactionsStore.saveTransactionDraft();

    settingsStore.setEnableApplicationLockWebAuthn(false);
    clearWebAuthnConfig();

    showInputPinCodeSheetForDisable.value = false;
}

watch(isEnableApplicationLockWebAuthn, (newValue) => {
    const userAppLockState = getUserAppLockState();

    if (newValue && userAppLockState && userStore.currentUserBasicInfo) {
        showLoading();

        registerWebAuthnCredential(
            userAppLockState,
            userStore.currentUserBasicInfo,
        ).then(({ id }) => {
            hideLoading();

            saveWebAuthnConfig(id);
            settingsStore.setEnableApplicationLockWebAuthn(true);
            showToast('You have enabled WebAuthn successfully');
        }).catch(error => {
            logger.error('failed to enable WebAuthn', error);

            hideLoading();

            if (error.notSupported) {
                showToast('WebAuth is not supported on this device');
            } else if (error.name === 'NotAllowedError') {
                showToast('User has canceled authentication');
            } else if (error.invalid) {
                showToast('Failed to enable WebAuthn');
            } else {
                showToast('User has canceled or this device does not support WebAuthn');
            }

            isEnableApplicationLockWebAuthn.value = false;
            settingsStore.setEnableApplicationLockWebAuthn(false);
            clearWebAuthnConfig();
        });
    } else {
        settingsStore.setEnableApplicationLockWebAuthn(false);
        clearWebAuthnConfig();
    }
});
</script>
