<template>
    <v-text-field ref="textInput" type="text" :class="extraClass" :density="density" :readonly="!!readonly" :disabled="!!disabled"
                  :label="label" :placeholder="placeholder"
                  :persistent-placeholder="!!persistentPlaceholder"
                  v-model="currentValue"
                  @keydown="onKeyUpDown" @keyup="onKeyUpDown" @paste="onPaste">
    </v-text-field>
</template>

<script setup lang="ts">
import { VTextField } from 'vuetify/components/VTextField';

import { computed, useTemplateRef } from 'vue';

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

const textInput = useTemplateRef<VTextField>('textInput');

const extraClass = computed<string>(() => {
    return props.class || '';
});

function focus(): void {
    textInput.value?.focus();
}

defineExpose({
    focus
});
</script>
