import icons from "../consts/icon.js";
import utils from "../lib/utils.js";

export default function (iconId) {
    if (utils.isNumber(iconId)) {
        iconId = iconId.toString();
    }

    if (!icons.allAccountIcons[iconId]) {
        return icons.defaultAccountIcon.icon;
    }

    return icons.allAccountIcons[iconId].icon;
}
