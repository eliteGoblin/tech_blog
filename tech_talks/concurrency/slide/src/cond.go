type Record struct {
	sync.Mutex
	data string

	cond *sync.Cond
}

func NewRecord() *Record {
	r := Record{}
	r.cond = sync.NewCond(&r)
	return &r
}

func main() {
	var wg sync.WaitGroup

	rec := NewRecord()
	wg.Add(1)
	go func(rec *Record) {
		defer wg.Done()
		rec.Lock()
		rec.cond.Wait()
		rec.Unlock()
		fmt.Println("Data: ", rec.data)
		return
	}(rec)

	time.Sleep(2 * time.Second)
	rec.Lock()
	rec.data = "gopher"
	rec.Unlock()

	rec.cond.Signal()

	wg.Wait() // wait till all goutine completes
}