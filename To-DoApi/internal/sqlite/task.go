package sqlite

import (
	"log"
	"strconv"
)

type Task struct {
	TaskID      int    `json:"id,omitempty" validate:"numeric,omitempty"`
	Title       string `json:"title,omitempty" validate:"required"`
	Description string `json:"description,omitempty"`
	Priorety    string `json:"priorety,omitempty"`
	Date        string `json:"date,omitempty"`
}

type Tasks struct {
	Total int    `json:"total"`
	Task  []Task `json:"tasks"`
}

func (d *Database) GetTasks(UserID int) (*Tasks, error) {
	ex, err := d.Connect.Query("SELECT id, title, description, priorety, date FROM tasks WHERE userid=?", UserID)
	if err != nil {
		return nil, err
	}

	var tasks Tasks
	for ex.Next() {
		var idData, title, description, priorety, date string

		err = ex.Scan(&idData, &title, &description, &priorety, &date)
		if err != nil {
			return nil, err
		}

		id, err := strconv.Atoi(idData)
		if err != nil {
			return nil, err
		}

		tasks.Total++
		tasks.Task = append(tasks.Task, Task{TaskID: id, Title: title, Description: description, Priorety: priorety, Date: date})
	}

	return &tasks, nil
}

func (d *Database) GetTask(UserID int, TaskID int) (*Task, error) {
	var title, description, priorety, date string

	err := d.Connect.QueryRow("SELECT title, description, priorety, date FROM tasks WHERE userid=? AND id=?",
		UserID, TaskID).Scan(&title, &description, &priorety, &date)
	if err != nil {
		return nil, err
	}

	return &Task{Title: title, Description: description, Priorety: priorety, Date: date}, nil
}

func (d *Database) AddTask(UserID int, task Task) (int, error) {
	ex, err := d.Connect.Exec("INSERT INTO tasks (UserID, title, description, priorety, date) VALUES (?, ?, ?, ?, ?)",
		UserID, task.Title, task.Description, task.Priorety, task.Date)
	if err != nil {
		return 0, err
	}

	id, err := ex.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (d *Database) UpdateTitle(TaskID int, UserID int, value string) error {
	_, err := d.Connect.Exec("UPDATE tasks SET title=? WHERE id=? AND userid=?", value, TaskID, UserID)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) UpdateDesc(TaskID int, UserID int, value string) error {
	_, err := d.Connect.Exec("UPDATE tasks SET description=? WHERE id=? AND userid=?", value, TaskID, UserID)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) UpdatePriorety(TaskID int, UserID int, value string) error {
	_, err := d.Connect.Exec("UPDATE tasks SET priorety=? WHERE id=? AND userid=?", value, TaskID, UserID)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func (d *Database) DeleteTask(TaskID int, UserID int) error {
	_, err := d.Connect.Exec("DELETE FROM tasks WHERE id=? AND userid=?", TaskID, UserID)
	if err != nil {
		return err
	}

	return nil
}
