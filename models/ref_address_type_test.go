package models

import "testing"

func NewAddressTest(t *testing.T) *Address {
	return NewAddress(newDbForTest(t))
}

func NewAddressTypeTest(t *testing.T) *AddressType {
	return NewAddressType(newDbForTest(t))
}

func NewCustAddressTest(t *testing.T) *CustomerAdd {
	return NewCustomerAdd(newDbForTest(t))
}

func TestNewAddress(t *testing.T) {
	t.Log("Address: Testing creating a new address")

	address := NewAddressTest(t)

	user := userForTesting(t)

	address.UserID = user.UserID

	address.Country = "USA"
	address.State = "Kentucky"
	address.City = "Louisvile"
	address.Zip = 40217
	address.Line1 = "822 Milton St"
	address.Line2 = "apt 1"

	err := address.CreateAddress(nil)

	if err != nil {
		t.Fatalf("Error Creating Address %v\n", err)
	}

	address_type := NewAddressTypeTest(t)

	address_type.AddressDesc = "Billling Address"
	address_type.AddressName = "Billing"

	err = address_type.CreateAddressType(nil)

	if err != nil {
		t.Fatalf("Error Creating Address Type %v\n", err)
	}

	cust_address := NewCustAddressTest(t)

	cust_address.AddressId = address.AddressID
	cust_address.AddressType = address_type.AddressTypeID

	cust_address.CreateCustomerAdd(nil)
	if err != nil {
		t.Fatalf("Error Creating Customer Address %v\n", err)
	}

	//cleanup
	_, err = address_type.DeleteById(nil, address_type.AddressTypeID)
	if err != nil {
		t.Fatalf("Deleting Address Type by id should not fail. Error: %v", err)
	}
	_, err = address.DeleteById(nil, address.AddressID)
	if err != nil {
		t.Fatalf("Deleting Adddress by id should not fail. Error: %v", err)
	}
	_, err = cust_address.DeleteById(nil, cust_address.CustAddID)
	if err != nil {
		t.Fatalf("Deleting Customer address by id should not fail. Error: %v", err)
	}
	_, err = user.DeleteById(nil, user.UserID)
	if err != nil {
		t.Fatalf("Deleting user by id should not fail. Error: %v", err)
	}

}
