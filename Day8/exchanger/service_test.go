package exchanger

import "testing"

func TestConvertAmount(t *testing.T) {
	service := NewExchangeService(nil)
	result := service.ConvertAmount(2.0, 100.0)
	expected := 200.0

	if result != expected {
		t.Errorf("ConvertAmount(2.0, 100.0) = %f, expected %f", result, expected)
	}
}
