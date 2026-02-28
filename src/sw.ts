import { clientsClaim } from 'workbox-core';
import type {
    WorkboxPlugin,
    CacheDidUpdateCallbackParam,
    CacheKeyWillBeUsedCallbackParam,
    CacheWillUpdateCallbackParam,
    CachedResponseWillBeUsedCallbackParam
} from 'workbox-core/types';
import { cleanupOutdatedCaches, precacheAndRoute } from 'workbox-precaching';
import { registerRoute } from 'workbox-routing';
import { CacheFirst, NetworkFirst, StaleWhileRevalidate } from 'workbox-strategies';

interface CacheTimestampEntry {
    request: Request;
    time: number;
}

class DynamicExpirationPlugin implements WorkboxPlugin {
    private static readonly SW_CACHE_TIME_HEADER: string = 'ezbookkeeping-sw-cache-time';
    private maxEntries: number;
    private maxAgeMilliseconds: number;
    private cleaningCache: boolean = false;

    constructor(maxEntries: number, maxAgeMilliseconds: number) {
        this.maxEntries = maxEntries;
        this.maxAgeMilliseconds = maxAgeMilliseconds;
    }

    public getMaxEntries(): number {
        return this.maxEntries;
    }

    public setMaxEntries(maxEntries: number): void {
        this.maxEntries = maxEntries;
    }

    public getMaxAgeMilliseconds(): number {
        return this.maxAgeMilliseconds;
    }

    public setMaxAgeMilliseconds(maxAgeMilliseconds: number): void {
        this.maxAgeMilliseconds = maxAgeMilliseconds;
    }

    public async cacheWillUpdate(param: CacheWillUpdateCallbackParam): Promise<Response | null> {
        const response = param.response;

        if (!response || response.status < 200 || response.status >= 300 || response.type === 'opaque') {
            return null;
        }

        const body = await response.blob();
        const headers = new Headers(response.headers);
        headers.set(DynamicExpirationPlugin.SW_CACHE_TIME_HEADER, Date.now().toString());

        return new Response(body, {
            status: response.status,
            statusText: response.statusText,
            headers: headers
        });
    }

    public async cachedResponseWillBeUsed(param: CachedResponseWillBeUsedCallbackParam): Promise<Response | null> {
        const cachedResponse = param.cachedResponse;

        if (!cachedResponse) {
            return null;
        }

        const cacheTime: string | null = cachedResponse.headers.get(DynamicExpirationPlugin.SW_CACHE_TIME_HEADER);

        if (!cacheTime) {
            return cachedResponse;
        }

        const age: number = Date.now() - Number(cacheTime);

        if (this.maxAgeMilliseconds > 0 && age >= this.maxAgeMilliseconds) {
            if (param.cacheName) {
                const cache = await caches.open(param.cacheName);
                await cache.delete(param.request);
            }

            return null;
        }

        return cachedResponse;
    }

    public async cacheDidUpdate(param: CacheDidUpdateCallbackParam): Promise<void> {
        if (this.cleaningCache || !param.cacheName) {
            return;
        }

        this.cleaningCache = true;

        const cache: Cache = await caches.open(param.cacheName);
        const requests: readonly Request[] = await cache.keys();

        if (requests.length <= this.maxEntries) {
            this.cleaningCache = false;
            return;
        }

        const entries: CacheTimestampEntry[] = [];

        for (const request of requests) {
            const response: Response | undefined = await cache.match(request);

            if (!response) {
                continue;
            }

            const cacheTime: string | null = response.headers.get(DynamicExpirationPlugin.SW_CACHE_TIME_HEADER);
            let time: number = cacheTime ? Number(cacheTime) : 0;

            if (Number.isFinite(time)) {
                const age: number = Date.now() - time;

                if (this.maxAgeMilliseconds > 0 && age >= this.maxAgeMilliseconds) {
                    await cache.delete(request);
                    continue;
                }
            } else {
                time = 0;
            }

            entries.push({
                request: request,
                time: time
            });
        }

        if (entries.length <= this.maxEntries) {
            this.cleaningCache = false;
            return;
        }

        entries.sort((a, b) => a.time - b.time);

        const removeCount: number = entries.length - this.maxEntries;

        for (let i = 0; i < removeCount; i++) {
            const entry = entries[i];

            if (entry && entry.request) {
                await cache.delete(entry.request);
            }
        }

        this.cleaningCache = false;
    }
}

class MapDataRequestStripTokenPlugin implements WorkboxPlugin {
    public async cacheKeyWillBeUsed(param: CacheKeyWillBeUsedCallbackParam): Promise<Request> {
        const url = new URL(param.request.url);

        if (url.searchParams.has('token')) {
            url.searchParams.delete('token');
            return new Request(url.href, param.request);
        }

        return param.request;
    }
}

interface MapCacheConfig {
    enabled: boolean;
    patterns: RegExp[];
    mapDataRequestStripTokenPlugin: MapDataRequestStripTokenPlugin;
    expirationPlugin: DynamicExpirationPlugin;
}

declare const self: ServiceWorkerGlobalScope;

const SW_ASSETS_CACHE_NAME: string = 'ezbookkeeping-assets-cache';
const SW_CODE_CACHE_NAME: string = 'ezbookkeeping-code-cache';
const SW_MAP_CACHE_NAME: string = 'ezbookkeeping-map-cache';

const SW_MESSAGE_TYPE_UPDATE_MAP_CACHE_CONFIG: string = 'UPDATE_MAP_CACHE_CONFIG';
const SW_MESSAGE_TYPE_UPDATE_MAP_CACHE_CONFIG_RESPONSE: string = 'UPDATE_MAP_CACHE_CONFIG_RESPONSE';

const DEFAULT_MAP_CACHE_MAX_ENTRIES: number = 1000;
const DEFAULT_MAP_CACHE_MAX_AGE_MILLISECONDS: number = 30 * 24 * 60 * 60 * 1000;

const mapCacheConfig: MapCacheConfig = {
    enabled: false,
    patterns: [],
    mapDataRequestStripTokenPlugin: new MapDataRequestStripTokenPlugin(),
    expirationPlugin: new DynamicExpirationPlugin(DEFAULT_MAP_CACHE_MAX_ENTRIES, DEFAULT_MAP_CACHE_MAX_AGE_MILLISECONDS)
};

self.skipWaiting();
clientsClaim();
precacheAndRoute(self.__WB_MANIFEST);
cleanupOutdatedCaches();

registerRoute(
    /.*\/img\/desktop\/.*\.(png|jpg|jpeg|gif|tiff|bmp|svg)/,
    new StaleWhileRevalidate({
        cacheName: SW_ASSETS_CACHE_NAME,
    })
);

registerRoute(
    /.*\/fonts\/.*\.(eot|ttf|svg|woff)/,
    new CacheFirst({
        cacheName: SW_ASSETS_CACHE_NAME,
    })
);

registerRoute(
    /.*\/(mobile|mobile\/|desktop|desktop\/)$/,
    new NetworkFirst({
        cacheName: SW_CODE_CACHE_NAME,
    })
);

registerRoute(
    /.*\/(mobile|mobile\/)#!\//,
    new NetworkFirst({
        cacheName: SW_CODE_CACHE_NAME,
    })
);

registerRoute(
    /.*\/(desktop|desktop\/)#\//,
    new NetworkFirst({
        cacheName: SW_CODE_CACHE_NAME,
    })
);

registerRoute(
    /.*\/(index\.html|mobile\.html|desktop\.html)/,
    new NetworkFirst({
        cacheName: SW_CODE_CACHE_NAME,
    })
);

registerRoute(
    /.*\/css\/.*\.css/,
    new CacheFirst({
        cacheName: SW_CODE_CACHE_NAME,
    })
);

registerRoute(
    /.*\/js\/.*\.js/,
    new CacheFirst({
        cacheName: SW_CODE_CACHE_NAME,
    })
);

registerRoute(
    ({ url }) => {
        if (!mapCacheConfig.enabled || mapCacheConfig.patterns.length < 1) {
            return false;
        }

        for (const pattern of mapCacheConfig.patterns) {
            if (pattern.test && pattern.test(url.href)) {
                return true;
            }
        }

        return false;
    },
    new CacheFirst({
        cacheName: SW_MAP_CACHE_NAME,
        plugins: [
            mapCacheConfig.mapDataRequestStripTokenPlugin,
            mapCacheConfig.expirationPlugin
        ]
    })
);

self.addEventListener('message', (event: ExtendableMessageEvent) => {
    try {
        if (event.data && event.data.type === SW_MESSAGE_TYPE_UPDATE_MAP_CACHE_CONFIG && 'payload' in event.data) {
            mapCacheConfig.enabled = !!event.data.payload['enabled'];
            mapCacheConfig.patterns = [];
            mapCacheConfig.expirationPlugin.setMaxEntries(event.data.payload['maxEntries'] ?? DEFAULT_MAP_CACHE_MAX_ENTRIES);
            mapCacheConfig.expirationPlugin.setMaxAgeMilliseconds(event.data.payload['maxAgeMilliseconds'] ?? DEFAULT_MAP_CACHE_MAX_AGE_MILLISECONDS);

            if (event.data.payload['patterns'] && Array.isArray(event.data.payload['patterns'])) {
                for (const pattern of event.data.payload['patterns']) {
                    if (pattern) {
                        mapCacheConfig.patterns.push(new RegExp(pattern as string));
                    }
                }
            }

            if (event.ports && event.ports[0] && typeof event.ports[0].postMessage === 'function') {
                event.ports[0].postMessage({
                    type: SW_MESSAGE_TYPE_UPDATE_MAP_CACHE_CONFIG_RESPONSE,
                    payload: {
                        enabled: mapCacheConfig.enabled,
                        patterns: event.data.payload['patterns'],
                        maxEntries: mapCacheConfig.expirationPlugin.getMaxEntries(),
                        maxAgeMilliseconds: mapCacheConfig.expirationPlugin.getMaxAgeMilliseconds()
                    }
                });
            }
        }
    } catch (ex) {
        console.error('failed to process message in service worker', ex);
    }
});
