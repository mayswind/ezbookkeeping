import type { Coordinate } from '@/core/coordinate.ts';

export interface MapProvider {
    getWebsite(): string;
    isSupportGetGeoLocationByClick(): boolean;
    asyncLoadAssets(language?: string): Promise<unknown>;
    createMapInstance(options: MapCreateOptions): MapInstance | null;
}

export interface MapInstance {
    dependencyLoaded: boolean;
    inited: boolean;
    initMapInstance(mapContainer: HTMLElement, options: MapInstanceInitOptions): void;
    getDefaultZoomLevel(): number;
    getMinZoomLevel(): number;
    getMaxZoomLevel(): number;
    getZoomLevel(): number;
    setMapCenterTo(center: Coordinate, zoomLevel: number): void;
    setMapCenterMarker(position: Coordinate): void;
    removeMapCenterMarker(): void;
    zoomIn(): void;
    zoomOut(): void;
}

export interface MapCreateOptions {
    readonly enableZoomControl?: boolean;
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
    readonly onZoomChange?: (level: number) => void;
}
