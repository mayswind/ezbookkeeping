import type { IconInfo, IconInfoWithId } from '@/core/icon.ts';

export function getIconsInRows(allIconInfos: Record<string, IconInfo>, itemPerRow: number): IconInfoWithId[][] {
    const ret: IconInfoWithId[][] = [];
    let rowCount = 0;

    for (const iconInfoId in allIconInfos) {
        if (!Object.prototype.hasOwnProperty.call(allIconInfos, iconInfoId)) {
            continue;
        }

        const iconInfo = allIconInfos[iconInfoId];

        if (!ret[rowCount]) {
            ret[rowCount] = [];
        } else if (ret[rowCount] && ret[rowCount].length >= itemPerRow) {
            rowCount++;
            ret[rowCount] = [];
        }

        ret[rowCount].push({
            id: iconInfoId,
            icon: iconInfo.icon
        });
    }

    return ret;
}
