<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler"
              :class="heightClass"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left"></div>
            <div class="right">
                <f7-link sheet-close :text="tt('Done')"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <f7-treeview>
                <f7-treeview-item item-toggle
                                  :opened="isPrimaryItemHasSecondaryValue(item)"
                                  :label="ti((primaryTitleField ? item[primaryTitleField] : item) as string, !!primaryTitleI18n)"
                                  :key="primaryKeyField ? item[primaryKeyField] : item"
                                  v-for="item in items"
                                  v-show="item && (!primaryHiddenField || !item[primaryHiddenField])">
                    <template #media>
                        <ItemIcon :icon-type="primaryIconType" :icon-id="item[primaryIconField]"
                                  :color="primaryColorField ? item[primaryColorField] : undefined" v-if="primaryIconField"></ItemIcon>
                    </template>

                    <f7-treeview-item selectable
                                      :selected="isSecondarySelected(subItem)"
                                      :label="ti((secondaryTitleField ? (subItem as Record<string, unknown>)[secondaryTitleField] : subItem) as string, !!secondaryTitleI18n)"
                                      :key="secondaryKeyField ? (subItem as Record<string, unknown>)[secondaryKeyField] : subItem"
                                      v-for="subItem in item[primarySubItemsField]"
                                      v-show="subItem && (!secondaryHiddenField || !(subItem as Record<string, unknown>)[secondaryHiddenField])"
                                      @click="onSecondaryItemClicked(subItem)">
                        <template #media>
                            <ItemIcon :icon-type="secondaryIconType" :icon-id="(subItem as Record<string, unknown>)[secondaryIconField]"
                                      :color="secondaryColorField ? (subItem as Record<string, unknown>)[secondaryColorField] : undefined" v-if="secondaryIconField"></ItemIcon>
                        </template>
                    </f7-treeview-item>
                </f7-treeview-item>
            </f7-treeview>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { isArray } from '@/lib/common.ts';
import { type Framework7Dom, scrollToSelectedItem } from '@/lib/ui/mobile.ts';

const props = defineProps<{
    modelValue: unknown;
    primaryKeyField?: string;
    primaryTitleField?: string;
    primaryTitleI18n?: boolean;
    primaryIconField?: string;
    primaryIconType?: string;
    primaryColorField?: string;
    primaryHiddenField?: string;
    primarySubItemsField: string;
    secondaryKeyField?: string;
    secondaryValueField?: string;
    secondaryTitleField?: string;
    secondaryTitleI18n?: boolean;
    secondaryIconField?: string;
    secondaryIconType?: string;
    secondaryColorField?: string;
    secondaryHiddenField?: string;
    items: Record<string, unknown>[];
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: unknown): void;
    (e: 'update:show', value: boolean): void;
}>();

const { tt, ti } = useI18n();

const currentValue = ref<unknown>(props.modelValue);

const heightClass = computed<string>(() => {
    let count = 0;

    if (isArray(props.items)) {
        count = props.items.length;
    } else {
        for (const field in props.items) {
            if (!Object.prototype.hasOwnProperty.call(props.items, field)) {
                continue;
            }

            count++;
        }
    }

    if (count > 6) {
        return 'tree-view-selection-huge-sheet';
    } else if (count > 2) {
        return 'tree-view-selection-large-sheet';
    } else {
        return '';
    }
});

function isPrimaryItemHasSecondaryValue(primaryItem: Record<string, unknown>): boolean {
    const subItems = primaryItem[props.primarySubItemsField] as unknown[];

    if (subItems.length < 1) {
        return false;
    }

    for (let i = 0; i < subItems.length; i++) {
        const secondaryItem = subItems[i];

        if (props.secondaryHiddenField && (secondaryItem as Record<string, unknown>)[props.secondaryHiddenField]) {
            continue;
        }

        if (props.secondaryValueField && (secondaryItem as Record<string, unknown>)[props.secondaryValueField] === currentValue.value) {
            return true;
        } else if (!props.secondaryValueField && secondaryItem === currentValue.value) {
            return true;
        }
    }

    return false;
}

function isSecondarySelected(subItem: unknown): boolean {
    if (props.secondaryValueField) {
        return currentValue.value === (subItem as Record<string, unknown>)[props.secondaryValueField];
    } else {
        return currentValue.value === subItem;
    }
}

function onSecondaryItemClicked(subItem: unknown): void {
    if (props.secondaryValueField) {
        currentValue.value = (subItem as Record<string, unknown>)[props.secondaryValueField];
    } else {
        currentValue.value = subItem;
    }

    emit('update:modelValue', currentValue.value);
    emit('update:show', false);
}

function onSheetOpen(event: { $el: Framework7Dom }): void {
    currentValue.value = props.modelValue;
    scrollToSelectedItem(event.$el, '.page-content', '.treeview-item .treeview-item-selected');
}

function onSheetClosed(): void {
    emit('update:show', false);
}
</script>

<style>
@media (min-height: 630px) {
    .tree-view-selection-large-sheet {
        height: 310px;
    }

    .tree-view-selection-huge-sheet {
        height: 400px;
    }
}
</style>
