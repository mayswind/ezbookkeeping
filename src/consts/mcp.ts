import type { NameAndDisplayName } from '@/core/base.ts';

export const ALL_MCP_TOOLS: NameAndDisplayName[] = [
    {
        name: 'query_all_accounts',
        displayName:'Query All Accounts'
    },
    {
        name: 'query_all_accounts_balance',
        displayName:'Query All Account Balances'
    },
    {
        name: 'query_all_transaction_categories',
        displayName:'Query All Transaction Categories'
    },
    {
        name: 'query_all_transaction_tags',
        displayName:'Query All Transaction Tags'
    },
    {
        name: 'query_latest_exchange_rates',
        displayName:'Query Latest Exchange Rates'
    },
    {
        name: 'query_transactions',
        displayName:'Query Transactions'
    },
    {
        name: 'add_transaction',
        displayName:'Add Transaction'
    },
    {
        name: 'delete_transaction',
        displayName:'Delete Transaction'
    }
];

export const ALL_MCP_TOOL_NAMES: string[] = ALL_MCP_TOOLS.map(tool => tool.name);

export const DEFAULT_ENABLED_MCP_TOOL_NAMES: string[] = [
    'query_all_accounts',
    'query_all_accounts_balance',
    'query_all_transaction_categories',
    'query_all_transaction_tags',
    'query_latest_exchange_rates',
    'query_transactions',
    'add_transaction'
]
