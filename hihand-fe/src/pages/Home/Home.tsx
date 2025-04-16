import { useEffect, useState } from 'react';
import { connectSSE } from '../../api/sse';
import { Order } from '../../types/Order';
import { fetchItems, createOrder } from '../../api';
import './Home.css';

const Home = () => {
    const [orders, setOrders] = useState<Order[]>([]);
    const [loading, setLoading] = useState<boolean>(true);
    const [newOrder, setNewOrder] = useState<Order>({
        id: '',
        recipient_name: 'Nguyen Van C',
        user_id: '12345',
        contact_phone: '0123456789',
        email: 'nguyenvana@example.com',
        address: '123 Đường ABC, Quận 1, TP. HCM',
        order_date: '',
        status: 'Created',
        total: 500000,
        note: '',
        details: [],
    });

    useEffect(() => {
        const getOrders = async () => {
            try {
                const fetchedOrders = await fetchItems();
                setOrders(fetchedOrders);
            } catch (error) {
                console.error('Failed to fetch orders:', error);
            } finally {
                setLoading(false);
            }
        };

        getOrders();
    }, []);

    useEffect(() => {
        if (!loading) {
            const source = connectSSE((msg: any) => {
                console.log('msg:', msg);
                if (msg.success && msg.order) {
                    setOrders((prev) => [msg.order, ...prev]);
                } else {
                    console.error('SSE error:', msg.message);
                }
            });

            return () => source.close();
        }
    }, [loading]);

    const handleCreateOrder = async () => {
        try {
            if (!newOrder.recipient_name || !newOrder.contact_phone || newOrder.details.length === 0) {
                alert('Vui lòng điền đầy đủ thông tin người nhận và ít nhất một sản phẩm.');
                return;
            }
            const message = await createOrder(newOrder);
            alert(message);

            // Reset lại đơn hàng sau khi tạo
            setNewOrder({
                id: '',
                recipient_name: 'Nguyen Van C',
                user_id: '12345',
                contact_phone: '0123456789',
                email: 'nguyenvana@example.com',
                address: '123 Đường ABC, Quận 1, TP. HCM',
                order_date: '',
                status: 'Created',
                total: 500000,
                note: '',
                details: [],
            });
        } catch (error: any) {
            alert(error.message || 'Tạo đơn hàng thất bại!');
            console.error('Failed to create order:', error);
        }
    };


    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { name, value } = e.target;
        setNewOrder((prevOrder) => ({ ...prevOrder, [name]: value }));
    };

    const handleOrderDetailsChange = (e: React.ChangeEvent<HTMLInputElement>, index: number) => {
        const { name, value } = e.target;
        const newDetails = [...newOrder.details];
        const updatedDetail = {
            ...newDetails[index],
            [name]: name === 'quantity' || name === 'price' ? Number(value) : value,
        };

        const quantity = name === 'quantity' ? Number(value) : updatedDetail.quantity;
        const price = name === 'price' ? Number(value) : updatedDetail.price;
        updatedDetail.total = quantity * price;

        newDetails[index] = updatedDetail;
        setNewOrder((prevOrder) => ({ ...prevOrder, details: newDetails }));
    };

    const handleAddProductDetail = () => {
        const newProductDetail = {
            id: '',
            name: 'Laptop 1',
            product_id: 'a1',
            quantity: 10,
            price: 0,
            total: 0,
        };
        setNewOrder((prevOrder) => ({
            ...prevOrder,
            details: [...prevOrder.details, newProductDetail],
        }));
    };

    if (loading) {
        return <div>Loading...</div>;
    }

    return (
        <div className="list-item">
            <h1>Orders</h1>

            <div className="create-order-form">
                <div className='order'>
                    <h2>Tạo đơn hàng mới</h2>
                    <div>
                        <input
                            type="text"
                            name="user_id"
                            value={newOrder.user_id}
                            onChange={handleInputChange}
                            placeholder="User ID"
                        />
                    </div>
                    <div>
                        <input
                            type="text"
                            name="recipient_name"
                            value={newOrder.recipient_name}
                            onChange={handleInputChange}
                            placeholder="Người nhận"
                        />
                    </div>
                    <div>
                        <input
                            type="text"
                            name="contact_phone"
                            value={newOrder.contact_phone}
                            onChange={handleInputChange}
                            placeholder="Điện thoại"
                        />
                    </div>
                    <div>
                        <input
                            type="email"
                            name="email"
                            value={newOrder.email}
                            onChange={handleInputChange}
                            placeholder="Email"
                        />
                    </div>
                    <div>
                        <input
                            type="text"
                            name="address"
                            value={newOrder.address}
                            onChange={handleInputChange}
                            placeholder="Địa chỉ"
                        />
                    </div>
                    <div>
                        <textarea
                            name="note"
                            value={newOrder.note}
                            onChange={handleInputChange}
                            placeholder="Ghi chú"
                        />
                    </div>
                </div>

                {/* Chi tiết sản phẩm */}
                <div className='order_details'>
                    <h3>Chi tiết sản phẩm</h3>
                    {newOrder.details.map((detail, index) => (
                        <div key={index} className="product-detail-row">
                            <div className="detail-field">
                                <label htmlFor={`name-${index}`}>Mã sản phẩm</label>
                                <input
                                    id={`name-${index}`}
                                    type="text"
                                    name="product_id"
                                    value={detail.product_id}
                                    onChange={(e) => handleOrderDetailsChange(e, index)}
                                    placeholder="Mã sản phẩm"
                                />
                            </div>
                            <div className="detail-field">
                                <label htmlFor={`name-${index}`}>Tên sản phẩm</label>
                                <input
                                    id={`name-${index}`}
                                    type="text"
                                    name="name"
                                    value={detail.name}
                                    onChange={(e) => handleOrderDetailsChange(e, index)}
                                    placeholder="Tên sản phẩm"
                                />
                            </div>
                            <div className="detail-field">
                                <label htmlFor={`quantity-${index}`}>Số lượng</label>
                                <input
                                    id={`quantity-${index}`}
                                    type="number"
                                    name="quantity"
                                    value={detail.quantity}
                                    onChange={(e) => handleOrderDetailsChange(e, index)}
                                    placeholder="Số lượng"
                                />
                            </div>
                            <div className="detail-field">
                                <label htmlFor={`price-${index}`}>Giá</label>
                                <input
                                    id={`price-${index}`}
                                    type="number"
                                    name="price"
                                    value={detail.price}
                                    onChange={(e) => handleOrderDetailsChange(e, index)}
                                    placeholder="Giá"
                                />
                            </div>
                            <div className="detail-field">
                                <label htmlFor={`total-${index}`}>Thành tiền</label>
                                <input
                                    id={`total-${index}`}
                                    type="number"
                                    name="total"
                                    value={detail.total}
                                    onChange={(e) => handleOrderDetailsChange(e, index)}
                                    placeholder="Thành tiền"
                                />
                            </div>
                        </div>
                    ))}

                    <button onClick={handleAddProductDetail}>+ Thêm sản phẩm</button>
                </div>
            </div>
            <button onClick={handleCreateOrder}>Tạo đơn hàng</button>
            <ul>
                {orders.map((item) => (
                    <li key={item.id}>
                        <h3>Người nhận: {item.recipient_name}</h3>
                        <p>Điện thoại: {item.contact_phone}</p>
                        <p>Địa chỉ: {item.address}</p>
                        <p>Tổng tiền: {item.total.toLocaleString()} VND</p>
                        <p>Trạng thái: {item.status}</p>
                        {item.details.length > 0 && (
                            <ul>
                                {item.details.map((detail) => (
                                    <li key={detail.id}>
                                        {detail.name} - SL: {detail.quantity} - Giá: {detail.price.toLocaleString()} VND
                                    </li>
                                ))}
                            </ul>
                        )}
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default Home;
