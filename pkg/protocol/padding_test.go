// Copyright (C) 2024  mieru authors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>

package protocol

import (
	crand "crypto/rand"
	mrand "math/rand"
	"testing"

	"github.com/enfein/mieru/v3/pkg/common"
	"github.com/enfein/mieru/v3/pkg/rng"
)

func TestMaxPaddingSize(t *testing.T) {
	testcases := []struct {
		mtu                 int
		transport           common.TransportProtocol
		fragmentSize        int
		existingPaddingSize int
		want                int
	}{
		{
			1400,
			common.StreamTransport,
			maxPDU,
			255,
			255,
		},
		{
			1400,
			common.PacketTransport,
			1400 - packetOverhead - 16,
			12,
			4,
		},
		{
			1400,
			common.UnknownTransport,
			0,
			255,
			255,
		},
	}
	for _, tc := range testcases {
		got := MaxPaddingSize(tc.mtu, tc.transport, tc.fragmentSize, tc.existingPaddingSize)
		if got != tc.want {
			t.Errorf("MaxPaddingSize() = %d, want %d", got, tc.want)
		}
	}
}

func TestNewASCIIPadding(t *testing.T) {
	maxPaddingLen := rng.Intn(256)
	minConsecutiveASCIILen := rng.IntRange(0, maxPaddingLen+1)
	padding := newPadding(paddingOpts{
		maxLen: maxPaddingLen,
		ascii: &asciiPaddingOpts{
			minConsecutiveASCIILen: minConsecutiveASCIILen,
		},
	})

	if len(padding) > maxPaddingLen {
		t.Errorf("Padding length %d is bigger than the maximum length %d", len(padding), maxPaddingLen)
	}
	if len(padding) < minConsecutiveASCIILen {
		t.Errorf("Padding length %d is smaller than the minimum length %d", len(padding), minConsecutiveASCIILen)
	}
	consecutive := common.MaxConsecutivePrintableLength(padding)
	if consecutive < minConsecutiveASCIILen {
		t.Errorf("Padding's consecutive printable length %d is smaller than the required minimum consecutive ASCII length %d", consecutive, minConsecutiveASCIILen)
	}
}

func TestNewEntropyPadding(t *testing.T) {
	var meet, miss int
	for i := 0; i < 1000; i++ {
		existingDataLen := rng.IntRange(0, 1024)
		existingData := make([]byte, existingDataLen)
		if _, err := crand.Read(existingData); err != nil {
			t.Fatalf("Failed to generate random data: %v", err)
		}
		maxPaddingLen := rng.Intn(256)
		targetPercentage := mrand.Intn(50) + 1
		targetProbability := float64(targetPercentage) / 100.0
		padding := newPadding(paddingOpts{
			maxLen: maxPaddingLen,
			entropy: &entropyPaddingOpts{
				targetProbability: targetProbability,
				existingData:      existingData,
			},
		})

		if len(padding) > maxPaddingLen {
			t.Errorf("Padding length %d is bigger than the maximum length %d", len(padding), maxPaddingLen)
		}
		combined := append(existingData, padding...)
		bitDistribution := common.ToBitDistribution(combined)
		if bitDistribution.Bit0 > targetProbability && bitDistribution.Bit1 > targetProbability {
			// Doesn't meet the target.
			// In this case, the padding bits are either all 0 or all 1.
			if !common.IsBitsAllZero(padding) && !common.IsBitsAllOne(padding) {
				t.Errorf("Can't meet target probability %f, but padding %v is not pure 0 bit or 1 bit", targetProbability, padding)
			} else {
				miss++
			}
		} else {
			meet++
		}
	}
	t.Logf("%d paddings meet the target probability; %d paddings miss", meet, miss)
}
