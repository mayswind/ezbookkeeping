<template>
    <div ref="mapContainer" :class="'map-view-container' + (mapClass ? ` ${mapClass}` : '')" :style="finalMapStyle"></div>
    <slot name="error-title"
          :mapSupported="mapSupported" :mapDependencyLoaded="mapDependencyLoaded"
          v-if="!mapSupported || !mapDependencyLoaded"></slot>
    <slot name="error-content"
          :mapSupported="mapSupported" :mapDependencyLoaded="mapDependencyLoaded"
          v-if="!mapSupported || !mapDependencyLoaded"></slot>
</template>

<script setup lang="ts">
import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import type { Coordinate } from '@/core/coordinate.ts';
import { isNumber } from '@/lib/common.ts';
import type { MapInstance } from '@/lib/map/base.ts';
import { createMapInstance } from '@/lib/map/index.ts';

const props = defineProps<{
    height?: string;
    mapClass?: string;
    mapStyle?: Record<string, string>;
    enableZoomControl?: boolean;
    geoLocation?: Coordinate;
}>();

const emit = defineEmits<{
    (e: 'click', geoLocation: Coordinate): void;
}>();

const { tt, getCurrentLanguageInfo } = useI18n();

const mapContainer = useTemplateRef<HTMLElement>('mapContainer');
const mapInstance = ref<MapInstance | null>(createMapInstance({
    enableZoomControl: props.enableZoomControl
}));
const initCenter = ref<Coordinate>({
    latitude: 0,
    longitude: 0
});
const zoomLevel = ref<number>(1);

const mapSupported = computed<boolean>(() => !!mapInstance.value);
const mapDependencyLoaded = computed<boolean>(() => mapInstance.value?.dependencyLoaded || false);

const finalMapStyle = computed<Record<string, string>>(() => {
    const styles: Record<string, string> = Object.assign({}, props.mapStyle);

    if (props.height) {
        styles['height'] = props.height;
    }

    if (!mapSupported.value || !mapDependencyLoaded.value) {
        styles['height'] = '0';
    }

    return styles;
});

function initMapView(): void {
    let isFirstInit = false;
    let centerChanged = false;

    if (!mapSupported.value || !mapDependencyLoaded.value || !mapInstance.value) {
        return;
    }

    if (props.geoLocation && (props.geoLocation.longitude || props.geoLocation.latitude)) {
        if (initCenter.value.latitude !== props.geoLocation.latitude || initCenter.value.longitude !== props.geoLocation.longitude) {
            initCenter.value.latitude = props.geoLocation.latitude;
            initCenter.value.longitude = props.geoLocation.longitude;
            zoomLevel.value = mapInstance.value.getDefaultZoomLevel();

            centerChanged = true;
        }
    } else if (!props.geoLocation || (!props.geoLocation.longitude && !props.geoLocation.latitude)) {
        if (initCenter.value.latitude || initCenter.value.longitude) {
            initCenter.value.latitude = 0;
            initCenter.value.longitude = 0;
            zoomLevel.value = mapInstance.value.getMinZoomLevel();

            centerChanged = true;
        }
    }

    if (!mapInstance.value.inited) {
        const languageInfo = getCurrentLanguageInfo();

        mapInstance.value.initMapInstance(mapContainer.value as HTMLElement, {
            language: languageInfo?.alternativeLanguageTag,
            initCenter: initCenter.value,
            zoomLevel: zoomLevel.value,
            text: {
                zoomIn: tt('Zoom in'),
                zoomOut: tt('Zoom out'),
            },
            onClick: (geoLocation: Coordinate) => {
                emit('click', geoLocation);
            },
            onZoomChange(level: number) {
                if (isNumber(level)) {
                    zoomLevel.value = level;
                } else if (mapInstance.value) {
                    zoomLevel.value = Math.round(mapInstance.value.getZoomLevel());
                }
            },
        });

        if (mapInstance.value.inited) {
            isFirstInit = true;
        }
    }

    if (isFirstInit || centerChanged) {
        mapInstance.value.setMapCenterTo(initCenter.value, zoomLevel.value);
    }

    if (centerChanged && zoomLevel.value > mapInstance.value.getMinZoomLevel()) {
        mapInstance.value.setMapCenterMarker(initCenter.value);
    } else if (centerChanged && zoomLevel.value <= mapInstance.value.getMinZoomLevel()) {
        mapInstance.value.removeMapCenterMarker();
    }
}

function setMarkerPosition(geoLocation?: Coordinate): void {
    if (!mapInstance.value) {
        return;
    }

    if (geoLocation) {
        mapInstance.value.setMapCenterMarker(geoLocation);
    }
}

function allowZoomIn(): boolean {
    if (!mapSupported.value || !mapDependencyLoaded.value || !mapInstance.value) {
        return false;
    }

    return zoomLevel.value < mapInstance.value.getMaxZoomLevel();
}

function allowZoomOut(): boolean {
    if (!mapSupported.value || !mapDependencyLoaded.value || !mapInstance.value) {
        return false;
    }

    return zoomLevel.value > mapInstance.value.getMinZoomLevel();
}

function zoomIn(): void {
    if (!mapInstance.value) {
        return;
    }

    mapInstance.value.zoomIn();
}

function zoomOut(): void {
    if (!mapInstance.value) {
        return;
    }

    mapInstance.value.zoomOut();
}

defineExpose({
    initMapView,
    setMarkerPosition,
    allowZoomIn,
    allowZoomOut,
    zoomIn,
    zoomOut
});
</script>

<style>
.map-view-container {
    width: 100%;
}
</style>
