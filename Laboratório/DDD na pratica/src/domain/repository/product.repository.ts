import Product from "../entity/product";
import Repository from "./repository.interface";

export interface ProductRepository extends Repository<Product> {
  findByName(name: string): Promise<Product>;
}