import { Sequelize } from "sequelize-typescript";
import Address from "../../domain/entity/address";
import Customer from "../../domain/entity/customer";
import Order from "../../domain/entity/order";
import OrderItem from "../../domain/entity/order_item";
import Product from "../../domain/entity/product";
import CustomerModel from "../db/sequelize/models/customer.model";
import { OrderItemModel, OrderModel } from "../db/sequelize/models/order.model";
import ProductModel from "../db/sequelize/models/product";
import CustomerRepository from "./customer.repository";
import OrderRepository from "./order.repository";
import ProductRepository from "./product.repository";

async function makeOrder(): Promise<Order> {
  const customerRepository = new CustomerRepository();
  const customer = new Customer("123", "Carlos");
  const address = new Address("Street 1", 1, "Zipcode", "City");
  customer.changeAddress(address);
  await customerRepository.create(customer);

  const productRepository = new ProductRepository();
  const product = new Product("123", "Product Name", 10);
  await productRepository.create(product);

  const orderItem = new OrderItem(
    "1",
    product.name,
    product.price,
    product.id,
    2
  );

  return new Order("123", customer.id, [orderItem]);
}

describe("Order repository test", () => {
  let sequelize: Sequelize;

  beforeEach(async () => {
    sequelize = new Sequelize({
      dialect: "sqlite",
      storage: ":memory:",
      logging: false,
      sync: { force: true },
    });

    await sequelize.addModels([
      CustomerModel,
      OrderModel,
      OrderItemModel,
      ProductModel,
    ]);
    await sequelize.sync();
  });

  afterEach(async () => {
    await sequelize.drop();
    await sequelize.close();
  });

  it("should create a new order", async () => {
    const order = await makeOrder();
    const orderRepository = new OrderRepository();
    await orderRepository.create(order);

    const orderModel = await OrderModel.findOne({
      where: { id: order.id },
      include: ["items"],
    });

    expect(orderModel.toJSON()).toStrictEqual({
      id: "123",
      customer_id: "123",
      total: order.total(),
      items: order.items.map((orderItem) => ({
        id: orderItem.id,
        name: orderItem.name,
        price: orderItem.price,
        quantity: orderItem.quantity,
        product_id: "123",
        order_id: "123",
      })),
    });
  });

  it("should find an order", async () => {
    const order = await makeOrder();
    const orderRepository = new OrderRepository();

    await orderRepository.create(order);

    const orderResult = await orderRepository.find(order.id);

    expect(orderResult).toStrictEqual(order);
  });

  it("should throw error when order not found", async () => {
    const orderRepository = new OrderRepository();

    expect(async () => {
      await orderRepository.find("123invalid");
    }).rejects.toThrow("Order not found");
  });

  it("should update an order", async () => {
    const order = await makeOrder();
    const orderRepository = new OrderRepository();

    await orderRepository.create(order);

    const customerRepository = new CustomerRepository();
    const customer = new Customer("333", "New Costumer");
    const address = new Address("Street 2", 2, "Zipcode 2", "City 2");
    customer.changeAddress(address);
    await customerRepository.create(customer);

    order.changeCustomer(customer);
    await orderRepository.update(order);

    const updatedOrder = await orderRepository.find(order.id);

    expect(updatedOrder.customerId).toStrictEqual(order.customerId);
  });
});
