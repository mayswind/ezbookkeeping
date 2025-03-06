<template>
    <f7-page no-navbar no-swipeback login-screen hide-toolbar-on-scroll>
        <f7-login-screen-title>
            <img alt="logo" class="login-page-logo" :src="APPLICATION_LOGO_PATH" />
            <f7-block class="login-page-tile margin-vertical-half">{{ tt('global.app.title') }}</f7-block>
        </f7-login-screen-title>

        <f7-list form>
            <f7-list-item class="no-padding no-margin">
                <template #inner>
                    <div class="display-flex justify-content-center full-line">{{ tt('Unlock Application') }}</div>
                </template>
            </f7-list-item>
            <f7-list-item class="list-item-pincode-input padding-horizontal margin-horizontal">
                <pin-code-input :secure="true" :length="6" :auto-confirm="true" v-model="pinCode" @pincode:confirm="unlockByPin" />
            </f7-list-item>
        </f7-list>

        <f7-list>
            <f7-list-button :class="{ 'disabled': !isPinCodeValid(pinCode) }" :text="tt('Unlock with PIN Code')" @click="unlockByPin"></f7-list-button>
            <f7-list-button v-if="isWebAuthnAvailable" :text="tt('Unlock with WebAuthn')" @click="unlockByWebAuthn"></f7-list-button>
            <f7-block-footer>
                <f7-link :text="tt('Re-login')" @click="relogin"></f7-link>
            </f7-block-footer>
            <f7-block-footer class="padding-bottom">
            </f7-block-footer>
        </f7-list>

        <language-select-button />

        <f7-list class="login-page-bottom">
            <f7-block-footer>
                <div class="login-page-powered-by">
                    <span>Powered by</span>
                    <f7-link external href="https://github.com/mayswind/ezbookkeeping" target="_blank">ezBookkeeping</f7-link>
                    <span>{{ version }}</span>
                </div>
            </f7-block-footer>
        </f7-list>

        <f7-toolbar class="login-page-fixed-bottom" tabbar bottom :outline="false">
            <div class="login-page-powered-by">
                <span>Powered by</span>
                <f7-link external href="https://github.com/mayswind/ezbookkeeping" target="_blank">ezBookkeeping</f7-link>
                <span>{{ version }}</span>
            </div>
        </f7-toolbar>
    </f7-page>
</template>

<script setup lang="ts">
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';
import { useUnlockPageBase } from '@/views/base/UnlockPageBase.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';

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
import { isModalShowing } from '@/lib/ui/mobile.ts';
import logger from '@/lib/logger.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt } = useI18n();
const { showToast, showConfirm } = useI18nUIComponents();
const { version, pinCode, isWebAuthnAvailable, isPinCodeValid, doAfterUnlocked, doRelogin } = useUnlockPageBase();

const settingsStore = useSettingsStore();
const userStore = useUserStore();

function unlockByWebAuthn(): void {
    const router = props.f7router;
    const webAuthnCredentialId = getWebAuthnCredentialId();

    if (!userStore.currentUserBasicInfo || !webAuthnCredentialId) {
        showToast('An error occurred');
        return;
    }

    if (!settingsStore.appSettings.applicationLockWebAuthn || !hasWebAuthnConfig()) {
        showToast('WebAuthn is not enabled');
        return;
    }

    if (!isWebAuthnSupported()) {
        showToast('WebAuth is not supported on this device');
        return;
    }

    showLoading();

    verifyWebAuthnCredential(
        userStore.currentUserBasicInfo,
        webAuthnCredentialId
    ).then(({ id, userName, userSecret }) => {
        hideLoading();

        unlockTokenByWebAuthn(id, userName, userSecret);
        doAfterUnlocked();

        router.refreshPage();
    }).catch(error => {
        hideLoading();
        logger.error('failed to use webauthn to verify', error);

        if (error.notSupported) {
            showToast('WebAuth is not supported on this device');
        } else if (error.name === 'NotAllowedError') {
            showToast('User has canceled authentication');
        } else if (error.invalid) {
            showToast('Failed to authenticate with WebAuthn');
        } else {
            showToast('User has canceled or this device does not support WebAuthn');
        }
    });
}

function unlockByPin(pinCode: string): void {
    if (!isPinCodeValid(pinCode)) {
        return;
    }

    if (isModalShowing()) {
        return;
    }

    const router = props.f7router;
    const user = userStore.currentUserBasicInfo;

    if (!user || !user.username) {
        showToast('An error occurred');
        return;
    }

    try {
        unlockTokenByPinCode(user.username, pinCode);
        doAfterUnlocked();

        router.refreshPage();
    } catch (ex) {
        logger.error('failed to unlock with pin code', ex);
        showToast('Incorrect PIN code');
    }
}

function relogin(): void {
    const router = props.f7router;

    showConfirm('Are you sure you want to re-login?', () => {
        doRelogin();

        router.navigate('/login', {
            clearPreviousHistory: true
        });
    });
}
</script>
