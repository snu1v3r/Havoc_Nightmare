package db

import (
	"Havoc/pkg/db"
	"Havoc/pkg/logger"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
)

func (d *DatabaseCli) ListenerList() error {
	var (
		list      = table.NewWriter()
		listeners []db.Listener
		err       error
	)

	if listeners, err = d.database.ListenerList(); err != nil {
		return err
	}

	if len(listeners) == 0 {
		logger.Warn("no listeners has been registered")
		return err
	}

	list.SetOutputMirror(os.Stdout)
	list.AppendHeader(table.Row{"name", "protocol", "status"})
	list.SetStyle(table.StyleColoredDark)

	for _, listener := range listeners {
		list.AppendRow(table.Row{listener.Name, listener.Protocol, listener.Status})
	}

	list.AppendFooter(table.Row{fmt.Sprintf("Listeners: %v", len(listeners)), "", ""})
	list.Render()

	return nil
}

func (d *DatabaseCli) ListenerRemove(uuid string) error {
	var (
		err error
	)

	if err = d.database.ListenerRemove(uuid); err != nil {
		logger.Error("failed to remove listener: %v", err)
		return err
	}

	logger.Info("listener %v successfully removed", uuid)

	return nil
}

func (d *DatabaseCli) ListenerClear() error {
	var (
		listeners []db.Listener
		err       error
		removed   int
	)

	if listeners, err = d.database.ListenerList(); err != nil {
		return err
	}

	if len(listeners) == 0 {
		logger.Warn("no listeners available or started")
		return err
	}

	for _, listener := range listeners {
		if err = d.database.ListenerRemove(listener.Name); err != nil {
			logger.Error("failed to remove listener: %v", err)
			continue
		}

		removed++
	}

	logger.Info("cleared %v listeners from the database", removed)

	return nil
}
