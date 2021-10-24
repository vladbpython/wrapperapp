package monitoring

import (
	"fmt"
	"strings"
	"time"

	"github.com/vladbpython/wrapperapp/interfaces"
	"github.com/vladbpython/wrapperapp/monitoring/adapters"
)

type Monitoring struct {
	AppName    string
	Events     []string
	Components []string
	Adapters   []interfaces.MonitoringAdapterInterface
}

func (m *Monitoring) AddEvents(events ...string) {
	for i, event := range events {
		m.Events[i] = strings.ToUpper(event)
	}
}

func (m *Monitoring) AddComponents(components ...string) {
	for i, component := range components {
		m.Components[i] = strings.ToLower(component)
	}
}

func (m *Monitoring) AddAdapter(adapter interfaces.MonitoringAdapterInterface) error {
	return adapter.Initializate()
}

func (m *Monitoring) SendData(event, appName, message string, dateTime time.Time, dateTimeLayout string) error {
	var mDateTime string
	if dateTimeLayout != "" {
		mDateTime = dateTime.Format(dateTimeLayout)
	} else {
		mDateTime = dateTime.String()
	}
	Message := fmt.Sprintf("<b>Event</b>: %s\n\r<b>App:</b> %s\n\r<b>Message:</b> %s\n\r<b>DateTime:</b> %s", event, fmt.Sprintf("%s %s", m.AppName, appName), message, mDateTime)

	for _, Event := range m.Events {
		if strings.ToUpper(event) == Event {
			for _, adapter := range m.Adapters {

				err := adapter.SendData(Message)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func NewMonitoringFromConfig(appName string, cfg *ConfigMinotiring) (*Monitoring, error) {
	var NAdapter uint
	monitor := &Monitoring{
		AppName:    appName,
		Events:     make([]string, len(cfg.Events)),
		Components: make([]string, len(cfg.Compontens)),
		Adapters:   make([]interfaces.MonitoringAdapterInterface, len(cfg.Adapters)),
	}
	monitor.AddComponents(cfg.Compontens...)
	monitor.AddEvents(cfg.Events...)

	for _, adapter := range cfg.Adapters {
		adapterName := fmt.Sprintf("monitor %s adapter", adapter.Adapter)
		if strings.ToLower(adapter.Adapter) == "telegram" {
			mAdapter, err := adapters.NewTelegramAdapter(adapterName, &adapter)
			if err != nil {
				return monitor, err
			}
			err = monitor.AddAdapter(mAdapter)
			if err != nil {
				return monitor, err
			}
			monitor.Adapters[NAdapter] = mAdapter
			NAdapter++
		}
	}

	return monitor, nil
}
