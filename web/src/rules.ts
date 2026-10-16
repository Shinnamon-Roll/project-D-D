// Display-only derivations (mirror internal/rules/classes.go). The server owns every maximum and every check;
// these only turn scores into the numbers a player adds to a physical die roll.
import type { Abilities, AbilityId, Class } from './protocol'
import { useStore } from './store'

export const ABILITIES: AbilityId[] = ['str', 'dex', 'con', 'int', 'wis', 'cha'] // sheet order

export const abilityMod = (score: number) => Math.floor((score - 10) / 2)

// "+3" / "−1" (true minus sign, so signs line up in tabular figures)
export const signed = (n: number) => (n < 0 ? `−${-n}` : `+${n}`)

// Proficiency bonus by level, from the server's table; level is clamped to 1..20.
export function profBonus(level: number) {
  const t = useStore.getState().classes?.proficiency
  return t?.[Math.min(20, Math.max(1, level)) - 1] ?? 2
}

// saves: the proficient saves (a Character's derived `saves`, or a class's for a roster card).
export const saveBonus = (a: Abilities, ab: AbilityId, pb: number, saves: AbilityId[]) => abilityMod(a[ab]) + (saves.includes(ab) ? pb : 0)

export const skillBonus = (a: Abilities, ab: AbilityId, id: string, pb: number, skills: string[], expertise: string[]) =>
  abilityMod(a[ab]) + (expertise.includes(id) ? 2 * pb : skills.includes(id) ? pb : 0)

// Level-1 HP for a roster card (a picked character carries the server's hpMax).
export const hpLevel1 = (c: Class, con: number) => Math.max(1, c.hitDie + abilityMod(con))

export const featuresAt = (c: Class, level: number) => c.features.filter((f) => f.level <= level)

export function useClass(id: string | undefined) {
  return useStore((s) => s.classes?.classes.find((c) => c.id === id))
}

// The 18 skills in sheet order; empty until /api/classes arrives.
export function useSkills() {
  return useStore((s) => s.classes?.skills) ?? []
}
