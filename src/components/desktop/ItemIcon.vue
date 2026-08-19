<template>
    <div class="item-img-icon-container" :class="classes" :style="style" v-if="!hiddenStatus && customIconUrl">
        <img class="item-img-icon" :src="customIconUrl" />
    </div>
    <i class="item-icon" :class="classes" :style="style" v-else-if="!hiddenStatus && !customIconUrl"></i>
    <v-badge class="right-bottom-icon" color="secondary" offset-y="4"
             :location="`bottom ${textDirection === TextDirection.LTR ? 'right' : 'left'}`"
             :icon="mdiEyeOffOutline" v-if="hiddenStatus">
        <div class="item-img-icon-container" :class="classes" :style="style" v-if="customIconUrl">
            <img class="item-img-icon" :src="customIconUrl" />
        </div>
        <i class="item-icon" :class="classes" :style="style" v-else-if="!customIconUrl"></i>
    </v-badge>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { type CommonIconProps, useItemIconBase } from '@/components/base/ItemIconBase.ts';

import { useI18n } from '@/locales/helpers.ts';

import { TextDirection } from '@/core/text.ts';

import {
    mdiEyeOffOutline
} from '@mdi/js';

interface DesktopItemIconProps extends CommonIconProps {
    class?: string;
    hiddenStatus?: boolean;
}

const props = defineProps<DesktopItemIconProps>();

const { getCurrentLanguageTextDirection } = useI18n();
const { style, customIconUrl, getAccountIcon, getCategoryIcon } = useItemIconBase(props);

const textDirection = computed<TextDirection>(() => getCurrentLanguageTextDirection());

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
.item-icon,
.item-img-icon,
.item-img-icon-container {
    display: inline-flex;
    vertical-align: middle;
    background-size: 100% auto;
    background-position: center;
    background-repeat: no-repeat;
    font-style: normal;
    position: relative;
}

.item-icon {
    font-size: var(--ebk-icon-font-size);
}

.item-img-icon {
    width: var(--ebk-icon-font-size);
    height: var(--ebk-icon-font-size);
}
</style>
