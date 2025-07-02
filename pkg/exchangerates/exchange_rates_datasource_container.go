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
		Container.Current = newCommonHttpExchangeRatesDataSource(&ReserveBankOfAustraliaDataSource{})
		return nil
	} else if config.ExchangeRatesDataSource == settings.BankOfCanadaDataSource {
		Container.Current = newCommonHttpExchangeRatesDataSource(&BankOfCanadaDataSource{})
		return nil
	} else if config.ExchangeRatesDataSource == settings.CzechNationalBankDataSource {
		Container.Current = newCommonHttpExchangeRatesDataSource(&CzechNationalBankDataSource{})
		return nil
	} else if config.ExchangeRatesDataSource == settings.DanmarksNationalbankDataSource {
		Container.Current = newCommonHttpExchangeRatesDataSource(&DanmarksNationalbankDataSource{})
		return nil
	} else if config.ExchangeRatesDataSource == settings.EuroCentralBankDataSource {
		Container.Current = newCommonHttpExchangeRatesDataSource(&EuroCentralBankDataSource{})
		return nil
	} else if config.ExchangeRatesDataSource == settings.NationalBankOfGeorgiaDataSource {
		Container.Current = newCommonHttpExchangeRatesDataSource(&NationalBankOfGeorgiaDataSource{})
		return nil
	} else if config.ExchangeRatesDataSource == settings.CentralBankOfHungaryDataSource {
		Container.Current = newCommonHttpExchangeRatesDataSource(&CentralBankOfHungaryDataSource{})
		return nil
	} else if config.ExchangeRatesDataSource == settings.BankOfIsraelDataSource {
		Container.Current = newCommonHttpExchangeRatesDataSource(&BankOfIsraelDataSource{})
		return nil
	} else if config.ExchangeRatesDataSource == settings.CentralBankOfMyanmarDataSource {
		Container.Current = newCommonHttpExchangeRatesDataSource(&CentralBankOfMyanmarDataSource{})
		return nil
	} else if config.ExchangeRatesDataSource == settings.NorgesBankDataSource {
		Container.Current = newCommonHttpExchangeRatesDataSource(&NorgesBankDataSource{})
		return nil
	} else if config.ExchangeRatesDataSource == settings.NationalBankOfPolandDataSource {
		Container.Current = newCommonHttpExchangeRatesDataSource(&NationalBankOfPolandDataSource{})
		return nil
	} else if config.ExchangeRatesDataSource == settings.NationalBankOfRomaniaDataSource {
		Container.Current = newCommonHttpExchangeRatesDataSource(&NationalBankOfRomaniaDataSource{})
		return nil
	} else if config.ExchangeRatesDataSource == settings.BankOfRussiaDataSource {
		Container.Current = newCommonHttpExchangeRatesDataSource(&BankOfRussiaDataSource{})
		return nil
	} else if config.ExchangeRatesDataSource == settings.SwissNationalBankDataSource {
		Container.Current = newCommonHttpExchangeRatesDataSource(&SwissNationalBankDataSource{})
		return nil
	} else if config.ExchangeRatesDataSource == settings.NationalBankOfUkraineDataSource {
		Container.Current = newCommonHttpExchangeRatesDataSource(&NationalBankOfUkraineDataSource{})
		return nil
	} else if config.ExchangeRatesDataSource == settings.CentralBankOfUzbekistanDataSource {
		Container.Current = newCommonHttpExchangeRatesDataSource(&CentralBankOfUzbekistanDataSource{})
		return nil
	} else if config.ExchangeRatesDataSource == settings.InternationalMonetaryFundDataSource {
		Container.Current = newCommonHttpExchangeRatesDataSource(&InternationalMonetaryFundDataSource{})
		return nil
	} else if config.ExchangeRatesDataSource == settings.UserCustomExchangeRatesDataSource {
		Container.Current = newUserCustomExchangeRatesDataSource()
		return nil
	} else if config.ExchangeRatesDataSource == settings.BankOfVietnamDataSource {
		Container.Current = newCommonHttpExchangeRatesDataSource(&InternationalTechcombankDataSource{})
		return nil
	}

	return errs.ErrInvalidExchangeRatesDataSource
}
