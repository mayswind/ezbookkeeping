import type { KeyAndName } from './base.ts';

export interface BrowserCacheStatistics {
    readonly totalCacheSize: number;
    readonly codeCacheSize: number;
    readonly assetsCacheSize: number;
    readonly mapCacheSize: number;
    readonly userAvatarCacheSize: number;
    readonly transactionPictureCacheSize: number;
    readonly customIconCacheSize: number;
    readonly othersCacheSize: number;
}

export interface SWMapCacheConfig {
    enabled: boolean;
    patterns: string[];
    maxEntries: number;
    maxAgeMilliseconds: number;
}

export interface SWTransactionPictureCacheConfig {
    enabled: boolean;
    maxEntries: number;
    maxAgeMilliseconds: number;
}

export class FileCacheType implements KeyAndName {
    private static readonly allInstances: FileCacheType[] = [];
    private static readonly allInstancesByKey: Record<string, FileCacheType> = {};

    public static readonly ApplicationCode = new FileCacheType('ezbookkeeping-code-cache', 'Application Code', ['workbox-runtime-']);
    public static readonly ResourceFiles = new FileCacheType('ezbookkeeping-assets-cache', 'Resource Files', ['workbox-precache-v2-']);
    public static readonly MapData = new FileCacheType('ezbookkeeping-map-cache', 'Map Data', []);
    public static readonly UserAvatars = new FileCacheType('ezbookkeeping-user-avatar-cache', 'User Avatars', []);
    public static readonly TransactionPictures = new FileCacheType('ezbookkeeping-transaction-picture-cache', 'Transaction Pictures', []);
    public static readonly CustomIcons = new FileCacheType('ezbookkeeping-custom-icon-cache', 'Custom Icons', []);

    public readonly key: string;
    public readonly name: string;
    public readonly cacheNamePrefixes: string[];

    private constructor(key: string, name: string, cacheNamePrefixes: string[]) {
        this.key = key;
        this.name = name;
        this.cacheNamePrefixes = cacheNamePrefixes;

        FileCacheType.allInstances.push(this);
        FileCacheType.allInstancesByKey[key] = this;
    }

    public static values(): FileCacheType[] {
        return FileCacheType.allInstances;
    }

    public static valueOf(key: string): FileCacheType | undefined {
        return FileCacheType.allInstancesByKey[key];
    }
}
