package stats

import (
	"github.com/rcrowley/go-librato"
	"log"
	"time"
)

type StatsSink struct {
	accountsCreated  int64
	errorsRendered   int64
	dbxAuthCancelled int64
	usersProcessed   int64
	batchesProcessed int64
	pageRenderErrors int64
	runs             int64
	sink             chan int
	metrics          librato.Metrics
	env              string
}

// The various events we track
const (
	ACCOUNT_CREATE    = 1
	RENDER_ERROR      = 2
	CANCEL_DBX_AUTH   = 3
	USER_PROCESSED    = 4
	BATCH_PROCESSED   = 5
	PAGE_RENDER_ERROR = 6
	RUN_COMPLETE      = 7
	ENV_WEB           = "web"
	ENV_WORKER        = "worker"
)

func NewStatsSink(user string, token string, environment string) *StatsSink {
	// Inital sink
	sink := &StatsSink{
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		make(chan int),
		librato.NewSimpleMetrics(user, token, environment),
		environment,
	}

	go sink.handle()
	go sink.flushTimer()

	return sink
}

// Report an event
func (s *StatsSink) Event(event int) {
	s.sink <- event
}

func (s *StatsSink) handle() {
	for {
		event, ok := <-s.sink

		if !ok {
			return
		}

		switch event {
		case ACCOUNT_CREATE:
			s.accountsCreated += 1
		case RENDER_ERROR:
			s.errorsRendered += 1
		case CANCEL_DBX_AUTH:
			s.dbxAuthCancelled += 1
		case USER_PROCESSED:
			s.usersProcessed += 1
		case BATCH_PROCESSED:
			s.batchesProcessed += 1
		case PAGE_RENDER_ERROR:
			s.pageRenderErrors += 1
		case RUN_COMPLETE:
			s.runs += 1
		}
	}
}

// This runs forever and just sends statistics on the given interval.
func (s *StatsSink) flushTimer() {
	for {
		<-time.After(120 * time.Second)
		log.Println("METRICS: Flushing statistics to Librato")

		if s.env == ENV_WEB {
			// Send the created users
			ac := s.metrics.GetGauge(s.env + "_accounts_created")
			ac <- s.accountsCreated

			// Send the rendered errors
			er := s.metrics.GetGauge(s.env + "_errors_rendered")
			er <- s.errorsRendered

			// end the cancelled dropbox auth
			dr := s.metrics.GetGauge(s.env + "_dropbox_auth_cancelled")
			dr <- s.dbxAuthCancelled

			// Zero out the counters
			s.accountsCreated = 0
			s.errorsRendered = 0
			s.dbxAuthCancelled = 0
		}

		if s.env == ENV_WORKER {
			us := s.metrics.GetGauge(s.env + "_users_processed")
			us <- s.usersProcessed

			ba := s.metrics.GetGauge(s.env + "_batches_processed")
			ba <- s.batchesProcessed

			pa := s.metrics.GetGauge(s.env + "_page_render_errors")
			pa <- s.pageRenderErrors

			ru := s.metrics.GetGauge(s.env + "_runs")
			ru <- s.runs

			// Zero out the counters

			s.usersProcessed = 0
			s.batchesProcessed = 0
			s.pageRenderErrors = 0
			s.runs = 0
		}
	}
}
