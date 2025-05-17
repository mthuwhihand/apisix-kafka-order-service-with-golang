// src/types/Order.ts

export interface OrderDetail {
    id: string;
    order_id: string;
    user_id: string;
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
    status: string;
    total: number;
    note: string
    details: OrderDetail[];
    created_at?: string;
    updated_at?: string;
}

export interface PaginatedOrdersResponse {
    data: Order[];
    total: number;
    page: number;
    page_size: number;
}

export interface OrderDetailRequest {
    user_id: string;
    product_id: string;
    name: string;
    price: number;
    quantity: number;
    total: number;
}

export interface OrderRequest {
    user_id: string;
    recipient_name: string;
    contact_phone: string;
    email: string;
    address: string;
    status: string;
    total: number;
    note: string
    details: OrderDetailRequest[];
    created_at?: string;
    updated_at?: string;
}