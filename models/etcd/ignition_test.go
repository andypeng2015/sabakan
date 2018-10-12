package etcd

import (
	"context"
	"fmt"
	"testing"

	"github.com/cybozu-go/sabakan"
)

func testTemplate(t *testing.T) {
	t.Parallel()

	d, _ := testNewDriver(t)

	_, err := d.GetTemplate(context.Background(), "cs", "1111")
	if err != sabakan.ErrNotFound {
		t.Fatal("unexpected error: ", err)
	}

	id, err := d.PutTemplate(context.Background(), "cs", "data", nil)
	if err != nil {
		t.Fatal(err)
	}

	ign, err := d.GetTemplate(context.Background(), "cs", id)
	if err != nil {
		t.Fatal(err)
	}
	if ign != "data" {
		t.Error("wrong template stored", ign)
	}

	err = d.DeleteTemplate(context.Background(), "cs", id)
	if err != nil {
		t.Fatal(err)
	}

	err = d.DeleteTemplate(context.Background(), "cs", id)
	if err != sabakan.ErrNotFound {
		t.Fatal("unexpected error: ", err)
	}
}

func testTemplateIDs(t *testing.T) {
	t.Parallel()

	d, _ := testNewDriver(t)

	_, err := d.GetTemplateMetadataList(context.Background(), "cs")
	if err != sabakan.ErrNotFound {
		t.Fatal("unexpected error: ", err)
	}

	_, err = d.PutTemplate(context.Background(), "cs", "data", map[string]string{"version": "20181010"})
	if err != nil {
		t.Fatal(err)
	}

	metadata, err := d.GetTemplateMetadataList(context.Background(), "cs")
	if err != nil {
		t.Fatal(err)
	}
	if len(metadata) != 1 {
		t.Error("wrong number of templates", len(metadata))
		fmt.Println(metadata)
	}
	if len(metadata[0]["id"]) == 0 {
		t.Error("id is empty")
	}
	if metadata[0]["version"] != "20181010" {
		t.Error("wrong version", metadata[0])
	}

	for i := 0; i < sabakan.MaxIgnitions+10; i++ {
		_, err = d.PutTemplate(context.Background(), "cs", "data", nil)
		if err != nil {
			t.Fatal(err)
		}
	}
	metadata, err = d.GetTemplateMetadataList(context.Background(), "cs")
	if err != nil {
		t.Fatal(err)
	}
	if len(metadata) != sabakan.MaxIgnitions {
		t.Error("wrong number of templates", len(metadata))
	}

}

func TestIgnitionTemplate(t *testing.T) {
	t.Run("Template", testTemplate)
	t.Run("TemplateIDs", testTemplateIDs)
}
