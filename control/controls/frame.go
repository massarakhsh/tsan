package controls

import (
	"bitbucket.org/961961/tsan/control"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/shaman/lik/likdom"
)

//	Дескриптор элемента отображения формата 1+4
type FrameControl struct {
	control.DataControl
	Layout		int
	SizeLeft	int
	DivGor		int
	DivVert		int
}

const LAY_ONE = 0
const LAY_LR = 1
const LAY_TB = 2
const LAY_TBB = 3
const LAY_LLR = 4
const LAY_TTB = 5
const LAY_LRR = 6
const LAY_FOUR = 7

func (it *FrameControl) SetLayoutOne(left int) {
	it.setLayout(LAY_ONE, left, 100,100)
}
func (it *FrameControl) SetLayoutLR(left int, gor int) {
	it.setLayout(LAY_LR, left, gor,100)
}
func (it *FrameControl) SetLayoutTB(left int, vert int) {
	it.setLayout(LAY_TB, left, 100, vert)
}
func (it *FrameControl) SetLayoutTBB(left int, gor int, vert int) {
	it.setLayout(LAY_TBB, left, gor, vert)
}
func (it *FrameControl) SetLayoutLLR(left int, gor int, vert int) {
	it.setLayout(LAY_LLR, left, gor, vert)
}
func (it *FrameControl) SetLayoutTTB(left int, gor int, vert int) {
	it.setLayout(LAY_TTB, left, gor, vert)
}
func (it *FrameControl) SetLayoutLRR(left int, gor int, vert int) {
	it.setLayout(LAY_LRR, left, gor, vert)
}
func (it *FrameControl) SetLayoutFour(left int, gor int, vert int) {
	it.setLayout(LAY_FOUR, left, gor, vert)
}

//	Установить макет экрана
func (it *FrameControl) setLayout(layout int, left int, gor int, vert int) {
	it.Layout = layout
	it.SizeLeft = left
	it.DivGor = gor
	it.DivVert = vert
}

//	Отображение окна
func (it *FrameControl) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	div := likdom.BuildDivClassId("roll_data deck_data", "area_data")
	tbl := div.BuildTableClass("fill")
	ndx := 2
	if it.SizeLeft == 0 { ndx-- }
	if it.Layout == LAY_ONE || it.Layout == LAY_TB { ndx-- }
	ndy := 2
	if it.Layout == LAY_ONE || it.Layout == LAY_LR { ndy-- }
	w0 := it.SizeLeft
	w1 := (sx - control.BD * ndx - w0) * it.DivGor / 100
	w2 := sx - control.BD * ndx - w0 - w1
	h1 := (sy - control.BD * ndy) * it.DivVert / 100
	h2 := sy - control.BD * ndy - h1
	row := tbl.BuildTr()
	if it.SizeLeft > 0 {
		if td, w, h := it.BuildSection(row, w0, sy); td != nil {
			if it.Layout != LAY_ONE && it.Layout != LAY_LR {
				td.SetAttr("rowspan=2")
			}
			if ctrl := it.FindControl("L"); ctrl != nil {
				td.AppendItem(ctrl.BuildShow(rule, w-control.BD, h-control.BD))
			}
		}
	}
	if true {
		rs,cs, ww,hh := "1","1", w1,h1
		if it.Layout == LAY_TBB {
			cs = "2"
			ww += control.BD + w2
		} else if it.Layout == LAY_LRR {
			rs = "2"
			hh += control.BD + h2
		}
		if td, w, h := it.BuildSection(row, ww, hh); td != nil {
			td.SetAttr("colspan", cs, "rowspan", rs)
			if ctrl := it.FindControl("LU"); ctrl != nil {
				td.AppendItem(ctrl.BuildShow(rule, w-control.BD, h-control.BD))
			}
		}
	}
	if it.Layout != LAY_ONE && it.Layout != LAY_TB && it.Layout != LAY_TBB {
		rs,cs, ww,hh := "1","1", w2,h1
		if it.Layout == LAY_LLR {
			rs = "2"
			hh += control.BD + h2
		}
		if td, w, h := it.BuildSection(row, ww, hh); td != nil {
			td.SetAttr("colspan", cs, "rowspan", rs)
			if ctrl := it.FindControl("RU"); ctrl != nil {
				td.AppendItem(ctrl.BuildShow(rule, w-control.BD, h-control.BD))
			}
		}
	}
	if it.Layout != LAY_ONE && it.Layout != LAY_LR {
		row := tbl.BuildTr()
		if it.Layout != LAY_LRR {
			rs, cs, ww, hh := "1", "1", w1, h2
			if it.Layout == LAY_TTB {
				cs = "2"
				ww += control.BD + w2
			}
			if td, w, h := it.BuildSection(row, ww, hh); td != nil {
				td.SetAttr("colspan", cs, "rowspan", rs)
				if ctrl := it.FindControl("LD"); ctrl != nil {
					td.AppendItem(ctrl.BuildShow(rule, w-control.BD, h-control.BD))
				}
			}
		}
		if it.Layout != LAY_TB && it.Layout != LAY_LLR && it.Layout != LAY_TTB {
			rs, cs, ww, hh := "1", "1", w2, h1
			if td, w, h := it.BuildSection(row, ww, hh); td != nil {
				td.SetAttr("colspan", cs, "rowspan", rs)
				if ctrl := it.FindControl("RD"); ctrl != nil {
					td.AppendItem(ctrl.BuildShow(rule, w-control.BD, h-control.BD))
				}
			}
		}
	}
	return div
}

