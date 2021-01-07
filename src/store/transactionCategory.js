import services from '../lib/services.js';
import logger from '../lib/logger.js';

import {
    UPDATE_TRANSACTION_CATEGORY_LIST_INVALID_STATE,
} from './mutations.js';

function addTransactionCategoryBatch(context, { categories }) {
    return new Promise((resolve, reject) => {
        services.addTransactionCategoryBatch({
            categories: categories
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                reject({ message: 'Unable to add preset categories' });
                return;
            }

            context.commit(UPDATE_TRANSACTION_CATEGORY_LIST_INVALID_STATE);

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to add preset categories', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                reject({ message: 'Unable to add preset categories' });
            } else {
                reject(error);
            }
        });
    });
}

export default {
    addTransactionCategoryBatch,
}
