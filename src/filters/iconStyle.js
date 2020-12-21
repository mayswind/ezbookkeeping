import defaultIconStyle from "./defaultIconStyle.js";
import accountIconStyle from "./accountIconStyle.js";
import categoryIconStyle from "./categoryIconStyle.js";

export default function (color, iconType, defaultColor) {
    if (iconType === 'account') {
        return accountIconStyle(color, defaultColor);
    } else if (iconType === 'category') {
        return categoryIconStyle(color, defaultColor);
    } else {
        return defaultIconStyle(color, defaultColor);
    }
}
