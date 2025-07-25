local es9Pan(from) = [

  // The ES-9 has an internal 8x8 mixer. I use the EC4 to control pan for the
  // mixes that pertain to my main and headphones output.
  //
  // https://www.expert-sleepers.co.uk/es9.html

  {

    name: 'c%d' % i,

    ec: {
      channel: 9,
      number: i,
      lower: 0,
      upper: 127,
      display: '±63',
      type: 'CCAb',
      mode: 'Acc1',
    },
    pb: {
      channel: 9,
      number: i,
      lower: 64,
      upper: 64,
      display: 'On',
      type: 'CC',
      mode: 'Key',
    },
  }
  for i in std.range(from, from + 7)
];

local sevenPC(from) = [

  // Program change messages for a Crumar Seven

  {

    name: '  %02d' % (i + 1),

    // encoders control volume
    ec: {
      channel: 1,
      number: 7,
      lower: 0,
      upper: 127,
      display: '127',
      type: 'CCAb',
      mode: 'Acc1',
    },

    // push buttons send program change message
    pb: {
      channel: 1,
      lower: i,
      upper: i,
      display: 'Off',
      type: 'PrgC',
      mode: 'Key',
    },
  }
  for i in std.range(from, from + 15)
];

local liveVolume(from) = [

  // Ableton Live track volume control

  {

    name: '  %02d' % (i + 1),

    // encoders control volume sliders
    ec: {
      channel: 13,
      number: i,
      lower: 0,
      upper: 127,
      display: '127',
      type: 'CCAb',
      mode: 'Acc3',
    },

    // push buttons toggle mute
    pb: {
      channel: 13,
      number: i + 64,
      lower: 0,
      upper: 1,
      display: 'Off',
      type: 'CC',
      mode: 'Togl',
    },
  }
  for i in std.range(from, from + 15)
];

local emptyGroup = {

  // List of 16 empty objects to account for the groups I don't use.
  //
  // ec4-config will set useless defaults.

  settings: [
    {}
    for i in std.range(1, 16)
  ],
};

{
  setups: [

    // Setup 1 is for an es-9...
    {
      name: 'es-9',
      groups: [
        {
          name: 'phon',
          settings: es9Pan(24) + es9Pan(88),
        },
        {
          name: 'main',
          settings: es9Pan(8) + es9Pan(72),
        },
      ] + [

        // Groups 3 to 16 are empty...
        emptyGroup
        for i in std.range(3, 16)
      ],
    },
  ] + [

    // Setups 2 to 4 are empty...
    {
      groups: [
        emptyGroup
        for i in std.range(1, 16)
      ],
    }
    for i in std.range(2, 4)
  ] + [

    // Setup 6 is for a crumar seven...
    {
      name: 'sevn',
      groups: [
        {
          name: '->16',
          settings: sevenPC(0),
        },
        {
          name: '->32',
          settings: sevenPC(16),
        },
      ] + [

        // Groups 3 to 16 are empty...
        emptyGroup
        for i in std.range(3, 16)
      ],
    },
  ] + [

    // Setups 6 to 12 are empty...
    {
      groups: [
        emptyGroup
        for i in std.range(1, 16)
      ],
    }
    for i in std.range(6, 12)
  ] + [

    // Setup 13 is for Live...
    {
      name: 'live',
      groups: [
        {
          // 1-16
          name: '->16',
          settings: liveVolume(0),
        },
        {
          // 17-32
          name: '->32',
          settings: liveVolume(16),
        },
        {
          // 33-48
          name: '->48',
          settings: liveVolume(32),
        },
        {
          // 49-64
          name: '->64',
          settings: liveVolume(48),
        },
      ] + [

        // Groups 5 to 16 are empty...
        emptyGroup
        for i in std.range(5, 16)
      ],
    },
  ] + [

    // Setups 14 to 16 are empty...
    {
      groups: [
        emptyGroup
        for i in std.range(1, 16)
      ],
    }
    for i in std.range(14, 16)
  ],
}
