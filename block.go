package cache

import "sort"

type block struct {
	start int
	total int
}

type sortBlocks []block

func (sb sortBlocks) Len() int {
	return len(sb)
}

func (sb sortBlocks) Swap(i, j int) {
	sb[i], sb[j] = sb[j], sb[i]
}

func (sb sortBlocks) Less(i, j int) bool {
	return sb[i].total < sb[j].total
}

func (sb *sortBlocks) add(start int, total int) {
	*sb = append(*sb, block{start: start, total: total})
}

func (sb *sortBlocks) getBlock(size int) (b block, ok bool) {
	sort.Sort(*sb)
	length := sb.Len()
	index := sort.Search(length, func(i int) bool {
		return (*sb)[i].total >= size
	})

	if index >= length || (*sb)[index].total < size {
		return
	}

	ok = true
	b = (*sb)[index]

	(*sb)[index] = (*sb)[length-1]
	(*sb) = (*sb)[:length-1]

	return
}
