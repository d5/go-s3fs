package s3fs_test

import (
	"path"
	"sort"
	"strings"
	"testing"

	"github.com/d5/go-s3fs"
	"github.com/stretchr/testify/assert"
)

func TestClient_ListAll(t *testing.T) {
	client := MustCreateS3Client()

	testDir := randomTestFilePath()

	testFilePaths := []string{
		path.Join(testDir, "file01"),
		path.Join(testDir, "file02"),
		path.Join(testDir, "file03"),
		path.Join(testDir, "sub1/file04"),
		path.Join(testDir, "sub1/file05"),
		path.Join(testDir, "sub1/file06"),
		path.Join(testDir, "sub1/sub2/file07"),
		path.Join(testDir, "sub1/sub2/file08"),
		path.Join(testDir, "sub3/file09"),
		path.Join(testDir, "sub3/file10"),
	}

	// list on empty dir
	items, err := client.ListAll(testDir)
	assert.Nil(t, err)
	assert.Equal(t, []s3fs.ListItem{}, items)

	// write test files
	for _, p := range testFilePaths {
		err := client.Write(&s3fs.File{Path: p, Data: []byte("test")})
		assert.Nil(t, err)
	}
	defer func() {
		err := client.DeleteAll(testDir)
		assert.Nil(t, err)
	}()

	// list test dir
	items, err = client.ListAll(testDir)
	assert.Nil(t, err)
	sortListItems(items)
	if assert.Equal(t, len(testFilePaths), len(items)) {
		for i := 0; i < len(testFilePaths); i++ {
			assert.Equal(t, testFilePaths[i], items[i].Path)
		}
	}

	// list sub dir (sub1)
	items, err = client.ListAll(path.Join(testDir, "sub1"))
	assert.Nil(t, err)
	sortListItems(items)
	assert.Equal(t, []s3fs.ListItem{
		{Path: path.Join(testDir, "sub1/file04")},
		{Path: path.Join(testDir, "sub1/file05")},
		{Path: path.Join(testDir, "sub1/file06")},
		{Path: path.Join(testDir, "sub1/sub2/file07")},
		{Path: path.Join(testDir, "sub1/sub2/file08")},
	}, items)

	// list sub dir (sub1/sub2)
	items, err = client.ListAll(path.Join(testDir, "sub1/sub2"))
	assert.Nil(t, err)
	sortListItems(items)
	assert.Equal(t, []s3fs.ListItem{
		{Path: path.Join(testDir, "sub1/sub2/file07")},
		{Path: path.Join(testDir, "sub1/sub2/file08")},
	}, items)

	// list sub dir (sub3)
	items, err = client.ListAll(path.Join(testDir, "sub3"))
	assert.Nil(t, err)
	sortListItems(items)
	assert.Equal(t, []s3fs.ListItem{
		{Path: path.Join(testDir, "sub3/file09")},
		{Path: path.Join(testDir, "sub3/file10")},
	}, items)

	// list non-existing sub dir (sub4)
	items, err = client.ListAll(path.Join(testDir, "sub4"))
	assert.Nil(t, err)
	assert.Equal(t, []s3fs.ListItem{}, items)
}

func sortListItems(items []s3fs.ListItem) {
	sort.Slice(items, func(i, j int) bool {
		c := strings.Compare(items[i].Path, items[j].Path)
		if c == 0 {
			return items[i].IsDir
		}

		return c < 0
	})
}
