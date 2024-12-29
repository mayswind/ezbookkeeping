<template>
    <f7-icon :f7="f7Icon" :icon="icon" :style="style">
        <slot></slot>
    </f7-icon>
</template>

<script>
import { ALL_ACCOUNT_ICONS, DEFAULT_ACCOUNT_ICON, ALL_CATEGORY_ICONS, DEFAULT_CATEGORY_ICON } from '@/consts/icon.ts';
import { DEFAULT_ICON_COLOR, DEFAULT_ACCOUNT_COLOR, DEFAULT_CATEGORY_COLOR } from '@/consts/color.ts';
import { isNumber } from '@/lib/common.ts';

export default {
    props: [
        'iconType',
        'iconId',
        'color',
        'defaultColor',
        'additionalColorAttr'
    ],
    computed: {
        f7Icon() {
            if (this.iconType === 'fixed-f7') {
                return this.iconId;
            } else {
                return '';
            }
        },
        icon() {
            if (this.iconType === 'account') {
                return this.getAccountIcon(this.iconId);
            } else if (this.iconType === 'category') {
                return this.getCategoryIcon(this.iconId);
            } else if (this.iconType === 'fixed') {
                return this.iconId;
            } else {
                return '';
            }
        },
        style() {
            let defaultColor = 'var(--default-icon-color)';

            if (this.defaultColor) {
                defaultColor = this.defaultColor;
            }

            if (this.iconType === 'account') {
                return this.getAccountIconStyle(this.color, defaultColor, this.additionalColorAttr);
            } else if (this.iconType === 'category') {
                return this.getCategoryIconStyle(this.color, defaultColor, this.additionalColorAttr);
            } else {
                return this.getDefaultIconStyle(this.color, defaultColor, this.additionalColorAttr);
            }
        }
    },
    methods: {
        getAccountIcon(iconId) {
            if (isNumber(iconId)) {
                iconId = iconId.toString();
            }

            if (!ALL_ACCOUNT_ICONS[iconId]) {
                return DEFAULT_ACCOUNT_ICON.icon;
            }

            return ALL_ACCOUNT_ICONS[iconId].icon;
        },
        getCategoryIcon(iconId) {
            if (isNumber(iconId)) {
                iconId = iconId.toString();
            }

            if (!ALL_CATEGORY_ICONS[iconId]) {
                return DEFAULT_CATEGORY_ICON.icon;
            }

            return ALL_CATEGORY_ICONS[iconId].icon;
        },
        getAccountIconStyle(color, defaultColor, additionalColorAttr) {
            if (color && color !== DEFAULT_ACCOUNT_COLOR) {
                color = '#' + color;
            } else {
                color = defaultColor;
            }

            const ret = {
                color: color
            };

            if (additionalColorAttr) {
                ret[additionalColorAttr] = color;
            }

            return ret;
        },
        getCategoryIconStyle(color, defaultColor, additionalColorAttr) {
            if (color && color !== DEFAULT_CATEGORY_COLOR) {
                color = '#' + color;
            } else {
                color = defaultColor;
            }

            const ret = {
                color: color
            };

            if (additionalColorAttr) {
                ret[additionalColorAttr] = color;
            }

            return ret;
        },
        getDefaultIconStyle(color, defaultColor, additionalColorAttr) {
            if (color && color !== DEFAULT_ICON_COLOR) {
                color = '#' + color;
            } else {
                color = defaultColor;
            }

            const ret = {
                color: color
            };

            if (additionalColorAttr) {
                ret[additionalColorAttr] = color;
            }

            return ret;
        }
    }
}
</script>
