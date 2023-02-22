// Copyright (c) Fortio Authors, All Rights Reserved
// See LICENSE for licensing terms. (Apache-2.0)

package sets_test

import (
	"testing"

	"fortio.org/assert"
	"fortio.org/sets"
)

func TestSetToString(t *testing.T) {
	s := sets.Set[string]{"z": {}, "a": {}, "c": {}, "b": {}}
	assert.Equal(t, "a,b,c,z", s.String())
	assert.Equal(t, s.Len(), 4)
	s.Clear()
	assert.Equal(t, "", s.String())
	assert.Equal(t, s.Len(), 0)
}

func TestArrayToSet(t *testing.T) {
	a := []string{"z", "a", "c", "b"}
	s := sets.FromSlice(a)
	assert.Equal(t, "a,b,c,z", s.String())
	assert.Equal(t, sets.Sort(s), []string{"a", "b", "c", "z"})
}

func TestRemoveCommon(t *testing.T) {
	setA := sets.New("a", "b", "c", "d")
	setB := sets.New("b", "d", "e", "f", "g")
	setAA := setA.Clone()
	setBB := setB.Clone()
	sets.RemoveCommon(setAA, setBB)
	assert.Equal(t, "a,c", setAA.String())   // removed
	assert.Equal(t, "e,f,g", setBB.String()) // added
	// Swap order to exercise the optimization on length of iteration
	// also check clone is not modifying the original etc
	setAA = setB.Clone() // putting B in AA on purpose and vice versa
	setBB = setA.Clone()
	assert.True(t, setAA.Equals(setB))
	assert.True(t, setB.Equals(setAA))
	assert.False(t, setAA.Equals(setA))
	assert.False(t, setB.Equals(setBB))
	sets.XOR(setAA, setBB)
	assert.Equal(t, "a,c", setBB.String())
	assert.Equal(t, "e,f,g", setAA.String())
	assert.True(t, setBB.Has("c"))
	setBB.Remove("c")
	assert.False(t, setBB.Has("c"))
}

func TestMinus(t *testing.T) {
	setA := sets.New("a", "b", "c", "d")
	setB := sets.New("b", "d", "e", "f", "g")
	setAB := setA.Clone().Minus(setB)
	setBA := setB.Clone().Minus(setA)
	assert.Equal(t, "a,c", setAB.String())
	assert.Equal(t, "e,f,g", setBA.String())
}

func TestPlus(t *testing.T) {
	setA := sets.New("a", "b", "c", "d")
	setB := sets.New("b", "d", "e", "f", "g")
	setAB := setA.Clone().Plus(setB)
	setBA := setB.Clone().Plus(setA)
	assert.Equal(t, "a,b,c,d,e,f,g", setAB.String())
	assert.Equal(t, "a,b,c,d,e,f,g", setBA.String())
}

func TestUnion(t *testing.T) {
	setA := sets.New("a", "b", "c", "d")
	setB := sets.New("b", "d", "e", "f", "g")
	setC := sets.Union(sets.Union[string](), setA, setB)
	assert.Equal(t, "a,b,c,d,e,f,g", setC.String())
}

func TestIntersection1(t *testing.T) {
	setA := sets.New("a", "b", "c", "d")
	setB := sets.New("b", "d", "e", "f", "g")
	setC := sets.Intersection(setA, setB)
	assert.Equal(t, "b,d", setC.String())
}

func TestIntersection2(t *testing.T) {
	assert.Equal(t, len(sets.Intersection[string]()), 0)
	setA := sets.New("a", "b", "c")
	setB := sets.New("d", "e", "f")
	// cover stop early when empty intersection is reached, ie 3rd set won't be looked at
	setC := sets.Intersection(setA, setB, setA)
	assert.Equal(t, "", setC.String())
}

func TestSubset(t *testing.T) {
	setA := sets.New("a", "b", "c", "d")
	setB := sets.New("b", "d", "e", "f", "g")
	setC := sets.New("b", "d")
	assert.True(t, setC.Subset(setA))
	assert.True(t, setA.Subset(setA))
	assert.False(t, setA.Subset(setC))
	assert.False(t, setA.Subset(setB))
	assert.False(t, setB.Subset(setA))
}
