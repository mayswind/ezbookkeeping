import chardet, { type Match } from 'chardet';

import { UTF_8, CHARDET_ENCODING_NAME_MAPPING } from '@/consts/file.ts';

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
