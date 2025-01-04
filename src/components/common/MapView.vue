<template>
    <div ref="mapContainer" style="width: 100%" :class="mapClass" :style="finalMapStyle"></div>
    <slot name="error-title"
          :mapSupported="mapSupported" :mapDependencyLoaded="mapDependencyLoaded"
          v-if="!mapSupported || !mapDependencyLoaded"></slot>
    <slot name="error-content"
          :mapSupported="mapSupported" :mapDependencyLoaded="mapDependencyLoaded"
          v-if="!mapSupported || !mapDependencyLoaded"></slot>
</template>

<script setup lang="ts">
import { type Ref, ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/lib/i18n.js';

import { copyObjectTo } from '@/lib/common.ts';
import type { MapInstance, MapPosition } from '@/lib/map/base.ts';
import { createMapInstance } from '@/lib/map/index.ts';

const props = defineProps<{
    height?: string;
    mapClass?: string;
    mapStyle?: Record<string, string>;
    geoLocation?: MapPosition;
}>();

const { tt, getCurrentLanguageInfo } = useI18n();

const mapContainer: Ref<HTMLElement> = useTemplateRef('mapContainer');
const mapInstance: Ref<MapInstance> = ref(createMapInstance());
const initCenter: Ref<MapPosition> = ref({
    latitude: 0,
    longitude: 0
});
const zoomLevel: Ref<number> = ref(1);

const mapSupported = computed<boolean>(() => {
    return !!mapInstance.value;
});

const mapDependencyLoaded = computed<boolean>(() => {
    return mapInstance.value && mapInstance.value.dependencyLoaded;
});

const finalMapStyle = computed<Record<string, string>>(() => {
    const styles = copyObjectTo(props.mapStyle, {});

    if (props.height) {
        styles.height = props.height;
    }

    if (!mapSupported.value || !mapDependencyLoaded.value) {
        styles.height = '0';
    }

    return styles;
});

function init() {
    let isFirstInit = false;
    let centerChanged = false;

    if (!mapSupported.value || !mapDependencyLoaded.value) {
        return;
    }

    if (props.geoLocation && (props.geoLocation.longitude || props.geoLocation.latitude)) {
        if (initCenter.value.latitude !== props.geoLocation.latitude || initCenter.value.longitude !== props.geoLocation.longitude) {
            initCenter.value.latitude = props.geoLocation.latitude;
            initCenter.value.longitude = props.geoLocation.longitude;
            zoomLevel.value = mapInstance.value.defaultZoomLevel;

            centerChanged = true;
        }
    } else if (!props.geoLocation || (!props.geoLocation.longitude && !props.geoLocation.latitude)) {
        if (initCenter.value.latitude || initCenter.value.longitude) {
            initCenter.value.latitude = 0;
            initCenter.value.longitude = 0;
            zoomLevel.value = mapInstance.value.minZoomLevel;

            centerChanged = true;
        }
    }

    if (!mapInstance.value.inited) {
        const languageInfo = getCurrentLanguageInfo();

        mapInstance.value.initMapInstance(mapContainer.value, {
            language: languageInfo?.alternativeLanguageTag,
            initCenter: initCenter.value,
            zoomLevel: zoomLevel.value,
            text: {
                zoomIn: tt('Zoom in'),
                zoomOut: tt('Zoom out'),
            }
        });

        if (mapInstance.value.inited) {
            isFirstInit = true;
        }
    }

    if (isFirstInit || centerChanged) {
        mapInstance.value.setMapCenterTo(initCenter.value, zoomLevel.value);
    }

    if (centerChanged && zoomLevel.value > mapInstance.value.minZoomLevel) {
        mapInstance.value.setMapCenterMarker(initCenter.value);
    } else if (centerChanged && zoomLevel.value <= mapInstance.value.minZoomLevel) {
        mapInstance.value.removeMapCenterMarker();
    }
}

defineExpose({
    init
});
</script>
