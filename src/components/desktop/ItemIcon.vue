<template>
    <i class="item-icon" :class="classes" :style="style" v-if="!hiddenStatus">
        <slot></slot>
    </i>
    <v-badge class="right-bottom-icon" color="secondary"
             location="bottom right" offset-y="4" :icon="icons.hide"
             v-if="hiddenStatus">
        <i class="item-icon" :class="classes" :style="style">
            <slot></slot>
        </i>
    </v-badge>
</template>

<script>
import { ALL_ACCOUNT_ICONS, DEFAULT_ACCOUNT_ICON, ALL_CATEGORY_ICONS, DEFAULT_CATEGORY_ICON } from '@/consts/icon.ts';
import { DEFAULT_ICON_COLOR, DEFAULT_ACCOUNT_COLOR, DEFAULT_CATEGORY_COLOR } from '@/consts/color.ts';
import { isNumber } from '@/lib/common.ts';

import {
    mdiEyeOffOutline
} from '@mdi/js';

export default {
    props: [
        'class',
        'iconType',
        'iconId',
        'color',
        'defaultColor',
        'additionalColorAttr',
        'size',
        'hiddenStatus'
    ],
    data() {
        return {
            icons: {
                hide: mdiEyeOffOutline
            }
        }
    },
    computed: {
        classes() {
            let allClasses = this.class ? (this.class + ' ') : '';

            if (this.iconType === 'account') {
                allClasses += this.getAccountIcon(this.iconId);
            } else if (this.iconType === 'category') {
                allClasses += this.getCategoryIcon(this.iconId);
            } else if (this.iconType === 'fixed') {
                allClasses += this.iconId;
            }

            return allClasses;
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

            if (this.size) {
                ret['font-size'] = this.size;
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

            if (this.size) {
                ret['font-size'] = this.size;
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

            if (this.size) {
                ret['font-size'] = this.size;
            }

            return ret;
        }
    }
}
</script>

<style>
.item-icon {
    font-size: var(--ebk-icon-font-size);
    display: inline-block;
    vertical-align: middle;
    background-size: 100% auto;
    background-position: center;
    background-repeat: no-repeat;
    font-style: normal;
    position: relative;
}
</style>
