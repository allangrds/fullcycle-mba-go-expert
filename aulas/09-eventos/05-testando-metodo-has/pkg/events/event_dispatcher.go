package events

import "errors"

var ErrHandlerAlreadyRegistered = errors.New("handler already registered for event")

type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	// Aqui é para verificar se o mesmo handler já foi registrado para este evento.
	// Para isso, usamos:
	// - Um if para checar se já existe uma lista de handlers para este eventName;
	// - Um map (ed.handlers) para checar se já existe uma lista de handlers para este eventName;
	// - Um for para percorrer os handlers já cadastrados desse evento;
	// - Uma comparação (h == handler) para detectar duplicata. Se encontrar, retorna ErrHandlerAlreadyRegistered.
	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return ErrHandlerAlreadyRegistered
			}
		}
	}

	ed.handlers[eventName] = append(ed.handlers[eventName], handler)

	return nil
}

func (ed *EventDispatcher) Clear() {
	ed.handlers = make(map[string][]EventHandlerInterface)
}

func (ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return true
			}
		}
	}

	return false
}
