import { EventInterface } from "../@shared/event.interface";

export default class ProductCreatedEvent implements EventInterface {
  dataTimeOccurred: Date;
  eventData: any;

  constructor(data: any) {
    this.dataTimeOccurred = new Date();
    this.eventData = data;
  }
}