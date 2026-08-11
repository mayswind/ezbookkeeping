<template>
    <v-btn-toggle class="segmented-button" density="compact" mandatory="force" size="small"
                  :disabled="disabled" :model-value="modelValue"
                  @update:model-value="updateValue">
        <v-btn key="false" variant="text" color="default" :ripple="false"
               :disabled="disabled" :value="false" @click="updateValue(!modelValue)">
            <span>{{ falseName }}</span>
        </v-btn>
        <v-btn key="true" variant="text" color="default" :ripple="false"
               :disabled="disabled" :value="true" @click="updateValue(!modelValue)">
            <span>{{ trueName }}</span>
        </v-btn>
    </v-btn-toggle>
</template>

<script setup lang="ts">
defineProps<{
    disabled?: boolean;
    falseName?: string;
    trueName?: string;
    modelValue: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: boolean): void;
}>();

function updateValue(value: boolean): void {
    emit('update:modelValue', value);
}
</script>

<style scoped>
.segmented-button.v-btn-toggle {
    --v-theme-overlay-multiplier: 0;
    width: auto;
    height: auto !important;
    background: rgba(var(--v-theme-on-surface), var(--v-focus-opacity));
    display: inline-flex;
    padding: 2px;
    border: none;
    border-radius: var(--ebk-radius-sm);
    overflow-x: auto;

    .v-btn--variant-text:not(.v-btn--icon) {
        padding-inline: 6px;
    }

    .v-btn {
        width: auto !important;
        height: 26px !important;
        min-width: 0;
        min-height: 26px;
        border-radius: calc(var(--ebk-radius-sm) - 2px) !important;
        color: rgba(var(--v-theme-on-surface), var(--v-medium-emphasis-opacity));
    }

    .v-btn:hover {
        color: rgba(var(--v-theme-on-surface), var(--v-high-emphasis-opacity));
    }

    .v-btn--active {
        color: rgba(var(--v-theme-on-surface), var(--v-high-emphasis-opacity));
        background: rgb(var(--v-theme-surface));
        box-shadow: 0 1px 3px rgba(var(--v-theme-on-surface), 0.16);

        &.v-btn--disabled {
            opacity: 0.45;
        }
    }
}
</style>
