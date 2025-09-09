import { entries } from '@/core/base.ts';
import type { IconInfo, IconInfoWithId } from '@/core/icon.ts';

export function getIconsInRows(allIconInfos: Record<string, IconInfo>, itemPerRow: number): IconInfoWithId[][] {
    const ret: IconInfoWithId[][] = [];
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
