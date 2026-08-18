<template>
    <main-page-layout no-navbar>
        <template #top-toolbar>
            <v-btn class="top-navigation-button" density="comfortable" variant="text" :icon="true"
                   color="primary" :aria-label="tt('Save')" @click="save">
                <v-icon :icon="mdiContentSaveOutline" size="24" />
                <v-tooltip activator="parent">{{ tt('Save') }}</v-tooltip>
            </v-btn>

            <v-btn class="top-navigation-button ms-1" density="comfortable" variant="text" :icon="true"
                   :aria-label="tt('Cancel')" @click="cancel">
                <v-icon :icon="mdiClose" size="24" />
                <v-tooltip activator="parent">{{ tt('Cancel') }}</v-tooltip>
            </v-btn>

            <v-btn class="top-navigation-button ms-1" density="comfortable" variant="text" :icon="true"
                   :aria-label="tt('More')">
                <v-icon :icon="mdiDotsVertical" size="24" />
                <v-tooltip activator="parent">{{ tt('More') }}</v-tooltip>
                <v-menu activator="parent">
                    <v-list>
                        <v-list-item :prepend-icon="mdiPlus" :title="tt('Add Widget')" @click="addWidget" />
                        <v-list-item :prepend-icon="mdiDeleteSweepOutline" :title="tt('Clear Layout')" @click="clearLayout" />
                        <v-divider class="my-2" />
                        <v-list-item :prepend-icon="mdiApplicationImport" :title="tt('Import Layout')" @click="importLayout" />
                        <v-list-item :prepend-icon="mdiApplicationExport" :title="tt('Export Layout')" @click="exportLayout" />
                        <v-divider class="my-2" />
                        <v-list-item :prepend-icon="mdiRestore" :title="tt('Reset to Default')" @click="resetLayout" />
                    </v-list>
                </v-menu>
            </v-btn>
        </template>

        <template #content>
            <overview-dashboard editing :layout="draftLayout" :loading="loadingOverview"
                                @update:layout="draftLayout = $event" @configure="configureWidget"
                                @add="addWidget" @remove="removeWidget" @refresh="reload(true)" />
        </template>
    </main-page-layout>

    <add-widget-dialog ref="addWidgetDialog" />
    <widget-settings-dialog ref="widgetSettingsDialog" />
    <json-import-dialog ref="layoutImportDialog" :title="tt('Import Layout')" :on-import="onImportLayout" />
    <json-export-dialog ref="layoutExportDialog" :title="tt('Export Layout')" :file-name="tt('dataExport.defaultOverviewLayoutFileName')" />

    <confirm-dialog ref="confirmDialog" />
    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import OverviewDashboard from './OverviewDashboard.vue';
import AddWidgetDialog from './dialogs/AddWidgetDialog.vue';
import WidgetSettingsDialog from './dialogs/WidgetSettingsDialog.vue';
import JsonImportDialog from '@/components/desktop/JsonImportDialog.vue';
import JsonExportDialog from '@/components/desktop/JsonExportDialog.vue';
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef, onMounted, onBeforeUnmount } from 'vue';
import { useRouter, onBeforeRouteLeave } from 'vue-router';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useOverviewStore } from '@/stores/overview.ts';

import { itemAndIndex } from '@/core/base.ts';
import {
    type OverviewWidgetType,
    type DesktopOverviewWidgetDefinition,
    type DesktopOverviewLayout,
    type DesktopOverviewWidgetLayout,
    OverviewWidgetDataRequirement
} from '@/core/overview_layout.ts';
import {
    DESKTOP_OVERVIEW_WIDGET_DEFINITIONS,
    DEFAULT_DESKTOP_OVERVIEW_LAYOUT
} from '@/consts/overview_layout.ts';

import {
    cloneOverviewLayout,
    findOverviewWidgetPosition,
    getOverviewDataRequirements,
    isDefaultDesktopOverviewLayout,
    normalizeDesktopOverviewLayout,
    parseDesktopOverviewLayout,
    serializeDesktopOverviewLayout
} from '@/lib/overview_layout.ts';
import { generateRandomUUID } from '@/lib/misc.ts';
import logger from '@/lib/logger.ts';

import {
    mdiApplicationExport,
    mdiApplicationImport,
    mdiClose,
    mdiContentSaveOutline,
    mdiDeleteSweepOutline,
    mdiDotsVertical,
    mdiPlus,
    mdiRestore
} from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof SnackBar>;
type AddWidgetDialogType = InstanceType<typeof AddWidgetDialog>;
type WidgetSettingsDialogType = InstanceType<typeof WidgetSettingsDialog>;
type JsonImportDialogType = InstanceType<typeof JsonImportDialog>;
type JsonExportDialogType = InstanceType<typeof JsonExportDialog>;

const router = useRouter();

const { tt } = useI18n();

const settingsStore = useSettingsStore();
const accountsStore = useAccountsStore();
const transactionCategoriesStore = useTransactionCategoriesStore();
const overviewStore = useOverviewStore();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');
const addWidgetDialog = useTemplateRef<AddWidgetDialogType>('addWidgetDialog');
const widgetSettingsDialog = useTemplateRef<WidgetSettingsDialogType>('widgetSettingsDialog');
const layoutImportDialog = useTemplateRef<JsonImportDialogType>('layoutImportDialog');
const layoutExportDialog = useTemplateRef<JsonExportDialogType>('layoutExportDialog');

const loadingOverview = ref<boolean>(true);
const leavingAfterAction = ref<boolean>(false);
const initialLayout = ref<DesktopOverviewLayout>(getInitialLayout());
const draftLayout = ref<DesktopOverviewLayout>(cloneOverviewLayout(initialLayout.value));
const initialJson = ref<string>(serializeDesktopOverviewLayout(initialLayout.value));

const isModified = computed<boolean>(() => serializeDesktopOverviewLayout(draftLayout.value) !== initialJson.value);

function getInitialLayout(): DesktopOverviewLayout {
    let initialLayout: DesktopOverviewLayout;

    try {
        initialLayout = parseDesktopOverviewLayout(settingsStore.appSettings.desktopOverviewPageLayout);
    } catch (error) {
        logger.warn('failed to parse desktop overview page layout in editor', error);
        initialLayout = cloneOverviewLayout(DEFAULT_DESKTOP_OVERVIEW_LAYOUT);
    }

    return initialLayout;
}

function reload(force: boolean): void {
    loadingOverview.value = true;

    const requirements = getOverviewDataRequirements(draftLayout.value);
    const promises: Promise<unknown>[] = [
        accountsStore.loadAllAccounts({ force: false }),
        transactionCategoriesStore.loadAllCategories({ force: false })
    ];

    if (requirements.includes(OverviewWidgetDataRequirement.TransactionOverview)) {
        promises.push(overviewStore.loadTransactionOverview({
            force: force,
            loadLast11Months: requirements.includes(OverviewWidgetDataRequirement.TransactionOverviewLast12Months)
        }));
    }

    Promise.all(promises).then(() => {
        loadingOverview.value = false;
        if (force) snackbar.value?.showMessage('Data has been updated');
    }).catch(error => {
        loadingOverview.value = false;
        if (!error.processed && !error.isUpToDate) snackbar.value?.showError(error);
    });
}

function addWidget(): void {
    addWidgetDialog.value?.open().then((type: OverviewWidgetType) => {
        const definition: DesktopOverviewWidgetDefinition = DESKTOP_OVERVIEW_WIDGET_DEFINITIONS[type];
        const position = findOverviewWidgetPosition(draftLayout.value.widgets, definition.defaultWidth, definition.defaultHeight);
        const newWidget: DesktopOverviewWidgetLayout = {
            id: generateRandomUUID(),
            type: type,
            ...position,
            w: definition.defaultWidth,
            h: definition.defaultHeight,
            settings: { ...definition.defaultSettings }
        };

        draftLayout.value.widgets.push(newWidget);
        reload(false);
    });
}

function removeWidget(id: string): void {
    for (const [widget, index] of itemAndIndex(draftLayout.value.widgets)) {
        if (widget.id === id) {
            draftLayout.value.widgets.splice(index, 1);
            return;
        }
    }
}

function configureWidget(widget: DesktopOverviewWidgetLayout): void {
    widgetSettingsDialog.value?.open(widget).then(updatedWidget => {
        for (const widget of draftLayout.value.widgets) {
            if (widget.id === updatedWidget.id) {
                widget.settings = updatedWidget.settings;
                break;
            }
        }
        reload(false);
    });
}

function clearLayout(): void {
    if (!draftLayout.value.widgets.length) {
        return;
    }

    confirmDialog.value?.open('Clear all widgets from the layout?').then(() => {
        draftLayout.value.widgets.length = 0;
    });
}

function resetLayout(): void {
    confirmDialog.value?.open('Reset the layout to its default value?').then(() => {
        draftLayout.value = cloneOverviewLayout(DEFAULT_DESKTOP_OVERVIEW_LAYOUT);
        reload(false);
    });
}

function importLayout(): void {
    layoutImportDialog.value?.open();
}

function exportLayout(): void {
    layoutExportDialog.value?.open({
        json: serializeDesktopOverviewLayout(draftLayout.value, true)
    });
}

function save(): void {
    const normalized = normalizeDesktopOverviewLayout(draftLayout.value);
    settingsStore.setDesktopOverviewPageLayout(isDefaultDesktopOverviewLayout(normalized) ? '' : serializeDesktopOverviewLayout(normalized));

    leavingAfterAction.value = true;
    router.push('/');
}

function cancel(): void {
    const leave = () => {
        leavingAfterAction.value = true;
        router.push('/app/settings/basic');
    };

    if (!isModified.value) {
        leave();
        return;
    }

    confirmDialog.value?.open('Discard unsaved layout changes?').then(leave);
}

function onImportLayout(data: string): boolean {
    try {
        draftLayout.value = parseDesktopOverviewLayout(data);
        reload(false);
        return true;
    } catch (error) {
        logger.error('Failed to import overview layout', error);
        snackbar.value?.showError('Layout import failed. Please make sure the layout is valid and try again.');
        return false;
    }
}

function onBeforeUnload(event: BeforeUnloadEvent): void {
    if (isModified.value) {
        event.preventDefault();
    }
}

onMounted(() => {
    window.addEventListener('beforeunload', onBeforeUnload);
});

onBeforeUnmount(() => {
    window.removeEventListener('beforeunload', onBeforeUnload)
});

onBeforeRouteLeave(() => {
    if (!leavingAfterAction.value && isModified.value) {
        return window.confirm(tt('Discard unsaved layout changes?'));
    }

    return true;
});

reload(false);
</script>
