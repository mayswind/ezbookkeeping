export interface BrowserCacheStatistics {
    readonly totalCacheSize: number;
    readonly codeCacheSize: number;
    readonly assetsCacheSize: number;
    readonly mapCacheSize: number;
    readonly othersCacheSize: number;
}

export interface SWMapCacheConfig {
    enabled: boolean;
    patterns: string[];
    maxEntries: number;
    maxAgeMilliseconds: number;
}
