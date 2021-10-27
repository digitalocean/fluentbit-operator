package copy

import (
	"os"
	"sort"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

func TestCopyFilesWithFilter(t *testing.T) {

	tests := []struct {
		name     string
		filterFn func(info os.FileInfo) bool
		expected []string
	}{
		{
			name: "all",
			filterFn: func(info os.FileInfo) bool {
				return true
			},
			expected: []string{"f0.conf", "f1.txt", "f2.conf"},
		},
		{
			name: "none",
			filterFn: func(info os.FileInfo) bool {
				return false
			},
			expected: []string{},
		},
		{
			name: "conf",
			filterFn: func(info os.FileInfo) bool {
				return strings.HasSuffix(info.Name(), ".conf")
			},
			expected: []string{"f0.conf", "f2.conf"},
		},
		{
			name: "txt",
			filterFn: func(info os.FileInfo) bool {
				return strings.HasSuffix(info.Name(), ".txt")
			},
			expected: []string{"f1.txt"},
		},
		{
			name: "gz",
			filterFn: func(info os.FileInfo) bool {
				return strings.HasSuffix(info.Name(), ".gz")
			},
			expected: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g := NewGomegaWithT(t)

			t.Cleanup(func() {
				_ = os.RemoveAll("testdata/dest")
			})

			err := CopyFilesWithFilter("testdata/src", "testdata/dest", test.filterFn)

			g.Expect(err).NotTo(HaveOccurred())

			f, err := os.Open("testdata/dest")
			g.Expect(err).NotTo(HaveOccurred())
			defer f.Close()

			fis, err := f.Readdir(-1)
			g.Expect(err).NotTo(HaveOccurred())

			actual := make([]string, 0, len(fis))
			for _, fi := range fis {
				name := fi.Name()

				actual = append(actual, name)
			}

			sort.Strings(actual)

			g.Expect(actual).To(Equal(test.expected))
		})
	}
}
