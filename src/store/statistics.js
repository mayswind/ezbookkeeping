import services from '../lib/services.js';
import logger from '../lib/logger.js';

import {
    LOAD_TRANSACTION_STATISTICS,
    INIT_TRANSACTION_STATISTICS_FILTER,
    UPDATE_TRANSACTION_STATISTICS_FILTER,
    UPDATE_TRANSACTION_STATISTICS_INVALID_STATE
} from "./mutations.js";

export function loadTransactionStatistics(context, { defaultCurrency }) {
    return new Promise((resolve, reject) => {
        services.getTransactionStatistics({
            startTime: context.state.transactionStatisticsFilter.startTime,
            endTime: context.state.transactionStatisticsFilter.endTime
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                reject({ message: 'Unable to get transaction statistics' });
                return;
            }

            context.commit(LOAD_TRANSACTION_STATISTICS, {
                statistics: data.result,
                defaultCurrency: defaultCurrency
            });

            if (context.state.transactionStatisticsStateInvalid) {
                context.commit(UPDATE_TRANSACTION_STATISTICS_INVALID_STATE, false);
            }

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to get transaction statistics', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                reject({ message: 'Unable to get transaction statistics' });
            } else {
                reject(error);
            }
        });
    });
}

export function initTransactionStatisticsFilter(context, filter) {
    context.commit(INIT_TRANSACTION_STATISTICS_FILTER, filter);
}

export function updateTransactionStatisticsFilter(context, filter) {
    context.commit(UPDATE_TRANSACTION_STATISTICS_FILTER, filter);
}
