// src/types/Order.ts
export interface OrderDetail {
    id: string;
    product_id: string;
    name: string;
    price: number;
    quantity: number;
    total: number;
}

export interface Order {
    id: string;
    user_id: string;
    recipient_name: string;
    contact_phone: string;
    email: string;
    address: string;
    order_date: string;
    status: string;
    total: number;
    note: string;
    details: OrderDetail[];
}

