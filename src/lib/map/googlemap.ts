// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-nocheck
import type { MapPosition } from '@/core/map.ts';
import type { MapProvider, MapInstance, MapInstanceInitOptions } from './base.ts';

import { asyncLoadAssets } from '@/lib/misc.ts';
import services from '@/lib/services.ts';

export class GoogleMapProvider implements MapProvider {
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

    public createMapInstance(): MapInstance | null {
        return new GoogleMapInstance();
    }
}

export class GoogleMapInstance implements MapInstance {
    public dependencyLoaded: boolean = false;
    public inited: boolean = false;

    public readonly defaultZoomLevel: number = 14;
    public readonly minZoomLevel: number = 1;

    private googleMapInstance: unknown = null;
    private googleMapCenterMarker: unknown | null;

    public constructor() {
        this.dependencyLoaded = !!GoogleMapProvider.GoogleMap;
    }

    public initMapInstance(mapContainer: HTMLElement, options: MapInstanceInitOptions): void {
        if (!GoogleMapProvider.GoogleMap) {
            return;
        }

        const googleMap = GoogleMapProvider.GoogleMap;

        this.googleMapInstance = new googleMap.Map(mapContainer, {
            zoom: options.zoomLevel,
            center: {
                lat: options.initCenter.latitude,
                lng: options.initCenter.longitude
            },
            maxZoom: 19,
            zoomControl: true,
            mapTypeControl: false,
            scaleControl: false,
            streetViewControl: false,
            rotateControl: false,
            fullscreenControl: false,
            gestureHandling: 'greedy',
            zoomControlOptions: {
                position: GoogleMapProvider.ControlPosition.LEFT_TOP
            }
        });

        this.googleMapInstance.addListener('click', function(e) {
            if (options.onClick) {
                options.onClick({
                    latitude: e.latLng.lat(),
                    longitude: e.latLng.lng()
                });
            }
        });

        this.inited = true;
    }

    public setMapCenterTo(center: MapPosition, zoomLevel: number): void {
        if (!GoogleMapProvider.GoogleMap || !this.googleMapInstance) {
            return;
        }

        this.googleMapInstance.setCenter({
            lat: center.latitude,
            lng: center.longitude
        });
        this.googleMapInstance.setZoom(zoomLevel);
    }

    public setMapCenterMarker(position: MapPosition): void {
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

    public removeMapCenterMarker(): void {
        if (!this.googleMapInstance || !this.googleMapCenterMarker) {
            return;
        }

        this.googleMapCenterMarker.setMap(null);
        this.googleMapCenterMarker = null;
    }
}
