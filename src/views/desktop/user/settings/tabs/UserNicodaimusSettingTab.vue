<template>
    <v-card>
        <v-card-title>{{ tt('nicodAImus API Key') }}</v-card-title>
        <v-card-text>
            <p class="mb-4">{{ tt('Link your nicodAImus API key to unlock AI-powered features like smart categorization and receipt scanning.') }}</p>

            <v-row>
                <v-col cols="12" md="8">
                    <v-text-field
                        type="password"
                        :label="tt('API Key')"
                        :placeholder="tt('Enter your nicodAImus API key')"
                        :disabled="linking"
                        v-model="apiKey"
                        @keyup.enter="handleLink"
                    />
                </v-col>
                <v-col cols="12" md="4" class="d-flex align-center">
                    <v-btn
                        color="primary"
                        :disabled="!apiKey || linking"
                        :loading="linking"
                        @click="handleLink"
                    >{{ tt('Link API Key') }}</v-btn>
                </v-col>
            </v-row>

            <v-alert v-if="successMessage" type="success" variant="tonal" class="mt-2">{{ successMessage }}</v-alert>
            <v-alert v-if="errorMessage" type="error" variant="tonal" class="mt-2">{{ errorMessage }}</v-alert>
        </v-card-text>
    </v-card>

    <v-card class="mt-4">
        <v-card-title>{{ tt('Subscription Tier') }}</v-card-title>
        <v-card-text>
            <v-chip :color="tierColor" size="large">{{ tierLabel }}</v-chip>

            <v-list class="mt-4" density="compact">
                <v-list-item>
                    <template v-slot:prepend>
                        <v-icon :color="canUseAI ? 'success' : 'grey'" :icon="canUseAI ? mdiCheckCircle : mdiCloseCircle" />
                    </template>
                    <v-list-item-title>{{ tt('AI Receipt Scanning & Smart Categorization') }}</v-list-item-title>
                    <v-list-item-subtitle v-if="!canUseAI">{{ tt('Requires alfred tier or higher') }}</v-list-item-subtitle>
                </v-list-item>
                <v-list-item>
                    <template v-slot:prepend>
                        <v-icon :color="canUseGroups ? 'success' : 'grey'" :icon="canUseGroups ? mdiCheckCircle : mdiCloseCircle" />
                    </template>
                    <v-list-item-title>{{ tt('Multi-User Groups') }}</v-list-item-title>
                    <v-list-item-subtitle v-if="!canUseGroups">{{ tt('Requires maurice tier or higher') }}</v-list-item-subtitle>
                </v-list-item>
            </v-list>
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { storeToRefs } from 'pinia';

import { useI18n } from '@/locales/helpers.ts';
import { useVaultStore } from '@/stores/vault.ts';

import {
    mdiCheckCircle,
    mdiCloseCircle
} from '@mdi/js';

const { tt } = useI18n();
const vaultStore = useVaultStore();
const { tier, canUseAI, canUseGroups } = storeToRefs(vaultStore);

const apiKey = ref('');
const linking = ref(false);
const successMessage = ref('');
const errorMessage = ref('');

const tierColor = computed(() => {
    switch (tier.value) {
        case 'jared': return 'purple';
        case 'maurice': return 'primary';
        case 'alfred': return 'success';
        default: return 'grey';
    }
});

const tierLabel = computed(() => {
    switch (tier.value) {
        case 'jared': return 'nicodAImus jared';
        case 'maurice': return 'nicodAImus maurice';
        case 'alfred': return 'nicodAImus alfred';
        default: return 'Free';
    }
});

async function handleLink(): Promise<void> {
    if (!apiKey.value || linking.value) return;

    linking.value = true;
    successMessage.value = '';
    errorMessage.value = '';

    try {
        const newTier = await vaultStore.linkApiKey(apiKey.value);
        successMessage.value = `API key linked successfully. Tier: ${newTier}`;
        apiKey.value = '';
    } catch (e: unknown) {
        errorMessage.value = e instanceof Error ? e.message : 'Failed to link API key';
    } finally {
        linking.value = false;
    }
}
</script>
