<template>
    <f7-button class="language-select-button" small popover-open=".lang-popover-menu" :disabled="disabled" :text="currentLanguageName"></f7-button>

    <f7-popover class="lang-popover-menu" @popover:open="onPopoverOpen">
        <f7-list dividers>
            <f7-list-item link="#" no-chevron popover-close
                          :class="{ 'list-item-selected': isLanguageSelected(lang.languageTag) }"
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
import { type LanguageSelectBaseProps, type LanguageSelectBaseEmits, useLanguageSelectButtonBase } from '@/components/base/LanguageSelectBase.ts';

import { scrollToSelectedItem } from '@/lib/ui/common.ts';
import { type Framework7Dom } from '@/lib/ui/mobile.ts';

const props = defineProps<LanguageSelectBaseProps>();
const emit = defineEmits<LanguageSelectBaseEmits>();

const {
    allLanguages,
    currentLanguageName,
    updateLanguage,
    isLanguageSelected
} = useLanguageSelectButtonBase(props, emit);

function onPopoverOpen(event: { $el: Framework7Dom }): void {
    scrollToSelectedItem(event.$el[0], '.popover-inner', '.popover-inner', 'li.list-item-selected');
}
</script>

<style>
.language-select-button {
    display: initial;
    padding: 8px 10px 8px 10px;
}
</style>
