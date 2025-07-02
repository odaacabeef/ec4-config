local es9Mixer(from) = [

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
    {
      name: 'mdlr',
      groups: [
        {
          name: 'phon',
          settings: es9Mixer(24) + es9Mixer(88),
        },
        {
          name: 'main',
          settings: es9Mixer(8) + es9Mixer(72),
        },
      ] + [

        // Fill in the other 14 groups of the first setup.

        emptyGroup
        for i in std.range(3, 16)
      ],
    },
  ] + [

    // Fill in the other 15 setups.

    {
      groups: [
        emptyGroup
        for i in std.range(1, 16)
      ],
    }
    for i in std.range(2, 16)
  ],
}
