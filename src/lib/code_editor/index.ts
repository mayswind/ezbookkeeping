import * as monaco from 'monaco-editor/editor';

import 'monaco-editor/features/register.all';

import 'monaco-editor/languages/definitions/javascript/register';

import { ScriptTarget, javascriptDefaults } from 'monaco-editor/languages/features/typescript/register';

import EditorWorker from 'monaco-editor/editor/editor.worker?worker';
import JavaScriptWorker from 'monaco-editor/languages/features/typescript/ts.worker?worker';

type MonacoEnvironmentContainer = typeof globalThis & {
    MonacoEnvironment?: {
        getWorker: (workerId: string, label: string) => Worker;
    };
};

const globalEnvironment = globalThis as MonacoEnvironmentContainer;

globalEnvironment.MonacoEnvironment = {
    getWorker: (_workerId: string, label: string): Worker => {
        if (label === 'javascript') {
            return new JavaScriptWorker();
        }

        return new EditorWorker();
    }
};

javascriptDefaults.setCompilerOptions({
    allowJs: true,
    allowNonTsExtensions: true,
    checkJs: true,
    noEmit: true,
    target: ScriptTarget.ES2020
});

javascriptDefaults.setDiagnosticsOptions({
    noSemanticValidation: false,
    noSyntaxValidation: false,
    noSuggestionDiagnostics: false,
    onlyVisible: true
});

monaco.languages.register({ id: 'json' });
monaco.languages.setMonarchTokensProvider('json', {
    tokenizer: {
        root: [
            [/[{}[\]]/, '@brackets'],
            [/[,:]/, 'delimiter'],
            [/"(?:[^"\\]|\\.)*"(?=\s*:)/, 'key'],
            [/"(?:[^"\\]|\\.)*"/, 'string'],
            [/-?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?/, 'number'],
            [/\b(?:true|false)\b/, 'keyword'],
            [/\bnull\b/, 'keyword'],
            [/\s+/, 'white']
        ]
    }
});

export {
    javascriptDefaults,
    monaco
};
