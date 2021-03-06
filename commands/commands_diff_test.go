package commands

import (
	"testing"
)

func TestDiffTags(t *testing.T) {

	tagsBefore := []Tag{
		Tag{Path: "0000/0000", Hash: hashing("0000000000000000000000000000000000000000000000000000000000000000"), Timestamp: 0x00000000},
		Tag{Path: "1111/1111", Hash: hashing("1111111111111111111111111111111111111111111111111111111111111111"), Timestamp: 0x11111111},
		Tag{Path: "2222/2222", Hash: hashing("2222222222222222222222222222222222222222222222222222222222222222"), Timestamp: 0x22222222},
		Tag{Path: "3333/3333", Hash: hashing("3333333333333333333333333333333333333333333333333333333333333333"), Timestamp: 0x33333333},
	}
	tagsAfter := []Tag{
		Tag{Path: "0000/0000", Hash: hashing("0000000000000000000000000000000000000000000000000000000000000000"), Timestamp: 0x00000000},
		Tag{Path: "1111/1111", Hash: hashing("1111111111111111111111111111111111111111111111111111111111111111"), Timestamp: 0x11111111},
		Tag{Path: "2222/2222", Hash: hashing("2222222222222222222222222222222222222222222222222222222222222222"), Timestamp: 0x22222222},
		Tag{Path: "3333/3333", Hash: hashing("3333333333333333333333333333333333333333333333333333333333333333"), Timestamp: 0x33333333},
		Tag{Path: "4444/4444", Hash: hashing("4444444444444444444444444444444444444444444444444444444444444444"), Timestamp: 0x44444444},
		Tag{Path: "5555/5555", Hash: hashing("5555555555555555555555555555555555555555555555555555555555555555"), Timestamp: 0x55555555},
	}

	// func diffTags(tagsBefore, tagsAfter []Tag) (deleted []Tag, added []Tag)
	t.Run("diffTags()", func(t *testing.T) {
		deleted, added := diffTags(tagsBefore, tagsAfter)
		if len(deleted) != 0 {
			t.Errorf("diffTags() return a slice %s, want the empty slice", deleted)
		}
		if len(added) != 2 ||
			added[0].String() != "4444444444444444444444444444444444444444444444444444444444444444 44444444 4444/4444" ||
			added[1].String() != "5555555555555555555555555555555555555555555555555555555555555555 55555555 5555/5555" {
			t.Errorf("diffTags() return a added slice %s", added)
		}

		tagsAfter = []Tag{
			Tag{Path: "0000/0000", Hash: hashing("0000000000000000000000000000000000000000000000000000000000000000"), Timestamp: 0x00000000},
			Tag{Path: "1111/ffff", Hash: hashing("1111111111111111111111111111111111111111111111111111111111111111"), Timestamp: 0x11111111},
			Tag{Path: "2222/2222", Hash: hashing("22222222222222222222222222222222ffffffffffffffffffffffffffffffff"), Timestamp: 0x22222222},
			Tag{Path: "3333/3333", Hash: hashing("3333333333333333333333333333333333333333333333333333333333333333"), Timestamp: 0x3333ffff},
		}
		deleted, added = diffTags(tagsBefore, tagsAfter)
		if len(deleted) != 3 ||
			deleted[0].String() != "1111111111111111111111111111111111111111111111111111111111111111 11111111 1111/1111" ||
			deleted[1].String() != "2222222222222222222222222222222222222222222222222222222222222222 22222222 2222/2222" ||
			deleted[2].String() != "3333333333333333333333333333333333333333333333333333333333333333 33333333 3333/3333" {
			t.Errorf("diffTags() return a deleted slice %s", deleted)
		}
		if len(added) != 3 ||
			added[0].String() != "1111111111111111111111111111111111111111111111111111111111111111 11111111 1111/ffff" ||
			added[1].String() != "22222222222222222222222222222222ffffffffffffffffffffffffffffffff 22222222 2222/2222" ||
			added[2].String() != "3333333333333333333333333333333333333333333333333333333333333333 3333ffff 3333/3333" {
			t.Errorf("diffTags() return a added slice %s", added)
		}
	})
}
