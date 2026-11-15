import { entries } from '@/core/base.ts';
import {
    type AccountIconItemType,
    type CategoryIconItemType,
    type CommonIconItemType,
    type SystemIconInfo,
    type SystemIconInfoWithId,
    type UserCustomIconInfo,
    IconType
} from '@/core/icon.ts';

import { type UserCustomIconInfoResponse } from '@/models/user_custom_icon.ts';

export function getAccountIconType(iconType: number | undefined): AccountIconItemType {
    if (iconType === IconType.UserCustom) {
        return 'user-custom';
    } else {
        return 'account';
    }
}

export function getCategoryIconType(iconType: number | undefined): CategoryIconItemType {
    if (iconType === IconType.UserCustom) {
        return 'user-custom';
    } else {
        return 'category';
    }
}

export function getIconType(type: string | undefined, iconType: unknown | undefined): CommonIconItemType {
    if (iconType === IconType.UserCustom) {
        return 'user-custom';
    } else if (type === 'account') {
        return 'account';
    } else if (type === 'category') {
        return 'category';
    } else {
        return 'fixed';
    }
}

export function getSystemIconsInRows(allIconInfos: Record<string, SystemIconInfo>, itemPerRow: number): SystemIconInfoWithId[][] {
    const ret: SystemIconInfoWithId[][] = [];
    let rowCount = 0;

    for (const [iconInfoId, iconInfo] of entries(allIconInfos)) {
        if (!ret[rowCount]) {
            ret[rowCount] = [];
        } else if (ret[rowCount] && ret[rowCount]!.length >= itemPerRow) {
            rowCount++;
            ret[rowCount] = [];
        }

        ret[rowCount]!.push({
            id: iconInfoId,
            icon: iconInfo.icon
        });
    }

    return ret;
}

export function getUserCustomIconsInRows(userCustomIconInfos: UserCustomIconInfoResponse[], itemPerRow: number): UserCustomIconInfo[][] {
    const ret: UserCustomIconInfo[][] = [];
    let rowCount = 0;

    for (const iconInfo of userCustomIconInfos) {
        if (!ret[rowCount]) {
            ret[rowCount] = [];
        } else if (ret[rowCount] && ret[rowCount]!.length >= itemPerRow) {
            rowCount++;
            ret[rowCount] = [];
        }

        ret[rowCount]!.push({
            id: iconInfo.id
        });
    }

    return ret;
}

export function isUserCustomIconEquals(customIcon1: UserCustomIconInfoResponse, customIcon2: UserCustomIconInfoResponse): boolean {
    if (customIcon1.id !== customIcon2.id || customIcon1.displayOrder !== customIcon2.displayOrder) {
        return false;
    }

    try {
        const originalUrl1 = new URL(customIcon1.originalUrl, window.location.href);
        const originalUrl2 = new URL(customIcon2.originalUrl, window.location.href);
        originalUrl1.searchParams.delete('token');
        originalUrl2.searchParams.delete('token');
        return originalUrl1.href === originalUrl2.href;
    } catch {
        return customIcon1.originalUrl === customIcon2.originalUrl;
    }
}

export function isUserCustomIconsEquals(customIcons1: UserCustomIconInfoResponse[], customIcons2: UserCustomIconInfoResponse[]): boolean {
    if (customIcons1.length !== customIcons2.length) {
        return false;
    }

    for (let i = 0; i < customIcons1.length; i++) {
        if (!isUserCustomIconEquals(customIcons1[i]!, customIcons2[i]!)) {
            return false;
        }
    }

    return true;
}
