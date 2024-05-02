package pkg

type KoboShelf struct {
	ID           string
	Name         string
	InternalName string
	Type         string
	IsDeleted    bool
}

type KoboShelfContent struct {
	ShelfName string
	ContentID string
	IsDeleted bool
}
