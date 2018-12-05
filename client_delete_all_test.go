package s3fs_test

import (
	"path"
	"testing"

	"github.com/d5/go-s3fs"
	"github.com/stretchr/testify/assert"
)

func TestClient_DeleteAll(t *testing.T) {
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

	// delete on empty dir
	err := client.DeleteAll(testDir)
	assert.Nil(t, err)

	// write test files
	for _, p := range testFilePaths {
		err := client.Write(&s3fs.File{Path: p, Data: []byte("test")})
		assert.Nil(t, err)
	}
	defer func() {
		err := client.DeleteAll(testDir)
		assert.Nil(t, err)
	}()
	assertEqualList(t, client, testDir, []s3fs.ListItem{
		{Path: path.Join(testDir, "file01")},
		{Path: path.Join(testDir, "file02")},
		{Path: path.Join(testDir, "file03")},
		{Path: path.Join(testDir, "sub1/file04")},
		{Path: path.Join(testDir, "sub1/file05")},
		{Path: path.Join(testDir, "sub1/file06")},
		{Path: path.Join(testDir, "sub1/sub2/file07")},
		{Path: path.Join(testDir, "sub1/sub2/file08")},
		{Path: path.Join(testDir, "sub3/file09")},
		{Path: path.Join(testDir, "sub3/file10")},
	})

	// delete sub3
	err = client.DeleteAll(path.Join(testDir, "sub3"))
	assert.Nil(t, err)
	assertEqualList(t, client, testDir, []s3fs.ListItem{
		{Path: path.Join(testDir, "file01")},
		{Path: path.Join(testDir, "file02")},
		{Path: path.Join(testDir, "file03")},
		{Path: path.Join(testDir, "sub1/file04")},
		{Path: path.Join(testDir, "sub1/file05")},
		{Path: path.Join(testDir, "sub1/file06")},
		{Path: path.Join(testDir, "sub1/sub2/file07")},
		{Path: path.Join(testDir, "sub1/sub2/file08")},
	})

	// delete sub1/sub2/file08
	err = client.DeleteAll(path.Join(testDir, "sub1/sub2/file08"))
	assert.Nil(t, err)
	assertEqualList(t, client, testDir, []s3fs.ListItem{
		{Path: path.Join(testDir, "file01")},
		{Path: path.Join(testDir, "file02")},
		{Path: path.Join(testDir, "file03")},
		{Path: path.Join(testDir, "sub1/file04")},
		{Path: path.Join(testDir, "sub1/file05")},
		{Path: path.Join(testDir, "sub1/file06")},
		{Path: path.Join(testDir, "sub1/sub2/file07")},
	})

	// delete sub1/sub2
	err = client.DeleteAll(path.Join(testDir, "sub1/sub2"))
	assert.Nil(t, err)
	assertEqualList(t, client, testDir, []s3fs.ListItem{
		{Path: path.Join(testDir, "file01")},
		{Path: path.Join(testDir, "file02")},
		{Path: path.Join(testDir, "file03")},
		{Path: path.Join(testDir, "sub1/file04")},
		{Path: path.Join(testDir, "sub1/file05")},
		{Path: path.Join(testDir, "sub1/file06")},
	})

	// delete sub1
	err = client.DeleteAll(path.Join(testDir, "sub1"))
	assert.Nil(t, err)
	assertEqualList(t, client, testDir, []s3fs.ListItem{
		{Path: path.Join(testDir, "file01")},
		{Path: path.Join(testDir, "file02")},
		{Path: path.Join(testDir, "file03")},
	})

	// delete /
	err = client.DeleteAll(testDir)
	assert.Nil(t, err)
	assertEqualList(t, client, testDir, []s3fs.ListItem{})

}

func assertEqualList(t *testing.T, client *s3fs.Client, pathPrefix string, expectedItems []s3fs.ListItem) {
	items, err := client.ListAll(pathPrefix)
	if !assert.Nil(t, err) {
		return
	}

	sortListItems(items)
	if assert.Equal(t, len(expectedItems), len(items)) {
		for i := 0; i < len(expectedItems); i++ {
			assert.Equal(t, expectedItems[i], items[i])
		}
	}
}
