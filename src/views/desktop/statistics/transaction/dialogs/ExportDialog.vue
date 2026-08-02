<template>
    <v-dialog width="1000" v-model="showState">
        <one-column-dialog-layout content-class="pa-0"
                                  :title="tt('Export Results')"
                                  :cancel-button-title="tt('Close')"
                                  @cancel="cancel">
            <template #after-title>
                <div ref="buttonContainer">
                    <v-btn density="compact" color="default" variant="text" class="ms-2" :icon="true"
                           :disabled="!exportedData" @click="copy">
                        <v-icon :icon="mdiContentCopy" size="20" />
                        <v-tooltip activator="parent">{{ tt('Copy') }}</v-tooltip>
                    </v-btn>
                    <v-btn density="compact" color="default" variant="text" class="ms-1" :icon="true"
                           @click="save()">
                        <v-icon :icon="mdiContentSaveOutline" size="22" />
                        <v-tooltip activator="parent">{{ tt('Save') }}</v-tooltip>
                    </v-btn>
                </div>
            </template>

            <template #toolbar>
                <toggle-button class="ms-2" :false-name="tt('Table')" :true-name="tt('Raw Data')"
                               v-model="showRawData"/>

                <v-btn density="compact" color="default" variant="text" class="ms-2" :icon="true">
                    <v-icon :icon="mdiDotsVertical" size="22" />
                    <v-menu activator="parent">
                        <v-list>
                            <v-list-subheader :title="tt('File Format')"/>
                            <v-list-item :prepend-icon="mdiComma"
                                         :append-icon="fileFormat === KnownFileType.CSV.extension ? mdiCheck : undefined"
                                         :title="tt('CSV (Comma-separated values) File')"
                                         @click="fileFormat = KnownFileType.CSV.extension"></v-list-item>
                            <v-list-item :prepend-icon="mdiKeyboardTab"
                                         :append-icon="fileFormat === KnownFileType.TSV.extension ? mdiCheck : undefined"
                                         :title="tt('TSV (Tab-separated values) File')"
                                         @click="fileFormat = KnownFileType.TSV.extension"></v-list-item>
                            <v-list-item :prepend-icon="extendMdiSemicolon"
                                         :append-icon="fileFormat === KnownFileType.SSV.extension ? mdiCheck : undefined"
                                         :title="tt('SSV (Semicolon-separated values) File')"
                                         @click="fileFormat = KnownFileType.SSV.extension"></v-list-item>
                            <v-list-item :prepend-icon="mdiLanguageMarkdownOutline"
                                         :append-icon="fileFormat === KnownFileType.MARKDOWN.extension ? mdiCheck : undefined"
                                         :title="tt('Markdown File')"
                                         @click="fileFormat = KnownFileType.MARKDOWN.extension"></v-list-item>
                            <v-list-item :prepend-icon="mdiCodeTags"
                                         :append-icon="fileFormat === KnownFileType.MERMAID.extension && mermaidChartType === ExportMermaidChartType.PieChart ? mdiCheck : undefined"
                                         :title="tt('Mermaid (Pie Chart)')"
                                         @click="fileFormat = KnownFileType.MERMAID.extension; mermaidChartType = ExportMermaidChartType.PieChart"
                                         v-if="supportedMermaidChartTypes[ExportMermaidChartType.PieChart]"></v-list-item>
                            <v-list-item :prepend-icon="mdiCodeTags"
                                         :append-icon="fileFormat === KnownFileType.MERMAID.extension && mermaidChartType === ExportMermaidChartType.XYChartBar ? mdiCheck : undefined"
                                         :title="tt('Mermaid (XY Chart)')"
                                         @click="fileFormat = KnownFileType.MERMAID.extension; mermaidChartType = ExportMermaidChartType.XYChartBar"
                                         v-if="supportedMermaidChartTypes[ExportMermaidChartType.XYChartBar]"></v-list-item>
                            <v-list-item :prepend-icon="mdiCodeTags"
                                         :append-icon="fileFormat === KnownFileType.MERMAID.extension && mermaidChartType === ExportMermaidChartType.XYChartLine ? mdiCheck : undefined"
                                         :title="tt('Mermaid (XY Chart)')"
                                         @click="fileFormat = KnownFileType.MERMAID.extension; mermaidChartType = ExportMermaidChartType.XYChartLine"
                                         v-if="supportedMermaidChartTypes[ExportMermaidChartType.XYChartLine]"></v-list-item>
                        </v-list>
                    </v-menu>
                </v-btn>
            </template>

            <template #content>
                <div class="d-flex flex-column flex-md-row flex-grow-1 overflow-y-auto" style="height: 530px">
                    <v-data-table
                        fixed-header
                        fixed-footer
                        multi-sort
                        density="compact"
                        :headers="dataTableHeaders"
                        :items="dataTableItems"
                        :hover="true"
                        :hide-default-footer="true"
                        :items-per-page="dataTableItems.length"
                        :no-data-text="tt('No data')"
                        v-if="!showRawData"
                    ></v-data-table>
                    <div class="w-100 h-100 code-container" v-if="showRawData">
                        <v-textarea no-resize class="w-100 h-100 ps-4 always-cursor-text"
                                    density="compact" variant="plain"
                                    :readonly="true" :rounded="false"
                                    :value="exportedData"></v-textarea>
                    </div>
                </div>
            </template>
        </one-column-dialog-layout>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';

import { type PartialRecord, itemAndIndex } from '@/core/base.ts';
import type { BigDecimal } from '@/core/numeral.ts';
import { KnownFileType } from '@/core/file.ts';
import { ExportMermaidChartType } from '@/core/statistics.ts';

import { replaceAll, arrayItemToObjectField } from '@/lib/common.ts';
import { BIG_DECIMAL_ZERO, parseBigDecimal } from '@/lib/numeral.ts';
import { copyTextToClipboard, startDownloadFile } from '@/lib/ui/common.ts';
import logger from '@/lib/logger.ts';

import {
    extendMdiSemicolon
} from '@/icons/desktop/extend_mdi_icons.ts';
import {
    mdiDotsVertical,
    mdiCheck,
    mdiComma,
    mdiKeyboardTab,
    mdiLanguageMarkdownOutline,
    mdiCodeTags,
    mdiContentCopy,
    mdiContentSaveOutline
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();

const userStore = useUserStore();

const buttonContainer = useTemplateRef<HTMLElement>('buttonContainer');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const headers = ref<string[]>([]);
const data = ref<string[][]>([]);
const fileFormat = ref<string>(KnownFileType.CSV.extension);
const supportedMermaidChartTypes = ref<PartialRecord<ExportMermaidChartType, boolean>>({});
const mermaidChartType = ref<ExportMermaidChartType | undefined>(undefined);
const showRawData = ref<boolean>(false);

const fileName = computed<string>(() => {
    const nickname = userStore.currentUserNickname;

    if (nickname) {
        return tt('dataExport.exportStatisticsFileName', {
            nickname: nickname
        });
    }

    return tt('dataExport.defaultExportStatisticsFileName');
});

const dataTableHeaders = computed<object[]>(() => {
    return headers.value.map((header, index) => ({
        key: index.toString(),
        value: `column${index}`,
        title: header,
        sortable: index > 0,
        nowrap: true,
        fixed: index === 0
    }));
});

const dataTableItems = computed<object[]>(() => {
    return data.value.map(row => {
        const item: Record<string, string> = {};

        row.forEach((value, index) => {
            item[`column${index}`] = value;
        });

        return item;
    });
});

const exportedData = computed<string>(() => {
    let ret = '';

    if (fileFormat.value === KnownFileType.CSV.extension || fileFormat.value === KnownFileType.TSV.extension || fileFormat.value === KnownFileType.SSV.extension) {
        let separator = ',';

        if (fileFormat.value === KnownFileType.TSV.extension) {
            separator = '\t';
        } else if (fileFormat.value === KnownFileType.SSV.extension) {
            separator = ';';
        }

        if (headers.value.length > 0) {
            ret += headers.value.map(item => replaceAll(item, separator, ' ')).join(separator);
        }

        for (const row of data.value) {
            ret += '\n';
            ret += row.map(item => replaceAll(item, separator, ' ')).join(separator);
        }
    } else if (fileFormat.value === KnownFileType.MARKDOWN.extension) {
        ret += '| ' + headers.value.map(item => replaceAll(item, '|', ' ')).join(' | ') + ' |';
        ret += '\n';
        ret += '| ' + headers.value.map(() => '---').join(' | ') + ' |';

        for (const row of data.value) {
            ret += '\n';
            ret += '| ' + row.map(item => replaceAll(item, '|', ' ')).join(' | ') + ' |';
        }
    } else if (fileFormat.value === KnownFileType.MERMAID.extension && mermaidChartType.value === ExportMermaidChartType.PieChart) {
        ret += 'pie';

        for (const row of data.value) {
            if (row.length > 0) {
                const lengendName: string = replaceAll(row[0] as string, '"', '\\"');
                const value: string = row[1] ?? '0';

                ret += `\n    "${lengendName}" : ${value}`;
            }
        }
    } else if (fileFormat.value === KnownFileType.MERMAID.extension && (mermaidChartType.value === ExportMermaidChartType.XYChartBar || mermaidChartType.value === ExportMermaidChartType.XYChartLine)) {
        ret += 'xychart';
        const lengendNames: string[] = [];
        const xAxisLabels: string[] = [];
        const yAxisValues: Record<number, string[]> = {};
        let minValue: BigDecimal = BIG_DECIMAL_ZERO;
        let maxValue: BigDecimal = BIG_DECIMAL_ZERO;

        for (const [header, index] of itemAndIndex(headers.value)) {
            if (index > 0) {
                lengendNames.push(header);
            }
        }

        for (const row of data.value) {
            for (const [item, index] of itemAndIndex(row)) {
                if (index === 0) {
                    xAxisLabels.push(`"${replaceAll(item, '"', '\\"')}"`);
                } else {
                    let values: string[] | undefined = yAxisValues[index - 1];

                    if (!values) {
                        values = [];
                        yAxisValues[index - 1] = values;
                    }

                    try {
                        const value: BigDecimal = parseBigDecimal(item);

                        if (value.greaterThan(maxValue)) {
                            maxValue = value;
                        } else if (value.lessThan(minValue)) {
                            minValue = value;
                        }
                    } catch (ex) {
                        logger.warn('cannot parse value, original value is ' + item, ex);
                    }

                    values.push(item);
                }
            }
        }

        ret += `\n    x-axis [${xAxisLabels.join(', ')}]`;

        if (minValue.isNegative() || maxValue.isPositive()) {
            ret += `\n    y-axis ${minValue} --> ${maxValue}`;
        }

        for (const [legendName, index] of itemAndIndex(lengendNames)) {
            const values = yAxisValues[index];

            if (values) {
                ret += `\n    %% ${legendName}`;

                if (mermaidChartType.value === ExportMermaidChartType.XYChartBar) {
                    ret += `\n    bar [${values.join(', ')}]`;
                } else if (mermaidChartType.value === ExportMermaidChartType.XYChartLine) {
                    ret += `\n    line [${values.join(', ')}]`;
                }
            }
        }
    }

    return ret;
});

function open(options: { headers: string[], data: string[][], supportedMermaidCharts?: ExportMermaidChartType[] }): void {
    headers.value = options.headers || [];
    data.value = options.data || [];
    fileFormat.value = KnownFileType.CSV.extension;
    supportedMermaidChartTypes.value = arrayItemToObjectField(options.supportedMermaidCharts || [], true);
    mermaidChartType.value = undefined;
    showRawData.value = false;
    showState.value = true;
}

function copy(): void {
    copyTextToClipboard(exportedData.value, buttonContainer.value);
    snackbar.value?.showMessage('Data copied');
}

function save(): void {
    let fileType = KnownFileType.parse(fileFormat.value);

    if (!fileType) {
        fileType = KnownFileType.CSV;
    }

    startDownloadFile(fileType.formatFileName(fileName.value), fileType.createBlob(exportedData.value));
}

function cancel(): void {
    showState.value = false;
}

defineExpose({
    open
});
</script>
