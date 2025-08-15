package keys

import (
	"github.com/adm87/finch-core/events"
	"github.com/adm87/finch-core/hash"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type KeyPhase uint8

func (p KeyPhase) Has(phase KeyPhase) bool {
	return p&phase != 0
}

const (
	KeyPhaseBegin KeyPhase = 1 << iota
	KeyPhasePress
	KeyPhaseRelease
)

var keyActions = make(map[ebiten.Key]map[KeyPhase]hash.HashSet[events.Command])

func RegisterAction(key ebiten.Key, phase KeyPhase, command events.Command) error {
	if _, exists := keyActions[key]; !exists {
		create_key_action(key)
	}

	if phase.Has(KeyPhaseBegin) {
		keyActions[key][KeyPhaseBegin].Add(command)
	}
	if phase.Has(KeyPhasePress) {
		keyActions[key][KeyPhasePress].Add(command)
	}
	if phase.Has(KeyPhaseRelease) {
		keyActions[key][KeyPhaseRelease].Add(command)
	}

	return nil
}

func UnregisterAction(key ebiten.Key, phase KeyPhase, command events.Command) error {
	if _, exists := keyActions[key]; !exists {
		return nil
	}

	if phase.Has(KeyPhaseBegin) {
		keyActions[key][KeyPhaseBegin].Remove(command)
	}
	if phase.Has(KeyPhasePress) {
		keyActions[key][KeyPhasePress].Remove(command)
	}
	if phase.Has(KeyPhaseRelease) {
		keyActions[key][KeyPhaseRelease].Remove(command)
	}

	if keyActions[key][KeyPhaseBegin].IsEmpty() && keyActions[key][KeyPhasePress].IsEmpty() && keyActions[key][KeyPhaseRelease].IsEmpty() {
		delete(keyActions, key)
	}

	return nil
}

func Poll() error {
	for key, phases := range keyActions {
		switch {
		case inpututil.IsKeyJustPressed(key) && !phases[KeyPhaseBegin].IsEmpty():
			if err := execute_phase_commands(phases[KeyPhaseBegin]); err != nil {
				return err
			}
		case inpututil.IsKeyJustReleased(key) && !phases[KeyPhaseRelease].IsEmpty():
			if err := execute_phase_commands(phases[KeyPhaseRelease]); err != nil {
				return err
			}
		case ebiten.IsKeyPressed(key) && !phases[KeyPhasePress].IsEmpty():
			if err := execute_phase_commands(phases[KeyPhasePress]); err != nil {
				return err
			}
		}
	}
	return nil
}

func create_key_action(key ebiten.Key) {
	keyActions[key] = map[KeyPhase]hash.HashSet[events.Command]{
		KeyPhaseBegin:   make(hash.HashSet[events.Command]),
		KeyPhasePress:   make(hash.HashSet[events.Command]),
		KeyPhaseRelease: make(hash.HashSet[events.Command]),
	}
}

func execute_phase_commands(commands hash.HashSet[events.Command]) error {
	for command := range commands {
		if err := command.Execute(); err != nil {
			return err
		}
	}
	return nil
}
