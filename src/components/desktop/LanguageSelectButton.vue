<template>
    <v-menu location="bottom" max-height="500" @update:model-value="onMenuStateChanged">
        <template #activator="{ props }">
            <v-btn variant="text" :disabled="disabled" v-bind="props">{{ currentLanguageName }}</v-btn>
        </template>
        <v-list ref="languageMenu">
            <v-list-item :key="lang.languageTag" :value="lang.languageTag"
                         :class="{ 'list-item-selected': isLanguageSelected(lang.languageTag) }"
                         v-for="lang in allLanguages">
                <v-list-item-title class="cursor-pointer" @click="updateLanguage(lang.languageTag)">
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
    allLanguages,
    currentLanguageName,
    updateLanguage,
    isLanguageSelected
} = useLanguageSelectButtonBase(props, emit);

const languageMenu = useTemplateRef<VList>('languageMenu');

function onMenuStateChanged(state: boolean): void {
    if (state) {
        nextTick(() => {
            scrollToSelectedItem(languageMenu.value?.$el, null, null, 'div.v-list-item.list-item-selected');
        });
    }
}
</script>
