import icons from "../consts/icon.js";
import utils from "../lib/utils.js";

export default function (iconId) {
    if (utils.isNumber(iconId)) {
        iconId = iconId.toString();
    }

    if (iconId <= icons.totalAccountIconCount) {
        return icons.allAccountIcons[iconId].f7Icon;
    }

    return icons.defaultAccountIcon.f7Icon;
}
