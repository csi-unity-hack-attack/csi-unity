package gounity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostById(t *testing.T) {
	ctx, err := newTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.tearDown()

	host, err := ctx.unity.GetHostById("Host_1")
	assert.Nil(t, err)

	assert.Equal(t, "Host_1", host.Id)
}

func TestGetHosts(t *testing.T) {
	ctx, err := newTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.tearDown()

	hosts, err := ctx.unity.GetHosts()
	assert.Nil(t, err)

	assert.Equal(t, 4, len(hosts))
	ids := []string{}
	for _, host := range hosts {
		ids = append(ids, host.Id)
	}
	assert.EqualValues(t, []string{"Host_1", "Host_2", "Host_3", "Host_4"}, ids)
}

func TestAttach(t *testing.T) {
	ctx, err := newTestContext()
	assert.Nil(t, err, "failed to setup rest client to mock server")
	defer ctx.tearDown()

	host, err := ctx.unity.GetHostById("Host_1")
	assert.Nil(t, err)

	lun, err := ctx.unity.GetLunById("sv_1")
	assert.Nil(t, err)

	hlu, err := host.Attach(lun)
	assert.Nil(t, err)

	assert.Equal(t, uint16(1), hlu)
}
