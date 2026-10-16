package room

// Token names (docs/features/npc-names.md): unique per scene, case-insensitive. "<base>" first, then "<base> #2",
// "#3"…; the next number is the highest of that base + 1. Visible tokens are numbered from visible tokens only and a
// hidden token that clashes is renamed instead, so players can never infer a hidden token from a gap in the numbers.

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// At most 9 digits so Atoi never overflows; a name past "#999999999" is read as a bare base, so nextName refuses
// to make one.
var numbered = regexp.MustCompile(`^(.+) #([1-9][0-9]{0,8})$`)

const maxNumber = 999999999

// splitName returns name's base and number: "Goblin #5" -> ("Goblin", 5), "Goblin" -> ("Goblin", 1).
func splitName(name string) (string, int) {
	if m := numbered.FindStringSubmatch(name); m != nil {
		n, _ := strconv.Atoi(m[2])
		return m[1], n
	}
	return name, 1
}

// nextName names a token made from src in a scene: src's base with number max(src's number, highest number of that
// base among the counted tokens + 1); number 1 is the bare base. Counted: every token of the scene but skip, hidden
// ones only when withHidden. Character tokens count like any other. Fails past maxNumber.
func (r *Room) nextName(sceneID, src, skip string, withHidden bool) (string, error) {
	base, n := splitName(src)
	for _, t := range r.st.Tokens {
		if t.SceneID != sceneID || t.ID == skip || (t.Hidden && !withHidden) {
			continue
		}
		if b, k := splitName(t.Name); strings.EqualFold(b, base) {
			n = max(n, k+1)
		}
	}
	if n > maxNumber {
		return "", fmt.Errorf("too many copies of %s", base)
	}
	if n == 1 {
		return base, nil
	}
	return base + " #" + strconv.Itoa(n), nil
}

// addNpcToken appends t (a new NPC token, ID and Name unset) named from src and returns its index. A visible token
// gets the next visible name and any hidden token holding that name is renamed; a hidden one gets the scene-wide next.
// On error the caller restores the tokens (see undo).
func (r *Room) addNpcToken(t Token, src string) (int, error) {
	name, err := r.nextName(t.SceneID, src, "", t.Hidden)
	if err != nil {
		return 0, err
	}
	t.ID, t.Name = r.id("t"), name
	r.st.Tokens = append(r.st.Tokens, t)
	if !t.Hidden {
		if err := r.displace(t.SceneID, t.Name, t.ID); err != nil {
			return 0, err
		}
	}
	return len(r.st.Tokens) - 1, nil
}

// undo returns a func that restores the tokens and pending attacks as they are now: a rename that fails midway
// (nextName past maxNumber) must not leave half the renames behind.
func (r *Room) undo() func() {
	toks, atks := slices.Clone(r.st.Tokens), r.attacks
	return func() { r.st.Tokens, r.attacks = toks, atks }
}

// displace renames every hidden token of the scene but keep whose name is name (case-insensitive) to its base's
// scene-wide next number. Old rooms may hold more than one such token, so it renames them all.
func (r *Room) displace(sceneID, name, keep string) error {
	for i := range r.st.Tokens {
		if t := &r.st.Tokens[i]; t.SceneID == sceneID && t.ID != keep && t.Hidden && strings.EqualFold(t.Name, name) {
			next, err := r.nextName(sceneID, t.Name, "", true)
			if err != nil {
				return err
			}
			r.rename(t, next)
		}
	}
	return nil
}

// claim makes way for keep, a character token just added to the scene: every NPC token there named name
// (case-insensitive) is renamed, a visible one to the next visible name (displacing a hidden holder of it), a hidden
// one to the scene-wide next.
func (r *Room) claim(sceneID, name, keep string) error {
	for i := range r.st.Tokens {
		t := &r.st.Tokens[i]
		if t.SceneID != sceneID || t.ID == keep || t.Hidden || t.CharacterID != "" || !strings.EqualFold(t.Name, name) {
			continue
		}
		next, err := r.nextName(sceneID, t.Name, t.ID, false)
		if err != nil {
			return err
		}
		r.rename(t, next)
		r.retarget(t)
		if err := r.displace(sceneID, t.Name, t.ID); err != nil {
			return err
		}
	}
	return r.displace(sceneID, name, keep)
}

// reveal gives a token that has just become visible the next visible name of its base (starting at its pool NPC's
// number, as npc.place would), renames a hidden token holding that name, and points pending attacks at the name.
// On error the caller restores the tokens (see undo).
func (r *Room) reveal(t *Token) error {
	base, _ := splitName(t.Name)
	src := base
	if n, ok := r.npcs.Get(t.NpcID); ok {
		if b, k := splitName(n.Name); strings.EqualFold(b, base) && k > 1 {
			src = base + " #" + strconv.Itoa(k)
		}
	}
	name, err := r.nextName(t.SceneID, src, t.ID, false)
	if err != nil {
		return err
	}
	if name != t.Name {
		r.rename(t, name)
	}
	if err := r.displace(t.SceneID, t.Name, t.ID); err != nil {
		return err
	}
	r.retarget(t)
	return nil
}

// retarget sets the declared name (the attacker's myAttack) of pending attacks on t to t's name. Only for a name
// players see: never from displace, a hidden token's new name must not reach the attacker.
func (r *Room) retarget(t *Token) {
	if slices.ContainsFunc(r.attacks, func(a Attack) bool { return a.TargetTokenID == t.ID }) {
		r.attacks = slices.Clone(r.attacks) // an undo may hold the old slice
		for i := range r.attacks {
			if r.attacks[i].TargetTokenID == t.ID {
				r.attacks[i].TargetName = t.Name
			}
		}
	}
}

// rename is an automatic rename; the DM log says so (DM only, it may name a hidden token).
func (r *Room) rename(t *Token, name string) {
	r.addLog("เปลี่ยนชื่อ "+t.Name+" เป็น "+name, "renamed "+t.Name+" to "+name)
	t.Name = name
}
