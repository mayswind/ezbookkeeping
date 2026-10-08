<template>
    <v-row>
        <v-col cols="12">
            <v-card :title="tt('Business Features')">
                <template #subtitle>
                    {{ tt('Turn on inventory, sales and customers if you run a business. They stay hidden otherwise.') }}
                </template>
                <v-card-text>
                    <v-switch color="primary" hide-details inset
                              :label="tt('Enable inventory, sales and customers')"
                              :disabled="worksForSomeoneElse"
                              :model-value="available"
                              @update:model-value="setEnabled(!!$event)" />

                    <v-alert class="mt-4" type="info" variant="tonal" density="compact" v-if="worksForSomeoneElse">
                        {{ tt('These features are on because you were invited to work in a business. They stay on while you are a member.') }}
                    </v-alert>

                    <div class="text-body-1 mt-4">{{ tt('When enabled you get:') }}</div>
                    <ul class="ms-6 mt-1">
                        <li>{{ tt('Inventory: items, stock levels and several locations') }}</li>
                        <li>{{ tt('Sales: a cart that records the sale in your accounts, also on credit') }}</li>
                        <li>{{ tt('Customers: who owes you money, and repayments') }}</li>
                        <li>{{ tt('Team: invite managers and staff to work in your business') }}</li>
                    </ul>

                    <div class="text-body-2 text-medium-emphasis mt-4">
                        {{ tt('This setting is remembered in this browser for your account.') }}
                    </div>
                </v-card-text>
            </v-card>
        </v-col>
    </v-row>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';

import { useExtI18n } from '@/ext/i18n.ts';
import { useBusiness } from '@/ext/business.ts';
import { useBusinessFeatures } from '@/ext/features.ts';

const { tt } = useExtI18n();
const { ensureLoaded } = useBusiness();
const { available, worksForSomeoneElse, setEnabled } = useBusinessFeatures();

onMounted(() => {
    ensureLoaded().catch(() => {
        // the switch still works without the list of businesses
    });
});
</script>
