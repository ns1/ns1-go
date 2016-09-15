package main

import (
	"encoding/json"
	"fmt"

	"github.com/ns1/ns1-go/rest/model/dns"
	"github.com/ns1/ns1-go/rest/model/filter"
)

func prettyPrint(header string, v interface{}) {
	fmt.Println(header)
	fmt.Printf("%#v \n", v)
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func main() {
	domain := "meta.com"

	record := dns.NewRecord(domain, "a", "A") // Construct the A record
	record.Meta.Priority = 100

	// Construct the PRIMARY answer(with BUFFALO feed).
	pAns := dns.NewAv4Answer("1.1.1.1")
	pAns.Meta.Priority = 1 // primary has higher priority
	// Connect the primary answer to the Buffalo feed.
	// pAns.Meta.Up = feeds["buf"].ID

	// Construct the SECONDARY answer(with LONDON feed).
	sAns := dns.NewAv4Answer("2.2.2.2")
	sAns.Meta.Priority = 2 // secondary has lower priority.
	// Connect the secondary answer to the London feed.
	// sAns.Connect(feeds["lon"].ID, "up")

	record.AddAnswer(pAns) // Add primary answer to record
	record.AddAnswer(sAns) // Add secondary answer to record

	// Construct and add both filters to the record(ORDER MATTERS)
	record.AddFilter(filter.NewUp())
	record.AddFilter(filter.NewSelFirstN(1))

	// Helper to show record in json before sending PUT
	prettyPrint("record :", record)

	// b, err := json.MarshalIndent(record.Meta, "", "  ")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(b))

	// // Generic answers to A record
	// recA := dns.NewRecord(domain, "a", "A")
	// printRecord("recA record: ", recA)
	// recA.AddAnswer(&dns.Answer{Rdata: []string{"1.1.1.1"}, Meta: metaBuff})
	// recA.AddAnswer(&dns.Answer{Rdata: []string{"9.9.9.9"}})
	// printRecord("recA WITH ANSWER: ", recA)
	// recA.Regions["newreg"] = dns.NewRegion()
	// recA.Regions["newreg"].Meta.SetUp(true)
	// recA.Regions["newreg"].Meta.SetConnections(4)
	// printRecord("recA WITH Region: ", recA)

	// _, err = client.Records.Create(recA)
	// checkRecordCreate(recA, err)

	// // ALIAS
	// recALIAS := dns.NewRecord(domain, "alias", "ALIAS")
	// recALIAS.AddAnswer(dns.NewALIASAnswer("username.alias"))
	// printRecord("recALIAS WITH ANSWER: ", recA)
	// _, err = client.Records.Create(recALIAS)
	// checkRecordCreate(recALIAS, err)

	// // Av4
	// recAv4 := dns.NewRecord(domain, "Av4", "A")
	// recAv4.AddAnswer(dns.NewAv4Answer("2.2.2.2"))
	// printRecord("recAv4 WITH ANSWER: ", recAv4)
	// _, err = client.Records.Create(recAv4)
	// checkRecordCreate(recAv4, err)

	// // Av6
	// recAv6 := dns.NewRecord(domain, "Av6", "AAAA")
	// recAv6.AddAnswer(dns.NewAv6Answer("dead::beef"))
	// printRecord("recAv6 WITH ANSWER: ", recAv6)
	// _, err = client.Records.Create(recAv6)
	// checkRecordCreate(recAv6, err)

	// // MX
	// recMX := dns.NewRecord(domain, "mx", "MX")
	// recMX.AddAnswer(dns.NewMXAnswer(20, "3.3.3.3"))
	// printRecord("mx WITH ANSWER: ", recMX)
	// _, err = client.Records.Create(recMX)
	// checkRecordCreate(recMX, err)

	// // CNAME
	// recCNAME := dns.NewRecord(domain, "cname", "CNAME")
	// recCNAME.AddAnswer(dns.NewCNAMEAnswer("Av4"))
	// printRecord("cname WITH ANSWER: ", recCNAME)
	// _, err = client.Records.Create(recCNAME)
	// checkRecordCreate(recCNAME, err)

	// // TXT
	// recTXT := dns.NewRecord(domain, "txt", "TXT")
	// recTXT.AddAnswer(dns.NewCNAMEAnswer("some super helpful text"))
	// printRecord("txt WITH ANSWER: ", recTXT)
	// _, err = client.Records.Create(recTXT)
	// checkRecordCreate(recTXT, err)

	// // SRV
	// recSRV := dns.NewRecord(domain, "srv", "SRV")
	// recSRV.AddAnswer(dns.NewSRVAnswer(5, 20, 81, "target.host"))
	// printRecord("srv WITH ANSWER: ", recSRV)
	// _, err = client.Records.Create(recSRV)
	// checkRecordCreate(recSRV, err)

	// _, err = client.Zones.Delete(domain)
	// if err != nil {
	// 	// Ignore if zone doesnt yet exist
	// 	if err != rest.ErrZoneMissing {
	// 		log.Fatal(err)
	// 	} else {
	// 		log.Printf("Delete %s: %s \n", z, err)
	// 	}
	// }
}
