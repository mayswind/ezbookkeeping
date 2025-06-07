import { type JestConfigWithTsJest, createDefaultEsmPreset } from 'ts-jest';

const presetConfig = createDefaultEsmPreset({
    tsconfig: '<rootDir>/tsconfig.jest.json'
});

const config: JestConfigWithTsJest = {
    ...presetConfig,
    clearMocks: true,
    collectCoverage: false,
    moduleNameMapper: {
        "^@/(.*)$": "<rootDir>/src/$1"
    },
    testEnvironment: "node",
    testMatch: [
        "**/__tests__/**/*.[jt]s?(x)",
        "!**/__tests__/*_gen.[jt]s?(x)"
    ]
};

export default config;
