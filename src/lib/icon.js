export function getIconsInRows(allIconInfos, itemPerRow) {
    const ret = [];
    let rowCount = 0;

    for (let iconInfoId in allIconInfos) {
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
