package asky

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
	choiceOptional  bool
	cursorIndicator string
	selectionMarker string
	disabledMarker  string
	pageSize        int
	selectedChoice  Choice
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
