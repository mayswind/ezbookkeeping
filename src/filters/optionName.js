import utils from '../lib/utils.js';

export default function (value, options, keyField, nameField, defaultName) {
    if (utils.isArray(options)) {
        if (keyField) {
            for (let i = 0; i < options.length; i++) {
                const option = options[i];

                if (option[keyField] === value) {
                    return option[nameField];
                }
            }
        } else {
            if (options[value]) {
                const option = options[value];

                return option[nameField];
            }
        }
    } else if (utils.isObject(options)) {
        if (keyField) {
            for (let key in options) {
                if (!Object.prototype.hasOwnProperty.call(options, key)) {
                    continue;
                }

                const option = options[key];

                if (option[keyField] === value) {
                    return option[nameField];
                }
            }
        } else {
            if (options[value]) {
                const option = options[value];

                return option[nameField];
            }
        }
    }

    return defaultName;
}
