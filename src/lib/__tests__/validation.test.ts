import { describe, expect, it } from 'vitest';

import { isValidEmail, isValidPassword } from '@/lib/validation.ts';

describe('isValidEmail', () => {
    it.each([
        'foo@bar.com',
        'foo@1.2.3.4',
        'foo_bar@foo.bar'
    ])('should accept a valid email address: %s', email => {
        expect(isValidEmail(email)).toBe(true);
    });

    it.each([
        '',
        'foo',
        '@bar',
        'foo@bar',
        'foo@bar.'
    ])('should reject an invalid email address: %s', email => {
        expect(isValidEmail(email)).toBe(false);
    });

    it('should reject an email address longer than 100 characters', () => {
        expect(isValidEmail(`${'a'.repeat(89)}@example.com`)).toBe(false);
    });
});

describe('isValidPassword', () => {
    it.each([
        '123456',
        'a'.repeat(128)
    ])('should accept a valid password', password => {
        expect(isValidPassword(password)).toBe(true);
    });

    it.each([
        '',
        '12345',
        'a'.repeat(129)
    ])('should reject a password outside the allowed length', password => {
        expect(isValidPassword(password)).toBe(false);
    });
});
