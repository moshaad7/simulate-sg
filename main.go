package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/couchbase/cbgt"
	"github.com/couchbase/clog"
)

func createManager(cfg cbgt.Cfg, uuid string) *cbgt.Manager {
	tags := []string{"feed", "janitor", "pindex", "planner"}
	options := make(map[string]string)
	options["managerLoadDataDir"] = "false"

	// server := "http://127.0.0.1:9000"
	// options["nsServerURL"] = server
	server := ""
	return cbgt.NewManagerEx(
		cbgt.VERSION, cfg, uuid, tags, "", 1, "", uuid, "", server, nil, options)
}

var indexType string
func init() {
	indexType = "useless-index"
	cbgt.RegisterPIndexImplType(indexType,
		&cbgt.PIndexImplType{
			New: func(indexType, indexParams, path string, restart func()) (cbgt.PIndexImpl, cbgt.Dest, error) {
				clog.Printf("[APP]: New PIndex created")
				return nil, nil, nil
			},
			Open: func(indexType, path string, restart func()) (cbgt.PIndexImpl, cbgt.Dest, error) {
				return nil, nil, errors.New("open PIndexImpl not supported")
			},
			OpenUsing: func(indexType, path, indexParams string, restart func()) (cbgt.PIndexImpl, cbgt.Dest, error) {
				return nil, nil, errors.New("openUsing PIndexImpl not supported")
			},
			Description: "useless index",
		})
}

func main() {
	loggerFunc := func(level, format string, args ...interface{}) string {
		ts := time.Now().Format("2006-01-02T15:04:05.000-07:00")
		prefix := ts + " [" + level + "] "
		if format != "" {
			return prefix + fmt.Sprintf(format, args...)
		}
		return prefix + fmt.Sprint(args...)
	}
	clog.SetLoggerCallback(loggerFunc)
	
	// Global in-memory cfg, it will be shared by all the managers routines.
	cfg := cbgt.NewCfgMem()

	numManagers := 4
	managers := make(map[int]*cbgt.Manager, numManagers)
	for i := 0; i < numManagers; i++ {
		managers[i] = createManager(cfg, fmt.Sprintf("Node#%d", i))
	}

	// To simulate Sync-Gateway startup behaviour, concurrently spawn multiple
	// routines, which will start the manager and immediately create the index. 
	// (same index on all the managers)
	for i := 0; i < numManagers; i++ {
		go func(i int) {
			err := managers[i].Start("wanted")
			if err != nil {
				clog.Printf("[APP]Node[%d]: Error starting, err: %v", i, err)
				return
			}

			clog.Printf("[APP]Node[%d] Started successfully", i)

			err = managers[i].CreateIndex(
				"nil", // sourcceType
				"",    // sourceName
				"",    // sourceUUID
				"{}",  // sourceParams
				indexType,
				"testIndex", // indexName
				"",          // indexParams
				cbgt.PlanParams{
					IndexPartitions:        16,
					MaxPartitionsPerPIndex: int(1024 / 16),
					NumReplicas:            0,
				},
				"", // prevIndexUUID - empty for new index
			)
			if err != nil {
				clog.Printf("[APP]Node[%v]: Error creating index, err: %v", i, err)
				return
			}
			clog.Printf("[APP]Node[%v]: Create Index successful", i)
			managers[i].Kick("NewIndexesCreated")
		}(i)
	}

	// Block infinitely
	<-make(chan struct{})
}
