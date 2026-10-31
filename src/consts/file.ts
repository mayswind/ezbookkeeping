export const SUPPORTED_IMAGE_EXTENSIONS: string = '.jpg,.jpeg,.png,.gif,.webp';

export const UTF_8 = 'utf-8';

export const SUPPORTED_FILE_ENCODINGS: string[] = [
    UTF_8, // UTF-8
    'utf-16le', // UTF-16 Little Endian
    'utf-16be', // UTF-16 Big Endian
    'utf-32le', // UTF-32 Little Endian
    'utf-32be', // UTF-32 Big Endian
    'cp437', // OEM United States (CP-437)
    'cp863', // OEM Canadian French (CP-863)
    'cp037', // IBM EBCDIC US/Canada (CP-037)
    'cp1047', // IBM EBCDIC Open Systems (CP-1047)
    'cp1140', // IBM EBCDIC US/Canada with Euro (CP-1140)
    "iso-8859-1", // Western European (ISO-8859-1)
    'cp850', // Western European (CP-850)
    'cp858', // Western European with Euro (CP-858)
    'windows-1252', // Western European (Windows-1252)
    'iso-8859-15', // Western European (ISO-8859-15)
    'iso-8859-4', // North European (ISO-8859-4)
    'iso-8859-10', // North European (ISO-8859-10)
    'cp865', // North European (CP-865)
    'iso-8859-2', // Central European (ISO-8859-2)
    'cp852', // Central European (CP-852)
    'windows-1250', // Central European (Windows-1250)
    'iso-8859-14', // Celtic (ISO-8859-14)
    'iso-8859-3', // South European (ISO-8859-3)
    'cp860', // Portuguese (CP-860)
    'iso-8859-7', // Greek (ISO-8859-7)
    'windows-1253', // Greek (Windows-1253)
    'iso-8859-9', // Turkish (ISO-8859-9)
    'windows-1254', // Turkish (Windows-1254)
    'iso-8859-13', // Baltic (ISO-8859-13)
    'windows-1257', // Baltic (Windows-1257)
    'iso-8859-16', // South-Eastern European (ISO-8859-16)
    'iso-8859-5', // Cyrillic (ISO-8859-5)
    'cp855', // Cyrillic (CP-855)
    'cp866', // Cyrillic (CP-866)
    'windows-1251', // Cyrillic (Windows-1251)
    'koi8r', // Cyrillic (KOI8-R)
    'koi8u', // Cyrillic (KOI8-U)
    'iso-8859-6', // Arabic (ISO-8859-6)
    'windows-1256', // Arabic (Windows-1256)
    'iso-8859-8', // Hebrew (ISO-8859-8)
    'cp862', // Hebrew (CP-862)
    'windows-1255', // Hebrew (Windows-1255)
    'windows-874', // Thai (Windows-874)
    'windows-1258', // Vietnamese (Windows-1258)
    'gb18030', // Chinese (Simplified, GB18030)
    'gbk', // Chinese (Simplified, GBK)
    'big5', // Chinese (Traditional, Big5)
    'euc-kr', // Korean (EUC-KR)
    'euc-jp', // Japanese (EUC-JP)
    'iso-2022-jp', // Japanese (ISO-2022-JP)
    'shift_jis', // Japanese (Shift_JIS)
];

export const CHARDET_ENCODING_NAME_MAPPING: Record<string, string> = {
    'UTF-8': UTF_8,
    'UTF-16LE': 'utf-16le',
    'UTF-16BE': 'utf-16be',
    'UTF-32LE': 'utf-32le',
    'UTF-32BE': 'utf-32be',
    'ISO-2022-JP': 'iso-2022-jp',
    // 'ISO-2022-KR': '', // not supported
    // 'ISO-2022-CN': '', // not supported
    'Shift_JIS': 'shift_jis',
    'Big5': 'big5',
    'EUC-JP': 'euc-jp',
    'EUC-KR': 'euc-kr',
    'GB18030': 'gb18030',
    'ISO-8859-1': 'iso-8859-1',
    'ISO-8859-2': 'iso-8859-2',
    'ISO-8859-5': 'iso-8859-5',
    'ISO-8859-6': 'iso-8859-6',
    'ISO-8859-7': 'iso-8859-7',
    'ISO-8859-8': 'iso-8859-8',
    'ISO-8859-9': 'iso-8859-9',
    'windows-874': 'windows-874',
    'windows-1250': 'windows-1250',
    'windows-1251': 'windows-1251',
    'windows-1252': 'windows-1252',
    'windows-1253': 'windows-1253',
    'windows-1254': 'windows-1254',
    'windows-1255': 'windows-1255',
    'windows-1256': 'windows-1256',
    'windows-1257': 'windows-1257',
    'windows-1258': 'windows-1258',
    'KOI8-R':'koi8r'
};
