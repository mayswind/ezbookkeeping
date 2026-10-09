<template>
    <v-dialog width="520" persistent :model-value="needsAcceptance">
        <v-card>
            <v-card-title>{{ tt('Terms of Service and Privacy Policy') }}</v-card-title>
            <v-card-text>
                <p class="mb-3">{{ tt('Please read and accept our terms to carry on. We ask again whenever they change.') }}</p>
                <ul class="ms-5 mb-4">
                    <li><a :href="LEGAL_TERMS_URL" target="_blank" rel="noopener">{{ tt('Terms of Service') }}</a></li>
                    <li><a :href="LEGAL_PRIVACY_URL" target="_blank" rel="noopener">{{ tt('Privacy Policy') }}</a></li>
                </ul>
                <v-checkbox density="compact" hide-details :disabled="saving" :label="tt('I have read and agree to the Terms of Service and the Privacy Policy')"
                            v-model="agreed" />
                <div class="text-error mt-2" v-if="problem">{{ problem }}</div>
            </v-card-text>
            <v-card-actions>
                <v-btn variant="text" :disabled="saving || loggingOut" @click="declineAndLogOut">{{ tt('I do not agree, log me out') }}</v-btn>
                <v-spacer />
                <v-btn color="primary" :loading="saving" :disabled="!agreed || loggingOut" @click="accept">{{ tt('Agree and continue') }}</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';

import { useExtI18n } from '@/ext/i18n.ts';
import { useRootStore } from '@/stores/index.ts';
import { useSettingsStore } from '@/stores/setting.ts';

import { describeError } from '@/ext/api.ts';
import { LEGAL_TERMS_URL, LEGAL_PRIVACY_URL } from '@/ext/legalVersion.ts';
import { useTerms } from '@/ext/terms.ts';

const { tt } = useExtI18n();
const router = useRouter();
const rootStore = useRootStore();
const settingsStore = useSettingsStore();
const { needsAcceptance, load, accept: acceptTerms } = useTerms();

const agreed = ref<boolean>(false);
const saving = ref<boolean>(false);
const loggingOut = ref<boolean>(false);
const problem = ref<string>('');

async function accept(): Promise<void> {
    saving.value = true;
    problem.value = '';

    try {
        await acceptTerms();
        agreed.value = false;
    } catch (error) {
        problem.value = describeError(error);
    } finally {
        saving.value = false;
    }
}

// Declining means not using the service: log out the way the menu does
async function declineAndLogOut(): Promise<void> {
    loggingOut.value = true;

    try {
        await rootStore.logout();
        settingsStore.clearAppSettings();
        await router.replace('/login');
    } catch (error) {
        problem.value = describeError(error);
    } finally {
        loggingOut.value = false;
    }
}

onMounted(() => {
    load().catch(() => {
        // if the answer cannot be fetched the app is not blocked; it asks again on the next page
    });
});
</script>
