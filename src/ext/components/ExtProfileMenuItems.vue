<template>
    <template v-if="working.length > 1">
        <v-divider class="my-1" />
        <v-list-subheader>{{ tt('Working in') }}</v-list-subheader>
        <v-list-item :prepend-icon="isCurrent(business) ? mdiCheckCircle : mdiCircleOutline"
                     :active="isCurrent(business)"
                     :key="business.ownerUid"
                     :title="business.name"
                     :subtitle="business.status === 'owner' ? tt('ext.yourOwnBusiness') : roleLabel(business.role)"
                     @click="selectBusiness(business)"
                     v-for="business in working">
        </v-list-item>
    </template>

    <v-divider class="my-1" v-if="teamAvailable" />
    <v-list-item :prepend-icon="mdiAccountGroupOutline" to="/ext/team" v-if="teamAvailable">
        <v-list-item-title>
            {{ tt('Team') }}
            <v-chip class="ms-2" size="x-small" color="error" v-if="invitations.length > 0">{{ invitations.length }}</v-chip>
        </v-list-item-title>
    </v-list-item>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';

import { useExtI18n } from '@/ext/i18n.ts';
import { installBusinessHeader, switchBusiness, useBusiness } from '@/ext/business.ts';
import { useBusinessFeatures } from '@/ext/features.ts';
import type { BusinessInfo } from '@/ext/types.ts';

import {
    mdiCheckCircle,
    mdiCircleOutline,
    mdiAccountGroupOutline
} from '@mdi/js';

const { tt, roleLabel } = useExtI18n();
const { working, invitations, current, refresh } = useBusiness();
const { teamAvailable } = useBusinessFeatures();

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
        // the menu still shows Team if the list cannot be loaded
    });
});
</script>
