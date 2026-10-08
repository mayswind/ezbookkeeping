import { useI18n as useVueI18n } from 'vue-i18n';

import { useI18n } from '@/locales/helpers.ts';

// Texts of the ext screens live in src/ext/locales/<language>.json, not in the app's own locale files, so merging
// upstream never touches them. As in the app, the English sentence is the key; missing languages and missing keys
// fall back to English. After adding a label, run scripts/ext-extract-i18n.py to register it.
const files = import.meta.glob<Record<string, unknown>>('./locales/*.json', { eager: true, import: 'default' });

let registered = false;

function registerMessages(): void {
    const i18n = useVueI18n();

    for (const [path, messages] of Object.entries(files)) {
        // "./locales/zh_Hans.json" -> "zh-Hans", the language key the app uses
        const language = path.replace('./locales/', '').replace('.json', '').replace('_', '-');
        i18n.mergeLocaleMessage(language, messages);
    }

    registered = true;
}

/** The app's i18n helpers, after making sure the ext messages are registered. Call it from setup(). */
export function useExtI18n() {
    if (!registered) {
        registerMessages();
    }

    const helpers = useI18n();

    return {
        ...helpers,
        /** "Owner", "Manager" or "Staff" */
        roleLabel: (role: string): string => helpers.tt(`ext.role.${role}`)
    };
}
