import { isString } from './common.js';

export function getFileExtension(filename) {
    if (!filename || !isString(filename)) {
        return '';
    }

    const parts = filename.split('.');
    return parts[parts.length - 1];
}

export function isFileExtensionSupported(filename, supportedExtensions) {
    if (!supportedExtensions) {
        return false;
    }

    const supportedExtensionsArray = supportedExtensions.split(',');
    const fileExtension = getFileExtension(filename).toLowerCase();

    for (let i = 0; i < supportedExtensionsArray.length; i++) {
        const supportedExtension = getFileExtension(supportedExtensionsArray[i]).toLowerCase();

        if (supportedExtension === fileExtension) {
            return true;
        }
    }

    return false;
}
