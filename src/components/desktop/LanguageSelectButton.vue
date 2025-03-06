<template>
    <v-menu location="bottom" max-height="500">
        <template #activator="{ props }">
            <v-btn variant="text" :disabled="disabled" v-bind="props">{{ currentLanguageName }}</v-btn>
        </template>
        <v-list>
            <v-list-item :key="lang.languageTag" :value="lang.languageTag" v-for="lang in allLanguages">
                <v-list-item-title class="cursor-pointer" @click="updateLanguage(lang.languageTag)">
                    <div class="d-flex align-center">
                        <span>{{ lang.nativeDisplayName }}</span>
                        <v-spacer style="min-width: 40px" />
                        <v-icon :icon="mdiCheck" v-if="isLanguageSelected(lang.languageTag)" />
                        <span class="text-field-append-text" v-if="!isLanguageSelected(lang.languageTag)">{{ lang.displayName }}</span>
                    </div>
                </v-list-item-title>
            </v-list-item>
        </v-list>
    </v-menu>
</template>

<script setup lang="ts">
import { type LanguageSelectButtonBaseProps, type LanguageSelectButtonBaseEmits, useLanguageSelectButtonBase } from '@/components/base/LanguageSelectButtonBase.ts';

import {
    mdiCheck
} from '@mdi/js';

const props = defineProps<LanguageSelectButtonBaseProps>();
const emit = defineEmits<LanguageSelectButtonBaseEmits>();

const {
    allLanguages,
    currentLanguageName,
    updateLanguage,
    isLanguageSelected
} = useLanguageSelectButtonBase(props, emit);
</script>
