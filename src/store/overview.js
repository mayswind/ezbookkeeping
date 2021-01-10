import services from '../lib/services.js';
import logger from '../lib/logger.js';

import {
    LOAD_TRANSACTION_OVERVIEW,
    UPDATE_TRANSACTION_OVERVIEW_INVALID_STATE,
} from './mutations.js';

export function loadTransactionOverview(context, { dateRange, force }) {
    if (!force && !context.state.transactionOverviewStateInvalid) {
        return new Promise((resolve) => {
            resolve(context.state.transactionOverview);
        });
    }

    return new Promise((resolve, reject) => {
        services.getTransactionOverview({
            today: dateRange.today,
            thisWeek: dateRange.thisWeek,
            thisMonth: dateRange.thisMonth,
            thisYear: dateRange.thisYear
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                reject({ message: 'Unable to get transaction overview' });
                return;
            }

            context.commit(LOAD_TRANSACTION_OVERVIEW, data.result);

            if (context.state.transactionOverviewStateInvalid) {
                context.commit(UPDATE_TRANSACTION_OVERVIEW_INVALID_STATE, false);
            }

            resolve(data.result);
        }).catch(error => {
            if (force) {
                logger.error('failed to force load transaction overview', error);
            } else {
                logger.error('failed to load transaction overview', error);
            }

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                reject({ message: 'Unable to get transaction overview' });
            } else {
                reject(error);
            }
        });
    });
}
