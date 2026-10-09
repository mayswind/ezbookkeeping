import { describe, expect, test } from 'vitest';

import details from '../../../legal/details.json';
import termsHtml from '../../../public/legal/terms.html?raw';
import privacyHtml from '../../../public/legal/privacy.html?raw';

import { LEGAL_PRIVACY_URL, LEGAL_TERMS_URL, LEGAL_VERSION } from '@/ext/legalVersion.ts';

// The pages are generated from legal/details.json (python3 scripts/build-legal.py). These tests catch a forgotten rebuild.
describe('legal pages', () => {
    test('the app asks people to accept the version that is published', () => {
        expect(LEGAL_VERSION).toBe(details.version);
        expect(termsHtml).toContain(`version ${details.version}`);
        expect(privacyHtml).toContain(`version ${details.version}`);
    });

    test('the links point at pages that exist', () => {
        expect(LEGAL_TERMS_URL).toBe('/legal/terms.html');
        expect(LEGAL_PRIVACY_URL).toBe('/legal/privacy.html');
        expect(termsHtml).toContain('Terms of Service');
        expect(privacyHtml).toContain('Privacy Policy');
    });

    test('nothing is left unfilled in the template syntax', () => {
        for (const page of [termsHtml, privacyHtml]) {
            expect(page).not.toContain('{{');
            expect(page).not.toContain('<!--IF');
        }
    });

    test('the pages explain the things a customer needs to know', () => {
        expect(privacyHtml).toContain('Amazon Web Services');
        expect(privacyHtml).toContain('Render');
        expect(privacyHtml).toContain(String(details.backup_days));
        expect(termsHtml).toContain('not accounting, tax, legal or financial advice');
        expect(termsHtml).toContain('export');
    });
});
