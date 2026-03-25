// core/emit_common.go v1
package core

func canEmitRCFTermFromRange(r Range) (RCFTerm, bool) {
	if !r.Inside {
		return NewRCFTerm(nil), false
	}

	loFloor, _ := floorQuoRem(r.Lo.Value.Num(), r.Lo.Value.Den())
	hiFloor, _ := floorQuoRem(r.Hi.Value.Num(), r.Hi.Value.Den())

	if loFloor.Cmp(hiFloor) != 0 {
		return NewRCFTerm(nil), false
	}

	return NewRCFTerm(loFloor), true
}

// core/emit_common.go v1
