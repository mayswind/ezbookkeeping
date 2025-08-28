<template>
    <v-pagination :density="density"
                  :disabled="disabled"
                  :total-visible="totalVisible ?? 7"
                  :length="totalPageCount"
                  v-model="currentPage">
        <template #item="{ key, page, isActive }">
            <v-btn variant="text"
                   :density="density"
                   :disabled="disabled"
                   :icon="true"
                   :color="isActive ? 'primary' : 'default'"
                   @click="currentPage = parseInt(page)"
                   v-if="page !== '...'"
            >
                <span>{{ getDisplayPage(page) }}</span>
            </v-btn>
            <v-btn variant="text"
                   color="default"
                   :density="density"
                   :disabled="disabled"
                   :icon="true"
                   v-if="page === '...'"
            >
                <span>{{ page }}</span>
                <v-menu activator="parent"
                        :disabled="disabled"
                        :close-on-content-click="false"
                        v-model="showMenus[key]">
                    <v-list>
                        <v-list-item class="text-sm" :density="density">
                            <v-list-item-title class="cursor-pointer">
                                <v-autocomplete width="100"
                                                item-title="name"
                                                item-value="value"
                                                auto-select-first="exact"
                                                :density="density"
                                                :items="allPages"
                                                :no-data-text="tt('No results')"
                                                v-model="currentPage"/>
                            </v-list-item-title>
                        </v-list-item>
                    </v-list>
                </v-menu>
            </v-btn>
        </template>
    </v-pagination>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import type { NameNumeralValue } from '@/core/base.ts';
import type { NumeralSystem } from '@/core/numeral.ts';
import type { ComponentDensity } from '@/lib/ui/desktop.ts';

const props = defineProps<{
    density?: ComponentDensity;
    disabled?: boolean;
    totalPageCount: number;
    totalVisible?: number;
    modelValue: number;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: number): void;
}>();

const { tt, getCurrentNumeralSystemType } = useI18n();

const showMenus = ref<Record<string, boolean>>({});

const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());

function getDisplayPage(page: number | string): string {
    return numeralSystem.value.replaceWesternArabicDigitsToLocalizedDigits(page.toString());
}

const allPages = computed<NameNumeralValue[]>(() => {
    const pages: NameNumeralValue[] = [];

    for (let i = 1; i <= props.totalPageCount; i++) {
        pages.push({ value: i, name: getDisplayPage(i) });
    }

    return pages;
});

const currentPage = computed<number>({
    get: () => props.modelValue,
    set: (value) => {
        if (value && value >= 1 && value <= props.totalPageCount) {
            emit('update:modelValue', value);

            for (const key in showMenus.value) {
                if (!Object.prototype.hasOwnProperty.call(showMenus.value, key)) {
                    continue;
                }

                showMenus.value[key] = false;
            }
        }
    }
});
</script>
