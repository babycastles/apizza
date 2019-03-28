package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/harrybrwn/apizza/cmd/internal/base"
	"github.com/harrybrwn/apizza/pkg/tests"
)

func testOrderNew(t *testing.T, buf *bytes.Buffer, cmds ...base.CliCommand) {
	cart, add := cmds[0], cmds[1]
	add.Cmd().ParseFlags([]string{"--name=testorder", "--products=12SCMEATZA"})
	err := add.Run(add.Cmd(), []string{})
	if err != nil {
		t.Error(err)
	}
	buf.Reset()

	if err := cart.Run(cart.Cmd(), []string{"testorder"}); err != nil {
		t.Error(err)
	}

	expected := `testorder
  Products:
    12SCMEATZA
  StoreID: 4336
  Method:  Carryout
  Address: 1600 Pennsylvania Ave NW
           Washington DC, 20500
`

	expected = `testorder
  products: 
    Medium (12") Hand Tossed MeatZZa™
      code:     12SCMEATZA
      options:  map[]
      quantity: 1
  storeID: 4336
  method:  Carryout
  address: 1600 Pennsylvania Ave NW
  	       Washington DC, 20500
`

	tests.Compare(t, buf.String(), strings.Replace(expected, "\t", "  ", -1))
}

func testAddOrder(t *testing.T, buf *bytes.Buffer, cmds ...base.CliCommand) {
	cart, add := cmds[0], cmds[1]
	if err := add.Run(add.Cmd(), []string{"testing"}); err != nil {
		t.Error(err)
	}
	if buf.String() != "" {
		t.Errorf("wrong output: should have no output: '%s'", buf.String())
	}
	buf.Reset()

	cart.Cmd().ParseFlags([]string{"-d"})
	if err := cart.Run(cart.Cmd(), []string{"testing"}); err != nil {
		t.Error(err)
	}
	buf.Reset()
}

func testOrderNewErr(t *testing.T, buf *bytes.Buffer, cmds ...base.CliCommand) {
	if err := cmds[0].Run(cmds[0].Cmd(), []string{}); err == nil {
		t.Error("expected error")
	}
}

func testOrderRunAdd(t *testing.T, buf *bytes.Buffer, cmds ...base.CliCommand) {
	cart := cmds[0]
	if err := cart.Run(cart.Cmd(), []string{}); err != nil {
		t.Error(err)
	}

	expected := `Your Orders:
  testorder
`
	tests.Compare(t, buf.String(), expected)
	buf.Reset()

	cart.Cmd().ParseFlags([]string{"--add", "W08PBNLW,W08PPLNW"})
	if err := cart.Run(cart.Cmd(), []string{"testorder"}); err != nil {
		t.Error(err)
	}
	tests.Compare(t, buf.String(), "order successfully updated.\n")
}

func testOrderPriceOutput(cart *cartCmd, buf *bytes.Buffer, t *testing.T) {
	cart.price = true
	if err := cart.Run(cart.Cmd(), []string{"testorder"}); err != nil {
		t.Error(err)
	}

	expected := "testorder\n  price: $34.07\n  products: \n    Medium (12\") Hand Tossed MeatZZa™\n      code:     12SCMEATZA\n      options:  map[]\n      quantity: 1\n    8-Piece Boneless Chicken\n      code:     W08PBNLW\n      options:  map[]\n      quantity: 1\n    8-piece Plain Wings\n      code:     W08PPLNW\n      options:  map[]\n      quantity: 1\n  storeID: 4336\n  method:  Carryout\n  address: 1600 Pennsylvania Ave NW\n           Washington DC, 20500\n"
	tests.Compare(t, buf.String(), expected)

	if err := cart.Run(cart.Cmd(), []string{"to-many", "args"}); err == nil {
		t.Error("expected error")
	}
}

func testOrderRunDelete(cart *cartCmd, buf *bytes.Buffer, t *testing.T) {
	cart.delete = true
	if err := cart.Run(cart.Cmd(), []string{"testorder"}); err != nil {
		t.Error(err)
	}
	tests.Compare(t, buf.String(), "testorder successfully deleted.\n")
	cart.delete = false
	buf.Reset()

	cart.Cmd().ParseFlags([]string{})
	if err := cart.Run(cart.Cmd(), []string{}); err != nil {
		t.Error(err)
	}
	expected := `No orders saved.
`
	tests.Compare(t, buf.String(), expected)
	buf.Reset()

	if err := cart.Run(cart.Cmd(), []string{"not_a_real_order"}); err == nil {
		t.Error("expected error")
	}
}