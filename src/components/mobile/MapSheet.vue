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
            <div ref="map" style="height: 400px; width: 100%"></div>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
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
            leaflet: null,
            tileLayer: null,
            zoomControl: null,
            attribution: null,
            marker: null,
            initCenter: [ 0, 0 ],
            zoomLevel: 1
        }
    },
    methods: {
        save() {
            this.$emit('update:show', false);
        },
        onSheetOpen() {
            let isFirstInit = false;
            let centerChanged = false;

            if (this.modelValue && (this.modelValue.longitude || this.modelValue.latitude)) {
                if (this.initCenter[0] !== this.modelValue.latitude || this.initCenter[1] !== this.modelValue.longitude) {
                    this.initCenter[0] = this.modelValue.latitude;
                    this.initCenter[1] = this.modelValue.longitude;
                    this.zoomLevel = 14;

                    centerChanged = true;
                }
            } else if (!this.modelValue || (!this.modelValue.longitude && !this.modelValue.latitude)) {
                if (this.initCenter[0] || this.initCenter[1]) {
                    this.initCenter[0] = 0;
                    this.initCenter[1] = 0;
                    this.zoomLevel = 1;

                    centerChanged = true;
                }
            }

            if (!this.leaflet) {
                const mapContainer = this.$refs.map;

                this.leaflet = this.$map.leaflet.map(mapContainer, {
                    attributionControl: false,
                    zoomControl: false
                });

                this.tileLayer = this.$map.leaflet.tileLayer(this.$map.generateOpenStreetMapTileImageUrl(), {
                    maxZoom: 19
                });
                this.tileLayer.addTo(this.leaflet);

                this.zoomControl = this.$map.leaflet.control.zoom({
                    zoomInTitle: this.$t('Zoom in'),
                    zoomOutTitle: this.$t('Zoom out'),
                });
                this.zoomControl.addTo(this.leaflet);

                this.attribution = this.$map.leaflet.control.attribution({
                    prefix: false
                });
                this.attribution.addAttribution('&copy; <a href="http://www.openstreetmap.org/copyright" class="external" target="_blank">OpenStreetMap</a>');
                this.attribution.addTo(this.leaflet);

                isFirstInit = true;
            }

            if (isFirstInit || centerChanged) {
                this.leaflet.setView(this.initCenter, this.zoomLevel);
            }

            if (centerChanged && this.zoomLevel > 1) {
                if (!this.marker) {
                    const markerIcon = this.$map.leaflet.icon({
                        iconUrl: 'img/map-marker-icon.png',
                        iconRetinaUrl: 'img/map-marker-icon-2x.png',
                        iconSize:    [25, 32],
                        iconAnchor:  [12, 32],
                        shadowUrl: 'img/map-marker-shadow.png',
                        shadowSize: [41, 32]
                    });
                    this.marker = this.$map.leaflet.marker(this.initCenter, {
                        icon: markerIcon
                    });
                    this.marker.addTo(this.leaflet);
                } else {
                    this.marker.setLatLng(this.initCenter);
                }
            } else if (centerChanged && this.zoomLevel <= 1) {
                if (this.marker) {
                    this.marker.remove();
                    this.marker = null;
                }
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
