import Product from "../entity/product";
import ProduceService from "./product.service";

describe("#ProduceService", () => {
  it("should change all products price", () => {
    const product1 = new Product("produc1", "Product 1", 10);
    const product2 = new Product("produc2", "Product 2", 20);

    ProduceService.increasePrice([product1, product2], 100);

    expect(product1.price).toBe(20);
    expect(product2.price).toBe(40);
  });
});