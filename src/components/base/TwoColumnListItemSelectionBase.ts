import { type Ref } from 'vue';

import {
    type TwoLevelItemSelectionBaseProps,
    useTwoLevelItemSelectionBase
} from '@/components/base/TwoLevelItemSelectionBase.ts';

import { getItemByKeyValue, getPrimaryValueBySecondaryValue } from '@/lib/common.ts';

export interface CommonTwoColumnListItemSelectionProps extends TwoLevelItemSelectionBaseProps {
    primaryValueField?: string;
    primaryHeaderField?: string;
    primaryHeaderI18n?: boolean;
    primaryFooterField?: string;
    primaryFooterI18n?: boolean;
    secondaryHeaderField?: string;
    secondaryHeaderI18n?: boolean;
    secondaryFooterField?: string;
    secondaryFooterI18n?: boolean;
}

export function useTwoColumnListItemSelectionBase(props: CommonTwoColumnListItemSelectionProps) {
    const {
        filterContent,
        visibleItemsCount,
        filteredItems,
        getFilteredSubItems,
        isSecondaryValueSelected,
        getSelectedSecondaryItem,
        updateCurrentSecondaryValue
    } = useTwoLevelItemSelectionBase(props);

    function getCurrentPrimaryValueBySecondaryValue(secondaryValue: unknown): unknown {
        return getPrimaryValueBySecondaryValue(props.items as Record<string, Record<string, unknown>[]>[], props.primarySubItemsField, props.primaryValueField, props.primaryHiddenField, props.secondaryValueField, props.secondaryHiddenField, secondaryValue);
    }

    function getSelectedPrimaryItem(currentPrimaryValue: unknown): unknown {
        if (props.primaryValueField) {
            return getItemByKeyValue(props.items, currentPrimaryValue, props.primaryValueField);
        } else {
            return currentPrimaryValue;
        }
    }

    function updateCurrentPrimaryValue(currentPrimaryValue: Ref<unknown>, item: unknown): void {
        if (props.primaryValueField) {
            currentPrimaryValue.value = (item as Record<string, unknown>)[props.primaryValueField];
        } else {
            currentPrimaryValue.value = item;
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
        getCurrentPrimaryValueBySecondaryValue,
        isSecondaryValueSelected,
        getSelectedPrimaryItem,
        getSelectedSecondaryItem,
        updateCurrentPrimaryValue,
        updateCurrentSecondaryValue
    };
}
