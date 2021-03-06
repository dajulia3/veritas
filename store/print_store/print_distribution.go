package print_store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/cloudfoundry-incubator/runtime-schema/models"

	"github.com/pivotal-cf-experimental/veritas/say"
	"github.com/pivotal-cf-experimental/veritas/veritas_models"
)

func PrintDistribution(tasks bool, lrps bool, clear bool, f io.Reader) error {
	decoder := json.NewDecoder(f)
	var dump veritas_models.StoreDump
	err := decoder.Decode(&dump)
	if err != nil {
		return err
	}

	printDistribution(dump, tasks, lrps, clear)

	return nil
}

func printDistribution(dump veritas_models.StoreDump, includeTasks bool, includeLRPS bool, clear bool) {
	cellIDs := []string{}
	for _, cells := range dump.Services.Cells {
		cellIDs = append(cellIDs, cells.CellID)
	}

	sort.Strings(cellIDs)

	nTasks := map[string]int{}
	nLRPsStarting := map[string]int{}
	nLRPsRunning := map[string]int{}

	for _, tasks := range dump.Tasks {
		for _, task := range tasks {
			nTasks[task.CellID]++
		}
	}

	for _, lrp := range dump.LRPS {
		for _, actuals := range lrp.ActualLRPsByIndex {
			for _, actual := range actuals {
				if actual.State == models.ActualLRPStateStarting {
					nLRPsStarting[actual.CellID]++
				} else {
					nLRPsRunning[actual.CellID]++
				}
			}
		}
	}

	buffer := &bytes.Buffer{}
	if clear {
		say.Fclear(buffer)
	}
	say.Fprintln(buffer, 0, "Distribution")
	for _, cellID := range cellIDs {
		numTasks := nTasks[cellID]
		numLRPs := nLRPsStarting[cellID] + nLRPsRunning[cellID]
		var content string
		if numTasks == 0 && numLRPs == 0 {
			content = say.Red("Empty")
		} else {
			content = fmt.Sprintf("%s%s%s", say.Yellow(strings.Repeat("•", nTasks[cellID])), say.Green(strings.Repeat("•", nLRPsRunning[cellID])), say.Gray(strings.Repeat("•", nLRPsStarting[cellID])))
		}
		say.Fprintln(buffer, 0, "%12s: %s", cellID, content)
	}

	buffer.WriteTo(os.Stdout)
}
