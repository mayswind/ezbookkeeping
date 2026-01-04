<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Account Category Order')"></f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link icon-f7="ellipsis" @click="showMoreActionSheet = true"></f7-link>
                <f7-link icon-f7="checkmark_alt" :class="{ 'disabled': !isDisplayOrderModified() }" @click="saveDisplayOrder"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list strong inset dividers sortable sortable-enabled class="margin-top"
                 @sortable:sort="onSort">
            <f7-list-item :id="getAccountCategoryDomId(accountCategory)"
                          :key="accountCategory.type"
                          :title="tt(accountCategory.name)"
                          v-for="accountCategory in accountCategories">
            </f7-list-item>
        </f7-list>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="resetToDefault()">{{ tt('Reset to Default') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';
import { useAccountCategoryDisplayOrderSettingsPageBase } from '@/views/base/settings/AccountCategoryDisplayOrderSettingsPageBase.ts';

import { AccountCategory } from '@/core/account.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt } = useI18n();
const { showToast } = useI18nUIComponents();

const {
    accountCategories,
    isDisplayOrderModified,
    loadDisplayOrderFromSettings,
    saveDisplayOrderToSettings,
    resetDisplayOrderToDefault
} = useAccountCategoryDisplayOrderSettingsPageBase();

const showMoreActionSheet = ref<boolean>(false);

function getAccountCategoryDomId(accountCategory: AccountCategory): string {
    return 'account_category_' + accountCategory.type;
}

function parseAccountCategoryTypeFromDomId(domId: string): string | null {
    if (!domId || domId.indexOf('account_category_') !== 0) {
        return null;
    }

    return domId.substring(17); // account_category_
}

function init(): void {
    loadDisplayOrderFromSettings();
}

function saveDisplayOrder(): void {
    saveDisplayOrderToSettings();
    showToast('Account category order saved');
    props.f7router.back();
}

function resetToDefault() {
    resetDisplayOrderToDefault();
    showMoreActionSheet.value = false;
}

function onSort(event: { el: { id: string }, from: number, to: number }): void {
    if (!event || !event.el || !event.el.id) {
        showToast('Unable to move account category');
        return;
    }

    const type = parseAccountCategoryTypeFromDomId(event.el.id);

    if (!type) {
        showToast('Unable to move account category');
        return;
    }

    let currentAccountCategory: AccountCategory | null = null;

    for (const accountCategory of accountCategories.value) {
        if (accountCategory.type.toString() === type) {
            currentAccountCategory = accountCategory;
            break;
        }
    }

    if (!currentAccountCategory || !accountCategories.value[event.to]) {
        showToast('Unable to move account category');
        return;
    }

    accountCategories.value.splice(event.to, 0, accountCategories.value.splice(event.from, 1)[0] as AccountCategory);
}

init();
</script>
