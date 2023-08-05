<template>
    <v-card :class="{ 'disabled': disabled }">
        <v-card-text class="d-flex align-center">
            <v-avatar color="secondary" size="38">
                <v-icon size="24" :icon="icon" />
            </v-avatar>
            <span class="text-base font-weight-bold ml-3">{{ title }}</span>
            <v-spacer/>
            <v-btn density="comfortable" color="default" variant="text" class="ml-2" :icon="true">
                <v-icon :icon="icons.more" />
                <v-menu activator="parent">
                    <v-list>
                        <slot name="menus"></slot>
                    </v-list>
                </v-menu>
            </v-btn>
        </v-card-text>
        <v-card-text class="mt-1 pb-2">
            <div class="font-weight-semibold text-truncate text-red text-h5 text-income me-2 mb-2" v-if="!loading || incomeAmount">{{ incomeAmount }}</div>
            <v-skeleton-loader class="skeleton-no-margin mt-4 mb-7" type="text" width="120px" :loading="true" v-else-if="loading && !incomeAmount"></v-skeleton-loader>
            <div class="text-truncate text-h6 text-expense" v-if="!loading || expenseAmount">{{ expenseAmount }}</div>
            <v-skeleton-loader class="skeleton-no-margin mb-2" type="text" width="120px" :loading="true" v-else-if="loading && !expenseAmount"></v-skeleton-loader>
        </v-card-text>
        <v-card-text class="mt-6">
            <span class="text-caption">{{ datetime }}</span>
        </v-card-text>
    </v-card>
</template>

<script>
import {
    mdiDotsVertical
} from '@mdi/js';

export default {
    props: [
        'loading',
        'disabled',
        'icon',
        'title',
        'expenseAmount',
        'incomeAmount',
        'datetime'
    ],
    data() {
        return {
            icons: {
                more: mdiDotsVertical
            }
        };
    }
}
</script>
