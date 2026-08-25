import type {
    ColorValue,
    ColorStyleValue,
    ColorInfo
} from '@/core/color.ts';

import {
    DEFAULT_ICON_COLOR,
    DEFAULT_ACCOUNT_COLOR,
    DEFAULT_CATEGORY_COLOR,
    DEFAULT_COLOR_STYLE_VARIABLE,
    DEFAULT_MOBILE_OVERVIEW_WIDGET_LIGHT_BACKGROUND_COLOR,
    DEFAULT_MOBILE_OVERVIEW_WIDGET_DARK_BACKGROUND_COLOR
} from '@/consts/color.ts';

export function getColorsInRows(allColorValues: ColorValue[], itemPerRow: number): ColorInfo[][] {
    const ret: ColorInfo[][] = [];
    let rowCount = -1;

    for (let i = 0; i < allColorValues.length; i++) {
        if (i % itemPerRow === 0) {
            ret[++rowCount] = [];
        }

        ret[rowCount]!.push({
            color: allColorValues[i] as ColorValue
        });
    }

    return ret;
}

export function getDisplayColor(color?: ColorValue): ColorStyleValue {
    if (color && color !== DEFAULT_ICON_COLOR) {
        return `#${color}`;
    } else {
        return DEFAULT_COLOR_STYLE_VARIABLE;
    }
}

export function getCategoryDisplayColor(color?: ColorValue): ColorStyleValue {
    if (color && color !== DEFAULT_CATEGORY_COLOR) {
        return `#${color}`;
    } else {
        return DEFAULT_COLOR_STYLE_VARIABLE;
    }
}
export function getAccountDisplayColor(color?: ColorValue): ColorStyleValue {
    if (color && color !== DEFAULT_ACCOUNT_COLOR) {
        return `#${color}`;
    } else {
        return DEFAULT_COLOR_STYLE_VARIABLE;
    }
}

export function getContrastTextColor(backgroundColor: ColorValue): ColorValue {
    const normalizedColor = backgroundColor.replace(/^#/, '');

    if (!/^[0-9a-fA-F]{6}$/.test(normalizedColor)) {
        return '000000';
    }

    const rgb = [0, 2, 4].map(offset => parseInt(normalizedColor.substring(offset, offset + 2), 16) / 255);
    const linearRgb = rgb.map(value => value <= 0.04045 ? value / 12.92 : Math.pow((value + 0.055) / 1.055, 2.4));
    const relativeLuminance = 0.2126 * linearRgb[0]! + 0.7152 * linearRgb[1]! + 0.0722 * linearRgb[2]!;

    return relativeLuminance > 0.179 ? '000000' : 'ffffff';
}

export function getContrastIconColor(backgroundColor: ColorValue): ColorValue {
    const normalizedColor = backgroundColor.replace(/^#/, '');

    if (!/^[0-9a-fA-F]{6}$/.test(normalizedColor) || normalizedColor.toLowerCase() === DEFAULT_MOBILE_OVERVIEW_WIDGET_LIGHT_BACKGROUND_COLOR) {
        return 'c67e48';
    } else if (normalizedColor.toLowerCase() === DEFAULT_MOBILE_OVERVIEW_WIDGET_DARK_BACKGROUND_COLOR) {
        return 'ffffff99';
    }

    const useDarkIcon = getContrastTextColor(normalizedColor) === '000000';
    const targetChannel = useDarkIcon ? 0 : 255;
    const targetRatio = useDarkIcon ? 0.32 : 1;

    return [0, 2, 4].map(offset => {
        const channel = parseInt(normalizedColor.substring(offset, offset + 2), 16);
        return Math.round(channel * (1 - targetRatio) + targetChannel * targetRatio).toString(16).padStart(2, '0');
    }).join('') + (useDarkIcon ? '' : '99');
}
