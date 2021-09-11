package common

import (
	"fmt"
	"strconv"
)

type Pagination struct {
	Page         int
	Rpp          int
	CurrentCount int
	TotalCount   int
}

func (p *Pagination) IsFirstPage() bool {
	return p.Page == 1
}

func (p *Pagination) IsLastPage() bool {
	return p.CurrentCount >= p.TotalCount
}

func (p *Pagination) AreAllRecordsReturned() bool {
	return p.Page*p.Rpp > p.TotalCount && p.Page == 1
}

func (p *Pagination) NextPage() {
	p.Page++
}

func (p *Pagination) PrevPage() {
	p.Page--
}

func (p *Pagination) UpdateCountsAndPrintCurrentInterval(totalCount string, resultsLength int) {
	p.TotalCount, _ = strconv.Atoi(totalCount)

	start := 1
	end := resultsLength
	if p.Page > 1 {
		start = p.Rpp * (p.Page - 1)
		end += start
	}
	if start == end {
		fmt.Println("No records found at this page")
	} else {
		fmt.Printf("Showing record(s) %d-%d out of %d record(s)\n", start, end, p.TotalCount)
	}

	p.CurrentCount = end
}
