import icons from "../consts/icon.js";
import utils from "../lib/utils.js";

export default function (token) {
    const ua = utils.parseUserAgent(token.userAgent);

    if (!ua || !ua.device) {
        return icons.deviceIcons.desktop.f7Icon;
    }

    if (ua.device.type === 'mobile') {
        return icons.deviceIcons.mobile.f7Icon;
    } else if (ua.device.type === 'wearable') {
        return icons.deviceIcons.wearable.f7Icon;
    } else if (ua.device.type === 'tablet') {
        return icons.deviceIcons.tablet.f7Icon;
    } else if (ua.device.type === 'smarttv') {
        return icons.deviceIcons.tv.f7Icon;
    } else {
        return icons.deviceIcons.desktop.f7Icon;
    }
}
