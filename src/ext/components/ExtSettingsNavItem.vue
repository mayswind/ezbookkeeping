<template>
    <li class="nav-link">
        <router-link to="/settings/business">
            <v-icon class="nav-item-icon" :icon="mdiStorefrontOutline"/>
            <span class="nav-item-title">{{ tt('Business Features') }}</span>
        </router-link>
    </li>
</template>

<script setup lang="ts">
import { useExtI18n } from '@/ext/i18n.ts';

import { mdiStorefrontOutline } from '@mdi/js';

const { tt } = useExtI18n();
</script>
