<template>
    <component
        :is="componentType"
        item-title="name"
        item-value="id"
        auto-select-first
        clear-on-select
        persistent-placeholder
        multiple
        :class="{ 'transaction-tag-auto-complete--select-only': !!selectOnly }"
        :chips="!selectOnly"
        :density="density"
        :variant="variant"
        :closable-chips="!readonly"
        :readonly="readonly"
        :disabled="disabled"
        :label="showLabel ? tt('Tags') : undefined"
        :placeholder="selectOnly ? tt('Tags') : tt('None')"
        :items="allTagsWithGroupHeader"
        :custom-filter="filterTag"
        :model-value="modelValue"
        v-model:search="tagSearchContent"
        v-model:menu="menuOpenState"
        @keydown.tab.capture="onTabKeyDown"
        @update:modelValue="updateModelValue"
    >
        <template #selection v-if="customSelectionText">
            <!-- hide the default selection when customSelectionText is provided -->
        </template>

        <template #prepend v-if="customSelectionText">
            <div class="d-flex align-center text-truncate cursor-pointer" @mousedown="menuOpenState = !menuOpenState">
                <span class="text-truncate" v-if="customSelectionText">{{ customSelectionText }}</span>
            </div>
        </template>

        <template #chip="{ props, internalItem }" v-if="!selectOnly">
            <v-chip :prepend-icon="mdiPound" :text="internalItem.title" v-bind="props"/>
        </template>

        <template #subheader="{ props }">
            <v-list-subheader class="text-body-small">{{ props['title'] }}</v-list-subheader>
        </template>

        <template #item="{ props, internalItem }">
            <v-list-item :value="internalItem.value" v-bind="props" v-if="internalItem.raw instanceof TransactionTag && !internalItem.raw.hidden">
                <template #title>
                    <v-list-item-title>
                        <div class="d-flex align-center">
                            <v-icon size="18" start :icon="mdiPound"/>
                            <span class="ms-n1">{{ internalItem.title }}</span>
                        </div>
                    </v-list-item-title>
                </template>
            </v-list-item>
            <v-list-item :disabled="true" v-bind="props" v-else-if="internalItem.raw instanceof TransactionTag && internalItem.raw.hidden">
                <template #title>
                    <v-list-item-title>
                        <div class="d-flex align-center">
                            <v-icon size="18" start :icon="mdiPound"/>
                            <span class="ms-n1">{{ internalItem.title }}</span>
                        </div>
                    </v-list-item-title>
                </template>
            </v-list-item>
        </template>

        <template #no-data>
            <v-list class="py-0">
                <v-list-item class="text-body-medium" v-if="tagSearchContent && allowAddNewTag" @click="saveNewTag(tagSearchContent)">{{ tt('format.misc.addNewTag', { tag: tagSearchContent }) }}</v-list-item>
                <v-list-item class="text-body-medium" v-else-if="!tagSearchContent || !allowAddNewTag">{{ tt('No available tag') }}</v-list-item>
            </v-list>
        </template>
    </component>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import { VAutocomplete } from 'vuetify/components/VAutocomplete';
import { VSelect } from "vuetify/components/VSelect";
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import {
    type TransactionTagWithGroupHeader,
    type CommonTransactionTagSelectionProps,
    useTransactionTagSelectionBase
} from '@/components/base/TransactionTagSelectionBase.ts';

import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import { NormalizedText } from '@/core/text.ts';
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
    selectOnly?: boolean;
    customSelectionText?: string;
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

const menuOpenState = ref<boolean>(false);
const componentType = computed<typeof VAutocomplete | typeof VSelect>(() => props.selectOnly ? VSelect : VAutocomplete);

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

function filterTag(value: string, query: string, item?: { value: unknown, raw: TransactionTagWithGroupHeader }): boolean {
    if (!item) {
        return false;
    }

    if ('type' in item.raw) {
        return true;
    }

    const normalizedFilterContent = NormalizedText.normalizeForSearch(query || '');

    if (!normalizedFilterContent) {
        return true;
    }

    return NormalizedText.normalizeForSearch(item.raw.name).indexOf(normalizedFilterContent) >= 0;
}

function onTabKeyDown(event: KeyboardEvent): void {
    const normalizedSearchContent = NormalizedText.normalizeForSearch(tagSearchContent.value);

    if (!normalizedSearchContent) {
        return;
    }

    const firstMatchedTag = allTagsWithGroupHeader.value.find(item => item instanceof TransactionTag && NormalizedText.normalizeForSearch(item.name).indexOf(normalizedSearchContent) >= 0);

    if (firstMatchedTag instanceof TransactionTag && !props.modelValue.includes(firstMatchedTag.id)) {
        event.preventDefault();
    }
}

function updateModelValue(newValue: string[]) {
    emit('update:modelValue', newValue);
}
</script>

<style>
.transaction-tag-auto-complete--select-only.v-input {
    .v-input__prepend {
        margin-inline-end: 0;
    }
}
</style>
