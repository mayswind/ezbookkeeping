import type { Coordinate } from '@/core/coordinate.ts';

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
    setMapCenterTo(center: Coordinate, zoomLevel: number): void;
    setMapCenterMarker(position: Coordinate): void;
    removeMapCenterMarker(): void;
}

export interface MapInstanceInitOptions {
    readonly language?: string;
    readonly initCenter: Coordinate;
    readonly zoomLevel: number;
    readonly text: {
        readonly zoomIn: string;
        readonly zoomOut: string;
    };
    readonly onClick?: (position: Coordinate) => void;
}
