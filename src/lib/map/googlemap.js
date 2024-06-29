import { asyncLoadAssets } from '@/lib/misc.js';
import services from '@/lib/services.js';

const googleMapHolder = {
    googleMap: null,
    ControlPosition: {
        LEFT_TOP: (window.google && window.google.maps && window.google.maps.ControlPosition) ? window.google.maps.ControlPosition.LEFT_TOP : 5
    }
};

export function getGoogleMapWebsite() {
    return 'https://maps.google.com';
}

export function loadGoogleMapAssets(language) {
    if (googleMapHolder.googleMap) {
        return;
    }

    if (!window.onGoogleMapCallback) {
        window.onGoogleMapCallback = () => {
            if (window.google) {
                googleMapHolder.googleMap = window.google.maps;
            }
        };
    }

    return asyncLoadAssets('js', services.generateGoogleMapJavascriptUrl(language, 'onGoogleMapCallback'));
}

export function createGoogleMapHolder() {
    return {
        mapProvider: 'googlemap',
        dependencyLoaded: !!googleMapHolder.googleMap,
        inited: false,
        defaultZoomLevel: 14,
        minZoomLevel: 1,
        googleMapInstance: null,
        googleMapCenterMarker: null
    };
}

export function createGoogleMapInstance(mapHolder, mapContainer, options) {
    if (!googleMapHolder.googleMap) {
        return null;
    }

    const googleMap = googleMapHolder.googleMap;

    mapHolder.googleMapInstance = new googleMap.Map(mapContainer, {
        zoom: options.zoomLevel,
        center: {
            lat: options.initCenter.latitude,
            lng: options.initCenter.longitude
        },
        maxZoom: 19,
        zoomControl: true,
        mapTypeControl: false,
        scaleControl: false,
        streetViewControl: false,
        rotateControl: false,
        fullscreenControl: false,
        gestureHandling: 'greedy',
        zoomControlOptions: {
            position: googleMapHolder.ControlPosition.LEFT_TOP
        }
    });
    mapHolder.inited = true;
}

export function setGoogleMapCenterTo(mapHolder, center, zoomLevel) {
    if (!googleMapHolder.googleMap || !mapHolder.googleMapInstance) {
        return;
    }

    mapHolder.googleMapInstance.setCenter({
        lat: center.latitude,
        lng: center.longitude
    });
    mapHolder.googleMapInstance.setZoom(zoomLevel);
}

export function setGoogleMapCenterMaker(mapHolder, position) {
    if (!googleMapHolder.googleMap || !mapHolder.googleMapInstance) {
        return;
    }

    const googleMap = googleMapHolder.googleMap;

    if (!mapHolder.googleMapCenterMarker) {
        mapHolder.googleMapCenterMarker = new googleMap.Marker({
            position: {
                lat: position.latitude,
                lng: position.longitude
            },
            map: mapHolder.googleMapInstance
        });
    } else {
        mapHolder.googleMapCenterMarker.setPosition({
            lat: position.latitude,
            lng: position.longitude
        });
    }
}

export function removeGoogleMapCenterMaker(mapHolder) {
    if (!mapHolder.googleMapInstance || !mapHolder.googleMapCenterMarker) {
        return;
    }

    mapHolder.googleMapCenterMarker.setMap(null);
    mapHolder.googleMapCenterMarker = null;
}
