package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gotk3/gotk3/gtk"
)

type Task struct {
	ID       int    `json:"id"`
	Note     string `json:"note"`
	Complete bool   `json:"complete"`
}

var tasks []Task

const filename = "tasks.json"

func saveTasks() {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	_ = ioutil.WriteFile(filename, data, 0644)
}

func loadTasks() {
	data, err := ioutil.ReadFile(filename)
	if err == nil {
		_ = json.Unmarshal(data, &tasks)
	}
}

func main() {
	gtk.Init(nil)
	loadTasks()

	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	win.SetTitle("GtkTask")
	win.SetDefaultSize(400, 300)
	win.Connect("destroy", func() { gtk.MainQuit() })

	vbox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	win.Add(vbox)

	listBox, _ := gtk.ListBoxNew()
	vbox.PackStart(listBox, true, true, 0)

	entry, _ := gtk.EntryNew()
	vbox.PackStart(entry, false, false, 0)

	// Функция обновления списка задач
	var updateList func()
	updateList = func() {
		// Удаляем все элементы из ListBox
		children := listBox.GetChildren()
		for children != nil {
			listBox.Remove(children.Data().(gtk.IWidget))
			children = children.Next()
		}

		// Добавляем задачи в список
		for i := 0; i < len(tasks); i++ {
			t := &tasks[i]

			hbox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
			label, _ := gtk.LabelNew(fmt.Sprintf("%d. %s (Done: %t)", t.ID, t.Note, t.Complete))
			hbox.PackStart(label, true, true, 0)

			// Кнопка "✔" (отметить как выполненную)
			doneButton, _ := gtk.ButtonNewWithLabel("✔")
			index := i // фиксируем индекс перед передачей в лямбду
			doneButton.Connect("clicked", func() {
				tasks[index].Complete = true
				saveTasks()
				updateList()
			})
			hbox.PackStart(doneButton, false, false, 0)

			// Кнопка "❌" (удалить)
			deleteButton, _ := gtk.ButtonNewWithLabel("❌")
			deleteButton.Connect("clicked", func() {
				tasks = append(tasks[:index], tasks[index+1:]...)
				saveTasks()
				updateList()
			})
			hbox.PackStart(deleteButton, false, false, 0)

			listBox.Add(hbox)
		}
		listBox.ShowAll()
	}

	// Кнопка добавления задачи
	addButton, _ := gtk.ButtonNewWithLabel("Add Task")
	addButton.Connect("clicked", func() {
		text, _ := entry.GetText()
		if text != "" {
			tasks = append(tasks, Task{ID: len(tasks) + 1, Note: text, Complete: false})
			entry.SetText("")
			saveTasks()
			updateList()
		}
	})

	vbox.PackStart(addButton, false, false, 0)
	win.ShowAll()
	updateList()
	gtk.Main()
}
