<template>
    <v-autocomplete
        item-title="displayName"
        item-value="currencyCode"
        auto-select-first
        persistent-placeholder
        :multiple="multiple"
        :chips="multiple"
        :closable-chips="multiple"
        :density="density"
        :variant="variant"
        :disabled="disabled"
        :label="label"
        :placeholder="placeholder"
        :items="allCurrencies"
        :no-data-text="tt('No results')"
        :custom-filter="filterCurrency"
        v-model="currentCurrencyValue"
    >
        <template #append-inner v-if="!multiple">
            <small class="text-field-append-text smaller text-no-wrap" v-if="currentCurrencyValue !== ACCOUNT_CURRENCY_NOT_SET_VALUE">{{ currentCurrencyValue }}</small>
        </template>

        <template #chip="{ props, item }" v-if="multiple">
            <v-chip v-bind="props" :text="item.currencyCode" />
        </template>

        <template #item="{ props, internalItem }">
            <v-list-item :value="internalItem.value" v-bind="props">
                <template #title>
                    <v-list-item-title>
                        <div class="d-flex align-center">
                            <span>{{ internalItem.title }}</span>
                            <v-spacer style="min-width: 40px" />
                            <v-icon :icon="mdiCheck" v-if="isSelected(internalItem.raw.currencyCode)" />
                            <small class="text-field-append-text" v-if="!isSelected(internalItem.raw.currencyCode)">{{ internalItem.raw.currencyCode }}</small>
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

import { NormalizedText } from '@/core/text.ts';
import type { LocalizedCurrencyInfo } from '@/core/currency.ts';
import { ACCOUNT_CURRENCY_NOT_SET_VALUE } from '@/consts/currency.ts';

import type { ComponentDensity, InputVariant } from '@/lib/ui/desktop.ts';

import {
    mdiCheck
} from '@mdi/js';

const props = defineProps<{
    density?: ComponentDensity;
    variant?: InputVariant;
    disabled?: boolean;
    label?: string;
    placeholder?: string;
    withNotSet?: boolean;
    multiple?: boolean;
    modelValue: string | string[];
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string | string[]): void;
}>();

const { tt, getAllCurrencies } = useI18n();

const allCurrencies = computed<LocalizedCurrencyInfo[]>(() => getAllCurrencies(props.withNotSet));

const currentCurrencyValue = computed<string | string[] | null>({
    get: () => props.modelValue,
    set: (value: string | string[] | null) => {
        if (value === null) {
            emit('update:modelValue', props.multiple ? [] : '');
        } else {
            emit('update:modelValue', value);
        }
    }
});

function isSelected(currencyCode: string): boolean {
    const value = currentCurrencyValue.value;
    return Array.isArray(value) ? value.includes(currencyCode) : value === currencyCode;
}

function filterCurrency(value: string, query: string, item?: { value: unknown, raw: LocalizedCurrencyInfo }): boolean {
    if (!item) {
        return false;
    }

    const normalizedFilterContent = NormalizedText.normalizeForSearch(query || '');

    if (!normalizedFilterContent) {
        return true;
    }

    return NormalizedText.normalizeForSearch(item.raw.displayName).indexOf(normalizedFilterContent) >= 0
        || NormalizedText.normalizeForSearch(item.raw.currencyCode).indexOf(normalizedFilterContent) >= 0;
}
</script>
