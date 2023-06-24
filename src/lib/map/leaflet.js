import mapConstants from '@/consts/map.js';
import { isMapDataFetchProxyEnabled, getTomTomMapAPIKey } from '@/lib/server_settings.js';
import services from '@/lib/services.js';

const leafletHolder = {
    leaflet: null
};

export function loadLeafletMapAssets() {
    return Promise.all([
        import('leaflet/dist/leaflet.css'),
        import('leaflet/dist/leaflet-src.esm.js').then(leaflet => leafletHolder.leaflet = leaflet)
    ]);
}

export function createLeafletMapHolder(mapProvider) {
    const mapTileSource = mapConstants.leafletTileSources[mapProvider];

    if (!mapTileSource) {
        return null;
    }

    return {
        mapProvider: mapProvider,
        dependencyLoaded: !!leafletHolder.leaflet,
        inited: false,
        defaultZoomLevel: mapTileSource.defaultZoomLevel,
        minZoomLevel: mapTileSource.minZoom,
        leafletInstance: null,
        leafletTileLayer: null,
        leafletZoomControl: null,
        leafletAttribution: null,
        leafletCenterMarker: null
    };
}

export function createLeafletMapInstance(mapHolder, mapContainer, options) {
    if (!leafletHolder.leaflet) {
        return null;
    }

    const leaflet = leafletHolder.leaflet;
    const leafletInstance = leaflet.map(mapContainer, {
        center: [ options.initCenter.latitude, options.initCenter.longitude ],
        zoom: options.zoomLevel,
        attributionControl: false,
        zoomControl: false
    });
    let mapTileSource = Object.assign({}, mapConstants.leafletTileSources[mapHolder.mapProvider]);

    if (isMapDataFetchProxyEnabled()) {
        mapTileSource.tileUrlFormat = services.generateMapProxyTileImageUrl(mapHolder.mapProvider, options.language);
        mapTileSource.tileUrlSubDomains = '';
    } else if (mapTileSource.tileUrlExtraParams) {
        const params = [];

        for (let i = 0; i < mapTileSource.tileUrlExtraParams.length; i++) {
            const param = mapTileSource.tileUrlExtraParams[i];

            if (param.paramValueType === 'tomtom_key') {
                params.push('key=' + getTomTomMapAPIKey());
            } else if (param.paramValueType === 'language' && options.language) {
                params.push('language=' + options.language);
            }
        }

        if (params.length) {
            mapTileSource.tileUrlFormat = mapTileSource.tileUrlFormat + '?' + params.join('&');
        }
    }

    const tileLayer = leaflet.tileLayer(mapTileSource.tileUrlFormat, {
        subdomains: mapTileSource.tileUrlSubDomains,
        maxZoom: mapTileSource.maxZoom
    });
    tileLayer.addTo(leafletInstance);

    const zoomControl = leaflet.control.zoom({
        zoomInTitle: options.text.zoomIn,
        zoomOutTitle: options.text.zoomOut
    });
    zoomControl.addTo(leafletInstance);

    if (mapTileSource.attribution) {
        const attribution = leaflet.control.attribution({
            prefix: false
        });
        attribution.addAttribution(mapTileSource.attribution);
        attribution.addTo(leafletInstance);
        mapHolder.leafletAttribution = attribution;
    }

    mapHolder.leafletInstance = leafletInstance;
    mapHolder.leafletTileLayer = tileLayer;
    mapHolder.leafletZoomControl = zoomControl;
    mapHolder.inited = true;
}

export function setLeafletMapCenterTo(mapHolder, center, zoomLevel) {
    if (!mapHolder.leafletInstance) {
        return;
    }

    mapHolder.leafletInstance.setView([ center.latitude, center.longitude ], zoomLevel);
}

export function setLeafletMapCenterMaker(mapHolder, position) {
    if (!leafletHolder.leaflet || !mapHolder.leafletInstance) {
        return;
    }

    const leaflet = leafletHolder.leaflet;

    if (!mapHolder.leafletCenterMarker) {
        const markerIcon = leaflet.icon({
            iconUrl: 'img/map-marker-icon.png',
            iconRetinaUrl: 'img/map-marker-icon-2x.png',
            iconSize: [25, 32],
            iconAnchor: [12, 32],
            shadowUrl: 'img/map-marker-shadow.png',
            shadowSize: [41, 32]
        });
        mapHolder.leafletCenterMarker = leaflet.marker([ position.latitude, position.longitude ], {
            icon: markerIcon
        });
        mapHolder.leafletCenterMarker.addTo(mapHolder.leafletInstance);
    } else {
        mapHolder.leafletCenterMarker.setLatLng([ position.latitude, position.longitude ]);
    }
}

export function removeLeafletMapCenterMaker(mapHolder) {
    if (!mapHolder.leafletInstance || !mapHolder.leafletCenterMarker) {
        return;
    }

    mapHolder.leafletCenterMarker.remove();
    mapHolder.leafletCenterMarker = null;
}
