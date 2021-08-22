package component

func (s ComponentType) String() string {
	return ComponentType_name[uint64(s)]
}

// Enum value maps for ComponentType.
var (
	ComponentType_name = map[uint64]string{
		MOVE_COMP.ComponentType():      "MOVE_COMP",
		AREA_COMP.ComponentType():      "AREA_COMP",
		POS_COMP.ComponentType():       "POS_COMP",
		ATTRIBUTE_COMP.ComponentType(): "ATTRIBUTE_COMP",
		PLAY_COMP.ComponentType():      "PLAY_COMP",
		FIGHTING_COMP.ComponentType():  "FIGHTING_COMP",
	}
)
