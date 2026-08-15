<template>
    <div class="icon item-img-icon-container" :style="style" v-if="customIconUrl">
        <img class="item-img-icon" :src="customIconUrl" />
        <slot></slot>
    </div>
    <f7-icon :f7="f7IconValue" :icon="icon" :style="style" v-else-if="!customIconUrl">
        <slot></slot>
    </f7-icon>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { type CommonIconProps, useItemIconBase } from '@/components/base/ItemIconBase.ts';

const props = defineProps<CommonIconProps>();
const { style, customIconUrl, getAccountIcon, getCategoryIcon } = useItemIconBase(props);

const f7IconValue = computed<string>(() => {
    if (props.iconType === 'fixed-f7') {
        return props.iconId.toString();
    } else {
        return '';
    }
});

const icon = computed<string>(() => {
    if (props.iconType === 'account') {
        return getAccountIcon(props.iconId);
    } else if (props.iconType === 'category') {
        return getCategoryIcon(props.iconId);
    } else if (props.iconType === 'fixed') {
        return props.iconId.toString();
    } else {
        return '';
    }
});
</script>

<style>
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

.item-img-icon {
    width: var(--ebk-icon-font-size);
    height: var(--ebk-icon-font-size);
}
</style>
