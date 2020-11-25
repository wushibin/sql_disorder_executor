package executor

type LoopInfo struct {
	TagIndex int
	Count    int
}

type Combinator struct {
	InstructionFlagList []int
	EOF                 bool
}

type CombinatorGenerator struct {
	InstructionFlagList []int
	UnSelectedIndexList []int
	LoopList            []LoopInfo
	LoopCount           int
	Result              chan *Combinator
}

func NewCombinatorGenerator(infoList []LoopInfo) *CombinatorGenerator {
	total := 0
	for _, info := range infoList {
		total += info.Count
	}

	incrIdxList := make([]int, total, total)
	for idx, _ := range incrIdxList {
		incrIdxList[idx] = idx
	}

	flagList := make([]int, total, total)

	generator := CombinatorGenerator{
		InstructionFlagList: flagList,
		UnSelectedIndexList: incrIdxList,
		LoopList:            infoList,
		LoopCount:           0,
		Result:              make(chan *Combinator),
	}

	go func() {
		generator.produce()
		generator.Result <- &Combinator{
			InstructionFlagList: nil,
			EOF:                 true,
		}
	}()

	return &generator
}

func DestroyCombinatorGenerator(generator *CombinatorGenerator) {
	close(generator.Result)
}

func (s *CombinatorGenerator) Generate() Combinator {
	more := <-s.Result
	return *more
}

func (s *CombinatorGenerator) produce() {
	loop := s.LoopList[s.LoopCount]

	if len(s.UnSelectedIndexList) == loop.Count {
		for _, flagIdx := range s.UnSelectedIndexList {
			s.InstructionFlagList[flagIdx] = loop.TagIndex
		}

		result := make([]int, len(s.InstructionFlagList), len(s.InstructionFlagList))
		for idx, val := range s.InstructionFlagList {
			result[idx] = val
		}

		// notify data is ready
		s.Result <- &Combinator{
			InstructionFlagList: result,
			EOF:                 false,
		}

		return
	}

	seqGenerator := NewSequenceGenerator(len(s.UnSelectedIndexList), loop.Count)
	defer DestroySequenceGenerator(seqGenerator)

	for {
		seq := seqGenerator.Next()
		if seq.EOF == true {
			return
		}

		for _, seq := range seq.Result {
			flagIdx := s.UnSelectedIndexList[seq]
			s.InstructionFlagList[flagIdx] = loop.TagIndex
		}

		nextUnSelectedIndexList := make([]int, len(s.UnSelectedIndexList)-loop.Count)
		cmpIdx := 0
		idx := 0
		for i := 0; i < len(s.UnSelectedIndexList); i++ {
			if int(i) == int(seq.Result[cmpIdx]) {
				if cmpIdx < len(seq.Result)-1 {
					cmpIdx++
				}
			} else {
				nextUnSelectedIndexList[idx] = s.UnSelectedIndexList[i]
				idx++
			}
		}

		nextCombinatorGenerator := CombinatorGenerator{
			InstructionFlagList: s.InstructionFlagList,
			UnSelectedIndexList: nextUnSelectedIndexList,
			LoopList:            s.LoopList,
			LoopCount:           s.LoopCount + 1,
			Result:              s.Result,
		}

		nextCombinatorGenerator.produce()
	}
}
