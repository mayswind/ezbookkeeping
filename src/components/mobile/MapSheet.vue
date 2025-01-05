<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left"></div>
            <div class="right">
                <f7-link :text="$t('Done')" @click="save"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content class="no-margin-vertical no-padding-vertical">
            <map-view ref="map" height="400px" :geo-location="geoLocation">
                <template #error-title="{ mapSupported, mapDependencyLoaded }">
                    <div class="display-flex padding justify-content-space-between align-items-center">
                        <div class="ebk-sheet-title" v-if="!mapSupported"><b>{{ $t('Unsupported Map Provider') }}</b></div>
                        <div class="ebk-sheet-title" v-else-if="!mapDependencyLoaded"><b>{{ $t('Cannot Initialize Map') }}</b></div>
                        <div class="ebk-sheet-title" v-else></div>
                    </div>
                </template>
                <template #error-content>
                    <div class="padding-horizontal padding-bottom">
                        <p class="no-margin">{{ $t('Please refresh the page and try again. If the error persists, ensure that the server\'s map settings are correctly configured.') }}</p>
                        <div class="margin-top text-align-center">
                            <f7-link @click="close" :text="$t('Close')"></f7-link>
                        </div>
                    </div>
                </template>
            </map-view>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { computed, useTemplateRef } from 'vue';
import type MapView from '@/components/common/MapView.vue';

import type { MapPosition } from '@/lib/map/base.ts';

const props = defineProps<{
    modelValue?: MapPosition;
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: MapPosition | undefined): void,
    (e: 'update:show', value: boolean): void
}>();

const map = useTemplateRef<MapView>('map');

const geoLocation = computed<MapPosition | undefined>({
    get: () => {
        return props.modelValue;
    },
    set: value => {
        emit('update:modelValue', value);
    }
});

function save() {
    emit('update:show', false);
}

function close() {
    emit('update:show', false);
}

function onSheetOpen() {
    if (map.value) {
        map.value.initMapView();
    }
}

function onSheetClosed() {
    close();
}
</script>
