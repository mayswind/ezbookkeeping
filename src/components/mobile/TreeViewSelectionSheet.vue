<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler"
              style="height: auto" :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left"></div>
            <div class="right">
                <f7-link sheet-close :text="tt('Done')"></f7-link>
            </div>
        </f7-toolbar>
        <f7-searchbar ref="searchbar" custom-searchs
                      :value="filterContent"
                      :placeholder="filterPlaceholder"
                      :disable-button="false"
                      v-if="enableFilter"
                      @input="filterContent = $event.target.value">
        </f7-searchbar>
        <f7-page-content :class="'no-padding-top ' + heightClass">
            <f7-list class="no-margin-top no-margin-bottom" v-if="!filteredItems || !filteredItems.length">
                <f7-list-item :title="filterNoItemsText"></f7-list-item>
            </f7-list>
            <f7-treeview>
                <f7-treeview-item item-toggle
                                  :opened="isPrimaryItemHasSecondaryValue(item)"
                                  :label="ti((primaryTitleField ? item[primaryTitleField] : item) as string, !!primaryTitleI18n)"
                                  :key="primaryKeyField ? item[primaryKeyField] : item"
                                  v-for="item in filteredItems">
                    <template #media>
                        <ItemIcon :icon-type="primaryIconType" :icon-id="item[primaryIconField]"
                                  :color="primaryColorField ? item[primaryColorField] : undefined" v-if="primaryIconField"></ItemIcon>
                    </template>

                    <f7-treeview-item selectable
                                      :selected="isSecondaryValueSelected(currentValue, subItem)"
                                      :label="ti((secondaryTitleField ? (subItem as Record<string, unknown>)[secondaryTitleField] : subItem) as string, !!secondaryTitleI18n)"
                                      :key="secondaryKeyField ? (subItem as Record<string, unknown>)[secondaryKeyField] : subItem"
                                      v-for="subItem in getFilteredSubItems(item)"
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
import { ref, computed, useTemplateRef } from 'vue';
import type { Searchbar } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { type TwoLevelItemSelectionBaseProps, useTwoLevelItemSelectionBase } from '@/components/base/TwoLevelItemSelectionBase.ts';

import { type Framework7Dom, scrollToSelectedItem } from '@/lib/ui/mobile.ts';

interface MobileTwoLevelItemSelectionBaseProps extends TwoLevelItemSelectionBaseProps {
    show: boolean;
}

const props = defineProps<MobileTwoLevelItemSelectionBaseProps>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: unknown): void;
    (e: 'update:show', value: boolean): void;
}>();

const { tt, ti } = useI18n();

const {
    filterContent,
    visibleItemsCount,
    filteredItems,
    getFilteredSubItems,
    isSecondaryValueSelected,
    updateCurrentSecondaryValue
} = useTwoLevelItemSelectionBase(props);

const searchbar = useTemplateRef<Searchbar.Searchbar>('searchbar');

const currentValue = ref<unknown>(props.modelValue);

const heightClass = computed<string>(() => {
    if (visibleItemsCount.value > 6) {
        return 'tree-view-selection-huge-sheet';
    } else if (visibleItemsCount.value > 2) {
        return 'tree-view-selection-large-sheet';
    } else {
        return 'tree-view-selection-default-sheet';
    }
});

function isPrimaryItemHasSecondaryValue(primaryItem: Record<string, unknown>): boolean {
    const subItems = primaryItem[props.primarySubItemsField] as unknown[];

    if (subItems.length < 1) {
        return false;
    }

    const lowerCaseFilterContent = filterContent.value?.toLowerCase() ?? '';

    for (let i = 0; i < subItems.length; i++) {
        const secondaryItem = subItems[i];

        if (props.secondaryHiddenField && (secondaryItem as Record<string, unknown>)[props.secondaryHiddenField]) {
            continue;
        }

        if (props.primaryTitleField && lowerCaseFilterContent) {
            const title = ti((secondaryItem as Record<string, unknown>)[props.primaryTitleField] as string, !!props.primaryTitleI18n);

            if (title.toLowerCase().indexOf(lowerCaseFilterContent) >= 0) {
                return true;
            }
        }

        if (props.secondaryValueField && (secondaryItem as Record<string, unknown>)[props.secondaryValueField] === currentValue.value) {
            return true;
        } else if (!props.secondaryValueField && secondaryItem === currentValue.value) {
            return true;
        }
    }

    return false;
}

function onSecondaryItemClicked(subItem: unknown): void {
    updateCurrentSecondaryValue(currentValue, subItem);
    emit('update:modelValue', currentValue.value);
    emit('update:show', false);
}

function onSheetOpen(event: { $el: Framework7Dom }): void {
    currentValue.value = props.modelValue;
    scrollToSelectedItem(event.$el, '.page-content', '.treeview-item .treeview-item-selected');
}

function onSheetClosed(): void {
    emit('update:show', false);
    filterContent.value = '';
    searchbar.value?.clear();
}
</script>

<style>
.tree-view-selection-default-sheet {
    height: 200px;
}

@media (min-height: 630px) {
    .tree-view-selection-large-sheet {
        height: 260px;
    }

    .tree-view-selection-huge-sheet {
        height: 380px;
    }
}

@media (max-height: 629px) {
    .tree-view-selection-large-sheet,
    .tree-view-selection-huge-sheet {
        height: 240px;
    }
}
</style>
