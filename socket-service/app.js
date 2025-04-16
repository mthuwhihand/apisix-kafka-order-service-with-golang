const express = require("express");
const http = require("http");
const { Server } = require("socket.io");
const { Kafka } = require("kafkajs");

const app = express();
const server = http.createServer(app);

// 👇 Socket.IO server setup với path là "/socket"
const io = new Server(server, {
  path: "/socket",
  cors: {
    origin: "*", // Có thể chỉ định origin cụ thể nếu muốn
    methods: ["GET", "POST"]
  }
});

// Kafka client setup
const kafka = new Kafka({
  clientId: "ws-service",
  brokers: ["localhost:9092"], // hoặc dùng địa chỉ container nếu chạy docker
});
const consumer = kafka.consumer({ groupId: "order_created_consumer_group" });

(async () => {
  await consumer.connect();
  await consumer.subscribe({ topic: "created_orders", fromBeginning: true });

  // Consume Kafka messages
  await consumer.run({
    eachMessage: async ({ topic, partition, message }) => {
      const clientId = message.key.toString();
      const rawBuffer = message.value;

      try {
        const decoded = rawBuffer.toString();
        const parsed = JSON.parse(decoded);

        console.log(`📨 Message for -${clientId}-:`, parsed);

        // Emit to specific clientId
        io.to(clientId).emit("order_created", parsed);
      } catch (err) {
        console.error("❌ Failed to parse message:", err);
      }
    },
  });

  // Socket.IO connection
  io.on("connection", (socket) => {
    console.log(`Socket connected: ${socket.id}`);

    socket.on("register", (clientId) => {
      socket.join(clientId);
      console.log(`Registered clientId ${clientId} to socket ${socket.id}`);
    });

    socket.on("disconnect", () => {
      console.log(`Disconnected: ${socket.id}`);
    });
  });

  // Start server
  server.listen(8082, () => {
    console.log("WebSocket service is running on port 8082");
  });
})();
