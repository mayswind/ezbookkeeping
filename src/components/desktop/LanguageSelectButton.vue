<template>
    <v-menu location="bottom" max-height="500" @update:model-value="onMenuStateChanged">
        <template #activator="{ props }">
            <v-btn class="language-select-button" variant="text" :disabled="disabled || loading || settingLanguage" :loading="loading || settingLanguage" v-bind="props">
                {{ currentLanguageName }}
                <template #loader>
                    <svg class="language-select-button-loading-border" aria-hidden="true">
                        <rect x="1" y="1" rx="4" pathLength="100" />
                    </svg>
                </template>
            </v-btn>
        </template>
        <v-list ref="languageMenu">
            <v-list-item :key="lang.languageTag" :value="lang.languageTag"
                         :class="{ 'list-item-selected': isLanguageSelected(lang.languageTag) }"
                         v-for="lang in allLanguages">
                <v-list-item-title class="cursor-pointer" @click="changeLanguage(lang.languageTag)">
                    <div class="d-flex align-center">
                        <span>{{ lang.nativeDisplayName }}</span>
                        <v-spacer style="min-width: 40px" />
                        <v-icon color="primary" :icon="mdiCheck" v-if="isLanguageSelected(lang.languageTag)" />
                        <span class="text-field-append-text" v-if="!isLanguageSelected(lang.languageTag)">{{ lang.displayName }}</span>
                    </div>
                </v-list-item-title>
            </v-list-item>
        </v-list>
    </v-menu>
</template>

<script setup lang="ts">
import { VList } from 'vuetify/components/VList';
import { type LanguageSelectBaseProps, type LanguageSelectBaseEmits, useLanguageSelectButtonBase } from '@/components/base/LanguageSelectBase.ts';

import { useTemplateRef, nextTick } from 'vue';

import { scrollToSelectedItem } from '@/lib/ui/common.ts';

import {
    mdiCheck
} from '@mdi/js';

const props = defineProps<LanguageSelectBaseProps>();
const emit = defineEmits<LanguageSelectBaseEmits>();

const {
    settingLanguage,
    allLanguages,
    currentLanguageName,
    updateLanguage,
    isLanguageSelected
} = useLanguageSelectButtonBase(props, emit);

const languageMenu = useTemplateRef<VList>('languageMenu');

function changeLanguage(languageTag: string): void {
    updateLanguage(languageTag).catch(error => {
        emit('error', error.message || error);
    });
}

function onMenuStateChanged(state: boolean): void {
    if (state) {
        nextTick(() => {
            scrollToSelectedItem(languageMenu.value?.$el, null, null, 'div.v-list-item.list-item-selected');
        });
    }
}
</script>

<style scoped>
.language-select-button.v-btn--loading :deep(.v-btn__content) {
    opacity: 1;
}

.language-select-button-loading-border {
    width: 100%;
    height: 100%;
    color: currentColor;
    pointer-events: none;
}

.language-select-button-loading-border rect {
    width: calc(100% - 4px);
    height: calc(100% - 2px);
    fill: none;
    stroke: currentColor;
    stroke-width: 2px;
    stroke-linecap: round;
    stroke-dasharray: 25 75;
    animation: language-select-button-loading-border 1.5s linear infinite;
}

@keyframes language-select-button-loading-border {
    to {
        stroke-dashoffset: -100;
    }
}

@media (prefers-reduced-motion: reduce) {
    .language-select-button-loading-border rect {
        animation: none;
    }
}
</style>
