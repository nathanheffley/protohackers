package main

import "testing"

func TestValidJson(t *testing.T) {
	json := []byte("{\"method\":\"isPrime\",\"number\":123}")
	err := ValidateJson(json)
	if err != nil {
		t.Fatalf(`ValidateJson(%s) should not return an error, %s`, json, err)
	}
}

func TestValidExtraJson(t *testing.T) {
	json := []byte("{\"method\":\"isPrime\",\"number\":123,\"extra\":\"who cares?\"}")
	err := ValidateJson(json)
	if err != nil {
		t.Fatalf(`ValidateJson(%s) should not return an error, %s`, json, err)
	}
}

func TestMalformedJson(t *testing.T) {
	json := []byte("not json")
	err := ValidateJson(json)
	if err == nil {
		t.Fatalf(`ValidateJson(%s) should return an error`, json)
	}
}

func TestJsonArrayIsInvalid(t *testing.T) {
	json := []byte("[\"method\", \"number\"]")
	err := ValidateJson(json)
	if err == nil {
		t.Fatalf(`ValidateJson(%s) should return an error`, json)
	}
}

func TestMissingMethodIsInvalid(t *testing.T) {
	json := []byte("{\"number\":123}")
	err := ValidateJson(json)
	if err == nil {
		t.Fatalf(`ValidateJson(%s) should return an error`, json)
	}
}

func TestWrongMethodIsInvalid(t *testing.T) {
	json := []byte("{\"method\":\"wrongMethod\",\"number\":123}")
	err := ValidateJson(json)
	if err == nil {
		t.Fatalf(`ValidateJson(%s) should return an error`, json)
	}
}

func TestMissingNumberIsInvalid(t *testing.T) {
	json := []byte("{\"method\":\"isPrime\"}")
	err := ValidateJson(json)
	if err == nil {
		t.Fatalf(`ValidateJson(%s) should return an error`, json)
	}
}

func TestWrongNumberIsInvalid(t *testing.T) {
	json := []byte("{\"method\":\"isPrime\",\"number\":\"123\"}")
	err := ValidateJson(json)
	if err == nil {
		t.Fatalf(`ValidateJson(%s) should return an error`, json)
	}
}
