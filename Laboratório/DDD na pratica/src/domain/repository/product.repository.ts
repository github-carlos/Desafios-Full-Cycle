import Product from "../entity/product";
import Repository from "./repository.interface";

export interface ProductRepositoryInterface extends Repository<Product> {
  findByName(name: string): Promise<Product>;
}