type Cond struct {
        // L is held while observing or changing the condition
        L Locker
        // contains filtered or unexported fields
}
type Locker interface {
    Lock()
    Unlock()
}

func (*Cond) Wait()

func (*Cond) Signal()

func (c *Cond) Broadcast()