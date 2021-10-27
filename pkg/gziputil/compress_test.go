package gziputil

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"
)

func TestIsCompressed(t *testing.T) {

	g := NewGomegaWithT(t)

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "compressed",
			input:    "testdata/fb.conf.gz",
			expected: true,
		},
		{
			name:     "uncompressed",
			input:    "testdata/fb.conf",
			expected: false,
		},
	}

	for _, test := range tests {
		actual, err := IsCompressed(test.input)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(actual).To(Equal(test.expected))
	}
}

func TestDecompress(t *testing.T) {

	g := NewGomegaWithT(t)

	newFile := "testdata/fb.conf.new"

	t.Cleanup(func() {
		_ = os.Remove(newFile)
	})

	expected, err := os.ReadFile("testdata/fb.conf")
	g.Expect(err).NotTo(HaveOccurred())

	err = Decompress("testdata/fb.conf.gz", newFile)
	g.Expect(err).NotTo(HaveOccurred())

	actual, err := os.ReadFile(newFile)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(actual).To(Equal(expected))
}

func TestCompress(t *testing.T) {

	g := NewGomegaWithT(t)

	newFile := "testdata/fb.conf.compress.test.gz"
	decFile := "testdata/fb.de.compress.test.conf"

	t.Cleanup(func() {
		_ = os.Remove(newFile)
		_ = os.Remove(decFile)
	})

	data, err := os.ReadFile("testdata/fb.conf")
	g.Expect(err).NotTo(HaveOccurred())

	compressed, err := CompressBytes(data)
	g.Expect(err).NotTo(HaveOccurred())

	err = os.WriteFile(newFile, compressed, 0644)
	g.Expect(err).NotTo(HaveOccurred())

	isCompressed, err := IsCompressed(newFile)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(isCompressed).To(BeTrue())

	err = Decompress(newFile, decFile)
	g.Expect(err).NotTo(HaveOccurred())

	decompressed, err := os.ReadFile(decFile)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(decompressed).To(Equal(data))
}
