package models

import "github.com/cloudfoundry-incubator/runtime-schema/models"

type VeritasServices struct {
	Executors   []models.ExecutorPresence
	FileServers []string
}
