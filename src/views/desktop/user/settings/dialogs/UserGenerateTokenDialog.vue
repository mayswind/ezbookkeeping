<template>
    <v-dialog width="800" :persistent="true" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-4">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <h4 class="text-h4">{{ tt('Generate Token') }}</h4>
                </div>
            </template>

            <v-card-text class="py-0 w-100 d-flex justify-center" v-if="tokenType === 'api' && generatedToken && serverUrl">
                <v-switch class="bidirectional-switch" color="secondary"
                          :label="tt('Example')"
                          v-model="showAPIExample"
                          @click="showAPIExample = !showAPIExample">
                    <template #prepend>
                        <span>{{ tt('Token') }}</span>
                    </template>
                </v-switch>
            </v-card-text>

            <v-card-text class="py-0 w-100 d-flex justify-center" v-if="tokenType === 'mcp' && generatedToken && serverUrl">
                <v-switch class="bidirectional-switch" color="secondary"
                          :label="tt('Configuration')"
                          v-model="showMCPConfiguration"
                          @click="showMCPConfiguration = !showMCPConfiguration">
                    <template #prepend>
                        <span>{{ tt('Token') }}</span>
                    </template>
                </v-switch>
            </v-card-text>

            <v-card-text class="mt-md-4 w-100 d-flex justify-center">
                <div class="w-100" v-if="!generatedToken">
                    <v-row>
                        <v-col cols="12" md="12">
                            <v-select
                                item-title="name"
                                item-value="value"
                                :disabled="generating"
                                :label="tt('Token Type')"
                                :placeholder="tt('Token Type')"
                                :items="tokenTypeOptions"
                                v-model="tokenType"
                            />
                        </v-col>
                        <v-col cols="12" :md="tokenExpirationTime >= 0 ? 12 : 6">
                            <v-select
                                item-title="displayName"
                                item-value="value"
                                :disabled="generating"
                                :label="tt('Expiration Time')"
                                :placeholder="tt('Expiration Time')"
                                :items="[
                                    { displayName: tt('No Expiration'), value: 0 },
                                    { displayName: tt('format.misc.nHour', { n: formatNumberToLocalizedNumerals(1) }), value: 3600 },
                                    { displayName: tt('format.misc.nDays', { n: formatNumberToLocalizedNumerals(1) }), value: 86400 },
                                    { displayName: tt('format.misc.nDays', { n: formatNumberToLocalizedNumerals(7) }), value: 604800 },
                                    { displayName: tt('format.misc.nDays', { n: formatNumberToLocalizedNumerals(30) }), value: 2592000 },
                                    { displayName: tt('format.misc.nDays', { n: formatNumberToLocalizedNumerals(90) }), value: 7776000 },
                                    { displayName: tt('format.misc.nDays', { n: formatNumberToLocalizedNumerals(180) }), value: 15552000 },
                                    { displayName: tt('format.misc.nDays', { n: formatNumberToLocalizedNumerals(365) }), value: 31536000 },
                                    { displayName: tt('Custom'), value: -1 }
                                ]"
                                v-model="tokenExpirationTime"
                            />
                        </v-col>
                        <v-col cols="12" md="6" v-if="tokenExpirationTime < 0">
                            <number-input
                                :persistent-placeholder="true"
                                :disabled="generating"
                                :label="tt('Custom Expiration Time (Seconds)')"
                                :placeholder="tt('Custom Expiration Time (Seconds)')"
                                v-model="tokenCustomExpirationTime"
                            />
                        </v-col>
                        <v-col cols="12" md="12">
                            <v-text-field
                                autocomplete="current-password"
                                type="password"
                                persistent-placeholder
                                :autofocus="true"
                                :disabled="generating"
                                :label="tt('Current Password')"
                                :placeholder="tt('Current Password')"
                                v-model="currentPassword"
                                @keyup.enter="generateToken"
                            />
                        </v-col>
                    </v-row>
                </div>

                <div class="w-100 code-container" v-if="generatedToken">
                    <v-textarea class="w-100 always-cursor-text" :readonly="true"
                                :rows="4" :value="generatedToken" v-if="(tokenType === 'api' && (!showAPIExample || !serverUrl)) || (tokenType === 'mcp' && (!showMCPConfiguration || !serverUrl))" />
                    <v-textarea class="w-100 always-cursor-text" :readonly="true"
                                :rows="5" :value="apiExample" v-if="tokenType === 'api' && showAPIExample && serverUrl" />
                    <v-textarea class="w-100 always-cursor-text" :readonly="true"
                                :rows="15" :value="mcpServerConfiguration" v-if="tokenType === 'mcp' && showMCPConfiguration && serverUrl" />
                </div>
            </v-card-text>

            <v-card-text class="w-100" v-if="(tokenExpirationTime === 0 || (tokenExpirationTime < 0 && tokenCustomExpirationTime === 0)) || tokenType === 'mcp'">
                <v-alert type="warning" variant="tonal">
                    <span v-if="tokenExpirationTime === 0 || (tokenExpirationTime < 0 && tokenCustomExpirationTime === 0)">{{ tt('Your token does not expire, please keep it secure.') }}</span>
                    <span v-if="tokenType === 'mcp'">{{ tt('When connecting to third-party apps, be aware that they and any large language models they use can access your private data.') }}</span>
                </v-alert>
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

import { type NameValue } from '@/core/base.ts';
import { type TokenGenerateAPIResponse, type TokenGenerateMCPResponse } from '@/models/token.ts';

import { isAPITokenEnabled, isMCPServerEnabled } from '@/lib/server_settings.ts';
import { copyTextToClipboard } from '@/lib/ui/common.ts';

type SnackBarType = InstanceType<typeof SnackBar>;

const { tt, formatNumberToLocalizedNumerals } = useI18n();

const tokensStore = useTokensStore();

const buttonContainer = useTemplateRef<HTMLElement>('buttonContainer');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

let resolveFunc: (() => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const showState = ref<boolean>(false);
const tokenType = ref<'api' | 'mcp'>(isAPITokenEnabled() ? 'api' : (isMCPServerEnabled() ? 'mcp' : 'api'));
const tokenExpirationTime = ref<number>(86400);
const tokenCustomExpirationTime = ref<number>(86400);
const currentPassword = ref<string>('');
const generating = ref<boolean>(false);
const showAPIExample = ref<boolean>(false);
const showMCPConfiguration = ref<boolean>(false);
const serverUrl = ref<string>('');
const generatedToken = ref<string>('');

const tokenTypeOptions = computed<NameValue[]>(() => {
    const options: NameValue[] = [];

    if (isAPITokenEnabled()) {
        options.push({ name: tt('API Token'), value: 'api' });
    }

    if (isMCPServerEnabled()) {
        options.push({ name: tt('MCP Token'), value: 'mcp' });
    }

    return options;
});

const apiExample = computed<string>(() => {
    return `curl -H 'Authorization: Bearer ${generatedToken.value}' '${serverUrl.value}/v1/users/profile/get.json'`;
});

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
    tokenType.value = isAPITokenEnabled() ? 'api' : (isMCPServerEnabled() ? 'mcp' : 'api');
    tokenExpirationTime.value = 86400;
    tokenCustomExpirationTime.value = 86400;
    generating.value = false;
    showAPIExample.value = false;
    showMCPConfiguration.value = false;
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

    tokensStore.generateToken({
        type: tokenType.value,
        expiresInSeconds: tokenExpirationTime.value >= 0 ? tokenExpirationTime.value : tokenCustomExpirationTime.value,
        password: currentPassword.value
    }).then(result => {
        generating.value = false;
        currentPassword.value = '';

        if (tokenType.value === 'api') {
            serverUrl.value = (result as TokenGenerateAPIResponse).apiBaseUrl;
        } else if (tokenType.value === 'mcp') {
            serverUrl.value = (result as TokenGenerateMCPResponse).mcpUrl;
        }

        generatedToken.value = result.token;
    }).catch(error => {
        generating.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function copy(): void {
    if (tokenType.value === 'api' && showAPIExample.value) {
        copyTextToClipboard(apiExample.value, buttonContainer.value);
    } else if (tokenType.value === 'mcp' && showMCPConfiguration.value) {
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
