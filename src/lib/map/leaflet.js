import { LEAFLET_TILE_SOURCES } from '@/consts/map.ts';
import {
    isMapDataFetchProxyEnabled,
    getCustomMapTileLayerUrl,
    getCustomMapAnnotationLayerUrl,
    isCustomMapAnnotationLayerDataFetchProxyEnabled,
    getCustomMapMinZoomLevel,
    getCustomMapMaxZoomLevel,
    getCustomMapDefaultZoomLevel,
    getTomTomMapAPIKey,
    getTianDiTuMapAPIKey
} from '@/lib/server_settings.ts';
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
    const mapTileSource = LEAFLET_TILE_SOURCES[mapProvider];

    if (mapProvider !== 'custom' && !mapTileSource) {
        return null;
    }

    return {
        mapProvider: mapProvider,
        dependencyLoaded: !!leafletHolder.leaflet,
        inited: false,
        defaultZoomLevel: mapProvider !== 'custom' ? mapTileSource.defaultZoomLevel : getCustomMapDefaultZoomLevel(),
        minZoomLevel: mapProvider !== 'custom' ? mapTileSource.minZoom : getCustomMapMinZoomLevel(),
        leafletInstance: null,
        leafletTileLayer: null,
        leafletAnnotationLayer: null,
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
    let mapTileSource = null;

    if (mapHolder.mapProvider !== 'custom') {
        mapTileSource = Object.assign({}, LEAFLET_TILE_SOURCES[mapHolder.mapProvider]);
    } else {
        mapTileSource = createCustomMapSource();
    }

    if (isMapDataFetchProxyEnabled()) {
        mapTileSource.tileUrlFormat = services.generateMapProxyTileImageUrl(mapHolder.mapProvider, options.language);
        mapTileSource.tileUrlSubDomains = '';
    } else if (mapTileSource.tileUrlExtraParams) {
        mapTileSource.tileUrlFormat = getFinalUrlFormat(mapTileSource.tileUrlFormat, mapTileSource.tileUrlExtraParams, options);
    }

    const tileLayer = leaflet.tileLayer(mapTileSource.tileUrlFormat, {
        subdomains: mapTileSource.tileUrlSubDomains,
        maxZoom: mapTileSource.maxZoom,
        minZoom: mapTileSource.minZoom
    });
    tileLayer.addTo(leafletInstance);

    if (mapTileSource.annotationUrlFormat || (mapHolder.mapProvider === 'custom' && isCustomMapAnnotationLayerDataFetchProxyEnabled())) {
        if (isMapDataFetchProxyEnabled()) {
            mapTileSource.annotationUrlFormat = services.generateMapProxyAnnotationImageUrl(mapHolder.mapProvider, options.language);
            mapTileSource.annotationUrlSubDomains = '';
        } else if (mapTileSource.annotationUrlExtraParams) {
            mapTileSource.annotationUrlFormat = getFinalUrlFormat(mapTileSource.annotationUrlFormat, mapTileSource.annotationUrlExtraParams, options);
        }

        const annotationLayer = leaflet.tileLayer(mapTileSource.annotationUrlFormat, {
            subdomains: mapTileSource.annotationUrlSubDomains,
            maxZoom: mapTileSource.maxZoom,
            minZoom: mapTileSource.minZoom
        });
        annotationLayer.addTo(leafletInstance);

        mapHolder.leafletAnnotationLayer = annotationLayer;
    }

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

function createCustomMapSource() {
    return {
        tileUrlFormat: getCustomMapTileLayerUrl(),
        tileUrlSubDomains: '',
        annotationUrlFormat: getCustomMapAnnotationLayerUrl(),
        annotationUrlSubDomains: '',
        minZoom: getCustomMapMinZoomLevel(),
        maxZoom: getCustomMapMaxZoomLevel(),
        defaultZoomLevel: getCustomMapDefaultZoomLevel()
    };
}

function getFinalUrlFormat(urlFormat, urlExtraParams, options) {
    const params = [];

    for (let i = 0; i < urlExtraParams.length; i++) {
        const param = urlExtraParams[i];

        if (param.paramValueType === 'tomtom_key') {
            params.push(param.paramName + '=' + getTomTomMapAPIKey());
        } else if (param.paramValueType === 'tianditu_key') {
            params.push(param.paramName + '=' + getTianDiTuMapAPIKey());
        } else if (param.paramValueType === 'language' && options.language) {
            params.push(param.paramName + '=' + options.language);
        }
    }

    if (params.length) {
        if (urlFormat.indexOf('?') >= 0) {
            urlFormat = urlFormat + '&' + params.join('&');
        } else {
            urlFormat = urlFormat + '?' + params.join('&');
        }
    }

    return urlFormat;
}
