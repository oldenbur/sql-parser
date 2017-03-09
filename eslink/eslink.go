package eslink

import (
//	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	log "github.com/cihub/seelog"
	es "github.com/mattbaird/elastigo/lib"
	util "github.com/oldenbur/gostree"
)

func conn() error {

	c := es.NewConn()
	resp, err := c.DoCommand("GET", "/_cat/indices", nil, nil)
	if err != nil {
		return log.Error("conn error in _cat/indices: ", err)
	}

	log.Debugf("/_cat/indices resp:\n%s", string(resp))

	catIndices := c.GetCatIndexInfo("")
	log.Debug("CatIndexInfo:")
	for _, i := range catIndices {
		log.Debugf("%s", i.Name)
	}

	return nil
}

type BulkRecord struct {
	Index string `json:"_index"`
	Type string `json:"_type"`
	Id string `json:"_id"`
}

type BulkIndex struct {
	Index BulkRecord `json:"index"`
}

func bulkJson() error {

	jsonFile := "/Users/paul.oldenburg/go/src/github.com/oldenbur/sql-parser/logs_search_sm.json"
	log.Debugf("jsonFile: %s", jsonFile)

	f, err := os.Open(jsonFile)
	if err != nil {
		return log.Error("bulkJson error in os.Open: ", err)
	}
	stree, err := util.NewSTreeJson(f)
	if err != nil {
		return log.Error("bulkJson error in NewSTreeJson: ", err)
	}

	log.Debugf("hits /total: %d", stree.IntVal("hits/total"))

	hits := stree.SliceVal("hits/hits")
	log.Debugf("len(hits): %d", len(hits))

	var b bytes.Buffer
	for _, hitVal := range hits {

		hit, err := util.ValueOf(hitVal)
		if err != nil {
			return log.Errorf("bulkJson error taking ValueOf STree: %v", err)
		}
//		hit := stree.STreeVal(fmt.Sprintf("hits/hits[%d]", i))
//		log.Debugf("hit[%d] - index: %s  type: %s  id: %s", i,
//			hit.StrVal("_index"), hit.StrVal("_type"), hit.StrVal("_id"))

		bi := &BulkIndex{
			Index: BulkRecord{
				Index: hit.StrVal("_index"),
				Type: hit.StrVal("_type"),
				Id: hit.StrVal("_id"),
			},
		}
		bij, err := json.Marshal(bi)
		if err != nil {
			return log.Errorf("bulkJson Marshal error: %v", err)
		}
		//		log.Debugf("bij: %s", string(bij))

//		lj, err := hit.MarshalJSON()
		lj, err := json.Marshal(hit)
		if err != nil {
			return log.Errorf("bulkJson STree.ToJson error: %v", err)
		}
//		log.Debugf("lj: %s", string(lj))

		fmt.Fprintf(&b, "%s\n%s\n", bij, lj)
	}

//	scanner := bufio.NewScanner(&b)
//	for scanner.Scan() {
//		fmt.Println(scanner.Text()) // Println will add back the final '\n'
//	}
//	if err := scanner.Err(); err != nil {
//		return log.Errorf("bulkJson error scanning bulk request: %v", err)
//	}

	c := es.NewConn()
	resp, err := c.DoCommand("POST", "/_bulk", nil, &b)
	if err != nil {
		return log.Error("conn error in _bulk: ", err)
	}
	log.Debugf("/_bulk resp:\n%s", string(resp))

	return nil
}