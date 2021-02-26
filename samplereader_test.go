package samplereader_test

import (
	"bytes"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang.org/x/sync/errgroup"

	"github.com/neilotoole/samplereader"
)

func TestSampleReader(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/sample1.csv")
	require.NoError(t, err)

	br := bytes.NewReader(b)
	sr := samplereader.New(br)
	require.NoError(t, err)

	r, err := sr.NewReader()
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, r.Close())
	}()

	gotB, err := ioutil.ReadAll(r)
	require.NoError(t, err)

	require.Equal(t, b, gotB)

}

func TestSampleReader2(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/sample1.csv")
	require.NoError(t, err)

	sr := samplereader.New(bytes.NewReader(b))
	require.NoError(t, err)

	g := &errgroup.Group{}
	for i := 0; i < 1000; i++ {
		i := i
		r, err := sr.NewReader()
		require.NoError(t, err)
		g.Go(func() error {
			// Add some jitter
			time.Sleep(time.Nanosecond * time.Duration(rand.Intn(300)))

			gotB, err := ioutil.ReadAll(r)
			if err != nil {
				assert.NoError(t, r.Close())
				return err
			}

			if bytes.Equal(b, gotB) {
				t.Logf("goroutine %d: SUCCESS", i)
				return nil
			}
			assert.Equal(t, b, gotB)
			return nil
		})
	}
	sr.Seal()

	err = g.Wait()
	assert.NoError(t, err)
}

func TestGenerateSample(t *testing.T) {
	generateSampleData(t)
}

func generateSampleData(t *testing.T) {
	f, err := os.Create(filepath.Join("testdata", "sample1.csv"))
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	const line = "A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z\n"
	for i := 0; i < 10000; i++ {
		_, err = f.WriteString(strconv.Itoa(i) + "," + line)
		if err != nil {
			t.Fatal(err)
		}
	}
}
