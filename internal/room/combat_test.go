package room

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"
)

// arena: party (p1 knight, p2 wizard, co-DM "co") playing on scene Beach (10x5) with a visible Wolf (hp 5) and a
// hidden Ghost (hp 5); a second scene Cave holds a visible Bat. Returns the room and token ids by name.
func arena(t *testing.T) (*Room, map[string]string) {
	t.Helper()
	r, _, _ := party(t)
	grid := m{"cellSize": 50, "cols": 10, "rows": 5}
	ok(t, r, "dm", "scene.add", m{"name": "Beach", "imageUrl": "/uploads/a.png", "grid": grid})
	ok(t, r, "dm", "scene.add", m{"name": "Cave", "imageUrl": "/uploads/b.png", "grid": grid})
	beach, cave := r.State().Scenes[0].ID, r.State().Scenes[1].ID
	ok(t, r, "dm", "scene.select", m{"sceneId": beach})
	withPool(t, r) // npc1 Wolf (visible), npc2 Ghost (hidden)
	ok(t, r, "dm", "npc.add", m{"name": "Bat", "hpMax": 3})
	ok(t, r, "dm", "npc.place", m{"npcId": "npc1", "sceneId": beach})
	ok(t, r, "dm", "npc.place", m{"npcId": "npc2", "sceneId": beach})
	ok(t, r, "dm", "npc.place", m{"npcId": "npc3", "sceneId": cave})
	ids := map[string]string{}
	for _, tk := range r.State().Tokens {
		ids[tk.Name] = tk.ID
	}
	ids["wizard"] = ids[r.State().Characters[1].Name.EN]
	ids["knight"] = ids[r.State().Characters[0].Name.EN]
	return r, ids
}

func tokenHP(r *Room, id string) int { return *r.token(id).HP }

func lastLog(r *Room) LogEntry { return r.dmLog[len(r.dmLog)-1] }

func TestCombatSet(t *testing.T) {
	r := creating(t)
	failWith(t, r, "dm", "combat.set", m{"on": true}, "not allowed while the room is creating")
	r, _ = arena(t)
	for _, who := range []string{"p1", "stranger"} {
		failWith(t, r, who, "combat.set", m{"on": true}, "only the DM can do that")
	}
	failWith(t, r, "dm", "combat.set", `{"on":"yes"}`, "malformed payload")
	ok(t, r, "co", "combat.set", m{"on": true})
	if !r.State().Combat || lastFeed(r).EN != "Combat started" || lastLog(r).Actor != "Cora" {
		t.Fatalf("combat on: %v %+v %+v", r.State().Combat, lastFeed(r), lastLog(r))
	}
	n := len(r.State().Feed)
	ok(t, r, "dm", "combat.set", m{"on": true}) // repeat: no change, no line
	if len(r.State().Feed) != n {
		t.Fatal("repeat combat.set added a line")
	}
	ok(t, r, "dm", "combat.set", m{"on": false})
	if r.State().Combat || lastFeed(r).EN != "Combat ended" {
		t.Fatalf("combat off: %+v", lastFeed(r))
	}
}

func TestAttackDeclare(t *testing.T) {
	r, ids := arena(t)
	failWith(t, r, "p1", "attack.declare", m{"targetTokenId": ids["Wolf"]}, "combat mode is off")
	ok(t, r, "dm", "combat.set", m{"on": true})
	for _, tc := range []struct {
		name, client string
		payload      any
		want         string
	}{
		{"dm has no character", "dm", m{"targetTokenId": ids["Wolf"]}, "you have no character"},
		{"co-DM has no character", "co", m{"targetTokenId": ids["Wolf"]}, "you have no character"},
		{"unseated", "stranger", m{"targetTokenId": ids["Wolf"]}, "you have no character"},
		{"missing token", "p1", m{"targetTokenId": "t999"}, "no such target"},
		{"hidden token = missing", "p1", m{"targetTokenId": ids["Ghost"]}, "no such target"},
		{"self", "p1", m{"targetTokenId": ids["knight"]}, "you cannot attack yourself"},
		{"other scene", "p1", m{"targetTokenId": ids["Bat"]}, "target is not in this scene"},
		{"malformed", "p1", `{"targetTokenId":5}`, "malformed payload"},
	} {
		t.Run(tc.name, func(t *testing.T) { failWith(t, r, tc.client, "attack.declare", tc.payload, tc.want) })
	}

	ok(t, r, "p1", "attack.declare", m{"targetTokenId": ids["Wolf"]})
	a := r.Attacks()
	if len(a) != 1 || a[0].ByCharacterID != r.State().Characters[0].ID || a[0].ByName.EN != "Aldric Vane" ||
		a[0].TargetTokenID != ids["Wolf"] || a[0].TargetName != "Wolf" || a[0].At == 0 || a[0].ID == "" {
		t.Fatalf("attacks %+v", a)
	}
	if my := r.MyAttack("p1"); my == nil || my.ID != a[0].ID || my.TargetName != "Wolf" {
		t.Fatalf("p1 myAttack %+v", my)
	}
	if r.MyAttack("p2") != nil || r.MyAttack("dm") != nil {
		t.Fatal("myAttack for someone else")
	}
	if l := lastLog(r); l.Actor != "Pat" || l.Text.EN != "Aldric Vane attacks Wolf" {
		t.Fatalf("log %+v", l)
	}
	if b, _ := json.Marshal(r.State()); strings.Contains(string(b), a[0].ID) {
		t.Fatal("attacks leaked into RoomState")
	}
	failWith(t, r, "p1", "attack.declare", m{"targetTokenId": ids["wizard"]}, "you already have a pending attack")
	ok(t, r, "p2", "attack.declare", m{"targetTokenId": ids["knight"]}) // one each: p2 still may

	// No token in the current scene: the DM removed it.
	r, ids = arena(t)
	ok(t, r, "dm", "combat.set", m{"on": true})
	ok(t, r, "dm", "token.remove", m{"tokenId": ids["knight"]})
	failWith(t, r, "p1", "attack.declare", m{"targetTokenId": ids["Wolf"]}, "your token is not in this scene")
}

func TestAttackCancel(t *testing.T) {
	r, ids := arena(t)
	ok(t, r, "dm", "combat.set", m{"on": true})
	failWith(t, r, "p1", "attack.cancel", m{}, "you have no pending attack")
	failWith(t, r, "dm", "attack.cancel", m{}, "you have no character")
	ok(t, r, "p1", "attack.declare", m{"targetTokenId": ids["Wolf"]})
	ok(t, r, "p2", "attack.declare", m{"targetTokenId": ids["Wolf"]})
	ok(t, r, "p1", "attack.cancel", m{})
	if a := r.Attacks(); len(a) != 1 || r.MyAttack("p1") != nil || r.MyAttack("p2") == nil {
		t.Fatalf("after cancel %+v", a)
	}
	if l := lastLog(r); l.Actor != "Pat" || l.Text.EN != "cancelled the attack on Wolf" {
		t.Fatalf("log %+v", l)
	}
}

func TestAttackResolve(t *testing.T) {
	dmg := func(n int) *int { return &n }
	wizName := func(r *Room) string { return r.State().Characters[1].Name.EN }
	for _, tc := range []struct {
		name     string
		target   string // token name in arena ids
		hit      bool
		damage   *int
		wantHP   int    // NPC target hp after (wizard targets: computed from the damage)
		wantFeed string // feed line before any death line; wizard targets: "" = computed
		death    string // "" = no death line
		logWant  string // a DM log line that must exist ("" = skip)
	}{
		{"miss", "Wolf", false, nil, 5, "Aldric Vane → Wolf: miss", "", ""},
		{"miss damage 0", "Wolf", false, dmg(0), 5, "Aldric Vane → Wolf: miss", "", ""},
		{"hit no damage", "Wolf", true, nil, 5, "Aldric Vane → Wolf: hit!", "", ""},
		{"hit npc 3", "Wolf", true, dmg(3), 2, "Aldric Vane → Wolf: hit!", "", "Aldric Vane → Wolf: hit, 3 damage (HP 2/5)"},
		{"hit npc clamps", "Wolf", true, dmg(999), 0, "Aldric Vane → Wolf: hit!", "💀 Wolf falls (HP 0)", "Aldric Vane → Wolf: hit, 999 damage (HP 0/5)"},
		{"hit character 3", "wizard", true, dmg(3), 0, "", "", ""},
		{"hit character clamps", "wizard", true, dmg(999), 0, "", "💀 %s is knocked out (HP 0)", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, ids := arena(t)
			ok(t, r, "dm", "combat.set", m{"on": true})
			ok(t, r, "p1", "attack.declare", m{"targetTokenId": ids[tc.target]})
			id := r.Attacks()[0].ID
			wiz := &r.st.Characters[1]
			hpBefore, hpMax := wiz.HP, wiz.HPMax
			p := m{"attackId": id, "hit": tc.hit}
			if tc.damage != nil {
				p["damage"] = *tc.damage
			}
			ok(t, r, "dm", "attack.resolve", p)
			want := tc.wantFeed
			if tc.target == "wizard" {
				hp := max(hpBefore-*tc.damage, 0)
				want = fmt.Sprintf("Aldric Vane → %s: hit! %s (HP %d/%d)", wizName(r), signed(hp-hpBefore), hp, hpMax)
				if wiz.HP != hp {
					t.Fatalf("wizard hp %d, want %d", wiz.HP, hp)
				}
			} else if got := tokenHP(r, ids[tc.target]); got != tc.wantHP {
				t.Fatalf("npc hp %d, want %d", got, tc.wantHP)
			}
			feed := r.State().Feed
			last := feed[len(feed)-1].Text.EN
			if tc.death != "" {
				d := strings.Replace(tc.death, "%s", wizName(r), 1)
				if last != d {
					t.Fatalf("death line %q, want %q", last, d)
				}
				last = feed[len(feed)-2].Text.EN
			}
			if last != want {
				t.Fatalf("feed %q, want %q", last, want)
			}
			if tc.target != "wizard" && strings.ContainsAny(last, "0123456789") {
				t.Fatalf("npc feed shows numbers: %q", last)
			}
			kind := map[bool]string{true: "hit", false: "miss"}[tc.hit]
			if b := r.Banner(); b == nil || b.Kind != kind || !strings.HasPrefix(b.Text.EN, "Aldric Vane → ") || b.Text.TH == "" {
				t.Fatalf("banner %+v", b)
			}
			if tc.logWant != "" && !hasLog(r, tc.logWant) {
				t.Fatalf("log missing %q: %+v", tc.logWant, r.dmLog)
			}
			if len(r.Attacks()) != 0 || r.MyAttack("p1") != nil {
				t.Fatal("attack still pending")
			}
			failWith(t, r, "dm", "attack.resolve", p, "no such attack") // resolve twice
		})
	}

	r, ids := arena(t)
	ok(t, r, "dm", "combat.set", m{"on": true})
	ok(t, r, "p1", "attack.declare", m{"targetTokenId": ids["Wolf"]})
	id := r.Attacks()[0].ID
	for _, tc := range []struct {
		client string
		p      any
		want   string
	}{
		{"p1", m{"attackId": id, "hit": true}, "only the DM can do that"},
		{"dm", m{"attackId": "a999", "hit": true}, "no such attack"},
		{"dm", m{"attackId": id, "hit": true, "damage": -1}, "damage must be 0-999"},
		{"dm", m{"attackId": id, "hit": true, "damage": 1000}, "damage must be 0-999"},
		{"dm", m{"attackId": id, "hit": false, "damage": 3}, "a miss deals no damage"},
		{"dm", `{"attackId":"x","hit":1}`, "malformed payload"},
	} {
		failWith(t, r, tc.client, "attack.resolve", tc.p, tc.want)
		if r.Banner() != nil {
			t.Fatal("banner left after a failed event")
		}
	}
}

// A target hidden after the declare: the result goes to the DM log only, with no feed line and no banner.
// A hit with damage on a target without HP still logs the number for the DM.
func TestAttackResolveHiddenOrNoHP(t *testing.T) {
	for _, tc := range []struct {
		name       string
		hide, noHP bool
		hit        bool
		damage     int
		wantFeed   string   // new public feed line, "" = none
		wantBanner bool     // banner sent to everyone
		wantLog    []string // DM log lines that must exist
		wantHP     int      // Wolf hp after (when it has HP)
	}{
		{"hidden miss", true, false, false, 0, "", false, []string{"Aldric Vane → Wolf: miss"}, 5},
		{"hidden hit no damage", true, false, true, 0, "", false, []string{"Aldric Vane → Wolf: hit!"}, 5},
		{"hidden hit kills", true, false, true, 5, "", false,
			[]string{"Aldric Vane → Wolf: hit, 5 damage (HP 0/5)", "💀 Wolf falls (HP 0)"}, 0},
		{"no HP hit with damage", false, true, true, 3, "Aldric Vane → Wolf: hit!", true,
			[]string{"Aldric Vane → Wolf: hit, 3 damage (no HP)"}, 0},
		{"hidden no HP hit with damage", true, true, true, 3, "", false,
			[]string{"Aldric Vane → Wolf: hit, 3 damage (no HP)"}, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, ids := arena(t)
			ok(t, r, "dm", "combat.set", m{"on": true})
			ok(t, r, "p1", "attack.declare", m{"targetTokenId": ids["Wolf"]})
			if tc.hide {
				ok(t, r, "dm", "token.setHidden", m{"tokenId": ids["Wolf"], "hidden": true})
			}
			if tc.noHP {
				tk := r.token(ids["Wolf"])
				tk.HP, tk.HPMax = nil, nil
			}
			n := len(r.State().Feed)
			ok(t, r, "dm", "attack.resolve", m{"attackId": r.Attacks()[0].ID, "hit": tc.hit, "damage": tc.damage})
			f := r.State().Feed[n:]
			if tc.wantFeed == "" && len(f) != 0 || tc.wantFeed != "" && (len(f) != 1 || f[0].Text.EN != tc.wantFeed) {
				t.Fatalf("feed %+v, want %q", f, tc.wantFeed)
			}
			if (r.Banner() != nil) != tc.wantBanner {
				t.Fatalf("banner %+v, want %v", r.Banner(), tc.wantBanner)
			}
			for _, l := range tc.wantLog {
				if !hasLog(r, l) {
					t.Fatalf("log missing %q: %+v", l, r.dmLog)
				}
			}
			if !tc.noHP && tokenHP(r, ids["Wolf"]) != tc.wantHP {
				t.Fatalf("wolf hp %d", tokenHP(r, ids["Wolf"]))
			}
		})
	}
}

// Moves of one token by one seat within moveQuiet share one log line (first from, latest to and time);
// another token, another seat or a later move gets its own line.
func TestDMLogMoveCollapse(t *testing.T) {
	r, ids := arena(t)
	clock := time.UnixMilli(1_000_000)
	r.now = func() time.Time { return clock }
	move := func(client, name string, col, row int) {
		t.Helper()
		ok(t, r, client, "token.move", m{"tokenId": ids[name], "col": col, "row": row})
	}
	moves := func() []string {
		var out []string
		for _, e := range r.dmLog {
			if strings.HasPrefix(e.Text.EN, "moved ") {
				out = append(out, e.Actor+": "+e.Text.EN)
			}
		}
		return out
	}
	knight := r.State().Characters[0].Name.EN
	tk := r.token(ids["Wolf"])
	start := cellText(tk.Col, tk.Row)
	kStart := cellText(r.token(ids["knight"]).Col, r.token(ids["knight"]).Row)

	move("dm", "Wolf", 1, 0)
	clock = clock.Add(time.Second)
	move("dm", "Wolf", 2, 0)
	if r.Fresh() != nil {
		t.Fatalf("a collapsed move must add no fresh line: %+v", r.Fresh())
	}
	move("p1", "knight", 3, 3) // another seat and token in between: Wolf's line still collapses
	clock = clock.Add(time.Second)
	move("dm", "Wolf", 4, 0)
	want := []string{"Dee: moved Wolf from " + start + " to (4,0)", "Pat: moved " + knight + " from " + kStart + " to (3,3)"}
	if got := moves(); !slices.Equal(got, want) {
		t.Fatalf("collapsed:\n got %q\nwant %q", got, want)
	}
	if e := r.dmLog[slices.IndexFunc(r.dmLog, func(e LogEntry) bool { return e.Actor == "Dee" && e.move != "" })]; e.At != clock.UnixMilli() {
		t.Fatalf("collapsed line time %d, want %d", e.At, clock.UnixMilli())
	}

	clock = clock.Add(time.Second) // 3 s after Wolf's first move: a new line
	move("dm", "Wolf", 5, 0)
	move("co", "Wolf", 6, 0) // another seat: its own line
	want = append(want, "Dee: moved Wolf from (4,0) to (5,0)", "Cora: moved Wolf from (5,0) to (6,0)")
	if got := moves(); !slices.Equal(got, want) {
		t.Fatalf("after window:\n got %q\nwant %q", got, want)
	}

	// Spamming moves cannot flush older lines out of the log.
	n := len(r.dmLog)
	for i := range maxLog + 10 {
		move("dm", "Wolf", i%10, 4)
	}
	if len(r.dmLog) != n {
		t.Fatalf("spam added %d lines", len(r.dmLog)-n)
	}
}

func hasLog(r *Room, en string) bool {
	for _, e := range r.dmLog {
		if e.Text.EN == en {
			return true
		}
	}
	return false
}

func TestCombatOffCancelsAttacks(t *testing.T) {
	r, ids := arena(t)
	ok(t, r, "dm", "combat.set", m{"on": true})
	ok(t, r, "p1", "attack.declare", m{"targetTokenId": ids["Wolf"]})
	ok(t, r, "p2", "attack.declare", m{"targetTokenId": ids["Wolf"]})
	id := r.Attacks()[0].ID
	ok(t, r, "dm", "combat.set", m{"on": false})
	if len(r.Attacks()) != 0 || r.MyAttack("p1") != nil || !hasLog(r, "cancelled 2 pending attack(s)") {
		t.Fatalf("attacks %+v", r.Attacks())
	}
	failWith(t, r, "dm", "attack.resolve", m{"attackId": id, "hit": true}, "no such attack")
	failWith(t, r, "p1", "attack.cancel", m{}, "you have no pending attack")
}

func TestTokenRemoveDropsAttacksOnIt(t *testing.T) {
	r, ids := arena(t)
	ok(t, r, "dm", "combat.set", m{"on": true})
	ok(t, r, "p1", "attack.declare", m{"targetTokenId": ids["Wolf"]})
	ok(t, r, "p2", "attack.declare", m{"targetTokenId": ids["knight"]})
	ok(t, r, "dm", "token.remove", m{"tokenId": ids["Wolf"]})
	if a := r.Attacks(); len(a) != 1 || a[0].TargetTokenID != ids["knight"] {
		t.Fatalf("attacks %+v", a)
	}
}

func TestDownedAndBackUp(t *testing.T) {
	r, ids := arena(t)
	knight := r.State().Characters[0].ID
	feedLen := func() int { return len(r.State().Feed) }

	n := feedLen()
	ok(t, r, "dm", "token.adjustHp", m{"tokenId": ids["Ghost"], "delta": -99}) // hidden: DM log only
	ok(t, r, "dm", "token.adjustHp", m{"tokenId": ids["Ghost"], "delta": 2})
	if feedLen() != n || !hasLog(r, "💀 Ghost falls (HP 0)") || !hasLog(r, "Ghost is back up") || !hasLog(r, "Ghost HP −5 (0/5)") {
		t.Fatalf("hidden ghost: feed grew %d, log %+v", feedLen()-n, r.dmLog)
	}
	ok(t, r, "dm", "token.adjustHp", m{"tokenId": ids["Wolf"], "delta": -5})
	if lastFeed(r).EN != "💀 Wolf falls (HP 0)" || lastFeed(r).TH != "💀 Wolf ล้มลง (HP 0)" {
		t.Fatalf("wolf %+v", lastFeed(r))
	}
	n = feedLen()
	ok(t, r, "dm", "token.adjustHp", m{"tokenId": ids["Wolf"], "delta": -1}) // already 0: no change, no line
	if feedLen() != n {
		t.Fatal("line for a no-change hp edit")
	}
	ok(t, r, "dm", "character.adjust", m{"characterId": knight, "hp": -999})
	if lastFeed(r).EN != "💀 Aldric Vane is knocked out (HP 0)" || lastFeed(r).TH != "💀 อัลดริก เวน หมดสติ (HP 0)" {
		t.Fatalf("knight %+v", lastFeed(r))
	}
	ok(t, r, "dm", "character.adjust", m{"characterId": knight, "hp": 1})
	if lastFeed(r).EN != "Aldric Vane is back up" || lastFeed(r).TH != "อัลดริก เวน ฟื้นขึ้น" {
		t.Fatalf("knight up %+v", lastFeed(r))
	}
}

func TestDMLog(t *testing.T) {
	r, ids := arena(t)
	for _, e := range []struct{ actor, en string }{
		{"Dee", "opened the room: Starter Set: Dragons of Stormwreck Isle"},
		{"Pat", "took a player seat"},
		{"Cora", "took the co-DM seat"},
		{"Pat", "is ready"},
		{"Dee", "started the game — character creation"},
		{"Pat", "Pat picked Aldric Vane (Knight)"},
		{"Dee", "added scene Beach (grid 10×5)"},
		{"Dee", "switched to scene Beach"},
		{"Dee", "added NPC Wolf to the pool"},
		{"Dee", "placed Ghost in Beach (hidden)"},
	} {
		if !hasLogBy(r, e.actor, e.en) {
			t.Fatalf("missing %s: %q in %+v", e.actor, e.en, r.dmLog)
		}
	}
	tk := r.token(ids["Wolf"])
	from := cellText(tk.Col, tk.Row)
	ok(t, r, "dm", "token.move", m{"tokenId": ids["Wolf"], "col": 4, "row": 3})
	if l := lastLog(r); l.Text.EN != "moved Wolf from "+from+" to (4,3)" || l.Text.TH != "ย้าย Wolf จาก "+from+" ไป (4,3)" {
		t.Fatalf("move log %+v", l)
	}
	for _, e := range []struct {
		client, typ string
		p           m
		en          string
	}{
		{"p1", "skill.addCustom", m{"name": "Kick", "desc": "", "cost": 1}, "added skill Kick to Aldric Vane"},
		{"dm", "token.setHidden", m{"tokenId": ids["Wolf"], "hidden": true}, "hid Wolf"},
		{"dm", "token.setHidden", m{"tokenId": ids["Wolf"], "hidden": false}, "revealed Wolf"},
		{"dm", "fog.add", m{"sceneId": r.State().Scenes[0].ID, "col": 1, "row": 1, "w": 2, "h": 1}, "added fog 2×1 at (1,1) in Beach"},
		{"dm", "roster.reset", m{"classId": "knight"}, "reset roster class knight"},
		{"dm", "party.rest", m{}, "DM: rest — mana refilled"},
	} {
		ok(t, r, e.client, e.typ, e.p)
		if l := lastLog(r); l.Text.EN != e.en {
			t.Fatalf("%s: log %+v, want %q", e.typ, l, e.en)
		}
	}

	// Failed events log nothing.
	n := len(r.dmLog)
	fail(t, r, "p1", "token.move", m{"tokenId": ids["Wolf"], "col": 0, "row": 0})
	if len(r.dmLog) != n || len(r.Fresh()) != 0 {
		t.Fatal("failed event logged")
	}
	if b, _ := json.Marshal(r.DMLog()); strings.Contains(string(b), `"p1"`) || strings.Contains(string(b), `"dm"`) {
		t.Fatalf("clientId in the log: %s", b)
	}
}

func hasLogBy(r *Room, actor, en string) bool {
	for _, e := range r.dmLog {
		if e.Actor == actor && e.Text.EN == en {
			return true
		}
	}
	return false
}

func TestDMLogActorBeforeEvent(t *testing.T) {
	r, _ := lobby(t)
	ok(t, r, "p1", "seat.leave", m{}) // the seat name is gone after the event
	if l := lastLog(r); l.Actor != "Pat" || l.Text.EN != "left their seat" {
		t.Fatalf("leave log %+v", l)
	}
}

func TestDMLogNoteThrottle(t *testing.T) {
	r, _ := lobby(t)
	clock := time.UnixMilli(1_000_000)
	r.now = func() time.Time { return clock }
	count := func() int {
		n := 0
		for _, e := range r.dmLog {
			if strings.HasPrefix(e.Text.EN, "edited note ") {
				n++
			}
		}
		return n
	}
	ok(t, r, "dm", "note.add", m{"title": "Twist"})
	ok(t, r, "dm", "note.update", m{"id": "n1", "title": "Twist", "body": "a"})
	ok(t, r, "dm", "note.add", m{"title": "Map"}) // another note: its own window
	if count() != 2 || lastLog(r).Text.EN != "edited note Map" || lastLog(r).Text.TH != "แก้โน้ต Map" {
		t.Fatalf("within 60s: %d %+v", count(), r.dmLog)
	}
	clock = clock.Add(61 * time.Second)
	ok(t, r, "dm", "note.update", m{"id": "n1", "title": "Plot", "body": "b"})
	if count() != 3 || lastLog(r).Text.EN != "edited note Plot" {
		t.Fatalf("after 60s: %d", count())
	}
	if b, _ := json.Marshal(r.DMLog()); strings.Contains(string(b), `"b"`) && strings.Contains(string(b), "body") {
		t.Fatal("note body in the log")
	}
}

func TestDMLogCap(t *testing.T) {
	r, _ := lobby(t)
	for i := range maxLog + 50 {
		ok(t, r, "p1", "seat.ready", m{"ready": i%2 == 0})
	}
	l := r.DMLog()
	if len(l) != maxLog || l[len(l)-1].Text.EN != "is not ready" || l[len(l)-1].Actor != "Pat" {
		t.Fatalf("len %d last %+v", len(l), l[len(l)-1])
	}
}
