import (
	"list"
	"strings"
)

setups: [..._#setup] & _#list16

_#setup: {
	name: _#name
	groups: [..._#group] & _#list16
}

_#group: {
	name: _#name
	settings: [..._#setting] & _#list16
}

_#setting: {
	name: _#name
	ec:   _#encoder
	pb:   _#pushButton
}

_#encoder: {
	channel: _#channel
	number:  _#cc
	lower:   _#cc
	upper:   _#cc
	display: *"Off" | "127" | "100" | "1000" | "±63" | "±50" | "±500" | "ONOF" | "9999"
	type:    "CCr1" | "CCr2" | *"CCAb" | "PrgC" | "CCAh" | "PBnd" | "AftT" | "Note" | "NRPN"
	mode:    "Div8" | "Div4" | "Div2" | "Acc0" | *"Acc1" | "Acc2" | "Acc3" | "LSp2" | "LSp4" | "LSp6"
}

_#pushButton: {
	channel: _#channel
	number:  _#cc
	lower:   _#cc
	upper:   _#cc
	display: "Off" | *"On"
	type:    "Off" | "Note" | *"CC" | "PrgC" | "PBnd" | "AftT" | "Grp" | "Set" | "Acc3" | "LSp6" | "Min" | "Max"
	mode:    *"Key" | "Togl"
}

_#list16:  list.MinItems(16) & list.MaxItems(16)
_#name:    strings.MaxRunes(4) | *"."
_#channel: >=1 & <=16 | *1
_#cc:      >=0 & <=127 | *0
