import { describe, expect, it } from 'vitest';

import { getContrastTextColor, getContrastIconColor } from '@/lib/color.ts';

describe('getContrastTextColor', () => {
    it('returns black for a light background', () => {
        expect(getContrastTextColor('ffcc00')).toBe('000000');
    });

    it('returns white for a dark background', () => {
        expect(getContrastTextColor('112233')).toBe('ffffff');
    });
});

describe('getContrastIconColor', () => {
    it('uses the original icon color for the default light background', () => {
        expect(getContrastIconColor('ffcc00')).toBe('c67e48');
    });

    it('lightens a dark background with the original default opacity', () => {
        expect(getContrastIconColor('c67e48')).toBe('ffffff99');
    });
});
