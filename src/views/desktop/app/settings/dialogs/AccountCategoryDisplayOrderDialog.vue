<template>
    <v-dialog width="600" :persistent="isDisplayOrderModified()" v-model="showState">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <div class="d-flex align-center">
                        <h4 class="text-h4">{{ tt('Account Category Order') }}</h4>
                    </div>
                    <v-spacer/>
                    <v-btn density="comfortable" color="default" variant="text" class="ms-2" :icon="true">
                        <v-icon :icon="mdiDotsVertical" />
                        <v-menu activator="parent">
                            <v-list>
                                <v-list-item :prepend-icon="mdiRestore"
                                             :title="tt('Reset to Default')"
                                             @click="resetDisplayOrderToDefault"></v-list-item>
                            </v-list>
                        </v-menu>
                    </v-btn>
                </div>
            </template>

            <v-card-text class="d-flex flex-column flex-md-row flex-grow-1 overflow-y-auto">
                <v-table hover density="comfortable" class="w-100 table-striped">
                    <draggable-list tag="tbody"
                                    item-key="id"
                                    handle=".drag-handle"
                                    ghost-class="dragging-item"
                                    v-model="accountCategories">
                        <template #item="{ element }">
                            <tr class="text-sm">
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
            </v-card-text>

            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                    <v-btn :disabled="!isDisplayOrderModified()" @click="saveDisplayOrder">{{ tt('Save') }}</v-btn>
                    <v-btn color="secondary" variant="tonal" @click="cancel">{{ tt('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useAccountCategoryDisplayOrderSettingsPageBase } from '@/views/base/settings/AccountCategoryDisplayOrderSettingsPageBase.ts';

import {
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
