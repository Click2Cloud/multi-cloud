package tests

import (
	"os"
	"testing"

	"github.com/opensds/multi-cloud/s3/pkg/helper"
	//log "github.com/sirupsen/logrus"
	bus "github.com/opensds/multi-cloud/s3/pkg/messagebus"
	. "gopkg.in/check.v1"
	"log"
)

func Test(t *testing.T) { TestingT(t) }

type MessageBusTestSuite struct {
	logger *log.Logger
	f      *os.File
}

var _ = Suite(&MessageBusTestSuite{})

func (mbs *MessageBusTestSuite) SetUpSuite(c *C) {
	var err error
	mbs.f, err = os.OpenFile("./test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	c.Assert(err, Equals, nil)
	helper.CONFIG.LogLevel = 20
	mbs.logger = log.New(mbs.f, "[yig]", log.LstdFlags, helper.CONFIG.LogLevel)
	helper.Logger = mbs.logger
}

func (mbs *MessageBusTestSuite) TearDownSuite(c *C) {
	mbs.f.Close()
}

func (mbs *MessageBusTestSuite) SetUpTest(c *C) {
	helper.CONFIG.MsgBus.Type = 0
	helper.CONFIG.MsgBus.RequestTimeoutMs = 3000
	helper.CONFIG.MsgBus.MessageTimeoutMs = 5000
	helper.CONFIG.MsgBus.SendMaxRetries = 2
	helper.CONFIG.MsgBus.Server = make(map[string]interface{})
	bus.ClearInit()
}

func (mbs *MessageBusTestSuite) TearDownTest(c *C) {
}