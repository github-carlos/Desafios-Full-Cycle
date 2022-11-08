import { Sequelize } from "sequelize-typescript";
import Product from "../../domain/entity/product";
import ProductModel from "../db/sequelize/models/product";
import ProductRepository from "./product.repository";

describe("#ProductRepository", () => {

  let sequelize: Sequelize;

  beforeEach(async() => {
    sequelize = new Sequelize({
      dialect: 'sqlite',
      storage: ':memory:',
      logging: false,
      sync: { force: true },
    })

    sequelize.addModels([ProductModel]);
    await sequelize.sync();
  });

  afterEach(async () => {
    await sequelize.close();
  });

  it("should create a new product", async () => {
    const repository = new ProductRepository();

    const newProduct = new Product('123', 'Product 1', 100);
    await repository.create(newProduct);

    const productFromDB = await ProductModel.findOne({where: {id: '123'}});

    expect(productFromDB.toJSON()).toStrictEqual({id: '123', name: 'Product 1', price: 100});
  });

  it("should update a product", async () => {
    const repository = new ProductRepository();
    const newProduct = new Product('1', 'Product 1', 100);
    
    await repository.create(newProduct);

    const productFromDB = await ProductModel.findOne({where: {id: '1'}});
    expect(productFromDB.toJSON()).toStrictEqual({id: '1', name: 'Product 1', price: 100});

    newProduct.changeName('Product New Name');
    newProduct.changePrice(200);
    await repository.update(newProduct);

    const updatedProductFromDB = await ProductModel.findOne({where: {id: newProduct.id}});
    expect(updatedProductFromDB.toJSON()).toStrictEqual({id: '1', name: 'Product New Name', price: 200});
  });
});