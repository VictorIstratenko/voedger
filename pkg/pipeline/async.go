/*
 * Copyright (c) 2021-present unTill Pro, Ltd.
 */

package pipeline

import "time"

func puller_async(wo *WiredOperator) {
	flushTimer := newFlushTimer(wo.FlushInterval)
	var open = true
	var work interface{}
	for open {
		select {
		case work, open = <-wo.Stdin:

			if !open {
				continue
			}

			workpiece := work.(IWorkpiece)

			if !wo.isActive() {
				p_release(workpiece)
				continue
			}

			if wo.forwardIfErrorAsync(workpiece) {
				continue
			}

			outWork, err := wo.doAsync(workpiece)

			if err != nil {
				wo.Stdout <- err
			} else {
				if outWork != nil {
					wo.Stdout <- outWork
				}
				flushTimer.reset()
			}
		case <-flushTimer.timer.C:
			flushTimer.ticked()
			p_flush(wo, placeFlushByTimer)
		}
	}

	p_flush(wo, placeFlushDisassembling)
	wo.Operator.Close()
	close(wo.Stdout)
	flushTimer.stop()
}

func p_flush(wo *WiredOperator, place string) {
	if !wo.isActive() {
		return
	}

	flushProc := func(work IWorkpiece) {
		if wo.isActive() {
			wo.Stdout <- work
		}
	}

	if err := wo.Operator.(IAsyncOperator).Flush(flushProc); err != nil {
		if wo.isActive() {
			wo.Stdout <- wo.NewError(err, nil, place)
		}
	}
}

func p_release(w IWorkpiece) {
	if w != nil {
		w.Release()
	}
}

type flushTimer struct {
	timer  *time.Timer
	intvl  time.Duration
	active bool
}

func newFlushTimer(interval time.Duration) *flushTimer {
	flush := flushTimer{
		intvl:  interval,
		active: true,
		timer:  time.NewTimer(interval),
	}
	flush.stop()
	return &flush
}

func (t *flushTimer) stop() {
	if t.active {
		if !t.timer.Stop() {
			<-t.timer.C
		}
		t.active = false
	}
}

func (t *flushTimer) reset() {
	if !t.active && t.intvl > 0 {
		t.timer.Reset(t.intvl)
		t.active = true
	}
}

func (t *flushTimer) ticked() {
	t.active = false
}
