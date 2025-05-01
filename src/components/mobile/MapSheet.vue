<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link :text="tt('Disable Click to Set Location')" @click="switchSetGeoLocationByClickMap(false)" v-if="isSupportGetGeoLocationByClick() && props.setGeoLocationByClickMap"></f7-link>
                <f7-link :text="tt('Enable Click to Set Location')" @click="switchSetGeoLocationByClickMap(true)" v-if="isSupportGetGeoLocationByClick() && !props.setGeoLocationByClickMap"></f7-link>
            </div>
            <div class="right">
                <f7-link :text="tt('Done')" @click="save"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content class="no-margin-vertical no-padding-vertical">
            <map-view ref="map" height="400px" :geo-location="geoLocation" @click="updateSpecifiedGeoLocation">
                <template #error-title="{ mapSupported, mapDependencyLoaded }">
                    <div class="display-flex padding justify-content-space-between align-items-center">
                        <div class="ebk-sheet-title" v-if="!mapSupported"><b>{{ tt('Unsupported Map Provider') }}</b></div>
                        <div class="ebk-sheet-title" v-else-if="!mapDependencyLoaded"><b>{{ tt('Cannot Initialize Map') }}</b></div>
                        <div class="ebk-sheet-title" v-else></div>
                    </div>
                </template>
                <template #error-content>
                    <div class="padding-horizontal padding-bottom">
                        <p class="no-margin">{{ tt('Please refresh the page and try again. If the error persists, ensure that the server\'s map settings are correctly configured.') }}</p>
                        <div class="margin-top text-align-center">
                            <f7-link @click="close" :text="tt('Close')"></f7-link>
                        </div>
                    </div>
                </template>
            </map-view>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { computed, useTemplateRef } from 'vue';
import MapView from '@/components/common/MapView.vue';

import { useI18n } from '@/locales/helpers.ts';

import type { MapPosition } from '@/core/map.ts';

import { isSupportGetGeoLocationByClick } from '@/lib/map/index.ts';

type MapViewType = InstanceType<typeof MapView>;

const props = defineProps<{
    modelValue?: MapPosition;
    setGeoLocationByClickMap?: boolean;
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: MapPosition | undefined): void;
    (e: 'update:setGeoLocationByClickMap', value: boolean): void;
    (e: 'update:show', value: boolean): void;
}>();

const { tt } = useI18n();

const map = useTemplateRef<MapViewType>('map');

const geoLocation = computed<MapPosition | undefined>({
    get: () => {
        return props.modelValue;
    },
    set: value => {
        emit('update:modelValue', value);
    }
});

function updateSpecifiedGeoLocation(mapPosition: MapPosition): void {
    if (isSupportGetGeoLocationByClick() && props.setGeoLocationByClickMap) {
        geoLocation.value = mapPosition;
        map.value?.setMarkerPosition(mapPosition);
    }
}

function switchSetGeoLocationByClickMap(value: boolean): void {
    emit('update:setGeoLocationByClickMap', value);
}

function save(): void {
    emit('update:show', false);
}

function close(): void {
    emit('update:show', false);
}

function onSheetOpen(): void {
    if (map.value) {
        map.value.initMapView();
    }
}

function onSheetClosed(): void {
    close();
}
</script>
