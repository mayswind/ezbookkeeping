<template>
    <f7-page ptr @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Device & Sessions')"></f7-nav-title>
            <f7-nav-right>
                <f7-link :class="{ 'disabled': sessions.length < 2 }" :text="tt('Logout All')" @click="revokeAll"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list strong inset dividers media-list class="margin-top skeleton-text" v-if="loading">
            <f7-list-item class="list-item-media-valign-middle"
                          title="Current"
                          text="Device Name (Browser xx.x.xxxx.xx)">
                <template #media>
                    <f7-icon f7="device_phone_portrait"></f7-icon>
                </template>
                <template #after>
                    <small>MM/DD/YYYY HH:mm:ss</small>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-list strong inset dividers media-list class="margin-top" v-else-if="!loading">
            <f7-list-item class="list-item-media-valign-middle" swipeout
                          :id="session.domId"
                          :title="session.deviceType === 'mcp' ? 'MCP' : (tt(session.isCurrent ? 'Current' : 'Other Device'))"
                          :text="session.deviceInfo"
                          :key="session.tokenId"
                          v-for="session in sessions">
                <template #media>
                    <f7-icon :f7="session.icon"></f7-icon>
                </template>
                <template #after>
                    <small>{{ session.lastSeenDateTime }}</small>
                </template>
                <f7-swipeout-actions right v-if="!session.isCurrent">
                    <f7-swipeout-button color="red" :text="tt('Log Out')" @click="revoke(session)"></f7-swipeout-button>
                </f7-swipeout-actions>
            </f7-list-item>
        </f7-list>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading, onSwipeoutDeleted } from '@/lib/ui/mobile.ts';

import { useTokensStore } from '@/stores/token.ts';

import { type TokenInfoResponse, SessionInfo } from '@/models/token.ts';

import { isEquals } from '@/lib/common.ts';
import { parseSessionInfo } from '@/lib/session.ts';

class MobilePageSessionInfo extends SessionInfo {
    public readonly domId: string;
    public readonly icon: string;
    public readonly lastSeenDateTime: string;

    public constructor(sessionInfo: SessionInfo) {
        super(sessionInfo.tokenId, sessionInfo.isCurrent, sessionInfo.deviceType, sessionInfo.deviceInfo, sessionInfo.createdByCli, sessionInfo.lastSeen);
        this.domId = getTokenDomId(sessionInfo.tokenId);
        this.icon = getTokenIcon(sessionInfo.deviceType);
        this.lastSeenDateTime = sessionInfo.lastSeen ? formatUnixTimeToLongDateTime(sessionInfo.lastSeen) : '-';
    }
}

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt, formatUnixTimeToLongDateTime } = useI18n();
const { showConfirm, showToast, routeBackOnError } = useI18nUIComponents();

const tokensStore = useTokensStore();

const tokens = ref<TokenInfoResponse[]>([]);
const loading = ref<boolean>(true);
const loadingError = ref<unknown | null>(null);

const sessions = computed<MobilePageSessionInfo[]>(() => {
    const sessions: MobilePageSessionInfo[] = [];

    if (!tokens.value) {
        return sessions;
    }

    for (let i = 0; i < tokens.value.length; i++) {
        const token = tokens.value[i];
        const sessionInfo = parseSessionInfo(token);
        sessions.push(new MobilePageSessionInfo(sessionInfo));
    }

    return sessions;
});

function getTokenIcon(deviceType: string): string {
    if (deviceType === 'phone') {
        return 'device_phone_portrait';
    } else if (deviceType === 'wearable') {
        return 'device_phone_portrait';
    } else if (deviceType === 'tablet') {
        return 'device_tablet_portrait';
    } else if (deviceType === 'tv') {
        return 'tv';
    } else if (deviceType === 'mcp') {
        return 'wand_stars';
    } else if (deviceType === 'cli') {
        return 'chevron_left_slash_chevron_right';
    } else {
        return 'device_desktop';
    }
}

function getTokenDomId(tokenId: string): string {
    return 'token_' + tokenId.replace(/:/g, '_');
}

function init(): void {
    loading.value = true;

    tokensStore.getAllTokens().then(response => {
        tokens.value = response;
        loading.value = false;
    }).catch(error => {
        if (error.processed) {
            loading.value = false;
        } else {
            loadingError.value = error;
            showToast(error.message || error);
        }
    });
}

function reload(done?: () => void): void {
    tokensStore.getAllTokens().then(response => {
        done?.();

        if (isEquals(response, tokens.value)) {
            showToast('Session list is up to date');
        } else {
            showToast('Session list has been updated');
        }

        tokens.value = response;
    }).catch(error => {
        done?.();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function revoke(session: SessionInfo): void {
    showConfirm('Are you sure you want to logout from this session?', () => {
        showLoading();

        tokensStore.revokeToken({
            tokenId: session.tokenId
        }).then(() => {
            hideLoading();

            onSwipeoutDeleted(getTokenDomId(session.tokenId), () => {
                for (let i = 0; i < tokens.value.length; i++) {
                    if (tokens.value[i].tokenId === session.tokenId) {
                        tokens.value.splice(i, 1);
                    }
                }
            });
        }).catch(error => {
            hideLoading();

            if (!error.processed) {
                showToast(error.message || error);
            }
        });
    });
}

function revokeAll(): void {
    if (tokens.value.length < 2) {
        return;
    }

    showConfirm('Are you sure you want to logout all other sessions?', () => {
        showLoading();

        tokensStore.revokeAllTokens().then(() => {
            hideLoading();

            for (let i = tokens.value.length - 1; i >= 0; i--) {
                if (!tokens.value[i].isCurrent) {
                    tokens.value.splice(i, 1);
                }
            }

            showToast('You have logged out all other sessions');
        }).catch(error => {
            hideLoading();

            if (!error.processed) {
                showToast(error.message || error);
            }
        });
    });
}

function onPageAfterIn(): void {
    routeBackOnError(props.f7router, loadingError);
}

init();
</script>
