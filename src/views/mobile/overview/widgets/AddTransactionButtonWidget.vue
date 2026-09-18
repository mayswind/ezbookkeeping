<template>
    <f7-segmented round class="add-transaction-button-group margin-horizontal" v-if="hasTransactionAddMenuItems">
        <f7-button round large fill icon-f7="plus" :text="tt('Add Transaction')"
                   @click="emit('navigate', MobileOverviewWidgetNavigationType.Url, '/transaction/add')" />
        <f7-button round large fill class="add-transaction-menu-button" icon-f7="chevron_down"
                   :popover-open="`.add-transaction-button-widget-popover-menu-${widgetId}`"
                   :aria-label="tt('More')"/>
    </f7-segmented>

    <f7-segmented round class="add-transaction-button-group margin-horizontal" v-else-if="!hasTransactionAddMenuItems" >
        <f7-button round large fill icon-f7="plus" :text="tt('Add Transaction')"
                   @click="emit('navigate', MobileOverviewWidgetNavigationType.Url, '/transaction/add')" />
    </f7-segmented>

    <f7-popover class="add-transaction-button-widget-popover-menu" :class="`add-transaction-button-widget-popover-menu-${widgetId}`">
        <f7-list dividers>
            <f7-list-item key="AIClipboardTextRecognition" link="#" no-chevron popover-close
                          :title="tt('AI Clipboard Text Recognition')"
                          @click="emit('navigate', MobileOverviewWidgetNavigationType.AIClipboardTextRecognition)"
                          v-if="isTransactionFromAITextRecognitionEnabled()">
                <template #media>
                    <f7-icon f7="wand_stars"></f7-icon>
                </template>
            </f7-list-item>
            <f7-list-item key="AIImageRecognition" link="#" no-chevron popover-close
                          :title="tt('AI Image Recognition')"
                          @click="emit('navigate', MobileOverviewWidgetNavigationType.AIImageRecognition)"
                          v-if="isTransactionFromAIImageRecognitionEnabled()">
                <template #media>
                    <f7-icon f7="wand_stars"></f7-icon>
                </template>
            </f7-list-item>
            <f7-list-item popover-close link="#" :key="template.id" :title="template.name"
                          @click="emit('navigate', MobileOverviewWidgetNavigationType.Url, '/transaction/add?templateId=' + template.id)"
                          v-for="template in allTransactionTemplates">
                <template #media>
                    <f7-icon f7="doc_plaintext"></f7-icon>
                </template>
            </f7-list-item>
        </f7-list>
    </f7-popover>
</template>

<script setup lang="ts">
import { computed} from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionTemplatesStore } from '@/stores/transactionTemplate.ts';

import { TemplateType } from '@/core/template.ts';
import { MobileOverviewWidgetNavigationType } from '@/core/overview_layout.ts';

import type { TransactionTemplate } from '@/models/transaction_template.ts';

import {
    isTransactionFromAITextRecognitionEnabled,
    isTransactionFromAIImageRecognitionEnabled
} from '@/lib/server_settings.ts';

defineProps<{
    widgetId: string;
}>();

const emit = defineEmits<{
    (e: 'navigate', type: MobileOverviewWidgetNavigationType, path?: string): void;
}>();

const { tt } = useI18n();

const transactionTemplatesStore = useTransactionTemplatesStore();

const allTransactionTemplates = computed<TransactionTemplate[]>(() => {
    const allTemplates = transactionTemplatesStore.allVisibleTemplates;
    return allTemplates[TemplateType.Normal.type] || [];
});

const hasTransactionAddMenuItems = computed<boolean>(() => isTransactionFromAITextRecognitionEnabled() || isTransactionFromAIImageRecognitionEnabled() || (allTransactionTemplates.value && allTransactionTemplates.value.length > 0));
</script>

<style>
.add-transaction-button-group .add-transaction-menu-button {
    flex: 0 0 var(--f7-button-large-height);
    width: var(--f7-button-large-height);
    border-left: 1px solid var(--f7-segmented-raised-divider-color);

    > .icon {
        padding-inline-end: 2px;
    }
}

.add-transaction-button-widget-popover-menu .popover-inner {
    max-height: 400px;
    overflow-y: auto;
}
</style>
