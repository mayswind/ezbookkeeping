<template>
    <v-data-table
        fixed-header
        fixed-footer
        density="compact"
        item-value="index"
        :class="{ 'import-transaction-table': true, 'disabled': !!disabled }"
        :height="parsedFileLinesTableHeight"
        :disable-sort="true"
        :headers="parsedFileLinesHeaders"
        :items="parsedFileLines"
        :no-data-text="tt('No data to import')"
        v-model:items-per-page="countPerPage"
        v-model:page="currentPage"
    >
        <template #headers="{ columns }">
            <tr>
                <th class="text-no-wrap" :key="column.key ?? undefined" v-for="column in columns">
                    <span v-if="!column.key || column.key === 'index'">{{ column.title }}</span>
                    <div class="py-1" v-if="column.key && column.key !== 'index'">
                        <span>{{ getParseDataMappedColumnDisplayName(parseInt(column.key)) }}</span>
                        <br/>
                        <span>({{ column.title }})</span>
                        <v-menu activator="parent" location="bottom" max-height="500">
                            <v-list>
                                <v-list-item :key="columnType.type"
                                             :append-icon="parsedFileDataColumnMapping.dataColumnMapping[columnType.type] === parseInt(column.key) ? mdiCheck : undefined"
                                             v-for="columnType in allImportTransactionColumnTypes"
                                             @click="parsedFileDataColumnMapping.toggleDataMappingColumn(parseInt(column.key), columnType.type)">
                                    <v-list-item-title class="cursor-pointer">
                                        {{ columnType.displayName }}
                                    </v-list-item-title>
                                </v-list-item>
                            </v-list>
                        </v-menu>
                    </div>
                </th>
            </tr>
        </template>
        <template #bottom>
            <div class="title-and-toolbar d-flex align-center text-no-wrap mt-2" v-if="parsedFileData">
                <v-btn color="secondary" density="compact" variant="outlined"
                       :append-icon="parsedFileDataColumnMapping.includeHeader ? mdiCheck : mdiClose"
                       @click="parsedFileDataColumnMapping.toggleIncludeHeader()">{{ tt('Include Header Line') }}</v-btn>
                <v-btn class="ms-2" color="secondary" density="compact" variant="outlined"
                       :disabled="!parsedFileDataColumnMapping || !parsedFileDataColumnMapping.isColumnMappingSet(ImportTransactionColumnType.TransactionType) || !parsedFileAllTransactionTypes">
                    <span>{{ tt('Transaction Type Mapping') }}</span>
                    <span class="ms-1" v-if="parsedFileDataColumnMapping && parsedFileDataColumnMapping.isColumnMappingSet(ImportTransactionColumnType.TransactionType) && parsedFileAllTransactionTypes">({{ getObjectOwnFieldCount(parsedFileValidMappedTransactionTypes) || tt('None') }})</span>
                    <v-menu eager activator="parent" location="bottom" max-height="500"
                            :close-on-content-click="false">
                        <v-list class="pa-0">
                            <v-list-item class="pa-0">
                                <v-table class="transaction-types-popup-menu">
                                    <tbody>
                                    <tr :key="typeName"
                                        v-for="typeName in parsedFileAllTransactionTypes">
                                        <td>{{ typeName }}</td>
                                        <td>
                                            <v-btn-toggle class="transaction-types-toggle" density="compact" variant="outlined"
                                                          mandatory="force" divided
                                                          v-model="parsedFileDataColumnMapping.transactionTypeMapping[typeName]">
                                                <v-btn :value="undefined">{{ tt('None') }}</v-btn>
                                                <v-btn :value="TransactionType.ModifyBalance">{{ tt('Modify Balance') }}</v-btn>
                                                <v-btn :value="TransactionType.Income">{{ tt('Income') }}</v-btn>
                                                <v-btn :value="TransactionType.Expense">{{ tt('Expense') }}</v-btn>
                                                <v-btn :value="TransactionType.Transfer">{{ tt('Transfer') }}</v-btn>
                                            </v-btn-toggle>
                                        </td>
                                    </tr>
                                    </tbody>
                                </v-table>
                            </v-list-item>
                        </v-list>
                    </v-menu>
                </v-btn>
                <v-btn class="ms-2" color="secondary" density="compact" variant="outlined"
                       :disabled="!parsedFileDataColumnMapping || !parsedFileDataColumnMapping.isColumnMappingSet(ImportTransactionColumnType.TransactionTime)">
                    <span>{{ tt('Time Format') }}</span>
                    <span class="ms-1" v-if="parsedFileDataColumnMapping && parsedFileDataColumnMapping.isColumnMappingSet(ImportTransactionColumnType.TransactionTime)">({{ parsedFileDataColumnMapping.timeFormat || parsedFileAutoDetectedTimeFormat || tt('Unknown') }})</span>
                    <v-menu eager activator="parent" location="bottom" max-height="500">
                        <v-list>
                            <v-list-item key="auto"
                                         :append-icon="parsedFileDataColumnMapping.timeFormat === '' ? mdiCheck : undefined"
                                         @click="parsedFileDataColumnMapping.timeFormat = ''">
                                <v-list-item-title class="cursor-pointer">
                                    <span>{{ tt('Auto detect') }}</span>
                                    <span class="ms-1" v-if="parsedFileAutoDetectedTimeFormat">({{ parsedFileAutoDetectedTimeFormat }})</span>
                                    <span class="ms-1" v-if="!parsedFileAutoDetectedTimeFormat">({{ tt('Unknown') }})</span>
                                </v-list-item-title>
                            </v-list-item>
                            <v-list-item :key="dateTimeFormat.format"
                                         :append-icon="parsedFileDataColumnMapping.timeFormat === dateTimeFormat.format ? mdiCheck : undefined"
                                         v-for="dateTimeFormat in KnownDateTimeFormat.values()"
                                         @click="parsedFileDataColumnMapping.timeFormat = dateTimeFormat.format">
                                <v-list-item-title class="cursor-pointer">
                                    {{ dateTimeFormat.format }}
                                </v-list-item-title>
                            </v-list-item>
                        </v-list>
                    </v-menu>
                </v-btn>
                <v-btn class="ms-2" color="secondary" density="compact" variant="outlined"
                       v-if="parsedFileDataColumnMapping && parsedFileDataColumnMapping.isColumnMappingSet(ImportTransactionColumnType.TransactionTimezone)">
                    <span>{{ tt('Timezone Format') }}</span>
                    <span class="ms-1" v-if="parsedFileDataColumnMapping && parsedFileDataColumnMapping.isColumnMappingSet(ImportTransactionColumnType.TransactionTimezone)">({{ KnownDateTimezoneFormat.valueOf(parsedFileDataColumnMapping.timezoneFormat || parsedFileAutoDetectedTimezoneFormat || '')?.name || tt('Unknown') }})</span>
                    <v-menu eager activator="parent" location="bottom" max-height="500">
                        <v-list>
                            <v-list-item key="auto"
                                         :append-icon="parsedFileDataColumnMapping.timezoneFormat === '' ? mdiCheck : undefined"
                                         @click="parsedFileDataColumnMapping.timezoneFormat = ''">
                                <v-list-item-title class="cursor-pointer">
                                    <span>{{ tt('Auto detect') }}</span>
                                    <span class="ms-1" v-if="parsedFileAutoDetectedTimezoneFormat && KnownDateTimezoneFormat.valueOf(parsedFileAutoDetectedTimezoneFormat || '')">({{ KnownDateTimezoneFormat.valueOf(parsedFileAutoDetectedTimezoneFormat || '')?.name }})</span>
                                    <span class="ms-1" v-if="!parsedFileAutoDetectedTimezoneFormat || !KnownDateTimezoneFormat.valueOf(parsedFileAutoDetectedTimezoneFormat || '')">({{ tt('Unknown') }})</span>
                                </v-list-item-title>
                            </v-list-item>
                            <v-list-item :key="timezoneFormat.value"
                                         :append-icon="parsedFileDataColumnMapping.timezoneFormat === timezoneFormat.value ? mdiCheck : undefined"
                                         v-for="timezoneFormat in KnownDateTimezoneFormat.values()"
                                         @click="parsedFileDataColumnMapping.timezoneFormat = timezoneFormat.value">
                                <v-list-item-title class="cursor-pointer">
                                    {{ timezoneFormat.name }}
                                </v-list-item-title>
                            </v-list-item>
                        </v-list>
                    </v-menu>
                </v-btn>
                <v-btn class="ms-2" color="secondary" density="compact" variant="outlined"
                       :disabled="!parsedFileDataColumnMapping || !parsedFileDataColumnMapping.isColumnMappingSet(ImportTransactionColumnType.Amount)">
                    <span>{{ tt('Amount Format') }}</span>
                    <span class="ms-1" v-if="parsedFileDataColumnMapping && parsedFileDataColumnMapping.isColumnMappingSet(ImportTransactionColumnType.Amount)">({{ KnownAmountFormat.valueOf(parsedFileDataColumnMapping.amountFormat || parsedFileAutoDetectedAmountFormat || '')?.format || tt('Unknown') }})</span>
                    <v-menu eager activator="parent" location="bottom" max-height="500">
                        <v-list>
                            <v-list-item key="auto"
                                         :append-icon="parsedFileDataColumnMapping.amountFormat === '' ? mdiCheck : undefined"
                                         @click="parsedFileDataColumnMapping.amountFormat = ''">
                                <v-list-item-title class="cursor-pointer">
                                    <span>{{ tt('Auto detect') }}</span>
                                    <span class="ms-1" v-if="parsedFileAutoDetectedAmountFormat && KnownAmountFormat.valueOf(parsedFileAutoDetectedAmountFormat || '')">({{ KnownAmountFormat.valueOf(parsedFileAutoDetectedAmountFormat || '')?.format }})</span>
                                    <span class="ms-1" v-if="!parsedFileAutoDetectedAmountFormat || !KnownAmountFormat.valueOf(parsedFileAutoDetectedAmountFormat || '')">({{ tt('Unknown') }})</span>
                                </v-list-item-title>
                            </v-list-item>
                            <v-list-item :key="amountFormat.type"
                                         :append-icon="parsedFileDataColumnMapping.amountFormat === amountFormat.type ? mdiCheck : undefined"
                                         v-for="amountFormat in KnownAmountFormat.values()"
                                         @click="parsedFileDataColumnMapping.amountFormat = amountFormat.type">
                                <v-list-item-title class="cursor-pointer">
                                    {{ amountFormat.format }}
                                </v-list-item-title>
                            </v-list-item>
                        </v-list>
                    </v-menu>
                </v-btn>
                <v-btn class="ms-2" color="secondary" density="compact" variant="outlined"
                       v-if="parsedFileDataColumnMapping && parsedFileDataColumnMapping.isColumnMappingSet(ImportTransactionColumnType.GeographicLocation)">
                    <span>{{ tt('Geographic Location Separator') }}</span>
                    <span class="ms-1" v-if="parsedFileDataColumnMapping.geoLocationOrder">({{ parsedFileDataColumnMapping.formatGeoLocation(tt('Latitude'), tt('Longitude')) }})</span>
                    <v-menu eager activator="parent" location="bottom" max-height="500"
                            :close-on-content-click="false">
                        <v-list class="pa-0">
                            <v-list-item class="pa-0">
                                <v-table class="transaction-types-popup-menu">
                                    <tbody>
                                    <tr :key="separator.value"
                                        v-for="separator in allSeparators">
                                        <td>{{ separator.name }} ({{separator.value}})</td>
                                        <td>
                                            <v-btn-toggle class="transaction-types-toggle" density="compact" variant="outlined"
                                                          mandatory="force" divided
                                                          v-model="parsedFileDataColumnMapping.geoLocationOrder"
                                                          v-if="parsedFileDataColumnMapping.geoLocationSeparator === separator.value">
                                                <v-btn value="latlon">{{ `${tt('Latitude')}${separator.value}${tt('Longitude')}` }}</v-btn>
                                                <v-btn value="lonlat">{{ `${tt('Longitude')}${separator.value}${tt('Latitude')}` }}</v-btn>
                                            </v-btn-toggle>
                                            <v-btn-group class="transaction-types-toggle" density="compact" variant="outlined"
                                                         divided v-if="parsedFileDataColumnMapping.geoLocationSeparator !== separator.value">
                                                <v-btn @click="parsedFileDataColumnMapping.setGeoLocationFormat(separator.value, 'latlon')">{{ `${tt('Latitude')}${separator.value}${tt('Longitude')}` }}</v-btn>
                                                <v-btn @click="parsedFileDataColumnMapping.setGeoLocationFormat(separator.value, 'lonlat')">{{ `${tt('Longitude')}${separator.value}${tt('Latitude')}` }}</v-btn>
                                            </v-btn-group>
                                        </td>
                                    </tr>
                                    </tbody>
                                </v-table>
                            </v-list-item>
                        </v-list>
                    </v-menu>
                </v-btn>
                <v-btn class="ms-2" color="secondary" density="compact" variant="outlined"
                       v-if="parsedFileDataColumnMapping && parsedFileDataColumnMapping.isColumnMappingSet(ImportTransactionColumnType.Tags)">
                    <span>{{ tt('Transaction Tags Separator') }}</span>
                    <span class="ms-1" v-if="parsedFileDataColumnMapping.tagSeparator">({{ parsedFileDataColumnMapping.tagSeparator }})</span>
                    <v-menu eager activator="parent" location="bottom" max-height="500">
                        <v-list>
                            <v-list-item :key="separator.value"
                                         :append-icon="parsedFileDataColumnMapping.tagSeparator === separator.value ? mdiCheck : undefined"
                                         v-for="separator in allSeparators"
                                         @click="parsedFileDataColumnMapping.tagSeparator = separator.value">
                                <v-list-item-title class="cursor-pointer">
                                    {{ separator.name }} ({{separator.value}})
                                </v-list-item-title>
                            </v-list-item>
                        </v-list>
                    </v-menu>
                </v-btn>
                <v-spacer/>
                <span>{{ tt('Lines Per Page') }}</span>
                <v-select class="ms-2" density="compact" max-width="100"
                          item-title="name"
                          item-value="value"
                          :disabled="!!disabled"
                          :items="parsedFileLinesTablePageOptions"
                          v-model="countPerPage"
                />
                <pagination-buttons density="compact"
                                    :disabled="!!disabled"
                                    :totalPageCount="Math.ceil((parsedFileLines ? parsedFileLines.length : 0) / countPerPage)"
                                    v-model="currentPage"></pagination-buttons>
            </div>
        </template>
    </v-data-table>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';
import PaginationButtons from '@/components/desktop/PaginationButtons.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { type NameValue, type NameNumeralValue, type TypeAndDisplayName, itemAndIndex, entries } from '@/core/base.ts';
import { type NumeralSystem, KnownAmountFormat } from '@/core/numeral.ts';
import { KnownDateTimeFormat } from '@/core/datetime.ts';
import { KnownDateTimezoneFormat } from '@/core/timezone.ts';
import { TransactionType } from '@/core/transaction.ts';
import { ImportTransactionColumnType, ImportTransactionDataMapping } from '@/core/import_transaction.ts';
import { KnownFileType } from '@/core/file.ts';

import {
    isNumber,
    isObjectEmpty,
    getObjectOwnFieldCount,
    findDisplayNameByType
} from '@/lib/common.ts';
import {
    openTextFileContent,
    startDownloadFile
} from '@/lib/ui/common.ts';
import logger from '@/lib/logger.ts';

import {
    mdiCheck,
    mdiClose,
    mdiFolderOpenOutline,
    mdiContentSaveOutline
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

interface ImportTransactionDefineColumnResult {
    includeHeader: boolean;
    columnMapping: Record<number, number>;
    transactionTypeMapping: Record<string, TransactionType>;
    timeFormat: string | undefined;
    timezoneFormat: string | undefined;
    amountDecimalSeparator: string | undefined;
    amountDigitGroupingSymbol: string | undefined;
    geoLocationSeparator: string;
    geoLocationOrder: string;
    tagSeparator: string;
}

interface ImportTransactionDefineColumnMenu {
    prependIcon: string;
    title: string;
    disabled?: boolean;
    onClick: () => void;
}

const props = defineProps<{
    parsedFileData?: string[][];
    disabled?: boolean;
}>();

const {
    tt,
    getCurrentNumeralSystemType,
    getAllImportTransactionColumnTypes
} = useI18n();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const currentPage = ref<number>(1);
const countPerPage = ref<number>(10);
const parsedFileDataColumnMapping = ref<ImportTransactionDataMapping>(ImportTransactionDataMapping.createEmpty());

const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());
const allImportTransactionColumnTypes = computed<TypeAndDisplayName[]>(() => getAllImportTransactionColumnTypes());

const menus = computed<ImportTransactionDefineColumnMenu[]>(() => [
    {
        prependIcon: mdiFolderOpenOutline,
        title: tt('Load Data Mapping File'),
        onClick: loadColumnMappingFile
    },
    {
        prependIcon: mdiContentSaveOutline,
        title: tt('Save Data Mapping File'),
        onClick: saveColumnMappingFile
    }
]);

const allSeparators = computed<NameValue[]>(() => {
    const separators: NameValue[] = [
        {
            name: tt('Space'),
            value: ' '
        },
        {
            name: tt('Comma'),
            value: ','
        },
        {
            name: tt('Semicolon'),
            value: ';'
        },
        {
            name: tt('Tab'),
            value: '\t'
        },
        {
            name: tt('Vertical Bar'),
            value: '|'
        }
    ];

    return separators;
});

const parsedFileLinesTableHeight = computed<number | undefined>(() => {
    if (countPerPage.value <= 10 || !parsedFileLines.value || parsedFileLines.value.length <= 10) {
        return undefined;
    } else {
        return 400;
    }
});

const parsedFileLinesHeaders = computed<object[]>(() => {
    let maxColumnCount = 0;

    if (props.parsedFileData) {
        for (const lineData of props.parsedFileData) {
            if (lineData.length > maxColumnCount) {
                maxColumnCount = lineData.length;
            }
        }
    }

    const firstLine: string[] = props.parsedFileData && props.parsedFileData.length > 0 ? (props.parsedFileData[0] as string[]) : [];

    const headers: object[] = [];
    headers.push({ key: 'index', value: 'index', title: '#', sortable: true, nowrap: true });

    for (let i = 0; i < maxColumnCount; i++) {
        let title = `#${i + 1}`;

        if (parsedFileDataColumnMapping.value.includeHeader && firstLine && firstLine[i]) {
            title = firstLine[i] as string;
        }

        headers.push({ key: i.toString(), value: `column${i + 1}`, title: title, sortable: true, nowrap: true });
    }

    return headers;
});

const parsedFileLines = computed<Record<string, string>[] | undefined>(() => {
    if (!props.parsedFileData) {
        return undefined;
    }

    const allLines: Record<string, string>[] = [];
    const startIndex = parsedFileDataColumnMapping.value.includeHeader ? 1 : 0;

    for (let i = startIndex, index = 1; i < props.parsedFileData.length; i++, index++) {
        const line: Record<string, string> = {};
        const columns = props.parsedFileData[i] as string[];

        for (const [data, columnIndex] of itemAndIndex(columns)) {
            line['index'] = index.toString();
            line[`column${columnIndex + 1}`] = data;
        }

        allLines.push(line);
    }

    return allLines;
});

const parsedFileLinesTablePageOptions = computed<NameNumeralValue[]>(() => getTablePageOptions(parsedFileLines.value?.length));

const parsedFileAllTransactionTypes = computed<string[]>(() => parsedFileDataColumnMapping.value.parseFileAllTransactionTypes(props.parsedFileData));
const parsedFileValidMappedTransactionTypes = computed<Record<string, TransactionType>>(() => parsedFileDataColumnMapping.value.parseFileValidMappedTransactionTypes(props.parsedFileData));
const parsedFileAutoDetectedTimeFormat = computed<string | undefined>(() => parsedFileDataColumnMapping.value.parseFileAutoDetectedTimeFormat(props.parsedFileData));
const parsedFileAutoDetectedTimezoneFormat = computed<string | undefined>(() => parsedFileDataColumnMapping.value.parseFileAutoDetectedTimezoneFormat(props.parsedFileData));
const parsedFileAutoDetectedAmountFormat = computed<string | undefined>(() => parsedFileDataColumnMapping.value.parseFileAutoDetectedAmountFormat(props.parsedFileData));

function getDisplayCount(count: number): string {
    return numeralSystem.value.formatNumber(count);
}

function getTablePageOptions(linesCount?: number): NameNumeralValue[] {
    const pageOptions: NameNumeralValue[] = [];

    if (!linesCount || linesCount < 1) {
        pageOptions.push({ value: -1, name: tt('All') });
        return pageOptions;
    }

    for (const count of [ 5, 10, 15, 20, 25, 30, 50 ]) {
        if (linesCount < count) {
            break;
        }

        pageOptions.push({ value: count, name: getDisplayCount(count) });
    }

    pageOptions.push({ value: -1, name: tt('All') });

    return pageOptions;
}

function getParseDataMappedColumnDisplayName(columnIndex: number): string {
    for (const [columnType, index] of entries(parsedFileDataColumnMapping.value.dataColumnMapping)) {
        if (index === columnIndex) {
            return findDisplayNameByType(allImportTransactionColumnTypes.value, parseInt(columnType)) || tt('Unspecified');
        }
    }

    return tt('Unspecified');
}

function generateResult(): ImportTransactionDefineColumnResult | undefined {
    const columnMapping: Record<number, number> = parsedFileDataColumnMapping.value.dataColumnMapping;
    const transactionTypeMapping: Record<string, TransactionType> = parsedFileValidMappedTransactionTypes.value;
    const includeHeader: boolean = parsedFileDataColumnMapping.value.includeHeader;
    const geoLocationSeparator: string = parsedFileDataColumnMapping.value.geoLocationSeparator;
    const geoLocationOrder: string = parsedFileDataColumnMapping.value.geoLocationOrder;
    const tagSeparator: string = parsedFileDataColumnMapping.value.tagSeparator;

    let timeFormat: string | undefined = parsedFileDataColumnMapping.value.timeFormat;
    let timezoneFormat: string | undefined = parsedFileDataColumnMapping.value.timezoneFormat;
    let amountFormat: string | undefined = parsedFileDataColumnMapping.value.amountFormat;
    let amountDecimalSeparator: string | undefined = undefined;
    let amountDigitGroupingSymbol: string | undefined = undefined;

    if (!columnMapping
        || !isNumber(columnMapping[ImportTransactionColumnType.TransactionTime.type])
        || !isNumber(columnMapping[ImportTransactionColumnType.TransactionType.type])
        || !isNumber(columnMapping[ImportTransactionColumnType.Amount.type])) {
        snackbar.value?.showError('Missing transaction time, transaction type, or amount column mapping');
        return undefined;
    }

    if (!transactionTypeMapping || isObjectEmpty(transactionTypeMapping)) {
        snackbar.value?.showError('Transaction type mapping is not set');
        return undefined;
    }

    if (!parsedFileDataColumnMapping.value.timeFormat) {
        timeFormat = parsedFileAutoDetectedTimeFormat.value;
    }

    if (!parsedFileDataColumnMapping.value.timezoneFormat) {
        timezoneFormat = parsedFileAutoDetectedTimezoneFormat.value;
    }

    if (!parsedFileDataColumnMapping.value.amountFormat) {
        amountFormat = parsedFileAutoDetectedAmountFormat.value;
    }

    if (amountFormat) {
        const knownAmountFormat = KnownAmountFormat.valueOf(amountFormat);

        if (knownAmountFormat) {
            amountDecimalSeparator = knownAmountFormat.decimalSeparator.symbol;
            amountDigitGroupingSymbol = knownAmountFormat.digitGroupingSymbol?.symbol;
        }
    }

    if (!timeFormat) {
        snackbar.value?.showError('Transaction time format is not set');
        return undefined;
    }

    if (!amountDecimalSeparator) {
        snackbar.value?.showError('Transaction amount format is not set');
        return undefined;
    }

    return {
        includeHeader: includeHeader,
        columnMapping: columnMapping,
        transactionTypeMapping: transactionTypeMapping,
        timeFormat: timeFormat,
        timezoneFormat: timezoneFormat,
        amountDecimalSeparator: amountDecimalSeparator,
        amountDigitGroupingSymbol: amountDigitGroupingSymbol,
        geoLocationSeparator: geoLocationSeparator,
        geoLocationOrder: geoLocationOrder,
        tagSeparator: tagSeparator
    };
}

function reset(): void {
    parsedFileDataColumnMapping.value.reset();
    currentPage.value = 1;
    countPerPage.value = 10;
}

function loadColumnMappingFile(): void {
    openTextFileContent({
        allowedExtensions: KnownFileType.JSON.contentType
    }).then(content => {
        const result = ImportTransactionDataMapping.parseFromJson(content);

        if (result) {
            parsedFileDataColumnMapping.value = result;
        } else {
            logger.error('Failed to parse data mapping file');
            snackbar.value?.showError('Data mapping file is invalid');
        }
    }).catch(error => {
        logger.error('Failed to open data mapping file', error);
        snackbar.value?.showError('Data mapping file is invalid');
    });
}

function saveColumnMappingFile(): void {
    const fileName = KnownFileType.JSON.formatFileName(tt('dataExport.defaultImportDataMappingFileName'));
    startDownloadFile(fileName, KnownFileType.JSON.createBlob(parsedFileDataColumnMapping.value.toJson()));
}

defineExpose({
    menus,
    generateResult,
    reset,
    loadColumnMappingFile,
    saveColumnMappingFile
});
</script>

<style>
.transaction-types-popup-menu .transaction-types-toggle {
    overflow-x: auto;
    white-space: nowrap;
}

.transaction-types-popup-menu .transaction-types-toggle.v-btn-toggle {
    height: auto !important;
    padding: 0;
    border: none;
}

.transaction-types-popup-menu .transaction-types-toggle.v-btn-toggle > .v-btn {
    border-color: rgba(var(--v-border-color), var(--v-border-opacity));
}

.transaction-types-popup-menu .transaction-types-toggle.v-btn-toggle > .v-btn:not(:first-child) {
    border-top-left-radius: 0;
    border-bottom-left-radius: 0;
    border-left: none;
}

.transaction-types-popup-menu .transaction-types-toggle.v-btn-toggle > .v-btn:not(:last-child) {
    border-top-right-radius: 0;
    border-bottom-right-radius: 0;
}

.transaction-types-popup-menu .transaction-types-toggle.v-btn-toggle > .v-btn {
    border: thin solid rgba(var(--v-border-color), var(--v-border-opacity));
}

.transaction-types-popup-menu .transaction-types-toggle.v-btn-toggle button.v-btn {
    width: auto !important;
}
</style>
