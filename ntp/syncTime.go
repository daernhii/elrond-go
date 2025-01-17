package ntp

import (
	"fmt"
	"sync"
	"time"

	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/logger"
	"github.com/beevik/ntp"
)

var log = logger.GetOrCreate("ntp")

// NTPOptions defines configuration options for an NTP query
type NTPOptions struct {
	Hosts        []string
	Version      int
	LocalAddress string
	Timeout      time.Duration
	Port         int
	HostIndex    int
}

// NewNTPGoogleConfig creates an NTPConfig object that configures NTP to use
// time.google.com. This is a useful default for tests, for example, to avoid
// loading a configuration file just to have an NTPConfig.
func NewNTPGoogleConfig() config.NTPConfig {
	return config.NTPConfig{
		Hosts:               []string{"time.google.com", "time.cloudflare.com", "time.apple.com", "time.windows.com"},
		Port:                123,
		Version:             0,
		TimeoutMilliseconds: 100}
}

// NewNTPOptions creates a new NTPOptions object.
func NewNTPOptions(ntpConfig config.NTPConfig, hostIndex int) NTPOptions {
	if ntpConfig.TimeoutMilliseconds <= 0 {
		ntpConfig.TimeoutMilliseconds = 100
	}
	timeout := time.Duration(ntpConfig.TimeoutMilliseconds) * time.Millisecond

	return NTPOptions{
		Hosts:        ntpConfig.Hosts,
		Port:         ntpConfig.Port,
		Version:      ntpConfig.Version,
		LocalAddress: "",
		Timeout:      timeout,
		HostIndex:    hostIndex,
	}
}

// queryNTP wraps beevikntp.QueryWithOptions, in order to use NTPOptions, which
// contains both Host and Port, unlike beevikntp.QueryOptions.
func queryNTP(options NTPOptions) (*ntp.Response, error) {
	queryOptions := ntp.QueryOptions{
		Timeout:      options.Timeout,
		Version:      options.Version,
		LocalAddress: options.LocalAddress,
		Port:         options.Port}

	log.Debug("ntp",
		"request", options.Hosts[options.HostIndex],
		"port", options.Port,
	)

	return ntp.QueryWithOptions(options.Hosts[options.HostIndex], queryOptions)
}

// syncTime defines an object for time synchronization
type syncTime struct {
	mut         sync.RWMutex
	clockOffset time.Duration
	syncPeriod  time.Duration
	ntpOptions  NTPOptions
	query       func(options NTPOptions) (*ntp.Response, error)
}

// NewSyncTime creates a syncTime object. The customQueryFunc argument allows
// the caller to set a different NTP-querying callback, if desired. If set to
// nil, then the default queryNTP is used.
func NewSyncTime(ntpConfig config.NTPConfig, syncPeriod time.Duration, customQueryFunc func(options NTPOptions) (*ntp.Response, error)) *syncTime {
	var queryFunc func(options NTPOptions) (*ntp.Response, error)
	if customQueryFunc == nil {
		queryFunc = queryNTP
	} else {
		queryFunc = customQueryFunc
	}

	hostIndex := checkNTPHost(ntpConfig, queryFunc)

	s := syncTime{
		clockOffset: 0,
		syncPeriod:  syncPeriod,
		query:       queryFunc,
		ntpOptions:  NewNTPOptions(ntpConfig, hostIndex)}

	return &s
}

func checkNTPHost(ntpConfig config.NTPConfig, customQueryFunc func(options NTPOptions) (*ntp.Response, error)) int {
	for hostIndex := range ntpConfig.Hosts {
		_, err := customQueryFunc(NewNTPOptions(ntpConfig, hostIndex))
		if err != nil {
			continue
		}
		log.Debug("using NTP server",
			"host", ntpConfig.Hosts[hostIndex],
		)
		return hostIndex
	}

	return 0
}

// StartSync method does the time synchronization at every syncPeriod time elapsed. This should be started
// as a go routine
func (s *syncTime) StartSync() {
	for {
		s.sync()
		time.Sleep(s.syncPeriod)
	}
}

// sync method does the time synchronization and sets the current offset difference between local time
// and server time with which it has done the synchronization
func (s *syncTime) sync() {
	clockOffsetSum := time.Duration(0)
	succeededRequests := 0

	for i := 0; i < len(s.ntpOptions.Hosts); i++ {
		r, err := s.query(s.ntpOptions)

		if err != nil {
			log.Debug("ntp", "error", err.Error())

			//Change host if the current host returns an error
			newHostIndex := checkNTPHost(config.NTPConfig{
				Hosts: s.ntpOptions.Hosts,
				Port:  s.ntpOptions.Port,
			}, s.query)
			s.ntpOptions.HostIndex = newHostIndex

			continue
		}

		log.Trace("ntp",
			"reading", r.Time.Format("Mon Jan 2 15:04:05 MST 2006"),
		)

		succeededRequests++
		clockOffsetSum += r.ClockOffset
	}

	if succeededRequests > 0 {
		averageClockOffset := time.Duration(int64(clockOffsetSum) / int64(succeededRequests))
		s.setClockOffset(averageClockOffset)
	}
}

// ClockOffset method gets the current time offset
func (s *syncTime) ClockOffset() time.Duration {
	s.mut.RLock()
	clockOffset := s.clockOffset
	s.mut.RUnlock()

	return clockOffset
}

func (s *syncTime) setClockOffset(clockOffset time.Duration) {
	s.mut.Lock()
	s.clockOffset = clockOffset
	s.mut.Unlock()
}

// FormattedCurrentTime method gets the formatted current time on which is added a given offset
func (s *syncTime) FormattedCurrentTime() string {
	return s.formatTime(s.CurrentTime())
}

// formatTime method gets the formatted time from a given time
func (s *syncTime) formatTime(time time.Time) string {
	str := fmt.Sprintf("%.4d-%.2d-%.2d %.2d:%.2d:%.2d.%.9d ", time.Year(), time.Month(), time.Day(), time.Hour(),
		time.Minute(), time.Second(), time.Nanosecond())
	return str
}

// CurrentTime method gets the current time on which is added the current offset
func (s *syncTime) CurrentTime() time.Time {
	s.mut.RLock()
	defer s.mut.RUnlock()
	return time.Now().Add(s.clockOffset)
}

// IsInterfaceNil returns true if there is no value under the interface
func (s *syncTime) IsInterfaceNil() bool {
	if s == nil {
		return true
	}
	return false
}
