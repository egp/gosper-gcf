// core/emit_common.go v2
package core

func canEmitRCFTermFromRange(r Range) (RCFTerm, bool) {
	if !r.Inside {
		return NewRCFTerm(nil), false
	}

	loFloor, _, err := floorQuoRemChecked(r.Lo.Value.Num(), r.Lo.Value.Den())
	if err != nil {
		return NewRCFTerm(nil), false
	}

	hiFloor, _, err := floorQuoRemChecked(r.Hi.Value.Num(), r.Hi.Value.Den())
	if err != nil {
		return NewRCFTerm(nil), false
	}

	if loFloor.Cmp(hiFloor) != 0 {
		return NewRCFTerm(nil), false
	}

	return NewRCFTerm(loFloor), true
}

// core/emit_common.go v2
