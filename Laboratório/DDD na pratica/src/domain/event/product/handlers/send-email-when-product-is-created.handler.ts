import EventHandlerInterface from "../../@shared/event-handler.interface";
import { EventInterface } from "../../@shared/event.interface";

export default class SendEmailWhenProductIsCreatedHandler implements EventHandlerInterface {
  handle(event: EventInterface): void {
    throw new Error("Method not implemented.");
  }
}