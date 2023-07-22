package scheduler

import (
	"context"
	"dagger.io/dagger"
	"github.com/excoriate/stiletto/internal/core/adapters"
	"github.com/excoriate/stiletto/internal/core/entities"
	"github.com/excoriate/stiletto/internal/errors"
	"github.com/excoriate/stiletto/internal/observability"
	"github.com/excoriate/stiletto/internal/utils"
	"go.uber.org/zap"
)

type ScheduledJobs struct {
	Id     string
	Client *entities.Client
	Jobs   []entities.Job

	Ctx    *context.Context
	Logger *zap.Logger

	DaggerClient *dagger.Client
}

type Builder struct {
	id           string
	client       *entities.Client
	jobs         []entities.Job
	daggerClient *dagger.Client

	// Generic inherited objects.
	error  error
	logger *zap.Logger
	ctx    *context.Context
}

type Scheduler interface {
	RunJobs(jobs []entities.Job) error
}

func (r *Builder) Build() (*ScheduledJobs, error) {
	if r.error != nil {
		return nil, r.error
	}

	return &ScheduledJobs{
		Id:           r.id,
		Client:       r.client,
		Jobs:         r.jobs,
		DaggerClient: r.daggerClient,
		Logger:       r.logger,
		Ctx:          r.ctx,
	}, nil
}

func (r *Builder) WithDaggerEngine() *Builder {
	ctx := r.client.Ctx

	daggerEngine, err := adapters.NewDaggerClient(ctx)

	if err != nil {
		schedulerErr := errors.NewConfigurationError(
			"Failed to create a new dagger engine as part of this scheduler", err)
		r.logger.Error(schedulerErr.Error())
		r.error = schedulerErr

		return r
	}

	r.daggerClient = daggerEngine
	return r
}

func (r *Builder) WithJobsToRun(job []entities.Job) *Builder {
	if r.client == nil {
		err := errors.NewConfigurationError("No client instance found. "+
			"The scheduler requires a valid 'Client' instance. "+
			"Ensure you're calling the WithClient() API/Function first!", nil)

		// FIXME: I'm not happy with this. I need to find a better way to handle this scenario.
		logger := observability.NewLogger()
		logger.Error(err.Error())

		r.error = err
		return r
	}

	r.jobs = job
	return r
}

func (r *Builder) WithClient(c *entities.Client) *Builder {
	r.client = c
	r.logger = c.Logger
	r.ctx = c.Ctx
	return r
}

func NewScheduler() *Builder {
	return &Builder{
		id:     utils.GetUUID(),
		client: nil,
		logger: nil,
		error:  nil,
		ctx:    nil,
	}
}
