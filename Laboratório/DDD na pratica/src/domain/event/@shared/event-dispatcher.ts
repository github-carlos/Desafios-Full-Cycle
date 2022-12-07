import EventDispatcherInterface from "./event-dispatcher.interface";
import EventHandlerInterface from "./event-handler.interface";
import { EventInterface } from "./event.interface";

export default class EventDispatcher implements EventDispatcherInterface {

  private eventHandlers: {[eventName: string]: EventHandlerInterface[]} = {};

  get getEventHandlers() {
    return this.eventHandlers;
  }

  notify(event: EventInterface): void {
    const handlers = this.eventHandlers[event.constructor.name]
    if (handlers) {
      handlers.forEach((handler) => {
        handler.handle(event);
      });
    }
  }
  register(eventName: string, handle: EventHandlerInterface): void {
    if (!this.eventHandlers[eventName]) {
      this.eventHandlers[eventName] = [];
    }
    this.eventHandlers[eventName].push(handle);
  }
  unregister(eventName: string, handle: EventHandlerInterface): void {
    const handlers = this.eventHandlers[eventName];

    if (!handlers || !handlers.length) {
      return;
    }

    const remainedHandlers = handlers.filter((registeredHandle) => registeredHandle !== handle);
    this.eventHandlers[eventName] = remainedHandlers;
  }
  unregisterAll(): void {
    this.eventHandlers = {};
  }

}