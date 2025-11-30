// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-nocheck
import type { Coordinate } from '@/core/coordinate.ts';

import { type LeafletTileSource, type LeafletTileSourceExtraParam, LEAFLET_TILE_SOURCES } from '@/consts/map.ts';

import type { MapProvider, MapInstance, MapCreateOptions, MapInstanceInitOptions } from './base.ts';

import { isFunction, isNumber } from '@/lib/common.ts';
import {
    isMapDataFetchProxyEnabled,
    getCustomMapTileLayerUrl,
    getCustomMapAnnotationLayerUrl,
    isCustomMapAnnotationLayerDataFetchProxyEnabled,
    getCustomMapMinZoomLevel,
    getCustomMapMaxZoomLevel,
    getCustomMapDefaultZoomLevel,
    getTomTomMapAPIKey,
    getTianDiTuMapAPIKey
} from '@/lib/server_settings.ts';
import services from '@/lib/services.ts';

export class LeafletMapProvider implements MapProvider {
    // https://leafletjs.com/reference.html#pan-options
    public static Leaflet: unknown = null;
    private readonly mapProvider: string;

    public constructor(mapProvider: string) {
        this.mapProvider = mapProvider;
    }

    public getWebsite(): string {
        if (this.mapProvider === 'custom') {
            return '';
        } else if (LEAFLET_TILE_SOURCES[this.mapProvider]) {
            return LEAFLET_TILE_SOURCES[this.mapProvider].website || '';
        } else {
            return '';
        }
    }

    public isSupportGetGeoLocationByClick(): boolean {
        return true;
    }

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    public asyncLoadAssets(language?: string): Promise<unknown> {
        return Promise.all([
            import('leaflet/dist/leaflet.css'),
            import('leaflet/dist/leaflet-src.esm.js').then(leaflet => LeafletMapProvider.Leaflet = leaflet)
        ]);
    }

    public createMapInstance(options: MapCreateOptions): MapInstance | null {
        const mapTileSource = LEAFLET_TILE_SOURCES[this.mapProvider];

        if (this.mapProvider !== 'custom' && !mapTileSource) {
            return null;
        }

        return new LeafletMapInstance(this.mapProvider, mapTileSource, options);
    }
}

export class LeafletMapInstance implements MapInstance {
    public dependencyLoaded: boolean = false;
    public inited: boolean = false;

    private readonly defaultZoomLevel: number;
    private readonly minZoomLevel: number;
    private readonly maxZoomLevel: number;

    private readonly mapProvider: string;
    private readonly presetMapTileSource: LeafletTileSource;
    private readonly mapCreateOptions: MapCreateOptions;

    private leafletInstance: unknown | null;
    private leafletTileLayer: unknown | null;
    private leafletAnnotationLayer: unknown | null;
    private leafletZoomControl: unknown | null;
    private leafletAttribution: unknown | null;
    private leafletCenterMarker: unknown | null;

    public constructor(mapProvider: string, mapTileSource: LeafletTileSource, options: MapCreateOptions) {
        this.dependencyLoaded = !!LeafletMapProvider.Leaflet;

        this.mapProvider = mapProvider;
        this.presetMapTileSource = mapTileSource;
        this.mapCreateOptions = options;

        this.defaultZoomLevel = this.presetMapTileSource?.defaultZoomLevel || getCustomMapDefaultZoomLevel();
        this.minZoomLevel = this.presetMapTileSource?.minZoom || getCustomMapMinZoomLevel();
        this.maxZoomLevel = this.presetMapTileSource?.maxZoom || getCustomMapMaxZoomLevel();
    }

    public initMapInstance(mapContainer: HTMLElement, options: MapInstanceInitOptions): void {
        if (!LeafletMapProvider.Leaflet) {
            return;
        }

        const leaflet = LeafletMapProvider.Leaflet;
        const leafletInstance = leaflet.map(mapContainer, {
            center: [ options.initCenter.latitude, options.initCenter.longitude ],
            zoom: options.zoomLevel,
            attributionControl: false,
            zoomControl: false
        });

        let tileUrlFormat, tileUrlSubDomains, annotationUrlFormat, annotationUrlSubDomains: string | undefined;
        let minZoom, maxZoom: number;

        if (this.mapProvider !== 'custom') {
            tileUrlFormat = this.presetMapTileSource?.tileUrlFormat;
            tileUrlSubDomains = this.presetMapTileSource?.tileUrlSubDomains;
            annotationUrlFormat = this.presetMapTileSource?.annotationUrlFormat;
            annotationUrlSubDomains = this.presetMapTileSource?.annotationUrlSubDomains;
            minZoom = this.presetMapTileSource?.minZoom;
            maxZoom = this.presetMapTileSource?.maxZoom;
        } else {
            tileUrlFormat = getCustomMapTileLayerUrl();
            annotationUrlFormat = getCustomMapAnnotationLayerUrl();
            minZoom = getCustomMapMinZoomLevel();
            maxZoom = getCustomMapMaxZoomLevel();
        }

        if (isMapDataFetchProxyEnabled()) {
            tileUrlFormat = services.generateMapProxyTileImageUrl(this.mapProvider, options.language);
            tileUrlSubDomains = '';
        } else if (this.presetMapTileSource && this.presetMapTileSource.tileUrlExtraParams) {
            tileUrlFormat = this.getFinalUrlFormat(this.presetMapTileSource.tileUrlFormat as string, this.presetMapTileSource.tileUrlExtraParams, options);
        }

        const tileLayer = leaflet.tileLayer(tileUrlFormat, {
            subdomains: tileUrlSubDomains,
            maxZoom: maxZoom,
            minZoom: minZoom
        });
        tileLayer.addTo(leafletInstance);

        if (annotationUrlFormat || (this.mapProvider === 'custom' && isCustomMapAnnotationLayerDataFetchProxyEnabled())) {
            if (isMapDataFetchProxyEnabled()) {
                annotationUrlFormat = services.generateMapProxyAnnotationImageUrl(this.mapProvider, options.language);
                annotationUrlSubDomains = '';
            } else if (this.presetMapTileSource && this.presetMapTileSource.annotationUrlExtraParams) {
                annotationUrlFormat = this.getFinalUrlFormat(this.presetMapTileSource.annotationUrlFormat as string, this.presetMapTileSource.annotationUrlExtraParams, options);
            }

            const annotationLayer = leaflet.tileLayer(annotationUrlFormat, {
                subdomains: annotationUrlSubDomains,
                maxZoom: maxZoom,
                minZoom: minZoom
            });
            annotationLayer.addTo(leafletInstance);

            this.leafletAnnotationLayer = annotationLayer;
        }

        if (this.mapCreateOptions.enableZoomControl) {
            this.leafletZoomControl = leaflet.control.zoom({
                zoomInTitle: options.text.zoomIn,
                zoomOutTitle: options.text.zoomOut
            });
            this.leafletZoomControl.addTo(leafletInstance);
        }

        if (this.presetMapTileSource && this.presetMapTileSource.attribution) {
            const attribution = leaflet.control.attribution({
                prefix: false
            });
            attribution.addAttribution(this.presetMapTileSource.attribution);
            attribution.addTo(leafletInstance);
            this.leafletAttribution = attribution;
        }

        leafletInstance.addEventListener('click', function(e) {
            if (options.onClick) {
                options.onClick({
                    latitude: e.latlng.lat,
                    longitude: e.latlng.lng
                });
            }
        });

        leafletInstance.addEventListener('zoomanim', function(e) {
            if (options.onZoomChange && 'zoom' in e && isNumber(e.zoom)) {
                options.onZoomChange(e.zoom);
            }
        });

        this.leafletInstance = leafletInstance;
        this.leafletTileLayer = tileLayer;
        this.inited = true;
    }

    public getDefaultZoomLevel(): number {
        return this.defaultZoomLevel;
    }

    public getMinZoomLevel(): number {
        return this.minZoomLevel;
    }

    public getMaxZoomLevel(): number {
        return this.maxZoomLevel;
    }

    public getZoomLevel(): number {
        if (!this.leafletInstance || !isFunction(this.leafletInstance.getZoom)) {
            return this.defaultZoomLevel;
        }

        return this.leafletInstance.getZoom() ?? this.defaultZoomLevel;
    }

    public setMapCenterTo(center: Coordinate, zoomLevel: number): void {
        if (!this.leafletInstance) {
            return;
        }

        this.leafletInstance.setView([ center.latitude, center.longitude ], zoomLevel);
    }

    public setMapCenterMarker(position: Coordinate): void {
        if (!LeafletMapProvider.Leaflet || !this.leafletInstance) {
            return;
        }

        const leaflet = LeafletMapProvider.Leaflet;

        if (!this.leafletCenterMarker) {
            const markerIcon = leaflet.icon({
                iconUrl: 'img/map-marker-icon.png',
                iconRetinaUrl: 'img/map-marker-icon-2x.png',
                iconSize: [25, 32],
                iconAnchor: [12, 32],
                shadowUrl: 'img/map-marker-shadow.png',
                shadowSize: [41, 32]
            });
            this.leafletCenterMarker = leaflet.marker([ position.latitude, position.longitude ], {
                icon: markerIcon
            });
            this.leafletCenterMarker.addTo(this.leafletInstance);
        } else {
            this.leafletCenterMarker.setLatLng([ position.latitude, position.longitude ]);
        }
    }

    public removeMapCenterMarker(): void {
        if (!this.leafletInstance || !this.leafletCenterMarker) {
            return;
        }

        this.leafletCenterMarker.remove();
        this.leafletCenterMarker = null;
    }

    public zoomIn(): void {
        if (!this.leafletInstance) {
            return;
        }

        this.leafletInstance.zoomIn();
    }

    public zoomOut(): void {
        if (!this.leafletInstance) {
            return;
        }

        this.leafletInstance.zoomOut();
    }

    private getFinalUrlFormat(urlFormat: string, urlExtraParams: LeafletTileSourceExtraParam[], options: MapInstanceInitOptions) {
        const params: string[] = [];

        for (const param of urlExtraParams) {
            if (param.paramValueType === 'tomtom_key') {
                params.push(param.paramName + '=' + getTomTomMapAPIKey());
            } else if (param.paramValueType === 'tianditu_key') {
                params.push(param.paramName + '=' + getTianDiTuMapAPIKey());
            } else if (param.paramValueType === 'language' && options.language) {
                params.push(param.paramName + '=' + options.language);
            }
        }

        if (params.length) {
            if (urlFormat.indexOf('?') >= 0) {
                urlFormat = urlFormat + '&' + params.join('&');
            } else {
                urlFormat = urlFormat + '?' + params.join('&');
            }
        }

        return urlFormat;
    }
}
