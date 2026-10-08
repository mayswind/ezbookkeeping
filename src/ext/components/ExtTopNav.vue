<template>
    <router-link to="/ext/inventory" :aria-current="isActive ? 'page' : undefined">
        <v-btn class="top-navigation-button ms-1" density="comfortable" variant="text"
               :aria-label="tt('Inventory')" :icon="true"
               :active="isActive" :color="isActive ? 'primary' : 'default'">
            <v-icon :icon="mdiPackageVariantClosed" size="24" />
            <v-tooltip activator="parent">{{ tt('Inventory') }}</v-tooltip>
        </v-btn>
    </router-link>

    <!-- always visible for people who belong to other businesses, so nobody forgets which one they are in -->
    <v-chip class="ms-3" size="small" label :color="workingForSomeoneElse ? 'warning' : undefined"
            :prepend-icon="mdiStorefrontOutline" v-if="current && working.length > 1">
        {{ tt('ext.chip', { name: current.name, role: roleLabel(current.role) }) }}
    </v-chip>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';

import { useExtI18n } from '@/ext/i18n.ts';
import { installBusinessHeader, useBusiness } from '@/ext/business.ts';

import { mdiPackageVariantClosed, mdiStorefrontOutline } from '@mdi/js';

const route = useRoute();
const { tt, roleLabel } = useExtI18n();
const { working, current, workingForSomeoneElse, ensureLoaded } = useBusiness();

const isActive = computed<boolean>(() => route.path.startsWith('/ext/inventory'));

installBusinessHeader();

onMounted(() => {
    ensureLoaded().catch(() => {
        // the app still works on the person's own business if the list cannot be loaded
    });
});
</script>
