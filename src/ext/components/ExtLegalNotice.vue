<template>
    <div class="ext-legal-notice">
        {{ tt('By continuing you accept our') }}
        <a :href="LEGAL_TERMS_URL" target="_blank" rel="noopener">{{ tt('Terms of Service') }}</a>
        {{ tt('and') }}
        <a :href="LEGAL_PRIVACY_URL" target="_blank" rel="noopener">{{ tt('Privacy Policy') }}</a>.
    </div>
</template>

<script setup lang="ts">
// Plain HTML only, so the same notice works in the desktop (Vuetify) and the mobile (Framework7) screens.
import { useExtI18n } from '@/ext/i18n.ts';

import { LEGAL_TERMS_URL, LEGAL_PRIVACY_URL } from '@/ext/legalVersion.ts';

const { tt } = useExtI18n();
</script>

<style>
.ext-legal-notice {
    text-align: center;
    font-size: 0.8rem;
    opacity: 0.8;
    margin: 0.75rem 1rem;
}
</style>
