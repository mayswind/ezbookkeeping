import { asyncLoadAssets } from '@/lib/misc.js';
import services from '@/lib/services.js';
import settings from '@/lib/settings.js';
import logger from '@/lib/logger.js';

const amapHolder = {
    AMap: null
};

export function loadAmapAssets() {
    if (amapHolder.AMap) {
        return;
    }

    if (!window._AMapSecurityConfig) {
        const amapSecurityConfig = {};

        if (settings.getAmapSecurityVerificationMethod() === 'internalproxy') {
            amapSecurityConfig.serviceHost = services.generateAmapApiInternalProxyUrl();
        } else if (settings.getAmapSecurityVerificationMethod() === 'externalproxy') {
            amapSecurityConfig.serviceHost = settings.getAmapApiExternalProxyUrl();
        } else if (settings.getAmapSecurityVerificationMethod() === 'plaintext') {
            amapSecurityConfig.securityJsCode = settings.getAmapApplicationSecret();
        }

        window._AMapSecurityConfig = amapSecurityConfig;
    }

    if (!window.onAMapCallback) {
        window.onAMapCallback = () => {
            amapHolder.AMap = window.AMap;
        };
    }

    return asyncLoadAssets('js', services.generateAmapJavascriptUrl('onAMapCallback'));
}

export function createAmapHolder() {
    return {
        mapProvider: 'amap',
        dependencyLoaded: !!amapHolder.AMap,
        inited: false,
        defaultZoomLevel: 14,
        minZoomLevel: 1,
        amapInstance: null,
        amapToolbar: null,
        amapCenterPosition: null,
        amapCenterMarker: null
    };
}

export function createAmapInstance(mapHolder, mapContainer, options) {
    if (!amapHolder.AMap) {
        return null;
    }

    const AMap = amapHolder.AMap;
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

    mapHolder.amapInstance = amapInstance;
    mapHolder.amapToolbar = amapToolbar;
    mapHolder.inited = true;
}

export function setAmapCenterTo(mapHolder, center, zoomLevel) {
    if (!amapHolder.AMap || !mapHolder.amapInstance) {
        return;
    }

    const AMap = amapHolder.AMap;

    if (amapHolder.amapCenterPosition
        && amapHolder.amapCenterPosition.originalLongitude === center.longitude
        && amapHolder.amapCenterPosition.originalLatitude === center.latitude
        && amapHolder.amapCenterPosition.convertedLongitude
        && amapHolder.amapCenterPosition.convertedLatitude
    ) {
        mapHolder.amapInstance.setZoomAndCenter(zoomLevel, new AMap.LngLat(amapHolder.amapCenterPosition.convertedLongitude, amapHolder.amapCenterPosition.convertedLatitude));
        return;
    }

    amapHolder.amapCenterPosition = {
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
            amapHolder.amapCenterPosition.convertedLongitude = convertedCenterPoint.getLng();
            amapHolder.amapCenterPosition.convertedLatitude = convertedCenterPoint.getLat();
        }

        mapHolder.amapInstance.setZoomAndCenter(zoomLevel, convertedCenterPoint);
    });
}

export function setAmapCenterMaker(mapHolder, position) {
    if (!amapHolder.AMap || !mapHolder.amapInstance) {
        return;
    }

    const AMap = amapHolder.AMap;
    const setMaker = function (point) {
        if (!mapHolder.amapCenterMarker) {
            mapHolder.amapCenterMarker = new AMap.Marker({
                position: point
            });
            mapHolder.amapInstance.add(mapHolder.amapCenterMarker);
        } else {
            mapHolder.amapCenterMarker.setPosition(point);
        }
    }

    if (amapHolder.amapCenterPosition
        && amapHolder.amapCenterPosition.originalLongitude === position.longitude
        && amapHolder.amapCenterPosition.originalLatitude === position.latitude
        && amapHolder.amapCenterPosition.convertedLongitude
        && amapHolder.amapCenterPosition.convertedLatitude
    ) {
        setMaker(new AMap.LngLat(amapHolder.amapCenterPosition.convertedLongitude, amapHolder.amapCenterPosition.convertedLatitude));
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

        setMaker(convertedMarkPoint);
    });
}

export function removeAmapCenterMaker(mapHolder) {
    if (!mapHolder.amapInstance || !mapHolder.amapCenterMarker) {
        return;
    }

    mapHolder.amapInstance.remove(mapHolder.amapCenterMarker);
    mapHolder.amapCenterMarker = null;
}
