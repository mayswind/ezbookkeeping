import services from "./services.js";
import settings from "./settings.js";

const leafletHolder = {
    leaflet: null
};

function loadLeafletMapAssets() {
    return Promise.all([
        import('leaflet/dist/leaflet.css'),
        import('leaflet/dist/leaflet-src.esm.js').then(leaflet => leafletHolder.leaflet = leaflet)
    ]);
}

function createLeafletMapInstance(mapHolder, mapContainer, options) {
    if (!leafletHolder.leaflet) {
        return null;
    }

    const leaflet = leafletHolder.leaflet;
    const leafletInstance = leaflet.map(mapContainer, {
        attributionControl: false,
        zoomControl: false
    });

    const mapTileImageUrl = services.generateOpenStreetMapTileImageUrl();

    const tileLayer = leaflet.tileLayer(mapTileImageUrl.url, {
        subdomains: mapTileImageUrl.subDomains,
        maxZoom: 19
    });
    tileLayer.addTo(leafletInstance);

    const zoomControl = leaflet.control.zoom({
        zoomInTitle: options.text.zoomIn,
        zoomOutTitle: options.text.zoomOut
    });
    zoomControl.addTo(leafletInstance);

    const attribution = leaflet.control.attribution({
        prefix: false
    });
    attribution.addAttribution('&copy; <a href="http://www.openstreetmap.org/copyright" class="external" target="_blank">OpenStreetMap</a>');
    attribution.addTo(leafletInstance);

    mapHolder.leafletInstance = leafletInstance;
    mapHolder.leafletTileLayer = tileLayer;
    mapHolder.leafletZoomControl = zoomControl;
    mapHolder.leafletAttribution = attribution;
    mapHolder.inited = true;
}

function setLeafletMapCenterTo(mapHolder, center, zoomLevel) {
    if (!mapHolder.leafletInstance) {
        return;
    }

    mapHolder.leafletInstance.setView([ center.latitude, center.longitude ], zoomLevel);
}

function setLeafletMapCenterMaker(mapHolder, position) {
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

function removeLeafletMapCenterMaker(mapHolder) {
    if (!mapHolder.leafletInstance || mapHolder.leafletCenterMarker) {
        mapHolder.leafletCenterMarker.remove();
        mapHolder.leafletCenterMarker = null;
    }
}

export function loadMapAssets() {
    if (settings.getMapProvider() === 'openstreetmap') {
        return loadLeafletMapAssets();
    }
}

export function createMapHolder() {
    if (settings.getMapProvider() === 'openstreetmap') {
        return {
            mapProvider: 'openstreetmap',
            dependencyLoaded: !!leafletHolder.leaflet,
            inited: false,
            defaultZoomLevel: 14,
            minZoomLevel: 1,
            leafletInstance: null,
            leafletTileLayer: null,
            leafletZoomControl: null,
            leafletAttribution: null,
            leafletCenterMarker: null
        }
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
    }
}

export function setMapCenterTo(mapHolder, center, zoomLevel) {
    if (!mapHolder) {
        return;
    }

    if (mapHolder.mapProvider === 'openstreetmap') {
        setLeafletMapCenterTo(mapHolder, center, zoomLevel);
    }
}

export function setMapCenterMarker(mapHolder, position) {
    if (!mapHolder) {
        return;
    }

    if (mapHolder.mapProvider === 'openstreetmap') {
        setLeafletMapCenterMaker(mapHolder, position);
    }
}

export function removeMapCenterMarker(mapHolder) {
    if (!mapHolder) {
        return;
    }

    if (mapHolder.mapProvider === 'openstreetmap') {
        removeLeafletMapCenterMaker(mapHolder);
    }
}
