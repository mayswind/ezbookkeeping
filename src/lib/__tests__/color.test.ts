import { describe, expect, it } from 'vitest';

import { getContrastTextColor, getContrastIconColor } from '@/lib/color.ts';

describe('getContrastTextColor', () => {
    it('returns black for a light background', () => {
        expect(getContrastTextColor('edddcd')).toBe('000000');
    });

    it('returns white for a dark background', () => {
        expect(getContrastTextColor('112233')).toBe('ffffff');
    });
});

describe('getContrastIconColor', () => {
    it('uses the original icon color for the default light background', () => {
        expect(getContrastIconColor('edddcd')).toBe('7f5e4b');
    });

    it('lightens a dark background with the original default opacity', () => {
        expect(getContrastIconColor('7f5e4b')).toBe('ffffff99');
    });
});
