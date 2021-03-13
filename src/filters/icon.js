import accountIcon from './accountIcon.js';
import categoryIcon from './categoryIcon.js';

export default function (iconId, iconType) {
    if (iconType === 'account') {
        return accountIcon(iconId);
    } else if (iconType === 'category') {
        return categoryIcon(iconId);
    } else {
        return '';
    }
}
