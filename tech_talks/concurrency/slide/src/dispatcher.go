package main

// START OMIT
func mainMsgDispatcher() {
LOOP:
	for {
		select {
		case evt, ok := <-chanA:
			if !ok {
				chanA = nil
			} else { // msg handler
			}
		case cmd, ok := <-chanB:
			if !ok {
				chanB = nil
			} else { // msg handler
			}
		case err, ok := <-chanC:
			if !ok {
				chanC = nil
			} else {
			}
		}
		if nil == chanA && nil == chanB && nil == chanC {
			break LOOP
		}
	}
}
// END OMIT