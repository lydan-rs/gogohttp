package response

import (
	"fmt"
	"math"
	// "http-protocol/internal/headers"
	// "strconv"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/require"
)

const LOREM_IPSUM = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum egestas ligula ut ligula vulputate, non malesuada tortor ullamcorper. Mauris pretium fringilla massa non lacinia. Curabitur suscipit tempor metus vel pulvinar. Nam aliquet pulvinar mattis. Proin tempor mauris a lorem vehicula, vitae finibus mi viverra. Donec id urna interdum, aliquam leo sed, congue turpis. Aliquam ac sapien velit. Ut mollis, tellus a viverra accumsan, nunc mi gravida dolor, non lobortis ipsum risus feugiat purus. Aenean vestibulum diam augue, eget tincidunt mi facilisis malesuada. Nam blandit ultricies ipsum sit amet ultrices. Sed pharetra tristique erat eu lacinia. Aliquam suscipit eget purus at dignissim. Nullam bibendum tortor vitae congue tempor.

Cras ante ex, gravida ut convallis sit amet, gravida quis eros. Sed pretium non turpis sit amet fermentum. Sed sit amet neque eros. In orci massa, aliquet egestas interdum ac, feugiat nec lorem. Nunc accumsan felis interdum elit aliquet, vel porta enim bibendum. Aliquam iaculis arcu vel lorem convallis ornare. Nam sodales ipsum lectus, a pulvinar odio sodales at. Vestibulum auctor odio id mollis venenatis. Nunc metus purus, elementum sed feugiat quis, condimentum eu urna. Cras ultricies rutrum mattis. Vestibulum consectetur justo vel ex lobortis, quis rutrum metus rutrum. Suspendisse erat sem, gravida id dignissim rhoncus, condimentum a ante. Sed in pellentesque lectus, ac euismod tellus. Cras id magna accumsan, ullamcorper turpis id, congue libero. Donec vel justo quam.

Ut non nibh et metus iaculis aliquet. Aenean interdum erat tellus, nec malesuada quam tempus a. Nulla at neque efficitur, varius mauris ut, efficitur nisl. Aliquam condimentum turpis eu metus luctus, id sollicitudin diam accumsan. Mauris sit amet leo sit amet sem porta pellentesque. In congue fringilla eros, mattis feugiat quam sollicitudin quis. Aenean laoreet facilisis faucibus. Vivamus in luctus dui. Morbi vel lorem ultricies, finibus sapien tincidunt, gravida libero. Etiam finibus gravida est, egestas accumsan arcu vestibulum non.

Etiam leo felis, facilisis eu augue a, auctor congue nisi. Duis quis pharetra sapien. Quisque pharetra leo elementum faucibus blandit. Donec et nisi ac felis lobortis suscipit id dignissim magna. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Proin vehicula, nibh ac imperdiet finibus, nisl quam congue elit, ac porta sem nunc vitae diam. Etiam sagittis fermentum orci vel posuere. Suspendisse mollis a nisl a ultricies. Suspendisse elementum laoreet tortor, vitae vestibulum urna aliquam vulputate. Nam nisi velit, elementum vitae faucibus ac, elementum elementum arcu. Vestibulum vel felis ac enim ornare maximus. Vivamus non augue diam.

Morbi facilisis arcu ac purus feugiat imperdiet. Fusce ornare sapien nunc, ac suscipit quam luctus at. Cras justo libero, sagittis eu ante molestie, placerat luctus mauris. Etiam venenatis velit non libero sollicitudin eleifend. Etiam ut auctor metus. Aliquam dignissim accumsan libero in pharetra. Phasellus pulvinar eleifend lacus eget blandit. Donec at libero dignissim, convallis massa ut, rutrum enim. Nam nec justo sed urna venenatis tempor. Duis ac nibh sem. Nulla non eros non eros efficitur eleifend. Maecenas congue ipsum ut tempus scelerisque.
`

func TestWriteChunks(t *testing.T) {
	var buffer bytes.Buffer

	chunkSize := 32
	n, err := writeChunks(&buffer, []byte(LOREM_IPSUM), uint(32))
	assert.NoError(t, err)
	assert.Equal(t, n, buffer.Len())
	assert.True(t, n > len(LOREM_IPSUM))
	fmt.Printf("len(LOREM_IPSUM): %d\n", len(LOREM_IPSUM))
	lenLorem := len(LOREM_IPSUM)
	nLines := int(math.Ceil(float64(lenLorem)/float64(chunkSize)))
	fmt.Printf("Lines: %d\n", nLines)
	lowerN := lenLorem + (nLines*5)
	upperN := lenLorem + nLines*(4+len(fmt.Sprintf("%x", chunkSize)))
	assert.True(t, (lowerN<=n)&&(n <= upperN), fmt.Sprintf("Condition: %d<=%d<=%d.", lowerN, n, upperN)) 
	fmt.Printf(buffer.String())
}
