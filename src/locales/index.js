import en from './en.js'
import zhHans from './zh_Hans.js'

export const defaultLanguage = 'en';

export const allLanguages = {
    'en': {
        name: 'English',
        displayName: 'English',
        content: en
    },
    'zh-Hans': {
        name: 'Simplified Chinese',
        displayName: '简体中文',
        aliases: ['zh-CHS', 'zh-CN', 'zh-SG'],
        content: zhHans
    }
};
