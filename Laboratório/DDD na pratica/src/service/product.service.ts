import Product from "../entity/product";

export default class ProduceService {
  static increasePrice(products: Array<Product>, percentage: number) {
    products.forEach(product => {
      const oldPrice = product.price;
      const newPrice = oldPrice + oldPrice * (percentage / 100);
      product.changePrice(newPrice);
    });
  }
}