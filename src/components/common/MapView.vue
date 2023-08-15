<template>
    <div ref="mapContainer" style="width: 100%" :style="{ 'height': mapHeight }"></div>
    <slot name="error-title"
          :mapSupported="mapSupported" :mapDependencyLoaded="mapDependencyLoaded"
          v-if="!mapSupported || !mapDependencyLoaded"></slot>
    <slot name="error-content"
          :mapSupported="mapSupported" :mapDependencyLoaded="mapDependencyLoaded"
          v-if="!mapSupported || !mapDependencyLoaded"></slot>
</template>

<script>
import {
    createMapHolder,
    initMapInstance,
    setMapCenterTo,
    setMapCenterMarker,
    removeMapCenterMarker
} from '@/lib/map/index.js';

export default {
    props: [
        'height',
        'geoLocation'
    ],
    expose: [
        'init'
    ],
    data() {
        this.mapHolder = createMapHolder();

        return {
            mapSupported: !!this.mapHolder,
            mapDependencyLoaded: this.mapHolder && this.mapHolder.dependencyLoaded,
            mapInited: false,
            initCenter: {
                latitude: 0,
                longitude: 0,
            },
            zoomLevel: 1
        }
    },
    computed: {
        mapHeight() {
            if (this.mapSupported && this.mapDependencyLoaded) {
                return this.height;
            }

            return '0';
        }
    },
    methods: {
        init() {
            let isFirstInit = false;
            let centerChanged = false;

            if (!this.mapSupported || !this.mapDependencyLoaded) {
                return;
            }

            if (this.geoLocation && (this.geoLocation.longitude || this.geoLocation.latitude)) {
                if (this.initCenter.latitude !== this.geoLocation.latitude || this.initCenter.longitude !== this.geoLocation.longitude) {
                    this.initCenter.latitude = this.geoLocation.latitude;
                    this.initCenter.longitude = this.geoLocation.longitude;
                    this.zoomLevel = this.mapHolder.defaultZoomLevel;

                    centerChanged = true;
                }
            } else if (!this.geoLocation || (!this.geoLocation.longitude && !this.geoLocation.latitude)) {
                if (this.initCenter.latitude || this.initCenter.longitude) {
                    this.initCenter.latitude = 0;
                    this.initCenter.longitude = 0;
                    this.zoomLevel = this.mapHolder.minZoomLevel;

                    centerChanged = true;
                }
            }

            if (!this.mapHolder.inited) {
                const languageInfo = this.$locale.getCurrentLanguageInfo();

                initMapInstance(this.mapHolder, this.$refs.mapContainer, {
                    language: languageInfo ? languageInfo.code : null,
                    initCenter: this.initCenter,
                    zoomLevel: this.zoomLevel,
                    text: {
                        zoomIn: this.$t('Zoom in'),
                        zoomOut: this.$t('Zoom out'),
                    }
                });

                if (this.mapHolder.inited) {
                    isFirstInit = true;
                }
            }

            if (isFirstInit || centerChanged) {
                setMapCenterTo(this.mapHolder, this.initCenter, this.zoomLevel);
            }

            if (centerChanged && this.zoomLevel > this.mapHolder.minZoomLevel) {
                setMapCenterMarker(this.mapHolder, this.initCenter);
            } else if (centerChanged && this.zoomLevel <= this.mapHolder.minZoomLevel) {
                removeMapCenterMarker(this.mapHolder);
            }
        }
    }
}
</script>
