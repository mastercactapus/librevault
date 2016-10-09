package librevault

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseSecret(t *testing.T) {
	testOwn := func(keyStr, ownKey string) {
		t.Run(fmt.Sprintf("'%s'.Owner='%s'", keyStr, ownKey), func(t *testing.T) {
			key, err := ParseSecret(keyStr)
			if err != nil {
				t.Fatal(err)
			}
			if key.Owner() != ownKey {
				t.Fail()
			}
		})
	}
	testRead := func(keyStr, readKey string) {
		t.Run(fmt.Sprintf("'%s'.ReadOnly='%s'", keyStr, readKey), func(t *testing.T) {
			key, err := ParseSecret(keyStr)
			if err != nil {
				t.Fatal(err)
			}
			if key.ReadOnly() != readKey {
				t.Fail()
			}
		})
	}
	testDown := func(keyStr, downKey string) {
		t.Run(fmt.Sprintf("'%s'.DownloadOnly='%s'", keyStr, downKey), func(t *testing.T) {
			key, err := ParseSecret(keyStr)
			if err != nil {
				t.Fatal(err)
			}
			if key.DownloadOnly() != downKey {
				t.Fail()
			}
		})
	}

	testData := []struct{ Own, Read, Down string }{
		{
			"A1fFfr3UMHoLqjoXPSaWHRySvijJrKJFPz3X8MtnNAzXTZ",
			"C1ETdSkHLVeNPWfqLTsUDWPCUZqKCzF5qjFJtys8KPT3wdQxgtkxk1WTuvZbZx2WJQ9Pd1DBgs6deoBsTNEgFyXNMh1",
			"D1AMcu13VWLTfKZfJNxkm18PeRQfJ3jfp19SirnurWzXfhV",
		},
		{
			"A1BnkZ49DFzBBsV1UiANUedYD4UpjtdB3Yg1wjpvj4dxUpQ",
			"C1CcaQPztiTd5rJx5iRhpHfYH3H8HKMx8WEynky63HKAHi8hHE8owygFBLuGNTypaakCACecUgFv7hTmWxFDRJNJFXw",
			"D1VvGmtwhfhQVahXMmrDjn21Dz3JoTTAjMWNEe39BwUhT3",
		},
		{
			"A16cSxDkq4MTqNSkHeVfifxFsbiXgrc7i4VFY2wA5MH3GNZ",
			"C1HcQnCwgcubRLoyqBxhzzUmYY5YprXZyHAfEodp3nA2Ayv9w2VAFcYcSeemNw6PsoHbETzgpnbbkpNFCSxMF2M2XM5",
			"D17wLRxZAabCfwkpFT996rtsXUsKsZDYeGxkGeaZUi71CPk",
		},
	}

	for _, data := range testData {
		testOwn(data.Own, data.Own)
		testOwn(data.Read, "")
		testOwn(data.Down, "")

		testRead(data.Own, data.Read)
		testRead(data.Read, data.Read)
		testRead(data.Down, "")

		testDown(data.Own, data.Down)
		testDown(data.Read, data.Down)
		testDown(data.Down, data.Down)

	}
}

func TestMarshalJSON(t *testing.T) {
	key, err := ParseSecret("A16cSxDkq4MTqNSkHeVfifxFsbiXgrc7i4VFY2wA5MH3GNZ")
	if err != nil {
		t.Fatal(err)
	}

	s := struct{ S *Secret }{S: key}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	expected := `{"S":"A16cSxDkq4MTqNSkHeVfifxFsbiXgrc7i4VFY2wA5MH3GNZ"}`
	if string(data) != expected {
		t.Errorf("expected '%s' but got '%s'", `{"S":"A16cSxDkq4MTqNSkHeVfifxFsbiXgrc7i4VFY2wA5MH3GNZ"}`, string(data))
	}
}

func TestUnmarshalJSON(t *testing.T) {
	data := []byte(`{"S":"A16cSxDkq4MTqNSkHeVfifxFsbiXgrc7i4VFY2wA5MH3GNZ"}`)

	var s struct {
		S *Secret
	}

	err := json.Unmarshal(data, &s)
	if err != nil {
		t.Fatal(err)
	}
	expected := "A16cSxDkq4MTqNSkHeVfifxFsbiXgrc7i4VFY2wA5MH3GNZ"
	if s.S.String() != expected {
		t.Errorf("expected '%s' but got %s", expected, s.S.String())
	}
}
