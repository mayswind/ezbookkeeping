<template>
    <v-btn class="top-navigation-button ms-1" density="comfortable" variant="text"
           :aria-label="tt('Business')" :icon="true"
           :active="isActive"
           :color="isActive || workingForSomeoneElse ? 'primary' : 'default'">
        <v-badge dot color="error" :model-value="invitations.length > 0">
            <v-icon :icon="isActive ? mdiStorefront : mdiStorefrontOutline" size="24" />
        </v-badge>
        <v-menu activator="parent" width="260" location="bottom end" offset="14px">
            <v-list>
                <v-list-subheader>{{ tt('Working in') }}</v-list-subheader>
                <v-list-item :prepend-icon="isCurrent(business) ? mdiCheck : undefined"
                             :key="business.ownerUid"
                             :title="business.status === 'owner' ? tt('My own business') : business.name"
                             :subtitle="tt(business.role)"
                             @click="selectBusiness(business)"
                             v-for="business in working">
                </v-list-item>

                <v-divider class="my-1" />
                <v-list-item :prepend-icon="mdiPackageVariantClosed" :title="tt('Inventory')" to="/ext/inventory"></v-list-item>
                <v-list-item :prepend-icon="mdiAccountGroupOutline" to="/ext/team">
                    <v-list-item-title>
                        {{ tt('Team') }}
                        <v-chip class="ms-2" size="x-small" color="error" v-if="invitations.length > 0">{{ invitations.length }}</v-chip>
                    </v-list-item-title>
                </v-list-item>
            </v-list>
        </v-menu>
    </v-btn>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';

import { useI18n } from '@/locales/helpers.ts';

import { installBusinessHeader, switchBusiness, useBusiness } from '@/ext/business.ts';
import type { BusinessInfo } from '@/ext/types.ts';

import {
    mdiCheck,
    mdiStorefront,
    mdiStorefrontOutline,
    mdiPackageVariantClosed,
    mdiAccountGroupOutline
} from '@mdi/js';

const route = useRoute();
const { tt } = useI18n();
const { working, invitations, current, workingForSomeoneElse, refresh } = useBusiness();

const isActive = computed<boolean>(() => route.path.startsWith('/ext/'));

function isCurrent(business: BusinessInfo): boolean {
    return !!current.value && current.value.ownerUid === business.ownerUid;
}

function selectBusiness(business: BusinessInfo): void {
    if (!isCurrent(business)) {
        switchBusiness(business.ownerUid);
    }
}

installBusinessHeader();

onMounted(() => {
    refresh().catch(() => {
        // the menu still works for the person's own business if the list cannot be loaded
    });
});
</script>
