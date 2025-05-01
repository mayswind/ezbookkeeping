import { LEAFLET_TILE_SOURCES } from '@/consts/map.ts';
import { getMapProvider } from '@/lib/server_settings.ts';

import type { MapProvider, MapInstance } from './base.ts';
import { LeafletMapProvider } from './leaflet.ts';
import { GoogleMapProvider } from './googlemap.ts';
import { BaiduMapProvider } from './baidumap.ts';
import { AmapMapProvider } from './amap.ts';

let mapProvider: MapProvider | null = null;

export function initMapProvider(language?: string): void {
    const mapProviderType = getMapProvider();

    if (LEAFLET_TILE_SOURCES[mapProviderType] || mapProviderType === 'custom') {
        mapProvider = new LeafletMapProvider(mapProviderType);
    } else if (mapProviderType === 'googlemap') {
        mapProvider = new GoogleMapProvider();
    } else if (mapProviderType === 'baidumap') {
        mapProvider = new BaiduMapProvider();
    } else if (mapProviderType === 'amap') {
        mapProvider = new AmapMapProvider();
    }

    if (mapProvider) {
        mapProvider.asyncLoadAssets(language);
    }
}

export function getMapWebsite(): string {
    return mapProvider?.getWebsite() || '';
}

export function isSupportGetGeoLocationByClick(): boolean {
    return mapProvider?.isSupportGetGeoLocationByClick() || false;
}

export function createMapInstance(): MapInstance | null {
    return mapProvider?.createMapInstance() || null;
}
