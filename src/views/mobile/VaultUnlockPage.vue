<template>
    <f7-page no-navbar no-swipeback login-screen hide-toolbar-on-scroll>
        <f7-login-screen-title>
            <img alt="logo" class="login-page-logo" :src="APPLICATION_LOGO_PATH" />
            <f7-block class="login-page-tile margin-vertical-half">{{ tt('global.app.title') }}</f7-block>
        </f7-login-screen-title>

        <f7-block-title>{{ tt('Unlock Vault') }}</f7-block-title>
        <f7-block>
            <p>{{ tt('Enter your passphrase to decrypt your data.') }}</p>
        </f7-block>

        <f7-list form dividers class="margin-bottom-half">
            <f7-list-input
                type="password"
                autocomplete="current-password"
                clear-button
                :disabled="unlocking"
                :label="tt('Passphrase')"
                :placeholder="tt('Enter your passphrase')"
                v-model:value="passphrase"
                @keyup.enter="handleUnlock"
            />
        </f7-list>

        <f7-list class="margin-top">
            <f7-list-button
                :disabled="!passphrase || unlocking"
                :loading="unlocking"
                @click="handleUnlock"
            >{{ tt('Unlock') }}</f7-list-button>
        </f7-list>

        <f7-block v-if="errorMessage">
            <p class="text-color-red">{{ errorMessage }}</p>
        </f7-block>

        <f7-block>
            <p class="text-color-gray text-align-center" style="font-size: 12px;">
                {{ tt('Forgot your passphrase? All encrypted data is permanently unrecoverable by design.') }}
            </p>
        </f7-block>
    </f7-page>
</template>

<script setup lang="ts">
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useVaultUnlockPageBase } from '@/views/base/VaultUnlockPageBase.ts';

import { APPLICATION_LOGO_PATH } from '@/consts/asset.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt } = useI18n();

const {
    passphrase,
    unlocking,
    errorMessage,
    doUnlock,
} = useVaultUnlockPageBase();

async function handleUnlock(): Promise<void> {
    const success = await doUnlock();
    if (success) {
        props.f7router.navigate('/', { clearPreviousHistory: true });
    }
}
</script>
