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
                            <h4 class="text-h4 mb-2">{{ tt('Create Encryption Passphrase') }}</h4>
                            <p class="mb-0">{{ tt('Your passphrase encrypts all financial data. Only you can read it.') }}</p>
                        </v-card-text>

                        <v-card-text class="pb-0">
                            <v-form @submit.prevent="handleSetup">
                                <v-row>
                                    <v-col cols="12">
                                        <v-text-field
                                            type="password"
                                            autocomplete="new-password"
                                            :autofocus="true"
                                            :disabled="creating"
                                            :label="tt('Passphrase')"
                                            :placeholder="tt('Choose a strong passphrase')"
                                            v-model="passphrase"
                                            @input="updateStrength"
                                        />
                                        <v-progress-linear
                                            v-if="strengthScore >= 0"
                                            :model-value="(strengthScore + 1) * 20"
                                            :color="strengthColor"
                                            class="mb-1"
                                        />
                                        <p v-if="strengthFeedback" class="text-caption text-medium-emphasis">{{ strengthFeedback }}</p>
                                        <p v-if="strengthScore >= 0 && strengthScore < 3" class="text-caption text-error">
                                            {{ tt('Passphrase is too weak. Please choose a stronger one.') }}
                                        </p>
                                    </v-col>

                                    <v-col cols="12">
                                        <v-text-field
                                            type="password"
                                            autocomplete="new-password"
                                            :disabled="creating"
                                            :label="tt('Confirm Passphrase')"
                                            :placeholder="tt('Enter passphrase again')"
                                            :error="!passphraseMatch"
                                            :error-messages="!passphraseMatch ? [tt('Passphrases do not match')] : []"
                                            v-model="confirmPassphrase"
                                            @keyup.enter="handleSetup"
                                        />
                                    </v-col>

                                    <v-col cols="12">
                                        <v-alert type="warning" variant="tonal" class="mb-4">
                                            <strong>{{ tt('Lost passphrase = lost data') }}</strong><br/>
                                            {{ tt('There is no recovery option. If you forget your passphrase, all encrypted data is permanently lost.') }}
                                        </v-alert>
                                        <v-checkbox
                                            :disabled="creating"
                                            :label="tt('I understand that a lost passphrase means all my data is unrecoverable')"
                                            v-model="warningAccepted"
                                        />
                                    </v-col>

                                    <v-col cols="12">
                                        <v-btn
                                            block
                                            color="primary"
                                            type="submit"
                                            :disabled="!canSubmit"
                                            :loading="creating"
                                        >
                                            {{ tt('Create Vault') }}
                                        </v-btn>
                                    </v-col>

                                    <v-col cols="12" v-if="errorMessage">
                                        <v-alert type="error" variant="tonal">{{ errorMessage }}</v-alert>
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
import { computed } from 'vue';
import { useRouter } from 'vue-router';

import { useI18n } from '@/locales/helpers.ts';
import { useVaultSetupPageBase } from '@/views/base/VaultSetupPageBase.ts';

import { APPLICATION_LOGO_PATH } from '@/consts/asset.ts';

const { tt } = useI18n();
const router = useRouter();

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

const strengthColor = computed(() => {
    if (strengthScore.value <= 1) return 'error';
    if (strengthScore.value === 2) return 'warning';
    if (strengthScore.value === 3) return 'success';
    return 'success';
});

async function handleSetup(): Promise<void> {
    const success = await doSetup();
    if (success) {
        router.replace('/');
    }
}
</script>
