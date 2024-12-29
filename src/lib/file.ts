import { isString } from './common.ts';

export function getFileExtension(filename: string): string {
    if (!filename || !isString(filename)) {
        return '';
    }

    const parts = filename.split('.');
    return parts[parts.length - 1];
}

export function isFileExtensionSupported(filename: string, supportedExtensions: string): boolean {
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
