package todo_test

import (
	todo "TodoList"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestList_Add(t *testing.T) {
	l := todo.List{}

	taskName := "New Task"
	l.Add(taskName)

	if l[0].Task != taskName {
		t.Errorf("Expect %q, got %q instead.", taskName, l[0].Task)
	}
	if l[0].Done {
		t.Errorf("New task should not be completed.")
	}
	minTime := time.Time{}
	if l[0].CompletedAt != minTime {
		t.Errorf("New task should not have a completed timestamp.")
	}
}

func TestList_Complete(t *testing.T) {
	l := todo.List{}

	taskName := "New Task"
	l.Add(taskName)

	err := l.Complete(1)

	if err != nil {
		t.Errorf("Completing task produced error %v", err)
	}
	if !l[0].Done {
		t.Errorf("Task should be completed.")
	}
	minTime := time.Time{}
	if l[0].CompletedAt == minTime {
		t.Errorf("CompletedAt should be set.")
	}
}

func TestList_Delete(t *testing.T) {
	l := todo.List{}

	tasks := []string{
		"New Task 1",
		"New Task 2",
		"New Task 3",
	}

	for _, v := range tasks {
		l.Add(v)
	}

	err := l.Delete(2)

	if err != nil {
		t.Errorf("Deleting task produced error %v", err)
	}

	if len(l) != 2 {
		t.Errorf("Expected List length %d, got %d instead.", 2, len(l))
	}

	if l[1].Task != tasks[2] {
		t.Errorf("Expected %q, got %q instead", tasks[2], l[1].Task)
	}
}

func TestList_SaveGet(t *testing.T) {
	l1 := todo.List{}
	l2 := todo.List{}

	taskName := "New Task"
	l1.Add(taskName)

	tf, err := os.CreateTemp("", "")

	if err != nil {
		t.Errorf("Error creating temp file.")
	}

	fmt.Println(tf.Name())

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Errorf("Error removing temp file.")
		}
	}(tf.Name())

	if err := l1.Save(tf.Name()); err != nil {
		t.Fatalf("Error saving List to file: %s", err)
	}
	if err := l2.Get(tf.Name()); err != nil {
		t.Fatalf("Error getting List from file: %s", err)
	}
	if l1[0].Task != l2[0].Task {
		t.Errorf("Task %q should match %q task", l1[0].Task, l2[0].Task)
	}
}
