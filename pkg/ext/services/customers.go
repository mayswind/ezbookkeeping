package extservices

import (
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	exterrs "github.com/mayswind/ezbookkeeping/pkg/ext/errors"
	extmodels "github.com/mayswind/ezbookkeeping/pkg/ext/models"
)

// CustomerInput is the editable part of a customer
type CustomerInput struct {
	Name  string
	Phone string
	Email string
	Note  string
}

// CustomerService manages customers and what they owe
type CustomerService struct{}

// Customers is the customer service singleton
var Customers = &CustomerService{}

func truncate(s string, max int) string {
	s = strings.TrimSpace(s)

	if len(s) > max {
		return s[:max]
	}

	return s
}

// Add creates a customer
func (s *CustomerService) Add(c core.Context, ownerUid int64, in CustomerInput) (*extmodels.Customer, error) {
	name := truncate(in.Name, 128)

	if name == "" {
		return nil, exterrs.ErrCustomerNameIsEmpty
	}

	now := nowUnix()
	customer := &extmodels.Customer{
		OwnerUid: ownerUid, Name: name, Phone: truncate(in.Phone, 32), Email: truncate(in.Email, 128), Note: truncate(in.Note, 255),
		CreatedUnixTime: now, UpdatedUnixTime: now,
	}
	_, err := ownerDB(ownerUid).NewSession(c).Insert(customer)

	return customer, err
}

// Modify updates a customer
func (s *CustomerService) Modify(c core.Context, ownerUid int64, customerId int64, in CustomerInput) (*extmodels.Customer, error) {
	name := truncate(in.Name, 128)

	if name == "" {
		return nil, exterrs.ErrCustomerNameIsEmpty
	}

	customer, err := s.Get(c, ownerUid, customerId)

	if err != nil {
		return nil, err
	}

	customer.Name, customer.Phone, customer.Email, customer.Note = name, truncate(in.Phone, 32), truncate(in.Email, 128), truncate(in.Note, 255)
	customer.UpdatedUnixTime = nowUnix()
	_, err = ownerDB(ownerUid).NewSession(c).ID(customerId).Cols("name", "phone", "email", "note", "updated_unix_time").Update(customer)

	return customer, err
}

// Delete soft-deletes a customer who owes nothing
func (s *CustomerService) Delete(c core.Context, ownerUid int64, customerId int64) error {
	customer, err := s.Get(c, ownerUid, customerId)

	if err != nil {
		return err
	}

	owed, err := s.Outstanding(c, ownerUid, customerId)

	if err != nil {
		return err
	} else if owed > 0 {
		return exterrs.ErrCustomerHasBalance
	}

	customer.Deleted = true
	customer.DeletedUnixTime = nowUnix()
	_, err = ownerDB(ownerUid).NewSession(c).ID(customerId).Cols("deleted", "deleted_unix_time").Update(customer)

	return err
}

// Get returns one active customer
func (s *CustomerService) Get(c core.Context, ownerUid int64, customerId int64) (*extmodels.Customer, error) {
	customer := &extmodels.Customer{}
	has, err := ownerDB(ownerUid).NewSession(c).Where("owner_uid=? AND customer_id=? AND deleted=?", ownerUid, customerId, false).Get(customer)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, exterrs.ErrCustomerNotFound
	}

	return customer, nil
}

// List returns the active customers of a business
func (s *CustomerService) List(c core.Context, ownerUid int64) ([]*extmodels.Customer, error) {
	customers := make([]*extmodels.Customer, 0)
	err := ownerDB(ownerUid).NewSession(c).Where("owner_uid=? AND deleted=?", ownerUid, false).OrderBy("name, customer_id").Find(&customers)

	return customers, err
}

// Outstanding returns how much a customer owes in total
func (s *CustomerService) Outstanding(c core.Context, ownerUid int64, customerId int64) (int64, error) {
	sums, err := ownerDB(ownerUid).NewSession(c).Where("owner_uid=? AND customer_id=? AND voided=?", ownerUid, customerId, false).SumsInt(&extmodels.Sale{}, "total", "paid")

	if err != nil {
		return 0, err
	}

	return sums[0] - sums[1], nil
}

// Balances returns what every customer with sales owes, keyed by customer id
func (s *CustomerService) Balances(c core.Context, ownerUid int64) (map[int64]int64, error) {
	rows := make([]*extmodels.SaleTotals, 0)
	err := ownerDB(ownerUid).NewSession(c).Table("ext_sale").
		Select("customer_id, SUM(total) AS total, SUM(paid) AS paid").
		Where("owner_uid=? AND voided=? AND customer_id>0", ownerUid, false).
		GroupBy("customer_id").Find(&rows)

	if err != nil {
		return nil, err
	}

	balances := make(map[int64]int64, len(rows))

	for _, row := range rows {
		balances[row.CustomerId] = row.Total - row.Paid
	}

	return balances, nil
}
