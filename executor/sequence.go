package executor

func NewSequenceGenerator(total, count int) *SequenceGenerator {

	choicer := &SequenceGenerator{Total: total, Count: count, Result: make(chan *Sequence), Values: make([]int, count, count)}

	go func() {
		choicer.generate(0, 0)
		choicer.Result <- &Sequence{
			Result: nil,
			EOF:    true,
		}
	}()

	return choicer
}

func DestroySequenceGenerator(s *SequenceGenerator) {
	close(s.Result)
}

type Sequence struct {
	Result []int
	EOF    bool
}

/**
 */
type SequenceGenerator struct {
	Result chan *Sequence
	Total  int
	Count  int
	Values []int
}

func (s *SequenceGenerator) Next() Sequence {
	more := <-s.Result
	return *more
}

func (s *SequenceGenerator) generate(start int, idx int) {
	// C(M, N) 具体实现
	rest := s.Count - idx
	if rest == 0 {
		result := make([]int, s.Count, s.Count)
		for idx, val := range s.Values {
			result[idx] = val
		}

		s.Result <- &Sequence{
			Result: result,
			EOF:    false,
		}
		return
	}

	for i := start; i <= s.Total-rest; i++ {
		s.Values[idx] = i
		s.generate(i+1, idx+1)
	}
}
