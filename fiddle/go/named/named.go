// named is a utility for dealing with named fiddles.
package named

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/skia-dev/glog"
	"go.skia.org/infra/fiddle/go/store"
)

var (
	// fiddleHashRe is used to validate fiddle hashes.
	fiddleHashRe = regexp.MustCompile("^[0-9a-zA-Z]{32}$")

	// fiddlNameRe is used to validate fiddle names.
	fiddleNameRe = regexp.MustCompile("^[0-9a-zA-Z_]+$")

	// trailingToMedia maps the end of each image URL to the store.Media type
	// that it corresponds to.
	trailingToMedia = map[string]store.Media{
		"_raster.png": store.CPU,
		"_gpu.png":    store.GPU,
		".pdf":        store.PDF,
		".skp":        store.SKP,
	}
)

// NameStore is an interface that store.Store conforms to that is just the
// methods that Named uses.
type NameStore interface {
	GetHashFromName(name string) (string, error)
	WriteName(name, hash, user string) error
}

// Named deals with creating and dereferencing named fiddles.
type Named struct {
	// cache maps fiddle names to fiddle hashes.
	cache map[string]string

	// mutex protects the cache.
	mutex sync.Mutex

	st NameStore
}

// New creates a new Named.
func New(st NameStore) *Named {
	return &Named{
		cache: map[string]string{},
		st:    st,
	}
}

// Add a named fiddle.
//
//   name - The name of the fidde, w/o the @ prefix.
//   hash - The fiddle hash.
//   user - The email of the user that created the name.
func (n *Named) Add(name, hash, user string) error {
	if !fiddleNameRe.MatchString(name) {
		return fmt.Errorf("Not a valid fiddle name %q", name)
	}
	if !fiddleHashRe.MatchString(hash) {
		return fmt.Errorf("Not a valid fiddle hash  %q", hash)
	}
	oldHash, err := n.DereferenceID("@" + name)
	if err == nil {
		if oldHash == hash {
			// Don't bother writing if the hash is already correct.
			return nil
		}
		glog.Infof("Named Fiddle Changed: %s %s -> %s by %s", name, oldHash, hash, user)
	} else {
		glog.Infof("Named Fiddle Created: %s %s by %s", name, hash, user)
	}
	if err := n.st.WriteName(name, hash, user); err != nil {
		return fmt.Errorf("Failed to write name: %s", err)
	}
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.cache[name] = hash
	return nil
}

// Dereference converts the id to a fiddlehash, where id could
// be either a fiddle name or a fiddle hash. Fiddle names are
// presumed to be prefixed with "@".
//
// Returns an error if the fiddle hash is not valid looking.
func (n *Named) DereferenceID(id string) (string, error) {
	if len(id) < 2 {
		return "", fmt.Errorf("Invalid ID")
	}
	if id[0] == '@' {
		id = id[1:]
		n.mutex.Lock()
		fiddleHash, ok := n.cache[id]
		n.mutex.Unlock()
		if !ok {
			// Look up the name in gs://skia-fiddle/named/<id>.
			var err error
			fiddleHash, err = n.st.GetHashFromName(id)
			if err != nil {
				return "", fmt.Errorf("Unknown name: %s", err)
			}
			if !fiddleHashRe.MatchString(fiddleHash) {
				return "", fmt.Errorf("Not a valid fiddle hash found in named file %q: %s", id, fiddleHash)
			}
			n.mutex.Lock()
			defer n.mutex.Unlock()
			// If found, add to cache.
			n.cache[id] = fiddleHash
		}
		return fiddleHash, nil
	} else {
		// match against regex?
		if !fiddleHashRe.MatchString(id) {
			return "", fmt.Errorf("Not a valid fiddle hash: %q", id)
		}
		return id, nil
	}
}

// DereferenceImageID converts the id of an image to the fiddle hash and the
// Media it represents, handling both original fiddle hashes and fiddle names
// as part of the image id. Fiddle names are
// presumed to be prefixed with "@".
//
// I.e. it converts "cbb8dee39e9f1576cd97c2d504db8eee_gpu.png" to
// "cbb8dee39e9f1576cd97c2d504db8eee" and store.GPU.  and, could also convert
// "@star_raster.png" to "cbb8dee39e9f1576cd97c2d504db8eee" and store.CPU.
func (n *Named) DereferenceImageID(id string) (string, store.Media, error) {
	// First strip off the "_raster.png" or ".pdf" from the end.
	trailing := filepath.Ext(id)
	if len(trailing) < 4 {
		return "", store.UNKNOWN, fmt.Errorf("Not a valid image id: %q", id)
	}
	id = id[:len(id)-len(trailing)]
	// If this is a .png then we need to strip off the trailing "_raster" or "_gpu".
	if trailing == ".png" {
		parts := strings.Split(id, "_")
		if len(parts) < 2 {
			return "", store.UNKNOWN, fmt.Errorf("Not a valid image id form: %q", id)
		}
		trailing = "_" + parts[len(parts)-1] + trailing
		id = strings.Join(parts[:len(parts)-1], "_")
	}
	media, ok := trailingToMedia[trailing]
	if !ok {
		return "", store.UNKNOWN, fmt.Errorf("Unknown media: %q", trailing)
	}
	// We are left with just the name or fiddle hash, dereference that.
	fiddleHash, err := n.DereferenceID(id)
	return fiddleHash, media, err
}
