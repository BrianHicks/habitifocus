package omnifocus

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

var allTasksScript = `
of = Application("OmniFocus");

function main() {
    var tasks = [];
    of.defaultDocument.flattenedTasks().forEach(function(task) {
		if (task.repetitionRule() !== null) { return; }

	    tasks.push({
		    "id": task.id(),
		    "name": task.name(),
			"done": task.completed()
		});
	});

	return JSON.stringify(tasks);
}

main();
`

// OFTask represents a few fields from an OmniFocus task
type OFTask struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Done bool   `json:"done"`
}

func (t *OFTask) String() string {
	return fmt.Sprintf("%s (done: %t)", t.Name, t.Done)
}

func GetTasks() (map[string]*OFTask, error) {
	command := exec.Command(
		"/usr/bin/osascript",
		"-l", "JavaScript",
		"-e", allTasksScript,
	)
	bytes, err := command.Output()
	if err != nil {
		log.Fatalf("could not read tasks: %s", err)
	}

	var tasks []*OFTask
	err = json.Unmarshal(bytes, &tasks)
	if err != nil {
		log.Fatalf("could not unmarshal: %s", err)
	}

	out := map[string]*OFTask{}
	for _, task := range tasks {
		out[task.ID] = task
	}
	return out, nil
}
