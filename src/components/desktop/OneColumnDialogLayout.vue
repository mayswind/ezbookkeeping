<template>
    <v-card class="dialog-layout-card">
        <v-card-title class="dialog-layout-card-title px-0 pt-3 pb-2">
            <div class="d-flex align-center justify-center ms-5">
                <div class="d-flex align-center">
                    <span class="dialog-layout-card-title-text text-title-medium text-wrap">{{ title }}</span>
                    <v-progress-circular indeterminate size="22" class="ms-2" v-if="loading"></v-progress-circular>
                </div>
                <slot name="after-title" />
                <v-spacer/>
                <slot name="toolbar" />
                <v-divider vertical class="ms-2" v-if="cancelButtonTitle" />
                <div class="d-flex align-items-center mx-3" v-if="cancelButtonTitle">
                    <v-icon size="22" :icon="mdiClose" :disabled="disabled" @click="emit('cancel')"></v-icon>
                    <v-tooltip activator="parent">{{ cancelButtonTitle }}</v-tooltip>
                </div>
            </div>
            <slot name="subtitle" />
        </v-card-title>
        <v-divider />
        <v-card-text class="flex-grow-1 overflow-y-auto" :class="contentClass" :style="contentStyle">
            <slot name="content" />
        </v-card-text>
        <v-divider v-if="$slots['footer']" />
        <v-card-text :class="footerClass ?? 'py-3'" :style="footerStyle" v-if="$slots['footer']">
            <div class="w-100 d-flex justify-center flex-wrap gap-4">
                <slot name="footer" />
            </div>
        </v-card-text>
    </v-card>
</template>

<script setup lang="ts">
import {
    mdiClose
} from '@mdi/js';

defineProps<{
    title: string;
    cancelButtonTitle?: string;
    loading?: boolean;
    disabled?: boolean;
    contentClass?: string;
    contentStyle?: string;
    footerClass?: string;
    footerStyle?: string;
}>();

const emit = defineEmits<{
    (e: 'cancel'): void;
}>();
</script>
