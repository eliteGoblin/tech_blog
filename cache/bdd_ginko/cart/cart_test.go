package cart

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCart(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shopping Cart Suite")
}

var _ = Describe("Shopping Cart", func() {
	itemA := Item{ID: "itemA", Name: "Item A", Price: 10.20, Qty: 0}
	itemB := Item{ID: "itemB", Name: "Item B", Price: 7.66, Qty: 0}

	Context("initially", func() {
		cart := Cart{}
		It("has 0 items", func() {
			Expect(cart.TotalUniqueItems()).Should(BeZero())
		})
		It("has 0 units", func() {
			Expect(cart.TotalUnits()).Should(BeZero())
		})
		Specify("the total amount is 0.00", func() {
			Expect(cart.TotalAmount()).Should(BeZero())
		})
	})
	Context("when a new item is added", func() {
		cart := Cart{}
		originalItemCount := cart.TotalUniqueItems()
		originalUnitCount := cart.TotalUnits()
		originalAmount := cart.TotalAmount()

		cart.AddItem(itemA)
		Context("the shopping cart", func() {
			It("has 1 more unique item than it had earlier", func() {
				Expect(cart.TotalUniqueItems()).Should(Equal(originalItemCount + 1))
			})
			It("has 1 more unit than it had earlier", func() {
				Expect(cart.TotalUnits()).Should(Equal(originalUnitCount + 1))
			})
			Specify("total amount increases by item price", func() {
				Expect(cart.TotalAmount).Should(Equal(originalAmount + itemA.Price))
			})
		})
	})
	Context("when an existing item is added", func() {
		Context("the shopping cart", func() {
			It("has the same number of unique items as earlier", func() {})
			It("has 1 more unit than it had earlier", func() {})
			Specify("total amount increases by item price", func() {})
		})
	})
	Context("that has 0 unit of item A", func() {
		Context("removing item A", func() {
			It("Should not change the number of items", func() {})
			It("Should not change the number of units", func() {})
			It("Should not change the amount", func() {})
		})
	})
})
