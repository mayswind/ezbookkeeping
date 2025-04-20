package exchangerates

import (
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/settings"
)

// ExchangeRatesDataSourceContainer contains the current exchange rates data source
type ExchangeRatesDataSourceContainer struct {
	Current ExchangeRatesDataSource
}

// Initialize a exchange rates data source container singleton instance
var (
	Container = &ExchangeRatesDataSourceContainer{}
)

// InitializeExchangeRatesDataSource initializes the current exchange rates data source according to the config
func InitializeExchangeRatesDataSource(config *settings.Config) error {
	if config.ExchangeRatesDataSource == settings.ReserveBankOfAustraliaDataSource {
		Container.Current = &ReserveBankOfAustraliaDataSource{}
		return nil
	} else if config.ExchangeRatesDataSource == settings.BankOfCanadaDataSource {
		Container.Current = &BankOfCanadaDataSource{}
		return nil
	} else if config.ExchangeRatesDataSource == settings.CzechNationalBankDataSource {
		Container.Current = &CzechNationalBankDataSource{}
		return nil
	} else if config.ExchangeRatesDataSource == settings.DanmarksNationalbankDataSource {
		Container.Current = &DanmarksNationalbankDataSource{}
		return nil
	} else if config.ExchangeRatesDataSource == settings.EuroCentralBankDataSource {
		Container.Current = &EuroCentralBankDataSource{}
		return nil
	} else if config.ExchangeRatesDataSource == settings.NationalBankOfGeorgiaDataSource {
		Container.Current = &NationalBankOfGeorgiaDataSource{}
		return nil
	} else if config.ExchangeRatesDataSource == settings.CentralBankOfHungaryDataSource {
		Container.Current = &CentralBankOfHungaryDataSource{}
		return nil
	} else if config.ExchangeRatesDataSource == settings.BankOfIsraelDataSource {
		Container.Current = &BankOfIsraelDataSource{}
		return nil
	} else if config.ExchangeRatesDataSource == settings.CentralBankOfMyanmarDataSource {
		Container.Current = &CentralBankOfMyanmarDataSource{}
		return nil
	} else if config.ExchangeRatesDataSource == settings.NorgesBankDataSource {
		Container.Current = &NorgesBankDataSource{}
		return nil
	} else if config.ExchangeRatesDataSource == settings.NationalBankOfPolandDataSource {
		Container.Current = &NationalBankOfPolandDataSource{}
		return nil
	} else if config.ExchangeRatesDataSource == settings.NationalBankOfRomaniaDataSource {
		Container.Current = &NationalBankOfRomaniaDataSource{}
		return nil
	} else if config.ExchangeRatesDataSource == settings.BankOfRussiaDataSource {
		Container.Current = &BankOfRussiaDataSource{}
		return nil
	} else if config.ExchangeRatesDataSource == settings.SwissNationalBankDataSource {
		Container.Current = &SwissNationalBankDataSource{}
		return nil
	} else if config.ExchangeRatesDataSource == settings.NationalBankOfUkraineDataSource {
		Container.Current = &NationalBankOfUkraineDataSource{}
		return nil
	} else if config.ExchangeRatesDataSource == settings.CentralBankOfUzbekistanDataSource {
		Container.Current = &CentralBankOfUzbekistanDataSource{}
		return nil
	} else if config.ExchangeRatesDataSource == settings.InternationalMonetaryFundDataSource {
		Container.Current = &InternationalMonetaryFundDataSource{}
		return nil
	}

	return errs.ErrInvalidExchangeRatesDataSource
}
