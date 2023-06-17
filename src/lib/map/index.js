import settings from "@/lib/settings.js";

import {
    loadLeafletMapAssets,
    createLeafletMapHolder,
    createLeafletMapInstance,
    setLeafletMapCenterTo,
    setLeafletMapCenterMaker,
    removeLeafletMapCenterMaker
} from './openstreetmap.js';

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

export function loadMapAssets(language) {
    if (settings.getMapProvider() === 'openstreetmap') {
        return loadLeafletMapAssets(language);
    } else if (settings.getMapProvider() === 'googlemap') {
        return loadGoogleMapAssets(language);
    } else if (settings.getMapProvider() === 'baidumap') {
        return loadBaiduMapAssets(language);
    }
}

export function createMapHolder() {
    if (settings.getMapProvider() === 'openstreetmap') {
        return createLeafletMapHolder();
    } else if (settings.getMapProvider() === 'googlemap') {
        return createGoogleMapHolder();
    } else if (settings.getMapProvider() === 'baidumap') {
        return createBaiduMapHolder();
    } else {
        return null;
    }
}

export function initMapInstance(mapHolder, mapContainer, options) {
    if (!mapHolder) {
        return;
    }

    if (mapHolder.mapProvider === 'openstreetmap') {
        createLeafletMapInstance(mapHolder, mapContainer, options);
    } else if (mapHolder.mapProvider === 'googlemap') {
        createGoogleMapInstance(mapHolder, mapContainer, options);
    } else if (mapHolder.mapProvider === 'baidumap') {
        createBaiduMapInstance(mapHolder, mapContainer, options);
    }
}

export function setMapCenterTo(mapHolder, center, zoomLevel) {
    if (!mapHolder) {
        return;
    }

    if (mapHolder.mapProvider === 'openstreetmap') {
        setLeafletMapCenterTo(mapHolder, center, zoomLevel);
    } else if (mapHolder.mapProvider === 'googlemap') {
        setGoogleMapCenterTo(mapHolder, center, zoomLevel);
    } else if (mapHolder.mapProvider === 'baidumap') {
        setBaiduMapCenterTo(mapHolder, center, zoomLevel);
    }
}

export function setMapCenterMarker(mapHolder, position) {
    if (!mapHolder) {
        return;
    }

    if (mapHolder.mapProvider === 'openstreetmap') {
        setLeafletMapCenterMaker(mapHolder, position);
    } else if (mapHolder.mapProvider === 'googlemap') {
        setGoogleMapCenterMaker(mapHolder, position);
    } else if (mapHolder.mapProvider === 'baidumap') {
        setBaiduMapCenterMaker(mapHolder, position);
    }
}

export function removeMapCenterMarker(mapHolder) {
    if (!mapHolder) {
        return;
    }

    if (mapHolder.mapProvider === 'openstreetmap') {
        removeLeafletMapCenterMaker(mapHolder);
    } else if (mapHolder.mapProvider === 'googlemap') {
        removeGoogleMapCenterMaker(mapHolder);
    } else if (mapHolder.mapProvider === 'baidumap') {
        removeBaiduMapCenterMaker(mapHolder);
    }
}
