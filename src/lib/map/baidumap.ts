// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-nocheck
import type { MapPosition } from '@/core/map.ts';
import type { MapProvider, MapInstance, MapInstanceInitOptions } from './base.ts';

import { asyncLoadAssets } from '@/lib/misc.ts';
import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

export class BaiduMapProvider implements MapProvider {
    public static BMap: unknown = null;
    public static BMAP_NAVIGATION_CONTROL_ZOOM: unknown = window.BMAP_NAVIGATION_CONTROL_ZOOM || 3;
    public static BMAP_ANCHOR_TOP_LEFT: unknown = window.BMAP_ANCHOR_TOP_LEFT || 0;
    public static COORDINATES_WGS84: unknown = window.COORDINATES_WGS84 || 1;
    public static COORDINATES_BD09: unknown = window.COORDINATES_BD09 || 5;

    public getWebsite(): string {
        return 'https://map.baidu.com';
    }

    public isSupportGetGeoLocationByClick(): boolean {
        return false;
    }

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    public asyncLoadAssets(language?: string): Promise<unknown> {
        if (BaiduMapProvider.BMap) {
            return Promise.resolve();
        }

        if (!window.onBMapCallback) {
            window.onBMapCallback = () => {
                BaiduMapProvider.BMap = window.BMap;
                BaiduMapProvider.BMAP_NAVIGATION_CONTROL_ZOOM = window.BMAP_NAVIGATION_CONTROL_ZOOM || 3;
                BaiduMapProvider.BMAP_ANCHOR_TOP_LEFT = window.BMAP_ANCHOR_TOP_LEFT || 0;
                BaiduMapProvider.COORDINATES_WGS84 = window.COORDINATES_WGS84 || 1;
                BaiduMapProvider.COORDINATES_BD09 = window.COORDINATES_BD09 || 5;
            };
        }

        return asyncLoadAssets('js', services.generateBaiduMapJavascriptUrl('onBMapCallback'));
    }

    public createMapInstance(): MapInstance | null {
        return new BaiduMapInstance();
    }
}

export class BaiduMapInstance implements MapInstance {
    public dependencyLoaded: boolean = false;
    public inited: boolean = false;

    public readonly defaultZoomLevel: number = 15;
    public readonly minZoomLevel: number = 1;

    private baiduMapInstance: unknown = null;
    private baiduMapConverter: unknown = null;
    private baiduMapNavigationControl: unknown = null;
    private baiduMapCenterPosition: unknown = null;
    private baiduMapCenterMarker: unknown | null;

    public constructor() {
        this.dependencyLoaded = !!BaiduMapProvider.BMap;
    }

    public initMapInstance(mapContainer: HTMLElement, options: MapInstanceInitOptions): void {
        if (!BaiduMapProvider.BMap) {
            return;
        }

        const BMap = BaiduMapProvider.BMap;
        const baiduMapInstance = new BMap.Map(mapContainer, {
            maxZoom: 19
        });
        baiduMapInstance.enableScrollWheelZoom();

        const baiduMapNavigationControl = new BMap.NavigationControl({
            type: BaiduMapProvider.BMAP_NAVIGATION_CONTROL_ZOOM,
            anchor: BaiduMapProvider.BMAP_ANCHOR_TOP_LEFT
        });
        baiduMapInstance.addControl(baiduMapNavigationControl);
        baiduMapInstance.centerAndZoom(new BMap.Point(options.initCenter.longitude, options.initCenter.latitude), options.zoomLevel);

        baiduMapInstance.addEventListener('click', function(e) {
            if (options.onClick) {
                options.onClick({
                    latitude: e.point.lat,
                    longitude: e.point.lng
                });
            }
        });

        this.baiduMapInstance = baiduMapInstance;
        this.baiduMapConverter = new BMap.Convertor();
        this.baiduMapNavigationControl = baiduMapNavigationControl;
        this.inited = true;
    }

    public setMapCenterTo(center: MapPosition, zoomLevel: number): void {
        if (!BaiduMapProvider.BMap || !this.baiduMapInstance) {
            return;
        }

        const BMap = BaiduMapProvider.BMap;

        if (this.baiduMapCenterPosition
            && this.baiduMapCenterPosition.originalLongitude === center.longitude
            && this.baiduMapCenterPosition.originalLatitude === center.latitude
            && this.baiduMapCenterPosition.convertedLongitude
            && this.baiduMapCenterPosition.convertedLatitude
        ) {
            this.baiduMapInstance.centerAndZoom(new BMap.Point(this.baiduMapCenterPosition.convertedLongitude, this.baiduMapCenterPosition.convertedLatitude), zoomLevel);
            return;
        }

        this.baiduMapCenterPosition = {
            originalLongitude: center.longitude,
            originalLatitude: center.latitude,
            convertedLongitude: null,
            convertedLatitude: null
        };

        const centerPoint = new BMap.Point(center.longitude, center.latitude);

        if (this.baiduMapConverter) {
            this.baiduMapConverter.translate([ centerPoint ], BaiduMapProvider.COORDINATES_WGS84, BaiduMapProvider.COORDINATES_BD09, data => {
                let convertedCenterPoint = centerPoint;

                if (data.status !== 0 || !data.points) {
                    logger.warn('baidu map geo position convert failed');
                } else {
                    convertedCenterPoint = data.points[0];
                    this.baiduMapCenterPosition.convertedLongitude = convertedCenterPoint.lng;
                    this.baiduMapCenterPosition.convertedLatitude = convertedCenterPoint.lat;
                }

                this.baiduMapInstance.centerAndZoom(convertedCenterPoint, zoomLevel);
            });
        } else {
            this.baiduMapInstance.centerAndZoom(centerPoint, zoomLevel);
        }
    }

    public setMapCenterMarker(position: MapPosition): void {
        if (!BaiduMapProvider.BMap || !this.baiduMapInstance) {
            return;
        }

        const BMap = BaiduMapProvider.BMap;

        if (this.baiduMapCenterPosition
            && this.baiduMapCenterPosition.originalLongitude === position.longitude
            && this.baiduMapCenterPosition.originalLatitude === position.latitude
            && this.baiduMapCenterPosition.convertedLongitude
            && this.baiduMapCenterPosition.convertedLatitude
        ) {
            this.setMaker(new BMap.Point(this.baiduMapCenterPosition.convertedLongitude, this.baiduMapCenterPosition.convertedLatitude));
            return;
        }

        const markerPoint = new BMap.Point(position.longitude, position.latitude);

        if (this.baiduMapConverter) {
            this.baiduMapConverter.translate([ markerPoint ], BaiduMapProvider.COORDINATES_WGS84, BaiduMapProvider.COORDINATES_BD09, data => {
                let convertedMarkPoint = markerPoint;

                if (data.status !== 0 || !data.points) {
                    logger.warn('baidu map geo position convert failed');
                } else {
                    convertedMarkPoint = data.points[0];
                }

                this.setMaker(convertedMarkPoint);
            });
        } else {
            this.setMaker(markerPoint);
        }
    }

    public removeMapCenterMarker(): void {
        if (!this.baiduMapInstance || !this.baiduMapCenterMarker) {
            return;
        }

        this.baiduMapInstance.removeOverlay(this.baiduMapCenterMarker);
        this.baiduMapCenterMarker = null;
    }

    private setMaker(point: unknown): void {
        if (!BaiduMapProvider.BMap || !this.baiduMapInstance) {
            return;
        }

        const BMap = BaiduMapProvider.BMap;

        if (!this.baiduMapCenterMarker) {
            this.baiduMapCenterMarker = new BMap.Marker(point);
            this.baiduMapInstance.addOverlay(this.baiduMapCenterMarker);
        } else {
            this.baiduMapCenterMarker.setPosition(point);
        }
    }
}
