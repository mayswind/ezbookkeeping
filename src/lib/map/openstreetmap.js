import services from "@/lib/services.js";

const leafletHolder = {
    leaflet: null
};

export function loadLeafletMapAssets() {
    return Promise.all([
        import('leaflet/dist/leaflet.css'),
        import('leaflet/dist/leaflet-src.esm.js').then(leaflet => leafletHolder.leaflet = leaflet)
    ]);
}

export function createLeafletMapHolder() {
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
