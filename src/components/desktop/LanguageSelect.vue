<template>
    <v-select
        item-title="nativeDisplayName"
        item-value="languageTag"
        persistent-placeholder
        :disabled="disabled"
        :label="label"
        :placeholder="placeholder"
        :items="allLanguages"
        v-model="currentLocaleValue"
    >
        <template #item="{ props, internalItem }">
            <v-list-item :value="internalItem.value" v-bind="props">
                <template #title>
                    <v-list-item-title>
                        <div class="d-flex align-center">
                            <span>{{ internalItem.title }}</span>
                            <v-spacer style="min-width: 40px" />
                            <v-icon :icon="mdiCheck" v-if="isLanguageSelected(internalItem.raw.languageTag)" />
                            <span class="text-field-append-text" v-if="!isLanguageSelected(internalItem.raw.languageTag)">{{ internalItem.raw.displayName }}</span>
                        </div>
                    </v-list-item-title>
                </template>
            </v-list-item>
        </template>
    </v-select>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { type LanguageSelectBaseProps, type LanguageSelectBaseEmits, useLanguageSelectButtonBase } from '@/components/base/LanguageSelectBase.ts';

import { useI18n } from '@/locales/helpers.ts';

import {
    mdiCheck
} from '@mdi/js';

interface DesktopLanguageSelectProps extends LanguageSelectBaseProps {
    label?: string;
    placeholder?: string;
}

const props = defineProps<DesktopLanguageSelectProps>();
const emit = defineEmits<LanguageSelectBaseEmits>();

const { getCurrentLanguageTag } = useI18n();

const {
    allLanguages,
    updateLanguage,
    isLanguageSelected
} = useLanguageSelectButtonBase(props, emit);

const currentLocaleValue = computed<string>({
    get: () => {
        if (props.useModelValue) {
            return props.modelValue ?? '';
        } else {
            return getCurrentLanguageTag()
        }
    },
    set: (value: string) => {
        updateLanguage(value);
    }
});
</script>
