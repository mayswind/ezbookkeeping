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

<script setup lang="ts">
import { computed } from 'vue';
import { type CommonIconProps, useItemIcon } from '@/components/base/itemIcon.ts';

import {
    mdiEyeOffOutline
} from '@mdi/js';

interface DesktopItemIconProps extends CommonIconProps {
    class?: string;
    hiddenStatus?: boolean;
}

const props = defineProps<DesktopItemIconProps>();
const { style, getAccountIcon, getCategoryIcon } = useItemIcon(props);

const icons = {
    hide: mdiEyeOffOutline
};

const classes = computed<string>(() => {
    let allClasses = props.class ? (props.class + ' ') : '';

    if (props.iconType === 'account') {
        allClasses += getAccountIcon(props.iconId);
    } else if (props.iconType === 'category') {
        allClasses += getCategoryIcon(props.iconId);
    } else if (props.iconType === 'fixed') {
        allClasses += props.iconId;
    }

    return allClasses;
});
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
