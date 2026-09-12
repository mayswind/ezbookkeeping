<template>
    <div class="code-editor w-100 h-100" :class="{ 'code-editor-ready': editorReady, 'code-editor-readonly': readonly }">
        <v-textarea
            no-resize
            hide-details
            class="code-editor-fallback w-100 h-100 code-textarea"
            :class="{ 'always-cursor-text': readonly }"
            :aria-label="ariaLabel"
            :readonly="readonly"
            :model-value="modelValue"
            @update:model-value="updateFallbackValue"
            v-if="!editorReady"
        />
        <div ref="editorContainer" class="code-editor-monaco w-100 h-100" :class="{ 'code-editor-monaco-ready': editorReady }"></div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, useTemplateRef, watch, nextTick, onMounted, onUnmounted } from 'vue';
import { useTheme } from 'vuetify';

import type * as Monaco from 'monaco-editor/editor';
import type { LanguageServiceDefaults } from 'monaco-editor/languages/features/typescript/register';

import { ThemeType } from '@/core/theme.ts';

import logger from '@/lib/logger.ts';

export type CodeEditorLanguage = 'javascript' | 'json';

export interface CodeEditorExtraLib {
    readonly content: string;
    readonly filePath?: string;
}

const props = defineProps<{
    ariaLabel?: string;
    readonly?: boolean;
    language: CodeEditorLanguage;
    lineNumbers?: boolean;
    extraLibs?: readonly CodeEditorExtraLib[];
    modelValue: string;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
}>();

const theme = useTheme();

const editorContainer = useTemplateRef<HTMLDivElement>('editorContainer');

let unmounted = false;
let editor: Monaco.editor.IStandaloneCodeEditor | undefined;
let model: Monaco.editor.ITextModel | undefined;
let javascriptDefaults: LanguageServiceDefaults | undefined;
let editorContentChanged: Monaco.IDisposable | undefined;
let extraLibDisposables: Monaco.IDisposable[] = [];

const editorReady = ref<boolean>(false);

const editorTheme = computed<string>(() => theme.global.name.value === ThemeType.Dark ? 'vs-dark' : 'vs');

function disposeExtraLibs(): void {
    for (const disposable of extraLibDisposables) {
        disposable.dispose();
    }

    extraLibDisposables = [];
}

function updateExtraLibs(): void {
    disposeExtraLibs();

    if (!javascriptDefaults || props.language !== 'javascript') {
        return;
    }

    if (props.extraLibs) {
        for (let i = 0; i < props.extraLibs.length; i++) {
            const extraLib = props.extraLibs[i];

            if (!extraLib) {
                continue;
            }

            extraLibDisposables.push(javascriptDefaults.addExtraLib(
                extraLib.content,
                extraLib.filePath || `inmemory://model/code-editor-extra-lib-${i}.d.ts`
            ));
        }
    }
}

function updateFallbackValue(value: string): void {
    if (!props.readonly) {
        emit('update:modelValue', value);
    }
}

onMounted(() => {
    import('@/lib/code_editor/index.ts').then(runtime => {
        if (unmounted || !editorContainer.value) {
            return;
        }

        javascriptDefaults = runtime.javascriptDefaults;
        updateExtraLibs();

        model = runtime.monaco.editor.createModel(props.modelValue, props.language);
        editor = runtime.monaco.editor.create(editorContainer.value, {
            ariaLabel: props.ariaLabel,
            automaticLayout: true,
            bracketPairColorization: { enabled: true },
            codeLens: false,
            colorDecorators: false,
            contextmenu: false,
            dragAndDrop: false,
            dropIntoEditor: {
                enabled: false
            },
            folding: !!props.lineNumbers,
            fontFamily: "Consolas, 'SFMono-Regular', Menlo, Monaco, 'Courier New', monospace",
            fontSize: 14,
            glyphMargin: false,
            inlayHints: {
                enabled: 'off'
            },
            inlineSuggest: {
                enabled: false
            },
            lineDecorationsWidth: props.lineNumbers ? 8 : (props.readonly ? 4 : 0),
            lineNumbers: props.lineNumbers ? 'on' : 'off',
            lineNumbersMinChars: props.lineNumbers ? 5 : 0,
            lightbulb: {
                enabled: runtime.monaco.editor.ShowLightbulbIconMode.Off
            },
            linkedEditing: false,
            minimap: { enabled: false },
            model: model,
            overviewRulerLanes: 0,
            pasteAs: {
                enabled: false
            },
            padding: {
                top: 2,
                bottom: 2
            },
            readOnly: props.readonly,
            renderLineHighlight: props.readonly ? 'none' : 'line',
            scrollBeyondLastLine: false,
            stickyScroll: {
                enabled: false
            },
            tabSize: 4,
            theme: editorTheme.value,
            unusualLineTerminators: 'off',
            wordWrap: 'off'
        });
        editor.addCommand(runtime.monaco.KeyCode.F1, () => undefined);

        editorContentChanged = editor.onDidChangeModelContent(() => {
            if (!props.readonly && editor) {
                emit('update:modelValue', editor.getValue());
            }
        });

        editorReady.value = true;

        nextTick(() => {
            editor?.layout()
        });
    }).catch(error => {
        logger.error('Failed to load code editor', error);
    });
});

onUnmounted(() => {
    unmounted = true;
    disposeExtraLibs();
    editorContentChanged?.dispose();
    editor?.dispose();
    model?.dispose();
    editorReady.value = false;
});

watch(() => props.modelValue, value => {
    if (editor && value !== editor.getValue()) {
        editor.setValue(value);
    }
});

watch(() => props.language, language => {
    if (model && editor) {
        const currentModel = model;

        import('@/lib/code_editor/index.ts').then(runtime => {
            runtime.monaco.editor.setModelLanguage(currentModel, language);
        }).catch(error => {
            logger.error('Failed to update code editor language', error);
        });

        updateExtraLibs();
    }
});

watch(() => props.extraLibs, () => {
    updateExtraLibs();
}, { deep: true });

watch(() => props.readonly, readonly => {
    editor?.updateOptions({
        lineDecorationsWidth: props.lineNumbers ? 8 : (readonly ? 4 : 0),
        readOnly: readonly,
        renderLineHighlight: readonly ? 'none' : 'line'
    });
});

watch(() => props.lineNumbers, lineNumbers => {
    editor?.updateOptions({
        lineDecorationsWidth: lineNumbers ? 8 : (props.readonly ? 4 : 0),
        lineNumbers: lineNumbers ? 'on' : 'off',
        lineNumbersMinChars: lineNumbers ? 3 : 0
    });
});

watch(editorTheme, value => {
    editor?.updateOptions({ theme: value });
});
</script>

<style scoped>
.code-editor {
    position: relative;
    min-height: 120px;
}

.code-editor-ready {
    overflow: hidden;
    border: thin solid rgba(var(--v-border-color), var(--v-border-opacity));
    border-radius: 6px;
}

.code-editor-readonly,
.code-editor-readonly :deep(.v-field),
.code-editor-readonly :deep(.monaco-editor),
.code-editor-readonly :deep(.monaco-editor-background),
.code-editor-readonly :deep(.monaco-editor .margin) {
    background: rgb(var(--ebk-code-background));
}

.code-editor-monaco :deep(.suggest-widget),
.code-editor-monaco :deep(.suggest-widget .monaco-list),
.code-editor-monaco :deep(.suggest-details),
.code-editor-monaco :deep(.monaco-hover) {
    background-color: rgb(var(--v-theme-surface)) !important;
}

.code-editor-monaco {
    position: absolute;
    inset: 0;
    visibility: hidden;
}

.code-editor-monaco-ready {
    visibility: visible;
}
</style>
