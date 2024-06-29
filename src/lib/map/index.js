import mapConstants from '@/consts/map.js';
import {
    getMapProvider
} from '@/lib/server_settings.js';

import {
    loadLeafletMapAssets,
    createLeafletMapHolder,
    createLeafletMapInstance,
    setLeafletMapCenterTo,
    setLeafletMapCenterMaker,
    removeLeafletMapCenterMaker
} from './leaflet.js';

import {
    getGoogleMapWebsite,
    loadGoogleMapAssets,
    createGoogleMapHolder,
    createGoogleMapInstance,
    setGoogleMapCenterTo,
    setGoogleMapCenterMaker,
    removeGoogleMapCenterMaker
} from './googlemap.js';

import {
    getBaiduMapWebsite,
    loadBaiduMapAssets,
    createBaiduMapHolder,
    createBaiduMapInstance,
    setBaiduMapCenterTo,
    setBaiduMapCenterMaker,
    removeBaiduMapCenterMaker
} from './baidumap.js';

import {
    getAmapWebsite,
    loadAmapAssets,
    createAmapHolder,
    createAmapInstance,
    setAmapCenterTo,
    setAmapCenterMaker,
    removeAmapCenterMaker
} from './amap.js';

export function getMapWebsite() {
    if (getMapProvider() === 'custom') {
        return '';
    } else if (mapConstants.leafletTileSources[getMapProvider()]) {
        return mapConstants.leafletTileSources[getMapProvider()].website;
    } else if (getMapProvider() === 'googlemap') {
        return getGoogleMapWebsite();
    } else if (getMapProvider() === 'baidumap') {
        return getBaiduMapWebsite();
    } else if (getMapProvider() === 'amap') {
        return getAmapWebsite();
    }
}

export function loadMapAssets(language) {
    if (mapConstants.leafletTileSources[getMapProvider()] || getMapProvider() === 'custom') {
        return loadLeafletMapAssets(language);
    } else if (getMapProvider() === 'googlemap') {
        return loadGoogleMapAssets(language);
    } else if (getMapProvider() === 'baidumap') {
        return loadBaiduMapAssets(language);
    } else if (getMapProvider() === 'amap') {
        return loadAmapAssets(language);
    }
}

export function createMapHolder() {
    if (mapConstants.leafletTileSources[getMapProvider()] || getMapProvider() === 'custom') {
        return createLeafletMapHolder(getMapProvider());
    } else if (getMapProvider() === 'googlemap') {
        return createGoogleMapHolder(getMapProvider());
    } else if (getMapProvider() === 'baidumap') {
        return createBaiduMapHolder(getMapProvider());
    } else if (getMapProvider() === 'amap') {
        return createAmapHolder(getMapProvider());
    } else {
        return null;
    }
}

export function initMapInstance(mapHolder, mapContainer, options) {
    if (!mapHolder) {
        return;
    }

    if (mapConstants.leafletTileSources[getMapProvider()] || getMapProvider() === 'custom') {
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

    if (mapConstants.leafletTileSources[getMapProvider()] || getMapProvider() === 'custom') {
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

    if (mapConstants.leafletTileSources[getMapProvider()] || getMapProvider() === 'custom') {
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

    if (mapConstants.leafletTileSources[getMapProvider()] || getMapProvider() === 'custom') {
        removeLeafletMapCenterMaker(mapHolder);
    } else if (mapHolder.mapProvider === 'googlemap') {
        removeGoogleMapCenterMaker(mapHolder);
    } else if (mapHolder.mapProvider === 'baidumap') {
        removeBaiduMapCenterMaker(mapHolder);
    } else if (mapHolder.mapProvider === 'amap') {
        removeAmapCenterMaker(mapHolder);
    }
}
