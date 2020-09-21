package rest_test

import (
	"net/http"
	"testing"

	"gopkg.in/ns1/ns1-go.v2/mockns1"
	api "gopkg.in/ns1/ns1-go.v2/rest"
	"gopkg.in/ns1/ns1-go.v2/rest/model/ipam"
)

func TestIPAMAddrs(t *testing.T) {
	mock, doer, err := mockns1.New(t)
	if err != nil {
		t.Fatalf("Error creating mock service: %v", err)
	}
	defer mock.Shutdown()

	client := api.NewClient(doer, api.SetEndpoint("https://"+mock.Address+"/v1/"))

	t.Run("List", func(t *testing.T) {
		t.Run("Pagination", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = true
			addrs := []ipam.Address{
				{Name: "a"},
				{Name: "b"},
				{Name: "c"},
				{Name: "d"},
			}
			err := mock.AddTestCase(http.MethodGet, "/ipam/address", http.StatusOK, nil, nil, "", addrs)
			if err != nil {
				t.Fatalf("error adding test case: %v", err)
			}

			respAddrs, _, err := client.IPAM.ListAddrs()
			if err != nil {
				t.Fatalf("error listing IPAM addresses: %v", err)
			}
			if len(respAddrs) != len(addrs) {
				t.Errorf("wrong length: want=%d, got=%d", len(addrs), len(respAddrs))
			}

			for i, addr := range respAddrs {
				if addr.Name != addrs[i].Name {
					t.Errorf("Incorrect name for zone %d: want=%q, got=%q", i, addrs[i].Name, addr.Name)
				}
			}
		})
		t.Run("No Pagination", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = false
			addrs := []ipam.Address{
				{Name: "a"},
				{Name: "b"},
			}

			linkHeader := `</address/0?after=1&limit=2>; rel="next"`
			header := http.Header{}
			header.Set("Link", linkHeader)

			err := mock.AddTestCase(http.MethodGet, "/ipam/address", http.StatusOK, nil, header, "", addrs)
			if err != nil {
				t.Fatalf("error adding test case: %v", err)
			}

			respAddrs, resp, err := client.IPAM.ListAddrs()
			if err != nil {
				t.Fatalf("error listing IPAM addresses: %v", err)
			}
			if len(respAddrs) != len(addrs) {
				t.Errorf("wrong length: want=%d, got=%d", len(addrs), len(respAddrs))
			}
			if l := resp.Header.Get("Link"); l != linkHeader {
				t.Errorf("wrong Link header found: want=%q, got=%q", linkHeader, l)
			}

			for i, addr := range respAddrs {
				if addr.Name != addrs[i].Name {
					t.Errorf("Incorrect name for zone %d: want=%q, got=%q", i, addrs[i].Name, addr.Name)
				}
			}
		})
	})

	t.Run("GetSubnet", func(t *testing.T) {
		defer mock.ClearTestCases()

		client.FollowPagination = false
		addr := ipam.Address{Name: "a"}

		err := mock.AddTestCase(http.MethodGet, "/ipam/address/1", http.StatusOK, nil, nil, "", addr)
		if err != nil {
			t.Fatalf("error adding test case: %v", err)
		}

		respAddr, _, err := client.IPAM.GetSubnet(1)
		if err != nil {
			t.Fatalf("error getting subnet: %v", err)
		}
		if respAddr.Name != addr.Name {
			t.Errorf("wrong address returned, want=%+v, got=%+v", addr, respAddr)
		}
	})

	t.Run("Children", func(t *testing.T) {
		t.Run("Pagination", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = true
			addrs := []ipam.Address{
				{Name: "a"},
				{Name: "b"},
				{Name: "c"},
				{Name: "d"},
			}
			err := mock.AddTestCase(http.MethodGet, "ipam/address/1/children", http.StatusOK, nil, nil, "", addrs)
			if err != nil {
				t.Fatalf("error adding test case: %v", err)
			}

			respAddrs, _, err := client.IPAM.GetChildren(1)
			if err != nil {
				t.Fatalf("error listing child subnets: %v", err)
			}
			if len(respAddrs) != len(addrs) {
				t.Errorf("wrong length: want=%d, got=%d", len(addrs), len(respAddrs))
			}

			for i, addr := range respAddrs {
				if addr.Name != addrs[i].Name {
					t.Errorf("Incorrect name for zone %d: want=%q, got=%q", i, addrs[i].Name, addr.Name)
				}
			}
		})
		t.Run("No Pagination", func(t *testing.T) {
			defer mock.ClearTestCases()

			client.FollowPagination = false
			addrs := []ipam.Address{
				{Name: "a"},
				{Name: "b"},
			}

			linkHeader := `</address/0?after=1&limit=2>; rel="next"`
			header := http.Header{}
			header.Set("Link", linkHeader)

			err := mock.AddTestCase(http.MethodGet, "ipam/address/1/children", http.StatusOK, nil, header, "", addrs)
			if err != nil {
				t.Fatalf("error adding test case: %v", err)
			}

			respAddrs, resp, err := client.IPAM.GetChildren(1)
			if err != nil {
				t.Fatalf("error listing child subnets %v", err)
			}
			if len(respAddrs) != len(addrs) {
				t.Errorf("wrong length: want=%d, got=%d", len(addrs), len(respAddrs))
			}
			if l := resp.Header.Get("Link"); l != linkHeader {
				t.Errorf("wrong Link header found: want=%q, got=%q", linkHeader, l)
			}

			for i, addr := range respAddrs {
				if addr.Name != addrs[i].Name {
					t.Errorf("Incorrect name for zone %d: want=%q, got=%q", i, addrs[i].Name, addr.Name)
				}
			}
		})
	})

	t.Run("GetParent", func(t *testing.T) {
		defer mock.ClearTestCases()

		addr := ipam.Address{Name: "a"}

		err := mock.AddTestCase(http.MethodGet, "/ipam/address/1/parent", http.StatusOK, nil, nil, "", addr)
		if err != nil {
			t.Fatalf("error adding test case: %v", err)
		}

		respAddr, _, err := client.IPAM.GetParent(1)
		if err != nil {
			t.Fatalf("error getting subnet: %v", err)
		}
		if respAddr.Name != addr.Name {
			t.Errorf("wrong address returned, want=%+v, got=%+v", addr, respAddr)
		}
	})

	t.Run("Create", func(t *testing.T) {
		defer mock.ClearTestCases()

		t.Run("RequiredParams", func(t *testing.T) {
			addr := &ipam.Address{Name: "a", Network: 1}
			_, _, err = client.IPAM.CreateSubnet(addr)
			if err == nil {
				t.Errorf("expected a missing prefix to result in an error")
			}
			addr = &ipam.Address{Name: "a", Prefix: "127.0.1.0/24"}
			_, _, err = client.IPAM.CreateSubnet(addr)
			if err == nil {
				t.Errorf("expected a missing network to result in an error")
			}
		})

		addr := &ipam.Address{Name: "a", Prefix: "127.0.1.0/24", Network: 1}
		err := mock.AddTestCase(http.MethodPut, "/ipam/address", http.StatusCreated, nil, nil, addr, addr)
		if err != nil {
			t.Fatalf("error adding test case: %v", err)
		}
		respAddr, _, err := client.IPAM.CreateSubnet(addr)
		if err != nil {
			t.Fatalf("error creating subnet: %v", err)
		}
		if respAddr.Name != addr.Name {
			t.Errorf("wrong address returned: want=%+v, got=%+v", addr, respAddr)
		}
	})

	t.Run("EditSubnet", func(t *testing.T) {
		t.Run("RequiredParams", func(t *testing.T) {
			addr := &ipam.Address{Name: "a"}
			_, _, _, err = client.IPAM.EditSubnet(addr, true)
			if err == nil {
				t.Errorf("expected a missing ID to result in an error")
			}
		})

		t.Run("No Parent", func(t *testing.T) {
			defer mock.ClearTestCases()

			addr := &ipam.Address{ID: 1, Name: "a"}
			err := mock.AddTestCase(http.MethodPost, "/ipam/address/1", http.StatusOK, nil, nil, addr, addr)
			if err != nil {
				t.Fatalf("error adding test case: %v", err)
			}
			respAddr, parent, _, err := client.IPAM.EditSubnet(addr, false)
			if err != nil {
				t.Fatalf("error editing subnet: %v", err)
			}
			if respAddr.Name != addr.Name {
				t.Errorf("wrong address returned: want=%+v, got=%+v", addr, respAddr)
			}
			if parent != nil {
				t.Errorf("non-nil parent unexpectedly returned")
			}
		})
		t.Run("With Parent", func(t *testing.T) {
			defer mock.ClearTestCases()

			addr := &ipam.Address{ID: 1, Name: "a"}
			err := mock.AddTestCase(http.MethodPost, "/ipam/address/1?parent=true", http.StatusOK, nil, nil, addr, addr)
			if err != nil {
				t.Fatalf("error adding test case: %v", err)
			}
			respAddr, parent, _, err := client.IPAM.EditSubnet(addr, true)
			if err != nil {
				t.Fatalf("error editing subnet: %v", err)
			}
			if respAddr.Name != addr.Name {
				t.Errorf("wrong address returned: want=%+v, got=%+v", addr, respAddr)
			}
			if parent == nil {
				t.Errorf("nil parent unexpectedly returned")
			}
		})
	})

	t.Run("SplitSubnet", func(t *testing.T) {
		defer mock.ClearTestCases()

		const prefix = 29
		addr := &ipam.Address{Name: "a", Prefix: "127.0.1.0/24", Network: 1}
		err := mock.AddTestCase(http.MethodPost, "/ipam/address/1/split", http.StatusCreated, nil, nil,
			struct {
				Prefix int `json:"prefix"`
			}{
				Prefix: prefix,
			},
			struct {
				RootAddr  int   `json:"root_address_id"`
				PrefixIDs []int `json:"prefix_ids"`
			}{
				RootAddr:  1,
				PrefixIDs: []int{1},
			},
		)
		if err != nil {
			t.Fatalf("error adding test case: %v", err)
		}
		rootAddr, prefixIDs, _, err := client.IPAM.SplitSubnet(1, prefix)
		if err != nil {
			t.Fatalf("error creating subnet: %v", err)
		}
		if len(prefixIDs) != 1 {
			t.Errorf("wrong number of new subnets returned: want=1, got=%d", len(prefixIDs))
		}
		if rootAddr != 1 {
			t.Errorf("wrong root addr ID returned: want=%d, got=%d", addr.ID, rootAddr)
		}
	})

	t.Run("MergeSubnet", func(t *testing.T) {
		defer mock.ClearTestCases()

		addr := &ipam.Address{Name: "a"}
		err := mock.AddTestCase(http.MethodPost, "/ipam/address/merge", http.StatusCreated, nil, nil,
			struct {
				Root  int `json:"root_address_id"`
				Merge int `json:"merged_address_id"`
			}{
				Root:  1,
				Merge: 2,
			},
			addr,
		)
		if err != nil {
			t.Fatalf("error adding test case: %v", err)
		}
		newAddr, _, err := client.IPAM.MergeSubnet(1, 2)
		if err != nil {
			t.Fatalf("error creating subnet: %v", err)
		}
		if newAddr.Name != addr.Name {
			t.Errorf("wrong addr returned: want=%+v, got=%+v", addr, newAddr)
		}
	})

	t.Run("DeleteSubnet", func(t *testing.T) {
		defer mock.ClearTestCases()

		err := mock.AddTestCase(http.MethodDelete, "/ipam/address/1", http.StatusCreated, nil, nil, "", nil)
		if err != nil {
			t.Fatalf("error adding test case: %v", err)
		}
		_, err = client.IPAM.DeleteSubnet(1)
		if err != nil {
			t.Fatalf("error creating subnet: %v", err)
		}
	})
}
