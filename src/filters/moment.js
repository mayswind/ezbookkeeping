import utils from '../lib/utils.js';

export default function (value, format) {
    if (utils.isNumber(value)) {
        return utils.formatUnixTime(value, format);
    } else {
        return utils.formatTime(value, format);
    }
}
