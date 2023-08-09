package japanese_dict

import (
	"fmt"
	"os"
	"testing"
)

func TestLoadJmdict(t *testing.T) {
	filepath := "./test_data/JMdict_e.xml"
	reader, err := os.Open(filepath)
	if err != nil {
		t.Errorf("Error opening file: %s", err)
	}

	jmdict, _, err := LoadJmdict(reader)
	if len(jmdict.Entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(jmdict.Entries))
	}
	if jmdict.Entries[0].Kanji[0].Expression != "日本" {
		t.Errorf("Expected 日本, got %s", jmdict.Entries[0].Kanji[0].Expression)
	}

	fmt.Println(jmdict)
}

func TestLoadJmdictMap(t *testing.T) {
	filepath := "./test_data/JMdict_e.xml"
	reader, err := os.Open(filepath)
	if err != nil {
		t.Errorf("Error opening file: %s", err)
	}

	dictMap, err := LoadJmdictMap(reader)
	if len(dictMap) != 1 {
		t.Errorf("Expected 1 entries, got %d", len(dictMap))
	}
	if dictMap["日本"][0].Kanji[0].Expression != "日本" {
		t.Errorf("Expected 日本, got %s", dictMap["日本"][0].Kanji[0].Expression)
	}

	fmt.Println(dictMap)
}
