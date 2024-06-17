package workflows

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
)

type stepRequest struct {
	stepRef string
	state   store.WorkflowExecution
}

// Engine handles the lifecycle of a single workflow and its executions.
type Engine struct {
	services.StateMachine
	logger               logger.Logger
	registry             core.CapabilitiesRegistry
	workflow             *workflow
	getLocalNode         func(ctx context.Context) (capabilities.Node, error)
	localNode            capabilities.Node
	executionStates      store.Store
	pendingStepRequests  chan stepRequest
	triggerEvents        chan capabilities.CapabilityResponse
	newWorkerCh          chan struct{}
	stepUpdateCh         chan store.WorkflowExecutionStep
	wg                   sync.WaitGroup
	stopCh               services.StopChan
	newWorkerTimeout     time.Duration
	maxExecutionDuration time.Duration

	// testing lifecycle hook to signal when an execution is finished.
	onExecutionFinished func(string)
	// testing lifecycle hook to signal initialization status
	afterInit func(success bool)
	// Used for testing to control the number of retries
	// we'll do when initializing the engine.
	maxRetries int
	// Used for testing to control the retry interval
	// when initializing the engine.
	retryMs int

	clock clockwork.Clock
}

func (e *Engine) Start(ctx context.Context) error {
	return e.StartOnce("Engine", func() error {
		// create a new context, since the one passed in via Start is short-lived.
		ctx, _ := e.stopCh.NewCtx()

		e.wg.Add(2)
		go e.init(ctx)
		go e.loop(ctx)

		return nil
	})
}

// resolveWorkflowCapabilities does the following:
//
// 1. Resolves the underlying capability for each trigger
// 2. Registers each step's capability to this workflow
func (e *Engine) resolveWorkflowCapabilities(ctx context.Context) error {
	//
	// Step 1. Resolve the underlying capability for each trigger
	//
	triggersInitialized := true
	for _, t := range e.workflow.triggers {
		tg, err := e.registry.GetTrigger(ctx, t.ID)
		if err != nil {
			e.logger.Errorf("capability id: %s failed to get trigger capability: %s", t.ID, err)
			// we don't immediately return here, since we want to retry all triggers
			// to notify the user of all errors at once.
			triggersInitialized = false
		} else {
			t.trigger = tg
		}
	}
	if !triggersInitialized {
		return newWorkflowError(errors.New("failed to resolve triggers"), e.workflow.id, "")
	}

	// Step 2. Walk the graph and register each step's capability to this workflow
	//
	// This means:
	// - fetching the capability
	// - register the capability to this workflow
	// - initializing the step's executionStrategy
	capabilityRegistrationErr := e.workflow.walkDo(workflows.KeywordTrigger, func(s *step) error {
		// The graph contains a dummy step for triggers, but
		// we handle triggers separately since there might be more than one
		// trigger registered to a workflow.
		if s.Ref == workflows.KeywordTrigger {
			return nil
		}

		err := e.initializeCapability(ctx, s)
		if err != nil {
			return newStepError(err, e.workflow.id, s.ID, s.Ref, "failed to initialize capability for step")
		}

		return nil
	})

	return capabilityRegistrationErr
}

func (e *Engine) initializeCapability(ctx context.Context, step *step) error {
	// If the capability already exists, that means we've already registered it
	if step.capability != nil {
		return nil
	}

	cp, err := e.registry.Get(ctx, step.ID)
	if err != nil {
		return newCapabilityError(err, e.workflow.id, step.ID, "failed to get capability")
	}

	info, err := cp.Info(ctx)
	if err != nil {
		return newCapabilityError(err, e.workflow.id, step.ID, "failed to get capability info")
	}

	// Special treatment for local targets - wrap into a transmission capability
	// If the DON is nil, this is a local target.
	if info.CapabilityType == capabilities.CapabilityTypeTarget && info.DON == nil {
		l := e.logger.With("capabilityID", step.ID)
		l.Debugw("wrapping capability in local transmission protocol")
		cp = transmission.NewLocalTargetCapability(
			e.logger,
			step.ID,
			e.localNode,
			cp.(capabilities.TargetCapability),
		)
	}

	// We configure actions, consensus and targets here, and
	// they all satisfy the `CallbackCapability` interface
	cc, ok := cp.(capabilities.CallbackCapability)
	if !ok {
		return newCapabilityError(errors.New("capability does not satisfy CallbackCapability"), e.workflow.id, step.ID, "")
	}

	if step.config == nil {
		configMap, newMapErr := values.NewMap(step.Config)
		if newMapErr != nil {
			return newCapabilityError(newMapErr, e.workflow.id, step.ID, "failed to convert config to values.Map")
		}
		step.config = configMap
	}

	registrationRequest := capabilities.RegisterToWorkflowRequest{
		Metadata: capabilities.RegistrationMetadata{
			WorkflowID: e.workflow.id,
		},
		Config: step.config,
	}

	err = cc.RegisterToWorkflow(ctx, registrationRequest)
	if err != nil {
		return newCapabilityError(
			err, e.workflow.id, step.ID,
			fmt.Sprintf("failed to register capability to workflow (%+v)", registrationRequest))
	}

	step.capability = cc
	return nil
}

// init does the following:
//
//  1. Resolves the LocalDON information
//  2. Resolves the underlying capability for each trigger
//  3. Registers each step's capability to this workflow
//  4. Registers for trigger events now that all capabilities are resolved
//
// Steps 1-3 are retried every 5 seconds until successful.
func (e *Engine) init(ctx context.Context) {
	defer e.wg.Done()

	retryErr := retryable(ctx, e.logger, e.retryMs, e.maxRetries, func() error {
		// first wait for localDON to return a non-error response; this depends
		// on the underlying peerWrapper returning the PeerID.
		node, err := e.getLocalNode(ctx)
		if err != nil {
			return fmt.Errorf("failed to get donInfo: %w", err)
		}
		e.localNode = node

		err = e.resolveWorkflowCapabilities(ctx)
		if err != nil {
			return newWorkflowError(err, e.workflow.id, "failed to resolve workflow capabilities")
		}
		return nil
	})

	if retryErr != nil {
		e.logger.Errorf("initialization failed: %s", retryErr)
		e.afterInit(false)
		return
	}

	e.logger.Debug("capabilities resolved, resuming in-progress workflows")
	err := e.resumeInProgressExecutions(ctx)
	if err != nil {
		e.logger.Errorf("failed to resume in-progress workflows: %v", err)
	}

	e.logger.Debug("registering triggers")
	for idx, t := range e.workflow.triggers {
		err := e.registerTrigger(ctx, t, idx)
		if err != nil {
			e.logger.Errorf("capability id: %s failed to register trigger: %s", t.ID, err)
		}
	}

	e.logger.Info("engine initialized")
	e.afterInit(true)
}

var (
	defaultOffset, defaultLimit = 0, 1_000
)

func (e *Engine) resumeInProgressExecutions(ctx context.Context) error {
	wipExecutions, err := e.executionStates.GetUnfinished(ctx, defaultOffset, defaultLimit)
	if err != nil {
		return err
	}

	// TODO: paginate properly
	if len(wipExecutions) >= defaultLimit {
		e.logger.Warnf("possible execution overflow during resumption, work in progress executions: %d >= %d", len(wipExecutions), defaultLimit)
	}

	// Cache the dependents associated with a step.
	// We may have to reprocess many executions, but should only
	// need to calculate the dependents of a step once since
	// they won't change.
	refToDeps := map[string][]*step{}
	for _, execution := range wipExecutions {
		for _, step := range execution.Steps {
			// NOTE: In order to determine what tasks need to be enqueued,
			// we look at any completed steps, and for each dependent,
			// check if they are ready to be enqueued.
			// This will also handle an execution that has stalled immediately on creation,
			// since we always create an execution with an initially completed trigger step.
			if step.Status != store.StatusCompleted {
				continue
			}

			sds, ok := refToDeps[step.Ref]
			if !ok {
				s, err := e.workflow.dependents(step.Ref)
				if err != nil {
					return err
				}

				sds = s
			}

			for _, sd := range sds {
				e.queueIfReady(execution, sd)
			}
		}
	}
	return nil
}

func generateTriggerId(workflowID string, triggerIdx int) string {
	return fmt.Sprintf("wf_%s_trigger_%d", workflowID, triggerIdx)
}

// registerTrigger is used during the initialization phase to bind a trigger to this workflow
func (e *Engine) registerTrigger(ctx context.Context, t *triggerCapability, triggerIdx int) error {
	triggerID := generateTriggerId(e.workflow.id, triggerIdx)
	triggerInputs, err := values.NewMap(
		map[string]any{
			"triggerId": triggerID, 
		},
	)
	if err != nil {
		return err
	}

	tc, err := values.NewMap(t.Config)
	if err != nil {
		return err
	}

	t.config = tc

	triggerRegRequest := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID:    e.workflow.id,
			WorkflowDonID: e.localNode.WorkflowDON.ID,
			WorkflowName:  e.workflow.name,
			WorkflowOwner: e.workflow.owner,
		},
		Config: tc,
		Inputs: triggerInputs,
	}
	eventsCh, err := t.trigger.RegisterTrigger(ctx, triggerRegRequest)
	if err != nil {
		// It's confusing that t.ID is different from triggerID, but
		// t.ID is the capability ID, and triggerID is the trigger ID.
		//
		// The capability ID is globally scoped, whereas the trigger ID
		// is scoped to this workflow.
		//
		// For example, t.ID might be "streams-trigger:network=mainnet@1.0.0"
		// and triggerID might be "wf_123_trigger_0"
		return newTriggerError(err, e.workflow.id, t.ID, triggerID, "failed to register trigger")
	}

	go func() {
		for event := range eventsCh {
			e.triggerEvents <- event
		}
	}()

	return nil
}

// loop is the synchronization goroutine for the engine, and is responsible for:
//   - dispatching new workers up to the limit specified (default = 100)
//   - starting a new execution when a trigger emits a message on `triggerEvents`
//   - updating the `executionState` with the outcome of a `step`.
//
// Note: `executionState` is only mutated by this loop directly.
//
// This is important to avoid data races, and any accesses of `executionState` by any other
// goroutine should happen via a `stepRequest` message containing a copy of the latest
// `executionState`.
//
// This works because a worker thread for a given step will only
// be spun up once all dependent steps have completed (guaranteeing that the state associated
// with those dependent steps will no longer change). Therefore as long this worker thread only
// accesses data from dependent states, the data will never be stale.
func (e *Engine) loop(ctx context.Context) {
	defer e.wg.Done()
	for {
		select {
		case <-ctx.Done():
			e.logger.Debugw("shutting down loop")
			return
		case resp, isOpen := <-e.triggerEvents:
			if !isOpen {
				e.logger.Errorf("trigger events channel is no longer open, skipping")
				continue
			}

			if resp.Err != nil {
				e.logger.Errorf("trigger event was an error %v; not executing", resp.Err)
				continue
			}

			te := &capabilities.TriggerEvent{}
			err := resp.Value.UnwrapTo(te)
			if err != nil {
				e.logger.Errorf("could not unwrap trigger event; error %v", resp.Err)
				continue
			}

			executionID, err := generateExecutionID(e.workflow.id, te.ID)
			if err != nil {
				e.logger.With("triggerid", te.ID).Errorf("could not generate execution ID: %v", err)
				continue
			}

			err = e.startExecution(ctx, executionID, resp.Value)
			if err != nil {
				e.logger.With("executionID", executionID).Errorf("failed to start execution: %v", err)
			}
		case pendingStepRequest := <-e.pendingStepRequests:
			// Wait for a new worker to be available before dispatching a new one.
			// We'll do this up to newWorkerTimeout. If this expires, we'll put the
			// message back on the queue and keep going.
			t := e.clock.NewTimer(e.newWorkerTimeout)
			select {
			case <-e.newWorkerCh:
				e.wg.Add(1)
				go e.workerForStepRequest(ctx, pendingStepRequest)
			case <-t.Chan():
				e.logger.With("executionID", pendingStepRequest.state.ExecutionID, "stepRef", pendingStepRequest.stepRef).
					Errorf("timed out when spinning off worker for pending step request %+v", pendingStepRequest)
				e.pendingStepRequests <- pendingStepRequest
			}
			t.Stop()
		case stepUpdate := <-e.stepUpdateCh:
			// Executed synchronously to ensure we correctly schedule subsequent tasks.
			err := e.handleStepUpdate(ctx, stepUpdate)
			if err != nil {
				e.logger.With("executionID", stepUpdate.ExecutionID, "stepRef", stepUpdate.Ref).
					Errorf("failed to update step state: %+v, %s", stepUpdate, err)
			}
		}
	}
}

func generateExecutionID(workflowID, eventID string) (string, error) {
	s := sha256.New()
	_, err := s.Write([]byte(workflowID))
	if err != nil {
		return "", err
	}

	_, err = s.Write([]byte(eventID))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(s.Sum(nil)), nil
}

// startExecution kicks off a new workflow execution when a trigger event is received.
func (e *Engine) startExecution(ctx context.Context, executionID string, event values.Value) error {
	e.logger.Debugw("executing on a trigger event", "event", event, "executionID", executionID)
	ec := &store.WorkflowExecution{
		Steps: map[string]*store.WorkflowExecutionStep{
			workflows.KeywordTrigger: {
				Outputs: store.StepOutput{
					Value: event,
				},
				Status:      store.StatusCompleted,
				ExecutionID: executionID,
				Ref:         workflows.KeywordTrigger,
			},
		},
		WorkflowID:  e.workflow.id,
		ExecutionID: executionID,
		Status:      store.StatusStarted,
	}

	err := e.executionStates.Add(ctx, ec)
	if err != nil {
		return err
	}

	// Find the tasks we need to fire when a trigger has fired and enqueue them.
	// This consists of a) nodes without a dependency and b) nodes which depend
	// on a trigger
	triggerDependents, err := e.workflow.dependents(workflows.KeywordTrigger)
	if err != nil {
		return err
	}

	for _, td := range triggerDependents {
		e.queueIfReady(*ec, td)
	}

	return nil
}

func (e *Engine) handleStepUpdate(ctx context.Context, stepUpdate store.WorkflowExecutionStep) error {
	state, err := e.executionStates.UpsertStep(ctx, &stepUpdate)
	if err != nil {
		return err
	}
	l := e.logger.With("executionID", state.ExecutionID, "stepRef", stepUpdate.Ref)

	switch stepUpdate.Status {
	case store.StatusCompleted:
		stepDependents, err := e.workflow.dependents(stepUpdate.Ref)
		if err != nil {
			return err
		}

		// There are no steps left to process in the current path, so let's check if
		// we've completed the workflow.
		if len(stepDependents) == 0 {
			workflowCompleted := true
			err := e.workflow.walkDo(workflows.KeywordTrigger, func(s *step) error {
				step, ok := state.Steps[s.Ref]
				// The step is missing from the state,
				// which means it hasn't been processed yet.
				// Let's mark `workflowCompleted` = false, and
				// continue.
				if !ok {
					workflowCompleted = false
					return nil
				}

				switch step.Status {
				case store.StatusCompleted, store.StatusErrored, store.StatusCompletedEarlyExit:
				default:
					workflowCompleted = false
				}
				return nil
			})
			if err != nil {
				return err
			}

			if workflowCompleted {
				return e.finishExecution(ctx, state.ExecutionID, store.StatusCompleted)
			}
		}

		// We haven't completed the workflow, but should we continue?
		// If we've been executing for too long, let's time the workflow out and stop here.
		if state.CreatedAt != nil && e.clock.Since(*state.CreatedAt) > e.maxExecutionDuration {
			l.Infow("execution timed out")
			return e.finishExecution(ctx, state.ExecutionID, store.StatusTimeout)
		}

		// Finally, since the workflow hasn't timed out or completed, let's
		// check for any dependents that are ready to process.
		for _, sd := range stepDependents {
			e.queueIfReady(state, sd)
		}
	case store.StatusCompletedEarlyExit:
		l.Infow("execution terminated early")
		// NOTE: even though this marks the workflow as completed, any branches of the DAG
		// that don't depend on the step that signaled for an early exit will still complete.
		// This is to ensure that any side effects are executed consistently, since otherwise
		// the async nature of the workflow engine would provide no guarantees.
		err := e.finishExecution(ctx, state.ExecutionID, store.StatusCompletedEarlyExit)
		if err != nil {
			return err
		}
	case store.StatusErrored:
		l.Infow("execution errored")
		err := e.finishExecution(ctx, state.ExecutionID, store.StatusErrored)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Engine) queueIfReady(state store.WorkflowExecution, step *step) {
	// Check if all dependencies are completed for the current step
	var waitingOnDependencies bool
	for _, dr := range step.Vertex.Dependencies {
		stepState, ok := state.Steps[dr]
		if !ok {
			waitingOnDependencies = true
			continue
		}

		// Unless the dependency is complete,
		// we'll mark waitingOnDependencies = true.
		// This includes cases where one of the dependent
		// steps has errored, since that means we shouldn't
		// schedule the step for execution.
		if stepState.Status != store.StatusCompleted {
			waitingOnDependencies = true
		}
	}

	// If all dependencies are completed, enqueue the step.
	if !waitingOnDependencies {
		e.logger.Debugw("step request enqueued", "stepRef", step.Ref, "executionID", state.ExecutionID, "state", copyState(state))
		e.pendingStepRequests <- stepRequest{
			state:   copyState(state),
			stepRef: step.Ref,
		}
	}
}

func (e *Engine) finishExecution(ctx context.Context, executionID string, status string) error {
	e.logger.Infow("finishing execution", "executionID", executionID, "status", status)
	err := e.executionStates.UpdateStatus(ctx, executionID, status)
	if err != nil {
		return err
	}

	e.onExecutionFinished(executionID)
	return nil
}

func (e *Engine) workerForStepRequest(ctx context.Context, msg stepRequest) {
	defer func() { e.newWorkerCh <- struct{}{} }()
	defer e.wg.Done()

	// Instantiate a child logger; in addition to the WorkflowID field the workflow
	// logger will already have, this adds the `stepRef` and `executionID`
	l := e.logger.With("stepRef", msg.stepRef, "executionID", msg.state.ExecutionID)

	l.Debugw("executing on a step event")
	stepState := &store.WorkflowExecutionStep{
		Outputs:     store.StepOutput{},
		ExecutionID: msg.state.ExecutionID,
		Ref:         msg.stepRef,
	}

	inputs, outputs, err := e.executeStep(ctx, msg)
	var stepStatus string
	switch {
	case errors.Is(capabilities.ErrStopExecution, err):
		l.Infow("step executed successfully with a termination")
		stepStatus = store.StatusCompletedEarlyExit
	case err != nil:
		l.Errorf("error executing step request: %s", err)
		stepStatus = store.StatusErrored
	default:
		l.Infow("step executed successfully", "outputs", outputs)
		stepStatus = store.StatusCompleted
	}

	stepState.Status = stepStatus
	stepState.Outputs.Value = outputs
	stepState.Outputs.Err = err
	stepState.Inputs = inputs

	// Let's try and emit the stepUpdate.
	// If the context is canceled, we'll just drop the update.
	// This means the engine is shutting down and the
	// receiving loop may not pick up any messages we emit.
	// Note: When full persistence support is added, any hanging steps
	// like this one will get picked up again and will be reprocessed.
	select {
	case <-ctx.Done():
		l.Errorf("context canceled before step update could be issued; error %v", err)
	case e.stepUpdateCh <- *stepState:
	}
}

// executeStep executes the referenced capability within a step and returns the result.
func (e *Engine) executeStep(ctx context.Context, msg stepRequest) (*values.Map, values.Value, error) {
	step, err := e.workflow.Vertex(msg.stepRef)
	if err != nil {
		return nil, nil, err
	}

	i, err := findAndInterpolateAllKeys(step.Inputs, msg.state)
	if err != nil {
		return nil, nil, err
	}

	inputs, err := values.NewMap(i.(map[string]any))
	if err != nil {
		return nil, nil, err
	}

	tr := capabilities.CapabilityRequest{
		Inputs: inputs,
		Config: step.config,
		Metadata: capabilities.RequestMetadata{
			WorkflowID:          msg.state.WorkflowID,
			WorkflowExecutionID: msg.state.ExecutionID,
			WorkflowOwner:       e.workflow.owner,
			WorkflowName:        e.workflow.name,
			WorkflowDonID:       e.localNode.WorkflowDON.ID,
		},
	}

	output, err := executeSyncAndUnwrapSingleValue(ctx, step.capability, tr)
	if err != nil {
		return inputs, nil, err
	}

	return inputs, output, err
}

func (e *Engine) deregisterTrigger(ctx context.Context, t *triggerCapability, triggerIdx int) error {
	triggerInputs, err := values.NewMap(
		map[string]any{
			"triggerId": generateTriggerId(e.workflow.id, triggerIdx),
		},
	)
	if err != nil {
		return err
	}
	deregRequest := capabilities.CapabilityRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID:    e.workflow.id,
			WorkflowDonID: e.localNode.WorkflowDON.ID,
			WorkflowName:  e.workflow.name,
			WorkflowOwner: e.workflow.owner,
		},
		Inputs: triggerInputs,
		Config: t.config,
	}

	// if t.trigger == nil, then we haven't initialized the workflow
	// yet, and can safely consider the trigger deregistered with
	// no further action.
	if t.trigger != nil {
		return t.trigger.UnregisterTrigger(ctx, deregRequest)
	}

	return nil
}

func (e *Engine) Close() error {
	return e.StopOnce("Engine", func() error {
		e.logger.Info("shutting down engine")
		ctx := context.Background()
		// To shut down the engine, we'll start by deregistering
		// any triggers to ensure no new executions are triggered,
		// then we'll close down any background goroutines,
		// and finally, we'll deregister any workflow steps.
		for idx, t := range e.workflow.triggers {
			err := e.deregisterTrigger(ctx, t, idx)
			if err != nil {
				return err
			}
		}

		close(e.stopCh)
		e.wg.Wait()

		err := e.workflow.walkDo(workflows.KeywordTrigger, func(s *step) error {
			if s.Ref == workflows.KeywordTrigger {
				return nil
			}

			reg := capabilities.UnregisterFromWorkflowRequest{
				Metadata: capabilities.RegistrationMetadata{
					WorkflowID: e.workflow.id,
				},
				Config: s.config,
			}

			// if capability is nil, then we haven't initialized
			// the workflow yet and can safely consider it deregistered
			// with no further action.
			if s.capability == nil {
				return nil
			}

			innerErr := s.capability.UnregisterFromWorkflow(ctx, reg)
			if innerErr != nil {
				return newStepError(innerErr, e.workflow.id, s.ID, s.Ref, fmt.Sprintf("failed to unregister capability from workflow: %+v", reg))
			}

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})
}

type Config struct {
	Spec                 string
	WorkflowID           string
	WorkflowOwner        string
	WorkflowName         string
	Lggr                 logger.Logger
	Registry             core.CapabilitiesRegistry
	MaxWorkerLimit       int
	QueueSize            int
	NewWorkerTimeout     time.Duration
	MaxExecutionDuration time.Duration
	GetLocalNode         func(ctx context.Context) (capabilities.Node, error)
	Store                store.Store

	// For testing purposes only
	maxRetries          int
	retryMs             int
	afterInit           func(success bool)
	onExecutionFinished func(weid string)
	clock               clockwork.Clock
}

const (
	defaultWorkerLimit          = 100
	defaultQueueSize            = 100000
	defaultNewWorkerTimeout     = 2 * time.Second
	defaultMaxExecutionDuration = 10 * time.Minute
)

func NewEngine(cfg Config) (engine *Engine, err error) {
	if cfg.Store == nil {
		return nil, newWorkflowError(errors.New("store is nil"), cfg.WorkflowID, "")
	}

	if cfg.MaxWorkerLimit == 0 {
		cfg.MaxWorkerLimit = defaultWorkerLimit
	}

	if cfg.QueueSize == 0 {
		cfg.QueueSize = defaultQueueSize
	}

	if cfg.NewWorkerTimeout == 0 {
		cfg.NewWorkerTimeout = defaultNewWorkerTimeout
	}

	if cfg.MaxExecutionDuration == 0 {
		cfg.MaxExecutionDuration = defaultMaxExecutionDuration
	}

	if cfg.GetLocalNode == nil {
		cfg.GetLocalNode = func(ctx context.Context) (capabilities.Node, error) {
			return capabilities.Node{}, nil
		}
	}

	if cfg.retryMs == 0 {
		cfg.retryMs = 5000
	}

	if cfg.afterInit == nil {
		cfg.afterInit = func(success bool) {}
	}

	if cfg.onExecutionFinished == nil {
		cfg.onExecutionFinished = func(weid string) {}
	}

	if cfg.clock == nil {
		cfg.clock = clockwork.NewRealClock()
	}

	// TODO: validation of the workflow spec
	// We'll need to check, among other things:
	// - that there are no step `ref` called `trigger` as this is reserved for any triggers
	// - that there are no duplicate `ref`s
	// - that the `ref` for any triggers is empty -- and filled in with `trigger`
	// - that the resulting graph is strongly connected (i.e. no disjointed subgraphs exist)
	// - etc.

	workflow, err := Parse(cfg.Spec)
	if err != nil {
		return nil, err
	}

	workflow.id = cfg.WorkflowID
	workflow.owner = cfg.WorkflowOwner
	workflow.name = hex.EncodeToString([]byte(cfg.WorkflowName))

	// Instantiate semaphore to put a limit on the number of workers
	newWorkerCh := make(chan struct{}, cfg.MaxWorkerLimit)
	for i := 0; i < cfg.MaxWorkerLimit; i++ {
		newWorkerCh <- struct{}{}
	}

	engine = &Engine{
		logger:               cfg.Lggr.Named("WorkflowEngine").With("workflowID", cfg.WorkflowID),
		registry:             cfg.Registry,
		workflow:             workflow,
		getLocalNode:         cfg.GetLocalNode,
		executionStates:      cfg.Store,
		pendingStepRequests:  make(chan stepRequest, cfg.QueueSize),
		newWorkerCh:          newWorkerCh,
		stepUpdateCh:         make(chan store.WorkflowExecutionStep),
		triggerEvents:        make(chan capabilities.CapabilityResponse),
		stopCh:               make(chan struct{}),
		newWorkerTimeout:     cfg.NewWorkerTimeout,
		maxExecutionDuration: cfg.MaxExecutionDuration,

		onExecutionFinished: cfg.onExecutionFinished,
		afterInit:           cfg.afterInit,
		maxRetries:          cfg.maxRetries,
		retryMs:             cfg.retryMs,
		clock:               cfg.clock,
	}
	return engine, nil
}

// ExecuteSyncAndUnwrapSingleValue is a convenience method that executes a capability synchronously and unwraps the
// result if it is a single value otherwise returns the list.
func executeSyncAndUnwrapSingleValue(ctx context.Context, cap capabilities.CallbackCapability, req capabilities.CapabilityRequest) (values.Value, error) {
	l, err := capabilities.ExecuteSync(ctx, cap, req)
	if err != nil {
		return nil, err
	}

	// `ExecuteSync` returns a `values.List` even if there was
	// just one return value. If that is the case, let's unwrap the
	// single value to make it easier to use in -- for example -- variable interpolation.
	if len(l.Underlying) > 1 {
		return l, nil
	}

	return l.Underlying[0], nil
}

// We don't define Unwrap since we don't want to expose the underlying error.
type WorkflowError struct {
	Err        error
	Reason     string
	WorkflowID string
}

func (e *WorkflowError) Error() string {
	if e.Reason == "" {
		return fmt.Sprintf("workflow id: %s %v", e.WorkflowID, e.Err)
	}

	return fmt.Sprintf("workflow id: %s %s: %v", e.WorkflowID, e.Reason, e.Err)
}

func newWorkflowError(err error, workflowID, reason string) *WorkflowError {
	return &WorkflowError{Err: err, WorkflowID: workflowID, Reason: reason}
}

type WorkflowExecutionError struct {
	*WorkflowError
	ExecutionID string
}

func (e *WorkflowExecutionError) Error() string {
	return fmt.Sprintf("execution id: %s %v", e.ExecutionID, e.WorkflowError)
}

func (e *WorkflowExecutionError) Unwrap() error {
	return e.WorkflowError
}

type CapabilityError struct {
	*WorkflowError
	CapabilityID string
}

func newCapabilityError(err error, workflowID, capabilityID, reason string) *CapabilityError {
	return &CapabilityError{
		WorkflowError: newWorkflowError(err, workflowID, reason),
		CapabilityID:  capabilityID,
	}
}

func (e *CapabilityError) Error() string {
	return fmt.Sprintf("capability id: %s %v", e.CapabilityID, e.WorkflowError)
}

func (e *CapabilityError) Unwrap() error {
	return e.WorkflowError
}

type TriggerError struct {
	*CapabilityError
	TriggerID string
}

func newTriggerError(err error, workflowID, capabilityID, triggerID, reason string) *TriggerError {
	return &TriggerError{
		CapabilityError: newCapabilityError(err, workflowID, capabilityID, reason),
		TriggerID:       triggerID,
	}
}

func (e *TriggerError) Error() string {
	return fmt.Sprintf("trigger id: %s %v", e.TriggerID, e.CapabilityError)
}

func (e *TriggerError) Unwrap() error {
	return e.CapabilityError
}

func newWorkflowExecutionError(err error, workflowID, executionID, reason string) *WorkflowExecutionError {
	return &WorkflowExecutionError{WorkflowError: newWorkflowError(err, workflowID, reason), ExecutionID: executionID}
}

type StepError struct {
	*WorkflowError
	CapabilityID string
	StepRef      string
}

func newStepError(err error, workflowID, capabilityID, stepRef, reason string) *StepError {
	return &StepError{
		WorkflowError: newWorkflowError(err, workflowID, reason),
		CapabilityID:  capabilityID,
		StepRef:       stepRef,
	}
}

func (e *StepError) Error() string {
	return fmt.Sprintf("step ref: %s capability id: %s %v", e.StepRef, e.CapabilityID, e.WorkflowError)
}

func (e *StepError) Unwrap() error {
	return e.WorkflowError
}

type StepExecutionError struct {
	*WorkflowExecutionError
	CapabilityID string
	StepRef      string
}

func newStepExecutionError(err error, workflowID, executionID, capabilityID, stepRef, reason string) *StepExecutionError {
	return &StepExecutionError{
		WorkflowExecutionError: newWorkflowExecutionError(err, workflowID, executionID, reason),
		CapabilityID:           capabilityID,
		StepRef:                stepRef,
	}
}

func (e *StepExecutionError) Error() string {
	return fmt.Sprintf("step ref: %s capability id: %s %v", e.StepRef, e.CapabilityID, e.WorkflowExecutionError)
}

// Unwrap method to allow errors.Is and errors.As to work
func (e *StepExecutionError) Unwrap() error {
	return e.WorkflowExecutionError
}
