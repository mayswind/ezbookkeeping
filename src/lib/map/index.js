import mapConstants from '@/consts/map.js';
import settings from '@/lib/settings.js';

import {
    loadLeafletMapAssets,
    createLeafletMapHolder,
    createLeafletMapInstance,
    setLeafletMapCenterTo,
    setLeafletMapCenterMaker,
    removeLeafletMapCenterMaker
} from './leaflet.js';

import {
    loadGoogleMapAssets,
    createGoogleMapHolder,
    createGoogleMapInstance,
    setGoogleMapCenterTo,
    setGoogleMapCenterMaker,
    removeGoogleMapCenterMaker
} from './googlemap.js';

import {
    loadBaiduMapAssets,
    createBaiduMapHolder,
    createBaiduMapInstance,
    setBaiduMapCenterTo,
    setBaiduMapCenterMaker,
    removeBaiduMapCenterMaker
} from './baidumap.js';

import {
    loadAmapAssets,
    createAmapHolder,
    createAmapInstance,
    setAmapCenterTo,
    setAmapCenterMaker,
    removeAmapCenterMaker
} from './amap.js';

export function loadMapAssets(language) {
    if (mapConstants.leafletTileSources[settings.getMapProvider()]) {
        return loadLeafletMapAssets(language);
    } else if (settings.getMapProvider() === 'googlemap') {
        return loadGoogleMapAssets(language);
    } else if (settings.getMapProvider() === 'baidumap') {
        return loadBaiduMapAssets(language);
    } else if (settings.getMapProvider() === 'amap') {
        return loadAmapAssets(language);
    }
}

export function createMapHolder() {
    if (mapConstants.leafletTileSources[settings.getMapProvider()]) {
        return createLeafletMapHolder(settings.getMapProvider());
    } else if (settings.getMapProvider() === 'googlemap') {
        return createGoogleMapHolder(settings.getMapProvider());
    } else if (settings.getMapProvider() === 'baidumap') {
        return createBaiduMapHolder(settings.getMapProvider());
    } else if (settings.getMapProvider() === 'amap') {
        return createAmapHolder(settings.getMapProvider());
    } else {
        return null;
    }
}

export function initMapInstance(mapHolder, mapContainer, options) {
    if (!mapHolder) {
        return;
    }

    if (mapConstants.leafletTileSources[settings.getMapProvider()]) {
        createLeafletMapInstance(mapHolder, mapContainer, options);
    } else if (mapHolder.mapProvider === 'googlemap') {
        createGoogleMapInstance(mapHolder, mapContainer, options);
    } else if (mapHolder.mapProvider === 'baidumap') {
        createBaiduMapInstance(mapHolder, mapContainer, options);
    } else if (mapHolder.mapProvider === 'amap') {
        createAmapInstance(mapHolder, mapContainer, options);
    }
}

export function setMapCenterTo(mapHolder, center, zoomLevel) {
    if (!mapHolder) {
        return;
    }

    if (mapConstants.leafletTileSources[settings.getMapProvider()]) {
        setLeafletMapCenterTo(mapHolder, center, zoomLevel);
    } else if (mapHolder.mapProvider === 'googlemap') {
        setGoogleMapCenterTo(mapHolder, center, zoomLevel);
    } else if (mapHolder.mapProvider === 'baidumap') {
        setBaiduMapCenterTo(mapHolder, center, zoomLevel);
    } else if (mapHolder.mapProvider === 'amap') {
        setAmapCenterTo(mapHolder, center, zoomLevel);
    }
}

export function setMapCenterMarker(mapHolder, position) {
    if (!mapHolder) {
        return;
    }

    if (mapConstants.leafletTileSources[settings.getMapProvider()]) {
        setLeafletMapCenterMaker(mapHolder, position);
    } else if (mapHolder.mapProvider === 'googlemap') {
        setGoogleMapCenterMaker(mapHolder, position);
    } else if (mapHolder.mapProvider === 'baidumap') {
        setBaiduMapCenterMaker(mapHolder, position);
    } else if (mapHolder.mapProvider === 'amap') {
        setAmapCenterMaker(mapHolder, position);
    }
}

export function removeMapCenterMarker(mapHolder) {
    if (!mapHolder) {
        return;
    }

    if (mapConstants.leafletTileSources[settings.getMapProvider()]) {
        removeLeafletMapCenterMaker(mapHolder);
    } else if (mapHolder.mapProvider === 'googlemap') {
        removeGoogleMapCenterMaker(mapHolder);
    } else if (mapHolder.mapProvider === 'baidumap') {
        removeBaiduMapCenterMaker(mapHolder);
    } else if (mapHolder.mapProvider === 'amap') {
        removeAmapCenterMaker(mapHolder);
    }
}
