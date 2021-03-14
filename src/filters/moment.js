import utils from '../lib/utils.js';

export default function (value, format, options) {
    if (!utils.isNumber(value)) {
        value = utils.getUnixTime(value);
    }

    let utcOffset = null;
    let currentUtcOffset = null;

    if (utils.isObject(options) && utils.isNumber(options.utcOffset)) {
        utcOffset = options.utcOffset;
    }

    if (utils.isObject(options) && utils.isNumber(options.currentUtcOffset)) {
        currentUtcOffset = options.currentUtcOffset;
    }

    return utils.formatUnixTime(value, format, utcOffset, currentUtcOffset);
}
