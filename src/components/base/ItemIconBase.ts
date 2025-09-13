import { computed } from 'vue';

import type {ColorStyleValue, ColorValue} from '@/core/color.ts';
import { ALL_ACCOUNT_ICONS, DEFAULT_ACCOUNT_ICON, ALL_CATEGORY_ICONS, DEFAULT_CATEGORY_ICON } from '@/consts/icon.ts';
import { DEFAULT_ICON_COLOR, DEFAULT_ACCOUNT_COLOR, DEFAULT_CATEGORY_COLOR, DEFAULT_COLOR_STYLE_VARIABLE } from '@/consts/color.ts';

import { isNumber } from '@/lib/common.ts';

type IconItemStyleName = string;
type IconItemStyleValue = ColorValue | string | number | undefined;
type CommonIconItemType = 'account' | 'category' | 'fixed';
type MobileIconItemType = 'fixed-f7';

export interface CommonIconProps {
    iconType: CommonIconItemType | MobileIconItemType;
    iconId: string | number;
    color?: ColorValue;
    defaultColor?: ColorStyleValue;
    additionalColorAttr?: string;
    size?: string | number;
}

export function useItemIconBase(props: CommonIconProps) {
    const style = computed<Record<IconItemStyleName, IconItemStyleValue>>(() => {
        let defaultColor: ColorStyleValue = DEFAULT_COLOR_STYLE_VARIABLE;

        if (props.defaultColor) {
            defaultColor = props.defaultColor;
        }

        if (props.iconType === 'account') {
            return getAccountIconStyle(props.color, defaultColor, props.additionalColorAttr);
        } else if (props.iconType === 'category') {
            return getCategoryIconStyle(props.color, defaultColor, props.additionalColorAttr);
        } else {
            return getDefaultIconStyle(props.color, defaultColor, props.additionalColorAttr);
        }
    });

    function getAccountIcon(iconId: string | number): string {
        if (isNumber(iconId)) {
            iconId = iconId.toString();
        }

        const iconInfo = ALL_ACCOUNT_ICONS[iconId];

        if (!iconInfo) {
            return DEFAULT_ACCOUNT_ICON.icon;
        }

        return iconInfo.icon;
    }

    function getCategoryIcon(iconId: string | number): string {
        if (isNumber(iconId)) {
            iconId = iconId.toString();
        }

        const iconInfo = ALL_CATEGORY_ICONS[iconId];

        if (!iconInfo) {
            return DEFAULT_CATEGORY_ICON.icon;
        }

        return iconInfo.icon;
    }

    function getAccountIconStyle(color?: ColorValue, defaultColor?: ColorStyleValue, additionalColorAttr?: string): Record<IconItemStyleName, IconItemStyleValue> {
        let displayColor: ColorStyleValue | undefined = defaultColor;

        if (color && color !== DEFAULT_ACCOUNT_COLOR) {
            displayColor = `#${color}`;
        }

        const ret: Record<IconItemStyleName, IconItemStyleValue> = {
            color: displayColor
        };

        if (additionalColorAttr) {
            ret[additionalColorAttr] = color;
        }

        if (props.size) {
            ret['font-size'] = props.size;
        }

        return ret;
    }

    function getCategoryIconStyle(color?: ColorValue, defaultColor?: ColorStyleValue, additionalColorAttr?: string): Record<IconItemStyleName, IconItemStyleValue> {
        let displayColor: ColorStyleValue | undefined = defaultColor;

        if (color && color !== DEFAULT_CATEGORY_COLOR) {
            displayColor = `#${color}`;
        }

        const ret: Record<IconItemStyleName, IconItemStyleValue> = {
            color: displayColor
        };

        if (additionalColorAttr) {
            ret[additionalColorAttr] = color;
        }

        if (props.size) {
            ret['font-size'] = props.size;
        }

        return ret;
    }

    function getDefaultIconStyle(color?: ColorValue, defaultColor?: ColorStyleValue, additionalColorAttr?: string): Record<IconItemStyleName, IconItemStyleValue> {
        let displayColor: ColorStyleValue | undefined = defaultColor;

        if (color && color !== DEFAULT_ICON_COLOR) {
            displayColor = `#${color}`;
        }

        const ret: Record<IconItemStyleName, IconItemStyleValue> = {
            color: displayColor
        };

        if (additionalColorAttr) {
            ret[additionalColorAttr] = color;
        }

        if (props.size) {
            ret['font-size'] = props.size;
        }

        return ret;
    }

    return {
        style,
        getAccountIcon,
        getCategoryIcon
    }
}
