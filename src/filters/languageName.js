import { allLanguages } from '../locales/index.js';

export default function (languageCode) {
    const lang = allLanguages[languageCode];

    if (!lang) {
        return '';
    }

    return lang.displayName;
}
