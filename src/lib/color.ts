import type { ColorValue, ColorInfo } from '@/core/color.ts';

export function getColorsInRows(allColorValues: ColorValue[], itemPerRow: number): ColorInfo[][] {
    const ret: ColorInfo[][] = [];
    let rowCount = -1;

    for (let i = 0; i < allColorValues.length; i++) {
        if (i % itemPerRow === 0) {
            ret[++rowCount] = [];
        }

        ret[rowCount].push({
            color: allColorValues[i]
        });
    }

    return ret;
}
