<template>
    <v-dialog width="800" :persistent="true" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-4">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <h4 class="text-h4">{{ tt('Generate MCP token') }}</h4>
                </div>
            </template>

            <v-card-text class="py-0 w-100 d-flex justify-center" v-if="generatedToken && serverUrl">
                <v-switch class="bidirectional-switch" color="secondary"
                          :label="tt('Configuration')"
                          v-model="showConfiguration"
                          @click="showConfiguration = !showConfiguration">
                    <template #prepend>
                        <span>{{ tt('Token') }}</span>
                    </template>
                </v-switch>
            </v-card-text>

            <v-card-text class="my-md-4 w-100 d-flex justify-center">
                <div class="w-100" v-if="!generatedToken">
                    <v-text-field
                        autocomplete="current-password"
                        type="password"
                        :autofocus="true"
                        :disabled="generating"
                        :label="tt('Current Password')"
                        :placeholder="tt('Current Password')"
                        v-model="currentPassword"
                        @keyup.enter="generateToken"
                    />
                </div>
                <div class="w-100 code-container" v-if="generatedToken">
                    <v-textarea class="w-100 always-cursor-text" :readonly="true"
                                :rows="4" :value="generatedToken" v-if="!showConfiguration || !serverUrl" />
                    <v-textarea class="w-100 always-cursor-text" :readonly="true"
                                :rows="15" :value="mcpServerConfiguration" v-if="showConfiguration && serverUrl" />
                </div>
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div ref="buttonContainer" class="w-100 d-flex justify-center gap-4">
                    <v-btn :disabled="generating || !currentPassword" @click="generateToken" v-if="!generatedToken">
                        {{ tt('Generate') }}
                        <v-progress-circular indeterminate size="22" class="ms-2" v-if="generating"></v-progress-circular>
                    </v-btn>
                    <v-btn color="secondary" variant="tonal" :disabled="generating"
                           @click="cancel" v-if="!generatedToken">{{ tt('Cancel') }}</v-btn>
                    <v-btn variant="tonal" @click="copy" v-if="generatedToken">{{ tt('Copy') }}</v-btn>
                    <v-btn color="secondary" variant="tonal" @click="close" v-if="generatedToken">{{ tt('Close') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTokensStore } from '@/stores/token.ts';

import { copyTextToClipboard } from '@/lib/ui/common.ts';

type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();

const tokensStore = useTokensStore();

const buttonContainer = useTemplateRef<HTMLElement>('buttonContainer');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

let resolveFunc: (() => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const showState = ref<boolean>(false);
const currentPassword = ref<string>('');
const generating = ref<boolean>(false);
const showConfiguration = ref<boolean>(false);
const serverUrl = ref<string>('');
const generatedToken = ref<string>('');

const mcpServerConfiguration = computed<string>(() => {
    return '{\n' +
        '    "mcpServers": {\n' +
        '        "ezbookkeeping-mcp": {\n' +
        '            "type": "streamable-http",\n' +
        '            "url": "' + serverUrl.value + '",\n' +
        '            "headers": {\n' +
        '                "Authorization": "Bearer ' + generatedToken.value + '"\n' +
        '            }\n' +
        '        }\n' +
        '    }\n' +
        '}'
});

function open(): Promise<void> {
    showState.value = true;
    currentPassword.value = '';
    generating.value = false;
    showConfiguration.value = false;
    serverUrl.value = '';
    generatedToken.value = '';

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function generateToken(): void {
    if (generating.value || !currentPassword.value) {
        return;
    }

    generating.value = true;

    tokensStore.generateMCPToken({
        password: currentPassword.value
    }).then(result => {
        generating.value = false;
        currentPassword.value = '';
        serverUrl.value = result.mcpUrl;
        generatedToken.value = result.token;
    }).catch(error => {
        generating.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function copy(): void {
    if (showConfiguration.value) {
        copyTextToClipboard(mcpServerConfiguration.value, buttonContainer.value);
    } else {
        copyTextToClipboard(generatedToken.value, buttonContainer.value);
    }

    snackbar.value?.showMessage('Data copied');
}

function cancel(): void {
    rejectFunc?.();
    showState.value = false;
}

function close(): void {
    resolveFunc?.();
    showState.value = false;
}

defineExpose({
    open
});
</script>
