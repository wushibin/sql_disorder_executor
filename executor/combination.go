package executor

type LoopInfo struct {
	TagIndex int
	Count    int
}

type Combinator struct {
	InstructionFlagList []int
	EOF                 bool
}

// 标识所有的执行序列：采用递归实现标识的过程。如：有3个SQL文件，分别有2, 3, 4条SQL语句，总数9条SQL语句，用长度为9的数组进行标识。
// 执行的过程： 第一步：先从数组中从选择2个元素，标识为文件1， 第二步：再从剩余7个未标识的元素中，选择3个标识为文件2，第三步： 最后剩余的4个标识为文件3.
//
// 第一步的选择有C（9, 2）,枚举每一次选择的结果，将数组对应的元素，标识为文件1
// 第一步标识基础上，再进行第二步的标识，第二步可以标识的选择有C(7, 3), 枚举每一次的结果，将相应的数组元素标识为文件2
// 最后进行剩余数组元素标识为文件3
// 总的选择枚举数量为C(9,2)*C(7, 3)*C(4,4)
type CombinatorGenerator struct {
	// 需要进行标识的数组，标识为N表示选择第N个文件
	InstructionFlagList []int
	// 没有进行标识的InstructionFlagList的数组下标
	UnSelectedIndexList []int
	LoopList            []LoopInfo
	// 当前进行第N个文件元素的选择
	LoopCount           int
	Result              chan *Combinator
}

// infoList []LoopInfo: 表示总共多少个Loop(文件)，第个Loop(文件)中有多少个元素（SQL语句）
func NewCombinatorGenerator(infoList []LoopInfo) *CombinatorGenerator {
	// 计算所有元素的数量
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

	// 启动go routine, 并从generator.Result中读取当次选择的结果
	go func() {
		generator.produce()
		// 检举完成所有的排列， EOF标识为true
		generator.Result <- &Combinator{
			InstructionFlagList: nil,
			EOF:                 true,
		}
	}()

	return &generator
}

func DestroyCombinatorGenerator(generator *CombinatorGenerator) {
	// 关闭channel
	close(generator.Result)
}

func (s *CombinatorGenerator) Generate() Combinator {
	// 获取当前的枚举结果
	more := <-s.Result
	return *more
}

func (s *CombinatorGenerator) produce() {
	loop := s.LoopList[s.LoopCount]

	// 如果未选择的元素与当前需要选择的元素数量一致, 则可以一次将剩余的元素标识完成
	if len(s.UnSelectedIndexList) == loop.Count {

		// 标识InstructionFlagList中未选择的元素为当前Loop的TagIndex
		for _, flagIdx := range s.UnSelectedIndexList {
			s.InstructionFlagList[flagIdx] = loop.TagIndex
		}

		// 复制枚举的结果
		result := make([]int, len(s.InstructionFlagList), len(s.InstructionFlagList))
		for idx, val := range s.InstructionFlagList {
			result[idx] = val
		}

		// EOF标识为false, 表示还有枚举的排列可以产生
		s.Result <- &Combinator{
			InstructionFlagList: result,
			EOF:                 false,
		}

		return
	}

	// 未选择元素s.UnSelectedIndexList中选择loop.Count个的枚举
	seqGenerator := NewSequenceGenerator(len(s.UnSelectedIndexList), loop.Count)
	defer DestroySequenceGenerator(seqGenerator)

	for {
		// 获取未选择元素中选择loop.Count的一次选择结果
		seq := seqGenerator.Next()
		if seq.EOF == true {
			return
		}

		// 标识相应元素为当前loop的Tag
		for _, seq := range seq.Result {
			flagIdx := s.UnSelectedIndexList[seq]
			s.InstructionFlagList[flagIdx] = loop.TagIndex
		}

		// 从s.UnSelectedIndexList剔除已经选择的元素, 生成新的nextUnSelectedIndexLists标识为所有未标识的元素
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

		// 递归将剩余未标识的元素标识
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
