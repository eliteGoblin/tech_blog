package books

type Book struct{
	Name string
	Title string
	Author string
	Pages int
}

func (b Book)CategoryByLength() string {
	if b.Pages > 300 {
		return "NOVEL"
	}
	return "SHORT STORY"
}