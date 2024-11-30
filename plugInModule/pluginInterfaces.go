package plugInModule

type ExtendCommandPlugin interface {
	Extend(line string, inCodeBlock, inTable, inList, isUsecaseOrExample bool) (string, bool, bool, bool, bool)
}

type ThemePlugin interface {
	Theme(html string) string
}
