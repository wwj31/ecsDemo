package graphUI

import (
	"testing"
)

type Int = int32

type STRU struct {
	I1 Int
	I2 Int
	I3 Int
}

func TestType(t *testing.T) {
	str := &STRU{123,5456,5665,}
	 fe := str.I2+123
	 FF(fe)
	FF(str.I1)
}

func FF(i int32)  {
	println(i)
}