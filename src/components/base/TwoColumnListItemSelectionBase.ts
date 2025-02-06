import { type Ref, ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { getItemByKeyValue, getPrimaryValueBySecondaryValue } from '@/lib/common.ts';

export interface CommonTwoColumnListItemSelectionProps {
    modelValue: unknown;
    primaryKeyField?: string;
    primaryValueField?: string;
    primaryTitleField?: string;
    primaryTitleI18n?: boolean;
    primaryHeaderField?: string;
    primaryHeaderI18n?: boolean;
    primaryFooterField?: string;
    primaryFooterI18n?: boolean;
    primaryIconField?: string;
    primaryIconType?: string;
    primaryColorField?: string;
    primaryHiddenField?: string;
    primarySubItemsField: string;
    secondaryKeyField?: string;
    secondaryValueField?: string;
    secondaryTitleField?: string;
    secondaryTitleI18n?: boolean;
    secondaryHeaderField?: string;
    secondaryHeaderI18n?: boolean;
    secondaryFooterField?: string;
    secondaryFooterI18n?: boolean;
    secondaryIconField?: string;
    secondaryIconType?: string;
    secondaryColorField?: string;
    secondaryHiddenField?: string;
    enableFilter?: boolean;
    filterPlaceholder?: string;
    filterNoItemsText?: string;
    items: Record<string, unknown>[];
}

export function useTwoColumnListItemSelectionBase(props: CommonTwoColumnListItemSelectionProps) {
    const { ti } = useI18n();

    const filterContent = ref<string>('');

    const filteredItems = computed<Record<string, unknown>[]>(() => {
        const finalItems: Record<string, unknown>[] = [];
        const items = props.items;

        for (const item of items) {
            if (props.primaryHiddenField && item[props.primaryHiddenField]) {
                continue;
            }

            if (!props.enableFilter || !filterContent.value) {
                finalItems.push(item);
                continue;
            }

            if (props.primaryTitleField) {
                const title = ti(item[props.primaryTitleField] as string, !!props.primaryTitleI18n);

                if (title.toLowerCase().indexOf(filterContent.value.toLowerCase()) >= 0) {
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

    function getCurrentPrimaryValueBySecondaryValue(secondaryValue: unknown): unknown {
        return getPrimaryValueBySecondaryValue(props.items as Record<string, Record<string, unknown>[]>[], props.primarySubItemsField, props.primaryValueField, props.primaryHiddenField, props.secondaryValueField, props.secondaryHiddenField, secondaryValue);
    }

    function isSecondaryValueSelected(currentSecondaryValue: unknown, subItem: unknown): boolean {
        if (props.secondaryValueField) {
            return currentSecondaryValue === (subItem as Record<string, unknown>)[props.secondaryValueField];
        } else {
            return currentSecondaryValue === subItem;
        }
    }

    function getSelectedPrimaryItem(currentPrimaryValue: unknown): unknown {
        if (props.primaryValueField) {
            return getItemByKeyValue(props.items, currentPrimaryValue, props.primaryValueField);
        } else {
            return currentPrimaryValue;
        }
    }

    function getSelectedSecondaryItem(currentSecondaryValue: unknown, selectedPrimaryItem: unknown): unknown {
        if (currentSecondaryValue && selectedPrimaryItem && (selectedPrimaryItem as Record<string, unknown>)[props.primarySubItemsField]) {
            return getItemByKeyValue((selectedPrimaryItem as Record<string, unknown>)[props.primarySubItemsField] as Record<string, unknown>[], currentSecondaryValue, props.secondaryValueField as string);
        } else {
            return null;
        }
    }

    function updateCurrentPrimaryValue(currentPrimaryValue: Ref<unknown>, item: unknown): void {
        if (props.primaryValueField) {
            currentPrimaryValue.value = (item as Record<string, unknown>)[props.primaryValueField];
        } else {
            currentPrimaryValue.value = item;
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
        filteredItems,
        // functions
        getFilteredSubItems,
        getCurrentPrimaryValueBySecondaryValue,
        isSecondaryValueSelected,
        getSelectedPrimaryItem,
        getSelectedSecondaryItem,
        updateCurrentPrimaryValue,
        updateCurrentSecondaryValue
    };
}
