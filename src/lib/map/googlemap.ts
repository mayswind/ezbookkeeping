// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-nocheck
import type { Coordinate } from '@/core/coordinate.ts';
import type { MapProvider, MapInstance, MapCreateOptions, MapInstanceInitOptions } from './base.ts';

import { isFunction } from '@/lib/common.ts';
import { asyncLoadAssets } from '@/lib/misc.ts';
import services from '@/lib/services.ts';

export class GoogleMapProvider implements MapProvider {
    // https://developers.google.com/maps/documentation/javascript/reference/map
    public static GoogleMap: unknown = null;
    public static ControlPosition = {
        LEFT_TOP: (window.google && window.google.maps && window.google.maps.ControlPosition) ? window.google.maps.ControlPosition.LEFT_TOP : 5
    };

    public getWebsite(): string {
        return 'https://maps.google.com';
    }

    public isSupportGetGeoLocationByClick(): boolean {
        return true;
    }

    public asyncLoadAssets(language?: string): Promise<unknown> {
        if (GoogleMapProvider.GoogleMap) {
            return Promise.resolve();
        }

        if (!window.onGoogleMapCallback) {
            window.onGoogleMapCallback = () => {
                if (window.google) {
                    GoogleMapProvider.GoogleMap = window.google.maps;
                    GoogleMapProvider.ControlPosition.LEFT_TOP = (window.google && window.google.maps && window.google.maps.ControlPosition) ? window.google.maps.ControlPosition.LEFT_TOP : 5;
                }
            };
        }

        return asyncLoadAssets('js', services.generateGoogleMapJavascriptUrl(language, 'onGoogleMapCallback'));
    }

    public createMapInstance(options: MapCreateOptions): MapInstance | null {
        return new GoogleMapInstance(options);
    }
}

export class GoogleMapInstance implements MapInstance {
    public dependencyLoaded: boolean = false;
    public inited: boolean = false;

    private readonly defaultZoomLevel: number = 14;
    private readonly minZoomLevel: number = 0;
    private readonly maxZoomLevel: number = 19;

    private readonly mapCreateOptions: MapCreateOptions;
    private googleMapInstance: unknown = null;
    private googleMapCenterMarker: unknown | null;

    public constructor(options: MapCreateOptions) {
        this.dependencyLoaded = !!GoogleMapProvider.GoogleMap;
        this.mapCreateOptions = options;
    }

    public initMapInstance(mapContainer: HTMLElement, options: MapInstanceInitOptions): void {
        if (!GoogleMapProvider.GoogleMap) {
            return;
        }

        const googleMap = GoogleMapProvider.GoogleMap;
        const googleMapInstance = new googleMap.Map(mapContainer, {
            zoom: options.zoomLevel,
            center: {
                lat: options.initCenter.latitude,
                lng: options.initCenter.longitude
            },
            maxZoom: 19,
            zoomControl: !!this.mapCreateOptions.enableZoomControl,
            mapTypeControl: false,
            scaleControl: false,
            streetViewControl: false,
            rotateControl: false,
            fullscreenControl: false,
            gestureHandling: 'greedy',
            zoomControlOptions: this.mapCreateOptions.enableZoomControl ? {
                position: GoogleMapProvider.ControlPosition.LEFT_TOP
            } : undefined
        });

        googleMapInstance.addListener('click', function(e) {
            if (options.onClick) {
                options.onClick({
                    latitude: e.latLng.lat(),
                    longitude: e.latLng.lng()
                });
            }
        });

        googleMapInstance.addListener('zoom_changed', function() {
            if (options.onZoomChange && isFunction(googleMapInstance.getZoom)) {
                options.onZoomChange(googleMapInstance.getZoom());
            }
        });

        this.googleMapInstance = googleMapInstance;
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
        if (!this.googleMapInstance || !isFunction(this.googleMapInstance.getZoom)) {
            return this.defaultZoomLevel;
        }

        return this.googleMapInstance.getZoom() ?? this.defaultZoomLevel;
    }

    public setMapCenterTo(center: Coordinate, zoomLevel: number): void {
        if (!GoogleMapProvider.GoogleMap || !this.googleMapInstance) {
            return;
        }

        this.googleMapInstance.setCenter({
            lat: center.latitude,
            lng: center.longitude
        });
        this.googleMapInstance.setZoom(zoomLevel);
    }

    public setMapCenterMarker(position: Coordinate): void {
        if (!GoogleMapProvider.GoogleMap || !this.googleMapInstance) {
            return;
        }

        const googleMap = GoogleMapProvider.GoogleMap;

        if (!this.googleMapCenterMarker) {
            this.googleMapCenterMarker = new googleMap.Marker({
                position: {
                    lat: position.latitude,
                    lng: position.longitude
                },
                map: this.googleMapInstance
            });
        } else {
            this.googleMapCenterMarker.setPosition({
                lat: position.latitude,
                lng: position.longitude
            });
        }
    }

    public zoomIn(): void {
        if (!this.googleMapInstance) {
            return;
        }

        this.googleMapInstance.setZoom(this.googleMapInstance.getZoom() + 1);
    }

    public zoomOut(): void {
        if (!this.googleMapInstance) {
            return;
        }

        this.googleMapInstance.setZoom(this.googleMapInstance.getZoom() - 1);
    }

    public removeMapCenterMarker(): void {
        if (!this.googleMapInstance || !this.googleMapCenterMarker) {
            return;
        }

        this.googleMapCenterMarker.setMap(null);
        this.googleMapCenterMarker = null;
    }
}
