import { useI18n as useVueI18n } from 'vue-i18n';

import { useI18n } from '@/locales/helpers.ts';

// English texts of the ext screens that contain placeholders. Texts without placeholders use the English sentence
// itself as the key, like the rest of the app. Placeholders only work for registered messages, so these have keys.
// Other languages fall back to English until they are translated.
const EXT_MESSAGES_EN = {
    ext: {
        role: {
            owner: 'Owner',
            manager: 'Manager',
            staff: 'Staff'
        },
        invitedAs: 'Invited as {role}',
        yourOwnBusiness: 'Your own business',
        confirmRemoveMember: 'Remove {name} from your business? They lose access immediately.',
        confirmLeave: 'Leave {name}? You will no longer have access to their business.',
        confirmDeleteItem: 'Delete {name}? Its stock history is kept.',
        confirmVoidSale: 'Void this sale of {total}? The stock goes back and the money is removed from your books.',
        banner: 'You are working in {name}\'s business as {role}. Everything you do here changes their books.',
        chip: '{name} ({role})'
    }
};

let registered = false;

/** The app's i18n helpers, after making sure the ext messages are registered. Call it from setup(). */
export function useExtI18n() {
    if (!registered) {
        useVueI18n().mergeLocaleMessage('en', EXT_MESSAGES_EN);
        registered = true;
    }

    const helpers = useI18n();

    return {
        ...helpers,
        /** "Owner", "Manager" or "Staff" */
        roleLabel: (role: string): string => helpers.tt(`ext.role.${role}`)
    };
}
