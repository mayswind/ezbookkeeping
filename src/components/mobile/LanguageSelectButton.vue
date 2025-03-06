<template>
    <f7-button small popover-open=".lang-popover-menu" :disabled="disabled" :text="currentLanguageName"></f7-button>

    <f7-popover class="lang-popover-menu">
        <f7-list dividers>
            <f7-list-item link="#" no-chevron popover-close
                          :key="lang.languageTag"
                          :title="lang.nativeDisplayName"
                          v-for="lang in allLanguages"
                          @click="updateLanguage(lang.languageTag)">
                <template #after>
                    <f7-icon class="list-item-checked-icon" f7="checkmark_alt" v-if="isLanguageSelected(lang.languageTag)"></f7-icon>
                    <span v-if="!isLanguageSelected(lang.languageTag)">{{ lang.displayName }}</span>
                </template>
            </f7-list-item>
        </f7-list>
    </f7-popover>
</template>

<script setup lang="ts">
import { type LanguageSelectButtonBaseProps, type LanguageSelectButtonBaseEmits, useLanguageSelectButtonBase } from '@/components/base/LanguageSelectButtonBase.ts';

const props = defineProps<LanguageSelectButtonBaseProps>();
const emit = defineEmits<LanguageSelectButtonBaseEmits>();

const {
    allLanguages,
    currentLanguageName,
    updateLanguage,
    isLanguageSelected
} = useLanguageSelectButtonBase(props, emit);
</script>
