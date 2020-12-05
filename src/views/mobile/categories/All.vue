<template>
    <f7-page>
        <f7-navbar :title="$t('Transaction Categories')" :back-link="$t('Back')"></f7-navbar>

        <f7-card>
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item :title="$t('Expense')" link="/category/list?type=1"></f7-list-item>
                    <f7-list-item :title="$t('Income')" link="/category/list?type=2"></f7-list-item>
                    <f7-list-item :title="$t('Transfer')" link="/category/list?type=3"></f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>
    </f7-page>
</template>
