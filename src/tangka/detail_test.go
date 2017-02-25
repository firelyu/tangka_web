package tangka_test

import (
	. "github.com/firelyu/tangka_web/src/tangka"
	"strconv"
	"testing"
)

func TestNewTangka(t *testing.T) {
	loopCount := 100
	var list []*Tangka

	// Test default name and author
	oneTangka := NewTangka("", "")
	if oneTangka.Name != DefaultTangkaName || oneTangka.Author != DefaultTangkaAuthor {
		t.Errorf("The default name(%s) or default author(%s) is wrong. Default name is %s, and default author is %s",
			oneTangka.Name, oneTangka.Author, DefaultTangkaName, DefaultTangkaAuthor)
	}

loop1:
	for i := 0; i < loopCount; i++ {
		name := "name" + strconv.Itoa(i)
		author := "author" + strconv.Itoa(i)
		tangka := NewTangka(name, author)

		// Test the id is unique
		for _, c := range []byte(tangka.Id) {
			if !('a' <= c && c <= 'z' || '0' <= c && c <= '9') {
				t.Errorf("The %s has invalid char %b", tangka.Id, c)
				break loop1
			}
		}

		// Test the default name and author
		if tangka.Name != name {
			t.Errorf("The Name(%s) is not the input one %s", tangka.Name, name)
			break loop1
		}

		if tangka.Author != author {
			t.Errorf("The Author(%s) is not the input one %s", tangka.Author, author)
			break loop1
		}

		list = append(list, tangka)
	}

loop2:
	// Test the id is unique
	for index, tangka := range list {
		for eachIndex, each := range list {
			// Don't compare self
			if eachIndex == index {
				break
			}
			if tangka.Id == each.Id {
				t.Errorf("The id of %v is the same as %v", tangka, each)
				break loop2
			}
		}
	}
}

func BenchmarkNewTangka(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewTangka(DefaultTangkaName+strconv.Itoa(i), DefaultTangkaAuthor+strconv.Itoa(i))
	}
}
