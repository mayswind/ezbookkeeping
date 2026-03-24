<template>
    <f7-page no-navbar no-swipeback login-screen hide-toolbar-on-scroll>
        <f7-login-screen-title>
            <img alt="logo" class="login-page-logo" :src="APPLICATION_LOGO_PATH" />
            <f7-block class="login-page-tile margin-vertical-half">{{ tt('global.app.title') }}</f7-block>
        </f7-login-screen-title>

        <f7-block-title>{{ tt('Create Encryption Passphrase') }}</f7-block-title>
        <f7-block>
            <p>{{ tt('Your passphrase encrypts all financial data. Only you can read it.') }}</p>
        </f7-block>

        <f7-list form dividers class="margin-bottom-half">
            <f7-list-input
                type="password"
                autocomplete="new-password"
                clear-button
                :disabled="creating"
                :label="tt('Passphrase')"
                :placeholder="tt('Choose a strong passphrase')"
                v-model:value="passphrase"
                @input="updateStrength"
            />
            <f7-list-input
                type="password"
                autocomplete="new-password"
                clear-button
                :disabled="creating"
                :label="tt('Confirm Passphrase')"
                :placeholder="tt('Enter passphrase again')"
                :error-message="!passphraseMatch ? tt('Passphrases do not match') : ''"
                :error-message-force="!passphraseMatch"
                v-model:value="confirmPassphrase"
            />
        </f7-list>

        <f7-block v-if="strengthScore >= 0 && strengthScore < 3">
            <p class="text-color-red">{{ tt('Passphrase is too weak. Please choose a stronger one.') }}</p>
        </f7-block>
        <f7-block v-if="strengthFeedback">
            <p class="text-color-gray">{{ strengthFeedback }}</p>
        </f7-block>

        <f7-block>
            <f7-block-header class="text-color-orange">
                <strong>{{ tt('Lost passphrase = lost data') }}</strong><br/>
                {{ tt('There is no recovery option. If you forget your passphrase, all encrypted data is permanently lost.') }}
            </f7-block-header>
        </f7-block>

        <f7-list>
            <f7-list-item
                checkbox
                :title="tt('I understand that a lost passphrase means all my data is unrecoverable')"
                :disabled="creating"
                v-model:checked="warningAccepted"
            />
        </f7-list>

        <f7-list class="margin-top">
            <f7-list-button
                :disabled="!canSubmit"
                :loading="creating"
                @click="handleSetup"
            >{{ tt('Create Vault') }}</f7-list-button>
        </f7-list>

        <f7-block v-if="errorMessage">
            <p class="text-color-red">{{ errorMessage }}</p>
        </f7-block>
    </f7-page>
</template>

<script setup lang="ts">
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useVaultSetupPageBase } from '@/views/base/VaultSetupPageBase.ts';

import { APPLICATION_LOGO_PATH } from '@/consts/asset.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt } = useI18n();

const {
    passphrase,
    confirmPassphrase,
    warningAccepted,
    creating,
    errorMessage,
    strengthScore,
    strengthFeedback,
    passphraseMatch,
    canSubmit,
    updateStrength,
    doSetup,
} = useVaultSetupPageBase();

async function handleSetup(): Promise<void> {
    const success = await doSetup();
    if (success) {
        props.f7router.navigate('/', { clearPreviousHistory: true });
    }
}
</script>
