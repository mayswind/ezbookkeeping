import { describe, expect, it } from 'vitest';

import { NormalizedText } from '@/core/text.ts';

describe('normalizeTextForSearch', () => {
	it('should return an empty string unchanged', () => {
		expect(NormalizedText.of('').normalizedText).toBe('');
	});

	it('should return lowercase ASCII text unchanged', () => {
		expect(NormalizedText.of('already lowercase 123').normalizedText).toBe('already lowercase 123');
	});

	it('should convert uppercase ASCII characters to lowercase', () => {
		expect(NormalizedText.of('Hello WORLD').normalizedText).toBe('hello world');
	});

	it('should strip diacritics from Unicode characters', () => {
		expect(NormalizedText.of('Café À LA').normalizedText).toBe('cafe a la');
	});

	it('should normalize Unicode compatibility characters', () => {
		expect(NormalizedText.of('ＦＯＯ ﬃ').normalizedText).toBe('foo ffi');
	});

	it('should apply Unicode case folding', () => {
		expect(NormalizedText.of('Straße').normalizedText).toBe('strasse');
	});

	it('should replace Japanese hiragana with katakana', () => {
		expect(NormalizedText.of('ひらかな').normalizedText).toBe('ヒラカナ');
	});

	it('should normalize mixed text', () => {
		expect(NormalizedText.of('Café ひらかな Straße').normalizedText).toBe('cafe ヒラカナ strasse');
	});
});
