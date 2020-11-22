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
	More                chan bool
	Run                 chan bool
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
		More:                make(chan bool),
		Run:                 make(chan bool),
	}

	go func() {
		generator.produce()
		generator.More <- false
		<-generator.Run
	}()

	return &generator
}

func DestroyCombinatorGenerator(generator *CombinatorGenerator) {
	close(generator.Run)
	close(generator.More)
}

func (s *CombinatorGenerator) Generate() Combinator {
	more := <-s.More

	size := len(s.InstructionFlagList)
	result := make([]int, size, size)
	if more == true {
		for idx, val := range s.InstructionFlagList {
			result[idx] = val
		}
	}

	// next loop
	s.Run <- true

	return Combinator{InstructionFlagList: result, EOF: !more}
}

func (s *CombinatorGenerator) produce() {
	loop := s.LoopList[s.LoopCount]

	if len(s.UnSelectedIndexList) == loop.Count {
		for _, flagIdx := range s.UnSelectedIndexList {
			s.InstructionFlagList[flagIdx] = loop.TagIndex
		}

		// notify data is ready
		s.More <- true

		// wait to run next round
		_ = <-s.Run
		return
	}

	seqGenerator := NewSequenceGenerator(len(s.UnSelectedIndexList), loop.Count)
	defer DestroySequenceGenerator(seqGenerator)

	for {
		seqs, more := seqGenerator.Next()
		if more == false {
			return
		}

		for _, seq := range seqs {
			flagIdx := s.UnSelectedIndexList[seq]
			s.InstructionFlagList[flagIdx] = loop.TagIndex
		}

		nextUnSelectedIndexList := make([]int, len(s.UnSelectedIndexList)-loop.Count)
		cmpIdx := 0
		idx := 0
		for i := 0; i < len(s.UnSelectedIndexList); i++ {
			if int(i) == int(seqs[cmpIdx]) {
				if cmpIdx < len(seqs)-1 {
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
			More:                s.More,
			Run:                 s.Run,
		}

		nextCombinatorGenerator.produce()
	}
}
