import utils from '../lib/utils.js';

export default function (value, format, options) {
    if (!utils.isNumber(value)) {
        value = utils.getUnixTime(value);
    }

    if (utils.isObject(options) && utils.isNumber(options.utcOffset)) {
        if (!utils.isNumber(options.currentUtcOffset)) {
            options.currentUtcOffset = utils.getTimezoneOffsetMinutes();
        }

        value = utils.getDummyUnixTimeForLocalDisplay(value, options.utcOffset, options.currentUtcOffset);
    }

    return utils.formatUnixTime(value, format);
}
