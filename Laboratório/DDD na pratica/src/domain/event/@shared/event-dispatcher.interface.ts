import EventHandlerInterface from "./event-handler.interface";
import { EventInterface } from "./event.interface";

export default interface EventDispatcherInterface {
  notify(event: EventInterface): void;
  register(eventName: string, handle: EventHandlerInterface): void;
  unregister(eventName: string, handle: EventHandlerInterface): void;
  unregisterAll(): void;
}