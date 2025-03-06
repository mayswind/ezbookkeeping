<template>
    <div class="layout-wrapper">
        <router-link to="/">
            <div class="auth-logo d-flex align-start gap-x-3">
                <img alt="logo" class="login-page-logo" :src="APPLICATION_LOGO_PATH" />
                <h1 class="font-weight-medium leading-normal text-2xl">{{ tt('global.app.title') }}</h1>
            </div>
        </router-link>
        <v-row no-gutters class="auth-wrapper">
            <v-col cols="12" md="8" class="d-none d-md-flex align-center justify-center position-relative">
                <div class="d-flex auth-img-footer" v-if="!isDarkMode">
                    <v-img src="img/desktop/background.svg"/>
                </div>
                <div class="d-flex auth-img-footer" v-if="isDarkMode">
                    <v-img src="img/desktop/background-dark.svg"/>
                </div>
                <div class="d-flex align-center justify-center w-100 pt-10">
                    <v-img max-width="600px" src="img/desktop/people3.svg"/>
                </div>
            </v-col>
            <v-col cols="12" md="4" class="auth-card d-flex flex-column">
                <div class="d-flex align-center justify-center h-100">
                    <v-card variant="flat" class="w-100 mt-0 px-4 pt-12" max-width="500">
                        <v-card-text>
                            <h4 class="text-h4 mb-2">{{ tt('Unlock Application') }}</h4>
                            <p class="mb-0" v-if="isWebAuthnAvailable">{{ tt('Please enter your PIN code or use WebAuthn to unlock application') }}</p>
                            <p class="mb-0" v-else-if="!isWebAuthnAvailable">{{ tt('Please enter your PIN code to unlock application') }}</p>
                        </v-card-text>

                        <v-card-text class="pb-0 mb-6">
                            <v-form>
                                <v-row>
                                    <v-col cols="12">
                                        <pin-code-input :disabled="verifyingByWebAuthn" :autofocus="true"
                                                        :secure="true" :length="6" :auto-confirm="true"
                                                        v-model="pinCode" @pincode:confirm="unlockByPin" />
                                    </v-col>

                                    <v-col cols="12">
                                        <v-btn block :disabled="!isPinCodeValid(pinCode) || verifyingByWebAuthn"
                                               @click="unlockByPin(pinCode)">
                                            {{ tt('Unlock with PIN Code') }}
                                        </v-btn>
                                    </v-col>

                                    <v-col cols="12" v-if="isWebAuthnAvailable">
                                        <v-btn block variant="tonal" :disabled="verifyingByWebAuthn"
                                               @click="unlockByWebAuthn">
                                            {{ tt('Unlock with WebAuthn') }}
                                            <v-progress-circular indeterminate size="22" class="ml-2" v-if="verifyingByWebAuthn"></v-progress-circular>
                                        </v-btn>
                                    </v-col>

                                    <v-col cols="12" class="text-center">
                                        <span class="me-1">{{ tt('Can\'t Unlock?') }}</span>
                                        <a class="text-primary" href="javascript:void(0);" @click="relogin">
                                            {{ tt('Re-login') }}
                                        </a>
                                    </v-col>
                                </v-row>
                            </v-form>
                        </v-card-text>
                    </v-card>
                </div>
                <v-spacer/>
                <div class="d-flex align-center justify-center">
                    <v-card variant="flat" class="w-100 px-4 pb-4" max-width="500">
                        <v-card-text class="pt-0">
                            <v-row>
                                <v-col cols="12" class="text-center">
                                    <language-select-button :disabled="verifyingByWebAuthn" />
                                </v-col>

                                <v-col cols="12" class="d-flex align-center pt-0">
                                    <v-divider />
                                </v-col>

                                <v-col cols="12" class="text-center text-sm">
                                    <span>Powered by </span>
                                    <a href="https://github.com/mayswind/ezbookkeeping" target="_blank">ezBookkeeping</a>&nbsp;<span>{{ version }}</span>
                                </v-col>
                            </v-row>
                        </v-card-text>
                    </v-card>
                </div>
            </v-col>
        </v-row>

        <confirm-dialog ref="confirmDialog"/>
        <snack-bar ref="snackbar" />
    </div>
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';
import { useRouter } from 'vue-router';
import { useTheme } from 'vuetify';

import { useI18n } from '@/locales/helpers.ts';
import { useUnlockPageBase } from '@/views/base/UnlockPageBase.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';

import { APPLICATION_LOGO_PATH } from '@/consts/asset.ts';
import { ThemeType } from '@/core/theme.ts';
import {
    isWebAuthnSupported,
    verifyWebAuthnCredential
} from '@/lib/webauthn.ts';
import {
    unlockTokenByWebAuthn,
    unlockTokenByPinCode,
    hasWebAuthnConfig,
    getWebAuthnCredentialId
} from '@/lib/userstate.ts';
import logger from '@/lib/logger.ts';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;

const router = useRouter();
const theme = useTheme();

const { tt } = useI18n();
const { version, pinCode, isWebAuthnAvailable, isPinCodeValid, doAfterUnlocked, doRelogin } = useUnlockPageBase();

const settingsStore = useSettingsStore();
const userStore = useUserStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const verifyingByWebAuthn = ref<boolean>(false);

const isDarkMode = computed<boolean>(() => theme.global.name.value === ThemeType.Dark);

function unlockByWebAuthn(): void {
    const webAuthnCredentialId = getWebAuthnCredentialId();

    if (!userStore.currentUserBasicInfo || !webAuthnCredentialId) {
        snackbar.value?.showMessage('An error occurred');
        return;
    }

    if (!settingsStore.appSettings.applicationLockWebAuthn || !hasWebAuthnConfig()) {
        snackbar.value?.showMessage('WebAuthn is not enabled');
        return;
    }

    if (!isWebAuthnSupported()) {
        snackbar.value?.showMessage('WebAuth is not supported on this device');
        return;
    }

    verifyingByWebAuthn.value = true;

    verifyWebAuthnCredential(
        userStore.currentUserBasicInfo,
        webAuthnCredentialId
    ).then(({ id, userName, userSecret }) => {
        verifyingByWebAuthn.value = false;

        unlockTokenByWebAuthn(id, userName, userSecret);
        doAfterUnlocked();

        router.replace('/');
    }).catch(error => {
        verifyingByWebAuthn.value = false;
        logger.error('failed to use webauthn to verify', error);

        if (error.notSupported) {
            snackbar.value?.showMessage('WebAuth is not supported on this device');
        } else if (error.name === 'NotAllowedError') {
            snackbar.value?.showMessage('User has canceled authentication');
        } else if (error.invalid) {
            snackbar.value?.showMessage('Failed to authenticate with WebAuthn');
        } else {
            snackbar.value?.showMessage('User has canceled or this device does not support WebAuthn');
        }
    });
}

function unlockByPin(pinCode: string): void {
    if (!isPinCodeValid(pinCode)) {
        return;
    }

    const user = userStore.currentUserBasicInfo;

    if (!user || !user.username) {
        snackbar.value?.showMessage('An error occurred');
        return;
    }

    try {
        unlockTokenByPinCode(user.username, pinCode);
        doAfterUnlocked();

        router.replace('/');
    } catch (ex) {
        logger.error('failed to unlock with pin code', ex);
        snackbar.value?.showMessage('Incorrect PIN code');
    }
}

function relogin(): void {
    confirmDialog.value?.open('Are you sure you want to re-login?').then(() => {
        doRelogin();

        router.replace('/login');
    });
}
</script>
