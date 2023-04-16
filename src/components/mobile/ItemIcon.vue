<template>
    <f7-icon :f7="f7Icon" :icon="icon" :style="style">
        <slot></slot>
    </f7-icon>
</template>

<script>
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
            if (this.$utilities.isNumber(iconId)) {
                iconId = iconId.toString();
            }

            if (!this.$constants.icons.allAccountIcons[iconId]) {
                return this.$constants.icons.defaultAccountIcon.icon;
            }

            return this.$constants.icons.allAccountIcons[iconId].icon;
        },
        getCategoryIcon(iconId) {
            if (this.$utilities.isNumber(iconId)) {
                iconId = iconId.toString();
            }

            if (!this.$constants.icons.allCategoryIcons[iconId]) {
                return this.$constants.icons.defaultCategoryIcon.icon;
            }

            return this.$constants.icons.allCategoryIcons[iconId].icon;
        },
        getAccountIconStyle(color, defaultColor, additionalColorAttr) {
            if (color && color !== this.$constants.colors.defaultAccountColor) {
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
            if (color && color !== this.$constants.colors.defaultCategoryColor) {
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
            if (color && color !== this.$constants.colors.defaultColor) {
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
