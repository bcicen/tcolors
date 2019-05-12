package state

const (
	NoChange Change = 1 << iota
	SelectedChanged
	HueChanged
	SaturationChanged
	ValueChanged
	ErrorMsgChanged
)

const AllChanged = SelectedChanged | HueChanged | SaturationChanged | ValueChanged

// Change represents currently pending state changes
type Change uint8

func (sc Change) Includes(other ...Change) bool {
	for _, x := range other {
		if sc&x == x {
			return true
		}
	}
	return false
}
