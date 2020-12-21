import colorConstants from '../consts/color.js';

export default function (color, defaultColor) {
    if (color && color !== colorConstants.defaultColor) {
        color = '#' + color;
    } else {
        color = defaultColor;
    }

    return {
        color: color
    };
}
