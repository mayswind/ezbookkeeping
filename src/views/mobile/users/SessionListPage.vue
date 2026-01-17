<template>
    <f7-page ptr @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :class="{ 'disabled': loading }" :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Device & Sessions')"></f7-nav-title>
            <f7-nav-right :class="{ 'disabled': loading }">
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
                          :title="tt(session.deviceName)"
                          :text="session.deviceInfo"
                          :key="session.tokenId"
                          v-for="session in sessions">
                <template #media>
                    <f7-icon :f7="session.icon"></f7-icon>
                </template>
                <template #after>
                    <small>{{ session.lastSeenDateTime }}</small>
                </template>
                <f7-swipeout-actions :left="textDirection === TextDirection.RTL"
                                     :right="textDirection === TextDirection.LTR"
                                     v-if="!session.isCurrent">
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

import { itemAndIndex, reversedItemAndIndex } from '@/core/base.ts';
import { TextDirection } from '@/core/text.ts';
import { type TokenInfoResponse, SessionDeviceType, SessionInfo } from '@/models/token.ts';

import { isEquals } from '@/lib/common.ts';
import { parseDateTimeFromUnixTime } from '@/lib/datetime.ts';
import { parseSessionInfo } from '@/lib/session.ts';

class MobilePageSessionInfo extends SessionInfo {
    public readonly domId: string;
    public readonly icon: string;
    public readonly lastSeenDateTime: string;

    public constructor(sessionInfo: SessionInfo) {
        super(sessionInfo.tokenId, sessionInfo.isCurrent, sessionInfo.deviceType, sessionInfo.deviceInfo, sessionInfo.deviceName, sessionInfo.lastSeen);
        this.domId = getTokenDomId(sessionInfo.tokenId);
        this.icon = getTokenIcon(sessionInfo.deviceType);
        this.lastSeenDateTime = sessionInfo.lastSeen ? formatDateTimeToLongDateTime(parseDateTimeFromUnixTime(sessionInfo.lastSeen)) : '-';
    }
}

const props = defineProps<{
    f7router: Router.Router;
}>();

const {
    tt,
    getCurrentLanguageTextDirection,
    formatDateTimeToLongDateTime
} = useI18n();

const { showConfirm, showToast, routeBackOnError } = useI18nUIComponents();

const tokensStore = useTokensStore();

const tokens = ref<TokenInfoResponse[]>([]);
const loading = ref<boolean>(true);
const loadingError = ref<unknown | null>(null);

const textDirection = computed<TextDirection>(() => getCurrentLanguageTextDirection());

const sessions = computed<MobilePageSessionInfo[]>(() => {
    const sessions: MobilePageSessionInfo[] = [];

    if (!tokens.value) {
        return sessions;
    }

    for (const token of tokens.value) {
        const sessionInfo = parseSessionInfo(token);
        sessions.push(new MobilePageSessionInfo(sessionInfo));
    }

    return sessions;
});

function getTokenIcon(deviceType: SessionDeviceType): string {
    if (deviceType === SessionDeviceType.Phone) {
        return 'device_phone_portrait';
    } else if (deviceType === SessionDeviceType.Wearable) {
        return 'device_phone_portrait';
    } else if (deviceType === SessionDeviceType.Tablet) {
        return 'device_tablet_portrait';
    } else if (deviceType === SessionDeviceType.TV) {
        return 'tv';
    } else if (deviceType === SessionDeviceType.Api) {
        return 'chevron_left_slash_chevron_right';
    } else if (deviceType === SessionDeviceType.MCP) {
        return 'sparkles';
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
                for (const [ token, index ] of itemAndIndex(tokens.value)) {
                    if (token.tokenId === session.tokenId) {
                        tokens.value.splice(index, 1);
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

            for (const [ token, index ] of reversedItemAndIndex(tokens.value)) {
                if (!token.isCurrent) {
                    tokens.value.splice(index, 1);
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
