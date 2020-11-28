package executor

// 构造从Total个元素中选择Count个元素的序列生成器
func NewSequenceGenerator(total, count int) *SequenceGenerator {

	choicer := &SequenceGenerator{Total: total, Count: count, Result: make(chan *Sequence), Values: make([]int, count, count)}

	// 启动go routine实现元素选择的算法，并通过choice.Result从go routine中同步读取当次选择的序列结果
	go func() {
		choicer.generate(0, 0)
		// 执行结束，标示EOF为True，表示所有序列已经枚举完成
		choicer.Result <- &Sequence{
			Result: nil,
			EOF:    true,
		}
	}()

	return choicer
}

// 关闭实例创建中的channel
func DestroySequenceGenerator(s *SequenceGenerator) {
	close(s.Result)
}

type Sequence struct {
	Result []int
	EOF    bool
}

// 枚举C(M, N)选择序列，如从3个元素中选择2个
// 依次会产生数组序列:（1, 2）, (1, 3), (2, 3)
type SequenceGenerator struct {
	Result chan *Sequence
	Total  int
	Count  int
	Values []int
}

// 获取C(M,N)当次选择的结果
func (s *SequenceGenerator) Next() Sequence {
	// 同步从s.Result中读取出当前选择的结果
	more := <-s.Result
	return *more
}

// 采用递归的方式实现C(M, N)的算法，(从M个元素中选择N个元素)
// start : M数组的起始位置
// idx: 选择第几个元素
func (s *SequenceGenerator) generate(start int, idx int) {
	rest := s.Count - idx
	// 是否选择足够count个元素
	if rest == 0 {
		// 复制当次C(M，N)的元素序列选择
		result := make([]int, s.Count, s.Count)
		for idx, val := range s.Values {
			result[idx] = val
		}

		// 通过channel发送选择序列的结果
		// EOF标识为false, 表示选择的C(M,N)还有其它的枚举序列可以产生
		s.Result <- &Sequence{
			Result: result,
			EOF:    false,
		}
		return
	}

	for i := start; i <= s.Total-rest; i++ {
		// 保存当前选择的元素
		s.Values[idx] = i
		// 递归选择下一个元素
		s.generate(i+1, idx+1)
	}
}
