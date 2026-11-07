import {
    type BrowserCacheStatistics,
    type SWMapCacheConfig,
    type SWTransactionPictureCacheConfig,
    FileCacheType
} from '@/core/cache.ts';

import {
    SW_SHARE_CACHE_NAME,
    SW_MESSAGE_TYPE_UPDATE_MAP_CACHE_CONFIG,
    SW_MESSAGE_TYPE_UPDATE_MAP_CACHE_CONFIG_RESPONSE,
    SW_MESSAGE_TYPE_UPDATE_TRANSACTION_PICTURE_CACHE_CONFIG,
    SW_MESSAGE_TYPE_UPDATE_TRANSACTION_PICTURE_CACHE_CONFIG_RESPONSE,
    TRANSACTION_PICTURE_CACHE_MAX_ENTRIES,
    MAP_CACHE_MAX_ENTRIES
} from '@/consts/cache.ts';

import { isFunction, isObject, isNumber } from './common.ts';
import services from './services.ts';
import logger from './logger.ts';

let controllerchangeListenerAdded: boolean = false;
let transactionPictureControllerchangeListenerAdded: boolean = false;

function findFirstCacheName(prefix: string): Promise<string> {
    if (!prefix) {
        return Promise.reject(new Error('prefix is empty'));
    }

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

function doUpdateMapCacheExpiration(expireSeconds: number): Promise<void> {
    const config: SWMapCacheConfig = {
        enabled: expireSeconds >= 0,
        patterns: services.getMapProxyTileImageAndAnnotationImageUrlPatterns(),
        maxEntries: MAP_CACHE_MAX_ENTRIES,
        maxAgeMilliseconds: expireSeconds * 1000
    };

    return new Promise((resolve, reject) => {
        if (!navigator.serviceWorker || !navigator.serviceWorker.controller) {
            reject(new Error('Service worker is not supported or not active'));
            return;
        }

        const controller = navigator.serviceWorker.controller;

        navigator.serviceWorker.ready.then(() => {
            const messageChannel = new MessageChannel();

            messageChannel.port1.onmessage = (event) => {
                if (event.data && event.data.type === SW_MESSAGE_TYPE_UPDATE_MAP_CACHE_CONFIG_RESPONSE) {
                    logger.info('Map cache config updated successfully in service worker: ' + JSON.stringify(event.data.payload));
                    resolve();
                } else {
                    logger.error('cannot update map cache config, invalid response from service worker', event);
                    reject(new Error('Invalid response from service worker'));
                }
            };

            controller.postMessage({
                type: SW_MESSAGE_TYPE_UPDATE_MAP_CACHE_CONFIG,
                payload: config
            }, [messageChannel.port2]);
        }).catch(error => {
            logger.error('failed to update map cache config', error);
            reject(error);
        });
    });
}

function doUpdateTransactionPictureCacheExpiration(expireSeconds: number): Promise<void> {
    const config: SWTransactionPictureCacheConfig = {
        enabled: expireSeconds >= 0,
        maxEntries: TRANSACTION_PICTURE_CACHE_MAX_ENTRIES,
        maxAgeMilliseconds: expireSeconds * 1000
    };

    return new Promise((resolve, reject) => {
        if (!navigator.serviceWorker || !navigator.serviceWorker.controller) {
            reject(new Error('Service worker is not supported or not active'));
            return;
        }

        const controller = navigator.serviceWorker.controller;

        navigator.serviceWorker.ready.then(() => {
            const messageChannel = new MessageChannel();

            messageChannel.port1.onmessage = (event) => {
                if (event.data && event.data.type === SW_MESSAGE_TYPE_UPDATE_TRANSACTION_PICTURE_CACHE_CONFIG_RESPONSE) {
                    logger.info('Transaction picture cache config updated successfully in service worker: ' + JSON.stringify(event.data.payload));
                    resolve();
                } else {
                    logger.error('cannot update transaction picture cache config, invalid response from service worker', event);
                    reject(new Error('Invalid response from service worker'));
                }
            };

            controller.postMessage({
                type: SW_MESSAGE_TYPE_UPDATE_TRANSACTION_PICTURE_CACHE_CONFIG,
                payload: config
            }, [messageChannel.port2]);
        }).catch(error => {
            logger.error('failed to update transaction picture cache config', error);
            reject(error);
        });
    });
}

function deleteAuthenticatedCacheEntry(cacheType: FileCacheType, originalUrl: string): Promise<boolean> {
    if (!window.caches || !originalUrl) {
        return Promise.resolve(false);
    }

    const url = new URL(originalUrl, window.location.href);
    url.searchParams.delete('token');

    return window.caches.open(cacheType.key).then(cache => cache.delete(url.href)).catch(error => {
        logger.warn(`failed to delete cache entry "${originalUrl}"`, error);
        return false;
    });
}

function deleteAuthenticatedCacheEntries(cacheType: FileCacheType, originalUrls: string[]): Promise<boolean> {
    if (!window.caches || !originalUrls || originalUrls.length < 1) {
        return Promise.resolve(false);
    }

    const urls: string[] = [];

    for (const originalUrl of originalUrls) {
        if (originalUrl) {
            const url = new URL(originalUrl, window.location.href);
            url.searchParams.delete('token');
            urls.push(url.href);
        }
    }

    if (urls.length < 1) {
        return Promise.resolve(false);
    }

    return window.caches.open(cacheType.key).then(cache => {
        const deletePromises: Promise<boolean>[] = [];

        for (const url of urls) {
            deletePromises.push(cache.delete(url).catch(error => {
                logger.warn(`failed to delete cache entry "${url}"`, error);
                return false;
            }));
        }

        return Promise.all(deletePromises).then(results => results.every(result => result));
    }).catch(error => {
        logger.warn(`failed to delete cache entries in cache "${cacheType.key}"`, error);
        return false;
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

export function getShareCacheImageBlob(): Promise<Blob | undefined> {
    if (!window.caches) {
        logger.error('caches API is not supported in this browser');
        return Promise.resolve(undefined);
    }

    return new Promise((resolve) => {
        window.caches.open(SW_SHARE_CACHE_NAME).then(cache => {
            cache.match(SW_SHARE_CACHE_NAME).then(response => {
                if (!response) {
                    resolve(undefined);
                    return;
                }

                response.blob().then(blob => {
                    cache.delete(SW_SHARE_CACHE_NAME).then(() => {
                        resolve(blob);
                    }).catch(error => {
                        logger.warn('failed to delete share cache image blob', error);
                        resolve(blob);
                    });
                }).catch(error => {
                    logger.error('failed to read share cache image blob', error);
                    resolve(undefined);
                });
            }).catch(error => {
                logger.error('failed to match share cache image blob', error);
                resolve(undefined);
            });
        }).catch(error => {
            logger.error('failed to open share cache', error);
            resolve(undefined);
        });
    });
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
            findFirstCacheName(FileCacheType.ResourceFiles.cacheNamePrefixes[0] ?? '').then(cacheName => getCacheTotalSize(cacheName)).catch(() => 0),
            findFirstCacheName(FileCacheType.ApplicationCode.cacheNamePrefixes[0] ?? '').then(cacheName => getCacheTotalSize(cacheName)).catch(() => 0),
            getCacheTotalSize(FileCacheType.ApplicationCode.key),
            getCacheTotalSize(FileCacheType.ResourceFiles.key),
            getCacheTotalSize(FileCacheType.MapData.key),
            getCacheTotalSize(FileCacheType.UserAvatars.key),
            getCacheTotalSize(FileCacheType.TransactionPictures.key),
            getCacheTotalSize(FileCacheType.CustomIcons.key)
        ]).then(([storageEstimate, precacheCacheSize, runtimeCacheSize, codeCacheSize, assetsCacheSize, mapCacheSize, userAvatarCacheSize, transactionPictureCacheSize, customIconCacheSize]) => {
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
                totalCacheSize = precacheCacheSize + runtimeCacheSize + codeCacheSize + assetsCacheSize + mapCacheSize + userAvatarCacheSize + transactionPictureCacheSize + customIconCacheSize;
            }

            let othersCacheSize: number = totalCacheSize - precacheCacheSize - runtimeCacheSize - codeCacheSize - assetsCacheSize - mapCacheSize - userAvatarCacheSize - transactionPictureCacheSize - customIconCacheSize;

            if (othersCacheSize < 0) {
                othersCacheSize = 0;
            }

            resolve({
                totalCacheSize: totalCacheSize,
                codeCacheSize: codeCacheSize + runtimeCacheSize,
                assetsCacheSize: assetsCacheSize + precacheCacheSize,
                mapCacheSize: mapCacheSize,
                userAvatarCacheSize: userAvatarCacheSize,
                transactionPictureCacheSize: transactionPictureCacheSize,
                customIconCacheSize: customIconCacheSize,
                othersCacheSize: othersCacheSize
            });
        }).catch(error => {
            logger.error("failed to clear cache", error);
            reject(error);
        });
    });
}

export function updateMapCacheExpiration(expireSeconds: number): void {
    if ('serviceWorker' in navigator) {
        if (!controllerchangeListenerAdded) {
            navigator.serviceWorker.addEventListener('controllerchange', () => {
                doUpdateMapCacheExpiration(expireSeconds);
            });
            controllerchangeListenerAdded = true;
        }

        if (navigator.serviceWorker.controller) {
            doUpdateMapCacheExpiration(expireSeconds);
        }
    }
}

export function updateTransactionPictureCacheExpiration(expireSeconds: number): void {
    if ('serviceWorker' in navigator) {
        if (!transactionPictureControllerchangeListenerAdded) {
            navigator.serviceWorker.addEventListener('controllerchange', () => {
                doUpdateTransactionPictureCacheExpiration(expireSeconds);
            });
            transactionPictureControllerchangeListenerAdded = true;
        }

        if (navigator.serviceWorker.controller) {
            doUpdateTransactionPictureCacheExpiration(expireSeconds);
        }
    }
}

export function deleteUserAvatarCacheEntry(avatarUrl: string): Promise<boolean> {
    return deleteAuthenticatedCacheEntry(FileCacheType.UserAvatars, avatarUrl);
}

export function deleteTransactionPictureCacheEntries(...pictureUrls: string[]): Promise<boolean> {
    return deleteAuthenticatedCacheEntries(FileCacheType.TransactionPictures, pictureUrls);
}

export function deleteCustomIconCacheEntry(iconUrl: string): Promise<boolean> {
    return deleteAuthenticatedCacheEntry(FileCacheType.CustomIcons, iconUrl);
}

export function clearCaches(cacheNames: string[], cacheNamePrefixes?: string[]): Promise<void> {
    if (!window.caches) {
        logger.error('caches API is not supported in this browser');
        return Promise.reject();
    }

    return new Promise((resolve) => {
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

        if (cacheNamePrefixes) {
            for (const prefix of cacheNamePrefixes) {
                promises.push(findFirstCacheName(prefix).then(cacheName => {
                    return window.caches.delete(cacheName).then(success => {
                        if (success) {
                            logger.info(`cache "${cacheName}" cleared successfully`);
                            return Promise.resolve(cacheName);
                        } else {
                            logger.warn(`failed to clear cache "${cacheName}"`);
                            return Promise.reject(cacheName);
                        }
                    });
                }).catch(error => {
                    logger.warn(`cache with prefix "${prefix}" not found`, error);
                    return Promise.resolve();
                }));
            }
        }

        Promise.all(promises).then(() => {
            resolve();
        }).catch(() => {
            resolve();
        });
    });
}

export function clearFileCaches(...cacheTypes: FileCacheType[]): Promise<void> {
    if (!cacheTypes || !cacheTypes.length) {
        return Promise.resolve();
    }

    const cacheNames: string[] = [];
    const cacheNamePrefixes: string[] = [];

    for (const cacheType of cacheTypes) {
        cacheNames.push(cacheType.key);
        cacheNamePrefixes.push(...cacheType.cacheNamePrefixes);
    }

    return clearCaches(cacheNames, cacheNamePrefixes);
}

export function clearTransactionPictureCache(): Promise<void> {
    return clearFileCaches(FileCacheType.TransactionPictures);
}

export function clearAllUserCaches(): Promise<void> {
    return clearFileCaches(
        FileCacheType.UserAvatars,
        FileCacheType.TransactionPictures,
        FileCacheType.CustomIcons
    );
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
            logger.error("failed to clear cache", error);
            reject(error);
        });
    });
}
