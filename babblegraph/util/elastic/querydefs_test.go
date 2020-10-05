package elastic

func TestInQuery(t *testing.T) {
	testQuery := InQuery{
		FieldName: "text",
		Values: ["abc", "123"],
	}
	expected := `{"terms":{"text":["abc","123"]}}`
	out, err := json.Marshal(testQuery)
	if err != nil {
		t.Errorf(err.Error())
	}
	if string(out) != expected {
		t.Errorf("Expected %s, got %s", expected, string(out))
	}
}
