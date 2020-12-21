import colorConstants from '../consts/color.js';

export default function (color, defaultColor) {
    if (color && color !== colorConstants.defaultCategoryColor) {
        color = '#' + color;
    } else {
        color = defaultColor;
    }

    return {
        color: color
    };
}
