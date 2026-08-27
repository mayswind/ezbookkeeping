<template>
    <div class="overview-widget__header">
        <v-icon class="overview-widget__icon" color="primary" :icon="icon" size="20" />
        <span class="overview-widget__title text-title-small" :title="title">{{ title }}</span>
        <slot />
    </div>
</template>

<script setup lang="ts">
defineProps<{
    title: string;
    icon: string;
}>();
</script>
