<template>
    <v-row>
        <v-col cols="12">
            <v-card :title="tt('Application Lock')">
                <v-card-text class="pb-0">
                    <p class="text-body-1 font-weight-semibold" v-if="!isEnableApplicationLock">
                        {{ tt('Application lock is not enabled') }}
                    </p>
                    <p class="text-body-1" v-if="isEnableApplicationLock">
                        {{ tt('Application lock has been enabled') }}
                    </p>
                </v-card-text>

                <v-card-text v-if="isEnableApplicationLock">
                    <v-switch :disabled="true"
                              :label="tt('Unlock with PIN Code')"
                              v-model="isEnableApplicationLock"/>
                    <v-switch :label="tt('Unlock with WebAuthn')"
                              :loading="enablingWebAuthn"
                              v-model="isEnableApplicationLockWebAuthn"
                              v-if="isSupportedWebAuthn"/>
                </v-card-text>

                <v-card-text class="pb-0">
                    <p class="text-body-1 font-weight-semibold" v-if="!isEnableApplicationLock">
                        {{ tt('Please enter a new 6-digit PIN code. The PIN code would encrypt your local data, so you need to enter it every time you open this app. If this PIN code is lost, you will need to log in again.') }}
                    </p>
                    <p class="text-body-1 font-weight-semibold" v-if="isEnableApplicationLock">
                        {{ tt('Your current PIN code is required to disable application lock.') }}
                    </p>
                </v-card-text>

                <v-card-text class="pb-0">
                    <v-row class="mb-3">
                        <v-col cols="12" md="12">
                            <div style="max-width: 428px">
                                <pin-code-input :secure="true" :length="6" v-model="pinCode" @pincode:confirm="confirm" />
                            </div>
                        </v-col>
                    </v-row>
                </v-card-text>

                <v-card-text>
                    <v-row>
                        <v-col cols="12" class="d-flex flex-wrap gap-4">
                            <v-btn :disabled="!pinCodeValid"
                                   v-if="!isEnableApplicationLock" @click="enable">
                                {{ tt('Enable Application Lock') }}
                            </v-btn>
                            <v-btn :disabled="!pinCodeValid"
                                   v-if="isEnableApplicationLock" @click="disable">
                                {{ tt('Disable Application Lock') }}
                            </v-btn>
                        </v-col>
                    </v-row>
                </v-card-text>
            </v-card>
        </v-col>
    </v-row>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
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

type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();
const { isSupportedWebAuthn, isEnableApplicationLock, isEnableApplicationLockWebAuthn } = useAppLockPageBase();

const settingsStore = useSettingsStore();
const userStore = useUserStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const pinCode = ref<string>('');
const enablingWebAuthn = ref<boolean>(false);
const transactionsStore = useTransactionsStore();

const pinCodeValid = computed<boolean>(() => {
    return pinCode.value?.length === 6 || false;
});

function confirm(): void {
    if (isEnableApplicationLock.value) {
        disable();
    } else {
        enable();
    }
}

function enable(): void {
    if (settingsStore.appSettings.applicationLock) {
        snackbar.value?.showMessage('Application lock has been enabled');
        return;
    }

    if (!pinCode.value || pinCode.value.length !== 6) {
        pinCode.value = '';
        snackbar.value?.showMessage('Invalid PIN code');
        return;
    }

    const user = userStore.currentUserBasicInfo;

    if (!user || !user.username) {
        pinCode.value = '';
        snackbar.value?.showMessage('An error occurred');
        return;
    }

    encryptToken(user.username, pinCode.value);
    settingsStore.setEnableApplicationLock(true);
    transactionsStore.saveTransactionDraft();

    settingsStore.setEnableApplicationLockWebAuthn(false);
    clearWebAuthnConfig();

    pinCode.value = '';
}

function disable(): void {
    if (!settingsStore.appSettings.applicationLock) {
        snackbar.value?.showMessage('Application lock is not enabled');
        return;
    }

    if (!isCorrectPinCode(pinCode.value)) {
        pinCode.value = '';
        snackbar.value?.showMessage('Incorrect PIN code');
        return;
    }

    pinCode.value = '';

    decryptToken();
    settingsStore.setEnableApplicationLock(false);
    transactionsStore.saveTransactionDraft();

    settingsStore.setEnableApplicationLockWebAuthn(false);
    clearWebAuthnConfig();
}

watch(isEnableApplicationLockWebAuthn, (newValue) => {
    const userAppLockState = getUserAppLockState();

    if (newValue && userAppLockState && userStore.currentUserBasicInfo) {
        enablingWebAuthn.value = true;

        registerWebAuthnCredential(
            userAppLockState,
            userStore.currentUserBasicInfo,
        ).then(({ id }) => {
            enablingWebAuthn.value = false;

            saveWebAuthnConfig(id);
            settingsStore.setEnableApplicationLockWebAuthn(true);
            snackbar.value?.showMessage('You have enabled WebAuthn successfully');
        }).catch(error => {
            logger.error('failed to enable WebAuthn', error);

            enablingWebAuthn.value = false;

            if (error.notSupported) {
                snackbar.value?.showMessage('WebAuth is not supported on this device');
            } else if (error.name === 'NotAllowedError') {
                snackbar.value?.showMessage('User has canceled authentication');
            } else if (error.invalid) {
                snackbar.value?.showMessage('Failed to enable WebAuthn');
            } else {
                snackbar.value?.showMessage('User has canceled or this device does not support WebAuthn');
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
