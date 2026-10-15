// Skill descriptions based on the SRD 5.2.1 Skills table (CC-BY-4.0), written in our own words.
export type SkillText = { name: string; desc: string; examples: string[] }

export const SRD_SKILLS: Record<string, { th: SkillText; en: SkillText }> = {
  athletics: {
    th: {
      name: 'กรีฑา',
      desc: 'ใช้พละกำลังของร่างกายเพื่อปีน กระโดด ว่ายน้ำ หรือออกแรงทำลายสิ่งของ',
      examples: ['กระโดดข้ามหุบเหวที่กว้างกว่าปกติ', 'ว่ายน้ำฝ่ากระแสน้ำเชี่ยว', 'ปีนกำแพงหินที่ลื่น', 'พังประตูไม้ให้เปิดออก'],
    },
    en: {
      name: 'Athletics',
      desc: 'Using raw physical power to climb, jump, swim, or force things apart.',
      examples: ['Leap across a wider-than-usual gap', 'Keep swimming in rough water', 'Climb a slippery rock wall', 'Smash open a wooden door'],
    },
  },
  acrobatics: {
    th: {
      name: 'กายกรรม',
      desc: 'ความว่องไวและการทรงตัว ช่วยให้ยืนหยัดได้ในจังหวะคับขันหรือทำท่าผาดโผน',
      examples: ['ทรงตัวเดินบนเชือกหรือขอบแคบ ๆ', 'ไม่ล้มบนพื้นน้ำแข็งหรือดาดฟ้าเรือที่โคลง', 'ตีลังกาหลบผ่านศัตรู'],
    },
    en: {
      name: 'Acrobatics',
      desc: 'Agility and balance that keep you on your feet in tricky spots or let you pull off a stunt.',
      examples: ['Walk along a narrow ledge or rope', 'Stay standing on ice or a rocking ship deck', 'Flip past an enemy'],
    },
  },
  'sleight-of-hand': {
    th: {
      name: 'ความคล่องแคล่วของมือ',
      desc: 'ใช้มือที่ไวและแนบเนียนทำสิ่งต่าง ๆ โดยไม่ให้ใครทันสังเกต',
      examples: ['ล้วงกระเป๋าเงินของคนอื่น', 'ซ่อนมีดเล่มเล็กไว้ในแขนเสื้อ', 'เล่นกลเสกเหรียญให้หายไป'],
    },
    en: {
      name: 'Sleight of Hand',
      desc: 'Quick, careful hands that do things without anyone noticing.',
      examples: ["Pick someone's pocket", 'Hide a small knife up your sleeve', 'Make a coin vanish with a trick'],
    },
  },
  stealth: {
    th: {
      name: 'การลอบเร้น',
      desc: 'เคลื่อนไหวอย่างเงียบเชียบและหลบซ่อนตัวไม่ให้ใครเห็นหรือได้ยิน',
      examples: ['ย่องผ่านยามที่กำลังเฝ้าประตู', 'หลบซ่อนหลังลังไม้ในเงามืด', 'แอบตามคนอื่นโดยไม่ให้รู้ตัว'],
    },
    en: {
      name: 'Stealth',
      desc: 'Moving quietly and staying hidden so nobody sees or hears you.',
      examples: ['Sneak past a guard at the gate', 'Hide behind crates in the shadows', 'Follow someone without being noticed'],
    },
  },
  arcana: {
    th: {
      name: 'ศาสตร์เวทมนตร์',
      desc: 'ความรู้เรื่องเวทมนตร์ คาถา ของวิเศษ และโลกหรือมิติอื่น ๆ',
      examples: ['จำได้ว่าคาถานี้คืออะไร', 'บอกได้ว่าอักขระบนประตูเป็นเวทแบบไหน', 'รู้เรื่องราวของมิติอื่น'],
    },
    en: {
      name: 'Arcana',
      desc: 'Knowledge about spells, magic items, and other planes of existence.',
      examples: ['Recognize which spell is being cast', 'Identify the magic in runes on a door', 'Recall facts about another plane'],
    },
  },
  history: {
    th: {
      name: 'ประวัติศาสตร์',
      desc: 'ความรู้เรื่องเหตุการณ์ในอดีต บุคคลสำคัญ อาณาจักร และวัฒนธรรมต่าง ๆ',
      examples: ['จำได้ว่าซากปราสาทนี้เคยเป็นของใคร', 'รู้ที่มาของสงครามครั้งเก่า', 'จำตราประจำตระกูลขุนนางได้'],
    },
    en: {
      name: 'History',
      desc: 'Knowledge of past events, famous people, nations, and cultures.',
      examples: ['Recall who once owned these castle ruins', 'Know how an old war started', "Recognize a noble family's crest"],
    },
  },
  investigation: {
    th: {
      name: 'การสืบสวน',
      desc: 'ค้นหาเบาะแสและใช้เหตุผลปะติดปะต่อว่าเกิดอะไรขึ้นหรือสิ่งของทำงานอย่างไร',
      examples: ['ค้นหาข้อมูลที่ซ่อนอยู่ในหนังสือเก่า', 'หากลไกเปิดประตูลับ', 'วิเคราะห์ว่ากับดักทำงานอย่างไร'],
    },
    en: {
      name: 'Investigation',
      desc: 'Searching for clues and reasoning out what happened or how something works.',
      examples: ['Dig up hidden facts in old books', 'Find the switch for a secret door', 'Work out how a trap is triggered'],
    },
  },
  nature: {
    th: {
      name: 'ธรรมชาติ',
      desc: 'ความรู้เรื่องภูมิประเทศ พืช สัตว์ และสภาพอากาศ',
      examples: ['รู้ว่าเห็ดชนิดนี้กินได้หรือมีพิษ', 'จำได้ว่าสัตว์ตัวนี้มีนิสัยอย่างไร', 'ดูท้องฟ้าแล้วรู้ว่าพายุกำลังมา'],
    },
    en: {
      name: 'Nature',
      desc: 'Knowledge of terrain, plants, animals, and weather.',
      examples: ['Tell if a mushroom is safe or poisonous', 'Recall how a wild beast behaves', 'Read the sky and see a storm coming'],
    },
  },
  religion: {
    th: {
      name: 'ศาสนา',
      desc: 'ความรู้เรื่องเทพเจ้า พิธีกรรมทางศาสนา และสัญลักษณ์ศักดิ์สิทธิ์',
      examples: ['จำได้ว่าสัญลักษณ์นี้เป็นของเทพองค์ไหน', 'รู้ความหมายของพิธีกรรมในวิหาร', 'รู้จักลัทธิลึกลับที่บูชาสิ่งชั่วร้าย'],
    },
    en: {
      name: 'Religion',
      desc: 'Knowledge of gods, religious rituals, and holy symbols.',
      examples: ['Identify which god a symbol belongs to', 'Understand a temple ritual', 'Recognize the signs of a dark cult'],
    },
  },
  'animal-handling': {
    th: {
      name: 'การดูแลสัตว์',
      desc: 'ทำให้สัตว์สงบลง ฝึกสัตว์ หรือทำให้สัตว์ทำตามที่เราต้องการ',
      examples: ['ปลอบม้าที่ตื่นตกใจให้สงบ', 'บังคับม้าให้กระโดดข้ามสิ่งกีดขวาง', 'ทำให้สุนัขเฝ้าบ้านไม่เห่า'],
    },
    en: {
      name: 'Animal Handling',
      desc: 'Calming, training, or guiding an animal to do what you want.',
      examples: ['Calm a frightened horse', 'Get your mount to jump an obstacle', 'Keep a guard dog from barking'],
    },
  },
  insight: {
    th: {
      name: 'การหยั่งรู้',
      desc: 'อ่านใจคนจากท่าทางและน้ำเสียง เพื่อรู้ว่าเขารู้สึกอย่างไรหรือคิดจะทำอะไร',
      examples: ['ดูออกว่าพ่อค้ากำลังโกหก', 'รู้ว่าคนตรงหน้ากำลังกลัวหรือโกรธ', 'เดาว่าอีกฝ่ายต้องการอะไรจริง ๆ'],
    },
    en: {
      name: 'Insight',
      desc: "Reading someone's body language and tone to tell how they feel or what they intend.",
      examples: ['Notice a merchant is lying', 'Tell whether someone is scared or angry', 'Guess what someone really wants'],
    },
  },
  medicine: {
    th: {
      name: 'การแพทย์',
      desc: 'ดูแลผู้บาดเจ็บ วินิจฉัยอาการป่วย และดูว่าใครตายเพราะอะไร',
      examples: ['ห้ามเลือดเพื่อนที่หมดสติ', 'บอกได้ว่าคนนี้ป่วยเป็นโรคอะไร', 'ตรวจศพแล้วรู้สาเหตุการตาย'],
    },
    en: {
      name: 'Medicine',
      desc: 'Caring for the wounded, diagnosing illness, and telling what killed someone.',
      examples: ['Stabilize a fallen ally', 'Figure out what illness someone has', 'Examine a body to find the cause of death'],
    },
  },
  perception: {
    th: {
      name: 'การรับรู้',
      desc: 'ใช้ตา หู และประสาทสัมผัสต่าง ๆ สังเกตสิ่งที่คนทั่วไปมักมองข้าม',
      examples: ['ได้ยินเสียงฝีเท้าย่องเข้ามา', 'เห็นศัตรูที่ซุ่มอยู่ในพุ่มไม้', 'สังเกตเห็นรอยแยกของประตูลับบนผนัง', 'ได้กลิ่นควันไฟจากที่ไกล'],
    },
    en: {
      name: 'Perception',
      desc: 'Using your senses to notice things that are easy to miss.',
      examples: ['Hear footsteps sneaking closer', 'Spot an enemy hiding in the bushes', 'Notice the outline of a secret door', 'Smell smoke from far away'],
    },
  },
  survival: {
    th: {
      name: 'การเอาชีวิตรอด',
      desc: 'ใช้ชีวิตในป่าเขา ตามรอย หาอาหาร และหลีกเลี่ยงอันตรายจากธรรมชาติ',
      examples: ['ตามรอยเท้าของสัตว์หรือศัตรู', 'หาอาหารและน้ำในป่า', 'หาทางกลับเมื่อหลงทาง', 'หลบเลี่ยงทรายดูดหรือหิมะถล่ม'],
    },
    en: {
      name: 'Survival',
      desc: 'Getting by in the wild: tracking, foraging, and avoiding natural dangers.',
      examples: ['Follow tracks of an animal or foe', 'Find food and water in the wild', 'Find your way when lost', 'Avoid quicksand or an avalanche'],
    },
  },
  deception: {
    th: {
      name: 'การหลอกลวง',
      desc: 'ทำให้คนอื่นเชื่อในสิ่งที่ไม่จริง ทั้งด้วยคำพูดและการปลอมตัว',
      examples: ['โกหกยามว่าเป็นแขกที่ได้รับเชิญ', 'ปลอมตัวเป็นคนอื่นอย่างแนบเนียน', 'แกล้งทำเป็นยอมแพ้เพื่อหลอกศัตรู'],
    },
    en: {
      name: 'Deception',
      desc: 'Making others believe something untrue, through words or disguise.',
      examples: ["Convince a guard you're an invited guest", 'Pass yourself off in a disguise', 'Fake a surrender to trick a foe'],
    },
  },
  intimidation: {
    th: {
      name: 'การข่มขู่',
      desc: 'ใช้ท่าทีน่าเกรงขามหรือคำขู่ เพื่อให้คนอื่นยอมทำตาม',
      examples: ['ขู่โจรให้บอกที่ซ่อนของหัวหน้า', 'จ้องตาจนอีกฝ่ายยอมถอยไป', 'ทำให้ฝูงชนแตกกระจายด้วยเสียงตะโกน'],
    },
    en: {
      name: 'Intimidation',
      desc: 'Using a fearsome presence or threats to make others do what you want.',
      examples: ["Threaten a bandit into revealing the boss's hideout", 'Stare someone down until they back off', 'Scatter a crowd with a shout'],
    },
  },
  performance: {
    th: {
      name: 'การแสดง',
      desc: 'สร้างความบันเทิงให้ผู้ชมด้วยการแสดง เล่าเรื่อง ร้องเพลง เล่นดนตรี หรือเต้นรำ',
      examples: ['เล่นดนตรีในโรงเตี๊ยมเพื่อแลกเงิน', 'เล่านิทานจนทุกคนตั้งใจฟัง', 'เต้นรำเบี่ยงความสนใจของยาม'],
    },
    en: {
      name: 'Performance',
      desc: 'Entertaining an audience by acting, storytelling, music, or dance.',
      examples: ['Play music at a tavern for coin', 'Tell a story that holds the room', 'Dance to distract the guards'],
    },
  },
  persuasion: {
    th: {
      name: 'การโน้มน้าว',
      desc: 'ใช้เหตุผลและความจริงใจพูดให้คนอื่นเห็นด้วยหรือยอมช่วย',
      examples: ['ขอให้ชาวบ้านช่วยชี้ทาง', 'ต่อรองราคากับพ่อค้า', 'พูดให้สองฝ่ายเลิกทะเลาะกัน'],
    },
    en: {
      name: 'Persuasion',
      desc: 'Using honest reasons and good manners to win someone over or get their help.',
      examples: ['Ask villagers to point the way', 'Haggle with a merchant', 'Talk two sides out of a fight'],
    },
  },
}
