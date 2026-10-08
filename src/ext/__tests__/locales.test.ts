import { describe, expect, test } from 'vitest';

import en from '@/ext/locales/en.json';

// Every label the ext screens show must be registered in src/ext/locales/en.json, otherwise it cannot be translated
// and texts with placeholders do not even render. This test reads the sources, so forgetting a label fails the build.
// Fix: run  python3 scripts/ext-extract-i18n.py
const sources = import.meta.glob<string>(['../**/*.vue', '../**/*.ts', '!../__tests__/**', '!../locales/**'], { query: '?raw', import: 'default', eager: true });

const LITERAL = /\btt\(\s*(['"])((?:\\.|(?!\1).)*)\1/g;
const OPEN_KEY = /\.open\(\s*'(ext\.[^']+)'/g;

function lookup(key: string): unknown {
    if (key in en) {
        return (en as Record<string, unknown>)[key];
    }

    // nested keys such as ext.role.owner
    return key.split('.').reduce<unknown>((node, part) => (node && typeof node === 'object' ? (node as Record<string, unknown>)[part] : undefined), en);
}

function usedKeys(pattern: RegExp, group: number): { key: string, file: string }[] {
    const used: { key: string, file: string }[] = [];

    for (const [file, text] of Object.entries(sources)) {
        for (const match of text.matchAll(pattern)) {
            used.push({ key: match[group]!.replace(/\\'/g, '\'').replace(/\\"/g, '"'), file });
        }
    }

    return used;
}

describe('ext locale file', () => {
    test('the sources were found', () => {
        expect(Object.keys(sources).length).toBeGreaterThan(10);
    });

    test('every label passed to tt() is registered', () => {
        const missing = usedKeys(LITERAL, 2).filter(({ key }) => typeof lookup(key) !== 'string').map(({ key, file }) => `${key}  (${file})`);

        expect(missing).toEqual([]);
    });

    test('every confirm dialog key is registered', () => {
        const missing = usedKeys(OPEN_KEY, 1).filter(({ key }) => typeof lookup(key) !== 'string').map(({ key, file }) => `${key}  (${file})`);

        expect(missing).toEqual([]);
    });

    test('the role names used by roleLabel exist', () => {
        for (const role of ['owner', 'manager', 'staff']) {
            expect(typeof lookup(`ext.role.${role}`)).toBe('string');
        }
    });

    test('messages avoid characters that the message format treats specially', () => {
        const flat: [string, string][] = [];

        const walk = (node: unknown, path: string): void => {
            if (typeof node === 'string') {
                flat.push([path, node]);
            } else if (node && typeof node === 'object') {
                Object.entries(node).forEach(([key, value]) => walk(value, path ? `${path}.${key}` : key));
            }
        };
        walk(en, '');

        expect(flat.filter(([, text]) => /[|@$]/.test(text)).map(([path]) => path)).toEqual([]);
    });

    test('placeholders of the key and the message match', () => {
        const placeholders = (text: string): string[] => (text.match(/\{\w+\}/g) ?? []).sort();

        const mismatched = Object.entries(en)
            .filter((entry): entry is [string, string] => typeof entry[1] === 'string')
            .filter(([key, text]) => JSON.stringify(placeholders(key)) !== JSON.stringify(placeholders(text)))
            .map(([key]) => key);

        expect(mismatched).toEqual([]);
    });
});
