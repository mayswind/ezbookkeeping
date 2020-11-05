<template>
    <f7-page>
        <f7-navbar :title="$t('New Transaction')" :back-link="$t('Back')">
            <f7-subnavbar>
                <f7-segmented raised>
                    <f7-button tab-link="#expense" :text="$t('Expense')" tab-link-active></f7-button>
                    <f7-button tab-link="#income" :text="$t('Income')"></f7-button>
                    <f7-button tab-link="#transfer" :text="$t('Transfer')"></f7-button>
                </f7-segmented>
            </f7-subnavbar>
        </f7-navbar>

        <f7-tabs>
            <f7-tab id="expense" tab-active class="page-content">
                <f7-list media-list class="skeleton-text">
                    <f7-list-item title="Placeholder"></f7-list-item>
                </f7-list>
            </f7-tab>
            <f7-tab id="income" class="page-content">
                <f7-list media-list class="skeleton-text">
                    <f7-list-item title="Placeholder"></f7-list-item>
                </f7-list>
            </f7-tab>
            <f7-tab id="transfer" class="page-content">
                <f7-list media-list class="skeleton-text">
                    <f7-list-item title="Placeholder"></f7-list-item>
                </f7-list>
            </f7-tab>
        </f7-tabs>
    </f7-page>
</template>
