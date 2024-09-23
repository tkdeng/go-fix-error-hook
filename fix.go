package fix

import "errors"

type errCB struct {
	cb       func(err error) bool
	matchAny bool
}

var errHandler = map[error][]errCB{}

// Hook a callback to be called as a potential fix to an error.
//
// @cb: return true if the fix was successful,
// return false, if it failed.
func Hook(err error, cb func(err error) bool) {
	if _, ok := errHandler[err]; !ok {
		errHandler[err] = []errCB{}
	}
	errHandler[err] = append(errHandler[err], errCB{cb: cb})
}

// HookAny is like Hook, but will use the error.Is method.
//
// @cb: return true if the fix was successful,
// return false, if it failed.
func HookAny(err error, cb func(err error) bool) {
	if _, ok := errHandler[err]; !ok {
		errHandler[err] = []errCB{}
	}
	errHandler[err] = append(errHandler[err], errCB{cb: cb, matchAny: true})
}

// Try to fix an error if a hook can fix it.
//
// @err will be updated with the new error, or with nil if fixed successfully.
//
// note: if @err is already nil, this method will be ignored.
func Try(err *error, retry func(err error) error) {
	if err == nil || *err == nil {
		return
	}

	// to prevent recursion
	errTried := []error{}

	for *err != nil {
		// check for recursion
		for _, e := range errTried {
			if *err == e {
				return
			}
		}
		errTried = append(errTried, *err)

		// try error handlers
		if handlers, ok := errHandler[*err]; ok {
			for _, handler := range handlers {
				if handler.cb(*err) { // if the handler did something (and its worth retrying)
					e := retry(*err)
					if e == nil || e != *err {
						// if we get a different error, stop running the current handlers.
						// we need to run different error handlers now.
						*err = e
						break
					}
				}
			}
		}

		// try with errors.Is
		for e, handlers := range errHandler {
			if errors.Is(*err, e) {
				for _, handler := range handlers {
					if handler.matchAny && handler.cb(*err) { // if the handler did something (and its worth retrying)
						e := retry(*err)
						if e == nil || e != *err {
							// if we get a different error, stop running the current handlers.
							// we need to run different error handlers now.
							*err = e
							break
						}
					}
				}
			}
		}
	}
}

// TryOnce will try to fix only one error.
//
// If a different error comes up, this method will stop trying,
// instead of trying to fix the next error.
//
// @err will be updated with the new error, or with nil if fixed successfully.
//
// note: if @err is already nil, this method will be ignored.
func TryOnce(err *error, retry func(err error) error) {
	if err == nil || *err == nil {
		return
	}

	// try error handlers
	if handlers, ok := errHandler[*err]; ok {
		for _, handler := range handlers {
			if handler.cb(*err) { // if the handler did something (and its worth retrying)
				e := retry(*err)
				if e == nil || e != *err {
					*err = e
					return
				}
			}
		}
	}

	// try with errors.Is
	for e, handlers := range errHandler {
		if errors.Is(*err, e) {
			for _, handler := range handlers {
				if handler.matchAny && handler.cb(*err) { // if the handler did something (and its worth retrying)
					e := retry(*err)
					if e == nil || e != *err {
						// if we get a different error, stop running the current handlers.
						// we need to run different error handlers now.
						*err = e
						break
					}
				}
			}
		}
	}
}
