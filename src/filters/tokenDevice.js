import utils from '../lib/utils.js';

export default function (token) {
    const ua = utils.parseUserAgent(token.userAgent);
    let result = '';

    if (ua.device.model) {
        result = ua.device.model;
    } else if (ua.os.name) {
        result = ua.os.name;

        if (ua.os.version) {
            result += ' ' + ua.os.version;
        }
    }

    if (ua.browser.name) {
        let browserInfo = ua.browser.name;

        if (ua.browser.version) {
            browserInfo += ' ' + ua.browser.version;
        }

        if (result) {
            result += ' (' + browserInfo + ')';
        } else {
            result = browserInfo;
        }
    }

    if (!result) {
        return 'Unknown Device';
    }

    return result;
}
