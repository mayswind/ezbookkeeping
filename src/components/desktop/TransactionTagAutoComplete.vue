<template>
    <v-autocomplete
        item-title="name"
        item-value="id"
        auto-select-first
        persistent-placeholder
        multiple
        chips
        :density="density"
        :variant="variant"
        :closable-chips="!readonly"
        :readonly="readonly"
        :disabled="disabled"
        :label="showLabel ? tt('Tags') : undefined"
        :placeholder="tt('None')"
        :items="allTagsWithGroupHeader"
        :model-value="modelValue"
        v-model:search="tagSearchContent"
        @update:modelValue="updateModelValue"
    >
        <template #chip="{ props, item }">
            <v-chip :prepend-icon="mdiPound" :text="item.title" v-bind="props"/>
        </template>

        <template #subheader="{ props }">
            <v-list-subheader>{{ props['title'] }}</v-list-subheader>
        </template>

        <template #item="{ props, item }">
            <v-list-item :value="item.value" v-bind="props" v-if="item.raw instanceof TransactionTag && !item.raw.hidden">
                <template #title>
                    <v-list-item-title>
                        <div class="d-flex align-center">
                            <v-icon size="20" start :icon="mdiPound"/>
                            <span>{{ item.title }}</span>
                        </div>
                    </v-list-item-title>
                </template>
            </v-list-item>
            <v-list-item :disabled="true" v-bind="props" v-else-if="item.raw instanceof TransactionTag && item.raw.hidden">
                <template #title>
                    <v-list-item-title>
                        <div class="d-flex align-center">
                            <v-icon size="20" start :icon="mdiPound"/>
                            <span>{{ item.title }}</span>
                        </div>
                    </v-list-item-title>
                </template>
            </v-list-item>
        </template>

        <template #no-data>
            <v-list class="py-0">
                <v-list-item v-if="tagSearchContent && allowAddNewTag" @click="saveNewTag(tagSearchContent)">{{ tt('format.misc.addNewTag', { tag: tagSearchContent }) }}</v-list-item>
                <v-list-item v-else-if="!tagSearchContent || !allowAddNewTag">{{ tt('No available tag') }}</v-list-item>
            </v-list>
        </template>
    </v-autocomplete>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { type CommonTransactionTagSelectionProps, useTransactionTagSelectionBase } from '@/components/base/TransactionTagSelectionBase.ts';

import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import { TransactionTag } from '@/models/transaction_tag.ts';

import type { ComponentDensity, InputVariant } from '@/lib/ui/desktop.ts';

import {
    mdiPound
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

interface DesktopTransactionTagSelectionProps extends CommonTransactionTagSelectionProps {
    density?: ComponentDensity;
    variant?: InputVariant;
    readonly?: boolean;
    disabled?: boolean;
    showLabel?: boolean;
}

const props = defineProps<DesktopTransactionTagSelectionProps>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string[]): void;
    (e: 'tag:saving', state: boolean, tagName: string): void;
}>();

const { tt } = useI18n();

const {
    tagSearchContent,
    allTagsWithGroupHeader
} = useTransactionTagSelectionBase(props, false);

const transactionTagsStore = useTransactionTagsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

function saveNewTag(tagName: string): void {
    emit('tag:saving', true, tagName);

    transactionTagsStore.saveTag({
        tag: TransactionTag.createNewTag(tagName)
    }).then(tag => {
        emit('tag:saving', false, tagName);

        if (tag && tag.id) {
            const newValue: string[] = [...props.modelValue];
            newValue.push(tag.id);
            updateModelValue(newValue);
        }
    }).catch(error => {
        emit('tag:saving', false, tagName);

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function updateModelValue(newValue: string[]) {
    emit('update:modelValue', newValue);
}
</script>
