const leafletTileSources = {
    'openstreetmap': {
        tileUrlFormat: 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
        tileUrlSubDomains: 'abc',
        attribution : '&copy; <a href="http://www.openstreetmap.org/copyright" class="external" target="_blank">OpenStreetMap</a>'
    }
}

export default {
    leafletTileSources: leafletTileSources
}
