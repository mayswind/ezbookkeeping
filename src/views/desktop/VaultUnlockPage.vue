<template>
    <div class="layout-wrapper">
        <div class="auth-logo d-flex align-start gap-x-3">
            <img alt="logo" class="login-page-logo" :src="APPLICATION_LOGO_PATH" />
            <h1 class="font-weight-medium leading-normal text-2xl">{{ tt('global.app.title') }}</h1>
        </div>
        <v-row no-gutters class="auth-wrapper">
            <v-col cols="12" md="4" class="auth-card d-flex flex-column mx-auto">
                <div class="d-flex align-center justify-center h-100">
                    <v-card variant="flat" class="w-100 mt-0 px-4 pt-12" max-width="500">
                        <v-card-text>
                            <h4 class="text-h4 mb-2">{{ tt('Unlock Vault') }}</h4>
                            <p class="mb-0">{{ tt('Enter your passphrase to decrypt your data.') }}</p>
                        </v-card-text>

                        <v-card-text class="pb-0">
                            <v-form @submit.prevent="handleUnlock">
                                <v-row>
                                    <v-col cols="12">
                                        <v-text-field
                                            type="password"
                                            autocomplete="current-password"
                                            :autofocus="true"
                                            :disabled="unlocking"
                                            :label="tt('Passphrase')"
                                            :placeholder="tt('Enter your passphrase')"
                                            v-model="passphrase"
                                            @keyup.enter="handleUnlock"
                                        />
                                    </v-col>

                                    <v-col cols="12">
                                        <v-btn
                                            block
                                            color="primary"
                                            type="submit"
                                            :disabled="!passphrase || unlocking"
                                            :loading="unlocking"
                                        >
                                            {{ tt('Unlock') }}
                                        </v-btn>
                                    </v-col>

                                    <v-col cols="12" v-if="errorMessage">
                                        <v-alert type="error" variant="tonal">{{ errorMessage }}</v-alert>
                                    </v-col>

                                    <v-col cols="12" class="text-center">
                                        <p class="text-caption text-medium-emphasis">
                                            {{ tt('Forgot your passphrase? All encrypted data is permanently unrecoverable by design.') }}
                                        </p>
                                    </v-col>
                                </v-row>
                            </v-form>
                        </v-card-text>
                    </v-card>
                </div>
            </v-col>
        </v-row>
    </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';

import { useI18n } from '@/locales/helpers.ts';
import { useVaultUnlockPageBase } from '@/views/base/VaultUnlockPageBase.ts';

import { APPLICATION_LOGO_PATH } from '@/consts/asset.ts';

const { tt } = useI18n();
const router = useRouter();

const {
    passphrase,
    unlocking,
    errorMessage,
    doUnlock,
} = useVaultUnlockPageBase();

async function handleUnlock(): Promise<void> {
    const success = await doUnlock();
    if (success) {
        router.replace('/');
    }
}
</script>
