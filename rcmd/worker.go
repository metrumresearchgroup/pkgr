package rcmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// InstallRequest provides information about the installation request
type InstallRequest struct {
	Package      string
	Path         string
	InstallArgs  *InstallArgs
	ExecSettings ExecSettings
	RSettings    RSettings
}

// InstallUpdate provides information about the Job in the queue
type InstallUpdate struct {
	Result       CmdResult
	BinaryPath   string
	Msg          string
	Err          error
	ShouldUpdate bool
}

// Worker does work
type Worker struct {
	ID          int
	WorkQueue   <-chan InstallRequest
	UpdateQueue chan<- InstallUpdate
	Quit        chan bool
	InstallFunc func(fs afero.Fs,
		tbp string, // tarball path
		args *InstallArgs,
		rs RSettings,
		es ExecSettings,
		lg *logrus.Logger) (CmdResult, string, error)
}

// InstallQueue represents a new install queue
type InstallQueue struct {
	WorkQueue   chan InstallRequest
	UpdateQueue chan InstallUpdate
	Workers     []Worker
}

// Start starts a worker
func (w *Worker) Start(lg *logrus.Logger) {
	lg.Infof("starting worker %v \n", w.ID)
	go func() {
		appFS := afero.NewOsFs()
		for {
			select {
			case ir := <-w.WorkQueue:
				// Receive a work request.
				res, bPath, err := w.InstallFunc(appFS, ir.Path, ir.InstallArgs, ir.RSettings, ir.ExecSettings, lg)
				w.UpdateQueue <- InstallUpdate{
					Result:       res,
					BinaryPath:   bPath,
					Msg:          "completed job",
					Err:          err,
					ShouldUpdate: true,
				}

			case <-w.Quit:
				// We have been asked to stop.
				lg.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}

//NewInstallQueue provides a new Installation queue with a set number of workers
func NewInstallQueue(n int,
	installFunc func(fs afero.Fs,
		tbp string, // tarball path
		args *InstallArgs,
		rs RSettings,
		es ExecSettings,
		lg *logrus.Logger) (CmdResult, string, error),
	updateFunc func(InstallUpdate), lg *logrus.Logger) *InstallQueue {
	wq := make(chan InstallRequest, 2000)
	uq := make(chan InstallUpdate, 50)
	iq := InstallQueue{
		WorkQueue:   wq,
		UpdateQueue: uq,
	}
	for i := 0; i < n; i++ {
		iq.RegisterNewWorker(i+1, installFunc, lg)
	}
	go iq.HandleUpdates(updateFunc)
	return &iq
}

// HandleUpdates handles updates
func (i *InstallQueue) HandleUpdates(fn func(InstallUpdate)) {
	for {
		iu := <-i.UpdateQueue
		if iu.ShouldUpdate {
			fn(iu)
		}
	}
}

// RegisterNewWorker registers new workers
func (i *InstallQueue) RegisterNewWorker(id int, fn func(fs afero.Fs,
	tbp string, // tarball path
	args *InstallArgs,
	rs RSettings,
	es ExecSettings,
	lg *logrus.Logger) (CmdResult, string, error), lg *logrus.Logger) {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		WorkQueue:   i.WorkQueue,
		UpdateQueue: i.UpdateQueue,
		Quit:        make(chan bool),
		InstallFunc: fn,
	}
	worker.Start(lg)
	i.Workers = append(i.Workers, worker)
	return
}

// Push adds work to the InstallQueue
func (i *InstallQueue) Push(r InstallRequest) {
	i.WorkQueue <- r
}
