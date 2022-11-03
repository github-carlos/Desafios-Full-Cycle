import Product from "./product";

describe("#Product", () => {

   it("should throw error when id is empty", () => {
    expect(() => new Product("", "Product1", 100)).toThrowError();
  });
  
  it("should throw error when name is empty", () => {
    expect(() => new Product("123", "", 100)).toThrowError();
  });

  it("should throw error when price is negative", () => {
    expect(() => new Product("123", "PoductName", -2)).toThrowError();
  })

  it("should change name when changeName is called", () => {
    const product = new Product("123", "ProduceName", 100);
    const newName = "NewProductName";
    product.changeName(newName);
    expect(product.name).toBe(newName);
  })
  it("should change price when changePrice is called", () => {
    const product = new Product("123", "ProduceName", 100);
    const newPrice = 200;
    product.changePrice(newPrice);
    expect(product.price).toBe(newPrice);
  })
});