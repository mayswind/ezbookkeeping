<template>
    <v-select
        class="icon-select"
        density="comfortable"
        item-title="icon"
        item-value="id"
        persistent-placeholder
        :disabled="disabled"
        :label="label"
        :menu-props="{ contentClass: 'icon-select-menu' }"
        v-model="icon"
        @update:menu="onMenuStateChanged"
    >
        <template #selection>
            <v-label class="cursor-pointer" style="padding-top: 3px">
                <ItemIcon :icon-type="getIconType(type, iconType)" :icon-id="icon" :color="color" />
            </v-label>
        </template>

        <template #no-data>
            <div class="icon-select-dropdown px-2" ref="dropdownMenu">
                <div class="icon-select-tabs-container">
                    <v-tabs grow density="compact" v-model="currentTab">
                        <v-tab value="system">{{ tt('System Icons') }}</v-tab>
                        <v-tab value="custom" v-if="isUserCustomIconEnabled()">{{ tt('Custom Icons') }}</v-tab>
                    </v-tabs>
                </div>
                <div v-if="currentTab === 'system'">
                    <div class="icon-item" :class="{ 'row-has-selected-item': hasSelectedIcon(row) }"
                         :style="`grid-template-columns: repeat(${itemPerRow}, minmax(0, 1fr));`"
                         :key="idx" v-for="(row, idx) in allSystemIconRows">
                        <div class="text-center" :key="iconInfo.id" v-for="iconInfo in row">
                            <div class="cursor-pointer" @click="updateIcon(IconType.System, iconInfo.id)">
                                <ItemIcon class="ma-2" icon-type="fixed" :icon-id="iconInfo.icon" :color="color" v-if="!modelValue || modelValue !== iconInfo.id" />
                                <v-badge class="right-bottom-icon" color="primary"
                                         offset-x="8" offset-y="10"
                                         :location="`bottom ${textDirection === TextDirection.LTR ? 'right' : 'left'}`"
                                         :icon="mdiCheck"
                                         v-if="modelValue && modelValue === iconInfo.id">
                                    <ItemIcon class="ma-2" icon-type="fixed" :icon-id="iconInfo.icon" :color="color" />
                                </v-badge>
                            </div>
                        </div>
                    </div>
                </div>
                <div v-if="currentTab === 'custom'">
                    <div class="icon-item" :class="{ 'row-has-selected-item': hasSelectedIcon(row) }"
                         :style="`grid-template-columns: repeat(${itemPerRow}, minmax(0, 1fr));`"
                         :key="idx" v-for="(row, idx) in allCustomIconRows">
                        <div class="text-center" :key="iconInfo.id" v-for="iconInfo in row">
                            <div class="cursor-pointer" @click="updateIcon(IconType.UserCustom, iconInfo.id)">
                                <ItemIcon class="ma-2" icon-type="user-custom" :icon-id="iconInfo.id" :color="color" v-if="!modelValue || modelValue !== iconInfo.id" />
                                <v-badge class="right-bottom-icon" color="primary"
                                         offset-x="8" offset-y="10"
                                         :location="`bottom ${textDirection === TextDirection.LTR ? 'right' : 'left'}`"
                                         :icon="mdiCheck"
                                         v-if="modelValue && modelValue === iconInfo.id">
                                    <ItemIcon class="ma-2" icon-type="user-custom" :icon-id="iconInfo.id" :color="color" />
                                </v-badge>
                            </div>
                        </div>
                    </div>

                    <div class="icon-item text-body-large py-2 ms-2" v-if="allCustomIconRows.length < 1">
                        {{ tt('No available custom icons') }}
                    </div>
                </div>
            </div>
        </template>
    </v-select>
</template>

<script setup lang="ts">
import { ref, computed, useTemplateRef, nextTick } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useUserCustomIconsStore } from '@/stores/userCustomIcon.ts';

import { TextDirection } from '@/core/text.ts';
import type { ColorValue } from '@/core/color.ts';
import { type SystemIconInfo, type SystemIconInfoWithId, type UserCustomIconInfo, IconType } from '@/core/icon.ts';

import { arrayContainsFieldValue } from '@/lib/common.ts';
import { getIconType, getSystemIconsInRows, getUserCustomIconsInRows } from '@/lib/icon.ts';
import { isUserCustomIconEnabled } from '@/lib/server_settings.ts';
import { scrollToSelectedItem } from '@/lib/ui/common.ts';

import {
    mdiCheck
} from '@mdi/js';

const props = defineProps<{
    modelValue: string;
    disabled?: boolean;
    label?: string;
    type: string;
    iconType: number;
    color: ColorValue;
    columnCount?: number;
    allSystemIconInfos: Record<string, SystemIconInfo>;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
    (e: 'update:iconType', value: IconType): void;
}>();

const { tt, getCurrentLanguageTextDirection } = useI18n();

const customIconsStore = useUserCustomIconsStore();

const dropdownMenu = useTemplateRef<HTMLElement>('dropdownMenu');

const currentTab = ref<'system' | 'custom'>(props.iconType === IconType.UserCustom ? 'custom' : 'system');
const itemPerRow = ref<number>(props.columnCount || 7);

const textDirection = computed<TextDirection>(() => getCurrentLanguageTextDirection());
const allSystemIconRows = computed<SystemIconInfoWithId[][]>(() => getSystemIconsInRows(props.allSystemIconInfos, itemPerRow.value));
const allCustomIconRows = computed<UserCustomIconInfo[][]>(() => getUserCustomIconsInRows(customIconsStore.allCustomIcons, itemPerRow.value));

const icon = computed<string>({
    get: () => props.modelValue,
    set: (value: string) => emit('update:modelValue', value)
});

function init(): void {
    customIconsStore.loadAllCustomIcons({ force: false });
}

function hasSelectedIcon(row: (SystemIconInfoWithId | UserCustomIconInfo)[]): boolean {
    return arrayContainsFieldValue(row, 'id', props.modelValue);
}

function updateIcon(type: IconType, id: string): void {
    icon.value = id;
    emit('update:iconType', type);
}

function onMenuStateChanged(state: boolean): void {
    if (state) {
        currentTab.value = props.iconType === IconType.UserCustom ? 'custom' : 'system';

        nextTick(() => {
            if (dropdownMenu.value && dropdownMenu.value.parentElement) {
                scrollToSelectedItem(dropdownMenu.value.parentElement, null, null, '.row-has-selected-item');
            }
        });
    }
}

init();
</script>

<style>
.icon-select:not(.v-input--disabled) .v-field__input,
.icon-select:not(.v-input--disabled) .v-label {
    opacity: 1;
}

.icon-select-menu > .v-sheet > .v-list {
    height: 300px;
    overflow-y: scroll;
}

.icon-select-dropdown {
    display: flow-root;
}

.icon-select-dropdown .icon-item {
    display: grid;
}

.icon-select-dropdown .icon-select-tabs-container {
    position: sticky;
    top: -8px;
    z-index: 1;
    margin: -8px -8px 8px;
    padding: 8px 8px 0;
    background-color: rgb(var(--v-theme-surface));
}
</style>
