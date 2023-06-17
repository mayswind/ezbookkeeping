import { asyncLoadAssets } from "@/lib/misc.js";
import services from "@/lib/services.js";
import logger from '@/lib/logger.js';

const baiduMapHolder = {
    BMap: null,
    BMAP_NAVIGATION_CONTROL_ZOOM: window.BMAP_NAVIGATION_CONTROL_ZOOM || 3,
    BMAP_ANCHOR_TOP_LEFT: window.BMAP_ANCHOR_TOP_LEFT || 0,
    COORDINATES_WGS84: window.COORDINATES_WGS84 || 1,
    COORDINATES_BD09: window.COORDINATES_BD09 || 5
};

export function loadBaiduMapAssets() {
    if (baiduMapHolder.BMap) {
        return;
    }

    if (!window.onBMapCallback) {
        window.onBMapCallback = () => {
            baiduMapHolder.BMap = window.BMap;
        };
    }

    return asyncLoadAssets('js', services.generateBaiduMapJavascriptUrl('onBMapCallback'));
}

export function createBaiduMapHolder() {
    return {
        mapProvider: 'baidumap',
        dependencyLoaded: !!baiduMapHolder.BMap,
        inited: false,
        defaultZoomLevel: 15,
        minZoomLevel: 1,
        baiduMapInstance: null,
        baiduMapConverter: null,
        baiduMapNavigationControl: null,
        baiduMapCenterMarker: null
    };
}

export function createBaiduMapInstance(mapHolder, mapContainer) {
    if (!baiduMapHolder.BMap) {
        return null;
    }

    const BMap = baiduMapHolder.BMap;
    const baiduMapInstance = new BMap.Map(mapContainer, {
        maxZoom: 19
    });
    baiduMapInstance.enableScrollWheelZoom();

    const baiduMapNavigationControl = new BMap.NavigationControl({
        type: baiduMapHolder.BMAP_NAVIGATION_CONTROL_ZOOM,
        anchor: baiduMapHolder.BMAP_ANCHOR_TOP_LEFT
    });
    baiduMapInstance.addControl(baiduMapNavigationControl);

    mapHolder.baiduMapInstance = baiduMapInstance;
    mapHolder.baiduMapConverter = new BMap.Convertor();
    mapHolder.inited = true;
}

export function setBaiduMapCenterTo(mapHolder, center, zoomLevel) {
    if (!baiduMapHolder.BMap || !mapHolder.baiduMapInstance) {
        return;
    }

    const BMap = baiduMapHolder.BMap;
    const centerPoint = new BMap.Point(center.longitude, center.latitude);

    if (mapHolder.baiduMapConverter) {
        mapHolder.baiduMapConverter.translate([ centerPoint ], baiduMapHolder.COORDINATES_WGS84, baiduMapHolder.COORDINATES_BD09, data => {
            if (data.status !== 0) {
                logger.warn('baidu map geo position convert failed');
            }

            const actualPoint = (data.status === 0 ? data.points[0] : centerPoint);
            mapHolder.baiduMapInstance.centerAndZoom(actualPoint, zoomLevel);
        });
    } else {
        mapHolder.baiduMapInstance.centerAndZoom(centerPoint, zoomLevel);
    }
}

export function setBaiduMapCenterMaker(mapHolder, position) {
    if (!baiduMapHolder.BMap || !mapHolder.baiduMapInstance) {
        return;
    }

    const BMap = baiduMapHolder.BMap;
    const markerPoint = new BMap.Point(position.longitude, position.latitude);

    mapHolder.baiduMapConverter.translate([ markerPoint ], baiduMapHolder.COORDINATES_WGS84, baiduMapHolder.COORDINATES_BD09, data => {
        if (data.status !== 0) {
            logger.warn('baidu map geo position convert failed');
        }

        const actualPoint = (data.status === 0 ? data.points[0] : markerPoint);

        if (!mapHolder.baiduMapCenterMarker) {
            mapHolder.baiduMapCenterMarker = new BMap.Marker(actualPoint);
            mapHolder.baiduMapInstance.addOverlay(mapHolder.baiduMapCenterMarker);
        } else {
            mapHolder.baiduMapCenterMarker.setPosition(actualPoint);
        }
    });
}

export function removeBaiduMapCenterMaker(mapHolder) {
    if (!mapHolder.baiduMapInstance || !mapHolder.baiduMapCenterMarker) {
        return;
    }

    mapHolder.baiduMapInstance.removeOverlay(mapHolder.baiduMapCenterMarker);
    mapHolder.baiduMapCenterMarker = null;
}
