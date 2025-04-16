// src/api/sse.ts
export const connectSSE = (onMessage: (msg: any) => void): EventSource => {
    const clientId = "client2";
    const source = new EventSource('http://localhost:9080/events/order_created?clientId=' + clientId);

    source.onmessage = (event) => {
        const data = JSON.parse(event.data);
        console.log("Type of event.data, ", typeof (event.data))
        // Kiểm tra status_code và xử lý dữ liệu theo tình huống
        if (data.status_code === '200') {
            console.log('Order created successfully:', data.value);
            onMessage({
                success: true,
                message: 'Order created successfully!',
                order: data.value // Đối tượng order
            });
        } else {
            console.log('Error occurred:', data.message);
            onMessage({
                success: false,
                message: 'Error: ' + data.message,
                error: data.value // Mô tả lỗi
            });
        }
    };

    source.onerror = (err) => {
        console.log('SSE Error:', err);
        source.close();
    };

    return source;
};
