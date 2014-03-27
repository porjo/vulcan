package roundrobin

import (
	timetools "github.com/mailgun/gotools-time"
	. "github.com/mailgun/vulcan/endpoint"
	. "launchpad.net/gocheck"
	"testing"
	"time"
)

func Test(t *testing.T) { TestingT(t) }

type RoundRobinSuite struct {
	tm *timetools.FreezedTime
}

var _ = Suite(&RoundRobinSuite{})

func (s *RoundRobinSuite) SetUpSuite(c *C) {
	s.tm = &timetools.FreezedTime{
		CurrentTime: time.Date(2012, 3, 4, 5, 6, 7, 0, time.UTC),
	}
}

func (s *RoundRobinSuite) newRR() *RoundRobin {
	r, err := NewRoundRobinWithOptions(s.tm, nil)
	if err != nil {
		panic(err)
	}
	return r
}

func (s *RoundRobinSuite) TestNoEndpoints(c *C) {
	r := s.newRR()
	_, err := r.NextEndpoint(nil)
	c.Assert(err, NotNil)
}

// Subsequent calls to load balancer with 1 endpoint are ok
func (s *RoundRobinSuite) TestSingleEndpoint(c *C) {
	r := s.newRR()

	u := MustParseUrl("http://localhost:5000")
	r.AddEndpoint(u)

	u2, err := r.NextEndpoint(nil)
	c.Assert(err, IsNil)
	c.Assert(u2, Equals, u)

	u3, err := r.NextEndpoint(nil)
	c.Assert(err, IsNil)
	c.Assert(u3, Equals, u)
}

// Make sure that load balancer round robins requests
func (s *RoundRobinSuite) TestMultipleEndpoints(c *C) {
	r := s.newRR()

	uA := MustParseUrl("http://localhost:5000")
	uB := MustParseUrl("http://localhost:5001")
	r.AddEndpoint(uA)
	r.AddEndpoint(uB)

	u, err := r.NextEndpoint(nil)
	c.Assert(err, IsNil)
	c.Assert(u, Equals, uA)

	u, err = r.NextEndpoint(nil)
	c.Assert(err, IsNil)
	c.Assert(u, Equals, uB)

	u, err = r.NextEndpoint(nil)
	c.Assert(err, IsNil)
	c.Assert(u, Equals, uA)
}

// Make sure that adding endpoints during load balancing works fine
func (s *RoundRobinSuite) TestAddEndpoints(c *C) {
	r := s.newRR()

	uA := MustParseUrl("http://localhost:5000")
	uB := MustParseUrl("http://localhost:5001")
	r.AddEndpoint(uA)

	u, err := r.NextEndpoint(nil)
	c.Assert(err, IsNil)
	c.Assert(u, Equals, uA)

	r.AddEndpoint(uB)

	// index was reset after altering endpoints
	u, err = r.NextEndpoint(nil)
	c.Assert(err, IsNil)
	c.Assert(u, Equals, uA)

	u, err = r.NextEndpoint(nil)
	c.Assert(err, IsNil)
	c.Assert(u, Equals, uB)
}

// Removing endpoints from the load balancer works fine as well
func (s *RoundRobinSuite) TestRemoveEndpoint(c *C) {
	r := s.newRR()

	uA := MustParseUrl("http://localhost:5000")
	uB := MustParseUrl("http://localhost:5001")
	r.AddEndpoint(uA)
	r.AddEndpoint(uB)

	u, err := r.NextEndpoint(nil)
	c.Assert(err, IsNil)
	c.Assert(u, Equals, uA)

	// Removing endpoint resets the counter
	r.RemoveEndpoint(uB)

	u, err = r.NextEndpoint(nil)
	c.Assert(err, IsNil)
	c.Assert(u, Equals, uA)
}

// Removing endpoints from the load balancer works fine as well
func (s *RoundRobinSuite) TestRemoveMultipleEndpoints(c *C) {
	r := s.newRR()

	uA := MustParseUrl("http://localhost:5000")
	uB := MustParseUrl("http://localhost:5001")
	uC := MustParseUrl("http://localhost:5002")
	r.AddEndpoint(uA)
	r.AddEndpoint(uB)
	r.AddEndpoint(uC)

	u, err := r.NextEndpoint(nil)
	c.Assert(err, IsNil)
	u, err = r.NextEndpoint(nil)
	c.Assert(err, IsNil)
	u, err = r.NextEndpoint(nil)
	c.Assert(err, IsNil)
	c.Assert(u, Equals, uC)

	// There's only one endpoint left
	r.RemoveEndpoint(uA)
	r.RemoveEndpoint(uB)
	u, err = r.NextEndpoint(nil)
	c.Assert(err, IsNil)
	c.Assert(u, Equals, uC)
}
