<template>
    <v-dialog width="600" :persistent="isDisplayOrderModified()" v-model="showState">
        <one-column-dialog-layout content-class="pa-0"
                                  :title="tt('Account Category Order')" :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #after-title>
                <v-btn density="compact" color="primary" variant="text" class="ms-2" :icon="true"
                       :disabled="!isDisplayOrderModified()" @click="saveDisplayOrder">
                    <v-icon :icon="mdiCheck" size="22" />
                    <v-tooltip activator="parent">{{ tt('Save') }}</v-tooltip>
                </v-btn>
            </template>

            <template #toolbar>
                <v-btn density="compact" color="default" variant="text" class="ms-2" :icon="true">
                    <v-icon :icon="mdiDotsVertical" size="22" />
                    <v-menu activator="parent">
                        <v-list>
                            <v-list-item :prepend-icon="mdiRestore"
                                         :title="tt('Reset to Default')"
                                         @click="resetDisplayOrderToDefault"></v-list-item>
                        </v-list>
                    </v-menu>
                </v-btn>
            </template>

            <template #content>
                <div class="d-flex flex-column flex-md-row flex-grow-1 overflow-y-auto">
                    <v-table hover density="comfortable" class="w-100 table-striped">
                        <draggable-list tag="tbody"
                                        item-key="id"
                                        handle=".drag-handle"
                                        ghost-class="dragging-item"
                                        v-model="accountCategories">
                            <template #item="{ element }">
                                <tr>
                                    <td>
                                        <div class="d-flex align-center">
                                            <div class="d-flex align-center">
                                                <span>{{ tt(element.name) }}</span>
                                            </div>

                                            <v-spacer/>

                                            <span class="ms-2">
                                            <v-icon class="drag-handle" :icon="mdiDrag"/>
                                            <v-tooltip activator="parent">{{ tt('Drag to Reorder') }}</v-tooltip>
                                        </span>
                                        </div>
                                    </td>
                                </tr>
                            </template>
                        </draggable-list>
                    </v-table>
                </div>
            </template>
        </one-column-dialog-layout>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useAccountCategoryDisplayOrderSettingsPageBase } from '@/views/base/settings/AccountCategoryDisplayOrderSettingsPageBase.ts';

import {
    mdiCheck,
    mdiDotsVertical,
    mdiRestore,
    mdiDrag
} from '@mdi/js';

const { tt } = useI18n();

const {
    accountCategories,
    isDisplayOrderModified,
    loadDisplayOrderFromSettings,
    saveDisplayOrderToSettings,
    resetDisplayOrderToDefault
} = useAccountCategoryDisplayOrderSettingsPageBase();

let resolveFunc: (() => void) | null = null;
let rejectFunc: (() => void) | null = null;

const showState = ref<boolean>(false);

function open(): Promise<void> {
    loadDisplayOrderFromSettings();
    showState.value = true;

    return new Promise<void>((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function saveDisplayOrder(): void {
    saveDisplayOrderToSettings();
    resolveFunc?.();
    showState.value = false;
}

function cancel(): void {
    rejectFunc?.();
    showState.value = false;
}

defineExpose({
    open
});
</script>
