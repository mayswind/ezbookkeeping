import pluginVue from 'eslint-plugin-vue';
import vueTsEslintConfig from '@vue/eslint-config-typescript';

export default [
    ...pluginVue.configs['flat/essential'],
    ...vueTsEslintConfig(),
    {
        languageOptions: {
            parserOptions: {
                projectService: true,
                tsconfigRootDir: import.meta.dirname,
            }
        },
    },
    {
        ignores: [
            'dist/**',
            '**/*.{js,jsx,cjs,mjs}'
        ]
    },
    {
        files: [
            '**/*.{vue,ts,tsx,mts,js,jsx,cjs,mjs}'
        ],
        rules: {
            '@typescript-eslint/no-this-alias': ['error', {
                allowedNames: ['self']
            }],
            'vue/valid-v-slot': ['error', {
                allowModifiers: true
            }],
            'vue/block-lang': ['error', {
                script: {
                    lang: ['ts', 'js']
                }
            }],
        }
    },
];
