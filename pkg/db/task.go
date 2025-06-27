package db

// Task представляет собой одну задачу в планировщике.
// Теги json используются для сериализации/десериализации данных при обмене с фронтендом.
type Task struct {
	ID      int64  `json:"id,omitempty"` // omitempty, чтобы не включать в JSON, если значение нулевое
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}
