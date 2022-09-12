// package anki_forvo_plugin
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileDownloader_Download(t *testing.T) {
	url := "https://apifree.forvo.com/audio/2i351b1f3o1b341j2c3426242c3239243o343p3h3k313e2m3l2g241i2e1l3329382q2e231p243a3g2m1k31223o3l1l2f371b2p2a1j3i3628381h3i2i273c333l3o342c2b2d3g2e1b292m1j3f3a2c1o2p2q1m1i3i223n1t1t_1b2d3n242o2o313k3d3n2j231f2q2i3f2c3c1h26292h1t1t"
	filepath := "test_output/dog_狗.mp3"
	d := NewFileDownloader()
	err := d.Download(url, filepath)
	assert.Nil(t, err)
}
