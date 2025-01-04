<template>
    <div class="btn-vertical-group d-flex flex-column">
        <v-btn border :key="idx"
               :color="value === button.value ? 'primary' : 'default'"
               :variant="value === button.value ? 'tonal' : 'outlined'" :disabled="disabled"
               v-for="(button, idx) in buttons"
               @click="value = button.value">
            {{ button.name }}
        </v-btn>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface Button {
    name: string,
    value: unknown
}

const props = defineProps<{
    modelValue: unknown;
    buttons: Button[];
    disabled?: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: unknown): void
}>();

const value = computed<unknown>({
    get: () => {
        return props.modelValue;
    },
    set: value => {
        if (value === props.modelValue) {
            return;
        }

        emit('update:modelValue', value);
    }
});
</script>

<style>
.btn-vertical-group .v-btn:not(:first-child) {
    border-top-left-radius: inherit;
    border-top-right-radius: inherit;
}

.btn-vertical-group .v-btn:not(:last-child) {
    border-bottom: 0;
    border-bottom-left-radius: inherit;
    border-bottom-right-radius: inherit;
}
</style>
