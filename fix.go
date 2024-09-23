package fix

var errHandler = map[error][]func() bool{}

// Hook a callback to be called as a potential fix to an error.
//
// @cb: return true if the fix was successful,
// return false, if it failed.
func Hook(err error, cb func() bool) {
	if _, ok := errHandler[err]; !ok {
		errHandler[err] = []func() bool{}
	}
	errHandler[err] = append(errHandler[err], cb)
}

// Try to fix an error if a hook can fix it.
//
// @err will be updated with the new error, or with nil if fixed successfully.
//
// note: if @err is already nil, this method will be ignored.
func Try(err *error, retry func() error) {
	if err == nil || *err == nil {
		return
	}

	if handlers, ok := errHandler[*err]; ok {
		for _, cb := range handlers {
			if cb() {
				e := retry()
				if e == nil || e != *err {
					*err = e
					return
				}
			}
		}
	}
}
