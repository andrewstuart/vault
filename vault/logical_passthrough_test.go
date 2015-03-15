package vault

import (
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
)

func TestPassthroughBackend_impl(t *testing.T) {
	var _ logical.Backend = new(PassthroughBackend)
}

func TestPassthroughBackend_RootPaths(t *testing.T) {
	var b PassthroughBackend
	root := b.RootPaths()
	if len(root) != 0 {
		t.Fatalf("unexpected: %v", root)
	}
}

func TestPassthroughBackend_Write(t *testing.T) {
	var b PassthroughBackend
	req := logical.TestRequest(t, logical.WriteOperation, "foo")
	req.Data["raw"] = "test"

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	out, err := req.Storage.Get("foo")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if out == nil {
		t.Fatalf("failed to write to view")
	}
}

func TestPassthroughBackend_Read(t *testing.T) {
	var b PassthroughBackend
	req := logical.TestRequest(t, logical.WriteOperation, "foo")
	req.Data["raw"] = "test"
	req.Data["lease"] = "1h"
	storage := req.Storage

	if _, err := b.HandleRequest(req); err != nil {
		t.Fatalf("err: %v", err)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "foo")
	req.Storage = storage

	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	expected := &logical.Response{
		IsSecret: true,
		Lease: &logical.Lease{
			Renewable:    false,
			Revokable:    false,
			Duration:     time.Hour,
			MaxDuration:  time.Hour,
			MaxIncrement: 0,
		},
		Data: map[string]interface{}{
			"raw":   "test",
			"lease": "1h",
		},
	}

	if !reflect.DeepEqual(resp, expected) {
		t.Fatalf("bad response.\n\nexpected: %#v\n\nGot: %#v", expected, resp)
	}
}

func TestPassthroughBackend_Delete(t *testing.T) {
	var b PassthroughBackend
	req := logical.TestRequest(t, logical.WriteOperation, "foo")
	req.Data["raw"] = "test"
	storage := req.Storage

	if _, err := b.HandleRequest(req); err != nil {
		t.Fatalf("err: %v", err)
	}

	req = logical.TestRequest(t, logical.DeleteOperation, "foo")
	req.Storage = storage
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}

	req = logical.TestRequest(t, logical.ReadOperation, "foo")
	req.Storage = storage
	resp, err = b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if resp != nil {
		t.Fatalf("bad: %v", resp)
	}
}

func TestPassthroughBackend_List(t *testing.T) {
	var b PassthroughBackend
	req := logical.TestRequest(t, logical.WriteOperation, "foo")
	req.Data["raw"] = "test"
	storage := req.Storage

	if _, err := b.HandleRequest(req); err != nil {
		t.Fatalf("err: %v", err)
	}

	req = logical.TestRequest(t, logical.ListOperation, "")
	req.Storage = storage
	resp, err := b.HandleRequest(req)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	expected := &logical.Response{
		IsSecret: false,
		Lease:    nil,
		Data: map[string]interface{}{
			"keys": []string{"foo"},
		},
	}

	if !reflect.DeepEqual(resp, expected) {
		t.Fatalf("bad response.\n\nexpected: %#v\n\nGot: %#v", expected, resp)
	}
}