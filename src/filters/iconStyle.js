import defaultIconStyle from "./defaultIconStyle.js";
import accountIconStyle from "./accountIconStyle.js";
import categoryIconStyle from "./categoryIconStyle.js";

export default function (color, iconType, defaultColor, additionalFieldName) {
    if (iconType === 'account') {
        return accountIconStyle(color, defaultColor, additionalFieldName);
    } else if (iconType === 'category') {
        return categoryIconStyle(color, defaultColor, additionalFieldName);
    } else {
        return defaultIconStyle(color, defaultColor, additionalFieldName);
    }
}
