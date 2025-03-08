<template>
    <v-autocomplete
        item-title="displayName"
        item-value="currencyCode"
        auto-select-first
        persistent-placeholder
        :disabled="disabled"
        :label="label"
        :placeholder="placeholder"
        :items="allCurrencies"
        :no-data-text="tt('No results')"
        :custom-filter="filterCurrency"
        v-model="currentCurrencyValue"
    >
        <template #append-inner>
            <small class="text-field-append-text smaller">{{ currentCurrencyValue }}</small>
        </template>

        <template #item="{ props, item }">
            <v-list-item :value="item.value" v-bind="props">
                <template #title>
                    <v-list-item-title>
                        <div class="d-flex align-center">
                            <span>{{ item.title }}</span>
                            <v-spacer style="min-width: 40px" />
                            <v-icon :icon="mdiCheck" v-if="currentCurrencyValue === item.raw.currencyCode" />
                            <small class="text-field-append-text" v-if="currentCurrencyValue !== item.raw.currencyCode">{{ item.raw.currencyCode }}</small>
                        </div>
                    </v-list-item-title>
                </template>
            </v-list-item>
        </template>
    </v-autocomplete>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from '@/locales/helpers.ts';

import type { LocalizedCurrencyInfo } from '@/core/currency.ts';

import {
    mdiCheck
} from '@mdi/js';

const props = defineProps<{
    disabled?: boolean;
    label?: string;
    placeholder?: string;
    modelValue: string;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
}>();

const { tt, getAllCurrencies } = useI18n();

const allCurrencies = computed<LocalizedCurrencyInfo[]>(() => getAllCurrencies());

const currentCurrencyValue = computed<string | null>({
    get: () => props.modelValue,
    set: (value: string | null) => {
        if (value === null) {
            emit('update:modelValue', '');
        } else {
            emit('update:modelValue', value);
        }
    }
});

function filterCurrency(value: string, query: string, item?: { value: unknown, raw: LocalizedCurrencyInfo }): boolean {
    if (!item) {
        return false;
    }

    const lowerCaseFilterContent = query.toLowerCase() || '';

    if (!lowerCaseFilterContent) {
        return true;
    }

    return item.raw.displayName.toLowerCase().indexOf(lowerCaseFilterContent) >= 0
        || item.raw.currencyCode.toLowerCase().indexOf(lowerCaseFilterContent) >= 0;
}
</script>
