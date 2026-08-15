import { ref } from 'vue';
import { defineStore } from 'pinia';

import { itemAndIndex } from '@/core/base.ts';
import type { UserCustomIconInfoResponse, UserCustomIconNewDisplayOrderRequest } from '@/models/user_custom_icon.ts';

import { isEquals } from '@/lib/common.ts';
import { deleteCustomIconCacheEntry } from '@/lib/cache.ts';

import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

export const useUserCustomIconsStore = defineStore('userCustomIcons', () => {
    const allCustomIcons = ref<UserCustomIconInfoResponse[]>([]);
    const customIconListStateInvalid = ref<boolean>(true);

    function loadCustomIconsList(customIcons: UserCustomIconInfoResponse[]): void {
        allCustomIcons.value = customIcons;
    }

    function addCustomIconToIconList(customIcon: UserCustomIconInfoResponse): void {
        allCustomIcons.value.push(customIcon);
    }

    function updateCustomIconDisplayOrderInIconList({ from, to }: { from: number, to: number }): void {
        allCustomIcons.value.splice(to, 0, allCustomIcons.value.splice(from, 1)[0] as UserCustomIconInfoResponse);
    }

    function removeCustomIconFromIconList(currentCustomIcon: UserCustomIconInfoResponse): void {
        for (const [customIcon, index] of itemAndIndex(allCustomIcons.value)) {
            if (customIcon.id === currentCustomIcon.id) {
                allCustomIcons.value.splice(index, 1);
                break;
            }
        }
    }

    function updateCustomIconListInvalidState(invalidState: boolean): void {
        customIconListStateInvalid.value = invalidState;
    }

    function resetCustomIcons(): void {
        allCustomIcons.value = [];
        customIconListStateInvalid.value = true;
    }

    function loadAllCustomIcons({ force }: { force?: boolean }): Promise<UserCustomIconInfoResponse[]> {
        if (!force && !customIconListStateInvalid.value) {
            return Promise.resolve(allCustomIcons.value);
        }

        return new Promise((resolve, reject) => {
            services.getAllUserCustomIcons().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve custom icon list' });
                    return;
                }

                if (customIconListStateInvalid.value) {
                    updateCustomIconListInvalidState(false);
                }

                if (force && data.result && isEquals(allCustomIcons.value, data.result)) {
                    reject({ message: 'Custom icon list is up to date', isUpToDate: true });
                    return;
                }

                loadCustomIconsList(data.result);

                resolve(allCustomIcons.value);
            }).catch(error => {
                if (force) {
                    logger.error('failed to force load custom icon list', error);
                } else {
                    logger.error('failed to load custom icon list', error);
                }

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve custom icon list' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function uploadCustomIcon({ iconFile, clientSessionId }: { iconFile: File, clientSessionId: string }): Promise<UserCustomIconInfoResponse> {
        return new Promise((resolve, reject) => {
            services.uploadUserCustomIcon({ iconFile, clientSessionId }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to upload custom icon' });
                    return;
                }

                addCustomIconToIconList(data.result);

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to upload custom icon', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to upload custom icon' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function changeCustomIconDisplayOrder({ iconId, from, to }: { iconId: string, from: number, to: number }): Promise<void> {
        return new Promise((resolve, reject) => {
            let currentIcon: UserCustomIconInfoResponse | null = null;

            for (const icon of allCustomIcons.value) {
                if (icon.id === iconId) {
                    currentIcon = icon;
                    break;
                }
            }

            if (!currentIcon || !allCustomIcons.value[to]) {
                reject({ message: 'Unable to move custom icon' });
                return;
            }

            if (!customIconListStateInvalid.value) {
                updateCustomIconListInvalidState(true);
            }

            updateCustomIconDisplayOrderInIconList({ from, to });

            resolve();
        });
    }

    function updateCustomIconDisplayOrders(): Promise<boolean> {
        const newDisplayOrders: UserCustomIconNewDisplayOrderRequest[] = [];

        for (const [icon, index] of itemAndIndex(allCustomIcons.value)) {
            newDisplayOrders.push({
                id: icon.id,
                displayOrder: index + 1
            });
        }

        return new Promise((resolve, reject) => {
            services.moveUserCustomIcons({
                newDisplayOrders: newDisplayOrders
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to move custom icon' });
                    return;
                }

                loadAllCustomIcons({ force: false }).finally(() => {
                    if (customIconListStateInvalid.value) {
                        updateCustomIconListInvalidState(false);
                    }

                    resolve(data.result);
                });
            }).catch(error => {
                logger.error('failed to save custom icons display order', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to move custom icon' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function deleteCustomIcon({ customIcon }: { customIcon: UserCustomIconInfoResponse }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.deleteUserCustomIcon({ id: customIcon.id }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to delete custom icon' });
                    return;
                }

                removeCustomIconFromIconList(customIcon);
                deleteCustomIconCacheEntry(services.getUserCustomIconUrlWithToken(customIcon.id));

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to delete custom icon', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to delete custom icon' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function getUserCustomIconUrlWithToken(customIcon: UserCustomIconInfoResponse): string {
        return services.getUserCustomIconUrlWithToken(customIcon.id);
    }

    return {
        // states
        allCustomIcons,
        customIconListStateInvalid,
        // functions
        updateCustomIconListInvalidState,
        resetCustomIcons,
        loadAllCustomIcons,
        uploadCustomIcon,
        changeCustomIconDisplayOrder,
        updateCustomIconDisplayOrders,
        deleteCustomIcon,
        getUserCustomIconUrlWithToken
    };
});
