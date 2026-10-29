package parser

import (
	"testing"

	"gbootcms/apps/admin/model/content"
)

func TestResolvePageTitleSinglePage(t *testing.T) {
	ctx := &Context{
		Sort: &content.ContentSort{
			Mcode: "1",
			Type:  1,
			Name:  "關於我們",
		},
		Content: &content.Content{
			Title: "關於我們",
		},
		Site: &content.Site{
			Title:    "GbootCMS",
			Subtitle: "高效能系統",
		},
	}

	title := resolvePageTitle(ctx)
	t.Logf("Generated Title: %s", title)
	expected := "關於我們-GbootCMS-高效能系統"
	if title != expected {
		t.Errorf("Expected %s, got %s", expected, title)
	}
}
