import { type Ref } from 'vue';

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
    items: unknown[];
}

export function useTwoColumnListItemSelectionBase(props: CommonTwoColumnListItemSelectionProps) {
    function getCurrentPrimaryValueBySecondaryValue(secondaryValue: unknown): unknown {
        return getPrimaryValueBySecondaryValue(props.items as Record<string, Record<string, unknown>[]>[] | Record<string, Record<string, Record<string, unknown>[]>>, props.primarySubItemsField, props.primaryValueField, props.primaryHiddenField, props.secondaryValueField, props.secondaryHiddenField, secondaryValue);
    }

    function isSecondaryValueSelected(currentSecondaryValue: unknown, subItem: unknown): boolean {
        if (props.secondaryValueField) {
            return currentSecondaryValue === (subItem as Record<string, unknown>)[props.secondaryValueField];
        } else {
            return currentSecondaryValue === subItem;
        }
    }

    function getSelectedPrimaryItem(currentPrimaryValue: unknown) {
        if (props.primaryValueField) {
            return getItemByKeyValue(props.items as Record<string, unknown>[] | Record<string, Record<string, unknown>>, currentPrimaryValue, props.primaryValueField);
        } else {
            return currentPrimaryValue;
        }
    }

    function getSelectedSecondaryItem(currentSecondaryValue: unknown, selectedPrimaryItem: unknown) {
        if (currentSecondaryValue && selectedPrimaryItem && (selectedPrimaryItem as Record<string, unknown>)[props.primarySubItemsField]) {
            return getItemByKeyValue((selectedPrimaryItem as Record<string, unknown>)[props.primarySubItemsField] as Record<string, unknown>[] | Record<string, Record<string, unknown>>, currentSecondaryValue, props.secondaryValueField as string);
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
        // functions
        getCurrentPrimaryValueBySecondaryValue,
        isSecondaryValueSelected,
        getSelectedPrimaryItem,
        getSelectedSecondaryItem,
        updateCurrentPrimaryValue,
        updateCurrentSecondaryValue
    };
}
