package executor

func NewSequenceGenerator(total, count int) *SequenceGenerator {

	choicer := &SequenceGenerator{Total: total, Count: count, Run: make(chan bool), More: make(chan bool), Values: make([]int, count, count)}

	go func() {
		choicer.generate(0, 0)
		choicer.More <- false
		<-choicer.Run
	}()

	return choicer
}

func DestroySequenceGenerator(s *SequenceGenerator) {
	close(s.Run)
	close(s.More)
}

/**
 */
type SequenceGenerator struct {
	Run    chan bool
	More   chan bool
	Total  int
	Count  int
	Values []int
}

func (s *SequenceGenerator) Next() ([]int, bool) {
	more := <-s.More

	result := make([]int, s.Count, s.Count)
	for idx, val := range s.Values {
		result[idx] = val
	}

	// trigger next round
	s.Run <- true

	return result, more
}

func (s *SequenceGenerator) generate(start int, idx int) {
	// C(M, N) 具体实现
	rest := s.Count - idx
	if rest == 0 {
		s.More <- true
		_ = <-s.Run
		return
	}

	for i := start; i <= s.Total-rest; i++ {
		s.Values[idx] = i
		s.generate(i+1, idx+1)
	}
}
