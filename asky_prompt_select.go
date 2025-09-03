package asky

type SelectionOption struct {
	Value    string
	Label    string
	Disabled bool
}

type Choice struct {
	Value    string
	Label    string
	Disabled bool
}

type SingleSelect struct {
	theme           Theme
	prefix          string
	label           string
	separator       string
	help            string
	choices         []Choice
	defaultChoice   int
	cursorIndicator string
	selectionMarker string
	pageSize        int
	selectedChoice  int
}

type MultiSelect struct {
	theme              Theme
	prefix             string
	label              string
	separator          string
	help               string
	choices            []Choice
	defaultChoices     []int
	minChoicesRequired int
	maxChoicesAllowed  int
	selectionMarker    string
	pageSize           int
	selectedChoices    []int
}
