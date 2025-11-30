<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" class="map-sheet" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link icon-f7="minus" :class="{ 'disabled': !map?.allowZoomOut() }" @click="map?.zoomOut()"></f7-link>
                <f7-link icon-f7="plus" :class="{ 'disabled': !map?.allowZoomIn() }" @click="map?.zoomIn()"></f7-link>
            </div>
            <div class="right map-sheet-toolbar-right">
                <f7-link :text="tt('Disable Click to Set Location')" @click="switchSetGeoLocationByClickMap(false)" v-if="!readonly && isSupportGetGeoLocationByClick() && props.setGeoLocationByClickMap"></f7-link>
                <f7-link class="map-sheet-toolbar-auto-hidden" :text="tt('Enable Click to Set Location')" @click="switchSetGeoLocationByClickMap(true)" v-if="!readonly && isSupportGetGeoLocationByClick() && !props.setGeoLocationByClickMap"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content class="no-margin no-padding">
            <map-view ref="map" height="400px"
                      :enable-zoom-control="false" :geo-location="geoLocation"
                      @click="updateSpecifiedGeoLocation">
                <template #error-title="{ mapSupported, mapDependencyLoaded }">
                    <div class="display-flex map-sheet-error-title padding justify-content-space-between align-items-center">
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

import type { Coordinate } from '@/core/coordinate.ts';

import { isSupportGetGeoLocationByClick } from '@/lib/map/index.ts';

type MapViewType = InstanceType<typeof MapView>;

const props = defineProps<{
    modelValue?: Coordinate;
    readonly?: boolean;
    setGeoLocationByClickMap?: boolean;
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: Coordinate | undefined): void;
    (e: 'update:setGeoLocationByClickMap', value: boolean): void;
    (e: 'update:show', value: boolean): void;
}>();

const { tt } = useI18n();

const map = useTemplateRef<MapViewType>('map');

const geoLocation = computed<Coordinate | undefined>({
    get: () => {
        return props.modelValue;
    },
    set: value => {
        emit('update:modelValue', value);
    }
});

function updateSpecifiedGeoLocation(coordinate: Coordinate): void {
    if (!props.readonly && isSupportGetGeoLocationByClick() && props.setGeoLocationByClickMap) {
        geoLocation.value = coordinate;
        map.value?.setMarkerPosition(coordinate);
    }
}

function switchSetGeoLocationByClickMap(value: boolean): void {
    emit('update:setGeoLocationByClickMap', value);
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

<style>
.map-sheet .map-sheet-error-title {
    margin-top: var(--f7-toolbar-height);
}

.map-sheet .map-sheet-toolbar-right {
    overflow: hidden;

    .map-sheet-toolbar-auto-hidden {
        overflow: hidden;

        > span {
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
        }
    }
}
</style>
