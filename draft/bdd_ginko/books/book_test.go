package books_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "bddtest/books"
)

var _ = Describe("Book", func() {
	var (
		longBook  Book
		shortBook Book
	)
	BeforeEach(func() {
		longBook = Book{
			Title:  "Les Miserables",
			Author: "Victor Hugo",
			Pages:  1488,
		}

		shortBook = Book{
			Title:  "Fox In Socks",
			Author: "Dr. Seuss",
			Pages:  24,
		}
		fmt.Println("level 1")
	})

	Describe("Categorizing book length", func() {
		//BeforeEach(func(){
		//	fmt.Println("mid level")
		//})
		Context("With more than 300 pages", func() {
			It("should be a novel", func() {
				Expect(longBook.CategoryByLength()).To(Equal("NOVEL"))
			})
		})
		//
		Context("With fewer than 300 pages", func() {
			BeforeEach(func(){
				fmt.Println("bottom level")
			})
			It("should be a short story", func() {
				Expect(shortBook.CategoryByLength()).To(Equal("SHORT STORY"))
			})
			It("should be a short story", func() {
				Expect(shortBook.CategoryByLength()).To(Equal("SHORT STORY"))
			})
		})
	})
	Describe("Categorizing book length2", func() {

		Context("With more than 300 pages", func() {
			It("should be a novel", func() {
				Expect(longBook.CategoryByLength()).To(Equal("NOVEL"))
			})
		})
		Context("With more than 300 pages", func() {
			It("should be a novel", func() {
				Expect(longBook.CategoryByLength()).To(Equal("NOVEL"))
			})
		})
		//Context("With more than 300 pages", func() {
		//	It("should be a novel", func() {
		//		Expect(longBook.CategoryByLength()).To(Equal("NOVEL"))
		//	})
		//})

	})
})
