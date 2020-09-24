// Copyright © 2020 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lighthousehttp_test

import (
	"context"
	"encoding/hex"
	"os"
	"strings"
	"testing"

	client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/lighthousehttp"
	"github.com/stretchr/testify/require"
)

// testValidatorIDProvider implements ValidatorIDProvider.
type testValidatorIDProvider struct {
	index  uint64
	pubKey string
}

func (t *testValidatorIDProvider) Index(ctx context.Context) (uint64, error) {
	return t.index, nil
}
func (t *testValidatorIDProvider) PubKey(ctx context.Context) ([]byte, error) {
	return hex.DecodeString(strings.TrimPrefix(t.pubKey, "0x"))
}

func TestProposerDuties(t *testing.T) {
	tests := []struct {
		name       string
		epoch      uint64
		validators []client.ValidatorIDProvider
		expected   int
	}{
		{
			name:     "Old",
			epoch:    1,
			expected: 32,
		},
		{
			name:     "Recent",
			epoch:    10989,
			expected: 32,
		},
		{
			name:  "GoodWithValidators",
			epoch: 4092,
			validators: []client.ValidatorIDProvider{
				&testValidatorIDProvider{
					index:  16056,
					pubKey: "0x9553a63a58d3a776a2483184e5af37aedf131b82ef1e0bcba7b3c01818f490371aac0c6f9a327fb7eb89190af7b085a5",
				},
				&testValidatorIDProvider{
					index:  35476,
					pubKey: "0x9216091f3e4fe0b0562a6c5bf6e8c35cf0c3b321b6f415de6631d7d12e58603e1e23c8d78f449b601f8d244d26f70aa7",
				},
			},
			expected: 2,
		},
	}

	service, err := lighthousehttp.New(context.Background(),
		lighthousehttp.WithAddress(os.Getenv("LIGHTHOUSEHTTP_ADDRESS")),
		lighthousehttp.WithTimeout(timeout),
	)
	require.NoError(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			duties, err := service.ProposerDuties(context.Background(), test.epoch, test.validators)
			require.NoError(t, err)
			require.NotNil(t, duties)
			require.Equal(t, test.expected, len(duties))
		})
	}
}