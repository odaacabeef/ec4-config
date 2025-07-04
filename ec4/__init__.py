from __future__ import absolute_import, division, print_function, unicode_literals
from _Generic.GenericScript import TRANSPORT_BUTTON_SPECIFICATIONS, GenericScript
from ableton.v3.control_surface import MapMode

def create_instance(c_instance):
    """
    "Generic" control surface config intended for Faderfox EC4
    https://github.com/gluon/AbletonLive12_MIDIRemoteScripts/tree/main/_Generic

    Volume Controls:
    - CC 0-63 control volume for tracks 1-64
    - All volume controls use channel 13

    Mute Controls:
    - CC 64-127 control mute for tracks 1-64
    - All mute controls use channel 13
    """

    # Map modes
    device_map_mode = MapMode.Absolute
    volume_map_mode = MapMode.Absolute

    # Device controls (not used in this implementation)
    device_controls = [(-1, -1)] * 16

    # Transport controls (not used in this implementation)
    transport_controls = {key.upper(): -1 for key in TRANSPORT_BUTTON_SPECIFICATIONS}

    # Volume controls: CC 0-63 for tracks 1-64, all on channel 13
    volume_controls = []
    for i in range(64):
        volume_controls.append((i, 12))  # (CC number, channel) - channel 13 = 12 (0-based)

    # Track arm controls (not used in this implementation)
    trackarm_controls = [-1] * 64  # Increased to 64 to match volume controls

    # Bank controls (not used in this implementation)
    bank_controls = {
        'ONOFF': -1,
        'TOGGLELOCK': -1,
        'NEXTBANK': -1,
        'PREVBANK': -1,
        'BANK1': -1,
        'BANK2': -1,
        'BANK3': -1,
        'BANK4': -1,
        'BANK5': -1,
        'BANK6': -1,
        'BANK7': -1,
        'BANK8': -1
    }

    # Controller descriptions
    controller_descriptions = {"CHANNEL": 12}  # Global channel 13 = 12 (0-based)

    # Mixer options
    mixer_options = {
        'NUMSENDS': 2,
        'INVERTMUTELEDS': True,
        'SEND1': [-1] * 64,
        'SEND2': [-1] * 64,
        'MUTE': [-1] * 64,
        'SOLO': [-1] * 64,
        'SELECT': [-1] * 64,
        'SENDMAPMODE': MapMode.Absolute,
        'MASTERVOLUME': -1,
        'MASTERVOLUMECHANNEL': -1,
        'CUEVOLUME': -1,
        'CUEVOLUMECHANNEL': -1,
        'CROSSFADER': -1,
        'CROSSFADERCHANNEL': -1,
        'CROSSFADERMAPMODE': MapMode.Absolute
    }

    # Mute controls: CC 64-127 for tracks 1-64, all on channel 13
    for i in range(64):
        mixer_options["MUTE"][i] = i + 64  # CC 64-127

    # Create and return the script instance
    return type("ec4", (GenericScript,), {})(c_instance,
        device_map_mode,
        volume_map_mode,
        tuple(device_controls),
        transport_controls,
        tuple(volume_controls),
        tuple(trackarm_controls),
        bank_controls,
        controller_descriptions,
        mixer_options)
