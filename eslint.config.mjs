import globals from 'globals';

import path from 'node:path';
import { fileURLToPath } from 'node:url';

import js from '@eslint/js';
import { FlatCompat } from '@eslint/eslintrc';
import { includeIgnoreFile } from '@eslint/compat';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const gitignorePath = path.resolve(__dirname, '.gitignore');

const compat = new FlatCompat({
    baseDirectory: __dirname,
    recommendedConfig: js.configs.recommended,
    allConfig: js.configs.all
});

export default [...compat.extends('eslint:recommended', 'plugin:vue/vue3-essential'),
    includeIgnoreFile(gitignorePath), {
    languageOptions: {
        globals: {
            ...globals.node,
        },
    },
    files: [
        "**/*.{vue,js,jsx,cjs,mjs}"
    ],
    rules: {
        'vue/no-use-v-if-with-v-for': 'off',
        'vue/valid-v-slot': ['error', {
            allowModifiers: true,
        }],
    },
}];
