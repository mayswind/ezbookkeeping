export interface PresetLeafletTileSource {
    readonly tileUrlFormat: string;
    readonly tileUrlSubDomains: string;
    readonly tileUrlExtraParams?: PresetLeafletTileSourceExtraParam[];
    readonly annotationUrlFormat?: string;
    readonly annotationUrlSubDomains?: string;
    readonly annotationUrlExtraParams?: PresetLeafletTileSourceExtraParam[];
    readonly minZoom: number;
    readonly maxZoom: number;
    readonly defaultZoomLevel: number;
    readonly website: string;
    readonly attribution: string;
}

export interface PresetLeafletTileSourceExtraParam {
    readonly paramName: string;
    readonly paramValueType: string;
}

export const LEAFLET_TILE_SOURCES: Record<string, PresetLeafletTileSource> = {
    'openstreetmap': {
        tileUrlFormat: 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
        tileUrlSubDomains: 'abc',
        minZoom: 1,
        maxZoom: 19,
        defaultZoomLevel: 14,
        website: 'https://www.openstreetmap.org',
        attribution : '&copy; <a href="https://www.openstreetmap.org/copyright" class="external" target="_blank">OpenStreetMap</a> contributors'
    },
    'openstreetmap-humanitarian': {
        tileUrlFormat: 'https://{s}.tile.openstreetmap.fr/hot/{z}/{x}/{y}.png',
        tileUrlSubDomains: 'abc',
        minZoom: 1,
        maxZoom: 19,
        defaultZoomLevel: 14,
        website: 'https://www.hotosm.org',
        attribution : '&copy; <a href="https://www.openstreetmap.org/copyright" class="external" target="_blank">OpenStreetMap</a> contributors, Tiles style by <a href="https://www.hotosm.org/" class="external" target="_blank">Humanitarian OpenStreetMap Team</a> hosted by <a href="https://openstreetmap.fr/" class="external" target="_blank">OpenStreetMap France</a>'
    },
    'opentopomap': {
        tileUrlFormat: 'https://{s}.tile.opentopomap.org/{z}/{x}/{y}.png',
        tileUrlSubDomains: 'abc',
        minZoom: 1,
        maxZoom: 17,
        defaultZoomLevel: 14,
        website: 'https://opentopomap.org',
        attribution : 'Map data: &copy; <a href="https://www.openstreetmap.org/copyright" class="external" target="_blank">OpenStreetMap</a> contributors, <a href="http://viewfinderpanoramas.org" class="external" target="_blank">SRTM</a> | Map style: &copy; <a href="https://opentopomap.org" class="external" target="_blank">OpenTopoMap</a> (<a href="https://creativecommons.org/licenses/by-sa/3.0/" class="external" target="_blank">CC-BY-SA</a>)'
    },
    'opnvkarte': {
        tileUrlFormat: 'https://tileserver.memomaps.de/tilegen/{z}/{x}/{y}.png',
        tileUrlSubDomains: '',
        minZoom: 1,
        maxZoom: 17,
        defaultZoomLevel: 14,
        website: 'https://memomaps.de',
        attribution : 'Map <a href="https://memomaps.de/" class="external" target="_blank">memomaps.de</a> <a href="http://creativecommons.org/licenses/by-sa/2.0/" class="external" target="_blank">CC-BY-SA</a>, map data &copy; <a href="https://www.openstreetmap.org/copyright" class="external" target="_blank">OpenStreetMap</a> contributors'
    },
    'cyclosm': {
        tileUrlFormat: 'https://{s}.tile-cyclosm.openstreetmap.fr/cyclosm/{z}/{x}/{y}.png',
        tileUrlSubDomains: 'abc',
        minZoom: 1,
        maxZoom: 19,
        defaultZoomLevel: 14,
        website: 'https://github.com/cyclosm/cyclosm-cartocss-style',
        attribution : '<a href="https://github.com/cyclosm/cyclosm-cartocss-style/releases" title="CyclOSM - Open Bicycle render" class="external" target="_blank">CyclOSM</a> | Map data: &copy; <a href="https://www.openstreetmap.org/copyright" class="external" target="_blank">OpenStreetMap</a> contributors'
    },
    'cartodb': {
        tileUrlFormat: 'https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}.png',
        tileUrlSubDomains: 'abcd',
        minZoom: 1,
        maxZoom: 20,
        defaultZoomLevel: 14,
        website: 'https://carto.com',
        attribution : '&copy; <a href="http://www.openstreetmap.org/copyright" class="external" target="_blank">OpenStreetMap</a>, &copy; <a href="https://carto.com/attributions" class="external" target="_blank">CARTO</a>'
    },
    'tomtom': {
        tileUrlFormat: 'https://{s}.api.tomtom.com/map/1/tile/basic/main/{z}/{x}/{y}.png',
        tileUrlSubDomains: 'abcd',
        tileUrlExtraParams: [
            {
                paramName: 'key',
                paramValueType: 'tomtom_key'
            },
            {
                paramName: 'language',
                paramValueType: 'language'
            }
        ],
        minZoom: 1,
        maxZoom: 19,
        defaultZoomLevel: 14,
        website: 'https://tomtom.com',
        attribution : '<a href="https://tomtom.com" class="external" target="_blank">&copy;  1992 - 2023 TomTom.</a>'
    },
    'tianditu': {
        tileUrlFormat: 'https://t{s}.tianditu.gov.cn/vec_w/wmts?SERVICE=WMTS&REQUEST=GetTile&VERSION=1.0.0&LAYER=vec&STYLE=default&TILEMATRIXSET=w&FORMAT=tiles&TILEMATRIX={z}&TILEROW={y}&TILECOL={x}',
        tileUrlSubDomains: '01234567',
        tileUrlExtraParams: [
            {
                paramName: 'tk',
                paramValueType: 'tianditu_key'
            }
        ],
        annotationUrlFormat: 'https://t{s}.tianditu.gov.cn/cva_w/wmts?SERVICE=WMTS&REQUEST=GetTile&VERSION=1.0.0&LAYER=cva&STYLE=default&TILEMATRIXSET=w&FORMAT=tiles&TILEMATRIX={z}&TILEROW={y}&TILECOL={x}',
        annotationUrlSubDomains: '01234567',
        annotationUrlExtraParams: [
            {
                paramName: 'tk',
                paramValueType: 'tianditu_key'
            }
        ],
        minZoom: 1,
        maxZoom: 18,
        defaultZoomLevel: 14,
        website: 'https://www.tianditu.gov.cn',
        attribution : '<a href="https://www.tianditu.gov.cn" class="external" target="_blank">天地图</a>'
    }
}
