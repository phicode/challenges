package rowcol

func ParseDirectionByte(direction byte) Direction {
	switch direction {
	case '>':
		return Right
	case '<':
		return Left
	case '^':
		return Up
	case 'v':
		return Down
	default:
		panic("invalid direction")
	}
}
