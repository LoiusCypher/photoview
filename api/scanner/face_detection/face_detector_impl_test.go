package face_detection

import (
	// "fmt"
	// "os"
	// "path/filepath"
	"math"
	"math/rand/v2"
	// "slices"
	// "strings"
	"testing"

	"github.com/Kagami/go-face"
	// "github.com/google/go-cmp/cmp"
	// "github.com/loiuscypher/photoview/api/graphql/models"
	//"github.com/loiuscypher/photoview/api/scanner/face_detection"
	"github.com/loiuscypher/photoview/api/test_utils"
	// scanner_utils "github.com/loiuscypher/photoview/api/test_utils/scanner"
)

func TestMain(m *testing.M) {
	test_utils.IntegrationTestRun(m)
}

func FillDescPattern(desc *face.Descriptor, val float32) {
	for j := range len(desc) {
		desc[j] = val
	}
}

func TestFaceDetectionImplCalcDist(t *testing.T) {
	type testCaseDist struct {
		desc1  face.Descriptor
		desc2  face.Descriptor
		expect float64
	}
	var zeroesDescriptor face.Descriptor
	var onesDescriptor face.Descriptor
	FillDescPattern(&onesDescriptor, 1.0)
	var twosDescriptor face.Descriptor
	FillDescPattern(&twosDescriptor, 2.0)
	randVal1 := rand.Float32()
	rand1Descriptor := zeroesDescriptor
	rand1Descriptor[rand.IntN(len(rand1Descriptor)-1)] = randVal1
	randVal2 := rand.Float32()
	rand2Descriptor := rand1Descriptor
	rand2Descriptor[len(rand2Descriptor)-1] = randVal2
	testCases := []testCaseDist{
		{zeroesDescriptor, zeroesDescriptor, 0.0},
		{zeroesDescriptor, onesDescriptor, float64(len(onesDescriptor))},
		{onesDescriptor, twosDescriptor, float64(len(zeroesDescriptor))},
		{zeroesDescriptor, rand1Descriptor, math.Pow(float64(randVal1), 2)},
		{zeroesDescriptor, rand2Descriptor, math.Pow(float64(randVal1), 2) + math.Pow(float64(randVal2), 2)},
	}
	for i, testCase := range testCases {
		dist := euclidean_distance(testCase.desc1, testCase.desc2)
		if dist != testCase.expect {
			t.Errorf("calculate euclidean distance %d: %f instead of %f", i, dist, testCase.expect)
		}
	}
}

func TestFaceDetectionImplMinDistIndexToIndexGroup(t *testing.T) {
	var zeroesDescriptor face.Descriptor
	var onesDescriptor face.Descriptor
	FillDescPattern(&onesDescriptor, 1.0)
	var twosDescriptor face.Descriptor
	FillDescPattern(&twosDescriptor, 2.0)
	randVal1 := rand.Float32() * 0.5
	rand1Descriptor := zeroesDescriptor
	rand1Descriptor[rand.IntN(len(rand1Descriptor)-1)] = randVal1
	randVal2 := rand.Float32()*0.5 + 0.5
	rand2Descriptor := rand1Descriptor
	rand2Descriptor[rand.IntN(len(rand1Descriptor)-1)] = randVal2

	fd := &faceDetector{
		faceDescriptors: []face.Descriptor{zeroesDescriptor, onesDescriptor, twosDescriptor, rand1Descriptor, rand2Descriptor},
	}

	type testCaseDist struct {
		index  int
		idxGrp []int
		expect float64
		expIdx int
	}
	testCases := []testCaseDist{
		{0, []int{}, 0.0, 0},
		{0, []int{1, 2, 3, 4}, math.Pow(float64(randVal1), 2), 3},
	}
	for i, testCase := range testCases {
		dist, distIdx := fd.calcMinDistIndex(testCase.index, testCase.idxGrp)
		if distIdx != testCase.expIdx {
			t.Errorf("calcMinDistIndex %d: Idx %d instead of %d", i, distIdx, testCase.expIdx)
		}
		if dist != testCase.expect {
			t.Errorf("calcMinDistIndex %d: dist %f instead of %f", i, dist, testCase.expect)
		}
	}
}

func TestFaceDetectionImplMinDistGroups(t *testing.T) {
	var zeroesDescriptor face.Descriptor
	var onesDescriptor face.Descriptor
	FillDescPattern(&onesDescriptor, 1.0)
	var twosDescriptor face.Descriptor
	FillDescPattern(&twosDescriptor, 2.0)
	randVal1 := rand.Float32() * 0.5
	rand1Descriptor := zeroesDescriptor
	rand1Descriptor[rand.IntN(len(rand1Descriptor)-1)] = randVal1
	randVal2 := rand.Float32()*0.5 + 0.5
	rand2Descriptor := rand1Descriptor
	rand2Descriptor[rand.IntN(len(rand1Descriptor)-1)] = randVal2

	fd := &faceDetector{
		faceDescriptors: []face.Descriptor{
			zeroesDescriptor,
			onesDescriptor,
			twosDescriptor,
			rand1Descriptor,
			rand2Descriptor},
	}

	type testCaseDist struct {
		idxGrp1 []int
		idxGrp2 []int
		expect  float64
	}
	testCases := []testCaseDist{
		{[]int{}, []int{}, 0.0},
		{[]int{0, 1}, []int{2, 3, 4}, math.Pow(float64(randVal1), 2)},
	}
	for i, testCase := range testCases {
		dist := fd.calcMinDistGroups(testCase.idxGrp1, testCase.idxGrp2)
		if dist != testCase.expect {
			t.Errorf("calcMinDistIndex %d: dist %f instead of %f", i, dist, testCase.expect)
		}
	}
}
