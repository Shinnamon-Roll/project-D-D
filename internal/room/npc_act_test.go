package room

import (
	"strings"
	"testing"

	"projectdnd/internal/rules"
)

func TestNpcActRefused(t *testing.T) {
	r, ids := arena(t)
	act := func(kind string) m { return m{"tokenId": ids["Wolf"], "targetId": ids["knight"], "kind": kind} }
	failWith(t, r, "dm", "npc.act", act("talk"), "combat mode is off")
	ok(t, r, "dm", "combat.set", m{"on": true})
	for _, tc := range []struct {
		name, client string
		payload      any
		want         string
	}{
		{"player", "p1", act("talk"), "only the DM can do that"},
		{"unseated", "stranger", act("talk"), "only the DM can do that"},
		{"malformed", "dm", `{"tokenId":5}`, "malformed payload"},
		{"missing token", "dm", m{"tokenId": "t999", "targetId": ids["knight"], "kind": "talk"}, "no such NPC token"},
		{"character as actor", "dm", m{"tokenId": ids["wizard"], "targetId": ids["knight"], "kind": "talk"}, "no such NPC token"},
		{"hidden NPC", "dm", m{"tokenId": ids["Ghost"], "targetId": ids["knight"], "kind": "talk"}, "reveal the NPC first"},
		{"NPC in other scene", "dm", m{"tokenId": ids["Bat"], "targetId": ids["knight"], "kind": "talk"}, "the NPC is not in this scene"},
		{"missing target", "dm", m{"tokenId": ids["Wolf"], "targetId": "t999", "kind": "talk"}, "no such target"},
		{"NPC target", "dm", m{"tokenId": ids["Wolf"], "targetId": ids["Wolf"], "kind": "talk"}, "no such target"},
		{"hidden NPC target", "dm", m{"tokenId": ids["Wolf"], "targetId": ids["Ghost"], "kind": "talk"}, "no such target"},
		{"bad kind", "dm", act("dance"), "kind must be talk, attack, gift or surrender"},
		{"no kind", "dm", m{"tokenId": ids["Wolf"], "targetId": ids["knight"]}, "kind must be talk, attack, gift or surrender"},
		{"talk too long", "dm", m{"tokenId": ids["Wolf"], "targetId": ids["knight"], "kind": "talk", "text": strings.Repeat("ก", 201)}, "talk text must be one line of at most 200 characters"},
		{"talk control char", "dm", m{"tokenId": ids["Wolf"], "targetId": ids["knight"], "kind": "talk", "text": "a\nb"}, "talk text must be one line of at most 200 characters"},
		{"damage -1", "dm", m{"tokenId": ids["Wolf"], "targetId": ids["knight"], "kind": "attack", "hit": true, "damage": -1}, "damage must be 0-999"},
		{"damage 1000", "dm", m{"tokenId": ids["Wolf"], "targetId": ids["knight"], "kind": "attack", "hit": true, "damage": 1000}, "damage must be 0-999"},
		{"damage float", "dm", m{"tokenId": ids["Wolf"], "targetId": ids["knight"], "kind": "attack", "hit": true, "damage": 1.5}, "malformed payload"},
		{"miss with damage", "dm", m{"tokenId": ids["Wolf"], "targetId": ids["knight"], "kind": "attack", "damage": 3}, "a miss deals no damage"},
		{"gift no item", "dm", act("gift"), "item name must be 1-40 characters"},
		{"gift long name", "dm", m{"tokenId": ids["Wolf"], "targetId": ids["knight"], "kind": "gift", "item": m{"name": strings.Repeat("x", 41), "qty": 1}}, "item name must be 1-40 characters"},
		{"gift qty 0", "dm", m{"tokenId": ids["Wolf"], "targetId": ids["knight"], "kind": "gift", "item": m{"name": "Key", "qty": 0}}, "item qty must be 1-99"},
		{"gift qty 100", "dm", m{"tokenId": ids["Wolf"], "targetId": ids["knight"], "kind": "gift", "item": m{"name": "Key", "qty": 100}}, "item qty must be 1-99"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			failWith(t, r, tc.client, "npc.act", tc.payload, tc.want)
			if r.Banner() != nil || len(r.Fresh()) != 0 {
				t.Fatal("banner or log left after a failed event")
			}
		})
	}

	r.token(ids["knight"]).SceneID = r.st.Scenes[1].ID
	failWith(t, r, "dm", "npc.act", act("talk"), "target is not in this scene")

	r, ids = arena(t)
	ok(t, r, "dm", "combat.set", m{"on": true})
	for i := len(r.st.Characters[0].Items); i < maxItemStacks; i++ {
		ok(t, r, "dm", "item.give", m{"characterId": r.st.Characters[0].ID, "name": "thing " + string(rune('A'+i)), "qty": 1})
	}
	failWith(t, r, "dm", "npc.act", m{"tokenId": ids["Wolf"], "targetId": ids["knight"], "kind": "gift", "item": m{"name": "Key", "qty": 1}},
		"at most 50 different items")
}

func TestNpcAct(t *testing.T) {
	const kn = "Aldric Vane"
	for _, tc := range []struct {
		name, client string
		npc          string // arena token name
		extra        m
		feed         rules.L // public line; the DM log line is the same
		kind         string
		bannerEN     string
		hpAfter      int // knight starts at 12
		death        bool
	}{
		{"talk", "dm", "Wolf", m{"kind": "talk", "text": "  Gold, now!  "},
			rules.L{TH: "Wolf คุยกับ อัลดริก เวน: “Gold, now!”", EN: "Wolf talks with Aldric Vane: “Gold, now!”"}, "talk", "", 12, false},
		{"talk empty", "co", "Wolf", m{"kind": "talk"},
			rules.L{TH: "Wolf คุยกับ อัลดริก เวน", EN: "Wolf talks with Aldric Vane"}, "talk", "", 12, false},
		{"talk with format verbs", "dm", "Wolf", m{"kind": "talk", "text": "100% %s %d"},
			rules.L{TH: "Wolf คุยกับ อัลดริก เวน: “100% %s %d”", EN: "Wolf talks with Aldric Vane: “100% %s %d”"}, "talk", "", 12, false},
		{"hit", "dm", "Wolf", m{"kind": "attack", "hit": true, "damage": 5},
			rules.L{TH: "Wolf โจมตี อัลดริก เวน — โดน! −5", EN: "Wolf attacks Aldric Vane — hit! −5"}, "hit", "Wolf → " + kn, 7, false},
		{"hit no damage", "dm", "Wolf", m{"kind": "attack", "hit": true},
			rules.L{TH: "Wolf โจมตี อัลดริก เวน — โดน!", EN: "Wolf attacks Aldric Vane — hit!"}, "hit", "Wolf → " + kn, 12, false},
		{"hit clamps and downs", "dm", "Wolf", m{"kind": "attack", "hit": true, "damage": 999},
			rules.L{TH: "Wolf โจมตี อัลดริก เวน — โดน! −999", EN: "Wolf attacks Aldric Vane — hit! −999"}, "hit", "Wolf → " + kn, 0, true},
		{"miss", "dm", "Wolf", m{"kind": "attack", "hit": false, "damage": 0},
			rules.L{TH: "Wolf โจมตี อัลดริก เวน — ไม่โดน", EN: "Wolf attacks Aldric Vane — miss"}, "miss", "Wolf → " + kn, 12, false},
		{"gift", "dm", "Wolf", m{"kind": "gift", "item": m{"name": " Rusty Key ", "qty": 2}},
			rules.L{TH: "Wolf มอบ Rusty Key ×2 ให้ อัลดริก เวน", EN: "Wolf gives Aldric Vane Rusty Key ×2"}, "gift", "", 12, false},
		{"surrender", "co", "Wolf", m{"kind": "surrender", "text": "ignored"},
			rules.L{TH: "Wolf ยอมแพ้ต่อ อัลดริก เวน", EN: "Wolf surrenders to Aldric Vane"}, "surrender", "", 12, false},
		{"second copy shows #n", "dm", "Wolf #2", m{"kind": "surrender"},
			rules.L{TH: "Wolf #2 ยอมแพ้ต่อ อัลดริก เวน", EN: "Wolf #2 surrenders to Aldric Vane"}, "surrender", "", 12, false},
		{"second copy hit", "dm", "Wolf #2", m{"kind": "attack", "hit": true, "damage": 2},
			rules.L{TH: "Wolf #2 โจมตี อัลดริก เวน — โดน! −2", EN: "Wolf #2 attacks Aldric Vane — hit! −2"}, "hit", "Wolf #2 → " + kn, 10, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, ids := arena(t)
			ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": r.st.CurrentSceneID}) // "Wolf #2"
			ids["Wolf #2"] = r.st.Tokens[len(r.st.Tokens)-1].ID
			ok(t, r, "dm", "combat.set", m{"on": true})
			items := len(r.st.Characters[0].Items)
			p := m{"tokenId": ids[tc.npc], "targetId": ids["knight"]}
			for k, v := range tc.extra {
				p[k] = v
			}
			ok(t, r, tc.client, "npc.act", p)

			feed := r.State().Feed
			line := feed[len(feed)-1].Text
			if tc.death {
				if line.EN != "💀 "+kn+" is knocked out (HP 0)" {
					t.Fatalf("death line %+v", line)
				}
				line = feed[len(feed)-2].Text
			}
			if line != tc.feed {
				t.Fatalf("feed %+v, want %+v", line, tc.feed)
			}
			if fr := r.Fresh(); len(fr) == 0 || fr[0].Text != tc.feed || fr[0].Actor == "" {
				t.Fatalf("dm log %+v, want %+v", fr, tc.feed)
			}
			b := r.Banner()
			wantText := tc.feed.EN
			if tc.bannerEN != "" {
				wantText = tc.bannerEN
			}
			if b == nil || b.Kind != tc.kind || b.Text.EN != wantText || !strings.HasPrefix(b.Text.TH, tc.npc+" ") {
				t.Fatalf("banner %+v, want %s %q", b, tc.kind, wantText)
			}
			if hp := r.st.Characters[0].HP; hp != tc.hpAfter {
				t.Fatalf("knight hp %d, want %d", hp, tc.hpAfter)
			}
			gotItems := r.st.Characters[0].Items
			if tc.kind == "gift" {
				if len(gotItems) != items+1 || gotItems[items] != (rules.RosterItem{Name: rules.L{TH: "Rusty Key", EN: "Rusty Key"}, Qty: 2}) {
					t.Fatalf("items %+v", gotItems)
				}
			} else if len(gotItems) != items {
				t.Fatalf("items changed: %+v", gotItems)
			}
		})
	}
}

// Nothing about hidden tokens reaches players: the player view after an npc.act names only the visible NPC.
func TestNpcActNoHiddenLeak(t *testing.T) {
	r, ids := arena(t)
	ok(t, r, "dm", "combat.set", m{"on": true})
	ok(t, r, "dm", "npc.act", m{"tokenId": ids["Wolf"], "targetId": ids["wizard"], "kind": "attack", "hit": true, "damage": 2})
	if s := string(mustJSON(r.PlayerState())); strings.Contains(s, "Ghost") {
		t.Fatalf("player state names the hidden NPC: %s", s)
	}
}
