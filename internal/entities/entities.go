package entities

// Task является структурой задачи.
type Task struct {
	Id      string `json:"id,omitempty" db:"id"`
	Date    string `json:"date,omitempty" db:"date"`
	Title   string `json:"title,omitempty" db:"title"`
	Comment string `json:"comment,omitempty" db:"comment"`
	Repeat  string `json:"repeat,omitempty" db:"repeat"`
}

// Result является структурой необходимой для сериализации http ответа сервера.
type Result struct {
	Tasks []Task `json:"tasks,omitempty"`
	Id    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
	Token string `json:"token,omitempty"`
}

var (
	TableName = "scheduler"
	UiDir     = "ui"
	DbFile    = "scheduler.db"
)
