<template>
    <div class="ext-business-banner d-flex align-center flex-wrap ga-2 px-4 py-2" role="status" v-if="workingForSomeoneElse && current">
        <v-icon :icon="mdiAlertOutline" size="22" />
        <span class="font-weight-medium">
            {{ tt('ext.banner', { name: current.name, role: roleLabel(current.role) }) }}
        </span>
        <v-spacer />
        <v-btn size="small" variant="flat" color="white" @click="backToOwn" v-if="own">
            {{ tt('Switch back to my business') }}
        </v-btn>
    </div>
</template>

<script setup lang="ts">
import { useExtI18n } from '@/ext/i18n.ts';
import { installBusinessHeader, switchBusiness, useBusiness } from '@/ext/business.ts';

import { mdiAlertOutline } from '@mdi/js';

const { tt, roleLabel } = useExtI18n();
const { own, current, workingForSomeoneElse, ensureLoaded } = useBusiness();

installBusinessHeader();
ensureLoaded().catch(() => {
    // without the list there is nothing to show; the app keeps working on the person's own business
});

function backToOwn(): void {
    if (own.value) {
        switchBusiness(own.value.ownerUid);
    }
}
</script>

<style scoped>
.ext-business-banner {
    position: sticky;
    top: 0;
    z-index: 1200;
    background: rgb(var(--v-theme-warning));
    color: rgb(var(--v-theme-on-warning));
}
</style>
