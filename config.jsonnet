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

local liveVolume(from) = [

  // Ableton Live volume slider control

  {

    name: 'v%02d' % (i - from + 1),

    ec: {
      channel: 13,
      number: i,
      lower: 0,
      upper: 127,
      display: '127',
      type: 'CCAb',
      mode: 'Acc3',
    },
    pb: {
      channel: 13,
      number: i + 64,
      lower: 0,
      upper: 1,
      display: 'Off',
      type: 'CC',
      mode: 'Key',
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

    // Setups 2 to 12 are empty...
    {
      groups: [
        emptyGroup
        for i in std.range(1, 16)
      ],
    }
    for i in std.range(2, 12)
  ] + [

    // Setup 13 is for Live...
    {
      name: 'live',
      groups: [
        {
          // Volume: 0-15
          // Mute: 64-79
          name: 'v+m',
          settings: liveVolume(0),
        },
      ] + [

        // Groups 2 to 16 are empty...
        emptyGroup
        for i in std.range(2, 16)
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
