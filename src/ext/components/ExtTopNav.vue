<template>
    <template v-if="available">
        <router-link to="/ext/sales" :aria-current="isActive('/ext/sales') ? 'page' : undefined">
            <v-btn class="top-navigation-button ms-1" density="comfortable" variant="text"
                   :aria-label="tt('Sales')" :icon="true"
                   :active="isActive('/ext/sales')" :color="isActive('/ext/sales') ? 'primary' : 'default'">
                <v-icon :icon="isActive('/ext/sales') ? mdiCart : mdiCartOutline" size="24" />
                <v-tooltip activator="parent">{{ tt('Sales') }}</v-tooltip>
            </v-btn>
        </router-link>

        <router-link to="/ext/customers" :aria-current="isActive('/ext/customers') ? 'page' : undefined">
            <v-btn class="top-navigation-button ms-1" density="comfortable" variant="text"
                   :aria-label="tt('Customers')" :icon="true"
                   :active="isActive('/ext/customers')" :color="isActive('/ext/customers') ? 'primary' : 'default'">
                <v-icon :icon="isActive('/ext/customers') ? mdiAccountCash : mdiAccountCashOutline" size="24" />
                <v-tooltip activator="parent">{{ tt('Customers') }}</v-tooltip>
            </v-btn>
        </router-link>

        <router-link to="/ext/inventory" :aria-current="isActive('/ext/inventory') ? 'page' : undefined">
            <v-btn class="top-navigation-button ms-1" density="comfortable" variant="text"
                   :aria-label="tt('Inventory')" :icon="true"
                   :active="isActive('/ext/inventory')" :color="isActive('/ext/inventory') ? 'primary' : 'default'">
                <v-icon :icon="mdiPackageVariantClosed" size="24" />
                <v-tooltip activator="parent">{{ tt('Inventory') }}</v-tooltip>
            </v-btn>
        </router-link>
    </template>

    <!-- always visible for people who belong to other businesses, so nobody forgets which one they are in -->
    <v-chip class="ms-3" size="small" label :color="workingForSomeoneElse ? 'warning' : undefined"
            :prepend-icon="mdiStorefrontOutline" v-if="current && working.length > 1">
        {{ tt('ext.chip', { name: current.name, role: roleLabel(current.role) }) }}
    </v-chip>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { useRoute } from 'vue-router';

import { useExtI18n } from '@/ext/i18n.ts';
import { installBusinessHeader, useBusiness } from '@/ext/business.ts';
import { useBusinessFeatures } from '@/ext/features.ts';

import {
    mdiAccountCash,
    mdiAccountCashOutline,
    mdiCart,
    mdiCartOutline,
    mdiPackageVariantClosed,
    mdiStorefrontOutline
} from '@mdi/js';

const route = useRoute();
const { tt, roleLabel } = useExtI18n();
const { working, current, workingForSomeoneElse, ensureLoaded } = useBusiness();
const { available } = useBusinessFeatures();

function isActive(prefix: string): boolean {
    return route.path.startsWith(prefix);
}

installBusinessHeader();

onMounted(() => {
    ensureLoaded().catch(() => {
        // the app still works on the person's own business if the list cannot be loaded
    });
});
</script>
