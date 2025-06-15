<template>
    <v-dialog width="1000" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-4">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <div class="d-flex w-100 align-center justify-center">
                        <h4 class="text-h4">{{ tt('Export Results') }}</h4>
                    </div>
                    <v-btn density="comfortable" color="default" variant="text" class="ml-2" :icon="true">
                        <v-icon :icon="mdiDotsVertical" />
                        <v-menu activator="parent">
                            <v-list>
                                <v-list-subheader :title="tt('File Format')"/>
                                <v-list-item :prepend-icon="mdiComma"
                                             :append-icon="fileFormat === 'csv' ? mdiCheck : undefined"
                                             :title="tt('CSV (Comma-separated values) File')"
                                             @click="fileFormat = 'csv'"></v-list-item>
                                <v-list-item :prepend-icon="mdiKeyboardTab"
                                             :append-icon="fileFormat === 'tsv' ? mdiCheck : undefined"
                                             :title="tt('TSV (Tab-separated values) File')"
                                             @click="fileFormat = 'tsv'"></v-list-item>
                                <v-list-item :prepend-icon="mdiLanguageMarkdownOutline"
                                             :append-icon="fileFormat === 'md' ? mdiCheck : undefined"
                                             :title="tt('Markdown File')"
                                             @click="fileFormat = 'md'"></v-list-item>
                            </v-list>
                        </v-menu>
                    </v-btn>
                </div>
            </template>

            <v-card-text class="py-0 w-100 d-flex justify-center">
                <v-switch class="export-data-display-switch" color="secondary"
                          :label="tt('Raw Data')"
                          v-model="showRawData"
                          @click="showRawData = !showRawData">
                    <template #prepend>
                        <span>{{ tt('Table') }}</span>
                    </template>
                </v-switch>
            </v-card-text>

            <v-card-text class="my-md-4 w-100 d-flex justify-center">
                <v-data-table
                    fixed-header
                    fixed-footer
                    multi-sort
                    density="compact"
                    height="365"
                    :headers="dataTableHeaders"
                    :items="dataTableItems"
                    :hover="true"
                    :hide-default-footer="true"
                    :items-per-page="dataTableItems.length"
                    :no-data-text="tt('No data')"
                    v-if="!showRawData"
                ></v-data-table>
                <div class="w-100 pl-2 code-container" v-if="showRawData">
                    <textarea class="w-100" style="outline: none; height: 360px" :readonly="true" :value="exportedData"></textarea>
                </div>
            </v-card-text>

            <v-card-text class="overflow-y-visible">
                <div ref="buttonContainer" class="w-100 d-flex justify-center gap-4">
                    <v-btn-group variant="tonal" density="comfortable">
                        <v-btn color="primary" :disabled="!exportedData" @click="copy">{{ tt('Copy') }}</v-btn>
                        <v-btn density="compact" color="primary" :disabled="!exportedData" :icon="true">
                            <v-icon :icon="mdiMenuDown" size="24" />
                            <v-menu activator="parent">
                                <v-list>
                                    <v-list-item :title="tt('Save')" @click="save()"></v-list-item>
                                </v-list>
                            </v-menu>
                        </v-btn>
                    </v-btn-group>
                    <v-btn color="secondary" variant="tonal" @click="cancel">{{ tt('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useUserStore } from '@/stores/user.ts';

import { replaceAll } from '@/lib/common.ts';
import { copyTextToClipboard, startDownloadFile } from '@/lib/ui/common.ts';

import {
    mdiDotsVertical,
    mdiCheck,
    mdiComma,
    mdiKeyboardTab,
    mdiLanguageMarkdownOutline,
    mdiMenuDown
} from '@mdi/js';

type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();

const userStore = useUserStore();

const buttonContainer = useTemplateRef<HTMLElement>('buttonContainer');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const headers = ref<string[]>([]);
const data = ref<string[][]>([]);
const fileFormat = ref<string>('csv');
const showRawData = ref<boolean>(false);

const dataTableHeaders = computed<object[]>(() => {
    return headers.value.map((header, index) => ({
        key: index.toString(),
        value: `column${index}`,
        title: header,
        sortable: index > 0,
        nowrap: true
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

    if (fileFormat.value === 'csv' || fileFormat.value === 'tsv') {
        let separator = ',';

        if (fileFormat.value === 'tsv') {
            separator = '\t';
        }

        if (headers.value.length > 0) {
            ret += headers.value.map(item => replaceAll(item, separator, ' ')).join(separator);
        }

        for (const row of data.value) {
            ret += '\n';
            ret += row.map(item => replaceAll(item, separator, ' ')).join(separator);
        }
    } else if (fileFormat.value === 'md') {
        ret += '| ' + headers.value.map(item => replaceAll(item, '|', ' ')).join(' | ') + ' |';
        ret += '\n';
        ret += '| ' + headers.value.map(() => '---').join(' | ') + ' |';

        for (const row of data.value) {
            ret += '\n';
            ret += '| ' + row.map(item => replaceAll(item, '|', ' ')).join(' | ') + ' |';
        }
    }


    return ret;
});

function getExportFileName(fileExtension: string): string {
    const nickname = userStore.currentUserNickname;

    if (nickname) {
        return tt('dataExport.exportStatisticsFileName', {
            nickname: nickname
        }) + '.' + fileExtension;
    }

    return tt('dataExport.defaultExportStatisticsFileName') + '.' + fileExtension;
}

function open(options: { headers: string[], data: string[][] }): void {
    headers.value = options.headers || [];
    data.value = options.data || [];
    fileFormat.value = 'csv';
    showRawData.value = false;
    showState.value = true;
}

function copy(): void {
    copyTextToClipboard(exportedData.value, buttonContainer.value);
    snackbar.value?.showMessage('Data copied');
}

function save(): void {
    let contentType = 'text/csv';

    if (fileFormat.value === 'tsv') {
        contentType = 'text/tab-separated-values';
    } else if (fileFormat.value === 'md') {
        contentType = 'text/markdown';
    }

    startDownloadFile(getExportFileName(fileFormat.value), new Blob([exportedData.value], { type: contentType }));
}

function cancel(): void {
    showState.value = false;
}

defineExpose({
    open
});
</script>

<style>
.export-data-display-switch {
    cursor: pointer;
}

.export-data-display-switch.v-input--horizontal .v-input__prepend {
    margin-right: 10px; /* same as the padding-left of `.v-switch .v-label` */
}
</style>
