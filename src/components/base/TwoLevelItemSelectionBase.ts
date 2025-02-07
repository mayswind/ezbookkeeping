import { type Ref, ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { getItemByKeyValue } from '@/lib/common.ts';

export interface TwoLevelItemSelectionBaseProps {
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
    enableFilter?: boolean;
    filterPlaceholder?: string;
    filterNoItemsText?: string;
    items: Record<string, unknown>[];
}

export function useTwoLevelItemSelectionBase(props: TwoLevelItemSelectionBaseProps) {
    const { ti } = useI18n();

    const filterContent = ref<string>('');

    const visibleItemsCount = computed<number>(() => {
        let count = 0;

        for (const item of props.items) {
            if (props.primaryHiddenField && item[props.primaryHiddenField]) {
                continue;
            }

            count++;
        }

        return count;
    });

    const filteredItems = computed<Record<string, unknown>[]>(() => {
        const finalItems: Record<string, unknown>[] = [];
        const items = props.items;
        const lowerCaseFilterContent = filterContent.value?.toLowerCase() ?? '';

        for (const item of items) {
            if (props.primaryHiddenField && item[props.primaryHiddenField]) {
                continue;
            }

            if (!props.enableFilter || !lowerCaseFilterContent) {
                finalItems.push(item);
                continue;
            }

            if (props.primaryTitleField) {
                const title = ti(item[props.primaryTitleField] as string, !!props.primaryTitleI18n);

                if (title.toLowerCase().indexOf(lowerCaseFilterContent) >= 0) {
                    finalItems.push(item);
                    continue;
                }
            }

            if (props.primarySubItemsField) {
                if (getFilteredSubItems(item).length > 0) {
                    finalItems.push(item);
                }
            }
        }

        return finalItems;
    });

    function getFilteredSubItems(selectedPrimaryItem: unknown): Record<string, unknown>[] {
        const finalItems: Record<string, unknown>[] = [];

        if (!selectedPrimaryItem || !props.primarySubItemsField) {
            return finalItems;
        }

        const subItems = (selectedPrimaryItem as Record<string, unknown>)[props.primarySubItemsField] as Record<string, unknown>[];
        let primaryTitleHasFilterContent = false;

        if (props.primaryTitleField) {
            const title = ti((selectedPrimaryItem as Record<string, unknown>)[props.primaryTitleField] as string, !!props.primaryTitleI18n);
            primaryTitleHasFilterContent = title.toLowerCase().indexOf(filterContent.value.toLowerCase()) >= 0;
        }

        for (const subItem of subItems) {
            if (props.secondaryHiddenField && subItem[props.secondaryHiddenField]) {
                continue;
            }

            if (!props.enableFilter || !filterContent.value) {
                finalItems.push(subItem);
                continue;
            }

            if (primaryTitleHasFilterContent) {
                finalItems.push(subItem);
                continue;
            }

            if (props.secondaryTitleField && filterContent.value) {
                const title = ti(subItem[props.secondaryTitleField] as string, !!props.secondaryTitleI18n);

                if (title.toLowerCase().indexOf(filterContent.value.toLowerCase()) >= 0) {
                    finalItems.push(subItem);
                }
            }
        }

        return finalItems;
    }

    function isSecondaryValueSelected(currentSecondaryValue: unknown, subItem: unknown): boolean {
        if (props.secondaryValueField) {
            return currentSecondaryValue === (subItem as Record<string, unknown>)[props.secondaryValueField];
        } else {
            return currentSecondaryValue === subItem;
        }
    }

    function getSelectedSecondaryItem(currentSecondaryValue: unknown, selectedPrimaryItem: unknown): unknown {
        if (currentSecondaryValue && selectedPrimaryItem && (selectedPrimaryItem as Record<string, unknown>)[props.primarySubItemsField]) {
            return getItemByKeyValue((selectedPrimaryItem as Record<string, unknown>)[props.primarySubItemsField] as Record<string, unknown>[], currentSecondaryValue, props.secondaryValueField as string);
        } else {
            return null;
        }
    }

    function updateCurrentSecondaryValue(currentSecondaryValue: Ref<unknown>, subItem: unknown): void {
        if (props.secondaryValueField) {
            currentSecondaryValue.value = (subItem as Record<string, unknown>)[props.secondaryValueField];
        } else {
            currentSecondaryValue.value = subItem;
        }
    }

    return {
        // states
        filterContent,
        // computed states
        visibleItemsCount,
        filteredItems,
        // functions
        getFilteredSubItems,
        isSecondaryValueSelected,
        getSelectedSecondaryItem,
        updateCurrentSecondaryValue
    };
}
