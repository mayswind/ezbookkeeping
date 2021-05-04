import utils from '../lib/utils.js';

export default function (value, options, keyName, valueName) {
    if (utils.isArray(options)) {
        for (let i = 0; i < options.length; i++) {
            const option = options[i];

            if (option[keyName] === value) {
                return option[valueName];
            }
        }
    } else if (utils.isObject(options)) {
        for (let key in options) {
            if (!Object.prototype.hasOwnProperty.call(options, key)) {
                continue;
            }

            const option = options[key];

            if (option[keyName] === value) {
                return option[valueName];
            }
        }
    }

    return '';
}
