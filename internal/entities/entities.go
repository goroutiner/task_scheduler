package entities

// Task является структкурой задачи.
type Task struct {
	Id      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// Result является структкурой необходимой для сериализации http ответа сревера.
type Result struct {
	Tasks []Task `json:"tasks"`
	Id    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
	Token string `json:"token,omitempty"`
}

// EnvMap словарь для хранения всех переменные окружения из .env.
var EnvMap map[string]string
