// Default roster: the 10 classes from the JSON block at the end of docs/design/roster.md (original text,
// not from any book or game). Keep both in sync; TestDefaultRosterMatchesDoc compares them.

package rules

// DefaultRoster returns a fresh copy of the 10 built-in characters.
func DefaultRoster() []RosterEntry {
	return []RosterEntry{
		{
			ID: "knight", Role: "tank", Art: "",
			HeroName:  L{"อัลดริก เวน", "Aldric Vane"},
			ClassName: L{"อัศวิน", "Knight"},
			Blurb:     L{"อัศวินที่เคยพลาดจนเสียเพื่อนไป จึงสาบานว่าจะไม่ยอมให้ใครล้มก่อนตัวเองอีก", "A knight who once lost a comrade and swore no friend would ever fall before him again."},
			Stats:     Stats{HPMax: 36, Str: 5, End: 8, Agi: 2, MPMax: 4},
			Skills: []RosterSkill{
				{"stand-guard", L{"ยืนบังให้", "Stand Guard"}, L{"เลือกเพื่อน 1 คนที่อยู่ติดกัน จนถึงตาถัดไปของเรา การโจมตีที่เล็งเพื่อนคนนั้นจะมาโดนเราแทน", "Pick an ally next to you. Until your next turn, any attack aimed at them hits you instead."}, 1},
				{"shield-bash", L{"กระแทกโล่", "Shield Bash"}, L{"เอาโล่กระแทกศัตรูที่อยู่ติดกัน ถ้าโดน ศัตรูกระเด็นถอยไป 2 ช่อง และเดินไม่ได้ในตาถัดไปของมัน", "Slam an adjacent enemy with your shield. On a hit it is knocked back 2 squares and can't move on its next turn."}, 2},
				{"last-stand", L{"ยืนหยัดครั้งสุดท้าย", "Last Stand"}, L{"จนจบรอบถัดไป เลือดของเราจะไม่ลดต่ำกว่า 1 ไม่ว่าจะโดนหนักแค่ไหน", "Until the end of the next round, your HP can't drop below 1, no matter how hard you are hit."}, 4},
			},
			Items: []RosterItem{
				{L{"โล่ตราสิงห์ตระกูลเวน", "Vane Family Lion Shield"}, 1},
				{L{"ดาบยาวอัศวิน", "Knight's Longsword"}, 1},
				{L{"หินลับดาบกับน้ำมันขัด", "Whetstone and Blade Oil"}, 1},
				{L{"ม้วนผ้าพันแผลสนาม", "Field Bandage Roll"}, 3},
			},
		},
		{
			ID: "wizard", Role: "control", Art: "",
			HeroName:  L{"โซรายา ควิลล์", "Soraya Quill"},
			ClassName: L{"จอมเวท", "Wizard"},
			Blurb:     L{"นักปราชญ์จากหอดูดาวกลางทะเลทราย ผู้อ่านดวงดาวได้คล่องเหมือนอ่านแผนที่สนามรบ", "A scholar from the desert observatories who reads the stars like a battle map."},
			Stats:     Stats{HPMax: 16, Str: 2, End: 3, Agi: 7, MPMax: 20},
			Skills: []RosterSkill{
				{"arcane-bolt", L{"ลูกศรเวท", "Arcane Bolt"}, L{"ยิงลำแสงเวทใส่ศัตรู 1 ตัวในระยะ 6 ช่อง ถ้าโดน เสียหาย d6", "Fire a bolt of light at one enemy up to 6 squares away. On a hit it takes d6 damage."}, 1},
				{"frozen-moment", L{"หยุดกาล", "Frozen Moment"}, L{"ตรึงเวลาศัตรู 1 ตัวในระยะ 6 ช่อง มันต้องข้ามตาถัดไปทั้งตา (บอสอาจต้านได้ DM ตัดสิน)", "Freeze one enemy within 6 squares in time: it skips its next turn entirely (a boss may resist, DM decides)."}, 3},
				{"starfall", L{"ดาวตก", "Starfall"}, L{"เรียกดาวตกลงพื้นที่ 3×3 ช่องในระยะ 8 ช่อง ทุกตัวในพื้นที่ ทั้งศัตรูและเพื่อน เสียหาย 2d6 โดยไม่ต้องทอยโดน", "Call falling stars onto a 3×3 area within 8 squares. Everyone inside, friend or foe, takes 2d6 damage, no hit roll needed."}, 5},
			},
			Items: []RosterItem{
				{L{"ไม้เท้าแก้วดารา", "Starglass Staff"}, 1},
				{L{"ตำราเวทปกหนังสีคราม", "Indigo-Bound Spellbook"}, 1},
				{L{"ขวดหมึกกับปากกาขนนกฮูก", "Ink Pot and Owl-Feather Quill"}, 1},
				{L{"ขวดน้ำยาจิตใส", "Clear-Mind Tonic"}, 2},
			},
		},
		{
			ID: "rogue", Role: "damage", Art: "",
			HeroName:  L{"ซารีฟ เคช", "Zarif Kesh"},
			ClassName: L{"โจร", "Rogue"},
			Blurb:     L{"หัวขโมยเจ้าเสน่ห์จากเมืองคลอง ที่ติดหนี้คนไว้ครึ่งเมือง", "A charming thief from the canal city who owes money to half the town."},
			Stats:     Stats{HPMax: 18, Str: 6, End: 3, Agi: 9, MPMax: 7},
			Skills: []RosterSkill{
				{"quick-fingers", L{"มือไว", "Quick Fingers"}, L{"สะเดาะกุญแจ หรือขโมยของชิ้นเล็ก 1 ชิ้นจากคนที่อยู่ติดกันได้เลยโดยไม่ต้องทอย ยกเว้นเขากำลังจ้องเราอยู่", "Pick a lock, or steal one small item from someone next to you, with no roll, unless they are watching you closely."}, 1},
				{"shadow-step", L{"หลบเข้าเงา", "Shadow Step"}, L{"กลืนหายไปในเงา ศัตรูเล็งเราไม่ได้จนกว่าเราจะโจมตีหรือจบตาถัดไป การโจมตีแรกจากที่ซ่อนแรงขึ้น +d6", "Melt into the shadows: enemies can't target you until you attack or your next turn ends. Your first attack from hiding deals +d6 damage."}, 2},
				{"venom-edge", L{"คมพิษ", "Venom Edge"}, L{"ทาพิษลงใบมีด ศัตรูตัวถัดไปที่เราฟันโดนจะเสียเลือด 2 แต้มทุกต้นตาของมัน นาน 3 ตา", "Coat your blade in poison. The next enemy you hit loses 2 HP at the start of each of its turns for 3 turns."}, 4},
			},
			Items: []RosterItem{
				{L{"มีดสั้นคู่เงินจันทรา", "Twin Moonsilver Daggers"}, 1},
				{L{"ชุดสะเดาะกุญแจ", "Lockpick Roll"}, 1},
				{L{"ลูกระเบิดควัน", "Smoke Pellet"}, 3},
				{L{"ตะขอเกี่ยวกับเชือกไหม", "Grappling Hook and Silk Rope"}, 1},
				{L{"ลูกเต๋าถ่วงน้ำหนัก", "Loaded Dice"}, 2},
			},
		},
		{
			ID: "ranger", Role: "damage", Art: "",
			HeroName:  L{"ซีกรุน แอชเวล", "Sigrun Ashvale"},
			ClassName: L{"พรานป่า", "Ranger"},
			Blurb:     L{"นักล่าเงียบขรึมจากป่าสนทางเหนือ ที่ไว้ใจเหยี่ยวคู่ใจมากกว่าคน", "A quiet huntress from the northern pines who trusts her falcon more than people."},
			Stats:     Stats{HPMax: 22, Str: 4, End: 6, Agi: 8, MPMax: 5},
			Skills: []RosterSkill{
				{"mark-quarry", L{"หมายตาเหยื่อ", "Mark the Quarry"}, L{"หมายตาศัตรู 1 ตัวที่มองเห็น จนกว่ามันจะล้มหรือเราพักผ่อน การโจมตีของเราใส่มันแรงขึ้น +2", "Mark one enemy you can see. Until it falls or you rest, your attacks against it deal +2 damage."}, 1},
				{"falcon-scout", L{"เหยี่ยวสอดแนม", "Falcon Scout"}, L{"ปล่อยเหยี่ยวบินล่วงหน้าไกลสูงสุด 20 ช่อง DM จะบอกว่ามันเห็นอะไร เช่น ศัตรู กับดัก หรือทางออก", "Send your falcon up to 20 squares ahead. The DM tells you what it sees: enemies, traps, exits."}, 2},
				{"deadeye-shot", L{"เล็งสังหาร", "Deadeye Shot"}, L{"ใช้ทั้งตายืนเล็งโดยไม่เดิน ตาถัดไปลูกธนูจะโดนแน่นอน เสียหาย 2d8", "Spend your whole turn aiming without moving. On your next turn the arrow hits automatically for 2d8 damage."}, 4},
			},
			Items: []RosterItem{
				{L{"ธนูยาวไม้แอช", "Ashwood Longbow"}, 1},
				{L{"ลูกธนูขนนกสีเทา", "Grey-Fletched Arrow"}, 20},
				{L{"นกหวีดกระดูกเรียกเหยี่ยว", "Bone Falcon Whistle"}, 1},
				{L{"เนื้อกวางรมควัน", "Smoked Venison Strip"}, 5},
			},
		},
		{
			ID: "cleric", Role: "healer", Art: "",
			HeroName:  L{"อิมารา โซเลนน์", "Imara Solenne"},
			ClassName: L{"นักบวช", "Cleric"},
			Blurb:     L{"หมอแห่งวิหารสุริยันที่หัวเราะเสียงดัง สวดเสียงเบา และไม่เคยทิ้งใครไว้ข้างหลัง", "A sun-temple healer who laughs loudly, prays quietly, and never leaves anyone behind."},
			Stats:     Stats{HPMax: 28, Str: 4, End: 5, Agi: 2, MPMax: 16},
			Skills: []RosterSkill{
				{"purify", L{"ชำระล้าง", "Purify"}, L{"แตะตัวเพื่อน 1 คน เพื่อล้างพิษ คำสาป หรืออาการหวาดกลัวออก 1 อย่าง", "Touch an ally to remove one poison, curse, or fear effect."}, 1},
				{"healing-touch", L{"สัมผัสเยียวยา", "Healing Touch"}, L{"วางมือบนเพื่อนที่อยู่ติดกัน เพื่อนฟื้นเลือด d6 + 2", "Lay a hand on an ally next to you: they regain d6 + 2 HP."}, 2},
				{"dawn-return", L{"คืนชีพยามรุ่ง", "Return at Dawn"}, L{"ปลุกเพื่อนที่อยู่ติดกันและเลือดเหลือ 0 ให้ลุกขึ้นมาพร้อมเลือดครึ่งหนึ่งของเลือดสูงสุด", "Bring an adjacent ally at 0 HP back to their feet with half their max HP."}, 5},
			},
			Items: []RosterItem{
				{L{"จี้สุริยันศักดิ์สิทธิ์", "Sunburst Holy Pendant"}, 1},
				{L{"กระบองเหล็ก", "Iron Mace"}, 1},
				{L{"ถุงสมุนไพรห้ามเลือด", "Styptic Herb Pouch"}, 4},
				{L{"ขวดน้ำมนต์", "Flask of Holy Water"}, 2},
				{L{"สร้อยลูกประคำ", "Prayer Beads"}, 1},
			},
		},
		{
			ID: "paladin", Role: "hybrid", Art: "",
			HeroName:  L{"เบรนนา ดอว์นเมียร์", "Brenna Dawnmere"},
			ClassName: L{"พาลาดิน", "Paladin"},
			Blurb:     L{"พาลาดินพูดตรงใจใหญ่ ผู้สาบานจะปกป้องคนอ่อนแอ และถือคำสาบานนั้นแบบตรงตัวทุกคำ", "A blunt, big-hearted paladin sworn to protect the weak, who takes that oath word for word."},
			Stats:     Stats{HPMax: 24, Str: 7, End: 7, Agi: 2, MPMax: 8},
			Skills: []RosterSkill{
				{"life-bond", L{"สายใยชีวิต", "Life Bond"}, L{"ถ่ายเลือดของตัวเองให้เพื่อนที่เราแตะอยู่ ได้สูงสุด 5 แต้ม", "Give up to 5 of your own HP to an ally you are touching."}, 1},
				{"radiant-smite", L{"ฟาดแสงศักดิ์สิทธิ์", "Radiant Smite"}, L{"การฟันโดนครั้งถัดไปแรงขึ้น +d8 และแรงเป็น 2 เท่าถ้าเป็นซากศพเดินได้หรือปีศาจ", "Your next melee hit deals +d8 damage, doubled against undead and demons."}, 2},
				{"holy-aegis", L{"โล่แสงศักดิ์สิทธิ์", "Holy Aegis"}, L{"สร้างโล่แสงคลุมเพื่อน 1 คนในระยะ 4 ช่อง ทุกครั้งที่เพื่อนโดนตี ความเสียหายลดลง 5 แต้ม จนถึงตาถัดไปของเรา", "Wrap one ally within 4 squares in a shield of light: every hit on them deals 5 less damage until your next turn."}, 4},
			},
			Items: []RosterItem{
				{L{"ดาบใหญ่แห่งคำปฏิญาณ", "Oathbound Greatsword"}, 1},
				{L{"เกราะอกขลิบเงิน", "Silver-Trimmed Breastplate"}, 1},
				{L{"ม้วนคำปฏิญาณผนึกครั่ง", "Wax-Sealed Oath Scroll"}, 1},
				{L{"ชาดอกตะวัน", "Sunpetal Tea"}, 2},
			},
		},
		{
			ID: "barbarian", Role: "damage", Art: "",
			HeroName:  L{"บาตู เรดเมน", "Batu Redmane"},
			ClassName: L{"นักรบเถื่อน", "Barbarian"},
			Blurb:     L{"นักรบจากทุ่งหญ้ากว้างที่ออกจากเผ่ามาตามหาศึกที่คู่ควรจะถูกแต่งเป็นเพลง", "A steppe warrior who left his clan to find a fight worthy of a song."},
			Stats:     Stats{HPMax: 32, Str: 9, End: 4, Agi: 2, MPMax: 6},
			Skills: []RosterSkill{
				{"war-cry", L{"คำรามศึก", "War Cry"}, L{"คำรามจนขนลุก ศัตรูอ่อนแอในระยะ 3 ช่องจะหนีหรือลังเลไป 1 ตา (DM ตัดสินว่าตัวไหนนับว่าอ่อนแอ)", "Let out a terrifying roar: weak enemies within 3 squares flee or hesitate for 1 turn (the DM decides who counts as weak)."}, 1},
				{"blood-rage", L{"โหมโทสะ", "Blood Rage"}, L{"นาน 3 ตา ทุกครั้งที่ฟันโดนแรงขึ้น +3 แต่ทอยที่ใช้ความว่องไว (AGI) −2", "For 3 turns your melee hits deal +3 damage, but your Agility rolls get −2."}, 2},
				{"unstoppable-charge", L{"พุ่งชนไม่ยั้ง", "Unstoppable Charge"}, L{"พุ่งเป็นเส้นตรงไกลสูงสุด 6 ช่อง ไม่มีอะไรขวางได้ แล้วฟันศัตรูที่ปลายทาง ถ้าโดนแรงขึ้น +d8", "Charge up to 6 squares in a straight line, nothing can block you, then strike the enemy at the end. On a hit it deals +d8 damage."}, 4},
			},
			Items: []RosterItem{
				{L{"ขวานยักษ์ตระกูลเรดเมน", "Redmane Greataxe"}, 1},
				{L{"เสื้อคลุมหนังหมาป่า", "Wolf-Pelt Cloak"}, 1},
				{L{"ขวานเล็กสำหรับขว้าง", "Throwing Hatchet"}, 3},
				{L{"ถุงนมม้าหมัก", "Skin of Fermented Mare's Milk"}, 1},
			},
		},
		{
			ID: "bard", Role: "support", Art: "",
			HeroName:  L{"โคซิโม ลาร์ก", "Cosimo Lark"},
			ClassName: L{"กวีพเนจร", "Bard"},
			Blurb:     L{"นักแต่งเพลงพเนจรที่สาบานว่าทุกเรื่องที่เล่าเป็นเรื่องจริง โดยเฉพาะเรื่องของตัวเอง", "A wandering songwriter who swears every tale he tells is true, especially the ones about himself."},
			Stats:     Stats{HPMax: 24, Str: 3, End: 4, Agi: 8, MPMax: 10},
			Skills: []RosterSkill{
				{"silver-tongue", L{"ปากหวาน", "Silver Tongue"}, L{"ทอยครั้งถัดไปเพื่อเกลี้ยกล่อม เอาใจ หรือหลอกใครสักคน ได้ +3", "Get +3 on your next roll to persuade, charm, or fool someone."}, 1},
				{"rousing-tune", L{"เพลงปลุกใจ", "Rousing Tune"}, L{"บรรเลงเพลงฮึกเหิม เพื่อนที่ได้ยินสูงสุด 2 คนได้ +2 ในการทอยครั้งถัดไป", "Play a stirring tune: up to 2 allies who can hear it get +2 on their next roll."}, 2},
				{"encore", L{"อีกสักรอบ!", "Encore!"}, L{"เลือกเพื่อน 1 คนที่เพิ่งเล่นจบตา เพื่อนได้เล่นอีก 1 ตาเต็มๆ ทันที", "Pick an ally who just finished their turn: they immediately take another full turn."}, 5},
			},
			Items: []RosterItem{
				{L{"พิณลูทสายเงิน", "Silver-Stringed Lute"}, 1},
				{L{"ดาบเรเพียร์", "Rapier"}, 1},
				{L{"สำรับไพ่ทำเครื่องหมาย", "Marked Deck of Cards"}, 1},
				{L{"หมวกขนนกสุดหรู", "Plumed Hat"}, 1},
				{L{"ขวดไวน์น้ำผึ้ง", "Bottle of Honeyed Wine"}, 2},
			},
		},
		{
			ID: "druid", Role: "hybrid", Art: "",
			HeroName:  L{"มาลัย แห่งพงไพร", "Malai of the Deepwood"},
			ClassName: L{"ดรูอิด", "Druid"},
			Blurb:     L{"ผู้พิทักษ์ป่าลึกที่พูดเสียงนุ่ม แต่กลายร่างเป็นหมีได้เมื่อจำเป็น", "A guardian of the deep jungle who speaks softly and becomes a bear when she must."},
			Stats:     Stats{HPMax: 20, Str: 7, End: 4, Agi: 3, MPMax: 14},
			Skills: []RosterSkill{
				{"thorn-lash", L{"แส้หนาม", "Thorn Lash"}, L{"ฟาดแส้เถาหนามใส่ศัตรูในระยะ 3 ช่อง ถ้าโดน เสียหาย d4 และถูกดึงเข้ามาหาเรา 2 ช่อง", "Lash an enemy within 3 squares with a thorny vine. On a hit it takes d4 damage and is pulled 2 squares toward you."}, 1},
				{"entangling-roots", L{"รากไม้พันธนาการ", "Entangling Roots"}, L{"รากไม้พุ่งขึ้นจากพื้นรัดศัตรู 1 ตัวในระยะ 5 ช่อง มันเดินไม่ได้ 2 ตา (แต่ยังโจมตีได้)", "Roots burst up around one enemy within 5 squares: it can't move for 2 turns (it can still attack)."}, 2},
				{"bear-form", L{"ร่างหมีป่า", "Bear Form"}, L{"แปลงร่างเป็นหมีป่าตัวใหญ่นาน 3 ตา ได้เลือดเสริม 10 แต้ม (หายเมื่อคืนร่าง) และนับ STR เป็น 10 ระหว่างเป็นหมีใช้สกิลอื่นไม่ได้", "Become a great bear for 3 turns: gain 10 extra HP (gone when you change back) and your Strength counts as 10. You can't use other skills as a bear."}, 4},
			},
			Items: []RosterItem{
				{L{"ไม้เท้าเถาวัลย์มีชีวิต", "Living Vine Staff"}, 1},
				{L{"ถุงเมล็ดพันธุ์ป่า", "Pouch of Wild Seeds"}, 1},
				{L{"เครื่องรางด้วงหยก", "Jade Beetle Charm"}, 1},
				{L{"มัดตะไคร่จันทร์", "Bundle of Moonmoss"}, 3},
			},
		},
		{
			ID: "necromancer", Role: "control", Art: "",
			HeroName:  L{"อีโว กริมเซล", "Ivo Grimsel"},
			ClassName: L{"เนโครแมนเซอร์", "Necromancer"},
			Blurb:     L{"ลูกชายคนเฝ้าสุสาน ที่ฟังเสียงคนตายเพราะคนเป็นไม่เคยฟังเขา", "A gravekeeper's son who listens to the dead because the living never listened to him."},
			Stats:     Stats{HPMax: 16, Str: 3, End: 7, Agi: 3, MPMax: 18},
			Skills: []RosterSkill{
				{"grave-whisper", L{"กระซิบจากหลุม", "Grave Whisper"}, L{"ถามศพ 1 คำถาม ศพจะตอบสั้นๆ ตามที่มันรู้ตอนยังมีชีวิต", "Ask a corpse one question. It answers in a few words with what it knew in life."}, 1},
				{"life-siphon", L{"ดูดชีวิต", "Life Siphon"}, L{"ดูดพลังชีวิตศัตรูในระยะ 4 ช่อง ถ้าโดน มันเสียหาย d6 และเราฟื้นเลือดเท่ากับที่ดูดได้", "Drain an enemy within 4 squares. On a hit it takes d6 damage and you regain that much HP."}, 2},
				{"raise-fallen", L{"ปลุกซากผู้ล่วงลับ", "Raise the Fallen"}, L{"ปลุกศพในระยะ 3 ช่องเป็นโครงกระดูกรับใช้ (เลือด 6 ตีแรง d4) ทำตามคำสั่งเราจนกว่าจะถูกทำลายหรือจบฉาก", "Raise a corpse within 3 squares as a skeleton servant (6 HP, hits for d4) that obeys you until destroyed or the scene ends."}, 5},
			},
			Items: []RosterItem{
				{L{"ตะเกียงกะโหลก", "Skull Lantern"}, 1},
				{L{"เคียวด้ามกระดูก", "Bone-Handled Sickle"}, 1},
				{L{"โถดินจากหลุมศพ", "Jar of Grave Dirt"}, 1},
				{L{"สมุดรายชื่อผู้วายชนม์", "Ledger of the Departed"}, 1},
				{L{"เสื้อคลุมขนอีกา", "Crow-Feather Cloak"}, 1},
			},
		},
	}
}
