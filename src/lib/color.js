export function getColorsInRows(allColorInfos, itemPerRow) {
    const ret = [];
    let rowCount = -1;

    for (let i = 0; i < allColorInfos.length; i++) {
        if (i % itemPerRow === 0) {
            ret[++rowCount] = [];
        }

        ret[rowCount].push({
            color: allColorInfos[i]
        });
    }

    return ret;
}
