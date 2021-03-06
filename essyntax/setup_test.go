package essyntax

import (
	"testing"
	"time"

	. "github.com/mattbaird/elastigo/lib"
)

/*
// elastigo Conn adapter to avoid a circular dependency
type conn interface {
	CreateIndex(name string) (interface{}, error)
	DeleteIndex(name string) (interface{}, error)

	Index(index string, _type string, id string, args map[string]interface{}, data interface{}) (interface{}, error)
}
*/

func newIndexWorker(c *Conn, t *testing.T) func(interface{}) {

	return func(d interface{}) {
		_, err := c.Index("oilers", "heyday", "", nil, d)
		if err != nil {
			t.Fatalf("Index failed: %s", err)
		}
	}
}

func PopulateTestDB(t *testing.T, c *Conn) {

	// it is not technically necessary to create an index here
	_, err := c.CreateIndex("oilers")
	if err != nil {
		t.Fatal("Error in CreateIndex", err)
	}

	// set the mapping for dob to be a date so it can be used for range searches
	_, err = c.DoCommand("PUT", "/oilers/heyday/_mapping?ignore_conflicts", nil,
		string(`{"heyday": {"properties": {
					"dob": {"type": "date", "format": "basic_date"},
					"pos": {"type": "string", "index": "not_analyzed"},
					"teams": {"type": "string", "index": "not_analyzed"}
				}}}`))
	if err != nil {
		t.Fatal("Error setting dob mapping", err)
	}

	idx := newIndexWorker(c, t)

	idx(`{"name": "Mark Messier",   "jersey": 11, "pos": "LW", "goals": 37, "PIM": 165, 
			"dob": "19610118", "teams": ["EDM", "NYR", "VAN"],
			"quote": "Really the team often will take on the personality of its coach."}`)
	idx(`{"name": "Wayne Gretzky",  "jersey": 99, "pos": "C",  "goals": 87,
			"dob": "19610126", "teams": ["EDM", "NYR", "STL"],
			"quote": "I skate to where the puck is going to be, not to where it has been."}`)
	idx(`{"name": "Paul Coffey",    "jersey": 7,  "pos": "D",  "goals": 40,
			"dob": "19610601", "teams": ["EDM", "DET"],
			"quote": "When we've got the puck, they can't score."}`)
	idx(`{"name": "Jari Kurri",     "jersey": 17, "pos": "RW", "goals": 52,
			"dob": "19600518", "teams": ["EDM", "VAN"],
			"quote": "My teammates make excellent decisions with the puck."}`)
	idx(`{"name": "Glenn Anderson", "jersey": 9,  "pos": "RW", "goals": 54,
			"dob": "19601002", "teams": ["EDM", "NYR", "TOR", "STL"],
			"quote": "The key is to hold the stick properly."}`)
	idx(`{"name": "Ken Linseman",   "jersey": 13, "pos": "C",  "goals": 18,
			"dob": "19580811", "teams": ["EDM", "TOR"]}`)
	idx(`{"name": "Pat Hughes",     "jersey": 16, "pos": "RW", "goals": 27,
			"dob": "19550325", "teams": ["EDM", "MTL", "PIT"]
			"quote": "Teams are fun with pucks."}`)
	idx(`{"name": "Dave Hunter",    "jersey": 12, "pos": "LW", "goals": 22,
			"dob": "19580101", "teams": ["EDM", "PIT"]}`)
	idx(`{"name": "Kevin Lowe",     "jersey": 4,  "pos": "D",  "goals": 4,
			"dob": "19590415", "teams": ["EDM", "NYR"],
			"quote": "Killer instinct is putting teams away and not making mistakes."}`)
	idx(`{"name": "Charlie Huddy",  "jersey": 22, "pos": "D",  "goals": 8,
			"dob": "19590602", "teams": ["EDM", "BUF", "STL"]}`)
	idx(`{"name": "Randy Gregg",    "jersey": 21, "pos": "D",  "goals": 13,
			"dob": "19560219", "teams": ["EDM", "VAN"],
			"quote": "Teaming up with Teammates makes a good TEAM."}`)
	idx(`{"name": "Dave Semenko",   "jersey": 27, "pos": "LW", "goals": 4, "PIM": 118,
			"dob": "19570712", "teams": ["EDM"],
			"quote": "Over here, Probie."}`)
	idx(`{"name": "Grant Fuhr",     "jersey": 31, "pos": "G",  "GAA": 3.91,
			"dob": "19620928", "teams": ["EDM", "TOR", "BUF", "STL"],
			"quote": "If I get run into again, I'm taking someone with me."}`)
	idx(`{"name": "Andy Moog",      "jersey": 35, "pos": "G",  "GAA": 3.77,
			"dob": "19600218", "teams": ["EDM", "BOS", "DAL", "MTL"]}`)

	// HACK to let the ES magic happen
	time.Sleep(time.Second)
}

func TearDownTestDB(c *Conn) {
	c.DeleteIndex("oilers")
}
