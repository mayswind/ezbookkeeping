<template>
    <v-col cols="12" md="12" v-if="fileTypeInfo?.subTypes">
        <v-select
            item-title="displayName"
            item-value="type"
            :disabled="disabled || disabledSubType"
            :label="tt('Format')"
            :placeholder="tt('Format')"
            :items="fileTypeInfo?.subTypes"
            v-model="fileSubType"
        />
    </v-col>

    <v-col cols="12" md="12" v-if="!fileTypeInfo?.dataFromTextbox && supportedEncodings">
        <v-select
            item-title="displayName"
            item-value="encoding"
            :disabled="disabled"
            :label="tt('File Encoding')"
            :placeholder="tt('File Encoding')"
            :items="supportedEncodings"
            v-model="fileEncoding"
        />
    </v-col>

    <v-col cols="12" md="12" v-if="isCustomFileFormat">
        <v-select
            item-title="displayName"
            item-value="type"
            :disabled="disabled || disabledHandlingMethod"
            :label="tt('Handling Method')"
            :placeholder="tt('Handling Method')"
            :items="[
                { displayName: tt('Column Mapping'), type: ImportCustomFileFormatProcessMethod.ColumnMapping },
                { displayName: tt('Custom Script'), type: ImportCustomFileFormatProcessMethod.CustomScript }
             ]"
            v-model="handlingMethod"
        />
    </v-col>

    <v-col cols="12" md="12" v-if="supportedAdditionalOptions">
        <v-select
            :disabled="disabled"
            :label="tt('Additional Options')"
            :placeholder="tt('Additional Options')"
            :model-value="fileTypeInfo?.type"
            v-model:menu="additionalOptionsMenuState"
        >
            <template #selection>
                <span class="cursor-pointer">{{ displaySelectedAdditionalOptions }}</span>
            </template>

            <template #no-data>
                <v-list class="py-0" select-strategy="classic"
                        :selected="objectFieldWithValueToArrayItem(additionalOptions, true as (boolean | undefined))"
                        @keydown.up.stop @keydown.down.stop>
                    <template v-for="item in ALL_IMPORT_FILE_ADDITIONAL_OPTIONS">
                        <v-list-item class="text-body-medium" :key="item.key" :value="item.key"
                                     :append-icon="additionalOptions[item.key] ? mdiCheck : undefined"
                                     @click="additionalOptions[item.key] = !additionalOptions[item.key]"
                                     v-if="isDefined(supportedAdditionalOptions[item.key])">{{ tt(item.name) }}</v-list-item>
                    </template>
                </v-list>
            </template>
        </v-select>
    </v-col>

    <v-col cols="12" md="12" v-if="supportedAIAdditionalPrompt">
        <v-textarea
            type="text"
            autocomplete="off"
            persistent-placeholder
            rows="2"
            :disabled="disabled"
            :placeholder="tt('Additional Prompt')"
            v-model="additionalPrompt"
        />
    </v-col>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import type {
    ImportFileTypeSupportedAdditionalOptions,
    LocalizedImportFileType,
    LocalizedImportFileTypeSupportedEncodings
} from '@/core/import_transaction.ts';
import { ImportCustomFileFormatProcessMethod } from '@/core/import_transaction.ts';
import { ALL_IMPORT_FILE_ADDITIONAL_OPTIONS } from '@/consts/import_transaction.ts';

import { isDefined, objectFieldWithValueToArrayItem } from '@/lib/common.ts';

import {
    mdiCheck
} from '@mdi/js';

const props = defineProps<{
    fileType: string;
    fileSubType: string;
    fileEncoding: string;
    handlingMethod: ImportCustomFileFormatProcessMethod;
    additionalOptions: ImportFileTypeSupportedAdditionalOptions;
    additionalPrompt: string;
    fileTypeInfo?: LocalizedImportFileType;
    supportedEncodings?: LocalizedImportFileTypeSupportedEncodings[];
    disabled?: boolean;
    disabledSubType?: boolean;
    disabledHandlingMethod?: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:fileType', value: string): void;
    (e: 'update:fileSubType', value: string): void;
    (e: 'update:fileEncoding', value: string): void;
    (e: 'update:handlingMethod', value: ImportCustomFileFormatProcessMethod): void;
    (e: 'update:additionalOptions', value: ImportFileTypeSupportedAdditionalOptions): void;
    (e: 'update:additionalPrompt', value: string): void;
}>();

const { tt, joinMultiText } = useI18n();

const additionalOptionsMenuState = ref<boolean>(false);

const isCustomFileFormat = computed<boolean>(() => props.fileTypeInfo?.type === 'dsv' || props.fileTypeInfo?.type === 'dsv_data' || props.fileTypeInfo?.type === 'excel');
const supportedAdditionalOptions = computed(() => props.fileTypeInfo?.supportedAdditionalOptions);
const supportedAIAdditionalPrompt = computed<boolean>(() => !!props.fileTypeInfo?.supportedAIAdditionalPrompt);

const fileSubType = computed<string>({
    get: () => props.fileSubType,
    set: value => emit('update:fileSubType', value)
});

const fileEncoding = computed<string>({
    get: () => props.fileEncoding,
    set: value => emit('update:fileEncoding', value)
});

const handlingMethod = computed<ImportCustomFileFormatProcessMethod>({
    get: () => props.handlingMethod,
    set: value => emit('update:handlingMethod', value)
});

const additionalOptions = computed<ImportFileTypeSupportedAdditionalOptions>({
    get: () => props.additionalOptions,
    set: value => emit('update:additionalOptions', value)
});

const additionalPrompt = computed<string>({
    get: () => props.additionalPrompt,
    set: value => emit('update:additionalPrompt', value)
});

const displaySelectedAdditionalOptions = computed<string>(() => {
    if (!supportedAdditionalOptions.value) {
        return tt('None');
    }

    const selectedOptions: string[] = [];

    for (const option of ALL_IMPORT_FILE_ADDITIONAL_OPTIONS) {
        if (isDefined(supportedAdditionalOptions.value[option.key]) && additionalOptions.value[option.key]) {
            selectedOptions.push(tt(option.name));
        }
    }

    if (selectedOptions.length < 1) {
        return tt('None');
    }

    return joinMultiText(selectedOptions);
});
</script>
