import colorConstants from '../consts/color.js';

export default function (color, defaultColor) {
    if (color && color !== colorConstants.defaultAccountColor) {
        color = '#' + color;
    } else {
        color = defaultColor;
    }

    return {
        color: color
    };
}
