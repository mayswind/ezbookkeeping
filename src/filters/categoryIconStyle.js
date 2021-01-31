import colorConstants from '../consts/color.js';

export default function (color, defaultColor, additionalFieldName) {
    if (color && color !== colorConstants.defaultCategoryColor) {
        color = '#' + color;
    } else {
        color = defaultColor;
    }

    const ret = {
        color: color
    };

    if (additionalFieldName) {
        ret[additionalFieldName] = color;
    }

    return ret;
}
