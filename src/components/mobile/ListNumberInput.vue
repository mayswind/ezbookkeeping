<template>
    <f7-list-input
        type="text"
        :readonly="!!readonly"
        :disabled="!!disabled"
        :label="label"
        :placeholder="placeholder"
        v-model:value="currentValue"
        @keydown="onKeyUpDown"
        @keyup="onKeyUpDown"
        @paste="onPaste"
    ></f7-list-input>
</template>

<script setup lang="ts">
import { type NumberInputProps, type NumberInputEmits, useNumberInputBase } from '@/components/base/NumberInputBase.ts';

const props = defineProps<NumberInputProps>();
const emit = defineEmits<NumberInputEmits>();

const {
    currentValue,
    onKeyUpDown,
    onPaste
} = useNumberInputBase(props, emit);
</script>
