<template>
    <code-editor class="w-100" :style="{ height: height || '360px' }" language="javascript"
                 :readonly="readonly" :rounded="true" :line-numbers="true"
                 :extra-libs="customScriptEditorExtraLibs" v-model="script" />
</template>

<script setup lang="ts">
import type { CodeEditorExtraLib } from '@/components/desktop/CodeEditor.vue';

import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

const props = defineProps<{
    modelValue: string;
    readonly?: boolean;
    height?: string;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
}>();

const { tt } = useI18n();

const script = computed<string>({
    get: () => props.modelValue,
    set: value => emit('update:modelValue', value)
});

const customScriptEditorExtraLibs = computed<CodeEditorExtraLib[]>(() => [{
    filePath: 'inmemory://model/ezbookkeeping-import-transaction.d.ts',
    content: `
interface ParsedDateTime {
    dateTime: string;
    format: string;
}

interface ParsedUtcOffset {
    name: string;
}

/** ${tt('sample.importTransactionCustomScript.functionReturnDescription')} */
interface ImportTransaction {
    /** ${tt('sample.importTransactionCustomScript.fieldTimeDescription')} */
    time: string | ParsedDateTime;
    /** ${tt('sample.importTransactionCustomScript.fieldUtcOffsetDescription')} */
    utcOffset: string | ParsedUtcOffset;
    /** ${tt('sample.importTransactionCustomScript.fieldTypeDescription')} */
    type: 'Income' | 'Expense' | 'Transfer';
    /** ${tt('sample.importTransactionCustomScript.fieldCategoryNameDescription')} */
    categoryName?: string;
    /** ${tt('sample.importTransactionCustomScript.fieldSourceAccountNameDescription')} */
    sourceAccountName?: string;
    /** ${tt('sample.importTransactionCustomScript.fieldDestinationAccountNameDescription')} */
    destinationAccountName?: string;
    /** ${tt('sample.importTransactionCustomScript.fieldSourceAmountDescription')} */
    sourceAmount: string;
    /** ${tt('sample.importTransactionCustomScript.fieldDestinationAmountDescription')} */
    destinationAmount?: string;
    /** ${tt('sample.importTransactionCustomScript.fieldGeoLocationDescription')} */
    geoLocation?: string;
    /** ${tt('sample.importTransactionCustomScript.fieldTagNamesDescription')} */
    tagNames?: string;
    /** ${tt('sample.importTransactionCustomScript.fieldCommentDescription')} */
    description?: string;
}

/** ${tt('sample.importTransactionCustomScript.fieldTypeDescription')} */
declare const TransactionType: Readonly<{
    /** ${tt('Income')} */
    Income: 'Income';
    /** ${tt('Expense')} */
    Expense: 'Expense';
    /** ${tt('Transfer')} */
    Transfer: 'Transfer';
}>;

/** ${tt('sample.importTransactionCustomScript.fieldTimeDescription')} */
declare function parseDateTime(dateTime: string, format: string): ParsedDateTime;

/** ${tt('sample.importTransactionCustomScript.fieldUtcOffsetDescription')} */
declare function parseUtcOffset(timezoneName: string): ParsedUtcOffset;`
}]);

</script>
