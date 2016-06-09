package awakengine

import (
	"reflect"
	"testing"

	"github.com/DrJosh9000/vec"
)

func TestViewInvalidation(t *testing.T) {
	v := &View{}
	v.compute()
	if !v.valid {
		t.Errorf("After compute, got valid value %t, want true", v.valid)
	}
	v.invalidate()
	if v.valid {
		t.Errorf("After invalidate, got valid value %t, want false", v.valid)
	}
	v.SetBounds(vec.Rect{})
	if v.valid {
		t.Errorf("After SetBounds, got valid value %t, want false", v.valid)
	}
	v.Bounds()
	if !v.valid {
		t.Errorf("After Bounds, got valid value %t, want true", v.valid)
	}
	v.SetOffset(vec.I2{})
	if v.valid {
		t.Errorf("After SetOffset, got valid value %t, want false", v.valid)
	}
	v.Offset()
	if !v.valid {
		t.Errorf("After Offset, got valid value %t, want true", v.valid)
	}
	v.SetZ(0)
	if v.valid {
		t.Errorf("After SetZ, got valid value %t, want false", v.valid)
	}
	v.Z()
	if !v.valid {
		t.Errorf("After Z, got valid value %t, want true", v.valid)
	}
	w := &View{}
	w.SetParent(v)
	if !v.valid {
		t.Errorf("After w.SetParent(v), got v.valid value %t, want true", v.valid)
	}
	if w.valid {
		t.Errorf("After w.SetParent(v), got w.valid value %t, want false", w.valid)
	}
	v.invalidate()
	if v.valid {
		t.Errorf("After v.invalidate, got v.valid value %t, want false", v.valid)
	}
	if w.valid {
		t.Errorf("After v.invalidate, got w.valid value %t, want false", w.valid)
	}
}

func TestSetParent(t *testing.T) {
	parent := &View{}
	child1 := &View{}
	child2 := &View{}
	child1.SetParent(parent)
	if got, want := parent.children, []*View{child1}; !reflect.DeepEqual(got, want) {
		t.Errorf("Got parent.children %v, want %v", got, want)
	}
	child2.SetParent(parent)
	if got, want := parent.children, []*View{child1, child2}; !reflect.DeepEqual(got, want) {
		t.Errorf("Got parent.children %v, want %v", got, want)
	}
	child1.SetParent(child2)
	if got, want := parent.children, []*View{child2}; !reflect.DeepEqual(got, want) {
		t.Errorf("Got parent.children %v, want %v", got, want)
	}
	if got, want := child2.children, []*View{child1}; !reflect.DeepEqual(got, want) {
		t.Errorf("Got parent.children %v, want %v", got, want)
	}
}
