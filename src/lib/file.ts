import chardet, { type Match } from 'chardet';

import type { ImportFileTypeAndExtensions } from '@/core/file.ts';

import { UTF_8, CHARDET_ENCODING_NAME_MAPPING } from '@/consts/file.ts';

import { isString } from './common.ts';

export function getFileExtension(filename: string): string {
    if (!filename || !isString(filename)) {
        return '';
    }

    const parts = filename.split('.');
    return parts[parts.length - 1] as string;
}

export function findExtensionByType(items: ImportFileTypeAndExtensions[] | undefined, type: string): string | undefined {
    if (!items || items.length < 1) {
        return undefined;
    }

    for (const item of items) {
        if (item.type === type) {
            return item.extensions;
        }
    }

    return undefined;
}

export function isFileExtensionSupported(filename: string, supportedExtensions: string): boolean {
    if (!supportedExtensions) {
        return false;
    }

    const supportedExtensionsArray = supportedExtensions.split(',');
    const fileExtension = getFileExtension(filename).toLowerCase();

    for (const supportedExtension of supportedExtensionsArray) {
        if (getFileExtension(supportedExtension).toLowerCase() === fileExtension) {
            return true;
        }
    }

    return false;
}

export function detectFileEncoding(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();

        reader.onload = () => {
            const arrayBuffer = reader.result as ArrayBuffer;
            const uint8Array = new Uint8Array(arrayBuffer);
            const possibleEncodings: Match[] = chardet.analyse(uint8Array);

            if (!possibleEncodings || possibleEncodings.length < 1) {
                reject(new Error('unable to detect file encoding'));
                return;
            }

            const mostPossibleEncoding: Match = possibleEncodings[0] as Match;

            if (!mostPossibleEncoding.name || mostPossibleEncoding.confidence < 50) {
                // check whether all characters are ASCII
                let isAllAscii = true;

                for (const byte of uint8Array) {
                    if (byte > 0x7F) {
                        isAllAscii = false;
                        break;
                    }
                }

                if (isAllAscii) {
                    resolve(UTF_8);
                    return;
                }

                reject(new Error('unable to detect file encoding'));
                return;
            }

            const encoding = CHARDET_ENCODING_NAME_MAPPING[mostPossibleEncoding.name];

            if (!encoding) {
                reject(new Error(`unsupported file encoding: ${mostPossibleEncoding.name}`));
                return;
            }

            resolve(encoding);
        };

        reader.onerror = () => {
            reject(new Error('failed to read file for encoding detection'));
        };

        reader.readAsArrayBuffer(file);
    });
}
