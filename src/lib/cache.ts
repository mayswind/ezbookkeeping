import type {
    BrowserCacheStatistics
} from '@/core/cache.ts';

import {
    SW_PRECACHE_CACHE_NAME_PREFIX,
    SW_RUNTIME_CACHE_NAME_PREFIX,
    SW_ASSETS_CACHE_NAME,
    SW_CODE_CACHE_NAME,
    SW_MAP_CACHE_NAME
} from '@/consts/cache.ts';

import { isFunction, isObject, isNumber } from './common.ts';
import logger from './logger.ts';

function findFirstCacheName(prefix: string): Promise<string> {
    if (!window.caches) {
        logger.error('caches API is not supported in this browser');
        return Promise.reject(new Error('caches API is not supported'));
    }

    return window.caches.keys().then(cacheNames => {
        for (const cacheName of cacheNames) {
            if (cacheName.startsWith(prefix)) {
                return cacheName;
            }
        }

        throw new Error(`cache with prefix "${prefix}" not found`);
    });
}

async function getCacheTotalSize(cacheName: string): Promise<number> {
    if (!window.caches) {
        logger.error('caches API is not supported in this browser');
        return Promise.reject(new Error('caches API is not supported'));
    }

    const cache = await window.caches.open(cacheName);
    const requests = await cache.keys();
    let totalSize = 0;

    for (const request of requests) {
        try {
            const response = await cache.match(request);

            if (response) {
                const blob = await response.clone().blob();
                totalSize += blob.size;
            }
        } catch (ex) {
            logger.warn(`failed to get size for request ${request.url} in cache ${cacheName}`, ex);
        }
    }

    return totalSize;
}

export function loadBrowserCacheStatistics(): Promise<BrowserCacheStatistics> {
    return new Promise((resolve, reject) => {
        const caches = window.caches;

        if (!caches) {
            logger.error('caches API is not supported in this browser');
            reject(new Error('caches API is not supported in this browser'));
            return;
        }

        return Promise.all([
            navigator && navigator.storage && isFunction(navigator.storage.estimate) ? navigator.storage.estimate() : Promise.resolve(undefined),
            findFirstCacheName(SW_PRECACHE_CACHE_NAME_PREFIX).then(cacheName => getCacheTotalSize(cacheName)).catch(() => 0),
            findFirstCacheName(SW_RUNTIME_CACHE_NAME_PREFIX).then(cacheName => getCacheTotalSize(cacheName)).catch(() => 0),
            getCacheTotalSize(SW_CODE_CACHE_NAME),
            getCacheTotalSize(SW_ASSETS_CACHE_NAME),
            getCacheTotalSize(SW_MAP_CACHE_NAME)
        ]).then(([storageEstimate, precacheCacheSize, runtimeCacheSize, codeCacheSize, assetsCacheSize, mapCacheSize]) => {
            let totalCacheSize: number = 0;

            if (storageEstimate) {
                const cachesUsage = 'usageDetails' in storageEstimate
                && isObject(storageEstimate.usageDetails)
                && 'caches' in storageEstimate.usageDetails
                    ? storageEstimate.usageDetails.caches : undefined;

                if (isNumber(cachesUsage)) {
                    totalCacheSize = cachesUsage;
                } else if (isNumber(storageEstimate.usage)) {
                    totalCacheSize = storageEstimate.usage;
                }
            }

            if (totalCacheSize < 1) {
                totalCacheSize = precacheCacheSize + runtimeCacheSize + codeCacheSize + assetsCacheSize + mapCacheSize;
            }

            let othersCacheSize: number = totalCacheSize - precacheCacheSize - runtimeCacheSize - codeCacheSize - assetsCacheSize - mapCacheSize;

            if (othersCacheSize < 0) {
                othersCacheSize = 0;
            }

            resolve({
                totalCacheSize: totalCacheSize,
                codeCacheSize: codeCacheSize + runtimeCacheSize,
                assetsCacheSize: assetsCacheSize + precacheCacheSize,
                mapCacheSize: mapCacheSize,
                othersCacheSize: othersCacheSize
            });
        }).catch(error => {
            logger.error("failed to clear cache", error);
            reject(error);
        });
    });
}

export function clearAllBrowserCaches(): Promise<void> {
    if (!window.caches) {
        logger.error('caches API is not supported in this browser');
        return Promise.reject();
    }

    return new Promise((resolve, reject) => {
        window.caches.keys().then(cacheNames => {
            const promises = [];

            for (const cacheName of cacheNames) {
                promises.push(window.caches.delete(cacheName).then(success => {
                    if (success) {
                        logger.info(`cache "${cacheName}" cleared successfully`);
                        return Promise.resolve(cacheName);
                    } else {
                        logger.warn(`failed to clear cache "${cacheName}"`);
                        return Promise.reject(cacheName);
                    }
                }));
            }

            Promise.all(promises).then(() => {
                logger.info("all caches cleared successfully");
                resolve();
            }).catch(() => {
                resolve();
            });
        }).catch(error => {
            logger.warn("failed to clear cache", error);
            reject(error);
        });
    });
}
