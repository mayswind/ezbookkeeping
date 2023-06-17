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
        <f7-page-content class="no-margin-vertical no-padding-vertical" v-if="mapHolder && dependencyLoaded">
            <div ref="map" style="height: 400px; width: 100%"></div>
        </f7-page-content>
        <f7-page-content class="no-margin-top no-padding-top" v-else-if="!mapHolder || !dependencyLoaded">
            <div class="display-flex padding justify-content-space-between align-items-center">
                <div class="ebk-sheet-title"><b>{{ mapErrorTitle }}</b></div>
            </div>
            <div class="padding-horizontal padding-bottom">
                <p class="no-margin">{{ mapErrorContent }}</p>
                <div class="margin-top text-align-center">
                    <f7-link @click="close" :text="$t('Close')"></f7-link>
                </div>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
import {
    createMapHolder,
    initMapInstance,
    setMapCenterTo,
    setMapCenterMarker,
    removeMapCenterMarker
} from '@/lib/map.js';

export default {
    props: [
        'modelValue',
        'show'
    ],
    emits: [
        'update:modelValue',
        'update:show'
    ],
    data() {
        return {
            mapHolder: createMapHolder(),
            initCenter: {
                latitude: 0,
                longitude: 0,
            },
            zoomLevel: 1
        }
    },
    computed: {
        dependencyLoaded() {
            return this.mapHolder && this.mapHolder.dependencyLoaded;
        },
        mapErrorTitle() {
            if (!this.mapHolder) {
                return this.$t('Unsupported Map Provider');
            }

            if (!this.dependencyLoaded) {
                return this.$t('Cannot Initialize Map');
            }

            return '';
        },
        mapErrorContent() {
            return this.$t('Please refresh the page and try again. If the error is still displayed, make sure that server map settings are set correctly.');
        }
    },
    methods: {
        save() {
            this.$emit('update:show', false);
        },
        onSheetOpen() {
            let isFirstInit = false;
            let centerChanged = false;

            if (!this.mapHolder || !this.mapHolder.dependencyLoaded) {
                return;
            }

            if (this.modelValue && (this.modelValue.longitude || this.modelValue.latitude)) {
                if (this.initCenter.latitude !== this.modelValue.latitude || this.initCenter.longitude !== this.modelValue.longitude) {
                    this.initCenter.latitude = this.modelValue.latitude;
                    this.initCenter.longitude = this.modelValue.longitude;
                    this.zoomLevel = this.mapHolder.defaultZoomLevel;

                    centerChanged = true;
                }
            } else if (!this.modelValue || (!this.modelValue.longitude && !this.modelValue.latitude)) {
                if (this.initCenter.latitude || this.initCenter.longitude) {
                    this.initCenter.latitude = 0;
                    this.initCenter.longitude = 0;
                    this.zoomLevel = this.mapHolder.minZoomLevel;

                    centerChanged = true;
                }
            }

            if (!this.mapHolder.inited) {
                initMapInstance(this.mapHolder, this.$refs.map, {
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
        },
        onSheetClosed() {
            this.close();
        },
        close() {
            this.$emit('update:show', false);
        }
    }
}
</script>
