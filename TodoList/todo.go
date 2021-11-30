package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// item struct represents a to-do item.
type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

// List represents a list of to-do items.
type List []item

// String prints out a formatted list.
// Implements the fmt.Stringer interface.
func (l *List) String() string {
	return l.ListDisplay(false, false)
}

func (l *List) ListDisplay(active bool, verbose bool) string {
	formatted := ""

	for k, t := range *l {
		if active && t.Done {
			continue
		}

		prefix := "  "
		if t.Done {
			prefix = "X "
		}

		suffix := ""
		if verbose {
			suffix = " (" + t.CreatedAt.Format(time.RFC822)
			if t.Done {
				suffix += " - " + t.CompletedAt.Format(time.RFC822)
			}
			suffix += ")"
		}

		formatted += fmt.Sprintf("%s%d: %s%s\n", prefix, k+1, t.Task, suffix)
	}

	return formatted
}

// Add creates a new to-do item and appends it to the List.
func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	*l = append(*l, t)
}

// Complete marks a to-do item as completed by setting
// Done = true and CompletedAt to the current time.
func (l *List) Complete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}
	ls[i-1].Done = true
	ls[i-1].CompletedAt = time.Now()

	return nil
}

// Delete removes a to-do item from the List.
func (l *List) Delete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}
	*l = append(ls[:i-1], ls[i:]...)

	return nil
}

// Save encodes the List as JSON and saves it using the provided file name.
func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, js, 0644)
}

// Get opens the file identified by filename,
// decodes the JSON data and parses it into a List.
func (l *List) Get(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if len(file) == 0 {
		return nil
	}
	return json.Unmarshal(file, l)
}
