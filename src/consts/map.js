const leafletTileSources = {
    'openstreetmap': {
        tileUrlFormat: 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
        tileUrlSubDomains: 'abc',
        minZoom: 1,
        maxZoom: 19,
        defaultZoomLevel: 14,
        attribution : '&copy; <a href="https://www.openstreetmap.org/copyright" class="external" target="_blank">OpenStreetMap</a> contributors'
    },
    'openstreetmap-humanitarian': {
        tileUrlFormat: 'https://{s}.tile.openstreetmap.fr/hot/{z}/{x}/{y}.png',
        tileUrlSubDomains: 'abc',
        minZoom: 1,
        maxZoom: 19,
        defaultZoomLevel: 14,
        attribution : '&copy; <a href="https://www.openstreetmap.org/copyright" class="external" target="_blank">OpenStreetMap</a> contributors, Tiles style by <a href="https://www.hotosm.org/" class="external" target="_blank">Humanitarian OpenStreetMap Team</a> hosted by <a href="https://openstreetmap.fr/" class="external" target="_blank">OpenStreetMap France</a>'
    },
    'opentopomap': {
        tileUrlFormat: 'https://{s}.tile.opentopomap.org/{z}/{x}/{y}.png',
        tileUrlSubDomains: 'abc',
        minZoom: 1,
        maxZoom: 17,
        defaultZoomLevel: 14,
        attribution : 'Map data: &copy; <a href="https://www.openstreetmap.org/copyright" class="external" target="_blank">OpenStreetMap</a> contributors, <a href="http://viewfinderpanoramas.org" class="external" target="_blank">SRTM</a> | Map style: &copy; <a href="https://opentopomap.org" class="external" target="_blank">OpenTopoMap</a> (<a href="https://creativecommons.org/licenses/by-sa/3.0/" class="external" target="_blank">CC-BY-SA</a>)'
    },
    'opnvkarte': {
        tileUrlFormat: 'https://tileserver.memomaps.de/tilegen/{z}/{x}/{y}.png',
        tileUrlSubDomains: '',
        minZoom: 1,
        maxZoom: 17,
        defaultZoomLevel: 14,
        attribution : 'Map <a href="https://memomaps.de/" class="external" target="_blank">memomaps.de</a> <a href="http://creativecommons.org/licenses/by-sa/2.0/" class="external" target="_blank">CC-BY-SA</a>, map data &copy; <a href="https://www.openstreetmap.org/copyright" class="external" target="_blank">OpenStreetMap</a> contributors'
    },
    'cyclosm': {
        tileUrlFormat: 'https://{s}.tile-cyclosm.openstreetmap.fr/cyclosm/{z}/{x}/{y}.png',
        tileUrlSubDomains: 'abc',
        minZoom: 1,
        maxZoom: 19,
        defaultZoomLevel: 14,
        attribution : '<a href="https://github.com/cyclosm/cyclosm-cartocss-style/releases" title="CyclOSM - Open Bicycle render" class="external" target="_blank">CyclOSM</a> | Map data: &copy; <a href="https://www.openstreetmap.org/copyright" class="external" target="_blank">OpenStreetMap</a> contributors'
    }
}

export default {
    leafletTileSources: leafletTileSources
}
