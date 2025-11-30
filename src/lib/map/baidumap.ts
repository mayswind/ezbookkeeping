// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-nocheck
import type { Coordinate } from '@/core/coordinate.ts';
import type { MapProvider, MapInstance, MapCreateOptions, MapInstanceInitOptions } from './base.ts';

import { isFunction } from '@/lib/common.ts';
import { asyncLoadAssets } from '@/lib/misc.ts';
import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

export class BaiduMapProvider implements MapProvider {
    // https://mapopen-pub-jsapi.bj.bcebos.com/jsapi/reference/jsapi_reference_3_0.html
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

    public createMapInstance(options: MapCreateOptions): MapInstance | null {
        return new BaiduMapInstance(options);
    }
}

export class BaiduMapInstance implements MapInstance {
    public dependencyLoaded: boolean = false;
    public inited: boolean = false;

    private readonly defaultZoomLevel: number = 15;
    private readonly minZoomLevel: number = 4;
    private readonly maxZoomLevel: number = 19;

    private readonly mapCreateOptions: MapCreateOptions;
    private baiduMapInstance: unknown = null;
    private baiduMapConverter: unknown = null;
    private baiduMapNavigationControl: unknown = null;
    private baiduMapCenterPosition: unknown = null;
    private baiduMapCenterMarker: unknown | null;

    public constructor(options: MapCreateOptions) {
        this.dependencyLoaded = !!BaiduMapProvider.BMap;
        this.mapCreateOptions = options;
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

        if (this.mapCreateOptions.enableZoomControl) {
            this.baiduMapNavigationControl = new BMap.NavigationControl({
                type: BaiduMapProvider.BMAP_NAVIGATION_CONTROL_ZOOM,
                anchor: BaiduMapProvider.BMAP_ANCHOR_TOP_LEFT
            });
            baiduMapInstance.addControl(this.baiduMapNavigationControl);
        }

        baiduMapInstance.centerAndZoom(new BMap.Point(options.initCenter.longitude, options.initCenter.latitude), options.zoomLevel);

        baiduMapInstance.addEventListener('click', function(e) {
            if (options.onClick) {
                options.onClick({
                    latitude: e.point.lat,
                    longitude: e.point.lng
                });
            }
        });

        baiduMapInstance.addEventListener('zoomend', function() {
            if (options.onZoomChange && isFunction(baiduMapInstance.getZoom)) {
                options.onZoomChange(baiduMapInstance.getZoom());
            }
        });

        this.baiduMapInstance = baiduMapInstance;
        this.baiduMapConverter = new BMap.Convertor();
        this.inited = true;
    }

    public getDefaultZoomLevel(): number {
        return this.defaultZoomLevel;
    }

    public getMinZoomLevel(): number {
        if (!this.baiduMapInstance || !isFunction(this.baiduMapInstance.getMapType)) {
            return this.minZoomLevel;
        }

        const mapType = this.baiduMapInstance.getMapType();

        if (!mapType || !isFunction(mapType.getMinZoom)) {
            return this.minZoomLevel;
        }

        return mapType.getMinZoom() ?? this.minZoomLevel;
    }

    public getMaxZoomLevel(): number {
        if (!this.baiduMapInstance || !isFunction(this.baiduMapInstance.getMapType)) {
            return this.maxZoomLevel;
        }

        const mapType = this.baiduMapInstance.getMapType();

        if (!mapType || !isFunction(mapType.getMaxZoom)) {
            return this.maxZoomLevel;
        }

        return mapType.getMaxZoom() ?? this.maxZoomLevel;
    }

    public getZoomLevel(): number {
        if (!this.baiduMapInstance || !isFunction(this.baiduMapInstance.getZoom)) {
            return this.defaultZoomLevel;
        }

        return this.baiduMapInstance.getZoom() ?? this.defaultZoomLevel;
    }

    public setMapCenterTo(center: Coordinate, zoomLevel: number): void {
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

    public setMapCenterMarker(position: Coordinate): void {
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

    public zoomIn(): void {
        if (!this.baiduMapInstance) {
            return;
        }

        this.baiduMapInstance.zoomIn();
    }

    public zoomOut(): void {
        if (!this.baiduMapInstance) {
            return;
        }

        this.baiduMapInstance.zoomOut();
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
