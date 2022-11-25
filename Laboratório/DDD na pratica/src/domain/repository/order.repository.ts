import Order from "../entity/order";
import Repository from "./repository.interface";

export interface OrderRepositoryInterface extends Repository<Order> {}