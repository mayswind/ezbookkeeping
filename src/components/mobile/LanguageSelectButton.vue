<template>
    <f7-button class="language-select-button" small
               :class="{ 'disabled': disabled || loading || settingLanguage }"
               :popover-open="disabled || loading || settingLanguage ? undefined : '#' + languagePopoverId"
               :text="currentLanguageName"></f7-button>

    <f7-popover :id="languagePopoverId" class="lang-popover-menu" @popover:open="onPopoverOpen">
        <f7-list dividers>
            <f7-list-item link="#" no-chevron popover-close
                          :class="{ 'list-item-selected': isLanguageSelected(lang.languageTag) }"
                          :key="lang.languageTag"
                          :title="lang.nativeDisplayName"
                          v-for="lang in allLanguages"
                          @click="changeLanguage(lang.languageTag)">
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
import { ref, useId } from 'vue';

import { scrollToSelectedItem } from '@/lib/ui/common.ts';
import { type Framework7Dom, showLoading, hideLoading } from '@/lib/ui/mobile.ts';

const props = defineProps<LanguageSelectBaseProps>();
const emit = defineEmits<LanguageSelectBaseEmits>();

const languagePopoverId = ref<string>('language-popover-' + useId());

const {
    settingLanguage,
    allLanguages,
    currentLanguageName,
    updateLanguage,
    isLanguageSelected
} = useLanguageSelectButtonBase(props, emit);

function changeLanguage(languageTag: string): void {
    if (props.disabled || props.loading || settingLanguage.value) {
        return;
    }

    showLoading(() => settingLanguage.value);

    updateLanguage(languageTag).catch(error => {
        emit('error', error.message || error);
    }).finally(() => {
        hideLoading();
    });
}

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
