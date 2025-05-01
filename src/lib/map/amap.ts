// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-nocheck
import type { MapPosition } from '@/core/map.ts';
import type { MapProvider, MapInstance, MapInstanceInitOptions } from './base.ts';

import { asyncLoadAssets } from '@/lib/misc.ts';
import services from '@/lib/services.ts';
import {
    getAmapSecurityVerificationMethod,
    getAmapApiExternalProxyUrl,
    getAmapApplicationSecret
} from '@/lib/server_settings.ts';
import logger from '@/lib/logger.ts';

export class AmapMapProvider implements MapProvider {
    public static AMap: unknown = null;

    public getWebsite(): string {
        return 'https://www.amap.com';
    }

    public isSupportGetGeoLocationByClick(): boolean {
        return false;
    }

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    public asyncLoadAssets(language?: string): Promise<unknown> {
        if (AmapMapProvider.AMap) {
            return Promise.resolve();
        }

        if (!window._AMapSecurityConfig) {
            const amapSecurityConfig = {};

            if (getAmapSecurityVerificationMethod() === 'internalproxy') {
                amapSecurityConfig.serviceHost = services.generateAmapApiInternalProxyUrl();
            } else if (getAmapSecurityVerificationMethod() === 'externalproxy') {
                amapSecurityConfig.serviceHost = getAmapApiExternalProxyUrl();
            } else if (getAmapSecurityVerificationMethod() === 'plaintext') {
                amapSecurityConfig.securityJsCode = getAmapApplicationSecret();
            }

            window._AMapSecurityConfig = amapSecurityConfig;
        }

        if (!window.onAMapCallback) {
            window.onAMapCallback = () => {
                AmapMapProvider.AMap = window.AMap;
            };
        }

        return asyncLoadAssets('js', services.generateAmapJavascriptUrl('onAMapCallback'));
    }

    public createMapInstance(): MapInstance | null {
        return new AmapMapInstance();
    }
}

export class AmapMapInstance implements MapInstance {
    public dependencyLoaded: boolean = false;
    public inited: boolean = false;

    public readonly defaultZoomLevel: number = 14;
    public readonly minZoomLevel: number = 1;

    private amapInstance: unknown = null;
    private amapToolbar: unknown = null;
    private amapCenterPosition: unknown = null;
    private amapCenterMarker: unknown | null;

    public constructor() {
        this.dependencyLoaded = !!AmapMapProvider.AMap;
    }

    public initMapInstance(mapContainer: HTMLElement, options: MapInstanceInitOptions): void {
        if (!AmapMapProvider.AMap) {
            return;
        }

        const AMap = AmapMapProvider.AMap;
        const amapInstance = new AMap.Map(mapContainer, {
            zoom: options.zoomLevel,
            center: [ options.initCenter.longitude, options.initCenter.latitude ],
            zooms: [ 1, 19 ],
            jogEnable: false
        });

        const amapToolbar = new AMap.ToolBar({
            position: 'LT'
        });
        amapInstance.addControl(amapToolbar);

        amapInstance.on('click', function(e) {
            if (options.onClick) {
                options.onClick({
                    latitude: e.lnglat.lat,
                    longitude: e.lnglat.lng
                });
            }
        });

        this.amapInstance = amapInstance;
        this.amapToolbar = amapToolbar;
        this.inited = true;
    }

    public setMapCenterTo(center: MapPosition, zoomLevel: number): void {
        if (!AmapMapProvider.AMap || !this.amapInstance) {
            return;
        }

        const AMap = AmapMapProvider.AMap;

        if (this.amapCenterPosition
            && this.amapCenterPosition.originalLongitude === center.longitude
            && this.amapCenterPosition.originalLatitude === center.latitude
            && this.amapCenterPosition.convertedLongitude
            && this.amapCenterPosition.convertedLatitude
        ) {
            this.amapInstance.setZoomAndCenter(zoomLevel, new AMap.LngLat(this.amapCenterPosition.convertedLongitude, this.amapCenterPosition.convertedLatitude));
            return;
        }

        this.amapCenterPosition = {
            originalLongitude: center.longitude,
            originalLatitude: center.latitude,
            convertedLongitude: null,
            convertedLatitude: null
        };

        const centerPoint = new AMap.LngLat(center.longitude, center.latitude);

        AMap.convertFrom(centerPoint, 'gps', (status, result) => {
            let convertedCenterPoint = centerPoint;

            if (result.info !== 'ok' || !result.locations) {
                logger.warn('amap geo position convert failed');
            } else {
                convertedCenterPoint = result.locations[0];
                this.amapCenterPosition.convertedLongitude = convertedCenterPoint.getLng();
                this.amapCenterPosition.convertedLatitude = convertedCenterPoint.getLat();
            }

            this.amapInstance.setZoomAndCenter(zoomLevel, convertedCenterPoint);
        });
    }

    public setMapCenterMarker(position: MapPosition): void {
        if (!AmapMapProvider.AMap || !this.amapInstance) {
            return;
        }

        const AMap = AmapMapProvider.AMap;

        if (this.amapCenterPosition
            && this.amapCenterPosition.originalLongitude === position.longitude
            && this.amapCenterPosition.originalLatitude === position.latitude
            && this.amapCenterPosition.convertedLongitude
            && this.amapCenterPosition.convertedLatitude
        ) {
            this.setMaker(new AMap.LngLat(this.amapCenterPosition.convertedLongitude, this.amapCenterPosition.convertedLatitude));
            return;
        }

        const markerPoint = new AMap.LngLat(position.longitude, position.latitude);

        AMap.convertFrom(markerPoint, 'gps', (status, result) => {
            let convertedMarkPoint = markerPoint;

            if (result.info !== 'ok' || !result.locations) {
                logger.warn('amap geo position convert failed');
            } else {
                convertedMarkPoint = result.locations[0];
            }

            this.setMaker(convertedMarkPoint);
        });
    }

    public removeMapCenterMarker(): void {
        if (!this.amapInstance || !this.amapCenterMarker) {
            return;
        }

        this.amapInstance.remove(this.amapCenterMarker);
        this.amapCenterMarker = null;
    }

    private setMaker(point: unknown): void {
        if (!AmapMapProvider.AMap || !this.amapInstance) {
            return;
        }

        const AMap = AmapMapProvider.AMap;

        if (!this.amapCenterMarker) {
            this.amapCenterMarker = new AMap.Marker({
                position: point
            });
            this.amapInstance.add(this.amapCenterMarker);
        } else {
            this.amapCenterMarker.setPosition(point);
        }
    }
}
