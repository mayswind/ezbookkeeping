export interface MapProvider {
    getWebsite(): string;
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
    }
}

export interface MapPosition {
    latitude: number;
    longitude: number;
}
