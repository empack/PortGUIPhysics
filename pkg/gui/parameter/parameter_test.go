package parameter

import (
	"fmt"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/test"
	"testing"
	"time"
)

const testTimeout = 100 * time.Millisecond // max time until change has to be propagated

const testName1 = "ParameterTestName"
const testName2 = "SomeOtherTestName"
const testDefault float64 = 10.0

var uut *Parameter
var nameValue binding.String
var minValue binding.Float
var valValue binding.Float
var maxValue binding.Float
var defaultValue binding.Float
var checkValue binding.Bool

func init() {
	nameValue = binding.NewString()
	minValue = binding.NewFloat()
	maxValue = binding.NewFloat()
	valValue = binding.NewFloat()
	defaultValue = binding.NewFloat()
	checkValue = binding.NewBool()

	uut = NewParameter(nameValue, defaultValue, valValue, minValue, maxValue, checkValue)
}
func TestParameterNameListener(t *testing.T) {
	t.Log("WARNING: Test with inconsistent determinacy") //TODO remove if better solution found
	wasNotified := false
	dataListener := binding.NewDataListener(func() {
		wasNotified = true
	})
	nameValue.AddListener(dataListener)
	if err := nameValue.Set(testName1); err != nil {
		t.Skip("Failed to change name binding data")
	}
	time.Sleep(testTimeout)
	if wasNotified == false {
		t.Errorf("ParameterNameNotification() failed. Listener was not called at binding name Set after  %dms", testTimeout/time.Millisecond)
	}
	wasNotified = false
	test.Type(uut.name, "T")
	time.Sleep(testTimeout)
	if wasNotified == false {
		t.Errorf("ParameterNameNotification() failed. Listener was not called at name entry text change after  %dms", testTimeout/time.Millisecond)
	}

	wasNotified = false
	nameValue.RemoveListener(dataListener)

	if err := nameValue.Set(testName2); err != nil {
		t.Skip("Failed to change name binding data")
	}
	time.Sleep(testTimeout)
	if wasNotified == true {
		t.Errorf("ParameterNameNotification() failed. Listener was called at binding name Set after removed in the next %dms", testTimeout/time.Millisecond)
	}

	test.Type(uut.name, "TestChars")
	time.Sleep(testTimeout)
	if wasNotified == true {
		t.Errorf("ParameterNameNotification() failed. Listener was called at name entry text change after removed in the next %dms", testTimeout/time.Millisecond)
		t.FailNow()
	}
}

func TestSetParameterName(t *testing.T) {
	t.Log("WARNING: Test with inconsistent determinacy") //TODO remove if better solution found

	if err := nameValue.Set(testName1); err != nil {
		t.Skip("Failed to change name binding data")
	}
	time.Sleep(testTimeout)
	if realName := uut.name.Text; testName1 != realName {
		t.Errorf("SetParameterName() failed. Expected %s, got %s after %dms", testName1, realName, testTimeout/time.Millisecond)
	}
}

func TestSetParameterDefault(t *testing.T) {
	t.Log("WARNING: Test with inconsistent determinacy") //TODO remove if better solution found

	if err := defaultValue.Set(testDefault); err != nil {
		t.Skip("Failed to change default binding data")
	}
	time.Sleep(testTimeout)
	if placeHolder := uut.val.PlaceHolder; fmt.Sprint(testDefault) != placeHolder {
		t.Errorf("TestSetParameterDefault() failed. Expected %s, got %s after %dms", fmt.Sprint(testDefault), placeHolder, testTimeout/time.Millisecond)
	}
}

func TestSetParameterCheck(t *testing.T) {
	t.Log("WARNING: Test with inconsistent determinacy") //TODO remove if better solution found

	if err := checkValue.Set(true); err != nil {
		t.Skip("Failed to change check binding data")
	}
	time.Sleep(testTimeout)
	if check := uut.check.Checked; true != check {
		t.Errorf("TestSetParameterCheck() failed. Expected true, got %t at binding Set after %dms", check, testTimeout/time.Millisecond)
	}

	prev := uut.check.Checked
	test.Tap(uut.check)
	time.Sleep(testTimeout)
	if check := uut.check.Checked; prev == check {
		t.Errorf("TestSetParameterCheck() failed. Expected %t, got %t at object interaction after %dms", !prev, check, testTimeout/time.Millisecond)
	}
}
