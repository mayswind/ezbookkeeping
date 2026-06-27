import { describe, expect, it } from 'vitest';

import { isTextualUUID } from '@/lib/common.ts';

describe('isTextualUUID', () => {
	it('should return false for null', () => {
		expect(isTextualUUID(null)).toBe(false);
	});

	it('should return false for undefined', () => {
		expect(isTextualUUID(undefined)).toBe(false);
	});

	it('should return false for number', () => {
		expect(isTextualUUID(12345)).toBe(false);
	});

	it('should return false for boolean', () => {
		expect(isTextualUUID(true)).toBe(false);
	});

	it('should return false for object', () => {
		expect(isTextualUUID({})).toBe(false);
	});

	it('should return false for empty string', () => {
		expect(isTextualUUID('')).toBe(false);
	});

	it('should return false for string without hyphens', () => {
		expect(isTextualUUID('notauuid')).toBe(false);
	});

	it('should return false for UUID with wrong number of segments', () => {
		expect(isTextualUUID('1234-5678')).toBe(false);
	});

	it('should return false for UUID with wrong segment lengths', () => {
		expect(isTextualUUID('1234567-1234-1234-1234-123456789012')).toBe(false);
	});

	it('should return false for UUID containing non-hex characters', () => {
		expect(isTextualUUID('1234567g-1234-1234-1234-123456789012')).toBe(false);
	});

	it('should return false for UUID with uppercase non-hex characters', () => {
		expect(isTextualUUID('1234567G-1234-1234-1234-123456789012')).toBe(false);
	});

	it('should return true for valid lowercase UUID', () => {
		expect(isTextualUUID('550e8400-e29b-41d4-a716-446655440000')).toBe(true);
	});

	it('should return true for valid uppercase UUID', () => {
		expect(isTextualUUID('550E8400-E29B-41D4-A716-446655440000')).toBe(true);
	});

	it('should return true for valid mixed-case UUID', () => {
		expect(isTextualUUID('550e8400-E29B-41d4-A716-446655440000')).toBe(true);
	});

	it('should return true for all-zero UUID', () => {
		expect(isTextualUUID('00000000-0000-0000-0000-000000000000')).toBe(true);
	});

	it('should return true for all-f UUID', () => {
		expect(isTextualUUID('ffffffff-ffff-ffff-ffff-ffffffffffff')).toBe(true);
	});
});
