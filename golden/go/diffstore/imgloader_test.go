package diffstore

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.skia.org/infra/go/gcs/test_gcsclient"
	"go.skia.org/infra/go/testutils"
	"go.skia.org/infra/go/testutils/unittest"
	"go.skia.org/infra/golden/go/diff"
	"go.skia.org/infra/golden/go/diffstore/common"
	diffstore_mocks "go.skia.org/infra/golden/go/diffstore/mocks"
	"go.skia.org/infra/golden/go/image/text"
	"go.skia.org/infra/golden/go/types"
)

const (
	gcsImageBaseDir = "dm-images-v1"

	// These digests correspond to the images below and are arbitrarily chosen.
	digest1 = types.Digest("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	digest2 = types.Digest("bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb")
	digest3 = types.Digest("cccccccccccccccccccccccccccccccc")

	// Fake GCS paths to the test images.
	image1GCSPath = gcsImageBaseDir + "/" + string(digest1) + ".png"
	image2GCSPath = gcsImageBaseDir + "/" + string(digest2) + ".png"
	image3GCSPath = gcsImageBaseDir + "/" + string(digest3) + ".png"

	// MD5 hashes of the PNG files.
	image1MD5Hash = "bde6b72edc996515916348e8f4dd406d" // = md5sum(aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.png)
	image2MD5Hash = "96f28080f8cebfdb463bb00724aba779" // = md5sum(bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb.png)
	image3MD5Hash = "fb463a46f01baa8ef0b9d1d87ba6e421" // = md5sum(cccccccccccccccccccccccccccccccc.png)

	skTextImage1 = `! SKTEXTSIMPLE
	1 5
	0x00000000
	0x01000000
	0x00010000
	0x00000100
	0x00000001`

	skTextImage2 = `! SKTEXTSIMPLE
	1 5
	0x01000000
	0x02000000
	0x00020000
	0x00000200
	0x00000002`

	skTextImage3 = `! SKTEXTSIMPLE
	1 5
	0x01000000
	0x03000000
	0x00010000
	0x00000200
	0x00000003`
)

// These images (of type *image.NRGBA) are created from the SKTEXTSIMPLE images defined above, and
// are assumed to be used in a read-only manner throughout the tests.
var image1 = skTextToImage(skTextImage1)
var image2 = skTextToImage(skTextImage2)
var image3 = skTextToImage(skTextImage3)

func TestImageLoaderExpectedMd5HashesAreCorrect(t *testing.T) {
	unittest.SmallTest(t)
	require.Equal(t, bytesToMd5HashString(imageToPng(image1).Bytes()), image1MD5Hash)
	require.Equal(t, bytesToMd5HashString(imageToPng(image2).Bytes()), image2MD5Hash)
	require.Equal(t, bytesToMd5HashString(imageToPng(image3).Bytes()), image3MD5Hash)
}

// Sets up the mock GCSClient and temp folder for images, and returns the test ImageLoader instance.
func setUp(t *testing.T) (*ImageLoader, *test_gcsclient.GCSClient, *diffstore_mocks.FailureStore) {
	// Build mock GCSClient.
	mockBucketClient := test_gcsclient.NewMockClient()

	// Only used for logging errors, which only some tests produce.
	mockBucketClient.On("Bucket").Return("test-bucket").Maybe()

	// Build mock FailureStore.
	mockFailureStore := &diffstore_mocks.FailureStore{}

	// Compute an arbitrary cache size.
	imgCacheCount, _ := getCacheCounts(10)

	// Create the ImageLoader instance.
	imageLoader, err := NewImgLoader(mockBucketClient, mockFailureStore, gcsImageBaseDir, imgCacheCount)
	require.NoError(t, err)

	return imageLoader, mockBucketClient, mockFailureStore
}

func TestImageLoaderGetSingleDigestFoundInBucket(t *testing.T) {
	unittest.SmallTest(t)
	imageLoader, mockClient, mockFailureStore := setUp(t)

	defer mockClient.AssertExpectations(t)
	defer mockFailureStore.AssertExpectations(t)

	// digest1 is present in the GCS bucket.
	expectImageWillBeRead(mockClient, image1GCSPath, image1MD5Hash, image1)

	// Get image.
	images, err := imageLoader.Get(1, types.DigestSlice{digest1})

	// Assert that the correct image was returned.
	require.NoError(t, err)
	require.Len(t, images, 1)
	require.Equal(t, images[0], imageToPng(image1).Bytes())
}

func TestImageLoaderGetSingleDigestNotFound(t *testing.T) {
	unittest.SmallTest(t)
	imageLoader, mockClient, mockFailureStore := setUp(t)

	defer mockClient.AssertExpectations(t)
	defer mockFailureStore.AssertExpectations(t)

	// digest1 is NOT present in the GCS bucket.
	var oa *storage.ObjectAttrs
	mockClient.On("GetFileObjectAttrs", testutils.AnyContext, image1GCSPath).Return(oa, errors.New("not found"))

	// Failure is stored.
	mockFailureStore.On("AddDigestFailure", testutils.AnyContext, diffFailureMatcher(digest1, "http_error")).Return(nil)

	// Get images.
	_, err := imageLoader.Get(1, types.DigestSlice{digest1})

	// Assert that retrieval failed.
	require.Error(t, err)
	require.Contains(t, err.Error(), "Unable to retrieve attributes")
}

func TestImageLoaderGetMultipleDigestsAllFoundInBucket(t *testing.T) {
	unittest.SmallTest(t)
	imageLoader, mockClient, mockFailureStore := setUp(t)

	defer mockClient.AssertExpectations(t)
	defer mockFailureStore.AssertExpectations(t)

	// digest1 and digest2 are present in the GCS bucket.
	expectImageWillBeRead(mockClient, image1GCSPath, image1MD5Hash, image1)
	expectImageWillBeRead(mockClient, image2GCSPath, image2MD5Hash, image2)

	// Get images.
	images, err := imageLoader.Get(1, types.DigestSlice{digest1, digest2})

	// Assert that the correct images were returned.
	require.NoError(t, err)
	require.Len(t, images, 2)
	require.Equal(t, images[0], imageToPng(image1).Bytes())
	require.Equal(t, images[1], imageToPng(image2).Bytes())
}

func TestImageLoaderGetMultipleDigestsDigest1FoundInBucketDigest2NotFound(t *testing.T) {
	unittest.SmallTest(t)
	imageLoader, mockClient, mockFailureStore := setUp(t)

	defer mockClient.AssertExpectations(t)
	defer mockFailureStore.AssertExpectations(t)

	// digest1 is present in the GCS bucket.
	expectImageWillBeRead(mockClient, image1GCSPath, image1MD5Hash, image1)

	// digest2 is NOT present in the GCS bucket.
	var oa2 *storage.ObjectAttrs = nil
	mockClient.On("GetFileObjectAttrs", testutils.AnyContext, image2GCSPath).Return(oa2, errors.New("not found"))

	// Failure is stored.
	mockFailureStore.On("AddDigestFailure", testutils.AnyContext, diffFailureMatcher(digest2, "http_error")).Return(nil)

	// Get images.
	_, err := imageLoader.Get(1, types.DigestSlice{digest1, digest2})

	// Assert that retrieval failed.
	require.Error(t, err)
	require.Contains(t, err.Error(), "Unable to retrieve attributes")
}

func TestImageLoaderWarm(t *testing.T) {
	unittest.SmallTest(t)
	imageLoader, mockClient, mockFailureStore := setUp(t)

	defer mockFailureStore.AssertExpectations(t)

	// digest1 is present in the GCS bucket.
	oa1 := &storage.ObjectAttrs{MD5: md5HashToBytes(image1MD5Hash)}
	mockClient.On("GetFileObjectAttrs", testutils.AnyContext, image1GCSPath).Return(oa1, nil).
		Once() // This ensures that Get doesn't hit GCS after a call to Warm for the same digest.

	// digest1 is read.
	reader1 := ioutil.NopCloser(imageToPng(image1))
	mockClient.On("FileReader", testutils.AnyContext, image1GCSPath).Return(reader1, nil).
		Once() // This ensures that Get doesn't hit GCS after a call to Warm for the same digest.

	// digest2 is present in the GCS bucket.
	oa2 := &storage.ObjectAttrs{MD5: md5HashToBytes(image2MD5Hash)}
	mockClient.On("GetFileObjectAttrs", testutils.AnyContext, image2GCSPath).Return(oa2, nil).
		Once() // This ensures that Get doesn't hit GCS after a call to Warm for the same digest.

	// digest2 is read.
	reader2 := ioutil.NopCloser(imageToPng(image2))
	mockClient.On("FileReader", testutils.AnyContext, image2GCSPath).Return(reader2, nil).
		Once() // This ensures that Get doesn't hit GCS after a call to Warm for the same digest.

	// Fetch both images from GCS and cache them in memory.
	imageLoader.Warm(1, types.DigestSlice{digest1, digest2}, true)

	// Assert that the mocked methods were called as expected.
	mockClient.AssertExpectations(t)

	// Assert that the images are in the cache.
	require.True(t, imageLoader.Contains(digest1))
	require.True(t, imageLoader.Contains(digest2))

	// Get cached images from memory. This shouldn't hit GCS. If it does, the mockClient will panic
	// as per the Once() calls.
	images, err := imageLoader.Get(1, types.DigestSlice{digest1, digest2})

	// Assert that the correct images were returned.
	require.NoError(t, err)
	require.Len(t, images, 2)
	require.Equal(t, images[0], imageToPng(image1).Bytes())
	require.Equal(t, images[1], imageToPng(image2).Bytes())
}

// TODO(lovisolo): Add test cases for multiple digests, and decide what to do about purgeGCS=false.
func TestImageLoaderPurgeImages(t *testing.T) {
	unittest.SmallTest(t)
	imageLoader, mockClient, mockFailureStore := setUp(t)

	defer mockClient.AssertExpectations(t)
	defer mockFailureStore.AssertExpectations(t)

	// digest1 is present in the GCS bucket.
	expectImageWillBeRead(mockClient, image1GCSPath, image1MD5Hash, image1)

	// Fetch digest from GCS and and cache it in memory.
	imageLoader.Warm(1, types.DigestSlice{digest1}, true)

	// Assert that the image is in the cache.
	require.True(t, imageLoader.Contains(digest1))

	// digest1 is deleted.
	mockClient.On("DeleteFile", testutils.AnyContext, image1GCSPath).Return(nil)

	// Purge image.
	err := imageLoader.PurgeImages(types.DigestSlice{digest1}, true)
	require.NoError(t, err)

	// Assert that the image was removed from the cache.
	require.False(t, imageLoader.Contains(digest1))
}

// Decodes an SKTEXTSIMPLE image.
func skTextToImage(s string) *image.NRGBA {
	buf := bytes.NewBufferString(s)
	img, err := text.Decode(buf)
	if err != nil {
		// This indicates an error with the static test data which is initialized before executing the
		// tests, thus we panic instead of asserting the absence of errors with require.NoError.
		panic(fmt.Sprintf("Failed to decode a valid image: %s", err))
	}
	return img.(*image.NRGBA)
}

// Takes an image and returns a PNG-encoded bytes.Buffer.
func imageToPng(image *image.NRGBA) *bytes.Buffer {
	buf := new(bytes.Buffer)
	err := common.EncodeImg(buf, image)
	if err != nil {
		// This indicates an error with the static test data which is initialized before executing the
		// tests, thus we panic instead of asserting the absence of errors with require.NoError.
		panic(fmt.Sprintf("Failed to encode image as PNG: %s", err))
	}
	return buf
}

// Takes a byte slice and returns its MD5 hash as a human-readable string.
func bytesToMd5HashString(bytes []byte) string {
	m := md5.New()
	m.Write(bytes)
	return hex.EncodeToString(m.Sum(nil))
}

// Takes a string with an MD5 hash and encodes it as a byte array.
func md5HashToBytes(md5Hash string) []byte {
	b, err := hex.DecodeString(md5Hash)
	if err != nil {
		// This indicates an error with the static test data which is initialized before executing the
		// tests, thus we panic instead of asserting the absence of errors with require.NoError.
		panic(fmt.Sprintf("Failed to encode digest as MD5 bytes: %s", err))
	}
	return b
}

// diffFailureMatcher is necessary due to the timestamp stored in field DigestFailure.TS.
func diffFailureMatcher(digest types.Digest, reason diff.DiffErr) interface{} {
	return mock.MatchedBy(func(failure *diff.DigestFailure) bool {
		return failure.Digest == digest && failure.Reason == reason
	})
}

func TestGetGSRelPath(t *testing.T) {
	unittest.SmallTest(t)

	digest := types.Digest("098f6bcd4621d373cade4e832627b4f6")
	expectedGSPath := string(digest + ".png")
	gsPath := getGCSRelPath(digest)
	require.Equal(t, expectedGSPath, gsPath)
}

// expectImageWillBeRead adds the mocked expectations for reading the given image with the
// given MD5 hash and will provide the given image data.
func expectImageWillBeRead(mgc *test_gcsclient.GCSClient, gcsPath, hash string, img *image.NRGBA) {
	oa := &storage.ObjectAttrs{MD5: md5HashToBytes(hash)}
	mgc.On("GetFileObjectAttrs", testutils.AnyContext, gcsPath).Return(oa, nil)
	reader := ioutil.NopCloser(imageToPng(img))
	mgc.On("FileReader", testutils.AnyContext, gcsPath).Return(reader, nil)
}
