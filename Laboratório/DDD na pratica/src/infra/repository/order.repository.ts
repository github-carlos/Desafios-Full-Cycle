import Order from "../../domain/entity/order";
import OrderItem from "../../domain/entity/order_item";
import { OrderRepositoryInterface } from "../../domain/repository/order.repository";
import { OrderItemModel, OrderModel } from "../db/sequelize/models/order.model";

export default class OrderRepository implements OrderRepositoryInterface {
  async find(id: string): Promise<Order> {
    try {
      const orderModel = await OrderModel.findOne({
        where: { id },
        include: [{ model: OrderItemModel }],
        rejectOnEmpty: true,
      });
      return new Order(
        orderModel.id,
        orderModel.customer_id,
        orderModel.items.map(
          (item) =>
            new OrderItem(
              item.id,
              item.name,
              item.price,
              item.product_id,
              item.quantity
            )
        )
      );
    } catch (err) {
      throw new Error("Order not found");
    }
  }
  findAll(): Promise<Order[]> {
    throw new Error("Method not implemented.");
  }
  async create(entity: Order): Promise<void> {
    await OrderModel.create(
      {
        id: entity.id,
        customer_id: entity.customerId,
        total: entity.total(),
        items: entity.items.map((item) => ({
          id: item.id,
          name: item.name,
          price: item.price,
          product_id: item.productId,
          quantity: item.quantity,
        })),
      },
      { include: [{ model: OrderItemModel }] }
    );
  }

  async update(entity: Order): Promise<void> {
    await OrderModel.update(
      {
        customer_id: entity.customerId,
        total: entity.total(),
      },
      {
        where: {
          id: entity.id,
        },
      }
    );
  }


  // async findAll(): Promise<Order[]> {
  //   const orderModels = await OrderModel.findAll();

  //   const orders = orderModels.map((orderModels) => {
  //     let order = new Order();

  //     return order;
  //   });

  //   return orders;
  // }
}
