import icons from '../consts/icon.js';
import utils from '../lib/utils.js';

export default function (iconId) {
    if (utils.isNumber(iconId)) {
        iconId = iconId.toString();
    }

    if (!icons.allCategoryIcons[iconId]) {
        return icons.defaultCategoryIcon.icon;
    }

    return icons.allCategoryIcons[iconId].icon;
}
