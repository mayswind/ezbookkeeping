import type { MapPosition } from '@/core/map.ts';

export interface MapProvider {
    getWebsite(): string;
    isSupportGetGeoLocationByClick(): boolean;
    asyncLoadAssets(language?: string): Promise<unknown>;
    createMapInstance(): MapInstance | null;
}

export interface MapInstance {
    dependencyLoaded: boolean;
    inited: boolean;
    readonly defaultZoomLevel: number;
    readonly minZoomLevel: number;
    initMapInstance(mapContainer: HTMLElement, options: MapInstanceInitOptions): void;
    setMapCenterTo(center: MapPosition, zoomLevel: number): void;
    setMapCenterMarker(position: MapPosition): void;
    removeMapCenterMarker(): void;
}

export interface MapInstanceInitOptions {
    readonly language?: string;
    readonly initCenter: MapPosition;
    readonly zoomLevel: number;
    readonly text: {
        readonly zoomIn: string;
        readonly zoomOut: string;
    };
    readonly onClick?: (position: MapPosition) => void;
}
