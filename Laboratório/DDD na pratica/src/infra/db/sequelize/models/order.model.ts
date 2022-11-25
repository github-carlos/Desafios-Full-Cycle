import {
  Table,
  Model,
  PrimaryKey,
  Column,
  ForeignKey,
  BelongsTo,
  HasMany
} from "sequelize-typescript";
import CustomerModel from "./customer.model";
import ProductModel from "./product";

@Table({
  tableName: "orders",
  timestamps: false,
})
export class OrderModel extends Model {
  @PrimaryKey
  @Column
  declare id: string;

  @ForeignKey(() => CustomerModel)
  @Column({allowNull: false})
  declare customer_id: string;

  @BelongsTo(() => CustomerModel)
  declare customer: CustomerModel;

  @Column({allowNull: false})
  declare total: number;

  @HasMany(() => OrderItemModel)
  declare items: OrderItemModel[];
}

@Table({
  tableName: "order_items",
  timestamps: false,
})
export class OrderItemModel extends Model {
  @PrimaryKey
  @Column
  declare id: string;

  @ForeignKey(() => ProductModel)
  @Column({allowNull: false})
  declare product_id: string;

  @BelongsTo(() => ProductModel)
  declare product: ProductModel;

  @ForeignKey(() => OrderModel)
  @Column({allowNull: false})
  declare order_id: string;

  @BelongsTo(() => OrderModel)
  declare order: OrderModel;

  @Column({allowNull: false})
  declare quantity: number;

  @Column({allowNull: false})
  declare name: string;

  @Column({allowNull: false})
  declare price: number;
}