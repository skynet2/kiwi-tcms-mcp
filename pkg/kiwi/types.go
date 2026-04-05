package kiwi

type TestCase struct {
	ID          int64  `json:"id"`
	Summary     string `json:"summary"`
	Notes       string `json:"notes"`
	Text        string `json:"text"`
	CaseStatus  int64  `json:"case_status"`
	Category    int64  `json:"category"`
	Priority    int64  `json:"priority"`
	Author      int64  `json:"author"`
	IsAutomated bool   `json:"is_automated"`
}

type TestCaseCreateInput struct {
	Summary     string `json:"summary"`
	Category    int64  `json:"category"`
	Priority    int64  `json:"priority"`
	CaseStatus  int64  `json:"case_status"`
	Notes       string `json:"notes,omitempty"`
	Text        string `json:"text,omitempty"`
	IsAutomated bool   `json:"is_automated,omitempty"`
}

type Category struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Product int64  `json:"product"`
}

type Priority struct {
	ID    int64  `json:"id"`
	Value string `json:"value"`
}

type Product struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Classification int64  `json:"classification"`
}

type Component struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Product int64  `json:"product"`
}

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type TestCaseStatus struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
