import { asyncLoadAssets } from '@/lib/misc.js';
import services from '@/lib/services.js';
import logger from '@/lib/logger.js';

const baiduMapHolder = {
    BMap: null,
    BMAP_NAVIGATION_CONTROL_ZOOM: window.BMAP_NAVIGATION_CONTROL_ZOOM || 3,
    BMAP_ANCHOR_TOP_LEFT: window.BMAP_ANCHOR_TOP_LEFT || 0,
    COORDINATES_WGS84: window.COORDINATES_WGS84 || 1,
    COORDINATES_BD09: window.COORDINATES_BD09 || 5
};

export function getBaiduMapWebsite() {
    return 'https://map.baidu.com';
}

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
        baiduMapCenterPosition: null,
        baiduMapCenterMarker: null
    };
}

export function createBaiduMapInstance(mapHolder, mapContainer, options) {
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
    baiduMapInstance.centerAndZoom(new BMap.Point(options.initCenter.longitude, options.initCenter.latitude), options.zoomLevel);

    mapHolder.baiduMapInstance = baiduMapInstance;
    mapHolder.baiduMapConverter = new BMap.Convertor();
    mapHolder.baiduMapNavigationControl = baiduMapNavigationControl;
    mapHolder.inited = true;
}

export function setBaiduMapCenterTo(mapHolder, center, zoomLevel) {
    if (!baiduMapHolder.BMap || !mapHolder.baiduMapInstance) {
        return;
    }

    const BMap = baiduMapHolder.BMap;

    if (baiduMapHolder.baiduMapCenterPosition
        && baiduMapHolder.baiduMapCenterPosition.originalLongitude === center.longitude
        && baiduMapHolder.baiduMapCenterPosition.originalLatitude === center.latitude
        && baiduMapHolder.baiduMapCenterPosition.convertedLongitude
        && baiduMapHolder.baiduMapCenterPosition.convertedLatitude
    ) {
        mapHolder.baiduMapInstance.centerAndZoom(new BMap.Point(baiduMapHolder.baiduMapCenterPosition.convertedLongitude, baiduMapHolder.baiduMapCenterPosition.convertedLatitude), zoomLevel);
        return;
    }

    baiduMapHolder.baiduMapCenterPosition = {
        originalLongitude: center.longitude,
        originalLatitude: center.latitude,
        convertedLongitude: null,
        convertedLatitude: null
    };

    const centerPoint = new BMap.Point(center.longitude, center.latitude);

    if (mapHolder.baiduMapConverter) {
        mapHolder.baiduMapConverter.translate([ centerPoint ], baiduMapHolder.COORDINATES_WGS84, baiduMapHolder.COORDINATES_BD09, data => {
            let convertedCenterPoint = centerPoint;

            if (data.status !== 0 || !data.points) {
                logger.warn('baidu map geo position convert failed');
            } else {
                convertedCenterPoint = data.points[0];
                baiduMapHolder.baiduMapCenterPosition.convertedLongitude = convertedCenterPoint.lng;
                baiduMapHolder.baiduMapCenterPosition.convertedLatitude = convertedCenterPoint.lat;
            }

            mapHolder.baiduMapInstance.centerAndZoom(convertedCenterPoint, zoomLevel);
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
    const setMaker = function (point) {
        if (!mapHolder.baiduMapCenterMarker) {
            mapHolder.baiduMapCenterMarker = new BMap.Marker(point);
            mapHolder.baiduMapInstance.addOverlay(mapHolder.baiduMapCenterMarker);
        } else {
            mapHolder.baiduMapCenterMarker.setPosition(point);
        }
    }

    if (baiduMapHolder.baiduMapCenterPosition
        && baiduMapHolder.baiduMapCenterPosition.originalLongitude === position.longitude
        && baiduMapHolder.baiduMapCenterPosition.originalLatitude === position.latitude
        && baiduMapHolder.baiduMapCenterPosition.convertedLongitude
        && baiduMapHolder.baiduMapCenterPosition.convertedLatitude
    ) {
        setMaker(new BMap.Point(baiduMapHolder.baiduMapCenterPosition.convertedLongitude, baiduMapHolder.baiduMapCenterPosition.convertedLatitude));
        return;
    }

    const markerPoint = new BMap.Point(position.longitude, position.latitude);

    if (mapHolder.baiduMapConverter) {
        mapHolder.baiduMapConverter.translate([ markerPoint ], baiduMapHolder.COORDINATES_WGS84, baiduMapHolder.COORDINATES_BD09, data => {
            let convertedMarkPoint = markerPoint;

            if (data.status !== 0 || !data.points) {
                logger.warn('baidu map geo position convert failed');
            } else {
                convertedMarkPoint = data.points[0];
            }

            setMaker(convertedMarkPoint);
        });
    } else {
        setMaker(markerPoint);
    }
}

export function removeBaiduMapCenterMaker(mapHolder) {
    if (!mapHolder.baiduMapInstance || !mapHolder.baiduMapCenterMarker) {
        return;
    }

    mapHolder.baiduMapInstance.removeOverlay(mapHolder.baiduMapCenterMarker);
    mapHolder.baiduMapCenterMarker = null;
}
