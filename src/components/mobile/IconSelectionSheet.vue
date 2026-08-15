<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler"
              :class="heightClass" :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar class="toolbar-with-swipe-handler">
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link sheet-close icon-f7="xmark"></f7-link>
            </div>
            <f7-segmented strong round class="width-100">
                <f7-button :active="currentTab === 'system'" @click="currentTab = 'system'">{{ tt('System Icons') }}</f7-button>
                <f7-button :active="currentTab === 'custom'" @click="currentTab = 'custom'" v-if="isUserCustomIconEnabled()">{{ tt('Custom Icons') }}</f7-button>
            </f7-segmented>
        </f7-toolbar>
        <f7-page-content>
            <f7-block class="margin-vertical no-padding" v-if="currentTab === 'system'">
                <div class="grid padding-vertical-half padding-horizontal-half"
                     :class="{ 'row-has-selected-item': hasSelectedIcon(row) }"
                     :style="`grid-template-columns: repeat(${itemPerRow}, minmax(0, 1fr));`"
                     :key="idx" v-for="(row, idx) in allSystemIconRows">
                    <div class="text-align-center" :key="iconInfo.id" v-for="iconInfo in row">
                        <ItemIcon icon-type="fixed" :icon-id="iconInfo.icon" :color="color" @click="onIconClicked(IconType.System, iconInfo.id)">
                            <f7-badge color="default" class="right-bottom-icon" v-if="currentValue && currentValue === iconInfo.id">
                                <f7-icon f7="checkmark_alt"></f7-icon>
                            </f7-badge>
                        </ItemIcon>
                    </div>
                </div>
            </f7-block>
            <f7-block class="margin-vertical no-padding" v-if="currentTab === 'custom'">
                <div class="grid padding-vertical-half padding-horizontal-half"
                     :class="{ 'row-has-selected-item': hasSelectedIcon(row) }"
                     :style="`grid-template-columns: repeat(${itemPerRow}, minmax(0, 1fr));`"
                     :key="idx" v-for="(row, idx) in allCustomIconRows">
                    <div class="text-align-center" :key="iconInfo.id" v-for="iconInfo in row">
                        <ItemIcon icon-type="user-custom" :icon-id="iconInfo.id" :color="color" @click="onIconClicked(IconType.UserCustom, iconInfo.id)">
                            <f7-badge color="default" class="right-bottom-icon" v-if="currentValue && currentValue === iconInfo.id">
                                <f7-icon f7="checkmark_alt"></f7-icon>
                            </f7-badge>
                        </ItemIcon>
                    </div>
                </div>

                <div class="padding-vertical-half padding-horizontal" v-if="allCustomIconRows.length < 1">
                    {{ tt('No available custom icons') }}
                </div>
            </f7-block>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useUserCustomIconsStore } from '@/stores/userCustomIcon.ts';

import { type SystemIconInfo, type SystemIconInfoWithId, type UserCustomIconInfo, IconType } from '@/core/icon.ts';

import { arrayContainsFieldValue } from '@/lib/common.ts';
import { getSystemIconsInRows, getUserCustomIconsInRows } from '@/lib/icon.ts';
import { scrollToSelectedItem } from '@/lib/ui/common.ts';
import { isUserCustomIconEnabled } from '@/lib/server_settings.ts';
import { type Framework7Dom } from '@/lib/ui/mobile.ts';

const props = defineProps<{
    modelValue: string;
    iconType: IconType;
    show: boolean;
    columnCount?: number;
    color: string;
    allSystemIconInfos: Record<string, SystemIconInfo>;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
    (e: 'update:iconType', value: IconType): void;
    (e: 'update:show', value: boolean): void;
}>();

const { tt } = useI18n();

const customIconsStore = useUserCustomIconsStore();

const currentTab = ref<'system' | 'custom'>(props.iconType === IconType.UserCustom ? 'custom' : 'system');
const currentValue = ref<string>(props.modelValue);
const itemPerRow = ref<number>(props.columnCount || 7);

const allSystemIconRows = computed<SystemIconInfoWithId[][]>(() => getSystemIconsInRows(props.allSystemIconInfos, itemPerRow.value));
const allCustomIconRows = computed<UserCustomIconInfo[][]>(() => getUserCustomIconsInRows(customIconsStore.allCustomIcons, itemPerRow.value));

const heightClass = computed<string>(() => {
    if (allSystemIconRows.value.length > 10 || allCustomIconRows.value.length > 10) {
        return 'icon-selection-huge-sheet';
    } else if (allSystemIconRows.value.length > 6 || allCustomIconRows.value.length > 6) {
        return 'icon-selection-large-sheet';
    } else {
        return 'icon-selection-default-sheet';
    }
});

function hasSelectedIcon(row: (SystemIconInfoWithId | UserCustomIconInfo)[]): boolean {
    return arrayContainsFieldValue(row, 'id', props.modelValue);
}

function onIconClicked(type: IconType, id: string): void {
    currentValue.value = id;
    emit('update:modelValue', currentValue.value);
    emit('update:iconType', type);
}

function onSheetOpen(event: { $el: Framework7Dom }): void {
    const promise = customIconsStore.loadAllCustomIcons({ force: false });
    currentTab.value = props.iconType === IconType.UserCustom ? 'custom' : 'system';

    if (props.iconType === IconType.UserCustom) {
        promise.finally(() => {
            scrollToSelectedItem(event.$el[0], '.sheet-modal-inner', '.page-content', '.row-has-selected-item');
        });
    } else {
        scrollToSelectedItem(event.$el[0], '.sheet-modal-inner', '.page-content', '.row-has-selected-item');
    }
}

function onSheetClosed(): void {
    emit('update:show', false);
}
</script>

<style>
.icon-selection-default-sheet {
    height: 310px;
}

@media (min-height: 630px) {
    .icon-selection-large-sheet {
        height: 370px;
    }

    .icon-selection-huge-sheet {
        height: 500px;
    }
}

@media (max-height: 629px) {
    .icon-selection-large-sheet,
    .icon-selection-huge-sheet {
        height: 320px;
    }
}
</style>
