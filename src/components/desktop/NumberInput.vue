<template>
    <v-text-field type="text" :class="extraClass" :density="density" :readonly="!!readonly" :disabled="!!disabled"
                  :label="label" :placeholder="placeholder"
                  :persistent-placeholder="!!persistentPlaceholder"
                  v-model="currentValue"
                  @keydown="onKeyUpDown" @keyup="onKeyUpDown" @paste="onPaste">
    </v-text-field>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { type NumberInputProps, type NumberInputEmits, useNumberInputBase } from '@/components/base/NumberInputBase.ts';

import type { ComponentDensity } from '@/lib/ui/desktop.ts';

interface DesktopNumberInputProps extends NumberInputProps {
    class?: string;
    density?: ComponentDensity;
    persistentPlaceholder?: boolean;
}

const props = defineProps<DesktopNumberInputProps>();
const emit = defineEmits<NumberInputEmits>();

const {
    currentValue,
    onKeyUpDown,
    onPaste
} = useNumberInputBase(props, emit);

const extraClass = computed<string>(() => {
    return props.class || '';
});
</script>
