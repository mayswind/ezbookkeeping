<template>
    <div>
        <div class="d-flex align-center ga-2">
            <v-select class="flex-1-1-0" item-title="name" item-value="value" :label="label" :items="filterTypes" v-model="filterType" />
            <amount-input class="flex-1-1-0" :currency="defaultCurrency" :show-currency="false"
                          :label="filterTypeAndAmounts?.filterType?.paramCount === 2 ? tt('Minimum Amount') : tt('Amount')"
                          v-model="amount1" v-if="filterTypeAndAmounts?.filterType" />
            <span v-if="filterTypeAndAmounts?.filterType?.paramCount === 2">{{ tt('format.misc.rangeSeparator') }}</span>
            <amount-input class="flex-1-1-0" :currency="defaultCurrency" :show-currency="false"
                          :label="filterTypeAndAmounts?.filterType?.paramCount === 2 ? tt('Maximum Amount') : tt('Amount')"
                          v-model="amount2" v-if="filterTypeAndAmounts?.filterType?.paramCount === 2" />
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';

import type { NameValue } from '@/core/base.ts';
import { AmountFilterType } from '@/core/numeral.ts';

const props = defineProps<{
    modelValue: string;
    label: string;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
}>();

const { tt } = useI18n();

const userStore = useUserStore();

const defaultCurrency = computed<string>(() => userStore.currentUserDefaultCurrency);

const filterTypes = computed<NameValue[]>(() => [
    { name: tt('All'), value: '' },
    ...AmountFilterType.values().map(type => ({ name: tt(type.name), value: type.type }))
]);

const filterTypeAndAmounts = computed<{ filterType: AmountFilterType, params: number[] } | undefined>(() => AmountFilterType.parseTextualFilter(props.modelValue));

const filterType = computed<string>({
    get: () => filterTypeAndAmounts.value?.filterType?.type ?? '',
    set: value => emit('update:modelValue', AmountFilterType.valueOf(value)?.toTextualFilter(amount1.value, amount2.value) ?? '')
});

const amount1 = computed<number>({
    get: () => filterTypeAndAmounts?.value?.params?.[0] ?? 0,
    set: value => emit('update:modelValue', filterTypeAndAmounts.value?.filterType?.toTextualFilter(value, amount2.value) ?? '')
});

const amount2 = computed<number>({
    get: () => filterTypeAndAmounts?.value?.params?.[1] ?? 0,
    set: value => emit('update:modelValue', filterTypeAndAmounts.value?.filterType?.toTextualFilter(amount1.value, value) ?? '')
});
</script>
