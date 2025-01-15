import type { PresetCategory } from '@/core/category.ts';

export const DEFAULT_EXPENSE_CATEGORIES: PresetCategory[] = [
    {
        name: 'Food & Drink',
        categoryIconId: '1',
        color: 'ff6b22',
        subCategories: [
            {
                name: 'Food',
                categoryIconId: '2',
                color: 'ff6b22'
            },
            {
                name: 'Drink',
                categoryIconId: '30',
                color: 'ff6b22'
            },
            {
                name: 'Fruit & Snack',
                categoryIconId: '70',
                color: 'ff6b22'
            }
        ]
    },
    {
        name: 'Clothing & Appearance',
        categoryIconId: '100',
        color: '673ab7',
        subCategories: [
            {
                name: 'Clothing',
                categoryIconId: '110',
                color: '673ab7'
            },
            {
                name: 'Jewelry',
                categoryIconId: '170',
                color: '673ab7'
            },
            {
                name: 'Cosmetic',
                categoryIconId: '180',
                color: '673ab7'
            },
            {
                name: 'Hair Cuts & Salon',
                categoryIconId: '190',
                color: '673ab7'
            }
        ]
    },
    {
        name: 'Housing & Houseware',
        categoryIconId: '200',
        color: '000000',
        subCategories: [
            {
                name: 'Houseware',
                categoryIconId: '210',
                color: '000000'
            },
            {
                name: 'Electronics',
                categoryIconId: '230',
                color: '000000'
            },
            {
                name: 'Repairs & Maintenance',
                categoryIconId: '250',
                color: '000000'
            },
            {
                name: 'Housekeeping Services',
                categoryIconId: '260',
                color: '000000'
            },
            {
                name: 'Utilities Expense',
                categoryIconId: '270',
                color: '000000'
            },
            {
                name: 'Rent & Mortgage',
                categoryIconId: '290',
                color: '000000'
            }
        ]
    },
    {
        name: 'Transportation',
        categoryIconId: '300',
        color: '009688',
        subCategories: [
            {
                name: 'Public Transit',
                categoryIconId: '310',
                color: '009688'
            },
            {
                name: 'Taxi & Car Rental',
                categoryIconId: '320',
                color: '009688'
            },
            {
                name: 'Personal Car Expense',
                categoryIconId: '330',
                color: '009688'
            },
            {
                name: 'Train Tickets',
                categoryIconId: '370',
                color: '009688'
            },
            {
                name: 'Airline Tickets',
                categoryIconId: '390',
                color: '009688'
            }
        ]
    },
    {
        name: 'Communication',
        categoryIconId: '400',
        color: '2196f3',
        subCategories: [
            {
                name: 'Telephone Bill',
                categoryIconId: '420',
                color: '2196f3'
            },
            {
                name: 'Internet Bill',
                categoryIconId: '430',
                color: '2196f3'
            },
            {
                name: 'Express Fee',
                categoryIconId: '480',
                color: '2196f3'
            }
        ]
    },
    {
        name: 'Entertainment',
        categoryIconId: '500',
        color: 'ff2d55',
        subCategories: [
            {
                name: 'Sports & Fitness',
                categoryIconId: '510',
                color: 'ff2d55'
            },
            {
                name: 'Party Expense',
                categoryIconId: '540',
                color: 'ff2d55'
            },
            {
                name: 'Movies & Shows',
                categoryIconId: '550',
                color: 'ff2d55'
            },
            {
                name: 'Toys & Games',
                categoryIconId: '560',
                color: 'ff2d55'
            },
            {
                name: 'Subscriptions',
                categoryIconId: '570',
                color: 'ff2d55'
            },
            {
                name: 'Pet Expense',
                categoryIconId: '580',
                color: 'ff2d55'
            },
            {
                name: 'Travelling',
                categoryIconId: '590',
                color: 'ff2d55'
            }
        ]
    },
    {
        name: 'Education & Studying',
        categoryIconId: '600',
        color: 'cddc39',
        subCategories: [
            {
                name: 'Books & Newspaper & Magazines',
                categoryIconId: '610',
                color: 'cddc39'
            },
            {
                name: 'Training Courses',
                categoryIconId: '660',
                color: 'cddc39'
            },
            {
                name: 'Certification & Examination',
                categoryIconId: '680',
                color: 'cddc39'
            }
        ]
    },
    {
        name: 'Gifts & Donations',
        categoryIconId: '700',
        color: '4cd964',
        subCategories: [
            {
                name: 'Gifts',
                categoryIconId: '710',
                color: '4cd964'
            },
            {
                name: 'Donations',
                categoryIconId: '780',
                color: '4cd964'
            }
        ]
    },
    {
        name: 'Medical & Healthcare',
        categoryIconId: '800',
        color: 'ff3b30',
        subCategories: [
            {
                name: 'Diagnosis & Treatment',
                categoryIconId: '840',
                color: 'ff3b30'
            },
            {
                name: 'Medications',
                categoryIconId: '860',
                color: 'ff3b30'
            },
            {
                name: 'Medical Devices',
                categoryIconId: '890',
                color: 'ff3b30'
            }
        ]
    },
    {
        name: 'Finance & Insurance',
        categoryIconId: '900',
        color: 'ff9500',
        subCategories: [
            {
                name: 'Tax Expense',
                categoryIconId: '910',
                color: 'ff9500'
            },
            {
                name: 'Service Charge',
                categoryIconId: '930',
                color: 'ff9500'
            },
            {
                name: 'Insurance Expense',
                categoryIconId: '950',
                color: 'ff9500'
            },
            {
                name: 'Interest Expense',
                categoryIconId: '970',
                color: 'ff9500'
            },
            {
                name: 'Compensation & Fine',
                categoryIconId: '990',
                color: 'ff9500'
            }
        ]
    },
    {
        name: 'Miscellaneous',
        categoryIconId: '1000',
        color: '8e8e93',
        subCategories: [
            {
                name: 'Other Expense',
                categoryIconId: '1010',
                color: '8e8e93'
            }
        ]
    }
];

export const DEFAULT_INCOME_CATEGORIES: PresetCategory[] = [
    {
        name: 'Occupational Earnings',
        categoryIconId: '2000',
        color: 'ff6b22',
        subCategories: [
            {
                name: 'Salary Income',
                categoryIconId: '2010',
                color: 'ff6b22'
            },
            {
                name: 'Bonus Income',
                categoryIconId: '2020',
                color: 'ff6b22'
            },
            {
                name: 'Overtime Pay',
                categoryIconId: '231',
                color: 'ff6b22'
            },
            {
                name: 'Side Job Income',
                categoryIconId: '2080',
                color: 'ff6b22'
            }
        ]
    },
    {
        name: 'Finance & Investment',
        categoryIconId: '900',
        color: 'ff9500',
        subCategories: [
            {
                name: 'Investment Income',
                categoryIconId: '2100',
                color: 'ff9500'
            },
            {
                name: 'Rental Income',
                categoryIconId: '290',
                color: 'ff9500'
            },
            {
                name: 'Interest Income',
                categoryIconId: '970',
                color: 'ff9500'
            }
        ]
    },
    {
        name: 'Miscellaneous',
        categoryIconId: '1000',
        color: '8e8e93',
        subCategories: [
            {
                name: 'Gift & Lucky Money',
                categoryIconId: '710',
                color: '8e8e93'
            },
            {
                name: 'Winnings Income',
                categoryIconId: '564',
                color: '8e8e93'
            },
            {
                name: 'Windfall',
                categoryIconId: '5200',
                color: '8e8e93'
            },
            {
                name: 'Other Income',
                categoryIconId: '3010',
                color: '8e8e93'
            }
        ]
    }
];

export const DEFAULT_TRANSFER_CATEGORIES: PresetCategory[] = [
    {
        name: 'General Transfer',
        categoryIconId: '4000',
        color: 'ff6b22',
        subCategories: [
            {
                name: 'Bank Transfer',
                categoryIconId: '900',
                color: 'ff6b22'
            },
            {
                name: 'Credit Card Repayment',
                categoryIconId: '980',
                color: 'ff6b22'
            },
            {
                name: 'Deposits & Withdrawals',
                categoryIconId: '981',
                color: 'ff6b22'
            }
        ]
    },
    {
        name: 'Loan & Debt',
        categoryIconId: '950',
        color: 'ff9500',
        subCategories: [
            {
                name: 'Borrowing Money',
                categoryIconId: '910',
                color: 'ff9500'
            },
            {
                name: 'Lending Money',
                categoryIconId: '290',
                color: 'ff9500'
            },
            {
                name: 'Repayment',
                categoryIconId: '930',
                color: 'ff9500'
            },
            {
                name: 'Debt Collection',
                categoryIconId: '5030',
                color: 'ff9500'
            }
        ]
    },

    {
        name: 'Miscellaneous',
        categoryIconId: '1000',
        color: '8e8e93',
        subCategories: [
            {
                name: 'Out-of-Pocket Expense',
                categoryIconId: '2010',
                color: '8e8e93'
            },
            {
                name: 'Reimbursement',
                categoryIconId: '920',
                color: '8e8e93'
            },
            {
                name: 'Other Transfer',
                categoryIconId: '4900',
                color: '8e8e93'
            }
        ]
    }
];
