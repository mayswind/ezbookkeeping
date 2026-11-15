<template>
    <div class="layout-wrapper">
        <router-link to="/">
            <div class="auth-logo d-flex align-start gap-x-3">
                <img alt="logo" class="login-page-logo" :src="APPLICATION_LOGO_PATH" />
                <span class="auth-app-title">{{ tt('global.app.title') }}</span>
            </div>
        </router-link>
        <v-row no-gutters class="auth-wrapper">
            <v-col cols="12" md="8" class="auth-image-background d-none d-md-flex align-center justify-center position-relative">
                <auth-illustration variant="unlock" />
            </v-col>
            <v-col cols="12" md="4" class="auth-card d-flex flex-column">
                <div class="d-flex align-center justify-center h-100">
                    <v-card variant="flat" class="w-100 mt-0 px-4 pt-12" max-width="500">
                        <v-card-text class="py-0">
                            <div class="text-headline-small mb-2">{{ tt('Unlock Application') }}</div>
                            <div class="auth-message text-body-large mb-0" v-if="isWebAuthnAvailable">{{ tt('Please enter your PIN code or use WebAuthn to unlock application') }}</div>
                            <div class="auth-message text-body-large mb-0" v-else-if="!isWebAuthnAvailable">{{ tt('Please enter your PIN code to unlock application') }}</div>
                        </v-card-text>

                        <v-card-text class="pb-0 mb-6">
                            <v-form>
                                <v-row>
                                    <v-col cols="12">
                                        <pin-code-input :disabled="unlocking || verifyingByWebAuthn" :autofocus="true"
                                                        :secure="true" :length="6" :auto-confirm="true"
                                                        v-model="pinCode" @pincode:confirm="unlockByPin" />
                                    </v-col>

                                    <v-col cols="12">
                                        <v-btn block :disabled="unlocking || verifyingByWebAuthn || !isPinCodeValid(pinCode)"
                                               @click="unlockByPin(pinCode)">
                                            {{ tt('Unlock with PIN Code') }}
                                        </v-btn>
                                    </v-col>

                                    <v-col cols="12" v-if="isWebAuthnAvailable">
                                        <v-btn block variant="tonal" :disabled="unlocking || verifyingByWebAuthn"
                                               @click="unlockByWebAuthn">
                                            {{ tt('Unlock with WebAuthn') }}
                                            <v-progress-circular indeterminate size="22" class="ms-2" v-if="verifyingByWebAuthn"></v-progress-circular>
                                        </v-btn>
                                    </v-col>

                                    <v-col cols="12" class="text-center text-body-large mt-2">
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
                    <v-card variant="flat" class="w-100 px-4 pb-3" max-width="500">
                        <v-card-text class="pt-0">
                            <div class="text-center">
                                <language-select-button :disabled="unlocking || verifyingByWebAuthn" @error="desktopPageStore.addErrorToast" />
                            </div>

                            <v-divider class="mt-2 mb-3" />

                            <div class="auth-powered-by text-center">
                                <span>Powered by </span>
                                <a href="https://github.com/mayswind/ezbookkeeping" target="_blank">ezBookkeeping</a>&nbsp;<span>{{ version }}</span>
                            </div>
                        </v-card-text>
                    </v-card>
                </div>
            </v-col>
        </v-row>

        <confirm-dialog ref="confirmDialog"/>
    </div>
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';

import { ref, useTemplateRef } from 'vue';
import { useRouter } from 'vue-router';

import { useI18n } from '@/locales/helpers.ts';
import { useUnlockPageBase } from '@/views/base/UnlockPageBase.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useDesktopPageStore } from '@/stores/desktopPage.ts';

import { APPLICATION_LOGO_PATH } from '@/consts/asset.ts';

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

const router = useRouter();

const { tt } = useI18n();
const { version, pinCode, isWebAuthnAvailable, isPinCodeValid, doAfterUnlocked, doRelogin } = useUnlockPageBase();

const settingsStore = useSettingsStore();
const userStore = useUserStore();
const desktopPageStore = useDesktopPageStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');

const unlocking = ref<boolean>(false);
const verifyingByWebAuthn = ref<boolean>(false);

function unlockByWebAuthn(): void {
    if (unlocking.value) {
        return;
    }

    const webAuthnCredentialId = getWebAuthnCredentialId();

    if (!userStore.currentUserBasicInfo || !webAuthnCredentialId) {
        desktopPageStore.addToast('An error occurred');
        return;
    }

    if (!settingsStore.appSettings.applicationLockWebAuthn || !hasWebAuthnConfig()) {
        desktopPageStore.addToast('WebAuthn is not enabled');
        return;
    }

    if (!isWebAuthnSupported()) {
        desktopPageStore.addToast('WebAuth is not supported on this device');
        return;
    }

    verifyingByWebAuthn.value = true;

    verifyWebAuthnCredential(
        userStore.currentUserBasicInfo,
        webAuthnCredentialId
    ).then(({ id, userName, userSecret }) => {
        verifyingByWebAuthn.value = false;

        if (unlockTokenByWebAuthn(id, userName, userSecret)) {
            unlocking.value = true;

            doAfterUnlocked().then(unlockHint => {
                desktopPageStore.addToast(unlockHint);
                router.replace('/');
            }).catch(error => {
                logger.error('failed to set language after unlocking', error);

                if (!error.processed) {
                    desktopPageStore.addErrorToast(error);
                }
            }).finally(() => {
                unlocking.value = false;
            });
        } else {
            logger.error('failed to unlock with webauthn');
            desktopPageStore.addToast('Failed to authenticate with WebAuthn');
        }
    }).catch(error => {
        verifyingByWebAuthn.value = false;
        logger.error('failed to use webauthn to verify', error);

        if (error.notSupported) {
            desktopPageStore.addToast('WebAuth is not supported on this device');
        } else if (error.name === 'NotAllowedError') {
            desktopPageStore.addToast('User has canceled authentication');
        } else if (error.invalid) {
            desktopPageStore.addToast('Failed to authenticate with WebAuthn');
        } else {
            desktopPageStore.addToast('User has canceled or this device does not support WebAuthn');
        }
    });
}

function unlockByPin(pinCode: string): void {
    if (unlocking.value || !isPinCodeValid(pinCode)) {
        return;
    }

    const user = userStore.currentUserBasicInfo;

    if (!user || !user.username) {
        desktopPageStore.addToast('An error occurred');
        return;
    }

    try {
        if (unlockTokenByPinCode(user.username, pinCode)) {
            unlocking.value = true;

            doAfterUnlocked().then(unlockHint => {
                desktopPageStore.addToast(unlockHint);
                router.replace('/');
            }).catch(error => {
                logger.error('failed to set language after unlocking', error);

                if (!error.processed) {
                    desktopPageStore.addErrorToast(error);
                }
            }).finally(() => {
                unlocking.value = false;
            });
        } else {
            logger.error('failed to unlock with pin code');
            desktopPageStore.addToast('Incorrect PIN code');
        }
    } catch (ex) {
        logger.error('failed to unlock with pin code', ex);
        desktopPageStore.addToast('Incorrect PIN code');
    }
}

function relogin(): void {
    confirmDialog.value?.open('Are you sure you want to re-login?').then(() => {
        doRelogin().then(() => {
            router.replace('/login');
        }).catch(error => {
            if (!error.processed) {
                desktopPageStore.addErrorToast(error);
            }
        });
    });
}
</script>
