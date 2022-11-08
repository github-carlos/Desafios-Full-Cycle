import Product from "../../domain/entity/product";
import { ProductRepositoryInterface } from "../../domain/repository/product.repository";
import ProductModel from "../db/sequelize/models/product";

export default class ProductRepository implements ProductRepositoryInterface {
  async update(data: Product): Promise<void> {
    await ProductModel.update({name: data.name, price: data.price}, {where: {id: data.id}});
  }
  findByName(name: string): Promise<Product> {
    throw new Error("Method not implemented.");
  }
  async create(data: Product): Promise<void> {
    await ProductModel.create({id: data.id, name: data.name, price: data.price});
  }
  find(id: string): Promise<Product> {
    throw new Error("Method not implemented.");
  }
  findAll(): Promise<Product[]> {
    throw new Error("Method not implemented.");
  }
}