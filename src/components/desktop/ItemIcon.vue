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
import iconConstants from '@/consts/icon.js';
import colorConstants from '@/consts/color.js';
import { isNumber } from '@/lib/common.js';

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

            if (!iconConstants.allAccountIcons[iconId]) {
                return iconConstants.defaultAccountIcon.icon;
            }

            return iconConstants.allAccountIcons[iconId].icon;
        },
        getCategoryIcon(iconId) {
            if (isNumber(iconId)) {
                iconId = iconId.toString();
            }

            if (!iconConstants.allCategoryIcons[iconId]) {
                return iconConstants.defaultCategoryIcon.icon;
            }

            return iconConstants.allCategoryIcons[iconId].icon;
        },
        getAccountIconStyle(color, defaultColor, additionalColorAttr) {
            if (color && color !== colorConstants.defaultAccountColor) {
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
            if (color && color !== colorConstants.defaultCategoryColor) {
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
            if (color && color !== colorConstants.defaultColor) {
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
