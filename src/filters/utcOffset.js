import utils from '../lib/utils.js';

export default function (utcOffsetMinutes) {
    const utcOffset = utils.getUtcOffsetByUtcOffsetMinutes(utcOffsetMinutes);
    return `(UTC${utcOffset})`;
}
