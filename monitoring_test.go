package nsone

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"reflect"
	"testing"
)

func TestUnmarshalMonitoringJobTypes(t *testing.T) {
	data, err := ioutil.ReadFile("test/data/monitoring_jobtypes.json")
	if err != nil {
		t.Error(err)
	}
	var m MonitoringJobTypes
	if err = json.Unmarshal(data, &m); err != nil {
		t.Error(err)
	}
	job, exists := m["ping"]
	if !exists {
		t.Error("'ping' job does not exist")
	}
	exp_shortdesc := "Ping (ICMP)"
	if job.ShortDesc != exp_shortdesc {
		t.Error(fmt.Sprintf("Expected job.ShortDesc to be '%s', but was actually '%s'", exp_shortdesc, job.ShortDesc))
	}
	exp_desc := "Ping a host using ICMP packets."
	if job.Desc != exp_desc {
		t.Error(fmt.Sprintf("Expected job.Desc to be '%s', but was actually '%s'", exp_desc, job.Desc))
	}
	results := job.Results
	if results == nil {
		t.Error("Job has no results")
	}
	r, exists := results["rtt"]
	if !exists {
		t.Error("Field 'ping' does not exist")
	}
	/*if !reflect.DeepEquals(r.Comparators, []string{"<", ">", "<=", ">=", "==", "!="}) {
		t.Error("comparators not as expected")
	}*/
	if !r.Metric {
		t.Error("Expected a metric")
	}
	if r.Validator != "number" {
		t.Error("validator not as expected")
	}
	r_exp_shortdesc := "Round trip time"
	if r.ShortDesc != r_exp_shortdesc {
		t.Error(fmt.Sprintf("Result ShortDesc wrong, expected '%s', got '%s'", r_exp_shortdesc, r.ShortDesc))
	}
	if r.Type != "number" {
		t.Error("Type is bad")
	}
	if r.Desc != "Average round trip time (in ms) of returned pings." {
		t.Error("Desc is bad")
	}
}

func TestUnmarshalMonitoringJobs(t *testing.T) {
	data, err := ioutil.ReadFile("test/data/monitoring_jobs.json")
	if err != nil {
		t.Error(err)
	}
	var m MonitoringJobs
	if err = json.Unmarshal(data, &m); err != nil {
		t.Error(err)
	}
	if len(m) != 1 {
		t.Error("Do not have any jobs")
	}
	j := m[0]
	if j.Id != "52a27d4397d5f07003fdbe7b" {
		t.Error("Wrong ID")
	}
	conf := j.Config
	if conf["host"] != "1.2.3.4" {
		t.Error("Wrong host")
	}
	status := j.Status["global"]
	if status.Since != 1389407609 {
		t.Error("since has unexpected value")
	}
	if status.Status != "up" {
		t.Error("Status is not up")
	}
	r := j.Rules[0]
	if r.Key != "rtt" {
		t.Error("RTT rule key is wrong")
	}
	if r.Value != 100 {
		t.Error("RTT rule value is wrong")
	}
	if r.Comparison != "<" {
		t.Error("RTT rule comparison is wrong")
	}
	if j.JobType != "ping" {
		t.Error("Jobtype is wrong")
	}
	if j.Regions[0] != "lga" {
		t.Error("First region is not lga")
	}
	if !j.Active {
		t.Error("Job is not active")
	}
	if j.Frequency != 60 {
		t.Error("Job frequency != 60")
	}
	if j.Policy != "quorum" {
		t.Error("Job policy is not quorum")
	}
	if j.RegionScope != "fixed" {
		t.Error("Job region scope is not fixed")
	}
}
